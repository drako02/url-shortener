package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
	"time"

	"github.com/drako02/url-shortener/config"
	"github.com/drako02/url-shortener/internal/models"
)
type CreateRequest struct{
    URL string `json:"url" binding:"required,url"`
}


func GenerateShortCode() string {
	b := make([]byte, 4)
	_,err := rand.Read(b)
	if(err != nil){
		panic(err)
	}
	return strings.TrimRight(base64.URLEncoding.EncodeToString(b), "=")[:6]


}

func CreateShortUrl(request CreateRequest)map[string] any {
	longUrl := request.URL

	dbEntry := models.URL{
		ShortCode: GenerateShortCode(),
		LongUrl: longUrl,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	config.DB.Create(&dbEntry)

	return map[string] any {
		"shortCode": dbEntry.ShortCode,
		"longUrl": dbEntry.LongUrl,
	}
}