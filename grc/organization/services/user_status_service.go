package services

import (
	"errors"
	"jethur/grc/models"
	"jethur/grc/organization/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type UserStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.UserStatus, error)
	GetByID(id uint) (*models.UserStatus, error)
	GetByType(status string) (*models.UserStatus, error) // Renamed method
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
	userStatus, err := uss.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return userStatus, nil

}

func (uss *UserStatusService) GetByType(status string) (*models.UserStatus, error) {
	statusResult, err := uss.repo.GetByType(status)
	if err != nil {
		return nil, err
	}

	return statusResult, nil
}
