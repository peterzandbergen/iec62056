package model

import "time"

// MeasurementRepo interface should be implemented by adapters providing
type MeasurementRepo interface {
	Put(*Measurement) error
	Get(key []byte) (*Measurement, error)
	GetPage(page, pagesize int) ([]*Measurement, error)
	GetAll() ([]*Measurement, error)
	Delete(*Measurement) error
}

// Measurement type contains a measurement for a meter.
type Measurement struct {
	Time           time.Time
	ManufacturerID string
	Identification string
	Readings       []DataSet
}

// DataSet is a measurement of a variable. Follows the OBIS scheme for the address.
type DataSet struct {
	Address string
	Value   string
	Unit    string
}

var addressMap = map[string]string{
	"1.8.1": "ConsumedEnergyTarif1",
	"1.8.2": "ConsumedEnergyTarif2",
	"2.8.1": "ProducedEnergyTarif1",
	"2.8.2": "ProducedEnergyTarif2",
}

// Address type for a Stringifier.
type Address string

func (a Address) String() string {
	if v, ok := addressMap[string(a)]; ok {
		return v
	}
	return string(a)
}
