package main

import (
	"github.com/anras5/todo-app-backend/internal/config"
	"github.com/go-chi/chi/v5"
	"testing"
)

func TestRoutes(t *testing.T) {
	var app config.Application

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
	// test passed
	default:
		t.Errorf("type is not *chi.Mux but %T", v)
	}
}
