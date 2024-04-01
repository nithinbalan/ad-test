package services

import (
	"jethur/grc/logs/services"
	"jethur/grc/settings/models"
	"jethur/grc/settings/repositories"
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
	return rs.repo.GetByID(id)
}
