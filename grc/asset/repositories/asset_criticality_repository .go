package repositories

import (
	"jethur/grc/asset/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AssetCriticalityRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(assetID uint, params *dtos.AssetCriticalityContributionListQueryParamsDTO) ([]*dtos.AssetCriticalityContributionListDTO, error)
	GetByID(id uint) (*models.AssetCriticalityContribution, error)
	Create(contribution *models.AssetCriticalityContribution, tx ...utility.DatabaseTransaction) (uint, error)
	Update(contribution *models.AssetCriticalityContribution, tx ...utility.DatabaseTransaction) error
	Delete(assetID uint, tx ...utility.DatabaseTransaction) error

	CalculateTotalScore(assetID uint, tx ...utility.DatabaseTransaction) (float64, error)
	GetAllContribution(assetID uint) ([]models.AssetCriticalityContribution, error)
}

type AssetCriticalityRepository struct {
	db *gorm.DB
}

func (repo *AssetCriticalityRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func NewAssetCriticalityRepository(db *gorm.DB) AssetCriticalityRepositoryInterface {
	return &AssetCriticalityRepository{
		db: db,
	}
}

func (repo *AssetCriticalityRepository) GetAll(assetID uint, params *dtos.AssetCriticalityContributionListQueryParamsDTO) ([]*dtos.AssetCriticalityContributionListDTO, error) {
	var contributions []*dtos.AssetCriticalityContributionListDTO

	// Build the query with joins to related tables
	query := repo.db.Table("asset_criticality_contributions").
		Select(`
			asset_criticality_contributions.id,
			asset_criticality_contributions.asset_register_id,
			asset_register.title_en AS asset_register_title,
			asset_criticality_contributions.asset_criticality_matrix_factor_id,
			asset_criticality_matrix_factors.title_en AS asset_criticality_matrix_factor_title_en,
			asset_criticality_matrix_factors.title_ar AS asset_criticality_matrix_factor_title_ar,
			asset_criticality_matrix_factors.description_en AS asset_criticality_matrix_factor_description_en,
			asset_criticality_matrix_factors.description_ar AS asset_criticality_matrix_factor_description_ar,
			asset_criticality_matrix_factors.weightage AS asset_criticality_matrix_factor_weightage,
			asset_criticality_matrix_factors.asset_criticality_matrix_id,
			asset_criticality_contributions.asset_criticality_matrix_factor_rating_id,
			asset_criticality_matrix_factor_ratings.title_en AS asset_criticality_matrix_factor_rating_title_en,
			asset_criticality_matrix_factor_ratings.title_ar AS asset_criticality_matrix_factor_rating_title_ar,
			asset_criticality_matrix_factor_ratings.score AS asset_criticality_matrix_factor_rating_score,
			asset_criticality_matrix_factor_ratings.type AS asset_criticality_matrix_factor_rating_type,
			asset_criticality_matrix_factor_ratings.color_code AS asset_criticality_matrix_factor_rating_color,
			asset_criticality_contributions.score,
			created_user.id AS created_by_id,
			created_user.first_name_en AS created_by_first_name_en,
			created_user.first_name_ar AS created_by_first_name_ar,
			created_user.last_name_en AS created_by_last_name_en,
			created_user.last_name_ar AS created_by_last_name_ar,
			created_user.initial AS created_by_initial,
			created_user.color AS created_by_color,
			created_user.email AS created_by_email,
			created_user.profile_image_path AS created_by_profile_image_path,
			asset_criticality_contributions.created_at,
			updated_user.id AS updated_by_id,
			updated_user.first_name_en AS updated_by_first_name_en,
			updated_user.first_name_ar AS updated_by_first_name_ar,
			updated_user.last_name_en AS updated_by_last_name_en,
			updated_user.last_name_ar AS updated_by_last_name_ar,
			updated_user.initial AS updated_by_initial,
			updated_user.color AS updated_by_color,
			updated_user.email AS updated_by_email,
			updated_user.profile_image_path AS updated_by_profile_image_path,
			asset_criticality_contributions.updated_at
		`).
		Joins("LEFT JOIN asset_register ON asset_register.id = asset_criticality_contributions.asset_register_id").
		Joins("LEFT JOIN asset_criticality_matrix_factors ON asset_criticality_matrix_factors.id = asset_criticality_contributions.asset_criticality_matrix_factor_id").
		Joins("LEFT JOIN asset_criticality_matrix_factor_ratings ON asset_criticality_matrix_factor_ratings.id = asset_criticality_contributions.asset_criticality_matrix_factor_rating_id").
		Joins("LEFT JOIN users AS created_user ON created_user.id = asset_criticality_contributions.created_by_id").
		Joins("LEFT JOIN users AS updated_user ON updated_user.id = asset_criticality_contributions.updated_by_id").
		Where("asset_criticality_contributions.asset_register_id = ?", assetID)

	// Apply keyword filtering if provided
	if params.Keyword != "" {
		query = query.Where(`
			(asset_criticality_matrix_factors.title_en LIKE ? OR
			asset_criticality_matrix_factor_ratings.title_en LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.last_name_en LIKE ?
			)`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	// Apply sorting based on parameters
	if params.SortBy == "" {
		params.SortBy = "asset_criticality_contributions.id" // Default sorting by ID
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default sorting order
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	// Execute the query
	if err := query.Find(&contributions).Error; err != nil {
		return nil, err
	}

	return contributions, nil
}

func (repo *AssetCriticalityRepository) GetByID(id uint) (*models.AssetCriticalityContribution, error) {
	var contribution models.AssetCriticalityContribution
	if err := repo.db.First(&contribution, id).Error; err != nil {
		return nil, err
	}
	return &contribution, nil
}

func (repo *AssetCriticalityRepository) Create(contribution *models.AssetCriticalityContribution, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// return db.Create(contribution).Error
	if err := db.Create(contribution).Error; err != nil {
		return 0, err
	}
	return contribution.ID, nil
}

func (repo *AssetCriticalityRepository) Update(contribution *models.AssetCriticalityContribution, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(contribution).Error
}

func (repo *AssetCriticalityRepository) Delete(assetID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete all contributions for the given asset register ID
	return db.Where("asset_register_id = ?", assetID).Delete(&models.AssetCriticalityContribution{}).Error
}

func (repo *AssetCriticalityRepository) CalculateTotalScore(assetID uint, tx ...utility.DatabaseTransaction) (float64, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var totalScore float64
	err := db.Model(&models.AssetCriticalityContribution{}).
		Where("asset_register_id = ?", assetID).
		Select("SUM(score)").
		Scan(&totalScore).Error
	if err != nil {
		return 0, err
	}
	return totalScore, nil
}

func (repo *AssetCriticalityRepository) GetAllContribution(assetID uint) ([]models.AssetCriticalityContribution, error) {
	var criticalityContribution []models.AssetCriticalityContribution

	// Fetch vulnerabilities with their related threats
	err := repo.db.
		Preload("AssetCriticalityMatrixFactor").
		Preload("AssetCriticalityMatrixFactorRating").
		Where("asset_register_id = ?", assetID).
		Find(&criticalityContribution).Error

	return criticalityContribution, err
}
