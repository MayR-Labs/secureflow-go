#!/usr/bin/env bash

# SecureFlow Installation Script
# This script installs the latest version of SecureFlow CLI
# Usage: 
#   Local install:  curl -sSL https://raw.githubusercontent.com/MayR-Labs/secureflow-go/main/install.sh | bash
#   Global install: curl -sSL https://raw.githubusercontent.com/MayR-Labs/secureflow-go/main/install.sh | bash -s -- --global

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Parse command line arguments
GLOBAL_INSTALL=false
for arg in "$@"; do
    case $arg in
        --global)
            GLOBAL_INSTALL=true
            shift
            ;;
    esac
done

# Default installation directory
if [ "$GLOBAL_INSTALL" = true ]; then
    INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
else
    INSTALL_DIR="$(pwd)"
fi

# GitHub repository
REPO="MayR-Labs/secureflow-go"
BINARY_NAME="secureflow"

echo -e "${BLUE}SecureFlow CLI Installer${NC}"
echo "========================="
if [ "$GLOBAL_INSTALL" = true ]; then
    echo -e "${BLUE}Mode:${NC} Global installation"
else
    echo -e "${BLUE}Mode:${NC} Local installation (current directory)"
fi
echo ""

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
            echo -e "${RED}Error: Unsupported operating system: $OS${NC}"
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
            echo -e "${RED}Error: Unsupported architecture: $ARCH${NC}"
            exit 1
            ;;
    esac
    
    echo -e "${GREEN}✓${NC} Detected platform: ${PLATFORM}-${ARCHITECTURE}"
}

# Get the latest release version
get_latest_version() {
    echo -e "${BLUE}→${NC} Fetching latest release..."
    
    if command -v curl >/dev/null 2>&1; then
        LATEST_VERSION=$(curl -sSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    elif command -v wget >/dev/null 2>&1; then
        LATEST_VERSION=$(wget -qO- "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    else
        echo -e "${RED}Error: Neither curl nor wget is available. Please install one of them.${NC}"
        exit 1
    fi
    
    if [ -z "$LATEST_VERSION" ]; then
        echo -e "${RED}Error: Failed to fetch latest version${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}✓${NC} Latest version: ${LATEST_VERSION}"
}

# Download and install binary
install_binary() {
    # Construct download URL
    if [ "$PLATFORM" = "windows" ]; then
        BINARY_FILE="${BINARY_NAME}-${PLATFORM}-${ARCHITECTURE}.exe"
    else
        BINARY_FILE="${BINARY_NAME}-${PLATFORM}-${ARCHITECTURE}"
    fi
    
    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${LATEST_VERSION}/${BINARY_FILE}"
    
    echo -e "${BLUE}→${NC} Downloading from: ${DOWNLOAD_URL}"
    
    # Create temp directory
    TMP_DIR=$(mktemp -d)
    TMP_FILE="${TMP_DIR}/${BINARY_FILE}"
    
    # Download binary
    if command -v curl >/dev/null 2>&1; then
        curl -sSL -o "${TMP_FILE}" "${DOWNLOAD_URL}"
    elif command -v wget >/dev/null 2>&1; then
        wget -q -O "${TMP_FILE}" "${DOWNLOAD_URL}"
    fi
    
    if [ ! -f "${TMP_FILE}" ]; then
        echo -e "${RED}Error: Failed to download binary${NC}"
        rm -rf "${TMP_DIR}"
        exit 1
    fi
    
    echo -e "${GREEN}✓${NC} Downloaded successfully"
    
    # Make binary executable
    chmod +x "${TMP_FILE}"
    
    # Determine final installation path
    if [ "$PLATFORM" = "windows" ]; then
        INSTALL_PATH="${INSTALL_DIR}/${BINARY_NAME}.exe"
    else
        INSTALL_PATH="${INSTALL_DIR}/${BINARY_NAME}"
    fi
    
    # Check if we need sudo (only for global install)
    if [ "$GLOBAL_INSTALL" = true ] && [ ! -w "$INSTALL_DIR" ]; then
        echo -e "${YELLOW}→${NC} Installing to ${INSTALL_PATH} (requires sudo)..."
        sudo mv "${TMP_FILE}" "${INSTALL_PATH}"
    else
        echo -e "${BLUE}→${NC} Installing to ${INSTALL_PATH}..."
        mv "${TMP_FILE}" "${INSTALL_PATH}"
    fi
    
    # Clean up
    rm -rf "${TMP_DIR}"
    
    echo -e "${GREEN}✓${NC} Installed to ${INSTALL_PATH}"
}

# Verify installation
verify_installation() {
    echo ""
    echo -e "${BLUE}→${NC} Verifying installation..."
    
    if [ "$GLOBAL_INSTALL" = true ]; then
        if command -v "${BINARY_NAME}" >/dev/null 2>&1; then
            VERSION_OUTPUT=$("${BINARY_NAME}" --version 2>&1 || echo "")
            echo -e "${GREEN}✓${NC} SecureFlow is installed globally and ready to use!"
            echo ""
            echo "Version: ${VERSION_OUTPUT}"
            echo ""
            echo "Usage:"
            echo "  secureflow init       # Initialize configuration"
            echo "  secureflow encrypt    # Encrypt files"
            echo "  secureflow decrypt    # Decrypt files"
            echo "  secureflow --help     # Show help"
            echo ""
            echo -e "For more information, visit: ${BLUE}https://github.com/${REPO}${NC}"
        else
            echo -e "${YELLOW}⚠${NC}  Installation complete, but ${BINARY_NAME} is not in your PATH"
            echo "   You may need to add ${INSTALL_DIR} to your PATH or restart your shell"
        fi
    else
        # Local installation
        if [ -f "${INSTALL_PATH}" ]; then
            VERSION_OUTPUT=$("${INSTALL_PATH}" --version 2>&1 || echo "")
            echo -e "${GREEN}✓${NC} SecureFlow is installed locally and ready to use!"
            echo ""
            echo "Version: ${VERSION_OUTPUT}"
            echo ""
            echo "Installed at: ${INSTALL_PATH}"
            echo ""
            echo "Usage:"
            echo "  ./secureflow init       # Initialize configuration"
            echo "  ./secureflow encrypt    # Encrypt files"
            echo "  ./secureflow decrypt    # Decrypt files"
            echo "  ./secureflow --help     # Show help"
            echo ""
            echo -e "${YELLOW}Note:${NC} For global installation, run: curl -sSL https://raw.githubusercontent.com/MayR-Labs/secureflow-go/main/install.sh | bash -s -- --global"
            echo ""
            echo -e "For more information, visit: ${BLUE}https://github.com/${REPO}${NC}"
        else
            echo -e "${RED}✗${NC} Installation failed: ${INSTALL_PATH} not found"
        fi
    fi
}

# Main installation flow
main() {
    detect_platform
    get_latest_version
    install_binary
    verify_installation
}

# Run main function
main
