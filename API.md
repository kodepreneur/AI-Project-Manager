# AI Project Manager — API Documentation

The backend exposes RESTful HTTP endpoints and WebSocket streams.

## Base URL
`http://localhost:8080/api`

---

## 1. Authentication

### Check Setup Status
* **`GET /api/auth/setup-status`**
* Returns whether an administrator account has been provisioned.

### Initial Administrator Setup
* **`POST /api/auth/setup`**
* Body:
  ```json
  {
    "username": "admin",
    "email": "admin@example.com",
    "password": "SecurePassword123"
  }
  ```

### Login
* **`POST /api/auth/login`**
* Body:
  ```json
  {
    "username_or_email": "admin",
    "password": "SecurePassword123"
  }
  ```
* Response: `{ "success": true, "data": { "token": "JWT_TOKEN", "user": { ... } } }`

### Current Profile
* **`GET /api/auth/me`** (Bearer Token required)

---

## 2. Projects

### List Projects
* **`GET /api/projects?status=Running&search=crm`**

### Create Project
* **`POST /api/projects`**
* Body:
  ```json
  {
    "name": "Toko Buah ERP",
    "slug": "toko-buah-erp",
    "description": "Retail ERP",
    "framework": "Laravel",
    "language": "PHP",
    "frontend_tech": "Vue 3",
    "database_type": "MySQL",
    "php_version": "8.3"
  }
  ```

### Project Health Check
* **`GET /api/projects/:id/health`**
* Validates directory, git, PHP, Node, npm, Composer, and .env.

---

## 3. Product Requirements Documents (PRD)

### Get PRD
* **`GET /api/projects/:id/prd`**

### Save / Update PRD
* **`POST /api/projects/:id/prd`**

### Generate 17-Section PRD with AI
* **`POST /api/projects/:id/prd/generate`**
* Body:
  ```json
  {
    "project_idea": "...",
    "target_users": "...",
    "business_goals": "...",
    "main_features": "..."
  }
  ```

### Approve PRD
* **`POST /api/projects/:id/prd/approve`**

---

## 4. Development Tasks & Kanban

### List Tasks
* **`GET /api/projects/:id/tasks`**

### Generate Plan from PRD
* **`POST /api/projects/:id/tasks/generate-from-prd`**
* Automatically parses approved PRD and builds phased tasks.

### Update Task Status
* **`PUT /api/tasks/:taskId/status`**
* Body: `{ "status": "in_progress" }` (Options: `backlog`, `todo`, `in_progress`, `review`, `done`, `cancelled`)

---

## 5. Antigravity AI Engine

### Get Session Status
* **`GET /api/projects/:id/antigravity/status`**

### Start AI Task
* **`POST /api/projects/:id/antigravity/start`**
* Body:
  ```json
  {
    "task_id": 12,
    "prompt": "Custom instruction"
  }
  ```

### Stop AI Task
* **`POST /api/projects/:id/antigravity/stop`**

---

## 6. Process Manager & Preview

### List Running Processes
* **`GET /api/projects/:id/processes`**

### Start Process
* **`POST /api/projects/:id/processes/start`**
* Body: `{ "name": "Vite", "command": "npm run dev", "port": 5173 }`

### Stop Process
* **`POST /api/processes/:procId/stop`**

### Capture Screenshot
* **`POST /api/projects/:id/preview/screenshot`**

---

## 7. Runtimes & Databases

* `GET /api/php/diagnostics`
* `GET /api/mysql/status`
* `POST /api/mysql/test`
* `POST /api/mysql/query`
* `GET /api/postgres/status`
* `POST /api/postgres/test`
* `POST /api/postgres/query`

---

## 8. WebSockets

* **`GET /ws/events`**: Real-time event notifications, log broadcasting, and Antigravity execution line updates.
* **`GET /ws/projects/:id/terminal`**: Interactive PTY pseudo-terminal stream powered by `creack/pty`.
