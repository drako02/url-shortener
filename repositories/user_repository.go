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

func (repo *UserRepository) UpdateById(id uint, fields map[string]string, ctx context.Context) (*models.User, error) {
	ok := isValidUserField(fields)
	if !ok {
		return nil, fmt.Errorf("invalid User field(s) to update %v", fields)
	}

	updateFields := make(map[string]interface{}, len(fields))
	for k, v := range fields {
		updateFields[k] = v
	}
	user := models.User{}
	err := repo.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Model(&user).Where("id = ?", id).Updates(updateFields).Error
	})

	if err != nil {
		return nil, fmt.Errorf("failed to update user with id %d: %w", id, err)
	}

	return &user, nil
}

var validUserFields = map[string]bool{
	"first_name": true,
	"last_name":  true,
	"email":      true,
}

func isValidUserField(fields map[string]string) bool {
	for f := range fields {
		if !validUserFields[f] {
			return false
		}
	}

	return true
}
