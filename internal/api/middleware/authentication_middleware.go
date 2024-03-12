package middleware

import (
	"net/http"
	"strings"

	"github.com/arielcr/payment-gateway/internal/config"
	"github.com/arielcr/payment-gateway/internal/utils"
	"github.com/gin-gonic/gin"
)

func Authenticate(config config.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		headerParts := strings.Split(authHeader, " ")
		clientToken := headerParts[1]
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No Authorization Header Provided"})
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(clientToken, config)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}
		c.Set("email", claims.Email)
		c.Set("first_name", claims.First_name)
		c.Set("last_name", claims.Last_name)
		c.Set("uid", claims.Uid)
		c.Set("user_type", claims.User_type)
		c.Next()
	}
}
