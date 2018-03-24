package iec62056

// DataMessage type contains the read meter information.
type DataMessage struct {
	MmanufacturerID string
	MeterID         string
	EnhancedID      string
	DataSets        []DataSet
}

// DataSet type contains the measurement returned by the meter.
type DataSet struct {
	Address string
	Value   string
	Unit    string
}
