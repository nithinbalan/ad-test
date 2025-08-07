package services

import (
	"os"

	libraryDTOs "jethur/grc/library/dtos"
	"jethur/grc/library/repositories"
	"jethur/utility"
)

type RiskLibraryServiceInterface interface {
	GetAll(params *libraryDTOs.RiskLibraryListQueryParamsDTO) (*libraryDTOs.RiskLibraryListResponseDTO, error)
	GetRiskCategoryCounts() ([]libraryDTOs.RiskCategoryCountDTO, error)
	GetRiskLikelihoodCounts() ([]libraryDTOs.RiskLikelihoodCountDTO, error)
	GetRiskImpactCounts() ([]libraryDTOs.RiskImpactCountDTO, error)
	ExportRiskLibrary(data []*libraryDTOs.RiskLibraryDTO, userID uint) (*os.File, error)
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

func (rls *RiskLibraryService) GetRiskLikelihoodCounts() ([]libraryDTOs.RiskLikelihoodCountDTO, error) {
	return rls.repo.GetRiskLikelihoodCounts()
}

func (rls *RiskLibraryService) GetRiskImpactCounts() ([]libraryDTOs.RiskImpactCountDTO, error) {
	return rls.repo.GetRiskImpactCounts()
}

func (rls *RiskLibraryService) ExportRiskLibrary(data []*libraryDTOs.RiskLibraryDTO, userID uint) (*os.File, error) {
	// Define headers
	headers := []string{
		"SL No.",
		"Title",
		"Description",
		"Category",
		"Likelihood",
		"Impact",
	}

	// Prepare data rows
	rows := make([][]interface{}, len(data))
	for i, risk := range data {
		rows[i] = []interface{}{
			i + 1,
			risk.Title,
			utility.GetStringValue(risk.Description),
			risk.Category,
			utility.GetStringValue(risk.Likelihood),
			utility.GetStringValue(risk.Impact),
		}
	}

	return utility.GenerateExcelFile(headers, rows)
}
