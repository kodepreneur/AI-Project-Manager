package system

import (
	"runtime"
	"time"

	"ai-project-manager/internal/database"
	"ai-project-manager/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
)

type SystemService struct{}

func NewSystemService() *SystemService {
	return &SystemService{}
}

func (s *SystemService) Metrics(c *fiber.Ctx) error {
	// CPU
	cpuPercentages, _ := cpu.Percent(200*time.Millisecond, false)
	cpuUsage := 0.0
	if len(cpuPercentages) > 0 {
		cpuUsage = cpuPercentages[0]
	}

	// Memory
	vm, _ := mem.VirtualMemory()
	memTotal := uint64(0)
	memUsed := uint64(0)
	memPercent := 0.0
	if vm != nil {
		memTotal = vm.Total / (1024 * 1024) // MB
		memUsed = vm.Used / (1024 * 1024)
		memPercent = vm.UsedPercent
	}

	// Disk
	d, _ := disk.Usage("/")
	diskTotal := uint64(0)
	diskUsed := uint64(0)
	diskPercent := 0.0
	if d != nil {
		diskTotal = d.Total / (1024 * 1024 * 1024) // GB
		diskUsed = d.Used / (1024 * 1024 * 1024)
		diskPercent = d.UsedPercent
	}

	// Load average
	l, _ := load.Avg()
	load1, load5, load15 := 0.0, 0.0, 0.0
	if l != nil {
		load1, load5, load15 = l.Load1, l.Load5, l.Load15
	}

	// Process count
	procs, _ := process.Pids()

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"cpu_percent":    cpuUsage,
			"num_cpu":        runtime.NumCPU(),
			"mem_total_mb":   memTotal,
			"mem_used_mb":    memUsed,
			"mem_percent":    memPercent,
			"disk_total_gb":  diskTotal,
			"disk_used_gb":   diskUsed,
			"disk_percent":   diskPercent,
			"load_1":         load1,
			"load_5":         load5,
			"load_15":        load15,
			"process_count":  len(procs),
			"os":             runtime.GOOS,
			"arch":           runtime.GOARCH,
			"server_time":    time.Now().Format("2006-01-02 15:04:05"),
		},
	})
}

func (s *SystemService) DashboardOverview(c *fiber.Ctx) error {
	var totalProjects, runningProjects, devProjects, completedProjects int64
	database.DB.Model(&models.Project{}).Count(&totalProjects)
	database.DB.Model(&models.Project{}).Where("status = ?", "Running").Count(&runningProjects)
	database.DB.Model(&models.Project{}).Where("status = ?", "Development").Count(&devProjects)
	database.DB.Model(&models.Project{}).Where("status = ?", "Completed").Count(&completedProjects)

	var totalTodos, completedTodos, inProgressTodos int64
	database.DB.Model(&models.Task{}).Count(&totalTodos)
	database.DB.Model(&models.Task{}).Where("status = ?", "done").Count(&completedTodos)
	database.DB.Model(&models.Task{}).Where("status = ?", "in_progress").Count(&inProgressTodos)

	var activeSessions int64
	database.DB.Model(&models.AntigravitySession{}).Where("status = ?", "running").Count(&activeSessions)

	// Recent logs/activity
	var recentLogs []models.ProjectLog
	database.DB.Order("created_at DESC").Limit(10).Find(&recentLogs)

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"projects": fiber.Map{
				"total":       totalProjects,
				"running":     runningProjects,
				"development": devProjects,
				"completed":   completedProjects,
			},
			"todos": fiber.Map{
				"total":       totalTodos,
				"completed":   completedTodos,
				"in_progress": inProgressTodos,
			},
			"antigravity": fiber.Map{
				"active_sessions": activeSessions,
			},
			"recent_activity": recentLogs,
		},
	})
}
