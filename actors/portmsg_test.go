package actors

import (
	"testing"

	"github.com/peterzandbergen/iec62056/model"
)

// type MeasurementRepo interface {
// 	Put(*Measurement) error
// 	Get(key []byte) (*Measurement, error)
// 	GetN(n int) ([]*Measurement, error)
// 	Delete(*Measurement) error
// }

type noopRepo struct {
}

func (r noopRepo) Put(*model.Measurement) error {
	return nil
}

func (r noopRepo) Get(key []byte) (*model.Measurement, error) {
	return nil, nil
}

func (r noopRepo) GetN(n int) ([]*model.Measurement, error) {
	return nil, nil
}

func (r noopRepo) Delete(*model.Measurement) error {
	return nil
}

// Mock repos
type mockMeterRepo struct {
	noopRepo
}

type mockLocalRepo struct {
	noopRepo
}

func TestIecMessageHandlerDo(t *testing.T) {
	meterRepo := &mockMeterRepo{}
	localRepo := &mockLocalRepo{}
	actor := IecMessageHandler{
		LocalRepo: localRepo,
		MeterRepo: meterRepo,
	}
	actor.Do()
}
