package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/drako02/url-shortener/models"
	"github.com/drako02/url-shortener/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// func TestURLHandler_Delete(t *testing.T) {
// 	type fields struct {
// 		svc *services.URLService
// 	}
// 	type args struct {
// 		c *gin.Context
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			h := &URLHandler{
// 				svc: tt.fields.svc,
// 			}
// 			h.Delete(tt.args.c)
// 		})
// 	}
// }

type mockURLManagementService struct {
	mock.Mock
}

func (s *mockURLManagementService) DeleteURL(id uint) (models.URL, error) {
	args := s.Called(id)
	return args.Get(0).(models.URL), args.Error(1)
}

func (s *mockURLManagementService) SetUrlActiveStatus(ctx context.Context, id uint, value bool) error {
	args := s.Called(ctx, id, value)
	return args.Error(0)
}

func TestURLHandler_SetActiveState(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type fields struct {
		svc services.URLManagementService
	}

	type args struct {
		requestBody map[string]any
	}

	tests := []struct {
		name               string
		fields             fields
		arg                args
		setupMock          func(mockSvc *mockURLManagementService)
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:   "Successfully activate URL",
			fields: fields{svc: &mockURLManagementService{}},
			arg:    args{requestBody: map[string]any{"id": 123, "value": true}},
			setupMock: func(mockSvc *mockURLManagementService) {
				mockSvc.On("SetUrlActiveStatus", mock.Anything, uint(123), true).Return(nil)
			},
			expectedStatusCode: 200,
			expectedBody:       `{"message":"URL status updated"}`,
		},
		{
			name: "Fail to update - service error",
			fields: fields{
				svc: &mockURLManagementService{},
			},
			arg: args{requestBody: map[string]any{"id": 2, "value": false}},
			setupMock: func(mockSvc *mockURLManagementService) {
				mockSvc.On("SetUrlActiveStatus", mock.Anything, uint(2), false).Return(
					errors.New("database error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       `{"error":"Failed to update URL status. Something went wrong"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			jsonBody, _ := json.Marshal(tt.arg.requestBody)
			fmt.Println("JSON body:", string(jsonBody))
			req, _ := http.NewRequest(http.MethodPut, "/activate", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			mockSvc := tt.fields.svc.(*mockURLManagementService)
			if tt.setupMock != nil {
				tt.setupMock(mockSvc)
			}

			hnd := NewURLHandler(mockSvc)
			hnd.UpdateUrlActiveStatus(c)

			assert.Equal(t, tt.expectedStatusCode, w.Code)

			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, w.Body.String())
			}

			mockSvc.AssertExpectations(t)

		})

	}
}
