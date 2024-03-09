// Package application provides functionality for managing the application lifecycle, including server setup and initialization.
// It includes methods for starting the server, loading configuration, initializing logger, storage, and router.
package application

import (
	"errors"
	"log"
	"log/slog"
	"os"

	"github.com/arielcr/payment-gateway/internal/api"
	"github.com/arielcr/payment-gateway/internal/api/handlers"
	"github.com/arielcr/payment-gateway/internal/config"
	"github.com/arielcr/payment-gateway/internal/storage"
)

// Setup contains application metadata
type Setup struct {
	Version    string
	BuildDate  string
	CommitHash string
}

// Server is the server of our application.
type Server struct {
	logger     *slog.Logger
	store      storage.Repository
	router     *api.Router
	config     config.Application
	version    string
	buildDate  string
	commitHash string
}

var (
	errStartingApplication = errors.New("unable to start application")
)

// NewServer creates a new instance of Server using the provided setup settings.
func NewServer(settings Setup) *Server {
	newServer := Server{
		version:    settings.Version,
		buildDate:  settings.BuildDate,
		commitHash: settings.CommitHash,
	}

	return &newServer
}

// Run starts the server by performing initialization steps such as loading configuration,
// initializing logger, storage, and router, and then starting the router to listen for incoming requests.
// It returns an error if any initialization step fails.
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

	slog.Info("initializing storage")
	storageError := s.initializeStorage()
	if storageError != nil {
		return errStartingApplication
	}

	slog.Info("initializing router")
	s.initializeRouter()
	s.router.Start()

	return nil
}

// loadConfiguration loads the application configuration using the config package.
// It returns an error if the configuration cannot be loaded.
func (s *Server) loadConfiguration() error {
	applicationConfig, err := config.Load()
	if err != nil {
		log.Println("level", "ERROR", "msg", "application setup could not be loaded", "error", err)

		return errors.New("application setup could not be loaded")
	}
	s.config = applicationConfig
	return nil
}

// initializeLogger initializes the logger based on the configuration settings.
// It sets the log level and creates a logger instance using the slog package.
// Returns an error if logger initialization fails.
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

// initializeStorage initializes the database connection and creates a storage instance.
// It uses the storage package to connect to MySQL and create a MySQLRepository.
// Returns an error if storage initialization fails.
func (s *Server) initializeStorage() error {
	db, err := storage.ConnectMySQL(s.config.Repository)
	if err != nil {
		return err
	}
	s.store = storage.NewMySQLRepository(db)
	return nil
}

// initializeRouter initializes the router with payment and refund handlers,
// configures the endpoints, and assigns the router to the server.
func (s *Server) initializeRouter() {
	paymentHandler := handlers.NewPaymentHandler(s.store, s.config)
	refundHandler := handlers.NewRefundHandler(s.store, s.config)
	router := api.NewRouter(s.config.ApplicationPort, paymentHandler, refundHandler)
	router.InitializeEndpoints()
	s.router = router
}
