package handlers

import (
	"github.com/drako02/url-shortener/config"
	"github.com/drako02/url-shortener/models"
)

func GetIdFromUid(uid string) (uint, error) {
	var user models.User
	res := config.DB.Where("uid = ?", uid).First(&user)
	if res.Error != nil {
		return 0, res.Error
	}
	return user.ID, nil

}
