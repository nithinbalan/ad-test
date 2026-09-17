package repositories

import (
	incidentDTOs "jethur/grc/kpi/dtos"
	"jethur/grc/models"
	"math"
	"time"

	"gorm.io/gorm"
)

type IncidentDashboardRepositoryInterface interface {
	GetIncidentStatusSummary(params *incidentDTOs.IncidentDashboardQueryParamsDTO) (*incidentDTOs.IncidentStatusSummaryDTO, error)
	GetIncidentSeverityBreakdown(params *incidentDTOs.IncidentDashboardQueryParamsDTO) (*incidentDTOs.IncidentSeverityBreakdownDTO, error)
	GetIncidentAcknowledgementSummary(params *incidentDTOs.IncidentDashboardQueryParamsDTO) (*incidentDTOs.IncidentAcknowledgementSummaryDTO, error)
	GetIncidentSLASummary(params *incidentDTOs.IncidentDashboardQueryParamsDTO) (*incidentDTOs.IncidentSLASummaryDTO, error)
	GetIncidentCategorySummary(params *incidentDTOs.IncidentDashboardQueryParamsDTO) ([]incidentDTOs.IncidentCategoryCountDTO, error)
	GetIncidentTypeSummary(params *incidentDTOs.IncidentDashboardQueryParamsDTO) ([]incidentDTOs.IncidentTypeCountDTO, error)
	GetIncidentTrend(params *incidentDTOs.IncidentDashboardQueryParamsDTO) ([]incidentDTOs.IncidentTrendDTO, error)
	GetIncidentActionPlanSummary(params *incidentDTOs.IncidentDashboardQueryParamsDTO) (*incidentDTOs.IncidentActionPlanSummaryDTO, error)
	GetOverdueActionPlans(params *incidentDTOs.IncidentDashboardQueryParamsDTO) ([]incidentDTOs.OverdueActionPlanDTO, error)
}

type IncidentDashboardRepository struct {
	db *gorm.DB
}

func NewIncidentDashboardRepository(db *gorm.DB) IncidentDashboardRepositoryInterface {
	return &IncidentDashboardRepository{db: db}
}

type incidentTrendCount struct {
	Month time.Time
	Count int64
}

func (repo *IncidentDashboardRepository) applyDateFilter(query *gorm.DB, params *incidentDTOs.IncidentDashboardQueryParamsDTO, table string) *gorm.DB {
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where(table+".reported_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	} else if params.Year > 0 {
		query = query.Where("YEAR("+table+".reported_at) = ?", params.Year)
	}
	return query
}

func (repo *IncidentDashboardRepository) applyTargetDateFilter(query *gorm.DB, params *incidentDTOs.IncidentDashboardQueryParamsDTO, table string) *gorm.DB {
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where(table+".target_date BETWEEN ? AND ?", params.DateFrom, params.DateTo)
	} else if params.Year > 0 {
		query = query.Where("YEAR("+table+".target_date) = ?", params.Year)
	}
	return query
}

func (repo *IncidentDashboardRepository) GetIncidentStatusSummary(params *incidentDTOs.IncidentDashboardQueryParamsDTO) (*incidentDTOs.IncidentStatusSummaryDTO, error) {
	summary := &incidentDTOs.IncidentStatusSummaryDTO{
		Statuses: []incidentDTOs.IncidentStatusCountDTO{},
	}

	totalQuery := repo.db.Model(&models.IncidentRegister{})
	totalQuery = repo.applyDateFilter(totalQuery, params, "incident_register")
	if err := totalQuery.Count(&summary.TotalIncidents).Error; err != nil {
		return nil, err
	}

	openQuery := repo.db.Table("incident_register ir").
		Joins("INNER JOIN incident_statuses ist ON ist.id = ir.incident_status_id").
		Where("ist.type <> ?", "closed")
	openQuery = repo.applyDateFilter(openQuery, params, "ir")
	if err := openQuery.Count(&summary.OpenIncidents).Error; err != nil {
		return nil, err
	}

	closedQuery := repo.db.Table("incident_register ir").
		Joins("INNER JOIN incident_statuses ist ON ist.id = ir.incident_status_id").
		Where("ist.type = ?", "closed")
	closedQuery = repo.applyDateFilter(closedQuery, params, "ir")
	if err := closedQuery.Count(&summary.ClosedIncidents).Error; err != nil {
		return nil, err
	}

	statusQuery := repo.db.Table("incident_statuses ist").
		Select(`
			ist.id,
			ist.title_en,
			ist.title_ar,
			ist.color_code,
			ist.type,
			COUNT(ir.id) AS count
		`).
		Joins("LEFT JOIN incident_register ir ON ir.incident_status_id = ist.id")
	statusQuery = repo.applyDateFilter(statusQuery, params, "ir")
	if err := statusQuery.Group("ist.id, ist.title_en, ist.title_ar, ist.color_code, ist.type").
		Order("ist.title_en ASC").
		Scan(&summary.Statuses).Error; err != nil {
		return nil, err
	}

	return summary, nil
}

