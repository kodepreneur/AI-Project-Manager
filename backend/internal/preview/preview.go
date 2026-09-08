package preview

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"ai-project-manager/internal/config"
	"ai-project-manager/internal/database"
	"ai-project-manager/internal/models"

	"github.com/gofiber/fiber/v2"
)

type PreviewService struct {
	cfg *config.Config
}

func NewPreviewService(cfg *config.Config) *PreviewService {
	return &PreviewService{cfg: cfg}
}

func (s *PreviewService) GetPreviewInfo(c *fiber.Ctx) error {
	projectID := c.Params("id")
	var project models.Project
	if err := database.DB.First(&project, projectID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{"code": "NOT_FOUND", "message": "Project not found"},
		})
	}

	targetURL := fmt.Sprintf("http://localhost:%d", project.PreviewPort)

	// Quick probe to check if dev server is responding
	client := http.Client{Timeout: 800 * time.Millisecond}
	resp, err := client.Get(targetURL)
	isResponding := err == nil && resp != nil
	if resp != nil {
		_ = resp.Body.Close()
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"project_id":     project.ID,
			"preview_url":    targetURL,
			"preview_port":   project.PreviewPort,
			"is_responding":  isResponding,
			"screenshot_url": project.ScreenshotURL,
		},
	})
}

func (s *PreviewService) CaptureScreenshot(c *fiber.Ctx) error {
	projectID := c.Params("id")
	var project models.Project
	if err := database.DB.First(&project, projectID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{"code": "NOT_FOUND", "message": "Project not found"},
		})
	}

	targetURL := fmt.Sprintf("http://localhost:%d", project.PreviewPort)
	screenshotFileName := fmt.Sprintf("project_%d.svg", project.ID)
	screenshotPath := filepath.Join(s.cfg.StorageRoot, "screenshots", screenshotFileName)

	// Check if chromium or google-chrome is available
	chromePath, err := exec.LookPath("chromium")
	if err != nil {
		chromePath, err = exec.LookPath("google-chrome")
	}

	if err == nil && chromePath != "" {
		// Attempt headless capture
		pngFile := filepath.Join(s.cfg.StorageRoot, "screenshots", fmt.Sprintf("project_%d.png", project.ID))
		cmd := exec.Command(chromePath, "--headless", "--disable-gpu", "--screenshot="+pngFile, "--window-size=1280,800", targetURL)
		if cmdErr := cmd.Run(); cmdErr == nil {
			project.ScreenshotURL = fmt.Sprintf("/api/storage/screenshots/project_%d.png", project.ID)
			database.DB.Save(&project)
			return c.JSON(fiber.Map{
				"success": true,
				"data":    fiber.Map{"screenshot_url": project.ScreenshotURL},
			})
		}
	}

	// Graceful fallback: Generate sleek modern SVG preview card
	svgContent := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 800 500" width="100%%" height="100%%">
  <defs>
    <linearGradient id="grad" x1="0%%" y1="0%%" x2="100%%" y2="100%%">
      <stop offset="0%%" stop-color="#0f172a" />
      <stop offset="50%%" stop-color="#1e1b4b" />
      <stop offset="100%%" stop-color="#090d16" />
    </linearGradient>
  </defs>
  <rect width="800" height="500" rx="16" fill="url(#grad)" />
  <circle cx="40" cy="30" r="6" fill="#ef4444" />
  <circle cx="60" cy="30" r="6" fill="#f59e0b" />
  <circle cx="80" cy="30" r="6" fill="#10b981" />
  <rect x="120" y="20" width="560" height="20" rx="6" fill="#1e293b" />
  <text x="400" y="34" fill="#94a3b8" font-family="-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif" font-size="11" text-anchor="middle">%s</text>
  
  <g transform="translate(400, 240)" text-anchor="middle">
    <circle cx="0" cy="-30" r="32" fill="#3b82f6" fill-opacity="0.15" stroke="#3b82f6" stroke-width="2" />
    <path d="M-8 -38 L8 -30 L-8 -22 Z" fill="#60a5fa" />
    <text x="0" y="30" fill="#f8fafc" font-family="-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif" font-size="22" font-weight="700">%s</text>
    <text x="0" y="58" fill="#94a3b8" font-family="-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif" font-size="14">%s • %s • %s</text>
    <rect x="-80" y="80" width="160" height="30" rx="6" fill="#3b82f6" fill-opacity="0.2" stroke="#60a5fa" stroke-width="1" />
    <text x="0" y="100" fill="#93c5fd" font-family="-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif" font-size="13" font-weight="500">Port :%d</text>
  </g>
</svg>`, targetURL, project.Name, project.Framework, project.Language, project.DatabaseType, project.PreviewPort)

	_ = os.WriteFile(screenshotPath, []byte(svgContent), 0644)
	project.ScreenshotURL = fmt.Sprintf("/api/storage/screenshots/%s", screenshotFileName)
	database.DB.Save(&project)

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"screenshot_url": project.ScreenshotURL,
		},
	})
}
