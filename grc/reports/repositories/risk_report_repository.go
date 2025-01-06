package repositories

import (
	reportDTOs "jethur/grc/reports/dtos"

	"gorm.io/gorm"
)

type RiskReportRepositoryInterface interface {
	GetKRIByCategory(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.KRIByCategoryDTO, error)
	GetKRIByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.KRIByResponsibleUnitDTO, error)
	GetKRIByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.KRIByResponsibleUserDTO, error)
	GetKRIByRiskLevel(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.KRIByRiskLevelDTO, error)
	GetKRIByReviewFrequency(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.KRIByReviewFrequencyDTO, error)

	GetRiskByCategory(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.RiskByCategoryDTO, error)
	GetRiskByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.RiskByResponsibleUnitDTO, error)
	GetRiskByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.RiskByResponsibleUserDTO, error)
	GetRiskByRiskSource(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.RiskByRiskSourceDTO, error)
	GetRiskByInherentRiskRating(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.RiskByRiskRatingDTO, error)
	GetRiskByResidualRiskRating(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.RiskByRiskRatingDTO, error)
	GetRiskByRiskStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.RiskByStatusDTO, error)
	GetRiskByRiskResponseStrategy(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.RiskByResponseStrategyDTO, error)
	GetRiskByRiskAppetiteToleranceLevel(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.RiskByRiskAppetiteToleranceLevelDTO, error)
	GetRiskByKRIEnabledStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.RiskByKRIEnabledStatusDTO, error)

	GetRiskActionPlansByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.RiskActionPlansByStatusDTO, error)
	GetRiskActionPlansByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.RiskActionPlansByResponsibleUserDTO, error)
	GetRiskActionPlansBySeverityLevel(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.RiskActionPlansBySeverityLevelDTO, error)
}

type RiskReportRepository struct {
	db *gorm.DB
}

func NewRiskReportRepository(db *gorm.DB) *RiskReportRepository {
	return &RiskReportRepository{db: db}
}

// Implementation of the `GetKRIByCategory` method
func (r *RiskReportRepository) GetKRIByCategory(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.KRIByCategoryDTO, error) {
	results := []*reportDTOs.KRIByCategoryDTO{}

	query := r.db.Table("kri_register").
		Select(`kri_categories.id AS kri_category_id, 
				kri_categories.title_en AS kri_category_title_en, 
				kri_categories.title_ar AS kri_category_title_ar, 
				COUNT(kri_register.id) AS count`).
		Joins("JOIN kri_categories ON kri_categories.id = kri_register.kri_category_id").
		Group("kri_categories.id, kri_categories.title_en, kri_categories.title_ar")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("kri_register.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (r *RiskReportRepository) GetKRIByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.KRIByResponsibleUnitDTO, error) {
	results := []*reportDTOs.KRIByResponsibleUnitDTO{}

	query := r.db.Table("kri_register").
		Select(`
			organization_units.id AS responsible_unit_id,
			organization_units.title_en AS responsible_unit_title_en,
			organization_units.title_ar AS responsible_unit_title_ar,
			organization_units.color_code AS responsible_unit_color_code,
			COUNT(kri_register.id) AS count
		`).
		Joins("JOIN organization_units ON organization_units.id = kri_register.responsible_unit_id").
		Group("organization_units.id, organization_units.title_en, organization_units.title_ar, organization_units.color_code")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("kri_register.created_at BETWEEN ? AND ?", params.DateFrom, params.DateTo)
	}

	err := query.Scan(&results).Error
	return results, err
}

func (r *RiskReportRepository) GetKRIByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.KRIByResponsibleUserDTO, error) {
	results := []*reportDTOs.KRIByResponsibleUserDTO{}

	query := r.db.Table("kri_register").
		Select(`
			users.id AS responsible_user_id,
			users.first_name_en AS responsible_user_first_name_en,
			users.first_name_ar AS responsible_user_first_name_ar,
			users.last_name_en AS responsible_user_last_name_en,
			users.last_name_ar AS responsible_user_last_name_ar,
			CONCAT(users.first_name_en, ' ', users.last_name_en) AS responsible_user_title_en,
            CONCAT(users.first_name_ar, ' ', users.last_name_ar) AS responsible_user_title_ar,
			users.initial AS responsible_user_initial,
			users.color AS responsible_user_color,
			users.email AS responsible_user_email,
			users.profile_image_path AS responsible_user_profile_image_path,
			COUNT(kri_register.id) AS count
		`).
		Joins("JOIN users ON users.id = kri_register.responsible_user_id").
		Group("users.id, users.first_name_en, users.first_name_ar, users.last_name_en, users.last_name_ar, users.initial, users.color, users.email, users.profile_image_path")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("kri_register.created_at BETWEEN ? AND ?", params.DateFrom, params.DateTo)
	}

	err := query.Scan(&results).Error
	return results, err
}

