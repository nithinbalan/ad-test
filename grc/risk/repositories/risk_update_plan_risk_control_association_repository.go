package repositories

import (
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskUpdatePlanRiskControlAssociationRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetRiskControlAssociationList(planID uint, isCritical bool) ([]*riskDTOs.RiskControlAssociationListDTO, error)
	CreateBulkRiskControls(mappings []*models.RiskControl, tx ...utility.DatabaseTransaction) error
	DeleteRiskControlsByRiskRegisterID(riskRegisterID uint, tx ...utility.DatabaseTransaction) error
	GetRiskControlsByRiskRegisterID(riskRegisterID uint) ([]*models.RiskControl, error)
	// Control Analysis for Risk Register Update Plan
	CreateBulkControlAnalyses(analyses []*models.RiskControlAnalysis, tx ...utility.DatabaseTransaction) error
	DeleteControlAnalysesByRiskRegisterIDs(riskRegisterIDs []uint, tx ...utility.DatabaseTransaction) error
	GetControlAnalysesByRiskRegisterID(riskRegisterID uint) ([]*models.RiskControlAnalysis, error)
	// Residual Evaluation for Risk Register Update Plan
	GetResidualEvaluationByRiskRegisterID(riskRegisterID uint, tx ...utility.DatabaseTransaction) (*models.RiskEvaluation, error)
	GetResidualEvaluationByRiskRegisterIDAndPlanID(riskRegisterID uint, planID uint, tx ...utility.DatabaseTransaction) (*models.RiskEvaluation, error)
	DeleteResidualEvaluationsByRiskRegisterIDs(riskRegisterIDs []uint, tx ...utility.DatabaseTransaction) error
	CreateBulkResidualEvaluations(evaluations []*models.RiskEvaluation, tx ...utility.DatabaseTransaction) error
	CreateResidualEvaluation(evaluation *models.RiskEvaluation, tx ...utility.DatabaseTransaction) error
	UpdateResidualEvaluation(evaluation *models.RiskEvaluation, tx ...utility.DatabaseTransaction) error
	CreateResidualEvaluationHistory(history *models.RiskResidualEvaluationHistory, tx ...utility.DatabaseTransaction) error
}

type RiskUpdatePlanRiskControlAssociationRepository struct {
	db *gorm.DB
}

func NewRiskUpdatePlanRiskControlAssociationRepository(db *gorm.DB) *RiskUpdatePlanRiskControlAssociationRepository {
	return &RiskUpdatePlanRiskControlAssociationRepository{
		db: db,
	}
}

func (repo *RiskUpdatePlanRiskControlAssociationRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *RiskUpdatePlanRiskControlAssociationRepository) CreateBulkRiskControls(mappings []*models.RiskControl, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if len(mappings) == 0 {
		return nil
	}

	return db.Create(&mappings).Error
}

func (repo *RiskUpdatePlanRiskControlAssociationRepository) DeleteRiskControlsByRiskRegisterID(riskRegisterID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("risk_register_id = ?", riskRegisterID).Delete(&models.RiskControl{}).Error
}

func (repo *RiskUpdatePlanRiskControlAssociationRepository) GetRiskControlsByRiskRegisterID(riskRegisterID uint) ([]*models.RiskControl, error) {
	var riskControls []*models.RiskControl
	err := repo.db.Where("risk_register_id = ?", riskRegisterID).Find(&riskControls).Error
	return riskControls, err
}

// Control Analysis Repository Methods for Risk Register Update Plan

func (repo *RiskUpdatePlanRiskControlAssociationRepository) CreateBulkControlAnalyses(analyses []*models.RiskControlAnalysis, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if len(analyses) == 0 {
		return nil
	}

	return db.Create(&analyses).Error
}

func (repo *RiskUpdatePlanRiskControlAssociationRepository) DeleteControlAnalysesByRiskRegisterIDs(riskRegisterIDs []uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Get risk control IDs for the given risk register IDs
	var riskControlIDs []uint
	err := db.Table("risk_controls").
		Where("risk_register_id IN ?", riskRegisterIDs).
		Pluck("id", &riskControlIDs).Error
	if err != nil {
		return err
	}

	if len(riskControlIDs) == 0 {
		return nil
	}

	return db.Where("risk_control_id IN ? AND type = ?", riskControlIDs, "inherent").Delete(&models.RiskControlAnalysis{}).Error
}

