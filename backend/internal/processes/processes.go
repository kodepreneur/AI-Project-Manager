package processes

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"ai-project-manager/internal/database"
	"ai-project-manager/internal/logs"
	"ai-project-manager/internal/models"

	"github.com/gofiber/fiber/v2"
)

type ProcessManager struct {
	runningCmds map[uint]*ActiveProcess
	mutex       sync.Mutex
}

type ActiveProcess struct {
	ProcessID uint
	ProjectID uint
	Cmd       *exec.Cmd
	Cancel    context.CancelFunc
}

func NewProcessManager() *ProcessManager {
	return &ProcessManager{
		runningCmds: make(map[uint]*ActiveProcess),
	}
}

type StartProcessRequest struct {
	Name    string `json:"name"`
	Command string `json:"command"`
	Port    int    `json:"port"`
}

func (m *ProcessManager) List(c *fiber.Ctx) error {
	projectID := c.Params("id")
	var procs []models.ProjectProcess
	database.DB.Where("project_id = ?", projectID).Find(&procs)

	// Check if PIDs are still alive
	for i := range procs {
		if procs[i].Status == "running" && procs[i].PID > 0 {
			if !isProcessAlive(procs[i].PID) {
				procs[i].Status = "stopped"
				procs[i].PID = 0
				database.DB.Save(&procs[i])
			}
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    procs,
	})
}

func (m *ProcessManager) Start(c *fiber.Ctx) error {
	projectID := c.Params("id")
	var pid uint
	fmt.Sscanf(projectID, "%d", &pid)

	var req StartProcessRequest
	if err := c.BodyParser(&req); err != nil || req.Command == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{"code": "INVALID_INPUT", "message": "Command is required"},
		})
	}

	if req.Name == "" {
		req.Name = "Development Server"
	}

	var project models.Project
	if err := database.DB.First(&project, pid).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{"code": "NOT_FOUND", "message": "Project not found"},
		})
	}

	proc := models.ProjectProcess{
		ProjectID: pid,
		Name:      req.Name,
		Command:   req.Command,
		Port:      req.Port,
		Status:    "running",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	database.DB.Create(&proc)

	ctx, cancel := context.WithCancel(context.Background())
	parts := strings.Fields(req.Command)
	var cmd *exec.Cmd
	if len(parts) > 1 {
		cmd = exec.CommandContext(ctx, parts[0], parts[1:]...)
	} else {
		cmd = exec.CommandContext(ctx, parts[0])
	}
	cmd.Dir = project.Path
	cmd.Env = os.Environ()

	if err := cmd.Start(); err != nil {
		proc.Status = "failed"
		database.DB.Save(&proc)
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{"code": "PROCESS_START_FAILED", "message": err.Error()},
		})
	}

	proc.PID = cmd.Process.Pid
	database.DB.Save(&proc)

	m.mutex.Lock()
	m.runningCmds[proc.ID] = &ActiveProcess{
		ProcessID: proc.ID,
		ProjectID: pid,
		Cmd:       cmd,
		Cancel:    cancel,
	}
	m.mutex.Unlock()

	// Update project status to Running
	project.Status = "Running"
	project.LastActivity = time.Now()
	database.DB.Save(&project)

	logs.Log.Info(&pid, "Project Process", fmt.Sprintf("Started process '%s' (PID %d)", proc.Name, proc.PID))

	go func() {
		_ = cmd.Wait()
		m.mutex.Lock()
		delete(m.runningCmds, proc.ID)
		m.mutex.Unlock()

		proc.Status = "stopped"
		proc.PID = 0
		database.DB.Save(&proc)
	}()

	return c.JSON(fiber.Map{
		"success": true,
		"data":    proc,
	})
}

func (m *ProcessManager) Stop(c *fiber.Ctx) error {
	procID := c.Params("procId")
	var proc models.ProjectProcess
	if err := database.DB.First(&proc, procID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{"code": "NOT_FOUND", "message": "Process not found"},
		})
	}

	m.mutex.Lock()
	if active, exists := m.runningCmds[proc.ID]; exists && active != nil {
		active.Cancel()
		if active.Cmd != nil && active.Cmd.Process != nil {
			_ = active.Cmd.Process.Kill()
		}
		delete(m.runningCmds, proc.ID)
	}
	m.mutex.Unlock()

	if proc.PID > 0 {
		if p, err := os.FindProcess(proc.PID); err == nil {
			_ = p.Kill()
		}
	}

	proc.Status = "stopped"
	proc.PID = 0
	database.DB.Save(&proc)

	logs.Log.Info(&proc.ProjectID, "Project Process", fmt.Sprintf("Stopped process '%s'", proc.Name))

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Process stopped successfully",
		"data":    proc,
	})
}

func isProcessAlive(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	// On Unix, signal 0 does not kill the process, but checks if it can receive signals
	return process.Signal(os.Signal(nil)) == nil
}
