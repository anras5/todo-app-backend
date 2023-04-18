package dbrepo

import (
	"errors"
	"github.com/anras5/todo-app-backend/internal/models"
)

func (m *testDBRepo) SelectTodos(completed ...bool) ([]*models.Todo, error) {
	return nil, nil
}

func (m *testDBRepo) SelectTodo(id int) (*models.Todo, error) {
	return nil, nil
}

func (m *testDBRepo) InsertTodo(todo models.Todo) (int, error) {
	// if the todos name is empty - fail
	if todo.Name == "" {
		return 2, errors.New("error")
	}
	return 1, nil
}

func (m *testDBRepo) UpdateTodo(todo models.Todo) error {
	return nil
}

func (m *testDBRepo) UpdateTodoCompleted(id int, completed bool) error {
	return nil
}

func (m *testDBRepo) DeleteTodo(id int) error {
	return nil
}
