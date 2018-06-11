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

func (a *PagerActor) GetFirst() (*model.Measurement, error) {
	var msm []*model.Measurement
	var err error
	if msm, err = a.Repo.GetAll(); err != nil {
		return nil, err
	}
	return msm[0], nil
}

func (a *PagerActor) GetLast() (*model.Measurement, error) {
	var msm []*model.Measurement
	var err error
	if msm, err = a.Repo.GetAll(); err != nil {
		return nil, err
	}
	return msm[len(msm)-1], nil
}
