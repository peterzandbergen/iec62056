package telegram

const verbose = false

type RequestMessage struct {
	deviceAddress string
}

// AcknowledgeMessage type needs documentation. TODO:
type AcknowledgeMessage struct {
	pcc ProtocolControlCharacter
	// baudrate
	modeCondtrol AcknowledgeMode
}
