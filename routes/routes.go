package routes

import (
	"github.com/drako02/url-shortener/handlers"
	"github.com/gin-gonic/gin"
)

type AppHandlers struct {
	URLHandler *handlers.URLHandler
	UserHandler *handlers.UserHandler
	// Update with other handler layers - eg userhandler
}

func NewAppHandler(urlHnd *handlers.URLHandler, userHnd *handlers.UserHandler) *AppHandlers{
	return &AppHandlers{ URLHandler: urlHnd, UserHandler: userHnd}
}

func SetUpRoutes(r *gin.Engine, appHandlers *AppHandlers) {

	if appHandlers.URLHandler != nil {
		RegisterUrlRoutes(r, appHandlers.URLHandler)

	}

	if appHandlers.UserHandler != nil {
		RegisterUserRoutes(r, appHandlers.UserHandler)
	}

	// handle later to cater for user handlers
	// RegisterUserRoutes(r)
}
