package main

import (
	"log/slog"
	"os"

	"github.com/arielcr/payment-gateway/internal/application"
)

var (
	Version    string
	BuildDate  string
	CommitHash string
)

func main() {
	slog.Info("starting application")

	app := newApplicationServer()

	if err := app.Run(); err != nil {
		slog.Error("unable to start service: %s", err)
		os.Exit(-1)
	}

	slog.Info("finishing application")
}

func newApplicationServer() *application.Server {
	settings := application.Setup{
		Version:    Version,
		BuildDate:  BuildDate,
		CommitHash: CommitHash,
	}

	return application.NewServer(settings)
}
