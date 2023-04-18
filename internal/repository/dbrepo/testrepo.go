package dbrepo

import (
	"errors"
	"github.com/anras5/todo-app-backend/internal/models"
)

func (m *testDBRepo) SelectTodos(completed ...bool) ([]*models.Todo, error) {
	// if completed is false, then fail
	if len(completed) > 0 && !completed[0] {
		return nil, errors.New("error")
	}
	return nil, nil
}

func (m *testDBRepo) SelectTodo(id int) (*models.Todo, error) {
	// if id is 2, then fail
	if id == 2 {
		return nil, errors.New("error")
	}
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
	if todo.ID == 2 {
		return errors.New("error")
	}
	return nil
}

func (m *testDBRepo) UpdateTodoCompleted(id int, completed bool) error {
	if id == 2 {
		return errors.New("error")
	}
	return nil
}

func (m *testDBRepo) DeleteTodo(id int) error {
	if id == 2 {
		return errors.New("error")
	}
	return nil
}
