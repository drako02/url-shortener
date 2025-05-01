package main

import (
	"strings"
	"time"

	config "github.com/drako02/url-shortener/config"
	// "github.com/drako02/url-shortener/middlewares"
	"github.com/drako02/url-shortener/repositories"
	routes "github.com/drako02/url-shortener/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDatabase()
	config.InitFirebase()
	config.InitFirebaseAuth()
	repositories.Migrate()
	config.InitKafkaProducer()

	defer config.KafkaProducer.Close()

	allowedSuffixes := []string{"vercel.app"}
	exactOrigins := []string{"http://localhost:3000"}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		AllowOriginFunc: func(origin string) bool {
			for _, allowedSuffix := range allowedSuffixes {
				if strings.HasSuffix(origin, allowedSuffix) {
					return true
				}
			}

			for _, exactOrigin := range exactOrigins {
				if exactOrigin == origin {
					return true
				}
			}
			return false
		},
	}))
	routes.SetUpRoutes(r)
	r.Run()
}