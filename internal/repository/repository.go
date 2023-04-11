package repository

import "github.com/anras5/todo-app-backend/internal/models"

type DatabaseRepo interface {
	SelectTodos() ([]*models.Todo, error)
	SelectTodo(id int) (*models.Todo, error)
	InsertTodo(todo models.Todo) (int, error)
}
