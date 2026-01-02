package repositories

import (
	"jethur/grc/models"
	reportDTOs "jethur/grc/reports/dtos"
	"time"

	"gorm.io/gorm"
)

type TPRMReportRepositoryInterface interface {
	GetThirdPartyByType(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ThirdPartyByTypeDTO, error)
	GetThirdPartyByCountry(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ThirdPartyByCountryDTO, error)
	GetThirdPartyByDepartment(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ThirdPartyByDepartmentDTO, error)
	GetThirdPartyByCriticality(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ThirdPartyByCriticalityDTO, error)
	GetThirdPartyByLatestAssessmentStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.TPAssessmentsByStatusDTO, error)
	GetThirdPartyByRiskRating(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ThirdPartyByRiskRatingDTO, error)
	GetThirdPartyServiceExpiry(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ThirdPartyServiceExpiryDTO, error)

	GetAssessmentsByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.TPAssessmentsByStatusDTO, error)
	GetAssessmentsByRiskRating(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.TPAssessmentsByRiskRatingDTO, error)
	GetAssessmentsByThirdPartyCategory(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.TPAssessmentsByCategoryDTO, error)
	GetAssessmentsByTemplate(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.TPAssessmentsByTemplateDTO, error)
	GetOverdueAssessments(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.TPOverdueAssessmentsDTO, error)
}

type TPRMReportRepository struct {
	db *gorm.DB
}

func NewTPRMReportRepository(db *gorm.DB) TPRMReportRepositoryInterface {
	return &TPRMReportRepository{
		db: db,
	}
}

