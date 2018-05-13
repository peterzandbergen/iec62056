package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	// ErrAlreadyRunning when service is already running.
	ErrAlreadyRunning = errors.New("service is already running")
	// ErrNotRunning when the service is not running.
	ErrNotRunning = errors.New("service not running")
)

// Service is the interface for services, for stopping and starting.
type Service interface {
	// Start starts the service and returns after the service has been successfully started.
	Start(ctx context.Context) error
	// Stop stops a running service.
	Stop(ctx context.Context) error
}

type svcError struct {
	err  string
	errs []error
}

func (se *svcError) Error() string {
	res := &strings.Builder{}
	fmt.Fprintf(res, "Error: %s\n", se.err)
	for _, e := range se.errs {
		fmt.Fprintf(res, "    err: %s\n", e.Error())
	}
	return res.String()
}

func (se *svcError) empty() bool {
	return len(se.errs) <= 0
}

func (se *svcError) add(err error) {
	se.errs = append(se.errs, err)
}

type serviceList struct {
	services []Service
}

func NewServicesList(services ...Service) Service {
	sl := &serviceList{
		services: make([]Service, len(services)),
	}
	for i, s := range services {
		sl.services[i] = s
	}
	return sl
}

func (s *serviceList) Add(svc Service) {
	s.services = append(s.services, svc)
}

func (s *serviceList) Start(ctx context.Context) error {
	se := &svcError{}
	for _, s := range s.services {
		if err := s.Start(ctx); err != nil {
			se.add(err)
		}
	}
	if se.empty() {
		return nil
	}
	return se
}

func (s *serviceList) Stop(ctx context.Context) error {
	se := &svcError{}
	for _, s := range s.services {
		tctx, cancel := context.WithTimeout(ctx, time.Duration(30)*time.Second)
		if err := s.Stop(tctx); err != nil {
			se.add(err)
		}
		cancel()
	}
	if se.empty() {
		return nil
	}
	return se
}
