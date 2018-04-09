// Package iecport digests the bytes from the serial port and produces a stream of
// measurement messages. The messages are equal to the messages in the model.
// The messages are presented on channel.
package iecstream

import (
	"sync"

	"github.com/peterzandbergen/iec62056/iec"
	"github.com/peterzandbergen/iec62056/model"
)

// IecStream type converts the messages from the serial port to Measurements.
type Stream struct {
	m       sync.Mutex
	running bool
	c       chan *model.Measurement
	p       iec.Port
}

func (i *Stream) OpenPort() error {
	i.m.Lock()
	defer i.m.Unlock()
	return nil
}

func (i *Stream) Start() error {
	i.m.Lock()
	defer i.m.Unlock()

	return nil
}
