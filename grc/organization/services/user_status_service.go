package services

import (
	"jethur/grc/organization/models"
	"jethur/grc/organization/repositories"
)

type UserStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.UserStatus, error)
	GetByID(id uint) (*models.UserStatus, error)
}

// Service and implementations
type UserStatusService struct {
	repo repositories.UserStatusRepositoryInterface
}

func NewUserStatusService(repo repositories.UserStatusRepositoryInterface) UserStatusServiceInterface {
	return &UserStatusService{
		repo: repo,
	}
}

func (uss *UserStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.UserStatus, error) {
	return uss.repo.GetAll(sort, sortBy, keyword, status)
}

func (uss UserStatusService) GetByID(id uint) (*models.UserStatus, error) {
	return uss.repo.GetByID(id)
}
