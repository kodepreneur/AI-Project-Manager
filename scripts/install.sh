#!/usr/bin/env bash
# AI Project Manager — 1-Click VPS Installer Script
set -e

echo "========================================================"
echo "    🚀 Installing AI Project Manager & Antigravity Panel "
echo "========================================================"

# Check root privileges
if [ "$EUID" -ne 0 ]; then
  echo "❌ Please run as root (use sudo): sudo bash $0 or curl ... | sudo bash"
  exit 1
fi

APP_DIR="/var/lib/ai-project-manager"
BIN_DIR="${APP_DIR}/bin"
STORAGE_DIR="${APP_DIR}/storage"
PROJECTS_DIR="${APP_DIR}/projects"
REPO_URL="https://github.com/kodepreneur/AI-Project-Manager.git"

# Detect source directory
SCRIPT_DIR=""
if [ -n "${BASH_SOURCE[0]}" ] && [ -f "${BASH_SOURCE[0]}" ]; then
  SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
fi

CLEANUP_SRC=false
if [ -d "./backend" ] && [ -d "./frontend" ]; then
  echo "📂 Using local repository in $(pwd)..."
  SRC_DIR="$(pwd)"
elif [ -n "$SCRIPT_DIR" ] && [ -d "$SCRIPT_DIR/../backend" ] && [ -d "$SCRIPT_DIR/../frontend" ]; then
  SRC_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
  echo "📂 Using repository in ${SRC_DIR}..."
else
  echo "📥 Standalone installer detected. Cloning repository..."
  CLEANUP_SRC=true
  SRC_DIR=$(mktemp -d /tmp/ai-project-manager-XXXXXX)

  if ! command -v git &> /dev/null; then
    echo "📦 Installing git..."
    apt-get update -y && apt-get install -y git
  fi

  git clone "${REPO_URL}" "${SRC_DIR}"
fi

# Ensure cleanup on exit if temp dir was created
if [ "$CLEANUP_SRC" = true ]; then
  trap 'rm -rf "${SRC_DIR}"' EXIT
fi

# Ensure system package manager and dependencies
if command -v apt-get &> /dev/null; then
  MISSING_PKGS=""
  if ! command -v curl &> /dev/null; then MISSING_PKGS="${MISSING_PKGS} curl"; fi
  if ! command -v gcc &> /dev/null; then MISSING_PKGS="${MISSING_PKGS} build-essential"; fi
  
  if [ -n "$MISSING_PKGS" ]; then
    echo "📦 Installing prerequisite packages:${MISSING_PKGS}..."
    apt-get update -y
    apt-get install -y ${MISSING_PKGS}
  fi
fi

# Ensure Go is installed
if ! command -v go &> /dev/null; then
  echo "📦 Installing Go compiler..."
  GO_VERSION="1.24.1"
  ARCH=$(uname -m)
  case "$ARCH" in
    x86_64) GO_ARCH="amd64" ;;
    aarch64|arm64) GO_ARCH="arm64" ;;
    *) GO_ARCH="amd64" ;;
  esac
  curl -fsSL "https://go.dev/dl/go${GO_VERSION}.linux-${GO_ARCH}.tar.gz" -o /tmp/go.tar.gz
  rm -rf /usr/local/go && tar -C /usr/local -xzf /tmp/go.tar.gz
  rm -f /tmp/go.tar.gz
  export PATH=$PATH:/usr/local/go/bin
  if ! grep -q '/usr/local/go/bin' /etc/profile; then
    echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
  fi
fi

# Ensure Node.js & npm are installed
if ! command -v node &> /dev/null || ! command -v npm &> /dev/null; then
  echo "📦 Installing Node.js LTS & npm..."
  if command -v apt-get &> /dev/null; then
    curl -fsSL https://deb.nodesource.com/setup_20.x | bash -
    apt-get install -y nodejs
  fi
fi

echo "📂 Creating application directories in ${APP_DIR}..."
mkdir -p "${BIN_DIR}"
mkdir -p "${STORAGE_DIR}/db"
mkdir -p "${STORAGE_DIR}/logs"
mkdir -p "${STORAGE_DIR}/screenshots"
mkdir -p "${PROJECTS_DIR}"

# Build Go backend
echo "🔨 Building Go backend binary..."
cd "${SRC_DIR}/backend"
CGO_ENABLED=1 go build -o "${BIN_DIR}/server" ./cmd/server
chmod +x "${BIN_DIR}/server"

# Build Frontend
echo "🌐 Building frontend UI bundle..."
cd "${SRC_DIR}/frontend"
npm install
npm run build
mkdir -p "${APP_DIR}/frontend"
rm -rf "${APP_DIR}/frontend/dist"
cp -r dist "${APP_DIR}/frontend/dist"

# Configure Systemd Service
echo "⚙️ Configuring systemd service..."
cat <<'EOF' > /etc/systemd/system/ai-project-manager.service
[Unit]
Description=AI Project Manager & Antigravity Control Panel
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/var/lib/ai-project-manager
ExecStart=/var/lib/ai-project-manager/bin/server
Restart=always
RestartSec=5
Environment=APP_NAME="AI Project Manager"
Environment=APP_ENV=production
Environment=APP_PORT=8080
Environment=DATABASE_URL=/var/lib/ai-project-manager/storage/db/app.db
Environment=PROJECT_ROOT=/var/lib/ai-project-manager/projects
Environment=STORAGE_ROOT=/var/lib/ai-project-manager/storage
Environment=ANTIGRAVITY_BINARY=antigravity
Environment=TERMINAL_ENABLED=true
Environment=SCREENSHOT_ENABLED=true

# Security Hardening
LimitNOFILE=65535
TimeoutStopSec=20

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable ai-project-manager
systemctl restart ai-project-manager
echo "✅ Service ai-project-manager started!"

# Detect public IP if possible
SERVER_IP=$(curl -s -4 ifconfig.me 2>/dev/null || hostname -I | awk '{print $1}' 2>/dev/null || echo "YOUR_SERVER_IP")

echo "========================================================"
echo "🎉 Installation Complete!"
echo "👉 Open http://${SERVER_IP}:8080 in your browser"
echo "👉 Follow the initial Administrator Setup wizard"
echo "👉 Service Status: sudo systemctl status ai-project-manager"
echo "========================================================"
