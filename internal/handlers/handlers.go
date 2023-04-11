package handlers

import (
	"database/sql"
	"github.com/anras5/todo-app-backend/internal/config"
	"github.com/anras5/todo-app-backend/internal/models"
	"github.com/anras5/todo-app-backend/internal/repository"
	"github.com/anras5/todo-app-backend/internal/repository/dbrepo"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.Application
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.Application, db *sql.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db, a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	payload := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Version string `json:"version"`
	}{
		Status:  "active",
		Message: "Go Todos API is up and running",
		Version: "1.0.0",
	}

	_ = m.App.WriteJSON(w, http.StatusOK, payload)
}

func (m *Repository) AllTodos(w http.ResponseWriter, r *http.Request) {

	todos, err := m.DB.SelectTodos()
	if err != nil {
		m.App.ErrorLog.Println(err)
		_ = m.App.ErrorJSON(w, err)
		return
	}

	_ = m.App.WriteJSON(w, http.StatusOK, todos)
}

func (m *Repository) OneTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	todoID, err := strconv.Atoi(id)
	if err != nil {
		m.App.ErrorLog.Println(err)
		_ = m.App.ErrorJSON(w, err)
		return
	}

	todo, err := m.DB.SelectTodo(todoID)
	if err != nil {
		m.App.ErrorLog.Println(err)
		_ = m.App.ErrorJSON(w, err)
		return
	}
	_ = m.App.WriteJSON(w, http.StatusOK, todo)
}

func (m *Repository) InsertTodo(w http.ResponseWriter, r *http.Request) {
	var todo models.Todo

	err := m.App.ReadJSON(w, r, &todo)
	if err != nil {
		m.App.ErrorLog.Println(err)
		_ = m.App.ErrorJSON(w, err)
		return
	}

	_, err = m.DB.InsertTodo(todo)
	if err != nil {
		m.App.ErrorLog.Println(err)
		_ = m.App.ErrorJSON(w, err)
		return
	}

	response := config.JSONResponse{
		Error:   false,
		Message: "todo inserted",
	}
	m.App.WriteJSON(w, http.StatusAccepted, response)

}

func (m *Repository) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	var todo models.Todo

	err := m.App.ReadJSON(w, r, &todo)
	if err != nil {
		_ = m.App.ErrorJSON(w, err)
		return
	}

	err = m.DB.UpdateTodo(todo)
	if err != nil {
		_ = m.App.ErrorJSON(w, err)
		return
	}

	response := config.JSONResponse{
		Error:   false,
		Message: "todo updated",
	}
	m.App.WriteJSON(w, http.StatusAccepted, response)
}

func (m *Repository) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	todoID, err := strconv.Atoi(id)
	if err != nil {
		_ = m.App.ErrorJSON(w, err)
		return
	}

	err = m.DB.DeleteTodo(todoID)
	if err != nil {
		_ = m.App.ErrorJSON(w, err)
		return
	}

	response := config.JSONResponse{
		Error:   false,
		Message: "todo deleted",
	}
	m.App.WriteJSON(w, http.StatusAccepted, response)

}
