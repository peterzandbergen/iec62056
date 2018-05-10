package service

import (
	"context"
	"net/http"
)

type Rest struct {
	s *http.Server
}

func NewRest() *Rest {
	return nil
}

func (r *Rest) Start(ctx context.Context) error {
	return nil
}

func (r *Rest) Stop(ctx context.Context) error {
	// Perform internal cleanup.
	// ShutDown the http server.
	return r.s.Shutdown(ctx)
}
