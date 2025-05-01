package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/drako02/url-shortener/services"
	"github.com/gin-gonic/gin"
)


func Create(c *gin.Context) {
	var request services.CreateRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := services.CreateShortUrl(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Location", fmt.Sprintf("/%s", result["short_code"]))
	c.JSON(http.StatusCreated, result)
}

func HandleRedirect(c *gin.Context) {
	shortCode := c.Param("shortCode")
	fmt.Println(shortCode)
	longUrl, err := services.GetLongUrl(shortCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Redirect(http.StatusFound, longUrl)
	services.WriteKafkaEvent("redirections", "shortCode", shortCode)
}

func GetUserUrls(c *gin.Context) {
	var request services.GetUserUrlRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}
	res, err := services.GetUserUrls(request.UID, request.Limit, request.Offset)
	if err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user urls "})
		return
	}

	recordCount, countError := services.GetTotalUrls(request.UID)
	if countError != nil {
		log.Print(countError.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get total url count for user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"urls": res, "recordCount": recordCount})
}

func QueryUrls(c *gin.Context){
	var request services.UrlQuery
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	log.Printf("request: %v", request)

	urls, count, err := services.UrlQueryResult(request)
	if err != nil {
		log.Printf("Failed to query urls for user %s: %v", request.UID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"urls": urls, "length" : count})
}



