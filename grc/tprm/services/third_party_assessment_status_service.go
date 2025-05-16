package services

import (
	"jethur/grc/models"
	"jethur/grc/tprm/repositories"
)

type ThirdPartyAssessmentStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ThirdPartyAssessmentStatus, error)
	GetByID(id uint) (*models.ThirdPartyAssessmentStatus, error)
	GetByType(status string) (*models.ThirdPartyAssessmentStatus, error) // Renamed method

}

// Service and implementations
type ThirdPartyAssessmentStatusService struct {
	repo repositories.ThirdPartyAssessmentStatusRepositoryInterface
}

func NewThirdPartyAssessmentStatusService(repo repositories.ThirdPartyAssessmentStatusRepositoryInterface) ThirdPartyAssessmentStatusServiceInterface {
	return &ThirdPartyAssessmentStatusService{
		repo: repo,
	}
}

func (lss *ThirdPartyAssessmentStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ThirdPartyAssessmentStatus, error) {
	return lss.repo.GetAll(sort, sortBy, keyword, status)
}

func (lss ThirdPartyAssessmentStatusService) GetByID(id uint) (*models.ThirdPartyAssessmentStatus, error) {
	return lss.repo.GetByID(id)
}

func (lss *ThirdPartyAssessmentStatusService) GetByType(status string) (*models.ThirdPartyAssessmentStatus, error) {
	// Fetch the thirdPartyAssessment status by its status
	statusResult, err := lss.repo.GetByType(status)
	if err != nil {
		// Handle error appropriately
		return nil, err
	}

	return statusResult, nil
}
