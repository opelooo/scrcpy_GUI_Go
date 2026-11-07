package service

import (
	"errors"
	"testing"

	"github.com/opelooo/scrcpy_GUI_Go/internal/domain"
)

type mockADBClient struct {
	devices      []domain.Device
	manufacturer string
	model        string
	ip           string
	executedArgs []string
	tcpipPort    int
	connectTarget string
}

func (m *mockADBClient) ListDevices() ([]domain.Device, error) {
	return m.devices, nil
}

func (m *mockADBClient) GetManufacturer(serial string) (string, error) {
	if serial == "" {
		return "", errors.New("empty serial")
	}
	return m.manufacturer, nil
}

func (m *mockADBClient) GetModel(serial string) (string, error) {
	return m.model, nil
}

func (m *mockADBClient) GetIP(serial string) (string, error) {
	if m.ip == "" {
		return "", errors.New("no ip")
	}
	return m.ip, nil
}

func (m *mockADBClient) ConfigureTCPIP(serial string, port int) error {
	m.tcpipPort = port
	return nil
}

func (m *mockADBClient) Connect(target string) error {
	m.connectTarget = target
	return nil
}

func (m *mockADBClient) ExecuteScrcpy(args []string) error {
	m.executedArgs = args
	return nil
}

func TestServiceListDevices(t *testing.T) {
	mock := &mockADBClient{
		devices: []domain.Device{
			{Serial: "DEV1", State: "device"},
			{Serial: "DEV2", State: "offline"},
		},
	}
	svc := NewScrcpyService(mock)
	active, err := svc.ListDevices()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(active) != 1 || active[0].Serial != "DEV1" {
		t.Fatalf("expected 1 active device, got %v", active)
	}
}

func TestServiceGetDeviceInfo(t *testing.T) {
	mock := &mockADBClient{
		manufacturer: "Google",
		model:        "Pixel 8",
	}
	svc := NewScrcpyService(mock)
	info, err := svc.GetDeviceInfo("DEV1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "Google (Pixel 8)"
	if info != expected {
		t.Fatalf("expected %q, got %q", expected, info)
	}
}

func TestServiceConnectTCPIP(t *testing.T) {
	mock := &mockADBClient{
		ip: "192.168.1.150",
	}
	svc := NewScrcpyService(mock)
	if err := svc.ConnectTCPIP("DEV1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mock.tcpipPort != 5555 {
		t.Fatalf("expected port 5555, got %d", mock.tcpipPort)
	}
	if mock.connectTarget != "192.168.1.150:5555" {
		t.Fatalf("expected target 192.168.1.150:5555, got %s", mock.connectTarget)
	}
}
