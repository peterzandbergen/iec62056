package actors

import (
	"log"

	"github.com/peterzandbergen/iec62056/model"
)

// IecMessageHandler reads a message from the port and stores it in the Repo
//
type IecMessageHandler struct {
	LocalRepo model.MeasurementRepo
	MeterRepo model.MeasurementRepo
}

// Do performs the actor task.
// Open the serial port repo, Get a message, and store it in the local cache.
// Performs a timeout on the Get to prevent blocking.
func (h *IecMessageHandler) Do() error {
	m, err := h.MeterRepo.Get(nil)
	if err != nil {
		// Log error
		log.Printf("Error getting measuerment from reader, error: %s", err.Error())
		return err
	}
	err = h.LocalRepo.Put(m)
	if err != nil {
		// Log error.
		log.Printf("Error storing to local cache, error: %s", err.Error())
		return err
	}
	log.Printf("Stored measurement from %s", m.Identification)
	return nil
}
