package dbrepo

import (
	"context"
	"database/sql"
	"github.com/anras5/todo-app-backend/internal/models"
	"time"
)

func (m *postgresDBRepo) SelectTodos(completed ...bool) ([]*models.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var todos []*models.Todo
	var query string
	var rows *sql.Rows
	var err error

	if len(completed) > 0 {
		// we want to filter the data
		query = `
SELECT ID, NAME, DESCRIPTION, DEADLINE, COMPLETED, CREATED_AT, UPDATED_AT
FROM TODO WHERE COMPLETED = $1 ORDER BY DEADLINE
`
		rows, err = m.DB.QueryContext(ctx, query, completed[0])
	} else {
		// we want to get all the data
		query = `
SELECT ID, NAME, DESCRIPTION, DEADLINE, COMPLETED, CREATED_AT, UPDATED_AT
FROM TODO ORDER BY DEADLINE
`
		rows, err = m.DB.QueryContext(ctx, query)
	}
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
		time.Now(),
		time.Now(),
	).Scan(&newID)
	if err != nil {
		return newID, err
	}
	return newID, nil
}

func (m *postgresDBRepo) UpdateTodo(todo models.Todo) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
update todo set name = $1, description = $2, deadline = $3, completed = $4, updated_at = $5
where id = $6
`
	_, err := m.DB.ExecContext(ctx, stmt,
		todo.Name,
		todo.Description,
		todo.Deadline,
		todo.Completed,
		time.Now(),
		todo.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) UpdateTodoCompleted(id int, completed bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
update todo set completed = $1, updated_at = $2 where id = $3
`
	_, err := m.DB.ExecContext(ctx, stmt, completed, time.Now(), id)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) DeleteTodo(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
delete from todo where id = $1
`

	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}
	return nil
}
