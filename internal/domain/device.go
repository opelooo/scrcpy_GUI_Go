package domain

// Device represents an attached Android device.
type Device struct {
	Serial string
	State  string
}

// ScrcpyOptions defines user parameters to launch scrcpy.
type ScrcpyOptions struct {
	DeviceSerial string
	MaxSize      string
	BitRate      string
	VideoOn      bool
	ScreenOn     bool
	StayAwake    bool
	RecordVideo  bool
	VideoFormat  string
	RecordAudio  bool
	AudioFormat  string
}

// ADBClient defines the low-level ADB/scrcpy integration interface.
type ADBClient interface {
	ListDevices() ([]Device, error)
	GetManufacturer(serial string) (string, error)
	GetModel(serial string) (string, error)
	GetIP(serial string) (string, error)
	ConfigureTCPIP(serial string, port int) error
	Connect(target string) error
	ExecuteScrcpy(args []string) error
}

// DeviceService defines business use-cases for device management and mirroring.
type DeviceService interface {
	ListDevices() ([]Device, error)
	GetDeviceInfo(serial string) (string, error)
	GetDeviceIP(serial string) (string, error)
	ConnectTCPIP(serial string) error
	LaunchScrcpy(opts ScrcpyOptions) error
}
