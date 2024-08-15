package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"

	"gorm.io/gorm"
)

type ProcessAssociatedItemRepositoryInterface interface {
	GetAllControls(processID uint, params *organizationDTOs.ProcessAssociatedItemListQueryParamsDTO) ([]*organizationDTOs.ProcessControlListDTO, error)
	GetAllRisks(processID uint, params *organizationDTOs.ProcessAssociatedItemListQueryParamsDTO) (*organizationDTOs.ProcessAssociatedItemRiskListResponseDTO, error)
	GetAllIncidents(processID uint, params *organizationDTOs.ProcessAssociatedItemListQueryParamsDTO) (*organizationDTOs.ProcessAssociatedItemIncidentListResponseDTO, error)
	GetAllAssociatedCounts(processID uint) (*organizationDTOs.ProcessAssociatedItemsCountResponseDTO, error)
}

type ProcessAssociatedItemRepository struct {
	db *gorm.DB
}

func NewProcessAssociatedItemRepository(db *gorm.DB) *ProcessAssociatedItemRepository {
	return &ProcessAssociatedItemRepository{
		db: db,
	}
}

func (repo *ProcessAssociatedItemRepository) GetAllControls(processID uint, params *organizationDTOs.ProcessAssociatedItemListQueryParamsDTO) ([]*organizationDTOs.ProcessControlListDTO, error) {
	var controlList []*organizationDTOs.ProcessControlListDTO

	query := repo.db.Table("process_controls").
		Select(`
			process_controls.id AS id,
			process_controls.process_id AS process_id,
			controls.id AS control_id,
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
            control_efficiency_statuses.type AS control_efficiency_status_type,
            control_efficiency_statuses.color_code AS control_efficiency_status_color_code,
            control_implementation_statuses.title_en AS control_implementation_status_title_en,
            control_implementation_statuses.title_ar AS control_implementation_status_title_ar,
            control_implementation_statuses.type AS control_implementation_status_type,
            control_implementation_statuses.color_code AS control_implementation_status_color_code,
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
			LEFT JOIN controls ON controls.id = process_controls.control_id
            LEFT JOIN control_categories ON controls.control_category_id = control_categories.id
            LEFT JOIN control_types ON controls.control_type_id = control_types.id
            LEFT JOIN control_efficiency_statuses ON controls.control_efficiency_status_id = control_efficiency_statuses.id
            LEFT JOIN control_implementation_statuses ON controls.control_implementation_status_id = control_implementation_statuses.id
			LEFT JOIN users AS owner ON controls.owner_id = owner.id

        `).Where("process_controls.process_id = ?", processID)

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

	if err := query.Find(&controlList).Error; err != nil {
		return nil, err
	}

	return controlList, nil
}

