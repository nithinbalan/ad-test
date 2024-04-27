package repositories

import (
	settingsModel "jethur/models/settings"

	"gorm.io/gorm"
)

type RoleRepositoryInterface interface {
	GetByID(id uint) (*settingsModel.Role, error)
}

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

	if result.Error != nil {
		return nil, result.Error
	}

	return role, result.Error
}
