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
	uid, exists := c.Get("uid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "could not authenticate user"})
		return
	}

	typedUid, ok := uid.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "could not authenticate user"})
		return
	}

	id, getIdError := services.GetIdFromUid(typedUid)
	if getIdError != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	request := services.ValidUserFields{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		log.Printf("%s Error binding JSON in Update for user with id %d: %v ", Loglevel.Error, id, err)
		return
	}

	user, err := hnd.Svc.UpdateUserInfo(id, request, c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user details"})
		log.Printf("%s Error updating user with id %d: %v", Loglevel.Error, id, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user updated successfully", "data": *user})
}
