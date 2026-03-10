#!/bin/bash
set -euo pipefail

# gtasks installer — downloads the latest release from GitHub
# Usage: curl -fsSL https://gtasks.sidv.dev/install | bash

REPO="BRO3886/gtasks"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
BINARY_NAME="gtasks"

info()  { printf "\033[36m%s\033[0m\n" "$*"; }
warn()  { printf "\033[33mWarning: %s\033[0m\n" "$*" >&2; }
error() { printf "\033[31mError: %s\033[0m\n" "$*" >&2; exit 1; }

# --- Detect OS and architecture ---

OS="$(uname -s)"
ARCH="$(uname -m)"

case "$OS" in
    Darwin) OS_KEY="mac" ;;
    Linux)  OS_KEY="linux" ;;
    *) error "Unsupported OS: $OS. For Windows, download the binary from https://github.com/${REPO}/releases" ;;
esac

case "$ARCH" in
    arm64|aarch64) ARCH_KEY="arm64" ;;
    x86_64)        ARCH_KEY="amd64" ;;
    *) error "Unsupported architecture: $ARCH" ;;
esac

if ! command -v curl >/dev/null 2>&1; then
    error "curl is required but not found"
fi

# --- Resolve latest version ---

info "Fetching latest release..."
LATEST=$(curl -sSL -H "Accept: application/vnd.github+json" \
    "https://api.github.com/repos/${REPO}/releases/latest" \
    | grep '"tag_name"' | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/')

if [ -z "$LATEST" ]; then
    error "Could not determine latest release"
fi

if ! echo "$LATEST" | grep -qE '^v[0-9]+\.[0-9]+\.[0-9]+'; then
    error "Unexpected version format: '$LATEST'"
fi

info "Latest version: $LATEST"

# --- Build asset name ---
# Archive naming convention:
#   gtasks_mac_arm64_<version>.tar.gz
#   gtasks_mac_amd64_<version>.tar.gz
#   gtasks_linux_arm64_<version>.tar.gz
#   gtasks_linux_amd64_<version>.tar.gz

ASSET_NAME="gtasks_${OS_KEY}_${ARCH_KEY}_${LATEST}.tar.gz"

DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${LATEST}/${ASSET_NAME}"

# --- Download and extract ---

TMPDIR_PATH=$(mktemp -d)
trap 'rm -rf "$TMPDIR_PATH"' EXIT

info "Downloading ${ASSET_NAME}..."
HTTP_CODE=$(curl -sSL -w "%{http_code}" -o "${TMPDIR_PATH}/${ASSET_NAME}" "$DOWNLOAD_URL")

if [ "$HTTP_CODE" != "200" ]; then
    error "Download failed (HTTP $HTTP_CODE). Asset '${ASSET_NAME}' may not exist for ${LATEST}."
fi

tar -xzf "${TMPDIR_PATH}/${ASSET_NAME}" -C "${TMPDIR_PATH}"

# --- Install ---

mkdir -p "$INSTALL_DIR"
if [ -w "$INSTALL_DIR" ]; then
    mv "${TMPDIR_PATH}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
else
    info "Requires sudo to install to ${INSTALL_DIR}"
    sudo mv "${TMPDIR_PATH}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
fi
chmod +x "${INSTALL_DIR}/${BINARY_NAME}"

info "Installed gtasks ${LATEST} to ${INSTALL_DIR}/${BINARY_NAME}"

# --- Verify ---

if command -v gtasks >/dev/null 2>&1; then
    info "Run 'gtasks --help' to get started"
else
    info "Add this to your shell profile: export PATH=\"\$PATH:${INSTALL_DIR}\""
fi

# --- Agent skill installation ---

printf '\n'
info "gtasks can install an AI agent skill that teaches Claude Code how to use it."
printf "Install agent skill now? [y/N] "
answer="n"
{ read -r answer < /dev/tty; } 2>/dev/null || true
if [ "$answer" = "y" ] || [ "$answer" = "Y" ]; then
    "${INSTALL_DIR}/${BINARY_NAME}" skills install || warn "Skill install failed — run 'gtasks skills install' manually"
fi
