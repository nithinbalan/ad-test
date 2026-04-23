package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/utility"

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

	// Risk Mapping (ERM Dept Evaluation)
	GetRisksByBusinessUnitID(businessUnitID uint) ([]uint, error)
	CreateERMDepartmentEvaluation(evaluation *models.RiskUpdatePlanErmDeptEvaluation, tx ...utility.DatabaseTransaction) error
	GetERMDepartmentEvaluation(planID uint, params *riskDTOs.ERMDeptEvaluationListQueryParamsDTO) (*riskDTOs.ERMDeptEvaluationListResponseDTO, error)
	GetExistingMappedRiskIDs(planID uint, tx ...utility.DatabaseTransaction) ([]uint, error)
	DeleteERMDepartmentEvaluation(id uint, tx ...utility.DatabaseTransaction) error
	DeleteERMDepartmentEvaluationsByRiskIDs(planID uint, riskIDs []uint, tx ...utility.DatabaseTransaction) error
	CheckRiskAlreadyMapped(planID uint, riskID uint) (bool, error)

	// Introductory Meeting Documents
	CreateIntroductoryMeetingDocument(document *models.RiskUpdatePlanIntroductoryMeetingDocument, tx ...utility.DatabaseTransaction) (uint, error)
	GetIntroductoryMeetingDocumentByID(id uint) (*models.RiskUpdatePlanIntroductoryMeetingDocument, error)
	GetIntroductoryMeetingDocumentsByPlanID(planID uint) ([]*models.RiskUpdatePlanIntroductoryMeetingDocument, error)
	DeleteIntroductoryMeetingDocument(id uint, tx ...utility.DatabaseTransaction) error

	// Business Context Operations
	CreateBusinessContext(businessContext *models.RiskUpdatePlanBusinessContext, tx ...utility.DatabaseTransaction) error
	GetBusinessContextByPlanID(planID uint) ([]*models.RiskUpdatePlanBusinessContext, error)
	GetBusinessContextByID(id uint) (*models.RiskUpdatePlanBusinessContext, error)
	UpdateBusinessContextStatus(id uint, status bool, tx ...utility.DatabaseTransaction) error
	GetBusinessUnitIDByPlanID(planID uint) (uint, error)
	GetBusinessContextIDsByBusinessUnit(businessUnitID uint, tx ...utility.DatabaseTransaction) ([]uint, error)
	GetProcessIDsByBusinessUnit(businessUnitID uint, tx ...utility.DatabaseTransaction) ([]uint, error)
	GetControlIDsByBusinessUnit(businessUnitID uint, tx ...utility.DatabaseTransaction) ([]uint, error)
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

// Risk Mapping (ERM Dept Evaluation) Methods

func (repo *RiskRegisterUpdatePlanRepository) GetRisksByBusinessUnitID(businessUnitID uint) ([]uint, error) {
	var riskIDs []uint

	err := repo.db.Table("risk_register").
		Select("id").
		Where("responsible_unit_id = ?", businessUnitID).
		Pluck("id", &riskIDs).Error

	if err != nil {
		return nil, err
	}

	return riskIDs, nil
}

