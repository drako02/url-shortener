package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	// "github.com/drako02/url-shortener/models"
	"github.com/drako02/url-shortener/repositories"
	"github.com/drako02/url-shortener/services"
	"github.com/drako02/url-shortener/utils"
	"github.com/gin-gonic/gin"
)

var Loglevel = struct {
	Error  string
	Log    string
	Info   string
	Debug  string
	Warning string
}{
	"[ERROR]",
	"[LOG]",
	"[INFO]",
	"[DEBUG]",
	"[WARNING]",
}

func Create(c *gin.Context) {
	var request utils.CreateRequest

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

// FIXME redirection time does not represent time it took accurately, time also not in seconds
func (h *URLHandler) HandleRedirect(c *gin.Context) {
	start := time.Now()
	shortcode := c.Param("shortCode")

	info := utils.RedirectInfo{
		Timestamp:      time.Now(),
		ShortCode:      shortcode,
		ClientIP:       c.ClientIP(),
		UserAgent:      c.Request.UserAgent(),
		Referer:        c.Request.Referer(),
		AcceptLanguage: c.GetHeader("Accept-Language"),
		StatusCode:     http.StatusFound,
		DurationMs:     int64(time.Since(start)),
	}

	infoJson, err := json.Marshal(info)
	if err != nil {
		log.Printf("%s Failed to marshal RedirectInfo: %v", Loglevel.Error, err)
		//Handle case later
	}

	infoStr := string(infoJson)

	log.Printf("%s Redirect Info: %+v", Loglevel.Info, infoStr)

	fmt.Println(shortcode)
	isActive, err := h.svc.URLIsActive(c.Request.Context(), nil, &shortcode)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "could not resolve url state"})
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Title":   "Server Error",
			"Message": "Sorry, we couldn’t resolve the URL state right now.",
		})
		return
	}

	if !isActive {
		// c.JSON(http.StatusBadRequest, gin.H{"error": "URL is not active"})
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Title":   "Bad Request",
			"Message": "URL has been deactivated by owner",
		})
		return
	}

	longUrl, err := services.GetLongUrl(shortcode)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "short code not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve URL"})
		}
		return
	}

	c.Redirect(http.StatusFound, longUrl)
	services.WriteKafkaEvent("redirections", "shortCode", infoStr)
}

func GetUserUrls(c *gin.Context) {
	var request utils.GetUserUrlRequest
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
	var request utils.UrlQuery
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
		log.Printf("%s Failed to bind JSON for UpdateUrlActiveStatus: %v", Loglevel.Error, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := h.svc.SetUrlActiveStatus(c.Request.Context(), req.Id, *req.Value)
	if err != nil {
		log.Printf("%s Failed to set active status of url with id %d: %v", Loglevel.Error, req.Id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update URL status. Something went wrong"})
		return
	}

	log.Printf("%s URL status updated successfully for ID %d", Loglevel.Info, req.Id)
	c.JSON(http.StatusOK, gin.H{"message": "URL status updated"})
}
