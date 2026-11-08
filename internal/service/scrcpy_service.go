package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/opelooo/scrcpy_GUI_Go/internal/domain"
	"github.com/opelooo/scrcpy_GUI_Go/internal/platform/adb"
)

// ScrcpyService coordinates device management and scrcpy invocation.
type ScrcpyService struct {
	client domain.ADBClient
}

// NewScrcpyService creates a new domain.DeviceService implementation.
func NewScrcpyService(client domain.ADBClient) *ScrcpyService {
	return &ScrcpyService{
		client: client,
	}
}

// ListDevices lists attached Android devices.
func (s *ScrcpyService) ListDevices() ([]domain.Device, error) {
	devices, err := s.client.ListDevices()
	if err != nil {
		return nil, fmt.Errorf("failed to list devices: %w", err)
	}

	// Filter to active devices
	var active []domain.Device
	for _, d := range devices {
		if d.State == "device" {
			active = append(active, d)
		}
	}
	return active, nil
}

// GetDeviceInfo queries and formats device manufacturer and model.
func (s *ScrcpyService) GetDeviceInfo(serial string) (string, error) {
	if serial == "" {
		return "", errors.New("please select device")
	}

	mfr, err := s.client.GetManufacturer(serial)
	if err != nil {
		return "", err
	}

	model, err := s.client.GetModel(serial)
	if err == nil && model != "" && model != mfr {
		return fmt.Sprintf("%s (%s)", mfr, model), nil
	}

	return mfr, nil
}

// GetDeviceIP queries device IP address for wireless connection.
func (s *ScrcpyService) GetDeviceIP(serial string) (string, error) {
	if serial == "" {
		return "", errors.New("please select device")
	}
	return s.client.GetIP(serial)
}

// ConnectTCPIP sets up wireless ADB and launches scrcpy over TCP/IP.
func (s *ScrcpyService) ConnectTCPIP(serial string) error {
	if serial == "" {
		return errors.New("please select device")
	}

	ip, err := s.client.GetIP(serial)
	if err != nil {
		return err
	}

	if err := s.client.ConfigureTCPIP(serial, 5555); err != nil {
		return err
	}

	time.Sleep(1200 * time.Millisecond)

	target := fmt.Sprintf("%s:5555", ip)
	if err := s.client.Connect(target); err != nil {
		return err
	}

	time.Sleep(500 * time.Millisecond)

	// Launch scrcpy --tcpip=<ip>:5555 with fallback to -s <target>
	err = s.client.ExecuteScrcpy([]string{fmt.Sprintf("--tcpip=%s", target)})
	if err != nil {
		return s.client.ExecuteScrcpy([]string{"-s", target})
	}

	return nil
}

// LaunchScrcpy builds arguments and executes scrcpy asynchronously.
func (s *ScrcpyService) LaunchScrcpy(opts domain.ScrcpyOptions) error {
	hasTarget := opts.DeviceSerial != "" || opts.SelectUSB || opts.SelectTCPIP || opts.TCPIPAddress != ""
	if !hasTarget && opts.TunnelHost == "" && opts.ADBServerHost == "" {
		return errors.New("please select device")
	}

	args := adb.BuildScrcpyArgs(opts)
	var extraEnv []string
	if host := strings.TrimSpace(opts.ADBServerHost); host != "" {
		if !strings.HasPrefix(host, "tcp:") {
			host = "tcp:" + host
		}
		extraEnv = append(extraEnv, "ADB_SERVER_SOCKET="+host)
	}
	return s.client.ExecuteScrcpyWithEnv(args, extraEnv)
}

// DirectConnect connects to an arbitrary IP:port over wireless ADB and launches scrcpy.
func (s *ScrcpyService) DirectConnect(opts domain.ConnectOptions) error {
	if opts.IP == "" {
		return errors.New("IP address required")
	}
	port := opts.Port
	if port <= 0 {
		port = 5555
	}
	target := fmt.Sprintf("%s:%d", opts.IP, port)

	if err := s.client.Connect(target); err != nil {
		return fmt.Errorf("connect to %s: %w", target, err)
	}

	tcpipArg := target
	if opts.ForceReconnect {
		tcpipArg = "+" + target
	}

	err := s.client.ExecuteScrcpy([]string{fmt.Sprintf("--tcpip=%s", tcpipArg)})
	if err != nil {
		return s.client.ExecuteScrcpy([]string{"-s", target})
	}
	return nil
}

// PairDevice pairs an Android 11+ device using its IP, pairing port, and 6-digit code.
func (s *ScrcpyService) PairDevice(opts domain.PairOptions) error {
	if opts.IP == "" {
		return errors.New("IP address required")
	}
	if opts.Port <= 0 {
		return errors.New("valid pairing port required")
	}
	if opts.Code == "" {
		return errors.New("pairing code required")
	}
	return s.client.Pair(opts.IP, opts.Port, opts.Code)
}

// DisconnectDevice disconnects a wireless device target.
func (s *ScrcpyService) DisconnectDevice(target string) error {
	if target == "" {
		return errors.New("device target required")
	}
	return s.client.Disconnect(target)
}

// GetScrcpyCapabilities queries the underlying client for scrcpy binary capabilities.
func (s *ScrcpyService) GetScrcpyCapabilities() (domain.ScrcpyCapabilities, error) {
	return s.client.GetScrcpyCapabilities()
}
