package repositories

import (
	"fmt"
	"math"
	"time"

	"gorm.io/gorm"
	auditDTOs "jethur/grc/audit/dtos"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/kpi/dtos"
	"jethur/grc/models"
)

type AuditDashboardRepositoryInterface interface {
	GetAuditStatusBreakdown(params *dtos.AuditDashboardQueryParamsDTO) (*dtos.AuditStatusBreakdownDTO, error)
	GetAuditSummary(params *dtos.AuditDashboardQueryParamsDTO) (*dtos.AuditSummaryDTO, error)
	GetAuditsByStatus(statusID uint, params *dtos.AuditDashboardListQueryParamsDTO) (*auditDTOs.AuditListResponseDTO, error)
	GetAuditFindingsSummary(params *dtos.AuditDashboardQueryParamsDTO) (*dtos.AuditFindingsSummaryDTO, error)
	GetAuditFindingsSeverityBreakdown(params *dtos.AuditDashboardQueryParamsDTO) (*dtos.AuditFindingsSeverityBreakdownDTO, error)
	GetAuditFindingsBySeverityLevel(severityLevelID uint, params *dtos.AuditDashboardListQueryParamsDTO) (*dtos.AuditFindingsBySeverityLevelListResponseDTO, error)
	GetAuditFindingsTrend() ([]dtos.AuditFindingsTrendMonthDTO, error)
	GetOverdueAuditActionPlans(params *dtos.AuditDashboardListQueryParamsDTO) (*dtos.AuditOverdueActionPlanListResponseDTO, error)

	// Helpers
	GetSeverityLevelTitle(severityLevelID uint) (string, error)
	GetAuditStatusTitle(statusID uint) (string, error)
}

type AuditDashboardRepository struct {
	db *gorm.DB
}

func NewAuditDashboardRepository(db *gorm.DB) AuditDashboardRepositoryInterface {
	return &AuditDashboardRepository{db: db}
}

func (repo *AuditDashboardRepository) applyDateFilter(query *gorm.DB, params *dtos.AuditDashboardQueryParamsDTO, table string) *gorm.DB {
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where(table+".created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	} else if params.Year > 0 {
		query = query.Where("YEAR("+table+".created_at) = ?", params.Year)
	}
	return query
}

func dateFilterParams(dateFrom, dateTo string, year int) *dtos.AuditDashboardQueryParamsDTO {
	return &dtos.AuditDashboardQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
		Year:     year,
	}
}

func (repo *AuditDashboardRepository) GetAuditStatusBreakdown(params *dtos.AuditDashboardQueryParamsDTO) (*dtos.AuditStatusBreakdownDTO, error) {
	summary := &dtos.AuditStatusBreakdownDTO{}

	excludedStatusTypes := []string{"draft", "in-review", "published", "in-progress"}

	totalQuery := repo.db.Model(&models.Audit{})
	totalQuery = repo.applyDateFilter(totalQuery, params, "audits")
	totalQuery = totalQuery.Where("audits.audit_status_id NOT IN (SELECT id FROM audit_statuses WHERE type IN ?)", excludedStatusTypes)
	if err := totalQuery.Count(&summary.Total).Error; err != nil {
		return nil, err
	}

	percentage := func(count int64) float64 {
		if summary.Total == 0 {
			return 0
		}
		return math.Round((float64(count)/float64(summary.Total))*100*100) / 100
	}

	// Breakdown by each real, currently-used audit status
	var auditStatuses []models.AuditStatus
	if err := repo.db.Where("type NOT IN ?", excludedStatusTypes).Find(&auditStatuses).Error; err != nil {
		return nil, err
	}

	// One grouped query for all per-status counts.
	type statusCountRow struct {
		AuditStatusID uint
		Cnt           int64
	}
	countQuery := repo.db.Model(&models.Audit{}).
		Select("audits.audit_status_id AS audit_status_id, COUNT(*) AS cnt").
		Group("audits.audit_status_id")
	countQuery = repo.applyDateFilter(countQuery, params, "audits")
	var countRows []statusCountRow
	if err := countQuery.Scan(&countRows).Error; err != nil {
		return nil, err
	}
	countByStatusID := make(map[uint]int64, len(countRows))
	for _, r := range countRows {
		countByStatusID[r.AuditStatusID] = r.Cnt
	}

	statusItems := make([]dtos.AuditStatusBreakdownItemDTO, 0, len(auditStatuses))
	for _, s := range auditStatuses {
		c := countByStatusID[uint(s.ID)]
		statusItems = append(statusItems, dtos.AuditStatusBreakdownItemDTO{
			Status: dtos.BreakdownReferenceDTO{
				ID:        uint(s.ID),
				Type:      s.Type,
				TitleEn:   s.TitleEn,
				TitleAr:   s.TitleAr,
				ColorCode: s.ColorCode,
			},
			Count:      c,
			Percentage: percentage(c),
		})
	}

	summary.Statuses = statusItems

	return summary, nil
}

