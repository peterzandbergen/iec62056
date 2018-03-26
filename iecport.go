package iec62056

import (
	"bufio"
	"errors"

	"github.com/peterzandbergen/iec62056/telegram"

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
	// Current mode.
	mode *serial.Mode
	// Buffered
	r *bufio.Reader
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
func New(settings *PortSettings) *Port {
	if settings == nil {
		settings = newDefaulSettings()
	}
	return &Port{
		BaudRateChangeDelay:    settings.BaudRateChangeDelay,
		InitialBaudRateModeABC: settings.InitialBaudRateModeABC,
		InitialBaudRateModeD:   settings.InitialBaudRateModeD,
		Timeout:                settings.Timeout,
		Verbose:                settings.Verbose,
	}
}

// Open the serial port using the settings.
// Each character consists of one start bit ( binary = 0 ), 7 data bits, normally one even parity bit and one stop bit ( binary = 1 )
func (p *Port) Open(portName string) error {
	p.mode = &serial.Mode{
		BaudRate: p.InitialBaudRateModeABC,
		DataBits: 7,
		Parity:   serial.EvenParity,
		StopBits: serial.OneStopBit,
	}
	var err error
	p.port, err = serial.Open(portName, p.mode)
	if err != nil {
		p.port = nil
		return err
	}
	// Create buffered IO for the port.
	p.r = bufio.NewReader(p.port)
	return nil
}

func (p *Port) Close() {
	if p.port == nil {
		return
	}
	p.port.Close()
	p.port = nil
	p.r = nil
}

func (p *Port) Read() (*DataMessage, error) {
	// Set the baudrate to 300
	p.mode.BaudRate = p.InitialBaudRateModeABC
	p.port.SetMode(p.mode)

	// Send a request command.
	_, err := telegram.SerializeRequestMessage(p.port, telegram.RequestMessage{})
	if err != nil {
		return nil, err
	}

	// Wait for the Identification Message.
	_, err = telegram.ParseIdentificationMessage(p.r)
	if err != nil {
		return nil, err
	}
	// Send ack.

	return nil, ErrPortOpenFailed
}
