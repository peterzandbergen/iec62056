package iecport

// RequestMessage
// DataMessage type contains the read meter information.
type DataMessage struct {
	ManufacturerID string
	MeterID        string
	EnhancedID     string
	DataSets       []DataSet
}

// DataSet type contains the measurement returned by the meter.
type DataSet struct {
	Address string
	Value   string
	Unit    string
}
