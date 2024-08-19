package services

import (
	"errors"
	"jethur/grc/audit/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditScheduleStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AuditScheduleStatus, error)
	GetByID(id uint) (*models.AuditScheduleStatus, error)
	GetByType(status string) (*models.AuditScheduleStatus, error)
}

// Service and implementations
type AuditScheduleStatusService struct {
	repo repositories.AuditScheduleStatusRepositoryInterface
}

func NewAuditScheduleStatusService(repo repositories.AuditScheduleStatusRepositoryInterface) AuditScheduleStatusServiceInterface {
	return &AuditScheduleStatusService{
		repo: repo,
	}
}

func (ass *AuditScheduleStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AuditScheduleStatus, error) {
	return ass.repo.GetAll(sort, sortBy, keyword, status)
}

func (ass AuditScheduleStatusService) GetByID(id uint) (*models.AuditScheduleStatus, error) {
	auditScheduleStatus, err := ass.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return auditScheduleStatus, nil
}

func (ass *AuditScheduleStatusService) GetByType(status string) (*models.AuditScheduleStatus, error) {
	statusResult, err := ass.repo.GetByType(status)
	if err != nil {
		// Handle error appropriately
		return nil, err
	}

	return statusResult, nil
}
