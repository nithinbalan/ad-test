package services

import (
	"errors"
	"jethur/grc/identity-access-review/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ReviewCampaignStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ReviewCampaignStatus, error)
	GetByID(id uint) (*models.ReviewCampaignStatus, error)
	GetByType(reviewCampaignStatusType string) (*models.ReviewCampaignStatus, error)
}

type ReviewCampaignStatusService struct {
	repo repositories.ReviewCampaignStatusRepositoryInterface
}

func NewReviewCampaignStatusService(
	repo repositories.ReviewCampaignStatusRepositoryInterface,
) ReviewCampaignStatusServiceInterface {
	return &ReviewCampaignStatusService{
		repo: repo,
	}
}

func (rcss *ReviewCampaignStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ReviewCampaignStatus, error) {
	return rcss.repo.GetAll(sort, sortBy, keyword, status)
}

func (rcss *ReviewCampaignStatusService) GetByID(id uint) (*models.ReviewCampaignStatus, error) {
	reviewCampaignStatus, err := rcss.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return reviewCampaignStatus, nil
}

func (rcss *ReviewCampaignStatusService) GetByType(reviewCampaignStatusType string) (*models.ReviewCampaignStatus, error) {
	reviewCampaignStatus, err := rcss.repo.GetByType(reviewCampaignStatusType)

	if err != nil {
		return nil, err
	}

	return reviewCampaignStatus, nil
}