func (r *RiskReportRepository) GetKRIByRiskLevel(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.KRIByRiskLevelDTO, error) {
	results := []*reportDTOs.KRIByRiskLevelDTO{}

	query := r.db.Table("kri_register").
		Select(`risk_levels.id AS risk_level_id,
				risk_levels.title_en AS risk_level_title_en,
				risk_levels.title_ar AS risk_level_title_ar,
				risk_levels.color_code AS color_code,
				COUNT(kri_register.id) AS count`).
		Joins("JOIN risk_levels ON risk_levels.id = kri_register.risk_level_id").
		Group("risk_levels.id, risk_levels.title_en, risk_levels.title_ar, risk_levels.color_code")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("kri_register.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (r *RiskReportRepository) GetKRIByReviewFrequency(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.KRIByReviewFrequencyDTO, error) {
	results := []*reportDTOs.KRIByReviewFrequencyDTO{}

	query := r.db.Table("kri_register").
		Select(`review_frequencies.id AS review_frequency_id,
				review_frequencies.title_en AS review_frequency_title_en,
				review_frequencies.title_ar AS review_frequency_title_ar,
				COUNT(kri_register.id) AS count`).
		Joins("JOIN review_frequencies ON review_frequencies.id = kri_register.review_frequency_id").
		Group("review_frequencies.id, review_frequencies.title_en, review_frequencies.title_ar")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("kri_register.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (r *RiskReportRepository) GetRiskByCategory(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.RiskByCategoryDTO, error) {
	results := []*reportDTOs.RiskByCategoryDTO{}

	query := r.db.Table("risk_register").
		Select(`risk_categories.id AS risk_category_id,
				risk_categories.title_en AS risk_category_title_en,
				risk_categories.title_ar AS risk_category_title_ar,
				COUNT(risk_register.id) AS count`).
		Joins("JOIN risk_categories ON risk_categories.id = risk_register.risk_category_id").
		Group("risk_categories.id, risk_categories.title_en, risk_categories.title_ar")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("risk_register.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (r *RiskReportRepository) GetRiskByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.RiskByResponsibleUnitDTO, error) {
	results := []*reportDTOs.RiskByResponsibleUnitDTO{}

	query := r.db.Table("risk_register").
		Select(`organization_units.id AS responsible_unit_id,
				organization_units.title_en AS responsible_unit_title_en,
				organization_units.title_ar AS responsible_unit_title_ar,
				organization_units.color_code AS responsible_unit_color_code,
				COUNT(risk_register.id) AS count`).
		Joins("JOIN organization_units ON organization_units.id = risk_register.responsible_unit_id").
		Group(`organization_units.id, 
				organization_units.title_en, 
				organization_units.title_ar, 
				organization_units.color_code`)

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("risk_register.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (r *RiskReportRepository) GetRiskByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.RiskByResponsibleUserDTO, error) {
	results := []*reportDTOs.RiskByResponsibleUserDTO{}

	query := r.db.Table("risk_register").
		Select(`users.id AS responsible_user_id,
			users.first_name_en AS responsible_user_first_name_en,
			users.first_name_ar AS responsible_user_first_name_ar,
			users.last_name_en AS responsible_user_last_name_en,
			users.last_name_ar AS responsible_user_last_name_ar,
			CONCAT(users.first_name_en, ' ', users.last_name_en) AS responsible_user_title_en,
            CONCAT(users.first_name_ar, ' ', users.last_name_ar) AS responsible_user_title_ar,
			users.initial AS responsible_user_initial,
			users.color AS responsible_user_color,
			users.email AS responsible_user_email,
			users.profile_image_path AS responsible_user_profile_image_path,
			COUNT(risk_register.id) AS count`).
		Joins("JOIN users ON users.id = risk_register.responsible_user_id").
		Group(`users.id, 
			users.first_name_en, 
			users.first_name_ar, 
			users.last_name_en, 
			users.last_name_ar, 
			users.initial, 
			users.color, 
			users.email, 
			users.profile_image_path`)

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("risk_register.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (r *RiskReportRepository) GetRiskByRiskSource(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.RiskByRiskSourceDTO, error) {
	results := []*reportDTOs.RiskByRiskSourceDTO{}

	query := r.db.Table("risk_register").
		Select(`risk_sources.id AS risk_source_id,
				risk_sources.title_en AS risk_source_title_en,
				risk_sources.title_ar AS risk_source_title_ar,
				COUNT(risk_register.id) AS count`).
		Joins("JOIN risk_sources ON risk_sources.id = risk_register.risk_source_id").
		Group("risk_sources.id, risk_sources.title_en, risk_sources.title_ar")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("risk_register.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (r *RiskReportRepository) GetRiskByInherentRiskRating(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.RiskByRiskRatingDTO, error) {
	results := []*reportDTOs.RiskByRiskRatingDTO{}

	query := r.db.Table("risk_register").
		Select(`risk_matrix_risk_ratings.id AS risk_rating_id,
				risk_matrix_risk_ratings.title_en AS risk_rating_title_en,
				risk_matrix_risk_ratings.title_ar AS risk_rating_title_ar,
				risk_matrix_risk_ratings.color_code AS color_code,
				COUNT(risk_register.id) AS count`).
		Joins("JOIN risk_matrix_risk_ratings ON risk_matrix_risk_ratings.id = risk_register.inherent_risk_rating_id").
		Group("risk_matrix_risk_ratings.id, risk_matrix_risk_ratings.title_en, risk_matrix_risk_ratings.title_ar, risk_matrix_risk_ratings.color_code")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("risk_register.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (r *RiskReportRepository) GetRiskByResidualRiskRating(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.RiskByRiskRatingDTO, error) {
	results := []*reportDTOs.RiskByRiskRatingDTO{}

	query := r.db.Table("risk_register").
		Select(`risk_matrix_risk_ratings.id AS risk_rating_id,
				risk_matrix_risk_ratings.title_en AS risk_rating_title_en,
				risk_matrix_risk_ratings.title_ar AS risk_rating_title_ar,
				risk_matrix_risk_ratings.color_code AS color_code,
				COUNT(risk_register.id) AS count`).
		Joins("JOIN risk_matrix_risk_ratings ON risk_matrix_risk_ratings.id = risk_register.residual_risk_rating_id").
		Group("risk_matrix_risk_ratings.id, risk_matrix_risk_ratings.title_en, risk_matrix_risk_ratings.title_ar, risk_matrix_risk_ratings.color_code")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("risk_register.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (r *RiskReportRepository) GetRiskByRiskStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.RiskByStatusDTO, error) {
	results := []*reportDTOs.RiskByStatusDTO{}

	query := r.db.Table("risk_register").
		Select(`risk_statuses.id AS risk_status_id,
				risk_statuses.title_en AS risk_status_title_en,
				risk_statuses.title_ar AS risk_status_title_ar,
				risk_statuses.color_code AS color_code,
				COUNT(risk_register.id) AS count`).
		Joins("JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id").
		Group("risk_statuses.id, risk_statuses.title_en, risk_statuses.title_ar, risk_statuses.color_code")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("risk_register.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (r *RiskReportRepository) GetRiskByRiskResponseStrategy(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.RiskByResponseStrategyDTO, error) {
	results := []*reportDTOs.RiskByResponseStrategyDTO{}

	query := r.db.Table("risk_register").
		Select(`risk_response_strategies.id AS response_strategy_id,
				risk_response_strategies.title_en AS response_strategy_title_en,
				risk_response_strategies.title_ar AS response_strategy_title_ar,
				COUNT(risk_register.id) AS count`).
		Joins("JOIN risk_response_strategies ON risk_response_strategies.id = risk_register.risk_response_strategy_id").
		Group("risk_response_strategies.id, risk_response_strategies.title_en, risk_response_strategies.title_ar")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("risk_register.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (r *RiskReportRepository) GetRiskByRiskAppetiteToleranceLevel(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.RiskByRiskAppetiteToleranceLevelDTO, error) {
	results := []*reportDTOs.RiskByRiskAppetiteToleranceLevelDTO{}

	query := r.db.Table("risk_appetite_alignments").
		Select(`risk_appetite_tolerance_levels.id AS tolerance_level_id,
				risk_appetite_tolerance_levels.title_en AS tolerance_level_title_en,
				risk_appetite_tolerance_levels.title_ar AS tolerance_level_title_ar,
				risk_appetite_tolerance_levels.color_code AS color_code,
				COUNT(risk_appetite_alignments.id) AS count`).
		Joins("JOIN risk_appetite_tolerance_levels ON risk_appetite_tolerance_levels.id = risk_appetite_alignments.risk_appetite_tolerance_level_id").
		Group(`risk_appetite_tolerance_levels.id, 
			   risk_appetite_tolerance_levels.title_en, 
			   risk_appetite_tolerance_levels.title_ar, 
			   risk_appetite_tolerance_levels.color_code`)

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("risk_appetite_alignments.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (r *RiskReportRepository) GetRiskByKRIEnabledStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.RiskByKRIEnabledStatusDTO, error) {
	results := []*reportDTOs.RiskByKRIEnabledStatusDTO{}

	// Build the subquery with date filtering applied
	subQuery := r.db.Table("risk_register").
		Select(`risk_register.id AS risk_id,
				CASE 
					WHEN EXISTS (SELECT 1 FROM risk_kris WHERE risk_kris.risk_register_id = risk_register.id) 
					THEN 'Enabled' 
					ELSE 'Not Enabled' 
				END AS kri_status`)

	// Apply date filtering in the subquery
	if params.DateFrom != "" && params.DateTo != "" {
		subQuery = subQuery.Where("risk_register.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	// Use the subquery to calculate counts grouped by KRI status
	query := r.db.Table("(?) AS sub", subQuery).
		Select(`kri_status, COUNT(*) AS count`).
		Group("kri_status")

	// Execute the query
	err := query.Scan(&results).Error
	return results, err
}

func (r *RiskReportRepository) GetRiskActionPlansByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.RiskActionPlansByStatusDTO, error) {
	results := []*reportDTOs.RiskActionPlansByStatusDTO{}

	query := r.db.Table("risk_action_plans").
		Select(`action_plan_statuses.id AS action_plan_status_id,
				action_plan_statuses.title_en AS action_plan_status_title_en,
				action_plan_statuses.title_ar AS action_plan_status_title_ar,
				action_plan_statuses.type AS action_plan_status_type,
				action_plan_statuses.color_code AS action_plan_status_color_code,
				COUNT(risk_action_plans.id) AS count`).
		Joins("JOIN action_plan_statuses ON action_plan_statuses.id = risk_action_plans.action_plan_status_id").
		Group("action_plan_statuses.id, action_plan_statuses.title_en, action_plan_statuses.title_ar, action_plan_statuses.type, action_plan_statuses.color_code")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("risk_action_plans.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (r *RiskReportRepository) GetRiskActionPlansByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.RiskActionPlansByResponsibleUserDTO, error) {
	results := []*reportDTOs.RiskActionPlansByResponsibleUserDTO{}

	query := r.db.Table("risk_action_plans").
		Select(`users.id AS responsible_user_id,
				users.first_name_en AS responsible_user_first_name_en,
				users.first_name_ar AS responsible_user_first_name_ar,
				users.last_name_en AS responsible_user_last_name_en,
				users.last_name_ar AS responsible_user_last_name_ar,
				CONCAT(users.first_name_en, ' ', users.last_name_en) AS responsible_user_title_en,
            	CONCAT(users.first_name_ar, ' ', users.last_name_ar) AS responsible_user_title_ar,
				users.initial AS responsible_user_initial,
				users.color AS responsible_user_color,
				users.email AS responsible_user_email,
				users.profile_image_path AS responsible_user_profile_image_path,
				COUNT(risk_action_plans.id) AS count`).
		Joins("JOIN users ON users.id = risk_action_plans.responsible_user_id").
		Group("users.id, users.first_name_en, users.first_name_ar, users.last_name_en, users.last_name_ar, users.initial, users.color, users.email, users.profile_image_path")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("risk_action_plans.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (r *RiskReportRepository) GetRiskActionPlansBySeverityLevel(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.RiskActionPlansBySeverityLevelDTO, error) {
	results := []*reportDTOs.RiskActionPlansBySeverityLevelDTO{}

	query := r.db.Table("risk_action_plans").
		Select(`severity_levels.id AS severity_level_id,
				severity_levels.title_en AS severity_level_title_en,
				severity_levels.title_ar AS severity_level_title_ar,
				severity_levels.color_code AS severity_level_color_code,
				severity_levels.type AS severity_level_type,
				COUNT(risk_action_plans.id) AS count`).
		Joins("JOIN severity_levels ON severity_levels.id = risk_action_plans.severity_level_id").
		Group("severity_levels.id, severity_levels.title_en, severity_levels.title_ar, severity_levels.color_code, severity_levels.type")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("risk_action_plans.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}
