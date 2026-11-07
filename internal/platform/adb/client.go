package adb

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/opelooo/scrcpy_GUI_Go/internal/domain"
)

// Client implements domain.ADBClient by executing system adb and scrcpy binaries.
type Client struct {
	ADBPath    string
	ScrcpyPath string
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

// NewClient initializes the ADB and scrcpy client, verifying both tools exist.
func NewClient() (*Client, error) {
	adb, err := findProgram("adb")
	if err != nil {
		return nil, fmt.Errorf("Cannot find adb in environment or folder!\nPlease install adb or place adb executable next to this app.\n%w", err)
	}

	scrcpy, err := findProgram("scrcpy")
	if err != nil {
		return nil, fmt.Errorf("Cannot find scrcpy in environment or folder!\nPlease install scrcpy or place scrcpy executable next to this app.\n%w", err)
	}

	return &Client{
		ADBPath:    adb,
		ScrcpyPath: scrcpy,
	}, nil
}

// ListDevices executes 'adb devices' and parses connected devices.
func (c *Client) ListDevices() ([]domain.Device, error) {
	cmd := exec.Command(c.ADBPath, "devices")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run adb devices: %w", err)
	}
	return ParseDevices(string(out)), nil
}

// GetManufacturer queries device manufacturer via getprop.
func (c *Client) GetManufacturer(serial string) (string, error) {
	if serial == "" {
		return "", errors.New("device serial required")
	}
	cmd := exec.Command(c.ADBPath, "-s", serial, "shell", "getprop", "ro.product.manufacturer")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get manufacturer: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

// GetModel queries device model via getprop.
func (c *Client) GetModel(serial string) (string, error) {
	if serial == "" {
		return "", errors.New("device serial required")
	}
	cmd := exec.Command(c.ADBPath, "-s", serial, "shell", "getprop", "ro.product.model")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get model: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

// GetIP queries device IP address using ip route on wlan interface.
func (c *Client) GetIP(serial string) (string, error) {
	if serial == "" {
		return "", errors.New("device serial required")
	}
	cmd := exec.Command(c.ADBPath, "-s", serial, "shell", "ip", "route")
	out, err := cmd.Output()
	if err == nil {
		if ip := ParseDeviceIP(string(out)); ip != "" {
			return ip, nil
		}
	}

	// Fallback to "ip -f inet addr show wlan0"
	cmdFallback := exec.Command(c.ADBPath, "-s", serial, "shell", "ip", "-f", "inet", "addr", "show", "wlan0")
	if outFallback, errFallback := cmdFallback.Output(); errFallback == nil {
		ipRe := regexp.MustCompile(`inet\s+((?:\d{1,3}\.){3}\d{1,3})`)
		if m := ipRe.FindStringSubmatch(string(outFallback)); len(m) > 1 {
			return m[1], nil
		}
	}

	return "", errors.New("could not find IP address for device (ensure Wi-Fi is connected)")
}

// ConfigureTCPIP restarts adb daemon on device listening on TCP/IP port.
func (c *Client) ConfigureTCPIP(serial string, port int) error {
	cmd := exec.Command(c.ADBPath, "-s", serial, "tcpip", fmt.Sprintf("%d", port))
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("adb tcpip failed: %s: %w", string(out), err)
	}
	return nil
}

// Connect connects to device over network target (e.g. 192.168.1.50:5555).
func (c *Client) Connect(target string) error {
	cmd := exec.Command(c.ADBPath, "connect", target)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("adb connect failed: %s: %w", string(out), err)
	}
	return nil
}

// ExecuteScrcpy executes scrcpy asynchronously with the given arguments.
func (c *Client) ExecuteScrcpy(args []string) error {
	cmd := exec.Command(c.ScrcpyPath, args...)
	return cmd.Start()
}
