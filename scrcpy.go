package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

// ADBRunner encapsulates adb and scrcpy command execution.
type ADBRunner struct {
	ADBPath    string
	ScrcpyPath string
}

// ScrcpyOptions defines user parameters to launch scrcpy.
type ScrcpyOptions struct {
	DeviceSerial string
	MaxSize      string
	BitRate      string
	VideoOn      bool
	ScreenOn     bool
	StayAwake    bool
	RecordVideo  bool
	VideoFormat  string
	RecordAudio  bool
	AudioFormat  string
}

// findProgram searches for an executable in PATH, working dir, and executable dir.
func findProgram(name string) (string, error) {
	if p, err := exec.LookPath(name); err == nil {
		return p, nil
	}

	if runtime.GOOS == "windows" && !strings.HasSuffix(strings.ToLower(name), ".exe") {
		if p, err := exec.LookPath(name + ".exe"); err == nil {
			return p, nil
		}
	}

	candidates := []string{
		name,
		"./" + name,
	}
	if runtime.GOOS == "windows" {
		candidates = append(candidates, name+".exe", "./"+name+".exe")
	}

	if exe, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exe)
		candidates = append(candidates, filepath.Join(exeDir, name))
		if runtime.GOOS == "windows" {
			candidates = append(candidates, filepath.Join(exeDir, name+".exe"))
		}
	}

	for _, cand := range candidates {
		if fi, err := os.Stat(cand); err == nil && !fi.IsDir() {
			return cand, nil
		}
	}

	return "", fmt.Errorf("%s not found in PATH or local directory", name)
}

// NewADBRunner initializes runner and validates that adb and scrcpy exist.
func NewADBRunner() (*ADBRunner, error) {
	adb, err := findProgram("adb")
	if err != nil {
		return nil, fmt.Errorf("Cannot find adb in environment or folder!\nPlease install adb or place adb executable next to this app.\n%w", err)
	}

	scrcpy, err := findProgram("scrcpy")
	if err != nil {
		return nil, fmt.Errorf("Cannot find scrcpy in environment or folder!\nPlease install scrcpy or place scrcpy executable next to this app.\n%w", err)
	}

	return &ADBRunner{
		ADBPath:    adb,
		ScrcpyPath: scrcpy,
	}, nil
}

// ListDevices executes 'adb devices' and returns a list of connected device IDs.
func (r *ADBRunner) ListDevices() ([]string, error) {
	cmd := exec.Command(r.ADBPath, "devices")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run adb devices: %w", err)
	}
	return parseDeviceList(string(out)), nil
}

func parseDeviceList(output string) []string {
	var devices []string
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" ||
			strings.HasPrefix(line, "*") ||
			strings.HasPrefix(line, "List of devices attached") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[1] == "device" {
			devices = append(devices, fields[0])
		} else if len(fields) == 1 {
			devices = append(devices, fields[0])
		}
	}
	return devices
}

// GetDeviceInfo queries device manufacturer (and model if available).
func (r *ADBRunner) GetDeviceInfo(serial string) (string, error) {
	if serial == "" {
		return "", errors.New("please select device")
	}
	cmd := exec.Command(r.ADBPath, "-s", serial, "shell", "getprop", "ro.product.manufacturer")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get manufacturer: %w", err)
	}
	mfr := strings.TrimSpace(string(out))

	cmdModel := exec.Command(r.ADBPath, "-s", serial, "shell", "getprop", "ro.product.model")
	if outModel, errModel := cmdModel.Output(); errModel == nil {
		model := strings.TrimSpace(string(outModel))
		if model != "" && model != mfr {
			return fmt.Sprintf("%s (%s)", mfr, model), nil
		}
	}
	return mfr, nil
}

// GetDeviceIP queries device IP address using ip route on wlan interface.
func (r *ADBRunner) GetDeviceIP(serial string) (string, error) {
	if serial == "" {
		return "", errors.New("please select device")
	}
	cmd := exec.Command(r.ADBPath, "-s", serial, "shell", "ip", "route")
	out, err := cmd.Output()
	if err == nil {
		if ip := parseDeviceIP(string(out)); ip != "" {
			return ip, nil
		}
	}

	cmdFallback := exec.Command(r.ADBPath, "-s", serial, "shell", "ip", "-f", "inet", "addr", "show", "wlan0")
	if outFallback, errFallback := cmdFallback.Output(); errFallback == nil {
		ipRe := regexp.MustCompile(`inet\s+((?:\d{1,3}\.){3}\d{1,3})`)
		if m := ipRe.FindStringSubmatch(string(outFallback)); len(m) > 1 {
			return m[1], nil
		}
	}

	return "", errors.New("could not find IP address for device (ensure Wi-Fi is connected)")
}

func parseDeviceIP(output string) string {
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

// BuildScrcpyArgs constructs the command-line arguments for scrcpy.
func (r *ADBRunner) BuildScrcpyArgs(opts ScrcpyOptions) []string {
	var args []string
	if opts.DeviceSerial != "" {
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

// RunScrcpy executes scrcpy with the specified options asynchronously.
func (r *ADBRunner) RunScrcpy(opts ScrcpyOptions) error {
	if opts.DeviceSerial == "" {
		return errors.New("please select device")
	}
	args := r.BuildScrcpyArgs(opts)
	cmd := exec.Command(r.ScrcpyPath, args...)
	return cmd.Start()
}

// ConnectTCPIP configures and connects to the device via TCP/IP and launches scrcpy.
func (r *ADBRunner) ConnectTCPIP(serial string) error {
	if serial == "" {
		return errors.New("please select device")
	}
	ip, err := r.GetDeviceIP(serial)
	if err != nil {
		return err
	}

	cmdTcpip := exec.Command(r.ADBPath, "-s", serial, "tcpip", "5555")
	if out, err := cmdTcpip.CombinedOutput(); err != nil {
		return fmt.Errorf("adb tcpip failed: %s: %w", string(out), err)
	}

	time.Sleep(1200 * time.Millisecond)

	target := fmt.Sprintf("%s:5555", ip)
	cmdConnect := exec.Command(r.ADBPath, "connect", target)
	if out, err := cmdConnect.CombinedOutput(); err != nil {
		return fmt.Errorf("adb connect failed: %s: %w", string(out), err)
	}

	time.Sleep(500 * time.Millisecond)

	cmdScrcpy := exec.Command(r.ScrcpyPath, fmt.Sprintf("--tcpip=%s", target))
	if err := cmdScrcpy.Start(); err != nil {
		cmdScrcpyFallback := exec.Command(r.ScrcpyPath, "-s", target)
		if errFallback := cmdScrcpyFallback.Start(); errFallback != nil {
			return fmt.Errorf("failed to start scrcpy: %w", err)
		}
	}

	return nil
}
