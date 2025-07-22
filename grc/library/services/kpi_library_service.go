package services

import (
	"os"

	libraryDTOs "jethur/grc/library/dtos"
	"jethur/grc/library/repositories"
	"jethur/utility"
)

type KpiLibraryServiceInterface interface {
	GetAllKpi(params *libraryDTOs.KPILibraryListQueryParamsDTO) (*libraryDTOs.KPILibraryListResponseDTO, error)
	Export(data []*libraryDTOs.KPILibraryDTO, userID uint) (*os.File, error)
	KPIGetCategoryCounts() ([]libraryDTOs.KPICategoryCountDTO, error)
}

// Service and implementations
type KpiLibraryService struct {
	repo repositories.KpiLibraryRepositoryInterface
}

func NewKpiLibraryService(
	repo repositories.KpiLibraryRepositoryInterface,
) KpiLibraryServiceInterface {
	return &KpiLibraryService{
		repo: repo,
	}
}

func (kls *KpiLibraryService) GetAllKpi(params *libraryDTOs.KPILibraryListQueryParamsDTO) (*libraryDTOs.KPILibraryListResponseDTO, error) {
	return kls.repo.GetAllKpi(params)
}

func (kls *KpiLibraryService) Export(data []*libraryDTOs.KPILibraryDTO, userID uint) (*os.File, error) {
	headers := []string{"SL No.", "Title", "Description", "Category", "Unit", "Frequency"}
	rows := make([][]interface{}, len(data))
	for i, k := range data {
		rows[i] = []interface{}{i + 1, k.Title, utility.GetStringValue(k.Description), k.Category, k.Unit, k.Frequency}
	}
	return utility.GenerateExcelFile(headers, rows)
}

func (s *KpiLibraryService) KPIGetCategoryCounts() ([]libraryDTOs.KPICategoryCountDTO, error) {
	return s.repo.KPIGetCategoryCounts()
}
