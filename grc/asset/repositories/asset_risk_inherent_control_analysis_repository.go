package repositories

import (
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type AssetRiskInherentControlAnalysisRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(riskID uint, params *riskDTOs.RiskControlAnalysisListQueryParamsDTO) ([]*riskDTOs.RiskControlAnalysisListDTO, error)
	GetFullDetailsByID(id uint) (*models.RiskControlAnalysis, error)
	GetByID(id uint) (*models.RiskControlAnalysis, error)
	Create(riskControlAnalysis *models.RiskControlAnalysis, tx ...utility.DatabaseTransaction) (uint, error)
	Update(riskControlAnalysis *models.RiskControlAnalysis, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	GetAllControlsByRiskID(riskID uint, tx ...utility.DatabaseTransaction) ([]*models.RiskControl, error)
	GetControlSummary(riskID uint) ([]models.ControlTest, error)
}

type AssetRiskInherentControlAnalysisRepository struct {
	db *gorm.DB
}

func NewAssetRiskInherentControlAnalysisRepository(db *gorm.DB) *AssetRiskInherentControlAnalysisRepository {
	return &AssetRiskInherentControlAnalysisRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *AssetRiskInherentControlAnalysisRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AssetRiskInherentControlAnalysisRepository) GetAll(riskID uint, params *riskDTOs.RiskControlAnalysisListQueryParamsDTO) ([]*riskDTOs.RiskControlAnalysisListDTO, error) {
	var riskControlAnalyses []*riskDTOs.RiskControlAnalysisListDTO

	query := repo.db.Table("risk_control_analyses").
		Select(`
            risk_control_analyses.id,
			risk_control_analyses.type,
			risk_controls.control_id,
			latest_control_test.id AS control_test_id,
			controls.reference_number AS control_reference_number,
			controls.code AS control_code,
			controls.title_en AS control_title_en,
			controls.title_ar AS control_title_ar,
			CAST(risk_control_analyses.score AS VARCHAR) AS score,
			controls.description AS control_description,
			control_types.title_en AS control_type_title_en,
			control_types.title_ar AS control_type_title_ar,
			control_types.type AS control_type,
			control_categories.title_en AS control_category_title_en,
			control_categories.title_ar AS control_category_title_ar,
			organization_units.title_en AS owner_unit_title_en,
			organization_units.title_ar AS owner_unit_title_ar,
			controls.responsible_user_id AS owner_id,
			owner.first_name_en AS owner_first_name_en,
			owner.first_name_ar AS owner_first_name_ar,
			owner.last_name_en AS owner_last_name_en,
			owner.last_name_ar AS owner_last_name_ar,
			owner.initial AS owner_initial,
			owner.color AS owner_color,
			owner.email AS owner_email,
			owner.profile_image_path AS owner_profile_image_path,
			risk_control_analyses.type,
			control_implementation_statuses.title_en AS control_implementation_status_title_en,
			control_implementation_statuses.title_ar AS control_implementation_status_title_ar,
			control_implementation_statuses.color_code AS control_implementation_status_color_code,
			control_efficiency_statuses.title_en AS control_efficiency_status_title_en,
			control_efficiency_statuses.title_ar AS control_efficiency_status_title_ar,
			control_efficiency_statuses.color_code AS control_efficiency_status_color_code,
			control_effectiveness_ratings.title_en AS control_effectiveness_rating_title_en,
			control_effectiveness_ratings.title_ar AS control_effectiveness_rating_title_ar,
			control_effectiveness_ratings.min_score AS control_effectiveness_rating_min_score,
			control_effectiveness_ratings.max_score AS control_effectiveness_rating_max_score,
			control_effectiveness_ratings.type AS control_effectiveness_rating_type,
			control_effectiveness_ratings.color_code AS control_effectiveness_rating_color_code,
			control_automation.id AS control_automation_id,
			control_automation.title_en AS control_automation_title_en,
			control_automation.title_ar AS control_automation_title_ar,
			control_automation.type AS control_automation_type,
			control_consistency.id AS control_consistency_id,
			control_consistency.title_en AS control_consistency_title_en,
			control_consistency.title_ar AS control_consistency_title_ar,
			control_consistency.type AS control_consistency_type,
			CONVERT(VARCHAR, latest_control_test.tested_at, 23) AS last_test_date,
			latest_control_test.title AS control_test_title,
			latest_control_test.description AS control_test_description,
			latest_control_test.comments AS control_test_comments,
			latest_control_test.score AS control_test_score,
			latest_control_test.control_effectiveness_score AS control_test_effectiveness_score,
			latest_control_test.risk_reduction AS control_test_risk_reduction,
			CONVERT(VARCHAR, latest_control_test.tested_at, 23) AS control_test_tested_at,
			tested_by.first_name_en AS control_test_tested_by_first_name_en,
			tested_by.last_name_en AS control_test_tested_by_last_name_en,
			tested_by.initial AS control_test_tested_by_initial,
			tested_by.color AS control_test_tested_by_color,
			control_statuses.id AS control_status_id,
			control_statuses.title_en AS control_status_title_en,
			control_statuses.title_ar AS control_status_title_ar,
			control_statuses.type AS control_status_type,
			control_statuses.color_code AS control_status_color_code,
            created_user.id AS created_user_id,
        	created_user.first_name_en AS created_user_first_name_en,
            created_user.first_name_ar AS created_user_first_name_ar,
            created_user.last_name_en AS created_user_last_name_en,
            created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
            created_user.color AS created_user_color,
            created_user.email AS created_user_email,
            created_user.profile_image_path AS created_user_profile_image_path,
			CONVERT(VARCHAR, risk_control_analyses.created_at, 23) AS created_at
        `).
		Joins(`
            LEFT JOIN risk_controls ON risk_controls.id = risk_control_analyses.risk_control_id
			LEFT JOIN controls ON controls.id = risk_controls.control_id
			LEFT JOIN control_types ON control_types.id = controls.control_type_id
			LEFT JOIN control_categories ON control_categories.id = controls.control_category_id
			LEFT JOIN organization_units ON organization_units.id = controls.responsible_unit_id
			LEFT JOIN users AS owner ON owner.id = controls.responsible_user_id
            LEFT JOIN control_efficiency_statuses ON control_efficiency_statuses.id = risk_control_analyses.control_efficiency_status_id
            LEFT JOIN control_implementation_statuses ON control_implementation_statuses.id = risk_control_analyses.control_implementation_status_id
            LEFT JOIN control_effectiveness_ratings ON control_effectiveness_ratings.id = risk_control_analyses.control_effectiveness_rating_id
            LEFT JOIN (
                SELECT ct.control_id, MAX(ct.id) as latest_test_id
                FROM control_tests ct
                GROUP BY ct.control_id
            ) AS latest_control_test_ids ON latest_control_test_ids.control_id = controls.id
			LEFT JOIN control_statuses ON control_statuses.id = controls.control_status_id
            LEFT JOIN control_tests AS latest_control_test ON latest_control_test.id = latest_control_test_ids.latest_test_id
            LEFT JOIN control_automations AS control_automation ON control_automation.id = latest_control_test.control_automation_id
            LEFT JOIN control_consistencies AS control_consistency ON control_consistency.id = latest_control_test.control_consistency_id
            LEFT JOIN users AS tested_by ON tested_by.id = latest_control_test.tested_by
            LEFT JOIN users AS created_user ON created_user.id = risk_control_analyses.created_by_id
        `).
		Where(`risk_controls.risk_register_id = ? AND risk_control_analyses.type = 'inherent'`, riskID)

	// Apply search keyword if exists
	if params.Keyword != "" {
		query = query.Where(`
            (risk_control_analyses.type LIKE ? OR
            risk_control_analyses.control_efficiency_status_id LIKE ? OR
            risk_control_analyses.control_implementation_status_id LIKE ?
            )`,
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
	if err := query.Find(&riskControlAnalyses).Error; err != nil {
		return nil, err
	}

	// Load control test documents for each analysis
	for _, analysis := range riskControlAnalyses {
		if analysis.ControlTestID != nil {
			var documents []organizationDTOs.ControlTestDocumentDTO

			err := repo.db.
				Table("control_test_documents").
				Select("id, control_test_id, url, name, original_name, token, type, size").
				Where("control_test_id = ?", *analysis.ControlTestID).
				Scan(&documents).Error

			if err != nil {
				return nil, err
			}

			if documents == nil {
				documents = []organizationDTOs.ControlTestDocumentDTO{}
			}

			analysis.ControlTestDocuments = documents
		} else {
			analysis.ControlTestDocuments = []organizationDTOs.ControlTestDocumentDTO{}
		}
	}

	return riskControlAnalyses, nil
}

func (repo *AssetRiskInherentControlAnalysisRepository) GetByID(id uint) (*models.RiskControlAnalysis, error) {
	var riskControlAnalysis *models.RiskControlAnalysis

	result := repo.db.
		Preload("RiskControl").
		First(&riskControlAnalysis, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return riskControlAnalysis, nil
}

func (repo *AssetRiskInherentControlAnalysisRepository) GetFullDetailsByID(id uint) (*models.RiskControlAnalysis, error) {
	var riskControlAnalysis *models.RiskControlAnalysis

	result := repo.db.
		Preload("RiskControl").
		Preload("ControlEfficiencyStatus").
		Preload("ControlImplementationStatus").
		Preload("ControlEffectivenessRating").
		First(&riskControlAnalysis, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return riskControlAnalysis, nil
}

func (repo *AssetRiskInherentControlAnalysisRepository) Create(riskControlAnalysis *models.RiskControlAnalysis, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(riskControlAnalysis).Error; err != nil {
		return 0, err
	}

	return riskControlAnalysis.ID, nil
}

func (repo *AssetRiskInherentControlAnalysisRepository) Update(riskControlAnalysis *models.RiskControlAnalysis, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(riskControlAnalysis).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AssetRiskInherentControlAnalysisRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Delete(&models.RiskControlAnalysis{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AssetRiskInherentControlAnalysisRepository) GetAllControlsByRiskID(riskID uint, tx ...utility.DatabaseTransaction) ([]*models.RiskControl, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var riskControls []*models.RiskControl

	err := db.Where("risk_register_id = ?", riskID).Find(&riskControls).Error

	if err != nil {
		return nil, err
	}

	return riskControls, nil
}

func (repo *AssetRiskInherentControlAnalysisRepository) GetControlSummary(riskID uint) ([]models.ControlTest, error) {
	var controlTests []models.ControlTest

	result := repo.db.
		Preload("ControlEfficiencyStatus").
		Preload("ControlImplementationStatus").
		Preload("ControlEffectivenessRating").
		Where("risk_register_id = ?", riskID).
		Where("tested_source = ?", "asset_risk").
		Find(&controlTests)
	if result.Error != nil {
		return nil, result.Error
	}

	return controlTests, nil
}
