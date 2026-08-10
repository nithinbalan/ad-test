package services

import (
	"errors"

	"jethur/grc/identity-access-review/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AccountTypeServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AccountType, error)
	GetByID(id uint) (*models.AccountType, error)
}

type AccountTypeService struct {
	repo repositories.AccountTypeRepositoryInterface
}

func NewAccountTypeService(repo repositories.AccountTypeRepositoryInterface) AccountTypeServiceInterface {
	return &AccountTypeService{repo: repo}
}

func (ats *AccountTypeService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AccountType, error) {
	return ats.repo.GetAll(sort, sortBy, keyword, status)
}

func (ats *AccountTypeService) GetByID(id uint) (*models.AccountType, error) {
	accountType, err := ats.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return accountType, nil
}
