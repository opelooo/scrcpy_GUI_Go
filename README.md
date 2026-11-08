# Minimal GUI Scrcpy (Go Version)

A lightweight, modern graphical user interface for [scrcpy](https://github.com/Genymobile/scrcpy) written in **Go** using [Fyne v2](https://fyne.io). Rewritten from the original Java/Swing GUI into a standalone native binary following standard Go project architecture conventions (`cmd/`, `internal/domain`, `internal/service`, `internal/platform`, `internal/ui`).

---

## Features

- **Zero JRE Dependency**: Built as a native compiled binary with fast startup and low memory footprint.
- **Idiomatic Go Architecture**: Separation of concerns with thin entry point, clean domain interfaces, and service orchestration.
- **Automatic Tool Detection**: Automatically detects `adb` and `scrcpy` from system `PATH` or the executable's folder.
- **Device Management**:
  - Automatically queries and lists connected Android devices (`adb devices`).
  - Background auto-detection of device connection/disconnection.
  - Transport badges: `[USB]` vs `[Wi-Fi]` visual indicators in the device list.
  - Double-click an attached device to immediately launch `scrcpy`.
  - Right-click context menu on any device:
    - **Show Device Info**: Displays device manufacturer and model.
    - **Show Device IP**: Displays device Wi-Fi IP address.
    - **Disconnect Device**: Disconnects wireless device (`adb disconnect <target>`).
- **Quick Action Buttons**:
  - **Refresh**: Rescans for attached devices.
  - **Info**: Shows device hardware details.
  - **Run**: Launches `scrcpy` with all selected options.
  - **TCP/IP**: Automates wireless connection over Wi-Fi (`adb tcpip 5555` $\rightarrow$ `adb connect <ip>:5555` $\rightarrow$ `scrcpy`).
  - **Connect IP**: Direct TCP/IP connect dialog without needing initial USB cable (`scrcpy --tcpip=IP[:PORT]`, with optional `+IP` force reconnect).
  - **Pair Device**: Android 11+ Wireless Debugging Pairing dialog (`adb pair <ip>:<port> <code>`).
- **Customizable Options**:
  - **Video Bit Rate**: Slider and text field synced from 1M to 50M (default: `2M`).
  - **Max Mirroring Size**: Input field (default: `1024`).
  - **Toggles**:
    - Video On / Off (`--no-video`)
    - Turn Screen On / Off (`--turn-screen-off`)
    - Stay Awake (`--stay-awake`)
  - **Recording**:
    - Record Video (`mp4`, `mkv`)
    - Record Audio (`m4a`, `wav`, `flac`, `opus`, `aac`, `mka`)
- **Dark Theme**: Clean, modern dark UI styling.

---

## Project Structure

```text
.
├── cmd/
│   └── scrcpy-gui/           # Thin entry point (flag parsing, wiring, startup)
│       ├── icon.png          # Embedded application icon
│       ├── lib/              # Dynamic library linker scripts for CGO on Linux
│       └── main.go
├── internal/
│   ├── domain/               # Domain models (Device, ScrcpyOptions) & interfaces
│   │   └── device.go
│   ├── platform/             # Low-level integrations & drivers
│   │   └── adb/              # ADB & Scrcpy CLI executor & output parsers
│   │       ├── client.go
│   │       ├── parser.go
│   │       └── parser_test.go
│   ├── service/              # Use-case business logic orchestrator
│   │   ├── scrcpy_service.go
│   │   └── scrcpy_service_test.go
│   └── ui/                   # Fyne UI presentation layer
│       ├── gui.go
│       └── gui_test.go
├── dist/                     # Compiled binaries
├── Makefile                  # Build, test, and package runner
├── go.mod
├── go.sum
└── README.md
```

---

## Prerequisites

1. Install [scrcpy](https://github.com/Genymobile/scrcpy) and `adb`:
   - **Linux (Ubuntu/Debian)**:
     ```bash
     sudo apt install scrcpy adb
     ```
   - **Windows**:
     Download the official scrcpy release containing `adb.exe` and `scrcpy.exe`. Put `scrcpy-gui.exe` in the same folder or add scrcpy to system `PATH`.

---

## Installation

The installer automatically detects your operating system (Linux or Windows) and sets up the application with desktop shortcuts and system integration.

### Linux / macOS
```bash
./install.sh
# or via make:
make install
```
- Installs binary to `~/.local/share/scrcpy-gui/scrcpy-gui`
- Creates command symlink at `~/.local/bin/scrcpy-gui`
- Installs application icon and desktop menu entry (`scrcpy-gui.desktop`)

To uninstall:
```bash
./uninstall.sh
# or:
make uninstall
```

### Windows
- **Via Command Prompt / File Explorer**: Double-click or run `install.bat`.
- **Via Git Bash / MSYS2**: Run `./install.sh`.
- Installs to `%LOCALAPPDATA%\scrcpy-gui`
- Creates shortcuts on Desktop and Start Menu
- Adds to User `PATH`

---

## Build & Run

### Linux
```bash
# Run tests & static analysis
make test
make vet

# Build binary into dist/scrcpy-gui
make build

# Run application
./dist/scrcpy-gui
```

### Windows (Cross-compilation using Zig)
```bash
make build-windows
# Binary created at dist/scrcpy-gui.exe
```

---

## License

GPL-3.0 License. See [LICENSE](LICENSE) for details.
