package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/utility"
	"time"

	"gorm.io/gorm"
)

type RiskRegisterUpdatePlanRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *riskDTOs.RiskRegisterUpdatePlanListQueryParamsDTO) (*riskDTOs.RiskRegisterUpdatePlanListResponseDTO, error)
	GetByID(id uint) (*models.RiskRegisterUpdatePlan, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskRegisterUpdatePlan, error)
	Create(plan *models.RiskRegisterUpdatePlan, tx ...utility.DatabaseTransaction) (uint, error)
	Update(plan *models.RiskRegisterUpdatePlan, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error

	// Team Members
	CreateTeamMembers(members []*models.RiskRegisterUpdatePlanTeamMember, tx ...utility.DatabaseTransaction) error
	DeleteTeamMembersByPlanID(planID uint, tx ...utility.DatabaseTransaction) error
	GetTeamMembersByPlanID(planID uint, tx ...utility.DatabaseTransaction) ([]*models.RiskRegisterUpdatePlanTeamMember, error)

	// ERM Department Evaluation - Risks
	GetRisksByBusinessUnitID(businessUnitID uint) ([]uint, error)
	GetERMDeptEvaluationRisks(timelineID uint, params *riskDTOs.ERMDeptEvaluationListQueryParamsDTO) (*riskDTOs.ERMDeptEvaluationListResponseDTO, error)
	CreateRiskMappings(risks []*models.RiskUpdatePlanRisk, tx ...utility.DatabaseTransaction) error
	DeleteAllRiskMappingsByTimelineID(timelineID uint, tx ...utility.DatabaseTransaction) error

	// ERM Department Evaluation - Processes
	GetProcessesByBusinessUnitID(businessUnitID uint) ([]uint, error)
	GetERMDeptEvaluationProcesses(timelineID uint, params *riskDTOs.ERMDeptEvaluationListQueryParamsDTO) (*riskDTOs.ERMDeptEvaluationProcessListResponseDTO, error)
	CreateProcessMappings(processes []*models.RiskUpdatePlanProcess, tx ...utility.DatabaseTransaction) error
	DeleteAllProcessMappingsByTimelineID(timelineID uint, tx ...utility.DatabaseTransaction) error

	// ERM Department Evaluation - Controls
	GetControlsByBusinessUnitID(businessUnitID uint) ([]uint, error)
	GetERMDeptEvaluationControls(timelineID uint, params *riskDTOs.ERMDeptEvaluationListQueryParamsDTO) (*riskDTOs.ERMDeptEvaluationControlListResponseDTO, error)
	CreateControlMappings(controls []*models.RiskUpdatePlanControl, tx ...utility.DatabaseTransaction) error
	DeleteAllControlMappingsByTimelineID(timelineID uint, tx ...utility.DatabaseTransaction) error

	// Timeline Operations
	GetTimelineByPlanIDAndTypeID(planID uint, timelineTypeID uint, tx ...utility.DatabaseTransaction) (*models.RiskUpdatePlanTimeline, error)
	UpdateTimelineCompletion(timelineID uint, currentUserID uint, tx ...utility.DatabaseTransaction) error

	// Introductory Meeting Documents
	CreateIntroductoryMeetingDocument(document *models.RiskUpdatePlanIntroductoryMeetingDocument, tx ...utility.DatabaseTransaction) (uint, error)
	GetIntroductoryMeetingDocumentByID(id uint) (*models.RiskUpdatePlanIntroductoryMeetingDocument, error)
	GetIntroductoryMeetingDocumentsByPlanID(planID uint) ([]*models.RiskUpdatePlanIntroductoryMeetingDocument, error)
	DeleteIntroductoryMeetingDocument(id uint, tx ...utility.DatabaseTransaction) error

	// ERM Department Evaluation Documents
	CreateERMDeptDocument(document *models.RiskUpdatePlanErmDeptDocument, tx ...utility.DatabaseTransaction) (uint, error)
	GetERMDeptDocumentByID(id uint) (*models.RiskUpdatePlanErmDeptDocument, error)
	GetERMDeptDocumentsByTimelineID(timelineID uint) ([]*models.RiskUpdatePlanErmDeptDocument, error)
	DeleteERMDeptDocument(id uint, tx ...utility.DatabaseTransaction) error

	// Business Context Operations
	CreateBusinessContext(businessContext *models.RiskUpdatePlanBusinessContext, tx ...utility.DatabaseTransaction) error
	GetBusinessContextByPlanID(planID uint) ([]*models.RiskUpdatePlanBusinessContext, error)
	GetBusinessContextByID(id uint) (*models.RiskUpdatePlanBusinessContext, error)
	UpdateBusinessContextStatus(id uint, status bool, tx ...utility.DatabaseTransaction) error
	GetBusinessUnitIDByPlanID(planID uint) (uint, error)
	GetBusinessContextIDsByBusinessUnit(businessUnitID uint, tx ...utility.DatabaseTransaction) ([]uint, error)
	GetProcessIDsByBusinessUnit(businessUnitID uint, tx ...utility.DatabaseTransaction) ([]uint, error)
	GetControlIDsByBusinessUnit(businessUnitID uint, tx ...utility.DatabaseTransaction) ([]uint, error)

	// Timeline Operations
	GetAllTimelineTypes() ([]*models.RiskUpdatePlanTimelineType, error)
	CreateTimeline(timeline *models.RiskUpdatePlanTimeline, tx ...utility.DatabaseTransaction) error
	CalculateAndUpdatePercentage(planID uint, tx ...utility.DatabaseTransaction) error
	DeleteTimelinesByPlanID(planID uint, tx ...utility.DatabaseTransaction) error
}

