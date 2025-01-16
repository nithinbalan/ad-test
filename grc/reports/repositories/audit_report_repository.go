package repositories

import (
	"fmt"
	reportDTOs "jethur/grc/reports/dtos"

	"gorm.io/gorm"
)

type AuditReportRepositoryInterface interface {
	GetAuditPlanByCategory(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditPlanByCategoryDTO, error)
	GetAuditPlanByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditPlanByStatusDTO, error)
	GetAuditPlanByTeam(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditPlanByTeamDTO, error)
	GetAuditPlanByTeamLead(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditPlanByTeamLeadDTO, error)
	GetAuditPlanByYear(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditPlanByYearDTO, error)

	GetAuditByCategory(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditByCategoryDTO, error)
	GetAuditByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditByStatusDTO, error)
	GetAuditByOpenVsClosed(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditByOpenVsClosedDTO, error)
	GetAuditByMonth(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditByMonthDTO, error)
	GetAuditByTeamLead(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditByTeamLeadDTO, error)

	GetAuditFindingByCategory(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditFindingByCategoryDTO, error)
	GetAuditFindingByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditFindingByStatusDTO, error)
	GetAuditFindingBySeverityLevel(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditFindingBySeverityLevelDTO, error)
	GetAuditFindingByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditFindingByResponsibleUnitDTO, error)
	GetAuditFindingByMonth(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditFindingByMonthDTO, error)
	GetAuditFindingByAuditor(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditFindingByAuditorDTO, error)

	GetAuditActionPlanByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditActionPlanByStatusDTO, error)
	GetAuditActionPlanBySeverityLevel(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditActionPlanBySeverityLevelDTO, error)
	GetAuditActionPlanByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditActionPlanByResponsibleUnitDTO, error)
	GetAuditActionPlanByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditActionPlanByResponsibleUserDTO, error)
	GetAuditActionPlanOnTimeVsDelayed(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditActionPlanOnTimeVsDelayedDTO, error)
}

type AuditReportRepository struct {
	db *gorm.DB
}

func NewAuditReportRepository(db *gorm.DB) *AuditReportRepository {
	return &AuditReportRepository{db: db}
}

