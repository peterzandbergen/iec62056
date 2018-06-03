package service

import (
	"net/http/httptest"
	"testing"
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
	r := httptest.NewRequest("GET", "http://localhost/measurements/?page=100", nil)
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
