package repositories

import (
	"math"

	"jethur/grc/kpi/dtos"
	"jethur/grc/models"

	"gorm.io/gorm"
)

type ProjectSecurityGateDashboardRepositoryInterface interface {
	GetSummary(params *dtos.ProjectSecurityGateDashboardQueryParamsDTO) (*dtos.ProjectSecurityGateDashboardSummaryDTO, error)
	GetExceptionSummary(params *dtos.ProjectSecurityGateDashboardQueryParamsDTO) (*dtos.ProjectSecurityGateExceptionSummaryDTO, error)
	GetChecklistsByResponseStatus(params *dtos.ProjectSecurityGateDashboardQueryParamsDTO) (*dtos.ProjectSecurityGateChecklistsByResponseStatusDTO, error)
	GetByReadinessStatus(params *dtos.ProjectSecurityGateDashboardQueryParamsDTO) (*dtos.ProjectSecurityGatesByReadinessStatusDTO, error)
	GetByClassification(params *dtos.ProjectSecurityGateDashboardQueryParamsDTO) (*dtos.ProjectSecurityGatesByClassificationDTO, error)
}

type ProjectSecurityGateDashboardRepository struct {
	db *gorm.DB
}

func NewProjectSecurityGateDashboardRepository(db *gorm.DB) ProjectSecurityGateDashboardRepositoryInterface {
	return &ProjectSecurityGateDashboardRepository{db: db}
}

