package actors

import (
	"errors"

	"github.com/peterzandbergen/iec62056/model"
)

type PagerActor struct {
	Repo model.MeasurementRepo
}

var (
	ErrNotImplemented = errors.New("not implemented")
	ErrBadArguments   = errors.New("bad arguments")
)

func (a *PagerActor) GetPage(page, pagesize int) ([]*model.Measurement, error) {
	// test the arguments.
	if page < 0 {
		return nil, ErrBadArguments
	}
	if pagesize <= 0 {
		return nil, ErrBadArguments
	}
	return a.Repo.GetPage(page, pagesize)
}

func (a *PagerActor) GetAll() ([]*model.Measurement, error) {
	return a.Repo.GetAll()
}
