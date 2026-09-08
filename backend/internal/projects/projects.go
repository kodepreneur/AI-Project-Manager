package projects

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"ai-project-manager/internal/config"
	"ai-project-manager/internal/database"
	"ai-project-manager/internal/logs"
	"ai-project-manager/internal/models"

	"github.com/gofiber/fiber/v2"
)

type ProjectService struct {
	cfg *config.Config
}

func NewProjectService(cfg *config.Config) *ProjectService {
	return &ProjectService{cfg: cfg}
}

type CreateProjectRequest struct {
	Name           string `json:"name"`
	Slug           string `json:"slug"`
	Description    string `json:"description"`
	Path           string `json:"path"`
	GitRepo        string `json:"git_repo"`
	GitBranch      string `json:"git_branch"`
	Framework      string `json:"framework"`
	Language       string `json:"language"`
	FrontendTech   string `json:"frontend_tech"`
	DatabaseType   string `json:"database_type"`
	PHPVersion     string `json:"php_version"`
	PreviewPort    int    `json:"preview_port"`
	PreviewCommand string `json:"preview_command"`
	TestCommand    string `json:"test_command"`
	Template       string `json:"template"`
}

func (s *ProjectService) List(c *fiber.Ctx) error {
	status := c.Query("status")
	search := c.Query("search")

	query := database.DB.Model(&models.Project{}).Preload("PRD").Preload("Tasks").Preload("Processes")

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if search != "" {
		query = query.Where("name LIKE ? OR description LIKE ? OR framework LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	var projects []models.Project
	query.Order("updated_at DESC").Find(&projects)

	return c.JSON(fiber.Map{
		"success": true,
		"data":    projects,
	})
}

func (s *ProjectService) Get(c *fiber.Ctx) error {
	id := c.Params("id")
	var project models.Project
	err := database.DB.Preload("PRD").Preload("Tasks").Preload("Processes").First(&project, id).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{"code": "NOT_FOUND", "message": "Project not found"},
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    project,
	})
}

func (s *ProjectService) Create(c *fiber.Ctx) error {
	var req CreateProjectRequest
	if err := c.BodyParser(&req); err != nil || req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{"code": "INVALID_INPUT", "message": "Project name is required"},
		})
	}

	if req.Slug == "" {
		req.Slug = strings.ToLower(strings.ReplaceAll(req.Name, " ", "-"))
	}

	// Check if slug exists; if so, append suffix
	var count int64
	database.DB.Model(&models.Project{}).Where("slug = ?", req.Slug).Count(&count)
	if count > 0 {
		req.Slug = fmt.Sprintf("%s-%d", req.Slug, time.Now().Unix()%10000)
	}

	if req.Path == "" {
		req.Path = filepath.Join(s.cfg.ProjectRoot, req.Slug)
	}
	if req.GitBranch == "" {
		req.GitBranch = "main"
	}
	if req.PreviewPort == 0 {
		req.PreviewPort = 8000 + int(time.Now().Unix()%1000)
	}
	if req.PreviewCommand == "" {
		req.PreviewCommand = "npm run dev"
	}
	if req.TestCommand == "" {
		req.TestCommand = "npm test"
	}

	// Create physical project directory
	_ = os.MkdirAll(req.Path, 0755)

	// Scaffold initial files if template provided
	s.scaffoldTemplate(req.Path, req.Template, req.Name)

	project := models.Project{
		Name:           req.Name,
		Slug:           req.Slug,
		Description:    req.Description,
		Path:           req.Path,
		GitRepo:        req.GitRepo,
		GitBranch:      req.GitBranch,
		Framework:      req.Framework,
		Language:       req.Language,
		FrontendTech:   req.FrontendTech,
		DatabaseType:   req.DatabaseType,
		PHPVersion:     req.PHPVersion,
		Status:         "Development",
		PreviewPort:    req.PreviewPort,
		PreviewCommand: req.PreviewCommand,
		TestCommand:    req.TestCommand,
		LastActivity:   time.Now(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := database.DB.Create(&project).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{"code": "DATABASE_ERROR", "message": err.Error()},
		})
	}

	logs.Log.Info(&project.ID, "Project", fmt.Sprintf("Project created: %s (%s)", project.Name, project.Path))

	return c.JSON(fiber.Map{
		"success": true,
		"data":    project,
	})
}

func (s *ProjectService) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var project models.Project
	if err := database.DB.First(&project, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{"code": "NOT_FOUND", "message": "Project not found"},
		})
	}

	if err := c.BodyParser(&project); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{"code": "INVALID_INPUT", "message": err.Error()},
		})
	}

	project.UpdatedAt = time.Now()
	project.LastActivity = time.Now()
	database.DB.Save(&project)

	logs.Log.Info(&project.ID, "Project", "Project updated: "+project.Name)

	return c.JSON(fiber.Map{
		"success": true,
		"data":    project,
	})
}

