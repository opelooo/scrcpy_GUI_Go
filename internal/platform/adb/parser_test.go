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
		AudioSource:     "playback",
		AudioDup:        true,
		AudioCodec:      "aac",
		AudioBitRate:    "192K",
		AudioBuffer:     "100ms",
		NoAudioPlayback: true,
	}

	args := BuildScrcpyArgs(opts)
	expectedSubstrings := []string{
		"-s", "DEV_AUDIO",
		"--audio-source=playback",
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
		RawKeyEvents:        true,
		NoKeyRepeat:         true,
		ShortcutMod:         "lctrl+lalt",
		OTGMode:             true,
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
		"--raw-key-events",
		"--no-key-repeat",
		"--shortcut-mod=lctrl+lalt",
		"--otg",
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

func TestBuildScrcpyArgsShortcutModDefault(t *testing.T) {
	opts := domain.ScrcpyOptions{
		DeviceSerial: "DEV_MOD",
		ShortcutMod:  "Default (Alt/Super)",
	}

	args := BuildScrcpyArgs(opts)
	joined := strings.Join(args, " ")
	if strings.Contains(joined, "--shortcut-mod") {
		t.Errorf("expected --shortcut-mod to be omitted when Default (Alt/Super), got %q", joined)
	}
}

func TestBuildScrcpyArgsGamepadFeatures(t *testing.T) {
	testCases := []struct {
		mode     string
		expected string
	}{
		{"uhid", "--gamepad=uhid"},
		{"aoa", "--gamepad=aoa"},
		{"disabled", "--gamepad=disabled"},
	}

	for _, tc := range testCases {
		opts := domain.ScrcpyOptions{
			DeviceSerial: "DEV_GP",
			GamepadMode:  tc.mode,
		}
		args := BuildScrcpyArgs(opts)
		joined := strings.Join(args, " ")
		if !strings.Contains(joined, tc.expected) {
			t.Errorf("expected %q in args, got %q", tc.expected, joined)
		}
	}
}

func TestBuildScrcpyArgsGamepadDefault(t *testing.T) {
	opts := domain.ScrcpyOptions{
		DeviceSerial: "DEV_GP",
		GamepadMode:  "default",
	}

	args := BuildScrcpyArgs(opts)
	joined := strings.Join(args, " ")
	if strings.Contains(joined, "--gamepad") {
		t.Errorf("expected --gamepad to be omitted when default, got %q", joined)
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

func TestBuildScrcpyArgsCameraDefaultSize(t *testing.T) {
	opts := domain.ScrcpyOptions{
		DeviceSerial: "DEV_CAM",
		VideoSource:  "camera",
		MaxSize:      "1024",
	}

	args := BuildScrcpyArgs(opts)
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "--camera-size=1920x1080") {
		t.Errorf("expected default camera size 1920x1080, got %q", joined)
	}
	if strings.Contains(joined, "-m 1024") {
		t.Errorf("expected -m to be omitted in camera mode, got %q", joined)
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

func TestConflictPreferTextAndRawKeyEvents(t *testing.T) {
	opts := domain.ScrcpyOptions{
		DeviceSerial: "DEV",
		PreferText:   true,
		RawKeyEvents: true,
	}

	args := BuildScrcpyArgs(opts)
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "--raw-key-events") {
		t.Errorf("expected --raw-key-events to be present, got %q", joined)
	}
	if strings.Contains(joined, "--prefer-text") {
		t.Errorf("expected --prefer-text to be omitted when raw-key-events is set, got %q", joined)
	}
}

func TestConflictOTGAndNoControl(t *testing.T) {
	opts := domain.ScrcpyOptions{
		DeviceSerial: "DEV",
		OTGMode:      true,
		NoControl:    true,
	}

	args := BuildScrcpyArgs(opts)
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "--otg") {
		t.Errorf("expected --otg to be present, got %q", joined)
	}
	if strings.Contains(joined, "--no-control") {
		t.Errorf("expected --no-control to be omitted in OTG mode, got %q", joined)
	}
}

func TestConflictOTGOverNetwork(t *testing.T) {
	opts := domain.ScrcpyOptions{
		DeviceSerial: "192.168.1.10:5555",
		OTGMode:      true,
	}

	args := BuildScrcpyArgs(opts)
	joined := strings.Join(args, " ")
	if strings.Contains(joined, "--otg") {
		t.Errorf("expected --otg to be omitted over network, got %q", joined)
	}
}

func TestConflictAudioDupPlaybackForced(t *testing.T) {
	opts := domain.ScrcpyOptions{
		DeviceSerial: "DEV",
		AudioOn:      true,
		AudioSource:  "output",
		AudioDup:     true,
	}

	args := BuildScrcpyArgs(opts)
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "--audio-source=playback") {
		t.Errorf("expected --audio-source=playback when audio-dup is set, got %q", joined)
	}
	if !strings.Contains(joined, "--audio-dup") {
		t.Errorf("expected --audio-dup to be present, got %q", joined)
	}
}

func TestConflictRecordWhenMediaDisabled(t *testing.T) {
	opts := domain.ScrcpyOptions{
		DeviceSerial: "DEV",
		VideoOn:      false,
		AudioOn:      false,
		RecordVideo:  true,
		VideoFormat:  "mp4",
	}

	args := BuildScrcpyArgs(opts)
	joined := strings.Join(args, " ")
	if strings.Contains(joined, "--record") {
		t.Errorf("expected --record to be omitted when both video and audio are off, got %q", joined)
	}
}

func TestParseScrcpyVersion(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		major   int
		minor   int
		patch   int
		wantErr bool
	}{
		{
			name:    "Standard Release 3.3.4",
			input:   "scrcpy 3.3.4 <https://github.com/Genymobile/scrcpy>",
			major:   3,
			minor:   3,
			patch:   4,
			wantErr: false,
		},
		{
			name:    "Two-segment 3.4",
			input:   "scrcpy 3.4 <https://github.com/Genymobile/scrcpy>",
			major:   3,
			minor:   4,
			patch:   0,
			wantErr: false,
		},
		{
			name:    "Prefixed v2.5.1",
			input:   "scrcpy v2.5.1",
			major:   2,
			minor:   5,
			patch:   1,
			wantErr: false,
		},
		{
			name:    "Invalid Output",
			input:   "some random text",
			major:   0,
			minor:   0,
			patch:   0,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			maj, min, ptch, err := ParseScrcpyVersion(tc.input)
			if (err != nil) != tc.wantErr {
				t.Fatalf("expected error: %v, got: %v", tc.wantErr, err)
			}
			if maj != tc.major || min != tc.minor || ptch != tc.patch {
				t.Errorf("expected (%d, %d, %d), got (%d, %d, %d)", tc.major, tc.minor, tc.patch, maj, min, ptch)
			}
		})
	}
}