type RiskRegisterUpdatePlanRepository struct {
	db *gorm.DB
}

func NewRiskRegisterUpdatePlanRepository(db *gorm.DB) *RiskRegisterUpdatePlanRepository {
	return &RiskRegisterUpdatePlanRepository{
		db: db,
	}
}

func (repo *RiskRegisterUpdatePlanRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *RiskRegisterUpdatePlanRepository) GetAll(params *riskDTOs.RiskRegisterUpdatePlanListQueryParamsDTO) (*riskDTOs.RiskRegisterUpdatePlanListResponseDTO, error) {
	var plans []*riskDTOs.RiskRegisterUpdatePlanListDTO
	var totalCount int64

	query := repo.db.Table("risk_register_update_plans").
		Select(`
			risk_register_update_plans.id,
			risk_register_update_plans.division_id,
			div.title_en AS division_title_en,
			div.title_ar AS division_title_ar,
			risk_register_update_plans.business_unit_id,
			bu.title_en AS business_unit_title_en,
			bu.title_ar AS business_unit_title_ar,
			risk_register_update_plans.risk_register_update_plan_status_id,
			status.title_en AS status_title_en,
			status.title_ar AS status_title_ar,
			status.type AS status_type,
			status.color_code AS status_color_code,
			risk_register_update_plans.number_of_operating_entities,
			risk_register_update_plans.workdays,
			risk_register_update_plans.risk_champion_id,
			champion.first_name_en AS risk_champion_first_name_en,
			champion.first_name_ar AS risk_champion_first_name_ar,
			champion.last_name_en AS risk_champion_last_name_en,
			champion.last_name_ar AS risk_champion_last_name_ar,
			champion.email AS risk_champion_email,
			CONVERT(VARCHAR, risk_register_update_plans.start_date, 23) AS start_date,
			CONVERT(VARCHAR, risk_register_update_plans.end_date, 23) AS end_date,
			risk_register_update_plans.notes,
			risk_register_update_plans.total_percentage_completed,
			risk_register_update_plans.created_at,
			risk_register_update_plans.updated_at,
			created_by.id AS created_user_id,
			created_by.first_name_en AS created_user_first_name_en,
			created_by.first_name_ar AS created_user_first_name_ar,
			created_by.last_name_en AS created_user_last_name_en,
			created_by.last_name_ar AS created_user_last_name_ar,
			created_by.initial AS created_user_initial,
			created_by.color AS created_user_color,
			created_by.email AS created_user_email,
			created_by.profile_image_path AS created_user_profile_image_path,
			updated_by.id AS updated_user_id,
			updated_by.first_name_en AS updated_user_first_name_en,
			updated_by.first_name_ar AS updated_user_first_name_ar,
			updated_by.last_name_en AS updated_user_last_name_en,
			updated_by.last_name_ar AS updated_user_last_name_ar,
			updated_by.initial AS updated_user_initial,
			updated_by.color AS updated_user_color,
			updated_by.email AS updated_user_email,
			updated_by.profile_image_path AS updated_user_profile_image_path
		`).
		Joins("LEFT JOIN organization_units div ON risk_register_update_plans.division_id = div.id").
		Joins("LEFT JOIN organization_units bu ON risk_register_update_plans.business_unit_id = bu.id").
		Joins("LEFT JOIN risk_register_update_plan_statuses status ON risk_register_update_plans.risk_register_update_plan_status_id = status.id").
		Joins("LEFT JOIN users champion ON risk_register_update_plans.risk_champion_id = champion.id").
		Joins("LEFT JOIN users created_by ON risk_register_update_plans.created_by_id = created_by.id").
		Joins("LEFT JOIN users updated_by ON risk_register_update_plans.updated_by_id = updated_by.id")

	// Apply filters
	for column, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where("risk_register_update_plans."+column+" IN (?)", values)
		}
	}

	// Apply keyword search
	if params.Keyword != "" {
		searchPattern := "%" + params.Keyword + "%"
		query = query.Where(`
			div.title_en LIKE ? OR 
			div.title_ar LIKE ? OR 
			bu.title_en LIKE ? OR 
			bu.title_ar LIKE ? OR
			champion.first_name_en LIKE ? OR
			champion.last_name_en LIKE ? OR
			risk_register_update_plans.notes LIKE ?
		`, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern)
	}

	// Count total records
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply sorting
	query = query.Order(params.SortBy + " " + params.Sort)

	// Apply pagination
	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	// Execute query
	if err := query.Find(&plans).Error; err != nil {
		return nil, err
	}

	totalPages := 1
	if !params.Export {
		totalPages = int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	}

	return &riskDTOs.RiskRegisterUpdatePlanListResponseDTO{
		Items: plans,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalCount:  totalCount,
			TotalPages:  totalPages,
		},
	}, nil
}

