package services

import (
	"errors"
	"jethur/grc/compliance/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ComplianceAssessmentRatingServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ComplianceAssessmentRating, error)
	GetByID(id uint) (*models.ComplianceAssessmentRating, error)
	GetByType(ratingType string) (*models.ComplianceAssessmentRating, error)
}

type ComplianceAssessmentRatingService struct {
	repo repositories.ComplianceAssessmentRatingRepositoryInterface
}

func NewComplianceAssessmentRatingService(repo repositories.ComplianceAssessmentRatingRepositoryInterface) ComplianceAssessmentRatingServiceInterface {
	return &ComplianceAssessmentRatingService{
		repo: repo,
	}
}

func (service *ComplianceAssessmentRatingService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ComplianceAssessmentRating, error) {
	return service.repo.GetAll(sort, sortBy, keyword, status)
}

func (service *ComplianceAssessmentRatingService) GetByID(id uint) (*models.ComplianceAssessmentRating, error) {
	rating, err := service.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return rating, nil
}

func (service *ComplianceAssessmentRatingService) GetByType(ratingType string) (*models.ComplianceAssessmentRating, error) {
	return service.repo.GetByType(ratingType)
}
