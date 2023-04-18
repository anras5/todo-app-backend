package handlers

import (
	"database/sql"
	"errors"
	"fmt"
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

func NewTestRepo(a *config.Application) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestingRepo(a),
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
		Version: "1.0.1",
	}

	_ = m.App.WriteJSON(w, http.StatusOK, payload)
}

func (m *Repository) AllTodos(w http.ResponseWriter, r *http.Request) {

	var todos []*models.Todo
	var err error
	var searchCompleted bool

	completed := r.URL.Query().Get("completed")
	if completed != "" {
		// we get a completed query params
		searchCompleted, err = strconv.ParseBool(completed)
		if err != nil {
			_ = m.App.ErrorJSON(w, errors.New("completed should be true or false"))
			return
		}
		todos, err = m.DB.SelectTodos(searchCompleted)
	} else {
		// we do not get any query params
		todos, err = m.DB.SelectTodos()
	}

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
		_ = m.App.ErrorJSON(w, err)
		return
	}

	todo, err := m.DB.SelectTodo(todoID)
	if err != nil {
		_ = m.App.ErrorJSON(w, err)
		return
	}
	_ = m.App.WriteJSON(w, http.StatusOK, todo)
}

func (m *Repository) InsertTodo(w http.ResponseWriter, r *http.Request) {
	var todo models.Todo

	err := m.App.ReadJSON(w, r, &todo)
	if err != nil {
		_ = m.App.ErrorJSON(w, err)
		return
	}

	_, err = m.DB.InsertTodo(todo)
	if err != nil {
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

func (m *Repository) UpdateTodoCompleted(w http.ResponseWriter, r *http.Request) {
	todoID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		_ = m.App.ErrorJSON(w, err)
		return
	}
	isCompleted := chi.URLParam(r, "complete")

	switch isCompleted {
	case "complete":
		err = m.DB.UpdateTodoCompleted(todoID, true)
		if err != nil {
			_ = m.App.ErrorJSON(w, err)
			return
		}
	case "incomplete":
		err = m.DB.UpdateTodoCompleted(todoID, false)
		if err != nil {
			_ = m.App.ErrorJSON(w, err)
			return
		}
	default:
		err = errors.New("should provide 'complete' or 'incomplete'")
		_ = m.App.ErrorJSON(w, err)
		return
	}

	response := config.JSONResponse{
		Error:   false,
		Message: fmt.Sprintf("todo updated to %s", isCompleted),
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
