package api

import (
	"log"
	"net/http"

	"github.com/arielcr/payment-gateway/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

type Router struct {
	Server         *gin.Engine
	Port           string
	PaymentHandler *handlers.PaymentHandler
}

func NewRouter(port string, paymentHandler *handlers.PaymentHandler) *Router {
	return &Router{
		Port:           port,
		PaymentHandler: paymentHandler,
	}
}

func (r *Router) InitializeEndpoints() {
	server := gin.Default()

	server.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusText(http.StatusOK),
		})
	})

	merchants := server.Group("/merchants")
	{
		merchants.POST("/payment/process", r.PaymentHandler.ProcessPayment)
	}

	payments := server.Group("/payments")
	{
		payments.GET("/:paymentID", r.PaymentHandler.GetPayment)
	}

	r.Server = server
}

func (r *Router) Start() {
	if err := r.Server.Run(r.Port); err != nil {
		log.Fatalln("error when server is initializing")
	}
}
