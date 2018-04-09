package telegram

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"

	"log"
)

// SerializeRequestMessage serializes the request message to w.
func SerializeRequestMessage(w io.Writer, rm RequestMessage) (int, error) {
	msg := string(StartChar) +
		string(RequestCommandChar) +
		rm.deviceAddress +
		string(EndChar) +
		string(CR) +
		string(LF)
	return w.Write([]byte(msg))
}

// DataMessage type captures the data message.
type DataMessage struct {
	DataSets *[]DataSet
	bcc      Bcc
}

func (d *DataMessage) String() string {
	b := &bytes.Buffer{}

	fmt.Fprintf(b, "bcc: %d %+v", d.bcc, *d.DataSets)
	return b.String()
}

// DataSet type.
type DataSet struct {
	Address string
	Value   string
	Unit    string
}

// Bcc type captures the checksum.
type Bcc byte

// Digest processes the next byte for the checksum.
func (bcc *Bcc) Digest(b ...byte) {
	for _, i := range b {
		*bcc = (*bcc) ^ Bcc(i)
	}
}

// ProtocolControlCharacter type.
type ProtocolControlCharacter byte

const (
	// ProtControlNormal value.
	ProtControlNormal = ProtocolControlCharacter(byte('0'))
	// ProtControlSecondary value.
	ProtControlSecondary = ProtocolControlCharacter(byte('1'))
)

// AcknowledgeMode type.
type AcknowledgeMode byte

// BaudrateIdentification type.
type BaudrateIdentification byte

// Baudrate returns the numeric baudrate from the code in the ID message.
func Baudrate(br BaudrateIdentification) int {
	switch rune(br) {
	case '0':
		return 300
	case 'A', '1':
		return 600
	case 'B', '2':
		return 1200
	case 'C', '3':
		return 2400
	case 'D', '4':
		return 4800
	case 'E', '5':
		return 9600
	case 'F', '6':
		return 19200
	}
	return 0
}

const (
	AckModeDataReadOut = AcknowledgeMode(byte('0'))
	AckModeProgramming = AcknowledgeMode(byte('1'))
	AckModeBinary      = AcknowledgeMode(byte('2'))
	AckModeReserved    = AcknowledgeMode(byte('3'))
	AckModeManufacture = AcknowledgeMode(byte('6'))
	AckModeIllegalMode = AcknowledgeMode(byte(' '))
)

const (
	CR                 = byte(0x0D)
	LF                 = byte(0x0A)
	FrontBoundaryChar  = byte('(')
	RearBoundaryChar   = byte(')')
	UnitSeparator      = byte('*')
	StartChar          = byte('/')
	RequestCommandChar = byte('?')
	EndChar            = byte('!')
	StxChar            = byte(0x02)
	EtxChar            = byte(0x03)
	SeqDelChar         = byte('\\')
)

// ValidTestDataMessage can be used for testing.
const ValidTestDataMessage = string(StxChar) +
	`1.1.1.1(12*kWh)` + `1.1.1.2(12*kWh)` + "\r\n" +
	string(EndChar) +
	string(CR) + string(LF) +
	string(EtxChar) +
	string(Bcc(0))

func ValidAddressChar(b byte) bool {
	switch b {
	case FrontBoundaryChar, RearBoundaryChar, StartChar, EndChar:
		return false
	default:
		return true
	}
}

func ValidValueChar(b byte) bool {
	switch b {
	case FrontBoundaryChar, UnitSeparator, RearBoundaryChar, StartChar, EndChar:
		return false
	default:
		return true
	}
}

func ValidUnitChar(b byte) bool {
	return ValidAddressChar(b)
}

// AcknowledgeModeFromByte returns the acknowledge mode from the given byte value.
func AcknowledgeModeFromByte(a byte) AcknowledgeMode {
	switch a {
	case 0, 1, 2:
		return AcknowledgeMode(a)
	case 3, 4, 5:
		return AckModeReserved
	}

	switch {
	case 6 <= a && a <= 9:
	case 'A' <= a && a <= 'Z':
		return AckModeManufacture
	}

	return AckModeIllegalMode
}

