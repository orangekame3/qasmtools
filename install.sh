#!/bin/bash

# qasmtools installation script
# Usage: curl -sSfL https://raw.githubusercontent.com/orangekame3/qasmtools/main/install.sh | sh

set -e

# Default settings
OWNER="orangekame3"
REPO="qasmtools"
BINARY_NAME="qasm"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Helper functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

# Detect OS and architecture
detect_os() {
    case "$(uname -s)" in
        Darwin*)    echo "Darwin" ;;
        Linux*)     echo "Linux" ;;
        CYGWIN*|MINGW*|MSYS*) echo "Windows" ;;
        *)          log_error "Unsupported OS: $(uname -s)" ;;
    esac
}

detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64)   echo "x86_64" ;;
        aarch64|arm64)  echo "arm64" ;;
        *)              log_error "Unsupported architecture: $(uname -m)" ;;
    esac
}

# Get the latest release version from GitHub API
get_latest_version() {
    local api_url="https://api.github.com/repos/${OWNER}/${REPO}/releases/latest"
    
    if command -v curl >/dev/null 2>&1; then
        curl -s "$api_url" | grep '"tag_name":' | sed -E 's/.*"tag_name": "([^"]+)".*/\1/'
    elif command -v wget >/dev/null 2>&1; then
        wget -qO- "$api_url" | grep '"tag_name":' | sed -E 's/.*"tag_name": "([^"]+)".*/\1/'
    else
        log_error "curl or wget is required to download qasmtools"
    fi
}

# Download and install binary
install_binary() {
    local version="$1"
    local os="$2"
    local arch="$3"
    
    # Remove 'v' prefix if present
    version="${version#v}"
    
    # Determine file extension and extraction method
    local archive_name="qasmtools_${os}_${arch}"
    local temp_dir="/tmp/qasmtools_install_$$"
    
    if [ "$os" = "Windows" ]; then
        archive_name="${archive_name}.zip"
    else
        archive_name="${archive_name}.tar.gz"
    fi
    
    local download_url="https://github.com/${OWNER}/${REPO}/releases/download/v${version}/${archive_name}"
    local temp_archive="/tmp/${archive_name}"
    
    log_info "Downloading ${BINARY_NAME} v${version} for ${os}/${arch}..."
    
    if command -v curl >/dev/null 2>&1; then
        curl -sSfL "$download_url" -o "$temp_archive"
    elif command -v wget >/dev/null 2>&1; then
        wget -q "$download_url" -O "$temp_archive"
    else
        log_error "curl or wget is required to download qasmtools"
    fi
    
    if [ ! -f "$temp_archive" ]; then
        log_error "Failed to download ${BINARY_NAME}"
    fi
    
    # Create temporary extraction directory
    mkdir -p "$temp_dir"
    
    # Extract archive
    log_info "Extracting ${archive_name}..."
    
    if [ "$os" = "Windows" ]; then
        if command -v unzip >/dev/null 2>&1; then
            unzip -q "$temp_archive" -d "$temp_dir"
        else
            log_error "unzip is required to extract Windows packages"
        fi
        binary_path="${temp_dir}/${BINARY_NAME}.exe"
    else
        if command -v tar >/dev/null 2>&1; then
            tar -xzf "$temp_archive" -C "$temp_dir"
        else
            log_error "tar is required to extract packages"
        fi
        binary_path="${temp_dir}/${BINARY_NAME}"
    fi
    
    # Check if binary exists
    if [ ! -f "$binary_path" ]; then
        log_error "Binary not found in archive: $binary_path"
    fi
    
    # Make executable
    chmod +x "$binary_path"
    
    # Create install directory if it doesn't exist
    if [ ! -d "$INSTALL_DIR" ]; then
        log_info "Creating install directory: $INSTALL_DIR"
        mkdir -p "$INSTALL_DIR" || log_error "Failed to create install directory. Try running with sudo or set INSTALL_DIR environment variable."
    fi
    
    # Install binary
    log_info "Installing ${BINARY_NAME} to ${INSTALL_DIR}..."
    
    local target_binary="${INSTALL_DIR}/${BINARY_NAME}"
    if [ "$os" = "Windows" ]; then
        target_binary="${target_binary}.exe"
    fi
    
    if mv "$binary_path" "$target_binary"; then
        log_info "Successfully installed ${BINARY_NAME} to ${target_binary}"
    else
        log_error "Failed to install ${BINARY_NAME}. Try running with sudo or set INSTALL_DIR environment variable."
    fi
    
    # Clean up
    rm -rf "$temp_dir" "$temp_archive"
    
    # Verify installation
    if command -v "${BINARY_NAME}" >/dev/null 2>&1; then
        local installed_version
        installed_version=$(${BINARY_NAME} --version 2>/dev/null | head -n1 || echo "unknown")
        log_info "Installation verified: ${installed_version}"
    else
        log_warn "${BINARY_NAME} installed to ${INSTALL_DIR} but not found in PATH"
        log_warn "You may need to add ${INSTALL_DIR} to your PATH or restart your shell"
    fi
}

# Main installation process
main() {
    log_info "Installing qasmtools..."
    
    # Check if running as root (for system-wide installation)
    if [ "$EUID" -eq 0 ] && [ "$INSTALL_DIR" = "/usr/local/bin" ]; then
        log_info "Running as root, installing system-wide"
    elif [ "$INSTALL_DIR" = "/usr/local/bin" ] && [ ! -w "$INSTALL_DIR" ]; then
        log_warn "No write permission to ${INSTALL_DIR}"
        log_warn "Consider running with sudo or setting INSTALL_DIR to a writable directory:"
        log_warn "  export INSTALL_DIR=\$HOME/.local/bin"
        log_warn "  curl -sSfL https://raw.githubusercontent.com/orangekame3/qasmtools/main/install.sh | sh"
        log_error "Installation aborted"
    fi
    
    # Detect system
    local os
    local arch
    os=$(detect_os)
    arch=$(detect_arch)
    
    log_info "Detected system: ${os}/${arch}"
    
    # Get latest version
    local version
    version=$(get_latest_version)
    
    if [ -z "$version" ]; then
        log_error "Failed to get latest version information"
    fi
    
    log_info "Latest version: ${version}"
    
    # Install binary
    install_binary "$version" "$os" "$arch"
    
    log_info "Installation complete!"
    log_info ""
    log_info "Usage:"
    log_info "  ${BINARY_NAME} fmt file.qasm        # Format QASM file"
    log_info "  ${BINARY_NAME} lint file.qasm       # Lint QASM file"
    log_info "  ${BINARY_NAME} --help               # Show help"
    log_info ""
    log_info "For more information, visit: https://github.com/${OWNER}/${REPO}"
}

# Handle command line arguments
case "${1:-}" in
    --help|-h)
        echo "qasmtools installation script"
        echo ""
        echo "Usage:"
        echo "  curl -sSfL https://raw.githubusercontent.com/orangekame3/qasmtools/main/install.sh | sh"
        echo ""
        echo "Environment variables:"
        echo "  INSTALL_DIR    Installation directory (default: /usr/local/bin)"
        echo ""
        echo "Examples:"
        echo "  # Install to custom directory"
        echo "  export INSTALL_DIR=\$HOME/.local/bin"
        echo "  curl -sSfL https://raw.githubusercontent.com/orangekame3/qasmtools/main/install.sh | sh"
        exit 0
        ;;
    --version|-v)
        echo "qasmtools installation script v1.0.0"
        exit 0
        ;;
esac

# Run main installation
main "$@"