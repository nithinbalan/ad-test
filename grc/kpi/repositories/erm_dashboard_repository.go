package repositories

import (
	"fmt"
	kpiDTOs "jethur/grc/kpi/dtos"
	"jethur/grc/models"
	"math"
	"time"

	"gorm.io/gorm"
)

type ERMDashboardRepositoryInterface interface {
	GetERMRiskOverview(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskOverviewDTO, error)
	GetERMRiskDistribution(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskDistributionDTO, error)
	GetERMOverdueRisks(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMOverdueRisksDTO, error)
	GetERMRisksAboveAppetite(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRisksAboveAppetiteDTO, error)
	GetERMLeadershipAttentionMetrics(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMLeadershipAttentionDTO, error)
	GetERMControlEfficiencySummary(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMControlEfficiencySummaryDTO, error)
	GetERMControlEfficiencyBreakdown(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMControlEfficiencyBreakdownDTO, error)
	GetERMControlsOverview(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMControlsOverviewDTO, error)
	GetERMRiskAppetiteOverview(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskAppetiteOverviewDTO, error)
	GetERMRiskAppetiteSummary(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskAppetiteSummarytDTO, error)
	GetERMRiskAgingStatus(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskAgingStatusDTO, error)
	GetERMRiskTreatmentPlanProgress(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskTreatmentPlanProgressDTO, error)
	GetERMStrategicObjectivesRiskOverview(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMStrategicObjectiveRiskOverviewDTO, error)
	GetERMStrategicObjectivesRiskSummary(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMStrategicObjectivesRiskSummaryDTO, error)
}

type ERMDashboardRepository struct {
	db *gorm.DB
}

func NewERMDashboardRepository(db *gorm.DB) *ERMDashboardRepository {
	return &ERMDashboardRepository{
		db: db,
	}
}

