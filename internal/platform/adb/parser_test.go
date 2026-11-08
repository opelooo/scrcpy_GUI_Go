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

func TestBuildScrcpyArgsVideoFeatures(t *testing.T) {
	opts := domain.ScrcpyOptions{
		DeviceSerial: "ABC456",
		VideoOn:      true,
		ScreenOn:     true,
		MaxFPS:       "60",
		VideoCodec:   "h265",
		Orientation:  "90",
		VideoBuffer:  "50ms",
		Crop:         "1080:1920:0:0",
		NoPlayback:   true,
	}

	args := BuildScrcpyArgs(opts)
	expectedSubstrings := []string{
		"-s", "ABC456",
		"--max-fps=60",
		"--video-codec=h265",
		"--orientation=90",
		"--video-buffer=50",
		"--crop=1080:1920:0:0",
		"--no-playback",
	}

	joined := strings.Join(args, " ")
	for _, sub := range expectedSubstrings {
		if !strings.Contains(joined, sub) {
			t.Errorf("expected args to contain %q, but got %q", sub, joined)
		}
	}
}

func TestBuildScrcpyArgsAudioFeatures(t *testing.T) {
	opts := domain.ScrcpyOptions{
		DeviceSerial:    "DEV_AUDIO",
		AudioOn:         true,
		AudioSource:     "mic",
		AudioDup:        true,
		AudioCodec:      "aac",
		AudioBitRate:    "192K",
		AudioBuffer:     "100ms",
		NoAudioPlayback: true,
	}

	args := BuildScrcpyArgs(opts)
	expectedSubstrings := []string{
		"-s", "DEV_AUDIO",
		"--audio-source=mic",
		"--audio-dup",
		"--audio-codec=aac",
		"--audio-bit-rate=192K",
		"--audio-buffer=100",
		"--no-audio-playback",
	}

	joined := strings.Join(args, " ")
	for _, sub := range expectedSubstrings {
		if !strings.Contains(joined, sub) {
			t.Errorf("expected args to contain %q, but got %q", sub, joined)
		}
	}
}

func TestBuildScrcpyArgsAudioDisabled(t *testing.T) {
	opts := domain.ScrcpyOptions{
		DeviceSerial: "DEV_NO_AUDIO",
		AudioOn:      false,
	}

	args := BuildScrcpyArgs(opts)
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "--no-audio") {
		t.Errorf("expected args to contain --no-audio, got %q", joined)
	}
}

func TestBuildScrcpyArgsControlFeatures(t *testing.T) {
	opts := domain.ScrcpyOptions{
		DeviceSerial:        "DEV_CTRL",
		NoControl:           false,
		KeyboardMode:        "uhid",
		MouseMode:           "uhid",
		NoMouseHover:        true,
		NoClipboardAutosync: true,
		LegacyPaste:         true,
		PushTarget:          "/sdcard/Movies/",
	}

	args := BuildScrcpyArgs(opts)
	expectedSubstrings := []string{
		"-s", "DEV_CTRL",
		"--keyboard=uhid",
		"--mouse=uhid",
		"--no-mouse-hover",
		"--no-clipboard-autosync",
		"--legacy-paste",
		"--push-target=/sdcard/Movies/",
	}

	joined := strings.Join(args, " ")
	for _, sub := range expectedSubstrings {
		if !strings.Contains(joined, sub) {
			t.Errorf("expected args to contain %q, but got %q", sub, joined)
		}
	}
}

func TestBuildScrcpyArgsControlDisabled(t *testing.T) {
	opts := domain.ScrcpyOptions{
		DeviceSerial: "DEV_RO",
		NoControl:    true,
	}

	args := BuildScrcpyArgs(opts)
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "--no-control") {
		t.Errorf("expected args to contain --no-control, got %q", joined)
	}
}

func TestBuildScrcpyArgsV4L2Features(t *testing.T) {
	opts := domain.ScrcpyOptions{
		DeviceSerial: "DEV_V4L2",
		V4L2Sink:     "/dev/video2",
		V4L2Buffer:   "50ms",
	}

	args := BuildScrcpyArgs(opts)
	joined := strings.Join(args, " ")
	expected := []string{"--v4l2-sink=/dev/video2", "--v4l2-buffer=50"}
	for _, sub := range expected {
		if !strings.Contains(joined, sub) {
			t.Errorf("expected args to contain %q, but got %q", sub, joined)
		}
	}
}

func TestBuildScrcpyArgsV4L2Empty(t *testing.T) {
	opts := domain.ScrcpyOptions{
		DeviceSerial: "DEV_V4L2",
		V4L2Sink:     "",
		V4L2Buffer:   "",
	}

	args := BuildScrcpyArgs(opts)
	joined := strings.Join(args, " ")
	if strings.Contains(joined, "--v4l2-sink") || strings.Contains(joined, "--v4l2-buffer") {
		t.Errorf("expected args not to contain v4l2 flags, got %q", joined)
	}
}

func TestBuildScrcpyArgsCameraSource(t *testing.T) {
	opts := domain.ScrcpyOptions{
		DeviceSerial: "DEV_CAM",
		VideoSource:  "camera",
		CameraFacing: "front",
		CameraSize:   "1920x1080",
		CameraFPS:    "60",
	}

	args := BuildScrcpyArgs(opts)
	joined := strings.Join(args, " ")
	expected := []string{
		"--video-source=camera",
		"--camera-facing=front",
		"--camera-size=1920x1080",
		"--camera-fps=60",
	}
	for _, sub := range expected {
		if !strings.Contains(joined, sub) {
			t.Errorf("expected args to contain %q, but got %q", sub, joined)
		}
	}
}

func TestBuildScrcpyArgsCameraDisabled(t *testing.T) {
	opts := domain.ScrcpyOptions{
		DeviceSerial: "DEV_CAM",
		VideoSource:  "display",
		CameraFacing: "front",
	}

	args := BuildScrcpyArgs(opts)
	joined := strings.Join(args, " ")
	if strings.Contains(joined, "--video-source=camera") || strings.Contains(joined, "--camera-facing") {
		t.Errorf("expected camera flags to be omitted when VideoSource is display, got %q", joined)
	}
}
