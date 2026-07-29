package services

import (
	"errors"
	bcmDTOs "jethur/grc/business-continuity/dtos"
	"jethur/grc/business-continuity/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type BCMBIAImpactLevelServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*bcmDTOs.BCMBIAImpactLevelDTO, error)
	GetByID(id uint) (*models.BCMBIAImpactLevel, error)
	GetFullDetailsByID(id uint) (*bcmDTOs.BCMBIAImpactLevelDetailDTO, error)
	GetByImpactAnalysisLevelID(impactAnalysisLevelID uint) (*models.BCMBIAImpactLevel, error)
	GetByImpactAnalysisLevelIDAndOutageBucketID(impactAnalysisLevelID uint, outageBucketID uint) (*models.BCMBIAImpactLevel, error)
}

// Service and implementations
type BCMBIAImpactLevelService struct {
	repo repositories.BCMBIAImpactLevelRepositoryInterface
}

func NewBCMBIAImpactLevelService(repo repositories.BCMBIAImpactLevelRepositoryInterface) BCMBIAImpactLevelServiceInterface {
	return &BCMBIAImpactLevelService{
		repo: repo,
	}
}

func (bbils *BCMBIAImpactLevelService) GetAll(sort string, sortBy string, keyword string, status string) ([]*bcmDTOs.BCMBIAImpactLevelDTO, error) {
	return bbils.repo.GetAll(sort, sortBy, keyword, status)
}

func (bbils *BCMBIAImpactLevelService) GetByID(id uint) (*models.BCMBIAImpactLevel, error) {
	bcmBiaImpactLevel, err := bbils.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return bcmBiaImpactLevel, nil
}

func (bbils *BCMBIAImpactLevelService) GetFullDetailsByID(id uint) (*bcmDTOs.BCMBIAImpactLevelDetailDTO, error) {
	var detailDTO bcmDTOs.BCMBIAImpactLevelDetailDTO

	bcmBiaImpactLevel, err := bbils.repo.GetFullDetailsByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	if err := mapstructure.Decode(bcmBiaImpactLevel, &detailDTO); err != nil {
		return nil, err
	}

	return &detailDTO, nil
}

func (bbils *BCMBIAImpactLevelService) GetByImpactAnalysisLevelID(impactAnalysisLevelID uint) (*models.BCMBIAImpactLevel, error) {
	bcmBiaImpactLevel, err := bbils.repo.GetByImpactAnalysisLevelID(impactAnalysisLevelID)

	if err != nil {
		return nil, err
	}

	return bcmBiaImpactLevel, nil
}

func (bbils *BCMBIAImpactLevelService) GetByImpactAnalysisLevelIDAndOutageBucketID(impactAnalysisLevelID uint, outageBucketID uint) (*models.BCMBIAImpactLevel, error) {
	bcmBiaImpactLevel, err := bbils.repo.GetByImpactAnalysisLevelIDAndOutageBucketID(impactAnalysisLevelID, outageBucketID)

	if err != nil {
		return nil, err
	}

	return bcmBiaImpactLevel, nil
}