// Implementation of the `GetAuditPlanByCategory` method
func (r *AuditReportRepository) GetAuditPlanByCategory(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditPlanByCategoryDTO, error) {
	results := []*reportDTOs.AuditPlanByCategoryDTO{}

	query := r.db.Table("audit_plans").
		Select(`audit_categories.id AS audit_category_id, 
				audit_categories.title_en AS audit_category_title_en, 
				audit_categories.title_ar AS audit_category_title_ar, 
				COUNT(audit_plans.id) AS count`).
		Joins("JOIN audit_categories ON audit_categories.id = audit_plans.audit_category_id").
		Group("audit_categories.id, audit_categories.title_en, audit_categories.title_ar")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("audit_plans.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (r *AuditReportRepository) GetAuditPlanByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditPlanByStatusDTO, error) {
	results := []*reportDTOs.AuditPlanByStatusDTO{}

	query := r.db.Table("audit_plans").
		Select(`audit_statuses.id AS audit_status_id,
				audit_statuses.title_en AS audit_status_title_en,
				audit_statuses.title_ar AS audit_status_title_ar,
				audit_statuses.color_code AS color_code,
				COUNT(audit_plans.id) AS count`).
		Joins("JOIN audit_statuses ON audit_statuses.id = audit_plans.audit_status_id").
		Group("audit_statuses.id, audit_statuses.title_en, audit_statuses.title_ar, audit_statuses.color_code")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("audit_plans.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (r *AuditReportRepository) GetAuditPlanByTeam(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditPlanByTeamDTO, error) {
	results := []*reportDTOs.AuditPlanByTeamDTO{}

	query := r.db.Table("audit_plans").
		Select(`audit_teams.id AS audit_team_id,
				audit_teams.title_en AS audit_team_title_en,
				audit_teams.title_ar AS audit_team_title_ar,
				COUNT(audit_plans.id) AS count`).
		Joins("JOIN audit_teams ON audit_teams.id = audit_plans.audit_team_id").
		Group("audit_teams.id, audit_teams.title_en, audit_teams.title_ar")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("audit_plans.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (r *AuditReportRepository) GetAuditPlanByTeamLead(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditPlanByTeamLeadDTO, error) {
	results := []*reportDTOs.AuditPlanByTeamLeadDTO{}

	query := r.db.Table("audit_plans").
		Select(`users.id AS team_lead_id,
			users.first_name_en AS team_lead_first_name_en,
			users.first_name_ar AS team_lead_first_name_ar,
			users.last_name_en AS team_lead_last_name_en,
			users.last_name_ar AS team_lead_last_name_ar,
			CONCAT(users.first_name_en, ' ', users.last_name_en) AS team_lead_title_en,
			CONCAT(users.first_name_ar, ' ', users.last_name_ar) AS team_lead_title_ar,
			users.initial AS team_lead_initial,
			users.color AS team_lead_color,
			users.email AS team_lead_email,
			users.profile_image_path AS team_lead_profile_image_path,
			COUNT(audit_plans.id) AS count`).
		Joins("JOIN users ON users.id = audit_plans.team_lead_id").
		Group("users.id, users.first_name_en, users.first_name_ar, users.last_name_en, users.last_name_ar, users.initial, users.color, users.email, users.profile_image_path")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("audit_plans.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (r *AuditReportRepository) GetAuditPlanByYear(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditPlanByYearDTO, error) {
	results := []*reportDTOs.AuditPlanByYearDTO{}

	// Fetch yearly audit plan count
	var rawResults []struct {
		Year  int  `gorm:"column:year"`
		Count uint `gorm:"column:count"`
	}

	query := r.db.Table("audit_plans").
		Select(`YEAR(audit_plans.created_at) AS year, COUNT(audit_plans.id) AS count`).
		Group("YEAR(audit_plans.created_at)").
		Order("YEAR(audit_plans.created_at) ASC")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("audit_plans.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&rawResults).Error
	if err != nil {
		return nil, err
	}

	// Convert year to English and Arabic title
	for _, record := range rawResults {
		results = append(results, &reportDTOs.AuditPlanByYearDTO{
			TitleEn: fmt.Sprintf("%d", record.Year), // English Title (Year)
			TitleAr: fmt.Sprintf("%d", record.Year), // Arabic Title (Same as Year)
			Count:   record.Count,
		})
	}

	return results, nil
}

func (r *AuditReportRepository) GetAuditByCategory(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditByCategoryDTO, error) {
	results := []*reportDTOs.AuditByCategoryDTO{}

	query := r.db.Table("audits").
		Select(`audit_categories.id AS audit_category_id, 
				audit_categories.title_en AS audit_category_title_en, 
				audit_categories.title_ar AS audit_category_title_ar, 
				COUNT(audits.id) AS count`).
		Joins("JOIN audit_plans ON audit_plans.id = audits.audit_plan_id").
		Joins("JOIN audit_categories ON audit_categories.id = audit_plans.audit_category_id").
		Group("audit_categories.id, audit_categories.title_en, audit_categories.title_ar")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("audits.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (r *AuditReportRepository) GetAuditByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditByStatusDTO, error) {
	results := []*reportDTOs.AuditByStatusDTO{}

	query := r.db.Table("audits").
		Select(`audit_statuses.id AS audit_status_id,
				audit_statuses.title_en AS audit_status_title_en,
				audit_statuses.title_ar AS audit_status_title_ar,
				audit_statuses.color_code AS color_code,
				COUNT(audits.id) AS count`).
		Joins("JOIN audit_statuses ON audit_statuses.id = audits.audit_status_id").
		Group("audit_statuses.id, audit_statuses.title_en, audit_statuses.title_ar, audit_statuses.color_code")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("audits.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (r *AuditReportRepository) GetAuditByOpenVsClosed(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditByOpenVsClosedDTO, error) {
	var rawResults []struct {
		ID        uint   `gorm:"column:id"`
		Type      string `gorm:"column:type"`
		TitleEn   string `gorm:"column:title_en"`
		TitleAr   string `gorm:"column:title_ar"`
		ColorCode string `gorm:"column:color_code"`
	}

	// Step 1: Fetch audit statuses with their details
	query := r.db.Table("audits").
		Select(`audit_statuses.id, audit_statuses.type, 
				audit_statuses.title_en, audit_statuses.title_ar, 
				audit_statuses.color_code`).
		Joins("JOIN audit_statuses ON audit_statuses.id = audits.audit_status_id")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("audits.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&rawResults).Error
	if err != nil {
		return nil, err
	}

	// Step 2: Process data in Go
	openCount := 0
	closedCount := 0
	var openStatus, closedStatus *reportDTOs.AuditByOpenVsClosedDTO

	for _, record := range rawResults {
		if record.Type == "published" || record.Type == "completed" {
			closedCount++
			if closedStatus == nil { // Initialize Closed status
				closedStatus = &reportDTOs.AuditByOpenVsClosedDTO{
					ID:        record.ID,
					TitleEn:   "Closed",
					TitleAr:   "مغلق", // Use Arabic title from DB
					Type:      "closed",
					ColorCode: record.ColorCode,
					Count:     0,
				}
			}
		} else {
			openCount++
			if openStatus == nil { // Initialize Open status
				openStatus = &reportDTOs.AuditByOpenVsClosedDTO{
					ID:        record.ID,
					TitleEn:   "Open",
					TitleAr:   "مفتوح", // Use Arabic title from DB
					ColorCode: record.ColorCode,
					Type:      "open",
					Count:     0,
				}
			}
		}
	}

	// Assign final counts
	if openStatus != nil {
		openStatus.Count = uint(openCount)
	}
	if closedStatus != nil {
		closedStatus.Count = uint(closedCount)
	}

	// Step 3: Prepare the response as a slice of pointers
	results := []*reportDTOs.AuditByOpenVsClosedDTO{}
	if openStatus != nil {
		results = append(results, openStatus)
	}
	if closedStatus != nil {
		results = append(results, closedStatus)
	}

	return results, nil
}

func (r *AuditReportRepository) GetAuditByMonth(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditByMonthDTO, error) {
	results := []*reportDTOs.AuditByMonthDTO{}

	// Define static Arabic month translations
	monthTranslations := map[string]string{
		"Jan": "يناير", "Feb": "فبراير", "Mar": "مارس", "Apr": "أبريل",
		"May": "مايو", "Jun": "يونيو", "Jul": "يوليو", "Aug": "أغسطس",
		"Sep": "سبتمبر", "Oct": "أكتوبر", "Nov": "نوفمبر", "Dec": "ديسمبر",
	}

	// Fetch last 12 months' audit count
	var rawResults []struct {
		MonthEn string `gorm:"column:title_en"`
		Count   uint   `gorm:"column:count"`
	}

	query := r.db.Table("audits").
		Select(`FORMAT(audits.created_at, 'MMM yyyy') AS title_en, COUNT(audits.id) AS count`).
		Where("audits.created_at >= DATEADD(MONTH, -11, GETDATE())").
		Group("FORMAT(audits.created_at, 'MMM yyyy')").
		Order("MIN(audits.created_at) ASC")

	err := query.Scan(&rawResults).Error
	if err != nil {
		return nil, err
	}

	// Convert English month names to Arabic and store results
	for _, record := range rawResults {
		monthShort := record.MonthEn[:3] // Extract three-letter month abbreviation (e.g., "Jan", "Feb")
		results = append(results, &reportDTOs.AuditByMonthDTO{
			TitleEn: record.MonthEn,
			TitleAr: monthTranslations[monthShort] + " " + record.MonthEn[4:], // Arabic + Year
			Count:   record.Count,
		})
	}

	return results, nil
}

func (r *AuditReportRepository) GetAuditByTeamLead(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditByTeamLeadDTO, error) {
	results := []*reportDTOs.AuditByTeamLeadDTO{}

	query := r.db.Table("audits").
		Select(`users.id AS team_lead_id,
			users.first_name_en AS team_lead_first_name_en,
			users.first_name_ar AS team_lead_first_name_ar,
			users.last_name_en AS team_lead_last_name_en,
			users.last_name_ar AS team_lead_last_name_ar,
			CONCAT(users.first_name_en, ' ', users.last_name_en) AS team_lead_title_en,
			CONCAT(users.first_name_ar, ' ', users.last_name_ar) AS team_lead_title_ar,
			users.initial AS team_lead_initial,
			users.color AS team_lead_color,
			users.email AS team_lead_email,
			users.profile_image_path AS team_lead_profile_image_path,
			COUNT(audit_plans.id) AS count`).
		Joins("JOIN audit_plans ON audit_plans.id = audits.audit_plan_id").
		Joins("JOIN users ON users.id = audit_plans.team_lead_id").
		Group("users.id, users.first_name_en, users.first_name_ar, users.last_name_en, users.last_name_ar, users.initial, users.color, users.email, users.profile_image_path")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("audits.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (r *AuditReportRepository) GetAuditFindingByCategory(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditFindingByCategoryDTO, error) {
	results := []*reportDTOs.AuditFindingByCategoryDTO{}

	query := r.db.Table("audit_findings").
		Select(`audit_finding_categories.id AS id, 
				audit_finding_categories.title_en AS title_en, 
				audit_finding_categories.title_ar AS title_ar, 
				COUNT(audit_findings.id) AS count`).
		Joins("JOIN audit_finding_categories ON audit_finding_categories.id = audit_findings.audit_finding_category_id").
		Group("audit_finding_categories.id, audit_finding_categories.title_en, audit_finding_categories.title_ar").
		Order("COUNT(audit_findings.id) DESC") // Order by highest count

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("audit_findings.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (r *AuditReportRepository) GetAuditFindingByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditFindingByStatusDTO, error) {
	results := []*reportDTOs.AuditFindingByStatusDTO{}

	query := r.db.Table("audit_findings").
		Select(`audit_finding_statuses.id,
                audit_finding_statuses.title_en,
                audit_finding_statuses.title_ar,
                audit_finding_statuses.color_code,
                COUNT(audit_findings.id) AS count`).
		Joins("JOIN audit_finding_statuses ON audit_finding_statuses.id = audit_findings.audit_finding_status_id").
		Group("audit_finding_statuses.id, audit_finding_statuses.title_en, audit_finding_statuses.title_ar, audit_finding_statuses.color_code")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("audit_findings.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (r *AuditReportRepository) GetAuditFindingBySeverityLevel(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditFindingBySeverityLevelDTO, error) {
	results := []*reportDTOs.AuditFindingBySeverityLevelDTO{}

	query := r.db.Table("audit_findings").
		Select(`severity_levels.id,
                severity_levels.title_en,
                severity_levels.title_ar,
                severity_levels.color_code,
                COUNT(audit_findings.id) AS count`).
		Joins("JOIN severity_levels ON severity_levels.id = audit_findings.severity_level_id").
		Group("severity_levels.id, severity_levels.title_en, severity_levels.title_ar, severity_levels.color_code")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("audit_findings.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (r *AuditReportRepository) GetAuditFindingByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditFindingByResponsibleUnitDTO, error) {
	results := []*reportDTOs.AuditFindingByResponsibleUnitDTO{}

	query := r.db.Table("audit_findings").
		Select(`
            organization_units.id AS responsible_unit_id,
            organization_units.title_en AS responsible_unit_title_en,
            organization_units.title_ar AS responsible_unit_title_ar,
            organization_units.color_code AS responsible_unit_color_code,
            COUNT(audit_findings.id) AS count
        `).
		Joins("JOIN auditable_items ON auditable_items.id = audit_findings.auditable_item_id").
		Joins("JOIN organization_units ON organization_units.id = auditable_items.responsible_unit_id").
		Group("organization_units.id, organization_units.title_en, organization_units.title_ar, organization_units.color_code")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("audit_findings.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (r *AuditReportRepository) GetAuditFindingByMonth(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditFindingByMonthDTO, error) {
	results := []*reportDTOs.AuditFindingByMonthDTO{}

	// Define static Arabic month translations
	monthTranslations := map[string]string{
		"Jan": "يناير", "Feb": "فبراير", "Mar": "مارس", "Apr": "أبريل",
		"May": "مايو", "Jun": "يونيو", "Jul": "يوليو", "Aug": "أغسطس",
		"Sep": "سبتمبر", "Oct": "أكتوبر", "Nov": "نوفمبر", "Dec": "ديسمبر",
	}

	// Fetch last 12 months' audit finding count
	var rawResults []struct {
		MonthEn string `gorm:"column:title_en"`
		Count   uint   `gorm:"column:count"`
	}

	query := r.db.Table("audit_findings").
		Select(`FORMAT(audit_findings.created_at, 'MMM yyyy') AS title_en, COUNT(audit_findings.id) AS count`).
		Where("audit_findings.created_at >= DATEADD(MONTH, -11, GETDATE())").
		Group("FORMAT(audit_findings.created_at, 'MMM yyyy')").
		Order("MIN(audit_findings.created_at) ASC")

	err := query.Scan(&rawResults).Error
	if err != nil {
		return nil, err
	}

	// Convert English month names to Arabic and store results
	for _, record := range rawResults {
		monthShort := record.MonthEn[:3] // Extract three-letter month abbreviation (e.g., "Jan", "Feb")
		results = append(results, &reportDTOs.AuditFindingByMonthDTO{
			TitleEn: record.MonthEn,
			TitleAr: monthTranslations[monthShort] + " " + record.MonthEn[4:], // Arabic + Year
			Count:   record.Count,
		})
	}

	return results, nil
}

func (r *AuditReportRepository) GetAuditFindingByAuditor(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditFindingByAuditorDTO, error) {
	results := []*reportDTOs.AuditFindingByAuditorDTO{}

	query := r.db.Table("audit_findings").
		Select(`
            users.id AS id,
            users.first_name_en AS user_first_name_en,
            users.first_name_ar AS user_first_name_ar,
            users.last_name_en AS user_last_name_en,
            users.last_name_ar AS user_last_name_ar,
            CONCAT(users.first_name_en, ' ', users.last_name_en) AS user_title_en,
            CONCAT(users.first_name_ar, ' ', users.last_name_ar) AS user_title_ar,
            users.initial AS user_initial,
            users.color AS user_color,
            users.email AS user_email,
            users.profile_image_path AS user_profile_image_path,
            COUNT(audit_findings.id) AS count
        `).
		Joins("JOIN users ON users.id = audit_findings.created_by_id").
		Group("users.id, users.first_name_en, users.first_name_ar, users.last_name_en, users.last_name_ar, users.initial, users.color, users.email, users.profile_image_path")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("audit_findings.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (r *AuditReportRepository) GetAuditActionPlanByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditActionPlanByStatusDTO, error) {
	results := []*reportDTOs.AuditActionPlanByStatusDTO{}

	query := r.db.Table("audit_finding_action_plans").
		Select(`
            action_plan_statuses.id,
            action_plan_statuses.title_en,
            action_plan_statuses.title_ar,
            action_plan_statuses.color_code,
            COUNT(audit_finding_action_plans.id) AS count
        `).
		Joins("JOIN action_plan_statuses ON action_plan_statuses.id = audit_finding_action_plans.action_plan_status_id").
		Group("action_plan_statuses.id, action_plan_statuses.title_en, action_plan_statuses.title_ar, action_plan_statuses.color_code")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("audit_finding_action_plans.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (r *AuditReportRepository) GetAuditActionPlanBySeverityLevel(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditActionPlanBySeverityLevelDTO, error) {
	results := []*reportDTOs.AuditActionPlanBySeverityLevelDTO{}

	query := r.db.Table("audit_finding_action_plans").
		Select(`severity_levels.id,
                severity_levels.title_en,
                severity_levels.title_ar,
                severity_levels.color_code,
                COUNT(audit_finding_action_plans.id) AS count`).
		Joins("JOIN severity_levels ON severity_levels.id = audit_finding_action_plans.severity_level_id").
		Group("severity_levels.id, severity_levels.title_en, severity_levels.title_ar, severity_levels.color_code")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("audit_finding_action_plans.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (r *AuditReportRepository) GetAuditActionPlanByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditActionPlanByResponsibleUnitDTO, error) {
	results := []*reportDTOs.AuditActionPlanByResponsibleUnitDTO{}

	query := r.db.Table("audit_finding_action_plans").
		Select(`
			organization_units.id AS responsible_unit_id,
			organization_units.title_en AS responsible_unit_title_en,
			organization_units.title_ar AS responsible_unit_title_ar,
			organization_units.color_code AS color_code,
			COUNT(audit_finding_action_plans.id) AS count
        `).
		Joins("JOIN audit_findings ON audit_findings.id = audit_finding_action_plans.audit_finding_id").
		Joins("JOIN auditable_items ON auditable_items.id = audit_findings.auditable_item_id").
		Joins("JOIN organization_units ON organization_units.id = auditable_items.responsible_unit_id").
		Group("organization_units.id, organization_units.title_en, organization_units.title_ar, organization_units.color_code")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("audit_finding_action_plans.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (r *AuditReportRepository) GetAuditActionPlanByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditActionPlanByResponsibleUserDTO, error) {
	results := []*reportDTOs.AuditActionPlanByResponsibleUserDTO{}

	query := r.db.Table("audit_finding_action_plans").
		Select(`
            users.id AS responsible_user_id,
            users.first_name_en AS responsible_user_first_name_en,
            users.first_name_ar AS responsible_user_first_name_ar,
            users.last_name_en AS responsible_user_last_name_en,
            users.last_name_ar AS responsible_user_last_name_ar,
            CONCAT(users.first_name_en, ' ', users.last_name_en) AS responsible_user_title_en,
            CONCAT(users.first_name_ar, ' ', users.last_name_ar) AS responsible_user_title_ar,
            users.initial AS responsible_user_initial,
            users.color AS responsible_user_color,
            users.email AS responsible_user_email,
            users.profile_image_path AS responsible_user_profile_image_path,
            COUNT(audit_finding_action_plans.id) AS count
        `).
		Joins("JOIN users ON users.id = audit_finding_action_plans.responsible_user_id").
		Group("users.id, users.first_name_en, users.first_name_ar, users.last_name_en, users.last_name_ar, users.initial, users.color, users.email, users.profile_image_path")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("audit_finding_action_plans.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (r *AuditReportRepository) GetAuditActionPlanOnTimeVsDelayed(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditActionPlanOnTimeVsDelayedDTO, error) {
	var results []*reportDTOs.AuditActionPlanOnTimeVsDelayedDTO

	// Query for on-time action plans
	var onTimeCount uint
	onTimeQuery := r.db.Table("audit_finding_action_plans").
		Select("COUNT(audit_finding_action_plans.id)"). // Fully qualify the id column
		Joins("JOIN action_plan_statuses ON action_plan_statuses.id = audit_finding_action_plans.action_plan_status_id").
		Where("action_plan_statuses.type IN ('completed', 'close')")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		onTimeQuery = onTimeQuery.Where("audit_finding_action_plans.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := onTimeQuery.Scan(&onTimeCount).Error
	if err != nil {
		return nil, err
	}

	// Add on-time result to the response
	results = append(results, &reportDTOs.AuditActionPlanOnTimeVsDelayedDTO{
		TitleEn: "On Time",
		TitleAr: "في الوقت المحدد",
		Count:   onTimeCount,
	})

	// Query for delayed action plans
	var delayedCount uint
	delayedQuery := r.db.Table("audit_finding_action_plans").
		Select("COUNT(audit_finding_action_plans.id)"). // Fully qualify the id column
		Joins("JOIN action_plan_statuses ON action_plan_statuses.id = audit_finding_action_plans.action_plan_status_id").
		Where("action_plan_statuses.type NOT IN ('completed', 'close')").
		Where("audit_finding_action_plans.target_date < GETDATE()")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		delayedQuery = delayedQuery.Where("audit_finding_action_plans.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err = delayedQuery.Scan(&delayedCount).Error
	if err != nil {
		return nil, err
	}

	// Add delayed result to the response
	results = append(results, &reportDTOs.AuditActionPlanOnTimeVsDelayedDTO{
		TitleEn: "Delayed",
		TitleAr: "متأخر",
		Count:   delayedCount,
	})

	return results, nil
}
