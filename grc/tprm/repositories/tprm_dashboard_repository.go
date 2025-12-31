package repositories

import (
	"jethur/grc/models"
	tprmDTOs "jethur/grc/tprm/dtos"
	"time"

	"gorm.io/gorm"
)

type TPRMDashboardRepositoryInterface interface {
	GetThirdPartyAssessmentProgress(params *tprmDTOs.TPRMDashboardQueryParamsDTO) (*tprmDTOs.TPRMThirdPartyAssessmentProgressDTO, error)
	GetThirdPartyRiskLevelDistribution(params *tprmDTOs.TPRMDashboardQueryParamsDTO) ([]tprmDTOs.TPRMRiskLevelDTO, error)
	GetAssessmentInsights(params *tprmDTOs.TPRMDashboardQueryParamsDTO) (*tprmDTOs.TPRMAssessmentInsightsDTO, error)
	GetAssessmentRiskLevelDistribution(params *tprmDTOs.TPRMDashboardQueryParamsDTO) ([]tprmDTOs.TPRMRiskLevelDTO, error)
	GetSystemStatistics(params *tprmDTOs.TPRMDashboardQueryParamsDTO) (*tprmDTOs.TPRMSystemStatisticsDTO, error)
	GetOngoingAssessments(params *tprmDTOs.TPRMDashboardQueryParamsDTO) ([]tprmDTOs.TPRMOngoingAssessmentDTO, error)
	GetThirdPartyServiceExpiryCounts(params *tprmDTOs.TPRMDashboardQueryParamsDTO) (*tprmDTOs.TPRMServiceExpiryCountsDTO, error)
	GetAssessmentByTemplateType(params *tprmDTOs.TPRMDashboardQueryParamsDTO) ([]tprmDTOs.TPRMAssessmentByTemplateDTO, error)
	GetThirdPartyByType(params *tprmDTOs.TPRMDashboardQueryParamsDTO) ([]tprmDTOs.TPRMThirdPartyByTypeDTO, error)
	GetAssessmentRiskMonthlyCounts(startDate time.Time, endDate time.Time, params *tprmDTOs.TPRMDashboardQueryParamsDTO) ([]tprmDTOs.TPRMAssessmentRiskMonthlyRowDTO, error)
}

type TPRMDashboardRepository struct {
	db *gorm.DB
}

func NewTPRMDashboardRepository(db *gorm.DB) TPRMDashboardRepositoryInterface {
	return &TPRMDashboardRepository{
		db: db,
	}
}

func (repo *TPRMDashboardRepository) GetThirdPartyAssessmentProgress(params *tprmDTOs.TPRMDashboardQueryParamsDTO) (*tprmDTOs.TPRMThirdPartyAssessmentProgressDTO, error) {

	result := tprmDTOs.TPRMThirdPartyAssessmentProgressDTO{}

	subQuery := repo.db.
		Table("third_party_register tpr").
		Select(`
			tpr.id AS third_party_id,
			MAX(CASE 
				WHEN tpas.type = 'approved' THEN 1 
				ELSE 0 
			END) AS is_assessed,
			MAX(CASE 
				WHEN tpas.type IN ('new', 'in-progress', 'submitted') THEN 1 
				ELSE 0 
			END) AS is_started
		`).
		Joins(`
			LEFT JOIN third_party_assessments tpa
				ON tpa.third_party_register_id = tpr.id
		`).
		Joins(`
			LEFT JOIN third_party_assessment_statuses tpas
				ON tpas.id = tpa.third_party_assessment_status_id
		`).
		Group("tpr.id")

	if params != nil && params.DateFrom != "" && params.DateTo != "" {
		subQuery = subQuery.Where(
			"tpr.created_at BETWEEN ? AND ?",
			params.DateFrom+" 00:00:00",
			params.DateTo+" 23:59:59",
		)
	}

	if params != nil {
		for col, values := range params.Filters {
			if len(values) > 0 {
				subQuery = subQuery.Where("tpr."+col+" IN (?)", values)
			}
		}
	}

	query := repo.db.
		Table("(?) AS tp_status", subQuery).
		Select(`
			COUNT(third_party_id) AS total,
			SUM(CASE WHEN is_assessed = 1 THEN 1 ELSE 0 END) AS assessed,
			SUM(
				CASE 
					WHEN is_assessed = 0 AND is_started = 1 THEN 1 
					ELSE 0 
				END
			) AS in_progress
		`)

	if err := query.Scan(&result).Error; err != nil {
		return nil, err
	}

	result.NotAssessed = result.Total - result.Assessed - result.InProgress

	return &result, nil
}

