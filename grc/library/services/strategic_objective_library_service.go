package services

import (
	"os"

	libraryDTOs "jethur/grc/library/dtos"
	"jethur/grc/library/repositories"
	"jethur/utility"
)

type StrategicObjectiveLibraryServiceInterface interface {
	GetAll(params *libraryDTOs.StrategicObjectiveLibraryListQueryParamsDTO) (*libraryDTOs.StrategicObjectiveLibraryListResponseDTO, error)
	Export(data []*libraryDTOs.StrategicObjectiveLibraryDTO, userID uint) (*os.File, error)
}

// Service and implementations
type StrategicObjectiveLibraryService struct {
	repo repositories.StrategicObjectiveLibraryRepositoryInterface
}

func NewStrategicObjectiveLibraryService(
	repo repositories.StrategicObjectiveLibraryRepositoryInterface,
) StrategicObjectiveLibraryServiceInterface {
	return &StrategicObjectiveLibraryService{
		repo: repo,
	}
}

func (sols *StrategicObjectiveLibraryService) GetAll(params *libraryDTOs.StrategicObjectiveLibraryListQueryParamsDTO) (*libraryDTOs.StrategicObjectiveLibraryListResponseDTO, error) {
	return sols.repo.GetAll(params)
}

func (sols *StrategicObjectiveLibraryService) Export(data []*libraryDTOs.StrategicObjectiveLibraryDTO, userID uint) (*os.File, error) {
	headers := []string{"SL No.", "Title", "Category"}
	rows := make([][]interface{}, len(data))
	for i, k := range data {
		rows[i] = []interface{}{i + 1, k.Title, k.Category}
	}
	return utility.GenerateExcelFile(headers, rows)
}
