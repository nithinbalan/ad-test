package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskUpdatePlanERMDeptRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	// ERM Department Evaluation - Risks
	GetRisksByBusinessUnitID(businessUnitID uint) ([]uint, error)
	GetERMDeptEvaluationRisks(timelineID uint, params *riskDTOs.CommonListQueryParamsDTO) (*riskDTOs.ERMDeptEvaluationListResponseDTO, error)
	CreateRiskMappings(risks []*models.RiskUpdatePlanRisk, tx ...utility.DatabaseTransaction) error
	DeleteAllRiskMappingsByTimelineID(timelineID uint, tx ...utility.DatabaseTransaction) error

	// ERM Department Evaluation - Processes
	GetProcessesByBusinessUnitID(businessUnitID uint) ([]uint, error)
	GetERMDeptEvaluationProcesses(timelineID uint, params *riskDTOs.CommonListQueryParamsDTO) (*riskDTOs.ERMDeptEvaluationProcessListResponseDTO, error)
	CreateProcessMappings(processes []*models.RiskUpdatePlanProcess, tx ...utility.DatabaseTransaction) error
	DeleteAllProcessMappingsByTimelineID(timelineID uint, tx ...utility.DatabaseTransaction) error

	// ERM Department Evaluation - Controls
	GetControlsByBusinessUnitID(businessUnitID uint) ([]uint, error)
	GetERMDeptEvaluationControls(timelineID uint, params *riskDTOs.CommonListQueryParamsDTO) (*riskDTOs.ERMDeptEvaluationControlListResponseDTO, error)
	CreateControlMappings(controls []*models.RiskUpdatePlanControl, tx ...utility.DatabaseTransaction) error
	DeleteAllControlMappingsByTimelineID(timelineID uint, tx ...utility.DatabaseTransaction) error

	// ERM Department Evaluation Documents
	CreateERMDeptDocument(document *models.RiskUpdatePlanErmDeptDocument, tx ...utility.DatabaseTransaction) (uint, error)
	GetERMDeptDocumentByID(id uint) (*models.RiskUpdatePlanErmDeptDocument, error)
	GetERMDeptDocumentsByTimelineID(timelineID uint) ([]*models.RiskUpdatePlanErmDeptDocument, error)
	DeleteERMDeptDocument(id uint, tx ...utility.DatabaseTransaction) error
}

type RiskUpdatePlanERMDeptRepository struct {
	db *gorm.DB
}

func NewRiskUpdatePlanERMDeptRepository(db *gorm.DB) *RiskUpdatePlanERMDeptRepository {
	return &RiskUpdatePlanERMDeptRepository{
		db: db,
	}
}

func (repo *RiskUpdatePlanERMDeptRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

// ============================================
// ERM Department Evaluation - Risks
// ============================================

func (repo *RiskUpdatePlanERMDeptRepository) GetRisksByBusinessUnitID(businessUnitID uint) ([]uint, error) {
	var riskIDs []uint
	err := repo.db.Table("risk_register").
		Select("risk_register.id").
		Joins("LEFT JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id").
		Joins("LEFT JOIN risk_register_update_plans ON risk_register_update_plans.id = risk_register.risk_register_update_plan_id").
		Where("risk_register.responsible_unit_id = ?", businessUnitID).
		Where("risk_statuses.type NOT IN (?)", []string{"identified", "submitted", "rejected"}).
		Where("(risk_register_update_plans.approved_by_id IS NOT NULL OR risk_register_update_plans.id IS NULL)").
		Pluck("risk_register.id", &riskIDs).Error

	return riskIDs, err
}

func (repo *RiskUpdatePlanERMDeptRepository) GetERMDeptEvaluationRisks(timelineID uint, params *riskDTOs.CommonListQueryParamsDTO) (*riskDTOs.ERMDeptEvaluationListResponseDTO, error) {
	var risks []*riskDTOs.ERMDeptEvaluationRiskListDTO
	var totalCount int64

	query := repo.db.Table("risk_update_plan_risks AS rupr").
		Select(`
			rupr.id,
			rupr.risk_update_plan_timeline_id,
			rupr.risk_register_id,
			rr.reference_number,
			rr.title,
			rr.description,
			rr.risk_code,
			rr.registration_date,
			rr.is_strategic,
			rr.is_operational,
			rr.expected_cause,
			rr.consequence,
			rr.is_internal,
			rr.is_external,
			rr.reported_at,
			rr.next_review_date,
			rr.current_risk_rating_value_id,
			rc.title_en AS risk_category_title_en,
			rc.title_ar AS risk_category_title_ar,
			rs.title_en AS risk_source_title_en,
			rs.title_ar AS risk_source_title_ar,
			ru.id AS responsible_unit_id,
			ru.title_en AS responsible_unit_title_en,
			ru.title_ar AS responsible_unit_title_ar,
			u.id AS responsible_user_id,
			u.first_name_en AS responsible_user_first_name_en,
			u.first_name_ar AS responsible_user_first_name_ar,
			u.last_name_en AS responsible_user_last_name_en,
			u.last_name_ar AS responsible_user_last_name_ar,
			u.initial AS responsible_user_initial,
			u.color AS responsible_user_color,
			u.email AS responsible_user_email,
			u.profile_image_path AS responsible_user_profile_image_path,
			rst.title_en AS risk_status_title_en,
			rst.title_ar AS risk_status_title_ar,
			rst.type AS risk_status_type,
			rst.color_code AS risk_status_color_code,
			rrs.title_en AS risk_response_strategy_title_en,
			rrs.title_ar AS risk_response_strategy_title_ar,
			rf.title_en AS review_frequency_title_en,
			rf.title_ar AS review_frequency_title_ar,
			sl.title_en AS severity_level_title_en,
			sl.title_ar AS severity_level_title_ar,
			cll.title_en AS current_likelihood_title_en,
			cll.title_ar AS current_likelihood_title_ar,
			cil.title_en AS current_impact_title_en,
			cil.title_ar AS current_impact_title_ar,
			crr.title_en AS current_risk_rating_title_en,
			crr.title_ar AS current_risk_rating_title_ar,
			crrv.value AS current_risk_rating_value,
			crr.color_code AS current_risk_rating_color_code,
			ill.title_en AS inherent_likelihood_title_en,
			ill.title_ar AS inherent_likelihood_title_ar,
			iil.title_en AS inherent_impact_title_en,
			iil.title_ar AS inherent_impact_title_ar,
			irr.title_en AS inherent_risk_rating_title_en,
			irr.title_ar AS inherent_risk_rating_title_ar,
			irrv.value AS inherent_risk_rating_value,
			irr.color_code AS inherent_risk_rating_color_code,
			rll.title_en AS residual_likelihood_title_en,
			rll.title_ar AS residual_likelihood_title_ar,
			ril.title_en AS residual_impact_title_en,
			ril.title_ar AS residual_impact_title_ar,
			rrr.title_en AS residual_risk_rating_title_en,
			rrr.title_ar AS residual_risk_rating_title_ar,
			rrrv.value AS residual_risk_rating_value,
			rrr.color_code AS residual_risk_rating_color_code,
			div.id AS division_id,
			div.title_en AS division_title_en,
			div.title_ar AS division_title_ar,
			sec.id AS section_id,
			sec.title_en AS section_title_en,
			sec.title_ar AS section_title_ar,
			ss.id AS sub_section_id,
			ss.title_en AS sub_section_title_en,
			ss.title_ar AS sub_section_title_ar
		`).
		Joins("INNER JOIN risk_register AS rr ON rr.id = rupr.risk_register_id").
		Joins("LEFT JOIN risk_categories AS rc ON rc.id = rr.risk_category_id").
		Joins("LEFT JOIN risk_sources AS rs ON rs.id = rr.risk_source_id").
		Joins("LEFT JOIN organization_units AS ru ON ru.id = rr.responsible_unit_id").
		Joins("LEFT JOIN users AS u ON u.id = rr.responsible_user_id").
		Joins("LEFT JOIN risk_statuses AS rst ON rst.id = rr.risk_status_id").
		Joins("LEFT JOIN risk_response_strategies AS rrs ON rrs.id = rr.risk_response_strategy_id").
		Joins("LEFT JOIN review_frequencies AS rf ON rf.id = rr.review_frequency_id").
		Joins("LEFT JOIN severity_levels AS sl ON sl.id = rr.severity_level_id").
		Joins("LEFT JOIN risk_matrix_likelihoods AS cll ON cll.id = rr.current_likelihood_id").
		Joins("LEFT JOIN risk_matrix_impacts AS cil ON cil.id = rr.current_impact_id").
		Joins("LEFT JOIN risk_matrix_risk_ratings AS crr ON crr.id = rr.current_risk_rating_id").
		Joins("LEFT JOIN risk_matrix_risk_rating_values AS crrv ON crrv.id = rr.current_risk_rating_value_id").
		Joins("LEFT JOIN risk_matrix_likelihoods AS ill ON ill.id = rr.inherent_likelihood_id").
		Joins("LEFT JOIN risk_matrix_impacts AS iil ON iil.id = rr.inherent_impact_id").
		Joins("LEFT JOIN risk_matrix_risk_ratings AS irr ON irr.id = rr.inherent_risk_rating_id").
		Joins("LEFT JOIN risk_matrix_risk_rating_values AS irrv ON irrv.id = rr.inherent_risk_rating_value_id").
		Joins("LEFT JOIN risk_matrix_likelihoods AS rll ON rll.id = rr.residual_likelihood_id").
		Joins("LEFT JOIN risk_matrix_impacts AS ril ON ril.id = rr.residual_impact_id").
		Joins("LEFT JOIN risk_matrix_risk_ratings AS rrr ON rrr.id = rr.residual_risk_rating_id").
		Joins("LEFT JOIN risk_matrix_risk_rating_values AS rrrv ON rrrv.id = rr.residual_risk_rating_value_id").
		Joins("LEFT JOIN organization_units AS div ON div.id = rr.division_id").
		Joins("LEFT JOIN organization_units AS sec ON sec.id = rr.section_id").
		Joins("LEFT JOIN organization_units AS ss ON ss.id = rr.sub_section_id").
		Where("rupr.risk_update_plan_timeline_id = ?", timelineID)

	// Apply filters
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("(rr.title LIKE ? OR rr.reference_number LIKE ?)", keyword, keyword)
	}

	// Count total
	countQuery := query
	if err := countQuery.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply sorting
	if params.SortBy != "" {
		sortOrder := "ASC"
		if params.Sort == "desc" {
			sortOrder = "DESC"
		}
		query = query.Order(params.SortBy + " " + sortOrder)
	} else {
		query = query.Order("rupr.created_at DESC")
	}

	// Apply pagination
	if params.PageSize > 0 {
		offset := (params.Page - 1) * params.PageSize
		query = query.Limit(params.PageSize).Offset(offset)
	}

	if err := query.Scan(&risks).Error; err != nil {
		return nil, err
	}

	pagination := commonDTOs.PaginationInfoDTO{
		CurrentPage: params.Page,
		PageSize:    params.PageSize,
		TotalPages:  int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize)),
		TotalCount:  totalCount,
	}

	return &riskDTOs.ERMDeptEvaluationListResponseDTO{
		Items:      risks,
		Pagination: pagination,
	}, nil
}

