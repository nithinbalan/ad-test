package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

type RoleRepositoryInterface interface {
	GetByID(id uint) (*models.Role, error)
}

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

	if result.Error != nil {
		return nil, result.Error
	}

	return role, result.Error
}
