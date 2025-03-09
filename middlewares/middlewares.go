package middlewares

import (
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
		if len(receivedToken) > 7 && receivedToken[:7] == "Bearer " {
			log.Printf("receivedToken: %s \n receivedToken[:7]: %s",receivedToken, receivedToken[:7])
            receivedToken = receivedToken[7:]
        }
		token, err := config.FirebaseAuth.VerifyIDToken(c.Request.Context(), receivedToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
			log.Println("Invalid Token")
			return
		}
		
        c.Set("uid", token.UID)
        log.Printf("User %s authorized", token.UID)
        
		c.Next()
	}
}
