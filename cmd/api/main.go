package main

import (
	"fmt"
	"github.com/anras5/todo-app-backend/internal/config"
	"github.com/anras5/todo-app-backend/internal/driver"
	"github.com/anras5/todo-app-backend/internal/handlers"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

const port = ":8080"

var app config.Application
var infoLog *log.Logger
var errorLog *log.Logger

func main() {

	// Load variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Set up loggers
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// Connect to database
	log.Println("Connecting to database")
	db, err := driver.ConnectSQL(fmt.Sprintf("host=localhost port=5432 dbname=todos user=%s password=%s sslmode=%s",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWD"), os.Getenv("SSL_MODE")))
	if err != nil {
		log.Fatal("Cannot connect to database! Dying...")
	} else {
		log.Println("Connected to database")
	}
	defer db.Close()

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	srv := &http.Server{
		Addr:    port,
		Handler: routes(&app),
	}
	app.InfoLog.Println("Listening on port 8080")
	err = srv.ListenAndServe()
	log.Fatal(err)

}
