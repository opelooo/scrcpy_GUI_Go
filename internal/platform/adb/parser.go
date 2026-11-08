package adb

import (
	"bufio"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/opelooo/scrcpy_GUI_Go/internal/domain"
)

// IsTCPIPDevice checks whether the serial represents a network TCP/IP endpoint (e.g. 192.168.1.5:5555).
func IsTCPIPDevice(serial string) bool {
	return strings.Contains(serial, ":")
}

// ParseDevices parses the raw stdout of 'adb devices' into domain.Device slice with transport type.
func ParseDevices(output string) []domain.Device {
	var devices []domain.Device
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" ||
			strings.HasPrefix(line, "*") ||
			strings.HasPrefix(line, "List of devices attached") {
			continue
		}
		fields := strings.Fields(line)
		connType := domain.ConnUSB
		if IsTCPIPDevice(fields[0]) {
			connType = domain.ConnTCPIP
		}

		if len(fields) >= 2 {
			devices = append(devices, domain.Device{
				Serial:         fields[0],
				State:          fields[1],
				ConnectionType: connType,
			})
		} else if len(fields) == 1 {
			devices = append(devices, domain.Device{
				Serial:         fields[0],
				State:          "device",
				ConnectionType: connType,
			})
		}
	}
	return devices
}

// ParseDeviceIP extracts the IPv4 address from 'ip route' or 'ip addr' output.
func ParseDeviceIP(output string) string {
	ipRe := regexp.MustCompile(`src\s+((?:\d{1,3}\.){3}\d{1,3})`)
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "wlan") {
			if m := ipRe.FindStringSubmatch(line); len(m) > 1 {
				return m[1]
			}
		}
	}

	scanner = bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		if m := ipRe.FindStringSubmatch(line); len(m) > 1 && !strings.HasPrefix(m[1], "127.") {
			return m[1]
		}
	}
	return ""
}

func appendConnectionArgs(args []string, opts domain.ScrcpyOptions) []string {
	if opts.TCPIPAddress != "" {
		return append(args, fmt.Sprintf("--tcpip=%s", opts.TCPIPAddress))
	}
	if opts.SelectUSB {
		return append(args, "-d")
	}
	if opts.SelectTCPIP {
		return append(args, "-e")
	}
	if opts.DeviceSerial != "" {
		return append(args, "-s", opts.DeviceSerial)
	}
	return args
}

func appendVideoAdvancedArgs(args []string, opts domain.ScrcpyOptions) []string {
	if fps := strings.TrimSpace(opts.MaxFPS); fps != "" && fps != "Default" {
		args = append(args, fmt.Sprintf("--max-fps=%s", fps))
	}
	if codec := strings.TrimSpace(opts.VideoCodec); codec != "" && codec != "default" && codec != "h264 (default)" {
		args = append(args, fmt.Sprintf("--video-codec=%s", codec))
	}
	if orient := strings.TrimSpace(opts.Orientation); orient != "" && orient != "default" && orient != "Default" {
		args = append(args, fmt.Sprintf("--orientation=%s", orient))
	}
	if buf := strings.TrimSuffix(strings.TrimSpace(opts.VideoBuffer), "ms"); buf != "" {
		args = append(args, fmt.Sprintf("--video-buffer=%s", buf))
	}
	if crop := strings.TrimSpace(opts.Crop); crop != "" {
		args = append(args, fmt.Sprintf("--crop=%s", crop))
	}
	if sink := strings.TrimSpace(opts.V4L2Sink); sink != "" {
		args = append(args, fmt.Sprintf("--v4l2-sink=%s", sink))
	}
	if vbuf := strings.TrimSuffix(strings.TrimSpace(opts.V4L2Buffer), "ms"); vbuf != "" {
		args = append(args, fmt.Sprintf("--v4l2-buffer=%s", vbuf))
	}
	return args
}

func appendVideoArgs(args []string, opts domain.ScrcpyOptions) []string {
	if strings.TrimSpace(opts.MaxSize) != "" && strings.TrimSpace(opts.VideoSource) != "camera" {
		args = append(args, "-m", strings.TrimSpace(opts.MaxSize))
	}
	if strings.TrimSpace(opts.BitRate) != "" {
		args = append(args, "-b", strings.TrimSpace(opts.BitRate))
	}
	if !opts.VideoOn {
		args = append(args, "--no-video")
	}
	if !opts.ScreenOn {
		args = append(args, "--turn-screen-off")
	}
	if opts.StayAwake {
		args = append(args, "--stay-awake")
	}
	if opts.NoPlayback {
		args = append(args, "--no-playback")
	}
	return appendVideoAdvancedArgs(args, opts)
}

