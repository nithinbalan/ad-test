package repositories

import (
	"fmt"
	"math"
	"time"

	"jethur/grc/kpi/dtos"
	"jethur/grc/models"

	"gorm.io/gorm"
)

type ProjectDashboardRepositoryInterface interface {
	GetProjectSummary(params *dtos.ProjectDashboardQueryParamsDTO) (*dtos.ProjectDashboardSummaryDTO, error)
	GetProjectsByPhase(params *dtos.ProjectDashboardQueryParamsDTO) (*dtos.ProjectsByPhaseDTO, error)
	GetProjectsByClassification(params *dtos.ProjectDashboardQueryParamsDTO) (*dtos.ProjectsByClassificationDTO, error)
	GetProjectsByStatus(params *dtos.ProjectDashboardQueryParamsDTO) (*dtos.ProjectsByStatusDTO, error)
	GetProjectCompletion(params *dtos.ProjectDashboardQueryParamsDTO) ([]dtos.ProjectCompletionListDTO, error)
	GetProjectsDueSoonAndOverdue(params *dtos.ProjectDashboardQueryParamsDTO) ([]dtos.ProjectDueSoonOverdueListDTO, error)
}

type ProjectDashboardRepository struct {
	db *gorm.DB
}

func NewProjectDashboardRepository(db *gorm.DB) ProjectDashboardRepositoryInterface {
	return &ProjectDashboardRepository{db: db}
}

type phaseCountRow struct {
	ProjectPhaseID uint  `gorm:"column:project_phase_id"`
	Cnt            int64 `gorm:"column:cnt"`
}

type classificationCountRow struct {
	ProjectClassificationID uint  `gorm:"column:project_classification_id"`
	Cnt                     int64 `gorm:"column:cnt"`
}

type statusCountRow struct {
	ProjectStatusID uint  `gorm:"column:project_status_id"`
	Cnt             int64 `gorm:"column:cnt"`
}

