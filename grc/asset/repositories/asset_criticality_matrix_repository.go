package repositories

import (
	// "jethur/grc/asset/dtos"
	"jethur/grc/asset/dtos"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AssetCriticalityMatrixRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *dtos.AssetCriticalityMatrixListQueryParamsDTO) (*dtos.AssetCriticalityMatrixListResponseDTO, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.AssetCriticalityMatrix, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AssetCriticalityMatrix, error)
	Create(matrix *models.AssetCriticalityMatrix, tx ...utility.DatabaseTransaction) (uint, error)
	Update(matrix *models.AssetCriticalityMatrix, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error

	CreateCategories(categories []*models.AssetCriticalityMatrixAssetCategory, tx ...utility.DatabaseTransaction) error
	DeleteCategories(matrixID uint, categoryIDs []uint, tx ...utility.DatabaseTransaction) error
	DeleteCategoriesByMatrixID(matrixID uint, tx ...utility.DatabaseTransaction) error
	// GetCriticalityMatrixByCategoryID(assetCategoryID uint) (*uint, error)
	GetCriticalityMatrixByCategoryID(categoryID uint) (uint, error)
}

type AssetCriticalityMatrixRepository struct {
	db *gorm.DB
}

func NewAssetCriticalityMatrixRepository(db *gorm.DB) AssetCriticalityMatrixRepositoryInterface {
	return &AssetCriticalityMatrixRepository{
		db: db,
	}
}

func (repo *AssetCriticalityMatrixRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AssetCriticalityMatrixRepository) GetAll(params *dtos.AssetCriticalityMatrixListQueryParamsDTO) (*dtos.AssetCriticalityMatrixListResponseDTO, error) {
	var matrices []*dtos.AssetCriticalityMatrixListDTO
	var totalCount int64

	query := repo.db.Table("asset_criticality_matrices").
		Select(`
			asset_criticality_matrices.id,
			asset_criticality_matrices.title_en,
			asset_criticality_matrices.title_ar,
			asset_criticality_matrices.status
		`)

	// Apply search filters
	if params.Keyword != "" {
		query = query.Where(`
			asset_criticality_matrices.title_en LIKE ? OR
			asset_criticality_matrices.title_ar LIKE ?`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	// Apply status filter
	if params.Status != "" {
		query = query.Where("asset_criticality_matrices.status = ?", params.Status)
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize)

	if err := query.Find(&matrices).Error; err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &dtos.AssetCriticalityMatrixListResponseDTO{
		Items: matrices,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}

	return response, nil
}

func (repo *AssetCriticalityMatrixRepository) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.AssetCriticalityMatrix, error) {
	var matrix models.AssetCriticalityMatrix
	db := repo.db
	if len(tx) > 0 {
		db = tx[0].DB()
	}

	if err := db.First(&matrix, id).Error; err != nil {
		return nil, err
	}
	return &matrix, nil
}

func (repo *AssetCriticalityMatrixRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AssetCriticalityMatrix, error) {
	var matrix models.AssetCriticalityMatrix
	db := repo.db
	if len(tx) > 0 {
		db = tx[0].DB()
	}

	if err := db.
		Preload("AssetCategories.AssetCategory").
		First(&matrix, id).Error; err != nil {
		return nil, err
	}
	return &matrix, nil
}

func (repo *AssetCriticalityMatrixRepository) Create(matrix *models.AssetCriticalityMatrix, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 {
		db = tx[0].DB()
	}
	if err := db.Create(matrix).Error; err != nil {
		return 0, err
	}
	return matrix.ID, nil
}

func (repo *AssetCriticalityMatrixRepository) Update(matrix *models.AssetCriticalityMatrix, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 {
		db = tx[0].DB()
	}
	return db.Save(matrix).Error
}

func (repo *AssetCriticalityMatrixRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 {
		db = tx[0].DB()
	}
	return db.Delete(&models.AssetCriticalityMatrix{}, id).Error
}

func (repo *AssetCriticalityMatrixRepository) CreateCategories(categories []*models.AssetCriticalityMatrixAssetCategory, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 {
		db = tx[0].DB()
	}
	return db.Create(&categories).Error
}

func (repo *AssetCriticalityMatrixRepository) DeleteCategoriesByMatrixID(matrixID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 {
		db = tx[0].DB()
	}
	return db.Where("asset_criticality_matrix_id = ?", matrixID).Delete(&models.AssetCriticalityMatrixAssetCategory{}).Error
}

func (repo *AssetCriticalityMatrixRepository) DeleteCategories(matrixID uint, categoryIDs []uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("asset_criticality_matrix_id = ? AND asset_category_id IN (?)", matrixID, categoryIDs).
		Delete(&models.AssetCriticalityMatrixAssetCategory{}).Error
}

// func (repo *AssetCriticalityMatrixRepository) GetCriticalityMatrixByCategoryID(assetCategoryID uint) (*uint, error) {
// 	var matrixID uint
// 	err := repo.db.
// 		Table("asset_criticality_matrix_asset_categories").
// 		Select("asset_criticality_matrix_id").
// 		Where("asset_category_id = ?", assetCategoryID).
// 		First(&matrixID).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &matrixID, nil
// }

func (repo *AssetCriticalityMatrixRepository) GetCriticalityMatrixByCategoryID(categoryID uint) (uint, error) {
	var matrixID uint
	err := repo.db.Table("asset_criticality_matrix_asset_categories").
		Select("asset_criticality_matrix_id").
		Where("asset_category_id = ?", categoryID).
		Order("id").
		Limit(1).
		Scan(&matrixID).Error

	if err != nil {
		return 0, err
	}
	return matrixID, nil
}