func (repo *AuditDashboardRepository) GetAuditSummary(params *dtos.AuditDashboardQueryParamsDTO) (*dtos.AuditSummaryDTO, error) {
	summary := &dtos.AuditSummaryDTO{}

	totalQuery := repo.db.Model(&models.Audit{})
	totalQuery = repo.applyDateFilter(totalQuery, params, "audits")
	totalQuery = totalQuery.Where("audits.audit_status_id NOT IN (SELECT id FROM audit_statuses WHERE type IN ('draft', 'in-review', 'published'))")
	if err := totalQuery.Count(&summary.TotalAudits).Error; err != nil {
		return nil, err
	}

	completedQuery := repo.db.Model(&models.Audit{}).
		Joins("JOIN audit_statuses ON audit_statuses.id = audits.audit_status_id").
		Where("audit_statuses.type = ?", "completed")
	completedQuery = repo.applyDateFilter(completedQuery, params, "audits")
	if err := completedQuery.Count(&summary.CompletedAudits).Error; err != nil {
		return nil, err
	}

	// in-progress: any audit that has started but isn't completed yet.
	inProgressQuery := repo.db.Model(&models.Audit{}).
		Joins("JOIN audit_statuses ON audit_statuses.id = audits.audit_status_id").
		Where("audit_statuses.type IN ?", []string{"in-progress", "planning", "fieldwork", "reporting"})
	inProgressQuery = repo.applyDateFilter(inProgressQuery, params, "audits")
	if err := inProgressQuery.Count(&summary.InProgressAudits).Error; err != nil {
		return nil, err
	}

	// not-started: an audit plan has been created and published, but no audit has been started for it yet
	notStartedQuery := repo.db.Model(&models.AuditPlan{}).
		Joins("JOIN audit_statuses ON audit_statuses.id = audit_plans.audit_status_id").
		Where("audit_statuses.type = ?", "published")
	notStartedQuery = repo.applyDateFilter(notStartedQuery, params, "audit_plans")
	if err := notStartedQuery.Count(&summary.NotStartedAudits).Error; err != nil {
		return nil, err
	}

	return summary, nil
}

