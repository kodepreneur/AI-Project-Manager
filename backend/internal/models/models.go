package models

import (
	"time"
)

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"uniqueIndex;size:100;not null" json:"username"`
	Email        string    `gorm:"uniqueIndex;size:150;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"`
	Role         string    `gorm:"size:50;default:'admin'" json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Project struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	Name            string         `gorm:"size:150;not null" json:"name"`
	Slug            string         `gorm:"uniqueIndex;size:150;not null" json:"slug"`
	Description     string         `gorm:"type:text" json:"description"`
	Path            string         `gorm:"size:255;not null" json:"path"`
	GitRepo         string         `gorm:"size:255" json:"git_repo"`
	GitBranch       string         `gorm:"size:100;default:'main'" json:"git_branch"`
	Framework       string         `gorm:"size:100;default:'Laravel'" json:"framework"`
	Language        string         `gorm:"size:100;default:'PHP'" json:"language"`
	FrontendTech    string         `gorm:"size:100;default:'Vue'" json:"frontend_tech"`
	DatabaseType    string         `gorm:"size:50;default:'MySQL'" json:"database_type"`
	PHPVersion      string         `gorm:"size:50;default:'8.3'" json:"php_version"`
	Status          string         `gorm:"size:50;default:'Development'" json:"status"` // Development, Running, Completed, Paused, Failed
	PreviewPort     int            `gorm:"default:8001" json:"preview_port"`
	PreviewCommand  string         `gorm:"size:255;default:'npm run dev'" json:"preview_command"`
	TestCommand     string         `gorm:"size:255;default:'npm test'" json:"test_command"`
	BuildCommand    string         `gorm:"size:255;default:'npm run build'" json:"build_command"`
	InstallCommand  string         `gorm:"size:255;default:'npm install'" json:"install_command"`
	ScreenshotURL   string         `gorm:"size:255" json:"screenshot_url"`
	LastActivity    time.Time      `json:"last_activity"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	PRD             *PRD           `gorm:"foreignKey:ProjectID" json:"prd,omitempty"`
	Tasks           []Task         `gorm:"foreignKey:ProjectID" json:"tasks,omitempty"`
	Processes       []ProjectProcess `gorm:"foreignKey:ProjectID" json:"processes,omitempty"`
}

type PRD struct {
	ID          uint         `gorm:"primaryKey" json:"id"`
	ProjectID   uint         `gorm:"uniqueIndex;not null" json:"project_id"`
	Title       string       `gorm:"size:200;not null" json:"title"`
	Content     string       `gorm:"type:text;not null" json:"content"`
	Status      string       `gorm:"size:50;default:'Draft'" json:"status"` // Draft, Review, Approved, Archived
	Version     int          `gorm:"default:1" json:"version"`
	ApprovedAt  *time.Time   `json:"approved_at"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	Versions    []PRDVersion `gorm:"foreignKey:PRDID" json:"versions,omitempty"`
}

type PRDVersion struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	PRDID     uint      `gorm:"index;not null" json:"prd_id"`
	Version   int       `gorm:"not null" json:"version"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	Summary   string    `gorm:"size:255" json:"summary"`
	CreatedAt time.Time `json:"created_at"`
}

type Phase struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProjectID uint      `gorm:"index;not null" json:"project_id"`
	Name      string    `gorm:"size:150;not null" json:"name"`
	Order     int       `gorm:"column:phase_order;default:0" json:"order"`
	CreatedAt time.Time `json:"created_at"`
}

type Task struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	ProjectID          uint      `gorm:"index;not null" json:"project_id"`
	PhaseID            *uint     `gorm:"index" json:"phase_id"`
	PhaseName          string    `gorm:"size:150" json:"phase_name"`
	Title              string    `gorm:"size:200;not null" json:"title"`
	Description        string    `gorm:"type:text" json:"description"`
	Priority           string    `gorm:"size:50;default:'medium'" json:"priority"` // low, medium, high, urgent
	Complexity         string    `gorm:"size:50;default:'medium'" json:"complexity"` // low, medium, high
	Status             string    `gorm:"size:50;default:'todo'" json:"status"` // backlog, todo, in_progress, review, done, cancelled
	AcceptanceCriteria string    `gorm:"type:text" json:"acceptance_criteria"`
	AIInstructions     string    `gorm:"type:text" json:"ai_instructions"`
	Dependencies       string    `gorm:"type:text" json:"dependencies"` // JSON array or comma separated
	FilesAffected      string    `gorm:"type:text" json:"files_affected"`
	Order              int       `gorm:"column:task_order;default:0" json:"order"`
	AIStatus           string    `gorm:"size:50;default:'idle'" json:"ai_status"` // idle, running, completed, failed
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type AntigravitySession struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	ProjectID   uint       `gorm:"index;not null" json:"project_id"`
	TaskID      *uint      `gorm:"index" json:"task_id"`
	PID         int        `json:"pid"`
	Status      string     `gorm:"size:50;default:'idle'" json:"status"` // idle, running, paused, completed, failed
	Command     string     `gorm:"type:text" json:"command"`
	Stdout      string     `gorm:"type:text" json:"stdout"`
	Stderr      string     `gorm:"type:text" json:"stderr"`
	ExitCode    int        `json:"exit_code"`
	StartTime   time.Time  `json:"start_time"`
	EndTime     *time.Time `json:"end_time"`
	CurrentTask string     `gorm:"size:255" json:"current_task"`
}

type ProjectProcess struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ProjectID   uint      `gorm:"index;not null" json:"project_id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Command     string    `gorm:"size:255;not null" json:"command"`
	PID         int       `json:"pid"`
	Port        int       `json:"port"`
	Status      string    `gorm:"size:50;default:'stopped'" json:"status"` // stopped, running, failed
	AutoRestart bool      `gorm:"default:false" json:"auto_restart"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type DatabaseConnection struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProjectID *uint     `gorm:"index" json:"project_id"`
	Driver    string    `gorm:"size:50;not null" json:"driver"` // mysql, postgres, sqlite
	Name      string    `gorm:"size:100;not null" json:"name"`
	Host      string    `gorm:"size:100;default:'127.0.0.1'" json:"host"`
	Port      int       `gorm:"default:3306" json:"port"`
	Database  string    `gorm:"size:100;not null" json:"database"`
	Username  string    `gorm:"size:100" json:"username"`
	Password  string    `gorm:"size:255" json:"-"` // Never expose in JSON
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProjectLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProjectID *uint     `gorm:"index" json:"project_id"`
	Category  string    `gorm:"size:50;index;not null" json:"category"` // Application, Antigravity, Terminal, PHP, Database, Project Process, System
	Level     string    `gorm:"size:20;default:'INFO'" json:"level"` // INFO, WARN, ERROR, DEBUG
	Message   string    `gorm:"type:text;not null" json:"message"`
	Metadata  string    `gorm:"type:text" json:"metadata"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`
}

type SystemSetting struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Key       string    `gorm:"uniqueIndex;size:100;not null" json:"key"`
	Value     string    `gorm:"type:text" json:"value"`
	Category  string    `gorm:"size:50;default:'general'" json:"category"`
	UpdatedAt time.Time `json:"updated_at"`
}
