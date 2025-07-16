package handlers

import (
	"log"
	"net/http"

	"github.com/drako02/url-shortener/services"
	"github.com/drako02/url-shortener/utils"
	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	var request utils.CreateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := services.CreateUser(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, res)
}

func GetUser(c *gin.Context) {
	uid := c.Param("uid")
	if uid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "uid is required"})
	}

	res, err := services.GetUser(uid)
	if err != nil {
		status := http.StatusNotFound
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func UserExists(c *gin.Context) {
	var request utils.ExistsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exists, err := services.UserExists(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]bool{"result": exists})
}

type UserHandler struct {
	Svc *services.UserService
}

func NewUserHandler(svc *services.UserService) *UserHandler {
	return &UserHandler{svc}
}

func (hnd *UserHandler) Update(c *gin.Context) {
	//TODO continue implementation
	// TODO look  into getting the user's id from the token in the middleware
	request := services.ValidUserFields{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		log.Printf("%s ")
		return
	}
}
