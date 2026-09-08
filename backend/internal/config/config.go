package config

import (
	"os"
	"path/filepath"
	"strconv"
)

type Config struct {
	AppName            string
	AppEnv             string
	Port               int
	DatabaseURL        string
	ProjectRoot        string
	StorageRoot        string
	AntigravityBinary  string
	JWTSecret          string
	TerminalEnabled    bool
	ScreenshotEnabled  bool
}

func Load() *Config {
	cwd, _ := os.Getwd()
	appRoot := cwd
	if filepath.Base(cwd) == "backend" {
		appRoot = filepath.Dir(cwd)
	}

	// Fallbacks
	projectRoot := getEnv("PROJECT_ROOT", filepath.Join(appRoot, "projects"))
	storageRoot := getEnv("STORAGE_ROOT", filepath.Join(appRoot, "storage"))

	_ = os.MkdirAll(projectRoot, 0755)
	_ = os.MkdirAll(filepath.Join(storageRoot, "logs"), 0755)
	_ = os.MkdirAll(filepath.Join(storageRoot, "screenshots"), 0755)
	_ = os.MkdirAll(filepath.Join(storageRoot, "db"), 0755)

	dbURL := getEnv("DATABASE_URL", filepath.Join(storageRoot, "db", "app.db"))

	port, err := strconv.Atoi(getEnv("APP_PORT", "8080"))
	if err != nil {
		port = 8080
	}

	return &Config{
		AppName:           getEnv("APP_NAME", "AI Project Manager"),
		AppEnv:            getEnv("APP_ENV", "development"),
		Port:              port,
		DatabaseURL:       dbURL,
		ProjectRoot:       projectRoot,
		StorageRoot:       storageRoot,
		AntigravityBinary: getEnv("ANTIGRAVITY_BINARY", "antigravity"),
		JWTSecret:         getEnv("JWT_SECRET", "super-secret-antigravity-dev-jwt-key-2026"),
		TerminalEnabled:   getEnv("TERMINAL_ENABLED", "true") == "true",
		ScreenshotEnabled: getEnv("SCREENSHOT_ENABLED", "true") == "true",
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}