func (repo *RiskRegisterUpdatePlanRepository) CreateERMDepartmentEvaluation(evaluation *models.RiskUpdatePlanErmDeptEvaluation, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(evaluation).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskRegisterUpdatePlanRepository) GetERMDepartmentEvaluation(planID uint, params *riskDTOs.ERMDeptEvaluationListQueryParamsDTO) (*riskDTOs.ERMDeptEvaluationListResponseDTO, error) {
	var mappedRisks []*riskDTOs.ERMDeptEvaluationRiskListDTO
	var totalCount int64

	query := repo.db.Table("risk_update_plan_erm_dept_evaluations").
		Select(`
			risk_update_plan_erm_dept_evaluations.id,
			risk_update_plan_erm_dept_evaluations.risk_register_update_plan_id,
			risk_update_plan_erm_dept_evaluations.risk_register_id,
			risk_register.reference_number,
			risk_register.title,
			risk_register.description,
			CONVERT(VARCHAR, risk_register.reported_at, 23) AS reported_at,
			CONVERT(VARCHAR, risk_update_plan_erm_dept_evaluations.created_at, 120) AS created_at,
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
			sections.title_ar AS section_title_ar
		`).
		Joins("LEFT JOIN risk_register ON risk_register.id = risk_update_plan_erm_dept_evaluations.risk_register_id").
		Joins("LEFT JOIN risk_categories ON risk_categories.id = risk_register.risk_category_id").
		Joins("LEFT JOIN risk_sources ON risk_sources.id = risk_register.risk_source_id").
		Joins("LEFT JOIN organization_units ON organization_units.id = risk_register.responsible_unit_id").
		Joins("LEFT JOIN users ON users.id = risk_register.responsible_user_id").
		Joins("LEFT JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id").
		Joins("LEFT JOIN risk_response_strategies ON risk_response_strategies.id = risk_register.risk_response_strategy_id").
		Joins("LEFT JOIN users AS created_user ON created_user.id = risk_register.created_by_id").
		Joins("LEFT JOIN users AS updated_user ON updated_user.id = risk_register.updated_by_id").
		Joins("LEFT JOIN review_frequencies ON review_frequencies.id = risk_register.review_frequency_id").
		Joins("LEFT JOIN severity_levels ON severity_levels.id = risk_register.severity_level_id").
		Joins("LEFT JOIN risk_matrix_likelihoods AS current_likelihood ON current_likelihood.id = risk_register.current_likelihood_id").
		Joins("LEFT JOIN risk_matrix_impacts AS current_impact ON current_impact.id = risk_register.current_impact_id").
		Joins("LEFT JOIN risk_matrix_risk_ratings AS current_risk_rating ON current_risk_rating.id = risk_register.current_risk_rating_id").
		Joins("LEFT JOIN risk_matrix_risk_rating_values AS current_risk_rating_value ON current_risk_rating_value.id = risk_register.current_risk_rating_value_id").
		Joins("LEFT JOIN risk_matrix_likelihoods AS inherent_likelihood ON inherent_likelihood.id = risk_register.inherent_likelihood_id").
		Joins("LEFT JOIN risk_matrix_impacts AS inherent_impact ON inherent_impact.id = risk_register.inherent_impact_id").
		Joins("LEFT JOIN risk_matrix_risk_ratings AS inherent_risk_rating ON inherent_risk_rating.id = risk_register.inherent_risk_rating_id").
		Joins("LEFT JOIN risk_matrix_risk_rating_values AS inherent_risk_rating_value ON inherent_risk_rating_value.id = risk_register.inherent_risk_rating_value_id").
		Joins("LEFT JOIN risk_matrix_likelihoods AS residual_likelihood ON residual_likelihood.id = risk_register.residual_likelihood_id").
		Joins("LEFT JOIN risk_matrix_impacts AS residual_impact ON residual_impact.id = risk_register.residual_impact_id").
		Joins("LEFT JOIN risk_matrix_risk_ratings AS residual_risk_rating ON residual_risk_rating.id = risk_register.residual_risk_rating_id").
		Joins("LEFT JOIN risk_matrix_risk_rating_values AS residual_risk_rating_value ON residual_risk_rating_value.id = risk_register.residual_risk_rating_value_id").
		Joins("LEFT JOIN organization_units AS divisions ON divisions.id = risk_register.division_id").
		Joins("LEFT JOIN organization_units AS sections ON sections.id = risk_register.section_id").
		Where("risk_update_plan_erm_dept_evaluations.risk_register_update_plan_id = ?", planID)

	if params.RiskAppetiteToleranceLevelID != 0 {
		query = query.Joins(`
			LEFT JOIN risk_evaluations ON risk_evaluations.risk_register_id = risk_register.id
			LEFT JOIN risk_appetite_alignments ON risk_appetite_alignments.risk_evaluation_id = risk_evaluations.id
		`).Where("risk_appetite_alignments.risk_appetite_tolerance_level_id = ?", params.RiskAppetiteToleranceLevelID)
	}

	// Apply dynamic filters
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	if params.ReportedAtFrom != "" {
		query = query.Where("risk_register.reported_at >= ?", params.ReportedAtFrom)
	}

	if params.ReportedAtTo != "" {
		query = query.Where("risk_register.reported_at <= ?", params.ReportedAtTo)
	}

	if params.NextReviewDateFrom != "" {
		query = query.Where("risk_register.next_review_date >= ?", params.NextReviewDateFrom)
	}

	if params.NextReviewDateTo != "" {
		query = query.Where("risk_register.next_review_date <= ?", params.NextReviewDateTo)
	}

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("risk_register.reported_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	// Apply quick filters
	for _, filter := range params.QuickFilter {
		switch filter {
		case "open":
			query = query.Where("risk_statuses.type = ?", "open")
		case "critical_risk":
			var highestRating string
			err := repo.db.Table("risk_matrix_risk_ratings").
				Select("type").
				Where("status = ?", "active").
				Order("CASE type WHEN 'very_high' THEN 1 WHEN 'high' THEN 2 WHEN 'medium' THEN 3 WHEN 'low' THEN 4 WHEN 'very_low' THEN 5 END").
				Limit(1).Pluck("type", &highestRating).Error
			if err == nil && highestRating != "" {
				query = query.Where("current_risk_rating.type = ?", highestRating)
			}
		case "strategic_risks":
			query = query.Where("risk_register.is_strategic = ?", true)
		case "operational_risks":
			query = query.Where("risk_register.is_operational = ?", true)
		}
	}

	// Apply search keyword
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(risk_register.title LIKE ? OR
			risk_register.description LIKE ? OR
			risk_register.reference_number LIKE ? OR
			risk_categories.title_en LIKE ? OR
			risk_categories.title_ar LIKE ? OR
			risk_sources.title_en LIKE ? OR
			risk_sources.title_ar LIKE ? OR
			users.first_name_en LIKE ? OR
			users.first_name_ar LIKE ? OR
			users.last_name_en LIKE ? OR
			users.last_name_ar LIKE ? OR
			risk_statuses.title_en LIKE ? OR
			risk_statuses.title_ar LIKE ? OR
			organization_units.title_en LIKE ? OR
			organization_units.title_ar LIKE ? OR
			current_risk_rating.title_en LIKE ? OR
			current_risk_rating.title_ar LIKE ? OR
			current_impact.title_en LIKE ? OR
			current_impact.title_ar LIKE ? OR
			current_likelihood.title_en LIKE ? OR
			current_likelihood.title_ar LIKE ? OR
			risk_register.reported_at LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.last_name_ar LIKE ? OR
			updated_user.first_name_en LIKE ? OR
			updated_user.first_name_ar LIKE ? OR
			updated_user.last_name_en LIKE ? OR
			updated_user.last_name_ar LIKE ? OR
			review_frequencies.title_en LIKE ? OR
			review_frequencies.title_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	if params.Value != 0 {
		query = query.Where("current_risk_rating_value.value = ?", params.Value)
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	// Count total results
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply pagination only if not exporting
	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	if err := query.Find(&mappedRisks).Error; err != nil {
		return nil, err
	}

	if params.Export {
		return &riskDTOs.ERMDeptEvaluationListResponseDTO{
			Items: mappedRisks,
		}, nil
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &riskDTOs.ERMDeptEvaluationListResponseDTO{
		Items: mappedRisks,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}

	return response, nil
}

func (repo *RiskRegisterUpdatePlanRepository) GetExistingMappedRiskIDs(planID uint, tx ...utility.DatabaseTransaction) ([]uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var riskIDs []uint
	err := db.Model(&models.RiskUpdatePlanErmDeptEvaluation{}).
		Where("risk_register_update_plan_id = ?", planID).
		Pluck("risk_register_id", &riskIDs).Error

	if err != nil {
		return nil, err
	}

	return riskIDs, nil
}

func (repo *RiskRegisterUpdatePlanRepository) DeleteERMDepartmentEvaluation(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Delete(&models.RiskUpdatePlanErmDeptEvaluation{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskRegisterUpdatePlanRepository) DeleteERMDepartmentEvaluationsByRiskIDs(planID uint, riskIDs []uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if len(riskIDs) == 0 {
		return nil
	}

	err := db.Where("risk_register_update_plan_id = ? AND risk_register_id IN ?", planID, riskIDs).
		Delete(&models.RiskUpdatePlanErmDeptEvaluation{}).Error

	return err
}

func (repo *RiskRegisterUpdatePlanRepository) CheckRiskAlreadyMapped(planID uint, riskID uint) (bool, error) {
	var count int64

	err := repo.db.Model(&models.RiskUpdatePlanErmDeptEvaluation{}).
		Where("risk_register_update_plan_id = ? AND risk_register_id = ?", planID, riskID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
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
		Preload("BusinessContext.Department").
		Preload("Process").
		Preload("Process.ResponsibleUnit").
		Preload("Control").
		Preload("Control.ResponsibleUnit").
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
		Preload("Process").
		Preload("Process.ResponsibleUnit").
		Preload("Control").
		Preload("Control.ResponsibleUnit").
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

	if err := db.Model(&models.RiskUpdatePlanBusinessContext{}).
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
