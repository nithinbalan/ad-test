package services

import (
	"errors"
	"jethur/grc/incident/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type IncidentInvestigationStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.IncidentInvestigationStatus, error)
	GetByID(id uint) (*models.IncidentInvestigationStatus, error)
}

// Service and implementations
type IncidentInvestigationStatusService struct {
	repo repositories.IncidentInvestigationStatusRepositoryInterface
}

func NewIncidentInvestigationStatusService(repo repositories.IncidentInvestigationStatusRepositoryInterface) IncidentInvestigationStatusServiceInterface {
	return &IncidentInvestigationStatusService{
		repo: repo,
	}
}

func (iiss *IncidentInvestigationStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.IncidentInvestigationStatus, error) {
	return iiss.repo.GetAll(sort, sortBy, keyword, status)
}

func (iiss IncidentInvestigationStatusService) GetByID(id uint) (*models.IncidentInvestigationStatus, error) {
	incidentInvestigationStatus, err := iiss.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return incidentInvestigationStatus, nil
}
