package handlers

import (
	"github.com/anras5/todo-app-backend/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"os"
	"testing"
)

var app config.Application

func TestMain(m *testing.M) {

	// -------------------------------------------------------------------------------------------- //
	// Set up loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// -------------------------------------------------------------------------------------------- //
	// set repo and handlers
	repo := NewTestRepo(&app)
	NewHandlers(repo)

	os.Exit(m.Run())
}

func getRoutes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	mux.Get("/", Repo.Home)
	mux.Get("/todos", Repo.AllTodos)
	mux.Post("/todos", Repo.InsertTodo)
	mux.Get("/todos/{id}", Repo.OneTodo)
	mux.Put("/todos/{id}", Repo.UpdateTodo)
	mux.Put("/todos/{id}/{complete}", Repo.UpdateTodoCompleted)
	mux.Delete("/todos/{id}", Repo.DeleteTodo)

	return mux
}
