package actors

import (
	"log"

	"github.com/peterzandbergen/iec62056/model"
)

// IecMessageHandler is to handle storing the new measurement.
type IecMessageHandler struct {
	Measurement *model.Measurement
	Repo        model.MeasurementRepo
}

// Do performst the actor task.
func (h *IecMessageHandler) Do() {
	err := h.Repo.Put(h.Measurement)
	if err != nil {
		// Log error.
	}
	log.Printf("Stored measurement from %s", h.Measurement.Identification)
}
