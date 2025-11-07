#!/usr/bin/env bash

# SecureFlow Launcher Script
# This script automatically selects the correct platform-specific executable
# and runs it with the provided arguments.

set -e

# Detect OS and architecture
detect_platform() {
    OS="$(uname -s)"
    ARCH="$(uname -m)"
    
    case "$OS" in
        Linux*)
            PLATFORM="linux"
            ;;
        Darwin*)
            PLATFORM="darwin"
            ;;
        MINGW*|MSYS*|CYGWIN*)
            PLATFORM="windows"
            ;;
        *)
            echo "Error: Unsupported operating system: $OS" >&2
            exit 1
            ;;
    esac
    
    case "$ARCH" in
        x86_64|amd64)
            ARCHITECTURE="amd64"
            ;;
        aarch64|arm64)
            ARCHITECTURE="arm64"
            ;;
        *)
            echo "Error: Unsupported architecture: $ARCH" >&2
            exit 1
            ;;
    esac
}

# Main
detect_platform

# Construct binary path
BINARY_NAME="secureflow-${PLATFORM}-${ARCHITECTURE}"
if [ "$PLATFORM" = "windows" ]; then
    BINARY_NAME="${BINARY_NAME}.exe"
fi

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BINARY_PATH="${SCRIPT_DIR}/.secureflow/${BINARY_NAME}"

# Check if binary exists
if [ ! -f "$BINARY_PATH" ]; then
    echo "Error: Binary not found at $BINARY_PATH" >&2
    echo "Platform: ${PLATFORM}-${ARCHITECTURE}" >&2
    echo "" >&2
    echo "Available binaries:" >&2
    ls -1 "${SCRIPT_DIR}/.secureflow/" 2>/dev/null || echo "  None found" >&2
    exit 1
fi

# Make sure binary is executable
chmod +x "$BINARY_PATH" 2>/dev/null || true

# Run the binary with all arguments passed to this script
exec "$BINARY_PATH" "$@"
