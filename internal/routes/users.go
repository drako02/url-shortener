package routes

import (
	"net/http"

	"github.com/drako02/url-shortener/internal/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine) {
	users := r.Group("/users")
	{
		users.GET("/:uid", getUser)
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

func getUser(c *gin.Context) {
	uid := c.Param("uid")
	if uid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "uid is required"})
	}
	// var request handlers.GetUserRequest
	// if err := c.ShouldBindJSON(&request); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": err.Error(),
	// 	})
	// 	return
	// }
	res, err := handlers.GetUser(uid)
	if err != nil {
		status := http.StatusNotFound
		if err.Error() == "use not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
