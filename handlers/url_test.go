package handlers

import (
	"testing"

	"github.com/drako02/url-shortener/services"
	"github.com/gin-gonic/gin"
)

func TestURLHandler_Delete(t *testing.T) {
	type fields struct {
		svc *services.URLService
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &URLHandler{
				svc: tt.fields.svc,
			}
			h.Delete(tt.args.c)
		})
	}
}
