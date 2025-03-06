package middlewares

import (
	"fmt"
	"log"
	"net/http"

	"github.com/drako02/url-shortener/config"
	"github.com/gin-gonic/gin"
)

func IsAuthenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestHeader := c.Request.Header
		receivedToken := requestHeader.Get("Authorization")
		if receivedToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No auth token provided"})
			log.Println("No auth token provided")
			return
		}
		token, err := config.FirebaseAuth.VerifyIDToken(c.Request.Context(), receivedToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
			log.Println("Invalid Token")
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": fmt.Sprintf("user %s authorized", token.UID)})
		c.Next()
	}
}
