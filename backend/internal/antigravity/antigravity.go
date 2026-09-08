package antigravity

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"ai-project-manager/internal/config"
	"ai-project-manager/internal/database"
	"ai-project-manager/internal/logs"
	"ai-project-manager/internal/models"
	"ai-project-manager/internal/websocket"

	"github.com/gofiber/fiber/v2"
)

type AntigravityService struct {
	cfg        *config.Config
	activeJobs map[uint]*ActiveSession
	mutex      sync.Mutex
}

type ActiveSession struct {
	SessionID uint
	ProjectID uint
	TaskID    uint
	Cmd       *exec.Cmd
	Cancel    context.CancelFunc
	StartTime time.Time
}

func NewAntigravityService(cfg *config.Config) *AntigravityService {
	return &AntigravityService{
		cfg:        cfg,
		activeJobs: make(map[uint]*ActiveSession),
	}
}

type StartTaskRequest struct {
	TaskID uint   `json:"task_id"`
	Prompt string `json:"prompt"`
}

type TestConnectionRequest struct {
	Binary string `json:"binary"`
}

// ResolveBinary locates the valid Antigravity CLI executable
func (s *AntigravityService) ResolveBinary(customBinary string) (string, error) {
	candidates := []string{}

	// 1. Explicit argument
	if customBinary = strings.TrimSpace(customBinary); customBinary != "" {
		candidates = append(candidates, customBinary)
	}

	// 2. Database SystemSetting override
	var dbSetting models.SystemSetting
	if database.DB != nil {
		if err := database.DB.Where("key = ?", "antigravity_binary").First(&dbSetting).Error; err == nil {
			val := strings.TrimSpace(dbSetting.Value)
			if val != "" {
				candidates = append(candidates, val)
			}
		}
	}

	// 3. Config fallback
	if s.cfg != nil && strings.TrimSpace(s.cfg.AntigravityBinary) != "" {
		candidates = append(candidates, strings.TrimSpace(s.cfg.AntigravityBinary))
	}

	// 4. Standard binary names to look up in PATH
	candidates = append(candidates, "agy", "antigravity", "antigravity-ide")

	// 5. Standard installation paths across Linux & macOS
	homeDir, _ := os.UserHomeDir()
	if homeDir != "" {
		candidates = append(candidates,
			filepath.Join(homeDir, ".local/bin/agy"),
			filepath.Join(homeDir, ".local/bin/antigravity"),
			filepath.Join(homeDir, ".antigravity/antigravity/bin/agy"),
			filepath.Join(homeDir, ".antigravity/antigravity/bin/antigravity"),
			filepath.Join(homeDir, ".antigravity-ide/antigravity-ide/bin/antigravity-ide"),
			filepath.Join(homeDir, ".antigravity-ide/antigravity-ide/bin/agy-ide"),
		)
	}

	// Linux standard & multi-user paths
	candidates = append(candidates,
		"/usr/local/bin/agy",
		"/usr/local/bin/antigravity",
		"/usr/bin/agy",
		"/usr/bin/antigravity",
		"/root/.local/bin/agy",
		"/root/.local/bin/antigravity",
	)

	// Check /home/* for user installations (e.g. /home/ubuntu/.local/bin/agy)
	if userDirs, err := filepath.Glob("/home/*/.local/bin/agy"); err == nil {
		candidates = append(candidates, userDirs...)
	}
	if userDirs, err := filepath.Glob("/home/*/.local/bin/antigravity"); err == nil {
		candidates = append(candidates, userDirs...)
	}

	// macOS standard paths
	candidates = append(candidates,
		"/opt/homebrew/bin/agy",
		"/opt/homebrew/bin/antigravity",
		"/Applications/Antigravity IDE.app/Contents/Resources/app/bin/antigravity-ide",
		"/Applications/Antigravity.app/Contents/MacOS/Antigravity",
	)

	// Test each candidate
	for _, candidate := range candidates {
		if candidate == "" {
			continue
		}

		// Check via LookPath
		if path, err := exec.LookPath(candidate); err == nil && path != "" {
			if info, statErr := os.Stat(path); statErr == nil && !info.IsDir() {
				return path, nil
			}
		}

		// Check direct file stat
		if info, statErr := os.Stat(candidate); statErr == nil && !info.IsDir() {
			if info.Mode()&0111 != 0 {
				return candidate, nil
			}
		}
	}

	return "", fmt.Errorf("no executable Antigravity CLI binary found in PATH or standard directories")
}

