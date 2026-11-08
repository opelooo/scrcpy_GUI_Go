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

// BuildScrcpyArgs constructs the CLI arguments for scrcpy execution.
func BuildScrcpyArgs(opts domain.ScrcpyOptions) []string {
	var args []string
	if opts.TCPIPAddress != "" {
		args = append(args, fmt.Sprintf("--tcpip=%s", opts.TCPIPAddress))
	} else if opts.SelectUSB {
		args = append(args, "-d")
	} else if opts.SelectTCPIP {
		args = append(args, "-e")
	} else if opts.DeviceSerial != "" {
		args = append(args, "-s", opts.DeviceSerial)
	}

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
	if opts.RecordVideo && opts.VideoFormat != "" {
		ts := time.Now().Format("20060102_150405")
		filename := fmt.Sprintf("scrcpy_record_%s.%s", ts, opts.VideoFormat)
		args = append(args, fmt.Sprintf("--record=%s", filename))
	} else if opts.RecordAudio && opts.AudioFormat != "" {
		ts := time.Now().Format("20060102_150405")
		filename := fmt.Sprintf("scrcpy_record_%s.%s", ts, opts.AudioFormat)
		args = append(args, fmt.Sprintf("--record=%s", filename))
	}
	return args
}