func (repo *IncidentDashboardRepository) GetIncidentSeverityBreakdown(params *incidentDTOs.IncidentDashboardQueryParamsDTO) (*incidentDTOs.IncidentSeverityBreakdownDTO, error) {
	result := &incidentDTOs.IncidentSeverityBreakdownDTO{
		SeverityLevels: []incidentDTOs.IncidentSeverityCountDTO{},
	}

	totalQuery := repo.db.Model(&models.IncidentRegister{})
	totalQuery = repo.applyDateFilter(totalQuery, params, "incident_register")
	if err := totalQuery.Count(&result.TotalIncidents).Error; err != nil {
		return nil, err
	}

	severityQuery := repo.db.Table("severity_levels sl").
		Select(`
			sl.id,
			sl.title_en,
			sl.title_ar,
			sl.color_code,
			sl.type,
			COUNT(ir.id) AS count
		`).
		Joins("LEFT JOIN incident_register ir ON ir.severity_level_id = sl.id")
	severityQuery = repo.applyDateFilter(severityQuery, params, "ir")
	if err := severityQuery.Group("sl.id, sl.title_en, sl.title_ar, sl.color_code, sl.type, sl.weight").
		Order("sl.weight DESC, sl.id").
		Scan(&result.SeverityLevels).Error; err != nil {
		return nil, err
	}

	for i := range result.SeverityLevels {
		if result.TotalIncidents > 0 {
			result.SeverityLevels[i].Percentage = float64(result.SeverityLevels[i].Count) * 100 / float64(result.TotalIncidents)
		}
	}

	return result, nil
}

func (repo *IncidentDashboardRepository) GetIncidentAcknowledgementSummary(params *incidentDTOs.IncidentDashboardQueryParamsDTO) (*incidentDTOs.IncidentAcknowledgementSummaryDTO, error) {
	summary := &incidentDTOs.IncidentAcknowledgementSummaryDTO{}

	totalQuery := repo.db.Model(&models.IncidentRegister{})
	totalQuery = repo.applyDateFilter(totalQuery, params, "incident_register")
	if err := totalQuery.Count(&summary.Total).Error; err != nil {
		return nil, err
	}

	ackQuery := repo.db.Model(&models.IncidentRegister{}).Where("is_acknowledge = ?", true)
	ackQuery = repo.applyDateFilter(ackQuery, params, "incident_register")
	if err := ackQuery.Count(&summary.Acknowledged).Error; err != nil {
		return nil, err
	}

	summary.NotAcknowledged = summary.Total - summary.Acknowledged

	return summary, nil
}

func (repo *IncidentDashboardRepository) GetIncidentSLASummary(params *incidentDTOs.IncidentDashboardQueryParamsDTO) (*incidentDTOs.IncidentSLASummaryDTO, error) {
	summary := &incidentDTOs.IncidentSLASummaryDTO{}

	achievedQuery := repo.db.Model(&models.IncidentRegister{}).Where("is_sla_time_achieved = ?", true)
	achievedQuery = repo.applyDateFilter(achievedQuery, params, "incident_register")
	if err := achievedQuery.Count(&summary.SLAAchieved).Error; err != nil {
		return nil, err
	}

	exceededQuery := repo.db.Model(&models.IncidentRegister{}).Where("is_sla_time_exceeded = ?", true)
	exceededQuery = repo.applyDateFilter(exceededQuery, params, "incident_register")
	if err := exceededQuery.Count(&summary.SLAExceeded).Error; err != nil {
		return nil, err
	}

	total := summary.SLAAchieved + summary.SLAExceeded
	if total > 0 {
		summary.SLACompliancePct = math.Round(float64(summary.SLAAchieved) / float64(total) * 100)
	}

	return summary, nil
}