// TestConnection verifies real connectivity with the Antigravity CLI binary
func (s *AntigravityService) TestConnection(c *fiber.Ctx) error {
	var req TestConnectionRequest
	_ = c.BodyParser(&req)

	binPath, err := s.ResolveBinary(req.Binary)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "CLI_NOT_FOUND",
				"message": "Antigravity CLI binary was not found.",
				"details": err.Error(),
				"hint":    "Make sure 'agy' or 'antigravity-ide' is in your PATH, or configure its absolute path in Settings.",
			},
		})
	}

	start := time.Now()

	// 1. Run --version with 5s timeout
	ctxVer, cancelVer := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelVer()

	verCmd := exec.CommandContext(ctxVer, binPath, "--version")
	verOut, verErr := verCmd.CombinedOutput()
	latencyMs := time.Since(start).Milliseconds()

	if verErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":       "CLI_EXECUTION_FAILED",
				"message":    fmt.Sprintf("Failed to execute '%s --version': %v", binPath, verErr),
				"output":     string(verOut),
				"binary_path": binPath,
			},
		})
	}

	versionLines := strings.Split(strings.TrimSpace(string(verOut)), "\n")
	version := ""
	commit := ""
	arch := ""
	if len(versionLines) > 0 {
		version = strings.TrimSpace(versionLines[0])
	}
	if len(versionLines) > 1 {
		commit = strings.TrimSpace(versionLines[1])
	}
	if len(versionLines) > 2 {
		arch = strings.TrimSpace(versionLines[2])
	}

	// 2. Run --status with 6s timeout to fetch workspace stats
	ctxStatus, cancelStatus := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancelStatus()

	statusCmd := exec.CommandContext(ctxStatus, binPath, "--status")
	statusOut, _ := statusCmd.CombinedOutput()
	rawStatus := string(statusOut)

	// Extract active workspaces from status output
	workspaces := []string{}
	wsRegex := regexp.MustCompile(`Folder \(([^)]+)\)`)
	matches := wsRegex.FindAllStringSubmatch(rawStatus, -1)
	for _, m := range matches {
		if len(m) > 1 && !contains(workspaces, m[1]) {
			workspaces = append(workspaces, m[1])
		}
	}

	summary := fmt.Sprintf("Antigravity CLI %s (%s)", version, arch)
	if len(workspaces) > 0 {
		summary += fmt.Sprintf(" • %d active workspace(s)", len(workspaces))
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"connected":      true,
			"binary_path":    binPath,
			"binary_name":    filepath.Base(binPath),
			"version":        version,
			"commit":         commit,
			"arch":           arch,
			"latency_ms":     latencyMs,
			"workspaces":     workspaces,
			"status_summary": summary,
			"raw_status":     strings.TrimSpace(rawStatus),
			"timestamp":      time.Now().Format(time.RFC3339),
		},
	})
}

func contains(arr []string, item string) bool {
	for _, a := range arr {
		if a == item {
			return true
		}
	}
	return false
}

