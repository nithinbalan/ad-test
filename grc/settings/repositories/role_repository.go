package repositories

import (
	"errors"
	settingsModel "jethur/models/settings"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type RoleRepositoryInterface interface {
	GetByID(id uint) (*settingsModel.Role, error)
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

func (repo *RoleRepository) GetByID(id uint) (*settingsModel.Role, error) {
	var role *settingsModel.Role
	result := repo.db.First(&role, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return role, result.Error
}
