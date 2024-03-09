// Package main serves as the entry point for the payment gateway application.
// It initializes the application server with version information and starts the application.
package main

import (
	"log/slog"
	"os"

	"github.com/arielcr/payment-gateway/internal/application"
)

// Version contains the version number of the application.
var (
	Version    string
	BuildDate  string
	CommitHash string
)

// main is the entry point of the application.
// It initializes the application server, starts it, and handles any errors that occur during startup.
func main() {
	slog.Info("starting application")

	app := newApplicationServer()

	if err := app.Run(); err != nil {
		slog.Error("unable to start service: %s", err)
		os.Exit(-1)
	}

	slog.Info("finishing application")
}

// newApplicationServer creates a new instance of the application server with the provided version information.
func newApplicationServer() *application.Server {
	settings := application.Setup{
		Version:    Version,
		BuildDate:  BuildDate,
		CommitHash: CommitHash,
	}

	return application.NewServer(settings)
}
