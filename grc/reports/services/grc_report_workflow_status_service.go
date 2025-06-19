package services

import (
	"errors"
	"jethur/grc/models"
	"jethur/grc/reports/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type GRCReportWorkflowStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.GRCReportWorkflowStatus, error)
	GetByID(id uint) (*models.GRCReportWorkflowStatus, error)
	GetByType(status string) (*models.GRCReportWorkflowStatus, error)
}

// Service and implementations
type GRCReportWorkflowStatusService struct {
	repo repositories.GRCReportWorkflowStatusRepositoryInterface
}

func NewGRCReportWorkflowStatusService(repo repositories.GRCReportWorkflowStatusRepositoryInterface) GRCReportWorkflowStatusServiceInterface {
	return &GRCReportWorkflowStatusService{
		repo: repo,
	}
}

func (dss *GRCReportWorkflowStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.GRCReportWorkflowStatus, error) {
	return dss.repo.GetAll(sort, sortBy, keyword, status)
}

func (dss GRCReportWorkflowStatusService) GetByID(id uint) (*models.GRCReportWorkflowStatus, error) {
	documentStatus, err := dss.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return documentStatus, nil
}

func (dss *GRCReportWorkflowStatusService) GetByType(status string) (*models.GRCReportWorkflowStatus, error) {
	documentStatus, err := dss.repo.GetByType(status)

	if err != nil {
		return nil, err
	}

	return documentStatus, nil
}
