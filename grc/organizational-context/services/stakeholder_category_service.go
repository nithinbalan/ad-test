package services

import (
	"errors"
	"jethur/grc/models"
	"jethur/grc/organizational-context/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type StakeholderCategoryServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.StakeholderCategory, error)
	GetByID(id uint) (*models.StakeholderCategory, error)
	GetByType(category string) (*models.StakeholderCategory, error)
}

type StakeholderCategoryService struct {
	repo repositories.StakeholderCategoryRepositoryInterface
}

func NewStakeholderCategoryService(repo repositories.StakeholderCategoryRepositoryInterface) StakeholderCategoryServiceInterface {
	return &StakeholderCategoryService{
		repo: repo,
	}
}

func (scs *StakeholderCategoryService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.StakeholderCategory, error) {
	return scs.repo.GetAll(sort, sortBy, keyword, status)
}

func (scs *StakeholderCategoryService) GetByID(id uint) (*models.StakeholderCategory, error) {
	category, err := scs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return category, nil
}

func (scs *StakeholderCategoryService) GetByType(category string) (*models.StakeholderCategory, error) {
	categoryResult, err := scs.repo.GetByType(category)
	if err != nil {
		return nil, err
	}
	return categoryResult, nil
}
