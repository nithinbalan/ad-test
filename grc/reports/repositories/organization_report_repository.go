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

	GetControlsByControlType(params *dtos.ReportListQueryParamsDTO) ([]*dtos.ControlsByControlTypeDTO, error)
	GetControlsByCategory(params *dtos.ReportListQueryParamsDTO) ([]*dtos.ControlsByControlCategoryDTO, error)
	GetControlsByComplianceRegister(params *dtos.ReportListQueryParamsDTO) ([]*dtos.ControlsByComplianceRegisterDTO, error)
	GetControlsByRiskRating(params *dtos.ReportListQueryParamsDTO) ([]*dtos.ControlsByRiskRatingDTO, error)
	GetControlsByImplementationStatus(params *dtos.ReportListQueryParamsDTO) ([]*dtos.ControlsByImplementationStatusDTO, error)
	GetControlsByEfficiencyStatus(params *dtos.ReportListQueryParamsDTO) ([]*dtos.ControlsByEfficiencyDTO, error)
	GetControlsByApplicability(params *dtos.ReportListQueryParamsDTO) ([]*dtos.ControlsByApplicabilityDTO, error)
}

type OrganizationReportRepository struct {
	db *gorm.DB
}

func NewOrganizationReportRepository(db *gorm.DB) *OrganizationReportRepository {
	return &OrganizationReportRepository{db: db}
}

