package services

import (
	"errors"
	"jethur/grc/audit/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditDraftReportStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string) ([]*models.AuditDraftReportStatus, error)
	GetByID(id uint) (*models.AuditDraftReportStatus, error)
	GetByType(status string) (*models.AuditDraftReportStatus, error)
}

// Service and implementations
type AuditDraftReportStatusService struct {
	repo repositories.AuditDraftReportStatusRepositoryInterface
}

func NewAuditDraftReportStatusService(repo repositories.AuditDraftReportStatusRepositoryInterface) AuditDraftReportStatusServiceInterface {
	return &AuditDraftReportStatusService{
		repo: repo,
	}
}

func (adrss *AuditDraftReportStatusService) GetAll(sort string, sortBy string, keyword string) ([]*models.AuditDraftReportStatus, error) {
	return adrss.repo.GetAll(sort, sortBy, keyword)
}

func (adrss *AuditDraftReportStatusService) GetByID(id uint) (*models.AuditDraftReportStatus, error) {
	auditDraftReportStatus, err := adrss.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return auditDraftReportStatus, nil
}

func (adrss *AuditDraftReportStatusService) GetByType(status string) (*models.AuditDraftReportStatus, error) {
	statusResult, err := adrss.repo.GetByType(status)
	if err != nil {
		return nil, err
	}

	return statusResult, nil
}
