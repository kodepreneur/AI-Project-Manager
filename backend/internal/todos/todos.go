package todos

import (
	"fmt"
	"time"

	"ai-project-manager/internal/database"
	"ai-project-manager/internal/logs"
	"ai-project-manager/internal/models"

	"github.com/gofiber/fiber/v2"
)

type TodoService struct{}

func NewTodoService() *TodoService {
	return &TodoService{}
}

type CreateTaskRequest struct {
	PhaseName          string `json:"phase_name"`
	Title              string `json:"title"`
	Description        string `json:"description"`
	Priority           string `json:"priority"`
	Complexity         string `json:"complexity"`
	Status             string `json:"status"`
	AcceptanceCriteria string `json:"acceptance_criteria"`
	AIInstructions     string `json:"ai_instructions"`
	Dependencies       string `json:"dependencies"`
	FilesAffected      string `json:"files_affected"`
}

type UpdateStatusRequest struct {
	Status string `json:"status"`
}

func (s *TodoService) List(c *fiber.Ctx) error {
	projectID := c.Params("id")
	tasks := []models.Task{}
	database.DB.Where("project_id = ?", projectID).Order("task_order ASC, created_at ASC").Find(&tasks)

	return c.JSON(fiber.Map{
		"success": true,
		"data":    tasks,
	})
}

func (s *TodoService) Create(c *fiber.Ctx) error {
	projectID := c.Params("id")
	var req CreateTaskRequest
	if err := c.BodyParser(&req); err != nil || req.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{"code": "INVALID_INPUT", "message": "Task title is required"},
		})
	}

	var pid uint
	fmt.Sscanf(projectID, "%d", &pid)

	if req.Status == "" {
		req.Status = "todo"
	}
	if req.Priority == "" {
		req.Priority = "medium"
	}
	if req.Complexity == "" {
		req.Complexity = "medium"
	}

	var count int64
	database.DB.Model(&models.Task{}).Where("project_id = ?", pid).Count(&count)

	task := models.Task{
		ProjectID:          pid,
		PhaseName:          req.PhaseName,
		Title:              req.Title,
		Description:        req.Description,
		Priority:           req.Priority,
		Complexity:         req.Complexity,
		Status:             req.Status,
		AcceptanceCriteria: req.AcceptanceCriteria,
		AIInstructions:     req.AIInstructions,
		Dependencies:       req.Dependencies,
		FilesAffected:      req.FilesAffected,
		Order:              int(count) + 1,
		AIStatus:           "idle",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if err := database.DB.Create(&task).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{"code": "CREATE_TASK_FAILED", "message": err.Error()},
		})
	}

	logs.Log.Info(&pid, "TODO", "Created task: "+task.Title)

	return c.JSON(fiber.Map{
		"success": true,
		"data":    task,
	})
}

func (s *TodoService) Update(c *fiber.Ctx) error {
	id := c.Params("taskId")
	var task models.Task
	if err := database.DB.First(&task, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{"code": "NOT_FOUND", "message": "Task not found"},
		})
	}

	if err := c.BodyParser(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{"code": "INVALID_INPUT", "message": err.Error()},
		})
	}

	task.UpdatedAt = time.Now()
	database.DB.Save(&task)

	return c.JSON(fiber.Map{
		"success": true,
		"data":    task,
	})
}

func (s *TodoService) UpdateStatus(c *fiber.Ctx) error {
	id := c.Params("taskId")
	var req UpdateStatusRequest
	if err := c.BodyParser(&req); err != nil || req.Status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{"code": "INVALID_INPUT", "message": "Status required"},
		})
	}

	var task models.Task
	if err := database.DB.First(&task, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{"code": "NOT_FOUND", "message": "Task not found"},
		})
	}

	task.Status = req.Status
	task.UpdatedAt = time.Now()
	database.DB.Save(&task)

	logs.Log.Info(&task.ProjectID, "TODO", fmt.Sprintf("Task '%s' moved to %s", task.Title, req.Status))

	return c.JSON(fiber.Map{
		"success": true,
		"data":    task,
	})
}

func (s *TodoService) Delete(c *fiber.Ctx) error {
	id := c.Params("taskId")
	database.DB.Delete(&models.Task{}, id)
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Task deleted",
	})
}

