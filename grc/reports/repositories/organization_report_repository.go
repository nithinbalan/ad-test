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
			user_statuses.id AS id,
			user_statuses.title_en AS user_status_title_en,
			user_statuses.title_ar AS user_status_title_ar,
			user_statuses.color_code AS user_status_color_code,
			user_statuses.type AS user_status_type,
			COUNT(users.id) AS count`).
		Joins("JOIN user_statuses ON user_statuses.id = users.user_status_id").
		Group("user_statuses.id, user_statuses.title_en, user_statuses.title_ar, user_statuses.color_code, user_statuses.type")

	// Apply date filters if provided
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
			roles.id AS id,
			roles.title_en AS role_title_en,
			roles.title_ar AS role_title_ar,
			COUNT(users.id) AS count`).
		Joins("JOIN users ON users.id = user_roles.user_id").
		Joins("JOIN roles ON roles.id = user_roles.role_id").
		Group("roles.id, roles.title_en, roles.title_ar")

	// Apply date filters if provided
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
			designations.id AS id,
			designations.title_en AS designation_title_en,
			designations.title_ar AS designation_title_ar,
			COUNT(users.id) AS count`).
		Joins("JOIN designations ON designations.id = users.designation_id").
		Group("designations.id, designations.title_en, designations.title_ar")

	// Apply date filters if provided
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
			COUNT(processes.id) AS count
		`).
		Joins("JOIN organization_units ON organization_units.id = processes.owner_unit_id").
		Group("processes.owner_unit_id, organization_units.title_en, organization_units.title_ar, organization_units.color_code, organization_units.status")

	// Apply date filters if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("processes.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}
