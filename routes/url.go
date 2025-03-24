package routes

import (
	"fmt"
	"net/http"

	"github.com/drako02/url-shortener/handlers"
	"github.com/gin-gonic/gin"
	"log"
)

func RegisterUrlRoutes(r *gin.Engine) {
	r.POST("/create", create)
	r.GET("/:shortCode", handleRedirect)
	r.POST("/user-urls", getUserUrls)
}

func create(c *gin.Context) {
	var request handlers.CreateRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := handlers.CreateShortUrl(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Location", fmt.Sprintf("/%s", result["short_code"]))
	c.JSON(http.StatusCreated, result)
}

func handleRedirect(c *gin.Context) {
	shortCode := c.Param("shortCode")
	fmt.Println(shortCode)
	longUrl, err := handlers.GetLongUrl(shortCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Redirect(http.StatusFound, longUrl)
	handlers.WriteKafkaEvent("redirections", "shortCode", shortCode)
}

func getUserUrls(c *gin.Context) {
	var request handlers.GetUserUrlRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}
	res, err := handlers.GetUserUrls(request.UID, request.Limit, request.Offset)
	if err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user urls "})
		return
	}

	recordCount, countError := handlers.GetTotalUrls(request.UID)
	if countError != nil {
		log.Print(countError.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error" : "Failed to get total url count for user"})
		return
	}


	c.JSON(http.StatusOK, gin.H{"urls": res, "recordCount": recordCount})
}
