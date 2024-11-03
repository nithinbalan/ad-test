package services

import (
	"errors"
	"jethur/grc/audit/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditFinalReportStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string) ([]*models.AuditFinalReportStatus, error)
	GetByID(id uint) (*models.AuditFinalReportStatus, error)
	GetByType(status string) (*models.AuditFinalReportStatus, error)
}

// Service and implementations
type AuditFinalReportStatusService struct {
	repo repositories.AuditFinalReportStatusRepositoryInterface
}

func NewAuditFinalReportStatusService(repo repositories.AuditFinalReportStatusRepositoryInterface) AuditFinalReportStatusServiceInterface {
	return &AuditFinalReportStatusService{
		repo: repo,
	}
}

func (afrss *AuditFinalReportStatusService) GetAll(sort string, sortBy string, keyword string) ([]*models.AuditFinalReportStatus, error) {
	return afrss.repo.GetAll(sort, sortBy, keyword)
}

func (afrss *AuditFinalReportStatusService) GetByID(id uint) (*models.AuditFinalReportStatus, error) {
	auditFinalReportStatus, err := afrss.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return auditFinalReportStatus, nil
}

func (afrss *AuditFinalReportStatusService) GetByType(status string) (*models.AuditFinalReportStatus, error) {
	statusResult, err := afrss.repo.GetByType(status)
	if err != nil {
		return nil, err
	}

	return statusResult, nil
}
