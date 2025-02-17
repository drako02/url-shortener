package routes

import (
	"fmt"
	"net/http"

	"github.com/drako02/url-shortener/internal/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine){
	r.POST("/create", create)
	r.GET("/:shortCode", handleRedirect)
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

func handleRedirect(c *gin.Context){
	shortCode := c.Param("shortCode");
	fmt.Println(shortCode)
	longUrl, err := handlers.GetLongUrl(shortCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Redirect(http.StatusFound, longUrl)
}