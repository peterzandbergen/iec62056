package actors

import (
	"log"

	"github.com/peterzandbergen/iec62056/model"
)

type CacheDumper struct {
	Repo         model.MeasurementRepo
	Measurements []*model.Measurement
}

// Do performst the actor task.
func (c *CacheDumper) Do() {
	// Get all entries from the repo.
	m, err := c.Repo.GetN(0)
	if err != nil {
		log.Printf("error reading the local cache: %s\n", err.Error())
		c.Measurements = nil
	} else {
		c.Measurements = m
	}
}