func (repo *ProjectDashboardRepository) applyDateFilter(query *gorm.DB, params *dtos.ProjectDashboardQueryParamsDTO, table string) *gorm.DB {
	if params != nil {
		if params.DateFrom != "" && params.DateTo != "" {
			query = query.Where(table+".created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
		} else if params.Year > 0 {
			query = query.Where("YEAR("+table+".created_at) = ?", params.Year)
		}
	}
	return query
}

func (repo *ProjectDashboardRepository) GetProjectSummary(params *dtos.ProjectDashboardQueryParamsDTO) (*dtos.ProjectDashboardSummaryDTO, error) {
	summary := &dtos.ProjectDashboardSummaryDTO{}

	query := repo.db.Table("projects p").
		Select(`
			COUNT(*) AS total_projects,
			COUNT(CASE WHEN ws.type = 'approve' AND (ph.type IS NULL OR (ph.type != 'go_live' AND ph.type != 'closure')) THEN 1 END) AS active_projects,
			COUNT(CASE WHEN ph.type IN ('go_live', 'closure') THEN 1 END) AS completed_projects,
			COUNT(CASE WHEN ws.type = 'approve' AND (ph.type IS NULL OR (ph.type != 'go_live' AND ph.type != 'closure')) AND ps.type = 'on_track' THEN 1 END) AS on_track,
			COUNT(CASE WHEN ws.type = 'approve' AND (ph.type IS NULL OR (ph.type != 'go_live' AND ph.type != 'closure')) AND ps.type = 'at_risk' THEN 1 END) AS at_risk,
			COUNT(CASE WHEN ws.type = 'approve' AND (ph.type IS NULL OR (ph.type != 'go_live' AND ph.type != 'closure')) AND ps.type = 'delayed' THEN 1 END) AS delayed,
			COUNT(CASE WHEN ws.type = 'approve' AND (ph.type IS NULL OR (ph.type != 'go_live' AND ph.type != 'closure')) AND p.end_date >= CAST(GETDATE() AS DATE) AND p.end_date <= DATEADD(day, 30, CAST(GETDATE() AS DATE)) THEN 1 END) AS projects_due_soon,
			COUNT(CASE WHEN ws.type = 'approve' AND (ph.type IS NULL OR (ph.type != 'go_live' AND ph.type != 'closure')) AND p.end_date < CAST(GETDATE() AS DATE) THEN 1 END) AS overdue_projects
		`).
		Joins("LEFT JOIN project_phases ph ON ph.id = p.project_phase_id").
		Joins("LEFT JOIN workflow_statuses ws ON ws.id = p.workflow_status_id").
		Joins("LEFT JOIN project_statuses ps ON ps.id = p.project_status_id")

	query = repo.applyDateFilter(query, params, "p")

	if err := query.Scan(summary).Error; err != nil {
		return nil, err
	}

	// Open issues (unresolved action plans across active projects)
	issuesQuery := repo.db.Table("project_action_plans").
		Joins("JOIN projects p ON p.id = project_action_plans.project_id").
		Joins("LEFT JOIN project_phases ph ON ph.id = p.project_phase_id").
		Joins("LEFT JOIN workflow_statuses ws ON ws.id = p.workflow_status_id").
		Joins("JOIN action_plan_statuses ON action_plan_statuses.id = project_action_plans.action_plan_status_id").
		Where("action_plan_statuses.is_closed = 0 AND action_plan_statuses.type NOT IN ('completed', 'closed')").
		Where("ws.type = 'approve' AND (ph.type IS NULL OR (ph.type != 'go_live' AND ph.type != 'closure'))")

	issuesQuery = repo.applyDateFilter(issuesQuery, params, "project_action_plans")

	if err := issuesQuery.Count(&summary.OpenIssues).Error; err != nil {
		return nil, err
	}

	// Open change requests (pending review or draft on active projects)
	crQuery := repo.db.Table("project_change_requests").
		Joins("JOIN projects p ON p.id = project_change_requests.project_id").
		Joins("LEFT JOIN project_phases ph ON ph.id = p.project_phase_id").
		Joins("LEFT JOIN workflow_statuses ws_proj ON ws_proj.id = p.workflow_status_id").
		Joins("LEFT JOIN workflow_statuses ON workflow_statuses.id = project_change_requests.workflow_status_id").
		Where("project_change_requests.workflow_status_id IS NULL OR workflow_statuses.type NOT IN ('approve', 'reject')").
		Where("ws_proj.type = 'approve' AND (ph.type IS NULL OR (ph.type != 'go_live' AND ph.type != 'closure'))")

	crQuery = repo.applyDateFilter(crQuery, params, "project_change_requests")

	if err := crQuery.Count(&summary.OpenChangeRequests).Error; err != nil {
		return nil, err
	}

	// Calculate completion percentage: projects in go live against all projects
	if summary.TotalProjects > 0 {
		summary.CompletionPercentage = math.Round((float64(summary.CompletedProjects)/float64(summary.TotalProjects))*100*100) / 100
		summary.AvgCompletion = summary.CompletionPercentage
	}

	// Calculate percentages against active projects
	if summary.ActiveProjects > 0 {
		summary.OpenChangeRequestsPercentage = math.Round((float64(summary.OpenChangeRequests)/float64(summary.ActiveProjects))*100*100) / 100
		summary.ProjectsDueSoonPercentage = math.Round((float64(summary.ProjectsDueSoon)/float64(summary.ActiveProjects))*100*100) / 100
		summary.OverdueProjectsPercentage = math.Round((float64(summary.OverdueProjects)/float64(summary.ActiveProjects))*100*100) / 100
	}

	return summary, nil
}

func (repo *ProjectDashboardRepository) GetProjectsByPhase(params *dtos.ProjectDashboardQueryParamsDTO) (*dtos.ProjectsByPhaseDTO, error) {
	summary := &dtos.ProjectsByPhaseDTO{
		Phases: []dtos.ProjectsByPhaseItemDTO{},
	}

	totalQuery := repo.db.Table("projects")
	totalQuery = repo.applyDateFilter(totalQuery, params, "projects")
	if err := totalQuery.Count(&summary.Total).Error; err != nil {
		return nil, err
	}

	var phases []models.ProjectPhase
	if err := repo.db.Where("status = ?", "active").Order("id ASC").Find(&phases).Error; err != nil {
		return nil, err
	}

	countQuery := repo.db.Table("projects").
		Select("projects.project_phase_id, COUNT(*) AS cnt").
		Where("projects.project_phase_id IS NOT NULL").
		Group("projects.project_phase_id")
	countQuery = repo.applyDateFilter(countQuery, params, "projects")

	var countRows []phaseCountRow
	if err := countQuery.Scan(&countRows).Error; err != nil {
		return nil, err
	}

	countByPhaseID := make(map[uint]int64, len(countRows))
	for _, r := range countRows {
		countByPhaseID[r.ProjectPhaseID] = r.Cnt
	}

	percentage := func(count int64) float64 {
		if summary.Total == 0 {
			return 0
		}
		return math.Round((float64(count)/float64(summary.Total))*100*100) / 100
	}

	for _, p := range phases {
		c := countByPhaseID[p.ID]
		summary.Phases = append(summary.Phases, dtos.ProjectsByPhaseItemDTO{
			ID:         p.ID,
			TitleEn:    p.TitleEn,
			TitleAr:    p.TitleAr,
			Type:       p.Type,
			Count:      c,
			Percentage: percentage(c),
		})
	}

	return summary, nil
}

func (repo *ProjectDashboardRepository) GetProjectsByClassification(params *dtos.ProjectDashboardQueryParamsDTO) (*dtos.ProjectsByClassificationDTO, error) {
	summary := &dtos.ProjectsByClassificationDTO{
		Classifications: []dtos.ProjectsByClassificationItemDTO{},
	}

	totalQuery := repo.db.Table("projects")
	totalQuery = repo.applyDateFilter(totalQuery, params, "projects")
	if err := totalQuery.Count(&summary.Total).Error; err != nil {
		return nil, err
	}

	var classifications []models.ProjectClassification
	if err := repo.db.Where("status = ?", "active").Order("id ASC").Find(&classifications).Error; err != nil {
		return nil, err
	}

	countQuery := repo.db.Table("projects").
		Select("projects.project_classification_id, COUNT(*) AS cnt").
		Where("projects.project_classification_id IS NOT NULL").
		Group("projects.project_classification_id")
	countQuery = repo.applyDateFilter(countQuery, params, "projects")

	var countRows []classificationCountRow
	if err := countQuery.Scan(&countRows).Error; err != nil {
		return nil, err
	}

	countByClassificationID := make(map[uint]int64, len(countRows))
	for _, r := range countRows {
		countByClassificationID[r.ProjectClassificationID] = r.Cnt
	}

	percentage := func(count int64) float64 {
		if summary.Total == 0 {
			return 0
		}
		return math.Round((float64(count)/float64(summary.Total))*100*100) / 100
	}

	for _, cl := range classifications {
		c := countByClassificationID[cl.ID]
		summary.Classifications = append(summary.Classifications, dtos.ProjectsByClassificationItemDTO{
			ID:         cl.ID,
			TitleEn:    cl.TitleEn,
			TitleAr:    cl.TitleAr,
			Count:      c,
			Percentage: percentage(c),
		})
	}

	return summary, nil
}

func (repo *ProjectDashboardRepository) GetProjectsByStatus(params *dtos.ProjectDashboardQueryParamsDTO) (*dtos.ProjectsByStatusDTO, error) {
	summary := &dtos.ProjectsByStatusDTO{
		Statuses: []dtos.ProjectsByStatusItemDTO{},
	}

	totalQuery := repo.db.Table("projects")
	totalQuery = repo.applyDateFilter(totalQuery, params, "projects")
	if err := totalQuery.Count(&summary.Total).Error; err != nil {
		return nil, err
	}

	var statuses []models.ProjectStatus
	if err := repo.db.Where("status = ?", "active").Order("id ASC").Find(&statuses).Error; err != nil {
		return nil, err
	}

	countQuery := repo.db.Table("projects").
		Select("projects.project_status_id, COUNT(*) AS cnt").
		Where("projects.project_status_id IS NOT NULL").
		Group("projects.project_status_id")
	countQuery = repo.applyDateFilter(countQuery, params, "projects")

	var countRows []statusCountRow
	if err := countQuery.Scan(&countRows).Error; err != nil {
		return nil, err
	}

	countByStatusID := make(map[uint]int64, len(countRows))
	for _, r := range countRows {
		countByStatusID[r.ProjectStatusID] = r.Cnt
	}

	percentage := func(count int64) float64 {
		if summary.Total == 0 {
			return 0
		}
		return math.Round((float64(count)/float64(summary.Total))*100*100) / 100
	}

	for _, s := range statuses {
		c := countByStatusID[s.ID]
		summary.Statuses = append(summary.Statuses, dtos.ProjectsByStatusItemDTO{
			ID:         s.ID,
			TitleEn:    s.TitleEn,
			TitleAr:    s.TitleAr,
			Type:       s.Type,
			ColorCode:  s.ColorCode,
			Count:      c,
			Percentage: percentage(c),
		})
	}

	return summary, nil
}

func (repo *ProjectDashboardRepository) GetProjectCompletion(params *dtos.ProjectDashboardQueryParamsDTO) ([]dtos.ProjectCompletionListDTO, error) {
	limit := 10
	if params != nil && params.Limit > 0 {
		limit = params.Limit
		if limit > 100 {
			limit = 100
		}
	}

	type completionRow struct {
		ID                    uint      `gorm:"column:id"`
		TitleEn               string    `gorm:"column:title_en"`
		TitleAr               string    `gorm:"column:title_ar"`
		PhaseID               *uint     `gorm:"column:phase_id"`
		PhaseTitleEn          *string   `gorm:"column:phase_title_en"`
		PhaseTitleAr          *string   `gorm:"column:phase_title_ar"`
		ClassificationID      *uint     `gorm:"column:classification_id"`
		ClassificationTitleEn *string   `gorm:"column:classification_title_en"`
		ClassificationTitleAr *string   `gorm:"column:classification_title_ar"`
		StartDate             time.Time `gorm:"column:start_date"`
		EndDate               time.Time `gorm:"column:end_date"`
		CompletionPercentage  *float64  `gorm:"column:completion_percentage"`
		StatusID              uint      `gorm:"column:status_id"`
		StatusTitleEn         string    `gorm:"column:status_title_en"`
		StatusTitleAr         string    `gorm:"column:status_title_ar"`
		StatusType            string    `gorm:"column:status_type"`
		StatusColorCode       string    `gorm:"column:status_color_code"`
	}

	query := repo.db.Table("projects p").
		Select(`
			p.id,
			p.title_en,
			p.title_ar,
			p.project_phase_id AS phase_id,
			ph.title_en AS phase_title_en,
			ph.title_ar AS phase_title_ar,
			p.project_classification_id AS classification_id,
			pc.title_en AS classification_title_en,
			pc.title_ar AS classification_title_ar,
			p.start_date,
			p.end_date,
			p.completion_percentage,
			ps.id AS status_id,
			ps.title_en AS status_title_en,
			ps.title_ar AS status_title_ar,
			ps.type AS status_type,
			ps.color_code AS status_color_code
		`).
		Joins("LEFT JOIN project_phases ph ON ph.id = p.project_phase_id").
		Joins("LEFT JOIN workflow_statuses ws ON ws.id = p.workflow_status_id").
		Joins("LEFT JOIN project_classifications pc ON pc.id = p.project_classification_id").
		Joins("JOIN project_statuses ps ON ps.id = p.project_status_id").
		Where("ws.type = 'approve' AND (ph.type IS NULL OR (ph.type != 'go_live' AND ph.type != 'closure'))")

	query = repo.applyDateFilter(query, params, "p")

	var rows []completionRow
	if err := query.Order("ISNULL(p.completion_percentage, 0) DESC, p.id DESC").Limit(limit).Scan(&rows).Error; err != nil {
		return nil, err
	}

	items := make([]dtos.ProjectCompletionListDTO, 0, len(rows))
	for _, r := range rows {
		comp := 0.0
		if r.CompletionPercentage != nil {
			comp = *r.CompletionPercentage
		}
		var startDate, startDateFormatted, endDate, endDateFormatted string
		if !r.StartDate.IsZero() {
			startDate = r.StartDate.Format("2006-01-02")
			startDateFormatted = r.StartDate.Format("02 Jan 2006")
		}
		if !r.EndDate.IsZero() {
			endDate = r.EndDate.Format("2006-01-02")
			endDateFormatted = r.EndDate.Format("02 Jan 2006")
		}
		items = append(items, dtos.ProjectCompletionListDTO{
			ID:                           r.ID,
			ReferenceNumber:              fmt.Sprintf("PRJ-%03d", r.ID),
			TitleEn:                      r.TitleEn,
			TitleAr:                      r.TitleAr,
			ProjectPhaseID:               r.PhaseID,
			ProjectPhaseTitleEn:          r.PhaseTitleEn,
			ProjectPhaseTitleAr:          r.PhaseTitleAr,
			ProjectClassificationID:      r.ClassificationID,
			ProjectClassificationTitleEn: r.ClassificationTitleEn,
			ProjectClassificationTitleAr: r.ClassificationTitleAr,
			StartDate:                    startDate,
			StartDateFormatted:           startDateFormatted,
			EndDate:                      endDate,
			EndDateFormatted:             endDateFormatted,
			CompletionPercentage:         comp,
			ProjectStatusID:              r.StatusID,
			ProjectStatusTitleEn:         r.StatusTitleEn,
			ProjectStatusTitleAr:         r.StatusTitleAr,
			ProjectStatusType:            r.StatusType,
			ProjectStatusColorCode:       r.StatusColorCode,
		})
	}

	return items, nil
}

func (repo *ProjectDashboardRepository) GetProjectsDueSoonAndOverdue(params *dtos.ProjectDashboardQueryParamsDTO) ([]dtos.ProjectDueSoonOverdueListDTO, error) {
	limit := 10
	if params != nil && params.Limit > 0 {
		limit = params.Limit
		if limit > 100 {
			limit = 100
		}
	}

	type scheduleRow struct {
		ID                    uint      `gorm:"column:id"`
		TitleEn               string    `gorm:"column:title_en"`
		TitleAr               string    `gorm:"column:title_ar"`
		PhaseID               *uint     `gorm:"column:phase_id"`
		PhaseTitleEn          *string   `gorm:"column:phase_title_en"`
		PhaseTitleAr          *string   `gorm:"column:phase_title_ar"`
		PhaseType             *string   `gorm:"column:phase_type"`
		ClassificationID      *uint     `gorm:"column:classification_id"`
		ClassificationTitleEn *string   `gorm:"column:classification_title_en"`
		ClassificationTitleAr *string   `gorm:"column:classification_title_ar"`
		StartDate             time.Time `gorm:"column:start_date"`
		EndDate               time.Time `gorm:"column:end_date"`
		StatusID              uint      `gorm:"column:status_id"`
		StatusTitleEn         string    `gorm:"column:status_title_en"`
		StatusTitleAr         string    `gorm:"column:status_title_ar"`
		StatusType            string    `gorm:"column:status_type"`
		StatusColorCode       string    `gorm:"column:status_color_code"`
	}

	query := repo.db.Table("projects p").
		Select(`
			p.id,
			p.title_en,
			p.title_ar,
			p.project_phase_id AS phase_id,
			ph.title_en AS phase_title_en,
			ph.title_ar AS phase_title_ar,
			ph.type AS phase_type,
			p.project_classification_id AS classification_id,
			pc.title_en AS classification_title_en,
			pc.title_ar AS classification_title_ar,
			p.start_date,
			p.end_date,
			ps.id AS status_id,
			ps.title_en AS status_title_en,
			ps.title_ar AS status_title_ar,
			ps.type AS status_type,
			ps.color_code AS status_color_code
		`).
		Joins("LEFT JOIN project_phases ph ON ph.id = p.project_phase_id").
		Joins("LEFT JOIN workflow_statuses ws ON ws.id = p.workflow_status_id").
		Joins("LEFT JOIN project_classifications pc ON pc.id = p.project_classification_id").
		Joins("JOIN project_statuses ps ON ps.id = p.project_status_id").
		Where("ws.type = 'approve' AND (ph.type IS NULL OR (ph.type != 'go_live' AND ph.type != 'closure'))").
		Where("p.end_date < CAST(GETDATE() AS DATE) OR (p.end_date >= CAST(GETDATE() AS DATE) AND p.end_date <= DATEADD(day, 30, CAST(GETDATE() AS DATE)))")

	query = repo.applyDateFilter(query, params, "p")

	var rows []scheduleRow
	orderClause := `
		CASE WHEN p.end_date < CAST(GETDATE() AS DATE) THEN 0 ELSE 1 END ASC,
		p.end_date ASC
	`
	if err := query.Order(orderClause).Limit(limit).Scan(&rows).Error; err != nil {
		return nil, err
	}

	today := time.Now().Truncate(24 * time.Hour)

	items := make([]dtos.ProjectDueSoonOverdueListDTO, 0, len(rows))
	for _, r := range rows {
		endDate := r.EndDate.Truncate(24 * time.Hour)
		daysDiff := int(math.Round(endDate.Sub(today).Hours() / 24))

		isOverdue := daysDiff < 0
		absDays := daysDiff
		if absDays < 0 {
			absDays = -absDays
		}

		var overdueDays, daysRemaining int
		var statusTitleEn, statusTitleAr, statusType, statusColorCode string

		if isOverdue {
			overdueDays = absDays
			daysRemaining = 0
			statusTitleEn = fmt.Sprintf("%dd overdue", absDays)
			statusTitleAr = fmt.Sprintf("متأخر %d يوم", absDays)
			statusType = "overdue"
			statusColorCode = "#DC3545" // red
		} else {
			overdueDays = 0
			daysRemaining = absDays
			if daysDiff == 0 {
				statusTitleEn = "Due today"
				statusTitleAr = "مستحق اليوم"
			} else {
				statusTitleEn = fmt.Sprintf("%dd remaining", absDays)
				statusTitleAr = fmt.Sprintf("متبقي %d يوم", absDays)
			}
			statusType = "due_soon"
			statusColorCode = "#FFC107" // yellow
		}

		var startDateStr, startDateFormatted, endDateStr, endDateFormatted string
		if !r.StartDate.IsZero() {
			startDateStr = r.StartDate.Format("2006-01-02")
			startDateFormatted = r.StartDate.Format("02 Jan 2006")
		}
		if !r.EndDate.IsZero() {
			endDateStr = r.EndDate.Format("2006-01-02")
			endDateFormatted = r.EndDate.Format("02 Jan 2006")
		}

		items = append(items, dtos.ProjectDueSoonOverdueListDTO{
			ID:                           r.ID,
			ReferenceNumber:              fmt.Sprintf("PRJ-%03d", r.ID),
			TitleEn:                      r.TitleEn,
			TitleAr:                      r.TitleAr,
			ProjectPhaseID:               r.PhaseID,
			ProjectPhaseTitleEn:          r.PhaseTitleEn,
			ProjectPhaseTitleAr:          r.PhaseTitleAr,
			ProjectPhaseType:             r.PhaseType,
			ProjectClassificationID:      r.ClassificationID,
			ProjectClassificationTitleEn: r.ClassificationTitleEn,
			ProjectClassificationTitleAr: r.ClassificationTitleAr,
			StartDate:                    startDateStr,
			StartDateFormatted:           startDateFormatted,
			EndDate:                      endDateStr,
			EndDateFormatted:             endDateFormatted,
			PlannedEndDate:               endDateStr,
			PlannedEndDateFormatted:      endDateFormatted,
			IsOverdue:                    isOverdue,
			OverdueDays:                  overdueDays,
			RemainingDays:                daysRemaining,
			DaysRemaining:                daysRemaining,
			ProjectStatusID:              r.StatusID,
			ProjectStatusTitleEn:         r.StatusTitleEn,
			ProjectStatusTitleAr:         r.StatusTitleAr,
			ProjectStatusType:            r.StatusType,
			ProjectStatusColorCode:       r.StatusColorCode,
			StatusTitleEn:                statusTitleEn,
			StatusTitleAr:                statusTitleAr,
			StatusType:                   statusType,
			StatusColorCode:              statusColorCode,
		})
	}

	return items, nil
}
