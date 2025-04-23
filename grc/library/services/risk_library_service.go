package services

import (
	libraryDTOs "jethur/grc/library/dtos"
	"jethur/grc/library/repositories"
)

type RiskLibraryServiceInterface interface {
	GetAll(params *libraryDTOs.RiskLibraryListQueryParamsDTO) (*libraryDTOs.RiskLibraryListResponseDTO, error)
	GetRiskCategoryCounts() ([]libraryDTOs.RiskCategoryCountDTO, error)
}

// Service and implementations
type RiskLibraryService struct {
	repo repositories.RiskLibraryRepositoryInterface
}

func NewRiskLibraryService(
	repo repositories.RiskLibraryRepositoryInterface,
) RiskLibraryServiceInterface {
	return &RiskLibraryService{
		repo: repo,
	}
}

func (rls *RiskLibraryService) GetAll(params *libraryDTOs.RiskLibraryListQueryParamsDTO) (*libraryDTOs.RiskLibraryListResponseDTO, error) {
	return rls.repo.GetAll(params)
}

func (rls *RiskLibraryService) GetRiskCategoryCounts() ([]libraryDTOs.RiskCategoryCountDTO, error) {
	return rls.repo.GetRiskCategoryCounts()
}
