package actors

import (
	"fmt"
	"io"
	"log"

	"github.com/peterzandbergen/iec62056/model"
)

// CacheDumper type reads all entries from the Repo and dumps it to the Writer.
type CacheDumper struct {
	Repo         model.MeasurementRepo
	Measurements []*model.Measurement
	Writer       io.Writer
}

// Do performst the actor task.
func (c *CacheDumper) Do() error {
	// Get all entries from the repo.
	m, err := c.Repo.GetN(0)
	if err != nil {
		log.Printf("error reading the local cache: %s\n", err.Error())
		c.Measurements = nil
		return err
	}
	c.Measurements = m
	fmt.Printf("retrieved %d measurements\n", len(m))
	for _, v := range m {
		fmt.Fprintf(c.Writer, "%+v", *v)
	}
	return nil
}
