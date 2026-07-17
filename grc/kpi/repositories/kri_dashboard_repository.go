package repositories

import (
	"fmt"
	kpiDTOs "jethur/grc/kpi/dtos"
	"jethur/grc/models"

	"gorm.io/gorm"
)

type KRIDashboardRepositoryInterface interface {
	GetKRIOverview(params *kpiDTOs.KriDashboardQueryParamsDTO) (*kpiDTOs.KRIOverviewDTO, error)
	GetKRIReviewDates(params *kpiDTOs.KriDashboardQueryParamsDTO) ([]kpiDTOs.KRIReviewItemDTO, error)
	GetKRIBreachedEvents() ([]kpiDTOs.KriBreachRawDTO, error)
	GetKRIStatusCategory(params *kpiDTOs.KriDashboardQueryParamsDTO) ([]kpiDTOs.KRIStatusCategoryItemDTO, int64, error)
	GetKRISeverityBreakdown(params *kpiDTOs.KriDashboardQueryParamsDTO) ([]kpiDTOs.KRISeverityItemDTO, int64, error)
	GetKRIExposure(params *kpiDTOs.KriDashboardQueryParamsDTO) ([]kpiDTOs.KriExposureRaw, error)
	GetKRISeverityByOrganizationUnit(params *kpiDTOs.KriDashboardQueryParamsDTO) ([]kpiDTOs.KriOrganizationUnitsRawDTO, error)
	GetKRISeverityByCategory(params *kpiDTOs.KriDashboardQueryParamsDTO) ([]kpiDTOs.KriCategorySeverityRawDTO, error)
}

type KRIDashboardRepository struct {
	db *gorm.DB
}

func NewKRIDashboardRepository(db *gorm.DB) KRIDashboardRepositoryInterface {
	return &KRIDashboardRepository{
		db: db,
	}
}