func (repo *RiskRegisterUpdatePlanRepository) GetByID(id uint) (*models.RiskRegisterUpdatePlan, error) {
	db := repo.db
	var plan models.RiskRegisterUpdatePlan
	if err := db.First(&plan, id).Error; err != nil {
		return nil, err
	}

	return &plan, nil
}

func (repo *RiskRegisterUpdatePlanRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskRegisterUpdatePlan, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var plan models.RiskRegisterUpdatePlan
	if err := db.Preload("Division").
		Preload("BusinessUnit").
		Preload("RiskRegisterUpdatePlanStatus").
		Preload("RiskChampion").
		Preload("TeamMembers.User").
		Preload("Timelines").
		Preload("Timelines.RiskUpdatePlanTimelineType").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("SubmittedByUser").
		Preload("ApprovedByUser").
		Preload("WorkflowStatus").
		First(&plan, id).Error; err != nil {
		return nil, err
	}

	return &plan, nil
}

func (repo *RiskRegisterUpdatePlanRepository) Create(plan *models.RiskRegisterUpdatePlan, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(plan).Error; err != nil {
		return 0, err
	}

	return plan.ID, nil
}

func (repo *RiskRegisterUpdatePlanRepository) Update(plan *models.RiskRegisterUpdatePlan, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Save(plan).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskRegisterUpdatePlanRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Delete(&models.RiskRegisterUpdatePlan{}, id).Error; err != nil {
		return err
	}

	return nil
}