func (repo *OrganizationReportRepository) GetUsersByStatus(params *dtos.ReportListQueryParamsDTO) ([]*dtos.UsersByStatusDTO, error) {
	results := []*dtos.UsersByStatusDTO{}

	query := repo.db.Table("user_statuses").
		Select(`
			user_statuses.id AS user_status_id,
			user_statuses.title_en AS user_status_title_en,
			user_statuses.title_ar AS user_status_title_ar,
			user_statuses.color_code AS user_status_color_code,
			user_statuses.type AS user_status_type,
			COALESCE(COUNT(users.id), 0) AS count`).
		Joins("LEFT JOIN users ON users.user_status_id = user_statuses.id").
		Group("user_statuses.id, user_statuses.title_en, user_statuses.title_ar, user_statuses.color_code, user_statuses.type")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("users.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *OrganizationReportRepository) GetUsersByRole(params *dtos.ReportListQueryParamsDTO) ([]*dtos.UsersByRoleDTO, error) {
	results := []*dtos.UsersByRoleDTO{}

	query := repo.db.Table("roles").
		Select(`
			roles.id AS role_id,
			roles.title_en AS role_title_en,
			roles.title_ar AS role_title_ar,
			COALESCE(COUNT(user_roles.user_id), 0) AS count`).
		Joins("LEFT JOIN user_roles ON user_roles.role_id = roles.id").
		Joins("LEFT JOIN users ON users.id = user_roles.user_id").
		Group("roles.id, roles.title_en, roles.title_ar")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("users.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *OrganizationReportRepository) GetUsersByDesignation(params *dtos.ReportListQueryParamsDTO) ([]*dtos.UsersByDesignationDTO, error) {
	results := []*dtos.UsersByDesignationDTO{}

	query := repo.db.Table("designations").
		Select(`
			designations.id AS designation_id,
			designations.title_en AS designation_title_en,
			designations.title_ar AS designation_title_ar,
			COALESCE(COUNT(users.id), 0) AS count`).
		Joins("LEFT JOIN users ON users.designation_id = designations.id").
		Group("designations.id, designations.title_en, designations.title_ar")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("users.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *OrganizationReportRepository) GetProcessesByResponsibleUnit(params *dtos.ReportListQueryParamsDTO) ([]*dtos.ProcessByUnitDTO, error) {
	results := []*dtos.ProcessByUnitDTO{}

	query := repo.db.Table("organization_units").
		Select(`
			organization_units.id AS owner_unit_id,
			organization_units.title_en AS owner_unit_title_en,
			organization_units.title_ar AS owner_unit_title_ar,
			organization_units.color_code AS owner_unit_color_code,
			organization_units.status AS owner_unit_status,
			COALESCE(COUNT(processes.id), 0) AS count`).
		Joins("LEFT JOIN processes ON processes.owner_unit_id = organization_units.id").
		Group("organization_units.id, organization_units.title_en, organization_units.title_ar, organization_units.color_code, organization_units.status")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("processes.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *OrganizationReportRepository) GetProcessesByRiskRating(params *dtos.ReportListQueryParamsDTO) ([]*dtos.ProcessByRiskRatingDTO, error) {
	results := []*dtos.ProcessByRiskRatingDTO{}

	query := repo.db.Table("risk_matrix_risk_ratings").
		Select(`
			risk_matrix_risk_ratings.id AS risk_rating_id,
			risk_matrix_risk_ratings.title_en AS risk_rating_title_en,
			risk_matrix_risk_ratings.title_ar AS risk_rating_title_ar,
			risk_matrix_risk_ratings.color_code AS risk_rating_color_code,
			risk_matrix_risk_ratings.type AS risk_rating_type,
			COALESCE(COUNT(processes.id), 0) AS count`).
		Joins("LEFT JOIN processes ON processes.risk_matrix_risk_rating_id = risk_matrix_risk_ratings.id").
		Group("risk_matrix_risk_ratings.id, risk_matrix_risk_ratings.title_en, risk_matrix_risk_ratings.title_ar, risk_matrix_risk_ratings.color_code, risk_matrix_risk_ratings.type")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("processes.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *OrganizationReportRepository) GetProcessesByResponsibleUser(params *dtos.ReportListQueryParamsDTO) ([]*dtos.ProcessByResponsibleUserDTO, error) {
	results := []*dtos.ProcessByResponsibleUserDTO{}

	query := repo.db.Table("users").
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
			COALESCE(COUNT(processes.id), 0) AS count`).
		Joins("LEFT JOIN processes ON processes.owner_id = users.id").
		Group("users.id, users.first_name_en, users.first_name_ar, users.last_name_en, users.last_name_ar, users.initial, users.color, users.email, users.profile_image_path")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("processes.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *OrganizationReportRepository) GetProcessesByProcessGroup(params *dtos.ReportListQueryParamsDTO) ([]*dtos.ProcessByProcessGroupDTO, error) {
	results := []*dtos.ProcessByProcessGroupDTO{}

	query := repo.db.Table("process_groups").
		Select(`
			process_groups.id AS process_group_id,
			process_groups.title_en AS process_group_title_en,
			process_groups.title_ar AS process_group_title_ar,
			COALESCE(COUNT(processes.id), 0) AS count`).
		Joins("LEFT JOIN processes ON processes.process_group_id = process_groups.id").
		Group("process_groups.id, process_groups.title_en, process_groups.title_ar")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("processes.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *OrganizationReportRepository) GetControlsByControlType(params *dtos.ReportListQueryParamsDTO) ([]*dtos.ControlsByControlTypeDTO, error) {
	results := []*dtos.ControlsByControlTypeDTO{}

	query := repo.db.Table("control_types").
		Select(`
			control_types.id AS control_type_id,
			control_types.title_en AS control_type_title_en,
			control_types.title_ar AS control_type_title_ar,
			COALESCE(COUNT(controls.id), 0) AS count`).
		Joins("LEFT JOIN controls ON controls.control_type_id = control_types.id").
		Group("control_types.id, control_types.title_en, control_types.title_ar")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("controls.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *OrganizationReportRepository) GetControlsByCategory(params *dtos.ReportListQueryParamsDTO) ([]*dtos.ControlsByControlCategoryDTO, error) {
	results := []*dtos.ControlsByControlCategoryDTO{}

	query := repo.db.Table("control_categories").
		Select(`
			control_categories.id AS control_category_id,
			control_categories.title_en AS control_category_title_en,
			control_categories.title_ar AS control_category_title_ar,
			COALESCE(COUNT(controls.id), 0) AS count`).
		Joins("LEFT JOIN controls ON controls.control_category_id = control_categories.id").
		Group("control_categories.id, control_categories.title_en, control_categories.title_ar")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("controls.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *OrganizationReportRepository) GetControlsByComplianceRegister(params *dtos.ReportListQueryParamsDTO) ([]*dtos.ControlsByComplianceRegisterDTO, error) {
	results := []*dtos.ControlsByComplianceRegisterDTO{}

	query := repo.db.Table("compliance_register").
		Select(`
			compliance_register.id AS compliance_register_id,
			compliance_register.title_en AS compliance_register_title_en,
			compliance_register.title_ar AS compliance_register_title_ar,
			COALESCE(COUNT(compliance_controls.control_id), 0) AS count`).
		Joins("LEFT JOIN compliance_controls ON compliance_controls.compliance_register_id = compliance_register.id").
		Joins("LEFT JOIN controls ON controls.id = compliance_controls.control_id").
		Group("compliance_register.id, compliance_register.title_en, compliance_register.title_ar")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("controls.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *OrganizationReportRepository) GetControlsByRiskRating(params *dtos.ReportListQueryParamsDTO) ([]*dtos.ControlsByRiskRatingDTO, error) {
	results := []*dtos.ControlsByRiskRatingDTO{}

	query := repo.db.Table("risk_matrix_risk_ratings").
		Select(`
			risk_matrix_risk_ratings.id AS risk_rating_id,
			risk_matrix_risk_ratings.title_en AS risk_rating_title_en,
			risk_matrix_risk_ratings.title_ar AS risk_rating_title_ar,
			risk_matrix_risk_ratings.color_code AS risk_rating_color_code,
			risk_matrix_risk_ratings.type AS risk_rating_type,
			COALESCE(COUNT(controls.id), 0) AS count`).
		Joins("LEFT JOIN controls ON controls.risk_matrix_risk_rating_id = risk_matrix_risk_ratings.id").
		Group("risk_matrix_risk_ratings.id, risk_matrix_risk_ratings.title_en, risk_matrix_risk_ratings.title_ar, risk_matrix_risk_ratings.color_code, risk_matrix_risk_ratings.type")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("controls.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *OrganizationReportRepository) GetControlsByImplementationStatus(params *dtos.ReportListQueryParamsDTO) ([]*dtos.ControlsByImplementationStatusDTO, error) {
	results := []*dtos.ControlsByImplementationStatusDTO{}

	query := repo.db.Table("control_implementation_statuses").
		Select(`
			control_implementation_statuses.id AS control_implementation_status_id,
			control_implementation_statuses.title_en AS control_implementation_status_title_en,
			control_implementation_statuses.title_ar AS control_implementation_status_title_ar,
			control_implementation_statuses.color_code AS control_implementation_status_color_code,
			control_implementation_statuses.type AS control_implementation_status_type,
			COALESCE(COUNT(controls.id), 0) AS count`).
		Joins("LEFT JOIN controls ON controls.control_implementation_status_id = control_implementation_statuses.id").
		Group("control_implementation_statuses.id, control_implementation_statuses.title_en, control_implementation_statuses.title_ar, control_implementation_statuses.color_code, control_implementation_statuses.type")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("controls.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *OrganizationReportRepository) GetControlsByEfficiencyStatus(params *dtos.ReportListQueryParamsDTO) ([]*dtos.ControlsByEfficiencyDTO, error) {
	results := []*dtos.ControlsByEfficiencyDTO{}

	query := repo.db.Table("control_efficiency_statuses").
		Select(`
			control_efficiency_statuses.id AS control_efficiency_status_id,
			control_efficiency_statuses.title_en AS control_efficiency_status_title_en,
			control_efficiency_statuses.title_ar AS control_efficiency_status_title_ar,
			control_efficiency_statuses.color_code AS control_efficiency_status_color_code,
			control_efficiency_statuses.type AS control_efficiency_status_type,
			COALESCE(COUNT(controls.id), 0) AS count`).
		Joins("LEFT JOIN controls ON controls.control_efficiency_status_id = control_efficiency_statuses.id").
		Group("control_efficiency_statuses.id, control_efficiency_statuses.title_en, control_efficiency_statuses.title_ar, control_efficiency_statuses.color_code, control_efficiency_statuses.type")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("controls.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *OrganizationReportRepository) GetControlsByApplicability(params *dtos.ReportListQueryParamsDTO) ([]*dtos.ControlsByApplicabilityDTO, error) {
	results := []*dtos.ControlsByApplicabilityDTO{}

	query := repo.db.Table("controls").
		Select(`
			controls.is_applicable AS is_applicable,
			COALESCE(COUNT(controls.id), 0) AS count`).
		Group("controls.is_applicable")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("controls.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}
