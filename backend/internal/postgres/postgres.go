package postgres

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type PostgresService struct{}

func NewPostgresService() *PostgresService {
	return &PostgresService{}
}

type ConnectRequest struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type QueryRequest struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	Query    string `json:"query"`
}

func (s *PostgresService) Status(c *fiber.Ctx) error {
	host := c.Query("host", "127.0.0.1")
	port := c.QueryInt("port", 5432)

	timeout := 1 * time.Second
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), timeout)
	running := err == nil
	if conn != nil {
		_ = conn.Close()
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"running": running,
			"host":    host,
			"port":    port,
			"engine":  "PostgreSQL",
		},
	})
}

func (s *PostgresService) TestConnection(c *fiber.Ctx) error {
	var req ConnectRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{"code": "INVALID_INPUT", "message": "Invalid request body"},
		})
	}

	if req.Host == "" {
		req.Host = "127.0.0.1"
	}
	if req.Port == 0 {
		req.Port = 5432
	}

	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", req.Host, req.Port), timeout)
	if err != nil {
		return c.JSON(fiber.Map{
			"success": false,
			"message": fmt.Sprintf("Cannot reach PostgreSQL server on %s:%d: %v", req.Host, req.Port, err),
		})
	}
	_ = conn.Close()

	return c.JSON(fiber.Map{
		"success": true,
		"message": fmt.Sprintf("Successfully reached PostgreSQL server at %s:%d", req.Host, req.Port),
	})
}

func (s *PostgresService) ExecuteQuery(c *fiber.Ctx) error {
	var req QueryRequest
	if err := c.BodyParser(&req); err != nil || req.Query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{"code": "INVALID_INPUT", "message": "SQL query is required"},
		})
	}

	if req.Host == "" {
		req.Host = "127.0.0.1"
	}
	if req.Port == 0 {
		req.Port = 5432
	}

	start := time.Now()
	psqlPath, err := exec.LookPath("psql")
	if err == nil && psqlPath != "" {
		cmd := exec.Command(psqlPath, "-h", req.Host, "-p", fmt.Sprintf("%d", req.Port), "-U", req.User, "-d", req.Database, "-c", req.Query)
		if req.Password != "" {
			cmd.Env = append(cmd.Environ(), "PGPASSWORD="+req.Password)
		}
		out, runErr := cmd.CombinedOutput()
		execDuration := time.Since(start).String()

		if runErr != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "QUERY_ERROR",
					"message": strings.TrimSpace(string(out)),
				},
			})
		}

		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		var rows [][]string
		for _, l := range lines {
			if strings.TrimSpace(l) != "" {
				rows = append(rows, []string{l})
			}
		}

		return c.JSON(fiber.Map{
			"success": true,
			"data": fiber.Map{
				"columns":        []string{"Output"},
				"rows":           rows,
				"affected_rows":  len(rows),
				"execution_time": execDuration,
			},
		})
	}

	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"success": false,
		"error": fiber.Map{
			"code":    "PSQL_CLI_UNAVAILABLE",
			"message": "PostgreSQL CLI tool ('psql') was not found on the host system to execute direct SQL.",
		},
	})
}