func TestBuildScrcpyArgsDeviceFeatures(t *testing.T) {
	opts := domain.ScrcpyOptions{
		DeviceSerial:     "DEV_DEV",
		ShowTouches:      true,
		PowerOffOnClose:  true,
		NoPowerOn:        true,
		ScreenOffTimeout: "120",
		StartApp:         "org.mozilla.firefox",
		KeepActive:       true,
	}

	args := BuildScrcpyArgs(opts)
	joined := strings.Join(args, " ")
	expectedFlags := []string{
		"-s DEV_DEV",
		"--show-touches",
		"--power-off-on-close",
		"--no-power-on",
		"--screen-off-timeout=120",
		"--start-app=org.mozilla.firefox",
		"--keep-active",
	}

	for _, flag := range expectedFlags {
		if !strings.Contains(joined, flag) {
			t.Errorf("expected flag %q to be in args, got %q", flag, joined)
		}
	}
}

func TestBuildScrcpyArgsDeviceFeaturesDefaults(t *testing.T) {
	opts := domain.ScrcpyOptions{
		DeviceSerial: "DEV_DEV",
	}

	args := BuildScrcpyArgs(opts)
	joined := strings.Join(args, " ")
	prohibitedFlags := []string{
		"--show-touches",
		"--power-off-on-close",
		"--no-power-on",
		"--screen-off-timeout",
		"--start-app",
		"--keep-active",
	}

	for _, flag := range prohibitedFlags {
		if strings.Contains(joined, flag) {
			t.Errorf("flag %q should not be present with defaults, got %q", flag, joined)
		}
	}
}

