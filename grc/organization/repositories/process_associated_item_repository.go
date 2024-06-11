package repositories

import (
	organizationDTOs "jethur/grc/organization/dtos"

	"gorm.io/gorm"
)

type ProcessAssociatedItemRepositoryInterface interface {
	GetAllControls(processID uint, params *organizationDTOs.ProcessAssociatedItemListQueryParamsDTO) (map[string]map[string][]*organizationDTOs.ProcessControlListDTO, error)
	GetAllRisks(processID uint, params *organizationDTOs.ProcessAssociatedItemListQueryParamsDTO) ([]*organizationDTOs.ProcessRiskListDTO, error)
	GetAllIncidents(processID uint, params *organizationDTOs.ProcessAssociatedItemListQueryParamsDTO) ([]*organizationDTOs.ProcessIncidentListDTO, error)
}

type ProcessAssociatedItemRepository struct {
	db *gorm.DB
}

func NewProcessAssociatedItemRepository(db *gorm.DB) *ProcessAssociatedItemRepository {
	return &ProcessAssociatedItemRepository{
		db: db,
	}
}

func (repo *ProcessAssociatedItemRepository) GetAllControls(processID uint, params *organizationDTOs.ProcessAssociatedItemListQueryParamsDTO) (map[string]map[string][]*organizationDTOs.ProcessControlListDTO, error) {
	var controlList []*organizationDTOs.ProcessControlListDTO

	query := repo.db.Table("process_controls").
		Select(`
			controls.id AS control_id,
            control_categories.title_en AS control_category_title_en,
            control_categories.title_ar AS control_category_title_ar,
            control_types.title_en AS control_type_title_en,
            control_types.title_ar AS control_type_title_ar,
            controls.title_en AS control_title_en,
            controls.title_ar AS control_title_ar,
            controls.description AS control_description,
            control_efficiency_statuses.title_en AS control_efficiency_status_en,
            control_efficiency_statuses.title_ar AS control_efficiency_status_ar,
            control_implementation_statuses.title_en AS control_implementation_status_en,
            control_implementation_statuses.title_ar AS control_implementation_status_ar
        `).
		Joins(`
			LEFT JOIN controls ON controls.id = process_controls.control_id
            LEFT JOIN control_categories ON controls.control_category_id = control_categories.id
            LEFT JOIN control_types ON controls.control_type_id = control_types.id
            LEFT JOIN control_efficiency_statuses ON controls.control_efficiency_status_id = control_efficiency_statuses.id
            LEFT JOIN control_implementation_statuses ON controls.control_implementation_status_id = control_implementation_statuses.id
        `).Where("process_controls.process_id = ?", processID)

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(controls.title_en LIKE ? OR
			controls.title_ar LIKE ? OR
			controls.reference_number LIKE ? OR
			control_categories.title_en LIKE ? OR
			control_categories.title_ar LIKE ? OR
			control_types.title_en LIKE ? OR
			control_types.title_ar LIKE ? OR
			control_efficiency_statuses.title_en LIKE ? OR
			control_efficiency_statuses.title_ar LIKE ? OR
			control_implementation_statuses.title_en LIKE ? OR
			control_implementation_statuses.title_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword)
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

	// Grouping controls by implementation status and then by efficiency status
	groupedControls := make(map[string]map[string][]*organizationDTOs.ProcessControlListDTO)

	for _, control := range controlList {
		implementationStatus := control.ControlImplementationStatusEn
		efficiencyStatus := control.ControlEfficiencyStatusEn

		if _, exists := groupedControls[implementationStatus]; !exists {
			groupedControls[implementationStatus] = make(map[string][]*organizationDTOs.ProcessControlListDTO)
		}
		if _, exists := groupedControls[implementationStatus][efficiencyStatus]; !exists {
			groupedControls[implementationStatus][efficiencyStatus] = []*organizationDTOs.ProcessControlListDTO{}
		}
		groupedControls[implementationStatus][efficiencyStatus] = append(groupedControls[implementationStatus][efficiencyStatus], control)
	}

	return groupedControls, nil
}

func (repo *ProcessAssociatedItemRepository) GetAllRisks(processID uint, params *organizationDTOs.ProcessAssociatedItemListQueryParamsDTO) ([]*organizationDTOs.ProcessRiskListDTO, error) {
	var riskList []*organizationDTOs.ProcessRiskListDTO

	query := repo.db.Table("risk_processes").
		Select(`
			risk_register.reference_number,
			risk_register.title,
			risk_register.description,
			risk_categories.title_en AS risk_category_title_en,
			risk_categories.title_ar AS risk_category_title_ar,
			risk_sources.title_en AS risk_source_title_en,
			risk_sources.title_ar AS risk_source_title_ar,
			organization_units.title_en AS owner_unit_title_en,
			organization_units.title_ar AS owner_unit_title_ar,
			users.first_name_en AS owner_first_name_en,
			users.first_name_ar AS owner_first_name_ar,
			users.last_name_en AS owner_last_name_en,
			users.last_name_ar AS owner_last_name_ar,
			risk_statuses.title_en AS risk_status_title_en,
			risk_statuses.title_ar AS risk_status_title_ar,
			risk_statuses.color_code AS risk_status_color_code
		`).
		Joins(`
			LEFT JOIN risk_register ON risk_register.id = risk_processes.risk_register_id
			LEFT JOIN risk_categories ON risk_categories.id = risk_register.risk_category_id
			LEFT JOIN risk_sources ON risk_sources.id = risk_register.risk_source_id
			LEFT JOIN organization_units ON organization_units.id = risk_register.owner_unit_id
			LEFT JOIN users ON users.id = risk_register.owner_id
			LEFT JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id
		`).Where("risk_processes.process_id = ?", processID)

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

	if params.SortBy == "" {
		params.SortBy = "risk_register.title"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	if err := query.Find(&riskList).Error; err != nil {
		return nil, err
	}

	return riskList, nil
}

func (repo *ProcessAssociatedItemRepository) GetAllIncidents(processID uint, params *organizationDTOs.ProcessAssociatedItemListQueryParamsDTO) ([]*organizationDTOs.ProcessIncidentListDTO, error) {
	var incidentList []*organizationDTOs.ProcessIncidentListDTO

	query := repo.db.Table("incident_processes").
		Select(`
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

	if err := query.Find(&incidentList).Error; err != nil {
		return nil, err
	}

	return incidentList, nil
}