var (
	ErrCRFound               = errors.New("End CR found")
	ErrNotImplementedYet     = errors.New("not implemented yet")
	ErrFormatError           = errors.New("format error")
	ErrFormatNoChars         = errors.New("no chars found")
	ErrEmptyDataLine         = errors.New("empty data line found")
	ErrUnexpectedEOF         = errors.New("unexpected end of file")
	ErrNoBlockEndChar        = errors.New("no block end char found")
	ErrNoStartChar           = errors.New("no StartChar found")
	ErrAddressTooLong        = errors.New("field too long")
	ErrValueTooLong          = errors.New("field too long")
	ErrUnitTooLong           = errors.New("field too long")
	ErrIdentificationTooLong = errors.New("identification field too long")
)

func ParseDataMessage(r *bufio.Reader) (*DataMessage, error) {
	var b byte
	var err error
	var res *[]DataSet
	var bcc = Bcc(0)

	if verbose {
		log.Println("Starting ParseDataMessage")
	}
	// Consume all bytes till a start of message is found.
	for {
		b, err = r.ReadByte()
		if err != nil {
			return nil, ErrUnexpectedEOF
		}
		if b == StxChar {
			break
		}
	}
	if verbose {
		log.Println("Found StxChar")
	}
	// Get the datasets.
	res, err = ParseDataBlock(r, &bcc)
	if err != nil {
		return nil, err
	}
	_, err = ParseDataMessageEnd(r, &bcc)
	if err != nil {
		return nil, err
	}

	return &DataMessage{
		DataSets: res,
		bcc:      bcc,
	}, nil
}

// ParseDataMessageEnd parses the end of a datamessage.
// ! CR LF ETX BCC
func ParseDataMessageEnd(r *bufio.Reader, bcc *Bcc) (*DataMessage, error) {
	var b byte
	var err error

	if verbose {
		log.Println("Starting ParseDataMessageEnd")
	}

	b, err = r.ReadByte()
	if err != nil {
		return nil, err
	}
	if b != EndChar {
		if verbose {
			log.Printf("ParseDataMessageEnd, error parsing EndChar, found %d", b)
		}
		return nil, ErrFormatError
	}
	bcc.Digest(b)

	b, err = r.ReadByte()
	if err != nil {
		return nil, err
	}
	if b != CR {
		if verbose {
			log.Println("ParseDataMessageEnd, error parsing CR")
		}
		return nil, ErrFormatError
	}
	bcc.Digest(b)

	b, err = r.ReadByte()
	if err != nil {
		return nil, err
	}
	if b != LF {
		if verbose {
			log.Println("ParseDataMessageEnd, error parsing LF")
		}
		return nil, ErrFormatError
	}
	bcc.Digest(b)

	b, err = r.ReadByte()
	if err != nil {
		return nil, err
	}
	if b != EtxChar {
		if verbose {
			log.Println("ParseDataMessageEnd, error parsing EtxChar")
		}
		return nil, ErrFormatError
	}
	bcc.Digest(b)

	b, err = r.ReadByte()
	if err != nil {
		if verbose {
			log.Println("ParseDataMessageEnd, error parsing Bcc")
		}
		return nil, err
	}

	return &DataMessage{
		bcc: *bcc,
	}, nil
}

// ParseDataBlock parses til no valid data lines can be parsed.
func ParseDataBlock(r *bufio.Reader, bcc *Bcc) (*[]DataSet, error) {
	var err error
	var res []DataSet

	if verbose {
		log.Println("Starting ParseDataBlock")
	}

	for {
		var ds []DataSet
		ds, err = ParseDataLine(r, bcc)
		if err != nil {
			if len(res) <= 0 {
				return nil, ErrEmptyDataLine
			}
			return &res, nil
		}
		res = append(res, ds...)
	}
}

// ParseDataMessage reads bytes from r till a new complete datamessage has been read or an error occured.
// func ParseDataMessage(r io.Reader) (*DataMessage, error) {
// 	return nil, ErrNotImplementedYet
// }

// ParseDataLine parses a DataSets till a CR LF has been detected.
// Data lines consist of one or more datasets.
func ParseDataLine(r *bufio.Reader, bcc *Bcc) ([]DataSet, error) {
	var b byte
	var err error
	var ds *DataSet
	var res []DataSet

	if verbose {
		log.Println("Starting ParseDataLine")
	}

	for {
		ds, err = ParseDataSet(r, bcc)
		if err != nil {
			r.UnreadByte()
			return nil, ErrFormatError
		}
		res = append(res, *ds)
		// Test if the next two chars are CR LF
		b, err = r.ReadByte()
		if err == nil && b == CR {
			bcc.Digest(b)
			b, err = r.ReadByte()
			if err == nil && b == LF {
				bcc.Digest(b)
				return res, nil
			}
			// Error, CR not followed by LF
			return nil, ErrFormatError
		}
		r.UnreadByte()
	}
}

