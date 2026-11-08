package adb

import (
	"bufio"
	"fmt"
	"regexp"
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
	return args
}

func appendVideoArgs(args []string, opts domain.ScrcpyOptions) []string {
	if strings.TrimSpace(opts.MaxSize) != "" {
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
	if src := strings.TrimSpace(opts.AudioSource); src != "" && src != "default" && src != "output (default)" {
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
	if opts.NoControl {
		return append(args, "--no-control")
	}
	if kbd := strings.TrimSpace(opts.KeyboardMode); kbd != "" && kbd != "default" && kbd != "Default" {
		args = append(args, fmt.Sprintf("--keyboard=%s", kbd))
	}
	if mouse := strings.TrimSpace(opts.MouseMode); mouse != "" && mouse != "default" && mouse != "Default" {
		args = append(args, fmt.Sprintf("--mouse=%s", mouse))
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
	if target := strings.TrimSpace(opts.PushTarget); target != "" && target != "/sdcard/Download/" {
		args = append(args, fmt.Sprintf("--push-target=%s", target))
	}
	return args
}

func appendRecordArgs(args []string, opts domain.ScrcpyOptions) []string {
	if opts.RecordVideo && opts.VideoFormat != "" {
		ts := time.Now().Format("20060102_150405")
		filename := fmt.Sprintf("scrcpy_record_%s.%s", ts, opts.VideoFormat)
		return append(args, fmt.Sprintf("--record=%s", filename))
	}
	if opts.RecordAudio && opts.AudioFormat != "" {
		ts := time.Now().Format("20060102_150405")
		filename := fmt.Sprintf("scrcpy_record_%s.%s", ts, opts.AudioFormat)
		return append(args, fmt.Sprintf("--record=%s", filename))
	}
	return args
}

// BuildScrcpyArgs constructs the CLI arguments for scrcpy execution.
func BuildScrcpyArgs(opts domain.ScrcpyOptions) []string {
	var args []string
	args = appendConnectionArgs(args, opts)
	args = appendVideoArgs(args, opts)
	args = appendAudioArgs(args, opts)
	args = appendControlArgs(args, opts)
	args = appendRecordArgs(args, opts)
	return args
}
