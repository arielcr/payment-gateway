// Package api provides functionality for setting up and managing API endpoints.
package api

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/arielcr/payment-gateway/internal/api/handlers"
	"github.com/arielcr/payment-gateway/internal/api/middleware"
	"github.com/arielcr/payment-gateway/internal/config"
	"github.com/gin-gonic/gin"
)

// Router manages the HTTP routes and handlers for the application.
type Router struct {
	Server         *gin.Engine
	logger         *slog.Logger
	Config         config.Application
	PaymentHandler *handlers.PaymentHandler
	RefundHandler  *handlers.RefundHandler
}

// NewRouter creates a new instance of Router with the provided port, payment handler, and refund handler.
func NewRouter(config config.Application, paymentHandler *handlers.PaymentHandler, refundHandler *handlers.RefundHandler, logger *slog.Logger) *Router {
	return &Router{
		Config:         config,
		PaymentHandler: paymentHandler,
		RefundHandler:  refundHandler,
		logger:         logger,
	}
}

// InitializeEndpoints sets up the HTTP routes and handlers for the server.
func (r *Router) InitializeEndpoints() {
	r.logger.Info("Initializing endpoints")

	server := gin.Default()

	// Add authentication middleware
	server.Use(middleware.Authenticate(r.Config))

	// Health check endpoint
	server.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusText(http.StatusOK),
		})
	})

	// Merchant endpoints for processing payments and refunds
	merchants := server.Group("/merchants")
	{
		merchants.POST("/payment/process", r.PaymentHandler.ProcessPayment)
		merchants.POST("/payment/:paymentID/refund", r.RefundHandler.RefundPayment)
	}

	// Payment endpoint for retrieving payment information by ID
	payments := server.Group("/payments")
	{
		payments.GET("/:paymentID", r.PaymentHandler.GetPayment)
	}

	r.Server = server
}

// Start starts the HTTP server on the specified port.
func (r *Router) Start() {
	r.logger.Info("Starting the web server")

	if err := r.Server.Run(r.Config.ApplicationPort); err != nil {
		log.Fatalln("error when server is initializing")
	}
}
