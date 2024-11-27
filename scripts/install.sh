#!/bin/bash

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Map architecture names
case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64)
        ARCH="arm64"
        ;;
    armv7l)
        ARCH="arm"
        ;;
    i386|i686)
        ARCH="386"
        ;;
esac

# Set installation directory
INSTALL_DIR="/usr/local/bin"
if [ ! -w "$INSTALL_DIR" ]; then
    echo "Error: Cannot write to $INSTALL_DIR"
    echo "Please run with sudo: sudo ./install.sh"
    exit 1
fi

# Find the appropriate binary
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
RELEASES_DIR="$SCRIPT_DIR/releases"
BINARY_PATH="$RELEASES_DIR/${OS}-${ARCH}/xdir"
if [ $OS = "windows" ]; then
    BINARY_PATH="${BINARY_PATH}.exe"
fi

echo "Looking for binary at: $BINARY_PATH"

# Check if binary exists
if [ ! -f "$BINARY_PATH" ]; then
    echo "Error: Binary not found for ${OS}/${ARCH}"
    echo "Please run generate-release.sh first or build from source"
    exit 1
fi

# Install binary
echo "Installing xdir for ${OS}/${ARCH}..."
cp "$BINARY_PATH" "$INSTALL_DIR/xdir"
chmod +x "$INSTALL_DIR/xdir"

echo "Installation complete! You can now use 'xdir' from anywhere."
echo "Try 'xdir help' for usage information."
