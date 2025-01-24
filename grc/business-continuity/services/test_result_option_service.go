package services

import (
	"errors"
	"jethur/grc/business-continuity/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type TestResultOptionServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.TestResultOption, error)
	GetByID(id uint) (*models.TestResultOption, error)
	GetByType(status string) (*models.TestResultOption, error)
}

// Service and implementations
type TestResultOptionService struct {
	repo repositories.TestResultOptionRepositoryInterface
}

func NewTestResultOptionService(repo repositories.TestResultOptionRepositoryInterface) TestResultOptionServiceInterface {
	return &TestResultOptionService{
		repo: repo,
	}
}

func (tess *TestResultOptionService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.TestResultOption, error) {
	return tess.repo.GetAll(sort, sortBy, keyword, status)
}

func (tess TestResultOptionService) GetByID(id uint) (*models.TestResultOption, error) {
	testResultOption, err := tess.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return testResultOption, nil
}

func (tess *TestResultOptionService) GetByType(status string) (*models.TestResultOption, error) {
	testResultOption, err := tess.repo.GetByType(status)

	if err != nil {
		return nil, err
	}

	return testResultOption, nil
}