func (repo *TPRMDashboardRepository) GetThirdPartyRiskLevelDistribution(params *tprmDTOs.TPRMDashboardQueryParamsDTO) ([]tprmDTOs.TPRMRiskLevelDTO, error) {

	var allRiskLevels []models.RiskLevel
	if err := repo.db.
		Table("risk_levels").
		Where("status = 'active'").
		Order("id").
		Find(&allRiskLevels).Error; err != nil {
		return nil, err
	}

	var counts []tprmDTOs.TPRMRiskLevelCountDTO

	query := repo.db.
		Table("third_party_register tpr").
		Select("tpr.risk_level_id AS risk_level_id, COUNT(*) AS count")

	if params != nil && params.DateFrom != "" && params.DateTo != "" {
		query = query.Where(
			"tpr.created_at BETWEEN ? AND ?",
			params.DateFrom+" 00:00:00",
			params.DateTo+" 23:59:59",
		)
	}

	if params != nil {
		for col, values := range params.Filters {
			if len(values) > 0 {
				query = query.Where("tpr."+col+" IN (?)", values)
			}
		}
	}

	query = query.Group("tpr.risk_level_id")

	if err := query.Scan(&counts).Error; err != nil {
		return nil, err
	}

	countMap := map[uint]int64{}
	for _, c := range counts {
		countMap[c.RiskLevelID] = c.Count
	}

	results := []tprmDTOs.TPRMRiskLevelDTO{}

	for _, rl := range allRiskLevels {
		results = append(results, tprmDTOs.TPRMRiskLevelDTO{
			ID:        rl.ID,
			TitleEn:   rl.TitleEn,
			TitleAr:   rl.TitleAr,
			ColorCode: rl.ColorCode,
			Count:     countMap[rl.ID],
		})
	}

	return results, nil
}

