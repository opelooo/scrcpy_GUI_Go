#!/usr/bin/env bash
set -eo pipefail

# ==============================================================================
# Scrcpy GUI - Cross-Platform Installer (Linux & Windows Auto-Detection)
# ==============================================================================

APP_NAME="Scrcpy GUI"
BINARY_NAME="scrcpy-gui"
APP_DIR_NAME="scrcpy-gui"
VERSION="1.1.0"

BOLD="\033[1m"
GREEN="\033[0;32m"
BLUE="\033[0;34m"
YELLOW="\033[0;33m"
RED="\033[0;31m"
RESET="\033[0m"

log_info()    { echo -e "${BLUE}[INFO]${RESET} $*"; }
log_success() { echo -e "${GREEN}[SUCCESS]${RESET} $*"; }
log_warn()    { echo -e "${YELLOW}[WARN]${RESET} $*"; }
log_error()   { echo -e "${RED}[ERROR]${RESET} $*"; }

# ------------------------------------------------------------------------------
# 1. OS Detection
# ------------------------------------------------------------------------------
detect_os() {
    local uname_s
    uname_s="$(uname -s 2>/dev/null || echo "")"

    case "$uname_s" in
        Linux*)
            echo "linux"
            ;;
        Darwin*)
            echo "darwin"
            ;;
        CYGWIN*|MINGW*|MSYS*|Windows_NT*)
            echo "windows"
            ;;
        *)
            if [[ "$OS" == "Windows_NT" ]] || [ -n "$WINDIR" ] || [ -n "$SYSTEMROOT" ]; then
                echo "windows"
            else
                echo "unknown"
            fi
            ;;
    esac
}

OS_TYPE="$(detect_os)"

# Optional override via argument (e.g. ./install.sh --os linux or --os windows)
if [[ "$1" == "--os" && -n "$2" ]]; then
    OS_TYPE="$2"
    shift 2
fi

log_info "Detected Operating System: ${BOLD}${OS_TYPE}${RESET}"

# ------------------------------------------------------------------------------
# 2. Linux Installation
# ------------------------------------------------------------------------------
install_linux() {
    echo "=== Installing $APP_NAME for Linux ==="

    local install_dir="$HOME/.local/share/$APP_DIR_NAME"
    local bin_dir="$HOME/.local/bin"
    local desktop_dir="$HOME/.local/share/applications"
    local icon_hicolor="$HOME/.local/share/icons/hicolor/256x256/apps"
    local icon_dir="$HOME/.local/share/icons"

    # Build binary if not already built in dist/
    if [ ! -f "dist/$BINARY_NAME" ]; then
        log_info "Compiling $APP_NAME for Linux..."
        mkdir -p dist
        go build -ldflags="-s -w" -o "dist/$BINARY_NAME" ./cmd/scrcpy-gui
        log_success "Build complete: dist/$BINARY_NAME"
    fi

    log_info "Creating target directories..."
    mkdir -p "$install_dir" "$bin_dir" "$desktop_dir" "$icon_hicolor" "$icon_dir"

    log_info "Installing binary and assets..."
    rm -f "$install_dir/$BINARY_NAME"
    cp "dist/$BINARY_NAME" "$install_dir/$BINARY_NAME"
    chmod +x "$install_dir/$BINARY_NAME"

    # Create symlink in ~/.local/bin
    ln -sf "$install_dir/$BINARY_NAME" "$bin_dir/$BINARY_NAME"

    # Copy icon
    if [ -f "icon.png" ]; then
        cp "icon.png" "$icon_hicolor/$BINARY_NAME.png"
        cp "icon.png" "$icon_dir/$BINARY_NAME.png"
    fi

    # Create desktop shortcut
    log_info "Creating desktop application entry..."
    cat <<EOF > "$desktop_dir/$BINARY_NAME.desktop"
[Desktop Entry]
Version=$VERSION
Type=Application
Name=$APP_NAME
GenericName=Android Screen Mirror
Comment=Minimalist GUI for scrcpy and adb
Exec=$install_dir/$BINARY_NAME
Icon=$BINARY_NAME
Terminal=false
Categories=Utility;Development;
Keywords=scrcpy;android;mirror;adb;
EOF
    chmod +x "$desktop_dir/$BINARY_NAME.desktop"

    # Update desktop database
    if command -v update-desktop-database &> /dev/null; then
        update-desktop-database "$desktop_dir" 2>/dev/null || true
    fi

    echo ""
    log_success "=== Installation Complete! ==="
    echo -e "Application: ${BOLD}$APP_NAME${RESET}"
    echo -e "Binary:      ${BOLD}$install_dir/$BINARY_NAME${RESET}"
    echo -e "Symlink:     ${BOLD}$bin_dir/$BINARY_NAME${RESET}"
    echo -e "Shortcut:    ${BOLD}$desktop_dir/$BINARY_NAME.desktop${RESET}"
    echo ""
    echo "You can launch '$APP_NAME' from your Application Menu or run '$BINARY_NAME' in terminal."
    if [[ ":$PATH:" != *":$HOME/.local/bin:"* ]]; then
        log_warn "Note: $HOME/.local/bin is not in your current PATH. Add it to ~/.bashrc or ~/.zshrc if needed."
    fi
}

