#!/usr/bin/env bash
set -euo pipefail

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m'

ok()   { echo -e "${GREEN}[OK]${NC} $*"; }
err()  { echo -e "${RED}[HATA]${NC} $*"; }
info() { echo -e "${CYAN}[*]${NC} $*"; }
warn() { echo -e "${YELLOW}[!]${NC} $*"; }

FAILED=0
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

echo ""
echo -e "${BOLD}========================================${NC}"
echo -e "${BOLD}       b0yzover Kurulum (Linux)        ${NC}"
echo -e "${BOLD}========================================${NC}"
echo ""

ensure_go() {
  if command -v go >/dev/null 2>&1; then
    ok "Go kurulu: $(go version)"
    return 0
  fi

  info "Go bulunamadı, paket yöneticisi ile kurulum deneniyor..."

  if command -v apt-get >/dev/null 2>&1; then
    if [ "$(id -u)" -eq 0 ]; then
      apt-get update -y && (apt-get install -y golang-go || apt-get install -y golang)
    else
      sudo apt-get update -y && (sudo apt-get install -y golang-go || sudo apt-get install -y golang)
    fi
  elif command -v dnf >/dev/null 2>&1; then
    if [ "$(id -u)" -eq 0 ]; then dnf install -y golang; else sudo dnf install -y golang; fi
  elif command -v yum >/dev/null 2>&1; then
    if [ "$(id -u)" -eq 0 ]; then yum install -y golang; else sudo yum install -y golang; fi
  elif command -v pacman >/dev/null 2>&1; then
    if [ "$(id -u)" -eq 0 ]; then pacman -Sy --noconfirm go; else sudo pacman -Sy --noconfirm go; fi
  elif command -v zypper >/dev/null 2>&1; then
    if [ "$(id -u)" -eq 0 ]; then zypper install -y go; else sudo zypper install -y go; fi
  else
    warn "Bilinen paket yöneticisi bulunamadı"
  fi

  hash -r 2>/dev/null || true
  export PATH="/usr/local/go/bin:${HOME}/go/bin:${PATH}"

  if command -v go >/dev/null 2>&1; then
    ok "Go kuruldu: $(go version)"
    return 0
  fi

  if command -v snap >/dev/null 2>&1; then
    info "snap ile Go kuruluyor..."
    if [ "$(id -u)" -eq 0 ]; then
      snap install go --classic || true
    else
      sudo snap install go --classic || true
    fi
    hash -r 2>/dev/null || true
    export PATH="/snap/bin:${PATH}"
  fi

  if command -v go >/dev/null 2>&1; then
    ok "Go kuruldu (snap): $(go version)"
    return 0
  fi

  err "Go otomatik kurulamadı"
  echo "  Manuel kurulum: https://go.dev/dl/  veya  https://golang.org/dl/"
  return 1
}

if ! ensure_go; then
  FAILED=1
fi

if [ "$FAILED" -eq 0 ]; then
  info "Proje dizini: $PROJECT_DIR"
  cd "$PROJECT_DIR"

  info "go mod download..."
  if go mod download; then
    ok "Bağımlılıklar indirildi"
  else
    err "go mod download başarısız"
    FAILED=1
  fi

  if [ "$FAILED" -eq 0 ]; then
    info "go build -o b0yzover ."
    if go build -o b0yzover .; then
      ok "Derleme tamam: ${PROJECT_DIR}/b0yzover"
    else
      err "go build başarısız"
      FAILED=1
    fi
  fi
fi

echo ""
echo -e "${BOLD}----------------------------------------${NC}"
if [ "$FAILED" -eq 0 ]; then
  ok "Kurulum başarıyla tamamlandı"
  exit 0
else
  err "Kurulum başarısız oldu"
  exit 1
fi
