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

	mux.Use(middleware.RequestID)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/todos", handlers.Repo.AllTodos)
	mux.Post("/todos", handlers.Repo.InsertTodo)
	mux.Get("/todos/{id}", handlers.Repo.OneTodo)
	mux.Put("/todos/{id}", handlers.Repo.UpdateTodo)
	mux.Delete("/todos/{id}", handlers.Repo.DeleteTodo)

	return mux
}
