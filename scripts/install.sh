#!/usr/bin/env bash
# AI Project Manager — 1-Click VPS Installer Script
set -e

echo "========================================================"
echo "    🚀 Installing AI Project Manager & Antigravity Panel "
echo "========================================================"

# Check root privileges
if [ "$EUID" -ne 0 ]; then
  echo "❌ Please run as root (use sudo)"
  exit 1
fi

APP_DIR="/var/lib/ai-project-manager"
BIN_DIR="${APP_DIR}/bin"
STORAGE_DIR="${APP_DIR}/storage"
PROJECTS_DIR="${APP_DIR}/projects"

echo "📂 Creating application directories in ${APP_DIR}..."
mkdir -p "${BIN_DIR}"
mkdir -p "${STORAGE_DIR}/db"
mkdir -p "${STORAGE_DIR}/logs"
mkdir -p "${STORAGE_DIR}/screenshots"
mkdir -p "${PROJECTS_DIR}"

# Build or copy binaries
if [ -f "./backend/server" ]; then
  echo "📦 Copying compiled server binary..."
  cp ./backend/server "${BIN_DIR}/server"
  chmod +x "${BIN_DIR}/server"
else
  echo "🔨 Building Go backend binary..."
  cd backend
  go build -o "${BIN_DIR}/server" ./cmd/server
  cd ..
fi

# Build or copy frontend dist
if [ -d "./frontend/dist" ]; then
  echo "🌐 Copying compiled frontend bundle..."
  mkdir -p "${APP_DIR}/frontend"
  cp -r ./frontend/dist "${APP_DIR}/frontend/dist"
fi

# Copy systemd service
echo "⚙️ Configuring systemd service..."
if [ -f "./deploy/ai-project-manager.service" ]; then
  cp ./deploy/ai-project-manager.service /etc/systemd/system/ai-project-manager.service
  systemctl daemon-reload
  systemctl enable ai-project-manager
  systemctl restart ai-project-manager
  echo "✅ Service ai-project-manager started!"
fi

echo "========================================================"
echo "🎉 Installation Complete!"
echo "👉 Open http://YOUR_SERVER_IP:8080 in your browser"
echo "👉 Follow the initial Administrator Setup wizard"
echo "========================================================"
