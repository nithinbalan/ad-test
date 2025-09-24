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
