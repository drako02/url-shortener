package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	// "github.com/drako02/url-shortener/models"
	"github.com/drako02/url-shortener/repositories"
	"github.com/drako02/url-shortener/services"
	"github.com/gin-gonic/gin"
)

var loglevel = struct {
	error  string
	log    string
	info   string
	debug  string
	waring string
}{
	"[ERROR]",
	"[LOG]",
	"[INFO]",
	"[DEBUG]",
	"[WARNING]",
}

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

func (h *URLHandler) HandleRedirect(c *gin.Context) {
	shortCode := c.Param("shortCode")
	fmt.Println(shortCode)
	isActive, err := h.svc.URLIsActive(c.Request.Context(), nil, &shortCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not resolve url state"})
		return
	}

	if !isActive {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL is not active"})
		return
	}

	longUrl, err := services.GetLongUrl(shortCode)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "short code not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve URL"})
		}
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

func QueryUrls(c *gin.Context) {
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

	c.JSON(http.StatusOK, gin.H{"urls": urls, "length": count})
}

type URLHandler struct {
	svc services.URLManagementService
}

type URLHandlerInterface interface {
	Delete(c *gin.Context)
}

func NewURLHandler(svc services.URLManagementService) *URLHandler {
	return &URLHandler{svc: svc}
}

func (h *URLHandler) Delete(c *gin.Context) {
	var request struct {
		Id uint `json:"id"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("%v, id: %d", err, request.Id)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	url, err := h.svc.DeleteURL(request.Id)
	if err != nil {
		log.Printf("Failed to delete URL with id %d: %v", request.Id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete URL"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"url": url})
}

type urlStatusChangeRequest struct {
	Id    uint  `json:"id" binding:"required"`
	Value *bool `json:"value" binding:"required"`
}

func (h *URLHandler) UpdateUrlActiveStatus(c *gin.Context) {
	var req urlStatusChangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("%s Failed to bind JSON for UpdateUrlActiveStatus: %v", loglevel.error, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := h.svc.SetUrlActiveStatus(c.Request.Context(), req.Id, *req.Value)
	if err != nil {
		log.Printf("%s Failed to set active status of url with id %d: %v", loglevel.error, req.Id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update URL status. Something went wrong"})
		return
	}

	log.Printf("%s URL status updated successfully for ID %d", loglevel.info, req.Id)
	c.JSON(http.StatusOK, gin.H{"message": "URL status updated"})
}
