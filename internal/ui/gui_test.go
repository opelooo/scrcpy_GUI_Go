package ui

import (
	"testing"

	"fyne.io/fyne/v2/test"

	"github.com/opelooo/scrcpy_GUI_Go/internal/domain"
)

type dummyService struct {
	devices []domain.Device
}

func (d *dummyService) ListDevices() ([]domain.Device, error) {
	return d.devices, nil
}

func (d *dummyService) GetDeviceInfo(serial string) (string, error) {
	return "Test Manufacturer", nil
}

func (d *dummyService) GetDeviceIP(serial string) (string, error) {
	return "192.168.1.55", nil
}

func (d *dummyService) ConnectTCPIP(serial string) error {
	return nil
}

func (d *dummyService) DirectConnect(opts domain.ConnectOptions) error {
	return nil
}

func (d *dummyService) PairDevice(opts domain.PairOptions) error {
	return nil
}

func (d *dummyService) DisconnectDevice(target string) error {
	return nil
}

func (d *dummyService) LaunchScrcpy(opts domain.ScrcpyOptions) error {
	return nil
}

func (d *dummyService) GetScrcpyCapabilities() (domain.ScrcpyCapabilities, error) {
	return domain.ScrcpyCapabilities{
		VersionRaw: "scrcpy 3.3.4",
		Major:      3,
		Minor:      3,
		Patch:      4,
	}, nil
}