func TestBuildScrcpyArgsWindowFeatures(t *testing.T) {
	opts := domain.ScrcpyOptions{
		DeviceSerial:       "DEV_WINDOW",
		AlwaysOnTop:        true,
		Fullscreen:         true,
		WindowBorderless:   true,
		DisableScreensaver: true,
		WindowTitle:        "My Phone Mirror",
	}

	args := BuildScrcpyArgs(opts)
	joined := strings.Join(args, " ")
	expectedFlags := []string{
		"-s DEV_WINDOW",
		"--always-on-top",
		"--fullscreen",
		"--window-borderless",
		"--disable-screensaver",
		"--window-title=My Phone Mirror",
	}

	for _, flag := range expectedFlags {
		if !strings.Contains(joined, flag) {
			t.Errorf("expected flag %q to be present, got %q", flag, joined)
		}
	}
}

func TestBuildScrcpyArgsWindowGeometry(t *testing.T) {
	opts := domain.ScrcpyOptions{
		DeviceSerial: "DEV_WINDOW",
		WindowX:      "100",
		WindowY:      "200",
		WindowWidth:  "1280",
		WindowHeight: "720",
	}

	args := BuildScrcpyArgs(opts)
	joined := strings.Join(args, " ")
	expectedFlags := []string{
		"--window-x=100",
		"--window-y=200",
		"--window-width=1280",
		"--window-height=720",
	}

	for _, flag := range expectedFlags {
		if !strings.Contains(joined, flag) {
			t.Errorf("expected geometry flag %q to be present, got %q", flag, joined)
		}
	}
}

func TestConflictNoWindowOverridesWindowArgs(t *testing.T) {
	opts := domain.ScrcpyOptions{
		DeviceSerial:       "DEV_WINDOW",
		NoWindow:           true,
		AlwaysOnTop:        true,
		Fullscreen:         true,
		WindowBorderless:   true,
		DisableScreensaver: true,
		WindowTitle:        "Should Be Ignored",
		WindowWidth:        "1280",
		WindowHeight:       "720",
	}

	args := BuildScrcpyArgs(opts)
	joined := strings.Join(args, " ")

	if !strings.Contains(joined, "--no-window") {
		t.Errorf("expected --no-window to be present, got %q", joined)
	}

	prohibitedFlags := []string{
		"--always-on-top",
		"--fullscreen",
		"--window-borderless",
		"--window-title",
		"--window-width",
		"--window-height",
	}

	for _, flag := range prohibitedFlags {
		if strings.Contains(joined, flag) {
			t.Errorf("flag %q should be suppressed when --no-window is active, got %q", flag, joined)
		}
	}
}

func TestConflictFullscreenOverridesGeometry(t *testing.T) {
	opts := domain.ScrcpyOptions{
		DeviceSerial: "DEV_WINDOW",
		Fullscreen:   true,
		WindowWidth:  "1280",
		WindowHeight: "720",
	}

	args := BuildScrcpyArgs(opts)
	joined := strings.Join(args, " ")

	if !strings.Contains(joined, "--fullscreen") {
		t.Errorf("expected --fullscreen to be present, got %q", joined)
	}
	if strings.Contains(joined, "--window-width") || strings.Contains(joined, "--window-height") {
		t.Errorf("window size flags should be suppressed when fullscreen is active, got %q", joined)
	}
}

func TestBuildScrcpyArgsWindowAdvanced(t *testing.T) {
	opts := domain.ScrcpyOptions{
		DeviceSerial:            "DEV_WINDOW",
		NoWindowAspectRatioLock: true,
		BackgroundColor:         "#123456",
		RenderFit:               "stretched",
	}

	args := BuildScrcpyArgs(opts)
	joined := strings.Join(args, " ")
	expectedFlags := []string{
		"--no-window-aspect-ratio-lock",
		"--background-color=#123456",
		"--render-fit=stretched",
	}

	for _, flag := range expectedFlags {
		if !strings.Contains(joined, flag) {
			t.Errorf("expected advanced window flag %q, got %q", flag, joined)
		}
	}
}
