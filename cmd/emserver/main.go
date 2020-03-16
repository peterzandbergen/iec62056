package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/peterzandbergen/iec62056/actors"
	"github.com/peterzandbergen/iec62056/adapters/cache"
	"github.com/peterzandbergen/iec62056/model"
	"github.com/peterzandbergen/iec62056/service"

	"github.com/spf13/pflag"
)

// Options for the program.
type options struct {
	LocalCache string
	ListenPort int
}

func (o *options) Parse() {
	if flag.Parsed() {
		return
	}
	pflag.StringVarP(&o.LocalCache, "local-cache-path", "l", "/tmp/emlog-cache", "Location of the local cache.")
	pflag.IntVarP(&o.ListenPort, "port", "p", 8080, "Port to listen on, can also be set using the PORT env var.")
	pflag.Parse()

	port := os.Getenv("PORT")
	if len(port) == 0 {
		return
	}
	if p, err := strconv.Atoi(port); err != nil {
		o.ListenPort = p
	}
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

func main() {
	log.Println("Starting emserver")

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

	// Create the services.

	// TODO: The status REST service.
	la := ":" + strconv.Itoa(o.ListenPort)
	log.Printf("Listening on %s", la)
	localRestSvc := service.NewHttpLocalService(la, localRepo)

	// Create services list.
	services := service.NewServicesList(localRestSvc)
	// services := service.NewServicesList(localRestSvc)

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
