package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/peterzandbergen/iec62056/iec/telegram"
	"go.bug.st/serial.v1"
)

func readCommands(wg *sync.WaitGroup, in io.Reader, p io.Writer) {
	// connect buffered reader and read bytes.
	var cr = bufio.NewReader(in)
	time.Sleep(time.Second)
	telegram.SerializeRequestMessage(p, telegram.RequestMessage{})
	for {
		b, err := cr.ReadByte()
		if err != nil {
			fmt.Printf("error reading command: %s\n", err.Error())
		}
		switch rune(b) {
		case 'r', 'R':
			// Send request.
			telegram.SerializeRequestMessage(p, telegram.RequestMessage{})
		case 'a', 'A':
			// Send request.
			ack := "\x06\x00\x00\x00\r\n"
			p.Write([]byte(ack))
		case 'q', 'Q':
			wg.Done()
			return
		}
	}
}

func writeResponses(wg *sync.WaitGroup, in *bufio.Reader, out io.Writer) {
	defer wg.Done()
	for {
		idm, err := telegram.ParseIdentificationMessage(in)
		if err != nil {
			fmt.Fprintf(out, "Error receiving ID message: %s\n", err.Error())
			continue
		}
		fmt.Fprintf(out, "%+v\n\n", *idm)
		dm, err := telegram.ParseDataMessage(in)
		if err != nil {
			fmt.Fprintf(out, "Error receiving data message: %s\n", err.Error())
			continue
		}
		fmt.Fprint(out, dm.String())
	}
}

func writeHexResponses(wg *sync.WaitGroup, in *bufio.Reader, out io.Writer) {
	defer wg.Done()
	for {
		b, err := in.ReadByte()
		if err != nil {
			return
		}
		fmt.Fprintf(out, "[%c]", b)
		if b == 0x0D {
			b = ' '
		}
		fmt.Fprint(out, rune(b))
	}
}

func main() {
	var wg = &sync.WaitGroup{}

	p, err := serial.Open("/dev/ttyUSB0",
		&serial.Mode{
			BaudRate: 300,
			DataBits: 7,
			Parity:   serial.EvenParity,
			StopBits: serial.OneStopBit,
		})
	if err != nil {
		fmt.Printf("error opening port: %s\n", err.Error())
		return
	}
	defer p.Close()

	br := bufio.NewReader(p)

	// Start the reader and the writer.
	wg.Add(2)
	go readCommands(wg, os.Stdin, p)
	go writeResponses(wg, br, os.Stdout)
	// go writeHexResponses(wg, br, os.Stdout)
	wg.Wait()
}
