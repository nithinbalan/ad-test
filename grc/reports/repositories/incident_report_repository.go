package repositories

import (
	"fmt"
	reportDTOs "jethur/grc/reports/dtos"

	"gorm.io/gorm"
)

type IncidentReportRepositoryInterface interface {
	GetIncidentByCategory(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentByCategoryDTO, error)
	GetIncidentByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentByStatusDTO, error)
	GetIncidentBySeverityLevel(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentBySeverityLevelDTO, error)
	GetIncidentByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentByResponsibleUnitDTO, error)
	GetIncidentByReportedUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentByReportedUserDTO, error)
	GetIncidentByReportedYear(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentByReportedYearDTO, error)

	GetIncidentByAsset(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentByAssetDTO, error)
	GetIncidentBySLA(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentBySLADTO, error)

	GetIncidentInvestigationSummary(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentInvestigationSummaryDTO, error)
	GetIncidentRCA(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentRCADTO, error)
	GetIncidentClosure(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentClosureDTO, error)
	GetOpenIncidentsWithRisk(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.OpenIncidentWithRiskDTO, error)

	GetIncidentPendingActionPlan(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentPendingActionPlanDTO, error)
	GetIncidentActionPlanByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentActionPlanByStatusDTO, error)
	GetIncidentActionPlanByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentActionPlanByResponsibleUnitDTO, error)
	GetIncidentActionPlanByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentActionPlanByResponsibleUserDTO, error)
}

type IncidentReportRepository struct {
	db *gorm.DB
}

func NewIncidentReportRepository(db *gorm.DB) *IncidentReportRepository {
	return &IncidentReportRepository{db: db}
}

