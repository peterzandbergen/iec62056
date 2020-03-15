package iec

import (
	"bufio"
	"errors"
	"log"

	"github.com/peterzandbergen/iec62056/iec/telegram"
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

// NewDefaulSettings returns portsettings with default settings.
func NewDefaultSettings() *PortSettings {
	return &PortSettings{
		BaudRateChangeDelay:    0,
		InitialBaudRateModeABC: 300,
		InitialBaudRateModeD:   2400,
		Timeout:                5000,
		Verbose:                false,
	}
}

// New creates a new port. If settings is nil, it uses the default settings.
func New(settings *PortSettings) *Port {
	if settings == nil {
		settings = NewDefaultSettings()
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
		log.Printf("cannot open serial port: %s", portName)
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

func readImmediateResponse(r *bufio.Reader) (*DataMessage, error) {
	// Wait for the Identification Message.
	im, err := telegram.ParseIdentificationMessage(r)
	if err != nil {
		return nil, err
	}

	// Wait for the Data.
	dm, err := telegram.ParseDataMessage(r)
	if err != nil {
		return nil, err
	}

	res := &DataMessage{
		ManufacturerID: im.ManID,
		MeterID:        im.Identification,
	}
	for _, m := range *dm.DataSets {
		var s = DataSet{
			Address: m.Address,
			Value:   m.Value,
			Unit:    m.Unit,
		}
		res.DataSets = append(res.DataSets, s)
	}
	return res, nil
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
	return readImmediateResponse(p.r)
}
