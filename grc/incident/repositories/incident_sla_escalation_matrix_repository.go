package repositories

import (
	"jethur/grc/incident/dtos"
	"jethur/grc/models"

	"gorm.io/gorm"
)

type IncidentSLAEscalationMatrixRepositoryInterface interface {
	GetAllEscalationMatrices(sort, sortBy, keyword, status string) ([]*dtos.IncidentEscalationMatrixListDTO, error)
	GetAllResolutionMatrices(sort, sortBy, keyword, status string) ([]*dtos.IncidentResolutionMatrixListDTO, error)
	GetBySeverityID(severityID uint, matrixType string) (*models.IncidentSLAEscalationMatrix, error)
	UpdateEscalationMatrix(matrix *models.IncidentSLAEscalationMatrix) error
	UpdateResolutionMatrix(matrix *models.IncidentSLAEscalationMatrix) error
}

type IncidentSLAEscalationMatrixRepository struct {
	db *gorm.DB
}

func NewIncidentSLAEscalationMatrixRepository(db *gorm.DB) *IncidentSLAEscalationMatrixRepository {
	return &IncidentSLAEscalationMatrixRepository{db: db}
}

func (repo *IncidentSLAEscalationMatrixRepository) GetAllEscalationMatrices(sort, sortBy, keyword, status string) ([]*dtos.IncidentEscalationMatrixListDTO, error) {
	var matrices []*dtos.IncidentEscalationMatrixListDTO

	query := repo.db.
		Table("incident_sla_escalation_matrix").
		Select(`
			incident_sla_escalation_matrix.id,
			incident_sla_escalation_matrix.severity_id,
			severity_levels.title_en AS severity_level_title_en,
			severity_levels.title_ar AS severity_level_title_ar,
			severity_levels.type AS severity_level_type,
			severity_levels.color_code AS severity_level_color_code,
			severity_levels.weight AS severity_level_weight,
			incident_sla_escalation_matrix.notify_within,
			incident_sla_escalation_matrix.notify_within_unit_id,
			notify_within_units.title_en AS notify_within_unit_title_en,
			notify_within_units.title_ar AS notify_within_unit_title_ar,
			incident_sla_escalation_matrix.notify_to_incident_owner,
			incident_sla_escalation_matrix.notify_to_response_team,
			incident_sla_escalation_matrix.notify_role_id,
			notify_roles.title_en AS notify_role_title_en,
			notify_roles.title_ar AS notify_role_title_ar,
			incident_sla_escalation_matrix.escalate_after_incident_owner,
			incident_sla_escalation_matrix.escalate_after_response_team,
			incident_sla_escalation_matrix.escalate_role_id,
			escalate_roles.title_en AS escalate_role_title_en,
			escalate_roles.title_ar AS escalate_role_title_ar,
			incident_sla_escalation_matrix.status
		`).
		Joins("LEFT JOIN severity_levels ON severity_levels.id = incident_sla_escalation_matrix.severity_id").
		Joins("LEFT JOIN business_continuity_duration_units AS notify_within_units ON notify_within_units.id = incident_sla_escalation_matrix.notify_within_unit_id").
		Joins("LEFT JOIN roles AS notify_roles ON notify_roles.id = incident_sla_escalation_matrix.notify_role_id").
		Joins("LEFT JOIN roles AS escalate_roles ON escalate_roles.id = incident_sla_escalation_matrix.escalate_role_id").
		Where("incident_sla_escalation_matrix.type = ?", "escalation_matrix")

	if keyword != "" {
		keyword = "%" + keyword + "%"

		query = query.Where(`
			incident_sla_escalation_matrix.status LIKE ? OR
			severity_levels.title_en LIKE ? OR
			severity_levels.title_ar LIKE ?
		`,
			keyword,
			keyword,
			keyword,
		)
	}

	if status == "active" || status == "inactive" {
		query = query.Where("incident_sla_escalation_matrix.status = ?", status)
	}

	if sortBy == "" {
		sortBy = "incident_sla_escalation_matrix.id"
	}

	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	query = query.Order(sortBy + " " + sort)

	if err := query.Find(&matrices).Error; err != nil {
		return nil, err
	}

	return matrices, nil
}

func (repo *IncidentSLAEscalationMatrixRepository) GetAllResolutionMatrices(sort, sortBy, keyword, status string) ([]*dtos.IncidentResolutionMatrixListDTO, error) {
	var matrices []*dtos.IncidentResolutionMatrixListDTO

	query := repo.db.
		Table("incident_sla_escalation_matrix").
		Select(`
			incident_sla_escalation_matrix.id,
			incident_sla_escalation_matrix.severity_id,
			severity_levels.title_en AS severity_level_title_en,
			severity_levels.title_ar AS severity_level_title_ar,
			severity_levels.type AS severity_level_type,
			severity_levels.color_code AS severity_level_color_code,
			severity_levels.weight AS severity_level_weight,
			incident_sla_escalation_matrix.close_within,
			incident_sla_escalation_matrix.close_within_unit_id,
			close_within_units.title_en AS close_within_unit_title_en,
			close_within_units.title_ar AS close_within_unit_title_ar,
			incident_sla_escalation_matrix.close_within_incident_owner,
			incident_sla_escalation_matrix.close_within_response_team,
			incident_sla_escalation_matrix.close_within_role_id,
			close_within_roles.title_en AS close_within_role_title_en,
			close_within_roles.title_ar AS close_within_role_title_ar,
			incident_sla_escalation_matrix.status
		`).
		Joins("LEFT JOIN severity_levels ON severity_levels.id = incident_sla_escalation_matrix.severity_id").
		Joins("LEFT JOIN business_continuity_duration_units AS close_within_units ON close_within_units.id = incident_sla_escalation_matrix.close_within_unit_id").
		Joins("LEFT JOIN roles AS close_within_roles ON close_within_roles.id = incident_sla_escalation_matrix.close_within_role_id").
		Where("incident_sla_escalation_matrix.type = ?", "resolution_matrix")

	if keyword != "" {
		keyword = "%" + keyword + "%"

		query = query.Where(`
			incident_sla_escalation_matrix.status LIKE ? OR
			severity_levels.title_en LIKE ? OR
			severity_levels.title_ar LIKE ?
		`,
			keyword,
			keyword,
			keyword,
		)
	}

	if status == "active" || status == "inactive" {
		query = query.Where("incident_sla_escalation_matrix.status = ?", status)
	}

	if sortBy == "" {
		sortBy = "incident_sla_escalation_matrix.id"
	}

	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	query = query.Order(sortBy + " " + sort)

	if err := query.Find(&matrices).Error; err != nil {
		return nil, err
	}

	return matrices, nil
}

func (repo *IncidentSLAEscalationMatrixRepository) GetBySeverityID(severityID uint, matrixType string) (*models.IncidentSLAEscalationMatrix, error) {
	var matrix models.IncidentSLAEscalationMatrix
	result := repo.db.Where("severity_id = ? AND type = ?", severityID, matrixType).First(&matrix)
	return &matrix, result.Error
}

func (repo *IncidentSLAEscalationMatrixRepository) UpdateEscalationMatrix(matrix *models.IncidentSLAEscalationMatrix) error {
	if err := repo.db.Save(matrix).Error; err != nil {
		return err
	}
	return nil
}

func (repo *IncidentSLAEscalationMatrixRepository) UpdateResolutionMatrix(matrix *models.IncidentSLAEscalationMatrix) error {
	if err := repo.db.Save(matrix).Error; err != nil {
		return err
	}
	return nil
}
