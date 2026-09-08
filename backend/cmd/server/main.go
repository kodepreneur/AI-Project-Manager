package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"ai-project-manager/internal/antigravity"
	"ai-project-manager/internal/auth"
	"ai-project-manager/internal/config"
	"ai-project-manager/internal/database"
	"ai-project-manager/internal/env"
	"ai-project-manager/internal/git"
	"ai-project-manager/internal/logs"
	"ai-project-manager/internal/mysql"
	"ai-project-manager/internal/php"
	"ai-project-manager/internal/postgres"
	"ai-project-manager/internal/preview"
	"ai-project-manager/internal/processes"
	"ai-project-manager/internal/prd"
	"ai-project-manager/internal/projects"
	"ai-project-manager/internal/system"
	"ai-project-manager/internal/terminal"
	"ai-project-manager/internal/todos"
	"ai-project-manager/internal/websocket"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	fws "github.com/gofiber/websocket/v2"
)

func main() {
	cfg := config.Load()
	database.Init(cfg)

	// Launch WebSocket Hub
	go websocket.GlobalHub.Run()

	app := fiber.New(fiber.Config{
		AppName: cfg.AppName,
	})

	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, HEAD, PUT, DELETE, PATCH, OPTIONS",
		AllowCredentials: false,
	}))

	// Services & Handlers
	authSvc := auth.NewAuthService(cfg)
	projectSvc := projects.NewProjectService(cfg)
	prdSvc := prd.NewPRDService()
	todoSvc := todos.NewTodoService()
	antiSvc := antigravity.NewAntigravityService(cfg)
	terminalSvc := terminal.NewTerminalService(cfg)
	processMgr := processes.NewProcessManager()
	previewSvc := preview.NewPreviewService(cfg)
	phpSvc := php.NewPHPService()
	mysqlSvc := mysql.NewMySQLService()
	pgSvc := postgres.NewPostgresService()
	gitSvc := git.NewGitService()
	envSvc := env.NewEnvService()
	sysSvc := system.NewSystemService()
	logHandler := logs.NewLogHandler()
	settingsHandler := system.NewSettingsHandler(cfg)

	// WebSocket Routes
	app.Use("/ws", func(c *fiber.Ctx) error {
		if fws.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/events", fws.New(func(c *fws.Conn) {
		websocket.GlobalHub.Register(c)
		defer websocket.GlobalHub.Unregister(c)
		for {
			if _, _, err := c.ReadMessage(); err != nil {
				break
			}
		}
	}))

	app.Get("/ws/projects/:id/terminal", terminalSvc.WebSocketHandler())

	// Static storage serving
	app.Static("/api/storage", cfg.StorageRoot)

	// REST API Routes
	api := app.Group("/api")

	// Public Auth
	api.Get("/auth/setup-status", authSvc.SetupStatus)
	api.Post("/auth/setup", authSvc.Setup)
	api.Post("/auth/login", authSvc.Login)

	// Protected routes
	protected := api.Group("", authSvc.JWTMiddleware())

	// Auth & Profile
	protected.Get("/auth/me", authSvc.Me)
	protected.Post("/auth/password", authSvc.ChangePassword)

	// Dashboard & System
	protected.Get("/system/metrics", sysSvc.Metrics)
	protected.Get("/system/overview", sysSvc.DashboardOverview)
	protected.Get("/settings", settingsHandler.Get)
	protected.Put("/settings", settingsHandler.Update)

	// Projects
	protected.Get("/projects", projectSvc.List)
	protected.Post("/projects", projectSvc.Create)
	protected.Get("/projects/:id", projectSvc.Get)
	protected.Put("/projects/:id", projectSvc.Update)
	protected.Delete("/projects/:id", projectSvc.Delete)
	protected.Post("/projects/:id/archive", projectSvc.Archive)
	protected.Get("/projects/:id/health", projectSvc.HealthCheck)

	// PRD
	protected.Get("/projects/:id/prd", prdSvc.Get)
	protected.Post("/projects/:id/prd", prdSvc.Save)
	protected.Post("/projects/:id/prd/generate", prdSvc.GenerateAI)
	protected.Post("/projects/:id/prd/approve", prdSvc.Approve)
	protected.Get("/projects/:id/prd/export", prdSvc.ExportMarkdown)

	// Tasks & Kanban
	protected.Get("/projects/:id/tasks", todoSvc.List)
	protected.Post("/projects/:id/tasks", todoSvc.Create)
	protected.Post("/projects/:id/tasks/generate-from-prd", todoSvc.GenerateFromPRD)
	protected.Put("/tasks/:taskId", todoSvc.Update)
	protected.Put("/tasks/:taskId/status", todoSvc.UpdateStatus)
	protected.Delete("/tasks/:taskId", todoSvc.Delete)

	// Antigravity AI
	protected.Get("/antigravity/test", antiSvc.TestConnection)
	protected.Post("/antigravity/test", antiSvc.TestConnection)
	protected.Get("/projects/:id/antigravity/status", antiSvc.GetStatus)
	protected.Post("/projects/:id/antigravity/test", antiSvc.TestConnection)
	protected.Post("/projects/:id/antigravity/start", antiSvc.Start)
	protected.Post("/projects/:id/antigravity/stop", antiSvc.Stop)

	// Processes
	protected.Get("/projects/:id/processes", processMgr.List)
	protected.Post("/projects/:id/processes/start", processMgr.Start)
	protected.Post("/processes/:procId/stop", processMgr.Stop)

	// Preview & Screenshots
	protected.Get("/projects/:id/preview", previewSvc.GetPreviewInfo)
	protected.Post("/projects/:id/preview/screenshot", previewSvc.CaptureScreenshot)

	// Git
	protected.Get("/projects/:id/git/status", gitSvc.Status)
	protected.Get("/projects/:id/git/diff", gitSvc.Diff)
	protected.Get("/projects/:id/git/commits", gitSvc.Commits)
	protected.Post("/projects/:id/git/commit", gitSvc.Commit)

	// Environment Variables
	protected.Get("/projects/:id/env", envSvc.Get)
	protected.Post("/projects/:id/env", envSvc.Save)

	// Runtimes & Infrastructure
	protected.Get("/php/diagnostics", phpSvc.Diagnostics)
	protected.Post("/php/run", phpSvc.RunCommand)
	protected.Get("/mysql/status", mysqlSvc.Status)
	protected.Post("/mysql/test", mysqlSvc.TestConnection)
	protected.Post("/mysql/query", mysqlSvc.ExecuteQuery)
	protected.Get("/postgres/status", pgSvc.Status)
	protected.Post("/postgres/test", pgSvc.TestConnection)
	protected.Post("/postgres/query", pgSvc.ExecuteQuery)

	// Logs
	protected.Get("/logs", logHandler.List)
	protected.Delete("/logs", logHandler.Clear)
	protected.Get("/logs/download", logHandler.Download)

	// Serve Frontend Static files if built
	candidates := []string{
		filepath.Join(filepath.Dir(cfg.ProjectRoot), "frontend", "dist"),
		filepath.Join(".", "frontend", "dist"),
		filepath.Join("..", "frontend", "dist"),
		"/var/lib/ai-project-manager/frontend/dist",
	}

	for _, distPath := range candidates {
		if info, err := os.Stat(distPath); err == nil && info.IsDir() {
			app.Static("/", distPath)
			app.Get("/*", func(c *fiber.Ctx) error {
				return c.SendFile(filepath.Join(distPath, "index.html"))
			})
			log.Printf("Serving frontend UI from: %s\n", distPath)
			break
		}
	}

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("🚀 %s running on http://localhost%s\n", cfg.AppName, addr)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
