package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/drako02/url-shortener/config"
	"github.com/drako02/url-shortener/models"
)

type CreateRequest struct {
	URL string `json:"url" binding:"required,url"`
	UID string `json:"uid" binding:"required"`
}

func GenerateShortCode() string {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return strings.TrimRight(base64.URLEncoding.EncodeToString(b), "=")[:6]

}

func CreateShortUrl(request CreateRequest) (map[string]any, error) {
	longUrl := request.URL
	uid := request.UID

	user := models.User{}
	result := config.DB.Where("uid = ?", uid).First(&user)
	if result.Error != nil {
		return nil, fmt.Errorf("user not found: %w", result.Error)
	}
	dbEntry := models.URL{
		ShortCode: GenerateShortCode(),
		LongUrl:   longUrl,
		// CreatedAt: time.Now(),
		// UpdatedAt: time.Now(),
		UserId: user.ID,
	}

	if err := config.DB.Create(&dbEntry).Error; err != nil {
		return nil, fmt.Errorf("failed to create shortened URL: %w", err)
	}

	return map[string]any{
		"short_code": dbEntry.ShortCode,
		"long_url":   dbEntry.LongUrl,
	}, nil
}
