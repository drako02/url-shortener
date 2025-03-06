package repositories

import (
	"log"

	"github.com/drako02/url-shortener/config"
	"github.com/drako02/url-shortener/models"
)

func Migrate(){
	err := config.DB.AutoMigrate(&models.URL{}, &models.User{})
	if (err != nil){
		log.Fatal("Migration failed")
	}
}