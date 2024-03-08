package api

import (
	"log"
	"net/http"

	"github.com/arielcr/payment-gateway/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

type Router struct {
	Server *gin.Engine
	Port   string
}

func NewRouter(port string) *Router {
	return &Router{
		Port: port,
	}
}

func (r *Router) InitializeEndpoints() {
	server := gin.Default()

	server.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusText(http.StatusOK),
		})
	})

	explorer := server.Group("/payment")
	{
		explorer.POST("/process", handlers.ProcessPayment)
	}

	r.Server = server
}

func (r *Router) Start() {
	if err := r.Server.Run(r.Port); err != nil {
		log.Fatalln("error when server is initializing")
	}
}
