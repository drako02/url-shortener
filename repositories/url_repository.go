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
