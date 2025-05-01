package routes

import (
	"github.com/drako02/url-shortener/handlers"
	"github.com/drako02/url-shortener/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterUrlRoutes(r *gin.Engine) {
	urls := r.Group("/")
	urls.Use(middlewares.IsAuthenticated())

	{
		urls.POST("/create", handlers.Create)
		urls.POST("/user-urls", handlers.GetUserUrls)
		urls.POST("/urls", handlers.QueryUrls)
	}
	r.GET("/:shortCode", handlers.HandleRedirect)

}
