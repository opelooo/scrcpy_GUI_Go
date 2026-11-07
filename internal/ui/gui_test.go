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

func (d *dummyService) LaunchScrcpy(opts domain.ScrcpyOptions) error {
	return nil
}

func TestGUIInitialization(t *testing.T) {
	app := test.NewApp()
	w := app.NewWindow("Test")
	svc := &dummyService{
		devices: []domain.Device{
			{Serial: "TEST_DEVICE_01", State: "device"},
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
}
