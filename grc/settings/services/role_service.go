package services

import (
	"errors"
	"jethur/grc/logs/services"
	"jethur/grc/settings/repositories"
	settingsModel "jethur/models/settings"
	"jethur/utility"

	"gorm.io/gorm"
)

type RoleServiceInterface interface {
	GetByID(id uint) (*settingsModel.Role, error)
}

// Service and implementations
type RoleService struct {
	repo           repositories.RoleRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewRoleService(repo repositories.RoleRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) RoleServiceInterface {
	return &RoleService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (rs RoleService) GetByID(id uint) (*settingsModel.Role, error) {
	role, err := rs.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return role, nil
}
