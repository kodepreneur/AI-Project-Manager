package env

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"ai-project-manager/internal/database"
	"ai-project-manager/internal/logs"
	"ai-project-manager/internal/models"

	"github.com/gofiber/fiber/v2"
)

type EnvService struct{}

func NewEnvService() *EnvService {
	return &EnvService{}
}

type EnvItem struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	IsSecret  bool   `json:"is_secret"`
	IsComment bool   `json:"is_comment"`
}

type UpdateEnvRequest struct {
	Items []EnvItem `json:"items"`
}

func (s *EnvService) Get(c *fiber.Ctx) error {
	projectID := c.Params("id")
	var project models.Project
	if err := database.DB.First(&project, projectID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "error": "Project not found"})
	}

	envPath := filepath.Join(project.Path, ".env")
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		// Fallback to .env.example
		envPath = filepath.Join(project.Path, ".env.example")
	}

	content, err := os.ReadFile(envPath)
	if err != nil {
		return c.JSON(fiber.Map{"success": true, "data": []EnvItem{}})
	}

	items := parseEnv(string(content))

	return c.JSON(fiber.Map{
		"success": true,
		"data":    items,
	})
}

func (s *EnvService) Save(c *fiber.Ctx) error {
	projectID := c.Params("id")
	var project models.Project
	if err := database.DB.First(&project, projectID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "error": "Project not found"})
	}

	var req UpdateEnvRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "Invalid body"})
	}

	var builder strings.Builder
	for _, item := range req.Items {
		if item.IsComment {
			builder.WriteString(item.Value + "\n")
		} else if item.Key != "" {
			builder.WriteString(fmt.Sprintf("%s=%s\n", item.Key, item.Value))
		}
	}

	envPath := filepath.Join(project.Path, ".env")
	if err := os.WriteFile(envPath, []byte(builder.String()), 0644); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}

	logs.Log.Info(&project.ID, "Environment", "Project .env file updated")

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Environment file saved",
	})
}

func parseEnv(content string) []EnvItem {
	var items []EnvItem
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(trimmed, "#") {
			items = append(items, EnvItem{
				Key:       "",
				Value:     trimmed,
				IsComment: true,
			})
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			k := strings.TrimSpace(parts[0])
			v := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
			upperK := strings.ToUpper(k)
			isSecret := strings.Contains(upperK, "PASS") ||
				strings.Contains(upperK, "SECRET") ||
				strings.Contains(upperK, "KEY") ||
				strings.Contains(upperK, "TOKEN")

			items = append(items, EnvItem{
				Key:      k,
				Value:    v,
				IsSecret: isSecret,
			})
		}
	}
	return items
}