func TestGUIInitialization(t *testing.T) {
	app := test.NewApp()
	w := app.NewWindow("Test")
	svc := &dummyService{
		devices: []domain.Device{
			{Serial: "TEST_DEVICE_01", State: "device", ConnectionType: domain.ConnUSB},
			{Serial: "192.168.1.50:5555", State: "device", ConnectionType: domain.ConnTCPIP},
		},
	}
	gui := NewGUI(app, w, svc)
	if gui == nil {
		t.Fatal("expected non-nil GUI")
	}

	if gui.selectedDevice != "TEST_DEVICE_01" {
		t.Fatalf("expected initial selected device to be TEST_DEVICE_01, got %s", gui.selectedDevice)
	}

	gui.selectDevice("ANOTHER_DEV")
	if gui.selectedDevice != "ANOTHER_DEV" {
		t.Fatalf("expected selectedDevice to be ANOTHER_DEV, got %s", gui.selectedDevice)
	}

	// Test clickable label badge formatting
	lbl := newClickableLabel(w.Canvas(), nil, nil, nil, nil, nil)
	lbl.updateDevice(domain.Device{Serial: "USB123", ConnectionType: domain.ConnUSB})
	if lbl.Text != "[USB] USB123" {
		t.Errorf("expected [USB] USB123, got %s", lbl.Text)
	}

	lbl.updateDevice(domain.Device{Serial: "192.168.1.99:5555", ConnectionType: domain.ConnTCPIP})
	if lbl.Text != "[Wi-Fi] 192.168.1.99:5555" {
		t.Errorf("expected [Wi-Fi] 192.168.1.99:5555, got %s", lbl.Text)
	}

	// Verify video engine controls initialization
	if gui.videoCodecCombo.Selected != "default" {
		t.Errorf("expected default codec, got %s", gui.videoCodecCombo.Selected)
	}
	if gui.maxFPSCombo.Selected != "Default" {
		t.Errorf("expected Default max FPS, got %s", gui.maxFPSCombo.Selected)
	}
	if gui.orientationCombo.Selected != "Default" {
		t.Errorf("expected Default orientation, got %s", gui.orientationCombo.Selected)
	}
	if gui.noPlaybackCheck.Checked {
		t.Errorf("expected noPlaybackCheck to be unchecked initially")
	}

	// Verify audio engine controls initialization
	if !gui.audioOnCheck.Checked {
		t.Errorf("expected audioOnCheck to be checked initially")
	}
	if gui.audioSourceCombo.Selected != "output" {
		t.Errorf("expected default audio source to be output, got %s", gui.audioSourceCombo.Selected)
	}
	if gui.audioCodecCombo.Selected != "opus" {
		t.Errorf("expected default audio codec to be opus, got %s", gui.audioCodecCombo.Selected)
	}
	if gui.audioBitRateCombo.Selected != "Default" {
		t.Errorf("expected default audio bitrate to be Default, got %s", gui.audioBitRateCombo.Selected)
	}
	if gui.audioDupCheck.Checked {
		t.Errorf("expected audioDupCheck to be unchecked initially")
	}
	if gui.noAudioPlaybackCheck.Checked {
		t.Errorf("expected noAudioPlaybackCheck to be unchecked initially")
	}

	// Verify control engine controls initialization
	if gui.noControlCheck.Checked {
		t.Errorf("expected noControlCheck to be unchecked initially")
	}
	if gui.otgModeCheck.Checked {
		t.Errorf("expected otgModeCheck to be unchecked initially")
	}
	if gui.controlOnlyCheck.Checked {
		t.Errorf("expected controlOnlyCheck to be unchecked initially")
	}
	if gui.keyboardModeCombo.Selected != "default" {
		t.Errorf("expected default keyboard mode, got %s", gui.keyboardModeCombo.Selected)
	}
	if gui.mouseModeCombo.Selected != "default" {
		t.Errorf("expected default mouse mode, got %s", gui.mouseModeCombo.Selected)
	}
	if gui.noMouseHoverCheck.Checked {
		t.Errorf("expected noMouseHoverCheck to be unchecked initially")
	}
	if gui.noClipboardAutosyncCheck.Checked {
		t.Errorf("expected noClipboardAutosyncCheck to be unchecked initially")
	}
	if gui.legacyPasteCheck.Checked {
		t.Errorf("expected legacyPasteCheck to be unchecked initially")
	}
	if gui.pushTargetEntry.Text != "/sdcard/Download/" {
		t.Errorf("expected default push target to be /sdcard/Download/, got %s", gui.pushTargetEntry.Text)
	}

	// Verify device features and backward compatibility
	if !gui.screenOnCheck.Checked {
		t.Errorf("expected screenOnCheck to be checked initially")
	}
	if gui.showTouchesCheck.Checked {
		t.Errorf("expected showTouchesCheck to be unchecked initially")
	}
	if gui.keepActiveCheck.Visible() {
		t.Errorf("expected keepActiveCheck to be hidden on scrcpy 3.3.4")
	}

	// Verify window widgets and backward compatibility
	if gui.alwaysOnTopCheck.Checked {
		t.Errorf("expected alwaysOnTopCheck to be unchecked initially")
	}
	if gui.fullscreenCheck.Checked {
		t.Errorf("expected fullscreenCheck to be unchecked initially")
	}
	if gui.noWindowCheck.Checked {
		t.Errorf("expected noWindowCheck to be unchecked initially")
	}
	if gui.noWindowAspectRatioLockCheck.Visible() {
		t.Errorf("expected noWindowAspectRatioLockCheck to be hidden on scrcpy 3.3.4")
	}

	// Test conflict callback: checking noWindow should uncheck fullscreen and alwaysOnTop
	gui.fullscreenCheck.SetChecked(true)
	gui.alwaysOnTopCheck.SetChecked(true)
	gui.noWindowCheck.SetChecked(true)

	if gui.fullscreenCheck.Checked {
		t.Errorf("expected fullscreenCheck to be unchecked when noWindow is checked")
	}
	if gui.alwaysOnTopCheck.Checked {
		t.Errorf("expected alwaysOnTopCheck to be unchecked when noWindow is checked")
	}

	// Verify recording controls initialization
	if gui.recordOrientationCombo.Selected != "default" {
		t.Errorf("expected default record orientation, got %s", gui.recordOrientationCombo.Selected)
	}
	if gui.recordPathEntry.Text != "" {
		t.Errorf("expected empty record path entry initially")
	}
	if gui.timeLimitEntry.Text != "" {
		t.Errorf("expected empty time limit entry initially")
	}

	// Verify display engine controls initialization
	if gui.newDisplayCheck.Checked {
		t.Errorf("expected newDisplayCheck to be unchecked initially")
	}
	if gui.noVDSystemDecorationsCheck.Checked {
		t.Errorf("expected noVDSystemDecorationsCheck to be unchecked initially")
	}
	if gui.noVDDestroyContentCheck.Checked {
		t.Errorf("expected noVDDestroyContentCheck to be unchecked initially")
	}
	if gui.displayIMEPolicyCombo.Selected != "default" {
		t.Errorf("expected default IME policy, got %s", gui.displayIMEPolicyCombo.Selected)
	}
	if gui.flexDisplayCheck.Visible() {
		t.Errorf("expected flexDisplayCheck to be hidden on scrcpy 3.3.4")
	}

	// Test conflict callback: checking newDisplay clears displayID
	gui.displayIDEntry.SetText("1")
	gui.newDisplayCheck.SetChecked(true)
	if gui.displayIDEntry.Text != "" {
		t.Errorf("expected displayIDEntry to be cleared when newDisplay is checked")
	}

	// Verify tunnel engine controls initialization
	if gui.forceAdbForwardCheck.Checked {
		t.Errorf("expected forceAdbForwardCheck to be unchecked initially")
	}
	if gui.tunnelHostEntry.Text != "" {
		t.Errorf("expected empty tunnelHostEntry initially")
	}
	if gui.tunnelPortEntry.Text != "" {
		t.Errorf("expected empty tunnelPortEntry initially")
	}
	if gui.adbServerHostEntry.Text != "" {
		t.Errorf("expected empty adbServerHostEntry initially")
	}
	if gui.portRangeEntry.Text != "" {
		t.Errorf("expected empty portRangeEntry initially")
	}
}
