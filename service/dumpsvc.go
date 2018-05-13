package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/peterzandbergen/iec62056/actors"
	"github.com/peterzandbergen/iec62056/model"
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
	FirstTime    time.Time
	LastTime     time.Time
	Measurements []*model.Measurement
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
	msm, err := a.GetAll()
	if err != nil {
		http.Error(w, fmt.Sprintf("internal error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	response := &MeasurementsResponse{
		FirstTime:    msm[0].Time,
		LastTime:     msm[len(msm)-1].Time,
		Measurements: msm,
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
