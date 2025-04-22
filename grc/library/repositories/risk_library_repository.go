package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	libraryDTOs "jethur/grc/library/dtos"

	"gorm.io/gorm"
)

type RiskLibraryRepositoryInterface interface {
	GetAll(params *libraryDTOs.RiskLibraryListQueryParamsDTO) (*libraryDTOs.RiskLibraryListResponseDTO, error)
	GetRiskCategoryCounts() ([]libraryDTOs.RiskCategoryCountDTO, error)
}

type RiskLibraryRepository struct {
	db *gorm.DB
}

func NewRiskLibraryRepository(db *gorm.DB) *RiskLibraryRepository {
	return &RiskLibraryRepository{
		db: db,
	}
}

func (repo *RiskLibraryRepository) GetAll(params *libraryDTOs.RiskLibraryListQueryParamsDTO) (*libraryDTOs.RiskLibraryListResponseDTO, error) {
	var RiskLibrary []*libraryDTOs.RiskLibraryDTO
	var totalCount int64

	query := repo.db.Table("risk_library").
		Select(`
			id,
			category,
			title,
			description,
			likelihood,
			impact
		`)

	// Filter by keyword
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			title LIKE ? OR
			description LIKE ? OR
			impact LIKE ? OR
			likelihood LIKE ? OR
			category LIKE ?
		`, keyword, keyword, keyword, keyword, keyword)
	}

	// Filter by category
	if params.Category != "" {
		query = query.Where("category = ?", params.Category)
	}

	// Count total for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Sorting
	sortBy := params.SortBy
	if sortBy == "" {
		sortBy = "title"
	}
	sortOrder := params.Sort
	if sortOrder != "desc" {
		sortOrder = "asc"
	}
	query = query.Order(sortBy + " " + sortOrder)

	// Pagination
	offset := (params.Page - 1) * params.PageSize
	err := query.Offset(offset).Limit(params.PageSize).Find(&RiskLibrary).Error
	if err != nil {
		return nil, err
	}

	// Response
	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	response := &libraryDTOs.RiskLibraryListResponseDTO{
		Items: RiskLibrary,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}

	return response, nil
}

func (repo *RiskLibraryRepository) GetRiskCategoryCounts() ([]libraryDTOs.RiskCategoryCountDTO, error) {
	var results []libraryDTOs.RiskCategoryCountDTO

	err := repo.db.
		Table("risk_library").
		Select("category, COUNT(*) as count").
		Group("category").
		Order("category ASC").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}
