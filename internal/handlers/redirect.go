package handlers

import (
	"github.com/drako02/url-shortener/config"
	"github.com/drako02/url-shortener/internal/models"
)
func GetLongUrl(shortCode string) (string, error) {
	var url models.URL;
	db := config.DB;
	result := db.Where("short_code = ?", shortCode).First(&url);
	if result.Error != nil {
		return "", result.Error
	}
	return url.LongUrl, nil
}