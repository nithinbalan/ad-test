package services

import (
	"errors"
	"jethur/grc/audit/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditRiskRatingServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AuditRiskRating, error)
	GetByID(id uint) (*models.AuditRiskRating, error)
	GetByType(riskRatingType string) (*models.AuditRiskRating, error)
}

// Service and implementations
type AuditRiskRatingService struct {
	repo repositories.AuditRiskRatingRepositoryInterface
}

func NewAuditRiskRatingService(repo repositories.AuditRiskRatingRepositoryInterface) AuditRiskRatingServiceInterface {
	return &AuditRiskRatingService{
		repo: repo,
	}
}

func (arrs *AuditRiskRatingService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AuditRiskRating, error) {
	return arrs.repo.GetAll(sort, sortBy, keyword, status)
}

func (arrs AuditRiskRatingService) GetByID(id uint) (*models.AuditRiskRating, error) {
	auditRiskRating, err := arrs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return auditRiskRating, nil
}

func (arrs *AuditRiskRatingService) GetByType(riskRatingType string) (*models.AuditRiskRating, error) {
	riskRatingResult, err := arrs.repo.GetByType(riskRatingType)
	if err != nil {
		// Handle error appropriately
		return nil, err
	}

	return riskRatingResult, nil
}
