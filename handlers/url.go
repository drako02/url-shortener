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

type GetUserUrlRequest struct {
	UID    string `json:"uid" binding:"required"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

func GetLongUrl(shortCode string) (string, error) {
	var url models.URL
	db := config.DB
	result := db.Where("short_code = ?", shortCode).First(&url)
	if result.Error != nil {
		return "", result.Error
	}
	return url.LongUrl, nil
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

func GetUserUrls(uid string, limit int, offset int) ([]models.URL, error) {
	var urls []models.URL
	id, err := GetIdFromUid(uid)
	if err != nil {
		return nil, fmt.Errorf("failed to find user with UID %s: %w", uid, err)

	}
	fmt.Printf("Found user ID: %d for UID: %s\n", id, uid)

	res := config.DB.Where("user_id = ?", id).Limit(limit).Offset(offset).Find(&urls)
	if res.Error != nil {
		return nil, res.Error
	}
	fmt.Printf("Retrieved %d URLs for user ID: %d\n", len(urls), id)

	if len(urls) > 0 {
		fmt.Printf("First URL: ShortCode=%s, LongUrl=%s\n",
			urls[0].ShortCode,
			truncateString(urls[0].LongUrl, 30))
	}

	return urls, nil
}

func GetTotalUrls(uid string) (int, error) {
	var count int64
	// var urls []models.URL
	userId, _ := GetIdFromUid(uid)
	res := config.DB.Table("urls").Where("user_id = ?", userId).Count(&count)
	if res.Error != nil {
		return 0, res.Error
	}
	return int(count), nil
}

func truncateString(s string, max int) string {
	if len(s) > max {
		return s[:max] + "..."
	}
	return s
}
