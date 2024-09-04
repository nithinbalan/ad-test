package services

import (
	"errors"
	"jethur/grc/audit/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditDraftFindingStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string) ([]*models.AuditDraftFindingStatus, error)
	GetByID(id uint) (*models.AuditDraftFindingStatus, error)
	GetByType(status string) (*models.AuditDraftFindingStatus, error)
}

// Service and implementations
type AuditDraftFindingStatusService struct {
	repo repositories.AuditDraftFindingStatusRepositoryInterface
}

func NewAuditDraftFindingStatusService(repo repositories.AuditDraftFindingStatusRepositoryInterface) AuditDraftFindingStatusServiceInterface {
	return &AuditDraftFindingStatusService{
		repo: repo,
	}
}

func (adfss *AuditDraftFindingStatusService) GetAll(sort string, sortBy string, keyword string) ([]*models.AuditDraftFindingStatus, error) {
	return adfss.repo.GetAll(sort, sortBy, keyword)
}

func (adfss *AuditDraftFindingStatusService) GetByID(id uint) (*models.AuditDraftFindingStatus, error) {
	auditDraftFindingStatus, err := adfss.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return auditDraftFindingStatus, nil
}

func (adfss *AuditDraftFindingStatusService) GetByType(status string) (*models.AuditDraftFindingStatus, error) {
	statusResult, err := adfss.repo.GetByType(status)
	if err != nil {
		return nil, err
	}

	return statusResult, nil
}