func (repo *IncidentReportRepository) GetIncidentByCategory(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentByCategoryDTO, error) {
	results := []*reportDTOs.IncidentByCategoryDTO{}

	query := repo.db.Table("incident_categories").
		Select(`incident_categories.id AS incident_category_id, 
				incident_categories.title_en AS incident_category_title_en, 
				incident_categories.title_ar AS incident_category_title_ar, 
				COUNT(incident_register.id) AS count`).
		Joins("LEFT JOIN incident_register ON incident_register.incident_category_id = incident_categories.id").
		Group("incident_categories.id, incident_categories.title_en, incident_categories.title_ar").
		Order("COUNT(incident_register.id) DESC") // Order by highest count

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("incident_register.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *IncidentReportRepository) GetIncidentByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentByStatusDTO, error) {
	results := []*reportDTOs.IncidentByStatusDTO{}

	query := repo.db.Table("incident_statuses").
		Select(`incident_statuses.id,
                incident_statuses.title_en,
                incident_statuses.title_ar,
                incident_statuses.color_code,
                COUNT(incident_register.id) AS count`).
		Joins("LEFT JOIN incident_register ON incident_register.incident_status_id = incident_statuses.id").
		Group("incident_statuses.id, incident_statuses.title_en, incident_statuses.title_ar, incident_statuses.color_code")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("incident_register.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *IncidentReportRepository) GetIncidentBySeverityLevel(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentBySeverityLevelDTO, error) {
	results := []*reportDTOs.IncidentBySeverityLevelDTO{}

	query := repo.db.Table("severity_levels").
		Select(`severity_levels.id,
                severity_levels.title_en,
                severity_levels.title_ar,
                severity_levels.color_code,
                COUNT(incident_register.id) AS count`).
		Joins("LEFT JOIN incident_register ON incident_register.severity_level_id = severity_levels.id").
		Group("severity_levels.id, severity_levels.title_en, severity_levels.title_ar, severity_levels.color_code")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("incident_register.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *IncidentReportRepository) GetIncidentByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentByResponsibleUnitDTO, error) {
	results := []*reportDTOs.IncidentByResponsibleUnitDTO{}

	query := repo.db.Table("incident_register").
		Select(`
            organization_units.id AS responsible_unit_id,
            organization_units.title_en AS responsible_unit_title_en,
            organization_units.title_ar AS responsible_unit_title_ar,
            organization_units.color_code AS responsible_unit_color_code,
            COUNT(incident_register.id) AS count
        `).
		Joins("JOIN organization_units ON organization_units.id = incident_register.responsible_unit_id").
		Group("organization_units.id, organization_units.title_en, organization_units.title_ar, organization_units.color_code")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("incident_register.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *IncidentReportRepository) GetIncidentByReportedUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentByReportedUserDTO, error) {
	results := []*reportDTOs.IncidentByReportedUserDTO{}

	query := repo.db.Table("incident_register").
		Select(`
            users.id AS reported_user_id,
            users.first_name_en AS reported_user_first_name_en,
            users.first_name_ar AS reported_user_first_name_ar,
            users.last_name_en AS reported_user_last_name_en,
            users.last_name_ar AS reported_user_last_name_ar,
            CONCAT(users.first_name_en, ' ', users.last_name_en) AS reported_user_title_en,
            CONCAT(users.first_name_ar, ' ', users.last_name_ar) AS reported_user_title_ar,
            users.initial AS reported_user_initial,
            users.color AS reported_user_color,
            users.email AS reported_user_email,
            users.profile_image_path AS reported_user_profile_image_path,
            COUNT(incident_register.id) AS count
        `).
		Joins("JOIN users ON users.id = incident_register.reported_by").
		Group("users.id, users.first_name_en, users.first_name_ar, users.last_name_en, users.last_name_ar, users.initial, users.color, users.email, users.profile_image_path")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("incident_register.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *IncidentReportRepository) GetIncidentByReportedYear(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentByReportedYearDTO, error) {
	results := []*reportDTOs.IncidentByReportedYearDTO{}

	// Fetch yearly incident report count
	var rawResults []struct {
		Year  int  `gorm:"column:year"`
		Count uint `gorm:"column:count"`
	}

	query := repo.db.Table("incident_register").
		Select(`YEAR(reported_at) AS year, COUNT(id) AS count`).
		Group("YEAR(reported_at)").
		Order("YEAR(reported_at) ASC")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("reported_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&rawResults).Error
	if err != nil {
		return nil, err
	}

	// Convert year to English and Arabic title
	for _, record := range rawResults {
		results = append(results, &reportDTOs.IncidentByReportedYearDTO{
			TitleEn: fmt.Sprintf("%d", record.Year), // English Title (Year)
			TitleAr: fmt.Sprintf("%d", record.Year), // Arabic Title (Same as Year)
			Count:   record.Count,
		})
	}

	return results, nil
}

func (repo *IncidentReportRepository) GetIncidentByAsset(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentByAssetDTO, error) {
	results := []*reportDTOs.IncidentByAssetDTO{}

	query := repo.db.Table("asset_register").
		Select(`asset_register.id,
                asset_register.title_en,
                asset_register.title_ar,
                COUNT(incident_register.id) AS count`).
		Joins("LEFT JOIN incident_assets ON incident_assets.asset_register_id = asset_register.id").
		Joins("LEFT JOIN incident_register ON incident_register.id = incident_assets.incident_register_id").
		Group("asset_register.id, asset_register.title_en, asset_register.title_ar")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("incident_register.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *IncidentReportRepository) GetIncidentBySLA(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentBySLADTO, error) {
	results := []*reportDTOs.IncidentBySLADTO{}

	query := repo.db.Table("sla_register").
		Select(`sla_register.id,
                sla_register.title_en,
                sla_register.title_ar,
                COUNT(incident_register.id) AS count`).
		Joins("LEFT JOIN incident_slas ON incident_slas.sla_register_id = sla_register.id").
		Joins("LEFT JOIN incident_register ON incident_register.id = incident_slas.incident_register_id").
		Group("sla_register.id, sla_register.title_en, sla_register.title_ar")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("incident_register.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *IncidentReportRepository) GetIncidentInvestigationSummary(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentInvestigationSummaryDTO, error) {
	results := []*reportDTOs.IncidentInvestigationSummaryDTO{}

	query := repo.db.Raw(`
        WITH FindingsCTE AS (
            SELECT
                findings.incident_investigation_id,
                STRING_AGG(findings.title, ', ') AS findings_summary
            FROM incident_investigation_findings AS findings
            GROUP BY findings.incident_investigation_id
        ),
        RootCauseCTE AS (
            SELECT
                root_cause.incident_register_id,
                STRING_AGG(CONCAT(root_cause.question, ': ', root_cause.answer), ', ') AS root_cause_analysis
            FROM incident_root_cause_analyses AS root_cause
            GROUP BY root_cause.incident_register_id
        ),
        InvestigationUsersCTE AS (
            SELECT
                investigation_users.incident_investigation_id,
                STRING_AGG(CONCAT(users.first_name_en, ' ', users.last_name_en), ', ') AS investigation_teams_en,
                STRING_AGG(CONCAT(users.first_name_ar, ' ', users.last_name_ar), ', ') AS investigation_teams_ar
            FROM incident_investigation_users AS investigation_users
            INNER JOIN users ON users.id = investigation_users.user_id
            GROUP BY investigation_users.incident_investigation_id
        )
        SELECT 
            incident_investigations.id AS investigation_id,
            incident_register.id AS incident_id,
            incident_register.title AS incident_title,
            team_lead.first_name_en AS investigator_first_name_en,
            team_lead.last_name_en AS investigator_last_name_en,
            team_lead.first_name_ar AS investigator_first_name_ar,
            team_lead.last_name_ar AS investigator_last_name_ar,
            incident_statuses.title_en AS incident_status_title_en,
            incident_statuses.title_ar AS incident_status_title_ar,
            COALESCE(FindingsCTE.findings_summary, '') AS findings_summary,
            COALESCE(RootCauseCTE.root_cause_analysis, '') AS root_cause_analysis,
            COALESCE(InvestigationUsersCTE.investigation_teams_en, '') AS investigation_teams_en,
            COALESCE(InvestigationUsersCTE.investigation_teams_ar, '') AS investigation_teams_ar
        FROM 
            incident_investigations
        INNER JOIN 
            incident_register ON incident_register.id = incident_investigations.incident_register_id
        INNER JOIN 
            incident_statuses ON incident_statuses.id = incident_register.incident_status_id
        LEFT JOIN 
            users AS team_lead ON team_lead.id = incident_investigations.team_lead_id
        LEFT JOIN 
            FindingsCTE ON FindingsCTE.incident_investigation_id = incident_investigations.id
        LEFT JOIN 
            RootCauseCTE ON RootCauseCTE.incident_register_id = incident_register.id
        LEFT JOIN 
            InvestigationUsersCTE ON InvestigationUsersCTE.incident_investigation_id = incident_investigations.id
        WHERE 
            incident_statuses.type IN ('investigating', 'resolved', 'open')
        GROUP BY 
            incident_investigations.id, incident_register.id, incident_register.title, team_lead.first_name_en, team_lead.last_name_en, team_lead.first_name_ar, team_lead.last_name_ar, incident_statuses.title_en, incident_statuses.title_ar, FindingsCTE.findings_summary, RootCauseCTE.root_cause_analysis, InvestigationUsersCTE.investigation_teams_en, InvestigationUsersCTE.investigation_teams_ar
    `)

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("incident_register.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *IncidentReportRepository) GetIncidentRCA(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentRCADTO, error) {
	results := []*reportDTOs.IncidentRCADTO{}

	query := repo.db.Raw(`
        WITH RootCauseCTE AS (
            SELECT
                root_cause.incident_register_id,
                STRING_AGG(CONCAT(root_cause.question, ': ', root_cause.answer), ', ') AS root_cause_identified
            FROM incident_root_cause_analyses AS root_cause
            GROUP BY root_cause.incident_register_id
        )
        SELECT 
            incident_register.id AS incident_id,
            incident_register.title AS incident_title,
            COALESCE(RootCauseCTE.root_cause_identified, '') AS root_cause_identified
        FROM 
            incident_register
        LEFT JOIN 
            RootCauseCTE ON RootCauseCTE.incident_register_id = incident_register.id
        WHERE 
            RootCauseCTE.root_cause_identified IS NOT NULL
        GROUP BY 
            incident_register.id, incident_register.title, RootCauseCTE.root_cause_identified
    `)

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("incident_register.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *IncidentReportRepository) GetIncidentClosure(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentClosureDTO, error) {
	results := []*reportDTOs.IncidentClosureDTO{}

	query := repo.db.Raw(`
        WITH LessonsCTE AS (
            SELECT
                incident_closure_id,
                STRING_AGG(lessons.lessons, ', ') AS associated_action_plans
            FROM incident_closure_lessons AS lessons
            GROUP BY incident_closure_id
        )
        SELECT 
            incident_register.id AS incident_id,
            incident_register.title AS incident_title,
            incident_closures.completion_date AS closure_date,
            COALESCE(incident_closures.resolution_summary, '') AS resolution_summary,
            incident_statuses.title_en AS final_status,
            COALESCE(LessonsCTE.associated_action_plans, '') AS associated_action_plans
        FROM 
            incident_closures
        INNER JOIN 
            incident_register ON incident_register.id = incident_closures.incident_register_id
        INNER JOIN 
            incident_statuses ON incident_statuses.id = incident_register.incident_status_id
        LEFT JOIN 
            LessonsCTE ON LessonsCTE.incident_closure_id = incident_closures.id
        WHERE 
            incident_statuses.type = 'closed'
        GROUP BY 
            incident_register.id, incident_register.title, incident_closures.completion_date, incident_closures.resolution_summary, incident_statuses.title_en, LessonsCTE.associated_action_plans
    `)

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("incident_closures.completion_date BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *IncidentReportRepository) GetIncidentPendingActionPlan(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentPendingActionPlanDTO, error) {
	results := []*reportDTOs.IncidentPendingActionPlanDTO{}

	query := repo.db.Raw(`
        SELECT 
            incident_register.id AS incident_id,
            incident_register.title AS incident_title,
            incident_action_plans.id AS action_plan_id,
            incident_action_plans.title AS action_plan_title,
            COALESCE(incident_action_plans.description, '') AS description,
            incident_action_plans.target_date AS due_date,
            action_plan_statuses.title_en AS status
        FROM 
            incident_action_plans
        INNER JOIN 
            incident_register ON incident_register.id = incident_action_plans.incident_register_id
        INNER JOIN 
            action_plan_statuses ON action_plan_statuses.id = incident_action_plans.action_plan_status_id
        WHERE 
            (action_plan_statuses.is_closed = 0 OR incident_action_plans.target_date < GETDATE())
        AND 
            incident_action_plans.target_date IS NOT NULL
        GROUP BY 
            incident_register.id, incident_register.title, incident_action_plans.id, incident_action_plans.title, 
            incident_action_plans.description, incident_action_plans.target_date, action_plan_statuses.title_en
    `)

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("incident_action_plans.target_date BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *IncidentReportRepository) GetOpenIncidentsWithRisk(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.OpenIncidentWithRiskDTO, error) {
	results := []*reportDTOs.OpenIncidentWithRiskDTO{}

	query := repo.db.Raw(`
        WITH AssociatedRisksCTE AS (
            SELECT 
                incident_risks.incident_register_id,
                STRING_AGG(risk_register.title, ', ') AS associated_risks
            FROM 
                incident_risks
            INNER JOIN 
                risk_register ON risk_register.id = incident_risks.risk_register_id
            GROUP BY 
                incident_risks.incident_register_id
        )
        SELECT 
            incident_register.id AS incident_id,
            incident_register.title AS incident_title,
            COALESCE(AssociatedRisksCTE.associated_risks, '') AS associated_risks,
            severity_levels.title_en AS severity,
            incident_statuses.title_en AS current_status,
            CONCAT(users.first_name_en, ' ', users.last_name_en) AS responsible_user
        FROM 
            incident_register
        INNER JOIN 
            severity_levels ON severity_levels.id = incident_register.severity_level_id
        INNER JOIN 
            incident_statuses ON incident_statuses.id = incident_register.incident_status_id
        LEFT JOIN 
            users ON users.id = incident_register.responsible_user_id
        LEFT JOIN 
            AssociatedRisksCTE ON AssociatedRisksCTE.incident_register_id = incident_register.id
        WHERE 
            incident_statuses.type = 'open' 
            AND COALESCE(AssociatedRisksCTE.associated_risks, '') != ''
        GROUP BY 
            incident_register.id, incident_register.title, AssociatedRisksCTE.associated_risks, 
            severity_levels.title_en, incident_statuses.title_en, users.first_name_en, users.last_name_en
    `)

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("incident_register.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *IncidentReportRepository) GetIncidentActionPlanByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentActionPlanByStatusDTO, error) {
	results := []*reportDTOs.IncidentActionPlanByStatusDTO{}

	query := repo.db.Table("action_plan_statuses").
		Select(`
            action_plan_statuses.id,
            action_plan_statuses.title_en,
            action_plan_statuses.title_ar,
            action_plan_statuses.color_code,
            COUNT(incident_action_plans.id) AS count
        `).
		Joins("LEFT JOIN incident_action_plans ON incident_action_plans.action_plan_status_id = action_plan_statuses.id").
		Group("action_plan_statuses.id, action_plan_statuses.title_en, action_plan_statuses.title_ar, action_plan_statuses.color_code")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("incident_action_plans.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *IncidentReportRepository) GetIncidentActionPlanByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentActionPlanByResponsibleUnitDTO, error) {
	results := []*reportDTOs.IncidentActionPlanByResponsibleUnitDTO{}

	query := repo.db.Table("incident_action_plans").
		Select(`
			organization_units.id AS responsible_unit_id,
			organization_units.title_en AS responsible_unit_title_en,
			organization_units.title_ar AS responsible_unit_title_ar,
			organization_units.color_code AS color_code,
			COUNT(incident_action_plans.id) AS count
        `).
		Joins("JOIN incident_register ON incident_register.id = incident_action_plans.incident_register_id").
		Joins("JOIN organization_units ON organization_units.id = incident_register.responsible_unit_id").
		Group("organization_units.id, organization_units.title_en, organization_units.title_ar, organization_units.color_code")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("incident_action_plans.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *IncidentReportRepository) GetIncidentActionPlanByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.IncidentActionPlanByResponsibleUserDTO, error) {
	results := []*reportDTOs.IncidentActionPlanByResponsibleUserDTO{}

	query := repo.db.Table("incident_action_plans").
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
            COUNT(incident_action_plans.id) AS count
        `).
		Joins("JOIN users ON users.id = incident_action_plans.responsible_user_id").
		Group("users.id, users.first_name_en, users.first_name_ar, users.last_name_en, users.last_name_ar, users.initial, users.color, users.email, users.profile_image_path")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("incident_action_plans.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}
