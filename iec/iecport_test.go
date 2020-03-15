package iec

import (
	"bufio"
	"bytes"
	"testing"
	"time"

	"github.com/peterzandbergen/iec62056/iec/telegram"
	"go.bug.st/serial.v1"
)

func TestNewPort(t *testing.T) {
	var settings = NewDefaultSettings()
	port := New(settings)
	if port.BaudRateChangeDelay != 0 {
		t.Errorf("port.BaudRateChangeDelay, expected %d, received %d", 0, port.BaudRateChangeDelay)
	}
	if port.InitialBaudRateModeABC != 300 {
		t.Errorf("port.InitialBaudRateModeABC, expected %d, received %d", 300, port.InitialBaudRateModeABC)
	}
	if port.InitialBaudRateModeD != 2400 {
		t.Errorf("port.InitialBaudRateModeD, expected %d, received %d", 2400, port.InitialBaudRateModeD)
	}
	if port.Timeout != 5000 {
		t.Errorf("port.Timeout, expected %d, received %d", 5000, port.Timeout)
	}
	if port.Verbose != false {
		t.Errorf("port.Verbose, expected %t, received %t", false, port.Verbose)
	}
}

func TestPortOpen(t *testing.T) {
	p := New(NewDefaultSettings())

	err := p.Open("/dev/ttyUSB0")
	if err != nil {
		t.Fatalf("Error opening port: %s", err.Error())
	}
	defer p.Close()
}

func TestReadIdentificationMessage(t *testing.T) {
	p := New(NewDefaultSettings())

	err := p.Open("/dev/ttyUSB0")
	if err != nil {
		t.Fatalf("Error opening port: %s", err.Error())
	}
	defer p.Close()

	// Set the baudrate to 300
	p.mode.BaudRate = p.InitialBaudRateModeABC
	p.port.SetMode(p.mode)

	time.Sleep(time.Second)

	// Send a request command.
	_, err = telegram.SerializeRequestMessage(p.port, telegram.RequestMessage{})
	if err != nil {
		t.Fatalf("error sending request message: %s", err.Error())
	}

	// time.Sleep(time.Second)

	// Wait for the Data Message.
	dm, err := telegram.ParseIdentificationMessage(p.r)
	if err != nil {
		t.Fatalf("error receiving idenfication message: %s", err.Error())
	}
	t.Logf("Identicatin message: %s", dm.String())
}

func TestReadDataMessage(t *testing.T) {
	p := New(NewDefaultSettings())

	err := p.Open("/dev/ttyUSB0")
	if err != nil {
		t.Fatalf("Error opening port: %s", err.Error())
	}
	defer p.Close()

	// Set the baudrate to 300
	p.mode.BaudRate = p.InitialBaudRateModeABC
	p.port.SetMode(p.mode)

	time.Sleep(time.Second)

	// Send a request command.
	_, err = telegram.SerializeRequestMessage(p.port, telegram.RequestMessage{})
	if err != nil {
		t.Fatalf("error sending request message: %s", err.Error())
	}

	time.Sleep(time.Second)

	// Wait for the Data Message.
	dm, err := telegram.ParseDataMessage(p.r)
	if err != nil {
		t.Fatalf("error receiving idenfication message: %s", err.Error())
	}
	t.Logf("Identicatin message: %s", dm.String())
}

