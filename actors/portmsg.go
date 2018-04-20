package actors

import (
	"log"
	"time"

	"github.com/peterzandbergen/iec62056/adapters/meter"
	"github.com/peterzandbergen/iec62056/iec"
	"github.com/peterzandbergen/iec62056/model"
)

// IecMessageHandler reads a message from the port and stores it in the Repo
//
type IecMessageHandler struct {
	Retries      int
	PortSettings iec.PortSettings
	Repo         model.MeasurementRepo
	err          error
}

type measurement struct {
	m   *model.Measurement
	err error
}

// Do performs the actor task.
// Open the serial port repo, Get a message, and store it in the local cache.
// Performs a timeout on the Get to prevent blocking.
func (h *IecMessageHandler) Do() {
	if h.Retries <= 0 {
		h.Retries = 1
	}
	for n := h.Retries; n > 0; n-- {
		if h.doOnce() == nil {
			return
		}
	}
	log.Print("number of retries exceeded")
}

func (h *IecMessageHandler) doOnce() error {
	// Open the serial port repo.
	m, err := meter.Open(h.PortSettings)
	if err != nil {
		// Log error.
		return err
	}
	defer m.Close()

	// Sleep to make sure the port is ready.
	time.Sleep(500 * time.Millisecond)
	// Channel to receive the measurement message.
	mc := make(chan *measurement)
	// Perform measurement in the background.
	go func() {
		msg := &measurement{}
		msg.m, msg.err = m.Get(nil)
		mc <- msg
	}()
	select {
	case <-time.NewTimer(time.Minute).C:
		// Handle the timeout.
		// Log timeout reading measurement.
		log.Printf("timeout reading a measurement")
		return err
	case msg := <-mc:
		// Measurement or error received.
		if msg.err != nil {
			// Log error reading measurement.
			log.Printf("error reading a measurement: %s", msg.err.Error())
			return msg.err
		}
		if err := h.Repo.Put(msg.m); err != nil {
			// Log error saving measurement.
			log.Printf("Error saving the measurement: %s", msg.err.Error())
			return err
		}
		log.Printf("Stored measurement from %s", msg.m.Identification)
	}
	return nil
}
