package services

import (
	"errors"
	"jethur/grc/audit/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditDocumentRequestStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AuditDocumentRequestStatus, error)
	GetByID(id uint) (*models.AuditDocumentRequestStatus, error)
	GetByType(status string) (*models.AuditDocumentRequestStatus, error)
}

// Service and implementations
type AuditDocumentRequestStatusService struct {
	repo repositories.AuditDocumentRequestStatusRepositoryInterface
}

func NewAuditDocumentRequestStatusService(repo repositories.AuditDocumentRequestStatusRepositoryInterface) AuditDocumentRequestStatusServiceInterface {
	return &AuditDocumentRequestStatusService{
		repo: repo,
	}
}

func (adrss *AuditDocumentRequestStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AuditDocumentRequestStatus, error) {
	return adrss.repo.GetAll(sort, sortBy, keyword, status)
}

func (adrss AuditDocumentRequestStatusService) GetByID(id uint) (*models.AuditDocumentRequestStatus, error) {
	auditDocumentRequestStatus, err := adrss.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return auditDocumentRequestStatus, nil
}

func (adrss *AuditDocumentRequestStatusService) GetByType(status string) (*models.AuditDocumentRequestStatus, error) {
	statusResult, err := adrss.repo.GetByType(status)
	if err != nil {
		// Handle error appropriately
		return nil, err
	}

	return statusResult, nil
}