// Team Members
func (repo *RiskRegisterUpdatePlanRepository) CreateTeamMembers(members []*models.RiskRegisterUpdatePlanTeamMember, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if len(members) == 0 {
		return nil
	}

	if err := db.Create(&members).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskRegisterUpdatePlanRepository) DeleteTeamMembersByPlanID(planID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("risk_register_update_plan_id = ?", planID).Delete(&models.RiskRegisterUpdatePlanTeamMember{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskRegisterUpdatePlanRepository) GetTeamMembersByPlanID(planID uint, tx ...utility.DatabaseTransaction) ([]*models.RiskRegisterUpdatePlanTeamMember, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var members []*models.RiskRegisterUpdatePlanTeamMember
	if err := db.Preload("User").Where("risk_register_update_plan_id = ?", planID).Find(&members).Error; err != nil {
		return nil, err
	}

	return members, nil
}

// ============================================
// ERM Department Evaluation - Risks
// ============================================

func (repo *RiskRegisterUpdatePlanRepository) GetRisksByBusinessUnitID(businessUnitID uint) ([]uint, error) {
	var riskIDs []uint
	err := repo.db.Table("risk_register").
		Select("id").
		Where("responsible_unit_id = ?", businessUnitID).
		Pluck("id", &riskIDs).Error

	return riskIDs, err
}

func (repo *RiskRegisterUpdatePlanRepository) GetERMDeptEvaluationRisks(timelineID uint, params *riskDTOs.ERMDeptEvaluationListQueryParamsDTO) (*riskDTOs.ERMDeptEvaluationListResponseDTO, error) {
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
			sec.title_ar AS section_title_ar
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
		Where("rupr.risk_update_plan_timeline_id = ?", timelineID)

	// Apply filters
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("(rr.title LIKE ? OR rr.reference_number LIKE ?)", keyword, keyword)
	}

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where("rr."+columnName+" IN (?)", values)
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

func (repo *RiskRegisterUpdatePlanRepository) CreateRiskMappings(risks []*models.RiskUpdatePlanRisk, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if len(risks) == 0 {
		return nil
	}

	return db.Create(&risks).Error
}

func (repo *RiskRegisterUpdatePlanRepository) DeleteAllRiskMappingsByTimelineID(timelineID uint, tx ...utility.DatabaseTransaction) error {
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

func (repo *RiskRegisterUpdatePlanRepository) GetProcessesByBusinessUnitID(businessUnitID uint) ([]uint, error) {
	var processIDs []uint
	err := repo.db.Table("processes").
		Select("id").
		Where("responsible_unit_id = ?", businessUnitID).
		Pluck("id", &processIDs).Error

	return processIDs, err
}

func (repo *RiskRegisterUpdatePlanRepository) GetERMDeptEvaluationProcesses(timelineID uint, params *riskDTOs.ERMDeptEvaluationListQueryParamsDTO) (*riskDTOs.ERMDeptEvaluationProcessListResponseDTO, error) {
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
			p.rto,
			p.rto_unit_id,
			rto_time_unit.title_en AS rto_unit_title_en,
			rto_time_unit.title_ar AS rto_unit_title_ar,
			p.mtpd,
			p.mtpd_unit_id,
			mtpd_time_unit.title_en AS mtpd_unit_title_en,
			mtpd_time_unit.title_ar AS mtpd_unit_title_ar,
			p.rpo,
			p.rpo_unit_id,
			rpo_time_unit.title_en AS rpo_unit_title_en,
			rpo_time_unit.title_ar AS rpo_unit_title_ar,
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
		Joins("LEFT JOIN time_units AS rto_time_unit ON rto_time_unit.id = p.rto_unit_id").
		Joins("LEFT JOIN time_units AS mtpd_time_unit ON mtpd_time_unit.id = p.mtpd_unit_id").
		Joins("LEFT JOIN time_units AS rpo_time_unit ON rpo_time_unit.id = p.rpo_unit_id").
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

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where("p."+columnName+" IN (?)", values)
		}
	}

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("p.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
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

func (repo *RiskRegisterUpdatePlanRepository) CreateProcessMappings(processes []*models.RiskUpdatePlanProcess, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if len(processes) == 0 {
		return nil
	}

	return db.Create(&processes).Error
}

func (repo *RiskRegisterUpdatePlanRepository) DeleteAllProcessMappingsByTimelineID(timelineID uint, tx ...utility.DatabaseTransaction) error {
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

func (repo *RiskRegisterUpdatePlanRepository) GetControlsByBusinessUnitID(businessUnitID uint) ([]uint, error) {
	var controlIDs []uint
	err := repo.db.Table("controls").
		Select("id").
		Where("responsible_unit_id = ?", businessUnitID).
		Pluck("id", &controlIDs).Error

	return controlIDs, err
}

func (repo *RiskRegisterUpdatePlanRepository) GetERMDeptEvaluationControls(timelineID uint, params *riskDTOs.ERMDeptEvaluationListQueryParamsDTO) (*riskDTOs.ERMDeptEvaluationControlListResponseDTO, error) {
	var controls []*riskDTOs.ERMDeptEvaluationControlListDTO
	var totalCount int64

	query := repo.db.Table("risk_update_plan_controls AS rupc").
		Select(`
			rupc.id,
			rupc.risk_update_plan_timeline_id,
			rupc.control_id,
			ct.title_en AS control_type_title_en,
			ct.title_ar AS control_type_title_ar,
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
			c.risk_reduction
		`).
		Joins("INNER JOIN controls AS c ON c.id = rupc.control_id").
		Joins("LEFT JOIN control_types AS ct ON ct.id = c.control_type_id").
		Joins("LEFT JOIN control_categories AS cc ON cc.id = c.control_category_id").
		Joins("LEFT JOIN organization_units AS ou ON ou.id = c.responsible_unit_id").
		Joins("LEFT JOIN users AS owner ON owner.id = c.responsible_user_id").
		Joins("LEFT JOIN control_automations AS ca ON ca.id = c.control_automation_id").
		Joins("LEFT JOIN control_consistencies AS ccon ON ccon.id = c.control_consistency_id").
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

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where("c."+columnName+" IN (?)", values)
		}
	}

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("c.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
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

func (repo *RiskRegisterUpdatePlanRepository) CreateControlMappings(controls []*models.RiskUpdatePlanControl, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if len(controls) == 0 {
		return nil
	}

	return db.Create(&controls).Error
}

func (repo *RiskRegisterUpdatePlanRepository) DeleteAllControlMappingsByTimelineID(timelineID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("risk_update_plan_timeline_id = ?", timelineID).
		Delete(&models.RiskUpdatePlanControl{}).Error
}

// Introductory Meeting Document Methods

func (repo *RiskRegisterUpdatePlanRepository) CreateIntroductoryMeetingDocument(document *models.RiskUpdatePlanIntroductoryMeetingDocument, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(document).Error; err != nil {
		return 0, err
	}

	return document.ID, nil
}

func (repo *RiskRegisterUpdatePlanRepository) GetIntroductoryMeetingDocumentByID(id uint) (*models.RiskUpdatePlanIntroductoryMeetingDocument, error) {
	var document models.RiskUpdatePlanIntroductoryMeetingDocument
	if err := repo.db.First(&document, id).Error; err != nil {
		return nil, err
	}

	return &document, nil
}

func (repo *RiskRegisterUpdatePlanRepository) GetIntroductoryMeetingDocumentsByPlanID(planID uint) ([]*models.RiskUpdatePlanIntroductoryMeetingDocument, error) {
	var documents []*models.RiskUpdatePlanIntroductoryMeetingDocument
	if err := repo.db.Where("risk_register_update_plan_id = ?", planID).Find(&documents).Error; err != nil {
		return nil, err
	}

	return documents, nil
}

func (repo *RiskRegisterUpdatePlanRepository) DeleteIntroductoryMeetingDocument(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Delete(&models.RiskUpdatePlanIntroductoryMeetingDocument{}, id).Error; err != nil {
		return err
	}

	return nil
}

// ERM Department Evaluation Documents

func (repo *RiskRegisterUpdatePlanRepository) CreateERMDeptDocument(document *models.RiskUpdatePlanErmDeptDocument, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(document).Error; err != nil {
		return 0, err
	}

	return document.ID, nil
}

func (repo *RiskRegisterUpdatePlanRepository) GetERMDeptDocumentByID(id uint) (*models.RiskUpdatePlanErmDeptDocument, error) {
	var document models.RiskUpdatePlanErmDeptDocument
	if err := repo.db.First(&document, id).Error; err != nil {
		return nil, err
	}

	return &document, nil
}

func (repo *RiskRegisterUpdatePlanRepository) GetERMDeptDocumentsByTimelineID(timelineID uint) ([]*models.RiskUpdatePlanErmDeptDocument, error) {
	var documents []*models.RiskUpdatePlanErmDeptDocument
	if err := repo.db.Where("risk_update_plan_timeline_id = ?", timelineID).Find(&documents).Error; err != nil {
		return nil, err
	}

	return documents, nil
}

func (repo *RiskRegisterUpdatePlanRepository) DeleteERMDeptDocument(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Delete(&models.RiskUpdatePlanErmDeptDocument{}, id).Error; err != nil {
		return err
	}

	return nil
}

// Business Context Methods

func (repo *RiskRegisterUpdatePlanRepository) CreateBusinessContext(businessContext *models.RiskUpdatePlanBusinessContext, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(businessContext).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskRegisterUpdatePlanRepository) GetBusinessContextByPlanID(planID uint) ([]*models.RiskUpdatePlanBusinessContext, error) {
	var businessContexts []*models.RiskUpdatePlanBusinessContext
	if err := repo.db.
		Preload("BusinessContext").
		Preload("RiskRegisterUpdatePlan").
		Preload("BusinessContext.Department").
		Preload("BusinessContext.CreatedByUser").
		Preload("BusinessContext.UpdatedByUser").
		Preload("Process").
		Preload("Process.ProcessGroup").
		Preload("Process.ResponsibleUnit").
		Preload("Process.ResponsibleUser").
		Preload("Process.LifecycleStatus").
		Preload("Process.Risks").
		Preload("Process.RiskMatrixRating").
		Preload("Process.RiskMatrixRatingValue").
		Preload("Process.RtoUnit").
		Preload("Process.MtpdUnit").
		Preload("Process.RpoUnit").
		Preload("Process.RiskLevel").
		Preload("Process.ProcessCriticalityLevel").
		Preload("Process.LastAssessedByUser").
		Preload("Process.ProcessActivities").
		Preload("Process.ProcessActivities.CreatedByUser").
		Preload("Control").
		Preload("Control.ControlType").
		Preload("Control.ControlCategory").
		Preload("Control.ResponsibleUnit").
		Preload("Control.ResponsibleUser").
		Preload("Control.ControlImplementationStatus").
		Preload("Control.ReviewFrequency").
		Preload("Control.CreatedByUser").
		Preload("Control.UpdatedByUser").
		Preload("Control.RiskMatrixRating").
		Preload("Control.RiskMatrixRatingValue").
		Preload("Control.ControlAutomation").
		Preload("Control.ControlConsistency").
		Preload("Control.ControlEffectivenessRating").
		Preload("CreatedByUser").
		Where("risk_register_update_plan_id = ?", planID).
		Find(&businessContexts).Error; err != nil {
		return nil, err
	}

	return businessContexts, nil
}

func (repo *RiskRegisterUpdatePlanRepository) GetBusinessContextByID(id uint) (*models.RiskUpdatePlanBusinessContext, error) {
	var businessContext models.RiskUpdatePlanBusinessContext
	if err := repo.db.
		Preload("BusinessContext").
		Preload("BusinessContext.Department").
		Preload("BusinessContext.CreatedByUser").
		Preload("BusinessContext.UpdatedByUser").
		Preload("Process").
		Preload("Process.ProcessGroup").
		Preload("Process.ResponsibleUnit").
		Preload("Process.ResponsibleUser").
		Preload("Process.ResponsibleUser.OrganizationUnit1").
		Preload("Process.LifecycleStatus").
		Preload("Process.Risks").
		Preload("Process.RiskMatrixRating").
		Preload("Process.RiskMatrixRatingValue").
		Preload("Process.CreatedByUser").
		Preload("Process.CreatedByUser.OrganizationUnit1").
		Preload("Process.UpdatedByUser").
		Preload("Process.UpdatedByUser.OrganizationUnit1").
		Preload("Process.RtoUnit").
		Preload("Process.MtpdUnit").
		Preload("Process.RpoUnit").
		Preload("Process.RiskLevel").
		Preload("Process.ProcessCriticalityLevel").
		Preload("Process.LastAssessedByUser").
		Preload("Process.ProcessActivities").
		Preload("Control").
		Preload("Control.ControlType").
		Preload("Control.ControlCategory").
		Preload("Control.ResponsibleUnit").
		Preload("Control.ResponsibleUser").
		Preload("Control.ResponsibleUser.OrganizationUnit1").
		Preload("Control.ControlImplementationStatus").
		Preload("Control.ControlEfficiencyStatus").
		Preload("Control.TestResponsibleUser").
		Preload("Control.TestResponsibleUser.OrganizationUnit1").
		Preload("Control.ReviewFrequency").
		Preload("Control.ControlReviewStatus").
		Preload("Control.TestFrequency").
		Preload("Control.LifecycleStatus").
		Preload("Control.ControlStatus").
		Preload("Control.CreatedByUser").
		Preload("Control.UpdatedByUser").
		Preload("Control.Risks").
		Preload("Control.RiskMatrixRating").
		Preload("Control.RiskMatrixRatingValue").
		Preload("Control.ControlAutomation").
		Preload("Control.ControlAutomation.CreatedByUser").
		Preload("Control.ControlConsistency").
		Preload("Control.ControlConsistency.CreatedByUser").
		Preload("Control.AssociatedDocuments").
		Preload("Control.AssociatedDocuments.Document").
		Preload("Control.ControlEffectivenessRating").
		Preload("CreatedByUser").
		First(&businessContext, id).Error; err != nil {
		return nil, err
	}

	return &businessContext, nil
}

func (repo *RiskRegisterUpdatePlanRepository) UpdateBusinessContextStatus(id uint, status bool, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Model(&models.RiskRegisterUpdatePlan{}).
		Where("id = ?", id).
		Update("is_business_context_updated", status).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskRegisterUpdatePlanRepository) GetBusinessUnitIDByPlanID(planID uint) (uint, error) {
	var businessUnitID uint
	if err := repo.db.Model(&models.RiskRegisterUpdatePlan{}).
		Where("id = ?", planID).
		Pluck("business_unit_id", &businessUnitID).Error; err != nil {
		return 0, err
	}

	return businessUnitID, nil
}

func (repo *RiskRegisterUpdatePlanRepository) GetBusinessContextIDsByBusinessUnit(businessUnitID uint, tx ...utility.DatabaseTransaction) ([]uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var ids []uint
	if err := db.Model(&models.OrganizationUnitBusinessContext{}).
		Where("department_id = ?", businessUnitID).
		Pluck("id", &ids).Error; err != nil {
		return nil, err
	}

	return ids, nil
}

func (repo *RiskRegisterUpdatePlanRepository) GetProcessIDsByBusinessUnit(businessUnitID uint, tx ...utility.DatabaseTransaction) ([]uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var ids []uint
	if err := db.Model(&models.Process{}).
		Where("responsible_unit_id = ?", businessUnitID).
		Pluck("id", &ids).Error; err != nil {
		return nil, err
	}

	return ids, nil
}

func (repo *RiskRegisterUpdatePlanRepository) GetControlIDsByBusinessUnit(businessUnitID uint, tx ...utility.DatabaseTransaction) ([]uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var ids []uint
	if err := db.Model(&models.Control{}).
		Where("responsible_unit_id = ?", businessUnitID).
		Pluck("id", &ids).Error; err != nil {
		return nil, err
	}

	return ids, nil
}

// Timeline Operations

func (repo *RiskRegisterUpdatePlanRepository) GetAllTimelineTypes() ([]*models.RiskUpdatePlanTimelineType, error) {
	var timelineTypes []*models.RiskUpdatePlanTimelineType
	if err := repo.db.Find(&timelineTypes).Error; err != nil {
		return nil, err
	}
	return timelineTypes, nil
}

func (repo *RiskRegisterUpdatePlanRepository) CreateTimeline(timeline *models.RiskUpdatePlanTimeline, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(timeline).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskRegisterUpdatePlanRepository) GetTimelineByPlanIDAndTypeID(planID uint, timelineTypeID uint, tx ...utility.DatabaseTransaction) (*models.RiskUpdatePlanTimeline, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var timeline models.RiskUpdatePlanTimeline
	if err := db.Where("risk_update_plan_id = ? AND risk_update_plan_timeline_type_id = ?", planID, timelineTypeID).
		First(&timeline).Error; err != nil {
		return nil, err
	}

	return &timeline, nil
}

func (repo *RiskRegisterUpdatePlanRepository) UpdateTimelineCompletion(timelineID uint, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	now := time.Now()
	updates := map[string]interface{}{
		"is_completed":  true,
		"updated_by_id": currentUserID,
		"updated_at":    now,
	}

	return db.Model(&models.RiskUpdatePlanTimeline{}).
		Where("id = ?", timelineID).
		Updates(updates).Error
}

func (repo *RiskRegisterUpdatePlanRepository) CalculateAndUpdatePercentage(planID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Count total timelines for this plan
	var totalCount int64
	if err := db.Model(&models.RiskUpdatePlanTimeline{}).
		Where("risk_update_plan_id = ?", planID).
		Count(&totalCount).Error; err != nil {
		return err
	}

	if totalCount == 0 {
		return nil
	}

	// Count completed timelines (is_completed = true OR is_not_applicable = true)
	var completedCount int64
	if err := db.Model(&models.RiskUpdatePlanTimeline{}).
		Where("risk_update_plan_id = ? AND is_completed = ?", planID, true).
		Count(&completedCount).Error; err != nil {
		return err
	}

	// Calculate percentage: (completedCount / totalCount) * 100
	percentage := int((float64(completedCount) / float64(totalCount)) * 100)

	// Update the plan's total_percentage_completed
	if err := db.Model(&models.RiskRegisterUpdatePlan{}).
		Where("id = ?", planID).
		Update("total_percentage_completed", percentage).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskRegisterUpdatePlanRepository) DeleteTimelinesByPlanID(planID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("risk_update_plan_id = ?", planID).Delete(&models.RiskUpdatePlanTimeline{}).Error; err != nil {
		return err
	}

	return nil
}
