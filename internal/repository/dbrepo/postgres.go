package dbrepo

import (
	"context"
	"github.com/anras5/todo-app-backend/internal/models"
	"time"
)

func (m *postgresDBRepo) SelectTodos() ([]*models.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var todos []*models.Todo

	query := `
SELECT ID, NAME, DESCRIPTION, DEADLINE, COMPLETED, CREATED_AT, UPDATED_AT
FROM TODO
`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		m.App.ErrorLog.Println(err)
		return nil, err
	}

	for rows.Next() {
		var todo models.Todo
		err := rows.Scan(
			&todo.ID,
			&todo.Name,
			&todo.Description,
			&todo.Deadline,
			&todo.Completed,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		)
		if err != nil {
			m.App.ErrorLog.Println(err)
			return nil, err
		}

		todos = append(todos, &todo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return todos, nil
}

func (m *postgresDBRepo) SelectTodo(id int) (*models.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var todo models.Todo

	query := `
select id, name, description, deadline, completed, created_at, updated_at
from todo
where id = $1
`

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&todo.ID,
		&todo.Name,
		&todo.Description,
		&todo.Deadline,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func (m *postgresDBRepo) InsertTodo(todo models.Todo) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
insert into todo (name, description, deadline, completed, created_at, updated_at)
values ($1, $2, $3, $4, $5, $6) returning id
`
	var newID int
	err := m.DB.QueryRowContext(ctx, stmt,
		todo.Name,
		todo.Description,
		todo.Deadline,
		todo.Completed,
		todo.CreatedAt,
		todo.UpdatedAt,
	).Scan(&newID)
	if err != nil {
		return newID, err
	}
	return newID, nil
}
