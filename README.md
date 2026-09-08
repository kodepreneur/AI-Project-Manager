# AI Project Manager

<p align="center">
  <img src="https://raw.githubusercontent.com/lucide-icons/lucide/main/icons/sparkles.svg" width="60" height="60" alt="Logo" />
</p>

A lightweight, modern web-based project management and development control panel designed to manage software projects through **Antigravity CLI**.

---

## 🌟 Key Capabilities

* **Project Lifecycle Management**: Multi-step project creation wizard with templates for Laravel, Vue, React, Go, and Node.js.
* **17-Section AI PRD Generator**: Generates comprehensive Product Requirements Documents before development starts, with revision tracking and markdown export.
* **PRD → TODO Breakdown Engine**: Analyzes approved PRDs to automatically generate phased development tasks with acceptance criteria.
* **Interactive Kanban Board**: Drag-and-drop task progression with backlog, todo, in-progress, review, and done states.
* **Antigravity CLI Session Manager**: Dispatches tasks directly to Antigravity CLI with full context injection (PRD, task criteria, project directory) and live WebSocket output streaming.
* **Interactive Web Terminal**: Full-featured in-browser terminal using **xterm.js** and Go **creack/pty**.
* **Live Preview & Screenshots**: Live application iframe embedding with automated preview screenshot generation.
* **Infrastructure Telemetry**: Host CPU, RAM, Disk, and load average metrics without heavy external agents.
* **Runtime & Database Tools**: PHP version diagnostics, MySQL/PostgreSQL status monitoring, and SQL console.
* **Developer First UX**: Sleek dark mode design inspired by Vercel and Linear, with a global Command Palette (`⌘K`).

---

## 🏗️ Architecture

```text
ai-project-manager/
├── backend/
│   ├── cmd/server/main.go        # Go Fiber server & router
│   ├── internal/
│   │   ├── antigravity/          # Antigravity CLI process manager
│   │   ├── auth/                 # JWT & first-time setup wizard
│   │   ├── database/             # GORM SQLite & migrations
│   │   ├── env/                  # Masked .env editor
│   │   ├── git/                  # Git status, diff viewer, commits
│   │   ├── logs/                 # Centralized structured logger
│   │   ├── mysql/ & postgres/    # Database status & SQL runner
│   │   ├── php/                  # PHP runtime detection
│   │   ├── prd/                  # 17-section AI PRD generator
│   │   ├── preview/              # Dev server preview & screenshots
│   │   ├── processes/            # Child process manager
│   │   ├── projects/             # Projects CRUD & templates
│   │   ├── system/               # Telemetry metrics
│   │   ├── terminal/             # PTY WebSocket terminal
│   │   ├── todos/                # Kanban tasks & PRD breakdown
│   │   └── websocket/            # Global event broadcaster
├── frontend/
│   ├── src/
│   │   ├── components/           # Kanban, Terminal, Preview, Git, PRD, CommandPalette
│   │   ├── layouts/              # AppLayout & AuthLayout
│   │   ├── pages/                # Dashboard, Projects, Detail, Kanban, Infra, Logs
│   │   ├── stores/               # Pinia state stores
│   │   └── router/               # Vue Router with setup guards
├── deploy/                       # Systemd, Nginx, and Caddy configurations
├── scripts/                      # 1-Click install script
└── storage/                      # SQLite DB, logs, and screenshots
```

---

## 🚀 Quick Start (Local Development)

### 1. Start Go Backend
```bash
cd backend
go run ./cmd/server
```
The server will start on `http://localhost:8080`.

### 2. Start Frontend Dev Server
```bash
cd frontend
npm run dev
```
Open `http://localhost:5173` in your browser.

---

## 📦 Single Binary / Production Build

You can serve the entire application as a single self-contained unit:

```bash
# 1. Build frontend
cd frontend
npm run build
cd ..

# 2. Build backend
cd backend
go build -o ../server ./cmd/server
cd ..

# 3. Run production server
./server
```
Access the application directly on `http://localhost:8080`.

---

## 🖥️ Production VPS Deployment

### 1-Click Installer
```bash
curl -fsSL https://raw.githubusercontent.com/kodepreneur/AI-Project-Manager/main/scripts/install.sh | sudo bash
```

### Systemd Service Management
```bash
sudo systemctl start ai-project-manager
sudo systemctl status ai-project-manager
sudo systemctl restart ai-project-manager
```

---

## 🔐 Security & Operations

* **First-Time Setup**: The panel requires creating a master administrator account on initial startup.
* **Authentication**: All API requests and WebSocket upgrades require valid JWT tokens.
* **Restricted Shell**: Web terminal sessions default to the project's folder with user environment variables.
* **Secret Masking**: Sensitive variables in `.env` (passwords, tokens, keys) are masked automatically in the UI.

---

## 📄 License
MIT License.
