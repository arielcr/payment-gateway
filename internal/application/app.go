package application

import (
	"errors"
	"log"
	"log/slog"
	"os"

	"github.com/arielcr/payment-gateway/internal/api"
	"github.com/arielcr/payment-gateway/internal/config"
)

// Setup contains application metadata
type Setup struct {
	Version    string
	BuildDate  string
	CommitHash string
}

// Server is the server of our application.
type Server struct {
	logger *slog.Logger
	// store      *stores.Store
	router     *api.Router
	config     config.Application
	version    string
	buildDate  string
	commitHash string
}

var (
	errStartingApplication = errors.New("unable to start application")
)

func NewServer(settings Setup) *Server {
	newServer := Server{
		version:    settings.Version,
		buildDate:  settings.BuildDate,
		commitHash: settings.CommitHash,
	}

	return &newServer
}

func (s *Server) Run() error {
	slog.Info("loading configuration")
	confError := s.loadConfiguration()
	if confError != nil {
		return errStartingApplication
	}

	slog.Info("initializing logger")
	loggerError := s.initializeLogger()
	if loggerError != nil {
		return errStartingApplication
	}

	slog.Info("initializing router")
	s.initializeRouter()
	s.router.Start()

	return nil
}

func (s *Server) loadConfiguration() error {
	applicationConfig, err := config.Load()
	if err != nil {
		log.Println("level", "ERROR", "msg", "application setup could not be loaded", "error", err)

		return errors.New("application setup could not be loaded")
	}
	s.config = applicationConfig
	return nil
}

func (s *Server) initializeLogger() error {
	logLevel := slog.LevelDebug

	if s.config.LogLevel == config.ProductionLog {
		slog.Info("setting log level", slog.String("level", "INFO"))
		logLevel = slog.LevelInfo
	}

	handlerOptions := &slog.HandlerOptions{
		Level: logLevel,
	}

	loggerHandler := slog.NewJSONHandler(os.Stdout, handlerOptions)
	logger := slog.New(loggerHandler)

	logger.Info(
		"logger has been initialized",
		slog.String("level", handlerOptions.Level.Level().String()),
	)

	slog.SetDefault(logger)

	s.logger = logger

	return nil
}

func (s *Server) initializeRouter() {
	router := api.NewRouter(s.config.ApplicationPort)
	router.InitializeEndpoints()
	s.router = router
}
