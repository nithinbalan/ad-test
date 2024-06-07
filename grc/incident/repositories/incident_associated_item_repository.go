package repositories

import (
	incidentDTOs "jethur/grc/incident/dtos"

	"gorm.io/gorm"
)

type IncidentAssociatedItemRepositoryInterface interface {
	GetAllControls(incidentID uint, params *incidentDTOs.IncidentAssociatedItemListQueryParamsDTO) ([]*incidentDTOs.IncidentControlListDTO, error)
	GetAllComplianceRequirements(incidentID uint, params *incidentDTOs.IncidentAssociatedItemListQueryParamsDTO) ([]*incidentDTOs.IncidentComplianceRequirementListDTO, error)
	GetAllRisks(incidentID uint, params *incidentDTOs.IncidentAssociatedItemListQueryParamsDTO) ([]*incidentDTOs.IncidentRiskListDTO, error)
}

type IncidentAssociatedItemRepository struct {
	db *gorm.DB
}

func NewIncidentAssociatedItemRepository(db *gorm.DB) *IncidentAssociatedItemRepository {
	return &IncidentAssociatedItemRepository{
		db: db,
	}
}

func (repo *IncidentAssociatedItemRepository) GetAllControls(incidentID uint, params *incidentDTOs.IncidentAssociatedItemListQueryParamsDTO) ([]*incidentDTOs.IncidentControlListDTO, error) {
	var controlList []*incidentDTOs.IncidentControlListDTO

	query := repo.db.Table("incident_controls").
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
			LEFT JOIN controls ON controls.id = incident_controls.control_id
            LEFT JOIN control_categories ON controls.control_category_id = control_categories.id
            LEFT JOIN control_types ON controls.control_type_id = control_types.id
            LEFT JOIN control_efficiency_statuses ON controls.control_efficiency_status_id = control_efficiency_statuses.id
            LEFT JOIN control_implementation_statuses ON controls.control_implementation_status_id = control_implementation_statuses.id
        `).
		Where("incident_controls.incident_register_id = ?", incidentID)

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

	return controlList, nil
}

func (repo *IncidentAssociatedItemRepository) GetAllComplianceRequirements(incidentID uint, params *incidentDTOs.IncidentAssociatedItemListQueryParamsDTO) ([]*incidentDTOs.IncidentComplianceRequirementListDTO, error) {
	var complianceRequirements []*incidentDTOs.IncidentComplianceRequirementListDTO

	query := repo.db.Table("incident_compliance_requirements").
		Select(`
			incident_compliance_requirements.compliance_requirement_id AS compliance_requirement_id,
            compliance_clauses.title AS compliance_clause_title,
            compliance_requirement_categories.title_en AS compliance_requirement_category_title_en,
            compliance_requirement_categories.title_ar AS compliance_requirement_category_title_ar,
            compliance_requirements.title AS compliance_requirement_title,
            compliance_statuses.title_en AS compliance_status_title_en,
            compliance_statuses.title_ar AS compliance_status_title_ar,
			compliance_statuses.color_code AS compliance_status_color_code

        `).
		Joins(`
			LEFT JOIN compliance_requirements ON compliance_requirements.id = incident_compliance_requirements.compliance_requirement_id
            LEFT JOIN compliance_clauses ON compliance_requirements.compliance_clause_id = compliance_clauses.id
            LEFT JOIN compliance_statuses ON compliance_requirements.compliance_status_id = compliance_statuses.id
            LEFT JOIN compliance_requirement_categories ON compliance_requirements.compliance_requirement_category_id = compliance_requirement_categories.id
        `).
		Where("incident_compliance_requirements.incident_register_id = ?", incidentID)

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(compliance_requirements.title LIKE ?)`,
			keyword)
	}

	if params.SortBy == "" {
		params.SortBy = "compliance_requirements.title"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	if err := query.Find(&complianceRequirements).Error; err != nil {
		return nil, err
	}

	return complianceRequirements, nil
}

func (repo *IncidentAssociatedItemRepository) GetAllRisks(incidentID uint, params *incidentDTOs.IncidentAssociatedItemListQueryParamsDTO) ([]*incidentDTOs.IncidentRiskListDTO, error) {
	var incidentRisk []*incidentDTOs.IncidentRiskListDTO

	query := repo.db.Table("incident_risks").
		Select(`
			incident_risks.risk_register_id AS risk_register_id,
			risk_register.title AS risk_register_title,
			risk_statuses.title_en AS risk_status_title_en,
			risk_statuses.title_ar AS risk_status_title_ar,
			risk_statuses.color_code AS risk_status_color_code
		`).
		Joins(`
			LEFT JOIN risk_register ON risk_register.id = incident_risks.risk_register_id
			LEFT JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id
		`).
		Where("incident_risks.incident_register_id = ?", incidentID)

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(risk_register.title LIKE ?)`,
			keyword)
	}

	if params.SortBy == "" {
		params.SortBy = "risk_register.title"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	if err := query.Find(&incidentRisk).Error; err != nil {
		return nil, err
	}

	return incidentRisk, nil
}
