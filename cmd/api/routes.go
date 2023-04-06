package main

import (
	"github.com/anras5/todo-app-backend/internal/config"
	"github.com/anras5/todo-app-backend/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func routes(app *config.Application) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)

	mux.Get("/", handlers.Repo.Home)

	return mux
}
