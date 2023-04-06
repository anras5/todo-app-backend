package config

import "log"

// Application holds the application config
type Application struct {
	Domain   string
	DSN      string
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}
