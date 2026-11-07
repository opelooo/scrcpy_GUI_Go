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
offline_dev	offline
unauth_dev	unauthorized

`
	got := ParseDevices(rawOutput)
	want := []domain.Device{
		{Serial: "RF8M12345XYZ", State: "device"},
		{Serial: "emulator-5554", State: "device"},
		{Serial: "offline_dev", State: "offline"},
		{Serial: "unauth_dev", State: "unauthorized"},
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
