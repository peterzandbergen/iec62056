package service

import (
	"context"
	"testing"
	"time"
)

func TestTimerStartStop(t *testing.T) {

	h := TimerHandleFunc(func(t time.Time) {
		time.Sleep(time.Duration(5) * time.Second)
	})
	svc := NewTimer(time.Second, h)
	_ = svc

	bctx := context.Background()
	ctx, _ := context.WithTimeout(bctx, time.Second)
	err := svc.Start(ctx)
	if err != nil {
		t.Errorf("Error starting service: %s", err.Error())
	}
	ctx, _ = context.WithTimeout(bctx, time.Second)
	err = svc.Stop(ctx)
	if err != nil {
		t.Errorf("Error stopping service: %s", err.Error())
	}
}
