package repositories

import (
	reportDTOs "jethur/grc/reports/dtos"

	"gorm.io/gorm"
)

type RiskReportRepositoryInterface interface {
	GetKRIByCategory(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.KRIByCategoryDTO, error)
	GetKRIByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.KRIByResponsibleUnitDTO, error)
	GetKRIByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.KRIByResponsibleUserDTO, error)
	GetKRIByRiskLevel(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.KRIByRiskLevelDTO, error)
	GetKRIByReviewFrequency(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.KRIByReviewFrequencyDTO, error)
}

type RiskReportRepository struct {
	db *gorm.DB
}

func NewRiskReportRepository(db *gorm.DB) *RiskReportRepository {
	return &RiskReportRepository{db: db}
}

// Implementation of the `GetKRIByCategory` method
func (r *RiskReportRepository) GetKRIByCategory(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.KRIByCategoryDTO, error) {
	results := []*reportDTOs.KRIByCategoryDTO{}

	query := r.db.Table("kri_register").
		Select(`kri_categories.id AS kri_category_id, 
				kri_categories.title_en AS kri_category_title_en, 
				kri_categories.title_ar AS kri_category_title_ar, 
				COUNT(kri_register.id) AS count`).
		Joins("JOIN kri_categories ON kri_categories.id = kri_register.kri_category_id").
		Group("kri_categories.id, kri_categories.title_en, kri_categories.title_ar")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("kri_register.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (r *RiskReportRepository) GetKRIByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.KRIByResponsibleUnitDTO, error) {
	results := []*reportDTOs.KRIByResponsibleUnitDTO{}

	query := r.db.Table("kri_register").
		Select(`
			organization_units.id AS responsible_unit_id,
			organization_units.title_en AS responsible_unit_title_en,
			organization_units.title_ar AS responsible_unit_title_ar,
			organization_units.color_code AS responsible_unit_color_code,
			COUNT(kri_register.id) AS count
		`).
		Joins("JOIN organization_units ON organization_units.id = kri_register.responsible_unit_id").
		Group("organization_units.id, organization_units.title_en, organization_units.title_ar, organization_units.color_code")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("kri_register.created_at BETWEEN ? AND ?", params.DateFrom, params.DateTo)
	}

	err := query.Scan(&results).Error
	return results, err
}

func (r *RiskReportRepository) GetKRIByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.KRIByResponsibleUserDTO, error) {
	results := []*reportDTOs.KRIByResponsibleUserDTO{}

	query := r.db.Table("kri_register").
		Select(`
			users.id AS responsible_user_id,
			users.first_name_en AS responsible_user_first_name_en,
			users.first_name_ar AS responsible_user_first_name_ar,
			users.last_name_en AS responsible_user_last_name_en,
			users.last_name_ar AS responsible_user_last_name_ar,
			users.initial AS responsible_user_initial,
			users.color AS responsible_user_color,
			users.email AS responsible_user_email,
			COUNT(kri_register.id) AS count
		`).
		Joins("JOIN users ON users.id = kri_register.responsible_user_id").
		Group("users.id, users.first_name_en, users.first_name_ar, users.last_name_en, users.last_name_ar, users.initial, users.color, users.email")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("kri_register.created_at BETWEEN ? AND ?", params.DateFrom, params.DateTo)
	}

	err := query.Scan(&results).Error
	return results, err
}

func (r *RiskReportRepository) GetKRIByRiskLevel(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.KRIByRiskLevelDTO, error) {
	results := []*reportDTOs.KRIByRiskLevelDTO{}

	query := r.db.Table("kri_register").
		Select(`risk_levels.id AS risk_level_id,
				risk_levels.title_en AS risk_level_title_en,
				risk_levels.title_ar AS risk_level_title_ar,
				risk_levels.color_code AS color_code,
				COUNT(kri_register.id) AS count`).
		Joins("JOIN risk_levels ON risk_levels.id = kri_register.risk_level_id").
		Group("risk_levels.id, risk_levels.title_en, risk_levels.title_ar, risk_levels.color_code")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("kri_register.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (r *RiskReportRepository) GetKRIByReviewFrequency(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.KRIByReviewFrequencyDTO, error) {
	results := []*reportDTOs.KRIByReviewFrequencyDTO{}

	query := r.db.Table("kri_register").
		Select(`review_frequencies.id AS review_frequency_id,
				review_frequencies.title_en AS review_frequency_title_en,
				review_frequencies.title_ar AS review_frequency_title_ar,
				COUNT(kri_register.id) AS count`).
		Joins("JOIN review_frequencies ON review_frequencies.id = kri_register.review_frequency_id").
		Group("review_frequencies.id, review_frequencies.title_en, review_frequencies.title_ar")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("kri_register.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}
