package meter

import (
	"errors"
	"log"
	"time"

	"github.com/peterzandbergen/iec62056/iec"
	"github.com/peterzandbergen/iec62056/model"
)

// Meter type represents a meter for reading energy measuerments.
// The timeout should be long enough for slow meters to return the measurement, default is 60 seconds.
type Meter struct {
	// PortSettings for the serial port, needs to be set.
	PortSettings *iec.PortSettings
	// PortName needs to be set.
	PortName string
	// TimeOut for reading the meter in seconds. Default is 60.
	TimeOut time.Duration
}

// Measurement is used internally.
type measurement struct {
	m   *model.Measurement
	err error
}

// Check if the interface has been fully implemented.
var _ model.MeasurementRepo = &Meter{}

// ErrTimeout indicates that reading the meter took too long.
var ErrTimeout = errors.New("timeout reading from meter")

// Copy one reading.
func copyReading(src iec.DataSet) (dst model.DataSet) {
	return model.DataSet{
		Address: src.Address,
		Value:   src.Value,
		Unit:    src.Unit,
	}
}

// Copy all readings.
func copyReadings(src []iec.DataSet) (dst []model.DataSet) {
	for _, s := range src {
		dst = append(dst, copyReading(s))
	}
	return dst
}

// Get returns a measurement from the meter.
// The Key parameter is ignored and can be set to nil.
func (m *Meter) Get(key []byte) (*model.Measurement, error) {
	t := time.Now()
	mm, err := m.readWithTimeout()
	if err != nil {
		return nil, err
	}
	mm.Time = t
	return mm, nil
}

// Put is a noop and should not be called.
// TODO: Return an unsupported error.
func (m *Meter) Put(*model.Measurement) error {
	return nil
}

// GetN returns one measurement.
func (m *Meter) GetN(n int) ([]*model.Measurement, error) {
	mm, err := m.Get(nil)
	if err != nil {
		return nil, err
	}
	return []*model.Measurement{
		mm,
	}, nil
}

func copyMsgToMsm(msg *iec.DataMessage) *model.Measurement {
	return &model.Measurement{
		Identification: msg.MeterID,
		ManufacturerID: msg.ManufacturerID,
		Readings:       copyReadings(msg.DataSets),
	}
}

// readWithTimeout opens the port with the correct settings and reads a value.
// Limits the time to read the measurement with the configured timeout.
func (m *Meter) readWithTimeout() (*model.Measurement, error) {
	if m.TimeOut <= 0 {
		m.TimeOut = 60
	}
	// Open the serial port repo.
	port := iec.New(m.PortSettings)
	err := port.Open(m.PortName)
	if err != nil {
		// Log error.
		return nil, err
	}
	// Close the port when done.
	defer port.Close()

	// Sleep to make sure the port is ready.
	time.Sleep(500 * time.Millisecond)
	// Channel to receive the measurement message.
	mc := make(chan *measurement)
	// Perform measurement in the background.
	go func() {
		msg := &measurement{}
		dm, err := port.Read()
		if err != nil {
			msg.err = err
			mc <- msg
			return
		}
		msg.m, msg.err = copyMsgToMsm(dm), nil
		mc <- msg
	}()
	// Wait for measurement or timeout.
	select {
	case <-time.NewTimer(m.TimeOut * time.Second).C:
		// Log timeout reading measurement.
		log.Printf("timeout reading a measurement")
		return nil, ErrTimeout
	case msg := <-mc:
		// Measurement or error received.
		if msg.err != nil {
			// Log error reading measurement.
			log.Printf("error reading a measurement: %s", msg.err.Error())
			return nil, msg.err
		}
		return msg.m, nil
	}
}

// Delete is a noop and should not be called.
// TODO: return an Unsupported Error.
func (m *Meter) Delete(*model.Measurement) error {
	return nil
}

// PortExists tests if the port exists and can be opened.
// Should only be called once for it will interfere with a Get call.
func (m *Meter) PortExists() bool {
	// Open the serial port repo.
	port := iec.New(m.PortSettings)
	err := port.Open(m.PortName)
	if err != nil {
		// Log error.
		return false
	}
	// Close the port when done.
	defer port.Close()
	return true
}