func (repo *AuditDashboardRepository) GetAuditsByStatus(statusID uint, params *dtos.AuditDashboardListQueryParamsDTO) (*auditDTOs.AuditListResponseDTO, error) {
	var audits []*auditDTOs.AuditListDTO
	var totalCount int64

	query := repo.db.Table("audits").
		Select(`
			audits.*,
			CONVERT(VARCHAR, audits.start_date, 23) AS start_date,
			CONVERT(VARCHAR, audits.end_date, 23) AS end_date,
			CONVERT(VARCHAR, audits.created_at, 23) AS created_at,
			audit_plans.title AS audit_plan_title,
			audit_plans.description AS audit_plan_description,
			audit_plans.purpose AS audit_plan_purpose,
			audit_plans.objective AS audit_plan_objective,
			audit_statuses.title_en AS audit_status_title_en,
			audit_statuses.title_ar AS audit_status_title_ar,
			audit_statuses.type AS audit_status_type,
			audit_statuses.color_code AS audit_status_color_code,
			audit_categories.title_en AS audit_category_title_en,
			audit_categories.title_ar AS audit_category_title_ar,
			team_lead.id AS team_lead_id,
			team_lead.first_name_en AS team_lead_first_name_en,
			team_lead.first_name_ar AS team_lead_first_name_ar,
			team_lead.last_name_en AS team_lead_last_name_en,
			team_lead.last_name_ar AS team_lead_last_name_ar,
			team_lead.email AS team_lead_email,
			team_lead.initial AS team_lead_initial,
			team_lead.color AS team_lead_color,
			team_lead.profile_image_path AS team_lead_profile_image_path,
			created_user.id AS created_user_id,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.email AS created_user_email,
			created_user.profile_image_path AS created_user_profile_image_path,
			updated_user.id AS updated_user_id,
			updated_user.first_name_en AS updated_user_first_name_en,
			updated_user.first_name_ar AS updated_user_first_name_ar,
			updated_user.last_name_en AS updated_user_last_name_en,
			updated_user.last_name_ar AS updated_user_last_name_ar,
			updated_user.initial AS updated_user_initial,
			updated_user.color AS updated_user_color,
			updated_user.email AS updated_user_email,
			updated_user.profile_image_path AS updated_user_profile_image_path
		`).
		Joins(`
			LEFT JOIN audit_statuses ON audit_statuses.id = audits.audit_status_id
			LEFT JOIN audit_plans ON audit_plans.id = audits.audit_plan_id
			LEFT JOIN audit_categories ON audit_categories.id = audit_plans.audit_category_id
			LEFT JOIN users AS team_lead ON team_lead.id = audit_plans.team_lead_id
			LEFT JOIN users AS created_user ON created_user.id = audits.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = audits.updated_by_id
		`)

	// Apply status filter.
	if statusID > 0 {
		query = query.Where("audits.audit_status_id = ?", statusID)
	} else {
		query = query.Where("audit_statuses.type NOT IN ?", []string{"draft", "in-review", "published"})
	}

	query = repo.applyDateFilter(query, dateFilterParams(params.DateFrom, params.DateTo, params.Year), "audits")

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(audits.start_date LIKE ? OR
			audits.end_date LIKE ? OR
			audit_plans.title LIKE ? OR
			audit_categories.title_en LIKE ? OR
			audit_categories.title_ar LIKE ? OR
			team_lead.first_name_en LIKE ? OR
			team_lead.first_name_ar LIKE ? OR
			team_lead.last_name_en LIKE ? OR
			team_lead.last_name_ar LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.last_name_ar LIKE ? OR
			updated_user.first_name_en LIKE ? OR
			updated_user.first_name_ar LIKE ? OR
			updated_user.last_name_en LIKE ? OR
			updated_user.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "audits.created_at"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "desc"
	}
	query = query.Order(params.SortBy + " " + params.Sort)

	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	// Execute the query
	if err := query.Find(&audits).Error; err != nil {
		return nil, err
	}

	totalPages := 1
	if !params.Export {
		totalPages = int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	}

	response := &auditDTOs.AuditListResponseDTO{
		Items: audits,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *AuditDashboardRepository) GetAuditFindingsSummary(params *dtos.AuditDashboardQueryParamsDTO) (*dtos.AuditFindingsSummaryDTO, error) {
	summary := &dtos.AuditFindingsSummaryDTO{}

	totalQuery := repo.db.Model(&models.AuditFinding{})
	totalQuery = repo.applyDateFilter(totalQuery, params, "audit_findings")
	if err := totalQuery.Count(&summary.FindingsTotal).Error; err != nil {
		return nil, err
	}

	openQuery := repo.db.Model(&models.AuditFinding{}).
		Joins("JOIN audit_finding_statuses ON audit_finding_statuses.id = audit_findings.audit_finding_status_id").
		Where("audit_finding_statuses.is_closed = ?", 0)
	openQuery = repo.applyDateFilter(openQuery, params, "audit_findings")
	if err := openQuery.Count(&summary.FindingsOpen).Error; err != nil {
		return nil, err
	}

	criticalQuery := repo.db.Model(&models.AuditFinding{}).
		Joins("JOIN severity_levels ON severity_levels.id = audit_findings.severity_level_id").
		Where("severity_levels.type IN ?", []string{"critical", "high"})
	criticalQuery = repo.applyDateFilter(criticalQuery, params, "audit_findings")
	if err := criticalQuery.Count(&summary.FindingsCriticalHigh).Error; err != nil {
		return nil, err
	}

	reassessmentQuery := repo.db.Model(&models.AuditFinding{}).
		Joins("JOIN audit_finding_statuses ON audit_finding_statuses.id = audit_findings.audit_finding_status_id").
		Where("audit_finding_statuses.type = ?", "re-assessment-pending")
	reassessmentQuery = repo.applyDateFilter(reassessmentQuery, params, "audit_findings")
	if err := reassessmentQuery.Count(&summary.ReassessmentPending).Error; err != nil {
		return nil, err
	}

	var closed int64
	closedQuery := repo.db.Model(&models.AuditFinding{}).
		Joins("JOIN audit_finding_statuses ON audit_finding_statuses.id = audit_findings.audit_finding_status_id").
		Where("audit_finding_statuses.is_closed = ?", 1)
	closedQuery = repo.applyDateFilter(closedQuery, params, "audit_findings")
	if err := closedQuery.Count(&closed).Error; err != nil {
		return nil, err
	}
	if summary.FindingsTotal > 0 {
		rate := float64(closed) / float64(summary.FindingsTotal) * 100
		summary.FindingClosureRate = math.Round(rate*100) / 100
	}

	return summary, nil
}

func (repo *AuditDashboardRepository) GetAuditFindingsSeverityBreakdown(params *dtos.AuditDashboardQueryParamsDTO) (*dtos.AuditFindingsSeverityBreakdownDTO, error) {
	result := &dtos.AuditFindingsSeverityBreakdownDTO{}

	// Breakdown by each severity level
	var severityLevels []models.SeverityLevel
	if err := repo.db.Order("weight").Find(&severityLevels).Error; err != nil {
		return nil, err
	}

	// One grouped query for all per-severity counts
	type severityCountRow struct {
		SeverityLevelID uint
		Cnt             int64
	}
	countQuery := repo.db.Model(&models.AuditFinding{}).
		Select("audit_findings.severity_level_id AS severity_level_id, COUNT(*) AS cnt").
		Group("audit_findings.severity_level_id")
	countQuery = repo.applyDateFilter(countQuery, params, "audit_findings")
	var countRows []severityCountRow
	if err := countQuery.Scan(&countRows).Error; err != nil {
		return nil, err
	}
	countBySeverityID := make(map[uint]int64, len(countRows))
	for _, r := range countRows {
		countBySeverityID[r.SeverityLevelID] = r.Cnt
		result.Total += r.Cnt
	}

	percentage := func(count int64) float64 {
		if result.Total == 0 {
			return 0
		}
		return math.Round((float64(count)/float64(result.Total))*100*100) / 100
	}

	result.SeverityLevels = make([]dtos.AuditFindingsSeverityBreakdownItemDTO, 0, len(severityLevels))
	for _, s := range severityLevels {
		c := countBySeverityID[s.ID]
		result.SeverityLevels = append(result.SeverityLevels, dtos.AuditFindingsSeverityBreakdownItemDTO{
			SeverityLevel: dtos.BreakdownReferenceDTO{
				ID:        s.ID,
				Type:      s.Type,
				TitleEn:   s.TitleEn,
				TitleAr:   s.TitleAr,
				ColorCode: s.ColorCode,
			},
			Count:      c,
			Percentage: percentage(c),
		})
	}
	return result, nil
}

func (repo *AuditDashboardRepository) GetAuditFindingsBySeverityLevel(severityLevelID uint, params *dtos.AuditDashboardListQueryParamsDTO) (*dtos.AuditFindingsBySeverityLevelListResponseDTO, error) {
	var items []*dtos.AuditFindingsBySeverityLevelListDTO
	var totalCount int64

	query := repo.db.Table("audit_findings").
		Select(`
			audit_findings.id,
			audit_findings.audit_id,
			audit_findings.title,
			audit_findings.reference_number,
			audit_findings.auditable_item_id,
			audit_findings.audit_finding_category_id,
			audit_findings.severity_level_id,
			audit_findings.audit_finding_status_id,
			CONVERT(VARCHAR, audit_findings.created_at, 23) AS created_at,
			CONVERT(VARCHAR, audit_findings.updated_at, 23) AS updated_at,
			audit_findings.created_by_id,
			audit_findings.updated_by_id,
			audit_plans.title AS audit_title,
			auditable_items.title_en AS auditable_item_title_en,
			auditable_items.title_ar AS auditable_item_title_ar,
			audit_finding_categories.title_en AS audit_finding_category_title_en,
			audit_finding_categories.title_ar AS audit_finding_category_title_ar,
			severity_levels.title_en AS severity_level_title_en,
			severity_levels.title_ar AS severity_level_title_ar,
			severity_levels.color_code AS severity_level_color_code,
			audit_finding_statuses.title_en AS audit_finding_status_title_en,
			audit_finding_statuses.title_ar AS audit_finding_status_title_ar,
			audit_finding_statuses.color_code AS audit_finding_status_color_code,
			created_user.id AS created_user_id,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.email AS created_user_email,
			created_user.profile_image_path AS created_user_profile_image_path,
			updated_user.id AS updated_user_id,
			updated_user.first_name_en AS updated_user_first_name_en,
			updated_user.first_name_ar AS updated_user_first_name_ar,
			updated_user.last_name_en AS updated_user_last_name_en,
			updated_user.last_name_ar AS updated_user_last_name_ar,
			updated_user.initial AS updated_user_initial,
			updated_user.color AS updated_user_color,
			updated_user.email AS updated_user_email,
			updated_user.profile_image_path AS updated_user_profile_image_path
		`).
		Joins(`
			LEFT JOIN audits ON audits.id = audit_findings.audit_id
			LEFT JOIN audit_plans ON audit_plans.id = audits.audit_plan_id
			LEFT JOIN auditable_items ON auditable_items.id = audit_findings.auditable_item_id
			LEFT JOIN audit_finding_categories ON audit_finding_categories.id = audit_findings.audit_finding_category_id
			LEFT JOIN severity_levels ON severity_levels.id = audit_findings.severity_level_id
			LEFT JOIN audit_finding_statuses ON audit_finding_statuses.id = audit_findings.audit_finding_status_id
			LEFT JOIN users AS created_user ON created_user.id = audit_findings.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = audit_findings.updated_by_id
		`)

	if severityLevelID > 0 {
		query = query.Where("audit_findings.severity_level_id = ?", severityLevelID)
	}

	query = repo.applyDateFilter(query, dateFilterParams(params.DateFrom, params.DateTo, params.Year), "audit_findings")

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(auditable_items.title_en LIKE ? OR
			audit_findings.title LIKE ? OR
			audit_findings.reference_number LIKE ? OR
			audit_finding_categories.title_en LIKE ? OR
			severity_levels.title_en LIKE ? OR
			audit_finding_statuses.title_en LIKE ? OR
			audit_finding_statuses.title_ar LIKE ? OR
			created_user.first_name_en LIKE ? OR
			updated_user.first_name_en LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "audit_findings.created_at"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}
	query = query.Order(params.SortBy + " " + params.Sort)

	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	// Execute the query
	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}

	totalPages := 1
	if !params.Export {
		totalPages = int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	}

	response := &dtos.AuditFindingsBySeverityLevelListResponseDTO{
		Items: items,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *AuditDashboardRepository) GetAuditFindingsTrend() ([]dtos.AuditFindingsTrendMonthDTO, error) {
	now := time.Now()
	months := make([]time.Time, 0, 6)
	for i := 5; i >= 0; i-- {
		months = append(months, time.Date(now.Year(), now.Month()-time.Month(i), 1, 0, 0, 0, 0, now.Location()))
	}
	start := months[0]

	englishMonths := []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
	arabicMonths := []string{"يناير", "فبراير", "مارس", "أبريل", "مايو", "يونيو", "يوليو", "أغسطس", "سبتمبر", "أكتوبر", "نوفمبر", "ديسمبر"}

	type monthCount struct {
		Ym  string
		Cnt int64
	}

	var newRows []monthCount
	if err := repo.db.Model(&models.AuditFinding{}).
		Select("FORMAT(created_at, 'yyyy-MM') AS ym, COUNT(*) AS cnt").
		Where("created_at >= ?", start).
		Group("FORMAT(created_at, 'yyyy-MM')").
		Scan(&newRows).Error; err != nil {
		return nil, err
	}

	var closedRows []monthCount
	if err := repo.db.Model(&models.AuditFinding{}).
		Joins("JOIN audit_finding_statuses ON audit_finding_statuses.id = audit_findings.audit_finding_status_id").
		Select("FORMAT(audit_findings.updated_at, 'yyyy-MM') AS ym, COUNT(*) AS cnt").
		Where("audit_finding_statuses.is_closed = ?", 1).
		Where("audit_findings.updated_at >= ?", start).
		Group("FORMAT(audit_findings.updated_at, 'yyyy-MM')").
		Scan(&closedRows).Error; err != nil {
		return nil, err
	}

	newMap := make(map[string]int64)
	for _, r := range newRows {
		newMap[r.Ym] = r.Cnt
	}
	closedMap := make(map[string]int64)
	for _, r := range closedRows {
		closedMap[r.Ym] = r.Cnt
	}

	result := make([]dtos.AuditFindingsTrendMonthDTO, 0, len(months))
	for _, m := range months {
		key := fmt.Sprintf("%04d-%02d", m.Year(), int(m.Month()))
		result = append(result, dtos.AuditFindingsTrendMonthDTO{
			Month:   key,
			MonthEn: fmt.Sprintf("%s %d", englishMonths[m.Month()-1], m.Year()),
			MonthAr: fmt.Sprintf("%s %d", arabicMonths[m.Month()-1], m.Year()),
			New:     newMap[key],
			Closed:  closedMap[key],
		})
	}
	return result, nil
}

func (repo *AuditDashboardRepository) GetOverdueAuditActionPlans(params *dtos.AuditDashboardListQueryParamsDTO) (*dtos.AuditOverdueActionPlanListResponseDTO, error) {
	var items []*dtos.AuditOverdueActionPlanListDTO
	var totalCount int64

	departmentSubqueryEn := `ISNULL((
		SELECT TOP 1 ou.title_en FROM organization_units ou
		WHERE ou.id IN (
			responsible_user.organization_unit_1_id,
			responsible_user.organization_unit_2_id,
			responsible_user.organization_unit_3_id,
			responsible_user.organization_unit_4_id,
			responsible_user.organization_unit_5_id
		)
		AND EXISTS (SELECT 1 FROM organization_levels ol WHERE ol.id = ou.organization_level_id AND ol.is_primary = 1)
	), '')`
	departmentSubqueryAr := `ISNULL((
		SELECT TOP 1 ou.title_ar FROM organization_units ou
		WHERE ou.id IN (
			responsible_user.organization_unit_1_id,
			responsible_user.organization_unit_2_id,
			responsible_user.organization_unit_3_id,
			responsible_user.organization_unit_4_id,
			responsible_user.organization_unit_5_id
		)
		AND EXISTS (SELECT 1 FROM organization_levels ol WHERE ol.id = ou.organization_level_id AND ol.is_primary = 1)
	), '')`

	query := repo.db.Table("audit_finding_action_plans").
		Select(`
        audit_finding_action_plans.id,
		audit_finding_action_plans.audit_finding_id,
		audit_finding_action_plans.action_plan_status_id,
		audit_finding_action_plans.severity_level_id,
		audit_finding_action_plans.responsible_user_id,
		audit_finding_action_plans.review_user_id,
		audit_finding_action_plans.created_by_id,
		audit_finding_action_plans.updated_by_id,
		audit_finding_action_plans.title,
		audit_finding_action_plans.target_date,
		audit_finding_action_plans.updated_at,
		audit_findings.title AS audit_finding_title,
		audit_findings.reference_number AS audit_finding_reference_number,
		audit_findings.description AS audit_finding_description,
		CONVERT(VARCHAR, audit_finding_action_plans.target_date, 23) AS target_date,
		CONVERT(VARCHAR, audit_finding_action_plans.created_at, 23) AS created_at,
        action_plan_statuses.title_en AS action_plan_status_title_en,
        action_plan_statuses.title_ar AS action_plan_status_title_ar,
		action_plan_statuses.color_code AS action_plan_status_color_code,
        severity_levels.title_en AS severity_level_title_en,
        severity_levels.title_ar AS severity_level_title_ar,
        severity_levels.color_code AS severity_level_color_code,
		responsible_user.id AS responsible_user_id,
        responsible_user.first_name_en AS responsible_user_first_name_en,
        responsible_user.first_name_ar AS responsible_user_first_name_ar,
        responsible_user.last_name_en AS responsible_user_last_name_en,
        responsible_user.last_name_ar AS responsible_user_last_name_ar,
		responsible_user.initial AS responsible_user_initial,
        responsible_user.color AS responsible_user_color,
        responsible_user.email AS responsible_user_email,
        responsible_user.profile_image_path AS responsible_user_profile_image_path,
		` + departmentSubqueryEn + ` AS responsible_user_department_title_en,
		` + departmentSubqueryAr + ` AS responsible_user_department_title_ar,
		CASE WHEN audit_finding_action_plans.target_date < GETDATE() THEN DATEDIFF(day, audit_finding_action_plans.target_date, GETDATE()) ELSE 0 END AS overdue_days,
		created_user.first_name_en AS created_user_first_name_en,
		created_user.first_name_ar AS created_user_first_name_ar,
		created_user.last_name_en AS created_user_last_name_en,
		created_user.last_name_ar AS created_user_last_name_ar,
		created_user.initial AS created_user_initial,
		created_user.color AS created_user_color,
		created_user.email AS created_user_email,
		created_user.profile_image_path AS created_user_profile_image_path
    `).
		Joins(`
        LEFT JOIN audit_findings ON audit_findings.id = audit_finding_action_plans.audit_finding_id
        LEFT JOIN action_plan_statuses ON action_plan_statuses.id = audit_finding_action_plans.action_plan_status_id
        LEFT JOIN severity_levels ON severity_levels.id = audit_finding_action_plans.severity_level_id
        LEFT JOIN users AS responsible_user ON responsible_user.id = audit_finding_action_plans.responsible_user_id
        LEFT JOIN users AS created_user ON created_user.id = audit_finding_action_plans.created_by_id
    `)

	// Base filter: only overdue action plans (not closed and past target date)
	query = query.Where("action_plan_statuses.is_closed = ?", 0)
	query = query.Where("audit_finding_action_plans.target_date IS NOT NULL")
	query = query.Where("audit_finding_action_plans.target_date < ?", time.Now())

	query = repo.applyDateFilter(query, dateFilterParams(params.DateFrom, params.DateTo, params.Year), "audit_finding_action_plans")

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
            (audit_finding_action_plans.title LIKE ? OR
            action_plan_statuses.title_en LIKE ? OR
            action_plan_statuses.title_ar LIKE ? OR
            audit_finding_action_plans.target_date LIKE ? OR
            responsible_user.first_name_en LIKE ? OR
            responsible_user.first_name_ar LIKE ? OR
            responsible_user.last_name_en LIKE ? OR
            responsible_user.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword,
		)
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "audit_finding_action_plans.title"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	// Execute the query
	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}

	totalPages := 1
	if !params.Export {
		totalPages = int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	}

	response := &dtos.AuditOverdueActionPlanListResponseDTO{
		Items: items,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *AuditDashboardRepository) GetSeverityLevelTitle(severityLevelID uint) (string, error) {
	var title string
	if err := repo.db.
		Table("severity_levels").
		Select("title_en").
		Where("id = ?", severityLevelID).
		Scan(&title).Error; err != nil {
		return "", err
	}
	return title, nil
}

func (repo *AuditDashboardRepository) GetAuditStatusTitle(statusID uint) (string, error) {
	var title string
	if err := repo.db.
		Table("audit_statuses").
		Select("title_en").
		Where("id = ?", statusID).
		Scan(&title).Error; err != nil {
		return "", err
	}
	return title, nil
}
