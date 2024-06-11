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
            risk_register.title AS risk_register_title,
            risk_register.description AS risk_register_description,
            risk_statuses.title_en AS risk_status_title_en,
            risk_statuses.title_ar AS risk_status_title_ar,
            risk_statuses.color_code AS risk_status_color_code
        `).
		Joins(`
            LEFT JOIN risk_register ON risk_register.id = risk_compliance_requirements.risk_register_id
            LEFT JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id
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
            incident_register.title AS incident_register_title,
            incident_register.description AS incident_register_description,
            incident_statuses.title_en AS incident_status_title_en,
            incident_statuses.title_ar AS incident_status_title_ar,
            incident_statuses.color_code AS incident_status_color_code
        `).
		Joins(`
            LEFT JOIN incident_register ON incident_register.id = incident_compliance_requirements.incident_register_id
            LEFT JOIN incident_statuses ON incident_statuses.id = incident_register.incident_status_id
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
