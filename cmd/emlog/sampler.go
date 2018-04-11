package main

import (
	"context"
	"sync"
	"time"

	"github.com/peterzandbergen/iec62056/adapters/meter"
	"github.com/peterzandbergen/iec62056/iec"
	"github.com/peterzandbergen/iec62056/model"
)

// SamplerHandler interface should be implemented by Measurement handlers.
type SamplerHandler interface {
	Handle(*model.Measurement)
}

// SamplerHandlerFunc allows to have a func to act as a handler for a Measurement.
type SamplerHandlerFunc func(*model.Measurement)

// Handle processes a measurement.
func (sh SamplerHandlerFunc) Handle(m *model.Measurement) {
	sh.Handle(m)
}

type sampler struct {
	mu       sync.Mutex
	interval time.Duration
	meter    *meter.Meter
	done     chan struct{}
	stopped  chan struct{}
	h        SamplerHandler
}

func NewSampler(port string, baudrate int, interval time.Duration) (*sampler, error) {
	ps := iec.NewDefaultSettings()
	m, err := meter.Open(*ps)
	if err != nil {
		return nil, err
	}
	s := &sampler{}
	s.meter = m
	s.done = make(chan struct{})
	s.stopped = make(chan struct{})
	s.interval = interval
	return s, nil
}

// Handle sets the handler for the
func (s *sampler) Handle(h SamplerHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.h = h
}

// Start blocks till Stop is called.
func (s *sampler) Start() {
	t := time.NewTicker(s.interval)
	for {
		select {
		case <-t.C:
			m, err := s.meter.Get(nil)
			if err != nil {
				// log error
			} else {
				s.h.Handle(m)
			}
		case <-s.done:
			t.Stop()
			s.meter.Close()
			close(s.stopped)
			return
		}
	}
}

// Stop stops the sampler process, blocks till it has been shut down.
func (s *sampler) Stop(ctx context.Context) error {
	closed := make(chan struct{})
	go func() {
		s.mu.Lock()
		// Signal the process to stop.
		close(s.done)
		s.mu.Unlock()
		<-s.stopped
		close(closed)
	}()

	// Wait for the
	select {
	case <-closed:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}

}
