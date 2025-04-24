package repositories

import (
	"fmt"
	reportDTOs "jethur/grc/reports/dtos"
	"time"

	"gorm.io/gorm"
)

type AuditInsightReportRepositoryInterface interface {
	GetAuditByStatus(params *reportDTOs.ReportListQueryParamsDTO, currentUserResponsibleUnitID uint) ([]*reportDTOs.AuditInsightAuditByStatusDTO, error)
	GetAuditByYear(params *reportDTOs.ReportListQueryParamsDTO, currentUserResponsibleUnitID uint) ([]*reportDTOs.AuditInsightAuditByYearDTO, error)
	GetAuditByMonth(params *reportDTOs.ReportListQueryParamsDTO, currentUserResponsibleUnitID uint) ([]*reportDTOs.AuditInsightAuditByMonthDTO, error)

	GetAuditPlanByCategory(params *reportDTOs.ReportListQueryParamsDTO, currentUserResponsibleUnitID uint) ([]*reportDTOs.AuditInsightAuditPlanByCategoryDTO, error)
	GetAuditPlanByTeamLead(params *reportDTOs.ReportListQueryParamsDTO, currentUserResponsibleUnitID uint) ([]*reportDTOs.AuditInsightAuditPlanByTeamLeadDTO, error)
	GetAuditPlanByAuditableItem(params *reportDTOs.ReportListQueryParamsDTO, currentUserResponsibleUnitID uint) ([]*reportDTOs.AuditInsightAuditPlanByAuditableItemDTO, error)

	GetAuditFindingByCategory(params *reportDTOs.ReportListQueryParamsDTO, currentUserResponsibleUnitID uint) ([]*reportDTOs.AuditInsightAuditFindingByCategoryDTO, error)
	GetAuditFindingByStatus(params *reportDTOs.ReportListQueryParamsDTO, currentUserResponsibleUnitID uint) ([]*reportDTOs.AuditInsightAuditFindingByStatusDTO, error)
	GetAuditFindingBySeverityLevel(params *reportDTOs.ReportListQueryParamsDTO, currentUserResponsibleUnitID uint) ([]*reportDTOs.AuditInsightAuditFindingBySeverityLevelDTO, error)
	GetAuditFindingByMonth(params *reportDTOs.ReportListQueryParamsDTO, currentUserResponsibleUnitID uint) ([]*reportDTOs.AuditInsightAuditFindingByMonthDTO, error)
	GetAuditFindingByYear(params *reportDTOs.ReportListQueryParamsDTO, currentUserResponsibleUnitID uint) ([]*reportDTOs.AuditInsightAuditFindingByYearDTO, error)

	GetAuditActionPlanByStatus(params *reportDTOs.ReportListQueryParamsDTO, currentUserResponsibleUnitID uint) ([]*reportDTOs.AuditInsightAuditActionPlanByStatusDTO, error)
	GetAuditActionPlanByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO, currentUserResponsibleUnitID uint) ([]*reportDTOs.AuditInsightAuditActionPlanByResponsibleUserDTO, error)

	GetAuditDraftReportByStatus(params *reportDTOs.ReportListQueryParamsDTO, currentUserResponsibleUnitID uint) ([]*reportDTOs.AuditInsightAuditDraftReportByStatusDTO, error)
	GetAuditFinalReportByStatus(params *reportDTOs.ReportListQueryParamsDTO, currentUserResponsibleUnitID uint) ([]*reportDTOs.AuditInsightAuditFinalReportByStatusDTO, error)
}

type AuditInsightReportRepository struct {
	db *gorm.DB
}

func NewAuditInsightReportRepository(db *gorm.DB) *AuditInsightReportRepository {
	return &AuditInsightReportRepository{db: db}
}

