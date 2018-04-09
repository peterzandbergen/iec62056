package meter

import (
	"testing"

	"github.com/peterzandbergen/iec62056/iec"
)

var portSettings = iec.PortSettings{
	PortName: "/dev/ttyUSB0",
}

func TestOpenAndClose(t *testing.T) {
	ps := portSettings
	m, err := Open(ps)
	if err != nil {
		t.Fatalf("Open, error: %s", err.Error())
	}
	defer m.Close()
}

func TestGet(t *testing.T) {

}
