package repositories

import (
	"jethur/grc/asset/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AssetControlAnalysisRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(assetID uint, params *dtos.AssetControlAnalysisListQueryParamsDTO) ([]*dtos.AssetControlAnalysisListDTO, error)
	GetByID(id uint) (*models.AssetControl, error)
	GetFullDetailsByID(id uint) (*models.AssetControl, error)
	Create(assetControlAnalysis *models.AssetControl) (uint, error)
	Update(assetControlAnalysis *models.AssetControl, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	// GetAllControlsByAssetID(assetID uint) ([]*models.AssetControl, error)
}

type AssetControlAnalysisRepository struct {
	db *gorm.DB
}

func NewAssetControlAnalysisRepository(db *gorm.DB) *AssetControlAnalysisRepository {
	return &AssetControlAnalysisRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *AssetControlAnalysisRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AssetControlAnalysisRepository) GetAll(assetID uint, params *dtos.AssetControlAnalysisListQueryParamsDTO) ([]*dtos.AssetControlAnalysisListDTO, error) {
	var assetControlAnalyses []*dtos.AssetControlAnalysisListDTO

	query := repo.db.Table("asset_controls").
		Select(`
            asset_controls.id,
			asset_controls.comments,
			controls.id AS control_id,
			controls.code AS control_code,
			controls.reference_number AS control_reference_number,
			controls.description AS control_description,
			controls.title_en AS control_title_en,
			controls.title_ar AS control_title_ar,
            control_types.title_en AS control_type_title_en,
			control_types.title_ar AS control_type_title_ar,
			control_categories.title_en AS control_category_title_en,
			control_categories.title_ar AS control_category_title_ar,
			organization_units.title_en AS owner_unit_title_en,
			organization_units.title_ar AS owner_unit_title_ar,
			owner.id AS owner_id,
			owner.first_name_en AS owner_first_name_en,
			owner.first_name_ar AS owner_first_name_ar,
			owner.last_name_en AS owner_last_name_en,
			owner.last_name_ar AS owner_last_name_ar,
			owner.initial AS owner_initial,
            owner.color AS owner_color,
            owner.email AS owner_email,
            owner.profile_image_path AS owner_profile_image_path,
			control_implementation_statuses.title_en AS control_implementation_status_title_en,
			control_implementation_statuses.title_ar AS control_implementation_status_title_ar,
			control_implementation_statuses.color_code AS control_implementation_status_color_code,
			control_efficiency_statuses.title_en AS control_efficiency_status_title_en,
			control_efficiency_statuses.title_ar AS control_efficiency_status_title_ar,
			control_efficiency_statuses.color_code AS control_efficiency_status_color_code,
            created_user.id AS created_user_id,
            created_user.first_name_en AS created_user_first_name_en,
            created_user.first_name_ar AS created_user_first_name_ar,
            created_user.last_name_en AS created_user_last_name_en,
            created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
            created_user.color AS created_user_color,
            created_user.email AS created_user_email,
            created_user.profile_image_path AS created_user_profile_image_path,
			asset_controls.created_at AS created_at,
			updated_user.id AS updated_user_id,
			updated_user.first_name_en AS updated_user_first_name_en,
			updated_user.first_name_ar AS updated_user_first_name_ar,
			updated_user.last_name_en AS updated_user_last_name_en,
			updated_user.last_name_ar AS updated_user_last_name_ar,
			updated_user.initial AS updated_user_initial,
			updated_user.color AS updated_user_color,
			updated_user.email AS updated_user_email,
			updated_user.profile_image_path AS updated_user_profile_image_path,   
			asset_controls.updated_at AS updated_at
        `).
		Joins(`
			LEFT JOIN controls ON controls.id = asset_controls.control_id
			LEFT JOIN control_types ON control_types.id = controls.control_type_id
			LEFT JOIN control_categories ON control_categories.id = controls.control_category_id
			LEFT JOIN organization_units ON organization_units.id = controls.owner_unit_id
			LEFT JOIN users AS owner ON owner.id = controls.owner_id
            LEFT JOIN control_efficiency_statuses ON control_efficiency_statuses.id = asset_controls.control_efficiency_status_id
            LEFT JOIN control_implementation_statuses ON control_implementation_statuses.id = asset_controls.control_implementation_status_id
            LEFT JOIN users AS created_user ON created_user.id = asset_controls.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = asset_controls.updated_by_id
        `).
		Where(`asset_controls.asset_register_id = ?`, assetID)

	// Apply search keyword if exists
	if params.Keyword != "" {
		query = query.Where(`
            (controls.title_en LIKE ? OR
            asset_controls.control_efficiency_status_id LIKE ? OR
            asset_controls.control_implementation_status_id LIKE ? )`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	if params.SortBy == "" {
		params.SortBy = "created_at"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	// Execute the query
	if err := query.Find(&assetControlAnalyses).Error; err != nil {
		return nil, err
	}

	return assetControlAnalyses, nil
}

func (repo *AssetControlAnalysisRepository) GetByID(id uint) (*models.AssetControl, error) {
	var assetControlAnalysis *models.AssetControl

	result := repo.db.
		First(&assetControlAnalysis, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return assetControlAnalysis, nil
}

func (repo *AssetControlAnalysisRepository) GetFullDetailsByID(id uint) (*models.AssetControl, error) {
	var assetControlAnalysis *models.AssetControl

	result := repo.db.
		Preload("Control").
		Preload("ControlEfficiencyStatus").
		Preload("ControlImplementationStatus").
		First(&assetControlAnalysis, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return assetControlAnalysis, nil
}

func (repo *AssetControlAnalysisRepository) Create(assetControlAnalysis *models.AssetControl) (uint, error) {
	if err := repo.db.Create(assetControlAnalysis).Error; err != nil {
		return 0, err
	}

	return assetControlAnalysis.ID, nil
}

func (repo *AssetControlAnalysisRepository) Update(assetControlAnalysis *models.AssetControl, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(assetControlAnalysis).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AssetControlAnalysisRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Delete(&models.AssetControl{}, id).Error; err != nil {
		return err
	}

	return nil
}

// func (repo *AssetControlAnalysisRepository) GetAllControlsByAssetID(assetID uint) ([]*models.AssetControl, error) {
// 	var assetControls []*models.AssetControl

// 	err := repo.db.Where("asset_register_id = ?", assetID).Find(&assetControls).Error

// 	if err != nil {
// 		return nil, err
// 	}

// 	return assetControls, nil
// }