// ParseDataSet reads bytes from r till a new complete dataset has been read or an error occured.
// A data message contains a list of data sets. Each data set consists of 3 fields "address", "value", and "unit".
// Each of these fields is optional an may thus be equal to the empty string.
// Data set ::= Address '(' Value(optional) ('*' unit)(optional) ')'
// Ignores CR and LF and reads up to the first !
func ParseDataSet(r *bufio.Reader, bcc *Bcc) (*DataSet, error) {
	// read chars til Front boundary.
	var b byte
	var err error
	var va [100]byte
	var v = va[:0]
	res := &DataSet{}

	// Read the address till FrontBoundaryChar == (
	if verbose {
		log.Println("Starting ParseDataSet")
	}

	if verbose {
		log.Println("Scanning for Address")
	}
ScanAddress:
	for {
		b, err = r.ReadByte()
		if err != nil {
			return nil, ErrFormatNoChars
		}
		switch b {
		case CR, LF:
			r.UnreadByte()
			return nil, ErrCRFound
		case FrontBoundaryChar:
			bcc.Digest(b)
			break ScanAddress
		default:
			bcc.Digest(b)
			if !ValidAddressChar(b) {
				return nil, ErrFormatError
			}
			v = append(v, b)
			if len(v) > 16 {
				return nil, ErrAddressTooLong
			}
		}
	}
	// Address read.
	res.Address = string(v)
	v = v[:0]

	// Scan for value till * or )
	if verbose {
		log.Println("Scanning for Value")
	}
ScanValue:
	for {
		b, err = r.ReadByte()
		if err != nil {
			return nil, ErrFormatError
		}
		bcc.Digest(b)
		switch b {
		case RearBoundaryChar, UnitSeparator:
			break ScanValue
		default:
			if !ValidValueChar(b) {
				return nil, ErrFormatError
			}
			v = append(v, b)
			if len(v) > 32 {
				return nil, ErrValueTooLong
			}
		}
	}
	res.Value = string(v)
	if b == RearBoundaryChar {
		res.Unit = ""
		return res, nil
	}
	v = v[:0]

	if verbose {
		log.Println("Scanning for Unit")
	}
ScanUnit:
	for {
		b, err = r.ReadByte()
		if err != nil {
			return nil, ErrFormatError
		}
		bcc.Digest(b)
		switch b {
		case RearBoundaryChar:
			break ScanUnit
		default:
			if !ValidValueChar(b) {
				return nil, ErrFormatError
			}
			v = append(v, b)
			if len(v) > 16 {
				return nil, ErrUnitTooLong
			}
		}
	}
	res.Unit = string(v)
	return res, nil
}

func ParseIdentificationMessage(r *bufio.Reader) (*IdentifcationMessage, error) {
	var b byte
	var err error
	var res = &IdentifcationMessage{}

	// StartChar
	b, err = r.ReadByte()
	if err != nil {
		return nil, ErrFormatError
	}
	if b != StartChar {
		return nil, ErrFormatError
	}

	// Manufacturer ID
	var id [3]byte
	// mID
	for i := 0; i < len(id); i++ {
		b, err = r.ReadByte()
		if err != nil {
			return nil, ErrFormatError
		}
		id[i] = b
	}
	res.ManID = string(id[:])

	// baudrate mode
	b, err = r.ReadByte()
	if err != nil {
		return nil, ErrFormatError
	}
	if Baudrate(BaudrateIdentification(b)) == 0 {
		return nil, ErrFormatError
	}
	res.BaudID = b

	var vt [33]byte
	var v = vt[:0]

	// \W or not
	b, err = r.ReadByte()
	if err != nil {
		return nil, ErrFormatError
	}
	if b == SeqDelChar {
		// Read a W
		b, err = r.ReadByte()
		if err != nil || b != byte('W') {
			return nil, ErrFormatError
		}
	} else {
		v = append(v, b)
	}

ScanIdenfication:
	for {
		b, err = r.ReadByte()
		if err != nil {
			return nil, ErrFormatError
		}
		switch {
		case b == CR:
			break ScanIdenfication
		default:
			v = append(v, b)
			if len(v) > 16 {
				return nil, ErrIdentificationTooLong
			}
		}
	}

	res.Identification = string(v)
	// Test if the last char is LF
	b, err = r.ReadByte()
	if err != nil || b != LF {
		return nil, ErrFormatError
	}
	return res, nil
}
