package services

import (
	"errors"
	"jethur/grc/models"
	"jethur/grc/settings/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type RolePermissionTypeServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string) ([]*models.RolePermissionType, error)
	GetByID(id uint) (*models.RolePermissionType, error)
	GetByType(permissionType string) (*models.RolePermissionType, error)
}

type RolePermissionTypeService struct {
	repo repositories.RolePermissionTypeRepositoryInterface
}

func NewRolePermissionTypeService(repo repositories.RolePermissionTypeRepositoryInterface) RolePermissionTypeServiceInterface {
	return &RolePermissionTypeService{
		repo: repo,
	}
}

func (rpts *RolePermissionTypeService) GetAll(sort string, sortBy string, keyword string) ([]*models.RolePermissionType, error) {
	return rpts.repo.GetAll(sort, sortBy, keyword)
}

func (rpts *RolePermissionTypeService) GetByID(id uint) (*models.RolePermissionType, error) {
	rolePermissionType, err := rpts.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return rolePermissionType, nil
}

func (rpts *RolePermissionTypeService) GetByType(permissionType string) (*models.RolePermissionType, error) {
	rolePermissionType, err := rpts.repo.GetByType(permissionType)
	if err != nil {
		return nil, err
	}
	return rolePermissionType, nil
}