func (repo *IncidentDashboardRepository) GetIncidentCategorySummary(params *incidentDTOs.IncidentDashboardQueryParamsDTO) ([]incidentDTOs.IncidentCategoryCountDTO, error) {
	result := []incidentDTOs.IncidentCategoryCountDTO{}

	query := repo.db.Table("incident_categories ic").
		Select(`
			ic.id AS category_id,
			ic.title_en,
			ic.title_ar,
			COUNT(ir.id) AS count
		`).
		Joins("LEFT JOIN incident_register ir ON ir.incident_category_id = ic.id")
	query = repo.applyDateFilter(query, params, "ir")
	err := query.Group("ic.id, ic.title_en, ic.title_ar").
		Order("count DESC").
		Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *IncidentDashboardRepository) GetIncidentTypeSummary(params *incidentDTOs.IncidentDashboardQueryParamsDTO) ([]incidentDTOs.IncidentTypeCountDTO, error) {
	result := []incidentDTOs.IncidentTypeCountDTO{}

	query := repo.db.Table("incident_types it").
		Select(`
			it.id AS type_id,
			it.title_en,
			it.title_ar,
			COUNT(ir.id) AS count
		`).
		Joins("LEFT JOIN incident_register ir ON ir.incident_type_id = it.id")
	query = repo.applyDateFilter(query, params, "ir")
	err := query.Group("it.id, it.title_en, it.title_ar").
		Order("count DESC").
		Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *IncidentDashboardRepository) GetIncidentTrend(params *incidentDTOs.IncidentDashboardQueryParamsDTO) ([]incidentDTOs.IncidentTrendDTO, error) {
	now := time.Now().UTC()
	year := now.Year()
	if params.Year > 0 {
		year = params.Year
	}
	startOfYear := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	nextMonth := startOfYear.AddDate(1, 0, 0)

	openedQuery := repo.db.Table("incident_register ir").
		Joins("INNER JOIN incident_statuses ist ON ist.id = ir.incident_status_id").
		Select("DATEFROMPARTS(YEAR(ir.reported_at), MONTH(ir.reported_at), 1) AS month, COUNT(ir.id) AS count").
		Where("ist.type != ? AND ir.reported_at >= ? AND ir.reported_at < ?", "closed", startOfYear, nextMonth)
	openedQuery = repo.applyDateFilter(openedQuery, params, "ir")
	openedCounts := []incidentTrendCount{}
	if err := openedQuery.Group("YEAR(ir.reported_at), MONTH(ir.reported_at)").
		Scan(&openedCounts).Error; err != nil {
		return nil, err
	}

	closedQuery := repo.db.Table("incident_register ir").
		Joins("INNER JOIN incident_statuses ist ON ist.id = ir.incident_status_id").
		Select("DATEFROMPARTS(YEAR(ir.reported_at), MONTH(ir.reported_at), 1) AS month, COUNT(ir.id) AS count").
		Where("ist.type = ? AND ir.reported_at >= ? AND ir.reported_at < ?", "closed", startOfYear, nextMonth)
	closedQuery = repo.applyDateFilter(closedQuery, params, "ir")
	closedCounts := []incidentTrendCount{}
	if err := closedQuery.Group("YEAR(ir.reported_at), MONTH(ir.reported_at)").
		Scan(&closedCounts).Error; err != nil {
		return nil, err
	}

	openedByMonth := make(map[string]int64, len(openedCounts))
	for _, item := range openedCounts {
		openedByMonth[item.Month.Format("2006-01")] = item.Count
	}
	closedByMonth := make(map[string]int64, len(closedCounts))
	for _, item := range closedCounts {
		closedByMonth[item.Month.Format("2006-01")] = item.Count
	}

	trend := make([]incidentDTOs.IncidentTrendDTO, 0, 12)
	for month := startOfYear; month.Before(nextMonth); month = month.AddDate(0, 1, 0) {
		key := month.Format("2006-01")
		trend = append(trend, incidentDTOs.IncidentTrendDTO{
			Month:  month.Format("Jan 06"),
			Opened: openedByMonth[key],
			Closed: closedByMonth[key],
		})
	}

	return trend, nil
}

