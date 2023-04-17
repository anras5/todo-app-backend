package handlers

import (
	"github.com/anras5/todo-app-backend/internal/config"
	"log"
	"os"
	"testing"
)

var app config.Application

func TestMain(m *testing.M) {

	// -------------------------------------------------------------------------------------------- //
	// Set up loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// -------------------------------------------------------------------------------------------- //
	// set repo and handlers
	repo := NewTestRepo(&app)
	NewHandlers(repo)

	os.Exit(m.Run())
}
