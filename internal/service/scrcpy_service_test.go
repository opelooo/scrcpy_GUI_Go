package service

import (
	"errors"
	"strings"
	"testing"

	"github.com/opelooo/scrcpy_GUI_Go/internal/domain"
)

type mockADBClient struct {
	devices       []domain.Device
	manufacturer  string
	model         string
	ip            string
	executedArgs  []string
	tcpipPort     int
	connectTarget string
	disconnectTgt string
	pairTarget    string
	pairCode      string
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

func (m *mockADBClient) Disconnect(target string) error {
	m.disconnectTgt = target
	return nil
}

func (m *mockADBClient) Pair(ip string, port int, code string) error {
	m.pairTarget = ip
	m.pairCode = code
	return nil
}

func (m *mockADBClient) ExecuteScrcpy(args []string) error {
	m.executedArgs = args
	return nil
}

func TestServiceListDevices(t *testing.T) {
	mock := &mockADBClient{
		devices: []domain.Device{
			{Serial: "DEV1", State: "device", ConnectionType: domain.ConnUSB},
			{Serial: "DEV2", State: "offline", ConnectionType: domain.ConnUSB},
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

func TestServiceDirectConnect(t *testing.T) {
	mock := &mockADBClient{}
	svc := NewScrcpyService(mock)

	// Standard connect
	err := svc.DirectConnect(domain.ConnectOptions{
		IP:             "192.168.1.200",
		Port:           5555,
		ForceReconnect: false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mock.connectTarget != "192.168.1.200:5555" {
		t.Errorf("expected target 192.168.1.200:5555, got %s", mock.connectTarget)
	}
	joined := strings.Join(mock.executedArgs, " ")
	if !strings.Contains(joined, "--tcpip=192.168.1.200:5555") {
		t.Errorf("expected --tcpip arg, got %s", joined)
	}

	// Force reconnect (+ prefix)
	err = svc.DirectConnect(domain.ConnectOptions{
		IP:             "192.168.1.200",
		Port:           5555,
		ForceReconnect: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	joinedForce := strings.Join(mock.executedArgs, " ")
	if !strings.Contains(joinedForce, "--tcpip=+192.168.1.200:5555") {
		t.Errorf("expected force reconnect --tcpip=+ arg, got %s", joinedForce)
	}
}

func TestServicePairAndDisconnect(t *testing.T) {
	mock := &mockADBClient{}
	svc := NewScrcpyService(mock)

	// Pair
	err := svc.PairDevice(domain.PairOptions{
		IP:   "192.168.1.100",
		Port: 42123,
		Code: "123456",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mock.pairTarget != "192.168.1.100" || mock.pairCode != "123456" {
		t.Errorf("pair failed to forward parameters")
	}

	// Disconnect
	err = svc.DisconnectDevice("192.168.1.100:5555")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mock.disconnectTgt != "192.168.1.100:5555" {
		t.Errorf("expected disconnect target 192.168.1.100:5555, got %s", mock.disconnectTgt)
	}
}