func (s *AntigravityService) GetStatus(c *fiber.Ctx) error {
	projectID := c.Params("id")
	var session models.AntigravitySession
	err := database.DB.Where("project_id = ?", projectID).Order("id DESC").First(&session).Error
	if err != nil {
		return c.JSON(fiber.Map{
			"success": true,
			"data": fiber.Map{
				"status": "idle",
			},
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    session,
	})
}

func (s *AntigravityService) Start(c *fiber.Ctx) error {
	projectID := c.Params("id")
	var pid uint
	fmt.Sscanf(projectID, "%d", &pid)

	var req StartTaskRequest
	_ = c.BodyParser(&req)

	var project models.Project
	if err := database.DB.Preload("PRD").First(&project, pid).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   fiber.Map{"code": "NOT_FOUND", "message": "Project not found"},
		})
	}

	// Verify Antigravity CLI binary exists before starting
	binPath, err := s.ResolveBinary("")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "CLI_NOT_FOUND",
				"message": "Cannot start AI task: Antigravity CLI binary not found.",
				"details": err.Error(),
				"hint":    "Please install 'agy' or configure the Antigravity binary path in Settings.",
			},
		})
	}

	var task models.Task
	if req.TaskID > 0 {
		database.DB.First(&task, req.TaskID)
	}

	s.mutex.Lock()
	if active, exists := s.activeJobs[pid]; exists && active != nil {
		s.mutex.Unlock()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   fiber.Map{"code": "ALREADY_RUNNING", "message": "An Antigravity session is already active for this project"},
		})
	}
	s.mutex.Unlock()

	// Ensure project physical directory exists
	_ = os.MkdirAll(project.Path, 0755)

	// Construct AI Context Package
	contextPayload := s.buildContextPayload(&project, &task, req.Prompt)

	session := models.AntigravitySession{
		ProjectID:   pid,
		TaskID:      &task.ID,
		Status:      "running",
		Command:     fmt.Sprintf("%s in %s", filepath.Base(binPath), project.Path),
		CurrentTask: task.Title,
		StartTime:   time.Now(),
	}
	database.DB.Create(&session)

	// Update task status to in_progress
	if task.ID > 0 {
		task.Status = "in_progress"
		task.AIStatus = "running"
		database.DB.Save(&task)
	}

	go s.runRealSession(pid, &session, &project, &task, binPath, contextPayload, req.Prompt)

	logs.Log.Info(&pid, "Antigravity", fmt.Sprintf("Started real Antigravity AI task: %s (PID will be assigned)", task.Title))

	return c.JSON(fiber.Map{
		"success": true,
		"data":    session,
	})
}

func (s *AntigravityService) Stop(c *fiber.Ctx) error {
	projectID := c.Params("id")
	var pid uint
	fmt.Sscanf(projectID, "%d", &pid)

	s.mutex.Lock()
	active, exists := s.activeJobs[pid]
	if exists && active != nil && active.Cancel != nil {
		active.Cancel()
		delete(s.activeJobs, pid)
	}
	s.mutex.Unlock()

	var session models.AntigravitySession
	if err := database.DB.Where("project_id = ? AND status = 'running'", pid).First(&session).Error; err == nil {
		now := time.Now()
		session.Status = "paused"
		session.EndTime = &now
		database.DB.Save(&session)
	}

	logs.Log.Warn(&pid, "Antigravity", "Antigravity session stopped by user")

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Antigravity session stopped",
	})
}

