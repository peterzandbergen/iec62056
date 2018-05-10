package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/peterzandbergen/iec62056/model"

	"github.com/peterzandbergen/iec62056/adapters/meter"

	"github.com/peterzandbergen/iec62056/actors"
	"github.com/peterzandbergen/iec62056/adapters/cache"
	"github.com/peterzandbergen/iec62056/iec"
	"github.com/peterzandbergen/iec62056/service"
	"github.com/spf13/pflag"
)

// Options for the program.
type options struct {
	DumpCache        bool
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
	pflag.BoolVarP(&o.DumpCache, "show-cache", "D", false, "Dump the content of the cache.")
	pflag.IntVarP(&o.Baudrate, "baudrate", "b", 300, "Baudrate of the serial port connected to the energy meter.")
	pflag.StringVarP(&o.Portname, "serial-port", "s", "/dev/ttyUSB0", "Device name of the serial port.")
	pflag.StringVarP(&o.LocalCache, "local-cache-path", "l", "/tmp/emlog-cache", "Location of the local cache.")
	pflag.StringVarP(&o.RemoteStorageURI, "remote-storage-uri", "R", "http://localhost:304725/emeterlog", "Remote Storage Service URI.")
	pflag.IntVarP(&o.Interval, "interval", "I", 300, "Interval for each measurement in seconds.")

	pflag.Parse()
}

func buildTimerHandler(meterRepo, localRepo model.MeasurementRepo) service.TimerHandler {
	return service.TimerHandleFunc(func(t time.Time) {
		a := actors.IecMessageHandler{
			LocalRepo: localRepo,
			MeterRepo: meterRepo,
		}
		a.Do()
	})
}

func buildMeterRepo(options *options) *meter.Meter {
	ps := iec.NewDefaultSettings()
	ps.PortName = options.Portname
	ps.InitialBaudRateModeABC = options.Baudrate
	mr := &meter.Meter{
		PortName:     options.Portname,
		PortSettings: ps,
	}
	return mr
}

func main() {
	log.Println("Starting emlog")

	// Catch ctrl-C and kill signal.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Process the options.
	o := &options{}
	o.Parse()

	// Create the repositories.
	// Local cache.
	localRepo, err := cache.Open(o.LocalCache)
	if err != nil {
		// Log and exit.
		os.Exit(1)
	}
	defer localRepo.Close()

	// Meter
	meterRepo := buildMeterRepo(o)
	if meterRepo == nil {
		log.Println("cannot open the meter repo")
		os.Exit(1)
	}

	// Create the services.

	// The measurement service.
	timerSvc := service.NewTimer(time.Duration(o.Interval)*time.Second, buildTimerHandler(meterRepo, localRepo))

	// TODO: The status REST service.
	// TODO: The save to cloud service.

	if o.DumpCache {
		// Create and start the CacheDumper.
		a := &actors.CacheDumper{
			Repo:   localRepo,
			Writer: os.Stdout,
		}
		a.Do()
		os.Exit(0)
	}

	// Create services list.
	services := service.NewServicesList(timerSvc)

	if err := services.Start(context.Background()); err != nil {
		log.Printf("error stating the services: %s", err.Error())
		os.Exit(1)
	}

	// Wait for the end
	log.Println("Services started")
	sig := <-c
	log.Printf("Received signal: %s\n", sig.String())

	// Perform clean up.
	log.Println("Stopping services")
	services.Stop(context.Background())

	log.Println("Services stopped")
}
