package services

import (
	"errors"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/settings/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type RoleServiceInterface interface {
	GetByID(id uint) (*models.Role, error)
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

func (rs RoleService) GetByID(id uint) (*models.Role, error) {
	role, err := rs.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return role, nil
}
