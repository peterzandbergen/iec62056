package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/peterzandbergen/iec62056/actors"
	"github.com/peterzandbergen/iec62056/model"
)

var (
	ErrBadParameter = errors.New("parameter error")
)

type HttpLocalService struct {
	listenAddress string
	localRepo     model.MeasurementRepo
	server        *http.Server
}

type GetAllHandler struct {
	server *HttpLocalService
}

type MeasurementsResponse struct {
	FirstTime    time.Time `json:"omitempty"`
	LastTime     time.Time `json:"omitempty"`
	Measurements []*model.Measurement
}

type errPagination struct {
	strings.Builder
}

func (s *errPagination) Error() string {
	return "bad pagination parameters\n" + s.String()
}

type pagination struct {
	page, size int
	err        *errPagination
}

func NewPagination(r *http.Request) *pagination {
	p := new(pagination)
	p.getParams(r)
	return p
}

func (p *pagination) getParams(r *http.Request) {
	page := r.FormValue("page")
	size := r.FormValue("size")

	p.page = 0
	p.size = 0

	serr := &errPagination{}
	if len(page) != 0 {
		if v, err := strconv.Atoi(r.FormValue("page")); err != nil {
			fmt.Fprintf(serr, "\tpage parameter error: %s\n", err.Error())
		} else {
			if v < 0 {
				fmt.Fprint(serr, "\tpage parameter cannog be negative\n")
			} else {
				p.page = v
			}
		}
	}
	if len(size) > 0 {
		if v, err := strconv.Atoi(r.FormValue("size")); err != nil {
			fmt.Fprintf(serr, "\tsize parameter error: %s\n", err.Error())
		} else {
			if v < 0 {
				fmt.Fprint(serr, "\tsize parameter cannot be negative\n")
			} else {
				p.size = v
			}
		}
	}
	if p.page > 0 && p.size == 0 {
		fmt.Fprint(serr, "\tnon zero page parameter requires non zero limit\n")
	}
	if serr.Len() > 0 {
		p.err = serr
	}
}

func (p *pagination) paginate() bool {
	return p.err == nil && p.size > 0
}

func NewHttpLocalService(address string, repo model.MeasurementRepo) Service {
	sm := &http.ServeMux{}
	svc := &HttpLocalService{
		listenAddress: address,
		localRepo:     repo,
		server: &http.Server{
			Handler: sm,
			Addr:    address,
		},
	}
	gah := &GetAllHandler{
		server: svc,
	}
	// Add handlers.
	sm.Handle("/measurements/", gah)
	return svc
}

// ServeHTTP reads all entries from the local repo and returns the JSON.
func (h *GetAllHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Create an actor.
	a := &actors.PagerActor{
		Repo: h.server.localRepo,
	}
	// Determine the pagination parameters.
	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("bad request error: %s", err.Error()), http.StatusBadRequest)
		return
	}
	var pag *pagination
	if pag = NewPagination(r); pag.err != nil {
		http.Error(w, fmt.Sprintf("bad request error: %s", pag.err.Error()), http.StatusBadRequest)
		return
	}
	var msm []*model.Measurement
	var err error
	if pag.paginate() {
		msm, err = a.GetPage(pag.page, pag.size)
	} else {
		msm, err = a.GetAll()
	}
	// Get the data.
	if err != nil {
		http.Error(w, fmt.Sprintf("internal error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	response := &MeasurementsResponse{
		Measurements: msm,
	}
	if !pag.paginate() {
		response.FirstTime = msm[0].Time
		response.LastTime = msm[len(msm)-1].Time
	}
	// Content type
	w.Header().Set("Content-Type", "application/json")
	// Take the output and serialize to the writer.
	j, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("internal error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	log.Printf("writing measurements response...")
	w.Write(j)
	log.Printf("writing measurements response...done")
}

// Start starts the HTTP server on the given address and port.
func (s *HttpLocalService) Start(ctx context.Context) error {
	var err error
	var done = make(chan struct{})
	go func() {
		err = s.server.ListenAndServe()
		close(done)
	}()
	select {
	case <-done:
		return err
	case <-time.After(time.Second):
		return nil
	}
}

func (s *HttpLocalService) Stop(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
