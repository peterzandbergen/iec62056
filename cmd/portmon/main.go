package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/peterzandbergen/iec62056/telegram"

	"go.bug.st/serial.v1"
)

func readCommands(wg *sync.WaitGroup, in io.Reader, p io.Writer) {
	// connect buffered reader and read bytes.
	var cr = bufio.NewReader(in)

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

func writeResponses(wg *sync.WaitGroup, in io.Reader, out io.Writer) {
	for {
		var buf = make([]byte, 1000)
		n, err := in.Read(buf)
		if err != nil {
			fmt.Fprintf(out, "Error %s, aborting\n", err.Error())
			wg.Done()
			return
		}
		n, err = out.Write(buf[:n])
	}
}

func main() {
	var wg = &sync.WaitGroup{}

	p, err := serial.Open("/dev/ttyUSB0",
		&serial.Mode{
			BaudRate: 9600,
			DataBits: 7,
			Parity:   serial.EvenParity,
			StopBits: serial.OneStopBit,
		})
	if err != nil {
		fmt.Printf("error opening port: %s\n", err.Error())
		return
	}
	defer p.Close()

	// Start the reader and the writer.
	wg.Add(2)
	go readCommands(wg, os.Stdin, p)
	go writeResponses(wg, p, os.Stdout)
	wg.Wait()
}
