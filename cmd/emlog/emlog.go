package main

import (
	"flag"
	"fmt"

	"github.com/spf13/pflag"
)

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

func main() {
	o := &options{}
	o.Parse()
	flag.PrintDefaults()
	fmt.Printf("Baudrate: %d\nPortname: %s\nLocalcache: %s\n", o.Baudrate, o.Portname, o.LocalCache)
}
