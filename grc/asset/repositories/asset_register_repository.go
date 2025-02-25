package repositories

import (
	"fmt"
	"jethur/grc/asset/dtos"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/utility"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// Interfaces
type AssetRegisterRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *dtos.AssetRegisterListQueryParamsDTO) (*dtos.AssetRegisterListResponseDTO, error)
	GetByID(id uint) (*models.AssetRegister, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AssetRegister, error)
	Create(assetRegister *models.AssetRegister, tx ...utility.DatabaseTransaction) (uint, error)
	Update(assetRegister *models.AssetRegister, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error

	GenerateReferenceNumber() (string, error)

	CreateAssetControls(controls []*models.AssetControl, tx ...utility.DatabaseTransaction) error
	DeleteAssetControls(controlIds []*uint, tx ...utility.DatabaseTransaction) error

	GetControlSummary(assetID uint) ([]models.AssetControl, error)
	CreateAssetDocuments(documents []*models.AssetDocument, tx ...utility.DatabaseTransaction) error
	DeleteAssetDocuments(documentIds []*uint, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.AssetDocument, error)
}

// Implementation
type AssetRegisterRepository struct {
	db *gorm.DB
}

func NewAssetRegisterRepository(db *gorm.DB) *AssetRegisterRepository {
	return &AssetRegisterRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *AssetRegisterRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AssetRegisterRepository) GetAll(params *dtos.AssetRegisterListQueryParamsDTO) (*dtos.AssetRegisterListResponseDTO, error) {
	var assets []*dtos.AssetRegisterListDTO
	var totalCount int64

	query := repo.db.Table("asset_register").
		Select(`
			asset_register.id,
			asset_register.reference_number,
			asset_register.title_en,
			asset_register.title_ar,
			asset_register.location,
			asset_register.description,
			CONVERT(VARCHAR, asset_register.acquisition_date, 23) AS acquisition_date,
			CONVERT(VARCHAR, asset_register.assignment_date, 23) AS assignment_date,
			CONVERT(VARCHAR, asset_register.warranty_expiry_date, 23) AS warranty_expiry_date,
			asset_register.manufacturer,
			asset_register.model,
			asset_register.serial_number,
			asset_types.title_en AS asset_type_title_en,
			asset_types.title_ar AS asset_type_title_ar,
			asset_categories.title_en AS asset_category_title_en,
			asset_categories.title_ar AS asset_category_title_ar,
			asset_sub_categories.title_en AS asset_sub_category_title_en,
			asset_sub_categories.title_ar AS asset_sub_category_title_ar,
			asset_statuses.title_en AS asset_status_title_en,
			asset_statuses.title_ar AS asset_status_title_ar,
			asset_statuses.color_code AS asset_status_color_code,
			asset_register.total_score,
			asset_criticality_matrix_rating_scales.title_en AS asset_criticality_matrix_rating_scale_title_en,
			asset_criticality_matrix_rating_scales.title_ar AS asset_criticality_matrix_rating_scale_title_ar,
			asset_criticality_matrix_rating_scales.color_code AS asset_criticality_matrix_rating_scale_color_code,
			organization_units.title_en AS responsible_unit_title_en,
			organization_units.title_ar AS responsible_unit_title_ar,
			responsible_user.id AS responsible_user_id,
			responsible_user.first_name_en AS responsible_user_first_name_en,
			responsible_user.first_name_ar AS responsible_user_first_name_ar,
			responsible_user.last_name_en AS responsible_user_last_name_en,
			responsible_user.last_name_ar AS responsible_user_last_name_ar,
			responsible_user.initial AS responsible_user_initial,
			responsible_user.color AS responsible_user_color,
			responsible_user.email AS responsible_user_email,
			responsible_user.profile_image_path AS responsible_user_profile_image_path,
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
		Joins("LEFT JOIN asset_types ON asset_types.id = asset_register.asset_type_id").
		Joins("LEFT JOIN asset_categories ON asset_categories.id = asset_register.asset_category_id").
		Joins("LEFT JOIN asset_sub_categories ON asset_sub_categories.id = asset_register.asset_sub_category_id").
		Joins("LEFT JOIN asset_statuses ON asset_statuses.id = asset_register.asset_status_id").
		Joins("LEFT JOIN asset_criticality_matrix_rating_scales ON asset_criticality_matrix_rating_scales.id = asset_register.asset_criticality_matrix_rating_scale_id").
		Joins("LEFT JOIN organization_units ON organization_units.id = asset_register.responsible_unit_id").
		Joins("LEFT JOIN users AS responsible_user ON responsible_user.id = asset_register.responsible_user_id").
		Joins("LEFT JOIN users AS created_user ON created_user.id = asset_register.created_by_id").
		Joins("LEFT JOIN users AS updated_user ON updated_user.id = asset_register.updated_by_id")

	// Apply filters
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where("asset_register."+columnName+" IN (?)", values)
		}
	}

	// Apply acquisition_date filters
	if params.AcquisitionDateFrom != "" {
		query = query.Where("asset_register.acquisition_date >= ?", params.AcquisitionDateFrom)
	}
	if params.AcquisitionDateTo != "" {
		query = query.Where("asset_register.acquisition_date <= ?", params.AcquisitionDateTo)
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		query = query.Where(`
			(asset_register.title_en LIKE ? OR
			asset_register.title_ar LIKE ? OR
			asset_register.reference_number LIKE ? OR
			asset_register.description LIKE ? OR
			asset_types.title_en LIKE ? OR
			asset_types.title_ar LIKE ? OR
			asset_register.manufacturer LIKE ? OR
			asset_register.model LIKE ? OR
			asset_register.serial_number LIKE ? OR
			asset_categories.title_en LIKE ? OR
			asset_categories.title_ar LIKE ? OR
			asset_statuses.title_en LIKE ? OR
			asset_statuses.title_ar LIKE ?)`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%")
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
	err := query.Find(&assets).Error
	if err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := 1
	if !params.Export {
		totalPages = int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	}

	// Prepare response
	response := &dtos.AssetRegisterListResponseDTO{
		Items: assets,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}

	return response, nil
}

func (repo *AssetRegisterRepository) GetByID(id uint) (*models.AssetRegister, error) {
	var assetRegister *models.AssetRegister
	result := repo.db.First(&assetRegister, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return assetRegister, nil
}

func (repo *AssetRegisterRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AssetRegister, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var asset *models.AssetRegister

	result := db.
		Preload("AssetType").
		Preload("AssetCategory").
		Preload("AssetSubCategory").
		Preload("AssetStatus").
		Preload("ResponsibleUser").
		Preload("ResponsibleUser.OrganizationUnit1").
		Preload("ResponsibleUnit").
		Preload("Controls").
		Preload("AssetCriticalityMatrixRatingScale").
		Preload("AssetCriticalityMatrix").
		Preload("Documents").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&asset, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return asset, nil
}

func (repo *AssetRegisterRepository) Create(assetRegister *models.AssetRegister, tx ...utility.DatabaseTransaction) (uint, error) {
	if err := repo.db.Create(assetRegister).Error; err != nil {
		return 0, err
	}
	return assetRegister.ID, nil
}

func (repo *AssetRegisterRepository) Update(assetRegister *models.AssetRegister, tx ...utility.DatabaseTransaction) error {
	if err := repo.db.Save(assetRegister).Error; err != nil {
		return err
	}
	return nil
}

func (repo *AssetRegisterRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var assetRegister *models.AssetRegister
	if err := db.First(&assetRegister, id).Error; err != nil {
		return err
	}
	return db.Delete(&models.AssetRegister{}, id).Error
}

func (repo *AssetRegisterRepository) GenerateReferenceNumber() (string, error) {
	var lastReferenceNumber string
	err := repo.db.Model(&models.AssetRegister{}).
		Select("reference_number").
		Order("reference_number desc").
		Limit(1).
		Pluck("reference_number", &lastReferenceNumber).
		Error

	if err != nil {
		return "", err
	}

	lastDigits := "0000"
	if lastReferenceNumber != "" && len(lastReferenceNumber) >= 4 {
		lastDigits = lastReferenceNumber[len(lastReferenceNumber)-4:]
	}

	currentYear := time.Now().Year()
	lastReferenceID, err := strconv.ParseUint(lastDigits, 10, 64)
	if err != nil {
		return "", err
	}
	newReferenceNumber := fmt.Sprintf("%s-%d-%04d", "AST", currentYear, lastReferenceID+1)

	return newReferenceNumber, nil
}

func (repo *AssetRegisterRepository) CreateAssetControls(controls []*models.AssetControl, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&controls).Error; err != nil {
		return err
	}
	return nil
}

func (repo *AssetRegisterRepository) DeleteAssetControls(controlIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete controls by IDs
	if err := db.Where("id IN (?)", controlIds).Delete(&models.AssetControl{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AssetRegisterRepository) GetControlSummary(assetID uint) ([]models.AssetControl, error) {
	var assetControls []models.AssetControl

	result := repo.db.
		Preload("Control").
		Preload("Control.ControlImplementationStatus").
		Preload("Control.ControlEfficiencyStatus").
		Where("asset_register_id = ?", assetID).
		Find(&assetControls)
	if result.Error != nil {
		return nil, result.Error
	}

	return assetControls, nil
}

func (repo *AssetRegisterRepository) CreateAssetDocuments(documents []*models.AssetDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&documents).Error; err != nil {
		return err
	}
	return nil
}

func (repo *AssetRegisterRepository) DeleteAssetDocuments(documentIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", documentIds).Delete(&models.AssetDocument{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AssetRegisterRepository) GetDocumentByID(id uint) (*models.AssetDocument, error) {
	var document models.AssetDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}
