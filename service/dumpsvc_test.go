package service

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/peterzandbergen/iec62056/model"
)

func TestAllRequest(t *testing.T) {
	// Create a test request.
	r := httptest.NewRequest("GET", "http://localhost/measurements/", nil)
	p := NewPagination(r)
	if p.err != nil {
		t.Fatalf("unexpected error: %s", p.err.Error())
	}
	if p.paginate() {
		t.Fatalf("did not expect pagination")
	}
}

func TestPageZero(t *testing.T) {
	// Create a test request.
	r := httptest.NewRequest("GET", "http://localhost/measurements/page=0", nil)
	p := NewPagination(r)
	if p.err != nil {
		t.Fatalf("unexpected error: %s", p.err.Error())
	}
	if p.paginate() {
		t.Fatalf("did not expect pagination")
	}
}

func TestLimitOnly(t *testing.T) {
	// Create a test request.
	r := httptest.NewRequest("GET", "http://localhost/measurements/?size=100", nil)
	p := NewPagination(r)
	if p.err != nil {
		t.Fatalf("unexpected error: %s", p.err.Error())
	}
	if !p.paginate() {
		t.Fatalf("expected pagination")
	}
	if p.size != 100 {
		t.Fatalf("expected %d, got %d", 100, p.size)
	}
}

func TestPageOnly(t *testing.T) {
	// Create a test request.
	r := httptest.NewRequest("GET", "http://localhost/measurements?page=100", nil)
	p := NewPagination(r)
	if p.err == nil {
		t.Fatalf("expected error")
	}
}

func TestPaginationOnly(t *testing.T) {
	// Create a test request.
	r := httptest.NewRequest("GET", "http://localhost/measurements/?page=200&size=100", nil)
	p := NewPagination(r)
	if p.err != nil {
		t.Fatalf("unexpected error: %s", p.err.Error())
	}
	if !p.paginate() {
		t.Fatalf("expected pagination")
	}
	if p.size != 100 {
		t.Fatalf("expected %d, got %d", 100, p.size)
	}
	if p.page != 200 {
		t.Fatalf("expected %d, got %d", 200, p.page)
	}
}

func TestContextAll(t *testing.T) {
	// Create a test request.
	r := httptest.NewRequest("GET", "http://localhost/measurements", nil)
	c := getContext(r)
	if c.err != nil {
		t.Errorf("unexptected error: %s", c.err.Error())
	}
	if c.first {
		t.Error("first is true")
	}
	if c.last {
		t.Error("last is true")
	}
	if c.pag.paginate() {
		t.Error("paginage is true")
	}
}

func TestContextAllSlash(t *testing.T) {
	// Create a test request.
	r := httptest.NewRequest("GET", "http://localhost/measurements/", nil)
	c := getContext(r)
	if c.err != nil {
		t.Errorf("unexptected error: %s", c.err.Error())
	}
	if c.first {
		t.Error("first is true")
	}
	if c.last {
		t.Error("last is true")
	}
	if c.pag.paginate() {
		t.Error("paginage is true")
	}
}

func TestContextFirst(t *testing.T) {
	// Create a test request.
	r := httptest.NewRequest("GET", "http://localhost/measurements/first", nil)
	c := getContext(r)
	if c.err != nil {
		t.Errorf("unexptected error: %s", c.err.Error())
	}
	if !c.first {
		t.Error("first is false")
	}
	if c.last {
		t.Error("last is true")
	}
	if c.pag != nil && c.pag.paginate() {
		t.Error("paginage is true")
	}
}

func TestContextLast(t *testing.T) {
	// Create a test request.
	r := httptest.NewRequest("GET", "http://localhost/measurements/last", nil)
	c := getContext(r)
	if c.err != nil {
		t.Errorf("unexptected error: %s", c.err.Error())
	}
	if c.first {
		t.Error("first is true")
	}
	if !c.last {
		t.Error("last is false")
	}
	if c.pag != nil && c.pag.paginate() {
		t.Error("paginage is true")
	}
}

func TestContextPaginate(t *testing.T) {
	// Create a test request.
	r := httptest.NewRequest("GET", "http://localhost/measurements?page=100&size=200", nil)
	c := getContext(r)
	if c.err != nil {
		t.Errorf("unexptected error: %s", c.err.Error())
	}
	if c.first {
		t.Error("first is true")
	}
	if c.last {
		t.Error("last is true")
	}
	if c.pag == nil {
		t.Error("c.pag is nil")
	}
	if c.pag != nil && !c.pag.paginate() {
		t.Error("paginage is false")
	}
	if c.pag.page != 100 {
		t.Error("page != 100")
	}
	if c.pag.size != 200 {
		t.Error("size != 200")
	}
}

const firstResponse = `{"First":{}}`

func TestMeasurementResponseMarshal(t *testing.T) {
	resp := &MeasurementsResponse{
		Data: &model.Measurement{
			Identification: "identification",
		},
	}
	js, err := json.Marshal(resp)
	if err == nil {
		t.Log(string(js))
		// t.Error()
	}
}
