package main

import (
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestRoutes(t *testing.T) {

	mux := routes()

	switch v := mux.(type) {
	case *chi.Mux:
	// test passed
	default:
		t.Errorf("type is not *chi.Mux but %T", v)
	}
}
