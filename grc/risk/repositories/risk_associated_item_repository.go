package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/utility"
	"sort"

	"gorm.io/gorm"
)

type RiskAssociatedItemRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAllControls(riskID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedControlListResponseDTO, error)
	GetAllComplianceRequirements(riskID uint, params *riskDTOs.RiskAssociatedItemListQueryParamsDTO) ([]riskDTOs.RiskComplianceRequirementListDTO, int64, error)
	GetAllProcesses(riskID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemProcessListResponseDTO, error)
	GetAllStrategicObjectives(riskID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemStrategicObjectiveListResponseDTO, error)
	GetAllSLAs(riskID uint, params *riskDTOs.RiskAssociatedItemListQueryParamsDTO) (*riskDTOs.RiskAssociatedItemSLAListResponseDTO, error)
	GetAllIncidents(riskID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedIncidentListResponseDTO, error)
	GetAllDocuments(riskID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error)
	GetAllAssets(riskID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemAssetListResponseDTO, error)
	GetAllThreats(riskID uint, params *riskDTOs.RiskAssociatedItemListQueryParamsDTO) (*riskDTOs.RiskAssociatedItemThreatListResponseDTO, error)
	GetAllClauses(riskID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemClauseListResponseDTO, error)
	GetAllRiskKPI(riskID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedKPIRegisterListResponseDTO, error)
	CreateRiskClauses(clauses []*models.RiskComplianceRequirement, tx ...utility.DatabaseTransaction) error
	DeleteRiskClauses(riskID uint, clauseIDs []uint, tx ...utility.DatabaseTransaction) error
	GetClausesByRiskID(riskID uint, tx ...utility.DatabaseTransaction) ([]uint, error)
	CreatRiskStrategicObjective(strategicObjectives []*models.RiskStrategicObjective, tx ...utility.DatabaseTransaction) error
	DeletetRiskStrategicObjective(strategicObjectiveID uint, riskID uint, tx ...utility.DatabaseTransaction) error
	CreatRiskKPI(kpiS []*models.KPIRisk, tx ...utility.DatabaseTransaction) error
	DeletetRiskKPI(kpiID uint, riskID uint, tx ...utility.DatabaseTransaction) error
	CreatRiskProcess(process []*models.RiskProcess, tx ...utility.DatabaseTransaction) error
	DeletetRiskProcess(riskID uint, processID uint, tx ...utility.DatabaseTransaction) error
	CreatRiskDocument(documents []*models.RiskAssociatedDocument, tx ...utility.DatabaseTransaction) error
	DeletetRiskDocument(riskID uint, documentID uint, tx ...utility.DatabaseTransaction) error
	CreatRiskIncident(incidents []*models.IncidentRisk, tx ...utility.DatabaseTransaction) error
	DeletetRiskIncident(riskID uint, incidentID uint, tx ...utility.DatabaseTransaction) error
	CreatRiskControl(incidents []*models.RiskControl, tx ...utility.DatabaseTransaction) error
	DeletetRiskControl(riskID uint, controlID uint, tx ...utility.DatabaseTransaction) error
	CreatRiskAsset(incidents []*models.RiskAsset, tx ...utility.DatabaseTransaction) error
	DeletetRiskAsset(riskID uint, assetID uint, tx ...utility.DatabaseTransaction) error

	GetAllRiskBCP(riskID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedBusinessContinuityPlanListResponseDTO, error)
	CreatRiskBCP(bcp []*models.BcpRisk, tx ...utility.DatabaseTransaction) error
	DeletetRiskBCP(riskID uint, bcpID uint, tx ...utility.DatabaseTransaction) error

	GetAllRiskProject(riskID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedProjectListResponseDTO, error)
	CreatRiskProject(projects []*models.ProjectRisk, tx ...utility.DatabaseTransaction) error
	DeletetRiskProject(riskID uint, projectID uint, tx ...utility.DatabaseTransaction) error
}

type RiskAssociatedItemRepository struct {
	db *gorm.DB
}

func NewRiskAssociatedItemRepository(db *gorm.DB) *RiskAssociatedItemRepository {
	return &RiskAssociatedItemRepository{
		db: db,
	}
}

