package services

import (
	"errors"
	"jethur/grc/models"
	"jethur/grc/organization/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type CommitteeMemberTypeServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.CommitteeMemberType, error)
	GetByID(id uint) (*models.CommitteeMemberType, error)
}

type CommitteeMemberTypeService struct {
	repo repositories.CommitteeMemberTypeRepositoryInterface
}

func NewCommitteeMemberTypeService(repo repositories.CommitteeMemberTypeRepositoryInterface) CommitteeMemberTypeServiceInterface {
	return &CommitteeMemberTypeService{
		repo: repo,
	}
}

func (cmts *CommitteeMemberTypeService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.CommitteeMemberType, error) {
	return cmts.repo.GetAll(sort, sortBy, keyword, status)
}

func (cmts *CommitteeMemberTypeService) GetByID(id uint) (*models.CommitteeMemberType, error) {
	committeeMemberType, err := cmts.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return committeeMemberType, nil
}
