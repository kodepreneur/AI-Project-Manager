package git

import (
	"os"
	"os/exec"
	"strings"

	"ai-project-manager/internal/database"
	"ai-project-manager/internal/logs"
	"ai-project-manager/internal/models"

	"github.com/gofiber/fiber/v2"
)

type GitService struct{}

func NewGitService() *GitService {
	return &GitService{}
}

type CommitRequest struct {
	Message string `json:"message"`
}

type BranchRequest struct {
	Branch string `json:"branch"`
}

func (s *GitService) Status(c *fiber.Ctx) error {
	projectID := c.Params("id")
	var project models.Project
	if err := database.DB.First(&project, projectID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "error": "Project not found"})
	}

	if _, err := os.Stat(project.Path); os.IsNotExist(err) {
		return c.JSON(fiber.Map{"success": false, "error": "Path does not exist"})
	}

	// Current branch
	branchCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	branchCmd.Dir = project.Path
	branchOut, _ := branchCmd.Output()
	currentBranch := strings.TrimSpace(string(branchOut))
	if currentBranch == "" {
		currentBranch = "main"
	}

	// Git status porcelain
	statusCmd := exec.Command("git", "status", "--porcelain")
	statusCmd.Dir = project.Path
	statusOut, _ := statusCmd.Output()

	lines := strings.Split(strings.TrimSpace(string(statusOut)), "\n")
	modified, added, deleted, untracked := 0, 0, 0, 0
	files := []fiber.Map{}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) < 3 {
			continue
		}
		statusCode := line[:2]
		fileName := strings.TrimSpace(line[3:])

		action := "modified"
		if strings.Contains(statusCode, "?") {
			untracked++
			action = "untracked"
		} else if strings.Contains(statusCode, "A") {
			added++
			action = "added"
		} else if strings.Contains(statusCode, "D") {
			deleted++
			action = "deleted"
		} else {
			modified++
		}

		files = append(files, fiber.Map{
			"action": action,
			"file":   fileName,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"current_branch": currentBranch,
			"modified":       modified,
			"added":          added,
			"deleted":        deleted,
			"untracked":      untracked,
			"total_changes":  len(files),
			"files":          files,
		},
	})
}

func (s *GitService) Diff(c *fiber.Ctx) error {
	projectID := c.Params("id")
	var project models.Project
	if err := database.DB.First(&project, projectID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "error": "Project not found"})
	}

	diffCmd := exec.Command("git", "diff")
	diffCmd.Dir = project.Path
	diffOut, _ := diffCmd.Output()

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"diff": string(diffOut),
		},
	})
}

func (s *GitService) Commits(c *fiber.Ctx) error {
	projectID := c.Params("id")
	var project models.Project
	if err := database.DB.First(&project, projectID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "error": "Project not found"})
	}

	logCmd := exec.Command("git", "log", "-n", "20", "--pretty=format:%H|%an|%ar|%s")
	logCmd.Dir = project.Path
	out, err := logCmd.Output()
	if err != nil {
		return c.JSON(fiber.Map{"success": true, "data": []fiber.Map{}})
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var commits []fiber.Map
	for _, l := range lines {
		parts := strings.Split(l, "|")
		if len(parts) >= 4 {
			commits = append(commits, fiber.Map{
				"hash":    parts[0][:8],
				"author":  parts[1],
				"time":    parts[2],
				"message": parts[3],
			})
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    commits,
	})
}

func (s *GitService) Commit(c *fiber.Ctx) error {
	projectID := c.Params("id")
	var project models.Project
	if err := database.DB.First(&project, projectID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "error": "Project not found"})
	}

	var req CommitRequest
	if err := c.BodyParser(&req); err != nil || req.Message == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "Commit message is required"})
	}

	addCmd := exec.Command("git", "add", "-A")
	addCmd.Dir = project.Path
	_ = addCmd.Run()

	commitCmd := exec.Command("git", "commit", "-m", req.Message)
	commitCmd.Dir = project.Path
	out, err := commitCmd.CombinedOutput()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": string(out)})
	}

	logs.Log.Info(&project.ID, "Git", "Git commit created: "+req.Message)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Committed successfully",
		"output":  string(out),
	})
}
