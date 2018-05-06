package meter

import (
	"testing"

	"github.com/peterzandbergen/iec62056/iec"
)

var portSettings = iec.PortSettings{
	PortName: "/dev/ttyUSB0",
}

func TestGet(t *testing.T) {
	ps := iec.NewDefaultSettings()
	m := &Meter{
		PortSettings: ps,
		PortName:     "/dev/ttyUSB0",
	}
	msm, err := m.Get(nil)
	if err != nil {
		t.Fatalf("Get failed, error: %s", err.Error())
	}
	t.Logf("Measurement: %v", *msm)
}
