package routes

import (
	"net/http"

	"github.com/drako02/url-shortener/internal/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine) {
	users := r.Group("/users")
	{
		users.GET("", getUsers)
		users.POST("", createUser)
	}

}

func createUser(c *gin.Context) {
	var request handlers.CreateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res := handlers.CreateUser(request)
	c.JSON(http.StatusCreated, res)
}

func getUsers(c *gin.Context) {
    // TODO: Implement get users handler
    c.JSON(http.StatusOK, gin.H{"message": "Not implemented"})
}
