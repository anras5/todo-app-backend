package repository

import "github.com/anras5/todo-app-backend/internal/models"

type DatabaseRepo interface {
	SelectTodos(completed ...bool) ([]*models.Todo, error)
	SelectTodo(id int) (*models.Todo, error)
	InsertTodo(todo models.Todo) (int, error)
	UpdateTodo(todo models.Todo) error
	UpdateTodoCompleted(id int, completed bool) error
	DeleteTodo(id int) error
}
