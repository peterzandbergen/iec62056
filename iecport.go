package iec62056

import (
	"errors"

	"go.bug.st/serial.v1"
)

var (
	// ErrPortOpenFailed is the error returned when the serial port cannot be opened.
	ErrPortOpenFailed = errors.New("Could not open the serial port")
)

// PortSettings contains the settings for opening a new port.
type PortSettings struct {
	BaudRateChangeDelay    int
	InitialBaudRateModeABC int
	InitialBaudRateModeD   int
	Timeout                int
	Verbose                bool
	// Name of the port.
	PortName string
}

// Port that can send and receive messages to the energy meter.
type Port struct {
	// Settings.
	BaudRateChangeDelay    int
	InitialBaudRateModeABC int
	InitialBaudRateModeD   int
	Timeout                int
	Verbose                bool

	// Serial port
	port serial.Port
}

// newDefaulSettings returns portsettings with default settings.
func newDefaulSettings() *PortSettings {
	return &PortSettings{
		BaudRateChangeDelay:    0,
		InitialBaudRateModeABC: 300,
		InitialBaudRateModeD:   2400,
		Timeout:                5000,
		Verbose:                false,
	}
}

// New creates a new port. If settings is nil, the uses the default settings.
func New(settings *PortSettings) (*Port, error) {
	if settings == nil {
		settings = newDefaulSettings()
	}
	return &Port{
		BaudRateChangeDelay:    settings.BaudRateChangeDelay,
		InitialBaudRateModeABC: settings.InitialBaudRateModeABC,
		InitialBaudRateModeD:   settings.InitialBaudRateModeD,
		Timeout:                settings.Timeout,
		Verbose:                settings.Verbose,
	}, nil
}

// Open the serial port using the settings.
// Each character consists of one start bit ( binary = 0 ), 7 data bits, normally one even parity bit and one stop bit ( binary = 1 )
func (p *Port) Open(portName string) error {
	mode := &serial.Mode{
		BaudRate: p.InitialBaudRateModeABC,
		DataBits: 7,
		Parity:   serial.EvenParity,
		StopBits: serial.OneStopBit,
	}
	var err error
	p.port, err = serial.Open(portName, mode)
	if err != nil {
		p.port = nil
	}
	return err
}

func (p *Port) Read() (*DataMessage, error) {
	// Set the correct baudrate.

	return nil, ErrPortOpenFailed
}