func (repo *TPRMReportRepository) GetThirdPartyByType(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ThirdPartyByTypeDTO, error) {

	results := []*reportDTOs.ThirdPartyByTypeDTO{}

	query := repo.db.Table("third_party_types tpt").
		Select(`
			tpt.id AS type_id,
			tpt.title_en,
			tpt.title_ar,
			COUNT(tpr.id) AS count
		`).
		Joins(`
			LEFT JOIN third_party_register tpr
			ON tpr.third_party_type_id = tpt.id
		`).
		Group("tpt.id, tpt.title_en, tpt.title_ar").
		Order("COUNT(tpr.id) DESC")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("tpr.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *TPRMReportRepository) GetThirdPartyByCountry(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ThirdPartyByCountryDTO, error) {

	results := []*reportDTOs.ThirdPartyByCountryDTO{}

	query := repo.db.Table("countries c").
		Select(`
			c.id AS country_id,
			c.title_en,
			c.title_ar,
			COUNT(tpr.id) AS count
		`).
		Joins(`
			LEFT JOIN third_party_register tpr
			ON tpr.country_id = c.id
		`).
		Group("c.id, c.title_en, c.title_ar").
		Having("COUNT(tpr.id) > 0").
		Order("COUNT(tpr.id) DESC")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("tpr.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *TPRMReportRepository) GetThirdPartyByDepartment(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ThirdPartyByDepartmentDTO, error) {

	results := []*reportDTOs.ThirdPartyByDepartmentDTO{}

	query := repo.db.Table("organization_units ou").
		Select(`
			ou.id AS department_id,
			ou.title_en,
			ou.title_ar,
			ou.color_code,
			COUNT(tpr.id) AS count
		`).
		Joins(`
			LEFT JOIN third_party_register tpr
			ON tpr.organization_unit_id = ou.id
		`).
		Group("ou.id, ou.title_en, ou.title_ar, ou.color_code").
		Order("COUNT(tpr.id) DESC")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("tpr.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *TPRMReportRepository) GetThirdPartyByCriticality(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ThirdPartyByCriticalityDTO, error) {

	results := []*reportDTOs.ThirdPartyByCriticalityDTO{}

	query := repo.db.Table("severity_levels sl").
		Select(`
			sl.id AS severity_id,
			sl.title_en,
			sl.title_ar,
			sl.color_code,
			COUNT(tpr.id) AS count
		`).
		Joins(`
			LEFT JOIN third_party_register tpr
			ON tpr.severity_level_id = sl.id
		`).
		Group("sl.id, sl.title_en, sl.title_ar, sl.color_code").
		Order("COUNT(tpr.id) DESC")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("tpr.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *TPRMReportRepository) GetThirdPartyByLatestAssessmentStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.TPAssessmentsByStatusDTO, error) {

	results := []*reportDTOs.TPAssessmentsByStatusDTO{}

	latestExpr := gorm.Expr(`
        (
            SELECT 
                tpa.id,
                tpa.third_party_assessment_status_id,
                tpa.created_at,
                ROW_NUMBER() OVER (
                    PARTITION BY tpa.third_party_register_id 
                    ORDER BY tpa.created_at DESC
                ) AS rn
            FROM third_party_assessments tpa
        ) AS latest
    `)

	query := repo.db.Table("third_party_assessment_statuses s").
		Select(`
            s.id AS status_id,
            s.title_en,
            s.title_ar,
            s.color_code,
            COUNT(latest.id) AS count
        `).
		Joins("LEFT JOIN ? ON latest.third_party_assessment_status_id = s.id AND latest.rn = 1", latestExpr).
		Where("s.type NOT IN ('in-progress','new')")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("latest.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	query = query.Group("s.id, s.title_en, s.title_ar, s.color_code").
		Order("COUNT(latest.id) DESC")

	err := query.Scan(&results).Error
	return results, err
}

func (repo *TPRMReportRepository) GetThirdPartyByRiskRating(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ThirdPartyByRiskRatingDTO, error) {

	results := []*reportDTOs.ThirdPartyByRiskRatingDTO{}

	latestExpr := gorm.Expr(`
        (
            SELECT 
                tpa.id,
                tpa.risk_level_id,
                tpa.created_at,
                ROW_NUMBER() OVER (
                    PARTITION BY tpa.third_party_register_id 
                    ORDER BY tpa.created_at DESC
                ) AS rn
            FROM third_party_assessments tpa
        ) AS latest
    `)

	query := repo.db.Table("risk_levels rl").
		Select(`
            rl.id AS risk_level_id,
            rl.title_en,
            rl.title_ar,
            rl.color_code,
            COUNT(latest.id) AS count
        `).
		Joins(`
            LEFT JOIN ? 
            ON latest.risk_level_id = rl.id
            AND latest.rn = 1
        `, latestExpr)

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("latest.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	query = query.Group("rl.id, rl.title_en, rl.title_ar, rl.color_code, rl.score_from").
		Order("rl.score_from DESC")

	err := query.Scan(&results).Error
	return results, err
}

func (repo *TPRMReportRepository) GetThirdPartyServiceExpiry(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ThirdPartyServiceExpiryDTO, error) {

	results := []*reportDTOs.ThirdPartyServiceExpiryDTO{}

	now := time.Now()
	exp30 := now.AddDate(0, 0, 30)
	exp60 := now.AddDate(0, 0, 60)

	queries := []struct {
		Category  string
		Condition string
		TimeArg   interface{}
	}{
		{"expired", "service_end_date < ?", now},
		{"expiring_30", "service_end_date BETWEEN ? AND ?", []time.Time{now, exp30}},
		{"expiring_60", "service_end_date BETWEEN ? AND ?", []time.Time{exp30, exp60}},
	}

	for _, q := range queries {
		var count int64

		query := repo.db.Model(&models.ThirdPartyRegister{})

		if params.DateFrom != "" && params.DateTo != "" {
			query = query.Where("created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
		}

		if times, ok := q.TimeArg.([]time.Time); ok {
			query = query.Where(q.Condition, times[0], times[1])
		} else {
			query = query.Where(q.Condition, q.TimeArg)
		}

		query.Count(&count)

		results = append(results, &reportDTOs.ThirdPartyServiceExpiryDTO{
			Category: q.Category,
			Count:    uint(count),
		})
	}

	return results, nil
}

func (repo *TPRMReportRepository) GetAssessmentsByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.TPAssessmentsByStatusDTO, error) {

	results := []*reportDTOs.TPAssessmentsByStatusDTO{}

	query := repo.db.Table("third_party_assessment_statuses s").
		Select(`
			s.id AS status_id,
			s.title_en,
			s.title_ar,
			s.color_code,
			COUNT(a.id) AS count
		`).
		Where("s.type NOT IN ('in-progress','new')").
		Joins(`
			LEFT JOIN third_party_assessments a
			ON a.third_party_assessment_status_id = s.id
		`).
		Group("s.id, s.title_en, s.title_ar, s.color_code").
		Order("COUNT(a.id) DESC")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where(
			"a.created_at BETWEEN ? AND ?",
			params.DateFrom+" 00:00:00",
			params.DateTo+" 23:59:59",
		)
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *TPRMReportRepository) GetAssessmentsByRiskRating(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.TPAssessmentsByRiskRatingDTO, error) {

	results := []*reportDTOs.TPAssessmentsByRiskRatingDTO{}

	query := repo.db.Table("risk_levels rl").
		Select(`
			rl.id AS risk_level_id,
			rl.title_en,
			rl.title_ar,
			rl.color_code,
			COUNT(a.id) AS count
		`).
		Joins(`
			LEFT JOIN third_party_assessments a
			ON a.risk_level_id = rl.id
		`).
		Group("rl.id, rl.title_en, rl.title_ar, rl.color_code,rl.score_from").
		Order("rl.score_from DESC")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("a.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *TPRMReportRepository) GetAssessmentsByThirdPartyCategory(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.TPAssessmentsByCategoryDTO, error) {

	results := []*reportDTOs.TPAssessmentsByCategoryDTO{}

	query := repo.db.Table("third_party_types tpt").
		Select(`
			tpt.id AS category_id,
			tpt.title_en,
			tpt.title_ar,
			COUNT(a.id) AS count
		`).
		Joins(`
			LEFT JOIN third_party_register tpr
			ON tpr.third_party_type_id = tpt.id
		`).
		Joins(`
			LEFT JOIN third_party_assessments a
			ON a.third_party_register_id = tpr.id
		`).
		Group("tpt.id, tpt.title_en, tpt.title_ar").
		Order("COUNT(a.id) DESC")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("a.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *TPRMReportRepository) GetAssessmentsByTemplate(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.TPAssessmentsByTemplateDTO, error) {

	results := []*reportDTOs.TPAssessmentsByTemplateDTO{}

	query := repo.db.Table("question_sets qs").
		Select(`
			qs.id AS template_id,
			qs.title_en,
			qs.title_ar,
			COUNT(a.id) AS count
		`).
		Joins(`
			LEFT JOIN third_party_assessments a
			ON a.question_set_id = qs.id
		`).
		Group("qs.id, qs.title_en, qs.title_ar").
		Order("COUNT(a.id) DESC")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("a.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *TPRMReportRepository) GetOverdueAssessments(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.TPOverdueAssessmentsDTO, error) {

	results := []*reportDTOs.TPOverdueAssessmentsDTO{}

	today := time.Now().Truncate(24 * time.Hour)
	next7 := today.AddDate(0, 0, 7)
	next14 := today.AddDate(0, 0, 14)

	checks := []struct {
		Category string
		Where    string
		Args     []interface{}
	}{
		{"overdue", "CAST(a.due_date AS date) < CAST(? AS date)", []interface{}{today}},
		{"upcoming_7", "CAST(a.due_date AS date) BETWEEN CAST(? AS date) AND CAST(? AS date)", []interface{}{today, next7}},
		{"upcoming_14", "CAST(a.due_date AS date) BETWEEN CAST(? AS date) AND CAST(? AS date)", []interface{}{next7, next14}},
	}

	for _, check := range checks {
		var count int64

		query := repo.db.Table("third_party_assessments a").Where(check.Where, check.Args...)

		if params.DateFrom != "" && params.DateTo != "" {
			query = query.Where(
				"a.created_at BETWEEN ? AND ?",
				params.DateFrom+" 00:00:00",
				params.DateTo+" 23:59:59",
			)
		}

		query.Count(&count)

		results = append(results, &reportDTOs.TPOverdueAssessmentsDTO{
			Category: check.Category,
			Count:    count,
		})
	}

	return results, nil
}
