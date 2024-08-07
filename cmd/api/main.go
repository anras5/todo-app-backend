package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/anras5/todo-app-backend/internal/config"
	"github.com/anras5/todo-app-backend/internal/driver"
	rpc "github.com/anras5/todo-app-backend/internal/grpc"
	"github.com/anras5/todo-app-backend/internal/handlers"
)

const port = ":8080"

var app config.Application
var infoLog *log.Logger
var errorLog *log.Logger

func main() {

	// -------------------------------------------------------------------------------------------- //
	// Set up loggers
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// -------------------------------------------------------------------------------------------- //
	// Connect to database
	app.InfoLog.Println("Connecting to database on port 5432")
	db, err := driver.ConnectSQL(fmt.Sprintf("host=postgres-db port=5432 dbname=todos user=%s password=%s sslmode=%s",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWD"), os.Getenv("SSL_MODE")))
	if err != nil {
		log.Fatal("Cannot connect to database! Dying...")
	} else {
		log.Println("Connected to database")
	}
	defer db.Close()

	// -------------------------------------------------------------------------------------------- //
	// set repo and handlers
	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	// -------------------------------------------------------------------------------------------- //
	// Start gRPC server
	grpcServer := rpc.NewTodoServer(db)
	go grpcServer.Run()

	srv := &http.Server{
		Addr:    port,
		Handler: routes(),
	}
	app.InfoLog.Println("Listening on port 8080")
	err = srv.ListenAndServe()
	log.Fatal(err)
}
