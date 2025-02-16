package routes

import (
	"net/http"

	"github.com/drako02/url-shortener/internal/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine){
	r.POST("/create", create)
}

func create(c *gin.Context){
	var request handlers.CreateRequest

	if err := c.ShouldBindJSON(&request); err!= nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result := handlers.CreateShortUrl(request)

	c.JSON(http.StatusCreated, result)
}