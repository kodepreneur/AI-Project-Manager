package database

import (
	"log"
	"os"
	"path/filepath"

	"ai-project-manager/internal/config"
	"ai-project-manager/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init(cfg *config.Config) *gorm.DB {
	dbDir := filepath.Dir(cfg.DatabaseURL)
	_ = os.MkdirAll(dbDir, 0755)

	db, err := gorm.Open(sqlite.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate tables
	err = db.AutoMigrate(
		&models.User{},
		&models.Project{},
		&models.PRD{},
		&models.PRDVersion{},
		&models.Phase{},
		&models.Task{},
		&models.AntigravitySession{},
		&models.ProjectProcess{},
		&models.DatabaseConnection{},
		&models.ProjectLog{},
		&models.SystemSetting{},
	)
	if err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	DB = db
	return db
}