func appendAudioArgs(args []string, opts domain.ScrcpyOptions) []string {
	if !opts.AudioOn {
		return append(args, "--no-audio")
	}
	src := strings.TrimSpace(opts.AudioSource)
	if opts.AudioDup {
		src = "playback"
	}
	if src != "" && src != "default" && src != "output (default)" {
		args = append(args, fmt.Sprintf("--audio-source=%s", src))
	}
	if opts.AudioDup {
		args = append(args, "--audio-dup")
	}
	if codec := strings.TrimSpace(opts.AudioCodec); codec != "" && codec != "default" && codec != "opus (default)" {
		args = append(args, fmt.Sprintf("--audio-codec=%s", codec))
	}
	if rate := strings.TrimSpace(opts.AudioBitRate); rate != "" && rate != "default" && rate != "Default" {
		args = append(args, fmt.Sprintf("--audio-bit-rate=%s", rate))
	}
	if buf := strings.TrimSuffix(strings.TrimSpace(opts.AudioBuffer), "ms"); buf != "" {
		args = append(args, fmt.Sprintf("--audio-buffer=%s", buf))
	}
	if opts.NoAudioPlayback {
		args = append(args, "--no-audio-playback")
	}
	return args
}

func appendControlArgs(args []string, opts domain.ScrcpyOptions) []string {
	isNetwork := opts.SelectTCPIP || opts.TCPIPAddress != "" || strings.Contains(opts.DeviceSerial, ":")
	if opts.OTGMode && !isNetwork {
		args = append(args, "--otg")
	} else if opts.NoControl {
		return append(args, "--no-control")
	}
	if kbd := strings.TrimSpace(opts.KeyboardMode); kbd != "" && kbd != "default" && kbd != "Default" {
		args = append(args, fmt.Sprintf("--keyboard=%s", kbd))
	}
	if mouse := strings.TrimSpace(opts.MouseMode); mouse != "" && mouse != "default" && mouse != "Default" {
		args = append(args, fmt.Sprintf("--mouse=%s", mouse))
	}
	if gp := strings.TrimSpace(opts.GamepadMode); gp != "" && !strings.HasPrefix(strings.ToLower(gp), "default") {
		args = append(args, fmt.Sprintf("--gamepad=%s", gp))
	}
	if opts.NoMouseHover {
		args = append(args, "--no-mouse-hover")
	}
	if opts.NoClipboardAutosync {
		args = append(args, "--no-clipboard-autosync")
	}
	if opts.LegacyPaste {
		args = append(args, "--legacy-paste")
	}
	if opts.RawKeyEvents {
		args = append(args, "--raw-key-events")
	} else if opts.PreferText {
		args = append(args, "--prefer-text")
	}
	if opts.NoKeyRepeat {
		args = append(args, "--no-key-repeat")
	}
	if sm := strings.TrimSpace(opts.ShortcutMod); sm != "" && !strings.HasPrefix(strings.ToLower(sm), "default") {
		args = append(args, fmt.Sprintf("--shortcut-mod=%s", sm))
	}
	if target := strings.TrimSpace(opts.PushTarget); target != "" && target != "/sdcard/Download/" {
		args = append(args, fmt.Sprintf("--push-target=%s", target))
	}
	return args
}

func determineRecordFilename(opts domain.ScrcpyOptions) string {
	if custom := strings.TrimSpace(opts.RecordPath); custom != "" {
		return custom
	}
	ts := time.Now().Format("20060102_150405")
	if opts.RecordVideo && opts.VideoFormat != "" {
		return fmt.Sprintf("scrcpy_record_%s.%s", ts, opts.VideoFormat)
	}
	if opts.RecordAudio && opts.AudioFormat != "" {
		return fmt.Sprintf("scrcpy_record_%s.%s", ts, opts.AudioFormat)
	}
	return ""
}

func appendRecordOutputArgs(args []string, opts domain.ScrcpyOptions) []string {
	if !opts.RecordVideo && !opts.RecordAudio {
		return args
	}

	filename := determineRecordFilename(opts)
	if filename != "" {
		args = append(args, fmt.Sprintf("--record=%s", filename))
	}
	if fmtVal := strings.TrimSpace(opts.RecordFormat); fmtVal != "" && fmtVal != "auto" && fmtVal != "Auto" {
		args = append(args, fmt.Sprintf("--record-format=%s", fmtVal))
	}
	return args
}

func appendRecordOrientationArgs(args []string, opts domain.ScrcpyOptions) []string {
	if !opts.RecordVideo || !opts.VideoOn {
		return args
	}
	ro := strings.TrimSpace(opts.RecordOrientation)
	if ro != "" && ro != "default" && ro != "Default" {
		args = append(args, fmt.Sprintf("--record-orientation=%s", ro))
	}
	return args
}

func appendTimeLimitArgs(args []string, opts domain.ScrcpyOptions) []string {
	tl := strings.TrimSpace(opts.TimeLimit)
	if tl != "" && tl != "0" {
		args = append(args, fmt.Sprintf("--time-limit=%s", tl))
	}
	return args
}

func appendRecordArgs(args []string, opts domain.ScrcpyOptions) []string {
	args = appendTimeLimitArgs(args, opts)
	if !opts.VideoOn && !opts.AudioOn {
		return args
	}
	args = appendRecordOutputArgs(args, opts)
	args = appendRecordOrientationArgs(args, opts)
	return args
}