func (repo *RiskUpdatePlanERMDeptRepository) CreateRiskMappings(risks []*models.RiskUpdatePlanRisk, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if len(risks) == 0 {
		return nil
	}

	return db.Create(&risks).Error
}

func (repo *RiskUpdatePlanERMDeptRepository) DeleteAllRiskMappingsByTimelineID(timelineID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("risk_update_plan_timeline_id = ?", timelineID).
		Delete(&models.RiskUpdatePlanRisk{}).Error
}

// ============================================
// ERM Department Evaluation - Processes
// ============================================

func (repo *RiskUpdatePlanERMDeptRepository) GetProcessesByBusinessUnitID(businessUnitID uint) ([]uint, error) {
	var processIDs []uint
	err := repo.db.Table("processes").
		Select("id").
		Where("responsible_unit_id = ?", businessUnitID).
		Pluck("id", &processIDs).Error

	return processIDs, err
}

func (repo *RiskUpdatePlanERMDeptRepository) GetERMDeptEvaluationProcesses(timelineID uint, params *riskDTOs.CommonListQueryParamsDTO) (*riskDTOs.ERMDeptEvaluationProcessListResponseDTO, error) {
	var processes []*riskDTOs.ERMDeptEvaluationProcessListDTO
	var totalCount int64

	query := repo.db.Table("risk_update_plan_processes AS rupp").
		Select(`
			rupp.id,
			rupp.risk_update_plan_timeline_id,
			rupp.process_id,
			p.reference_number,
			p.title_en,
			p.title_ar,
			p.description,
			p.marked_as_auditable_item,
			pg.title_en AS process_group_title_en,
			pg.title_ar AS process_group_title_ar,
			p.rto_from_value as rto_from_value,
			p.rto_to_value as rto_to_value,
			p.rto_unit_id as rto_unit_id,
			rto_bc_unit.title_en AS rto_unit_title_en,
			rto_bc_unit.title_ar AS rto_unit_title_ar,
			p.rto_condition as rto_condition,
			p.mtpd_from_value as mtpd_from_value,
			p.mtpd_to_value as mtpd_to_value,
			p.mtpd_unit_id as mtpd_unit_id,
			mtpd_bc_unit.title_en AS mtpd_unit_title_en,
			mtpd_bc_unit.title_ar AS mtpd_unit_title_ar,
			p.mtpd_condition as mtpd_condition,
			p.rpo_from_value as rpo_from_value,
			p.rpo_to_value as rpo_to_value,
			p.rpo_unit_id as rpo_unit_id,
			rpo_bc_unit.title_en AS rpo_unit_title_en,
			rpo_bc_unit.title_ar AS rpo_unit_title_ar,
			p.rpo_condition as rpo_condition,
			p.rsp_from_value as rsp_from_value,
			p.rsp_to_value as rsp_to_value,
			p.rsp_unit_id as rsp_unit_id,
			rsp_bc_unit.title_en AS rsp_unit_title_en,
			rsp_bc_unit.title_ar AS rsp_unit_title_ar,
			p.rsp_condition as rsp_condition,
			p.impact_analysis_level_id AS impact_analysis_level_id,
			bcm_impact_analysis_levels.title_en AS impact_analysis_level_title_en,
			bcm_impact_analysis_levels.title_ar AS impact_analysis_level_title_ar,
			bcm_impact_analysis_levels.color_code AS impact_analysis_level_color_code,
			bcm_impact_analysis_levels.type AS impact_analysis_level_type,
			p.is_bia_performed,
			p.last_assessed_at,
			last_assessed_user.id AS last_assessed_user_id,
			last_assessed_user.first_name_en AS last_assessed_user_first_name_en,
			last_assessed_user.first_name_ar AS last_assessed_user_first_name_ar,
			last_assessed_user.last_name_en AS last_assessed_user_last_name_en,
			last_assessed_user.last_name_ar AS last_assessed_user_last_name_ar,
			last_assessed_user.initial AS last_assessed_user_initial,
			last_assessed_user.color AS last_assessed_user_color,
			last_assessed_user.email AS last_assessed_user_email,
			last_assessed_user.profile_image_path AS last_assessed_user_profile_image_path,
			ou.title_en AS owner_unit_title_en,
			ou.title_ar AS owner_unit_title_ar,
			u.id AS owner_id,
			u.first_name_en AS owner_first_name_en,
			u.first_name_ar AS owner_first_name_ar,
			u.last_name_en AS owner_last_name_en,
			u.last_name_ar AS owner_last_name_ar,
			u.initial AS owner_initial,
			u.color AS owner_color,
			u.email AS owner_email,
			u.profile_image_path AS owner_profile_image_path,
			ls.title_en AS lifecycle_status_title_en,
			ls.title_ar AS lifecycle_status_title_ar,
			ls.color_code AS lifecycle_status_color_code,
			pcl.id AS process_criticality_level_id,
			pcl.title_en AS process_criticality_level_title_en,
			pcl.title_ar AS process_criticality_level_title_ar,
			cu.id AS created_user_id,
			cu.first_name_en AS created_user_first_name_en,
			cu.first_name_ar AS created_user_first_name_ar,
			cu.last_name_en AS created_user_last_name_en,
			cu.last_name_ar AS created_user_last_name_ar,
			cu.initial AS created_user_initial,
			cu.color AS created_user_color,
			cu.email AS created_user_email,
			cu.profile_image_path AS created_user_profile_image_path,
			uu.id AS updated_user_id,
			uu.first_name_en AS updated_user_first_name_en,
			uu.first_name_ar AS updated_user_first_name_ar,
			uu.last_name_en AS updated_user_last_name_en,
			uu.last_name_ar AS updated_user_last_name_ar,
			uu.initial AS updated_user_initial,
			uu.color AS updated_user_color,
			uu.email AS updated_user_email,
			uu.profile_image_path AS updated_user_profile_image_path,
			p.created_at,
			p.updated_at,
			rmrr.title_en AS risk_rating_title_en,
			rmrr.title_ar AS risk_rating_title_ar,
			rmrr.color_code AS risk_rating_color_code,
			rmrrv.value AS risk_rating_value,
			rl.id AS risk_level_id,
			rl.title_en AS risk_level_title_en,
			rl.title_ar AS risk_level_title_ar,
			rl.color_code AS risk_level_color_code,
			rl.type AS risk_level_type
		`).
		Joins("INNER JOIN processes AS p ON p.id = rupp.process_id").
		Joins("LEFT JOIN process_groups AS pg ON pg.id = p.process_group_id").
		Joins("LEFT JOIN organization_units AS ou ON ou.id = p.responsible_unit_id").
		Joins("LEFT JOIN users AS u ON u.id = p.responsible_user_id").
		Joins("LEFT JOIN lifecycle_statuses AS ls ON ls.id = p.lifecycle_status_id").
		Joins("LEFT JOIN users AS cu ON cu.id = p.created_by_id").
		Joins("LEFT JOIN users AS uu ON uu.id = p.updated_by_id").
		Joins("LEFT JOIN users AS last_assessed_user ON last_assessed_user.id = p.last_assessed_by_id").
		Joins("LEFT JOIN business_continuity_duration_units AS rto_bc_unit ON rto_bc_unit.id = p.rto_unit_id").
		Joins("LEFT JOIN business_continuity_duration_units AS mtpd_bc_unit ON mtpd_bc_unit.id = p.mtpd_unit_id").
		Joins("LEFT JOIN business_continuity_duration_units AS rpo_bc_unit ON rpo_bc_unit.id = p.rpo_unit_id").
		Joins("LEFT JOIN business_continuity_duration_units AS rsp_bc_unit ON rsp_bc_unit.id = p.rsp_unit_id").
		Joins("LEFT JOIN bcm_impact_analysis_levels ON bcm_impact_analysis_levels.id = p.impact_analysis_level_id").
		Joins("LEFT JOIN risk_matrix_risk_ratings AS rmrr ON rmrr.id = p.risk_matrix_risk_rating_id").
		Joins("LEFT JOIN risk_matrix_risk_rating_values AS rmrrv ON rmrrv.id = p.risk_matrix_risk_rating_value_id").
		Joins("LEFT JOIN risk_levels AS rl ON rl.id = p.risk_level_id").
		Joins("LEFT JOIN process_criticality_levels AS pcl ON pcl.id = p.process_criticality_level_id").
		Where("rupp.risk_update_plan_timeline_id = ?", timelineID)

	// Apply filters
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`(p.reference_number LIKE ? OR 
			p.title_en LIKE ? OR 
			p.title_ar LIKE ? OR 
			p.description LIKE ? OR 
			pg.title_en LIKE ? OR 
			pg.title_ar LIKE ? OR 
			ou.title_en LIKE ? OR 
			ou.title_ar LIKE ? OR 
			u.first_name_en LIKE ? OR 
			u.first_name_ar LIKE ? OR 
			ls.title_en LIKE ? OR 
			ls.title_ar LIKE ?)`,
			keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword)
	}

	// Count total
	countQuery := query
	if err := countQuery.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply sorting
	if params.SortBy != "" {
		sortOrder := "ASC"
		if params.Sort == "desc" {
			sortOrder = "DESC"
		}
		query = query.Order(params.SortBy + " " + sortOrder)
	} else {
		query = query.Order("p.created_at DESC")
	}

	// Apply pagination
	if params.PageSize > 0 {
		offset := (params.Page - 1) * params.PageSize
		query = query.Limit(params.PageSize).Offset(offset)
	}

	if err := query.Scan(&processes).Error; err != nil {
		return nil, err
	}

	// Fetch activities for each process
	for _, process := range processes {
		var activities []riskDTOs.ProcessActivityDTO

		err := repo.db.
			Table("process_activities").
			Select("id, title_en, title_ar, description_en, description_ar, sequence, estimated_duration, expected_output, created_at, updated_at").
			Where("process_id = ?", process.ProcessID).
			Order("sequence ASC").
			Scan(&activities).Error

		if err != nil {
			return nil, err
		}

		process.Activities = activities
	}

	pagination := commonDTOs.PaginationInfoDTO{
		CurrentPage: params.Page,
		PageSize:    params.PageSize,
		TotalPages:  int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize)),
		TotalCount:  totalCount,
	}

	return &riskDTOs.ERMDeptEvaluationProcessListResponseDTO{
		Items:      processes,
		Pagination: pagination,
	}, nil
}

