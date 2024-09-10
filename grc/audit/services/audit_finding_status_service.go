package services

import (
	"errors"
	"jethur/grc/audit/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditFindingStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string) ([]*models.AuditFindingStatus, error)
	GetByID(id uint) (*models.AuditFindingStatus, error)
	GetByType(status string) (*models.AuditFindingStatus, error)
}

// Service and implementations
type AuditFindingStatusService struct {
	repo repositories.AuditFindingStatusRepositoryInterface
}

func NewAuditFindingStatusService(repo repositories.AuditFindingStatusRepositoryInterface) AuditFindingStatusServiceInterface {
	return &AuditFindingStatusService{
		repo: repo,
	}
}

func (afss *AuditFindingStatusService) GetAll(sort string, sortBy string, keyword string) ([]*models.AuditFindingStatus, error) {
	return afss.repo.GetAll(sort, sortBy, keyword)
}

func (afss *AuditFindingStatusService) GetByID(id uint) (*models.AuditFindingStatus, error) {
	auditFindingStatus, err := afss.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return auditFindingStatus, nil
}

func (afss *AuditFindingStatusService) GetByType(status string) (*models.AuditFindingStatus, error) {
	statusResult, err := afss.repo.GetByType(status)
	if err != nil {
		return nil, err
	}

	return statusResult, nil
}
