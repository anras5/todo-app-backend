package dbrepo

import (
	"context"
	"github.com/anras5/todo-app-backend/internal/models"
	"time"
)

func (m *postgresDBRepo) AllTodos() ([]*models.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var todos []*models.Todo

	query := `
SELECT ID, NAME, DESCRIPTION, DEADLINE, COMPLETED, CREATED_AT, UPDATED_AT
FROM TODO
`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		m.App.ErrorLog.Println("postgres.go: Error executing QueryContext in AllTodos")
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
			m.App.ErrorLog.Println("postgres.go: Error scanning rows in AllTodos")
			return nil, err
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return todos, nil
}