func (repo *ERMDashboardRepository) GetERMRiskOverview(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskOverviewDTO, error) {

	var totalRisks, openRisks, closedRisks, unresolvedHighRisks int64

	row := repo.db.Table("risk_register").
		Joins("JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id").
		Where("risk_statuses.type IN ?", []string{"open", "closed"}).
		Select(`
        SUM(CASE WHEN risk_statuses.type = 'open' THEN 1 ELSE 0 END) AS open_risks,
        SUM(CASE WHEN risk_statuses.type = 'closed' THEN 1 ELSE 0 END) AS closed_risks,
        COUNT(*) AS total_risks
    `)

	if params.DateFrom != "" && params.DateTo != "" {
		row = row.Where(
			"risk_register.reported_at BETWEEN ? AND ?",
			params.DateFrom+" 00:00:00",
			params.DateTo+" 23:59:59",
		)
	}

	err := row.Row().Scan(&openRisks, &closedRisks, &totalRisks)
	if err != nil {
		return nil, err
	}

	// Determine the highest available risk rating dynamically
	var highestRating string
	err = repo.db.Table("risk_matrix_risk_ratings").
		Select("type").
		Where("status = ?", "active").
		Order("CASE type WHEN 'very_high' THEN 1 WHEN 'high' THEN 2 WHEN 'medium' THEN 3 WHEN 'low' THEN 4 WHEN 'very_low' THEN 5 END").
		Limit(1).
		Pluck("type", &highestRating).Error
	if err != nil || highestRating == "" {
		return nil, fmt.Errorf("failed to fetch highest risk rating: %v", err)
	}

	// Unresolved high risks
	highRiskQuery := repo.db.Model(&models.RiskRegister{}).
		Joins("JOIN risk_matrix_risk_ratings ON risk_matrix_risk_ratings.id = risk_register.current_risk_rating_id").
		Joins("JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id").
		Where("risk_matrix_risk_ratings.type = ?", highestRating).
		Where("risk_statuses.type = ?", "open")
	if params.DateFrom != "" && params.DateTo != "" {
		highRiskQuery = highRiskQuery.Where("risk_register.reported_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}
	err = highRiskQuery.Count(&unresolvedHighRisks).Error
	if err != nil {
		return nil, err
	}

	return &kpiDTOs.ERMRiskOverviewDTO{
		TotalRisks:          totalRisks,
		OpenRisks:           openRisks,
		ClosedRisks:         closedRisks,
		UnresolvedHighRisks: unresolvedHighRisks,
	}, nil
}

func (repo *ERMDashboardRepository) GetERMRiskDistribution(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskDistributionDTO, error) {
	summary := &kpiDTOs.ERMRiskDistributionDTO{}

	// Query to calculate total risks registered
	riskQuery := repo.db.Model(&models.RiskRegister{}).
		Joins("JOIN risk_matrix_risk_ratings ON risk_register.current_risk_rating_id = risk_matrix_risk_ratings.id").
		Joins("JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id").
		Where("risk_statuses.type IN ?", []string{"open", "closed"})

	if params.DateFrom != "" && params.DateTo != "" {
		riskQuery = riskQuery.Where("risk_register.reported_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	var totalRisksRegistered int64
	err := riskQuery.Count(&totalRisksRegistered).Error
	if err != nil {
		return nil, err
	}

	// Fetch active risk ratings with details
	var riskRatings []models.RiskMatrixRiskRating
	err = repo.db.Model(&models.RiskMatrixRiskRating{}).
		Where("status = ?", "active").
		Find(&riskRatings).Error
	if err != nil {
		return nil, err
	}

	riskCountsByRating := make([]kpiDTOs.ERMRiskRatingCountDTO, 0)

	// Loop over risk ratings and count associated risks
	for _, rating := range riskRatings {
		var count, openCount int64

		// clone riskQuery for this iteration
		ratingQuery := riskQuery.Session(&gorm.Session{})

		// add rating-specific condition
		ratingQuery = ratingQuery.Where("risk_register.current_risk_rating_id = ?", rating.ID)

		// total count
		if err := ratingQuery.Count(&count).Error; err != nil {
			return nil, err
		}

		// open count
		if err := ratingQuery.
			Where("risk_statuses.type = ?", "open").
			Count(&openCount).Error; err != nil {
			return nil, err
		}

		var percentage float64
		if totalRisksRegistered > 0 {
			percentage = math.Round(float64(count) / float64(totalRisksRegistered) * 100)
		}

		riskCountsByRating = append(riskCountsByRating, kpiDTOs.ERMRiskRatingCountDTO{
			TitleEn:    rating.TitleEn,
			TitleAr:    rating.TitleAr,
			ColorCode:  rating.ColorCode,
			Type:       rating.Type,
			Count:      count,
			OpenCount:  openCount,
			Percentage: percentage,
		})
	}

	summary.RiskCountsByRating = riskCountsByRating

	return summary, nil
}

func (repo *ERMDashboardRepository) GetERMOverdueRisks(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMOverdueRisksDTO, error) {
	currentTime := time.Now()

	// Base query with rating details
	baseQuery := repo.db.Model(&models.RiskRegister{}).
		Select(`
			risk_matrix_risk_ratings.id,
			risk_matrix_risk_ratings.type,
			risk_matrix_risk_ratings.title_en,
			risk_matrix_risk_ratings.title_ar,
			risk_matrix_risk_ratings.color_code,
			COUNT(*) AS count
		`).
		Joins("JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id").
		Joins("JOIN risk_matrix_risk_ratings ON risk_matrix_risk_ratings.id = risk_register.current_risk_rating_id").
		Where("next_review_date < ?", currentTime).
		Where("next_review_date IS NOT NULL").
		Where("risk_statuses.type = ?", "open").
		Group(`
			risk_matrix_risk_ratings.id,
			risk_matrix_risk_ratings.type,
			risk_matrix_risk_ratings.title_en,
			risk_matrix_risk_ratings.title_ar,
			risk_matrix_risk_ratings.color_code
		`)

	if params.DateFrom != "" && params.DateTo != "" {
		baseQuery = baseQuery.Where(
			"reported_at BETWEEN ? AND ?",
			params.DateFrom+" 00:00:00",
			params.DateTo+" 23:59:59",
		)
	}

	type result struct {
		ID        uint
		Type      string
		TitleEn   string
		TitleAr   string
		ColorCode string
		Count     int64
	}

	var results []result
	if err := baseQuery.Scan(&results).Error; err != nil {
		return nil, err
	}

	// Build the overdue details + total
	var totalRisksOverdue int64
	risksOverdueDetails := []kpiDTOs.ERMOverdueRiskDetailDTO{}

	for _, r := range results {
		risksOverdueDetails = append(risksOverdueDetails, kpiDTOs.ERMOverdueRiskDetailDTO{
			ID:        r.ID,
			Type:      r.Type,
			TitleEn:   r.TitleEn,
			TitleAr:   r.TitleAr,
			ColorCode: r.ColorCode,
			Count:     r.Count,
		})
		totalRisksOverdue += r.Count
	}

	return &kpiDTOs.ERMOverdueRisksDTO{
		TotalRisksOverdue:   totalRisksOverdue,
		OverdueByRiskLevels: risksOverdueDetails,
	}, nil
}

func (repo *ERMDashboardRepository) GetERMRisksAboveAppetite(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRisksAboveAppetiteDTO, error) {
	// Base query with risk levels
	baseQuery := repo.db.Model(&models.RiskRegister{}).
		Select(`
			risk_matrix_risk_ratings.id,
			risk_matrix_risk_ratings.type,
			risk_matrix_risk_ratings.title_en,
			risk_matrix_risk_ratings.title_ar,
			risk_matrix_risk_ratings.color_code,
			COUNT(*) AS count
		`).
		Joins("JOIN risk_evaluations ON risk_evaluations.risk_register_id = risk_register.id").
		Joins("JOIN risk_appetite_alignments ON risk_appetite_alignments.risk_evaluation_id = risk_evaluations.id").
		Joins("JOIN risk_appetite_tolerance_levels ON risk_appetite_tolerance_levels.id = risk_appetite_alignments.risk_appetite_tolerance_level_id").
		Joins("LEFT JOIN risk_matrix_risk_ratings ON risk_matrix_risk_ratings.id = risk_register.current_risk_rating_id").
		Where("risk_appetite_tolerance_levels.type = ?", "above").
		Group(`
			risk_matrix_risk_ratings.id,
			risk_matrix_risk_ratings.type,
			risk_matrix_risk_ratings.title_en,
			risk_matrix_risk_ratings.title_ar,
			risk_matrix_risk_ratings.color_code
		`)

	// Apply date filter if given
	if params.DateFrom != "" && params.DateTo != "" {
		baseQuery = baseQuery.Where(
			"risk_register.reported_at BETWEEN ? AND ?",
			params.DateFrom+" 00:00:00",
			params.DateTo+" 23:59:59",
		)
	}

	type result struct {
		ID        uint
		Type      string
		TitleEn   string
		TitleAr   string
		ColorCode string
		Count     int64
	}

	var results []result
	if err := baseQuery.Scan(&results).Error; err != nil {
		return nil, err
	}

	var totalRisksAboveAppetite int64
	risksByLevels := []kpiDTOs.ERMRisksAboveAppetiteLevelDTO{}

	for _, r := range results {
		risksByLevels = append(risksByLevels, kpiDTOs.ERMRisksAboveAppetiteLevelDTO{
			ID:        r.ID,
			Type:      r.Type,
			TitleEn:   r.TitleEn,
			TitleAr:   r.TitleAr,
			ColorCode: r.ColorCode,
			Count:     r.Count,
		})
		totalRisksAboveAppetite += r.Count
	}

	return &kpiDTOs.ERMRisksAboveAppetiteDTO{
		TotalRisksAboveAppetite:    totalRisksAboveAppetite,
		AboveAppetiteRisksByLevels: risksByLevels,
	}, nil
}

func (repo *ERMDashboardRepository) GetERMLeadershipAttentionMetrics(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMLeadershipAttentionDTO, error) {
	summary := &kpiDTOs.ERMLeadershipAttentionDTO{}

	// Base query function
	buildBaseQuery := func() *gorm.DB {
		query := repo.db.Model(&models.RiskRegister{}).
			Joins("JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id").
			Where("risk_statuses.type = ?", "open")

		if params.DateFrom != "" && params.DateTo != "" {
			query = query.Where("risk_register.reported_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
		}
		return query
	}

	// Strategic objectives at risk
	strategicQuery := buildBaseQuery().
		Joins("JOIN risk_strategic_objectives ON risk_strategic_objectives.risk_register_id = risk_register.id")
	if err := strategicQuery.Count(&summary.StrategicObjectivesAtRisk).Error; err != nil {
		return nil, err
	}

	// Critical risks uncontrolled
	criticalQuery := buildBaseQuery().
		Joins("JOIN risk_matrix_risk_ratings ON risk_matrix_risk_ratings.id = risk_register.current_risk_rating_id").
		Where("risk_matrix_risk_ratings.type IN ?", []string{"high", "very_high"}).
		Where(`NOT EXISTS (
			SELECT 1 FROM risk_controls
			JOIN controls ON controls.id = risk_controls.control_id
			JOIN control_efficiency_statuses ON control_efficiency_statuses.id = controls.control_efficiency_status_id
			WHERE risk_controls.risk_register_id = risk_register.id
			AND control_efficiency_statuses.type = 'efficient'
		)`)
	if err := criticalQuery.Count(&summary.CriticalRisksUncontrolled).Error; err != nil {
		return nil, err
	}

	// Overdue high-risk actions
	currentTime := time.Now()
	actionQuery := repo.db.Model(&models.RiskActionPlan{}).
		Joins("JOIN severity_levels ON severity_levels.id = risk_action_plans.severity_level_id").
		Joins("JOIN action_plan_statuses ON action_plan_statuses.id = risk_action_plans.action_plan_status_id").
		Where("severity_levels.type = ?", "high").
		Where("severity_levels.status = ?", "active").
		Where("action_plan_statuses.type != ?", "completed").
		Where("risk_action_plans.target_date < ?", currentTime)

	// Add date range filtering if parameters are provided
	if params.DateFrom != "" && params.DateTo != "" {
		actionQuery = actionQuery.Joins("JOIN risk_register ON risk_register.id = risk_action_plans.risk_register_id").
			Where("risk_register.reported_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}
	if err := actionQuery.Count(&summary.OverdueHighRiskActions).Error; err != nil {
		return nil, err
	}

	// Compliance-linked high risks
	complianceQuery := buildBaseQuery().
		Distinct("risk_compliance_requirements.compliance_requirement_id").
		Joins("JOIN risk_compliance_requirements ON risk_compliance_requirements.risk_register_id = risk_register.id").
		Joins("JOIN risk_matrix_risk_ratings ON risk_matrix_risk_ratings.id = risk_register.current_risk_rating_id").
		Where("risk_matrix_risk_ratings.type IN ?", []string{"high", "very_high"})

	if err := complianceQuery.Count(&summary.ComplianceLinkedHighRisks).Error; err != nil {
		return nil, err
	}

	return summary, nil
}

func (repo *ERMDashboardRepository) GetERMControlEfficiencySummary(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMControlEfficiencySummaryDTO, error) {
	type result struct {
		TotalControls     int64
		EfficientControls int64
	}

	var res result
	query := repo.db.Table("controls").
		Select(`
            COUNT(DISTINCT CASE WHEN controls.is_applicable = 1 THEN controls.id END) AS total_controls,
            COUNT(DISTINCT CASE WHEN controls.is_applicable = 1 AND control_efficiency_statuses.type = 'efficient' THEN controls.id END) AS efficient_controls
        `).
		Joins("JOIN risk_controls ON risk_controls.control_id = controls.id").
		Joins("JOIN control_efficiency_statuses ON controls.control_efficiency_status_id = control_efficiency_statuses.id").
		Where("controls.is_applicable = ?", true)

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Joins("JOIN risk_register ON risk_register.id = risk_controls.risk_register_id").
			Where("risk_register.reported_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&res).Error
	if err != nil {
		return nil, err
	}

	summary := &kpiDTOs.ERMControlEfficiencySummaryDTO{}
	if res.TotalControls > 0 {
		summary.OverallControlEfficiency = math.Round((float64(res.EfficientControls)/float64(res.TotalControls))*100*100) / 100
	} else {
		summary.OverallControlEfficiency = 0
	}

	return summary, nil
}

func (repo *ERMDashboardRepository) GetERMControlEfficiencyBreakdown(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMControlEfficiencyBreakdownDTO, error) {
	summary := &kpiDTOs.ERMControlEfficiencyBreakdownDTO{}

	query := repo.db.Model(&models.RiskControl{}).
		Select(`
			control_efficiency_statuses.id,
			control_efficiency_statuses.title_en,
			control_efficiency_statuses.title_ar,
			control_efficiency_statuses.color_code,
			control_efficiency_statuses.type,
			COUNT(DISTINCT controls.id) AS count
		`).
		Joins("JOIN controls ON risk_controls.control_id = controls.id").
		Joins("JOIN control_efficiency_statuses ON controls.control_efficiency_status_id = control_efficiency_statuses.id").
		Where("controls.is_applicable = ?", true).
		Group(`
        control_efficiency_statuses.id,
        control_efficiency_statuses.title_en,
        control_efficiency_statuses.title_ar,
        control_efficiency_statuses.color_code,
        control_efficiency_statuses.type
    `).
		Order("control_efficiency_statuses.id")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Joins("JOIN risk_register ON risk_register.id = risk_controls.risk_register_id").
			Where("risk_register.reported_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&summary.ControlsByEfficiency).Error
	if err != nil {
		return nil, err
	}

	return summary, nil
}

func (repo *ERMDashboardRepository) GetERMControlsOverview(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMControlsOverviewDTO, error) {
	summary := &kpiDTOs.ERMControlsOverviewDTO{}

	type controlsResult struct {
		TotalControls       int64
		ApplicableControls  int64
		ImplementedControls int64
		EfficientControls   int64
	}
	var cr controlsResult

	err := repo.db.Model(&models.RiskControl{}).
		Joins("JOIN controls ON risk_controls.control_id = controls.id").
		Joins("LEFT JOIN control_implementation_statuses ON controls.control_implementation_status_id = control_implementation_statuses.id").
		Joins("LEFT JOIN control_efficiency_statuses ON controls.control_efficiency_status_id = control_efficiency_statuses.id").
		Select(`
			COUNT(DISTINCT controls.id) AS total_controls,
			COUNT(DISTINCT CASE WHEN controls.is_applicable = 1 THEN controls.id END) AS applicable_controls,
			COUNT(DISTINCT CASE WHEN control_implementation_statuses.type = 'implemented' THEN controls.id END) AS implemented_controls,
			COUNT(DISTINCT CASE WHEN control_efficiency_statuses.type = 'efficient' THEN controls.id END) AS efficient_controls
		`).Scan(&cr).Error
	if err != nil {
		return nil, err
	}

	summary.TotalControls = cr.TotalControls
	summary.ApplicableControls = cr.ApplicableControls
	summary.ImplementedControls = cr.ImplementedControls

	summary.TotalControlsPercent = 100.0
	if cr.TotalControls > 0 {
		summary.ApplicableControlsPercent = math.Round(float64(cr.ApplicableControls)/float64(cr.TotalControls)*100*100) / 100
		summary.ImplementedControlsPercent = math.Round(float64(cr.ImplementedControls)/float64(cr.TotalControls)*100*100) / 100
	}

	if cr.ApplicableControls > 0 {
		summary.OverallControlEfficiency = math.Round(float64(cr.EfficientControls)/float64(cr.ApplicableControls)*100*100) / 100
	} else {
		summary.OverallControlEfficiency = 0
	}

	// High risk metrics
	type highRiskResult struct {
		ControlsMitigatingHighRisks      int64
		HighRisksWithNoControls          int64
		HighRisksWithInefficientControls int64
	}
	var hr highRiskResult

	err = repo.db.Table("risk_register").
		Joins("LEFT JOIN risk_controls ON risk_controls.risk_register_id = risk_register.id").
		Joins("LEFT JOIN risk_control_analyses ON risk_controls.id = risk_control_analyses.risk_control_id").
		Joins("LEFT JOIN control_efficiency_statuses ON risk_control_analyses.control_efficiency_status_id = control_efficiency_statuses.id").
		Joins("JOIN risk_matrix_risk_ratings ON risk_register.current_risk_rating_id = risk_matrix_risk_ratings.id").
		Joins("JOIN risk_statuses ON risk_register.risk_status_id = risk_statuses.id").
		Select(`
			COUNT(DISTINCT CASE WHEN risk_controls.control_id IS NOT NULL THEN risk_controls.control_id END) AS controls_mitigating_high_risks,
			COUNT(DISTINCT CASE WHEN risk_controls.id IS NULL THEN risk_register.id END) AS high_risks_with_no_controls,
			COUNT(DISTINCT CASE WHEN control_efficiency_statuses.type = 'inefficient' THEN risk_controls.control_id END) AS high_risks_with_inefficient_controls
		`).
		Where("risk_matrix_risk_ratings.type IN ?", []string{"high", "very_high"}).
		Where("risk_statuses.type = ?", "open").
		Scan(&hr).Error
	if err != nil {
		return nil, err
	}

	summary.ControlsMitigatingHighRisks = hr.ControlsMitigatingHighRisks
	summary.HighRisksWithNoControls = hr.HighRisksWithNoControls
	summary.HighRisksWithInefficientControls = hr.HighRisksWithInefficientControls

	if cr.TotalControls > 0 {
		summary.HighRisksWithInefficientPercent = math.Round(
			(float64(hr.HighRisksWithInefficientControls)/float64(cr.TotalControls)*100)*100,
		) / 100
	}

	return summary, nil
}

func (repo *ERMDashboardRepository) GetERMRiskAppetiteOverview(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskAppetiteOverviewDTO, error) {
	var results []struct {
		Department string `json:"department"`
		Above      int64  `json:"above"`
		Within     int64  `json:"within"`
		Below      int64  `json:"below"`
		TotalRisk  int64  `json:"total_risk"`
	}

	// Base query
	query := repo.db.Table("risk_register").
		Joins("LEFT JOIN organization_units ru_unit ON ru_unit.id = risk_register.responsible_unit_id").
		Joins("LEFT JOIN organization_units parent_unit ON COALESCE(ru_unit.organization_unit_id, ru_unit.id) = parent_unit.id").
		Joins("JOIN organization_levels ON organization_levels.id = ru_unit.organization_level_id AND organization_levels.is_primary = ?", true).
		Joins(`
			JOIN risk_evaluations
			ON risk_evaluations.risk_register_id = risk_register.id
			AND risk_evaluations.type = 'residual'
		`).
		Joins("JOIN risk_appetite_alignments ON risk_appetite_alignments.risk_evaluation_id = risk_evaluations.id").
		Joins("JOIN risk_appetite_tolerance_levels ON risk_appetite_tolerance_levels.id = risk_appetite_alignments.risk_appetite_tolerance_level_id AND risk_appetite_tolerance_levels.status = 'active'").
		Where("parent_unit.status = ?", "active").
		Select(`
			parent_unit.title_en AS department,
			SUM(CASE WHEN risk_appetite_tolerance_levels.type = 'above' THEN 1 ELSE 0 END) AS above,
			SUM(CASE WHEN risk_appetite_tolerance_levels.type = 'within' THEN 1 ELSE 0 END) AS within,
			SUM(CASE WHEN risk_appetite_tolerance_levels.type = 'below' THEN 1 ELSE 0 END) AS below,
			COUNT(*) AS total_risk
		`).
		Group("parent_unit.id, parent_unit.title_en")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where(
			"rr.reported_at BETWEEN ? AND ?",
			params.DateFrom+" 00:00:00",
			params.DateTo+" 23:59:59",
		)
	}

	// Execute query
	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	// Map results
	overview := make([]kpiDTOs.ERMRiskAppetiteOverviewItemDTO, len(results))
	for i, r := range results {
		overview[i] = kpiDTOs.ERMRiskAppetiteOverviewItemDTO{
			Department: r.Department,
			Above:      r.Above,
			Within:     r.Within,
			Below:      r.Below,
			TotalRisk:  r.TotalRisk,
		}
	}

	return &kpiDTOs.ERMRiskAppetiteOverviewDTO{
		RiskAppetiteByDepartments: overview,
	}, nil
}

func (repo *ERMDashboardRepository) GetERMRiskAppetiteSummary(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskAppetiteSummarytDTO, error) {
	type result struct {
		ID        uint
		Type      string
		TitleEn   string
		TitleAr   string
		ColorCode string
		Count     int64
	}

	var results []result
	var totalRisks int64

	baseQuery := repo.db.Table("risk_register").
		Joins(`
			JOIN risk_evaluations
			ON risk_evaluations.risk_register_id = risk_register.id
			AND risk_evaluations.type = 'residual'
		`).
		Joins("JOIN risk_appetite_alignments ON risk_appetite_alignments.risk_evaluation_id = risk_evaluations.id").
		Joins(`
			JOIN risk_appetite_tolerance_levels 
			ON risk_appetite_tolerance_levels.id = risk_appetite_alignments.risk_appetite_tolerance_level_id 
			AND risk_appetite_tolerance_levels.status = 'active'
		`).
		Select(`
			risk_appetite_tolerance_levels.id,
			risk_appetite_tolerance_levels.type,
			risk_appetite_tolerance_levels.title_en,
			risk_appetite_tolerance_levels.title_ar,
			risk_appetite_tolerance_levels.color_code,
			COUNT(*) AS count
		`).
		Group(`
			risk_appetite_tolerance_levels.id, 
			risk_appetite_tolerance_levels.type,
			risk_appetite_tolerance_levels.title_en,
			risk_appetite_tolerance_levels.title_ar,
			risk_appetite_tolerance_levels.color_code
		`)

	if params.DateFrom != "" && params.DateTo != "" {
		baseQuery = baseQuery.Where(
			"rr.reported_at BETWEEN ? AND ?",
			params.DateFrom+" 00:00:00",
			params.DateTo+" 23:59:59",
		)
	}

	if err := baseQuery.Scan(&results).Error; err != nil {
		return nil, err
	}

	levelsSummary := make([]kpiDTOs.ERMRiskAppetiteLevelDTO, len(results))
	for i, r := range results {
		levelsSummary[i] = kpiDTOs.ERMRiskAppetiteLevelDTO{
			ID:        r.ID,
			Type:      r.Type,
			TitleEn:   r.TitleEn,
			TitleAr:   r.TitleAr,
			ColorCode: r.ColorCode,
			Count:     r.Count,
		}
		totalRisks += r.Count
	}

	return &kpiDTOs.ERMRiskAppetiteSummarytDTO{
		RiskAppetiteLevelsSummary: levelsSummary,
		Total:                     totalRisks,
	}, nil
}

func (repo *ERMDashboardRepository) GetERMRiskAgingStatus(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskAgingStatusDTO, error) {

	var risks0To90Open, risks0To90Closed int64
	var risks91To180Open, risks91To180Closed int64
	var risks181To365Open, risks181To365Closed int64
	var risks1To2YearsOpen, risks1To2YearsClosed int64
	var risksAbove2YearsOpen, risksAbove2YearsClosed int64

	// Base query
	baseQuery := repo.db.Table("risk_register").
		Joins("JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id").
		Where("risk_statuses.type IN ?", []string{"open", "closed"})

	// Apply date filters if provided
	if params.DateFrom != "" && params.DateTo != "" {
		baseQuery = baseQuery.Where(
			"risk_register.reported_at BETWEEN ? AND ?",
			params.DateFrom+" 00:00:00",
			params.DateTo+" 23:59:59",
		)
	}

	// Calculate aging periods using DATEDIFF from reported_at to current date
	row := baseQuery.Select(`
        SUM(CASE WHEN risk_statuses.type = 'open' AND DATEDIFF(day, risk_register.reported_at, GETDATE()) BETWEEN 0 AND 90 THEN 1 ELSE 0 END) AS risks_0_to_90_open,
        SUM(CASE WHEN risk_statuses.type = 'closed' AND DATEDIFF(day, risk_register.reported_at, GETDATE()) BETWEEN 0 AND 90 THEN 1 ELSE 0 END) AS risks_0_to_90_closed,
        SUM(CASE WHEN risk_statuses.type = 'open' AND DATEDIFF(day, risk_register.reported_at, GETDATE()) BETWEEN 91 AND 180 THEN 1 ELSE 0 END) AS risks_91_to_180_open,
        SUM(CASE WHEN risk_statuses.type = 'closed' AND DATEDIFF(day, risk_register.reported_at, GETDATE()) BETWEEN 91 AND 180 THEN 1 ELSE 0 END) AS risks_91_to_180_closed,
        SUM(CASE WHEN risk_statuses.type = 'open' AND DATEDIFF(day, risk_register.reported_at, GETDATE()) BETWEEN 181 AND 365 THEN 1 ELSE 0 END) AS risks_181_to_365_open,
        SUM(CASE WHEN risk_statuses.type = 'closed' AND DATEDIFF(day, risk_register.reported_at, GETDATE()) BETWEEN 181 AND 365 THEN 1 ELSE 0 END) AS risks_181_to_365_closed,
        SUM(CASE WHEN risk_statuses.type = 'open' AND DATEDIFF(day, risk_register.reported_at, GETDATE()) BETWEEN 366 AND 730 THEN 1 ELSE 0 END) AS risks_1_to_2_years_open,
        SUM(CASE WHEN risk_statuses.type = 'closed' AND DATEDIFF(day, risk_register.reported_at, GETDATE()) BETWEEN 366 AND 730 THEN 1 ELSE 0 END) AS risks_1_to_2_years_closed,
        SUM(CASE WHEN risk_statuses.type = 'open' AND DATEDIFF(day, risk_register.reported_at, GETDATE()) > 730 THEN 1 ELSE 0 END) AS risks_above_2_years_open,
        SUM(CASE WHEN risk_statuses.type = 'closed' AND DATEDIFF(day, risk_register.reported_at, GETDATE()) > 730 THEN 1 ELSE 0 END) AS risks_above_2_years_closed
    `)

	err := row.Row().Scan(
		&risks0To90Open, &risks0To90Closed,
		&risks91To180Open, &risks91To180Closed,
		&risks181To365Open, &risks181To365Closed,
		&risks1To2YearsOpen, &risks1To2YearsClosed,
		&risksAbove2YearsOpen, &risksAbove2YearsClosed,
	)
	if err != nil {
		return nil, err
	}

	agingPeriods := []kpiDTOs.RiskAgingPeriodDTO{
		{
			RiskAging: "0-90 Days",
			Open:      risks0To90Open,
			Closed:    risks0To90Closed,
			Total:     risks0To90Open + risks0To90Closed,
		},
		{
			RiskAging: "91-180 Days",
			Open:      risks91To180Open,
			Closed:    risks91To180Closed,
			Total:     risks91To180Open + risks91To180Closed,
		},
		{
			RiskAging: "181-365 Days",
			Open:      risks181To365Open,
			Closed:    risks181To365Closed,
			Total:     risks181To365Open + risks181To365Closed,
		},
		{
			RiskAging: "1-2 Years",
			Open:      risks1To2YearsOpen,
			Closed:    risks1To2YearsClosed,
			Total:     risks1To2YearsOpen + risks1To2YearsClosed,
		},
		{
			RiskAging: "Above 2 Years",
			Open:      risksAbove2YearsOpen,
			Closed:    risksAbove2YearsClosed,
			Total:     risksAbove2YearsOpen + risksAbove2YearsClosed,
		},
	}

	summary := &kpiDTOs.ERMRiskAgingStatusDTO{
		RiskAgingPeriods: agingPeriods,
	}

	return summary, nil
}

func (repo *ERMDashboardRepository) GetERMRiskTreatmentPlanProgress(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskTreatmentPlanProgressDTO, error) {
	var results []struct {
		ID        uint   `json:"id"`
		TitleEn   string `json:"title_en"`
		TitleAr   string `json:"title_ar"`
		ColorCode string `json:"color_code"`
		Type      string `json:"type"`
		OnTime    int64  `json:"on_time"`
		Overdue   int64  `json:"overdue"`
	}

	query := repo.db.Table("risk_action_plans").
		Joins("JOIN action_plan_statuses ON action_plan_statuses.id = risk_action_plans.action_plan_status_id").
		Where("action_plan_statuses.status = ?", "active").
		Select(`
			action_plan_statuses.id,
			action_plan_statuses.title_en,
			action_plan_statuses.title_ar,
			action_plan_statuses.color_code,
			action_plan_statuses.type,
			SUM(CASE 
				WHEN risk_action_plans.target_date IS NULL OR risk_action_plans.target_date >= GETDATE() OR action_plan_statuses.is_closed = 1 
				THEN 1 
				ELSE 0 
			END) as on_time,
			SUM(CASE 
				WHEN risk_action_plans.target_date IS NOT NULL AND risk_action_plans.target_date < GETDATE() AND action_plan_statuses.is_closed = 0 
				THEN 1 
				ELSE 0 
			END) as overdue
		`).
		Group(`
			action_plan_statuses.id,
			action_plan_statuses.title_en,
			action_plan_statuses.title_ar,
			action_plan_statuses.color_code,
			action_plan_statuses.type
		`).
		Order(`
			CASE action_plan_statuses.type
				WHEN 'new' THEN 1
				WHEN 'in_progress' THEN 2  
				WHEN 'completed' THEN 3
				WHEN 'close' THEN 4
				ELSE 5
			END
		`)

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where(
			"risk_action_plans.created_at BETWEEN ? AND ?",
			params.DateFrom+" 00:00:00",
			params.DateTo+" 23:59:59",
		)
	}

	err := query.Scan(&results).Error
	if err != nil {
		return nil, err
	}

	treatmentPlanProgress := make([]kpiDTOs.ERMRiskTreatmentPlanProgressItemDTO, len(results))
	for i, result := range results {
		treatmentPlanProgress[i] = kpiDTOs.ERMRiskTreatmentPlanProgressItemDTO{
			ID:        result.ID,
			TitleEn:   result.TitleEn,
			TitleAr:   result.TitleAr,
			ColorCode: result.ColorCode,
			Type:      result.Type,
			OnTime:    result.OnTime,
			Overdue:   result.Overdue,
		}
	}

	return &kpiDTOs.ERMRiskTreatmentPlanProgressDTO{
		RiskTreatmentPlanProgress: treatmentPlanProgress,
	}, nil
}

func (repo *ERMDashboardRepository) GetERMStrategicObjectivesRiskOverview(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMStrategicObjectiveRiskOverviewDTO, error) {
	summary := &kpiDTOs.ERMStrategicObjectiveRiskOverviewDTO{}

	baseQuery := repo.db.Table("strategic_objectives so").
		Select(`
            so.title_en as strategic_objective,
            (
				SELECT STRING_AGG(rr.reference_number, ', ')
				WITHIN GROUP (ORDER BY rr.reference_number)
				FROM risk_register rr
				JOIN risk_strategic_objectives rso
				ON rr.id = rso.risk_register_id
				WHERE rso.strategic_objective_id = so.id
			) as linked_risks,
            COUNT(CASE WHEN rmrr.type IN ('high', 'very_high') THEN 1 END) as high_risk_count,
            CASE 
                WHEN MAX(CASE WHEN rmrr.type = 'very_high' THEN 5 
                             WHEN rmrr.type = 'high' THEN 4 
                             WHEN rmrr.type = 'medium' THEN 3 
                             WHEN rmrr.type = 'low' THEN 2 
                             WHEN rmrr.type = 'very_low' THEN 1 
                             ELSE 0 END) = 5 THEN 'Very High'
                WHEN MAX(CASE WHEN rmrr.type = 'very_high' THEN 5 
                             WHEN rmrr.type = 'high' THEN 4 
                             WHEN rmrr.type = 'medium' THEN 3 
                             WHEN rmrr.type = 'low' THEN 2 
                             WHEN rmrr.type = 'very_low' THEN 1 
                             ELSE 0 END) = 4 THEN 'High'
                WHEN MAX(CASE WHEN rmrr.type = 'very_high' THEN 5 
                             WHEN rmrr.type = 'high' THEN 4 
                             WHEN rmrr.type = 'medium' THEN 3 
                             WHEN rmrr.type = 'low' THEN 2 
                             WHEN rmrr.type = 'very_low' THEN 1 
                             ELSE 0 END) = 3 THEN 'Medium'
                WHEN MAX(CASE WHEN rmrr.type = 'very_high' THEN 5 
                             WHEN rmrr.type = 'high' THEN 4 
                             WHEN rmrr.type = 'medium' THEN 3 
                             WHEN rmrr.type = 'low' THEN 2 
                             WHEN rmrr.type = 'very_low' THEN 1 
                             ELSE 0 END) = 2 THEN 'Low'
                WHEN MAX(CASE WHEN rmrr.type = 'very_high' THEN 5 
                             WHEN rmrr.type = 'high' THEN 4 
                             WHEN rmrr.type = 'medium' THEN 3 
                             WHEN rmrr.type = 'low' THEN 2 
                             WHEN rmrr.type = 'very_low' THEN 1 
                             ELSE 0 END) = 1 THEN 'Very Low'
                ELSE 'No Rating'
            END as risk_rating,
            CASE 
                WHEN COUNT(CASE WHEN ratl.type IN ('above', 'significantly_above') THEN 1 END) > 0 THEN 'Yes'
                ELSE 'No'
            END as appetite_breach
        `).
		Joins("JOIN risk_strategic_objectives rso ON so.id = rso.strategic_objective_id").
		Joins("LEFT JOIN risk_register rr ON rr.id = rso.risk_register_id").
		Joins("LEFT JOIN risk_matrix_risk_ratings rmrr ON rmrr.id = rr.current_risk_rating_id").
		Joins(`
            LEFT JOIN risk_evaluations re
            ON re.risk_register_id = rr.id
            AND re.type = 'residual'
        `).
		Joins("LEFT JOIN risk_appetite_alignments raa ON raa.risk_evaluation_id = re.id").
		Joins("LEFT JOIN risk_appetite_tolerance_levels ratl ON ratl.id = raa.risk_appetite_tolerance_level_id AND ratl.status = 'active'").
		Where("so.status = ? AND rso.strategic_objective_id IS NOT NULL", "active")

	if params.DateFrom != "" && params.DateTo != "" {
		baseQuery = baseQuery.Where(
			"rr.reported_at BETWEEN ? AND ?",
			params.DateFrom+" 00:00:00",
			params.DateTo+" 23:59:59",
		)
	}

	// Group by strategic objective - only include objectives that have linked risks
	err := baseQuery.Group("so.id, so.title_en").Scan(&summary.StrategicObjectiveRisks).Error
	if err != nil {
		return nil, err
	}

	return summary, nil
}

func (repo *ERMDashboardRepository) GetERMStrategicObjectivesRiskSummary(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMStrategicObjectivesRiskSummaryDTO, error) {

	summary := &kpiDTOs.ERMStrategicObjectivesRiskSummaryDTO{}

	baseQuery := repo.db.Table("strategic_objectives so")

	// Apply date filters if provided
	if params.DateFrom != "" && params.DateTo != "" {
		baseQuery = baseQuery.Where(
			"so.created_at BETWEEN ? AND ?",
			params.DateFrom+" 00:00:00",
			params.DateTo+" 23:59:59",
		)
	}

	row := baseQuery.Select(`
        COUNT(DISTINCT CASE WHEN rso.strategic_objective_id IS NOT NULL THEN so.id END) as objectives_with_linked_risks,
        COUNT(DISTINCT CASE WHEN rmrr.type IN ('high', 'very_high') THEN so.id END) as objectives_with_high_rated_risks,
        COUNT(DISTINCT CASE WHEN rso.strategic_objective_id IS NULL THEN so.id END) as objectives_with_no_risk_linked
    `).
		Joins("LEFT JOIN risk_strategic_objectives rso ON so.id = rso.strategic_objective_id").
		Joins("LEFT JOIN risk_register rr ON rr.id = rso.risk_register_id").
		Joins("LEFT JOIN risk_matrix_risk_ratings rmrr ON rmrr.id = rr.current_risk_rating_id").
		Where("so.status = ?", "active")

	err := row.Row().Scan(
		&summary.ObjectivesWithLinkedRisks,
		&summary.ObjectivesWithHighRatedRisks,
		&summary.ObjectivesWithNoRiskLinked,
	)
	if err != nil {
		return nil, err
	}

	return summary, nil
}
