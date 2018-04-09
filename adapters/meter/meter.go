package meter

import (
	"time"

	"github.com/peterzandbergen/iec62056/iec"
	"github.com/peterzandbergen/iec62056/model"
)

// Meter type represents a meter for reading energy measuerments.
type Meter struct {
	port *iec.Port
}

var _ model.MeasurementRepo = &Meter{}

// Open returns a new meter on the port.
// It will open the port.
func Open(ps iec.PortSettings) (*Meter, error) {
	p := iec.New(&ps)
	err := p.Open(ps.PortName)
	if err != nil {
		return nil, err
	}
	return &Meter{
		port: p,
	}, nil
}

// Close closes the meter. Must be called to prevent resource leaking.
func (m *Meter) Close() {
	m.port.Close()
}

func copyReadings(src []iec.DataSet) (dst []model.DataSet) {
	for _, s := range src {
		d := model.DataSet{
			Address: s.Address,
			Value:   s.Value,
			Unit:    s.Unit,
		}
		dst = append(dst, d)
	}
	return dst
}

// Get returns a measurement from the meter.
// Key is not used and should be set to nil.
func (m *Meter) Get(key []byte) (*model.Measurement, error) {
	t := time.Now()
	dm, err := m.port.Read()
	if err != nil {
		return nil, err
	}

	return &model.Measurement{
		Time:           t,
		ManufacturerID: dm.ManufacturerID,
		Identification: dm.MeterID,
		Readings:       copyReadings(dm.DataSets),
	}, nil
}

// Put is a noop.
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

// Delete is a noop and should not be called.
// TODO: return an Unsupported Error.
func (m *Meter) Delete(*model.Measurement) error {
	return nil
}