func (repo *RiskUpdatePlanRiskControlAssociationRepository) GetControlAnalysesByRiskRegisterID(riskRegisterID uint) ([]*models.RiskControlAnalysis, error) {
	var controlAnalyses []*models.RiskControlAnalysis

	// Get risk control IDs for the given risk register ID
	var riskControlIDs []uint
	err := repo.db.Table("risk_controls").
		Where("risk_register_id = ?", riskRegisterID).
		Pluck("id", &riskControlIDs).Error
	if err != nil {
		return nil, err
	}

	if len(riskControlIDs) == 0 {
		return controlAnalyses, nil
	}

	err = repo.db.Where("risk_control_id IN ? AND type = ?", riskControlIDs, "inherent").Find(&controlAnalyses).Error
	return controlAnalyses, err
}

// Residual Evaluation Repository Methods for Risk Register Update Plan

func (repo *RiskUpdatePlanRiskControlAssociationRepository) DeleteResidualEvaluationsByRiskRegisterIDs(riskRegisterIDs []uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Get the evaluation IDs for the given risk register IDs
	var evaluationIDs []uint
	if err := db.Model(&models.RiskEvaluation{}).Where("risk_register_id IN ? AND type = ?", riskRegisterIDs, "residual").Pluck("id", &evaluationIDs).Error; err != nil {
		return err
	}

	if len(evaluationIDs) == 0 {
		return nil
	}

	// Delete impact analysis guidelines first
	if err := db.Where("risk_impact_analysis_id IN (SELECT id FROM risk_impact_analyses WHERE risk_evaluation_id IN ?)", evaluationIDs).Delete(&models.RiskImpactAnalysisGuideline{}).Error; err != nil {
		return err
	}

	// Delete impact analyses
	if err := db.Where("risk_evaluation_id IN ?", evaluationIDs).Delete(&models.RiskImpactAnalysis{}).Error; err != nil {
		return err
	}

	// Delete risk appetite alignments if they exist
	if err := db.Where("risk_evaluation_id IN ?", evaluationIDs).Delete(&models.RiskAppetiteAlignment{}).Error; err != nil {
		return err
	}

	// Delete the evaluations
	return db.Where("risk_register_id IN ? AND type = ?", riskRegisterIDs, "residual").Delete(&models.RiskEvaluation{}).Error
}

func (repo *RiskUpdatePlanRiskControlAssociationRepository) CreateBulkResidualEvaluations(evaluations []*models.RiskEvaluation, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if len(evaluations) == 0 {
		return nil
	}

	return db.Create(&evaluations).Error
}

