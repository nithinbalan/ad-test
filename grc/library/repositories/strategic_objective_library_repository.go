package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	libraryDTOs "jethur/grc/library/dtos"

	"gorm.io/gorm"
)

type StrategicObjectiveLibraryRepositoryInterface interface {
	GetAll(params *libraryDTOs.StrategicObjectiveLibraryListQueryParamsDTO) (*libraryDTOs.StrategicObjectiveLibraryListResponseDTO, error)
	GetCategoryCounts() ([]libraryDTOs.StrategicCategoryCountDTO, error)
}

type StrategicObjectiveLibraryRepository struct {
	db *gorm.DB
}

func NewStrategicObjectiveLibraryRepository(db *gorm.DB) *StrategicObjectiveLibraryRepository {
	return &StrategicObjectiveLibraryRepository{
		db: db,
	}
}

func (repo *StrategicObjectiveLibraryRepository) GetAll(params *libraryDTOs.StrategicObjectiveLibraryListQueryParamsDTO) (*libraryDTOs.StrategicObjectiveLibraryListResponseDTO, error) {
	var StrategicObjective []*libraryDTOs.StrategicObjectiveLibraryDTO
	var totalCount int64

	query := repo.db.Table("strategic_objective_library").
		Select("id, title, category")

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			title LIKE ? OR
			category LIKE ?
		`, keyword, keyword)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	sortBy := params.SortBy
	if sortBy == "" {
		sortBy = "title"
	}
	sortOrder := params.Sort
	if sortOrder != "desc" {
		sortOrder = "asc"
	}
	query = query.Order(sortBy + " " + sortOrder)

	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	err := query.Find(&StrategicObjective).Error
	if err != nil {
		return nil, err
	}
	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	response := &libraryDTOs.StrategicObjectiveLibraryListResponseDTO{
		Items: StrategicObjective,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *StrategicObjectiveLibraryRepository) GetCategoryCounts() ([]libraryDTOs.StrategicCategoryCountDTO, error) {
	var results []libraryDTOs.StrategicCategoryCountDTO

	err := repo.db.
		Table("strategic_objective_library").
		Select("category, COUNT(*) as count").
		Group("category").
		Order("category ASC").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}
