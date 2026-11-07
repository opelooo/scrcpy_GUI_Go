package main

import (
	"reflect"
	"strings"
	"testing"

	"fyne.io/fyne/v2/test"
)

func TestParseDeviceList(t *testing.T) {
	rawOutput := `* daemon not running; starting now at tcp:5037
* daemon started successfully
List of devices attached
RF8M12345XYZ	device
emulator-5554	device
offline_dev	offline
unauth_dev	unauthorized

`
	got := parseDeviceList(rawOutput)
	want := []string{"RF8M12345XYZ", "emulator-5554"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseDeviceList() = %v, want %v", got, want)
	}
}

func TestParseDeviceListEmpty(t *testing.T) {
	rawOutput := `List of devices attached

`
	got := parseDeviceList(rawOutput)
	if len(got) != 0 {
		t.Fatalf("expected empty slice, got %v", got)
	}
}

func TestParseDeviceIP(t *testing.T) {
	rawRoute := `default via 192.168.1.1 dev wlan0 proto static metric 100 
192.168.1.0/24 dev wlan0 proto kernel scope link src 192.168.1.142 metric 100 
`
	got := parseDeviceIP(rawRoute)
	want := "192.168.1.142"

	if got != want {
		t.Fatalf("parseDeviceIP() = %q, want %q", got, want)
	}
}

func TestBuildScrcpyArgs(t *testing.T) {
	r := &ADBRunner{ScrcpyPath: "scrcpy"}
	opts := ScrcpyOptions{
		DeviceSerial: "XYZ123",
		MaxSize:      "1024",
		BitRate:      "4M",
		VideoOn:      false,
		ScreenOn:     false,
		StayAwake:    true,
	}

	args := r.BuildScrcpyArgs(opts)
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

func TestBuildScrcpyArgsDefault(t *testing.T) {
	r := &ADBRunner{ScrcpyPath: "scrcpy"}
	opts := ScrcpyOptions{
		DeviceSerial: "DEV1",
		MaxSize:      "1024",
		BitRate:      "2M",
		VideoOn:      true,
		ScreenOn:     true,
		StayAwake:    false,
	}

	args := r.BuildScrcpyArgs(opts)
	joined := strings.Join(args, " ")

	if strings.Contains(joined, "--no-video") {
		t.Errorf("args should not contain --no-video when VideoOn is true")
	}
	if strings.Contains(joined, "--turn-screen-off") {
		t.Errorf("args should not contain --turn-screen-off when ScreenOn is true")
	}
	if strings.Contains(joined, "--stay-awake") {
		t.Errorf("args should not contain --stay-awake when StayAwake is false")
	}
}

func TestGUIInitialization(t *testing.T) {
	app := test.NewApp()
	w := app.NewWindow("Test")
	runner := &ADBRunner{ADBPath: "adb", ScrcpyPath: "scrcpy"}
	gui := NewGUI(app, w, runner)
	if gui == nil {
		t.Fatal("expected non-nil GUI")
	}
	gui.selectDevice("TEST_DEV_123")
	if gui.selectedDevice != "TEST_DEV_123" {
		t.Fatalf("expected selectedDevice to be TEST_DEV_123, got %s", gui.selectedDevice)
	}
}