func TestReadResponse(t *testing.T) {
	p := New(NewDefaultSettings())

	err := p.Open("/dev/ttyUSB0")
	if err != nil {
		t.Fatalf("Error opening port: %s", err.Error())
	}
	defer p.Close()

	// Set the baudrate to 300
	p.mode.BaudRate = p.InitialBaudRateModeABC
	p.port.SetMode(p.mode)

	time.Sleep(time.Second)

	// Send a request command.
	_, err = telegram.SerializeRequestMessage(p.port, telegram.RequestMessage{})
	if err != nil {
		t.Fatalf("error sending request message: %s", err.Error())
	}

	b, err := p.r.ReadByte()
	if err != nil {
		t.Fatalf("Error reading from port: %s", err.Error())
	}
	t.Logf("Bytes: '%s'", string(rune(b)))

	b, err = p.r.ReadByte()
	if err != nil {
		t.Fatalf("Error reading from port: %s", err.Error())
	}
	t.Logf("Bytes: '%s'", string(rune(b)))

	b, err = p.r.ReadByte()
	if err != nil {
		t.Fatalf("Error reading from port: %s", err.Error())
	}
	t.Logf("Bytes: '%s'", string(rune(b)))

	b, err = p.r.ReadByte()
	if err != nil {
		t.Fatalf("Error reading from port: %s", err.Error())
	}
	t.Logf("Bytes: '%s'", string(rune(b)))

	b, err = p.r.ReadByte()
	if err != nil {
		t.Fatalf("Error reading from port: %s", err.Error())
	}
	t.Logf("Bytes: '%s'", string(rune(b)))

	b, err = p.r.ReadByte()
	if err != nil {
		t.Fatalf("Error reading from port: %s", err.Error())
	}
	t.Logf("Bytes: '%s'", string(rune(b)))

	b, err = p.r.ReadByte()
	if err != nil {
		t.Fatalf("Error reading from port: %s", err.Error())
	}
	t.Logf("Bytes: '%s'", string(rune(b)))

	b, err = p.r.ReadByte()
	if err != nil {
		t.Fatalf("Error reading from port: %s", err.Error())
	}
	t.Logf("Bytes: '%s'", string(rune(b)))

}

func TestRawPort(t *testing.T) {
	p, err := serial.Open("/dev/ttyUSB0",
		&serial.Mode{
			BaudRate: 300,
			DataBits: 7,
			Parity:   serial.EvenParity,
			StopBits: serial.OneStopBit,
		})
	if err != nil {
		t.Fatalf("error opening port: %s", err.Error())
	}
	defer p.Close()
	p.ResetInputBuffer()
	p.ResetOutputBuffer()

	p.Write([]byte("/?!\r\n"))

	time.Sleep(time.Second)
	var buf = make([]byte, 1000)
	n, err := p.Read(buf)
	if err != nil {
		t.Fatalf("error opening port: %s", err.Error())
	}
	t.Logf("response: %s", string(buf[:n]))
}

func TestRawPortSerReq(t *testing.T) {
	p, err := serial.Open("/dev/ttyUSB0",
		&serial.Mode{
			BaudRate: 300,
			DataBits: 7,
			Parity:   serial.EvenParity,
			StopBits: serial.OneStopBit,
		})
	if err != nil {
		t.Fatalf("error opening port: %s", err.Error())
	}
	defer p.Close()
	p.ResetInputBuffer()
	p.ResetOutputBuffer()

	// p.Write([]byte("/?!\r\n"))
	telegram.SerializeRequestMessage(p, telegram.RequestMessage{})

	time.Sleep(time.Second)
	var buf = make([]byte, 1000)
	n, err := p.Read(buf)
	if err != nil {
		t.Fatalf("error opening port: %s", err.Error())
	}
	t.Logf("response: %s", string(buf[:n]))
}

func TestRawPortSerReqIDMsg(t *testing.T) {
	p, err := serial.Open("/dev/ttyUSB0",
		&serial.Mode{
			BaudRate: 300,
			DataBits: 7,
			Parity:   serial.EvenParity,
			StopBits: serial.OneStopBit,
		})
	if err != nil {
		t.Fatalf("error opening port: %s", err.Error())
	}
	defer p.Close()
	p.ResetInputBuffer()
	p.ResetOutputBuffer()

	// p.Write([]byte("/?!\r\n"))
	telegram.SerializeRequestMessage(p, telegram.RequestMessage{})

	time.Sleep(time.Second)
	var buf = make([]byte, 1000)
	n, err := p.Read(buf)
	if err != nil {
		t.Fatalf("error opening port: %s", err.Error())
	}
	t.Logf("response: %s", string(buf[:n]))
}

const identicationMessage = string(telegram.StartChar) +
	"MAN" +
	"A" +
	"identification" +
	string(telegram.CR) + string(telegram.LF)

const immediateResponse = identicationMessage + telegram.ValidTestDataMessage

func TestRead(t *testing.T) {
	r := bufio.NewReader(bytes.NewReader([]byte(immediateResponse)))
	m, err := readImmediateResponse(r)
	if err != nil {
		t.Fatalf("error reading immediate response: %s", err.Error())
	}
	t.Logf("message: %+v", m)
}
