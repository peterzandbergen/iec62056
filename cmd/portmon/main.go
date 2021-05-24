package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
	"bytes"

	"github.com/peterzandbergen/iec62056/iec/telegram"
	"go.bug.st/serial.v1"
)

func readCommands(wg *sync.WaitGroup, in io.Reader, p io.Writer) {
	// connect buffered reader and read bytes.
	b := bytes.Buffer{}
	telegram.SerializeRequestMessage(&b, telegram.RequestMessage{})
	reqMsg := b.Bytes()
	var cr = bufio.NewReader(in)
	time.Sleep(time.Second)
	// telegram.SerializeRequestMessage(p, telegram.RequestMessage{})
	for {
		b, err := cr.ReadByte()
		if err != nil {
			fmt.Printf("error reading command: %s\n", err.Error())
		}
		switch rune(b) {
		case 'r', 'R':
			// Send request.
			fmt.Println("sending request")
			// n, err := telegram.SerializeRequestMessage(p, reqMsg)
			n, err := p.Write(reqMsg)
			if err != nil {
				fmt.Printf("error: %s\n", err)
				break
			} 
			fmt.Printf("sent request bytes: %d\n", n)
		case 'a', 'A':
			// Send request.
			fmt.Println("sending ACK")
			ack := "\x06\x00\x00\x00\r\n"
			n, err := p.Write([]byte(ack))
			if err != nil {
				fmt.Printf("error: %s\n", err)
				break
			} 
			fmt.Printf("sent ack bytes: %d\n", n)
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
		fmt.Fprintf(out, "[%x]", b)
		if b == 0x0D {
			b = ' '
		}
		fmt.Fprint(out, string([]byte{b}))
	}
}

func writeLineResponses(wg *sync.WaitGroup, in *bufio.Reader, out io.Writer) {
	defer wg.Done()
	for {
		b, err := in.ReadByte()
		if err != nil {
			return
		}
		fmt.Fprint(out, string([]byte{b}))
	}
}

func main() {
	var wg = &sync.WaitGroup{}

	p, err := serial.Open("/dev/ttyUSB0",
		&serial.Mode{
			BaudRate: 115200,
			DataBits: 8,
			Parity:   serial.NoParity,
			StopBits: serial.OneStopBit,
		})
	if err != nil {
		fmt.Printf("error opening port: %s\n", err.Error())
		return
	}
	defer p.Close()

	br := bufio.NewReader(p)

	// Start the reader and the writer.
	fmt.Println("starting loop")
	wg.Add(2)
	// p.SetRTS(true)
	go readCommands(wg, os.Stdin, p)
	// go writeResponses(wg, br, os.Stdout)
	// go writeHexResponses(wg, br, os.Stdout)
	go writeLineResponses(wg, br, os.Stdout)
	p.SetRTS(false)
	// p.SetDTR(true)
	wg.Wait()
}
