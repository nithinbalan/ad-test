package services

import (
	"errors"
	"jethur/grc/audit/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AuditStatus, error)
	GetByID(id uint) (*models.AuditStatus, error)
	GetByType(status string) (*models.AuditStatus, error)
}

// Service and implementations
type AuditStatusService struct {
	repo repositories.AuditStatusRepositoryInterface
}

func NewAuditStatusService(repo repositories.AuditStatusRepositoryInterface) AuditStatusServiceInterface {
	return &AuditStatusService{
		repo: repo,
	}
}

func (ass *AuditStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AuditStatus, error) {
	return ass.repo.GetAll(sort, sortBy, keyword, status)
}

func (ass AuditStatusService) GetByID(id uint) (*models.AuditStatus, error) {
	auditStatus, err := ass.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return auditStatus, nil
}

func (ass *AuditStatusService) GetByType(status string) (*models.AuditStatus, error) {
	statusResult, err := ass.repo.GetByType(status)
	if err != nil {
		// Handle error appropriately
		return nil, err
	}

	return statusResult, nil
}
