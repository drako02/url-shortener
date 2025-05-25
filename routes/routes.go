package routes

import (
	"github.com/drako02/url-shortener/handlers"
	"github.com/gin-gonic/gin"
)

type AppHandlers struct {
	URLHandler *handlers.URLHandler
	// Update with other handler layers - eg userhandler
}

func NewAppHandler(urlHnd *handlers.URLHandler) *AppHandlers{
	return &AppHandlers{ URLHandler: urlHnd}
}

func SetUpRoutes(r *gin.Engine, appHandlers *AppHandlers) {

	if appHandlers.URLHandler != nil {
		RegisterUrlRoutes(r, appHandlers.URLHandler)

	}

	// handle later to cater for user handlers
	RegisterUserRoutes(r)
}
