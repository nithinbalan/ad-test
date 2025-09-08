package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ComplianceAssociatedItemRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAllControls(complianceID uint, params *complianceDTOs.ComplianceAssociatedItemListQueryParamsDTO) (*complianceDTOs.ComplianceAssociatedItemControlListResponseDTO, error)
	GetAllRisks(complianceID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error)
	GetAllIncidents(complianceID uint, params *complianceDTOs.ComplianceAssociatedItemListQueryParamsDTO) (*complianceDTOs.ComplianceAssociatedItemIncidentListResponseDTO, error)
	GetAllAssociatedCounts(complianceID uint) (*complianceDTOs.ComplianceAssociatedItemsCountResponseDTO, error)
	GetComplianceRequirementCountByStatus(complianceID uint) ([]complianceDTOs.ComplianceRequirementStatusCountResponseDTO, error)
	GetAllStrategicObjectives(complianceID uint, params *complianceDTOs.ComplianceAssociatedItemListQueryParamsDTO) (*complianceDTOs.ComplianceAssociatedItemStrategicObjectiveListResponseDTO, error)
	CreateComplianceControls(controls []*models.ComplianceControl, tx ...utility.DatabaseTransaction) error
	DeleteComplianceControlMapping(complianceID, controlID uint, tx ...utility.DatabaseTransaction) error
	IsControlMappedToCompliance(complianceID, controlID uint) (bool, error)
	CreateComplianceRisk(risks []*models.RiskComplianceRequirement, tx ...utility.DatabaseTransaction) error
	DeleteComplianceRisk(complianceID uint, riskID uint, tx ...utility.DatabaseTransaction) error
	GetAllDocuments(complianceID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error)
	CreateComplianceDocument(risks []*models.DocumentComplianceRequirement, tx ...utility.DatabaseTransaction) error
	DeleteComplianceDocument(complianceID uint, riskID uint, tx ...utility.DatabaseTransaction) error
	GetComplianceControlCountByStatus(complianceID uint) ([]complianceDTOs.ComplianceRequirementStatusCountResponseDTO, error)
}

type ComplianceAssociatedItemRepository struct {
	db *gorm.DB
}

func NewComplianceAssociatedItemRepository(db *gorm.DB) *ComplianceAssociatedItemRepository {
	return &ComplianceAssociatedItemRepository{
		db: db,
	}
}

func (repo *ComplianceAssociatedItemRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ComplianceAssociatedItemRepository) GetAllControls(complianceID uint, params *complianceDTOs.ComplianceAssociatedItemListQueryParamsDTO) (*complianceDTOs.ComplianceAssociatedItemControlListResponseDTO, error) {
	var controlList []*commonDTOs.AssociatedControlListDTO
	var totalCount int64

	query := repo.db.Table("compliance_controls").
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
			LEFT JOIN controls ON controls.id = compliance_controls.control_id
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
		`).Where("compliance_controls.compliance_register_id = ?", complianceID)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

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

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&controlList).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &complianceDTOs.ComplianceAssociatedItemControlListResponseDTO{
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

func (repo *ComplianceAssociatedItemRepository) GetAllRisks(complianceID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error) {
	var totalCount int64
	var riskList []*commonDTOs.AssociatedRiskListDTO

	subQuery := repo.db.
		Table("risk_compliance_requirements").
		Select("DISTINCT risk_register.id").
		Joins("LEFT JOIN risk_register ON risk_register.id = risk_compliance_requirements.risk_register_id").
		Where("risk_compliance_requirements.compliance_register_id = ?", complianceID)

	query := repo.db.Table("risk_register").
		Select(`
		    risk_register.id AS risk_register_id,
            risk_register.reference_number,
            risk_register.title,
            risk_register.description,
			CONVERT(VARCHAR, risk_register.reported_at, 23) AS reported_at,
            risk_register.created_at,
            risk_register.updated_at,
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
			CONVERT(VARCHAR, risk_register.next_review_date, 23) AS next_review_date,
            risk_register.last_review_note,
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
            residual_risk_rating.color_code AS residual_risk_rating_color_code
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
            LEFT JOIN risk_matrix_likelihoods AS current_likelihood ON current_likelihood.id = risk_register.current_likelihood_id
            LEFT JOIN risk_matrix_impacts AS current_impact ON current_impact.id = risk_register.current_impact_id
            LEFT JOIN risk_matrix_risk_ratings AS current_risk_rating ON current_risk_rating.id = risk_register.current_risk_rating_id
            LEFT JOIN risk_matrix_risk_rating_values AS current_risk_rating_value ON current_risk_rating_value.id = risk_register.current_risk_rating_value_id
            LEFT JOIN risk_matrix_likelihoods AS inherent_likelihood ON inherent_likelihood.id = risk_register.inherent_likelihood_id
            LEFT JOIN risk_matrix_impacts AS inherent_impact ON inherent_impact.id = risk_register.inherent_impact_id
            LEFT JOIN risk_matrix_risk_ratings AS inherent_risk_rating ON inherent_risk_rating.id = risk_register.inherent_risk_rating_id
            LEFT JOIN risk_matrix_risk_rating_values AS inherent_risk_rating_value ON inherent_risk_rating_value.id = risk_register.inherent_risk_rating_value_id
            LEFT JOIN risk_matrix_likelihoods AS residual_likelihood ON residual_likelihood.id = risk_register.residual_likelihood_id
            LEFT JOIN risk_matrix_impacts AS residual_impact ON residual_impact.id = risk_register.residual_impact_id
            LEFT JOIN risk_matrix_risk_ratings AS residual_risk_rating ON residual_risk_rating.id = risk_register.residual_risk_rating_id
            LEFT JOIN risk_matrix_risk_rating_values AS residual_risk_rating_value ON residual_risk_rating_value.id = risk_register.residual_risk_rating_value_id
        `).
		Where("risk_register.id IN (?)", subQuery)

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("risk_register.title LIKE ?", keyword)
	}

	if params.Value != 0 {
		query = query.Where("current_risk_rating_value.value = ?", params.Value)
	}

	if params.SortBy == "" {
		params.SortBy = "risk_register.title"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	// Ensure that the column used for sorting is included in the SELECT clause
	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&riskList).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &commonDTOs.AssociatedRiskListResponseDTO{
		Items: riskList,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *ComplianceAssociatedItemRepository) GetAllIncidents(complianceID uint, params *complianceDTOs.ComplianceAssociatedItemListQueryParamsDTO) (*complianceDTOs.ComplianceAssociatedItemIncidentListResponseDTO, error) {
	var totalCount int64
	var incidentList []*complianceDTOs.ComplianceIncidentListDTO

	// Second Query: Fetch compliance requirement incidents using the compliance requirement IDs
	query := repo.db.Table("incident_compliance_requirements").
		Select(`
            DISTINCT incident_compliance_requirements.incident_register_id AS incident_register_id,
            incident_register.reference_number AS incident_reference_number,
            incident_register.title AS incident_register_title,
            incident_register.description AS incident_register_description,
            incident_categories.title_en AS incident_category_title_en,
            incident_categories.title_ar AS incident_category_title_ar,
            incident_sub_categories.title_en AS incident_sub_category_title_en,
            incident_sub_categories.title_ar AS incident_sub_category_title_ar,
            incident_types.title_en AS incident_type_title_en,
            incident_types.title_ar AS incident_type_title_ar,
			severity_levels.title_en AS incident_severity_level_title_en,
			severity_levels.title_ar AS incident_severity_level_title_ar,
			severity_levels.color_code AS incident_severity_level_color_code,
			incident_statuses.title_en AS incident_status_title_en,
            incident_statuses.title_ar AS incident_status_title_ar,
            incident_statuses.color_code AS incident_status_color_code,
			reported_user.id AS incident_reported_by_id,
            reported_user.first_name_en AS incident_reported_user_first_name_en,
            reported_user.first_name_ar AS incident_reported_user_first_name_ar,
            reported_user.last_name_en AS incident_reported_user_last_name_en,
            reported_user.last_name_ar AS incident_reported_user_last_name_ar,
			reported_user.initial AS reported_user_initial,
        	reported_user.color AS reported_user_color,
        	reported_user.email AS reported_user_email,
        	reported_user.profile_image_path AS reported_user_profile_image_path
        `).
		Joins(`
            LEFT JOIN incident_register ON incident_register.id = incident_compliance_requirements.incident_register_id
			LEFT JOIN incident_categories ON incident_categories.id = incident_register.incident_category_id
            LEFT JOIN incident_sub_categories ON incident_sub_categories.id = incident_register.incident_sub_category_id
            LEFT JOIN incident_types ON incident_types.id = incident_register.incident_type_id
            LEFT JOIN incident_statuses ON incident_statuses.id = incident_register.incident_status_id
		    LEFT JOIN severity_levels ON severity_levels.id = incident_register.severity_level_id
            LEFT JOIN users AS reported_user ON reported_user.id = incident_register.reported_by
        `).
		Where("incident_compliance_requirements.compliance_register_id = ?", complianceID)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("incident_register.title LIKE ?", keyword)
	}

	if params.SortBy == "" {
		params.SortBy = "incident_register.title"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	err := query.Offset(offset).Limit(params.PageSize).Find(&incidentList).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &complianceDTOs.ComplianceAssociatedItemIncidentListResponseDTO{
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

func (repo *ComplianceAssociatedItemRepository) GetAllStrategicObjectives(complianceID uint, params *complianceDTOs.ComplianceAssociatedItemListQueryParamsDTO) (*complianceDTOs.ComplianceAssociatedItemStrategicObjectiveListResponseDTO, error) {
	var strategicObjectiveList []*complianceDTOs.ComplianceStrategicObjectiveListDTO
	var totalCount int64

	query := repo.db.Table("compliance_requirement_strategic_objectives").
		Select(`
			strategic_objectives.id AS strategic_objective_id,
			strategic_objectives.title_en AS strategic_objective_title_en,
			strategic_objectives.title_ar AS strategic_objective_title_ar,
			strategic_objectives.status AS strategic_objective_status
		`).
		Joins(`
			LEFT JOIN strategic_objectives ON strategic_objectives.id = compliance_requirement_strategic_objectives.strategic_objective_id
		`).Where("compliance_requirement_strategic_objectives.compliance_register_id = ?", complianceID)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(strategic_objectives.title_en LIKE ? OR
			strategic_objectives.title_ar LIKE ?)`,
			keyword, keyword)
	}

	if params.SortBy == "" {
		params.SortBy = "strategic_objectives.title_en"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&strategicObjectiveList).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &complianceDTOs.ComplianceAssociatedItemStrategicObjectiveListResponseDTO{
		Items: strategicObjectiveList,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *ComplianceAssociatedItemRepository) GetAllAssociatedCounts(complianceID uint) (*complianceDTOs.ComplianceAssociatedItemsCountResponseDTO, error) {
	var risksCount, controlsCount, incidentsCount, requirementsCount int64
	var err error

	// Query for controls count
	err = repo.db.Model(&models.ComplianceControl{}).
		Where("compliance_register_id = ?", complianceID).
		Count(&controlsCount).Error
	if err != nil {
		return nil, err
	}

	// Query for risks count
	err = repo.db.Model(&models.RiskComplianceRequirement{}).
		Joins("JOIN compliance_requirements ON risk_compliance_requirements.compliance_requirement_id = compliance_requirements.id").
		Joins("JOIN compliance_clauses ON compliance_requirements.compliance_clause_id = compliance_clauses.id").
		Where("compliance_clauses.compliance_register_id = ?", complianceID).
		Count(&risksCount).Error
	if err != nil {
		return nil, err
	}

	// Query for incidents count
	err = repo.db.Model(&models.IncidentComplianceRequirement{}).
		Joins("JOIN compliance_requirements ON incident_compliance_requirements.compliance_requirement_id = compliance_requirements.id").
		Joins("JOIN compliance_clauses ON compliance_requirements.compliance_clause_id = compliance_clauses.id").
		Where("compliance_clauses.compliance_register_id = ?", complianceID).
		Count(&incidentsCount).Error
	if err != nil {
		return nil, err
	}
	// Query for compliance requirements count
	err = repo.db.Model(&models.ComplianceRequirement{}).
		Joins("JOIN compliance_clauses ON compliance_requirements.compliance_clause_id = compliance_clauses.id").
		Where("compliance_clauses.compliance_register_id = ?", complianceID).
		Count(&requirementsCount).Error
	if err != nil {
		return nil, err
	}

	return &complianceDTOs.ComplianceAssociatedItemsCountResponseDTO{
		Controls:               int(controlsCount),
		Risks:                  int(risksCount),
		Incidents:              int(incidentsCount),
		ComplianceRequirements: int(requirementsCount),
	}, nil
}

