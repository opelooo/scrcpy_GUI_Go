# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0] - 2026-11-08

### Added
- **Interactive Shortcuts Cheatsheet**: Added interactive modal cheatsheet dialog displaying all official scrcpy keyboard and mouse shortcuts (Window, Navigation, System, Clipboard, Mouse) with dynamic active modifier substitution.
- **Shortcuts Modifier Customization**: Expanded `--shortcut-mod` presets with `Default (Alt/Super)`, `rctrl`, `lctrl+lalt`, `lalt`, `ralt`, `lsuper`, `rsuper`, and `lctrl+lsuper`.
- **Camera Mirroring Suite (`camera.md`)**: Full support for Android 12+ Camera API via `--video-source=camera`, `--camera-id`, `--camera-facing`, `--camera-size`, `--camera-ar`, `--camera-fps`, and `--camera-high-speed`.
- **Cross-Platform V4L2 Isolation**: Conditional V4L2 Sink & Buffer activation on Linux (`runtime.GOOS == "linux"`); automatically suppresses V4L2 flags on Windows/macOS with helpful UI guidance.
- **Clean OTG Mode & Control-Only Mode (`otg.md`)**:
  - Implemented strict argument isolation for `--otg` hardware USB keyboard/mouse mode to prevent scrcpy fatal errors from forbidden mirroring/media flags.
  - Added `ControlOnly` mode (`--no-video --no-audio`) for wireless UHID remote control.
- **Network Tunnels & Remote ADB Server (`tunnels.md`)**: Added support for `--tunnel-host`, `--tunnel-port`, `--force-adb-forward`, `-p, --port`, and `ADB_SERVER_SOCKET` environment variable.
- **Virtual Display Management (`virtual-display.md`)**: Added `--new-display` secondary screen creation with custom resolution/DPI, `--no-vd-system-decorations`, `--no-vd-destroy-content`, `--display-ime-policy`, and version-gated `--flex-display`.
- **Recording Enhancements (`recording.md`)**: Added custom record output filename path, `--record-format`, `--record-orientation`, and `--time-limit`.
- **Window Management Suite (`window.md`)**: Added `--always-on-top`, `--fullscreen`, `--window-borderless`, `--disable-screensaver`, `--no-window`, window title, geometry placement (X, Y, W, H), and version-gated `--no-window-aspect-ratio-lock`, `--background-color`, and `--render-fit`.
- **Device Management Suite (`device.md`)**: Added `--show-touches`, `--power-off-on-close`, `--screen-off-timeout`, `--start-app`, `--no-power-on`, and version-gated `--keep-active`.
- **Scrcpy Version Detection**: Added automated semver capability probe with backward compatibility fallback checks.

### Changed
- **UI Tab Consolidation**: Merged 8 scattered tabs into 4 cohesive, logically grouped tabs:
  - `Media & Record`: Video Engine, Camera, Audio Engine, and Recording.
  - `Display & Window`: Window Modes, Geometry, and Virtual Secondary Display.
  - `Control`: Keyboard, Mouse, Gamepad, Shortcuts Cheatsheet, OTG Mode, Control Only, and Input Modifiers.
  - `Device & Network`: Device Power/Screen Behaviors, Auto-Start App, and Network/ADB Tunnels.
- **Responsive Layout**: Replaced rigid fixed-size grid wrappers with flexible Border layout, allowing the options panel to dynamically fill window space.
- **Window Dimensions**: Expanded default main window dimensions to 1060 x 560 to eliminate tab overflow and allow flexible user resizing.
- **Shortcuts Table Spacing**: Adjusted column widths (Action 180px, Shortcut 330px, Description flexible) to completely eliminate text overlapping.

### Fixed
- Fixed tab header overflow that forced options into an overflow menu (`...`).
- Fixed text collision and overlap in shortcuts cheatsheet dialog.
- Fixed `-m` substring conflict with `--mouse=disabled` in parser tests.
- Fixed argument conflicts between `--camera-id` and `--camera-facing`.
- Fixed argument conflicts between `--camera-size` and `--camera-ar` / `-m`.

---

## [1.0.0] - 2026-11-08

### Added
- Complete rewrite of scrcpy GUI into Go using Fyne v2.
- Clean Architecture layout conforming to idiomatic Go guidelines.
- Device discovery and connection management: USB, Wi-Fi TCP/IP, direct IP connect, and Android 11+ wireless pairing.
- Transport badges (`[USB]`, `[Wi-Fi]`) on attached devices list.
- Basic mirroring controls: Bitrate slider, Max size, Video ON/OFF.
- Audio forwarding controls: Audio ON/OFF, Audio source, Audio codec, Audio bitrate.
- Control engine: Keyboard mode, Mouse mode, Gamepad mode, Read-only mode.
- Linux and Windows cross-platform installation scripts (`install.sh`, `install.bat`, `Makefile`).

