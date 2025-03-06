package handlers

import (
	"fmt"
	"time"

	"github.com/drako02/url-shortener/config"
	"github.com/drako02/url-shortener/models"
	"gorm.io/gorm"
)

type CreateUserRequest struct {
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	UID       string  `json:"uid" binding:"required"`
	Email string `json:"email" binding:"required"`
}

type GetUserRequest struct {
	UID string `json:"uid" binding:"required"`
}

// var db *gorm.DB = config.DB
func CreateUser(payload CreateUserRequest) (*models.User, error) {
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
		Email: payload.Email,
		
	}

	result := db.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil

}

func GetUser(uid string) (*models.User, error) {
	db := config.DB

	var user models.User
	result := db.Where("uid=?", uid).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return &user, nil

}

type ExistsRequest struct{
	Email string `json:"email" binding:"required"`
}

func UserExists(request ExistsRequest) (bool, error) {
	db := config.DB
	var user models.User
	result := db.Where("email=?", request.Email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound{
			return false, nil
		}
		return false, result.Error
	}
	return true, nil
}
