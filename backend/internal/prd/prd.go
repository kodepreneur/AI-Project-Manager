package prd

import (
	"fmt"
	"strings"
	"time"

	"ai-project-manager/internal/database"
	"ai-project-manager/internal/logs"
	"ai-project-manager/internal/models"

	"github.com/gofiber/fiber/v2"
)

type PRDService struct{}

func NewPRDService() *PRDService {
	return &PRDService{}
}

type GeneratePRDRequest struct {
	ProjectIdea           string `json:"project_idea"`
	TargetUsers           string `json:"target_users"`
	BusinessGoals         string `json:"business_goals"`
	MainFeatures          string `json:"main_features"`
	TechnicalRequirements string `json:"technical_requirements"`
	AdditionalNotes       string `json:"additional_notes"`
}

type SavePRDRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Status  string `json:"status"`
}

func (s *PRDService) Get(c *fiber.Ctx) error {
	projectID := c.Params("id")
	var prd models.PRD
	err := database.DB.Preload("Versions").Where("project_id = ?", projectID).First(&prd).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{"code": "NOT_FOUND", "message": "PRD not found for this project"},
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    prd,
	})
}

func (s *PRDService) Save(c *fiber.Ctx) error {
	projectID := c.Params("id")
	var req SavePRDRequest
	if err := c.BodyParser(&req); err != nil || req.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{"code": "INVALID_INPUT", "message": "PRD content is required"},
		})
	}

	if req.Title == "" {
		req.Title = "Product Requirements Document"
	}
	if req.Status == "" {
		req.Status = "Draft"
	}

	var existingPRD models.PRD
	err := database.DB.Where("project_id = ?", projectID).First(&existingPRD).Error

	if err != nil {
		// Create new PRD
		var pid uint
		fmt.Sscanf(projectID, "%d", &pid)
		newPRD := models.PRD{
			ProjectID: pid,
			Title:     req.Title,
			Content:   req.Content,
			Status:    req.Status,
			Version:   1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		database.DB.Create(&newPRD)

		// Create first version history
		database.DB.Create(&models.PRDVersion{
			PRDID:     newPRD.ID,
			Version:   1,
			Content:   newPRD.Content,
			Summary:   "Initial PRD creation",
			CreatedAt: time.Now(),
		})

		logs.Log.Info(&pid, "PRD", "PRD created for project")
		return c.JSON(fiber.Map{"success": true, "data": newPRD})
	}

	// Update existing PRD and archive version
	existingPRD.Version++
	existingPRD.Title = req.Title
	existingPRD.Content = req.Content
	existingPRD.Status = req.Status
	existingPRD.UpdatedAt = time.Now()

	database.DB.Save(&existingPRD)
	database.DB.Create(&models.PRDVersion{
		PRDID:     existingPRD.ID,
		Version:   existingPRD.Version,
		Content:   existingPRD.Content,
		Summary:   fmt.Sprintf("Revision %d update", existingPRD.Version),
		CreatedAt: time.Now(),
	})

	logs.Log.Info(&existingPRD.ProjectID, "PRD", fmt.Sprintf("PRD updated to version %d", existingPRD.Version))

	return c.JSON(fiber.Map{
		"success": true,
		"data":    existingPRD,
	})
}

func (s *PRDService) Approve(c *fiber.Ctx) error {
	projectID := c.Params("id")
	var prd models.PRD
	if err := database.DB.Where("project_id = ?", projectID).First(&prd).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{"code": "NOT_FOUND", "message": "PRD not found"},
		})
	}

	now := time.Now()
	prd.Status = "Approved"
	prd.ApprovedAt = &now
	database.DB.Save(&prd)

	// Update project last activity
	database.DB.Model(&models.Project{}).Where("id = ?", projectID).Update("last_activity", now)

	logs.Log.Info(&prd.ProjectID, "PRD", "PRD officially approved. Ready for development task generation.")

	return c.JSON(fiber.Map{
		"success": true,
		"message": "PRD approved successfully",
		"data":    prd,
	})
}

func (s *PRDService) ExportMarkdown(c *fiber.Ctx) error {
	projectID := c.Params("id")
	var prd models.PRD
	if err := database.DB.Where("project_id = ?", projectID).First(&prd).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("PRD not found")
	}

	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"PRD-Project-%s.md\"", projectID))
	c.Set("Content-Type", "text/markdown; charset=utf-8")
	return c.SendString(prd.Content)
}

func (s *PRDService) GenerateAI(c *fiber.Ctx) error {
	projectID := c.Params("id")
	var req GeneratePRDRequest
	if err := c.BodyParser(&req); err != nil || req.ProjectIdea == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{"code": "INVALID_INPUT", "message": "Project idea is required"},
		})
	}

	var project models.Project
	database.DB.First(&project, projectID)

	techStack := fmt.Sprintf("%s / %s / %s", project.Framework, project.FrontendTech, project.DatabaseType)
	if req.TechnicalRequirements != "" {
		techStack += fmt.Sprintf(" (%s)", req.TechnicalRequirements)
	}

	// Generate structured 17-section PRD
	generatedMarkdown := s.constructStructuredPRD(project.Name, req, techStack)

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"title":   fmt.Sprintf("PRD — %s", project.Name),
			"content": generatedMarkdown,
		},
	})
}

