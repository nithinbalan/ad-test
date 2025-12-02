package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	libraryDTOs "jethur/grc/library/dtos"

	"gorm.io/gorm"
)

type KpiLibraryRepositoryInterface interface {
	GetAllKpi(params *libraryDTOs.KPILibraryListQueryParamsDTO) (*libraryDTOs.KPILibraryListResponseDTO, error)
	KPIGetCategoryCounts() ([]libraryDTOs.KPICategoryCountDTO, error)
}

type KpiLibraryRepository struct {
	db *gorm.DB
}

func NewKpiLibraryRepository(db *gorm.DB) *KpiLibraryRepository {
	return &KpiLibraryRepository{
		db: db,
	}
}

func (repo *KpiLibraryRepository) GetAllKpi(params *libraryDTOs.KPILibraryListQueryParamsDTO) (*libraryDTOs.KPILibraryListResponseDTO, error) {
	var KPIs []*libraryDTOs.KPILibraryDTO
	var totalCount int64

	query := repo.db.Table("kpi_library").
		Select("id, category, title, description, unit, frequency")

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			title LIKE ? OR
			description LIKE ? OR
			category LIKE ?
		`, keyword, keyword, keyword)
	}

	// Filter by category
	if len(params.Categories) > 0 {
		query = query.Where("category IN ?", params.Categories)
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

	err := query.Find(&KPIs).Error
	if err != nil {
		return nil, err
	}
	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	response := &libraryDTOs.KPILibraryListResponseDTO{
		Items: KPIs,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *KpiLibraryRepository) KPIGetCategoryCounts() ([]libraryDTOs.KPICategoryCountDTO, error) {
	var results []libraryDTOs.KPICategoryCountDTO

	err := repo.db.
		Table("kpi_library").
		Select("category, COUNT(*) as count").
		Group("category").
		Order("category ASC").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}
