package telegram

import (
	"bufio"
	"bytes"
	"io"
	"testing"
)

const dataMessage1 = ``
const dataLine1 = ``
const dataSetEmpty = ``
const dataSetNoValueNoUnit = `1.1.1.1()`
const dataSetFailingRear = `1.1.1.1(`
const dataSetValueNoUnit = `1.1.1.1(12)`
const dataSetValueUnit = `1.1.1.1(12*kWh)`
const dataSetValueUnitCRLF = `1.1.1.1(12*kWh)` + "\r\n"

const validDataLine = `1.1.1.1(12*kWh)` + `1.1.1.1(12*kWh)` + "\r\n"

const validDataBlock = validDataLine + validDataLine + string(EndChar)

const testBcc1 = 83 ^ 0
const testBcc2 = testBcc1 ^ 0xff

// func TestParseDataMessage1(t *testing.T) {
// 	m, err := ParseDataMessage(bytes.NewBufferString(dataMessage1))
// 	if err != nil {
// 		t.Errorf("Error parsing data message: %s", err.Error())
// 	}
// 	_ = m
// }

func TestParseDataSetNoValueNoUnit(t *testing.T) {
	var bcc Bcc
	bcc = Bcc(0)
	// var r io.Reader = bytes.NewBufferString(dataSetNoValueNoUnit)
	m, err := ParseDataSet(bufio.NewReader(bytes.NewBufferString(dataSetNoValueNoUnit)), &bcc)
	if err != nil {
		t.Fatalf("Error parsing data set: %s", err.Error())
	}
	if m.address != "1.1.1.1" {
		t.Errorf("Error parsing address, expected %s, received %s", "1.1.1.1", m.address)
	}
	if m.value != "" {
		t.Errorf("Error parsing value, expected \"%s\", received \"%s\"", "", m.value)
	}
	if m.unit != "" {
		t.Errorf("Error parsing unit, expected \"%s\", received \"%s\"", "", m.unit)
	}
}

func TestParseDataSetValueNoUnit(t *testing.T) {
	var bcc Bcc
	bcc = Bcc(0)
	m, err := ParseDataSet(bufio.NewReader(bytes.NewBufferString(dataSetValueNoUnit)), &bcc)
	if err != nil {
		t.Fatalf("Error parsing data set: %s", err.Error())
	}
	if m.address != "1.1.1.1" {
		t.Errorf("Error parsing address, expected %s, received %s", "1.1.1.1", m.address)
	}
	if m.value != "12" {
		t.Errorf("Error parsing value, expected \"%s\", received \"%s\"", "12", m.value)
	}
	if m.unit != "" {
		t.Errorf("Error parsing unit, expected \"%s\", received \"%s\"", "", m.unit)
	}
}

func TestParseDataSetValueUnit(t *testing.T) {
	var bcc Bcc
	bcc = Bcc(0)

	var r io.Reader = bytes.NewBufferString(dataSetValueUnit)
	m, err := ParseDataSet(bufio.NewReader(r), &bcc)
	if err != nil {
		t.Errorf("Error parsing data set: %s", err.Error())
	}
	if m.address != "1.1.1.1" {
		t.Errorf("Error parsing address, expected %s, received %s", "1.1.1.1", m.address)
	}
	if m.value != "12" {
		t.Errorf("Error parsing value, expected \"%s\", received \"%s\"", "12", m.value)
	}
	if m.unit != "kWh" {
		t.Errorf("Error parsing unit, expected \"%s\", received \"%s\"", "kWh", m.unit)
	}
}

func TestParseDataSetValueUnitCRLF(t *testing.T) {
	var bcc Bcc
	bcc = Bcc(0)

	var r = bufio.NewReader(bytes.NewBufferString(dataSetValueUnitCRLF))
	m, err := ParseDataSet(r, &bcc)
	if err != nil {
		t.Errorf("Error parsing data set: %s", err.Error())
	}
	if m.address != "1.1.1.1" {
		t.Errorf("Error parsing address, expected %s, received %s", "1.1.1.1", m.address)
	}
	if m.value != "12" {
		t.Errorf("Error parsing value, expected \"%s\", received \"%s\"", "12", m.value)
	}
	if m.unit != "kWh" {
		t.Errorf("Error parsing unit, expected \"%s\", received \"%s\"", "kWh", m.unit)
	}
	if b, err := r.ReadByte(); err == nil && b != CR {
		t.Errorf("Expected CR, received %s", rune(b))
	}
}

func TestParseDataSetEmpty(t *testing.T) {
	var bcc Bcc
	bcc = Bcc(0)

	var r io.Reader = bytes.NewBufferString(dataSetEmpty)
	m, err := ParseDataSet(bufio.NewReader(r), &bcc)
	if err == nil {
		t.Fatal("Expected an error but received none.")
	}
	if m != nil {
		t.Error("Message is not nil.")
	}
}

func TestFailingRear(t *testing.T) {
	var bcc Bcc
	bcc = Bcc(0)

	m, err := ParseDataSet(bufio.NewReader(bytes.NewBufferString(dataSetFailingRear)), &bcc)
	if err == nil {
		t.Fatal("Expected an error but received none.")
	}
	if m != nil {
		t.Error("Message is not nil.")
	}
}

func TestBcc(t *testing.T) {
	var bcc = Bcc(83)
	bcc.Digest(0)
	if bcc != Bcc(testBcc1) {
		t.Errorf("Bcc1: expected %d, received %d", testBcc1, bcc)
	}
	bcc.Digest(0xff)
	if bcc != testBcc2 {
		t.Errorf("Bcc2: expected %d, received %d", testBcc2, bcc)
	}
}

func TestDataLine(t *testing.T) {
	var bcc Bcc
	bcc = Bcc(0)

	var r = bufio.NewReader(bytes.NewBufferString(validDataLine))
	l, err := ParseDataLine(r, &bcc)
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}
	if len(l) != 2 {
		t.Errorf("Expected 2 datasets, received %d", len(l))
	}
}

func TestDataBlock(t *testing.T) {
	var bcc Bcc
	bcc = Bcc(0)

	var r = bufio.NewReader(bytes.NewBufferString(validDataBlock))
	l, err := ParseDataBlock(r, &bcc)
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}
	if len(l) != 4 {
		t.Errorf("Expected 4 datasets, received %d", len(l))
	}
}
