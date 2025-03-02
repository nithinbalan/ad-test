package services

import (
	"errors"
	"jethur/grc/business-continuity/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type TestAndExerciseStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.TestAndExerciseStatus, error)
	GetByID(id uint) (*models.TestAndExerciseStatus, error)
	GetByType(status string) (*models.TestAndExerciseStatus, error)
}

// Service and implementations
type TestAndExerciseStatusService struct {
	repo repositories.TestAndExerciseStatusRepositoryInterface
}

func NewTestAndExerciseStatusService(repo repositories.TestAndExerciseStatusRepositoryInterface) TestAndExerciseStatusServiceInterface {
	return &TestAndExerciseStatusService{
		repo: repo,
	}
}

func (tess *TestAndExerciseStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.TestAndExerciseStatus, error) {
	return tess.repo.GetAll(sort, sortBy, keyword, status)
}

func (tess TestAndExerciseStatusService) GetByID(id uint) (*models.TestAndExerciseStatus, error) {
	testAndExerciseStatus, err := tess.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return testAndExerciseStatus, nil
}

func (tess *TestAndExerciseStatusService) GetByType(status string) (*models.TestAndExerciseStatus, error) {
	testAndExerciseStatus, err := tess.repo.GetByType(status)

	if err != nil {
		return nil, err
	}

	return testAndExerciseStatus, nil
}
