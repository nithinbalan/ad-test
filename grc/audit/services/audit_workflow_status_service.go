package services

import (
	"errors"
	"jethur/grc/audit/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditWorkflowStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string) ([]*models.AuditWorkflowStatus, error)
	GetByID(id uint) (*models.AuditWorkflowStatus, error)
	GetByType(status string) (*models.AuditWorkflowStatus, error)
}

// Service and implementations
type AuditWorkflowStatusService struct {
	repo repositories.AuditWorkflowStatusRepositoryInterface
}

func NewAuditWorkflowStatusService(repo repositories.AuditWorkflowStatusRepositoryInterface) AuditWorkflowStatusServiceInterface {
	return &AuditWorkflowStatusService{
		repo: repo,
	}
}

func (awss *AuditWorkflowStatusService) GetAll(sort string, sortBy string, keyword string) ([]*models.AuditWorkflowStatus, error) {
	return awss.repo.GetAll(sort, sortBy, keyword)
}

func (awss *AuditWorkflowStatusService) GetByID(id uint) (*models.AuditWorkflowStatus, error) {
	auditWorkflowStatus, err := awss.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return auditWorkflowStatus, nil
}

func (awss *AuditWorkflowStatusService) GetByType(status string) (*models.AuditWorkflowStatus, error) {
	statusResult, err := awss.repo.GetByType(status)
	if err != nil {
		return nil, err
	}

	return statusResult, nil
}
