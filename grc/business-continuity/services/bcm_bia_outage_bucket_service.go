package services

import (
	"errors"
	"jethur/grc/business-continuity/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type BCMBIAOutageBucketServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.BCMBIAOutageBucket, error)
	GetByID(id uint) (*models.BCMBIAOutageBucket, error)
}

// Service and implementations
type BCMBIAOutageBucketService struct {
	repo repositories.BCMBIAOutageBucketRepositoryInterface
}

func NewBCMBIAOutageBucketService(repo repositories.BCMBIAOutageBucketRepositoryInterface) BCMBIAOutageBucketServiceInterface {
	return &BCMBIAOutageBucketService{
		repo: repo,
	}
}

func (bbobs *BCMBIAOutageBucketService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.BCMBIAOutageBucket, error) {
	return bbobs.repo.GetAll(sort, sortBy, keyword, status)
}

func (bbobs BCMBIAOutageBucketService) GetByID(id uint) (*models.BCMBIAOutageBucket, error) {
	bcmBiaOutageBucket, err := bbobs.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return bcmBiaOutageBucket, nil
}