func (repo *RiskUpdatePlanERMDeptRepository) CreateProcessMappings(processes []*models.RiskUpdatePlanProcess, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if len(processes) == 0 {
		return nil
	}

	return db.Create(&processes).Error
}

func (repo *RiskUpdatePlanERMDeptRepository) DeleteAllProcessMappingsByTimelineID(timelineID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("risk_update_plan_timeline_id = ?", timelineID).
		Delete(&models.RiskUpdatePlanProcess{}).Error
}

// ============================================
// ERM Department Evaluation - Controls
// ============================================

func (repo *RiskUpdatePlanERMDeptRepository) GetControlsByBusinessUnitID(businessUnitID uint) ([]uint, error) {
	var controlIDs []uint
	err := repo.db.Table("controls").
		Select("id").
		Where("responsible_unit_id = ?", businessUnitID).
		Pluck("id", &controlIDs).Error

	return controlIDs, err
}

func (repo *RiskUpdatePlanERMDeptRepository) GetERMDeptEvaluationControls(timelineID uint, params *riskDTOs.CommonListQueryParamsDTO) (*riskDTOs.ERMDeptEvaluationControlListResponseDTO, error) {
	var controls []*riskDTOs.ERMDeptEvaluationControlListDTO
	var totalCount int64

	query := repo.db.Table("risk_update_plan_controls AS rupc").
		Select(`
			rupc.id,
			rupc.risk_update_plan_timeline_id,
			rupc.control_id,
			ct.title_en AS control_type_title_en,
			ct.title_ar AS control_type_title_ar,
			ct.type AS control_type,
			cc.title_en AS control_category_title_en,
			cc.title_ar AS control_category_title_ar,
			ou.title_en AS owner_unit_title_en,
			ou.title_ar AS owner_unit_title_ar,
			owner.id AS owner_id,
			owner.first_name_en AS owner_first_name_en,
			owner.first_name_ar AS owner_first_name_ar,
			owner.last_name_en AS owner_last_name_en,
			owner.last_name_ar AS owner_last_name_ar,
			owner.initial AS owner_initial,
			owner.color AS owner_color,
			owner.email AS owner_email,
			owner.profile_image_path AS owner_profile_image_path,
			cis.title_en AS control_implementation_status_title_en,
			cis.title_ar AS control_implementation_status_title_ar,
			cis.color_code AS control_implementation_status_color_code,
			ces.title_en AS control_efficiency_status_title_en,
			ces.title_ar AS control_efficiency_status_title_ar,
			ces.color_code AS control_efficiency_status_color_code,
			c.reference_number,
			c.code,
			c.title_en,
			c.title_ar,
			c.description,
			c.objective,
			rf.title_en AS review_frequency_title_en,
			rf.title_ar AS review_frequency_title_ar,
			c.last_review_date,
			c.next_review_date,
			c.last_review_note,
			tru.id AS test_responsible_user_id,
			tru.first_name_en AS test_responsible_user_first_name_en,
			tru.first_name_ar AS test_responsible_user_first_name_ar,
			tru.last_name_en AS test_responsible_user_last_name_en,
			tru.last_name_ar AS test_responsible_user_last_name_ar,
			tru.initial AS test_responsible_user_initial,
			tru.color AS test_responsible_user_color,
			tru.email AS test_responsible_user_email,
			tru.profile_image_path AS test_responsible_user_profile_image_path,
			tf.title_en AS test_frequency_title_en,
			tf.title_ar AS test_frequency_title_ar,
			c.last_test_date,
			c.next_test_date,
			c.implemented_at,
			c.is_applicable,
			c.justification,
			c.implementation_guideline,
			c.applicability_change_date,
			ls.title_en AS lifecycle_status_title_en,
			ls.title_ar AS lifecycle_status_title_ar,
			ls.color_code AS lifecycle_status_color_code,
			cs.id AS control_status_id,
			cs.title_en AS control_status_title_en,
			cs.title_ar AS control_status_title_ar,
			cs.color_code AS control_status_color_code,
			cs.type AS control_status_type,
			cu.id AS created_user_id,
			cu.first_name_en AS created_user_first_name_en,
			cu.first_name_ar AS created_user_first_name_ar,
			cu.last_name_en AS created_user_last_name_en,
			cu.last_name_ar AS created_user_last_name_ar,
			cu.initial AS created_user_initial,
			cu.color AS created_user_color,
			cu.email AS created_user_email,
			cu.profile_image_path AS created_user_profile_image_path,
			uu.id AS updated_user_id,
			uu.first_name_en AS updated_user_first_name_en,
			uu.first_name_ar AS updated_user_first_name_ar,
			uu.last_name_en AS updated_user_last_name_en,
			uu.last_name_ar AS updated_user_last_name_ar,
			uu.initial AS updated_user_initial,
			uu.color AS updated_user_color,
			uu.email AS updated_user_email,
			uu.profile_image_path AS updated_user_profile_image_path,
			c.created_at,
			c.updated_at,
			rmrr.title_en AS risk_rating_title_en,
			rmrr.title_ar AS risk_rating_title_ar,
			rmrr.color_code AS risk_rating_color_code,
			rmrrv.value AS risk_rating_value,
			ca.id AS control_automation_id,
			ca.title_en AS control_automation_title_en,
			ca.title_ar AS control_automation_title_ar,
			ca.type AS control_automation_type,
			ca.score AS control_automation_score,
			ca.status AS control_automation_status,
			ccon.id AS control_consistency_id,
			ccon.title_en AS control_consistency_title_en,
			ccon.title_ar AS control_consistency_title_ar,
			ccon.type AS control_consistency_type,
			ccon.score AS control_consistency_score,
			ccon.status AS control_consistency_status,
			ccon.description_en AS control_consistency_description_en,
			ccon.description_ar AS control_consistency_description_ar,
			c.control_effectiveness_score,
			cer.id AS control_effectiveness_rating_id,
			cer.title_en AS control_effectiveness_rating_title_en,
			cer.title_ar AS control_effectiveness_rating_title_ar,
			cer.min_score AS control_effectiveness_rating_min_score,
			cer.max_score AS control_effectiveness_rating_max_score,
			cer.type AS control_effectiveness_rating_type,
			cer.color_code AS control_effectiveness_rating_color_code,
			c.risk_reduction,
			latest_ct.id AS control_test_id,
			latest_ct.title AS control_test_title,
			latest_ct.description AS control_test_description,
			latest_ct.comments AS control_test_comments,
			latest_ct.score AS control_test_score,
			latest_ct.control_effectiveness_score AS control_test_effectiveness_score,
			latest_ct.risk_reduction AS control_test_risk_reduction,
			CONVERT(VARCHAR, latest_ct.tested_at, 23) AS control_test_tested_at,
			ct_tested_by.id AS control_test_tested_by_id,
			ct_tested_by.first_name_en AS control_test_tested_by_first_name_en,
			ct_tested_by.first_name_ar AS control_test_tested_by_first_name_ar,
			ct_tested_by.last_name_en AS control_test_tested_by_last_name_en,
			ct_tested_by.last_name_ar AS control_test_tested_by_last_name_ar,
			ct_tested_by.initial AS control_test_tested_by_initial,
			ct_tested_by.color AS control_test_tested_by_color,
			ct_tested_by.email AS control_test_tested_by_email,
			ct_tested_by.profile_image_path AS control_test_tested_by_profile_image_path
		`).
		Joins("INNER JOIN controls AS c ON c.id = rupc.control_id").
		Joins("LEFT JOIN control_types AS ct ON ct.id = c.control_type_id").
		Joins("LEFT JOIN control_categories AS cc ON cc.id = c.control_category_id").
		Joins("LEFT JOIN organization_units AS ou ON ou.id = c.responsible_unit_id").
		Joins("LEFT JOIN users AS owner ON owner.id = c.responsible_user_id").
		Joins("LEFT JOIN control_implementation_statuses AS cis ON cis.id = c.control_implementation_status_id").
		Joins("LEFT JOIN control_efficiency_statuses AS ces ON ces.id = c.control_efficiency_status_id").
		Joins("LEFT JOIN review_frequencies AS rf ON rf.id = c.review_frequency_id").
		Joins("LEFT JOIN users AS tru ON tru.id = c.test_responsible_user_id").
		Joins("LEFT JOIN review_frequencies AS tf ON tf.id = c.test_frequency_id").
		Joins("LEFT JOIN control_statuses AS cs ON cs.id = c.control_status_id").
		Joins("LEFT JOIN lifecycle_statuses AS ls ON ls.id = c.lifecycle_status_id").
		Joins("LEFT JOIN users AS cu ON cu.id = c.created_by_id").
		Joins("LEFT JOIN users AS uu ON uu.id = c.updated_by_id").
		Joins("LEFT JOIN risk_matrix_risk_ratings AS rmrr ON rmrr.id = c.risk_matrix_risk_rating_id").
		Joins("LEFT JOIN risk_matrix_risk_rating_values AS rmrrv ON rmrrv.id = c.risk_matrix_risk_rating_value_id").
		Joins("LEFT JOIN control_effectiveness_ratings AS cer ON cer.id = c.control_effectiveness_rating_id").
		Joins(`LEFT JOIN (
			SELECT ct.control_id, MAX(ct.id) as latest_test_id
			FROM control_tests ct
			GROUP BY ct.control_id
		) AS latest_ct_ids ON latest_ct_ids.control_id = c.id`).
		Joins("LEFT JOIN control_tests AS latest_ct ON latest_ct.id = latest_ct_ids.latest_test_id").
		Joins("LEFT JOIN control_automations AS ca ON ca.id = latest_ct.control_automation_id").
		Joins("LEFT JOIN control_consistencies AS ccon ON ccon.id = latest_ct.control_consistency_id").
		Joins("LEFT JOIN users AS ct_tested_by ON ct_tested_by.id = latest_ct.tested_by").
		Where("rupc.risk_update_plan_timeline_id = ?", timelineID)

	// Apply filters
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`(c.reference_number LIKE ? OR 
			c.code LIKE ? OR 
			c.title_en LIKE ? OR 
			c.title_ar LIKE ? OR 
			c.description LIKE ? OR 
			ct.title_en LIKE ? OR 
			ct.title_ar LIKE ? OR 
			cc.title_en LIKE ? OR 
			cc.title_ar LIKE ? OR 
			ou.title_en LIKE ? OR 
			ou.title_ar LIKE ? OR 
			owner.first_name_en LIKE ? OR 
			owner.first_name_ar LIKE ? OR 
			ls.title_en LIKE ? OR 
			ls.title_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword)
	}

	// Apply is_tested filter
	if params.IsTested != nil {
		if *params.IsTested {
			// Show only tested controls (controls with control tests)
			query = query.Where("latest_ct.id IS NOT NULL")
		} else {
			// Show only untested controls (controls without control tests)
			query = query.Where("latest_ct.id IS NULL")
		}
	}

	// Count total
	countQuery := query
	if err := countQuery.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply sorting
	if params.SortBy != "" {
		sortOrder := "ASC"
		if params.Sort == "desc" {
			sortOrder = "DESC"
		}
		query = query.Order(params.SortBy + " " + sortOrder)
	} else {
		query = query.Order("c.created_at DESC")
	}

	// Apply pagination
	if params.PageSize > 0 {
		offset := (params.Page - 1) * params.PageSize
		query = query.Limit(params.PageSize).Offset(offset)
	}

	if err := query.Scan(&controls).Error; err != nil {
		return nil, err
	}

	// Load control test documents for each control
	for _, control := range controls {
		if control.ControlTestID != nil {
			var documents []organizationDTOs.ControlTestDocumentDTO

			err := repo.db.
				Table("control_test_documents").
				Select("id, control_test_id, url, name, original_name, token, type, size").
				Where("control_test_id = ?", *control.ControlTestID).
				Scan(&documents).Error

			if err != nil {
				return nil, err
			}

			if documents == nil {
				documents = []organizationDTOs.ControlTestDocumentDTO{}
			}

			control.ControlTestDocuments = documents
		} else {
			control.ControlTestDocuments = []organizationDTOs.ControlTestDocumentDTO{}
		}
	}

	pagination := commonDTOs.PaginationInfoDTO{
		CurrentPage: params.Page,
		PageSize:    params.PageSize,
		TotalPages:  int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize)),
		TotalCount:  totalCount,
	}

	return &riskDTOs.ERMDeptEvaluationControlListResponseDTO{
		Items:      controls,
		Pagination: pagination,
	}, nil
}

func (repo *RiskUpdatePlanERMDeptRepository) CreateControlMappings(controls []*models.RiskUpdatePlanControl, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if len(controls) == 0 {
		return nil
	}

	return db.Create(&controls).Error
}

func (repo *RiskUpdatePlanERMDeptRepository) DeleteAllControlMappingsByTimelineID(timelineID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("risk_update_plan_timeline_id = ?", timelineID).
		Delete(&models.RiskUpdatePlanControl{}).Error
}

// ERM Department Evaluation Documents

func (repo *RiskUpdatePlanERMDeptRepository) CreateERMDeptDocument(document *models.RiskUpdatePlanErmDeptDocument, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(document).Error; err != nil {
		return 0, err
	}

	return document.ID, nil
}

func (repo *RiskUpdatePlanERMDeptRepository) GetERMDeptDocumentByID(id uint) (*models.RiskUpdatePlanErmDeptDocument, error) {
	var document models.RiskUpdatePlanErmDeptDocument
	if err := repo.db.First(&document, id).Error; err != nil {
		return nil, err
	}

	return &document, nil
}

func (repo *RiskUpdatePlanERMDeptRepository) GetERMDeptDocumentsByTimelineID(timelineID uint) ([]*models.RiskUpdatePlanErmDeptDocument, error) {
	var documents []*models.RiskUpdatePlanErmDeptDocument
	if err := repo.db.Where("risk_update_plan_timeline_id = ?", timelineID).Find(&documents).Error; err != nil {
		return nil, err
	}

	return documents, nil
}

func (repo *RiskUpdatePlanERMDeptRepository) DeleteERMDeptDocument(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Delete(&models.RiskUpdatePlanErmDeptDocument{}, id).Error; err != nil {
		return err
	}

	return nil
}
