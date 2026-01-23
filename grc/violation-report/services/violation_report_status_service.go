package services

import (
	"errors"
	"jethur/grc/models"
	"jethur/grc/violation-report/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type ViolationReportStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ViolationReportStatus, error)
	GetByID(id uint) (*models.ViolationReportStatus, error)
	GetByType(status string) (*models.ViolationReportStatus, error)
}

type ViolationReportStatusService struct {
	repo repositories.ViolationReportStatusRepositoryInterface
}

func NewViolationReportStatusService(repo repositories.ViolationReportStatusRepositoryInterface) ViolationReportStatusServiceInterface {
	return &ViolationReportStatusService{
		repo: repo,
	}
}

func (vrss *ViolationReportStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ViolationReportStatus, error) {
	return vrss.repo.GetAll(sort, sortBy, keyword, status)
}

func (vrss *ViolationReportStatusService) GetByID(id uint) (*models.ViolationReportStatus, error) {
	violationReportStatus, err := vrss.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return violationReportStatus, nil
}

func (vrss *ViolationReportStatusService) GetByType(typeValue string) (*models.ViolationReportStatus, error) {
	statusResult, err := vrss.repo.GetByType(typeValue)
	if err != nil {
		return nil, err
	}
	return statusResult, nil
}
