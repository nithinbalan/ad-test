package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/models"

	"gorm.io/gorm"
)

type ComplianceAssociatedItemRepositoryInterface interface {
	GetAllControls(complianceID uint, params *complianceDTOs.ComplianceAssociatedItemListQueryParamsDTO) (*complianceDTOs.ComplianceAssociatedItemControlListResponseDTO, error)
	GetAllRisks(complianceID uint, params *complianceDTOs.ComplianceAssociatedItemListQueryParamsDTO) (*complianceDTOs.ComplianceAssociatedItemRiskListResponseDTO, error)
	GetAllIncidents(complianceID uint, params *complianceDTOs.ComplianceAssociatedItemListQueryParamsDTO) (*complianceDTOs.ComplianceAssociatedItemIncidentListResponseDTO, error)
	GetAllAssociatedCounts(complianceID uint) (*complianceDTOs.ComplianceAssociatedItemsCountResponseDTO, error)
	GetComplianceRequirementCountByStatus(complianceID uint) ([]complianceDTOs.ComplianceRequirementStatusCountResponseDTO, error)
	GetAllStrategicObjectives(complianceID uint, params *complianceDTOs.ComplianceAssociatedItemListQueryParamsDTO) (*complianceDTOs.ComplianceAssociatedItemStrategicObjectiveListResponseDTO, error)
}

type ComplianceAssociatedItemRepository struct {
	db *gorm.DB
}

func NewComplianceAssociatedItemRepository(db *gorm.DB) *ComplianceAssociatedItemRepository {
	return &ComplianceAssociatedItemRepository{
		db: db,
	}
}

