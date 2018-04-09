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

const validDataBlock = validDataLine +
	validDataLine +
	string(EndChar)

const validDataBlockNoEnd = validDataLine +
	validDataLine

const validDataMessageBcc = validDataBlock +
	string(EndChar) +
	string(CR) + string(LF) +
	string(EtxChar)

const validDataMessage = string(StxChar) +
	validDataBlockNoEnd +
	string(EndChar) +
	string(CR) + string(LF) +
	string(EtxChar) +
	string(Bcc(0))

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
	if m.Address != "1.1.1.1" {
		t.Errorf("Error parsing address, expected %s, received %s", "1.1.1.1", m.Address)
	}
	if m.Value != "" {
		t.Errorf("Error parsing value, expected \"%s\", received \"%s\"", "", m.Value)
	}
	if m.Unit != "" {
		t.Errorf("Error parsing unit, expected \"%s\", received \"%s\"", "", m.Unit)
	}
}

func TestParseDataSetValueNoUnit(t *testing.T) {
	var bcc Bcc
	bcc = Bcc(0)
	m, err := ParseDataSet(bufio.NewReader(bytes.NewBufferString(dataSetValueNoUnit)), &bcc)
	if err != nil {
		t.Fatalf("Error parsing data set: %s", err.Error())
	}
	if m.Address != "1.1.1.1" {
		t.Errorf("Error parsing address, expected %s, received %s", "1.1.1.1", m.Address)
	}
	if m.Value != "12" {
		t.Errorf("Error parsing value, expected \"%s\", received \"%s\"", "12", m.Value)
	}
	if m.Unit != "" {
		t.Errorf("Error parsing unit, expected \"%s\", received \"%s\"", "", m.Unit)
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
	if m.Address != "1.1.1.1" {
		t.Errorf("Error parsing address, expected %s, received %s", "1.1.1.1", m.Address)
	}
	if m.Value != "12" {
		t.Errorf("Error parsing value, expected \"%s\", received \"%s\"", "12", m.Value)
	}
	if m.Unit != "kWh" {
		t.Errorf("Error parsing unit, expected \"%s\", received \"%s\"", "kWh", m.Unit)
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
	if m.Address != "1.1.1.1" {
		t.Errorf("Error parsing address, expected %s, received %s", "1.1.1.1", m.Address)
	}
	if m.Value != "12" {
		t.Errorf("Error parsing value, expected \"%s\", received \"%s\"", "12", m.Value)
	}
	if m.Unit != "kWh" {
		t.Errorf("Error parsing unit, expected \"%s\", received \"%s\"", "kWh", m.Unit)
	}
	if b, err := r.ReadByte(); err == nil && b != CR {
		t.Errorf("Expected CR, received %s", string(rune(b)))
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

const testBcc1 = 83 ^ 0
const testBcc2 = testBcc1 ^ 0xff

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

	bcc = 12
	bcc.Digest(23)
	bcc.Digest(23)
	if bcc != 12 {
		t.Errorf("Bcc not equal")
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
	if len(*l) != 4 {
		t.Errorf("Expected 4 datasets, received %d", len(*l))
	}
}

func TestDataBlockNoEnd(t *testing.T) {
	var bcc Bcc
	bcc = Bcc(0)

	var r = bufio.NewReader(bytes.NewBufferString(validDataBlockNoEnd))
	l, err := ParseDataBlock(r, &bcc)
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}
	if len(*l) != 4 {
		t.Errorf("Expected 4 datasets, received %d", len(*l))
	}
}

func TestDataMessage(t *testing.T) {
	r := bufio.NewReader(bytes.NewBufferString(validDataMessage))
	dm, err := ParseDataMessage(r)
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}
	_ = dm
	if len(*(dm.DataSets)) != 4 {
		t.Errorf("Expected 4, received %d", len(*(dm.DataSets)))
	}
	var bcc Bcc
	bcc.Digest([]byte(validDataMessageBcc)...)
	if bcc != dm.bcc {
		t.Errorf("Bcc: Expected %d, received %d", bcc, dm.bcc)
	}
	if len(*(dm.DataSets)) != 4 {
		t.Errorf("Expected 4, received %d", len(*(dm.DataSets)))
	}
}

const dataSetValue32 = `1.1.1.1(12345678901234567890123456789012)`

func TestParseDataSetValue32(t *testing.T) {
	var bcc Bcc
	bcc = Bcc(0)
	// var r io.Reader = bytes.NewBufferString(dataSetNoValueNoUnit)
	m, err := ParseDataSet(bufio.NewReader(bytes.NewBufferString(dataSetValue32)), &bcc)
	if err != nil {
		t.Fatalf("Error parsing data set: %s", err.Error())
	}
	if m.Address != "1.1.1.1" {
		t.Errorf("Error parsing address, expected %s, received %s", "1.1.1.1", m.Address)
	}
	if m.Value != "12345678901234567890123456789012" {
		t.Errorf("Error parsing value, expected \"%s\", received \"%s\"", "", m.Value)
	}
	if m.Unit != "" {
		t.Errorf("Error parsing unit, expected \"%s\", received \"%s\"", "", m.Unit)
	}
}

const dataSetNoValue33 = `1.1.1.1(123456789012345678901234567890123)`

func TestParseDataSetValue33(t *testing.T) {
	var bcc Bcc
	bcc = Bcc(0)
	// var r io.Reader = bytes.NewBufferString(dataSetNoValueNoUnit)
	_, err := ParseDataSet(bufio.NewReader(bytes.NewBufferString(dataSetNoValue33)), &bcc)
	if err == nil {
		t.Fatal("Exptected error.")
	}
	if err != ErrValueTooLong {
		t.Fatalf("Exptected %s, received %s.", ErrValueTooLong.Error(), err.Error())
	}
}

const dataSetUnit16 = `1.1.1.1(12345678901234567890123456789012*1234567890123456)`

func TestParseDataSetUnit16(t *testing.T) {
	var bcc Bcc
	bcc = Bcc(0)
	// var r io.Reader = bytes.NewBufferString(dataSetNoValueNoUnit)
	m, err := ParseDataSet(bufio.NewReader(bytes.NewBufferString(dataSetUnit16)), &bcc)
	if err != nil {
		t.Fatalf("Error parsing data set: %s", err.Error())
	}
	if m.Address != "1.1.1.1" {
		t.Errorf("Error parsing address, expected %s, received %s", "1.1.1.1", m.Address)
	}
	if m.Value != "12345678901234567890123456789012" {
		t.Errorf("Error parsing value, expected \"%s\", received \"%s\"", "12345678901234567890123456789012", m.Value)
	}
	if m.Unit != "1234567890123456" {
		t.Errorf("Error parsing unit, expected \"%s\", received \"%s\"", "1234567890123456", m.Unit)
	}
}

const dataSetUnit17 = `1.1.1.1(12345678901234567890123456789012*12345678901234567)`

func TestParseDataSetUnit17(t *testing.T) {
	var bcc Bcc
	bcc = Bcc(0)
	// var r io.Reader = bytes.NewBufferString(dataSetNoValueNoUnit)
	_, err := ParseDataSet(bufio.NewReader(bytes.NewBufferString(dataSetUnit17)), &bcc)
	if err == nil {
		t.Fatal("Exptected error.")
	}
	if err != ErrUnitTooLong {
		t.Fatalf("Exptected %s, received %s.", ErrUnitTooLong.Error(), err.Error())
	}
}

const identicationMessage = string(StartChar) +
	"MAN" +
	"A" +
	"identification" +
	string(CR) + string(LF)

func TestParstIdenticationMessage(t *testing.T) {
	im, err := ParseIdentificationMessage(bufio.NewReader(bytes.NewBufferString(identicationMessage)))
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}
	if im.BaudID != 'A' {
		t.Errorf("baudrateID, expected %c, received %c", 'A', im.BaudID)
	}
	if im.ManID != "MAN" {
		t.Errorf("mID, expected %s, received %s", "MAN", im.ManID)
	}
	if im.Identification != "identification" {
		t.Errorf("identification, expected %s, received %s", "idenfication", im.Identification)
	}
}

const identicationMessageBackslashW = string(StartChar) +
	"MAN" +
	"A" +
	"\\W" +
	"identification" +
	string(CR) + string(LF)

func TestParstIdenticationMessageBackslashW(t *testing.T) {
	im, err := ParseIdentificationMessage(bufio.NewReader(bytes.NewBufferString(identicationMessageBackslashW)))
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}
	if im.BaudID != 'A' {
		t.Errorf("baudrateID, expected %c, received %c", 'A', im.BaudID)
	}
	if im.ManID != "MAN" {
		t.Errorf("mID, expected %s, received %s", "MAN", im.ManID)
	}
	if im.Identification != "identification" {
		t.Errorf("identification, expected %s, received %s", "idenfication", im.Identification)
	}
}

const identicationMessageBad = string(StartChar) +
	"asdasdMAN" +
	"A" +
	"identification" +
	string(CR) + string(LF)

func TestParstIdenticationMessageBad(t *testing.T) {
	_, err := ParseIdentificationMessage(bufio.NewReader(bytes.NewBufferString(identicationMessageBad)))
	if err == nil {
		t.Fatal("Expected error.")
	}
}

const identicationMessageTooLong = string(StartChar) +
	"MAN" +
	"A" +
	"12345678901234567" +
	string(CR) + string(LF)

func TestParstIdenticationMessageTooLong(t *testing.T) {
	_, err := ParseIdentificationMessage(bufio.NewReader(bytes.NewBufferString(identicationMessageTooLong)))
	if err == nil {
		t.Fatal("Expected error.")
	}
}

func TestRequestMessage(t *testing.T) {
	b := &bytes.Buffer{}
	SerializeRequestMessage(b, RequestMessage{})
	if b.String() != "/?!\r\n" {
		t.Fatalf("bad request message: b.String()")
	}
}
