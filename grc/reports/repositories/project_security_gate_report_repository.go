package repositories

import (
	reportDTOs "jethur/grc/reports/dtos"

	"gorm.io/gorm"
)

type ProjectSecurityGateReportRepositoryInterface interface {
	GetProjectSecurityGateByStage(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectSecurityGateByStageDTO, error)
	GetProjectSecurityGateByReadinessStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectSecurityGateByReadinessStatusDTO, error)
	GetProjectSecurityGateByClassification(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectSecurityGateByClassificationDTO, error)
}

type ProjectSecurityGateReportRepository struct {
	db *gorm.DB
}

func NewProjectSecurityGateReportRepository(db *gorm.DB) *ProjectSecurityGateReportRepository {
	return &ProjectSecurityGateReportRepository{db: db}
}

func (repo *ProjectSecurityGateReportRepository) GetProjectSecurityGateByStage(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectSecurityGateByStageDTO, error) {
	results := []*reportDTOs.ProjectSecurityGateByStageDTO{}

	query := repo.db.Table("project_security_gate_stages").
		Select(`project_security_gate_stages.id,
				project_security_gate_stages.title_en,
				project_security_gate_stages.title_ar,
				project_security_gate_stages.color_code,
				COUNT(project_security_gates.id) AS count`).
		Joins("LEFT JOIN project_security_gates ON project_security_gates.stage_id = project_security_gate_stages.id").
		Group("project_security_gate_stages.id, project_security_gate_stages.title_en, project_security_gate_stages.title_ar, project_security_gate_stages.color_code")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("project_security_gates.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *ProjectSecurityGateReportRepository) GetProjectSecurityGateByReadinessStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectSecurityGateByReadinessStatusDTO, error) {
	results := []*reportDTOs.ProjectSecurityGateByReadinessStatusDTO{}

	query := repo.db.Table("project_security_gate_readiness_statuses").
		Select(`project_security_gate_readiness_statuses.id,
				project_security_gate_readiness_statuses.title_en,
				project_security_gate_readiness_statuses.title_ar,
				project_security_gate_readiness_statuses.color_code,
				COUNT(project_security_gates.id) AS count`).
		Joins("LEFT JOIN project_security_gates ON project_security_gates.readiness_status_id = project_security_gate_readiness_statuses.id").
		Group("project_security_gate_readiness_statuses.id, project_security_gate_readiness_statuses.title_en, project_security_gate_readiness_statuses.title_ar, project_security_gate_readiness_statuses.color_code")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("project_security_gates.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *ProjectSecurityGateReportRepository) GetProjectSecurityGateByClassification(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectSecurityGateByClassificationDTO, error) {
	results := []*reportDTOs.ProjectSecurityGateByClassificationDTO{}

	// Counted through the mapping table, a gate can carry several classifications
	query := repo.db.Table("project_security_gate_classifications").
		Select(`project_security_gate_classifications.id,
				project_security_gate_classifications.title_en,
				project_security_gate_classifications.title_ar,
				project_security_gate_classifications.color_code,
				COUNT(project_security_gates.id) AS count`).
		Joins(`LEFT JOIN project_security_gate_classification_mappings
					ON project_security_gate_classification_mappings.classification_id = project_security_gate_classifications.id`).
		Joins(`LEFT JOIN project_security_gates
					ON project_security_gates.id = project_security_gate_classification_mappings.security_gate_id`).
		Group("project_security_gate_classifications.id, project_security_gate_classifications.title_en, project_security_gate_classifications.title_ar, project_security_gate_classifications.color_code")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("project_security_gates.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}