func (repo *RiskAssociatedItemRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}
func (repo *RiskAssociatedItemRepository) GetAllControls(riskID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedControlListResponseDTO, error) {
	var controlList []*commonDTOs.AssociatedControlListDTO
	var totalCount int64

	query := repo.db.Table("risk_controls").
		Select(`
		DISTINCT controls.id AS control_id,
		controls.*,
		CONVERT(VARCHAR, controls.last_test_date, 23) AS control_last_test_date,
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
		review_frequencies.title_en AS review_frequency_title_en,
		review_frequencies.title_ar AS review_frequency_title_ar,
		test_responsible_user.id AS test_responsible_user_id,
		test_responsible_user.first_name_en AS test_responsible_user_first_name_en,
		test_responsible_user.first_name_ar AS test_responsible_user_first_name_ar,
		test_responsible_user.last_name_en AS test_responsible_user_last_name_en,
		test_responsible_user.last_name_ar AS test_responsible_user_last_name_ar,
		test_responsible_user.initial AS test_responsible_user_initial,
		test_responsible_user.color AS test_responsible_user_color,
		test_responsible_user.email AS test_responsible_user_email,
		test_responsible_user.profile_image_path AS test_responsible_user_profile_image_path,
		test_frequencies.title_en AS test_frequency_title_en,
		test_frequencies.title_ar AS test_frequency_title_ar,
		lifecycle_statuses.title_en AS lifecycle_status_title_en,
		lifecycle_statuses.title_ar AS lifecycle_status_title_ar,
		lifecycle_statuses.color_code AS lifecycle_status_color_code,
		control_implementation_statuses.title_en AS control_implementation_status_title_en,
		control_implementation_statuses.title_ar AS control_implementation_status_title_ar,
		control_implementation_statuses.color_code AS control_implementation_status_color_code,
		control_efficiency_statuses.title_en AS control_efficiency_status_title_en,
		control_efficiency_statuses.title_ar AS control_efficiency_status_title_ar,
		control_efficiency_statuses.color_code AS control_efficiency_status_color_code,
		risk_matrix_risk_ratings.title_en AS risk_rating_title_en,
		risk_matrix_risk_ratings.title_ar AS risk_rating_title_ar,
		risk_matrix_risk_ratings.color_code AS risk_rating_color_code,
		risk_matrix_risk_rating_values.value AS risk_rating_value,
		created_user.id AS created_user_id,
		created_user.first_name_en AS created_user_first_name_en,
		created_user.first_name_ar AS created_user_first_name_ar,
		created_user.last_name_en AS created_user_last_name_en,
		created_user.last_name_ar AS created_user_last_name_ar,
		created_user.initial AS created_user_initial,
		created_user.color AS created_user_color,
		created_user.email AS created_user_email,
		created_user.profile_image_path AS created_user_profile_image_path
	`).
		Joins(`
		LEFT JOIN controls ON controls.id = risk_controls.control_id
		LEFT JOIN control_types ON controls.control_type_id = control_types.id
		LEFT JOIN control_categories ON controls.control_category_id = control_categories.id
		LEFT JOIN organization_units ON controls.responsible_unit_id = organization_units.id
		LEFT JOIN users AS owner ON controls.responsible_user_id = owner.id
		LEFT JOIN review_frequencies ON controls.review_frequency_id = review_frequencies.id
		LEFT JOIN users AS test_responsible_user ON controls.test_responsible_user_id = test_responsible_user.id
		LEFT JOIN review_frequencies AS test_frequencies ON controls.test_frequency_id = test_frequencies.id
		LEFT JOIN lifecycle_statuses ON controls.lifecycle_status_id = lifecycle_statuses.id
		LEFT JOIN control_implementation_statuses ON controls.control_implementation_status_id = control_implementation_statuses.id
		LEFT JOIN control_efficiency_statuses ON controls.control_efficiency_status_id = control_efficiency_statuses.id
		LEFT JOIN risk_matrix_risk_ratings ON risk_matrix_risk_ratings.id = controls.risk_matrix_risk_rating_id
		LEFT JOIN risk_matrix_risk_rating_values ON risk_matrix_risk_rating_values.id = controls.risk_matrix_risk_rating_value_id
		LEFT JOIN users AS created_user ON created_user.id = controls.created_by_id
	`).Where("risk_controls.risk_register_id = ?", riskID)

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(controls.title_en LIKE ? OR
			controls.title_ar LIKE ? OR
			controls.reference_number LIKE ? OR
			controls.code LIKE ? OR
			control_categories.title_en LIKE ? OR
			control_categories.title_ar LIKE ? OR
			control_types.title_en LIKE ? OR
			control_types.title_ar LIKE ? OR
			control_efficiency_statuses.title_en LIKE ? OR
			control_efficiency_statuses.title_ar LIKE ? OR
			control_implementation_statuses.title_en LIKE ? OR
			control_implementation_statuses.title_ar LIKE ? OR
			owner.first_name_en LIKE ? OR
			owner.first_name_ar LIKE ? OR
			owner.last_name_en LIKE ? OR
			owner.last_name_ar LIKE ? OR
			owner.email LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	if params.SortBy == "" {
		params.SortBy = "controls.title_en"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&controlList).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &commonDTOs.AssociatedControlListResponseDTO{
		Items: controlList,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil

}

func (repo *RiskAssociatedItemRepository) GetAllComplianceRequirements(riskID uint, params *riskDTOs.RiskAssociatedItemListQueryParamsDTO) ([]riskDTOs.RiskComplianceRequirementListDTO, int64, error) {
	var complianceRequirementList []riskDTOs.RiskComplianceRequirementListDTO
	var totalCount int64

	query := repo.db.Table("risk_compliance_requirements").
		Select(`
			risk_compliance_requirements.compliance_register_id AS compliance_register_id,
			compliance_register.title_en AS compliance_register_title_en,
			compliance_register.title_ar AS compliance_register_title_ar,
			risk_compliance_requirements.compliance_requirement_id AS compliance_requirement_id,
			compliance_clauses.title_en AS compliance_clause_title_en,
			compliance_clauses.title_ar AS compliance_clause_title_ar,
			compliance_clauses.clause_no AS compliance_clause_no,
			compliance_requirement_categories.title_en AS compliance_requirement_category_title_en,
			compliance_requirement_categories.title_ar AS compliance_requirement_category_title_ar,
			compliance_requirements.title_en AS compliance_requirement_title_en,
			compliance_requirements.title_ar AS compliance_requirement_title_ar,
			compliance_statuses.title_en AS compliance_status_title_en,
			compliance_statuses.title_ar AS compliance_status_title_ar,
			compliance_statuses.color_code AS compliance_status_color_code
		`).
		Joins(`
			LEFT JOIN compliance_requirements ON compliance_requirements.id = risk_compliance_requirements.compliance_requirement_id
			LEFT JOIN compliance_clauses ON compliance_requirements.compliance_clause_id = compliance_clauses.id
			LEFT JOIN compliance_statuses ON compliance_requirements.compliance_status_id = compliance_statuses.id
			LEFT JOIN compliance_requirement_categories ON compliance_requirements.compliance_requirement_category_id = compliance_requirement_categories.id
			LEFT JOIN compliance_register ON compliance_register.id = risk_compliance_requirements.compliance_register_id
		`).
		Where("risk_compliance_requirements.risk_register_id = ?", riskID)

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`(compliance_requirements.title_en LIKE ?)`, keyword)
	}

	if params.SortBy == "" {
		params.SortBy = "compliance_requirements.title_en"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&complianceRequirementList).Error
	if err != nil {
		return nil, 0, err
	}

	return complianceRequirementList, totalCount, nil
}

