package services

import (
	"errors"
	"jethur/grc/models"
	"jethur/grc/reports/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type GRCReportCommentStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.GRCReportCommentStatus, error)
	GetByID(id uint) (*models.GRCReportCommentStatus, error)
	GetByType(status string) (*models.GRCReportCommentStatus, error)
}

// Service and implementations
type GRCReportCommentStatusService struct {
	repo repositories.GRCReportCommentStatusRepositoryInterface
}

func NewGRCReportCommentStatusService(repo repositories.GRCReportCommentStatusRepositoryInterface) GRCReportCommentStatusServiceInterface {
	return &GRCReportCommentStatusService{
		repo: repo,
	}
}

func (dss *GRCReportCommentStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.GRCReportCommentStatus, error) {
	return dss.repo.GetAll(sort, sortBy, keyword, status)
}

func (dss GRCReportCommentStatusService) GetByID(id uint) (*models.GRCReportCommentStatus, error) {
	documentStatus, err := dss.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return documentStatus, nil
}

func (dss *GRCReportCommentStatusService) GetByType(status string) (*models.GRCReportCommentStatus, error) {
	documentStatus, err := dss.repo.GetByType(status)

	if err != nil {
		return nil, err
	}

	return documentStatus, nil
}
