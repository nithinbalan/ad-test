package services

import (
	"errors"
	"jethur/grc/models"
	"jethur/grc/organization/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type CommitteeTermServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string, typeFilter string) ([]*models.CommitteeTerm, error)
	GetByID(id uint) (*models.CommitteeTerm, error)
	GetByType(typeFilter string) (*models.CommitteeTerm, error)
}

// Service and implementations
type CommitteeTermService struct {
	repo repositories.CommitteeTermRepositoryInterface
}

func NewCommitteeTermService(repo repositories.CommitteeTermRepositoryInterface) CommitteeTermServiceInterface {
	return &CommitteeTermService{
		repo: repo,
	}
}

func (cts *CommitteeTermService) GetAll(sort string, sortBy string, keyword string, status string, typeFilter string) ([]*models.CommitteeTerm, error) {
	return cts.repo.GetAll(sort, sortBy, keyword, status, typeFilter)
}

func (cts CommitteeTermService) GetByID(id uint) (*models.CommitteeTerm, error) {
	committeeTerm, err := cts.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return committeeTerm, nil

}

func (cts *CommitteeTermService) GetByType(typeFilter string) (*models.CommitteeTerm, error) {
	committeeTerm, err := cts.repo.GetByType(typeFilter)
	if err != nil {
		return nil, err
	}

	return committeeTerm, nil
}