# ------------------------------------------------------------------------------
# 3. Windows Installation (Git Bash / MSYS2 / WSL / MINGW)
# ------------------------------------------------------------------------------
install_windows() {
    echo "=== Installing $APP_NAME for Windows ==="

    local win_exe="${BINARY_NAME}.exe"
    local dist_exe="dist/$win_exe"

    # Build binary if not already built in dist/
    if [ ! -f "$dist_exe" ]; then
        log_info "Compiling $APP_NAME for Windows..."
        mkdir -p dist
        if command -v zig &>/dev/null && [[ "$(uname -s)" == Linux* ]]; then
            CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC="zig cc -target x86_64-windows-gnu" go build -ldflags="-H=windowsgui -s -w" -o "$dist_exe" ./cmd/scrcpy-gui
        else
            go build -ldflags="-H=windowsgui -s -w" -o "$dist_exe" ./cmd/scrcpy-gui
        fi
        log_success "Build complete: $dist_exe"
    fi

    # Determine LocalAppData path
    local target_dir=""
    if [ -n "$LOCALAPPDATA" ]; then
        target_dir="$LOCALAPPDATA/$APP_DIR_NAME"
    elif command -v powershell.exe &>/dev/null; then
        local raw_localappdata
        raw_localappdata="$(powershell.exe -NoProfile -Command "[Environment]::GetFolderPath('LocalApplicationData')" | tr -d '\r')"
        target_dir="$(cygpath -u "$raw_localappdata" 2>/dev/null || echo "$raw_localappdata")/$APP_DIR_NAME"
    else
        target_dir="$USERPROFILE/AppData/Local/$APP_DIR_NAME"
    fi

    log_info "Target directory: $target_dir"
    mkdir -p "$target_dir"

    log_info "Copying executable and assets..."
    cp "$dist_exe" "$target_dir/$win_exe"
    if [ -f "icon.png" ]; then
        cp "icon.png" "$target_dir/icon.png"
    fi

    # Create Windows Shortcuts via PowerShell if available
    if command -v powershell.exe &>/dev/null; then
        log_info "Creating Desktop and Start Menu shortcuts via PowerShell..."
        powershell.exe -NoProfile -ExecutionPolicy Bypass -Command "
            \$WshShell = New-Object -ComObject WScript.Shell

            # Paths
            \$installDir = [Environment]::ExpandEnvironmentVariables(\"%LOCALAPPDATA%\\$APP_DIR_NAME\")
            \$targetExe   = Join-Path \$installDir '$win_exe'
            \$desktopPath = [Environment]::GetFolderPath('Desktop')
            \$startMenu   = [Environment]::GetFolderPath('Programs')

            # Desktop Shortcut
            \$desktopLink = Join-Path \$desktopPath '$APP_NAME.lnk'
            \$sc1 = \$WshShell.CreateShortcut(\$desktopLink)
            \$sc1.TargetPath = \$targetExe
            \$sc1.WorkingDirectory = \$installDir
            \$sc1.Save()

            # Start Menu Shortcut
            \$startLink = Join-Path \$startMenu '$APP_NAME.lnk'
            \$sc2 = \$WshShell.CreateShortcut(\$startLink)
            \$sc2.TargetPath = \$targetExe
            \$sc2.WorkingDirectory = \$installDir
            \$sc2.Save()

            # Add to User PATH if not present
            \$userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
            if (\$userPath -notlike \"*\$installDir*\") {
                [Environment]::SetEnvironmentVariable('Path', \$userPath + \";\$installDir\", 'User')
            }
        " || log_warn "Failed to create PowerShell shortcuts. Manual shortcut creation may be required."
    fi

    echo ""
    log_success "=== Installation Complete! ==="
    echo -e "Application: ${BOLD}$APP_NAME${RESET}"
    echo -e "Installed to: ${BOLD}$target_dir${RESET}"
    echo "Shortcuts have been placed on your Desktop and Start Menu."
}

# ------------------------------------------------------------------------------
# Dispatch based on detected OS
# ------------------------------------------------------------------------------
case "$OS_TYPE" in
    linux)
        install_linux
        ;;
    windows)
        install_windows
        ;;
    darwin)
        log_warn "macOS detected. Running Unix-style installation..."
        install_linux
        ;;
    *)
        log_error "Unsupported or undetected Operating System: '$OS_TYPE'"
        echo "Usage: ./install.sh [--os linux|windows]"
        exit 1
        ;;
esac

