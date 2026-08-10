package services

import (
	"errors"

	"jethur/grc/identity-access-review/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AccountDecisionServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AccountDecision, error)
	GetByID(id uint) (*models.AccountDecision, error)
}

type AccountDecisionService struct {
	repo repositories.AccountDecisionRepositoryInterface
}

func NewAccountDecisionService(repo repositories.AccountDecisionRepositoryInterface) AccountDecisionServiceInterface {
	return &AccountDecisionService{repo: repo}
}

func (ads *AccountDecisionService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AccountDecision, error) {
	return ads.repo.GetAll(sort, sortBy, keyword, status)
}

func (ads *AccountDecisionService) GetByID(id uint) (*models.AccountDecision, error) {
	accountDecision, err := ads.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return accountDecision, nil
}
