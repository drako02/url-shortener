package main

import (
	"github.com/drako02/url-shortener/config"
	"github.com/drako02/url-shortener/internal/repositories"
	"github.com/drako02/url-shortener/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDatabase()
	repositories.Migrate()

	r := gin.Default()
	routes.SetUpRoutes(r)
	r.Run()
}
