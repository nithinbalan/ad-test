package repositories

import (
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskInherentControlAnalysisRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(riskID uint, params *riskDTOs.RiskControlAnalysisListQueryParamsDTO) ([]*riskDTOs.RiskControlAnalysisListDTO, error)
	GetFullDetailsByID(id uint) (*models.RiskControlAnalysis, error)
	GetByID(id uint) (*models.RiskControlAnalysis, error)
	Create(riskControlAnalysis *models.RiskControlAnalysis) (uint, error)
	Update(riskControlAnalysis *models.RiskControlAnalysis, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	GetAllControlsByRiskID(riskID uint) ([]*models.RiskControl, error)
	GetControlSummary(riskID uint) ([]models.RiskControlAnalysis, error)
}

type RiskInherentControlAnalysisRepository struct {
	db *gorm.DB
}

func NewRiskInherentControlAnalysisRepository(db *gorm.DB) *RiskInherentControlAnalysisRepository {
	return &RiskInherentControlAnalysisRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *RiskInherentControlAnalysisRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *RiskInherentControlAnalysisRepository) GetAll(riskID uint, params *riskDTOs.RiskControlAnalysisListQueryParamsDTO) ([]*riskDTOs.RiskControlAnalysisListDTO, error) {
	var riskControlAnalyses []*riskDTOs.RiskControlAnalysisListDTO

	query := repo.db.Table("risk_control_analyses").
		Select(`
            risk_control_analyses.id,
			risk_control_analyses.type,
			risk_control_analyses.control_test_id,
			risk_control_analyses.score,
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
            created_user.color AS created_userr_color,
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

	return riskControlAnalyses, nil
}

func (repo *RiskInherentControlAnalysisRepository) GetByID(id uint) (*models.RiskControlAnalysis, error) {
	var riskControlAnalysis *models.RiskControlAnalysis

	result := repo.db.
		Preload("RiskControl").
		First(&riskControlAnalysis, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return riskControlAnalysis, nil
}

func (repo *RiskInherentControlAnalysisRepository) GetFullDetailsByID(id uint) (*models.RiskControlAnalysis, error) {
	var riskControlAnalysis *models.RiskControlAnalysis

	result := repo.db.
		Preload("RiskControl").
		Preload("ControlEfficiencyStatus").
		Preload("ControlImplementationStatus").
		First(&riskControlAnalysis, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return riskControlAnalysis, nil
}

func (repo *RiskInherentControlAnalysisRepository) Create(riskControlAnalysis *models.RiskControlAnalysis) (uint, error) {
	if err := repo.db.Create(riskControlAnalysis).Error; err != nil {
		return 0, err
	}

	return riskControlAnalysis.ID, nil
}

func (repo *RiskInherentControlAnalysisRepository) Update(riskControlAnalysis *models.RiskControlAnalysis, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(riskControlAnalysis).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskInherentControlAnalysisRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Delete(&models.RiskControlAnalysis{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskInherentControlAnalysisRepository) GetAllControlsByRiskID(riskID uint) ([]*models.RiskControl, error) {
	var riskControls []*models.RiskControl

	err := repo.db.Where("risk_register_id = ?", riskID).Find(&riskControls).Error

	if err != nil {
		return nil, err
	}

	return riskControls, nil
}

func (repo *RiskInherentControlAnalysisRepository) GetControlSummary(riskID uint) ([]models.RiskControlAnalysis, error) {
	var riskControlAnalyses []models.RiskControlAnalysis

	result := repo.db.
		Preload("ControlEfficiencyStatus").
		Preload("ControlImplementationStatus").
		Joins("JOIN risk_controls ON risk_controls.id = risk_control_analyses.risk_control_id").
		Where("risk_controls.risk_register_id = ?", riskID).
		Where("risk_control_analyses.type = ?", "inherent").
		Find(&riskControlAnalyses)
	if result.Error != nil {
		return nil, result.Error
	}

	return riskControlAnalyses, nil
}
