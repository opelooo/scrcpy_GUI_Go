#!/usr/bin/env bash
set -eo pipefail

APP_NAME="Scrcpy GUI"
BINARY_NAME="scrcpy-gui"
APP_DIR_NAME="scrcpy-gui"

BOLD="\033[1m"
GREEN="\033[0;32m"
BLUE="\033[0;34m"
RESET="\033[0m"

detect_os() {
    local uname_s
    uname_s="$(uname -s 2>/dev/null || echo "")"
    case "$uname_s" in
        Linux*) echo "linux" ;;
        Darwin*) echo "darwin" ;;
        CYGWIN*|MINGW*|MSYS*|Windows_NT*) echo "windows" ;;
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

echo "=== Uninstalling $APP_NAME ($OS_TYPE) ==="

if [[ "$OS_TYPE" == "linux" || "$OS_TYPE" == "darwin" ]]; then
    rm -rf "$HOME/.local/share/$APP_DIR_NAME"
    rm -f "$HOME/.local/bin/$BINARY_NAME"
    rm -f "$HOME/.local/share/applications/$BINARY_NAME.desktop"
    rm -f "$HOME/.local/share/icons/hicolor/256x256/apps/$BINARY_NAME.png"
    rm -f "$HOME/.local/share/icons/$BINARY_NAME.png"

    if command -v update-desktop-database &> /dev/null; then
        update-desktop-database "$HOME/.local/share/applications" 2>/dev/null || true
    fi
elif [[ "$OS_TYPE" == "windows" ]]; then
    target_dir="$LOCALAPPDATA/$APP_DIR_NAME"
    rm -rf "$target_dir"
    if command -v powershell.exe &>/dev/null; then
        powershell.exe -NoProfile -Command "
            \$desktop = [Environment]::GetFolderPath('Desktop')
            \$programs = [Environment]::GetFolderPath('Programs')
            Remove-Item -Force -ErrorAction SilentlyContinue (Join-Path \$desktop '$APP_NAME.lnk')
            Remove-Item -Force -ErrorAction SilentlyContinue (Join-Path \$programs '$APP_NAME.lnk')
        " || true
    fi
fi

echo -e "${GREEN}[SUCCESS]${RESET} $APP_NAME has been uninstalled."

