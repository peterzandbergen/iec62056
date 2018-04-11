package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/peterzandbergen/iec62056/actors"
	"github.com/peterzandbergen/iec62056/adapters/cache"
	"github.com/peterzandbergen/iec62056/model"
	"github.com/spf13/pflag"
)

// Service type is used for long running services.
type Service interface {
	// Start starts the service. The call should return immediately.
	Start()
	// Stop stops the service and wait till the service stops.
	Stop(ctx context.Context) error
}

// Options for the program.
type options struct {
	Baudrate         int
	Portname         string
	LocalCache       string
	RemoteStorageURI string
	Interval         int
}

func (o *options) Parse() {
	if flag.Parsed() {
		return
	}
	pflag.IntVarP(&o.Baudrate, "baudrate", "b", 300, "Baudrate of the serial port connected to the energy meter.")
	pflag.StringVarP(&o.Portname, "serial-port", "s", "/dev/ttyUSB0", "Device name of the serial port.")
	pflag.StringVarP(&o.LocalCache, "local-cache-path", "l", "~/.emlogcache", "Location of the local cache.")
	pflag.StringVarP(&o.RemoteStorageURI, "remote-storage-uri", "R", "http://localhost:304725/emeterlog", "Remote Storage Service URI.")
	pflag.IntVarP(&o.Interval, "interval", "I", 300, "Interval for each measurement in seconds.")

	pflag.Parse()
}

// MeasurementHandler type is the handler for measurements.
type MeasurementHandler struct {
	repo model.MeasurementRepo
}

// Handle creates the actor and passes the measurement.
func (h *MeasurementHandler) Handle(m *model.Measurement) {
	// Create the actor to handle the message and store it in the repo.
	a := actors.IecMessageHandler{
		Measurement: m,
		Repo:        h.repo,
	}
	// Call the actor Do function.
	a.Do()
}

// BuildSamplerService builds a service using the given options.
func BuildSamplerService(o options, repo model.MeasurementRepo) (Service, error) {
	s, err := NewSampler(o.Portname, o.Baudrate, time.Duration(o.Interval)*time.Second)
	if err != nil {
		return nil, err
	}
	// Create the repo.
	h := &MeasurementHandler{
		repo: repo,
	}
	s.Handle(h)
	return s, nil
}

func main() {
	log.Println("Starting emlog")
	var services []Service

	// Catch ctrl-C and kill signal.
	c := make(chan os.Signal, 1)
	signal.Notify(c)

	// Process the options.
	o := &options{}
	o.Parse()

	// Create the repositories.
	localRepo, err := cache.Open(o.LocalCache)
	if err != nil {
		// Log and exit.
		os.Exit(1)
	}
	defer localRepo.Close()

	// Create the adapters.
	var s Service
	s, err = BuildSamplerService(*o, localRepo)
	if err != nil {
		log.Printf("cannot create sampler service: %s", err.Error())
		os.Exit(1)
	}
	services = append(services, s)

	// Start the services.
	for _, s := range services {
		go s.Start()
	}

	// Wait for the end
	log.Println("Services started")
	sig := <-c
	log.Printf("Received signal: %s\n", sig.String())

	// Perform clean up.
	log.Println("Stopping services")
	for _, s := range services {
		ctx, done := context.WithTimeout(context.Background(), time.Second*time.Duration(10))
		err := s.Stop(ctx)
		done()
		if err != nil {
			// TODO: Log error
		}
	}
	log.Println("Services stopped")
}
