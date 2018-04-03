package main

import (
	"flag"
	"fmt"
)

type options struct {
	Baudrate   int
	Portname   string
	LocalCache string
}

func (o *options) Parse() {
	if flag.Parsed() {
		return
	}
	flag.IntVar(&o.Baudrate, "baudrate", 300, "Baudrate of the serial port connected to the energy meter.")
	flag.StringVar(&o.Portname, "serial-port", "/dev/ttyUSB0", "Device name of the serial port")
	flag.StringVar(&o.LocalCache, "local-cache-path", "/var/lib/emlogcache", "Location of the local cache")
	flag.Parse()
}

func main() {
	o := &options{}
	o.Parse()
	flag.PrintDefaults()
	fmt.Printf("Baudrate: %d\nPortname: %s\nLocalcache: %s\n", o.Baudrate, o.Portname, o.LocalCache)
}
