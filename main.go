package main

import (
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

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000",}, 
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	routes.SetUpRoutes(r)
	r.Run()
}
