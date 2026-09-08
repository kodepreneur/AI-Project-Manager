package terminal

import (
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"ai-project-manager/internal/config"
	"ai-project-manager/internal/database"
	"ai-project-manager/internal/logs"
	"ai-project-manager/internal/models"

	"github.com/creack/pty"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type TerminalService struct {
	cfg *config.Config
}

func NewTerminalService(cfg *config.Config) *TerminalService {
	return &TerminalService{cfg: cfg}
}

type ResizeMsg struct {
	Type string `json:"type"`
	Cols uint16 `json:"cols"`
	Rows uint16 `json:"rows"`
}

func (s *TerminalService) WebSocketHandler() fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		projectID := c.Params("id")
		var project models.Project
		if err := database.DB.First(&project, projectID).Error; err != nil {
			_ = c.WriteMessage(websocket.TextMessage, []byte("Error: Project not found\r\n"))
			c.Close()
			return
		}

		workDir := project.Path
		if _, err := os.Stat(workDir); os.IsNotExist(err) {
			_ = os.MkdirAll(workDir, 0755)
		}

		// Choose shell
		shell := os.Getenv("SHELL")
		if shell == "" {
			shell = "/bin/bash"
			if _, err := os.Stat(shell); os.IsNotExist(err) {
				shell = "/bin/sh"
			}
		}

		cmd := exec.Command(shell)
		cmd.Dir = workDir
		cmd.Env = append(os.Environ(),
			"TERM=xterm-256color",
			"COLORTERM=truecolor",
			"PROJECT_NAME="+project.Name,
			"PROJECT_PATH="+project.Path,
		)

		ptmx, err := pty.Start(cmd)
		if err != nil {
			_ = c.WriteMessage(websocket.TextMessage, []byte("Error starting PTY: "+err.Error()+"\r\n"))
			c.Close()
			return
		}
		defer func() {
			_ = ptmx.Close()
			_ = cmd.Process.Kill()
		}()

		logs.Log.Info(&project.ID, "Terminal", "Interactive terminal session started for: "+filepath.Base(workDir))

		var once sync.Once
		closeBoth := func() {
			once.Do(func() {
				_ = ptmx.Close()
				_ = c.Close()
			})
		}

		// Read from PTY -> Write to WebSocket
		go func() {
			defer closeBoth()
			buf := make([]byte, 4096)
			for {
				n, err := ptmx.Read(buf)
				if err != nil {
					return
				}
				if err := c.WriteMessage(websocket.BinaryMessage, buf[:n]); err != nil {
					return
				}
			}
		}()

		// Read from WebSocket -> Write to PTY
		for {
			msgType, msg, err := c.ReadMessage()
			if err != nil {
				break
			}

			// Check if message is a JSON resize event
			if msgType == websocket.TextMessage && len(msg) > 0 && msg[0] == '{' {
				var resize ResizeMsg
				if err := json.Unmarshal(msg, &resize); err == nil && resize.Type == "resize" {
					_ = pty.Setsize(ptmx, &pty.Winsize{
						Cols: resize.Cols,
						Rows: resize.Rows,
					})
					continue
				}
			}

			// Forward raw input to PTY
			if _, err := ptmx.Write(msg); err != nil {
				if err == io.EOF {
					break
				}
			}
		}
	})
}