func (repo *ComplianceAssociatedItemRepository) GetAllControls(complianceID uint, params *complianceDTOs.ComplianceAssociatedItemListQueryParamsDTO) (*complianceDTOs.ComplianceAssociatedItemControlListResponseDTO, error) {
	var controlList []*complianceDTOs.ComplianceControlListDTO
	var totalCount int64

	query := repo.db.Table("compliance_controls").
		Select(`
			DISTINCT controls.id AS control_id,
			controls.reference_number AS control_reference_number,
            controls.code AS control_code,
			CONVERT(VARCHAR, controls.last_test_date, 23) AS control_last_test_date,
            control_categories.title_en AS control_category_title_en,
            control_categories.title_ar AS control_category_title_ar,
            control_types.title_en AS control_type_title_en,
            control_types.title_ar AS control_type_title_ar,
            controls.title_en AS control_title_en,
            controls.title_ar AS control_title_ar,
            controls.description AS control_description,
            control_efficiency_statuses.title_en AS control_efficiency_status_title_en,
            control_efficiency_statuses.title_ar AS control_efficiency_status_title_ar,
            control_efficiency_statuses.color_code AS control_efficiency_status_color_code,
            control_efficiency_statuses.type AS control_efficiency_status_type,
            control_implementation_statuses.title_en AS control_implementation_status_title_en,
            control_implementation_statuses.title_ar AS control_implementation_status_title_ar,
            control_implementation_statuses.type AS control_implementation_status_type,
            control_implementation_statuses.color_code AS control_implementation_status_color_code,
			owner.id AS control_owner_id,
			owner.first_name_en AS control_owner_first_name_en,
			owner.first_name_ar AS control_owner_first_name_ar,
			owner.last_name_en AS control_owner_last_name_en,
			owner.last_name_ar AS control_owner_last_name_ar,
			owner.initial AS control_owner_initial,
            owner.color AS control_owner_color,
            owner.email AS control_owner_email,
            owner.profile_image_path AS control_owner_profile_image_path
        `).
		Joins(`
			LEFT JOIN controls ON controls.id = compliance_controls.control_id
            LEFT JOIN control_categories ON controls.control_category_id = control_categories.id
            LEFT JOIN control_types ON controls.control_type_id = control_types.id
            LEFT JOIN control_efficiency_statuses ON controls.control_efficiency_status_id = control_efficiency_statuses.id
            LEFT JOIN control_implementation_statuses ON controls.control_implementation_status_id = control_implementation_statuses.id
			LEFT JOIN users AS owner ON controls.owner_id = owner.id
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

func (repo *ComplianceAssociatedItemRepository) GetAllRisks(complianceID uint, params *complianceDTOs.ComplianceAssociatedItemListQueryParamsDTO) (*complianceDTOs.ComplianceAssociatedItemRiskListResponseDTO, error) {
	var totalCount int64
	var riskList []*complianceDTOs.ComplianceRiskListDTO

	// Second Query: Fetch compliance requirement risks using the compliance requirement IDs
	query := repo.db.Table("risk_compliance_requirements").
		Select(`
            DISTINCT risk_compliance_requirements.risk_register_id AS risk_register_id,
            risk_register.reference_number AS risk_reference_number,
            risk_register.title AS risk_register_title,
            risk_register.description AS risk_register_description,
			risk_categories.title_en AS risk_category_title_en,
            risk_categories.title_ar AS risk_category_title_ar,
            risk_sources.title_en AS risk_source_title_en,
            risk_sources.title_ar AS risk_source_title_ar,
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
            risk_statuses.color_code AS risk_status_color_code,
            review_frequencies.title_en AS risk_review_frequency_title_en,
            review_frequencies.title_ar AS risk_review_frequency_title_ar,
			risk_matrix_likelihoods.title_en AS risk_matrix_likelihood_title_en,
            risk_matrix_likelihoods.title_ar AS risk_matrix_likelihood_title_ar,
            risk_matrix_likelihoods.score AS risk_matrix_likelihood_score,
			risk_matrix_impacts.title_en AS risk_matrix_impact_title_en,
            risk_matrix_impacts.title_ar AS risk_matrix_impact_title_ar,
            risk_matrix_impacts.score AS risk_matrix_impact_score,
			risk_matrix_risk_ratings.title_en AS risk_matrix_risk_rating_title_en,
            risk_matrix_risk_ratings.title_ar AS risk_matrix_risk_rating_title_ar,
            risk_matrix_risk_ratings.type AS risk_matrix_risk_rating_type,
            risk_matrix_risk_ratings.color_code AS risk_matrix_risk_rating_color_code,
			risk_matrix_risk_rating_values.value AS risk_matrix_risk_rating_value,
            risk_register.next_review_date AS risk_next_review_date,
            risk_statuses.title_en AS risk_status_title_en,
            risk_statuses.title_ar AS risk_status_title_ar,
            risk_statuses.color_code AS risk_status_color_code
        `).
		Joins(`
            LEFT JOIN risk_register ON risk_register.id = risk_compliance_requirements.risk_register_id
			LEFT JOIN risk_categories ON risk_categories.id = risk_register.risk_category_id
            LEFT JOIN risk_sources ON risk_sources.id = risk_register.risk_source_id
            LEFT JOIN organization_units ON organization_units.id = risk_register.responsible_unit_id
            LEFT JOIN users ON users.id = risk_register.responsible_user_id
            LEFT JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id
            LEFT JOIN risk_matrix_likelihoods ON risk_matrix_likelihoods.id = risk_register.current_likelihood_id
            LEFT JOIN risk_matrix_impacts ON risk_matrix_impacts.id = risk_register.current_impact_id
            LEFT JOIN risk_matrix_risk_ratings ON risk_matrix_risk_ratings.id = risk_register.current_risk_rating_id
			LEFT JOIN risk_matrix_risk_rating_values ON risk_matrix_risk_rating_values.id = risk_register.current_risk_rating_value_id
            LEFT JOIN review_frequencies ON review_frequencies.id = risk_register.review_frequency_id
        `).
		Where("risk_compliance_requirements.compliance_register_id = ?", complianceID)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("risk_register.title LIKE ?", keyword)
	}

	if params.Value != 0 {
		query = query.Where("risk_matrix_risk_rating_values.value = ?", params.Value)
	}

	if params.SortBy == "" {
		params.SortBy = "risk_register.title"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
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

	response := &complianceDTOs.ComplianceAssociatedItemRiskListResponseDTO{
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
