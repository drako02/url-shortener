package routes

import (
	"net/http"

	"github.com/drako02/url-shortener/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine) {
	users := r.Group("/users")

	{
		users.GET("/:uid", getUser)
		users.POST("/exists", userExists)
		users.POST("", createUser)
	}

}

func createUser(c *gin.Context) {
	var request handlers.CreateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := handlers.CreateUser(request)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
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
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func userExists(c *gin.Context) {
	var request handlers.ExistsRequest
	if err := c.ShouldBindJSON(&request); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exists, err := handlers.UserExists(request)
	if(err != nil){
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]bool{"result": exists})
}
