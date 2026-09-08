package system

import (
	"ai-project-manager/internal/config"
	"ai-project-manager/internal/database"
	"ai-project-manager/internal/models"

	"github.com/gofiber/fiber/v2"
)

type SettingsHandler struct {
	cfg *config.Config
}

func NewSettingsHandler(cfg *config.Config) *SettingsHandler {
	return &SettingsHandler{cfg: cfg}
}

func (h *SettingsHandler) Get(c *fiber.Ctx) error {
	var settings []models.SystemSetting
	database.DB.Find(&settings)

	settingsMap := fiber.Map{
		"app_name":           h.cfg.AppName,
		"app_env":            h.cfg.AppEnv,
		"project_root":       h.cfg.ProjectRoot,
		"storage_root":       h.cfg.StorageRoot,
		"antigravity_binary": h.cfg.AntigravityBinary,
		"terminal_enabled":   h.cfg.TerminalEnabled,
		"screenshot_enabled": h.cfg.ScreenshotEnabled,
	}

	for _, s := range settings {
		settingsMap[s.Key] = s.Value
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    settingsMap,
	})
}

func (h *SettingsHandler) Update(c *fiber.Ctx) error {
	var body map[string]string
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "Invalid body"})
	}

	for k, v := range body {
		var setting models.SystemSetting
		if err := database.DB.Where("key = ?", k).First(&setting).Error; err == nil {
			setting.Value = v
			database.DB.Save(&setting)
		} else {
			database.DB.Create(&models.SystemSetting{
				Key:   k,
				Value: v,
			})
		}

		if k == "antigravity_binary" && v != "" {
			h.cfg.AntigravityBinary = v
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Settings updated",
	})
}
