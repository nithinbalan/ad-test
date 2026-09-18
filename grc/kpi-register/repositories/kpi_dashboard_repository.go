package repositories

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/kpi-register/dtos"

	"gorm.io/gorm"
)

type KPIDashboardRepositoryInterface interface {
	GetKPISummary(params *dtos.KPIDashboardQueryParamsDTO) (*dtos.KPISummaryDTO, error)
	GetKPIStatusSummary(params *dtos.KPIDashboardQueryParamsDTO) ([]*dtos.KPIStatusCountDTO, error)
	GetKPIPassFailSummary(params *dtos.KPIDashboardQueryParamsDTO) (*dtos.KPIPassFailSummaryResponseDTO, error)
	GetTopPerformingKPIs(params *dtos.KPIDashboardQueryParamsDTO) ([]*dtos.PerformanceSummaryDTO, error)
	GetBottomPerformingKPIs(params *dtos.KPIDashboardQueryParamsDTO) ([]*dtos.PerformanceSummaryDTO, error)
	GetCategorywiseKPIStats(params *dtos.KPIDashboardQueryParamsDTO) ([]*dtos.KPIPerformanceCategoryStatsDTO, error)
	GetDepartmentwiseKPIStats(params *dtos.KPIDashboardQueryParamsDTO) ([]*dtos.KPIDepartmentStatsDTO, error)
	GetOverdueForAchievementSummary(params *dtos.KPIDashboardQueryParamsDTO) (*dtos.KPIOverdueSummaryDTO, error)
	GetTopOverdueForAchievementKPIs(params *dtos.KPIDashboardQueryParamsDTO) ([]*dtos.KPIOverdueAchievementListDTO, error)
	GetOverdueReviewSummary(params *dtos.KPIDashboardQueryParamsDTO) (*dtos.KPIOverdueSummaryDTO, error)
	GetTopOverdueReviewKPIs(params *dtos.KPIDashboardQueryParamsDTO) ([]*dtos.KPIOverdueReviewListDTO, error)

	GetKPIStatusOverview(params *dtos.KPIDashboardQueryParamsDTO) ([]*dtos.KPIStatusOverviewDTO, error)
	GetAchievementSummary(params *dtos.KPIDashboardQueryParamsDTO) (*dtos.KPIAchievementSummaryDTO, error)
	GetKPIByReviewStatus(params *dtos.KPIDashboardQueryParamsDTO) ([]*dtos.KPIReviewStatusDistributionDTO, error)
	GetKPIAchievementByDomain(params *dtos.KPIDashboardQueryParamsDTO) ([]*dtos.KPIAchievementByDomainDTO, error)
	GetAllKPIs(params *dtos.KPIDashboardQueryParamsDTO) (*dtos.KPIRegisterListResponseDTO, error)
	GetTopKPIs(params *dtos.KPIDashboardQueryParamsDTO) ([]dtos.KPIRegisterListDTO, error)
	GetBottomKPIs(params *dtos.KPIDashboardQueryParamsDTO) ([]dtos.KPIRegisterListDTO, error)
}

type KPIDashboardRepository struct {
	db *gorm.DB
}

func NewKPIDashboardRepository(db *gorm.DB) KPIDashboardRepositoryInterface {
	return &KPIDashboardRepository{
		db: db,
	}
}

