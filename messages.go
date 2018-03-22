package iec62056

type DataMessage struct {
	DataSets []DataSet
}

type DataSet struct {
}

type AcknowledgeMessage struct {
}

type AcknowledgeMode byte

const (
	DataReadOut = AcknowledgeMode(byte('0'))
	Programming = AcknowledgeMode(byte('1'))
	Binary      = AcknowledgeMode(byte('2'))
	Reserved    = AcknowledgeMode(byte('3'))
	Manufacture = AcknowledgeMode(byte('6'))
	IllegalMode = AcknowledgeMode(byte(' '))
)

// AcknowledgeModeFromByte
func AcknowledgeModeFromByte(a byte) AcknowledgeMode {
	switch a {
	case 0, 1, 2:
		return AcknowledgeMode(a)
	case 3, 4, 5:
		return Reserved
	}

	switch {
	case 6 <= a && a <= 9:
	case 'A' <= a && a <= 'Z':
		return Manufacture
	}

	return IllegalMode
}
