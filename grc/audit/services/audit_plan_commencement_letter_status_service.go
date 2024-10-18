package services

import (
	"errors"
	"jethur/grc/audit/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditPlanCommencementLetterStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AuditPlanCommencementLetterStatus, error)
	GetByID(id uint) (*models.AuditPlanCommencementLetterStatus, error)
	GetByType(status string) (*models.AuditPlanCommencementLetterStatus, error)
}

// Service and implementations
type AuditPlanCommencementLetterStatusService struct {
	repo repositories.AuditPlanCommencementLetterStatusRepositoryInterface
}

func NewAuditPlanCommencementLetterStatusService(repo repositories.AuditPlanCommencementLetterStatusRepositoryInterface) AuditPlanCommencementLetterStatusServiceInterface {
	return &AuditPlanCommencementLetterStatusService{
		repo: repo,
	}
}

func (aps *AuditPlanCommencementLetterStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AuditPlanCommencementLetterStatus, error) {
	return aps.repo.GetAll(sort, sortBy, keyword, status)
}

func (aps *AuditPlanCommencementLetterStatusService) GetByID(id uint) (*models.AuditPlanCommencementLetterStatus, error) {
	status, err := aps.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return status, nil
}

func (aps *AuditPlanCommencementLetterStatusService) GetByType(status string) (*models.AuditPlanCommencementLetterStatus, error) {
	statusResult, err := aps.repo.GetByType(status)
	if err != nil {
		return nil, err
	}
	return statusResult, nil
}
