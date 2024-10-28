package repositories

import (
	"jethur/grc/reports/dtos"

	"gorm.io/gorm"
)

type OrganizationReportRepositoryInterface interface {
	GetUsersByStatus(params *dtos.ReportListQueryParamsDTO) ([]*dtos.UsersByStatusDTO, error)
	GetUsersByRole(params *dtos.ReportListQueryParamsDTO) ([]*dtos.UsersByRoleDTO, error)
	GetUsersByDesignation(params *dtos.ReportListQueryParamsDTO) ([]*dtos.UsersByDesignationDTO, error)

	GetProcessesByResponsibleUnit(params *dtos.ReportListQueryParamsDTO) ([]*dtos.ProcessByUnitDTO, error)
	GetProcessesByRiskRating(params *dtos.ReportListQueryParamsDTO) ([]*dtos.ProcessByRiskRatingDTO, error)
	GetProcessesByResponsibleUser(params *dtos.ReportListQueryParamsDTO) ([]*dtos.ProcessByResponsibleUserDTO, error)
	GetProcessesByProcessGroup(params *dtos.ReportListQueryParamsDTO) ([]*dtos.ProcessByProcessGroupDTO, error)
}

type OrganizationReportRepository struct {
	db *gorm.DB
}

func NewOrganizationReportRepository(db *gorm.DB) *OrganizationReportRepository {
	return &OrganizationReportRepository{db: db}
}

func (repo *OrganizationReportRepository) GetUsersByStatus(params *dtos.ReportListQueryParamsDTO) ([]*dtos.UsersByStatusDTO, error) {
	results := []*dtos.UsersByStatusDTO{}

	query := repo.db.Table("users").
		Select(`
			user_statuses.id AS user_status_id,
			user_statuses.title_en AS user_status_title_en,
			user_statuses.title_ar AS user_status_title_ar,
			user_statuses.color_code AS user_status_color_code,
			user_statuses.type AS user_status_type,
			COUNT(users.id) AS count`).
		Joins("JOIN user_statuses ON user_statuses.id = users.user_status_id").
		Group("user_statuses.id, user_statuses.title_en, user_statuses.title_ar, user_statuses.color_code, user_statuses.type")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("users.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *OrganizationReportRepository) GetUsersByRole(params *dtos.ReportListQueryParamsDTO) ([]*dtos.UsersByRoleDTO, error) {
	results := []*dtos.UsersByRoleDTO{}

	query := repo.db.Table("user_roles").
		Select(`
			roles.id AS role_id,
			roles.title_en AS role_title_en,
			roles.title_ar AS role_title_ar,
			COUNT(users.id) AS count`).
		Joins("JOIN users ON users.id = user_roles.user_id").
		Joins("JOIN roles ON roles.id = user_roles.role_id").
		Group("roles.id, roles.title_en, roles.title_ar")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("users.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *OrganizationReportRepository) GetUsersByDesignation(params *dtos.ReportListQueryParamsDTO) ([]*dtos.UsersByDesignationDTO, error) {
	results := []*dtos.UsersByDesignationDTO{}

	query := repo.db.Table("users").
		Select(`
			designations.id AS designation_id,
			designations.title_en AS designation_title_en,
			designations.title_ar AS designation_title_ar,
			COUNT(users.id) AS count`).
		Joins("JOIN designations ON designations.id = users.designation_id").
		Group("designations.id, designations.title_en, designations.title_ar")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("users.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *OrganizationReportRepository) GetProcessesByResponsibleUnit(params *dtos.ReportListQueryParamsDTO) ([]*dtos.ProcessByUnitDTO, error) {
	results := []*dtos.ProcessByUnitDTO{}

	query := repo.db.Table("processes").
		Select(`
			processes.owner_unit_id AS owner_unit_id,
			organization_units.title_en AS owner_unit_title_en,
			organization_units.title_ar AS owner_unit_title_ar,
			organization_units.color_code AS owner_unit_color_code,
			organization_units.status AS owner_unit_status,
			COUNT(processes.id) AS count`).
		Joins("JOIN organization_units ON organization_units.id = processes.owner_unit_id").
		Group("processes.owner_unit_id, organization_units.title_en, organization_units.title_ar, organization_units.color_code, organization_units.status")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("processes.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *OrganizationReportRepository) GetProcessesByRiskRating(params *dtos.ReportListQueryParamsDTO) ([]*dtos.ProcessByRiskRatingDTO, error) {
	results := []*dtos.ProcessByRiskRatingDTO{}

	query := repo.db.Table("processes").
		Select(`
			risk_matrix_risk_ratings.id AS risk_rating_id,
			risk_matrix_risk_ratings.title_en AS risk_rating_title_en,
			risk_matrix_risk_ratings.title_ar AS risk_rating_title_ar,
			risk_matrix_risk_ratings.color_code AS risk_rating_color_code,
			risk_matrix_risk_ratings.type AS risk_rating_type,
			COUNT(processes.id) AS count`).
		Joins("JOIN risk_matrix_risk_ratings ON risk_matrix_risk_ratings.id = processes.risk_matrix_risk_rating_id").
		Group("risk_matrix_risk_ratings.id, risk_matrix_risk_ratings.title_en, risk_matrix_risk_ratings.title_ar, risk_matrix_risk_ratings.color_code, risk_matrix_risk_ratings.type")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("processes.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *OrganizationReportRepository) GetProcessesByResponsibleUser(params *dtos.ReportListQueryParamsDTO) ([]*dtos.ProcessByResponsibleUserDTO, error) {
	results := []*dtos.ProcessByResponsibleUserDTO{}

	query := repo.db.Table("processes").
		Select(`
			users.id AS responsible_user_id,
			users.first_name_en AS responsible_user_first_name_en,
			users.first_name_ar AS responsible_user_first_name_ar,
			users.last_name_en AS responsible_user_last_name_en,
			users.last_name_ar AS responsible_user_last_name_ar,
			users.initial AS responsible_user_initial,
			users.color AS responsible_user_color,
			users.email AS responsible_user_email,
			users.profile_image_path AS responsible_user_profile_image_path,
			COUNT(processes.id) AS count`).
		Joins("JOIN users ON users.id = processes.owner_id").
		Group("users.id, users.first_name_en, users.first_name_ar, users.last_name_en, users.last_name_ar, users.initial, users.color, users.email, users.profile_image_path")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("processes.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *OrganizationReportRepository) GetProcessesByProcessGroup(params *dtos.ReportListQueryParamsDTO) ([]*dtos.ProcessByProcessGroupDTO, error) {
	results := []*dtos.ProcessByProcessGroupDTO{}

	query := repo.db.Table("processes").
		Select(`
			process_groups.id AS process_group_id,
			process_groups.title_en AS process_group_title_en,
			process_groups.title_ar AS process_group_title_ar,
			COUNT(processes.id) AS count`).
		Joins("JOIN process_groups ON process_groups.id = processes.process_group_id").
		Group("process_groups.id, process_groups.title_en, process_groups.title_ar")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("processes.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}