func (repo *KRIDashboardRepository) GetKRIOverview(params *kpiDTOs.KriDashboardQueryParamsDTO) (*kpiDTOs.KRIOverviewDTO, error) {

	result := &kpiDTOs.KRIOverviewDTO{}

	totalQuery := repo.db.Table("kri_register").Select("COUNT(*)")
	if params.DateFrom != "" && params.DateTo != "" {
		totalQuery = totalQuery.Where("created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}
	if err := totalQuery.Row().Scan(&result.TotalKRIs); err != nil {
		return nil, err
	}

	// Determine the highest available risk rating dynamically
	var highestRating string
	err := repo.db.Table("risk_matrix_risk_ratings").
		Select("type").
		Where("status = ?", "active").
		Order("CASE type WHEN 'very_high' THEN 1 WHEN 'high' THEN 2 WHEN 'medium' THEN 3 WHEN 'low' THEN 4 WHEN 'very_low' THEN 5 END").
		Limit(1).
		Pluck("type", &highestRating).Error
	if err != nil || highestRating == "" {
		return nil, fmt.Errorf("failed to fetch highest risk rating: %v", err)
	}

	highRiskQuery := repo.db.Table("risk_kris rk").
		Select("COUNT(DISTINCT rk.kri_register_id) AS high_count").
		Joins("JOIN risk_register rr ON rr.id = rk.risk_register_id").
		Joins("JOIN risk_matrix_risk_ratings rmr ON rmr.id = rr.current_risk_rating_id").
		Joins("JOIN risk_statuses rs ON rs.id = rr.risk_status_id").
		Where("rmr.type = ?", highestRating).
		Where("rs.type NOT IN (?)", models.PendingRiskStatusTypes).
		Where("rr.is_cyber_risk = ?", false)

	if params.DateFrom != "" && params.DateTo != "" {
		highRiskQuery = highRiskQuery.Where("rk.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	if err := highRiskQuery.Row().Scan(&result.HighRiskLinkedKRIs); err != nil {
		return nil, err
	}

	return result, nil
}

func (repo *KRIDashboardRepository) GetKRIReviewDates(params *kpiDTOs.KriDashboardQueryParamsDTO) ([]kpiDTOs.KRIReviewItemDTO, error) {

	var list []kpiDTOs.KRIReviewItemDTO

	query := repo.db.Table("kri_register").
		Select("id, next_review_date")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where(
			"created_at BETWEEN ? AND ?",
			params.DateFrom+" 00:00:00",
			params.DateTo+" 23:59:59",
		)
	}

	err := query.Scan(&list).Error
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (repo *KRIDashboardRepository) GetKRIBreachedEvents() ([]kpiDTOs.KriBreachRawDTO, error) {

	results := []kpiDTOs.KriBreachRawDTO{}

	query := repo.db.Table("kri_frequency_values kfv").
		Select(`
            kr.id AS kri_id,
            kr.title AS kri_title,
            kfv.actual_value AS value,
            kfv.created_at AS reviewed_at,
            krts.type AS risk_tolerance_status_type
        `).
		Joins("JOIN kri_register kr ON kr.id = kfv.kri_register_id").
		Joins("LEFT JOIN kri_risk_tolerance_statuses krts ON krts.id = kfv.risk_tolerance_status_id").
		Where("kfv.actual_value IS NOT NULL").
		Where("kfv.created_at >= DATEADD(MONTH, -6, GETDATE())")

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *KRIDashboardRepository) GetKRIStatusCategory(params *kpiDTOs.KriDashboardQueryParamsDTO) ([]kpiDTOs.KRIStatusCategoryItemDTO, int64, error) {

	categories := []kpiDTOs.KRIStatusCategoryItemDTO{}
	var total int64

	totalQuery := repo.db.Table("kri_register").Select("COUNT(*) AS total")

	if params.DateFrom != "" && params.DateTo != "" {
		totalQuery = totalQuery.Where("created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	if err := totalQuery.Row().Scan(&total); err != nil {
		return nil, 0, err
	}

	query := repo.db.Table("kri_register").
		Select(`
			kc.title_en AS title_en,
			kc.title_ar AS title_ar,
			COUNT(*) AS count
		`).
		Joins("JOIN kri_categories kc ON kc.id = kri_register.kri_category_id").
		Group("kc.title_en, kc.title_ar").
		Order("COUNT(*) DESC")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("kri_register.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	if err := query.Scan(&categories).Error; err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

func (repo *KRIDashboardRepository) GetKRISeverityBreakdown(params *kpiDTOs.KriDashboardQueryParamsDTO) ([]kpiDTOs.KRISeverityItemDTO, int64, error) {

	severities := []kpiDTOs.KRISeverityItemDTO{}
	var total int64

	totalQuery := repo.db.Table("kri_register").Select("COUNT(*) AS total")

	if params.DateFrom != "" && params.DateTo != "" {
		totalQuery = totalQuery.Where("created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	if err := totalQuery.Row().Scan(&total); err != nil {
		return nil, 0, err
	}

	query := repo.db.Table("kri_register").
		Select(`
			rl.title_en AS title_en,
			rl.title_ar AS title_ar,
			rl.color_code AS color_code,
			rl.type AS type,
			COUNT(*) AS count
		`).
		Joins("JOIN risk_levels rl ON rl.id = kri_register.risk_level_id").
		Group("rl.type, rl.title_en, rl.title_ar, rl.color_code")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("kri_register.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	if err := query.Scan(&severities).Error; err != nil {
		return nil, 0, err
	}

	return severities, total, nil
}

func (repo *KRIDashboardRepository) GetKRIExposure(params *kpiDTOs.KriDashboardQueryParamsDTO) ([]kpiDTOs.KriExposureRaw, error) {

	list := []kpiDTOs.KriExposureRaw{}

	query := repo.db.Table("kri_register kr").
		Select(`
            kr.id,
            kr.title,
            kr.current_value,
            kr.minimum_susceptibility_percent,
            kr.maximum_susceptibility_percent,
            kr.maximum_tolerable_percent,
            kr.is_lower_better,
            kr.is_higher_better
        `)

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("kr.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&list).Error
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (repo *KRIDashboardRepository) GetKRISeverityByOrganizationUnit(params *kpiDTOs.KriDashboardQueryParamsDTO) ([]kpiDTOs.KriOrganizationUnitsRawDTO, error) {

	results := []kpiDTOs.KriOrganizationUnitsRawDTO{}

	query := repo.db.Table("kri_register kr").
		Select(`
            kr.responsible_unit_id AS organization_unit_id,
            ou.title_en AS organization_unit_title_en,
            ou.title_ar AS organization_unit_title_ar,
            rl.type AS risk_level_type,
            rl.title_en AS risk_level_title_en,
            rl.title_ar AS risk_level_title_ar,
            rl.color_code AS risk_level_color_code,
            COUNT(*) AS count
        `).
		Joins("JOIN organization_units ou ON ou.id = kr.responsible_unit_id").
		Joins("LEFT JOIN risk_levels rl ON rl.id = kr.risk_level_id").
		Where("kr.risk_level_id IS NOT NULL").
		Group(`
            kr.responsible_unit_id,
            ou.title_en, ou.title_ar,
            rl.type, rl.title_en, rl.title_ar, rl.color_code
        `)

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where(
			"kr.created_at BETWEEN ? AND ?",
			params.DateFrom+" 00:00:00",
			params.DateTo+" 23:59:59",
		)
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *KRIDashboardRepository) GetKRISeverityByCategory(params *kpiDTOs.KriDashboardQueryParamsDTO) ([]kpiDTOs.KriCategorySeverityRawDTO, error) {

	results := []kpiDTOs.KriCategorySeverityRawDTO{}

	query := repo.db.Table("kri_register kr").
		Select(`
            kr.kri_category_id AS category_id,
            kc.title_en AS category_title_en,
            kc.title_ar AS category_title_ar,
            rl.type AS risk_level_type,
            rl.title_en AS risk_level_title_en,
            rl.title_ar AS risk_level_title_ar,
            rl.color_code AS risk_level_color_code,
            COUNT(*) AS count
        `).
		Joins("JOIN kri_categories kc ON kc.id = kr.kri_category_id").
		Joins("LEFT JOIN risk_levels rl ON rl.id = kr.risk_level_id").
		Where("kr.risk_level_id IS NOT NULL").
		Group(`
            kr.kri_category_id,
            kc.title_en, kc.title_ar,
            rl.type, rl.title_en, rl.title_ar, rl.color_code
        `)

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where(
			"kr.created_at BETWEEN ? AND ?",
			params.DateFrom+" 00:00:00",
			params.DateTo+" 23:59:59",
		)
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}