// applyDateFilter scopes every metric by the security gate creation date so the
// checklist level counts stay aligned with the security gate level ones.
func (repo *ProjectSecurityGateDashboardRepository) applyDateFilter(query *gorm.DB, params *dtos.ProjectSecurityGateDashboardQueryParamsDTO, table string) *gorm.DB {
	if params != nil {
		if params.DateFrom != "" && params.DateTo != "" {
			query = query.Where(table+".created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
		} else if params.Year > 0 {
			query = query.Where("YEAR("+table+".created_at) = ?", params.Year)
		}
	}
	return query
}

func (repo *ProjectSecurityGateDashboardRepository) percentageOf(count int64, total int64) float64 {
	if total == 0 {
		return 0
	}
	return math.Round((float64(count)/float64(total))*100*100) / 100
}

func (repo *ProjectSecurityGateDashboardRepository) checklistQuery(params *dtos.ProjectSecurityGateDashboardQueryParamsDTO) *gorm.DB {
	query := repo.db.Table("project_security_gate_checklists").
		Joins(`
			JOIN project_security_gates
				ON project_security_gates.id = project_security_gate_checklists.security_gate_id
		`)

	return repo.applyDateFilter(query, params, "project_security_gates")
}

func (repo *ProjectSecurityGateDashboardRepository) GetSummary(params *dtos.ProjectSecurityGateDashboardQueryParamsDTO) (*dtos.ProjectSecurityGateDashboardSummaryDTO, error) {
	summary := &dtos.ProjectSecurityGateDashboardSummaryDTO{}

	type gateRow struct {
		TotalSecurityGates     int64 `gorm:"column:total_security_gates"`
		CompletedSecurityGates int64 `gorm:"column:completed_security_gates"`
	}

	gateQuery := repo.db.Table("project_security_gates").
		Select(`
			COUNT(*) AS total_security_gates,
			COALESCE(SUM(CASE WHEN readiness_status.type = 'go_live' THEN 1 ELSE 0 END), 0) AS completed_security_gates
		`).
		Joins(`
			LEFT JOIN project_security_gate_readiness_statuses AS readiness_status
				ON readiness_status.id = project_security_gates.readiness_status_id
		`)
	gateQuery = repo.applyDateFilter(gateQuery, params, "project_security_gates")

	var gates gateRow
	if err := gateQuery.Scan(&gates).Error; err != nil {
		return nil, err
	}

	type checklistRow struct {
		TotalChecklists     int64 `gorm:"column:total_checklists"`
		CompletedChecklists int64 `gorm:"column:completed_checklists"`
		CriticalChecklists  int64 `gorm:"column:critical_checklists"`
		RisksRaised         int64 `gorm:"column:risks_raised"`
	}

	checklistQuery := repo.checklistQuery(params).
		Select(`
			COUNT(*) AS total_checklists,
			COALESCE(SUM(CASE
				WHEN response_status.type IN ('implemented', 'not_applicable')
					OR (project_security_gate_checklists.is_exception = 1 AND project_security_gate_checklists.approved_by_id IS NOT NULL)
				THEN 1 ELSE 0 END), 0) AS completed_checklists,
			COALESCE(SUM(CASE WHEN requirements.is_critical = 1 THEN 1 ELSE 0 END), 0) AS critical_checklists,
			COALESCE(SUM(CASE WHEN risks.checklist_id IS NOT NULL THEN 1 ELSE 0 END), 0) AS risks_raised
		`).
		Joins(`
			JOIN project_security_gate_checklist_template_requirements AS requirements
				ON requirements.id = project_security_gate_checklists.requirement_id
			LEFT JOIN project_security_gate_checklist_response_statuses AS response_status
				ON response_status.id = project_security_gate_checklists.response_status_id
			LEFT JOIN (
				SELECT DISTINCT checklist_id FROM project_security_gate_checklist_risks
			) AS risks ON risks.checklist_id = project_security_gate_checklists.id
		`)

	var checklists checklistRow
	if err := checklistQuery.Scan(&checklists).Error; err != nil {
		return nil, err
	}

	summary.TotalSecurityGates = gates.TotalSecurityGates
	summary.CompletedSecurityGates = gates.CompletedSecurityGates
	summary.CompletedSecurityGatesPercentage = repo.percentageOf(gates.CompletedSecurityGates, gates.TotalSecurityGates)

	summary.TotalChecklists = checklists.TotalChecklists
	summary.CompletedChecklists = checklists.CompletedChecklists
	summary.CompletedChecklistsPercentage = repo.percentageOf(checklists.CompletedChecklists, checklists.TotalChecklists)
	summary.RisksRaised = checklists.RisksRaised
	summary.RisksRaisedPercentage = repo.percentageOf(checklists.RisksRaised, checklists.TotalChecklists)
	summary.CriticalChecklists = checklists.CriticalChecklists
	summary.CriticalChecklistsPercentage = repo.percentageOf(checklists.CriticalChecklists, checklists.TotalChecklists)

	return summary, nil
}

func (repo *ProjectSecurityGateDashboardRepository) GetExceptionSummary(params *dtos.ProjectSecurityGateDashboardQueryParamsDTO) (*dtos.ProjectSecurityGateExceptionSummaryDTO, error) {
	summary := &dtos.ProjectSecurityGateExceptionSummaryDTO{}

	type exceptionRow struct {
		TotalChecklists    int64 `gorm:"column:total_checklists"`
		TotalExceptions    int64 `gorm:"column:total_exceptions"`
		ExceptionsApproved int64 `gorm:"column:exceptions_approved"`
		ExceptionsPending  int64 `gorm:"column:exceptions_pending"`
	}

	query := repo.checklistQuery(params).
		Select(`
			COUNT(*) AS total_checklists,
			COALESCE(SUM(CASE WHEN project_security_gate_checklists.is_exception = 1 THEN 1 ELSE 0 END), 0) AS total_exceptions,
			COALESCE(SUM(CASE WHEN project_security_gate_checklists.is_exception = 1 AND project_security_gate_checklists.approved_by_id IS NOT NULL THEN 1 ELSE 0 END), 0) AS exceptions_approved,
			COALESCE(SUM(CASE WHEN project_security_gate_checklists.is_exception = 1 AND project_security_gate_checklists.approved_by_id IS NULL THEN 1 ELSE 0 END), 0) AS exceptions_pending
		`)

	var row exceptionRow
	if err := query.Scan(&row).Error; err != nil {
		return nil, err
	}

	summary.TotalChecklists = row.TotalChecklists
	summary.TotalExceptions = row.TotalExceptions
	summary.TotalExceptionsPercentage = repo.percentageOf(row.TotalExceptions, row.TotalChecklists)
	summary.ExceptionsApproved = row.ExceptionsApproved
	summary.ExceptionsApprovedPercentage = repo.percentageOf(row.ExceptionsApproved, row.TotalExceptions)
	summary.ExceptionsPending = row.ExceptionsPending
	summary.ExceptionsPendingPercentage = repo.percentageOf(row.ExceptionsPending, row.TotalExceptions)

	return summary, nil
}

func (repo *ProjectSecurityGateDashboardRepository) GetChecklistsByResponseStatus(params *dtos.ProjectSecurityGateDashboardQueryParamsDTO) (*dtos.ProjectSecurityGateChecklistsByResponseStatusDTO, error) {
	breakdown := &dtos.ProjectSecurityGateChecklistsByResponseStatusDTO{
		Statuses: []dtos.ProjectSecurityGateBreakdownItemDTO{},
	}

	// Only answered checklists make up the total so the percentages add up to 100.
	totalQuery := repo.checklistQuery(params).
		Where("project_security_gate_checklists.response_status_id IS NOT NULL")
	if err := totalQuery.Count(&breakdown.TotalChecklists).Error; err != nil {
		return nil, err
	}

	var statuses []models.ProjectSecurityGateChecklistResponseStatus
	if err := repo.db.Where("status = ?", "active").Order("id ASC").Find(&statuses).Error; err != nil {
		return nil, err
	}

	countQuery := repo.checklistQuery(params).
		Select("project_security_gate_checklists.response_status_id AS id, COUNT(*) AS cnt").
		Where("project_security_gate_checklists.response_status_id IS NOT NULL").
		Group("project_security_gate_checklists.response_status_id")

	type countRow struct {
		ID  uint  `gorm:"column:id"`
		Cnt int64 `gorm:"column:cnt"`
	}

	var countRows []countRow
	if err := countQuery.Scan(&countRows).Error; err != nil {
		return nil, err
	}

	countByStatusID := make(map[uint]int64, len(countRows))
	for _, r := range countRows {
		countByStatusID[r.ID] = r.Cnt
	}

	for _, s := range statuses {
		c := countByStatusID[s.ID]
		breakdown.Statuses = append(breakdown.Statuses, dtos.ProjectSecurityGateBreakdownItemDTO{
			ID:         s.ID,
			TitleEn:    s.TitleEn,
			TitleAr:    s.TitleAr,
			Type:       s.Type,
			ColorCode:  s.ColorCode,
			Count:      c,
			Percentage: repo.percentageOf(c, breakdown.TotalChecklists),
		})
	}

	return breakdown, nil
}

func (repo *ProjectSecurityGateDashboardRepository) GetByReadinessStatus(params *dtos.ProjectSecurityGateDashboardQueryParamsDTO) (*dtos.ProjectSecurityGatesByReadinessStatusDTO, error) {
	breakdown := &dtos.ProjectSecurityGatesByReadinessStatusDTO{
		Statuses: []dtos.ProjectSecurityGateBreakdownItemDTO{},
	}

	totalQuery := repo.applyDateFilter(repo.db.Table("project_security_gates"), params, "project_security_gates")
	if err := totalQuery.Count(&breakdown.TotalSecurityGates).Error; err != nil {
		return nil, err
	}

	var statuses []models.ProjectSecurityGateReadinessStatus
	if err := repo.db.Where("status = ?", "active").Order("id ASC").Find(&statuses).Error; err != nil {
		return nil, err
	}

	countQuery := repo.db.Table("project_security_gates").
		Select("project_security_gates.readiness_status_id AS id, COUNT(*) AS cnt").
		Group("project_security_gates.readiness_status_id")
	countQuery = repo.applyDateFilter(countQuery, params, "project_security_gates")

	type countRow struct {
		ID  uint  `gorm:"column:id"`
		Cnt int64 `gorm:"column:cnt"`
	}

	var countRows []countRow
	if err := countQuery.Scan(&countRows).Error; err != nil {
		return nil, err
	}

	countByStatusID := make(map[uint]int64, len(countRows))
	for _, r := range countRows {
		countByStatusID[r.ID] = r.Cnt
	}

	for _, s := range statuses {
		c := countByStatusID[s.ID]
		breakdown.Statuses = append(breakdown.Statuses, dtos.ProjectSecurityGateBreakdownItemDTO{
			ID:         s.ID,
			TitleEn:    s.TitleEn,
			TitleAr:    s.TitleAr,
			Type:       s.Type,
			ColorCode:  s.ColorCode,
			Count:      c,
			Percentage: repo.percentageOf(c, breakdown.TotalSecurityGates),
		})
	}

	return breakdown, nil
}

func (repo *ProjectSecurityGateDashboardRepository) GetByClassification(params *dtos.ProjectSecurityGateDashboardQueryParamsDTO) (*dtos.ProjectSecurityGatesByClassificationDTO, error) {
	breakdown := &dtos.ProjectSecurityGatesByClassificationDTO{
		Classifications: []dtos.ProjectSecurityGateBreakdownItemDTO{},
	}

	totalQuery := repo.applyDateFilter(repo.db.Table("project_security_gates"), params, "project_security_gates")
	if err := totalQuery.Count(&breakdown.TotalSecurityGates).Error; err != nil {
		return nil, err
	}

	var classifications []models.ProjectSecurityGateClassification
	if err := repo.db.Where("status = ?", "active").Order("id ASC").Find(&classifications).Error; err != nil {
		return nil, err
	}

	countQuery := repo.db.Table("project_security_gate_classification_mappings").
		Select("project_security_gate_classification_mappings.classification_id AS id, COUNT(*) AS cnt").
		Joins(`
			JOIN project_security_gates
				ON project_security_gates.id = project_security_gate_classification_mappings.security_gate_id
		`).
		Group("project_security_gate_classification_mappings.classification_id")
	countQuery = repo.applyDateFilter(countQuery, params, "project_security_gates")

	type countRow struct {
		ID  uint  `gorm:"column:id"`
		Cnt int64 `gorm:"column:cnt"`
	}

	var countRows []countRow
	if err := countQuery.Scan(&countRows).Error; err != nil {
		return nil, err
	}

	// A security gate can carry more than one classification, so the percentages
	// are shares of the classifications assigned, not of the security gates.
	countByClassificationID := make(map[uint]int64, len(countRows))
	for _, r := range countRows {
		countByClassificationID[r.ID] = r.Cnt
	}

	for _, cl := range classifications {
		breakdown.TotalClassifications += countByClassificationID[cl.ID]
	}

	for _, cl := range classifications {
		c := countByClassificationID[cl.ID]
		breakdown.Classifications = append(breakdown.Classifications, dtos.ProjectSecurityGateBreakdownItemDTO{
			ID:         cl.ID,
			TitleEn:    cl.TitleEn,
			TitleAr:    cl.TitleAr,
			Type:       cl.Type,
			ColorCode:  cl.ColorCode,
			Count:      c,
			Percentage: repo.percentageOf(c, breakdown.TotalClassifications),
		})
	}

	return breakdown, nil
}
