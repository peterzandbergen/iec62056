package telegram

import (
	"bufio"
	"errors"
	"io"
)

type DataSet struct {
	address string
	value   string
	unit    string
}

type Bcc byte

func (bcc *Bcc) Digest(b byte) {
	*bcc = (*bcc) ^ Bcc(b)
}

// IdentifcationMessage type is the message from the meter in response to the read command.
type IdentifcationMessage struct{}

// ProtocolControlCharacter
type ProtocolControlCharacter byte

const (
	// ProtocolNormal
	ProtControlNormal = ProtocolControlCharacter(byte('0'))
	// ProtControlSecondary
	ProtControlSecondary = ProtocolControlCharacter(byte('1'))
)

// AcknowledgeMessage type needs documentation. TODO:
type AcknowledgeMessage struct {
	ProtocolControlCharacter
}

type AcknowledgeMode byte

const (
	AckModeDataReadOut = AcknowledgeMode(byte('0'))
	AckModeProgramming = AcknowledgeMode(byte('1'))
	AckModeBinary      = AcknowledgeMode(byte('2'))
	AckModeReserved    = AcknowledgeMode(byte('3'))
	AckModeManufacture = AcknowledgeMode(byte('6'))
	AckModeIllegalMode = AcknowledgeMode(byte(' '))
)

const (
	CR                = byte(13)
	LF                = byte(10)
	FrontBoundaryChar = byte('(')
	RearBoundaryChar  = byte(')')
	UnitSeparator     = byte('*')
	StartChar         = byte('/')
	EndChar           = byte('!')
)

func ValidAddresChar(b byte) bool {
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
	return ValidAddresChar(b)
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
	ErrCRFound           = errors.New("End CR found")
	ErrNotImplementedYet = errors.New("Not implemented yet.")
	ErrFormatError       = errors.New("Format error.")
	ErrFormatNoChars     = errors.New("No chars found.")
	ErrEmptyDataLine     = errors.New("Empty data line found.")
	ErrUnexpectedEOF     = errors.New("Unexpected end of file.")
	ErrNoBlockEndChar    = errors.New("No block end char found.")
)

func ParseDataBlock(r *bufio.Reader, bcc *Bcc) ([]DataSet, error) {
	var b byte
	var err error
	var res []DataSet

	for {
		var ds []DataSet
		ds, err = ParseDataLine(r, bcc)
		if err != nil {
			return nil, err
		}
		if len(ds) <= 0 {
			return nil, ErrEmptyDataLine
		}
		res = append(res, ds...)
		b, err = r.ReadByte()
		if err != nil {
			return nil, ErrUnexpectedEOF
		}
		if b == EndChar {
			r.UnreadByte()
			return res, nil
		}
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

	for {
		ds, err = ParseDataSet(r, bcc)
		if err != nil {
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
			} else {
				return nil, ErrFormatError
			}
		} else {
			r.UnreadByte()
		}
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
			if !ValidAddresChar(b) {
				return nil, ErrFormatError
			}
			v = append(v, b)
		}
	}
	// Address read.
	res.address = string(v)
	v = v[:0]

	// Scan for value till * or )
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
		}
	}
	res.value = string(v)
	if b == RearBoundaryChar {
		res.unit = ""
		return res, nil
	}
	v = v[:0]

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
		}
	}
	res.unit = string(v)
	return res, nil
}

func ParseIdentificationMessage(r io.Reader) (*IdentifcationMessage, error) {
	return nil, ErrNotImplementedYet
}
