package dbrepo

import (
	"database/sql"
	"github.com/anras5/todo-app-backend/internal/config"
	"github.com/anras5/todo-app-backend/internal/repository"
)

type postgresDBRepo struct {
	App *config.Application
	DB  *sql.DB
}

type testDBRepo struct {
	App *config.Application
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, app *config.Application) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: app,
		DB:  conn,
	}
}

func NewTestingRepo(a *config.Application) repository.DatabaseRepo {
	return &testDBRepo{
		App: a,
	}
}
