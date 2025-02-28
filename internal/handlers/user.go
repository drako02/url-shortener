package handlers

import (
	"fmt"
	"time"

	"github.com/drako02/url-shortener/config"
	"github.com/drako02/url-shortener/internal/models"
)

type CreateUserRequest struct {
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	UID       string  `json:"uid" binding:"required"`
}

type GetUserRequest struct {
	UID string `json:"uid" binding:"required"`
}

func CreateUser(payload CreateUserRequest) map[string]any {
	db := config.DB

	var firstName, lastName string
	if payload.FirstName != nil {
		firstName = *payload.FirstName
	}
	if payload.LastName != nil {
		lastName = *payload.LastName
	}
	user := models.User{
		FirstName: firstName,
		LastName:  lastName,
		UID:       payload.UID,
		JoinedAt:  time.Now(),
	}

	result := db.Create(&user)
	if result.Error != nil {
		return map[string]any{
			"error": result.Error.Error(),
		}
	}
	return map[string]any{
		"ID":  user.ID,
		"UID": payload.UID,
	}

}

func GetUser(uid string) (map[string]any, error) {
	var user models.User
	db := config.DB
	result := db.Where("uid=?", uid).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return map[string]any{
		"id": user.ID,
		"uid": user.UID,
		"firstName": user.FirstName,
		"lastName": user.LastName,
		"joinedAt": user.JoinedAt,
	}, nil

}
