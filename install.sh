#!/bin/sh
set -e

# ponysay-go installer script
# Supports macOS and Linux (amd64, arm64)

REPO="bradly0cjw/ponysay-go"

# 1. Detect OS
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "$OS" in
    darwin) OS="darwin" ;;
    linux) OS="linux" ;;
    *)
        echo "Error: Unsupported operating system: $OS" >&2
        exit 1
        ;;
esac

# 2. Detect Architecture
ARCH="$(uname -m)"
case "$ARCH" in
    x86_64|amd64) ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *)
        echo "Error: Unsupported architecture: $ARCH" >&2
        exit 1
        ;;
esac

BINARY_NAME="ponysay-${OS}-${ARCH}"
DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/${BINARY_NAME}"

echo "Detected OS: ${OS}, Arch: ${ARCH}"
echo "Downloading ${BINARY_NAME}..."

TMP_DIR="$(mktemp -d 2>/dev/null || mktemp -d -t 'ponysay')"
TMP_BINARY="${TMP_DIR}/ponysay"
cleanup() {
    rm -rf "$TMP_DIR"
}
trap cleanup EXIT INT TERM

# 3. Download using curl or wget
if command -v curl >/dev/null 2>&1; then
    curl -fsSL "$DOWNLOAD_URL" -o "$TMP_BINARY"
elif command -v wget >/dev/null 2>&1; then
    wget -qO "$TMP_BINARY" "$DOWNLOAD_URL"
else
    echo "Error: Neither curl nor wget was found. Please install curl or wget." >&2
    exit 1
fi

chmod +x "$TMP_BINARY"

# 4. Determine installation path
# Global install if running as root or if /usr/local/bin is writable, otherwise user local (~/.local/bin)
if [ "$(id -u)" -eq 0 ] || [ -w "/usr/local/bin" ]; then
    INSTALL_DIR="/usr/local/bin"
    IS_GLOBAL=1
else
    INSTALL_DIR="${HOME}/.local/bin"
    IS_GLOBAL=0
fi

mkdir -p "$INSTALL_DIR"

echo "Installing ponysay to ${INSTALL_DIR}..."
mv "$TMP_BINARY" "${INSTALL_DIR}/ponysay"
ln -sf "${INSTALL_DIR}/ponysay" "${INSTALL_DIR}/ponythink"

echo "Successfully installed ponysay & ponythink to ${INSTALL_DIR}!"

# 5. Check PATH for local install
if [ "$IS_GLOBAL" -eq 0 ]; then
    case ":$PATH:" in
        *":${INSTALL_DIR}:"*) ;;
        *)
            echo ""
            echo "--------------------------------------------------------"
            echo "WARNING: ${INSTALL_DIR} is not in your current PATH!"
            echo "To run 'ponysay' and 'ponythink' from anywhere, add this line"
            echo "to your ~/.bashrc, ~/.zshrc, or profile:"
            echo ""
            echo "    export PATH=\"\$HOME/.local/bin:\$PATH\""
            echo "--------------------------------------------------------"
            ;;
    esac
fi

echo ""
echo "Try running:"
echo "    ponysay 'I am just the cutest pony!'"
