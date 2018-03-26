package iec62056

import (
	"testing"
)

func TestNewPort(t *testing.T) {
	var settings = newDefaulSettings()
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
	p := New(newDefaulSettings())

	err := p.Open("/dev/ttyUSB0")
	if err != nil {
		t.Fatalf("Error opening port: %s", err.Error())
	}
	defer p.Close()
}
