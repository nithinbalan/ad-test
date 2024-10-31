package repositories

import (
	reportDTOs "jethur/grc/reports/dtos"

	"gorm.io/gorm"
)

type ComplianceReportRepositoryInterface interface {
	GetComplianceByComplianceCategory(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ComplianceByCategoryDTO, error)
	GetComplianceBySeverityLevel(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ComplianceBySeverityLevelDTO, error)
	GetComplianceByComplianceStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ComplianceByStatusDTO, error)
	GetComplianceByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ComplianceByResponsibleUnitDTO, error)
	GetComplianceByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ComplianceByResponsibleUserDTO, error)

	GetComplianceRequirementByComplianceStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ComplianceRequirementByComplianceStatusDTO, error)
	GetComplianceRequirementBySeverity(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ComplianceRequirementBySeverityLevelDTO, error)
	GetComplianceRequirementByCategory(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ComplianceRequirementByCategoryDTO, error)
	GetComplianceRequirementByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ComplianceRequirementByResponsibleUserDTO, error)
}

type ComplianceReportRepository struct {
	db *gorm.DB
}

func NewComplianceReportRepository(db *gorm.DB) *ComplianceReportRepository {
	return &ComplianceReportRepository{
		db: db,
	}
}

func (repo *ComplianceReportRepository) GetComplianceByComplianceCategory(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ComplianceByCategoryDTO, error) {
	results := []*reportDTOs.ComplianceByCategoryDTO{}
	query := repo.db.Table("compliance_register").
		Select(`
			compliance_categories.id AS compliance_category_id,
			compliance_categories.title_en AS compliance_category_title_en,
			compliance_categories.title_ar AS compliance_category_title_ar,
			COUNT(compliance_register.id) AS count
		`).
		Joins("JOIN compliance_categories ON compliance_categories.id = compliance_register.compliance_category_id").
		Group("compliance_categories.id, compliance_categories.title_en, compliance_categories.title_ar")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("compliance_register.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *ComplianceReportRepository) GetComplianceBySeverityLevel(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ComplianceBySeverityLevelDTO, error) {
	results := []*reportDTOs.ComplianceBySeverityLevelDTO{}
	query := repo.db.Table("compliance_register").
		Select(`
			severity_levels.id AS severity_level_id,
			severity_levels.title_en AS severity_level_title_en,
			severity_levels.title_ar AS severity_level_title_ar,
			severity_levels.color_code AS severity_level_color_code,
			severity_levels.type AS severity_level_type,
			COUNT(compliance_register.id) AS count
		`).
		Joins("JOIN severity_levels ON severity_levels.id = compliance_register.severity_level_id").
		Group("severity_levels.id, severity_levels.title_en, severity_levels.title_ar, severity_levels.color_code, severity_levels.type")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("compliance_register.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *ComplianceReportRepository) GetComplianceByComplianceStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ComplianceByStatusDTO, error) {
	results := []*reportDTOs.ComplianceByStatusDTO{}
	query := repo.db.Table("compliance_register").
		Select(`
			compliance_statuses.id AS compliance_status_id,
			compliance_statuses.title_en AS compliance_status_title_en,
			compliance_statuses.title_ar AS compliance_status_title_ar,
			compliance_statuses.color_code AS compliance_status_color_code,
			compliance_statuses.type AS compliance_status_type,
			COUNT(compliance_register.id) AS count
		`).
		Joins("JOIN compliance_statuses ON compliance_statuses.id = compliance_register.compliance_status_id").
		Group("compliance_statuses.id, compliance_statuses.title_en, compliance_statuses.title_ar, compliance_statuses.color_code, compliance_statuses.type")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("compliance_register.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *ComplianceReportRepository) GetComplianceByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ComplianceByResponsibleUnitDTO, error) {
	results := []*reportDTOs.ComplianceByResponsibleUnitDTO{}
	query := repo.db.Table("compliance_register").
		Select(`
			organization_units.id AS responsible_unit_id,
			organization_units.title_en AS responsible_unit_title_en,
			organization_units.title_ar AS responsible_unit_title_ar,
			organization_units.color_code AS responsible_unit_color_code,
			COUNT(compliance_register.id) AS count
		`).
		Joins("JOIN organization_units ON organization_units.id = compliance_register.responsible_unit_id").
		Group("organization_units.id, organization_units.title_en, organization_units.title_ar, organization_units.color_code")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("compliance_register.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *ComplianceReportRepository) GetComplianceByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ComplianceByResponsibleUserDTO, error) {
	results := []*reportDTOs.ComplianceByResponsibleUserDTO{}
	query := repo.db.Table("compliance_register").
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
			COUNT(compliance_register.id) AS count
		`).
		Joins("JOIN users ON users.id = compliance_register.responsible_user_id").
		Group("users.id, users.first_name_en, users.first_name_ar, users.last_name_en, users.last_name_ar, users.initial, users.color, users.email, users.profile_image_path")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("compliance_register.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *ComplianceReportRepository) GetComplianceRequirementByComplianceStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ComplianceRequirementByComplianceStatusDTO, error) {
	var results []*reportDTOs.ComplianceRequirementByComplianceStatusDTO
	query := repo.db.Table("compliance_requirements").
		Select(`
            compliance_statuses.id AS compliance_status_id,
            compliance_statuses.title_en AS compliance_status_title_en,
            compliance_statuses.title_ar AS compliance_status_title_ar,
            COUNT(compliance_requirements.id) AS count
        `).
		Joins("JOIN compliance_statuses ON compliance_statuses.id = compliance_requirements.compliance_status_id").
		Group("compliance_statuses.id, compliance_statuses.title_en, compliance_statuses.title_ar")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("compliance_requirements.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *ComplianceReportRepository) GetComplianceRequirementBySeverity(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ComplianceRequirementBySeverityLevelDTO, error) {
	var results []*reportDTOs.ComplianceRequirementBySeverityLevelDTO
	query := repo.db.Table("compliance_requirements").
		Select(`
            severity_levels.id AS severity_level_id,
            severity_levels.title_en AS severity_level_title_en,
            severity_levels.title_ar AS severity_level_title_ar,
            severity_levels.color_code AS severity_level_color_code,
            severity_levels.type AS severity_level_type,
            COUNT(compliance_requirements.id) AS count
        `).
		Joins("JOIN severity_levels ON severity_levels.id = compliance_requirements.severity_level_id").
		Group("severity_levels.id, severity_levels.title_en, severity_levels.title_ar, severity_levels.color_code, severity_levels.type")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("compliance_requirements.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *ComplianceReportRepository) GetComplianceRequirementByCategory(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ComplianceRequirementByCategoryDTO, error) {
	var results []*reportDTOs.ComplianceRequirementByCategoryDTO
	query := repo.db.Table("compliance_requirements").
		Select(`
            compliance_requirement_categories.id AS compliance_requirement_category_id,
            compliance_requirement_categories.title_en AS compliance_requirement_category_title_en,
            compliance_requirement_categories.title_ar AS compliance_requirement_category_title_ar,
            COUNT(compliance_requirements.id) AS count
        `).
		Joins("JOIN compliance_requirement_categories ON compliance_requirement_categories.id = compliance_requirements.compliance_requirement_category_id").
		Group("compliance_requirement_categories.id, compliance_requirement_categories.title_en, compliance_requirement_categories.title_ar")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("compliance_requirements.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *ComplianceReportRepository) GetComplianceRequirementByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ComplianceRequirementByResponsibleUserDTO, error) {
	var results []*reportDTOs.ComplianceRequirementByResponsibleUserDTO
	query := repo.db.Table("compliance_requirements").
		Select(`
            users.id AS responsible_user_id,
            users.first_name_en AS responsible_user_first_name_en,
            users.first_name_ar AS responsible_user_first_name_ar,
            users.last_name_en AS responsible_user_last_name_en,
            users.last_name_ar AS responsible_user_last_name_ar,
            users.email AS responsible_user_email,
            users.color AS responsible_user_color,
            users.profile_image_path AS responsible_user_profile_image_path,
            COUNT(compliance_requirements.id) AS count
        `).
		Joins("JOIN users ON users.id = compliance_requirements.responsible_user_id").
		Group("users.id, users.first_name_en, users.first_name_ar, users.last_name_en, users.last_name_ar, users.email, users.color, users.profile_image_path")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("compliance_requirements.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}
