package main

import (
	"github.com/drako02/url-shortener/config"
	"github.com/drako02/url-shortener/internal/repositories"
	"github.com/drako02/url-shortener/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func main() {
	config.ConnectDatabase()
	repositories.Migrate()

	r := gin.Default()
	r.Use(cors.Default())
	routes.SetUpRoutes(r)
	r.Run()
}