func (s *PRDService) constructStructuredPRD(name string, req GeneratePRDRequest, techStack string) string {
	features := strings.Split(req.MainFeatures, "\n")
	formattedFeatures := ""
	for _, f := range features {
		f = strings.TrimSpace(f)
		if f != "" {
			formattedFeatures += fmt.Sprintf("- **%s**: Detailed specification with edge-case handling.\n", f)
		}
	}
	if formattedFeatures == "" {
		formattedFeatures = "- **Core CRUD**: Intuitive management interface with data validation.\n- **Authentication & Roles**: Secure role-based access control.\n- **Dashboard & Reporting**: Real-time metrics and status telemetry.\n"
	}

	targetUsers := req.TargetUsers
	if targetUsers == "" {
		targetUsers = "Software engineers, system administrators, and product end users."
	}

	businessGoals := req.BusinessGoals
	if businessGoals == "" {
		businessGoals = "Streamline operational overhead, improve automation velocity, and provide a single pane of glass for workflow execution."
	}

	return fmt.Sprintf(`# Product Requirements Document (PRD)

## 1. Product Overview
**%s** is designed to deliver a modern, resilient solution for:
%s

## 2. Problem Statement
Users currently face fragmented workflows, complex manual setup procedures, and lack of real-time visibility into project and development states.

## 3. Goals
%s

## 4. Target Users
%s

## 5. User Roles
- **Administrator**: Complete control over project configurations, secrets, deployments, and database instances.
- **Developer / Contributor**: Creates development tasks, executes AI prompts, manages branches, and reviews test runs.
- **Viewer / Stakeholder**: Reviews PRD specifications and inspects application previews.

## 6. Functional Requirements
%s

## 7. Modules
1. **Core Runtime Engine**: Lifecycle and process orchestration.
2. **AI Development Pipeline**: Integration with Antigravity CLI for automated coding and validation.
3. **Database & Persistence**: Schema migrations, connection pooling, and live query consoles.
4. **Interactive Dev Terminal**: Low-latency PTY session streaming directly to browser.
5. **Preview & Diagnostics**: Real-time telemetry, automated screenshot generation, and health checks.

## 8. User Flows
1. **Project Initiation**: Create project via template wizard -> Generate AI PRD -> Review and approve PRD.
2. **Task Execution**: PRD converts into phased development plan -> Send task context to Antigravity CLI -> Live terminal streaming -> Automated test run -> Git diff human review -> Mark Done.
3. **Preview & Deployment**: Start runtime server -> Verify preview -> Check system diagnostics.

## 9. Database Requirements
- Primary engine: %s
- Automated relational migrations with rollback capability.
- Indexed lookups on foreign keys and slugs.
- Soft deletion where audit history is mandated.

## 10. API Requirements
- Clean RESTful design with structured JSON responses (e.g. { "success": true, "data": ... }).
- WebSocket endpoint for low-latency streaming (terminal I/O, process logs, AI stdout).
- Strict bearer token authentication with role enforcement.

## 11. UI Requirements
- Dark-first aesthetic engineered for developer tooling (Vercel/Linear style).
- Responsive layout with desktop command palette (⌘K) and collapsible drawer for mobile.
- High-contrast typography, interactive status badges, and drag-and-drop Kanban board.

## 12. Non Functional Requirements
- High availability and lightweight memory footprint (< 50MB RAM idle).
- Fast cold-start time (< 500ms).
- Graceful shutdown for all child processes and active PTY terminals.

## 13. Security Requirements
- Robust password hashing using bcrypt.
- Masking of sensitive environment credentials (API keys, DB secrets).
- Strict path traversal sanitization on all project file endpoints.

## 14. Performance Requirements
- Sub-50ms API latency for standard queries.
- Zero-buffer WebSocket streaming for terminal commands.
- Paginated log views with client-side filtering.

## 15. Deployment Requirements
- Systemd service daemon support.
- Production reverse-proxy templates for Nginx and Caddy with WebSocket upgrade.
- Single-binary executable deployment alongside static web assets.

## 16. Acceptance Criteria
- [ ] PRD approved by project stakeholder.
- [ ] Development tasks systematically broken down into Kanban backlog.
- [ ] Antigravity CLI connects to project context and executes tasks.
- [ ] Automated tests pass with zero critical regressions.
- [ ] Live preview and health checks report 100%% operational status.

## 17. Development Phases
- **Phase 1 — Project Foundation**: Repository setup, database schema, and authentication.
- **Phase 2 — Core Engine & APIs**: Business logic, CRUD handlers, and WebSocket pipelines.
- **Phase 3 — UI & Developer Experience**: Kanban board, interactive terminal, and metrics dashboard.
- **Phase 4 — AI Automation & Integrations**: Antigravity CLI orchestration, validation hooks, and preview.
- **Phase 5 — Hardening & Release**: Unit test coverage, reverse proxy configs, and production installer.
`, name, req.ProjectIdea, businessGoals, targetUsers, formattedFeatures, techStack)
}
