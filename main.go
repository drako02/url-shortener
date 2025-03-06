package main

import (
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
	r.Use(cors.Default())
	routes.SetUpRoutes(r)
	r.Run()
}
