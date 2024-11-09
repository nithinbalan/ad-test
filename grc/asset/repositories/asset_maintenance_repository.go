package repositories

import (
	"jethur/grc/asset/dtos"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type AssetMaintenanceRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *dtos.AssetMaintenanceListQueryParamsDTO) (*dtos.AssetMaintenanceListResponseDTO, error)
	GetByID(id uint) (*models.AssetMaintenance, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AssetMaintenance, error)
	Create(assetMaintenance *models.AssetMaintenance, tx ...utility.DatabaseTransaction) (uint, error)
	Update(assetMaintenance *models.AssetMaintenance, tx ...utility.DatabaseTransaction) error
	Delete(id uint) error
}

// Implementation
type AssetMaintenanceRepository struct {
	db *gorm.DB
}

func NewAssetMaintenanceRepository(db *gorm.DB) *AssetMaintenanceRepository {
	return &AssetMaintenanceRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *AssetMaintenanceRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AssetMaintenanceRepository) GetAll(params *dtos.AssetMaintenanceListQueryParamsDTO) (*dtos.AssetMaintenanceListResponseDTO, error) {
	var maintenances []*dtos.AssetMaintenanceListDTO
	var totalCount int64

	query := repo.db.Table("asset_maintenance").
		Select(`
			asset_maintenance.id,
			asset_maintenance.type,
			CONVERT(VARCHAR, asset_maintenance.performed_date, 23) AS performed_date,
			asset_maintenance.title,
			asset_maintenance.description,
			asset_maintenance.outcome,
			asset_register.id AS asset_register_id,
			performed_user.id AS performed_by_id,
			performed_user.first_name_en AS performed_by_first_name_en,
			performed_user.first_name_ar AS performed_by_first_name_ar,
			performed_user.last_name_en AS performed_by_last_name_en,
			performed_user.last_name_ar AS performed_by_last_name_ar,
			performed_user.initial AS performed_by_initial,
			performed_user.color AS performed_by_color,
			performed_user.email AS performed_by_email,
			performed_user.profile_image_path AS performed_by_profile_image_path,
			created_user.id AS created_user_id,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.email AS created_user_email,
			created_user.profile_image_path AS created_user_profile_image_path,
			updated_user.id AS updated_user_id,
			updated_user.first_name_en AS updated_user_first_name_en,
			updated_user.first_name_ar AS updated_user_first_name_ar,
			updated_user.last_name_en AS updated_user_last_name_en,
			updated_user.last_name_ar AS updated_user_last_name_ar,
			updated_user.initial AS updated_user_initial,
			updated_user.color AS updated_user_color,
			updated_user.email AS updated_user_email,
			updated_user.profile_image_path AS updated_user_profile_image_path
		`).
		Joins("LEFT JOIN asset_register ON asset_register.id = asset_maintenance.asset_register_id").
		Joins("LEFT JOIN users AS performed_user ON performed_user.id = asset_maintenance.performed_by_id").
		Joins("LEFT JOIN users AS created_user ON created_user.id = asset_maintenance.created_by_id").
		Joins("LEFT JOIN users AS updated_user ON updated_user.id = asset_maintenance.updated_by_id")

	// Apply filters
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where("asset_maintenance."+columnName+" IN (?)", values)
		}
	}

	// Apply performed_date filters
	if params.PerformedDateFrom != "" {
		query = query.Where("asset_maintenance.performed_date >= ?", params.PerformedDateFrom)
	}
	if params.PerformedDateTo != "" {
		query = query.Where("asset_maintenance.performed_date <= ?", params.PerformedDateTo)
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		query = query.Where(`
			(asset_maintenance.title LIKE ? OR
			asset_maintenance.description LIKE ? OR
			performed_user.first_name_en LIKE ? OR
			performed_user.first_name_ar LIKE ? OR
			performed_user.last_name_en LIKE ? OR
			performed_user.last_name_ar LIKE ?)`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply sorting and pagination
	query = query.Order(params.SortBy + " " + params.Sort)
	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	// Execute query
	err := query.Find(&maintenances).Error
	if err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := 1
	if !params.Export {
		totalPages = int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	}

	// Prepare response
	response := &dtos.AssetMaintenanceListResponseDTO{
		Items: maintenances,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}

	return response, nil
}

func (repo *AssetMaintenanceRepository) GetByID(id uint) (*models.AssetMaintenance, error) {
	var assetMaintenance *models.AssetMaintenance
	result := repo.db.First(&assetMaintenance, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return assetMaintenance, nil
}

func (repo *AssetMaintenanceRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AssetMaintenance, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var maintenance *models.AssetMaintenance

	result := db.
		Preload("AssetRegister").
		Preload("PerformedBy").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&maintenance, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return maintenance, nil
}

func (repo *AssetMaintenanceRepository) Create(assetMaintenance *models.AssetMaintenance, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(assetMaintenance).Error; err != nil {
		return 0, err
	}
	return assetMaintenance.ID, nil
}

func (repo *AssetMaintenanceRepository) Update(assetMaintenance *models.AssetMaintenance, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Save(assetMaintenance).Error; err != nil {
		return err
	}
	return nil
}

func (repo *AssetMaintenanceRepository) Delete(id uint) error {
	var assetMaintenance *models.AssetMaintenance
	if err := repo.db.First(&assetMaintenance, id).Error; err != nil {
		return err
	}
	return repo.db.Delete(&models.AssetMaintenance{}, id).Error
}