func (repo *ProcessAssociatedItemRepository) GetAllRisks(processID uint, params *organizationDTOs.ProcessAssociatedItemListQueryParamsDTO) (*organizationDTOs.ProcessAssociatedItemRiskListResponseDTO, error) {
	var riskList []*organizationDTOs.ProcessRiskListDTO
	var totalCount int64

	query := repo.db.Table("risk_processes").
		Select(`
			risk_processes.id,
			risk_processes.risk_register_id,
			risk_processes.process_id,
			risk_register.reference_number,
			risk_register.title,
			risk_register.description,
			risk_categories.title_en AS risk_category_title_en,
			risk_categories.title_ar AS risk_category_title_ar,
			risk_sources.title_en AS risk_source_title_en,
			risk_sources.title_ar AS risk_source_title_ar,
			organization_units.title_en AS responsible_unit_title_en,
			organization_units.title_ar AS responsible_unit_title_ar,
			users.first_name_en AS responsible_user_first_name_en,
			users.first_name_ar AS responsible_user_first_name_ar,
			users.last_name_en AS responsible_user_last_name_en,
			users.last_name_ar AS responsible_user_last_name_ar,
			users.initial AS responsible_user_initial,
			users.color AS responsible_user_color,
			users.email AS responsible_user_email,
			users.profile_image_path AS responsible_user_profile_image_path,
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
			risk_statuses.title_en AS risk_status_title_en,
			risk_statuses.title_ar AS risk_status_title_ar,
			risk_statuses.color_code AS risk_status_color_code
		`).
		Joins(`
			LEFT JOIN risk_register ON risk_register.id = risk_processes.risk_register_id
			LEFT JOIN risk_categories ON risk_categories.id = risk_register.risk_category_id
			LEFT JOIN risk_sources ON risk_sources.id = risk_register.risk_source_id
			LEFT JOIN organization_units ON organization_units.id = risk_register.responsible_unit_id
			LEFT JOIN users ON users.id = risk_register.responsible_user_id
			LEFT JOIN risk_matrix_likelihoods ON risk_matrix_likelihoods.id = risk_register.current_likelihood_id
            LEFT JOIN risk_matrix_impacts ON risk_matrix_impacts.id = risk_register.current_impact_id
            LEFT JOIN risk_matrix_risk_ratings ON risk_matrix_risk_ratings.id = risk_register.current_risk_rating_id
			LEFT JOIN risk_matrix_risk_rating_values ON risk_matrix_risk_rating_values.id = risk_register.current_risk_rating_value_id
			LEFT JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id
		`).Where("risk_processes.process_id = ?", processID)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(risk_register.reference_number LIKE ? OR
			risk_register.title LIKE ? OR
			risk_register.description LIKE ? OR
			risk_categories.title_en LIKE ? OR
			risk_categories.title_ar LIKE ? OR
			risk_sources.title_en LIKE ? OR
			risk_sources.title_ar LIKE ? OR
			organization_units.title_en LIKE ? OR
			organization_units.title_ar LIKE ? OR
			users.first_name_en LIKE ? OR
			users.first_name_ar LIKE ? OR
			users.last_name_en LIKE ? OR
			users.last_name_ar LIKE ? OR
			risk_statuses.title_en LIKE ? OR
			risk_statuses.title_ar LIKE ? OR
			risk_statuses.color_code LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
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

	query = query.Order(params.SortBy + " " + params.Sort)
	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&riskList).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &organizationDTOs.ProcessAssociatedItemRiskListResponseDTO{
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

func (repo *ProcessAssociatedItemRepository) GetAllIncidents(processID uint, params *organizationDTOs.ProcessAssociatedItemListQueryParamsDTO) (*organizationDTOs.ProcessAssociatedItemIncidentListResponseDTO, error) {
	var incidentList []*organizationDTOs.ProcessIncidentListDTO
	var totalCount int64

	query := repo.db.Table("incident_processes").
		Select(`
			incident_processes.id,
			incident_processes.incident_register_id,
			incident_processes.process_id,
			incident_register.reference_number,
			incident_register.title,
			incident_register.description,
			incident_register.incident_at,
			incident_register.location,
			incident_register.reported_at,
			incident_register.quick_action,
			incident_categories.title_en AS incident_category_title_en,
			incident_categories.title_ar AS incident_category_title_ar,
			incident_sub_categories.title_en AS incident_sub_category_title_en,
			incident_sub_categories.title_ar AS incident_sub_category_title_ar,
			incident_types.title_en AS incident_type_title_en,
			incident_types.title_ar AS incident_type_title_ar,
			severity_levels.title_en AS severity_level_title_en,
			severity_levels.title_ar AS severity_level_title_ar,
			severity_levels.color_code AS severity_level_color_code,
			incident_statuses.title_en AS incident_status_title_en,
			incident_statuses.title_ar AS incident_status_title_ar,
			incident_statuses.color_code AS incident_status_color_code
		`).
		Joins(`
			LEFT JOIN incident_register ON incident_register.id = incident_processes.incident_register_id
			LEFT JOIN incident_categories ON incident_categories.id = incident_register.incident_category_id
			LEFT JOIN incident_sub_categories ON incident_sub_categories.id = incident_register.incident_sub_category_id
			LEFT JOIN incident_types ON incident_types.id = incident_register.incident_type_id
			LEFT JOIN incident_statuses ON incident_statuses.id = incident_register.incident_status_id
			LEFT JOIN severity_levels ON severity_levels.id = incident_register.severity_level_id
		`).
		Where("incident_processes.process_id = ?", processID)
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

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

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&incidentList).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &organizationDTOs.ProcessAssociatedItemIncidentListResponseDTO{
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

func (repo *ProcessAssociatedItemRepository) GetAllAssociatedCounts(processID uint) (*organizationDTOs.ProcessAssociatedItemsCountResponseDTO, error) {
	var risksCount, controlsCount, incidentsCount int64
	var err error

	// Query for controls count
	err = repo.db.Model(&models.ProcessControl{}).
		Where("process_id = ?", processID).
		Count(&controlsCount).Error
	if err != nil {
		return nil, err
	}

	// Query for risks count
	err = repo.db.Model(&models.RiskProcess{}).
		Where("process_id = ?", processID).
		Count(&risksCount).Error
	if err != nil {
		return nil, err
	}

	// Query for incidents count
	err = repo.db.Model(&models.IncidentProcess{}).
		Where("process_id = ?", processID).
		Count(&incidentsCount).Error
	if err != nil {
		return nil, err
	}

	return &organizationDTOs.ProcessAssociatedItemsCountResponseDTO{
		Controls:  int(controlsCount),
		Risks:     int(risksCount),
		Incidents: int(incidentsCount),
	}, nil
}
