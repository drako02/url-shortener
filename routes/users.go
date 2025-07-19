package routes

import (
	"github.com/drako02/url-shortener/handlers"
	"github.com/drako02/url-shortener/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine, handler *handlers.UserHandler) {
	users := r.Group("/users")

	{
		users.POST("/exists", handlers.UserExists)
		users.POST("", handlers.CreateUser)
	}

	protected := r.Group("/")

	protected.Use(middlewares.IsAuthenticated(), middlewares.LogRequest())
	{
		protected.GET("/users/:uid", handlers.GetUser)
		protected.PATCH("/user/", handler.Update) // TODO Test this endpoint
	}

}