func (repo *AuditInsightReportRepository) GetAuditByStatus(params *reportDTOs.ReportListQueryParamsDTO, currentUserResponsibleUnitID uint) ([]*reportDTOs.AuditInsightAuditByStatusDTO, error) {
	results := []*reportDTOs.AuditInsightAuditByStatusDTO{}

	query := repo.db.Table("audit_statuses").
		Select(`
			audit_statuses.id AS audit_status_id,
			audit_statuses.title_en AS audit_status_title_en,
			audit_statuses.title_ar AS audit_status_title_ar,
			audit_statuses.color_code AS color_code,
			COALESCE(COUNT(audits.id), 0) AS count
		`).
		Joins(`
			LEFT JOIN audits ON audits.audit_status_id = audit_statuses.id
			LEFT JOIN audit_plan_auditable_items ON audit_plan_auditable_items.audit_plan_id = audits.audit_plan_id
			LEFT JOIN auditable_items ON auditable_items.id = audit_plan_auditable_items.auditable_item_id
		`).
		Where(`(auditable_items.responsible_unit_id = ? OR audits.id IS NULL)`, currentUserResponsibleUnitID).
		Group("audit_statuses.id, audit_statuses.title_en, audit_statuses.title_ar, audit_statuses.color_code")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("audits.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *AuditInsightReportRepository) GetAuditByYear(params *reportDTOs.ReportListQueryParamsDTO, currentUserResponsibleUnitID uint) ([]*reportDTOs.AuditInsightAuditByYearDTO, error) {
	results := []*reportDTOs.AuditInsightAuditByYearDTO{}

	// Fetch yearly audit plan count
	var rawResults []struct {
		Year  int  `gorm:"column:year"`
		Count uint `gorm:"column:count"`
	}

	query := repo.db.Table("audits").
		Select(`
			YEAR(audits.created_at) AS year,
			COUNT(DISTINCT audits.id) AS count
		`).
		Joins(`
			LEFT JOIN audit_plan_auditable_items ON audits.audit_plan_id = audit_plan_auditable_items.audit_plan_id
			LEFT JOIN auditable_items ON auditable_items.id = audit_plan_auditable_items.auditable_item_id
		`).
		Where("auditable_items.responsible_unit_id = ?", currentUserResponsibleUnitID).
		Group("YEAR(audits.created_at)").
		Order("YEAR(audits.created_at) ASC")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("audits.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&rawResults).Error
	if err != nil {
		return nil, err
	}

	// Convert year to English and Arabic title
	for _, record := range rawResults {
		results = append(results, &reportDTOs.AuditInsightAuditByYearDTO{
			TitleEn: fmt.Sprintf("%d", record.Year), // English Title (Year)
			TitleAr: fmt.Sprintf("%d", record.Year), // Arabic Title (Same as Year)
			Count:   record.Count,
		})
	}

	return results, nil
}

func (repo *AuditInsightReportRepository) GetAuditByMonth(params *reportDTOs.ReportListQueryParamsDTO, currentUserResponsibleUnitID uint) ([]*reportDTOs.AuditInsightAuditByMonthDTO, error) {
	results := []*reportDTOs.AuditInsightAuditByMonthDTO{}

	// Define static Arabic month translations
	monthTranslations := map[string]string{
		"Jan": "يناير", "Feb": "فبراير", "Mar": "مارس", "Apr": "أبريل",
		"May": "مايو", "Jun": "يونيو", "Jul": "يوليو", "Aug": "أغسطس",
		"Sep": "سبتمبر", "Oct": "أكتوبر", "Nov": "نوفمبر", "Dec": "ديسمبر",
	}

	// Generate the last 12 months
	last12Months := make(map[string]*reportDTOs.AuditInsightAuditByMonthDTO)
	now := time.Now()
	for i := 0; i < 12; i++ {
		month := now.AddDate(0, -i, 0).Format("Jan 2006") // e.g., "Jan 2025"
		monthShort := month[:3]                           // Extract short form (e.g., "Jan")
		last12Months[month] = &reportDTOs.AuditInsightAuditByMonthDTO{
			TitleEn: month,
			TitleAr: monthTranslations[monthShort] + " " + month[4:], // Arabic + Year
			Count:   0,                                               // Default count
		}
	}

	// Fetch data for the last 12 months
	var rawResults []struct {
		MonthEn string `gorm:"column:title_en"`
		Count   uint   `gorm:"column:count"`
	}

	query := repo.db.Table("audits").
		Select(`FORMAT(audits.created_at, 'MMM yyyy') AS title_en, COUNT(audits.id) AS count`).
		Joins(`
			LEFT JOIN audit_plan_auditable_items ON audits.audit_plan_id = audit_plan_auditable_items.audit_plan_id
			LEFT JOIN auditable_items ON auditable_items.id = audit_plan_auditable_items.auditable_item_id
		`).
		Where("auditable_items.responsible_unit_id = ?", currentUserResponsibleUnitID).
		Where("audits.created_at >= DATEADD(MONTH, -11, GETDATE())").
		Group("FORMAT(audits.created_at, 'MMM yyyy')").
		Order("MIN(audits.created_at) ASC")

	err := query.Scan(&rawResults).Error
	if err != nil {
		return nil, err
	}

	// Merge database results into the last 12 months
	for _, record := range rawResults {
		if dto, exists := last12Months[record.MonthEn]; exists {
			dto.Count = record.Count // Update count for the matching month
		}
	}

	// Convert map to a slice in chronological order
	for i := 11; i >= 0; i-- {
		month := now.AddDate(0, -i, 0).Format("Jan 2006") // Same key format
		results = append(results, last12Months[month])
	}

	return results, nil
}

func (repo *AuditInsightReportRepository) GetAuditPlanByCategory(params *reportDTOs.ReportListQueryParamsDTO, currentUserResponsibleUnitID uint) ([]*reportDTOs.AuditInsightAuditPlanByCategoryDTO, error) {
	results := []*reportDTOs.AuditInsightAuditPlanByCategoryDTO{}

	// Start the query from audit_categories
	query := repo.db.Table("audit_categories").
		Select(`
			audit_categories.id AS audit_category_id, 
			audit_categories.title_en AS audit_category_title_en, 
			audit_categories.title_ar AS audit_category_title_ar, 
			COUNT(DISTINCT audit_plans.id) AS count
		`).
		Joins(`
			LEFT JOIN audit_plans ON audit_plans.audit_category_id = audit_categories.id
			LEFT JOIN audit_plan_auditable_items ON audit_plan_auditable_items.audit_plan_id = audit_plans.id
			LEFT JOIN auditable_items ON auditable_items.id = audit_plan_auditable_items.auditable_item_id 
				AND auditable_items.responsible_unit_id = ?
		`, currentUserResponsibleUnitID).
		Group("audit_categories.id, audit_categories.title_en, audit_categories.title_ar")

	// Apply date filters if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("audit_plans.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	// Execute the query
	err := query.Scan(&results).Error
	return results, err
}

func (repo *AuditInsightReportRepository) GetAuditPlanByTeamLead(params *reportDTOs.ReportListQueryParamsDTO, currentUserResponsibleUnitID uint) ([]*reportDTOs.AuditInsightAuditPlanByTeamLeadDTO, error) {
	results := []*reportDTOs.AuditInsightAuditPlanByTeamLeadDTO{}

	query := repo.db.Table("audit_plans").
		Select(`
			users.id AS team_lead_id,
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
			COUNT(audit_plans.id) AS count
		`).
		Joins(`
			JOIN users ON users.id = audit_plans.team_lead_id
			JOIN audit_plan_auditable_items ON audit_plan_auditable_items.audit_plan_id = audit_plans.id
			JOIN auditable_items ON auditable_items.id = audit_plan_auditable_items.auditable_item_id
		`).
		Where("auditable_items.responsible_unit_id = ?", currentUserResponsibleUnitID).
		Group("users.id, users.first_name_en, users.first_name_ar, users.last_name_en, users.last_name_ar, users.initial, users.color, users.email, users.profile_image_path")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("audit_plans.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *AuditInsightReportRepository) GetAuditPlanByAuditableItem(params *reportDTOs.ReportListQueryParamsDTO, currentUserResponsibleUnitID uint) ([]*reportDTOs.AuditInsightAuditPlanByAuditableItemDTO, error) {
	results := []*reportDTOs.AuditInsightAuditPlanByAuditableItemDTO{}

	// Query to fetch audit plans grouped by auditable items
	query := repo.db.Table("auditable_items").
		Select(`
			auditable_items.id AS auditable_item_id, 
			COALESCE(CASE 
				WHEN auditable_items.type = 'control' THEN controls.title_en
				WHEN auditable_items.type = 'process' THEN processes.title_en
				WHEN auditable_items.type = 'risk' THEN risk_register.title
				WHEN auditable_items.type = 'compliance_requirement' THEN compliance_requirements.title_en
				WHEN auditable_items.type = 'incident' THEN incident_register.title
				WHEN auditable_items.type = 'user' THEN CONCAT(users.first_name_en, ' ', users.last_name_en)
				WHEN auditable_items.type = 'document' THEN documents.title_en
				ELSE auditable_items.title_en
			END, '') AS auditable_item_title_en,
			
			COALESCE(CASE 
				WHEN auditable_items.type = 'control' THEN controls.title_ar
				WHEN auditable_items.type = 'process' THEN processes.title_ar
				WHEN auditable_items.type = 'compliance_requirement' THEN compliance_requirements.title_ar
				WHEN auditable_items.type = 'document' THEN documents.title_ar
				ELSE auditable_items.title_ar
			END, '') AS auditable_item_title_ar,

			COALESCE(COUNT(DISTINCT audit_plans.id), 0) AS count
		`).
		Joins(`LEFT JOIN audit_plan_auditable_items ON audit_plan_auditable_items.auditable_item_id = auditable_items.id`).
		Joins(`LEFT JOIN audit_plans ON audit_plans.id = audit_plan_auditable_items.audit_plan_id`).
		// Correct table names
		Joins(`LEFT JOIN controls ON controls.id = auditable_items.control_id`).
		Joins(`LEFT JOIN processes ON processes.id = auditable_items.process_id`).
		Joins(`LEFT JOIN risk_register ON risk_register.id = auditable_items.risk_register_id`).
		Joins(`LEFT JOIN compliance_requirements ON compliance_requirements.id = auditable_items.compliance_requirement_id`).
		Joins(`LEFT JOIN incident_register ON incident_register.id = auditable_items.incident_register_id`).
		Joins(`LEFT JOIN users ON users.id = auditable_items.user_id`).
		Joins(`LEFT JOIN documents ON documents.id = auditable_items.document_id`).
		Where("auditable_items.responsible_unit_id = ?", currentUserResponsibleUnitID).
		Group("auditable_items.id, auditable_items.type, auditable_items.title_en, auditable_items.title_ar, controls.title_en, controls.title_ar, processes.title_en, processes.title_ar, risk_register.title, compliance_requirements.title_en, compliance_requirements.title_ar, incident_register.title, users.first_name_en, users.last_name_en, users.first_name_ar, users.last_name_ar, documents.title_en, documents.title_ar")

	// Apply date filters if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("audit_plans.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	// Execute the query
	err := query.Scan(&results).Error
	return results, err
}

func (repo *AuditInsightReportRepository) GetAuditFindingByCategory(params *reportDTOs.ReportListQueryParamsDTO, currentUserResponsibleUnitID uint) ([]*reportDTOs.AuditInsightAuditFindingByCategoryDTO, error) {
	results := []*reportDTOs.AuditInsightAuditFindingByCategoryDTO{}

	query := repo.db.Table("audit_finding_categories").
		Select(`audit_finding_categories.id AS audit_finding_category_id, 
				audit_finding_categories.title_en AS audit_finding_category_title_en, 
				audit_finding_categories.title_ar AS audit_finding_category_title_ar, 
				COUNT(audit_findings.id) AS count`).
		Joins(`
			LEFT JOIN audit_findings ON audit_findings.audit_finding_category_id = audit_finding_categories.id
			LEFT JOIN auditable_items ON auditable_items.id = audit_findings.auditable_item_id
		`).
		Where("auditable_items.responsible_unit_id = ?", currentUserResponsibleUnitID).
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

func (repo *AuditInsightReportRepository) GetAuditFindingByStatus(params *reportDTOs.ReportListQueryParamsDTO, currentUserResponsibleUnitID uint) ([]*reportDTOs.AuditInsightAuditFindingByStatusDTO, error) {
	results := []*reportDTOs.AuditInsightAuditFindingByStatusDTO{}

	query := repo.db.Table("audit_finding_statuses").
		Select(`audit_finding_statuses.id,
                audit_finding_statuses.title_en,
                audit_finding_statuses.title_ar,
                audit_finding_statuses.color_code,
                COALESCE(COUNT(audit_findings.id), 0) AS count`).
		Joins(`
			LEFT JOIN audit_findings ON audit_findings.audit_finding_status_id = audit_finding_statuses.id
			LEFT JOIN auditable_items ON auditable_items.id = audit_findings.auditable_item_id
		`).
		Where(`(auditable_items.responsible_unit_id = ? OR audit_findings.id IS NULL)`, currentUserResponsibleUnitID).
		Group("audit_finding_statuses.id, audit_finding_statuses.title_en, audit_finding_statuses.title_ar, audit_finding_statuses.color_code")

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

func (repo *AuditInsightReportRepository) GetAuditFindingBySeverityLevel(params *reportDTOs.ReportListQueryParamsDTO, currentUserResponsibleUnitID uint) ([]*reportDTOs.AuditInsightAuditFindingBySeverityLevelDTO, error) {
	results := []*reportDTOs.AuditInsightAuditFindingBySeverityLevelDTO{}

	query := repo.db.Table("severity_levels").
		Select(`severity_levels.id,
                severity_levels.title_en,
                severity_levels.title_ar,
                severity_levels.color_code,
                COUNT(audit_findings.id) AS count`).
		Joins(`
			LEFT JOIN audit_findings ON audit_findings.severity_level_id = severity_levels.id
			LEFT JOIN auditable_items ON auditable_items.id = audit_findings.auditable_item_id
		`).
		Where("auditable_items.responsible_unit_id = ?", currentUserResponsibleUnitID).
		Group("severity_levels.id, severity_levels.title_en, severity_levels.title_ar, severity_levels.color_code")

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

func (repo *AuditInsightReportRepository) GetAuditFindingByMonth(params *reportDTOs.ReportListQueryParamsDTO, currentUserResponsibleUnitID uint) ([]*reportDTOs.AuditInsightAuditFindingByMonthDTO, error) {
	results := []*reportDTOs.AuditInsightAuditFindingByMonthDTO{}

	// Define static Arabic month translations
	monthTranslations := map[string]string{
		"Jan": "يناير", "Feb": "فبراير", "Mar": "مارس", "Apr": "أبريل",
		"May": "مايو", "Jun": "يونيو", "Jul": "يوليو", "Aug": "أغسطس",
		"Sep": "سبتمبر", "Oct": "أكتوبر", "Nov": "نوفمبر", "Dec": "ديسمبر",
	}

	// Generate the last 12 months
	last12Months := make(map[string]*reportDTOs.AuditInsightAuditFindingByMonthDTO)
	now := time.Now()
	for i := 0; i < 12; i++ {
		month := now.AddDate(0, -i, 0).Format("Jan 2006") // e.g., "Jan 2025"
		monthShort := month[:3]                           // Extract short form (e.g., "Jan")
		last12Months[month] = &reportDTOs.AuditInsightAuditFindingByMonthDTO{
			TitleEn: month,
			TitleAr: monthTranslations[monthShort] + " " + month[4:], // Arabic + Year
			Count:   0,                                               // Default count
		}
	}

	// Fetch data for the last 12 months
	var rawResults []struct {
		MonthEn string `gorm:"column:title_en"`
		Count   uint   `gorm:"column:count"`
	}

	query := repo.db.Table("audit_findings").
		Select(`FORMAT(audit_findings.created_at, 'MMM yyyy') AS title_en, COUNT(audit_findings.id) AS count`).
		Joins(`LEFT JOIN auditable_items ON auditable_items.id = audit_findings.auditable_item_id`).
		Where("auditable_items.responsible_unit_id = ?", currentUserResponsibleUnitID).
		Where("audit_findings.created_at >= DATEADD(MONTH, -11, GETDATE())").
		Group("FORMAT(audit_findings.created_at, 'MMM yyyy')").
		Order("MIN(audit_findings.created_at) ASC")

	err := query.Scan(&rawResults).Error
	if err != nil {
		return nil, err
	}

	// Merge database results into the last 12 months
	for _, record := range rawResults {
		if dto, exists := last12Months[record.MonthEn]; exists {
			dto.Count = record.Count // Update count for the matching month
		}
	}

	// Convert map to a slice in chronological order
	for i := 11; i >= 0; i-- {
		month := now.AddDate(0, -i, 0).Format("Jan 2006") // Same key format
		results = append(results, last12Months[month])
	}

	return results, nil
}

func (repo *AuditInsightReportRepository) GetAuditFindingByYear(params *reportDTOs.ReportListQueryParamsDTO, currentUserResponsibleUnitID uint) ([]*reportDTOs.AuditInsightAuditFindingByYearDTO, error) {
	results := []*reportDTOs.AuditInsightAuditFindingByYearDTO{}

	// Fetch yearly audit finding count
	var rawResults []struct {
		Year  int  `gorm:"column:year"`
		Count uint `gorm:"column:count"`
	}

	query := repo.db.Table("audit_findings").
		Select(`
			YEAR(audit_findings.created_at) AS year,
			COUNT(DISTINCT audit_findings.id) AS count
		`).
		Joins(`
			LEFT JOIN auditable_items ON auditable_items.id = audit_findings.auditable_item_id
		`).
		Where("auditable_items.responsible_unit_id = ?", currentUserResponsibleUnitID).
		Group("YEAR(audit_findings.created_at)").
		Order("YEAR(audit_findings.created_at) ASC")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("audit_findings.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&rawResults).Error
	if err != nil {
		return nil, err
	}

	// Convert year to English and Arabic title
	for _, record := range rawResults {
		results = append(results, &reportDTOs.AuditInsightAuditFindingByYearDTO{
			TitleEn: fmt.Sprintf("%d", record.Year), // English Title (Year)
			TitleAr: fmt.Sprintf("%d", record.Year), // Arabic Title (Same as Year)
			Count:   record.Count,
		})
	}

	return results, nil
}

func (repo *AuditInsightReportRepository) GetAuditActionPlanByStatus(params *reportDTOs.ReportListQueryParamsDTO, currentUserResponsibleUnitID uint) ([]*reportDTOs.AuditInsightAuditActionPlanByStatusDTO, error) {
	results := []*reportDTOs.AuditInsightAuditActionPlanByStatusDTO{}

	query := repo.db.Table("action_plan_statuses").
		Select(`
            action_plan_statuses.id,
            action_plan_statuses.title_en,
            action_plan_statuses.title_ar,
            action_plan_statuses.color_code,
            COALESCE(COUNT(audit_finding_action_plans.id), 0) AS count
        `).
		Joins(`
			LEFT JOIN audit_finding_action_plans ON audit_finding_action_plans.action_plan_status_id = action_plan_statuses.id
			LEFT JOIN audit_findings ON audit_findings.id = audit_finding_action_plans.audit_finding_id
			LEFT JOIN auditable_items ON auditable_items.id = audit_findings.auditable_item_id
		`).
		Where("auditable_items.responsible_unit_id = ?", currentUserResponsibleUnitID).
		Group("action_plan_statuses.id, action_plan_statuses.title_en, action_plan_statuses.title_ar, action_plan_statuses.color_code")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("audit_finding_action_plans.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *AuditInsightReportRepository) GetAuditActionPlanByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO, currentUserResponsibleUnitID uint) ([]*reportDTOs.AuditInsightAuditActionPlanByResponsibleUserDTO, error) {
	results := []*reportDTOs.AuditInsightAuditActionPlanByResponsibleUserDTO{}

	query := repo.db.Table("audit_finding_action_plans").
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
		Joins("JOIN audit_findings ON audit_findings.id = audit_finding_action_plans.audit_finding_id").
		Joins("JOIN auditable_items ON auditable_items.id = audit_findings.auditable_item_id").
		Joins("JOIN users ON users.id = auditable_items.responsible_user_id").
		Where("auditable_items.responsible_unit_id = ?", currentUserResponsibleUnitID).
		Group("users.id, users.first_name_en, users.last_name_en, users.first_name_ar, users.last_name_ar, users.initial, users.color, users.email, users.profile_image_path")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("audit_finding_action_plans.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *AuditInsightReportRepository) GetAuditDraftReportByStatus(params *reportDTOs.ReportListQueryParamsDTO, currentUserResponsibleUnitID uint) ([]*reportDTOs.AuditInsightAuditDraftReportByStatusDTO, error) {
	results := []*reportDTOs.AuditInsightAuditDraftReportByStatusDTO{}

	query := repo.db.Table("audit_draft_report_statuses").
		Select(`
			audit_draft_report_statuses.id,
			audit_draft_report_statuses.title_en,
			audit_draft_report_statuses.title_ar,
			audit_draft_report_statuses.color_code,
			COALESCE(COUNT(audit_draft_reports.id), 0) AS count
		`).
		Joins(`
			LEFT JOIN audit_draft_reports ON audit_draft_reports.audit_draft_report_status_id = audit_draft_report_statuses.id
			LEFT JOIN audits ON audits.id = audit_draft_reports.audit_id
			LEFT JOIN audit_plan_auditable_items ON audit_plan_auditable_items.audit_plan_id = audits.audit_plan_id
			LEFT JOIN auditable_items ON auditable_items.id = audit_plan_auditable_items.auditable_item_id
		`).
		Where(`(auditable_items.responsible_unit_id = ? OR audit_draft_reports.id IS NULL)`, currentUserResponsibleUnitID).
		Group("audit_draft_report_statuses.id, audit_draft_report_statuses.title_en, audit_draft_report_statuses.title_ar, audit_draft_report_statuses.color_code")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("audit_draft_reports.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *AuditInsightReportRepository) GetAuditFinalReportByStatus(params *reportDTOs.ReportListQueryParamsDTO, currentUserResponsibleUnitID uint) ([]*reportDTOs.AuditInsightAuditFinalReportByStatusDTO, error) {
	results := []*reportDTOs.AuditInsightAuditFinalReportByStatusDTO{}

	query := repo.db.Table("audit_final_report_statuses").
		Select(`
			audit_final_report_statuses.id,
			audit_final_report_statuses.title_en,
			audit_final_report_statuses.title_ar,
			audit_final_report_statuses.color_code,
			COALESCE(COUNT(audit_final_reports.id), 0) AS count
		`).
		Joins(`
			LEFT JOIN audit_final_reports ON audit_final_reports.audit_final_report_status_id = audit_final_report_statuses.id
			LEFT JOIN audits ON audits.id = audit_final_reports.audit_id
			LEFT JOIN audit_plan_auditable_items ON audit_plan_auditable_items.audit_plan_id = audits.audit_plan_id
			LEFT JOIN auditable_items ON auditable_items.id = audit_plan_auditable_items.auditable_item_id
		`).
		Where(`(auditable_items.responsible_unit_id = ? OR audit_final_reports.id IS NULL)`, currentUserResponsibleUnitID).
		Group("audit_final_report_statuses.id, audit_final_report_statuses.title_en, audit_final_report_statuses.title_ar, audit_final_report_statuses.color_code")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("audit_final_reports.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}
