package main

import (
	"github.com/anras5/todo-app-backend/internal/config"
	"github.com/anras5/todo-app-backend/internal/handlers"
	"log"
	"net/http"
)

const port = ":8080"

var app config.Application

func main() {

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	srv := &http.Server{
		Addr:    port,
		Handler: routes(&app),
	}
	err := srv.ListenAndServe()
	log.Fatal(err)

}