func (repo *RiskAssociatedItemRepository) GetAllProcesses(riskID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemProcessListResponseDTO, error) {
	var processList []*commonDTOs.AssociatedProcessListDTO
	var totalCount int64

	query := repo.db.Table("risk_processes").
		Select(`
		risk_processes.id,
		risk_processes.process_id,
		processes.reference_number,
		processes.title_en,
		processes.title_ar,
		processes.description,
		processes.marked_as_auditable_item,
		process_groups.title_en AS process_group_title_en,
		process_groups.title_ar AS process_group_title_ar,
		processes.rto as rto,
		processes.rto_unit_id as rto_unit_id,
		rto_time_unit.title_en AS rto_unit_title_en,
		rto_time_unit.title_ar AS rto_unit_title_ar,
		processes.mtpd as mtpd,
		processes.mtpd_unit_id as mtpd_unit_id,
		mtpd_time_unit.title_en AS mtpd_unit_title_en,
		mtpd_time_unit.title_ar AS mtpd_unit_title_ar,
		processes.rpo as rpo,
		processes.rpo_unit_id as rpo_unit_id,
		rpo_time_unit.title_en AS rpo_unit_title_en,
		rpo_time_unit.title_ar AS rpo_unit_title_ar,
		processes.is_bia_performed,
		processes.last_assessed_at,
		last_assessed_user.id AS last_assessed_user_id,
		last_assessed_user.first_name_en AS last_assessed_user_first_name_en,
		last_assessed_user.first_name_ar AS last_assessed_user_first_name_ar,
		last_assessed_user.last_name_en AS last_assessed_user_last_name_en,
		last_assessed_user.last_name_ar AS last_assessed_user_last_name_ar,
		last_assessed_user.initial AS last_assessed_user_initial,
		last_assessed_user.color AS last_assessed_user_color,
		last_assessed_user.email AS last_assessed_user_email,
		last_assessed_user.profile_image_path AS last_assessed_user_profile_image_path,
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
		lifecycle_statuses.title_en AS lifecycle_status_title_en,
		lifecycle_statuses.title_ar AS lifecycle_status_title_ar,
		lifecycle_statuses.color_code AS lifecycle_status_color_code,
		process_criticality_levels.id AS process_criticality_level_id,
		process_criticality_levels.title_en AS process_criticality_level_title_en,
		process_criticality_levels.title_ar AS process_criticality_level_title_ar,
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
		processes.created_at,
		processes.updated_at,
		risk_matrix_risk_ratings.title_en AS risk_rating_title_en,
		risk_matrix_risk_ratings.title_ar AS risk_rating_title_ar,
		risk_matrix_risk_ratings.color_code AS risk_rating_color_code,
		risk_matrix_risk_rating_values.value AS risk_rating_value,
		risk_levels.title_en AS risk_level_title_en,
        risk_levels.title_ar AS risk_level_title_ar,
		risk_levels.color_code AS risk_level_color_code,
		risk_levels.type AS risk_level_type
		`).
		Joins(`
			LEFT JOIN processes ON processes.id = risk_processes.process_id
			LEFT JOIN process_groups ON process_groups.id = processes.process_group_id
			LEFT JOIN organization_units ON organization_units.id = processes.responsible_unit_id
			LEFT JOIN users AS owner ON owner.id = processes.responsible_user_id
			LEFT JOIN lifecycle_statuses ON lifecycle_statuses.id = processes.lifecycle_status_id
			LEFT JOIN users AS created_user ON created_user.id = processes.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = processes.updated_by_id
			LEFT JOIN users AS last_assessed_user ON last_assessed_user.id = processes.last_assessed_by_id
			LEFT JOIN time_units AS rto_time_unit ON rto_time_unit.id = processes.rto_unit_id
			LEFT JOIN time_units AS mtpd_time_unit ON mtpd_time_unit.id = processes.mtpd_unit_id
			LEFT JOIN time_units AS rpo_time_unit ON rpo_time_unit.id = processes.rpo_unit_id
			LEFT JOIN risk_matrix_risk_ratings ON risk_matrix_risk_ratings.id = processes.risk_matrix_risk_rating_id
			LEFT JOIN risk_matrix_risk_rating_values ON risk_matrix_risk_rating_values.id = processes.risk_matrix_risk_rating_value_id
			LEFT JOIN risk_levels ON risk_levels.id = processes.risk_level_id
			LEFT JOIN process_criticality_levels ON process_criticality_levels.id = processes.process_criticality_level_id
		`).Where("risk_processes.risk_register_id = ?", riskID)

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(processes.title_en LIKE ? OR
			processes.title_ar LIKE ? OR
			processes.reference_number LIKE ? OR
			process_groups.title_en LIKE ? OR
			process_groups.title_ar LIKE ? OR
			organization_units.title_en LIKE ? OR
			organization_units.title_ar LIKE ? OR
			owner.first_name_en LIKE ? OR
			owner.first_name_ar LIKE ? OR
			owner.last_name_en LIKE ? OR
			owner.last_name_ar LIKE ? OR
			lifecycle_statuses.title_en LIKE ? OR
			lifecycle_statuses.title_ar LIKE ? OR
			lifecycle_statuses.color_code LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword)
	}

	if params.SortBy == "" {
		params.SortBy = "processes.title_en"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&processList).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &commonDTOs.AssociatedItemProcessListResponseDTO{
		Items: processList,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *RiskAssociatedItemRepository) GetAllStrategicObjectives(riskID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemStrategicObjectiveListResponseDTO, error) {
	var objectiveList []*commonDTOs.AssociatedStrategicObjectiveListDTO
	var totalCount int64

	query := repo.db.Table("risk_strategic_objectives").
		Select(`
			strategic_objectives.*,
			risk_strategic_objectives.strategic_objective_id,
			strategic_objective_categories.title_en AS strategic_objective_category_title_en,
			strategic_objective_categories.title_ar AS strategic_objective_category_title_ar,
			organization_units.title_en AS responsible_unit_title_en,
			organization_units.title_ar AS responsible_unit_title_ar,
			responsible_user.first_name_en AS responsible_user_first_name_en,
			responsible_user.last_name_en AS responsible_user_last_name_en,
			responsible_user.first_name_ar AS responsible_user_first_name_ar,
			responsible_user.last_name_ar AS responsible_user_last_name_ar,
			responsible_user.initial AS responsible_user_initial,
			responsible_user.email AS responsible_user_email,
			responsible_user.color AS responsible_user_color,
			responsible_user.profile_image_path AS responsible_user_profile_image_path,
			review_user.first_name_en AS review_user_first_name_en,
			review_user.first_name_ar AS review_user_first_name_ar,
			review_user.last_name_en AS review_user_last_name_en,
			review_user.last_name_ar AS review_user_last_name_ar,
			review_user.initial AS review_user_initial,
			review_user.color AS review_user_color,
			review_user.email AS review_user_email,
			review_user.profile_image_path AS review_user_profile_image_path,

			review_frequencies.title_en AS review_frequency_title_en,
			review_frequencies.title_ar AS review_frequency_title_ar,

			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.email AS created_user_email,
			created_user.profile_image_path AS created_user_profile_image_path,

			updated_user.first_name_en AS updated_user_first_name_en,
			updated_user.first_name_ar AS updated_user_first_name_ar,
			updated_user.last_name_en AS updated_user_last_name_en,
			updated_user.last_name_ar AS updated_user_last_name_ar,
			updated_user.initial AS updated_user_initial,
			updated_user.color AS updated_user_color,
			updated_user.email AS updated_user_email,
			updated_user.profile_image_path AS updated_user_profile_image_path
        `).
		Joins(`
			LEFT JOIN strategic_objectives ON strategic_objectives.id = risk_strategic_objectives.strategic_objective_id
			LEFT JOIN strategic_objective_categories ON strategic_objective_categories.id = strategic_objectives.strategic_objective_category_id
			LEFT JOIN organization_units ON organization_units.id = strategic_objectives.responsible_unit_id
			LEFT JOIN users AS responsible_user ON responsible_user.id = strategic_objectives.responsible_user_id
			LEFT JOIN users AS review_user ON review_user.id = strategic_objectives.review_user_id
			LEFT JOIN review_frequencies ON review_frequencies.id = strategic_objectives.review_frequency_id
			LEFT JOIN users AS created_user ON created_user.id = strategic_objectives.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = strategic_objectives.updated_by_id
        `).Where("risk_strategic_objectives.risk_register_id = ?", riskID)

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(strategic_objectives.title_en LIKE ? OR
			strategic_objectives.title_ar LIKE ? OR
			strategic_objectives.status LIKE ? )`,
			keyword, keyword, keyword)
	}

	if params.SortBy == "" {
		params.SortBy = "strategic_objectives.title_en"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&objectiveList).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &commonDTOs.AssociatedItemStrategicObjectiveListResponseDTO{
		Items: objectiveList,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}

	return response, nil
}

func (repo *RiskAssociatedItemRepository) GetAllSLAs(riskID uint, params *riskDTOs.RiskAssociatedItemListQueryParamsDTO) (*riskDTOs.RiskAssociatedItemSLAListResponseDTO, error) {
	var slaList []*riskDTOs.RiskSLAListDTO
	var totalCount int64

	query := repo.db.Table("risk_slas").
		Select(`
			sla_register.reference_number,
			sla_register.title_en,
			sla_register.title_ar,
			sla_register.description,
			sla_register.client,
			sla_categories.title_en AS sla_category_title_en,
            sla_categories.title_ar AS sla_category_title_ar,
            severity_levels.title_en AS severity_level_title_en,
            severity_levels.title_ar AS severity_level_title_ar,
            organization_units.title_en AS responsible_unit_title_en,
            organization_units.title_ar AS responsible_unit_title_ar,
			responsible_users.id AS responsible_user_id,
            responsible_users.first_name_en AS responsible_user_first_name_en,
            responsible_users.first_name_ar AS responsible_user_first_name_ar,
            responsible_users.last_name_en AS responsible_user_last_name_en,
            responsible_users.last_name_ar AS responsible_user_last_name_ar,
			responsible_users.initial AS responsible_user_initial,
            responsible_users.color AS responsible_user_color,
            responsible_users.email AS responsible_user_email,
            responsible_users.profile_image_path AS responsible_user_profile_image_path,
            lifecycle_statuses.title_en AS lifecycle_status_title_en,
            lifecycle_statuses.title_ar AS lifecycle_status_title_ar,
			lifecycle_statuses.color_code AS lifecycle_status_color_code,
			compliance_statuses.title_en AS compliance_status_title_en,
			compliance_statuses.title_ar AS compliance_status_title_ar,
			compliance_statuses.color_code AS compliance_status_color_code
		`).
		Joins(`
			LEFT JOIN sla_register ON sla_register.id = risk_slas.sla_register_id
			LEFT JOIN sla_categories ON sla_register.sla_category_id = sla_categories.id
            LEFT JOIN severity_levels ON sla_register.severity_level_id = severity_levels.id
            LEFT JOIN organization_units ON sla_register.responsible_unit_id = organization_units.id
            LEFT JOIN users AS responsible_users ON sla_register.responsible_user_id = responsible_users.id
			LEFT JOIN lifecycle_statuses ON sla_register.lifecycle_status_id = lifecycle_statuses.id
		    LEFT JOIN compliance_statuses ON sla_register.compliance_status_id = compliance_statuses.id
		`).Where("risk_slas.risk_register_id = ?", riskID)

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		query = query.Where(`
			(sla_register.title_en LIKE ? OR
			sla_register.title_ar LIKE ? OR
			sla_register.reference_number LIKE ? OR
			sla_register.description LIKE ? OR
			sla_register.client LIKE ? OR
			sla_categories.title_en LIKE ? OR
			sla_categories.title_ar LIKE ? OR
			organization_units.title_en LIKE ? OR
			organization_units.title_ar LIKE ? OR
			responsible_users.first_name_en LIKE ? OR
			responsible_users.first_name_ar LIKE ? OR
			responsible_users.last_name_en LIKE ? OR
			responsible_users.last_name_ar LIKE ? OR
			responsible_users.email LIKE ? OR
            responsible_users.initial LIKE ? OR
            responsible_users.color LIKE ? OR
            responsible_users.profile_image_path LIKE ? OR
			severity_levels.title_en LIKE ? OR
			severity_levels.title_ar LIKE ? OR
			compliance_statuses.title_en LIKE ? OR
			compliance_statuses.title_ar LIKE ? OR
			lifecycle_statuses.title_en LIKE ? OR
			lifecycle_statuses.title_ar LIKE ?)`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	if params.SortBy == "" {
		params.SortBy = "sla_register.title_en"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&slaList).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &riskDTOs.RiskAssociatedItemSLAListResponseDTO{
		Items: slaList,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}

	return response, nil
}

func (repo *RiskAssociatedItemRepository) GetAllIncidents(riskID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedIncidentListResponseDTO, error) {
	var incidentList []*commonDTOs.AssociatedIncidentListDTO
	var totalCount int64

	query := repo.db.Table("incident_risks").
		Select(`
			incident_register.*,
			incident_risks.id as id,
			incident_risks.incident_register_id,
            incident_categories.title_en AS incident_category_title_en,
            incident_categories.title_ar AS incident_category_title_ar,
            incident_types.title_en AS incident_type_title_en,
            incident_types.title_ar AS incident_type_title_ar,
            incident_statuses.title_en AS incident_status_title_en,
            incident_statuses.title_ar AS incident_status_title_ar,
			incident_statuses.color_code AS incident_status_color_code,
			severity_levels.title_en AS severity_level_title_en,
			severity_levels.title_ar AS severity_level_title_ar,
			severity_levels.color_code AS severity_level_color_code,
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
			reported_user.id AS reported_by_id,
            reported_user.first_name_en AS reported_user_first_name_en,
            reported_user.first_name_ar AS reported_user_first_name_ar,
            reported_user.last_name_en AS reported_user_last_name_en,
            reported_user.last_name_ar AS reported_user_last_name_ar,
			reported_user.initial AS reported_user_initial,
            reported_user.color AS reported_user_color,
            reported_user.email AS reported_user_email,
            reported_user.profile_image_path AS reported_user_profile_image_path,
			CONVERT(VARCHAR, incident_investigations.initiated_at, 23) AS investigation_initiated_at,
			incident_investigations.description AS investigation_description,
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
		Joins(`
			LEFT JOIN incident_register ON incident_register.id = incident_risks.incident_register_id
			LEFT JOIN incident_categories ON incident_categories.id = incident_register.incident_category_id
			LEFT JOIN incident_sub_categories ON incident_sub_categories.id = incident_register.incident_sub_category_id
			LEFT JOIN incident_types ON incident_types.id = incident_register.incident_type_id
			LEFT JOIN incident_statuses ON incident_statuses.id = incident_register.incident_status_id
			LEFT JOIN severity_levels ON severity_levels.id = incident_register.severity_level_id
			LEFT JOIN organization_units ON organization_units.id = incident_register.responsible_unit_id
		    LEFT JOIN users AS responsible_user ON responsible_user.id = incident_register.responsible_user_id
			LEFT JOIN users AS reported_user ON reported_user.id = incident_register.reported_by
			LEFT JOIN incident_investigations ON incident_investigations.incident_register_id = incident_register.id
			LEFT JOIN users AS created_user ON created_user.id = incident_register.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = incident_register.updated_by_id
		`).
		Where("incident_risks.risk_register_id = ?", riskID)

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(incident_register.reference_number LIKE ? OR
			incident_register.title LIKE ? OR
			incident_register.description LIKE ? OR
			incident_register.location LIKE ? OR
			incident_categories.title_en LIKE ? OR
			incident_categories.title_ar LIKE ? OR
			incident_sub_categories.title_en LIKE ? OR
			incident_sub_categories.title_ar LIKE ? OR
			incident_types.title_en LIKE ? OR
			incident_types.title_ar LIKE ? OR
			severity_levels.title_en LIKE ? OR
			severity_levels.title_ar LIKE ? OR
			incident_statuses.title_en LIKE ? OR
			incident_statuses.title_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	if params.SortBy == "" {
		params.SortBy = "incident_register.title"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&incidentList).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &commonDTOs.AssociatedIncidentListResponseDTO{
		Items: incidentList,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *RiskAssociatedItemRepository) GetAllDocuments(riskID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error) {
	var documentList []*commonDTOs.AssociatedDocumentListDTO
	var totalCount int64

	query := repo.db.Table("risk_associated_documents").
		Select(`
			risk_associated_documents.id AS id,
			risk_associated_documents.document_id,
			documents.reference_number AS reference_number,
			documents.document_type_id AS document_type_id,
			documents.title_en AS title_en,
			documents.title_ar AS title_ar,
			documents.last_checkout_at AS last_checkout_at,
			documents.last_checkin_at AS last_checkin_at,
			documents.is_locked AS is_locked,
			documents.description_en AS description_en,
			documents.description_ar AS description_ar,
			documents.responsible_user_id,
			documents.created_at,
			documents.updated_at,
			documents.published_at,
			document_categories.title_en AS document_category_title_en,
			document_categories.title_ar AS document_category_title_ar,
			document_types.title_en AS document_type_title_en,
			document_types.title_ar AS document_type_title_ar,
			document_access_types.title_en AS document_access_type_title_en,
			document_access_types.title_ar AS document_type_access_title_ar,
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
			document_versions.url AS document_url,
			document_versions.name AS document_name,
			document_versions.size AS document_size,
			document_versions.token AS document_token,
			document_versions.original_name AS document_original_name,
			document_versions.type AS document_type,
			document_versions.version AS document_version,
			document_statuses.title_en AS document_status_title_en,
			document_statuses.title_ar AS document_status_title_ar,
			document_statuses.color_code AS document_status_color_code,
			document_statuses.type AS document_status_type,
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
			published_user.first_name_en AS published_user_first_name_en,
			published_user.first_name_ar AS published_user_first_name_ar,
			published_user.last_name_en AS published_user_last_name_en,
			published_user.last_name_ar AS published_user_last_name_ar,
			published_user.initial AS published_user_initial,
			published_user.color AS published_user_color,
			published_user.email AS published_user_email,
			published_user.profile_image_path AS published_user_profile_image_path
		`).
		Joins(`
			LEFT JOIN documents ON documents.id = risk_associated_documents.document_id
			LEFT JOIN document_categories ON document_categories.id = documents.document_category_id
			LEFT JOIN document_types ON document_types.id = documents.document_type_id
			LEFT JOIN document_access_types ON document_access_types.id = documents.document_access_type_id
			LEFT JOIN organization_units ON organization_units.id = documents.responsible_unit_id
			LEFT JOIN users AS responsible_user ON responsible_user.id = documents.responsible_user_id
			LEFT JOIN document_versions ON document_versions.document_id = documents.id AND document_versions.is_latest_version = 1
			LEFT JOIN document_statuses ON document_statuses.id = documents.document_status_id
			LEFT JOIN users AS created_user ON created_user.id = documents.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = documents.updated_by_id
			LEFT JOIN users AS published_user ON published_user.id = documents.published_by_id
		`).
		Where("risk_associated_documents.risk_register_id = ?", riskID)

	for _, filter := range params.QuickFilter {
		switch filter {
		case "policy-procedure":
			query = query.Where("document_types.type IN (?)", []string{"policy", "procedure"})
		case "others":
			query = query.Where("document_types.type = ?", "others")
		}

	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(documents.title_en LIKE ? OR
			documents.title_ar LIKE ? OR
			documents.reference_number LIKE ? OR
			document_categories.title_en LIKE ? OR
			document_categories.title_ar LIKE ? OR
			responsible_user.first_name_en LIKE ? OR
			responsible_user.first_name_ar LIKE ? OR
			responsible_user.last_name_en LIKE ? OR
			responsible_user.last_name_ar LIKE ? OR
			organization_units.title_en LIKE ? OR
			organization_units.title_ar LIKE ? OR
			document_statuses.title_en LIKE ? OR
			document_statuses.title_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	if params.SortBy == "" {
		params.SortBy = "documents.title_en"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&documentList).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &commonDTOs.AssociatedDocumentListResponseDTO{
		Items: documentList,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *RiskAssociatedItemRepository) GetAllAssets(riskID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemAssetListResponseDTO, error) {
	var riskAssets []*commonDTOs.AssociatedAssetListDTO
	var totalCount int64

	query := repo.db.Table("risk_assets").
		Select(`
		    risk_assets.id,
			risk_assets.asset_register_id,
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
		Joins("LEFT JOIN asset_register ON asset_register.id = risk_assets.asset_register_id").
		Joins("LEFT JOIN asset_categories ON asset_categories.id = asset_register.asset_category_id").
		Joins("LEFT JOIN asset_sub_categories ON asset_sub_categories.id = asset_register.asset_sub_category_id").
		Joins("LEFT JOIN asset_statuses ON asset_statuses.id = asset_register.asset_status_id").
		Joins("LEFT JOIN asset_criticality_matrix_rating_scales ON asset_criticality_matrix_rating_scales.id = asset_register.asset_criticality_matrix_rating_scale_id").
		Joins("LEFT JOIN organization_units ON organization_units.id = asset_register.responsible_unit_id").
		Joins("LEFT JOIN users AS responsible_user ON responsible_user.id = asset_register.responsible_user_id").
		Joins("LEFT JOIN users AS created_user ON created_user.id = asset_register.created_by_id").
		Joins("LEFT JOIN users AS updated_user ON updated_user.id = asset_register.updated_by_id").
		Where("risk_assets.risk_register_id = ?", riskID)

	// Apply search keyword if exists
	if params.Keyword != "" {
		query = query.Where(`
			(asset_register.title_en LIKE ? OR
			asset_register.title_ar LIKE ? OR
			asset_register.reference_number LIKE ? OR
			asset_register.description LIKE ? OR
			asset_categories.title_en LIKE ? OR
			asset_categories.title_ar LIKE ? OR
			asset_statuses.title_en LIKE ? OR
			asset_statuses.title_ar LIKE ?)`,

			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	if params.SortBy == "" {
		params.SortBy = "asset_register.title_en"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&riskAssets).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &commonDTOs.AssociatedItemAssetListResponseDTO{
		Items: riskAssets,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}

	return response, nil
}

func (repo *RiskAssociatedItemRepository) GetAllThreats(riskID uint, params *riskDTOs.RiskAssociatedItemListQueryParamsDTO) (*riskDTOs.RiskAssociatedItemThreatListResponseDTO, error) {
	var riskThreats []*riskDTOs.ThreatListDTO
	var totalCount int64

	// Query the database
	query := repo.db.Table("risk_threats").
		Select(`
			risk_threats.id,
			threats.title_en,
			threats.title_ar,
			threat_types.title_en AS threat_type_title_en,
			threat_types.title_ar AS threat_type_title_ar,
			severity_levels.title_en AS severity_level_title_en,
			severity_levels.title_ar AS severity_level_title_ar,
			severity_levels.color_code AS severity_level_color_code,
			threat_statuses.title_en AS threat_status_title_en,
			threat_statuses.title_ar AS threat_status_title_ar,
			threat_statuses.color_code AS threat_status_color_code,
			threat_statuses.type AS threat_status_type,
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
			threats.description,
			threats.source,
			CONVERT(VARCHAR, threats.detection_date, 23) AS detection_date,
			CONVERT(VARCHAR, threats.resolution_date, 23) AS resolution_date,
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
		Joins("LEFT JOIN threats ON threats.id = risk_threats.threat_id").
		Joins("LEFT JOIN threat_types ON threat_types.id = threats.threat_type_id").
		Joins("LEFT JOIN severity_levels ON severity_levels.id = threats.severity_level_id").
		Joins("LEFT JOIN threat_statuses ON threat_statuses.id = threats.threat_status_id").
		Joins("LEFT JOIN organization_units ON organization_units.id = threats.responsible_unit_id").
		Joins("LEFT JOIN users AS responsible_user ON responsible_user.id = threats.responsible_user_id").
		Joins("LEFT JOIN users AS created_user ON created_user.id = threats.created_by_id").
		Joins("LEFT JOIN users AS updated_user ON updated_user.id = threats.updated_by_id").
		Where("risk_threats.risk_register_id = ?", riskID)

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply search filters
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(threats.title_en LIKE ? OR
			threats.title_ar LIKE ? OR
			threat_types.title_en LIKE ? OR
			threat_types.title_ar LIKE ? OR
			severity_levels.title_en LIKE ? OR
			severity_levels.title_ar LIKE ? OR
			threat_statuses.title_en LIKE ? OR
			threat_statuses.title_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
		)
	}

	// Apply sorting
	if params.SortBy == "" {
		params.SortBy = "threats.title_en"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}
	query = query.Order(params.SortBy + " " + params.Sort)

	// Apply pagination
	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize)

	// Fetch the data
	if err := query.Find(&riskThreats).Error; err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	// Build the response
	response := &riskDTOs.RiskAssociatedItemThreatListResponseDTO{
		Items: riskThreats,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}

	return response, nil
}

func (repo *RiskAssociatedItemRepository) GetAllClauses(riskID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemClauseListResponseDTO, error) {

	var clauseList []commonDTOs.AssociatedClauseRowDTO
	var totalCount int64
	query := repo.db.Table("risk_compliance_requirements").
		Select(`
			compliance_clauses.title_en AS clause_title_en,
			compliance_clauses.title_ar AS clause_title_ar,
			compliance_clauses.clause_no,
			compliance_clauses.compliance_register_id,
			compliance_clauses.id AS compliance_clause_id,
			compliance_register.title_en AS compliance_register_title_en,
			compliance_register.title_ar AS compliance_register_title_ar
		`).
		Joins(`
			INNER JOIN compliance_clauses ON risk_compliance_requirements.compliance_clause_id = compliance_clauses.id
			LEFT JOIN compliance_register ON risk_compliance_requirements.compliance_register_id = compliance_register.id
		`).
		Where("risk_compliance_requirements.compliance_clause_id IS NOT NULL").
		Where("risk_compliance_requirements.risk_register_id = ?", riskID)
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`(
		compliance_register.title_en LIKE ? OR
		compliance_register.title_ar LIKE ? OR
		compliance_clauses.clause_no LIKE ? OR
		compliance_clauses.title_en LIKE ? OR
		compliance_clauses.title_ar LIKE ?)`, keyword, keyword, keyword, keyword, keyword)
	}

	if params.SortBy == "" {
		params.SortBy = "compliance_register.title_en"
	}

	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	if err := query.Find(&clauseList).Error; err != nil {
		return nil, err
	}

	clausesByComplianceRegister := make(map[uint]*commonDTOs.AssociatedClauseListDTO)
	for _, row := range clauseList {
		if _, exists := clausesByComplianceRegister[row.ComplianceRegisterID]; !exists {
			clausesByComplianceRegister[row.ComplianceRegisterID] = &commonDTOs.AssociatedClauseListDTO{
				ComplianceRegisterID:      row.ComplianceRegisterID,
				ComplianceRegisterTitleEn: row.ComplianceRegisterTitleEn,
				ComplianceRegisterTitleAr: row.ComplianceRegisterTitleAr,
				Clauses:                   []commonDTOs.AssociatedClauseDTO{},
			}
		}
		clausesByComplianceRegister[row.ComplianceRegisterID].Clauses = append(clausesByComplianceRegister[row.ComplianceRegisterID].Clauses, commonDTOs.AssociatedClauseDTO{
			ComplianceClauseID: row.ComplianceClauseID,
			ClauseNo:           row.ClauseNo,
			ClauseTitleEn:      row.ClauseTitleEn,
			ClauseTitleAr:      row.ClauseTitleAr,
		})
	}

	var groupedClauses []*commonDTOs.AssociatedClauseListDTO
	for _, v := range clausesByComplianceRegister {
		// Sort clauses inside the group by ClauseNo
		sort.Slice(v.Clauses, func(i, j int) bool {
			return v.Clauses[i].ClauseNo < v.Clauses[j].ClauseNo
		})
		groupedClauses = append(groupedClauses, v)
	}

	if groupedClauses == nil {
		groupedClauses = []*commonDTOs.AssociatedClauseListDTO{}
	}

	totalCount = int64(len(groupedClauses))
	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	start := (params.Page - 1) * params.PageSize
	if start > int(totalCount) {
		start = int(totalCount)
	}

	end := start + params.PageSize
	if end > int(totalCount) {
		end = int(totalCount)
	}

	pagedClauses := groupedClauses[start:end]
	response := &commonDTOs.AssociatedItemClauseListResponseDTO{
		Items: pagedClauses,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}

	return response, nil
}

func (repo *RiskAssociatedItemRepository) GetClausesByRiskID(riskID uint, tx ...utility.DatabaseTransaction) ([]uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var clauseIDs []uint
	if err := db.
		Model(&models.RiskComplianceRequirement{}).
		Where("risk_register_id = ?", riskID).
		Pluck("compliance_clause_id", &clauseIDs).Error; err != nil {
		return nil, err
	}

	return clauseIDs, nil
}

func (repo *RiskAssociatedItemRepository) GetAllRiskKPI(riskID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedKPIRegisterListResponseDTO, error) {
	var kpiRegister []*commonDTOs.AssociatedKPIRegisterListDTO
	var totalCount int64

	query := repo.db.Table("kpi_risks").
		Select(`
            kpi_register.*,
			kpi_risks.kpi_register_id,
			kpi_categories.title_en AS kpi_category_title_en,
			kpi_categories.title_ar AS kpi_category_title_ar,
			CONVERT(VARCHAR, kpi_register.start_date, 23) AS start_date,
			CONVERT(VARCHAR, kpi_register.end_date, 23) AS target_date,
			responsible_user.id AS responsible_user_id,
			responsible_user.first_name_en AS responsible_user_first_name_en,
			responsible_user.first_name_ar AS responsible_user_first_name_ar,
			responsible_user.last_name_en AS responsible_user_last_name_en,
			responsible_user.last_name_ar AS responsible_user_last_name_ar,
			responsible_user.initial AS responsible_user_initial,
			responsible_user.color AS responsible_user_color,
			responsible_user.email AS responsible_user_email,
			responsible_user.profile_image_path AS responsible_user_profile_image_path,
			kpi_statuses.title_en AS kpi_status_title_en,
			kpi_statuses.title_ar AS kpi_status_title_ar,
			kpi_statuses.type AS kpi_status_type,
			kpi_statuses.color_code AS kpi_status_color_code,
			review_frequencies.title_en AS review_frequency_title_en,
			review_frequencies.title_ar AS review_frequency_title_ar,
			review_frequencies.type AS review_frequency_type,
			calculation_types.title_en AS calculation_type_title_en,
			calculation_types.title_ar AS calculation_type_title_ar,
			calculation_types.type AS calculation_type_type,
			units.title_en AS unit_title_en,
			units.title_ar AS unit_title_ar,
			units.code_en AS unit_code_en,
			units.code_ar AS unit_code_ar,
			units.color_code AS unit_color_code,
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
		Joins(`
		    LEFT JOIN kpi_register ON kpi_register.id = kpi_risks.kpi_register_id
			LEFT JOIN kpi_categories ON kpi_categories.id = kpi_register.kpi_category_id
			LEFT JOIN kpi_statuses ON kpi_statuses.id = kpi_register.kpi_status_id
			LEFT JOIN calculation_types ON calculation_types.id = kpi_register.calculation_type_id
			LEFT JOIN review_frequencies ON review_frequencies.id = kpi_register.review_frequency_id
			LEFT JOIN units ON units.id = kpi_register.unit_id
			LEFT JOIN users AS responsible_user ON responsible_user.id = kpi_register.responsible_user_id
			LEFT JOIN users AS created_user ON created_user.id = kpi_register.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = kpi_register.updated_by_id
        `).Where("kpi_risks.risk_register_id = ?", riskID)

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(kpi_register.title LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.last_name_ar LIKE ? OR
    		updated_user.first_name_en LIKE ? OR
			updated_user.first_name_ar LIKE ? OR
			updated_user.last_name_en LIKE ? OR
			updated_user.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	if params.SortBy == "" {
		params.SortBy = "kpi_register.title"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}
	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	query = query.Order(params.SortBy + " " + params.Sort)

	// Apply pagination
	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize)

	// Fetch the data
	if err := query.Find(&kpiRegister).Error; err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &commonDTOs.AssociatedKPIRegisterListResponseDTO{
		Items: kpiRegister,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *RiskAssociatedItemRepository) CreateRiskClauses(clauses []*models.RiskComplianceRequirement, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(&clauses).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskAssociatedItemRepository) DeleteRiskClauses(riskID uint, clauseIDs []uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.
		Where("risk_register_id = ? AND compliance_clause_id IN ?", riskID, clauseIDs).
		Delete(&models.RiskComplianceRequirement{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskAssociatedItemRepository) CreatRiskStrategicObjective(strategicObjectives []*models.RiskStrategicObjective, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&strategicObjectives).Error
}

func (repo *RiskAssociatedItemRepository) DeletetRiskStrategicObjective(strategicObjectiveID uint, riskId uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.
		Where("risk_register_id = ? AND strategic_objective_id = ?", riskId, strategicObjectiveID).
		Delete(&models.RiskStrategicObjective{}).
		Error
}

func (repo *RiskAssociatedItemRepository) CreatRiskKPI(kpiS []*models.KPIRisk, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&kpiS).Error
}

func (repo *RiskAssociatedItemRepository) DeletetRiskKPI(kpiID uint, riskId uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.
		Where("risk_register_id = ? AND kpi_register_id = ?", riskId, kpiID).
		Delete(&models.KPIRisk{}).
		Error
}

func (repo *RiskAssociatedItemRepository) CreatRiskProcess(process []*models.RiskProcess, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&process).Error
}

func (repo *RiskAssociatedItemRepository) DeletetRiskProcess(riskID uint, processID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.
		Where("risk_register_id = ? AND process_id = ?", riskID, processID).
		Delete(&models.RiskProcess{}).
		Error
}

func (repo *RiskAssociatedItemRepository) CreatRiskDocument(documents []*models.RiskAssociatedDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&documents).Error
}

func (repo *RiskAssociatedItemRepository) DeletetRiskDocument(riskID uint, documentID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.
		Where("risk_register_id = ? AND document_id = ?", riskID, documentID).
		Delete(&models.RiskAssociatedDocument{}).
		Error
}
func (repo *RiskAssociatedItemRepository) CreatRiskIncident(incidents []*models.IncidentRisk, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&incidents).Error
}

func (repo *RiskAssociatedItemRepository) DeletetRiskIncident(riskID uint, incidentID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.
		Where("risk_register_id = ? AND incident_register_id = ?", riskID, incidentID).
		Delete(&models.IncidentRisk{}).
		Error
}

func (repo *RiskAssociatedItemRepository) CreatRiskControl(controls []*models.RiskControl, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&controls).Error
}

func (repo *RiskAssociatedItemRepository) DeletetRiskControl(riskID uint, controlID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.
		Where("risk_register_id = ? AND control_id = ?", riskID, controlID).
		Delete(&models.RiskControl{}).
		Error
}

func (repo *RiskAssociatedItemRepository) CreatRiskAsset(assetes []*models.RiskAsset, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&assetes).Error
}

func (repo *RiskAssociatedItemRepository) DeletetRiskAsset(riskID uint, assetID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.
		Where("risk_register_id = ? AND asset_register_id = ?", riskID, assetID).
		Delete(&models.RiskAsset{}).
		Error
}

func (repo *RiskAssociatedItemRepository) GetAllRiskBCP(riskID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedBusinessContinuityPlanListResponseDTO, error) {
	var businessContinuityPlans []*commonDTOs.AssociatedBusinessContinuityPlanListDTO
	var totalCount int64

	query := repo.db.Table("bcp_risks").
		Select(`
            business_continuity_plans.*,
			business_functions.title AS business_function_title,
			bcp_statuses.title_en AS bcp_status_title_en,
			bcp_statuses.title_ar AS bcp_status_title_ar,
			bcp_statuses.color_code AS bcp_status_color_code,
			bcp_statuses.type AS bcp_status_type,
			rto_time_unit.title_en AS rto_unit_title_en,
			rto_time_unit.title_ar AS rto_unit_title_ar,
			mtpd_time_unit.title_en AS mtpd_unit_title_en,
			mtpd_time_unit.title_ar AS mtpd_unit_title_ar,
			rpo_time_unit.title_en AS rpo_unit_title_en,
			rpo_time_unit.title_ar AS rpo_unit_title_ar,
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
		Joins(`
		    LEFT JOIN business_continuity_plans ON business_continuity_plans.id = bcp_risks.business_continuity_plan_id
		    LEFT JOIN business_functions ON business_functions.id = business_continuity_plans.business_function_id
			LEFT JOIN bcp_statuses ON bcp_statuses.id = business_continuity_plans.bcp_status_id
			LEFT JOIN users AS responsible_user ON responsible_user.id = business_continuity_plans.responsible_user_id
			LEFT JOIN time_units AS rto_time_unit ON rto_time_unit.id = business_continuity_plans.rto_unit_id
			LEFT JOIN time_units AS mtpd_time_unit ON mtpd_time_unit.id = business_continuity_plans.mtpd_unit_id
			LEFT JOIN time_units AS rpo_time_unit ON rpo_time_unit.id = business_continuity_plans.rpo_unit_id
			LEFT JOIN users AS created_user ON created_user.id = business_continuity_plans.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = business_continuity_plans.updated_by_id
        `).Where("bcp_risks.risk_register_id = ?", riskID)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(business_continuity_plans.reference_number LIKE ? OR
			business_functions.title LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.last_name_ar LIKE ? OR
    		updated_user.first_name_en LIKE ? OR
			updated_user.first_name_ar LIKE ? OR
			updated_user.last_name_en LIKE ? OR
			updated_user.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword)
	}

	if params.SortBy == "" {
		params.SortBy = "business_functions.title"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&businessContinuityPlans).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &commonDTOs.AssociatedBusinessContinuityPlanListResponseDTO{
		Items: businessContinuityPlans,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *RiskAssociatedItemRepository) CreatRiskBCP(bcps []*models.BcpRisk, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&bcps).Error
}

func (repo *RiskAssociatedItemRepository) DeletetRiskBCP(riskID uint, bcpID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.
		Where("risk_register_id = ? AND business_continuity_plan_id = ?", riskID, bcpID).
		Delete(&models.BcpRisk{}).
		Error
}
func (repo *RiskAssociatedItemRepository) GetAllRiskProject(riskID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedProjectListResponseDTO, error) {
	var projects []*commonDTOs.AssociatedProjectListDTO
	var totalCount int64

	query := repo.db.Table("project_risks").
		Select(`
			projects.id,
			projects.title_en,
			projects.title_ar,
			projects.description_en,
			projects.description_ar,
			CONVERT(VARCHAR, projects.start_date, 23) AS start_date,
			CONVERT(VARCHAR, projects.end_date, 23) AS end_date,
			project_types.id AS project_type_id,
			project_types.title_en AS project_type_title_en,
			project_types.title_ar AS project_type_title_ar,
			project_statuses.id AS project_status_id,
			project_statuses.title_en AS project_status_title_en,
			project_statuses.title_ar AS project_status_title_ar,
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
		Joins("LEFT JOIN projects ON projects.id = project_risks.project_id").
		Joins("LEFT JOIN project_types ON project_types.id = projects.project_type_id").
		Joins("LEFT JOIN project_statuses ON project_statuses.id = projects.project_status_id").
		Joins("LEFT JOIN users AS responsible_user ON responsible_user.id = projects.responsible_user_id").
		Joins("LEFT JOIN users AS created_user ON created_user.id = projects.created_by_id").
		Joins("LEFT JOIN users AS updated_user ON updated_user.id = projects.updated_by_id").
		Where("project_risks.risk_register_id = ?", riskID)

	// Keyword search
	if params.Keyword != "" {
		like := "%" + params.Keyword + "%"
		query = query.Where(`
			projects.title_en LIKE ? OR
			projects.title_ar LIKE ? OR
			projects.description_en LIKE ? OR
			projects.description_ar LIKE ? OR
			project_types.title_en LIKE ? OR
			project_types.title_ar LIKE ? OR
			project_statuses.title_en LIKE ? OR
			project_statuses.title_ar LIKE ? OR
			responsible_user.first_name_en LIKE ? OR
			responsible_user.first_name_ar LIKE ? OR
			responsible_user.last_name_en LIKE ? OR
			responsible_user.last_name_ar LIKE ? OR
			responsible_user.email LIKE ?`,
			like, like, like, like,
			like, like, like, like,
			like, like, like, like,
			like)
	}
	if params.SortBy == "" {
		params.SortBy = "projects.title_en"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&projects).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &commonDTOs.AssociatedProjectListResponseDTO{
		Items: projects,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil

}

func (repo *RiskAssociatedItemRepository) CreatRiskProject(projects []*models.ProjectRisk, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&projects).Error
}

func (repo *RiskAssociatedItemRepository) DeletetRiskProject(riskID uint, projectID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.
		Where("risk_register_id = ? AND project_id = ?", riskID, projectID).
		Delete(&models.ProjectRisk{}).
		Error
}
