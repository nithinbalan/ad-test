package repositories

import (
	"fmt"
	reportDTOs "jethur/grc/reports/dtos"

	"gorm.io/gorm"
)

type ProjectReportRepositoryInterface interface {
	GetProjectByType(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectByTypeDTO, error)
	GetProjectByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectByStatusDTO, error)
	GetProjectByPhase(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectByPhaseDTO, error)
	GetProjectByClassification(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectByClassificationDTO, error)
	GetProjectByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectByResponsibleUnitDTO, error)
	GetProjectByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectByResponsibleUserDTO, error)
	GetProjectByYear(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectByYearDTO, error)
	GetProjectActionPlanByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectActionPlanByStatusDTO, error)
	GetProjectActionPlanByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectActionPlanByResponsibleUserDTO, error)
	GetProjectActionPlanByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectActionPlanByResponsibleUnitDTO, error)
}

type ProjectReportRepository struct {
	db *gorm.DB
}

func NewProjectReportRepository(db *gorm.DB) *ProjectReportRepository {
	return &ProjectReportRepository{db: db}
}

func (repo *ProjectReportRepository) GetProjectByType(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectByTypeDTO, error) {
	results := []*reportDTOs.ProjectByTypeDTO{}

	query := repo.db.Table("project_types").
		Select(`project_types.id,
				project_types.title_en,
				project_types.title_ar,
				COUNT(projects.id) AS count`).
		Joins("LEFT JOIN projects ON projects.project_type_id = project_types.id").
		Group("project_types.id, project_types.title_en, project_types.title_ar").
		Order("COUNT(projects.id) DESC")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("projects.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *ProjectReportRepository) GetProjectByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectByStatusDTO, error) {
	results := []*reportDTOs.ProjectByStatusDTO{}

	query := repo.db.Table("project_statuses").
		Select(`project_statuses.id,
				project_statuses.title_en,
				project_statuses.title_ar,
				project_statuses.color_code,
				COUNT(projects.id) AS count`).
		Joins("LEFT JOIN projects ON projects.project_status_id = project_statuses.id").
		Group("project_statuses.id, project_statuses.title_en, project_statuses.title_ar, project_statuses.color_code")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("projects.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *ProjectReportRepository) GetProjectByPhase(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectByPhaseDTO, error) {
	results := []*reportDTOs.ProjectByPhaseDTO{}

	query := repo.db.Table("project_phases").
		Select(`project_phases.id,
				project_phases.title_en,
				project_phases.title_ar,
				project_phases.color_code,
				COUNT(projects.id) AS count`).
		Joins("LEFT JOIN projects ON projects.project_phase_id = project_phases.id").
		Group("project_phases.id, project_phases.title_en, project_phases.title_ar, project_phases.color_code").
		Order("COUNT(projects.id) DESC")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("projects.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *ProjectReportRepository) GetProjectByClassification(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectByClassificationDTO, error) {
	results := []*reportDTOs.ProjectByClassificationDTO{}

	query := repo.db.Table("project_classifications").
		Select(`project_classifications.id,
				project_classifications.title_en,
				project_classifications.title_ar,
				COUNT(projects.id) AS count`).
		Joins("LEFT JOIN projects ON projects.project_classification_id = project_classifications.id").
		Group("project_classifications.id, project_classifications.title_en, project_classifications.title_ar").
		Order("COUNT(projects.id) DESC")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("projects.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *ProjectReportRepository) GetProjectByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectByResponsibleUnitDTO, error) {
	results := []*reportDTOs.ProjectByResponsibleUnitDTO{}

	query := repo.db.Table("projects").
		Select(`
			organization_units.id AS responsible_unit_id,
			organization_units.title_en AS responsible_unit_title_en,
			organization_units.title_ar AS responsible_unit_title_ar,
			organization_units.color_code AS responsible_unit_color_code,
			COUNT(projects.id) AS count
		`).
		Joins("JOIN users ON users.id = projects.responsible_user_id").
		Joins("JOIN organization_units ON organization_units.id = users.organization_unit_1_id").
		Group("organization_units.id, organization_units.title_en, organization_units.title_ar, organization_units.color_code").
		Order("COUNT(projects.id) DESC")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("projects.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *ProjectReportRepository) GetProjectByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectByResponsibleUserDTO, error) {
	results := []*reportDTOs.ProjectByResponsibleUserDTO{}

	query := repo.db.Table("projects").
		Select(`
			users.id AS responsible_user_id,
			users.first_name_en AS responsible_user_first_name_en,
			users.first_name_ar AS responsible_user_first_name_ar,
			users.last_name_en AS responsible_user_last_name_en,
			users.last_name_ar AS responsible_user_last_name_ar,
			CONCAT(users.first_name_en, ' ', users.last_name_en) AS responsible_user_title_en,
			CONCAT(users.first_name_ar, ' ', users.last_name_ar) AS responsible_user_title_ar,
			users.initial AS responsible_user_initial,
			users.color AS responsible_user_color,
			users.email AS responsible_user_email,
			users.profile_image_path AS responsible_user_profile_image_path,
			COUNT(projects.id) AS count
		`).
		Joins("JOIN users ON users.id = projects.responsible_user_id").
		Group("users.id, users.first_name_en, users.first_name_ar, users.last_name_en, users.last_name_ar, users.initial, users.color, users.email, users.profile_image_path").
		Order("COUNT(projects.id) DESC")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("projects.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *ProjectReportRepository) GetProjectByYear(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectByYearDTO, error) {
	results := []*reportDTOs.ProjectByYearDTO{}

	var rawResults []struct {
		Year  int  `gorm:"column:year"`
		Count uint `gorm:"column:count"`
	}

	query := repo.db.Table("projects").
		Select(`YEAR(start_date) AS year, COUNT(id) AS count`).
		Group("YEAR(start_date)").
		Order("YEAR(start_date) ASC")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("start_date BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&rawResults).Error
	if err != nil {
		return nil, err
	}

	for _, record := range rawResults {
		results = append(results, &reportDTOs.ProjectByYearDTO{
			TitleEn: fmt.Sprintf("%d", record.Year),
			TitleAr: fmt.Sprintf("%d", record.Year),
			Count:   record.Count,
		})
	}

	return results, nil
}

func (repo *ProjectReportRepository) GetProjectActionPlanByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectActionPlanByStatusDTO, error) {
	results := []*reportDTOs.ProjectActionPlanByStatusDTO{}

	query := repo.db.Table("action_plan_statuses").
		Select(`
			action_plan_statuses.id,
			action_plan_statuses.title_en,
			action_plan_statuses.title_ar,
			action_plan_statuses.color_code,
			COUNT(project_action_plans.id) AS count
		`).
		Joins("LEFT JOIN project_action_plans ON project_action_plans.action_plan_status_id = action_plan_statuses.id").
		Group("action_plan_statuses.id, action_plan_statuses.title_en, action_plan_statuses.title_ar, action_plan_statuses.color_code")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("project_action_plans.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *ProjectReportRepository) GetProjectActionPlanByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectActionPlanByResponsibleUserDTO, error) {
	results := []*reportDTOs.ProjectActionPlanByResponsibleUserDTO{}

	query := repo.db.Table("project_action_plans").
		Select(`
			users.id AS responsible_user_id,
			users.first_name_en AS responsible_user_first_name_en,
			users.first_name_ar AS responsible_user_first_name_ar,
			users.last_name_en AS responsible_user_last_name_en,
			users.last_name_ar AS responsible_user_last_name_ar,
			CONCAT(users.first_name_en, ' ', users.last_name_en) AS responsible_user_title_en,
			CONCAT(users.first_name_ar, ' ', users.last_name_ar) AS responsible_user_title_ar,
			users.initial AS responsible_user_initial,
			users.color AS responsible_user_color,
			users.email AS responsible_user_email,
			users.profile_image_path AS responsible_user_profile_image_path,
			COUNT(project_action_plans.id) AS count
		`).
		Joins("JOIN users ON users.id = project_action_plans.responsible_user_id").
		Group("users.id, users.first_name_en, users.first_name_ar, users.last_name_en, users.last_name_ar, users.initial, users.color, users.email, users.profile_image_path").
		Order("COUNT(project_action_plans.id) DESC")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("project_action_plans.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *ProjectReportRepository) GetProjectActionPlanByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectActionPlanByResponsibleUnitDTO, error) {
	results := []*reportDTOs.ProjectActionPlanByResponsibleUnitDTO{}

	query := repo.db.Table("project_action_plans").
		Select(`
			organization_units.id AS responsible_unit_id,
			organization_units.title_en AS responsible_unit_title_en,
			organization_units.title_ar AS responsible_unit_title_ar,
			organization_units.color_code AS color_code,
			COUNT(project_action_plans.id) AS count
		`).
		Joins("JOIN users ON users.id = project_action_plans.responsible_user_id").
		Joins("JOIN organization_units ON organization_units.id = users.organization_unit_1_id").
		Group("organization_units.id, organization_units.title_en, organization_units.title_ar, organization_units.color_code").
		Order("COUNT(project_action_plans.id) DESC")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("project_action_plans.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}
