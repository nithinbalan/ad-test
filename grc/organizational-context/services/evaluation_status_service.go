package services

import (
	"errors"
	"jethur/grc/models"
	"jethur/grc/organizational-context/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type EvaluationStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.EvaluationStatus, error)
	GetByID(id uint) (*models.EvaluationStatus, error)
	GetByType(status string) (*models.EvaluationStatus, error)
}

type EvaluationStatusService struct {
	repo repositories.EvaluationStatusRepositoryInterface
}

func NewEvaluationStatusService(repo repositories.EvaluationStatusRepositoryInterface) EvaluationStatusServiceInterface {
	return &EvaluationStatusService{
		repo: repo,
	}
}

func (ess *EvaluationStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.EvaluationStatus, error) {
	return ess.repo.GetAll(sort, sortBy, keyword, status)
}

func (ess *EvaluationStatusService) GetByID(id uint) (*models.EvaluationStatus, error) {
	status, err := ess.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return status, nil
}

func (ess *EvaluationStatusService) GetByType(status string) (*models.EvaluationStatus, error) {
	statusResult, err := ess.repo.GetByType(status)
	if err != nil {
		return nil, err
	}
	return statusResult, nil
}
