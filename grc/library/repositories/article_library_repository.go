package repositories

import (
	"errors"

	commonDTOs "jethur/grc/common/dtos"
	libraryDTOs "jethur/grc/library/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ArticleLibraryRepositoryInterface interface {
	GetAll(params *libraryDTOs.ArticleLibraryListQueryParamsDTO) (*libraryDTOs.ArticleLibraryListResponseDTO, error)
	GetByID(id uint) (*models.ArticleLibrary, error)
}

type ArticleLibraryRepository struct {
	db *gorm.DB
}

func NewArticleLibraryRepository(db *gorm.DB) *ArticleLibraryRepository {
	return &ArticleLibraryRepository{
		db: db,
	}
}

func (repo *ArticleLibraryRepository) GetAll(params *libraryDTOs.ArticleLibraryListQueryParamsDTO) (*libraryDTOs.ArticleLibraryListResponseDTO, error) {
	var Article []*libraryDTOs.ArticleLibraryDTO
	var totalCount int64

	query := repo.db.Table("article_library").
		Select(`
			article_library.id,
			article_library.article_category_id AS category_id,
			article_category.title_en AS category_title_en,
			article_category.title_ar AS category_title_ar,
			article_library.image_path,
			article_library.title_en,
			article_library.title_ar,
			article_library.created_by,
			article_library.created_at
		`).
		Joins("LEFT JOIN article_categories AS article_category ON article_library.article_category_id = article_category.id")

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			article_library.title_en LIKE ? OR
			article_library.title_ar LIKE ? OR
			article_category.title_en LIKE ? OR
			article_category.title_ar LIKE ?  
		`, keyword, keyword, keyword, keyword)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	sortBy := params.SortBy
	if sortBy == "" {
		sortBy = "title_en"
	}
	sortOrder := params.Sort
	if sortOrder != "desc" {
		sortOrder = "asc"
	}
	query = query.Order(sortBy + " " + sortOrder)

	err := query.Find(&Article).Error
	if err != nil {
		return nil, err
	}
	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	response := &libraryDTOs.ArticleLibraryListResponseDTO{
		Items: Article,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *ArticleLibraryRepository) GetByID(id uint) (*models.ArticleLibrary, error) {
	var article *models.ArticleLibrary

	result := repo.db.
		Preload("ArticleCategory").
		First(&article, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}
	return article, nil
}
