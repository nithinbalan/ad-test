package repositories

import (
	// "jethur/grc/asset/dtos"

	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AssetCriticalityMatrixRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll() ([]*models.AssetCriticalityMatrix, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.AssetCriticalityMatrix, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AssetCriticalityMatrix, error)
	Create(matrix *models.AssetCriticalityMatrix, tx ...utility.DatabaseTransaction) (uint, error)
	Update(matrix *models.AssetCriticalityMatrix, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error

	CreateCategories(categories []*models.AssetCriticalityMatrixAssetCategory, tx ...utility.DatabaseTransaction) error
	DeleteCategories(matrixID uint, categoryIDs []uint, tx ...utility.DatabaseTransaction) error
	DeleteCategoriesByMatrixID(matrixID uint, tx ...utility.DatabaseTransaction) error
	ExistsByCategoryID(categoryID uint, excludeMatrixID uint, tx ...utility.DatabaseTransaction) (bool, error)
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

func (repo *AssetCriticalityMatrixRepository) GetAll() ([]*models.AssetCriticalityMatrix, error) {
	var matrices []*models.AssetCriticalityMatrix

	err := repo.db.
		Preload("Factors").
		Preload("FactorRatings").
		Preload("RatingScales").
		Preload("AssetCategories.AssetCategory").
		Find(&matrices).Error

	if err != nil {
		return nil, err
	}

	return matrices, nil
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
		Preload("Factors").
		Preload("FactorRatings").
		Preload("RatingScales").
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

func (repo *AssetCriticalityMatrixRepository) ExistsByCategoryID(categoryID uint, excludeMatrixID uint, tx ...utility.DatabaseTransaction) (bool, error) {
	db := repo.db
	if len(tx) > 0 {
		db = tx[0].DB()
	}

	var count int64
	query := db.Table("asset_criticality_matrix_asset_categories").
		Where("asset_category_id = ?", categoryID)

	if excludeMatrixID > 0 {
		query = query.Where("asset_criticality_matrix_id != ?", excludeMatrixID)
	}

	err := query.Count(&count).Error

	return count > 0, err
}

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