func (repo *ComplianceAssociatedItemRepository) GetComplianceRequirementCountByStatus(complianceID uint) ([]complianceDTOs.ComplianceRequirementStatusCountResponseDTO, error) {
	requirementStatusCounts := []complianceDTOs.ComplianceRequirementStatusCountResponseDTO{}

	query := `
        SELECT 
            compliance_statuses.id, 
            compliance_statuses.type, 
            compliance_statuses.title_en, 
            compliance_statuses.title_ar, 
            compliance_statuses.color_code,
            COUNT(compliance_requirements.id) AS Count
        FROM 
            compliance_requirements
        LEFT JOIN 
            compliance_statuses ON compliance_statuses.id = compliance_requirements.compliance_status_id
        LEFT JOIN 
            compliance_clauses ON compliance_requirements.compliance_clause_id = compliance_clauses.id
        WHERE 
            compliance_clauses.compliance_register_id = ?
        GROUP BY 
            compliance_statuses.id, compliance_statuses.type, compliance_statuses.title_en, compliance_statuses.title_ar, compliance_statuses.color_code
    `

	err := repo.db.Raw(query, complianceID).Scan(&requirementStatusCounts).Error
	if err != nil {
		return nil, err
	}

	return requirementStatusCounts, nil
}

func (repo *ComplianceAssociatedItemRepository) CreateComplianceControls(controls []*models.ComplianceControl, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&controls).Error
}

