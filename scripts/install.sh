#!/bin/sh
set -e

REPO="Kocoro-lab/ptengine-cli"
BINARY_NAME="ptengine-cli"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Detect OS
OS="$(uname -s)"
case "$OS" in
  Linux*)  OS="linux" ;;
  Darwin*) OS="darwin" ;;
  *)       echo "Error: unsupported OS: $OS"; exit 1 ;;
esac

# Detect architecture
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64|amd64)    ARCH="amd64" ;;
  aarch64|arm64)   ARCH="arm64" ;;
  *)               echo "Error: unsupported architecture: $ARCH"; exit 1 ;;
esac

# Get latest version tag
echo "Fetching latest release..."
VERSION=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"v([^"]+)".*/\1/')

if [ -z "$VERSION" ]; then
  echo "Error: failed to determine latest version"
  exit 1
fi

# Check if already installed and up to date
if command -v "$BINARY_NAME" >/dev/null 2>&1; then
  CURRENT=$("$BINARY_NAME" version 2>/dev/null | grep -o '"version":"[^"]*"' | cut -d'"' -f4 || echo "")
  if [ "$CURRENT" = "$VERSION" ]; then
    echo "${BINARY_NAME} v${VERSION} is already installed and up to date."
    exit 0
  fi
  if [ -n "$CURRENT" ]; then
    echo "Upgrading ${BINARY_NAME} from v${CURRENT} to v${VERSION}..."
  fi
fi

# Download
ARCHIVE="${BINARY_NAME}_${VERSION}_${OS}_${ARCH}.tar.gz"
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/v${VERSION}/${ARCHIVE}"

echo "Downloading ${BINARY_NAME} v${VERSION} for ${OS}/${ARCH}..."
TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

if ! curl -fsSL "$DOWNLOAD_URL" -o "${TMP_DIR}/${ARCHIVE}"; then
  echo "Error: failed to download ${DOWNLOAD_URL}"
  exit 1
fi
tar xzf "${TMP_DIR}/${ARCHIVE}" -C "$TMP_DIR"

# Install
if [ -w "$INSTALL_DIR" ]; then
  cp "${TMP_DIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
  chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
else
  echo "Installing to ${INSTALL_DIR} (requires sudo)..."
  sudo cp "${TMP_DIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
  sudo chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
fi

echo ""
echo "${BINARY_NAME} v${VERSION} installed to ${INSTALL_DIR}/${BINARY_NAME}"

# Only show getting-started hints on fresh install
if [ -z "$CURRENT" ]; then
  echo ""
  echo "Get started:"
  echo "  ${BINARY_NAME} config set --api-key pt-your-api-key"
  echo "  ${BINARY_NAME} heatmap describe"
  echo "  ${BINARY_NAME} heatmap query --help"
fi
