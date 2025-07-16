package services

import (
	"context"
	"fmt"
	"time"

	"github.com/drako02/url-shortener/config"
	"github.com/drako02/url-shortener/models"
	"github.com/drako02/url-shortener/repositories"
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

type UserService struct {
	Repo *repositories.UserRepository
}

func NewUserService(svc *repositories.UserRepository) *UserService {
	return &UserService{svc}
}

type ValidUserFields struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func (f *ValidUserFields) ToMap() map[string]string {
    result := make(map[string]string)
    if f.FirstName != "" {
        result["first_name"] = f.FirstName
    }
    if f.LastName != "" {
        result["last_name"] = f.LastName
    }
    if f.Email != "" {
        result["email"] = f.Email
    }
    return result
}

func (svc *UserService) UpdateUserInfo(id uint, fields ValidUserFields , ctx context.Context) error {
	err := svc.Repo.UpdateById(id, fields.ToMap(), ctx)
	if err != nil {
		return fmt.Errorf("failed to update user info: %v", err)
	}
	return nil
}
