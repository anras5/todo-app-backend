package main

import (
	"net/http"

	"github.com/anras5/todo-app-backend/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(EnableCORS)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/todos", handlers.Repo.AllTodos)
	mux.Post("/todos", handlers.Repo.InsertTodo)
	mux.Get("/todos/{id}", handlers.Repo.OneTodo)
	mux.Put("/todos/{id}", handlers.Repo.UpdateTodo)
	mux.Put("/todos/{id}/{complete}", handlers.Repo.UpdateTodoCompleted)
	mux.Delete("/todos/{id}", handlers.Repo.DeleteTodo)

	return mux
}