func appendCameraArgs(args []string, opts domain.ScrcpyOptions) []string {
	src := strings.TrimSpace(opts.VideoSource)
	if src != "camera" {
		return args
	}
	args = append(args, "--video-source=camera")
	if facing := strings.TrimSpace(opts.CameraFacing); facing != "" && facing != "default" && facing != "Default" {
		args = append(args, fmt.Sprintf("--camera-facing=%s", facing))
	}
	size := strings.TrimSpace(opts.CameraSize)
	if size == "" {
		size = "1920x1080"
	}
	args = append(args, fmt.Sprintf("--camera-size=%s", size))
	if fps := strings.TrimSpace(opts.CameraFPS); fps != "" && fps != "default" && fps != "Default" {
		args = append(args, fmt.Sprintf("--camera-fps=%s", fps))
	}
	return args
}

func appendDeviceArgs(args []string, opts domain.ScrcpyOptions) []string {
	if opts.ShowTouches {
		args = append(args, "--show-touches")
	}
	if opts.PowerOffOnClose {
		args = append(args, "--power-off-on-close")
	}
	if opts.NoPowerOn {
		args = append(args, "--no-power-on")
	}
	if timeout := strings.TrimSpace(opts.ScreenOffTimeout); timeout != "" && timeout != "0" {
		args = append(args, fmt.Sprintf("--screen-off-timeout=%s", timeout))
	}
	if app := strings.TrimSpace(opts.StartApp); app != "" {
		args = append(args, fmt.Sprintf("--start-app=%s", app))
	}
	if opts.KeepActive {
		args = append(args, "--keep-active")
	}
	return args
}

var scrcpyVersionRe = regexp.MustCompile(`scrcpy\s+v?(\d+)\.(\d+)(?:\.(\d+))?`)

// ParseScrcpyVersion parses scrcpy version string from scrcpy --version output.
func ParseScrcpyVersion(output string) (int, int, int, error) {
	m := scrcpyVersionRe.FindStringSubmatch(output)
	if len(m) < 3 {
		return 0, 0, 0, errors.New("invalid scrcpy version output")
	}
	major, err := strconv.Atoi(m[1])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid major version: %w", err)
	}
	minor, err := strconv.Atoi(m[2])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid minor version: %w", err)
	}
	patch := 0
	if len(m) > 3 && m[3] != "" {
		if p, err := strconv.Atoi(m[3]); err == nil {
			patch = p
		}
	}
	return major, minor, patch, nil
}

func appendWindowGeometryArgs(args []string, opts domain.ScrcpyOptions) []string {
	if x := strings.TrimSpace(opts.WindowX); x != "" {
		args = append(args, fmt.Sprintf("--window-x=%s", x))
	}
	if y := strings.TrimSpace(opts.WindowY); y != "" {
		args = append(args, fmt.Sprintf("--window-y=%s", y))
	}
	if w := strings.TrimSpace(opts.WindowWidth); w != "" {
		args = append(args, fmt.Sprintf("--window-width=%s", w))
	}
	if h := strings.TrimSpace(opts.WindowHeight); h != "" {
		args = append(args, fmt.Sprintf("--window-height=%s", h))
	}
	return args
}

func appendWindowAdvancedArgs(args []string, opts domain.ScrcpyOptions) []string {
	if opts.NoWindowAspectRatioLock {
		args = append(args, "--no-window-aspect-ratio-lock")
	}
	if bg := strings.TrimSpace(opts.BackgroundColor); bg != "" {
		args = append(args, fmt.Sprintf("--background-color=%s", bg))
	}
	if fit := strings.TrimSpace(opts.RenderFit); fit != "" && fit != "default" && fit != "Default" {
		args = append(args, fmt.Sprintf("--render-fit=%s", fit))
	}
	return args
}

func appendWindowArgs(args []string, opts domain.ScrcpyOptions) []string {
	if opts.NoWindow {
		return append(args, "--no-window")
	}
	if opts.AlwaysOnTop {
		args = append(args, "--always-on-top")
	}
	if opts.Fullscreen {
		args = append(args, "--fullscreen")
	}
	if opts.WindowBorderless {
		args = append(args, "--window-borderless")
	}
	if opts.DisableScreensaver {
		args = append(args, "--disable-screensaver")
	}
	if title := strings.TrimSpace(opts.WindowTitle); title != "" {
		args = append(args, fmt.Sprintf("--window-title=%s", title))
	}
	if !opts.Fullscreen {
		args = appendWindowGeometryArgs(args, opts)
	}
	return appendWindowAdvancedArgs(args, opts)
}

// BuildScrcpyArgs constructs the CLI arguments for scrcpy execution.
func BuildScrcpyArgs(opts domain.ScrcpyOptions) []string {
	var args []string
	args = appendConnectionArgs(args, opts)
	args = appendVideoArgs(args, opts)
	args = appendCameraArgs(args, opts)
	args = appendAudioArgs(args, opts)
	args = appendControlArgs(args, opts)
	args = appendDeviceArgs(args, opts)
	args = appendWindowArgs(args, opts)
	args = appendRecordArgs(args, opts)
	return args
}
