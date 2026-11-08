#!/usr/bin/env bash

# ==============================================================================
# Scrcpy & Android Device Diagnostic Tool
# Checks ADB, Scrcpy, USB Hardware, Kernel logs, and Wireless connectivity
# ==============================================================================

BOLD="\033[1m"
GREEN="\033[0;32m"
BLUE="\033[0;34m"
YELLOW="\033[0;33m"
RED="\033[0;31m"
RESET="\033[0m"

log_info()    { echo -e "${BLUE}[INFO]${RESET} $*"; }
log_ok()      { echo -e "${GREEN}[OK]${RESET} $*"; }
log_warn()    { echo -e "${YELLOW}[WARN]${RESET} $*"; }
log_fail()    { echo -e "${RED}[FAIL]${RESET} $*"; }

find_adb() {
    if command -v adb &>/dev/null; then
        command -v adb
        return
    fi
    local candidates=(
        "$HOME/Android/Sdk/platform-tools/adb"
        "$ANDROID_HOME/platform-tools/adb"
        "$ANDROID_SDK_ROOT/platform-tools/adb"
        "/usr/lib/android-sdk/platform-tools/adb"
        "/opt/android-sdk/platform-tools/adb"
    )
    for cand in "${candidates[@]}"; do
        if [ -x "$cand" ]; then
            echo "$cand"
            return
        fi
    done
    echo ""
}

find_scrcpy() {
    command -v scrcpy 2>/dev/null || echo ""
}

run_diagnostics() {
    echo -e "${BOLD}======================================================${RESET}"
    echo -e "${BOLD}         Android & Scrcpy Device Diagnostics          ${RESET}"
    echo -e "${BOLD}======================================================${RESET}"
    echo ""

    # 1. Check ADB
    local adb_bin
    adb_bin="$(find_adb)"
    if [ -n "$adb_bin" ]; then
        local adb_ver
        adb_ver="$("$adb_bin" version 2>/dev/null | head -n 1 || echo "unknown")"
        log_ok "ADB found at: $adb_bin"
        log_info "ADB Version: $adb_ver"
    else
        log_fail "ADB not found in PATH or standard Android SDK directories!"
        return 1
    fi

    # 2. Check Scrcpy
    local scrcpy_bin
    scrcpy_bin="$(find_scrcpy)"
    if [ -n "$scrcpy_bin" ]; then
        local scrcpy_ver
        scrcpy_ver="$("$scrcpy_bin" --version 2>/dev/null | head -n 1 || echo "unknown")"
        log_ok "Scrcpy found at: $scrcpy_bin"
        log_info "Scrcpy Version: $scrcpy_ver"
    else
        log_warn "Scrcpy binary not found in PATH! (install with: sudo apt install scrcpy)"
    fi

    echo ""
    echo -e "${BOLD}--- 1. ADB Connected Devices ---${RESET}"
    local raw_devices
    raw_devices="$("$adb_bin" devices -l 2>/dev/null || true)"
    echo "$raw_devices"

    local device_count
    device_count="$(echo "$raw_devices" | grep -v "List of devices attached" | grep -v "^$" | wc -l)"

    if [ "$device_count" -gt 0 ]; then
        echo ""
        log_ok "Detected $device_count device(s) via ADB."
        if echo "$raw_devices" | grep -q "unauthorized"; then
            log_warn "Device is 'unauthorized'. Check phone screen and allow USB debugging."
        elif echo "$raw_devices" | grep -q "offline"; then
            log_warn "Device is 'offline'. Try reconnecting cable or restart ADB."
        fi
        return 0
    fi

    log_warn "No devices recognized by ADB."
    echo ""

    # 3. Check USB Hardware Enumeration (lsusb)
    echo -e "${BOLD}--- 2. USB Hardware Detection (lsusb) ---${RESET}"
    local usb_matches
    usb_matches="$(lsusb 2>/dev/null | grep -i -E "xiaomi|vivo|samsung|oppo|realme|google|huawei|transsion|infinix|android" || true)"

    if [ -n "$usb_matches" ]; then
        log_ok "Found phone on USB bus:"
        echo "$usb_matches"
        echo ""
        log_info "Phone is physically detected on USB bus, but ADB cannot communicate with it."
        log_info "Checklist on phone:"
        echo "  1. Set USB mode to 'File Transfer' (MTP), not 'Charging only'."
        echo "  2. Go to Settings -> Developer Options -> Enable 'USB Debugging'."
        echo "  3. Look at phone screen for 'Allow USB Debugging' popup and tap Allow."
    else
        log_fail "No Android phone detected on USB bus by lsusb."
    fi

    echo ""

    # 4. Kernel USB logs (dmesg)
    echo -e "${BOLD}--- 3. Kernel USB Connection Logs (dmesg) ---${RESET}"
    local kernel_logs
    kernel_logs="$(dmesg 2>/dev/null | tail -n 35 | grep -i -E "usb|disconnect|Cannot enable|unable to enumerate|error -71" || true)"

    if [ -n "$kernel_logs" ]; then
        echo "$kernel_logs"
        echo ""
        if echo "$kernel_logs" | grep -q "Cannot enable. Maybe the USB cable is bad?"; then
            log_fail "KERNEL DETECTED HARDWARE/SIGNAL ERROR: 'Maybe the USB cable is bad?'"
            echo "  -> The USB cable or port keeps disconnecting during hardware handshake."
            echo "  -> Recommendation: Change USB cable or plug into a rear motherboard USB port."
        fi
        if echo "$kernel_logs" | grep -q "error -71"; then
            log_fail "KERNEL DETECTED PROTOCOL ERROR -71 (EPROTO):"
            echo "  -> Electrical signal dropped or packet transmission failed on USB bus."
        fi
    else
        log_info "No recent USB kernel errors logged."
    fi

    echo ""
    echo -e "${BOLD}--- 4. Troubleshooting Steps ---${RESET}"
    echo "1. Change the USB cable (use a known good data cable)."
    echo "2. Plug directly into the rear USB port of the PC."
    echo "3. Restart ADB server using: $0 --restart"
    echo "4. Or use Wireless Debugging: In Scrcpy GUI, click 'Pair Device...' then 'Connect IP...'"
}

restart_adb() {
    local adb_bin
    adb_bin="$(find_adb)"
    if [ -z "$adb_bin" ]; then
        log_fail "ADB not found."
        exit 1
    fi
    log_info "Restarting ADB daemon..."
    "$adb_bin" kill-server 2>/dev/null || true
    sleep 1
    "$adb_bin" start-server
    log_ok "ADB daemon restarted."
}

# Command dispatch
case "$1" in
    --restart|-r)
        restart_adb
        run_diagnostics
        ;;
    --watch|-w)
        log_info "Entering watch mode (refreshes every 2 seconds). Press Ctrl+C to exit."
        while true; do
            clear
            run_diagnostics
            sleep 2
        done
        ;;
    *)
        run_diagnostics
        ;;
esac
