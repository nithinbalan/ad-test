package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

type ActionPlanDelegateeUserRepositoryInterface interface {
	GetUserByID(userID uint) (*models.User, error)
	HasViewPermission(userID uint, moduleID uint) (bool, error)
}

type ActionPlanDelegateeUserRepository struct {
	db *gorm.DB
}

func NewActionPlanDelegateeUserRepository(db *gorm.DB) *ActionPlanDelegateeUserRepository {
	return &ActionPlanDelegateeUserRepository{
		db: db,
	}
}

func (repo *ActionPlanDelegateeUserRepository) GetUserByID(userID uint) (*models.User, error) {
	var user *models.User

	result := repo.db.First(&user, userID)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (repo *ActionPlanDelegateeUserRepository) HasViewPermission(userID uint, moduleID uint) (bool, error) {
	var count int64

	err := repo.db.Table("user_roles ur").
		Joins("JOIN role_permissions rp ON ur.role_id = rp.role_id").
		Joins("JOIN role_permission_types rpt ON rp.role_permission_type_id = rpt.id").
		Where("ur.user_id = ?", userID).
		Where("rp.module_id = ?", moduleID).
		Where("rp.is_granted = ?", true).
		Where("rpt.type = ?", "view").
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
