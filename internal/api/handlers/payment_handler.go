package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ProcessPayment(context *gin.Context) {

	context.JSON(http.StatusAccepted, gin.H{"payment": "accepted"})
}
