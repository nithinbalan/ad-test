package repositories

import (
	"errors"
	"jethur/grc/settings/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type RoleRepositoryInterface interface {
	GetByID(id uint) (*models.Role, error)
}

// implementation
type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{
		db: db,
	}
}

func (repo *RoleRepository) GetByID(id uint) (*models.Role, error) {
	var role *models.Role
	result := repo.db.First(&role, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return role, result.Error
}
