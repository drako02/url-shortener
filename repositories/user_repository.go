package repositories

import (
	"context"
	"fmt"

	"github.com/drako02/url-shortener/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) *UserRepository {
	return &UserRepository{DB}
}

func (repo *UserRepository) UpdateById(id uint, fields map[string]string, ctx context.Context) error {
	ok := isValidUserField(fields)
	if !ok {
		return fmt.Errorf("invalid User field(s) to update %v", fields)
	}

	err := repo.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Updates(fields).Error
	})

	if err != nil {
		return fmt.Errorf("failed to update user with id %d: %w", id, err)
	}

	return nil
}

var validUserFields = map[string]bool{
	"first_name": true,
	"last_name":  true,
	"email":     true,
}

func isValidUserField(fields map[string]string) bool {
	for f := range fields {
		if !validUserFields[f] {
			return false
		}
	}

	return true
}
