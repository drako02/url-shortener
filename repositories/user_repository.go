package repositories

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) *UserRepository {
	return &UserRepository{DB}
}

func (repo *UserRepository) UpdateById(id uint, field any, value string, ctx context.Context) error {
	ok := isValidUserField(field)
	if !ok {
		return fmt.Errorf("Invalid User field(s) to update %w", field)
	}

	switch typedField := field.(type) {
	case string:
		err := repo.DB.WithContext(ctx).Where("id = ?", id).Update(typedField, value).Error
		if err != nil {
			return fmt.Errorf("Failed to update user field %s, error: %w", typedField, err)
		}
		return nil

	case []string:
		err := repo.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			for _, f := range typedField {
				// tx.WithContext(ctx).Where("id = ?", id).Update(f,)
				// FIXME We might have to change the arguments of the main function to use a struct
				// which specifies the fields the client wants to updare
			}
		})
	}

}

var validUserFields = map[string]bool{
	"firstName": true,
	"lastName":  true,
	"email":     true,
}

func isValidUserField(field any) bool {

	switch typedField := field.(type) {
	case string:
		return validUserFields[typedField]

	case []string:
		for _, f := range typedField {
			if !validFields[f] {
				return false
			}
		}
		return true

	default:
		return false
	}
}