func (repo *ComplianceAssociatedItemRepository) DeleteComplianceControlMapping(complianceID, controlID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.
		Where("compliance_register_id = ? AND control_id = ?", complianceID, controlID).
		Delete(&models.ComplianceControl{}).
		Error
}

func (repo *ComplianceAssociatedItemRepository) IsControlMappedToCompliance(complianceID, controlID uint) (bool, error) {
	var count int64
	err := repo.db.
		Model(&models.ComplianceControl{}).
		Where("compliance_register_id = ? AND control_id = ?", complianceID, controlID).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repo *ComplianceAssociatedItemRepository) CreateComplianceRisk(risks []*models.RiskComplianceRequirement, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&risks).Error
}

func (repo *ComplianceAssociatedItemRepository) DeleteComplianceRisk(complianceID uint, riskId uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.
		Where("risk_register_id = ? AND compliance_register_id = ?", riskId, complianceID).
		Delete(&models.RiskComplianceRequirement{}).
		Error
}

func (repo *ComplianceAssociatedItemRepository) GetAllDocuments(complianceID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error) {
	var documentList []*commonDTOs.AssociatedDocumentListDTO
	var totalCount int64

	query := repo.db.Table("document_compliance_requirements").
		Select(`
			document_compliance_requirements.id AS id,
			document_compliance_requirements.document_id,
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
			LEFT JOIN documents ON documents.id = document_compliance_requirements.document_id
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
		Where("document_compliance_requirements.compliance_register_id = ?", complianceID)

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

func (repo *ComplianceAssociatedItemRepository) CreateComplianceDocument(documents []*models.DocumentComplianceRequirement, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&documents).Error
}

func (repo *ComplianceAssociatedItemRepository) DeleteComplianceDocument(complianceID uint, documentId uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.
		Where("document_id = ? AND compliance_register_id = ?", documentId, complianceID).
		Delete(&models.DocumentComplianceRequirement{}).
		Error
}

func (repo *ComplianceAssociatedItemRepository) GetComplianceControlCountByStatus(complianceID uint) ([]complianceDTOs.ComplianceRequirementStatusCountResponseDTO, error) {
	controlStatusCounts := []complianceDTOs.ComplianceRequirementStatusCountResponseDTO{}

	query := `
        SELECT 
            cs.id,
            cs.type,
            cs.title_en,
            cs.title_ar,
            cs.color_code,
            COUNT(c.id) AS Count
        FROM 
            compliance_controls cc
        JOIN 
            controls c ON cc.control_id = c.id
        LEFT JOIN 
            compliance_statuses cs ON cs.id = c.compliance_status_id
        WHERE 
            cc.compliance_register_id = ?
        GROUP BY 
            cs.id, cs.type, cs.title_en, cs.title_ar, cs.color_code
    `

	err := repo.db.Raw(query, complianceID).Scan(&controlStatusCounts).Error
	if err != nil {
		return nil, err
	}

	return controlStatusCounts, nil
}