func (s *AntigravityService) runRealSession(pid uint, session *models.AntigravitySession, project *models.Project, task *models.Task, binPath string, contextPayload string, userPrompt string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.mutex.Lock()
	s.activeJobs[pid] = &ActiveSession{
		SessionID: session.ID,
		ProjectID: pid,
		TaskID:    task.ID,
		Cancel:    cancel,
		StartTime: time.Now(),
	}
	s.mutex.Unlock()

	defer func() {
		s.mutex.Lock()
		delete(s.activeJobs, pid)
		s.mutex.Unlock()
	}()

	// Write task context file to project path for Antigravity reference
	contextDir := filepath.Join(project.Path, ".antigravity")
	_ = os.MkdirAll(contextDir, 0755)
	contextFile := filepath.Join(contextDir, fmt.Sprintf("task_%d_prompt.md", task.ID))
	_ = os.WriteFile(contextFile, []byte(contextPayload), 0644)

	// Broadcast task startup banner
	websocket.GlobalHub.Broadcast("antigravity_output", &project.ID, fiber.Map{
		"session_id": session.ID,
		"line":       fmt.Sprintf("🚀 Dispatched task to Antigravity CLI (%s): %s", filepath.Base(binPath), task.Title),
		"timestamp":  time.Now().Format("15:04:05"),
	})
	websocket.GlobalHub.Broadcast("antigravity_output", &project.ID, fiber.Map{
		"session_id": session.ID,
		"line":       fmt.Sprintf("📂 Working Directory: %s", project.Path),
		"timestamp":  time.Now().Format("15:04:05"),
	})

	// Construct agent task prompt
	agentPrompt := fmt.Sprintf("Implement task: %s\nAcceptance Criteria: %s", task.Title, task.AcceptanceCriteria)
	if userPrompt != "" {
		agentPrompt += fmt.Sprintf("\nUser Instructions: %s", userPrompt)
	}

	// Prepare command arguments for Antigravity CLI
	// Runs chat session in agent mode with the prompt and attached context file
	args := []string{"chat", "-m", "agent", agentPrompt, "-a", contextFile}
	cmd := exec.CommandContext(ctx, binPath, args...)
	cmd.Dir = project.Path
	cmd.Env = append(os.Environ(), "ANTIGRAVITY_AUTOPILOT=1")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		s.finishSession(session, task, "failed", 1, fmt.Sprintf("Failed to initialize stdout pipe: %v", err))
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		s.finishSession(session, task, "failed", 1, fmt.Sprintf("Failed to initialize stderr pipe: %v", err))
		return
	}

	if err := cmd.Start(); err != nil {
		s.finishSession(session, task, "failed", 1, fmt.Sprintf("Failed to launch Antigravity CLI process: %v", err))
		return
	}

	session.PID = cmd.Process.Pid
	database.DB.Save(session)

	websocket.GlobalHub.Broadcast("antigravity_output", &project.ID, fiber.Map{
		"session_id": session.ID,
		"line":       fmt.Sprintf("⚡ Process spawned with PID %d", session.PID),
		"timestamp":  time.Now().Format("15:04:05"),
	})

	var outputBuf strings.Builder
	reader := io.MultiReader(stdout, stderr)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()
		outputBuf.WriteString(line + "\n")

		// Broadcast real output lines
		websocket.GlobalHub.Broadcast("antigravity_output", &project.ID, fiber.Map{
			"session_id": session.ID,
			"line":       line,
			"timestamp":  time.Now().Format("15:04:05"),
		})
	}

	exitCode := 0
	if err := cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = 1
		}
	}

	session.Stdout = outputBuf.String()
	if exitCode == 0 {
		websocket.GlobalHub.Broadcast("antigravity_output", &project.ID, fiber.Map{
			"session_id": session.ID,
			"line":       "✅ Antigravity CLI task completed successfully (Exit Code 0).",
			"timestamp":  time.Now().Format("15:04:05"),
		})
		s.finishSession(session, task, "completed", 0, "")
	} else {
		errMsg := fmt.Sprintf("Antigravity process finished with exit code %d", exitCode)
		websocket.GlobalHub.Broadcast("antigravity_output", &project.ID, fiber.Map{
			"session_id": session.ID,
			"line":       fmt.Sprintf("❌ %s", errMsg),
			"timestamp":  time.Now().Format("15:04:05"),
		})
		s.finishSession(session, task, "failed", exitCode, errMsg)
	}
}

func (s *AntigravityService) finishSession(session *models.AntigravitySession, task *models.Task, status string, exitCode int, errMsg string) {
	now := time.Now()
	session.Status = status
	session.ExitCode = exitCode
	session.EndTime = &now
	if errMsg != "" {
		session.Stderr = errMsg
	}
	database.DB.Save(session)

	if task != nil && task.ID > 0 {
		task.AIStatus = status
		if status == "completed" {
			task.Status = "review" // Move to human review
		}
		database.DB.Save(task)
	}

	logs.Log.Info(&session.ProjectID, "Antigravity", fmt.Sprintf("Session %d finished with status: %s", session.ID, status))

	websocket.GlobalHub.Broadcast("antigravity_finished", &session.ProjectID, fiber.Map{
		"session_id": session.ID,
		"status":     status,
		"task_id":    task.ID,
	})
}

func (s *AntigravityService) buildContextPayload(project *models.Project, task *models.Task, userPrompt string) string {
	prdContent := "No PRD available."
	if project.PRD != nil {
		prdContent = project.PRD.Content
	}

	return fmt.Sprintf(`# PROJECT CONTEXT
- **Project**: %s (%s)
- **Framework**: %s | **Language**: %s | **Frontend**: %s | **DB**: %s
- **Path**: %s

## TASK CONTEXT
- **Title**: %s
- **Phase**: %s
- **Priority**: %s
- **Complexity**: %s

### Description:
%s

### Acceptance Criteria:
%s

### AI Instructions:
%s

### Custom User Instructions:
%s

## PRD SUMMARY
%s
`, project.Name, project.Slug, project.Framework, project.Language, project.FrontendTech, project.DatabaseType, project.Path,
		task.Title, task.PhaseName, task.Priority, task.Complexity, task.Description, task.AcceptanceCriteria, task.AIInstructions, userPrompt, prdContent)
}
