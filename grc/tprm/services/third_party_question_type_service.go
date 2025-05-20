package services

import (
	"jethur/grc/models"
	"jethur/grc/tprm/repositories"
)

type ThirdPartyQuestionTypeServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ThirdPartyQuestionType, error)
	GetByID(id uint) (*models.ThirdPartyQuestionType, error)
	GetByType(questionType string) (*models.ThirdPartyQuestionType, error)
}

// Service and implementations
type ThirdPartyQuestionTypeService struct {
	repo repositories.ThirdPartyQuestionTypeRepositoryInterface
}

func NewThirdPartyQuestionTypeService(repo repositories.ThirdPartyQuestionTypeRepositoryInterface) ThirdPartyQuestionTypeServiceInterface {
	return &ThirdPartyQuestionTypeService{
		repo: repo,
	}
}

func (tpqts *ThirdPartyQuestionTypeService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ThirdPartyQuestionType, error) {
	return tpqts.repo.GetAll(sort, sortBy, keyword, status)
}

func (tpqts *ThirdPartyQuestionTypeService) GetByID(id uint) (*models.ThirdPartyQuestionType, error) {
	return tpqts.repo.GetByID(id)
}

func (tpqts *ThirdPartyQuestionTypeService) GetByType(questionType string) (*models.ThirdPartyQuestionType, error) {
	result, err := tpqts.repo.GetByType(questionType)
	if err != nil {
		return nil, err
	}
	return result, nil
}
