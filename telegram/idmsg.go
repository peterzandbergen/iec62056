package telegram

import "fmt"

// IdentifcationMessage type is the message from the meter in response to the read command.
type IdentifcationMessage struct {
	mID            string
	baudID         byte
	identification string
}

func (i *IdentifcationMessage) String() string {
	return fmt.Sprintf("mID: %s, baudID: %c, identification: %s", i.mID, i.baudID, i.identification)
}