func (s *TodoService) GenerateFromPRD(c *fiber.Ctx) error {
	projectID := c.Params("id")
	var pid uint
	fmt.Sscanf(projectID, "%d", &pid)

	var project models.Project
	if err := database.DB.First(&project, pid).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{"code": "NOT_FOUND", "message": "Project not found"},
		})
	}

	var prd models.PRD
	if err := database.DB.Where("project_id = ?", pid).First(&prd).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{"code": "NO_PRD", "message": "Project does not have a PRD. Generate and approve a PRD first."},
		})
	}

	// Preset plan tasks tailored to the project's framework and tech
	tasks := s.buildDefaultTasksForProject(&project)

	for i, t := range tasks {
		t.ProjectID = pid
		t.Order = i + 1
		t.CreatedAt = time.Now()
		t.UpdatedAt = time.Now()
		database.DB.Create(&t)
	}

	logs.Log.Info(&pid, "TODO", fmt.Sprintf("Generated %d development tasks from approved PRD", len(tasks)))

	return c.JSON(fiber.Map{
		"success": true,
		"message": fmt.Sprintf("Generated %d tasks from PRD", len(tasks)),
		"count":   len(tasks),
	})
}

func (s *TodoService) buildDefaultTasksForProject(p *models.Project) []models.Task {
	return []models.Task{
		// Phase 1
		{
			PhaseName:          "Phase 1 — Project Foundation",
			Title:              "Initialize Repository & Environment Configuration",
			Description:        fmt.Sprintf("Configure %s project structure, setup .env variables, and initialize repository tracking.", p.Framework),
			Priority:           "high",
			Complexity:         "low",
			Status:             "todo",
			AcceptanceCriteria: "- Project directory contains valid configuration\n- .env exists with required app variables\n- Git repository initialized",
			AIInstructions:     "Verify all environment files and dependencies are properly wired for local execution.",
		},
		{
			PhaseName:          "Phase 1 — Project Foundation",
			Title:              "Database Connection & Initial Migrations",
			Description:        fmt.Sprintf("Configure %s connection and create base database migrations.", p.DatabaseType),
			Priority:           "high",
			Complexity:         "medium",
			Status:             "todo",
			AcceptanceCriteria: "- Database connection test passes\n- Base migrations execute cleanly without errors",
			AIInstructions:     "Create database models and verify schema against the PRD requirements.",
		},
		// Phase 2
		{
			PhaseName:          "Phase 2 — Core Engine & APIs",
			Title:              "Authentication & Role-Based Access Control",
			Description:        "Implement user login, token issuance, and authorization middleware.",
			Priority:           "urgent",
			Complexity:         "medium",
			Status:             "backlog",
			AcceptanceCriteria: "- Login endpoint returns valid JWT\n- Protected routes return 401 when unauthenticated",
			AIInstructions:     "Use secure password hashing (bcrypt) and adhere to standard REST authorization header format.",
		},
		{
			PhaseName:          "Phase 2 — Core Engine & APIs",
			Title:              "Implement Primary Resource Models & CRUD",
			Description:        "Implement core REST API endpoints and business logic validation as defined in the PRD.",
			Priority:           "high",
			Complexity:         "high",
			Status:             "backlog",
			AcceptanceCriteria: "- Full CRUD operations functional\n- Input validation with descriptive error payloads",
			AIInstructions:     "Build domain models, database queries, and route handlers with consistent JSON responses.",
		},
		// Phase 3
		{
			PhaseName:          "Phase 3 — UI & Frontend Integration",
			Title:              "Build Responsive Frontend Layout & Components",
			Description:        fmt.Sprintf("Construct %s components with Tailwind CSS styling.", p.FrontendTech),
			Priority:           "medium",
			Complexity:         "medium",
			Status:             "backlog",
			AcceptanceCriteria: "- Mobile-responsive layout\n- Interactive states, modals, and toasts implemented",
			AIInstructions:     "Ensure dark mode compatibility and sleek developer tooling aesthetic.",
		},
		{
			PhaseName:          "Phase 3 — UI & Frontend Integration",
			Title:              "Connect Frontend with Backend REST Endpoints",
			Description:        "Integrate API client, state management, and real-time data binding.",
			Priority:           "medium",
			Complexity:         "medium",
			Status:             "backlog",
			AcceptanceCriteria: "- Live data rendering from API\n- Optimistic UI updates with error handling",
			AIInstructions:     "Wire API services to frontend views using Axios and store state in Pinia.",
		},
		// Phase 4
		{
			PhaseName:          "Phase 4 — Testing & Quality Assurance",
			Title:              "Execute Automated Test Suite & Validation",
			Description:        fmt.Sprintf("Run test command: '%s' and resolve any regressions.", p.TestCommand),
			Priority:           "high",
			Complexity:         "medium",
			Status:             "backlog",
			AcceptanceCriteria: "- All unit and feature tests pass\n- Zero unhandled exceptions in logs",
			AIInstructions:     "Execute test suite, parse error stack traces, and apply fixes until 100% green.",
		},
	}
}
