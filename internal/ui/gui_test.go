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
}
