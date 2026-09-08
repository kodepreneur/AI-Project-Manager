package mysql

import (
	"database/sql"
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type MySQLService struct{}

func NewMySQLService() *MySQLService {
	return &MySQLService{}
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

func (s *MySQLService) Status(c *fiber.Ctx) error {
	host := c.Query("host", "127.0.0.1")
	port := c.QueryInt("port", 3306)

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
			"engine":  "MySQL / MariaDB",
		},
	})
}

func (s *MySQLService) TestConnection(c *fiber.Ctx) error {
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
		req.Port = 3306
	}

	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", req.Host, req.Port), timeout)
	if err != nil {
		return c.JSON(fiber.Map{
			"success": false,
			"message": fmt.Sprintf("Cannot connect to MySQL server on %s:%d: %v", req.Host, req.Port, err),
		})
	}
	_ = conn.Close()

	return c.JSON(fiber.Map{
		"success": true,
		"message": fmt.Sprintf("Successfully reached MySQL server at %s:%d", req.Host, req.Port),
	})
}

func (s *MySQLService) ExecuteQuery(c *fiber.Ctx) error {
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
		req.Port = 3306
	}

	start := time.Now()
	mysqlPath, err := exec.LookPath("mysql")
	if err == nil && mysqlPath != "" {
		args := []string{
			"-h", req.Host,
			"-P", fmt.Sprintf("%d", req.Port),
			"-u", req.User,
		}
		if req.Password != "" {
			args = append(args, "-p"+req.Password)
		}
		if req.Database != "" {
			args = append(args, req.Database)
		}
		args = append(args, "-e", req.Query)

		cmd := exec.Command(mysqlPath, args...)
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
		var columns []string
		var rows [][]string

		if len(lines) > 0 {
			columns = strings.Split(lines[0], "\t")
			for _, line := range lines[1:] {
				if strings.TrimSpace(line) != "" {
					rows = append(rows, strings.Split(line, "\t"))
				}
			}
		}

		return c.JSON(fiber.Map{
			"success": true,
			"data": fiber.Map{
				"columns":        columns,
				"rows":           rows,
				"affected_rows":  len(rows),
				"execution_time": execDuration,
			},
		})
	}

	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"success": false,
		"error": fiber.Map{
			"code":    "MYSQL_CLI_UNAVAILABLE",
			"message": "MySQL command line tool was not found on the host system to execute direct SQL.",
		},
	})
}

func parseRows(rows *sql.Rows) ([]string, [][]interface{}, error) {
	cols, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}

	var results [][]interface{}
	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}
		if err := rows.Scan(columnPointers...); err != nil {
			return nil, nil, err
		}
		results = append(results, columns)
	}
	return cols, results, nil
}
