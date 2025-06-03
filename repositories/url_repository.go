package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/drako02/url-shortener/models"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("url not found")

type URLRepository struct {
	DB *gorm.DB
}

func NewURLRepository(db *gorm.DB) *URLRepository {
	return &URLRepository{DB: db}
}

type Deleter interface {
	Delete(ctx context.Context, id uint) (models.URL, error)
}

type Updater interface {
	UpdateById(ctx context.Context, id uint, data Data) error
}

type Getter interface {
	GetByShortCode(ctx context.Context, shortCode string) (models.URL, error)
}

type RepoInterface interface {
	Deleter
	Updater
	Getter
}

var _ Deleter = (*URLRepository)(nil)

func (r *URLRepository) Delete(ctx context.Context, id uint) (models.URL, error) {
	var url models.URL

	// First find the URL
	if err := r.DB.WithContext(ctx).Where("id = ?", id).First(&url).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return url, fmt.Errorf("%w: %d", ErrNotFound, id)
		}
		return url, fmt.Errorf("failed to find url %q: %w", id, err)
	}

	if err := r.DB.WithContext(ctx).Delete(&models.URL{}, "id = ?", id).Error; err != nil {
		return url, fmt.Errorf("failed to delete url %q: %w", id, err)
	}

	return url, nil
}

const (
	// ShortCode string = "short_code"
	LongUrl string = "long_url"
	Active  string = "active"
)

var validFields = map[string]bool{
	LongUrl: true,
	Active:  true,
}

func IsValidUpdateField(field string) bool {
	return validFields[field]
}

type Data struct {
	Field string `json:"field,omitempty"`
	Value any    `json:"value,omitempty"`
}

func (r *URLRepository) UpdateById(ctx context.Context, id uint, data Data) error {
	valid := IsValidUpdateField(data.Field)
	if !valid {
		return fmt.Errorf("invalid field")
	}

	db := r.DB.WithContext(ctx).Model(&models.URL{}).Where("id = ?", id).Update(data.Field, data.Value)

	if db.Error != nil {
		return fmt.Errorf("failed to update urls field %s with %v for id %d", data.Field, data.Value, id)
	}

	if db.RowsAffected == 0 {
		return fmt.Errorf("%w: %d", ErrNotFound, id)
	}

	return nil

}

func (r *URLRepository) GetByShortCode(ctx context.Context, shortCode string) (models.URL, error){
	var url models.URL 

	err := r.DB.WithContext(ctx).Model(&url).Where("short_code = ?", shortCode).First(&url).Error
	if err != nil {
		return url, err
	}

	return url, nil
}
