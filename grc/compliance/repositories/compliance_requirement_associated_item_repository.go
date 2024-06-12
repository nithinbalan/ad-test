package repositories

import (
	complianceDTOs "jethur/grc/compliance/dtos"

	"gorm.io/gorm"
)

type ComplianceRequirementAssociatedItemRepositoryInterface interface {
	GetAllRisks(complianceID uint, params *complianceDTOs.ComplianceRequirementAssociatedItemListQueryParamsDTO) ([]*complianceDTOs.ComplianceRequirementRiskListDTO, error)
	GetAllIncidents(complianceID uint, params *complianceDTOs.ComplianceRequirementAssociatedItemListQueryParamsDTO) ([]*complianceDTOs.ComplianceRequirementIncidentListDTO, error)
}

type ComplianceRequirementAssociatedItemRepository struct {
	db *gorm.DB
}

func NewComplianceRequirementAssociatedItemRepository(db *gorm.DB) *ComplianceRequirementAssociatedItemRepository {
	return &ComplianceRequirementAssociatedItemRepository{
		db: db,
	}
}

func (repo *ComplianceRequirementAssociatedItemRepository) GetAllRisks(complianceID uint, params *complianceDTOs.ComplianceRequirementAssociatedItemListQueryParamsDTO) ([]*complianceDTOs.ComplianceRequirementRiskListDTO, error) {
	var complianceRequirementIDs []uint

	// First Query: Fetch all compliance requirement IDs by joining compliance_clauses and compliance_requirements
	err := repo.db.Table("compliance_requirements").
		Joins("JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id").
		Where("compliance_clauses.compliance_register_id = ?", complianceID).
		Pluck("compliance_requirements.id", &complianceRequirementIDs).Error

	if err != nil {
		return nil, err
	}

	// Check if we found any compliance requirement IDs
	if len(complianceRequirementIDs) == 0 {
		return []*complianceDTOs.ComplianceRequirementRiskListDTO{}, nil
	}

	var riskList []*complianceDTOs.ComplianceRequirementRiskListDTO

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
            organization_units.title_en AS risk_owner_unit_title_en,
            organization_units.title_ar AS risk_owner_unit_title_ar,
            users.first_name_en AS risk_owner_first_name_en,
            users.first_name_ar AS risk_owner_first_name_ar,
            users.last_name_en AS risk_owner_last_name_en,
            users.last_name_ar AS risk_owner_last_name_ar,
            risk_statuses.title_en AS risk_status_title_en,
            risk_statuses.title_ar AS risk_status_title_ar,
            risk_statuses.color_code AS risk_status_color_code,
            review_frequencies.title_en AS risk_review_frequency_title_en,
            review_frequencies.title_ar AS risk_review_frequency_title_ar,
            risk_register.next_review_date AS risk_next_review_date,
            risk_statuses.title_en AS risk_status_title_en,
            risk_statuses.title_ar AS risk_status_title_ar,
            risk_statuses.color_code AS risk_status_color_code
        `).
		Joins(`
            LEFT JOIN risk_register ON risk_register.id = risk_compliance_requirements.risk_register_id
			LEFT JOIN risk_categories ON risk_categories.id = risk_register.risk_category_id
            LEFT JOIN risk_sources ON risk_sources.id = risk_register.risk_source_id
            LEFT JOIN organization_units ON organization_units.id = risk_register.owner_unit_id
            LEFT JOIN users ON users.id = risk_register.owner_id
            LEFT JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id
            LEFT JOIN review_frequencies ON review_frequencies.id = risk_register.review_frequency_id
        `).
		Where("risk_compliance_requirements.compliance_requirement_id IN (?)", complianceRequirementIDs)

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("risk_register.title LIKE ?", keyword)
	}

	if params.SortBy == "" {
		params.SortBy = "risk_register.title"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	// Ensure that the column used for sorting is included in the SELECT clause
	query = query.Order(params.SortBy + " " + params.Sort)

	// Add Debug() to print the generated SQL
	query = query.Debug()

	if err := query.Find(&riskList).Error; err != nil {
		return nil, err
	}

	return riskList, nil
}

func (repo *ComplianceRequirementAssociatedItemRepository) GetAllIncidents(complianceID uint, params *complianceDTOs.ComplianceRequirementAssociatedItemListQueryParamsDTO) ([]*complianceDTOs.ComplianceRequirementIncidentListDTO, error) {
	var complianceRequirementIDs []uint

	// First Query: Fetch all compliance requirement IDs by joining compliance_clauses and compliance_requirements
	err := repo.db.Table("compliance_requirements").
		Joins("JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id").
		Where("compliance_clauses.compliance_register_id = ?", complianceID).
		Pluck("compliance_requirements.id", &complianceRequirementIDs).Error

	if err != nil {
		return nil, err
	}
	// Check if we found any compliance requirement IDs
	if len(complianceRequirementIDs) == 0 {
		return []*complianceDTOs.ComplianceRequirementIncidentListDTO{}, nil
	}

	var incidentList []*complianceDTOs.ComplianceRequirementIncidentListDTO

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
            reported_user.first_name_en AS incident_reported_by_first_name_en,
            reported_user.first_name_ar AS incident_reported_by_first_name_ar,
            reported_user.last_name_en AS incident_reported_by_last_name_en,
            reported_user.last_name_ar AS incident_reported_by_last_name_ar
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
		Where("incident_compliance_requirements.compliance_requirement_id IN (?)", complianceRequirementIDs)

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

	if err := query.Find(&incidentList).Error; err != nil {
		return nil, err
	}

	return incidentList, nil
}
