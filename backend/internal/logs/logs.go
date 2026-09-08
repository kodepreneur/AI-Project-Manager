package logs

import (
	"encoding/json"
	"fmt"
	"time"

	"ai-project-manager/internal/database"
	"ai-project-manager/internal/models"
	"ai-project-manager/internal/websocket"
)

type Logger struct{}

var Log = &Logger{}

func (l *Logger) Record(projectID *uint, category, level, message string, metadata interface{}) {
	metaStr := ""
	if metadata != nil {
		if b, err := json.Marshal(metadata); err == nil {
			metaStr = string(b)
		}
	}

	logEntry := models.ProjectLog{
		ProjectID: projectID,
		Category:  category,
		Level:     level,
		Message:   message,
		Metadata:  metaStr,
		CreatedAt: time.Now(),
	}

	// Persist to DB
	if database.DB != nil {
		_ = database.DB.Create(&logEntry)
	}

	fmt.Printf("[%s] [%s] %s: %s\n", time.Now().Format("15:04:05"), level, category, message)

	// Broadcast via WebSocket
	websocket.GlobalHub.Broadcast("log", projectID, logEntry)
}

func (l *Logger) Info(projectID *uint, category, message string, meta ...interface{}) {
	var m interface{}
	if len(meta) > 0 {
		m = meta[0]
	}
	l.Record(projectID, category, "INFO", message, m)
}

func (l *Logger) Warn(projectID *uint, category, message string, meta ...interface{}) {
	var m interface{}
	if len(meta) > 0 {
		m = meta[0]
	}
	l.Record(projectID, category, "WARN", message, m)
}

func (l *Logger) Error(projectID *uint, category, message string, meta ...interface{}) {
	var m interface{}
	if len(meta) > 0 {
		m = meta[0]
	}
	l.Record(projectID, category, "ERROR", message, m)
}
