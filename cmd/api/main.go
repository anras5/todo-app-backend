package main

import (
	"github.com/anras5/todo-app-backend/internal/config"
	"github.com/anras5/todo-app-backend/internal/handlers"
	"log"
	"net/http"
	"os"
)

const port = ":8080"

var app config.Application
var infoLog *log.Logger
var errorLog *log.Logger

func main() {

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	srv := &http.Server{
		Addr:    port,
		Handler: routes(&app),
	}
	app.InfoLog.Println("Listening on port 8080")
	err := srv.ListenAndServe()
	log.Fatal(err)

}
