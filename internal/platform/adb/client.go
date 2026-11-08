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

func findInPath(name string) (string, bool) {
	if p, err := exec.LookPath(name); err == nil {
		return p, true
	}
	if runtime.GOOS == "windows" && !strings.HasSuffix(strings.ToLower(name), ".exe") {
		if p, err := exec.LookPath(name + ".exe"); err == nil {
			return p, true
		}
	}
	return "", false
}

func checkFileCandidate(cand string) (string, bool) {
	if cand == "" {
		return "", false
	}
	fi, err := os.Stat(cand)
	if err == nil && !fi.IsDir() {
		return cand, true
	}
	if runtime.GOOS == "windows" && !strings.HasSuffix(strings.ToLower(cand), ".exe") {
		candExe := cand + ".exe"
		if fiExe, errExe := os.Stat(candExe); errExe == nil && !fiExe.IsDir() {
			return candExe, true
		}
	}
	return "", false
}

func findInLocalDirs(name string) (string, bool) {
	candidates := []string{
		name,
		"./" + name,
	}
	if exe, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exe)
		candidates = append(candidates, filepath.Join(exeDir, name))
	}
	for _, cand := range candidates {
		if path, ok := checkFileCandidate(cand); ok {
			return path, true
		}
	}
	return "", false
}

func getCommonCandidateDirs() []string {
	var dirs []string
	if env := os.Getenv("ANDROID_HOME"); env != "" {
		dirs = append(dirs, filepath.Join(env, "platform-tools"))
	}
	if env := os.Getenv("ANDROID_SDK_ROOT"); env != "" {
		dirs = append(dirs, filepath.Join(env, "platform-tools"))
	}
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		dirs = append(dirs,
			filepath.Join(home, "Android", "Sdk", "platform-tools"),
			filepath.Join(home, "Android", "sdk", "platform-tools"),
			filepath.Join(home, ".android-sdk", "platform-tools"),
			filepath.Join(home, "Library", "Android", "sdk", "platform-tools"),
			filepath.Join(home, ".local", "bin"),
			filepath.Join(home, "bin"),
		)
	}
	if localApp := os.Getenv("LOCALAPPDATA"); localApp != "" {
		dirs = append(dirs, filepath.Join(localApp, "Android", "Sdk", "platform-tools"))
	}
	dirs = append(dirs,
		"/opt/android-sdk/platform-tools",
		"/usr/lib/android-sdk/platform-tools",
		"/usr/local/bin",
		"/opt/homebrew/bin",
	)
	return dirs
}

func findInCommonDirs(name string) (string, bool) {
	for _, dir := range getCommonCandidateDirs() {
		cand := filepath.Join(dir, name)
		if path, ok := checkFileCandidate(cand); ok {
			return path, true
		}
	}
	return "", false
}

func findViaWhereis(name string) (string, bool) {
	if runtime.GOOS == "windows" {
		return "", false
	}
	out, err := exec.Command("whereis", "-b", name).Output()
	if err != nil {
		return "", false
	}
	parts := strings.Fields(string(out))
	if len(parts) <= 1 {
		return "", false
	}
	for _, cand := range parts[1:] {
		if path, ok := checkFileCandidate(cand); ok {
			return path, true
		}
	}
	return "", false
}

// findProgram searches for an executable in PATH, working dir, common SDK locations, and whereis.
func findProgram(name string) (string, error) {
	if p, ok := findInPath(name); ok {
		return p, nil
	}
	if p, ok := findInLocalDirs(name); ok {
		return p, nil
	}
	if p, ok := findInCommonDirs(name); ok {
		return p, nil
	}
	if p, ok := findViaWhereis(name); ok {
		return p, nil
	}
	return "", fmt.Errorf("%s not found in PATH, standard SDK folders, or local directory", name)
}

// NewClient initializes the ADB and scrcpy client, verifying both tools exist.
func NewClient() (*Client, error) {
	adb, err := findProgram("adb")
	if err != nil {
		return nil, fmt.Errorf("Cannot find adb in environment or folder!\nPlease install adb or place adb executable next to this app.\n%w", err)
	}

	// Ensure adb directory is in PATH for downstream subprocesses (e.g. scrcpy)
	adbDir := filepath.Dir(adb)
	currPath := os.Getenv("PATH")
	if !strings.Contains(currPath, adbDir) {
		_ = os.Setenv("PATH", adbDir+string(os.PathListSeparator)+currPath)
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

// Disconnect disconnects a device over network target (e.g. 192.168.1.50:5555).
func (c *Client) Disconnect(target string) error {
	cmd := exec.Command(c.ADBPath, "disconnect", target)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("adb disconnect failed: %s: %w", string(out), err)
	}
	return nil
}

// Pair pairs a device using Android 11+ wireless pairing code.
func (c *Client) Pair(ip string, port int, code string) error {
	target := fmt.Sprintf("%s:%d", ip, port)
	cmd := exec.Command(c.ADBPath, "pair", target, code)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("adb pair failed: %s: %w", string(out), err)
	}
	return nil
}

// ExecuteScrcpy executes scrcpy asynchronously with the given arguments.
func (c *Client) ExecuteScrcpy(args []string) error {
	cmd := exec.Command(c.ScrcpyPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if runtime.GOOS == "linux" {
		const sysSDL2 = "/usr/lib/x86_64-linux-gnu/libSDL2-2.0.so.0"
		if fi, err := os.Stat(sysSDL2); err == nil && !fi.IsDir() {
			cmd.Env = append(os.Environ(), "LD_PRELOAD="+sysSDL2)
		}
	}
	return cmd.Start()
}

// GetScrcpyCapabilities queries scrcpy version and determines supported flags.
func (c *Client) GetScrcpyCapabilities() (domain.ScrcpyCapabilities, error) {
	if c.ScrcpyPath == "" {
		return domain.ScrcpyCapabilities{}, errors.New("scrcpy binary not found")
	}

	cmdVer := exec.Command(c.ScrcpyPath, "--version")
	outVer, err := cmdVer.Output()
	if err != nil {
		return domain.ScrcpyCapabilities{}, fmt.Errorf("failed to get scrcpy version: %w", err)
	}

	raw := strings.TrimSpace(string(outVer))
	major, minor, patch, err := ParseScrcpyVersion(raw)
	if err != nil {
		return domain.ScrcpyCapabilities{VersionRaw: raw}, nil
	}

	caps := domain.ScrcpyCapabilities{
		VersionRaw: raw,
		Major:      major,
		Minor:      minor,
		Patch:      patch,
	}

	if major > 3 || (major == 3 && minor >= 4) {
		caps.HasKeepActive = true
		caps.HasWindowAdvanced = true
		caps.HasFlexDisplay = true
		return caps, nil
	}

	cmdHelp := exec.Command(c.ScrcpyPath, "--help")
	if outHelp, errHelp := cmdHelp.Output(); errHelp == nil {
		helpText := string(outHelp)
		if strings.Contains(helpText, "--keep-active") {
			caps.HasKeepActive = true
		}
		if strings.Contains(helpText, "--render-fit") {
			caps.HasWindowAdvanced = true
		}
		if strings.Contains(helpText, "--flex-display") {
			caps.HasFlexDisplay = true
		}
	}

	return caps, nil
}