func (repo *TPRMDashboardRepository) GetAssessmentInsights(params *tprmDTOs.TPRMDashboardQueryParamsDTO) (*tprmDTOs.TPRMAssessmentInsightsDTO, error) {

	result := tprmDTOs.TPRMAssessmentInsightsDTO{}

	query := repo.db.Table("third_party_assessments tpa").
		Select(`
			COUNT(tpa.id) AS total_assessments,

			SUM(
				CASE 
					WHEN tpas.type IN ('new', 'in-progress', 'submitted') 
					THEN 1 ELSE 0 
				END
			) AS in_progress,

			SUM(
				CASE 
					WHEN tpas.type = 'approved'
					THEN 1 ELSE 0 
				END
			) AS completed,

			SUM(
				CASE 
					WHEN tpa.due_date IS NOT NULL
					 AND tpa.due_date < CAST(GETDATE() AS DATE)
					 AND tpas.type != 'approved'
					THEN 1 ELSE 0
				END
			) AS overdue
		`).
		Joins("LEFT JOIN third_party_assessment_statuses tpas ON tpas.id = tpa.third_party_assessment_status_id").
		Joins("LEFT JOIN third_party_register tpr ON tpr.id = tpa.third_party_register_id")

	if params != nil && params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("tpa.created_at BETWEEN ? AND ?",
			params.DateFrom+" 00:00:00",
			params.DateTo+" 23:59:59",
		)
	}

	for col, vals := range params.Filters {
		if len(vals) > 0 {
			query = query.Where("tpr."+col+" IN (?)", vals)
		}
	}

	if err := query.Scan(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

func (repo *TPRMDashboardRepository) GetAssessmentRiskLevelDistribution(params *tprmDTOs.TPRMDashboardQueryParamsDTO) ([]tprmDTOs.TPRMRiskLevelDTO, error) {

	items := []tprmDTOs.TPRMRiskLevelDTO{}

	query := repo.db.Table("third_party_assessments tpa").
		Select(`
			rl.id,
			rl.title_en,
			rl.title_ar,
			rl.color_code,
			COUNT(tpa.id) AS count
		`).
		Joins("JOIN risk_levels rl ON rl.id = tpa.risk_level_id").
		Joins("JOIN third_party_register tpr ON tpr.id = tpa.third_party_register_id").
		Where("rl.status = 'active'").
		Group(`
			rl.id,
			rl.title_en,
			rl.title_ar,
			rl.color_code
		`)

	if params != nil && params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("tpa.created_at BETWEEN ? AND ?",
			params.DateFrom+" 00:00:00",
			params.DateTo+" 23:59:59",
		)
	}

	for col, vals := range params.Filters {
		if len(vals) > 0 {
			query = query.Where("tpr."+col+" IN (?)", vals)
		}
	}

	if err := query.Scan(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *TPRMDashboardRepository) GetSystemStatistics(params *tprmDTOs.TPRMDashboardQueryParamsDTO) (*tprmDTOs.TPRMSystemStatisticsDTO, error) {

	var totalQuestions int64
	var templates int64

	q1 := repo.db.Table("third_party_question_libraries")

	if params != nil && params.DateFrom != "" && params.DateTo != "" {
		q1 = q1.Where("created_at BETWEEN ? AND ?",
			params.DateFrom+" 00:00:00",
			params.DateTo+" 23:59:59",
		)
	}

	if err := q1.Count(&totalQuestions).Error; err != nil {
		return nil, err
	}

	q2 := repo.db.Table("question_sets")

	if params != nil && params.DateFrom != "" && params.DateTo != "" {
		q2 = q2.Where("created_at BETWEEN ? AND ?",
			params.DateFrom+" 00:00:00",
			params.DateTo+" 23:59:59",
		)
	}

	if err := q2.Count(&templates).Error; err != nil {
		return nil, err
	}

	return &tprmDTOs.TPRMSystemStatisticsDTO{
		TotalQuestions: totalQuestions,
		Templates:      templates,
	}, nil
}

func (repo *TPRMDashboardRepository) GetOngoingAssessments(params *tprmDTOs.TPRMDashboardQueryParamsDTO) ([]tprmDTOs.TPRMOngoingAssessmentDTO, error) {

	items := []tprmDTOs.TPRMOngoingAssessmentDTO{}

	query := repo.db.Table("third_party_assessments tpa").
		Select(`
			tpa.reference_number AS id,
			tpa.title,
			ou.title_en AS organization_unit_title_en,
			ou.title_ar AS organization_unit_title_ar,
			tpa.due_date
		`).
		Joins("JOIN third_party_assessment_statuses tpas ON tpas.id = tpa.third_party_assessment_status_id").
		Joins("JOIN third_party_register tpr ON tpr.id = tpa.third_party_register_id").
		Joins("LEFT JOIN organization_units ou ON ou.id = tpr.organization_unit_id").
		Where("tpas.type IN ('new', 'in-progress', 'submitted')").
		Where("tpa.due_date IS NOT NULL").
		Order("tpa.due_date ASC")

	if params != nil && params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("tpa.created_at BETWEEN ? AND ?",
			params.DateFrom+" 00:00:00",
			params.DateTo+" 23:59:59",
		)
	}

	for col, vals := range params.Filters {
		if len(vals) > 0 {
			query = query.Where("tpr."+col+" IN (?)", vals)
		}
	}

	if err := query.Scan(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *TPRMDashboardRepository) GetThirdPartyServiceExpiryCounts(params *tprmDTOs.TPRMDashboardQueryParamsDTO) (*tprmDTOs.TPRMServiceExpiryCountsDTO, error) {

	result := tprmDTOs.TPRMServiceExpiryCountsDTO{}

	query := repo.db.Table("third_party_register tpr").
		Select(`
			SUM(
				CASE 
					WHEN tpr.service_end_date >= CAST(GETDATE() AS DATE)
					 AND tpr.service_end_date <= DATEADD(day, 7, CAST(GETDATE() AS DATE))
					THEN 1 ELSE 0
				END
			) AS due_within_7d,

			SUM(
				CASE 
					WHEN tpr.service_end_date < CAST(GETDATE() AS DATE)
					 AND tpr.service_end_date >= DATEADD(day, -30, CAST(GETDATE() AS DATE))
					THEN 1 ELSE 0
				END
			) AS overdue_lt_30d,

			SUM(
				CASE 
					WHEN tpr.service_end_date < DATEADD(day, -30, CAST(GETDATE() AS DATE))
					THEN 1 ELSE 0
				END
			) AS overdue_gt_30d
		`).
		Where("tpr.service_end_date IS NOT NULL")

	if params != nil && params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("tpr.created_at BETWEEN ? AND ?",
			params.DateFrom+" 00:00:00",
			params.DateTo+" 23:59:59",
		)
	}

	for col, vals := range params.Filters {
		if len(vals) > 0 {
			query = query.Where("tpr."+col+" IN (?)", vals)
		}
	}

	if err := query.Scan(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

func (repo *TPRMDashboardRepository) GetAssessmentByTemplateType(params *tprmDTOs.TPRMDashboardQueryParamsDTO) ([]tprmDTOs.TPRMAssessmentByTemplateDTO, error) {

	results := []tprmDTOs.TPRMAssessmentByTemplateDTO{}

	query := repo.db.
		Table("third_party_assessments tpa").
		Select(`
			qs.title_en,
			qs.title_ar,
			COUNT(tpa.id) AS count
		`).
		Joins("JOIN question_sets qs ON qs.id = tpa.question_set_id").
		Joins("JOIN third_party_register tpr ON tpr.id = tpa.third_party_register_id").
		Group(`
			qs.id,
			qs.title_en,
			qs.title_ar
		`).
		Order("count DESC")

	if params != nil && params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("tpa.created_at BETWEEN ? AND ?",
			params.DateFrom+" 00:00:00",
			params.DateTo+" 23:59:59",
		)
	}

	for col, vals := range params.Filters {
		if len(vals) > 0 {
			query = query.Where("tpr."+col+" IN (?)", vals)
		}
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *TPRMDashboardRepository) GetThirdPartyByType(params *tprmDTOs.TPRMDashboardQueryParamsDTO) ([]tprmDTOs.TPRMThirdPartyByTypeDTO, error) {

	results := []tprmDTOs.TPRMThirdPartyByTypeDTO{}

	query := repo.db.
		Table("third_party_register tpr").
		Select(`
			tpt.title_en,
			tpt.title_ar,
			COUNT(tpr.id) AS count
		`).
		Joins(`
			JOIN third_party_types tpt
				ON tpt.id = tpr.third_party_type_id
		`).
		Where("tpt.status = 'active'").
		Group(`
			tpt.id,
			tpt.title_en,
			tpt.title_ar
		`).
		Order("count DESC")

	if params != nil && params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("tpr.created_at BETWEEN ? AND ?",
			params.DateFrom+" 00:00:00",
			params.DateTo+" 23:59:59",
		)
	}

	for col, vals := range params.Filters {
		if len(vals) > 0 {
			query = query.Where("tpr."+col+" IN (?)", vals)
		}
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *TPRMDashboardRepository) GetAssessmentRiskMonthlyCounts(startDate time.Time, endDate time.Time, params *tprmDTOs.TPRMDashboardQueryParamsDTO) ([]tprmDTOs.TPRMAssessmentRiskMonthlyRowDTO, error) {

	rows := []tprmDTOs.TPRMAssessmentRiskMonthlyRowDTO{}

	query := repo.db.
		Table("third_party_assessments tpa").
		Select(`
			FORMAT(tpa.created_at, 'yyyy-MM') AS month,
			rl.type AS risk_type,
			COUNT(tpa.id) AS count
		`).
		Joins("JOIN risk_levels rl ON rl.id = tpa.risk_level_id AND rl.status = 'active'").
		Joins("JOIN third_party_register tpr ON tpr.id = tpa.third_party_register_id").
		Where("tpa.risk_level_id IS NOT NULL").
		Where("tpa.created_at >= ? AND tpa.created_at < ?", startDate, endDate).
		Group(`
			FORMAT(tpa.created_at, 'yyyy-MM'),
			rl.type
		`)

	for col, vals := range params.Filters {
		if len(vals) > 0 {
			query = query.Where("tpr."+col+" IN (?)", vals)
		}
	}

	if err := query.Scan(&rows).Error; err != nil {
		return nil, err
	}

	return rows, nil
}
