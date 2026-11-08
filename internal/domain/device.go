package domain

// ConnectionType represents the device transport method.
type ConnectionType string

const (
	ConnUSB   ConnectionType = "USB"
	ConnTCPIP ConnectionType = "TCP/IP"
)

// Device represents an attached Android device.
type Device struct {
	Serial         string
	State          string
	ConnectionType ConnectionType
}

// ConnectOptions defines parameters for direct TCP/IP connection.
type ConnectOptions struct {
	IP             string
	Port           int
	ForceReconnect bool
}

// PairOptions defines parameters for Android 11+ wireless pairing.
type PairOptions struct {
	IP   string
	Port int
	Code string
}

// ScrcpyOptions defines user parameters to launch scrcpy.
type ScrcpyOptions struct {
	DeviceSerial string
	MaxSize      string
	BitRate      string
	VideoOn      bool
	ScreenOn     bool
	StayAwake    bool
	RecordVideo       bool
	VideoFormat       string
	RecordAudio       bool
	AudioFormat       string
	RecordPath        string
	RecordFormat      string
	RecordOrientation string
	TimeLimit         string
	SelectUSB    bool
	SelectTCPIP  bool
	TCPIPAddress string
	MaxFPS       string
	VideoCodec   string
	Orientation  string
	VideoBuffer  string
	Crop            string
	NoPlayback      bool
	AudioOn         bool
	AudioSource     string
	AudioDup        bool
	AudioCodec      string
	AudioBitRate    string
	AudioBuffer     string
	NoAudioPlayback      bool
	NoControl            bool
	KeyboardMode         string
	MouseMode            string
	NoMouseHover         bool
	NoClipboardAutosync  bool
	LegacyPaste          bool
	PushTarget           string
	V4L2Sink             string
	V4L2Buffer           string
	VideoSource          string
	CameraFacing         string
	CameraSize           string
	CameraFPS            string
	PreferText           bool
	RawKeyEvents         bool
	NoKeyRepeat          bool
	ShortcutMod          string
	OTGMode              bool
	ControlOnly          bool
	GamepadMode          string
	ShowTouches          bool
	PowerOffOnClose      bool
	NoPowerOn            bool
	ScreenOffTimeout     string
	StartApp             string
	KeepActive              bool
	AlwaysOnTop             bool
	Fullscreen              bool
	WindowBorderless        bool
	WindowTitle             string
	WindowX                 string
	WindowY                 string
	WindowWidth             string
	WindowHeight            string
	DisableScreensaver      bool
	NoWindow                bool
	NoWindowAspectRatioLock bool
	BackgroundColor         string
	RenderFit               string
	DisplayID               string
	NewDisplay              bool
	NewDisplayResolution    string
	NoVDSystemDecorations   bool
	NoVDDestroyContent      bool
	DisplayIMEPolicy        string
	FlexDisplay             bool
	TunnelHost              string
	TunnelPort              string
	ForceAdbForward         bool
	PortRange               string
	ADBServerHost           string
}

// ScrcpyCapabilities defines supported scrcpy binary capabilities for backward compatibility.
type ScrcpyCapabilities struct {
	VersionRaw        string
	Major             int
	Minor             int
	Patch             int
	HasKeepActive     bool
	HasWindowAdvanced bool
	HasFlexDisplay    bool
}

// ADBClient defines the low-level ADB/scrcpy integration interface.
type ADBClient interface {
	ListDevices() ([]Device, error)
	GetManufacturer(serial string) (string, error)
	GetModel(serial string) (string, error)
	GetIP(serial string) (string, error)
	ConfigureTCPIP(serial string, port int) error
	Connect(target string) error
	Disconnect(target string) error
	Pair(ip string, port int, code string) error
	ExecuteScrcpy(args []string) error
	ExecuteScrcpyWithEnv(args []string, extraEnv []string) error
	GetScrcpyCapabilities() (ScrcpyCapabilities, error)
}

// DeviceService defines business use-cases for device management and mirroring.
type DeviceService interface {
	ListDevices() ([]Device, error)
	GetDeviceInfo(serial string) (string, error)
	GetDeviceIP(serial string) (string, error)
	ConnectTCPIP(serial string) error
	DirectConnect(opts ConnectOptions) error
	PairDevice(opts PairOptions) error
	DisconnectDevice(target string) error
	LaunchScrcpy(opts ScrcpyOptions) error
	GetScrcpyCapabilities() (ScrcpyCapabilities, error)
}
