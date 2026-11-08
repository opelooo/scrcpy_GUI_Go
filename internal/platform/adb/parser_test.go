package adb

import (
	"reflect"
	"strings"
	"testing"

	"github.com/opelooo/scrcpy_GUI_Go/internal/domain"
)

func TestParseDevices(t *testing.T) {
	rawOutput := `* daemon not running; starting now at tcp:5037
* daemon started successfully
List of devices attached
RF8M12345XYZ	device
emulator-5554	device
192.168.1.105:5555	device
offline_dev	offline
unauth_dev	unauthorized

`
	got := ParseDevices(rawOutput)
	want := []domain.Device{
		{Serial: "RF8M12345XYZ", State: "device", ConnectionType: domain.ConnUSB},
		{Serial: "emulator-5554", State: "device", ConnectionType: domain.ConnUSB},
		{Serial: "192.168.1.105:5555", State: "device", ConnectionType: domain.ConnTCPIP},
		{Serial: "offline_dev", State: "offline", ConnectionType: domain.ConnUSB},
		{Serial: "unauth_dev", State: "unauthorized", ConnectionType: domain.ConnUSB},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ParseDevices() = %+v, want %+v", got, want)
	}
}

func TestParseDeviceIP(t *testing.T) {
	rawRoute := `default via 192.168.1.1 dev wlan0 proto static metric 100 
192.168.1.0/24 dev wlan0 proto kernel scope link src 192.168.1.142 metric 100 
`
	got := ParseDeviceIP(rawRoute)
	want := "192.168.1.142"

	if got != want {
		t.Fatalf("ParseDeviceIP() = %q, want %q", got, want)
	}
}

func TestBuildScrcpyArgs(t *testing.T) {
	opts := domain.ScrcpyOptions{
		DeviceSerial: "XYZ123",
		MaxSize:      "1024",
		BitRate:      "4M",
		VideoOn:      false,
		ScreenOn:     false,
		StayAwake:    true,
	}

	args := BuildScrcpyArgs(opts)
	expectedSubstrings := []string{
		"-s", "XYZ123",
		"-m", "1024",
		"-b", "4M",
		"--no-video",
		"--turn-screen-off",
		"--stay-awake",
	}

	joined := strings.Join(args, " ")
	for _, sub := range expectedSubstrings {
		if !strings.Contains(joined, sub) {
			t.Errorf("expected args to contain %q, but got %q", sub, joined)
		}
	}
}

func TestBuildScrcpyArgsConnectionSelectors(t *testing.T) {
	// 1. TCPIP address
	tcpOpts := domain.ScrcpyOptions{
		TCPIPAddress: "192.168.1.55:5555",
	}
	tcpArgs := strings.Join(BuildScrcpyArgs(tcpOpts), " ")
	if !strings.Contains(tcpArgs, "--tcpip=192.168.1.55:5555") {
		t.Errorf("expected --tcpip flag, got %s", tcpArgs)
	}

	// 2. Select USB
	usbOpts := domain.ScrcpyOptions{
		SelectUSB: true,
	}
	usbArgs := strings.Join(BuildScrcpyArgs(usbOpts), " ")
	if !strings.Contains(usbArgs, "-d") {
		t.Errorf("expected -d flag, got %s", usbArgs)
	}

	// 3. Select TCPIP
	tcpipOpts := domain.ScrcpyOptions{
		SelectTCPIP: true,
	}
	tcpipArgs := strings.Join(BuildScrcpyArgs(tcpipOpts), " ")
	if !strings.Contains(tcpipArgs, "-e") {
		t.Errorf("expected -e flag, got %s", tcpipArgs)
	}
}