func (s *ProjectService) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	var project models.Project
	if err := database.DB.First(&project, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{"code": "NOT_FOUND", "message": "Project not found"},
		})
	}

	// Delete related tasks, PRD, sessions
	database.DB.Where("project_id = ?", project.ID).Delete(&models.Task{})
	database.DB.Where("project_id = ?", project.ID).Delete(&models.PRD{})
	database.DB.Where("project_id = ?", project.ID).Delete(&models.AntigravitySession{})
	database.DB.Where("project_id = ?", project.ID).Delete(&models.ProjectProcess{})
	database.DB.Delete(&project)

	logs.Log.Warn(nil, "Project", "Project deleted: "+project.Name)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Project deleted successfully",
	})
}

func (s *ProjectService) Archive(c *fiber.Ctx) error {
	id := c.Params("id")
	var project models.Project
	if err := database.DB.First(&project, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{"code": "NOT_FOUND", "message": "Project not found"},
		})
	}

	if project.Status == "Archived" {
		project.Status = "Development"
	} else {
		project.Status = "Archived"
	}
	database.DB.Save(&project)

	return c.JSON(fiber.Map{
		"success": true,
		"data":    project,
	})
}

func (s *ProjectService) HealthCheck(c *fiber.Ctx) error {
	id := c.Params("id")
	var project models.Project
	if err := database.DB.First(&project, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{"code": "NOT_FOUND", "message": "Project not found"},
		})
	}

	results := []fiber.Map{}

	// 1. Directory check
	_, dirErr := os.Stat(project.Path)
	results = append(results, fiber.Map{
		"name":        "Project Directory",
		"passed":      dirErr == nil,
		"description": project.Path,
	})

	// 2. Git repo check
	_, gitErr := os.Stat(filepath.Join(project.Path, ".git"))
	results = append(results, fiber.Map{
		"name":        "Git Repository",
		"passed":      gitErr == nil,
		"description": "Git tracking initialized in project folder",
	})

	// 3. Node & NPM check
	_, nodeErr := exec.LookPath("node")
	_, npmErr := exec.LookPath("npm")
	results = append(results, fiber.Map{
		"name":        "Node.js & NPM Runtime",
		"passed":      nodeErr == nil && npmErr == nil,
		"description": "JavaScript package manager available",
	})

	// 4. PHP & Composer check (if PHP project)
	if strings.Contains(strings.ToLower(project.Language), "php") || strings.Contains(strings.ToLower(project.Framework), "laravel") {
		_, phpErr := exec.LookPath("php")
		_, compErr := exec.LookPath("composer")
		results = append(results, fiber.Map{
			"name":        "PHP & Composer Runtime",
			"passed":      phpErr == nil,
			"description": fmt.Sprintf("PHP CLI (%s) and Composer", project.PHPVersion),
			"composer":    compErr == nil,
		})
	}

	// 5. Environment configuration (.env)
	_, envErr := os.Stat(filepath.Join(project.Path, ".env"))
	results = append(results, fiber.Map{
		"name":        "Environment File (.env)",
		"passed":      envErr == nil,
		"description": "Project environment variables configured",
	})

	// 6. PRD Status
	var prd models.PRD
	hasPrd := database.DB.Where("project_id = ?", project.ID).First(&prd).Error == nil
	results = append(results, fiber.Map{
		"name":        "Approved PRD",
		"passed":      hasPrd && prd.Status == "Approved",
		"description": "Product Requirements Document reviewed and approved",
	})

	return c.JSON(fiber.Map{
		"success": true,
		"data":    results,
	})
}

func (s *ProjectService) scaffoldTemplate(path, template, name string) {
	readmeContent := fmt.Sprintf("# %s\n\nGenerated by AI Project Manager.\n", name)
	_ = os.WriteFile(filepath.Join(path, "README.md"), []byte(readmeContent), 0644)

	envContent := fmt.Sprintf("APP_NAME=\"%s\"\nAPP_ENV=local\nAPP_DEBUG=true\nAPP_PORT=8001\n", name)
	_ = os.WriteFile(filepath.Join(path, ".env.example"), []byte(envContent), 0644)
	_ = os.WriteFile(filepath.Join(path, ".env"), []byte(envContent), 0644)

	// If git not initialized, run git init
	if _, err := os.Stat(filepath.Join(path, ".git")); os.IsNotExist(err) {
		cmd := exec.Command("git", "init")
		cmd.Dir = path
		_ = cmd.Run()
	}
}