func (repo *IncidentDashboardRepository) GetIncidentActionPlanSummary(params *incidentDTOs.IncidentDashboardQueryParamsDTO) (*incidentDTOs.IncidentActionPlanSummaryDTO, error) {
	summary := &incidentDTOs.IncidentActionPlanSummaryDTO{}

	totalQuery := repo.db.Model(&models.IncidentActionPlan{})
	totalQuery = repo.applyTargetDateFilter(totalQuery, params, "incident_action_plans")
	if err := totalQuery.Count(&summary.Total).Error; err != nil {
		return nil, err
	}

	completedQuery := repo.db.Table("incident_action_plans iap").
		Joins("INNER JOIN action_plan_statuses aps ON aps.id = iap.action_plan_status_id").
		Where("aps.is_closed = ?", true)
	completedQuery = repo.applyTargetDateFilter(completedQuery, params, "iap")
	if err := completedQuery.Count(&summary.Completed).Error; err != nil {
		return nil, err
	}

	summary.Open = summary.Total - summary.Completed

	overdueQuery := repo.db.Table("incident_action_plans iap").
		Joins("INNER JOIN action_plan_statuses aps ON aps.id = iap.action_plan_status_id").
		Where("aps.is_closed = ? AND iap.target_date IS NOT NULL AND iap.target_date < CAST(GETDATE() AS date)", false)
	overdueQuery = repo.applyTargetDateFilter(overdueQuery, params, "iap")
	if err := overdueQuery.Count(&summary.Overdue).Error; err != nil {
		return nil, err
	}

	return summary, nil
}

func (repo *IncidentDashboardRepository) GetOverdueActionPlans(params *incidentDTOs.IncidentDashboardQueryParamsDTO) ([]incidentDTOs.OverdueActionPlanDTO, error) {
	result := []incidentDTOs.OverdueActionPlanDTO{}

	query := repo.db.Table("incident_action_plans iap").
		Select(`
			iap.id AS action_plan_id,
			iap.title,
			ir.id AS incident_register_id,
			ir.reference_number AS incident_reference_number,
			ir.title AS incident_title,
			COALESCE(sl.title_en, '') AS severity_level_title_en,
			COALESCE(sl.title_ar, '') AS severity_level_title_ar,
			COALESCE(sl.color_code, '') AS severity_level_color_code,
			COALESCE(sl.type, '') AS severity_level_type,
			u.id AS responsible_user_id,
			u.first_name_en AS responsible_user_first_name_en,
			u.first_name_ar AS responsible_user_first_name_ar,
			u.last_name_en AS responsible_user_last_name_en,
			u.last_name_ar AS responsible_user_last_name_ar,
			u.initial AS responsible_user_initial,
			u.color AS responsible_user_color,
			u.email AS responsible_user_email,
			u.profile_image_path AS responsible_user_profile_image_path,
			CONVERT(VARCHAR, iap.target_date, 23) AS target_date,
			DATEDIFF(day, iap.target_date, CAST(GETDATE() AS date)) AS overdue_days
		`).
		Joins("INNER JOIN incident_register ir ON ir.id = iap.incident_register_id").
		Joins("INNER JOIN action_plan_statuses aps ON aps.id = iap.action_plan_status_id").
		Joins("LEFT JOIN severity_levels sl ON sl.id = iap.severity_level_id").
		Joins("INNER JOIN users u ON u.id = iap.responsible_user_id").
		Where("aps.is_closed = ? AND iap.target_date IS NOT NULL AND iap.target_date < CAST(GETDATE() AS date)", false)
	query = repo.applyTargetDateFilter(query, params, "iap")
	err := query.Order("iap.target_date ASC").
		Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
