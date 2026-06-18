package services

import (
	"errors"
	"jethur/grc/awareness-reporting/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AwarenessCampaignStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AwarenessCampaignStatus, error)
	GetByID(id uint) (*models.AwarenessCampaignStatus, error)
	GetByType(status string) (*models.AwarenessCampaignStatus, error)
}


type AwarenessCampaignStatusService struct {
	repo repositories.AwarenessCampaignStatusRepositoryInterface
}

func NewAwarenessCampaignStatusService(
	repo repositories.AwarenessCampaignStatusRepositoryInterface,
) AwarenessCampaignStatusServiceInterface {
	return &AwarenessCampaignStatusService{
		repo: repo,
	}
}

func (acss *AwarenessCampaignStatusService) GetAll(sort string,sortBy string,keyword string,status string,) ([]*models.AwarenessCampaignStatus, error) {
	return acss.repo.GetAll(sort, sortBy, keyword, status)
}

func (acss AwarenessCampaignStatusService) GetByID(id uint) (*models.AwarenessCampaignStatus, error) {
	awarenessCampaignStatus, err := acss.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return awarenessCampaignStatus, nil
}

func (acss *AwarenessCampaignStatusService) GetByType(status string) (*models.AwarenessCampaignStatus, error) {
	awarenessCampaignStatus, err := acss.repo.GetByType(status)

	if err != nil {
		return nil, err
	}

	return awarenessCampaignStatus, nil
}