// GetRiskControlAssociationList retrieves risks from both risk_register and risk_update_plan_risks tables
func (repo *RiskUpdatePlanRiskControlAssociationRepository) GetRiskControlAssociationList(planID uint, isCritical bool) ([]*riskDTOs.RiskControlAssociationListDTO, error) {
	var risks []*riskDTOs.RiskControlAssociationListDTO

	// Build UNION query to get risk IDs from both sources
	unionQuery := `
		SELECT DISTINCT rr.id
		FROM risk_register AS rr
		WHERE rr.risk_register_update_plan_id = ?
		
		UNION
		
		SELECT DISTINCT rr.id
		FROM risk_update_plan_risks AS rupr
		INNER JOIN risk_register AS rr ON rr.id = rupr.risk_register_id
		INNER JOIN risk_update_plan_timelines AS rupt ON rupt.id = rupr.risk_update_plan_timeline_id
		INNER JOIN risk_update_plan_timeline_types AS ruptt ON ruptt.id = rupt.risk_update_plan_timeline_type_id
		WHERE rupt.risk_update_plan_id = ? AND ruptt.type = 'erm_department_evaluation'
	`

	var riskIDs []uint
	if err := repo.db.Raw(unionQuery, planID, planID).Pluck("id", &riskIDs).Error; err != nil {
		return nil, err
	}

	if len(riskIDs) == 0 {
		return []*riskDTOs.RiskControlAssociationListDTO{}, nil
	}

	// Build main query with all joins, preferring plan draft inherent evaluation
	query := repo.db.Table("risk_register").
		Select(`
			risk_register.id,
			risk_register.reference_number,
			risk_register.title,
			risk_register.description,
			CONVERT(VARCHAR, risk_register.reported_at, 23) AS reported_at,
			risk_register.created_at,
			risk_register.updated_at,
			risk_register.risk_code,
			CONVERT(VARCHAR, risk_register.registration_date, 23) AS registration_date,
			risk_register.is_strategic,
			risk_register.is_operational,
			risk_register.expected_cause,
			risk_register.consequence,
			risk_register.is_internal,
			risk_register.is_external,
			risk_register.current_risk_rating_value_id,
			risk_categories.title_en AS risk_category_title_en,
			risk_categories.title_ar AS risk_category_title_ar,
			risk_sources.title_en AS risk_source_title_en,
			risk_sources.title_ar AS risk_source_title_ar,
			organization_units.id AS responsible_unit_id,
			organization_units.title_en AS responsible_unit_title_en,
			organization_units.title_ar AS responsible_unit_title_ar,
			users.id AS responsible_user_id,
			users.first_name_en AS responsible_user_first_name_en,
			users.first_name_ar AS responsible_user_first_name_ar,
			users.last_name_en AS responsible_user_last_name_en,
			users.last_name_ar AS responsible_user_last_name_ar,
			users.initial AS responsible_user_initial,
			users.color AS responsible_user_color,
			users.email AS responsible_user_email,
			users.profile_image_path AS responsible_user_profile_image_path,
			risk_statuses.title_en AS risk_status_title_en,
			risk_statuses.title_ar AS risk_status_title_ar,
			risk_response_strategies.id AS risk_response_strategy_id,
			risk_response_strategies.title_en AS risk_response_strategy_title_en,
			risk_response_strategies.title_ar AS risk_response_strategy_title_ar,
			risk_statuses.color_code AS risk_status_color_code,
			risk_statuses.type AS risk_status_type,
			review_frequencies.title_en AS review_frequency_title_en,
			review_frequencies.title_ar AS review_frequency_title_ar,
			severity_levels.title_en AS severity_level_title_en,
			severity_levels.title_ar AS severity_level_title_ar,
			CONVERT(VARCHAR, risk_register.next_review_date, 23) AS next_review_date,
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
			updated_user.profile_image_path AS updated_user_profile_image_path,
			current_likelihood.title_en AS current_likelihood_title_en,
			current_likelihood.title_ar AS current_likelihood_title_ar,
			current_impact.title_en AS current_impact_title_en,
			current_impact.title_ar AS current_impact_title_ar,
			current_risk_rating.title_en AS current_risk_rating_title_en,
			current_risk_rating.title_ar AS current_risk_rating_title_ar,
			current_risk_rating_value.value AS current_risk_rating_value,
			current_risk_rating.color_code AS current_risk_rating_color_code,
			inherent_likelihood.title_en AS inherent_likelihood_title_en,
			inherent_likelihood.title_ar AS inherent_likelihood_title_ar,
			inherent_impact.title_en AS inherent_impact_title_en,
			inherent_impact.title_ar AS inherent_impact_title_ar,
			inherent_risk_rating.title_en AS inherent_risk_rating_title_en,
			inherent_risk_rating.title_ar AS inherent_risk_rating_title_ar,
			inherent_risk_rating_value.value AS inherent_risk_rating_value,
			inherent_risk_rating.color_code AS inherent_risk_rating_color_code,
			residual_likelihood.title_en AS residual_likelihood_title_en,
			residual_likelihood.title_ar AS residual_likelihood_title_ar,
			residual_impact.title_en AS residual_impact_title_en,
			residual_impact.title_ar AS residual_impact_title_ar,
			residual_risk_rating.title_en AS residual_risk_rating_title_en,
			residual_risk_rating.title_ar AS residual_risk_rating_title_ar,
			residual_risk_rating_value.value AS residual_risk_rating_value,
			residual_risk_rating.color_code AS residual_risk_rating_color_code,
			divisions.id AS division_id,
			divisions.title_en AS division_title_en,
			divisions.title_ar AS division_title_ar,
			sections.id AS section_id,
			sections.title_en AS section_title_en,
			sections.title_ar AS section_title_ar,
			sub_sections.id AS sub_section_id,
			sub_sections.title_en AS sub_section_title_en,
			sub_sections.title_ar AS sub_section_title_ar
		`).
		Joins(`
			LEFT JOIN risk_categories ON risk_categories.id = risk_register.risk_category_id
			LEFT JOIN risk_sources ON risk_sources.id = risk_register.risk_source_id
			LEFT JOIN organization_units ON organization_units.id = risk_register.responsible_unit_id
			LEFT JOIN users ON users.id = risk_register.responsible_user_id
			LEFT JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id
			LEFT JOIN risk_response_strategies ON risk_response_strategies.id = risk_register.risk_response_strategy_id
			LEFT JOIN users AS created_user ON created_user.id = risk_register.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = risk_register.updated_by_id
			LEFT JOIN review_frequencies ON review_frequencies.id = risk_register.review_frequency_id
			LEFT JOIN severity_levels ON severity_levels.id = risk_register.severity_level_id
			LEFT JOIN risk_matrix_likelihoods AS current_likelihood ON current_likelihood.id = risk_register.current_likelihood_id
			LEFT JOIN risk_matrix_impacts AS current_impact ON current_impact.id = risk_register.current_impact_id
			LEFT JOIN risk_matrix_risk_ratings AS current_risk_rating ON current_risk_rating.id = risk_register.current_risk_rating_id
			LEFT JOIN risk_matrix_risk_rating_values AS current_risk_rating_value ON current_risk_rating_value.id = risk_register.current_risk_rating_value_id
			LEFT JOIN risk_evaluations AS plan_inherent_eval ON plan_inherent_eval.id = (SELECT MAX(id) FROM risk_evaluations WHERE risk_register_id = risk_register.id AND type = 'inherent' AND risk_register_update_plan_id = ? AND is_latest = 0)
			LEFT JOIN risk_evaluations AS live_inherent_eval ON live_inherent_eval.id = (
				SELECT MAX(id)
				FROM risk_evaluations
				WHERE risk_register_id = risk_register.id
					AND type = 'inherent'
					AND risk_register_update_plan_id IS NULL
					AND is_latest = 1
			)
			LEFT JOIN risk_matrix_likelihoods AS inherent_likelihood ON inherent_likelihood.id = COALESCE(plan_inherent_eval.risk_matrix_likelihood_id, live_inherent_eval.risk_matrix_likelihood_id, risk_register.inherent_likelihood_id)
			LEFT JOIN risk_matrix_impacts AS inherent_impact ON inherent_impact.id = COALESCE(plan_inherent_eval.risk_matrix_impact_id, live_inherent_eval.risk_matrix_impact_id, risk_register.inherent_impact_id)
			LEFT JOIN risk_matrix_risk_ratings AS inherent_risk_rating ON inherent_risk_rating.id = COALESCE(plan_inherent_eval.risk_matrix_risk_rating_id, live_inherent_eval.risk_matrix_risk_rating_id, risk_register.inherent_risk_rating_id)
			LEFT JOIN risk_matrix_risk_rating_values AS inherent_risk_rating_value ON inherent_risk_rating_value.id = COALESCE(plan_inherent_eval.risk_matrix_risk_rating_value_id, live_inherent_eval.risk_matrix_risk_rating_value_id, risk_register.inherent_risk_rating_value_id)
			LEFT JOIN risk_evaluations AS plan_residual_eval ON plan_residual_eval.id = (SELECT MAX(id) FROM risk_evaluations WHERE risk_register_id = risk_register.id AND type = 'residual' AND risk_register_update_plan_id = ? AND is_latest = 0)
			LEFT JOIN risk_evaluations AS live_residual_eval ON live_residual_eval.id = (
				SELECT MAX(id)
				FROM risk_evaluations
				WHERE risk_register_id = risk_register.id
					AND type = 'residual'
					AND risk_register_update_plan_id IS NULL
					AND is_latest = 1
			)
			LEFT JOIN risk_matrix_likelihoods AS residual_likelihood ON residual_likelihood.id = COALESCE(plan_residual_eval.risk_matrix_likelihood_id, live_residual_eval.risk_matrix_likelihood_id, risk_register.residual_likelihood_id)
			LEFT JOIN risk_matrix_impacts AS residual_impact ON residual_impact.id = COALESCE(plan_residual_eval.risk_matrix_impact_id, live_residual_eval.risk_matrix_impact_id, risk_register.residual_impact_id)
			LEFT JOIN risk_matrix_risk_ratings AS residual_risk_rating ON residual_risk_rating.id = COALESCE(plan_residual_eval.risk_matrix_risk_rating_id, live_residual_eval.risk_matrix_risk_rating_id, risk_register.residual_risk_rating_id)
			LEFT JOIN risk_matrix_risk_rating_values AS residual_risk_rating_value ON residual_risk_rating_value.id = COALESCE(plan_residual_eval.risk_matrix_risk_rating_value_id, live_residual_eval.risk_matrix_risk_rating_value_id, risk_register.residual_risk_rating_value_id)
			LEFT JOIN organization_units AS divisions ON divisions.id = risk_register.division_id
			LEFT JOIN organization_units AS sections ON sections.id = risk_register.section_id
			LEFT JOIN organization_units AS sub_sections ON sub_sections.id = risk_register.sub_section_id
		`, planID, planID).
		Where("risk_register.id IN ?", riskIDs)

	// Apply critical filter if needed
	if isCritical {
		query = query.Where("residual_risk_rating.type IN ?", []string{"high", "critical"})
	}

	if err := query.Find(&risks).Error; err != nil {
		return nil, err
	}

	// Fetch control analysis and action plans for each risk
	for _, risk := range risks {
		// Get first control analysis (inherent type)
		var controlAnalysis riskDTOs.RiskControlAssociationAnalysisDTO
		err := repo.db.Table("risk_control_analyses AS rca").
			Select(`
				rca.id,
				rca.type,
				rca.risk_control_id,
				rca.score,
				rca.control_effectiveness_score,
				rca.control_effectiveness_rating_id,
				rca.risk_reduction,
				rca.control_efficiency_status_id,
				rca.control_implementation_status_id,
				c.id AS control_id,
				c.title_en AS control_title_en,
				c.title_ar AS control_title_ar,
				ces.title_en AS control_efficiency_status_title_en,
				ces.title_ar AS control_efficiency_status_title_ar,
				cis.title_en AS control_implementation_status_title_en,
				cis.title_ar AS control_implementation_status_title_ar,
				cer.title_en AS control_effectiveness_rating_title_en,
				cer.title_ar AS control_effectiveness_rating_title_ar
			`).
			Joins("INNER JOIN risk_controls AS rc ON rc.id = rca.risk_control_id").
			Joins("INNER JOIN controls AS c ON c.id = rc.control_id").
			Joins("LEFT JOIN control_efficiency_statuses AS ces ON ces.id = rca.control_efficiency_status_id").
			Joins("INNER JOIN control_implementation_statuses AS cis ON cis.id = rca.control_implementation_status_id").
			Joins("LEFT JOIN control_effectiveness_ratings AS cer ON cer.id = rca.control_effectiveness_rating_id").
			Where("rc.risk_register_id = ? AND rca.type = ?", risk.ID, "inherent").
			Order("rca.id ASC").
			Limit(1).
			Scan(&controlAnalysis).Error

		if err == nil && controlAnalysis.ID != 0 {
			risk.ControlAnalysis = &controlAnalysis
		}

		// Get all action plans for this risk with complete details
		var actionPlans []riskDTOs.RiskControlAssociationActionPlanDTO
		err = repo.db.Table("risk_action_plans AS rap").
			Select(`
				rap.id,
				rap.risk_register_id,
				rap.title,
				rap.description,
				rap.severity_level_id,
				sl.title_en AS severity_level_title_en,
				sl.title_ar AS severity_level_title_ar,
				rap.responsible_unit_id,
				ou.title_en AS responsible_unit_title_en,
				ou.title_ar AS responsible_unit_title_ar,
				rap.responsible_user_id,
				ru.first_name_en AS responsible_user_first_name_en,
				ru.first_name_ar AS responsible_user_first_name_ar,
				ru.last_name_en AS responsible_user_last_name_en,
				ru.last_name_ar AS responsible_user_last_name_ar,
				ru.initial AS responsible_user_initial,
				ru.color AS responsible_user_color,
				ru.email AS responsible_user_email,
				ru.profile_image_path AS responsible_user_profile_image_path,
				rap.comments,
				CONVERT(VARCHAR, rap.start_date, 23) AS start_date,
				CONVERT(VARCHAR, rap.target_date, 23) AS target_date,
				rap.cost,
				rap.action_plan_status_id,
				aps.title_en AS action_plan_status_title_en,
				aps.title_ar AS action_plan_status_title_ar,
				aps.type AS action_plan_status_type,
				aps.color_code AS action_plan_status_color_code,
				rap.review_user_id,
				rvu.first_name_en AS review_user_first_name_en,
				rvu.first_name_ar AS review_user_first_name_ar,
				rvu.last_name_en AS review_user_last_name_en,
				rvu.last_name_ar AS review_user_last_name_ar,
				rvu.initial AS review_user_initial,
				rvu.color AS review_user_color,
				rvu.email AS review_user_email,
				rvu.profile_image_path AS review_user_profile_image_path,
				rap.percentage,
				rap.notes,
				rap.delegatee_user_id,
				du.first_name_en AS delegatee_user_first_name_en,
				du.first_name_ar AS delegatee_user_first_name_ar,
				du.last_name_en AS delegatee_user_last_name_en,
				du.last_name_ar AS delegatee_user_last_name_ar,
				du.initial AS delegatee_user_initial,
				du.color AS delegatee_user_color,
				du.email AS delegatee_user_email,
				du.profile_image_path AS delegatee_user_profile_image_path,
				rap.created_at,
				rap.updated_at,
				rap.created_by_id AS created_user_id,
				cu.first_name_en AS created_user_first_name_en,
				cu.first_name_ar AS created_user_first_name_ar,
				cu.last_name_en AS created_user_last_name_en,
				cu.last_name_ar AS created_user_last_name_ar,
				cu.initial AS created_user_initial,
				cu.color AS created_user_color,
				cu.email AS created_user_email,
				cu.profile_image_path AS created_user_profile_image_path,
				rap.updated_by_id AS updated_user_id,
				uu.first_name_en AS updated_user_first_name_en,
				uu.first_name_ar AS updated_user_first_name_ar,
				uu.last_name_en AS updated_user_last_name_en,
				uu.last_name_ar AS updated_user_last_name_ar,
				uu.initial AS updated_user_initial,
				uu.color AS updated_user_color,
				uu.email AS updated_user_email,
				uu.profile_image_path AS updated_user_profile_image_path
			`).
			Joins("LEFT JOIN severity_levels AS sl ON sl.id = rap.severity_level_id").
			Joins("LEFT JOIN organization_units AS ou ON ou.id = rap.responsible_unit_id").
			Joins("INNER JOIN users AS ru ON ru.id = rap.responsible_user_id").
			Joins("INNER JOIN action_plan_statuses AS aps ON aps.id = rap.action_plan_status_id").
			Joins("LEFT JOIN users AS rvu ON rvu.id = rap.review_user_id").
			Joins("LEFT JOIN users AS du ON du.id = rap.delegatee_user_id").
			Joins("INNER JOIN users AS cu ON cu.id = rap.created_by_id").
			Joins("LEFT JOIN users AS uu ON uu.id = rap.updated_by_id").
			Where("rap.risk_register_id = ?", risk.ID).
			Find(&actionPlans).Error

		if err == nil {
			risk.ActionPlans = actionPlans
		} else {
			risk.ActionPlans = []riskDTOs.RiskControlAssociationActionPlanDTO{}
		}
	}

	return risks, nil
}