func (repo *KPIDashboardRepository) GetKPISummary(params *dtos.KPIDashboardQueryParamsDTO) (*dtos.KPISummaryDTO, error) {
	result := &dtos.KPISummaryDTO{}

	query := repo.db.Table("kpi_register kr").
		Select(`
		SUM(CASE WHEN kr.end_date > GETDATE() THEN 1 ELSE 0 END) AS total_active_kpi_count,
		SUM(CASE WHEN ks.type = 'achieved' THEN 1 ELSE 0 END) AS passed_kpi_count,
		SUM(CASE WHEN (CAST(kr.end_date AS DATE) < CAST(GETDATE() AS DATE) AND ks.type != 'achieved') THEN 1 ELSE 0 END) AS failed_kpi_count,
		SUM(CASE WHEN 
			(kr.end_date < GETDATE() AND ks.type != 'achieved')
			THEN 1 ELSE 0 END
		) AS overdue_kpi_count,
		SUM(CASE WHEN (ks.type = 'achieved' AND kr.actual_value > kr.target_value) THEN 1 ELSE 0 END) AS kpi_achieved_ahead_of_target_count,

		CASE 
			WHEN COUNT(kr.id) = 0 THEN 0
			ELSE
				CAST(
					(SUM(CASE WHEN ks.type = 'achieved' THEN 1 ELSE 0 END) * 100.0) / COUNT(kr.id) AS DECIMAL(10,2)
				)
		END AS overall_kpi_achievement_percentage
	`).
		Joins("LEFT JOIN kpi_statuses ks ON ks.id = kr.kpi_status_id")

	if params != nil && params.DomainID != nil {
		query = query.Where("kr.kpi_category_id = ?", params.DomainID)
	}

	if params != nil && params.IsKRI != nil {
		query = query.Where("kr.is_kri = ?", *params.IsKRI)
	}

	if params != nil {
		for col, vals := range params.Filters {
			if len(vals) > 0 {
				query = query.Where("kr."+col+" IN (?)", vals)
			}
		}
	}

	if err := query.Scan(result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (repo *KPIDashboardRepository) GetKPIStatusSummary(params *dtos.KPIDashboardQueryParamsDTO) ([]*dtos.KPIStatusCountDTO, error) {

	db := repo.db

	var statuses []*dtos.KPIStatusCountDTO
	if err := db.
		Table("kpi_statuses").
		Select("id, title_en, title_ar, color_code, type").
		Where("type IN ('new','achieved','on-track','at-risk','not-achieved')").
		Order(`
			CASE type
				WHEN 'new' THEN 1
				WHEN 'achieved' THEN 2
				WHEN 'on-track' THEN 3
				WHEN 'at-risk' THEN 4
				WHEN 'not-achieved' THEN 5
				ELSE 6
			END
		`).
		Scan(&statuses).Error; err != nil {
		return nil, err
	}

	type StatusCount struct {
		KPIStatusID uint
		Count       int64
	}

	var counts []StatusCount

	query := db.
		Table("kpi_register kr").
		Select("kr.kpi_status_id AS kpi_status_id, COUNT(kr.id) AS count").
		Group("kr.kpi_status_id")

	if params != nil {

		if params.DomainID != nil {
			query = query.Where("kr.kpi_category_id = ?", *params.DomainID)
		}

		if params.IsKRI != nil {
			query = query.Where("kr.is_kri = ?", *params.IsKRI)
		}

		for col, vals := range params.Filters {
			if len(vals) > 0 {
				query = query.Where("kr."+col+" IN (?)", vals)
			}
		}
	}

	if err := query.Scan(&counts).Error; err != nil {
		return nil, err
	}

	countMap := make(map[uint]int64)
	for _, c := range counts {
		countMap[c.KPIStatusID] = c.Count
	}

	for i := range statuses {
		statuses[i].Count = countMap[statuses[i].ID]
	}

	return statuses, nil
}

func (repo *KPIDashboardRepository) GetKPIPassFailSummary(params *dtos.KPIDashboardQueryParamsDTO) (*dtos.KPIPassFailSummaryResponseDTO, error) {

	var result dtos.KPIPassFailSummaryResponseDTO

	query := repo.db.
		Table("kpi_register kr").
		Select(`
			COUNT(kr.id) AS total,
			SUM(CASE WHEN ks.type = 'achieved' THEN 1 ELSE 0 END) AS passed,
			SUM(CASE WHEN (ks.type != 'achieved' AND CAST(kr.end_date AS DATE) < CAST(GETDATE() AS DATE)) THEN 1 ELSE 0 END) AS failed
		`).
		Joins("JOIN kpi_statuses ks ON ks.id = kr.kpi_status_id")

	if params != nil && params.DomainID != nil {
		query = query.Where("kr.kpi_category_id = ?", params.DomainID)
	}

	if params != nil && params.IsKRI != nil {
		query = query.Where("kr.is_kri = ?", *params.IsKRI)
	}

	if params != nil {
		for col, vals := range params.Filters {
			if len(vals) > 0 {
				query = query.Where("kr."+col+" IN (?)", vals)
			}
		}
	}

	if err := query.Scan(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

func (repo *KPIDashboardRepository) GetTopPerformingKPIs(params *dtos.KPIDashboardQueryParamsDTO) ([]*dtos.PerformanceSummaryDTO, error) {

	var results []*dtos.PerformanceSummaryDTO

	query := repo.db.
		Table("kpi_register kr").
		Select(`
			kr.id,
			kr.title,
			kr.responsible_user_id,

			u.first_name_en  AS responsible_user_first_name_en,
			u.last_name_en   AS responsible_user_last_name_en,
			u.first_name_ar  AS responsible_user_first_name_ar,
			u.last_name_ar   AS responsible_user_last_name_ar,
			u.initial        AS responsible_user_initial,
			u.email          AS responsible_user_email,
			u.color          AS responsible_user_color,
			u.profile_image_path AS responsible_user_profile_image_path,

			role_data.title_en AS responsible_user_role_title_en,
			role_data.title_ar AS responsible_user_role_title_ar,

			CAST(
				CASE 
					WHEN kr.achievement > 100 THEN 100
					WHEN kr.achievement < 0 THEN 0
					ELSE kr.achievement
				END AS DECIMAL(10,2)
			) AS performance_percentage,

			ks.title_en AS status_title_en,
			ks.title_ar AS status_title_ar,
			ks.color_code AS status_color_code,
			ks.type AS status_type,

			CASE 
				WHEN kr.achievement >= 95 THEN 'Outstanding'
				WHEN kr.achievement >= 90 THEN 'Excellent'
				WHEN kr.achievement >= 80 THEN 'Good'
				ELSE 'Needs Attention'
			END AS performance_label
		`).
		Joins(`
			LEFT JOIN kpi_statuses ks 
				ON ks.id = kr.kpi_status_id
		`).
		Joins(`
			LEFT JOIN users u
				ON u.id = kr.responsible_user_id
		`).
		Joins(`
			OUTER APPLY (
				SELECT TOP 1 r.title_en, r.title_ar
				FROM user_roles ur
				JOIN roles r ON r.id = ur.role_id
				WHERE ur.user_id = u.id
			) role_data
		`).
		Where(`
			kr.target_value > 0
			AND ks.type = 'achieved'
		`)

	if params != nil && params.DomainID != nil {
		query = query.Where("kr.kpi_category_id = ?", *params.DomainID)
	}

	if params != nil && params.IsKRI != nil {
		query = query.Where("kr.is_kri = ?", *params.IsKRI)
	}

	if params != nil {
		for col, vals := range params.Filters {
			if len(vals) == 0 {
				continue
			}

			if col == "kpi_category_id" && params.DomainID != nil {
				continue
			}

			query = query.Where(fmt.Sprintf("kr.%s IN ?", col), vals)
		}
	}

	query = query.
		Order(`
			CASE 
				WHEN kr.achievement > 100 THEN 100
				ELSE kr.achievement
			END DESC
		`).
		Limit(5)

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *KPIDashboardRepository) GetBottomPerformingKPIs(params *dtos.KPIDashboardQueryParamsDTO) ([]*dtos.PerformanceSummaryDTO, error) {

	var results []*dtos.PerformanceSummaryDTO

	query := repo.db.
		Table("kpi_register kr").
		Select(`
			kr.id,
			kr.title,
			kr.responsible_user_id,

			u.first_name_en  AS responsible_user_first_name_en,
			u.last_name_en   AS responsible_user_last_name_en,
			u.first_name_ar  AS responsible_user_first_name_ar,
			u.last_name_ar   AS responsible_user_last_name_ar,
			u.initial        AS responsible_user_initial,
			u.email          AS responsible_user_email,
			u.color          AS responsible_user_color,
			u.profile_image_path AS responsible_user_profile_image_path,

			role_data.title_en AS responsible_user_role_title_en,
			role_data.title_ar AS responsible_user_role_title_ar,

			CAST(
				CASE 
					WHEN kr.achievement > 100 THEN 100
					WHEN kr.achievement < 0 THEN 0
					ELSE kr.achievement
				END AS DECIMAL(10,2)
			) AS performance_percentage,

			ks.title_en AS status_title_en,
			ks.title_ar AS status_title_ar,
			ks.color_code AS status_color_code,
			ks.type AS status_type,

			CASE 
				WHEN (kr.actual_value * 100.0 / NULLIF(kr.target_value,0)) < 50 THEN 'Poor'
				WHEN (kr.actual_value * 100.0 / NULLIF(kr.target_value,0)) < 60 THEN 'Fair'
				WHEN (kr.actual_value * 100.0 / NULLIF(kr.target_value,0)) < 75 THEN 'Needs Attention'
				ELSE 'Good'
			END AS performance_label
		`).
		Joins(`
			LEFT JOIN kpi_statuses ks 
				ON ks.id = kr.kpi_status_id
		`).
		Joins(`
			LEFT JOIN users u
				ON u.id = kr.responsible_user_id
		`).
		Joins(`
			OUTER APPLY (
				SELECT TOP 1 r.title_en, r.title_ar
				FROM user_roles ur
				JOIN roles r ON r.id = ur.role_id
				WHERE ur.user_id = u.id
			) role_data
		`).
		Where(`
			kr.target_value > 0
			AND NOT (ks.type = 'achieved')
		`)

	if params != nil && params.DomainID != nil {
		query = query.Where("kr.kpi_category_id = ?", *params.DomainID)
	}

	if params != nil && params.IsKRI != nil {
		query = query.Where("kr.is_kri = ?", *params.IsKRI)
	}

	if params != nil {
		for col, vals := range params.Filters {
			if len(vals) == 0 {
				continue
			}

			if col == "kpi_category_id" && params.DomainID != nil {
				continue
			}

			query = query.Where(fmt.Sprintf("kr.%s IN ?", col), vals)
		}
	}

	query = query.
		Order("(kr.actual_value * 100.0 / NULLIF(kr.target_value,0)) ASC").
		Limit(5)

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *KPIDashboardRepository) GetCategorywiseKPIStats(params *dtos.KPIDashboardQueryParamsDTO) ([]*dtos.KPIPerformanceCategoryStatsDTO, error) {

	type Category struct {
		ID      uint
		TitleEn string
		TitleAr string
	}

	var categories []Category

	if err := repo.db.
		Table("kpi_categories").
		Select("id, title_en, title_ar").
		Where("status = ?", "active").
		Order("id").
		Find(&categories).Error; err != nil {
		return nil, err
	}

	type CategoryCount struct {
		CategoryID uint
		Count      int64
	}

	var counts []CategoryCount

	query := repo.db.
		Table("kpi_register kr").
		Select("kr.kpi_category_id AS category_id, COUNT(kr.id) AS count").
		Group("kr.kpi_category_id")

	query = query.Where("kr.target_value > 0")

	if params != nil && params.DomainID != nil {
		query = query.Where("kr.kpi_category_id = ?", *params.DomainID)
	}

	if params != nil && params.IsKRI != nil {
		query = query.Where("kr.is_kri = ?", *params.IsKRI)
	}

	if params != nil {
		for col, vals := range params.Filters {
			if len(vals) > 0 {
				query = query.Where("kr."+col+" IN ?", vals)
			}
		}
	}

	if err := query.Scan(&counts).Error; err != nil {
		return nil, err
	}

	countMap := make(map[uint]int64)
	var total int64

	for _, c := range counts {
		countMap[c.CategoryID] = c.Count
		total += c.Count
	}

	results := []*dtos.KPIPerformanceCategoryStatsDTO{}

	for _, cat := range categories {

		count := countMap[cat.ID]

		var percentage float64
		if total > 0 {
			percentage = (float64(count) / float64(total)) * 100
		}

		results = append(results, &dtos.KPIPerformanceCategoryStatsDTO{
			ID:         cat.ID,
			TitleEn:    cat.TitleEn,
			TitleAr:    cat.TitleAr,
			Count:      count,
			Percentage: percentage,
		})
	}

	return results, nil
}

func (repo *KPIDashboardRepository) GetDepartmentwiseKPIStats(params *dtos.KPIDashboardQueryParamsDTO) ([]*dtos.KPIDepartmentStatsDTO, error) {

	type OrgLevel struct {
		ID    uint
		Level int
	}

	var primaryLevel OrgLevel

	if err := repo.db.
		Table("organization_levels").
		Select("id, level").
		Where("is_primary = ?", true).
		First(&primaryLevel).Error; err != nil {
		return nil, err
	}

	primaryField := fmt.Sprintf("u.organization_unit_%d_id", primaryLevel.Level)

	type Department struct {
		ID        uint
		TitleEn   string
		TitleAr   string
		ColorCode string
	}

	var departments []Department

	if err := repo.db.
		Table("organization_units ou").
		Select("ou.id, ou.title_en, ou.title_ar, ou.color_code").
		Where("ou.organization_level_id = ? AND ou.status = ?", primaryLevel.ID, "active").
		Order("ou.id").
		Find(&departments).Error; err != nil {
		return nil, err
	}

	type DepartmentCount struct {
		DepartmentID uint
		Count        int64
	}

	var counts []DepartmentCount

	query := repo.db.
		Table("kpi_register kr").
		Select(fmt.Sprintf(`
			%s AS department_id,
			COUNT(DISTINCT kr.id) AS count
		`, primaryField)).
		Joins("JOIN users u ON u.id = kr.responsible_user_id").
		Where(fmt.Sprintf("%s IS NOT NULL", primaryField)).
		Group(primaryField)

	if params != nil {

		if params.DomainID != nil {
			query = query.Where("kr.kpi_category_id = ?", *params.DomainID)
		}

		if params.IsKRI != nil {
			query = query.Where("kr.is_kri = ?", *params.IsKRI)
		}

		for col, vals := range params.Filters {
			if len(vals) > 0 {
				query = query.Where("kr."+col+" IN ?", vals)
			}
		}
	}

	if err := query.Scan(&counts).Error; err != nil {
		return nil, err
	}

	countMap := make(map[uint]int64)
	var total int64

	for _, c := range counts {
		countMap[c.DepartmentID] = c.Count
		total += c.Count
	}

	results := []*dtos.KPIDepartmentStatsDTO{}

	for _, dept := range departments {

		count := countMap[dept.ID]

		var percentage float64
		if total > 0 {
			percentage = (float64(count) / float64(total)) * 100
		}

		results = append(results, &dtos.KPIDepartmentStatsDTO{
			ID:         dept.ID,
			TitleEn:    dept.TitleEn,
			TitleAr:    dept.TitleAr,
			ColorCode:  dept.ColorCode,
			Count:      count,
			Percentage: percentage,
		})
	}

	return results, nil
}

func (repo *KPIDashboardRepository) GetOverdueForAchievementSummary(params *dtos.KPIDashboardQueryParamsDTO) (*dtos.KPIOverdueSummaryDTO, error) {

	var result dtos.KPIOverdueSummaryDTO

	query := repo.db.
		Table("kpi_register kr").
		Select(`
			COUNT(kr.id) AS total_count,
			ISNULL(AVG(ABS(DATEDIFF(DAY, kr.end_date, GETDATE()))), 0) AS average_days_overdue
		`).
		Joins("LEFT JOIN kpi_statuses ks ON ks.id = kr.kpi_status_id").
		Where(`
			ks.type != 'achieved'
			AND kr.end_date < GETDATE()
		`)

	if params != nil && params.DomainID != nil {
		query = query.Where("kr.kpi_category_id = ?", *params.DomainID)
	}

	if params != nil && params.IsKRI != nil {
		query = query.Where("kr.is_kri = ?", *params.IsKRI)
	}

	if params != nil {
		for col, vals := range params.Filters {
			if len(vals) > 0 {
				query = query.Where("kr."+col+" IN ?", vals)
			}
		}
	}

	if err := query.Scan(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

func (repo *KPIDashboardRepository) GetTopOverdueForAchievementKPIs(params *dtos.KPIDashboardQueryParamsDTO) ([]*dtos.KPIOverdueAchievementListDTO, error) {

	var results []*dtos.KPIOverdueAchievementListDTO

	query := repo.db.
		Table("kpi_register kr").
		Select(`
			kr.id,
			kr.title,
			kr.responsible_user_id,
			u.first_name_en  AS responsible_user_first_name_en,
			u.last_name_en   AS responsible_user_last_name_en,
			u.first_name_ar  AS responsible_user_first_name_ar,
			u.last_name_ar   AS responsible_user_last_name_ar,
			u.initial        AS responsible_user_initial,
			u.color          AS responsible_user_color,
			u.profile_image_path AS responsible_user_profile_image_path,
			kr.target_value AS target_value,
			kr.actual_value AS actual_value,
			DATEDIFF(DAY, kr.end_date, GETDATE()) AS days_overdue
		`).
		Joins("LEFT JOIN users u ON u.id = kr.responsible_user_id").
		Joins("LEFT JOIN kpi_statuses ks ON ks.id = kr.kpi_status_id").
		Where(`
			ks.type != 'achieved'
			AND kr.end_date < GETDATE()
		`)

	if params != nil && params.DomainID != nil {
		query = query.Where("kr.kpi_category_id = ?", *params.DomainID)
	}

	if params != nil && params.IsKRI != nil {
		query = query.Where("kr.is_kri = ?", *params.IsKRI)
	}

	if params != nil {
		for col, vals := range params.Filters {
			if len(vals) > 0 {
				query = query.Where("kr."+col+" IN ?", vals)
			}
		}
	}

	query = query.
		Order("days_overdue DESC").
		Limit(5)

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *KPIDashboardRepository) GetOverdueReviewSummary(params *dtos.KPIDashboardQueryParamsDTO) (*dtos.KPIOverdueSummaryDTO, error) {

	var result dtos.KPIOverdueSummaryDTO

	query := repo.db.
		Table("kpi_register kr").
		Select(`
			COUNT(kr.id) AS total_count,
			ISNULL(AVG(ABS(DATEDIFF(DAY, kr.end_date, GETDATE()))), 0) AS average_days_overdue
		`).
		Where(`
			kr.next_review_date < GETDATE() AND (kr.last_reviewed_date IS NULL OR kr.last_reviewed_date != kr.next_review_date)
		`)

	if params != nil && params.DomainID != nil {
		query = query.Where("kr.kpi_category_id = ?", *params.DomainID)
	}

	if params != nil && params.IsKRI != nil {
		query = query.Where("kr.is_kri = ?", *params.IsKRI)
	}

	if params != nil {
		for col, vals := range params.Filters {
			if len(vals) > 0 {
				query = query.Where("kr."+col+" IN ?", vals)
			}
		}
	}

	if err := query.Scan(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

func (repo *KPIDashboardRepository) GetTopOverdueReviewKPIs(params *dtos.KPIDashboardQueryParamsDTO) ([]*dtos.KPIOverdueReviewListDTO, error) {

	var results []*dtos.KPIOverdueReviewListDTO

	query := repo.db.
		Table("kpi_register kr").
		Select(`
			kr.id,
			kr.title,
			kr.responsible_user_id,
			u.first_name_en  AS responsible_user_first_name_en,
			u.last_name_en   AS responsible_user_last_name_en,
			u.first_name_ar  AS responsible_user_first_name_ar,
			u.last_name_ar   AS responsible_user_last_name_ar,
			u.initial        AS responsible_user_initial,
			u.color          AS responsible_user_color,
			u.profile_image_path AS responsible_user_profile_image_path,
			kr.next_review_date As target_review_date,
			DATEDIFF(DAY, kr.next_review_date, GETDATE()) AS days_overdue
		`).
		Joins("LEFT JOIN users u ON u.id = kr.responsible_user_id").
		Where(`
			kr.next_review_date < GETDATE() AND (kr.last_reviewed_date IS NULL OR kr.last_reviewed_date != kr.next_review_date)
		`)

	if params != nil && params.DomainID != nil {
		query = query.Where("kr.kpi_category_id = ?", *params.DomainID)
	}

	if params != nil && params.IsKRI != nil {
		query = query.Where("kr.is_kri = ?", *params.IsKRI)
	}

	if params != nil {
		for col, vals := range params.Filters {
			if len(vals) > 0 {
				query = query.Where("kr."+col+" IN ?", vals)
			}
		}
	}

	query = query.
		Order("days_overdue DESC").
		Limit(5)

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *KPIDashboardRepository) GetKPIStatusOverview(params *dtos.KPIDashboardQueryParamsDTO) ([]*dtos.KPIStatusOverviewDTO, error) {

	var results []*dtos.KPIStatusOverviewDTO

	rankedValuesSubQuery := repo.db.
		Table("kpi_value_updates").
		Select(`
			kpi_register_id,
			value,
			ROW_NUMBER() OVER (
				PARTITION BY kpi_register_id 
				ORDER BY reviewed_date DESC, created_at DESC
			) AS rn
		`).
		Where("reviewed_date IS NOT NULL")

	if params != nil && params.StartDate != nil && params.EndDate != nil {
		rankedValuesSubQuery = rankedValuesSubQuery.
			Where("reviewed_date BETWEEN ? AND ?", *params.StartDate, *params.EndDate)
	}

	kpiValuesQuery := repo.db.
		Table("(?) ranked_values", rankedValuesSubQuery).
		Select(`
			kpi_register_id,
			MAX(CASE WHEN rn = 1 THEN value END) AS latest_value,
			MAX(CASE WHEN rn = 2 THEN value END) AS prev_value
		`).
		Group("kpi_register_id")

	krJoinConditions := "kr.kpi_status_id = ks.id"
	krJoinArgs := []interface{}{}

	if params != nil && params.IsKRI != nil {
		krJoinConditions += " AND kr.is_kri = ?"
		krJoinArgs = append(krJoinArgs, *params.IsKRI)
	}

	if params != nil && params.Filters != nil {
		for columnName, values := range params.Filters {
			if len(values) > 0 {
				krJoinConditions += " AND kr." + columnName + " IN (?)"
				krJoinArgs = append(krJoinArgs, values)
			}
		}
	}

	if params != nil && len(params.QuickFilter) > 0 {
		krJoinConditions += " AND kr.type IN (?)"
		krJoinArgs = append(krJoinArgs, params.QuickFilter)
	}

	query := repo.db.
		Table("kpi_statuses ks").
		Select(`
			ks.id,
			ks.title_en,
			ks.title_ar,
			ks.color_code,
			ks.type,
			ks.label,

			COUNT(DISTINCT kr.id) AS count,

			COALESCE(AVG(
				CASE 
					WHEN kpi_values.latest_value IS NOT NULL AND kr.target_value IS NOT NULL THEN
						CASE 
							WHEN kal.type = 'lower' THEN 
								CASE WHEN kpi_values.latest_value = 0 THEN 0 ELSE (kr.target_value / kpi_values.latest_value) * 100 END
							ELSE 
								CASE WHEN kr.target_value = 0 THEN 0 ELSE (kpi_values.latest_value / kr.target_value) * 100 END
						END
				END
			), 0) AS average_achievement,

			COALESCE(AVG(
				CASE 
					WHEN kpi_values.latest_value IS NOT NULL AND kpi_values.prev_value IS NOT NULL
					THEN ((kpi_values.latest_value - kpi_values.prev_value) / NULLIF(kpi_values.prev_value, 0)) * 100
				END
			), 0) AS average_trend
		`).
		Joins("LEFT JOIN kpi_register kr ON "+krJoinConditions, krJoinArgs...).
		Joins("LEFT JOIN (?) kpi_values ON kpi_values.kpi_register_id = kr.id", kpiValuesQuery).
		Joins("LEFT JOIN review_frequencies rf ON rf.id = kr.review_frequency_id").
		Joins("LEFT JOIN kpi_achievement_logics kal ON kal.id = kr.kpi_achievement_logic_id").
		Where("ks.type <> ?", "new")

	if params != nil && params.PeriodType != nil {
		var freqTypes []string

		switch *params.PeriodType {
		case "this-month":
			freqTypes = []string{"daily", "weekly", "monthly"}
		case "this-quarter":
			freqTypes = []string{"daily", "weekly", "monthly", "quarterly"}
		case "this-half-year":
			freqTypes = []string{"daily", "weekly", "monthly", "quarterly", "half-yearly"}
		case "this-year", "custom-range":
			freqTypes = []string{"daily", "weekly", "monthly", "quarterly", "half-yearly", "annually"}
		}

		if len(freqTypes) > 0 {
			query = query.Where("rf.type IN ? OR rf.type IS NULL", freqTypes)
		}
	}

	query = query.
		Group("ks.id, ks.title_en, ks.title_ar, ks.color_code, ks.type, ks.label, ks.[order]").
		Order("ks.[order]")

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *KPIDashboardRepository) GetAchievementSummary(params *dtos.KPIDashboardQueryParamsDTO) (*dtos.KPIAchievementSummaryDTO, error) {

	result := &dtos.KPIAchievementSummaryDTO{}

	latestValueSubQuery := repo.db.
		Table("kpi_value_updates").
		Select(`
			id,
			kpi_register_id,
			value,
			reviewed_date,
			ROW_NUMBER() OVER (
				PARTITION BY kpi_register_id 
				ORDER BY reviewed_date DESC, id DESC
			) AS rn
		`).
		Where("reviewed_date IS NOT NULL")

	if params != nil && params.StartDate != nil && params.EndDate != nil {
		latestValueSubQuery = latestValueSubQuery.
			Where("reviewed_date BETWEEN ? AND ?", *params.StartDate, *params.EndDate)
	}

	query := repo.db.Table("kpi_register kr").
		Joins("LEFT JOIN (?) kvu ON kvu.kpi_register_id = kr.id AND kvu.rn = 1", latestValueSubQuery).
		Joins("LEFT JOIN review_frequencies rf ON rf.id = kr.review_frequency_id").
		Joins("LEFT JOIN kpi_achievement_logics kal ON kal.id = kr.kpi_achievement_logic_id").
		Select(`
			COUNT(DISTINCT kr.id) AS total_kpi_count,
			COALESCE(
				CAST(
					SUM(
						CASE 
							WHEN kvu.value IS NOT NULL AND kr.target_value IS NOT NULL THEN
								CASE 
									WHEN kal.type = 'lower' THEN 
										CASE 
											WHEN kvu.value = 0 THEN 0
											WHEN (kr.target_value / kvu.value) * 100 > 100 THEN 100
											ELSE (kr.target_value / kvu.value) * 100
										END
									ELSE 
										CASE 
											WHEN kr.target_value = 0 THEN 0
											WHEN (kvu.value / kr.target_value) * 100 > 100 THEN 100
											ELSE (kvu.value / kr.target_value) * 100
										END
								END
							ELSE 0
						END
					) / NULLIF(COUNT(DISTINCT kr.id), 0)
				AS DECIMAL(10,2)),
			0) AS overall_achievement
		`)

	if params != nil && params.IsKRI != nil {
		query = query.Where("kr.is_kri = ?", *params.IsKRI)
	}

	if params != nil && params.Filters != nil {
		for columnName, values := range params.Filters {
			if len(values) > 0 {
				query = query.Where("kr."+columnName+" IN (?)", values)
			}
		}
	}

	if params != nil && len(params.QuickFilter) > 0 {
		query = query.Where("kr.type IN (?)", params.QuickFilter)
	}

	if params != nil && params.PeriodType != nil {
		var freqTypes []string

		switch *params.PeriodType {
		case "this-month":
			freqTypes = []string{"daily", "weekly", "monthly"}
		case "this-quarter":
			freqTypes = []string{"daily", "weekly", "monthly", "quarterly"}
		case "this-half-year":
			freqTypes = []string{"daily", "weekly", "monthly", "quarterly", "half-yearly"}
		case "this-year", "custom-range":
			freqTypes = []string{"daily", "weekly", "monthly", "quarterly", "half-yearly", "annually"}
		}

		if len(freqTypes) > 0 {
			query = query.Where("rf.type IN ? OR rf.type IS NULL", freqTypes)
		}
	}

	if err := query.Scan(result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (repo *KPIDashboardRepository) GetKPIByReviewStatus(params *dtos.KPIDashboardQueryParamsDTO) ([]*dtos.KPIReviewStatusDistributionDTO, error) {

	results := []*dtos.KPIReviewStatusDistributionDTO{}

	krSubQuery := repo.db.
		Table("kpi_register kr").
		Select("kr.id, kr.kpi_review_status_id").
		Joins("LEFT JOIN review_frequencies rf ON rf.id = kr.review_frequency_id")

	if params != nil && params.DomainID != nil {
		krSubQuery = krSubQuery.Where("kr.kpi_category_id = ?", *params.DomainID)
	}

	if params != nil && params.IsKRI != nil {
		krSubQuery = krSubQuery.Where("kr.is_kri = ?", *params.IsKRI)
	}

	if params != nil && params.Filters != nil {
		for columnName, values := range params.Filters {
			if len(values) > 0 {
				krSubQuery = krSubQuery.Where("kr."+columnName+" IN (?)", values)
			}
		}
	}

	if params != nil && len(params.QuickFilter) > 0 {
		krSubQuery = krSubQuery.Where("kr.type IN (?)", params.QuickFilter)
	}

	if params != nil && params.StartDate != nil && params.EndDate != nil {
		krSubQuery = krSubQuery.Where("kr.start_date <= ? AND kr.end_date >= ?", *params.EndDate, *params.StartDate)
	}

	if params != nil && params.PeriodType != nil {
		var freqTypes []string

		switch *params.PeriodType {
		case "this-month":
			freqTypes = []string{"daily", "weekly", "biweekly", "monthly"}
		case "this-quarter":
			freqTypes = []string{"daily", "weekly", "biweekly", "monthly", "quarterly"}
		case "this-half-year":
			freqTypes = []string{"daily", "weekly", "biweekly", "monthly", "quarterly", "half-yearly"}
		case "this-year", "custom-range":
			freqTypes = []string{"daily", "weekly", "biweekly", "monthly", "quarterly", "half-yearly", "annually"}
		}

		if len(freqTypes) > 0 {
			krSubQuery = krSubQuery.Where("rf.type IN ? OR rf.type IS NULL", freqTypes)
		}
	}

	query := repo.db.
		Table("kpi_review_statuses krs").
		Select(`
			krs.id,
			krs.title_en,
			krs.title_ar,
			krs.color_code,
			krs.type,
			COUNT(DISTINCT kr.id) AS count
		`).
		Where("krs.status = ?", "active").
		Joins("LEFT JOIN (?) kr ON kr.kpi_review_status_id = krs.id", krSubQuery).
		Group("krs.id, krs.title_en, krs.title_ar, krs.color_code, krs.type").
		Order("krs.id")

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *KPIDashboardRepository) GetKPIAchievementByDomain(params *dtos.KPIDashboardQueryParamsDTO) ([]*dtos.KPIAchievementByDomainDTO, error) {
	var results []*dtos.KPIAchievementByDomainDTO

	latestValueSubQuery := repo.db.
		Table("kpi_value_updates").
		Select(`
			id,
			kpi_register_id,
			value,
			reviewed_date,
			ROW_NUMBER() OVER (
				PARTITION BY kpi_register_id 
				ORDER BY reviewed_date DESC, id DESC
			) AS rn
		`).
		Where("reviewed_date IS NOT NULL")

	if params != nil && params.StartDate != nil && params.EndDate != nil {
		latestValueSubQuery = latestValueSubQuery.
			Where("reviewed_date BETWEEN ? AND ?", *params.StartDate, *params.EndDate)
	}

	kpiScoreSubQuery := repo.db.
		Table("kpi_register kr").
		Select(`
			kr.id,
			kr.kpi_category_id,
			CASE 
				WHEN kvu.value IS NOT NULL AND kr.target_value IS NOT NULL THEN
					CASE 
						WHEN kal.type = 'lower' THEN 
							CASE 
								WHEN kvu.value = 0 THEN 0
								ELSE (kr.target_value / kvu.value) * 100
							END
						ELSE 
							CASE 
								WHEN kr.target_value = 0 THEN 0
								ELSE (kvu.value / kr.target_value) * 100
							END
					END
				ELSE NULL
			END AS score
		`).
		Joins("LEFT JOIN (?) kvu ON kvu.kpi_register_id = kr.id AND kvu.rn = 1", latestValueSubQuery).
		Joins("LEFT JOIN kpi_achievement_logics kal ON kal.id = kr.kpi_achievement_logic_id")

	if params != nil && params.IsKRI != nil {
		kpiScoreSubQuery = kpiScoreSubQuery.Where("kr.is_kri = ?", *params.IsKRI)
	}

	if params != nil && params.Filters != nil {
		for columnName, values := range params.Filters {
			if columnName == "kpi_category_id" {
				continue
			}
			if len(values) > 0 {
				kpiScoreSubQuery = kpiScoreSubQuery.Where("kr."+columnName+" IN (?)", values)
			}
		}
	}

	if params != nil && len(params.QuickFilter) > 0 {
		kpiScoreSubQuery = kpiScoreSubQuery.Where("kr.type IN (?)", params.QuickFilter)
	}

	query := repo.db.
		Table("kpi_categories kc").
		Select(`
			kc.id AS category_id,
			kc.title_en AS category_name,
			kc.type AS category_type,
			COALESCE(AVG(kpi_scores.score), 0) AS average_score
		`).
		Where("kc.status = ?", "active").
		Joins("LEFT JOIN (?) kpi_scores ON kpi_scores.kpi_category_id = kc.id", kpiScoreSubQuery)

	if params != nil && params.Filters != nil {
		if categoryIDs, ok := params.Filters["kpi_category_id"]; ok && len(categoryIDs) > 0 {
			query = query.Where("kc.id IN (?)", categoryIDs)
		}
	}

	if params != nil && params.PeriodType != nil {
		var freqTypes []string

		switch *params.PeriodType {
		case "this-month":
			freqTypes = []string{"daily", "weekly", "monthly"}
		case "this-quarter":
			freqTypes = []string{"daily", "weekly", "monthly", "quarterly"}
		case "this-half-year":
			freqTypes = []string{"daily", "weekly", "monthly", "quarterly", "half-yearly"}
		case "this-year", "custom-range":
			freqTypes = []string{"daily", "weekly", "monthly", "quarterly", "half-yearly", "annually"}
		}

		if len(freqTypes) > 0 {
			query = query.Joins("LEFT JOIN kpi_register kr ON kr.kpi_category_id = kc.id").
				Joins("LEFT JOIN review_frequencies rf ON rf.id = kr.review_frequency_id").
				Where("rf.type IN ? OR rf.type IS NULL", freqTypes)
		}
	}

	query = query.
		Group("kc.id, kc.title_en, kc.type").
		Order("kc.id")

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *KPIDashboardRepository) GetAllKPIs(params *dtos.KPIDashboardQueryParamsDTO) (*dtos.KPIRegisterListResponseDTO, error) {

	var kpiRegister []*dtos.KPIRegisterListDTO
	var totalCount int64

	rankedValuesQuery := repo.db.Table(`
		(
			SELECT 
				kpi_register_id,
				value,
				ROW_NUMBER() OVER (
					PARTITION BY kpi_register_id 
					ORDER BY reviewed_date DESC, created_at DESC
				) AS rn
			FROM kpi_value_updates
			WHERE reviewed_date IS NOT NULL
		) ranked_values
	`)

	if params.StartDate != nil && params.EndDate != nil {
		rankedValuesQuery = repo.db.Table(`
			(
				SELECT 
					kpi_register_id,
					value,
					ROW_NUMBER() OVER (
						PARTITION BY kpi_register_id 
						ORDER BY reviewed_date DESC, created_at DESC
					) AS rn
				FROM kpi_value_updates
				WHERE reviewed_date BETWEEN ? AND ?
			) ranked_values
		`, *params.StartDate, *params.EndDate)
	}

	kpiValuesQuery := repo.db.
		Table("(?) kpi_values", rankedValuesQuery).
		Select(`
			kpi_register_id,
			MAX(CASE WHEN rn = 1 THEN value END) AS latest_value,
			MAX(CASE WHEN rn = 2 THEN value END) AS prev_value
		`).
		Group("kpi_register_id")

	query := repo.db.Table("kpi_register AS kr").
		Joins("LEFT JOIN (?) kpi_values ON kpi_values.kpi_register_id = kr.id", kpiValuesQuery).
		Joins("LEFT JOIN kpi_categories kc ON kc.id = kr.kpi_category_id").
		Joins("LEFT JOIN kpi_statuses ks ON ks.id = kr.kpi_status_id").
		Joins("LEFT JOIN kpi_review_statuses krs ON krs.id = kr.kpi_review_status_id").
		Joins("LEFT JOIN review_frequencies rf ON rf.id = kr.review_frequency_id").
		Joins("LEFT JOIN calculation_types ct ON ct.id = kr.calculation_type_id").
		Joins("LEFT JOIN units u ON u.id = kr.unit_id").
		Joins("LEFT JOIN users ru ON ru.id = kr.responsible_user_id").
		Joins("LEFT JOIN users cu ON cu.id = kr.created_by_id").
		Joins("LEFT JOIN users uu ON uu.id = kr.updated_by_id").
		Joins("LEFT JOIN kpi_achievement_logics kal ON kal.id = kr.kpi_achievement_logic_id")

	query = query.Where("ks.status = ?", "active")
	query = query.Where("ks.type <> ?", "new")

	if params != nil && params.IsKRI != nil {
		query = query.Where("kr.is_kri = ?", *params.IsKRI)
	}

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where("kr."+columnName+" IN (?)", values)
		}
	}

	for _, filter := range params.QuickFilter {
		switch filter {
		case "top-achievement":
			query = query.
				Where("kpi_values.latest_value IS NOT NULL").
				Where("ks.is_not_achieved = ?", false)

		case "bottom-achievement":
			query = query.
				Where("kpi_values.latest_value IS NOT NULL").
				Where("ks.is_not_achieved = ?", true)
		}
	}

	if params.PeriodType != nil {
		var freqTypes []string

		switch *params.PeriodType {
		case "this-month":
			freqTypes = []string{"daily", "weekly", "monthly"}
		case "this-quarter":
			freqTypes = []string{"daily", "weekly", "monthly", "quarterly"}
		case "this-half-year":
			freqTypes = []string{"daily", "weekly", "monthly", "quarterly", "half-yearly"}
		case "this-year", "custom-range":
			freqTypes = []string{"daily", "weekly", "monthly", "quarterly", "half-yearly", "annually"}
		}

		if len(freqTypes) > 0 {
			query = query.Where("rf.type IN ? OR rf.type IS NULL", freqTypes)
		}
	}

	query = query.Select(`
		kr.*,

		CAST(ROUND(kpi_values.latest_value, 2) AS DECIMAL(10,2)) AS actual_value,
		CAST(ROUND(kpi_values.prev_value, 2) AS DECIMAL(10,2)) AS previous_value,

		ROUND(
			CASE 
				WHEN kpi_values.latest_value IS NOT NULL AND kpi_values.prev_value IS NOT NULL
				THEN ((kpi_values.latest_value - kpi_values.prev_value) / NULLIF(kpi_values.prev_value, 0)) * 100
				ELSE 0
			END,
		2) AS trend,

		CONVERT(VARCHAR, kr.start_date, 23) AS start_date,
		CONVERT(VARCHAR, kr.end_date, 23) AS end_date,

		kc.title_en AS kpi_category_title_en,
		kc.title_ar AS kpi_category_title_ar,

		ks.title_en AS kpi_status_title_en,
		ks.title_ar AS kpi_status_title_ar,
		ks.type AS kpi_status_type,
		ks.color_code AS kpi_status_color_code,

		krs.title_en AS kpi_review_status_title_en,
		krs.title_ar AS kpi_review_status_title_ar,
		krs.type AS kpi_review_status_type,
		krs.color_code AS kpi_review_status_color_code,

		rf.title_en AS review_frequency_title_en,
		rf.title_ar AS review_frequency_title_ar,
		rf.type AS review_frequency_type,

		ct.title_en AS calculation_type_title_en,
		ct.title_ar AS calculation_type_title_ar,
		ct.type AS calculation_type_type,

		u.title_en AS unit_title_en,
		u.title_ar AS unit_title_ar,
		u.code_en AS unit_code_en,
		u.code_ar AS unit_code_ar,
		u.color_code AS unit_color_code,

		kal.title_en AS kpi_achievement_logic_title_en,
		kal.title_ar AS kpi_achievement_logic_title_ar,
		kal.type AS kpi_achievement_logic_type,

		ru.id AS responsible_user_id,
		ru.first_name_en AS responsible_user_first_name_en,
		ru.last_name_en AS responsible_user_last_name_en,
		ru.first_name_ar AS responsible_user_first_name_ar,
		ru.last_name_ar AS responsible_user_last_name_ar,
		ru.initial AS responsible_user_initial,
		ru.email AS responsible_user_email,
		ru.color AS responsible_user_color,
		ru.profile_image_path AS responsible_user_profile_image_path,

		cu.id AS created_user_id,
		cu.first_name_en AS created_user_first_name_en,
		cu.first_name_ar AS created_user_first_name_ar,
		cu.last_name_en AS created_user_last_name_en,
		cu.last_name_ar AS created_user_last_name_ar,
		cu.initial AS created_user_initial,
		cu.color AS created_user_color,
		cu.email AS created_user_email,
		cu.profile_image_path AS created_user_profile_image_path,

		uu.id AS updated_user_id,
		uu.first_name_en AS updated_user_first_name_en,
		uu.first_name_ar AS updated_user_first_name_ar,
		uu.last_name_en AS updated_user_last_name_en,
		uu.last_name_ar AS updated_user_last_name_ar,
		uu.initial AS updated_user_initial,
		uu.color AS updated_user_color,
		uu.email AS updated_user_email,
		uu.profile_image_path AS updated_user_profile_image_path
	`)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize)

	if err := query.Find(&kpiRegister).Error; err != nil {
		return nil, err
	}

	totalPages := 1
	totalPages = int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	return &dtos.KPIRegisterListResponseDTO{
		Items: kpiRegister,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}, nil
}

func (repo *KPIDashboardRepository) GetTopKPIs(params *dtos.KPIDashboardQueryParamsDTO) ([]dtos.KPIRegisterListDTO, error) {

	results := []dtos.KPIRegisterListDTO{}

	subQuery := `
		SELECT 
			kpi_register_id,
			MAX(CASE WHEN rn = 1 THEN value END) AS latest_value,
			MAX(CASE WHEN rn = 2 THEN value END) AS prev_value
		FROM (
			SELECT 
				kpi_register_id,
				value,
				ROW_NUMBER() OVER (
					PARTITION BY kpi_register_id 
					ORDER BY reviewed_date DESC, created_at DESC
				) AS rn
			FROM kpi_value_updates
			WHERE reviewed_date IS NOT NULL
	`

	args := []interface{}{}

	if params != nil && params.StartDate != nil && params.EndDate != nil {
		subQuery += " AND reviewed_date BETWEEN ? AND ?"
		args = append(args, *params.StartDate, *params.EndDate)
	}

	subQuery += `) t GROUP BY kpi_register_id`

	query := repo.db.Table("kpi_register AS kr").
		Select(`
			kr.id,
			kr.reference_number,
			kr.title,
			kr.type,
			kr.description,
			kr.kpi_category_id,
			kr.responsible_user_id,
			kr.kpi_status_id,
			kr.kpi_review_status_id,
			kr.target_value,
			kr.remarks,

			CAST(ROUND(t.latest_value, 2) AS DECIMAL(10,2)) AS actual_value,

			CAST(
				ROUND(
					CASE 
						WHEN t.latest_value IS NOT NULL AND kr.target_value IS NOT NULL THEN
							CASE 
								WHEN kal.type = 'lower' THEN 
									CASE WHEN t.latest_value = 0 THEN 0 ELSE (kr.target_value / t.latest_value) * 100 END
								ELSE 
									(t.latest_value / kr.target_value) * 100
							END
					END, 2)
			AS DECIMAL(10,2)) AS achievement,

			CAST(ROUND(t.prev_value, 2) AS DECIMAL(10,2)) AS previous_value,

			ROUND(
				CASE 
					WHEN t.latest_value IS NOT NULL AND t.prev_value IS NOT NULL
					THEN ((t.latest_value - t.prev_value) / NULLIF(t.prev_value, 0)) * 100
					ELSE 0
				END, 2
			) AS trend,

			CONVERT(VARCHAR, kr.start_date, 23) AS start_date,
			CONVERT(VARCHAR, kr.end_date, 23) AS end_date,

			kc.title_en AS kpi_category_title_en,
			kc.title_ar AS kpi_category_title_ar,

			ks.title_en AS kpi_status_title_en,
			ks.title_ar AS kpi_status_title_ar,
			ks.type AS kpi_status_type,
			ks.color_code AS kpi_status_color_code,

			krs.title_en AS kpi_review_status_title_en,
			krs.title_ar AS kpi_review_status_title_ar,
			krs.type AS kpi_review_status_type,
			krs.color_code AS kpi_review_status_color_code,

			rf.title_en AS review_frequency_title_en,
			rf.title_ar AS review_frequency_title_ar,
			rf.type AS review_frequency_type,

			ls.title_en AS status_title_en,
			ls.title_ar AS status_title_ar,
			ls.type AS status_type,
			ls.color_code AS status_color_code,

			ct.title_en AS calculation_type_title_en,
			ct.title_ar AS calculation_type_title_ar,
			ct.type AS calculation_type_type,

			u.title_en AS unit_title_en,
			u.title_ar AS unit_title_ar,
			u.code_en AS unit_code_en,
			u.code_ar AS unit_code_ar,
			u.color_code AS unit_color_code,

			ru.id AS responsible_user_id,
			ru.first_name_en AS responsible_user_first_name_en,
			ru.last_name_en AS responsible_user_last_name_en,
			ru.first_name_ar AS responsible_user_first_name_ar,
			ru.last_name_ar AS responsible_user_last_name_ar,
			ru.initial AS responsible_user_initial,
			ru.email AS responsible_user_email,
			ru.color AS responsible_user_color,
			ru.profile_image_path AS responsible_user_profile_image_path,

			cu.id AS created_user_id,
			cu.first_name_en AS created_user_first_name_en,
			cu.first_name_ar AS created_user_first_name_ar,
			cu.last_name_en AS created_user_last_name_en,
			cu.last_name_ar AS created_user_last_name_ar,
			cu.initial AS created_user_initial,
			cu.color AS created_user_color,
			cu.email AS created_user_email,
			cu.profile_image_path AS created_user_profile_image_path,

			uu.id AS updated_user_id,
			uu.first_name_en AS updated_user_first_name_en,
			uu.first_name_ar AS updated_user_first_name_ar,
			uu.last_name_en AS updated_user_last_name_en,
			uu.last_name_ar AS updated_user_last_name_ar,
			uu.initial AS updated_user_initial,
			uu.color AS updated_user_color,
			uu.email AS updated_user_email,
			uu.profile_image_path AS updated_user_profile_image_path
		`).
		Joins("LEFT JOIN ("+subQuery+") AS t ON t.kpi_register_id = kr.id", args...).
		Joins("LEFT JOIN kpi_categories kc ON kc.id = kr.kpi_category_id").
		Joins("LEFT JOIN kpi_statuses ks ON ks.id = kr.kpi_status_id").
		Joins("LEFT JOIN kpi_review_statuses krs ON krs.id = kr.kpi_review_status_id").
		Joins("LEFT JOIN review_frequencies rf ON rf.id = kr.review_frequency_id").
		Joins("LEFT JOIN calculation_types ct ON ct.id = kr.calculation_type_id").
		Joins("LEFT JOIN units u ON u.id = kr.unit_id").
		Joins("LEFT JOIN users ru ON ru.id = kr.responsible_user_id").
		Joins("LEFT JOIN users cu ON cu.id = kr.created_by_id").
		Joins("LEFT JOIN users uu ON uu.id = kr.updated_by_id").
		Joins("LEFT JOIN kpi_achievement_logics kal ON kal.id = kr.kpi_achievement_logic_id").
		Joins("LEFT JOIN lifecycle_statuses ls ON ls.id = kr.status_id")

	if params != nil && params.IsKRI != nil {
		query = query.Where("kr.is_kri = ?", *params.IsKRI)
	}

	if params != nil && params.Filters != nil {
		for col, values := range params.Filters {
			if len(values) > 0 {
				query = query.Where("kr."+col+" IN (?)", values)
			}
		}
	}

	if params != nil && len(params.QuickFilter) > 0 {
		query = query.Where("kr.type IN (?)", params.QuickFilter)
	}

	if params != nil && params.PeriodType != nil {
		var freqTypes []string

		switch *params.PeriodType {
		case "this-month":
			freqTypes = []string{"daily", "weekly", "monthly"}
		case "this-quarter":
			freqTypes = []string{"daily", "weekly", "monthly", "quarterly"}
		case "this-half-year":
			freqTypes = []string{"daily", "weekly", "monthly", "quarterly", "half-yearly"}
		case "this-year", "custom-range":
			freqTypes = []string{"daily", "weekly", "monthly", "quarterly", "half-yearly", "annually"}
		}

		if len(freqTypes) > 0 {
			query = query.Where("rf.type IN ? OR rf.type IS NULL", freqTypes)
		}
	}

	query = query.
		Where("t.latest_value IS NOT NULL AND kr.target_value IS NOT NULL").
		Where("ks.is_not_achieved = ?", false).
		Order("achievement DESC").
		Limit(5)

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *KPIDashboardRepository) GetBottomKPIs(params *dtos.KPIDashboardQueryParamsDTO) ([]dtos.KPIRegisterListDTO, error) {

	results := []dtos.KPIRegisterListDTO{}

	subQuery := `
		SELECT 
			kpi_register_id,
			MAX(CASE WHEN rn = 1 THEN value END) AS latest_value,
			MAX(CASE WHEN rn = 2 THEN value END) AS prev_value
		FROM (
			SELECT 
				kpi_register_id,
				value,
				ROW_NUMBER() OVER (
					PARTITION BY kpi_register_id 
					ORDER BY reviewed_date DESC, created_at DESC
				) AS rn
			FROM kpi_value_updates
			WHERE reviewed_date IS NOT NULL
	`

	args := []interface{}{}

	if params != nil && params.StartDate != nil && params.EndDate != nil {
		subQuery += " AND reviewed_date BETWEEN ? AND ?"
		args = append(args, *params.StartDate, *params.EndDate)
	}

	subQuery += `) t GROUP BY kpi_register_id`

	query := repo.db.Table("kpi_register AS kr").
		Select(`
			kr.id,
			kr.reference_number,
			kr.title,
			kr.type,
			kr.description,
			kr.kpi_category_id,
			kr.responsible_user_id,
			kr.kpi_status_id,
			kr.kpi_review_status_id,
			kr.target_value,
			kr.remarks,

			CAST(ROUND(t.latest_value, 2) AS DECIMAL(10,2)) AS actual_value,

			CAST(
				ROUND(
					CASE 
						WHEN t.latest_value IS NOT NULL AND kr.target_value IS NOT NULL THEN
							CASE 
								WHEN kal.type = 'lower' THEN 
									CASE WHEN t.latest_value = 0 THEN 0 ELSE (kr.target_value / t.latest_value) * 100 END
								ELSE 
									(t.latest_value / kr.target_value) * 100
							END
					END, 2)
			AS DECIMAL(10,2)) AS achievement,

			CAST(ROUND(t.prev_value, 2) AS DECIMAL(10,2)) AS previous_value,

			ROUND(
				CASE 
					WHEN t.latest_value IS NOT NULL AND t.prev_value IS NOT NULL
					THEN ((t.latest_value - t.prev_value) / NULLIF(t.prev_value, 0)) * 100
					ELSE 0
				END, 2
			) AS trend,

			CONVERT(VARCHAR, kr.start_date, 23) AS start_date,
			CONVERT(VARCHAR, kr.end_date, 23) AS end_date,

			kc.title_en AS kpi_category_title_en,
			kc.title_ar AS kpi_category_title_ar,

			ks.title_en AS kpi_status_title_en,
			ks.title_ar AS kpi_status_title_ar,
			ks.type AS kpi_status_type,
			ks.color_code AS kpi_status_color_code,

			krs.title_en AS kpi_review_status_title_en,
			krs.title_ar AS kpi_review_status_title_ar,
			krs.type AS kpi_review_status_type,
			krs.color_code AS kpi_review_status_color_code,

			rf.title_en AS review_frequency_title_en,
			rf.title_ar AS review_frequency_title_ar,
			rf.type AS review_frequency_type,

			ct.title_en AS calculation_type_title_en,
			ct.title_ar AS calculation_type_title_ar,
			ct.type AS calculation_type_type,

			u.title_en AS unit_title_en,
			u.title_ar AS unit_title_ar,
			u.code_en AS unit_code_en,
			u.code_ar AS unit_code_ar,
			u.color_code AS unit_color_code,

			ru.id AS responsible_user_id,
			ru.first_name_en AS responsible_user_first_name_en,
			ru.last_name_en AS responsible_user_last_name_en,
			ru.first_name_ar AS responsible_user_first_name_ar,
			ru.last_name_ar AS responsible_user_last_name_ar,
			ru.initial AS responsible_user_initial,
			ru.email AS responsible_user_email,
			ru.color AS responsible_user_color,
			ru.profile_image_path AS responsible_user_profile_image_path
		`).
		Joins("LEFT JOIN ("+subQuery+") AS t ON t.kpi_register_id = kr.id", args...).
		Joins("LEFT JOIN kpi_categories kc ON kc.id = kr.kpi_category_id").
		Joins("LEFT JOIN kpi_statuses ks ON ks.id = kr.kpi_status_id").
		Joins("LEFT JOIN kpi_review_statuses krs ON krs.id = kr.kpi_review_status_id").
		Joins("LEFT JOIN review_frequencies rf ON rf.id = kr.review_frequency_id").
		Joins("LEFT JOIN calculation_types ct ON ct.id = kr.calculation_type_id").
		Joins("LEFT JOIN units u ON u.id = kr.unit_id").
		Joins("LEFT JOIN users ru ON ru.id = kr.responsible_user_id").
		Joins("LEFT JOIN kpi_achievement_logics kal ON kal.id = kr.kpi_achievement_logic_id")

	if params != nil && params.IsKRI != nil {
		query = query.Where("kr.is_kri = ?", *params.IsKRI)
	}

	if params != nil && params.Filters != nil {
		for col, values := range params.Filters {
			if len(values) > 0 {
				query = query.Where("kr."+col+" IN (?)", values)
			}
		}
	}

	if params != nil && len(params.QuickFilter) > 0 {
		query = query.Where("kr.type IN (?)", params.QuickFilter)
	}

	if params != nil && params.PeriodType != nil {
		var freqTypes []string

		switch *params.PeriodType {
		case "this-month":
			freqTypes = []string{"daily", "weekly", "monthly"}
		case "this-quarter":
			freqTypes = []string{"daily", "weekly", "monthly", "quarterly"}
		case "this-half-year":
			freqTypes = []string{"daily", "weekly", "monthly", "quarterly", "half-yearly"}
		case "this-year", "custom-range":
			freqTypes = []string{"daily", "weekly", "monthly", "quarterly", "half-yearly", "annually"}
		}

		if len(freqTypes) > 0 {
			query = query.Where("rf.type IN ? OR rf.type IS NULL", freqTypes)
		}
	}

	query = query.
		Where("t.latest_value IS NOT NULL AND kr.target_value IS NOT NULL").
		Where("ks.is_not_achieved = ?", true).
		Order("achievement ASC").
		Limit(5)

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}
