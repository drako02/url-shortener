
package main
import (
	"github.com/gin-gonic/gin"
	"github.com/drako02/url-shortener/internal/routes"
	"github.com/drako02/url-shortener/config"
	"github.com/drako02/url-shortener/internal/repositories"
)

func main () {
	config.ConnectDatabase()
	repositories.Migrate()

	r := gin.Default();
	// r.GET("/ping", func( c *gin.Context){
	// 	c.JSON(200, gin.H{
	// 		"message": "pong",
	// 	})
	// })
	routes.RegisterRoutes(r)
	r.Run()
}