func (repo *RiskUpdatePlanRiskControlAssociationRepository) GetResidualEvaluationByRiskRegisterID(riskRegisterID uint, tx ...utility.DatabaseTransaction) (*models.RiskEvaluation, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var evaluation models.RiskEvaluation
	err := db.Where("risk_register_id = ? AND type = ?", riskRegisterID, "residual").First(&evaluation).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &evaluation, nil
}

func (repo *RiskUpdatePlanRiskControlAssociationRepository) GetResidualEvaluationByRiskRegisterIDAndPlanID(riskRegisterID uint, planID uint, tx ...utility.DatabaseTransaction) (*models.RiskEvaluation, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var evaluation models.RiskEvaluation
	err := db.Where("risk_register_id = ? AND type = ? AND risk_register_update_plan_id = ?", riskRegisterID, "residual", planID).First(&evaluation).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &evaluation, nil
}

func (repo *RiskUpdatePlanRiskControlAssociationRepository) CreateResidualEvaluation(evaluation *models.RiskEvaluation, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(evaluation).Error
}

func (repo *RiskUpdatePlanRiskControlAssociationRepository) UpdateResidualEvaluation(evaluation *models.RiskEvaluation, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Save(evaluation).Error
}

func (repo *RiskUpdatePlanRiskControlAssociationRepository) CreateResidualEvaluationHistory(history *models.RiskResidualEvaluationHistory, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(history).Error
}
