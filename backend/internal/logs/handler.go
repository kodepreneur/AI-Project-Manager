package logs

import (
	"fmt"
	"strings"

	"ai-project-manager/internal/database"
	"ai-project-manager/internal/models"

	"github.com/gofiber/fiber/v2"
)

type LogHandler struct{}

func NewLogHandler() *LogHandler {
	return &LogHandler{}
}

func (h *LogHandler) List(c *fiber.Ctx) error {
	category := c.Query("category")
	level := c.Query("level")
	search := c.Query("search")
	projectID := c.Query("project_id")
	limit := c.QueryInt("limit", 100)

	query := database.DB.Model(&models.ProjectLog{})

	if category != "" && category != "All" {
		query = query.Where("category = ?", category)
	}
	if level != "" && level != "All" {
		query = query.Where("level = ?", level)
	}
	if projectID != "" {
		query = query.Where("project_id = ?", projectID)
	}
	if search != "" {
		query = query.Where("message LIKE ? OR category LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	var logsList []models.ProjectLog
	query.Order("created_at DESC").Limit(limit).Find(&logsList)

	return c.JSON(fiber.Map{
		"success": true,
		"data":    logsList,
	})
}

func (h *LogHandler) Clear(c *fiber.Ctx) error {
	projectID := c.Query("project_id")
	if projectID != "" {
		database.DB.Where("project_id = ?", projectID).Delete(&models.ProjectLog{})
	} else {
		database.DB.Exec("DELETE FROM project_logs")
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Logs cleared successfully",
	})
}

func (h *LogHandler) Download(c *fiber.Ctx) error {
	var logsList []models.ProjectLog
	database.DB.Order("created_at ASC").Find(&logsList)

	var builder strings.Builder
	for _, l := range logsList {
		builder.WriteString(fmt.Sprintf("[%s] [%s] [%s] %s\n", l.CreatedAt.Format("2006-01-02 15:04:05"), l.Level, l.Category, l.Message))
	}

	c.Set("Content-Disposition", "attachment; filename=\"ai-project-manager-logs.txt\"")
	c.Set("Content-Type", "text/plain; charset=utf-8")
	return c.SendString(builder.String())
}
