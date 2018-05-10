package service

import (
	"context"
	"sync"
	"time"
)

type Timer struct {
	lock     sync.Mutex
	h        TimerHandler
	interval time.Duration
	done     chan *doneMsg
}

// doneMsg is send when the service needs to stop.
type doneMsg struct {
	ack chan struct{}
}

type TimerHandler interface {
	Handle(time.Time)
}

type TimerHandleFunc func(time.Time)

func (f TimerHandleFunc) Handle(t time.Time) {
	f(t)
}

// NewTimer creates a new timer and returns a Service.
func NewTimer(interval time.Duration, h TimerHandler) Service {
	return &Timer{
		h:        h,
		interval: interval,
	}
}

func (t *Timer) serve() {
	go func() {
		tick := time.NewTicker(t.interval)
		defer tick.Stop()
		for {
			t.h.Handle(time.Now())
			select {
			case <-tick.C:
			case msg := <-t.done:
				close(msg.ack)
				return
			}
		}
	}()
}

func (t *Timer) Start(ctx context.Context) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	if t.done != nil {
		// Already running.
		return ErrAlreadyRunning
	}
	t.done = make(chan *doneMsg)
	t.serve()
	return nil
}

// Stop the timer.
func (t *Timer) Stop(ctx context.Context) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	if t.done == nil {
		return ErrNotRunning
	}
	sm := &doneMsg{
		ack: make(chan struct{}),
	}
	// Send the done message to done and wait for the ack channel to close.
	t.done <- sm
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-sm.ack:
		t.done = nil
		return nil
	}
}
