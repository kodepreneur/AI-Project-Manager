package php

import (
	"os/exec"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type PHPService struct{}

func NewPHPService() *PHPService {
	return &PHPService{}
}

type PHPVersionInfo struct {
	Version   string `json:"version"`
	Path      string `json:"path"`
	Installed bool   `json:"installed"`
	IsDefault bool   `json:"is_default"`
}

func (s *PHPService) Diagnostics(c *fiber.Ctx) error {
	candidateVersions := []string{"8.1", "8.2", "8.3", "8.4"}
	var versions []PHPVersionInfo

	defaultCLI, _ := exec.LookPath("php")
	defaultVerStr := ""
	if defaultCLI != "" {
		if out, err := exec.Command(defaultCLI, "-r", "echo PHP_VERSION;").Output(); err == nil {
			defaultVerStr = string(out)
		}
	}

	for _, ver := range candidateVersions {
		binName := "php" + ver
		path, err := exec.LookPath(binName)
		isInst := err == nil
		isDef := false

		if isInst && defaultCLI != "" && strings.HasPrefix(defaultVerStr, ver) {
			isDef = true
		}

		versions = append(versions, PHPVersionInfo{
			Version:   ver,
			Path:      path,
			Installed: isInst,
			IsDefault: isDef,
		})
	}

	// Also add the default if not caught
	if defaultCLI != "" && defaultVerStr != "" {
		hasDefault := false
		for _, v := range versions {
			if v.IsDefault {
				hasDefault = true
				break
			}
		}
		if !hasDefault {
			versions = append([]PHPVersionInfo{{
				Version:   defaultVerStr,
				Path:      defaultCLI,
				Installed: true,
				IsDefault: true,
			}}, versions...)
		}
	}

	// Get modules
	var modules []string
	if defaultCLI != "" {
		if out, err := exec.Command(defaultCLI, "-m").Output(); err == nil {
			lines := strings.Split(string(out), "\n")
			for _, l := range lines {
				l = strings.TrimSpace(l)
				if l != "" && !strings.HasPrefix(l, "[") {
					modules = append(modules, l)
				}
			}
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"default_version": defaultVerStr,
			"default_path":    defaultCLI,
			"versions":        versions,
			"modules":         modules,
			"installed":       defaultCLI != "",
		},
	})
}

func (s *PHPService) RunCommand(c *fiber.Ctx) error {
	type RunReq struct {
		Command string `json:"command"`
	}
	var req RunReq
	if err := c.BodyParser(&req); err != nil || req.Command == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{"code": "INVALID_INPUT", "message": "PHP command or script required"},
		})
	}

	phpPath, err := exec.LookPath("php")
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{"code": "PHP_NOT_FOUND", "message": "PHP CLI is not installed on this server"},
		})
	}

	cmd := exec.Command(phpPath, "-r", req.Command)
	out, runErr := cmd.CombinedOutput()

	return c.JSON(fiber.Map{
		"success": runErr == nil,
		"data": fiber.Map{
			"output": string(out),
		},
	})
}
