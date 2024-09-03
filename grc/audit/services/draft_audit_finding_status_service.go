package services

import (
	"errors"
	"jethur/grc/audit/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type DraftAuditFindingStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string) ([]*models.DraftAuditFindingStatus, error)
	GetByID(id uint) (*models.DraftAuditFindingStatus, error)
	GetByType(status string) (*models.DraftAuditFindingStatus, error)
}

// Service and implementations
type DraftAuditFindingStatusService struct {
	repo repositories.DraftAuditFindingStatusRepositoryInterface
}

func NewDraftAuditFindingStatusService(repo repositories.DraftAuditFindingStatusRepositoryInterface) DraftAuditFindingStatusServiceInterface {
	return &DraftAuditFindingStatusService{
		repo: repo,
	}
}

func (dafss *DraftAuditFindingStatusService) GetAll(sort string, sortBy string, keyword string) ([]*models.DraftAuditFindingStatus, error) {
	return dafss.repo.GetAll(sort, sortBy, keyword)
}

func (dafss *DraftAuditFindingStatusService) GetByID(id uint) (*models.DraftAuditFindingStatus, error) {
	draftAuditFindingStatus, err := dafss.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return draftAuditFindingStatus, nil
}

func (dafss *DraftAuditFindingStatusService) GetByType(status string) (*models.DraftAuditFindingStatus, error) {
	statusResult, err := dafss.repo.GetByType(status)
	if err != nil {
		return nil, err
	}

	return statusResult, nil
}
