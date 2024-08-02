package dbrepo

import (
	"database/sql"

	"github.com/anras5/todo-app-backend/internal/repository"
)

type postgresDBRepo struct {
	DB *sql.DB
}

type testDBRepo struct {
	DB *sql.DB
}

func NewPostgresRepo(conn *sql.DB) repository.DatabaseRepo {
	return &postgresDBRepo{
		DB: conn,
	}
}

func NewTestingRepo() repository.DatabaseRepo {
	return &testDBRepo{}
}
