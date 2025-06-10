package services

import (
	"fmt"
	"time"

	"github.com/drako02/url-shortener/config"
	"github.com/drako02/url-shortener/models"
	"github.com/drako02/url-shortener/utils"
	"gorm.io/gorm"
)

func GetIdFromUid(uid string) (uint, error) {
	var user models.User
	res := config.DB.Where("uid = ?", uid).First(&user)
	if res.Error != nil {
		return 0, res.Error
	}
	return user.ID, nil

}

func CreateUser(payload utils.CreateUserRequest) (*models.User, error) {
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
		Email:     payload.Email,
	}

	result := db.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil

}

// func DeleteUser (){
// 	db := config.DB
// }

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

func UserExists(request utils.ExistsRequest) (bool, error) {
	db := config.DB
	var user models.User
	result := db.Where("email=?", request.Email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, result.Error
	}
	return true, nil
}
