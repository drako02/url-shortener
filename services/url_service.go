package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"github.com/drako02/url-shortener/config"
	"github.com/drako02/url-shortener/models"
	"github.com/drako02/url-shortener/repositories"
	"gorm.io/gorm"
)

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
		UserId:    user.ID,
		Active: true,
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
	var res *gorm.DB
	if limit <= 0 && offset <= 0 {
		res = config.DB.Where("user_id = ?", id).Order("created_at DESC").Find(&urls)
	} else {
		res = config.DB.Where("user_id = ?", id).Order("created_at DESC").Limit(limit).Offset(offset).Find(&urls)
	}

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
	userId, err := GetIdFromUid(uid)
	if err != nil {
		return 0, fmt.Errorf("failed to get user id from uid %s: %w", uid, err)
	}

	res := config.DB.Model(&models.URL{}).Where("user_id = ?", userId).Count(&count)
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

//----------------------------------------------------

const (
	OperatorEqual       FilterOperator = "eq"
	OperatorNotEqual    FilterOperator = "ne"
	OperatorGreaterThan FilterOperator = "gt"
	OperatorLessThan    FilterOperator = "lt"
	OperatorContains    FilterOperator = "contains"
	OperatorStartsWith  FilterOperator = "starts_with"
	OperatorEndsWith    FilterOperator = "ends_with"
	OperatorBetween     FilterOperator = "between"
	OperatorFulltext    FilterOperator = "fulltext"
)

func UrlQueryResult(query UrlQuery) ([]models.URL, int, error) {
	if query.Limit <= 0 {
		query.Limit = 10
	}
	userID, err := GetIdFromUid(query.UID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find user with UID %s: %w", query.UID, err)
	}

	//base query
	db := config.DB.Model(&models.URL{}).Where("user_id = ?", userID)

	log.Printf("operator %v", query.Filters[0].Operator)
	for _, filter := range query.Filters {

		if filter.Operator != OperatorFulltext {
			// 	if !isValidField(filter.Fields) {
			// 		continue
			// 	}
			// } else {
			if !isValidField(filter.Field) {
				continue
			}
		}

		switch filter.Operator {
		case OperatorContains:
			log.Print("contains")

			db = db.Where(filter.Field+" LIKE ?", "%"+fmt.Sprint(filter.Value)+"%")

		case OperatorFulltext:
			log.Print("fulltext")

			var fields []string
			for _, str := range filter.Fields {
				if validFields[str] {
					fields = append(fields, str)
				}
			}

			if len(filter.Fields) == 0 {
				// Default to searching in common fields
				filter.Fields = []string{"short_code", "long_url"}
			}
			concatArgs := strings.Join(fields, ", ' ', ")
			log.Printf("concatArgs: %s", concatArgs)

			db = db.Where(fmt.Sprintf("CONCAT(%s) LIKE ?", concatArgs), "%"+fmt.Sprint(filter.Value)+"%")

			// Implement other cases later
		}

	}

	// total matching records
	var count int64
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// sorting
	if query.SortBy != "" && isValidField(query.SortBy) {
		direction := "ASC"
		if strings.ToUpper(query.SortOrder) == "DESC" {
			direction = "DESC"
		}
		db = db.Order(fmt.Sprintf("%s %s", query.SortBy, direction))
	} else {
		db = db.Order("created_at DESC")
	}

	// pagination
	db = db.Limit(query.Limit).Offset(query.Offset)

	var urls []models.URL
	if err := db.Find(&urls).Error; err != nil {
		return nil, 0, err
	}

	return urls, int(count), nil

}

var validFields = map[string]bool{
	"id":         true,
	"short_code": true,
	"long_url":   true,
	"created_at": true,
	"updated_at": true,
	"clicks":     true,
}

func isValidField(field any) bool {

	switch typedField := field.(type) {
	case string:
		return validFields[typedField]
	case []string:
		for _, str := range typedField {
			if !validFields[str] {
				return false
			}
		}
		return true
	default:
		return false

	}
}

type URLService struct {
	// repo repositories.RepoInterface
	updater repositories.Updater
	deleter repositories.Deleter
}

type URLManagementService interface {
	DeleteURL(id uint) (models.URL, error)
	SetUrlActiveStatus(ctx context.Context, id uint, value bool) error
}

func NewURLService(repo repositories.RepoInterface) *URLService {
	return &URLService{updater: repo, deleter: repo}
}

func (s *URLService) DeleteURL(id uint) (models.URL, error) {
	ctx := context.Background()
	return s.deleter.Delete(ctx, id)
}

// func _DeleteUrl(shortCode string) error {
//     if err := config.DB.
//         Delete(&models.URL{}, "short_code = ?", shortCode).
//         Error; err != nil {
//         return fmt.Errorf("failed to delete url %q: %w", shortCode, err)
//     }
//     return nil
// }

func (s *URLService) SetUrlActiveStatus(ctx context.Context, id uint, value bool) error {
	err := s.updater.UpdateById(ctx, id, repositories.Data{Field: "active", Value: value})
	if err != nil {
		return fmt.Errorf("failed to activate url with id %d: %v", id, err)
	}

	return nil
}

