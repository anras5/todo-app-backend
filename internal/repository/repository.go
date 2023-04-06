package repository

import "github.com/anras5/todo-app-backend/internal/models"

type DatabaseRepo interface {
	AllTodos() ([]*models.Todo, error)
}
