package repositories

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	kpiDTOs "jethur/grc/kpi/dtos"
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"
	"math"
	"time"

	"gorm.io/gorm"
)

type ERMDashboardRepositoryInterface interface {
	GetERMRiskOverview(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskOverviewDTO, error)
	GetERMRiskDistribution(params *kpiDTOs.ERMDashboardQueryParamsDTO) ([]kpiDTOs.ERMRiskDistributionDTO, error)
	GetERMOverdueRisks(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMOverdueRisksDTO, error)
	GetERMRisksAboveAppetite(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRisksAboveAppetiteDTO, error)
	GetERMLeadershipAttentionMetrics(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMLeadershipAttentionDTO, error)
	GetERMControlEfficiencySummary(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMControlEfficiencySummaryDTO, error)
	GetERMControlEfficiencyBreakdown(params *kpiDTOs.ERMDashboardQueryParamsDTO) ([]kpiDTOs.ERMRiskSummaryItemDTO, error)
	GetERMControlsOverview(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMControlsOverviewDTO, error)
	GetERMRiskAppetiteOverview(params *kpiDTOs.ERMDashboardQueryParamsDTO) ([]kpiDTOs.ERMRiskAppetiteOverviewDTO, error)
	GetERMRiskAppetiteSummary(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskAppetiteSummarytDTO, error)
	GetERMRiskAgingStatus(params *kpiDTOs.ERMDashboardQueryParamsDTO) ([]kpiDTOs.ERMRiskAgingStatusDTO, error)
	GetERMRiskTreatmentPlanProgress(params *kpiDTOs.ERMDashboardQueryParamsDTO) ([]kpiDTOs.ERMRiskTreatmentPlanProgressDTO, error)
	GetERMStrategicObjectivesRiskOverview(params *kpiDTOs.ERMDashboardListQueryParamsDTO) (*kpiDTOs.ERMStrategicObjectiveRiskOverviewDTO, error)
	GetERMStrategicObjectivesRiskSummary(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMStrategicObjectivesRiskSummaryDTO, error)
	GetERMRiskSummaryByOrganizationUnit(params *kpiDTOs.ERMDashboardListQueryParamsDTO) (*kpiDTOs.ERMRiskOrganizationUnitSummaryDTO, error)

	// ERM Risk Details
	GetERMRiskCategorySummary(params *kpiDTOs.ERMDashboardQueryParamsDTO) ([]kpiDTOs.ERMRiskCategorySummaryDTO, error)
	GetERMRiskSourceSummary(params *kpiDTOs.ERMDashboardQueryParamsDTO) ([]kpiDTOs.ERMRiskSourceSummaryDTO, error)
	GetERMRiskStatusSummary(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskStatusSummaryDTO, error)
	GetERMRiskOrganizationUnitSummary(params *kpiDTOs.ERMDashboardQueryParamsDTO) ([]kpiDTOs.ERMRiskSummaryByOrganizationUnitDTO, error)
	GetERMRisksIdentifiedSummary() ([]kpiDTOs.ERMRisksIdentifiedSummaryDTO, error)
	GetERMRiskResponseStrategySummary(params *kpiDTOs.ERMDashboardQueryParamsDTO) ([]kpiDTOs.ERMRiskResponseStrategySummaryDTO, error)
	GetERMRiskControlEfficiencySummary(params *kpiDTOs.ERMDashboardQueryParamsDTO) ([]kpiDTOs.ERMRiskControlEfficiencySummaryDTO, error)
	GetERMRiskMitigationPlanStatusSummary(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskMitigationPlanStatusSummaryDTO, error)
	GetERMOpenRisksSummary(params *kpiDTOs.ERMDashboardListQueryParamsDTO) (*kpiDTOs.ERMOpenRisksSummaryDTO, error)

	GetTopRiskRegisters(params *kpiDTOs.ERMDashboardQueryParamsDTO) ([]*riskDTOs.RiskRegisterListDTO, error)
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

	summary := &kpiDTOs.ERMRiskOverviewDTO{}

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

	err := row.Row().Scan(&summary.OpenRisks, &summary.ClosedRisks, &summary.TotalRisks)
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
	err = highRiskQuery.Count(&summary.UnresolvedHighRisks).Error
	if err != nil {
		return nil, err
	}

	return summary, nil
}

func (repo *ERMDashboardRepository) GetERMRiskDistribution(params *kpiDTOs.ERMDashboardQueryParamsDTO) ([]kpiDTOs.ERMRiskDistributionDTO, error) {

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

	result := make([]kpiDTOs.ERMRiskDistributionDTO, 0)

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

		result = append(result, kpiDTOs.ERMRiskDistributionDTO{
			TitleEn:    rating.TitleEn,
			TitleAr:    rating.TitleAr,
			ColorCode:  rating.ColorCode,
			Type:       rating.Type,
			Count:      count,
			OpenCount:  openCount,
			Percentage: percentage,
		})
	}

	return result, nil
}

func (repo *ERMDashboardRepository) GetERMOverdueRisks(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMOverdueRisksDTO, error) {

	results := []kpiDTOs.ERMRiskSummaryItemDTO{}
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

	if err := baseQuery.Scan(&results).Error; err != nil {
		return nil, err
	}

	// Build the overdue details + total
	var totalRisksOverdue int64
	risksOverdueDetails := []kpiDTOs.ERMRiskSummaryItemDTO{}

	for _, r := range results {
		risksOverdueDetails = append(risksOverdueDetails, kpiDTOs.ERMRiskSummaryItemDTO{
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

	results := []kpiDTOs.ERMRiskSummaryItemDTO{}

	// Base query with risk levels
	baseQuery := repo.db.Model(&models.RiskRegister{}).
		Select(`
        risk_matrix_risk_ratings.id,
        risk_matrix_risk_ratings.type,
        risk_matrix_risk_ratings.title_en,
        risk_matrix_risk_ratings.title_ar,
        risk_matrix_risk_ratings.color_code,
        COUNT(DISTINCT risk_register.id) AS count
    `).
		Joins("JOIN risk_evaluations ON risk_evaluations.risk_register_id = risk_register.id").
		Joins("JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id").
		Joins("JOIN risk_appetite_alignments ON risk_appetite_alignments.risk_evaluation_id = risk_evaluations.id").
		Joins("JOIN risk_appetite_tolerance_levels ON risk_appetite_tolerance_levels.id = risk_appetite_alignments.risk_appetite_tolerance_level_id").
		Joins("LEFT JOIN risk_matrix_risk_ratings ON risk_matrix_risk_ratings.id = risk_register.current_risk_rating_id").
		Where("risk_statuses.type IN ?", []string{"open", "closed"}).
		Where("risk_appetite_tolerance_levels.type = ?", "above").
		Where("risk_evaluations.type = ?", "inherent").
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

	if err := baseQuery.Scan(&results).Error; err != nil {
		return nil, err
	}

	var totalRisksAboveAppetite int64
	risksByLevels := []kpiDTOs.ERMRiskSummaryItemDTO{}

	for _, r := range results {
		risksByLevels = append(risksByLevels, kpiDTOs.ERMRiskSummaryItemDTO{
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

	// Determine the highest available risk rating dynamically
	var highestRating string
	err := repo.db.Table("risk_matrix_risk_ratings").
		Select("type").
		Where("status = ?", "active").
		Order("CASE type WHEN 'very_high' THEN 1 WHEN 'high' THEN 2 WHEN 'medium' THEN 3 WHEN 'low' THEN 4 WHEN 'very_low' THEN 5 END").
		Limit(1).
		Pluck("type", &highestRating).Error
	if err != nil || highestRating == "" {
		return nil, fmt.Errorf("failed to fetch highest risk rating: %v", err)
	}

	// Base query function
	buildBaseQuery := func() *gorm.DB {
		query := repo.db.Model(&models.RiskRegister{}).
			Joins("JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id").
			Joins("JOIN risk_matrix_risk_ratings ON risk_matrix_risk_ratings.id = risk_register.current_risk_rating_id").
			Where("risk_statuses.type = ?", "open").
			Where("risk_matrix_risk_ratings.type = ?", highestRating)

		if params.DateFrom != "" && params.DateTo != "" {
			query = query.Where("risk_register.reported_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
		}
		return query
	}

	// Strategic objectives at risk
	strategicQuery := buildBaseQuery().
		Joins("JOIN risk_strategic_objectives ON risk_strategic_objectives.risk_register_id = risk_register.id").
		Distinct("risk_register.id")
	if err := strategicQuery.Count(&summary.StrategicObjectivesAtRisk).Error; err != nil {
		return nil, err
	}

	// Critical risks uncontrolled
	criticalQuery := buildBaseQuery().
		Where(`NOT EXISTS (
			SELECT 1 FROM risk_controls
			JOIN controls ON controls.id = risk_controls.control_id
			JOIN control_efficiency_statuses ON control_efficiency_statuses.id = controls.control_efficiency_status_id
			WHERE risk_controls.risk_register_id = risk_register.id
		)`)
	if err = criticalQuery.Count(&summary.CriticalRisksUncontrolled).Error; err != nil {
		return nil, err
	}

	// Overdue high-risk actions
	currentTime := time.Now()
	actionQuery := buildBaseQuery().
		Joins("JOIN risk_action_plans ON risk_action_plans.risk_register_id = risk_register.id").
		Joins("JOIN action_plan_statuses ON action_plan_statuses.id = risk_action_plans.action_plan_status_id").
		Where("action_plan_statuses.type != ?", "completed").
		Where("risk_action_plans.target_date < ?", currentTime).
		Distinct("risk_register.id")

	if err := actionQuery.Count(&summary.OverdueHighRiskActions).Error; err != nil {
		return nil, err
	}

	// Compliance-linked high risks
	complianceQuery := buildBaseQuery().
		Joins("JOIN risk_compliance_requirements ON risk_compliance_requirements.risk_register_id = risk_register.id").
		Distinct("risk_register.id")

	if err := complianceQuery.Count(&summary.ComplianceLinkedHighRisks).Error; err != nil {
		return nil, err
	}

	return summary, nil
}

func (repo *ERMDashboardRepository) GetERMControlEfficiencySummary(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMControlEfficiencySummaryDTO, error) {

	var totalControls, efficientControls int64

	query := repo.db.Table("risk_controls").
		Joins("JOIN controls ON risk_controls.control_id = controls.id").
		Joins("JOIN risk_register ON risk_register.id = risk_controls.risk_register_id").
		Joins("JOIN risk_statuses ON risk_register.risk_status_id = risk_statuses.id").
		Joins("LEFT JOIN control_efficiency_statuses ON controls.control_efficiency_status_id = control_efficiency_statuses.id").
		Where("risk_statuses.type IN ?", []string{"open", "closed"}).
		Select(`
			COUNT(DISTINCT controls.id) AS total_controls,
			COUNT(DISTINCT CASE WHEN control_efficiency_statuses.type = 'efficient' THEN controls.id END) AS efficient_controls
		`)

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Joins("JOIN risk_register ON risk_register.id = risk_controls.risk_register_id").
			Where("risk_register.reported_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Row().Scan(&totalControls, &efficientControls)
	if err != nil {
		return nil, err
	}

	summary := &kpiDTOs.ERMControlEfficiencySummaryDTO{}
	if totalControls > 0 {
		summary.OverallControlEfficiency = math.Round((float64(efficientControls)/float64(totalControls))*100*100) / 100
	} else {
		summary.OverallControlEfficiency = 0
	}

	return summary, nil
}

func (repo *ERMDashboardRepository) GetERMControlEfficiencyBreakdown(params *kpiDTOs.ERMDashboardQueryParamsDTO) ([]kpiDTOs.ERMRiskSummaryItemDTO, error) {
	summary := []kpiDTOs.ERMRiskSummaryItemDTO{}

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

	err := query.Scan(&summary).Error
	if err != nil {
		return nil, err
	}

	return summary, nil
}

func (repo *ERMDashboardRepository) GetERMControlsOverview(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMControlsOverviewDTO, error) {
	summary := &kpiDTOs.ERMControlsOverviewDTO{}

	controlsResult := &kpiDTOs.ERMControlsResultDTO{}
	highRiskControlResult := &kpiDTOs.ERMHighRiskControlsResultDTO{}

	err := repo.db.Model(&models.RiskControl{}).
		Joins("JOIN controls ON risk_controls.control_id = controls.id").
		Joins("JOIN risk_register ON risk_controls.risk_register_id = risk_register.id").
		Joins("JOIN risk_statuses ON risk_register.risk_status_id = risk_statuses.id").
		Joins("LEFT JOIN control_implementation_statuses ON controls.control_implementation_status_id = control_implementation_statuses.id").
		Joins("LEFT JOIN control_efficiency_statuses ON controls.control_efficiency_status_id = control_efficiency_statuses.id").
		Where("risk_statuses.type IN ?", []string{"open", "closed"}).
		Select(`
			COUNT(DISTINCT controls.id) AS total_controls,
			COUNT(DISTINCT CASE WHEN control_implementation_statuses.type = 'implemented' THEN controls.id END) AS implemented_controls,
			COUNT(DISTINCT CASE WHEN control_efficiency_statuses.type = 'efficient' THEN controls.id END) AS efficient_controls
		`).Scan(&controlsResult).Error
	if err != nil {
		return nil, err
	}

	summary.TotalControls = controlsResult.TotalControls
	summary.ImplementedControls = controlsResult.ImplementedControls
	summary.EfficientControls = controlsResult.EfficientControls

	if controlsResult.TotalControls > 0 {
		summary.ImplementedControlsPercent = math.Round(float64(controlsResult.ImplementedControls)/float64(controlsResult.TotalControls)*100*100) / 100
	}

	if controlsResult.TotalControls > 0 {
		summary.EfficientControlsPercent = math.Round(float64(controlsResult.EfficientControls)/float64(controlsResult.TotalControls)*100*100) / 100
	}

	if controlsResult.TotalControls > 0 {
		summary.OverallControlEfficiency = math.Round(float64(controlsResult.EfficientControls)/float64(controlsResult.TotalControls)*100*100) / 100
	} else {
		summary.OverallControlEfficiency = 0
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

	err = repo.db.Table("risk_register").
		Joins("LEFT JOIN risk_controls ON risk_controls.risk_register_id = risk_register.id").
		Joins("LEFT JOIN controls ON risk_controls.control_id = controls.id").
		Joins("LEFT JOIN control_efficiency_statuses ON controls.control_efficiency_status_id = control_efficiency_statuses.id").
		Joins("JOIN risk_matrix_risk_ratings ON risk_register.current_risk_rating_id = risk_matrix_risk_ratings.id").
		Joins("JOIN risk_statuses ON risk_register.risk_status_id = risk_statuses.id").
		Select(`
			COUNT(DISTINCT CASE WHEN risk_controls.control_id IS NOT NULL THEN risk_controls.control_id END) AS controls_mitigating_high_risks,
			COUNT(DISTINCT CASE WHEN risk_controls.id IS NULL THEN risk_register.id END) AS high_risks_with_no_controls,
			COUNT(DISTINCT CASE WHEN control_efficiency_statuses.type = 'inefficient' THEN risk_controls.control_id END) AS high_risks_with_inefficient_controls
		`).
		Where("risk_matrix_risk_ratings.type = ?", highestRating).
		Where("risk_statuses.type = ?", "open").
		Scan(&highRiskControlResult).Error
	if err != nil {
		return nil, err
	}

	summary.ControlsMitigatingHighRisks = highRiskControlResult.ControlsMitigatingHighRisks
	summary.HighRisksWithNoControls = highRiskControlResult.HighRisksWithNoControls
	summary.HighRisksWithInefficientControls = highRiskControlResult.HighRisksWithInefficientControls

	if controlsResult.TotalControls > 0 {
		summary.HighRisksWithInefficientControlsPercent = math.Round(
			(float64(highRiskControlResult.HighRisksWithInefficientControls)/float64(controlsResult.TotalControls)*100)*100,
		) / 100
	}

	return summary, nil
}

func (repo *ERMDashboardRepository) GetERMRiskAppetiteOverview(params *kpiDTOs.ERMDashboardQueryParamsDTO) ([]kpiDTOs.ERMRiskAppetiteOverviewDTO, error) {
	results := []kpiDTOs.ERMRiskAppetiteOverviewQueryResultDTO{}

	query := repo.db.Model(&models.RiskRegister{}).
		Joins("LEFT JOIN organization_units ru ON ru.id = risk_register.responsible_unit_id").
		Joins("JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id").
		Joins("LEFT JOIN organization_units parent_unit ON COALESCE(ru.organization_unit_id, ru.id) = parent_unit.id").
		Joins("JOIN organization_levels ol ON ol.id = ru.organization_level_id AND ol.is_primary = ?", true).
		Joins("JOIN risk_evaluations re ON re.risk_register_id = risk_register.id AND re.type = 'inherent'").
		Joins("JOIN risk_appetite_alignments raa ON raa.risk_evaluation_id = re.id").
		Joins("JOIN risk_appetite_tolerance_levels tl ON tl.id = raa.risk_appetite_tolerance_level_id AND tl.status = 'active'").
		Where("risk_statuses.type IN ?", []string{"open", "closed"}).
		Where("parent_unit.status = ?", "active").
		Select(`
			parent_unit.title_en AS organization_unit,
			tl.id AS level_id, tl.title_en, tl.title_ar, tl.color_code, tl.type,
			COUNT(*) AS count
		`).
		Group("parent_unit.id, parent_unit.title_en, tl.id, tl.title_en, tl.title_ar, tl.color_code, tl.type")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("risk_register.reported_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	// Fetch total risks per organization unit
	totalRisksMap := make(map[string]int64)
	totalQuery := repo.db.Model(&models.RiskRegister{}).
		Joins("LEFT JOIN organization_units ru ON ru.id = risk_register.responsible_unit_id").
		Joins("JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id").
		Joins("LEFT JOIN organization_units parent_unit ON COALESCE(ru.organization_unit_id, ru.id) = parent_unit.id").
		Joins("JOIN organization_levels ol ON ol.id = ru.organization_level_id AND ol.is_primary = ?", true).
		Where("risk_statuses.type IN ?", []string{"open", "closed"}).
		Where("parent_unit.status = ?", "active")

	if params.DateFrom != "" && params.DateTo != "" {
		totalQuery = totalQuery.Where("risk_register.reported_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	var totalResults []struct {
		OrganizationUnit string
		TotalRisk        int64
	}

	if err := totalQuery.Select("parent_unit.title_en AS organization_unit, COUNT(DISTINCT risk_register.id) AS total_risk").
		Group("parent_unit.id, parent_unit.title_en").
		Scan(&totalResults).Error; err != nil {
		return nil, err
	}

	for _, tr := range totalResults {
		totalRisksMap[tr.OrganizationUnit] = tr.TotalRisk
	}

	// Step 3: Map results → Organization units with dynamic levels
	orgMap := make(map[string]*kpiDTOs.ERMRiskAppetiteOverviewDTO)
	for _, r := range results {
		if _, exists := orgMap[r.OrganizationUnit]; !exists {
			orgMap[r.OrganizationUnit] = &kpiDTOs.ERMRiskAppetiteOverviewDTO{
				OrganizationUnit:   r.OrganizationUnit,
				RiskAppetiteLevels: []kpiDTOs.ERMRiskSummaryItemDTO{},
				TotalRisk:          totalRisksMap[r.OrganizationUnit],
			}
		}
		org := orgMap[r.OrganizationUnit]
		org.RiskAppetiteLevels = append(org.RiskAppetiteLevels, kpiDTOs.ERMRiskSummaryItemDTO{
			ID:        r.LevelID,
			TitleEn:   r.TitleEn,
			TitleAr:   r.TitleAr,
			ColorCode: r.ColorCode,
			Type:      r.Type,
			Count:     r.Count,
		})
	}

	// Convert map → slice
	overview := make([]kpiDTOs.ERMRiskAppetiteOverviewDTO, 0, len(orgMap))
	for _, org := range orgMap {
		overview = append(overview, *org)
	}

	return overview, nil
}

func (repo *ERMDashboardRepository) GetERMRiskAppetiteSummary(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskAppetiteSummarytDTO, error) {
	results := []kpiDTOs.ERMRiskSummaryItemDTO{}
	var totalRisks int64

	baseQuery := repo.db.Table("risk_register").
		Joins(`
			JOIN risk_evaluations
			ON risk_evaluations.risk_register_id = risk_register.id
			AND risk_evaluations.type = 'inherent'
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

	levelsSummary := make([]kpiDTOs.ERMRiskSummaryItemDTO, len(results))
	for i, r := range results {
		levelsSummary[i] = kpiDTOs.ERMRiskSummaryItemDTO{
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

func (repo *ERMDashboardRepository) GetERMRiskAgingStatus(params *kpiDTOs.ERMDashboardQueryParamsDTO) ([]kpiDTOs.ERMRiskAgingStatusDTO, error) {

	result := &kpiDTOs.ERMRiskAgingStatusQuereyResultDTO{}

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

	err := row.Scan(&result).Error
	if err != nil {
		return nil, err
	}

	agingPeriods := []kpiDTOs.ERMRiskAgingStatusDTO{
		{
			RiskAging: "0-90 Days",
			Open:      result.Risks0To90Open,
			Closed:    result.Risks0To90Closed,
			Total:     result.Risks0To90Open + result.Risks0To90Closed,
		},
		{
			RiskAging: "91-180 Days",
			Open:      result.Risks91To180Open,
			Closed:    result.Risks91To180Closed,
			Total:     result.Risks91To180Open + result.Risks91To180Closed,
		},
		{
			RiskAging: "181-365 Days",
			Open:      result.Risks181To365Open,
			Closed:    result.Risks181To365Closed,
			Total:     result.Risks181To365Open + result.Risks181To365Closed,
		},
		{
			RiskAging: "1-2 Years",
			Open:      result.Risks1To2YearsOpen,
			Closed:    result.Risks1To2YearsClosed,
			Total:     result.Risks1To2YearsOpen + result.Risks1To2YearsClosed,
		},
		{
			RiskAging: "Above 2 Years",
			Open:      result.RisksAbove2YearsOpen,
			Closed:    result.RisksAbove2YearsClosed,
			Total:     result.RisksAbove2YearsOpen + result.RisksAbove2YearsClosed,
		},
	}

	return agingPeriods, nil
}

func (repo *ERMDashboardRepository) GetERMRiskTreatmentPlanProgress(params *kpiDTOs.ERMDashboardQueryParamsDTO) ([]kpiDTOs.ERMRiskTreatmentPlanProgressDTO, error) {
	results := []kpiDTOs.ERMRiskTreatmentPlanProgressDTO{}

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

	return results, nil
}

func (repo *ERMDashboardRepository) GetERMStrategicObjectivesRiskOverview(params *kpiDTOs.ERMDashboardListQueryParamsDTO) (*kpiDTOs.ERMStrategicObjectiveRiskOverviewDTO, error) {
	strategicObjectiveRisks := []kpiDTOs.ERMStregicObjectiveRiskItemDTO{}
	var totalCount int64

	// Determine the highest available risk rating dynamically
	var highestRating string
	err := repo.db.Table("risk_matrix_risk_ratings").
		Select("type").
		Where("status = ?", "active").
		Order("CASE type WHEN 'very_high' THEN 1 WHEN 'high' THEN 2 WHEN 'medium' THEN 3 WHEN 'low' THEN 4 WHEN 'very_low' THEN 5 END").
		Limit(1).
		Pluck("type", &highestRating).Error
	if err != nil || highestRating == "" {
		return nil, fmt.Errorf("failed to fetch highest risk rating: %v", err)
	}

	// Format the SQL query with the highest rating
	selectQuery := fmt.Sprintf(`
            so.title_en as strategic_objective_title_en,
			so.title_ar as strategic_objective_title_ar,
            (
				SELECT STRING_AGG(rr.reference_number, ', ')
				WITHIN GROUP (ORDER BY rr.reference_number)
				FROM risk_register rr
				JOIN risk_strategic_objectives rso
				ON rr.id = rso.risk_register_id
				WHERE rso.strategic_objective_id = so.id
			) as linked_risks,
            COUNT(CASE WHEN rmrr.type = '%s' THEN 1 END) as high_risk_count,
            COALESCE(
                (SELECT TOP 1 title_en FROM risk_matrix_risk_ratings 
                 WHERE CASE 
                     WHEN type = 'very_high' THEN 5
                     WHEN type = 'high' THEN 4
                     WHEN type = 'medium' THEN 3
                     WHEN type = 'low' THEN 2
                     WHEN type = 'very_low' THEN 1
                     ELSE 0 
                 END = MAX(CASE 
                     WHEN rmrr.type = 'very_high' THEN 5
                     WHEN rmrr.type = 'high' THEN 4
                     WHEN rmrr.type = 'medium' THEN 3
                     WHEN rmrr.type = 'low' THEN 2
                     WHEN rmrr.type = 'very_low' THEN 1
                     ELSE 0 
                 END)),
                'No Rating'
            ) AS risk_rating_title_en,
            COALESCE(
                (SELECT TOP 1 title_ar FROM risk_matrix_risk_ratings 
                 WHERE CASE 
                     WHEN type = 'very_high' THEN 5
                     WHEN type = 'high' THEN 4
                     WHEN type = 'medium' THEN 3
                     WHEN type = 'low' THEN 2
                     WHEN type = 'very_low' THEN 1
                     ELSE 0 
                 END = MAX(CASE 
                     WHEN rmrr.type = 'very_high' THEN 5
                     WHEN rmrr.type = 'high' THEN 4
                     WHEN rmrr.type = 'medium' THEN 3
                     WHEN rmrr.type = 'low' THEN 2
                     WHEN rmrr.type = 'very_low' THEN 1
                     ELSE 0 
                 END)),
                'No Rating'
            ) AS risk_rating_title_ar,
			 COALESCE(
                (SELECT TOP 1 color_code FROM risk_matrix_risk_ratings 
                 WHERE CASE 
                     WHEN type = 'very_high' THEN 5
                     WHEN type = 'high' THEN 4
                     WHEN type = 'medium' THEN 3
                     WHEN type = 'low' THEN 2
                     WHEN type = 'very_low' THEN 1
                     ELSE 0 
                 END = MAX(CASE 
                     WHEN rmrr.type = 'very_high' THEN 5
                     WHEN rmrr.type = 'high' THEN 4
                     WHEN rmrr.type = 'medium' THEN 3
                     WHEN rmrr.type = 'low' THEN 2
                     WHEN rmrr.type = 'very_low' THEN 1
                     ELSE 0 
                 END)),
                ''
            ) AS risk_rating_color_code,
            CASE 
                WHEN COUNT(CASE WHEN ratl.type IN ('above', 'significantly_above') THEN 1 END) > 0 THEN 'Yes'
                ELSE 'No'
            END as appetite_breach
        `, highestRating)

	baseQuery := repo.db.Table("strategic_objectives so").
		Select(selectQuery).
		Joins("JOIN risk_strategic_objectives rso ON so.id = rso.strategic_objective_id").
		Joins("LEFT JOIN risk_register rr ON rr.id = rso.risk_register_id").
		Joins("LEFT JOIN risk_matrix_risk_ratings rmrr ON rmrr.id = rr.current_risk_rating_id").
		Joins(`
            LEFT JOIN risk_evaluations re
            ON re.risk_register_id = rr.id
            AND re.type = 'inherent'
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

	baseQuery = baseQuery.Group("so.id, so.title_en, so.title_ar")

	if err := baseQuery.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	offset := (params.Page - 1) * params.PageSize

	err = baseQuery.Offset(offset).Limit(params.PageSize).Find(&strategicObjectiveRisks).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &kpiDTOs.ERMStrategicObjectiveRiskOverviewDTO{
		Items: strategicObjectiveRisks,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}

	return response, nil
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

func (repo *ERMDashboardRepository) GetERMRiskSummaryByOrganizationUnit(params *kpiDTOs.ERMDashboardListQueryParamsDTO) (*kpiDTOs.ERMRiskOrganizationUnitSummaryDTO, error) {
	results := []kpiDTOs.ERMRiskOrganizationUnitSummaryItemDTO{}
	var totalCount int64

	query := repo.db.Table("risk_register rr").
		Joins("LEFT JOIN organization_units ru ON ru.id = rr.responsible_unit_id").
		Joins("LEFT JOIN organization_units parent_unit ON COALESCE(ru.organization_unit_id, ru.id) = parent_unit.id").
		Joins("JOIN organization_levels ol ON ol.id = ru.organization_level_id AND ol.is_primary = ?", true).
		Joins("LEFT JOIN risk_matrix_risk_ratings current_rmrr ON current_rmrr.id = rr.current_risk_rating_id AND current_rmrr.status = 'active'").
		Joins("LEFT JOIN risk_matrix_risk_ratings residual_rmrr ON residual_rmrr.id = rr.residual_risk_rating_id AND residual_rmrr.status = 'active'").
		Joins("LEFT JOIN risk_statuses rs ON rs.id = rr.risk_status_id AND rs.status = ?", "active").
		Joins(`LEFT JOIN (
			SELECT rap.risk_register_id, COUNT(*) AS count
			FROM risk_action_plans rap
			JOIN action_plan_statuses aps ON aps.id = rap.action_plan_status_id AND aps.status = 'active'
			WHERE rap.target_date < GETDATE() AND aps.is_closed = 0
			GROUP BY rap.risk_register_id
		) overdue_ap ON overdue_ap.risk_register_id = rr.id`).
		Where("ru.status = ?", "active").
		Where("rs.type IN ?", []string{"open", "closed"})

	// Date filtering
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where(
			"rr.reported_at BETWEEN ? AND ?",
			params.DateFrom+" 00:00:00",
			params.DateTo+" 23:59:59",
		)
	}

	// Select aggregated data
	query = query.Select(`
		parent_unit.title_en AS organization_unit_title_en,
		parent_unit.title_ar AS organization_unit_title_ar,
		COUNT(*) AS total_risks,
		-- High risks based on current risk rating
		COUNT(CASE WHEN current_rmrr.type IN ('high', 'very_high') THEN 1 END) AS high_risks,
		COUNT(CASE WHEN rs.type = 'open' THEN 1 END) AS open_risks,
		COUNT(CASE WHEN rs.type = 'closed' THEN 1 END) AS closed_risks,
		COALESCE(SUM(overdue_ap.count), 0) AS overdue_action_plans,
		-- Residual risk level based on residual risk rating (only for risks that have residual rating)
		MAX(
			CASE 
				WHEN residual_rmrr.type = 'very_high' THEN 5
				WHEN residual_rmrr.type = 'high' THEN 4
				WHEN residual_rmrr.type = 'medium' THEN 3
				WHEN residual_rmrr.type = 'low' THEN 2
				WHEN residual_rmrr.type = 'very_low' THEN 1
				ELSE 0
			END
		) AS residual_risk_score,

		-- pick corresponding fields using the same numeric logic
		MAX(CASE WHEN residual_rmrr.type = 'very_high' THEN residual_rmrr.title_en
				WHEN residual_rmrr.type = 'high' THEN residual_rmrr.title_en
				WHEN residual_rmrr.type = 'medium' THEN residual_rmrr.title_en
				WHEN residual_rmrr.type = 'low' THEN residual_rmrr.title_en
				WHEN residual_rmrr.type = 'very_low' THEN residual_rmrr.title_en
			END) AS residual_risk_rating_title_en,

		MAX(CASE WHEN residual_rmrr.type = 'very_high' THEN residual_rmrr.title_ar
				WHEN residual_rmrr.type = 'high' THEN residual_rmrr.title_ar
				WHEN residual_rmrr.type = 'medium' THEN residual_rmrr.title_ar
				WHEN residual_rmrr.type = 'low' THEN residual_rmrr.title_ar
				WHEN residual_rmrr.type = 'very_low' THEN residual_rmrr.title_ar
			END) AS residual_risk_rating_title_ar,

		MAX(CASE WHEN residual_rmrr.type = 'very_high' THEN residual_rmrr.color_code
				WHEN residual_rmrr.type = 'high' THEN residual_rmrr.color_code
				WHEN residual_rmrr.type = 'medium' THEN residual_rmrr.color_code
				WHEN residual_rmrr.type = 'low' THEN residual_rmrr.color_code
				WHEN residual_rmrr.type = 'very_low' THEN residual_rmrr.color_code
			END) AS residual_risk_rating_color_code
	`).
		Group("parent_unit.id, parent_unit.title_en, parent_unit.title_ar")

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	offset := (params.Page - 1) * params.PageSize

	err := query.Offset(offset).Limit(params.PageSize).Find(&results).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	overview := make([]kpiDTOs.ERMRiskOrganizationUnitSummaryItemDTO, len(results))
	for i, r := range results {
		overview[i] = kpiDTOs.ERMRiskOrganizationUnitSummaryItemDTO{
			OrganizationUnitTitleEn:     r.OrganizationUnitTitleEn,
			OrganizationUnitTitleAr:     r.OrganizationUnitTitleAr,
			TotalRisks:                  r.TotalRisks,
			HighRisks:                   r.HighRisks,
			OpenRisks:                   r.OpenRisks,
			ClosedRisks:                 r.ClosedRisks,
			OverdueActionPlans:          r.OverdueActionPlans,
			ResidualRiskRatingTitleEn:   r.ResidualRiskRatingTitleEn,
			ResidualRiskRatingTitleAr:   r.ResidualRiskRatingTitleAr,
			ResidualRiskRatingColorCode: r.ResidualRiskRatingColorCode,
		}
	}

	response := &kpiDTOs.ERMRiskOrganizationUnitSummaryDTO{
		Items: overview,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}

	return response, nil
}

func (repo *ERMDashboardRepository) GetERMRiskCategorySummary(params *kpiDTOs.ERMDashboardQueryParamsDTO) ([]kpiDTOs.ERMRiskCategorySummaryDTO, error) {
	results := []kpiDTOs.ERMRiskCategorySummaryDTO{}

	query := repo.db.Model(&models.RiskRegister{}).
		Select(`
			rc.title_en AS category_title_en,
			rc.title_ar AS category_title_ar, 
			COUNT(risk_register.id) AS count
		`).
		Joins("JOIN risk_categories rc ON rc.id = risk_register.risk_category_id AND rc.status = ?", "active").
		Joins("JOIN risk_statuses rs ON rs.id = risk_register.risk_status_id AND rs.status = ?", "active").
		Where("rs.type IN ?", []string{"open", "closed"}).
		Group("rc.id, rc.title_en, rc.title_ar")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("risk_register.reported_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *ERMDashboardRepository) GetERMRiskSourceSummary(params *kpiDTOs.ERMDashboardQueryParamsDTO) ([]kpiDTOs.ERMRiskSourceSummaryDTO, error) {
	results := []kpiDTOs.ERMRiskSourceSummaryDTO{}

	query := repo.db.Model(&models.RiskRegister{}).
		Select(`
			rs.title_en AS source_title_en,
			rs.title_ar AS source_title_ar, 
			COUNT(risk_register.id) AS count
		`).
		Joins("JOIN risk_sources rs ON rs.id = risk_register.risk_source_id AND rs.status = ?", "active").
		Joins("JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id AND risk_statuses.status = ?", "active").
		Where("risk_statuses.type IN ?", []string{"open", "closed"}).
		Group("rs.id, rs.title_en, rs.title_ar")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("risk_register.reported_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *ERMDashboardRepository) GetERMRiskStatusSummary(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskStatusSummaryDTO, error) {
	result := kpiDTOs.ERMRiskStatusSummaryDTO{}

	query := repo.db.Model(&models.RiskRegister{}).
		Select(`
			rs.id, 
			rs.title_en, 
			rs.title_ar, 
			rs.color_code, 
			rs.type, 
			COUNT(risk_register.id) AS count
		`).
		Joins(`
			JOIN risk_statuses rs
			ON rs.id = risk_register.risk_status_id
			AND rs.status = ?
			AND rs.type IN ?`, "active", []string{"open", "closed"})

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("risk_register.reported_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	query = query.Group("rs.id, rs.title_en, rs.title_ar, rs.color_code, rs.type")

	statuses := []kpiDTOs.ERMRiskSummaryItemDTO{}

	if err := query.Scan(&statuses).Error; err != nil {
		return nil, err
	}

	var total int64
	for _, s := range statuses {
		total += s.Count
		result.RisksByStatus = append(result.RisksByStatus, kpiDTOs.ERMRiskSummaryItemDTO{
			ID:        s.ID,
			TitleEn:   s.TitleEn,
			TitleAr:   s.TitleAr,
			ColorCode: s.ColorCode,
			Type:      s.Type,
			Count:     s.Count,
		})
	}
	result.Total = total

	return &result, nil
}

func (repo *ERMDashboardRepository) GetERMRiskOrganizationUnitSummary(params *kpiDTOs.ERMDashboardQueryParamsDTO) ([]kpiDTOs.ERMRiskSummaryByOrganizationUnitDTO, error) {
	results := []kpiDTOs.ERMRiskSummaryByOrganizationUnitDTO{}

	query := repo.db.Model(&models.RiskRegister{}).
		Joins("LEFT JOIN organization_units ru ON ru.id = risk_register.responsible_unit_id").
		Joins("LEFT JOIN organization_units parent_unit ON COALESCE(ru.organization_unit_id, ru.id) = parent_unit.id").
		Joins("JOIN organization_levels ol ON ol.id = ru.organization_level_id AND ol.is_primary = ?", true).
		Joins("JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id AND risk_statuses.status = ?", "active").
		Where("risk_statuses.type IN ?", []string{"open", "closed"}).
		Where("parent_unit.status = ?", "active").
		Select(`
			parent_unit.title_en AS organization_unit_title_en,
			parent_unit.title_ar AS organization_unit_title_ar, 
			parent_unit.color_code As organization_unit_color_code,
			COUNT(risk_register.id) AS count
		`).
		Group("parent_unit.id, parent_unit.title_en, parent_unit.title_ar, parent_unit.color_code")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("risk_register.reported_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *ERMDashboardRepository) GetERMRisksIdentifiedSummary() ([]kpiDTOs.ERMRisksIdentifiedSummaryDTO, error) {
	results := []kpiDTOs.ERMRisksIdentifiedSummaryDTO{}

	currentYear := time.Now().Year()
	startYear := currentYear - 6

	query := repo.db.Model(&models.RiskRegister{}).
		Joins("JOIN risk_statuses rs ON rs.id = risk_register.risk_status_id AND rs.status = ?", "active").
		Where("rs.type IN ?", []string{"open", "closed"}).
		Where("reported_at IS NOT NULL").
		Where("YEAR(reported_at) BETWEEN ? AND ?", startYear, currentYear).
		Select("YEAR(reported_at) AS year, COUNT(*) AS count").
		Group("YEAR(reported_at)").
		Order("YEAR(reported_at) ASC")

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *ERMDashboardRepository) GetERMRiskResponseStrategySummary(params *kpiDTOs.ERMDashboardQueryParamsDTO) ([]kpiDTOs.ERMRiskResponseStrategySummaryDTO, error) {
	results := []kpiDTOs.ERMRiskResponseStrategySummaryDTO{}

	query := repo.db.Model(&models.RiskRegister{}).
		Joins("JOIN risk_statuses rs ON rs.id = risk_register.risk_status_id AND rs.status = ?", "active").
		Joins("JOIN risk_response_strategies rrs ON rrs.id = risk_register.risk_response_strategy_id AND rrs.status = ?", "active").
		Where("rs.type IN ?", []string{"open", "closed"}).
		Select(`
			rrs.title_en AS response_strategy_title_en,
			rrs.title_ar AS response_strategy_title_ar, 
			COUNT(risk_register.id) AS count
		`).
		Group("rrs.id, rrs.title_en, rrs.title_ar").
		Order("count DESC")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("risk_register.reported_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *ERMDashboardRepository) GetERMRiskControlEfficiencySummary(params *kpiDTOs.ERMDashboardQueryParamsDTO) ([]kpiDTOs.ERMRiskControlEfficiencySummaryDTO, error) {
	var totalControls int64

	results := []kpiDTOs.ERMRiskControlEfficiencySummaryDTO{}

	// Base query for total count
	totalQuery := repo.db.Model(&models.RiskControlAnalysis{}).
		Joins("JOIN risk_controls rc ON rc.id = risk_control_analyses.risk_control_id").
		Joins("JOIN risk_register rr ON rr.id = rc.risk_register_id").
		Where("risk_control_analyses.type = ?", "inherent")

	if params.DateFrom != "" && params.DateTo != "" {
		totalQuery = totalQuery.Where("rr.reported_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	if err := totalQuery.Count(&totalControls).Error; err != nil {
		return nil, err
	}

	// Query grouped by efficiency ranges
	query := repo.db.Table("risk_control_analyses rca").
		Select(`
            CASE
                WHEN rca.score >= 80 THEN '> 80% Efficiency'
                WHEN rca.score >= 50 THEN '50% - 80% Efficiency'
                ELSE '< 50% Efficiency'
            END AS efficiency_range,
            COUNT(*) as count
        `).
		Joins("JOIN risk_controls rc ON rc.id = rca.risk_control_id").
		Joins("JOIN risk_register rr ON rr.id = rc.risk_register_id").
		Where("rca.type = ?", "inherent")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("rr.reported_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	query = query.Group("CASE WHEN rca.score >= 80 THEN '> 80% Efficiency' WHEN rca.score >= 50 THEN '50% - 80% Efficiency' ELSE '< 50% Efficiency' END")

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	// Map results
	summary := make([]kpiDTOs.ERMRiskControlEfficiencySummaryDTO, len(results))
	for i, r := range results {
		percent := 0.0
		if totalControls > 0 {
			percent = (float64(r.Percentage) / float64(totalControls)) * 100
		}
		summary[i] = kpiDTOs.ERMRiskControlEfficiencySummaryDTO{
			EfficiencyRange: r.EfficiencyRange,
			Percentage:      percent,
		}
	}

	return summary, nil
}

func (repo *ERMDashboardRepository) GetERMRiskMitigationPlanStatusSummary(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMRiskMitigationPlanStatusSummaryDTO, error) {
	results := []kpiDTOs.ERMRiskMitigationPlanStatusSummaryItemDTO{}

	query := repo.db.Table("risk_action_plans rap").
		Select(`
			aps.id,
			aps.title_en,
			aps.title_ar,
			aps.color_code,
			aps.type,
			COUNT(rap.id) AS count,
			SUM(COUNT(rap.id)) OVER() AS total
		`).
		Joins("JOIN action_plan_statuses aps ON aps.id = rap.action_plan_status_id AND aps.status = 'active'")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("risk_register.reported_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	query = query.Group("aps.id, aps.title_en, aps.title_ar, aps.color_code, aps.type").
		Order("count DESC")

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	// Map results to DTO
	summary := make([]kpiDTOs.ERMRiskSummaryItemDTO, len(results))
	var total int64
	for i, r := range results {
		summary[i] = kpiDTOs.ERMRiskSummaryItemDTO{
			ID:        r.ID,
			TitleEn:   r.TitleEn,
			TitleAr:   r.TitleAr,
			ColorCode: r.ColorCode,
			Type:      r.Type,
			Count:     r.Count,
		}
		total = r.Total
	}

	return &kpiDTOs.ERMRiskMitigationPlanStatusSummaryDTO{
		Total:                  total,
		MitigationPlanStatuses: summary,
	}, nil
}

func (repo *ERMDashboardRepository) GetERMOpenRisksSummary(params *kpiDTOs.ERMDashboardListQueryParamsDTO) (*kpiDTOs.ERMOpenRisksSummaryDTO, error) {
	results := []kpiDTOs.ERMOpenRiskSummaryItemDTO{}
	var totalCount int64

	query := repo.db.Table("risk_register rr").
		Select(`
			rr.title AS risk_title,
			rc.title_en AS category,
			parent_unit.title_en AS organization_unit,
			COUNT(DISTINCT rc_control.id) AS num_controls,
			AVG(COALESCE(rca.score, 0)) AS control_efficiency_percent,
			CASE
				WHEN AVG(COALESCE(rca.score, 0)) >= 80 THEN 'Efficient'
				WHEN AVG(COALESCE(rca.score, 0)) >= 50 THEN 'Moderate'
				ELSE 'Inefficient'
			END AS control_efficiency,
			irms.title_en AS inherent_rating,
			iein.value AS inherent_rating_value,
			rms.title_en AS current_rating,
			cein.value AS current_rating_value,
			CASE 
				WHEN MAX(CASE WHEN ratl.type = 'significantly_above' THEN 1 ELSE 0 END) = 1 THEN 
					MAX(CASE WHEN ratl.type = 'significantly_above' THEN ratl.title_en END)
				WHEN MAX(CASE WHEN ratl.type = 'above' THEN 1 ELSE 0 END) = 1 THEN 
					MAX(CASE WHEN ratl.type = 'above' THEN ratl.title_en END)
				WHEN MAX(CASE WHEN ratl.type = 'approaching' THEN 1 ELSE 0 END) = 1 THEN 
					MAX(CASE WHEN ratl.type = 'approaching' THEN ratl.title_en END)
				WHEN MAX(CASE WHEN ratl.type = 'within' THEN 1 ELSE 0 END) = 1 THEN 
					MAX(CASE WHEN ratl.type = 'within' THEN ratl.title_en END)
				WHEN MAX(CASE WHEN ratl.type = 'marginally_below' THEN 1 ELSE 0 END) = 1 THEN 
					MAX(CASE WHEN ratl.type = 'marginally_below' THEN ratl.title_en END)
				WHEN MAX(CASE WHEN ratl.type = 'below' THEN 1 ELSE 0 END) = 1 THEN 
					MAX(CASE WHEN ratl.type = 'below' THEN ratl.title_en END)
				ELSE MAX(ratl.title_en)
			END AS risk_appetite_level_title_en,
			CASE 
				WHEN MAX(CASE WHEN ratl.type = 'significantly_above' THEN 1 ELSE 0 END) = 1 THEN 
					MAX(CASE WHEN ratl.type = 'significantly_above' THEN ratl.title_ar END)
				WHEN MAX(CASE WHEN ratl.type = 'above' THEN 1 ELSE 0 END) = 1 THEN 
					MAX(CASE WHEN ratl.type = 'above' THEN ratl.title_ar END)
				WHEN MAX(CASE WHEN ratl.type = 'approaching' THEN 1 ELSE 0 END) = 1 THEN 
					MAX(CASE WHEN ratl.type = 'approaching' THEN ratl.title_ar END)
				WHEN MAX(CASE WHEN ratl.type = 'within' THEN 1 ELSE 0 END) = 1 THEN 
					MAX(CASE WHEN ratl.type = 'within' THEN ratl.title_ar END)
				WHEN MAX(CASE WHEN ratl.type = 'marginally_below' THEN 1 ELSE 0 END) = 1 THEN 
					MAX(CASE WHEN ratl.type = 'marginally_below' THEN ratl.title_ar END)
				WHEN MAX(CASE WHEN ratl.type = 'below' THEN 1 ELSE 0 END) = 1 THEN 
					MAX(CASE WHEN ratl.type = 'below' THEN ratl.title_ar END)
				ELSE MAX(ratl.title_ar)
			END AS risk_appetite_level_title_ar,
			CASE 
				WHEN MAX(CASE WHEN ratl.type = 'significantly_above' THEN 1 ELSE 0 END) = 1 THEN 
					MAX(CASE WHEN ratl.type = 'significantly_above' THEN ratl.color_code END)
				WHEN MAX(CASE WHEN ratl.type = 'above' THEN 1 ELSE 0 END) = 1 THEN 
					MAX(CASE WHEN ratl.type = 'above' THEN ratl.color_code END)
				WHEN MAX(CASE WHEN ratl.type = 'approaching' THEN 1 ELSE 0 END) = 1 THEN 
					MAX(CASE WHEN ratl.type = 'approaching' THEN ratl.color_code END)
				WHEN MAX(CASE WHEN ratl.type = 'within' THEN 1 ELSE 0 END) = 1 THEN 
					MAX(CASE WHEN ratl.type = 'within' THEN ratl.color_code END)
				WHEN MAX(CASE WHEN ratl.type = 'marginally_below' THEN 1 ELSE 0 END) = 1 THEN 
					MAX(CASE WHEN ratl.type = 'marginally_below' THEN ratl.color_code END)
				WHEN MAX(CASE WHEN ratl.type = 'below' THEN 1 ELSE 0 END) = 1 THEN 
					MAX(CASE WHEN ratl.type = 'below' THEN ratl.color_code END)
				ELSE MAX(ratl.color_code)
			END AS risk_appetite_level_color_code,
			DATEDIFF(day, rr.reported_at, GETDATE()) AS days_open
		`).
		Joins("JOIN risk_statuses rs ON rs.id = rr.risk_status_id AND rs.status = ?", "active").
		Joins("LEFT JOIN risk_categories rc ON rc.id = rr.risk_category_id AND rc.status = ?", "active").
		Joins("LEFT JOIN organization_units ru ON ru.id = rr.responsible_unit_id").
		Joins("LEFT JOIN organization_units parent_unit ON COALESCE(ru.organization_unit_id, ru.id) = parent_unit.id").
		Joins("JOIN organization_levels ol ON ol.id = ru.organization_level_id AND ol.is_primary = ?", true).
		Joins("LEFT JOIN risk_controls rc_control ON rc_control.risk_register_id = rr.id").
		Joins("LEFT JOIN risk_control_analyses rca ON rca.risk_control_id = rc_control.id AND rca.type = ?", "inherent").
		Joins("LEFT JOIN control_efficiency_statuses ces ON ces.id = rca.control_efficiency_status_id AND ces.status = ?", "active").
		Joins("LEFT JOIN risk_evaluations re ON re.risk_register_id = rr.id AND re.type = ?", "inherent").
		Joins("LEFT JOIN risk_matrix_risk_ratings irms ON irms.id = rr.inherent_risk_rating_id AND irms.status = ?", "active").
		Joins("LEFT JOIN risk_matrix_risk_rating_values iein ON iein.id = rr.inherent_risk_rating_value_id").
		Joins("LEFT JOIN risk_matrix_risk_ratings rms ON rms.id = re.risk_matrix_risk_rating_id AND rms.status = ?", "active").
		Joins("LEFT JOIN risk_matrix_risk_rating_values cein ON cein.id = re.risk_matrix_risk_rating_value_id").
		Joins("LEFT JOIN risk_appetite_alignments raa ON raa.risk_evaluation_id = re.id").
		Joins("LEFT JOIN risk_appetite_tolerance_levels ratl ON ratl.id = raa.risk_appetite_tolerance_level_id AND ratl.status = ?", "active").
		Where("rs.type = ?", "open").
		Where("parent_unit.status = ?", "active").
		Group("rr.id, rr.title, rc.title_en, parent_unit.title_en, irms.title_en, iein.value, rms.title_en, cein.value, rr.reported_at").
		Having("COUNT(DISTINCT rc_control.id) > 0")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("rr.reported_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	offset := (params.Page - 1) * params.PageSize

	err := query.Offset(offset).Limit(params.PageSize).Find(&results).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &kpiDTOs.ERMOpenRisksSummaryDTO{
		Items: results,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}

	return response, nil
}

func (repo *ERMDashboardRepository) GetTopRiskRegisters(params *kpiDTOs.ERMDashboardQueryParamsDTO) ([]*riskDTOs.RiskRegisterListDTO, error) {
	topRiskRegisters := []*riskDTOs.RiskRegisterListDTO{}

	query := repo.db.Table("risk_register").
		Select(`
            risk_register.id,
            risk_register.reference_number,
            risk_register.title,
            risk_register.description,
			CONVERT(VARCHAR, risk_register.reported_at, 23) AS reported_at,
            risk_register.created_at,
            risk_register.updated_at,
			risk_register.current_risk_rating_value_id,
            risk_categories.title_en AS risk_category_title_en,
            risk_categories.title_ar AS risk_category_title_ar,
            risk_sources.title_en AS risk_source_title_en,
            risk_sources.title_ar AS risk_source_title_ar,
			organization_units.id AS responsible_unit_id,
            organization_units.title_en AS responsible_unit_title_en,
            organization_units.title_ar AS responsible_unit_title_ar,
			users.id AS responsible_user_id,
            users.first_name_en AS responsible_user_first_name_en,
            users.first_name_ar AS responsible_user_first_name_ar,
            users.last_name_en AS responsible_user_last_name_en,
            users.last_name_ar AS responsible_user_last_name_ar,
			users.initial AS responsible_user_initial,
        	users.color AS responsible_user_color,
        	users.email AS responsible_user_email,
        	users.profile_image_path AS responsible_user_profile_image_path,
            risk_statuses.title_en AS risk_status_title_en,
            risk_statuses.title_ar AS risk_status_title_ar,
			risk_response_strategies.title_en AS risk_response_strategy_title_en,
            risk_response_strategies.title_ar AS risk_response_strategy_title_ar,
            risk_statuses.color_code AS risk_status_color_code,
			risk_statuses.type AS risk_status_type,
            review_frequencies.title_en AS review_frequency_title_en,
            review_frequencies.title_ar AS review_frequency_title_ar,
			CONVERT(VARCHAR, risk_register.next_review_date, 23) AS next_review_date,
            risk_register.last_review_note,
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
        	updated_user.profile_image_path AS updated_user_profile_image_path,
            current_likelihood.title_en AS current_likelihood_title_en,
            current_likelihood.title_ar AS current_likelihood_title_ar,
            current_impact.title_en AS current_impact_title_en,
            current_impact.title_ar AS current_impact_title_ar,
            current_risk_rating.title_en AS current_risk_rating_title_en,
            current_risk_rating.title_ar AS current_risk_rating_title_ar,
            current_risk_rating_value.value AS current_risk_rating_value,
            current_risk_rating.color_code AS current_risk_rating_color_code,
            inherent_likelihood.title_en AS inherent_likelihood_title_en,
            inherent_likelihood.title_ar AS inherent_likelihood_title_ar,
            inherent_impact.title_en AS inherent_impact_title_en,
            inherent_impact.title_ar AS inherent_impact_title_ar,
            inherent_risk_rating.title_en AS inherent_risk_rating_title_en,
            inherent_risk_rating.title_ar AS inherent_risk_rating_title_ar,
            inherent_risk_rating_value.value AS inherent_risk_rating_value,
            inherent_risk_rating.color_code AS inherent_risk_rating_color_code,
            residual_likelihood.title_en AS residual_likelihood_title_en,
            residual_likelihood.title_ar AS residual_likelihood_title_ar,
            residual_impact.title_en AS residual_impact_title_en,
            residual_impact.title_ar AS residual_impact_title_ar,
            residual_risk_rating.title_en AS residual_risk_rating_title_en,
            residual_risk_rating.title_ar AS residual_risk_rating_title_ar,
            residual_risk_rating_value.value AS residual_risk_rating_value,
            residual_risk_rating.color_code AS residual_risk_rating_color_code
        `).
		Joins(`
            LEFT JOIN risk_categories ON risk_categories.id = risk_register.risk_category_id
            LEFT JOIN risk_sources ON risk_sources.id = risk_register.risk_source_id
            LEFT JOIN organization_units ON organization_units.id = risk_register.responsible_unit_id
            LEFT JOIN users ON users.id = risk_register.responsible_user_id
            LEFT JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id
			LEFT JOIN risk_response_strategies ON risk_response_strategies.id = risk_register.risk_response_strategy_id
            LEFT JOIN users AS created_user ON created_user.id = risk_register.created_by_id
            LEFT JOIN users AS updated_user ON updated_user.id = risk_register.updated_by_id
            LEFT JOIN review_frequencies ON review_frequencies.id = risk_register.review_frequency_id
            LEFT JOIN risk_matrix_likelihoods AS current_likelihood ON current_likelihood.id = risk_register.current_likelihood_id
            LEFT JOIN risk_matrix_impacts AS current_impact ON current_impact.id = risk_register.current_impact_id
            LEFT JOIN risk_matrix_risk_ratings AS current_risk_rating ON current_risk_rating.id = risk_register.current_risk_rating_id
            LEFT JOIN risk_matrix_risk_rating_values AS current_risk_rating_value ON current_risk_rating_value.id = risk_register.current_risk_rating_value_id
            LEFT JOIN risk_matrix_likelihoods AS inherent_likelihood ON inherent_likelihood.id = risk_register.inherent_likelihood_id
            LEFT JOIN risk_matrix_impacts AS inherent_impact ON inherent_impact.id = risk_register.inherent_impact_id
            LEFT JOIN risk_matrix_risk_ratings AS inherent_risk_rating ON inherent_risk_rating.id = risk_register.inherent_risk_rating_id
            LEFT JOIN risk_matrix_risk_rating_values AS inherent_risk_rating_value ON inherent_risk_rating_value.id = risk_register.inherent_risk_rating_value_id
            LEFT JOIN risk_matrix_likelihoods AS residual_likelihood ON residual_likelihood.id = risk_register.residual_likelihood_id
            LEFT JOIN risk_matrix_impacts AS residual_impact ON residual_impact.id = risk_register.residual_impact_id
            LEFT JOIN risk_matrix_risk_ratings AS residual_risk_rating ON residual_risk_rating.id = risk_register.residual_risk_rating_id
            LEFT JOIN risk_matrix_risk_rating_values AS residual_risk_rating_value ON residual_risk_rating_value.id = risk_register.residual_risk_rating_value_id
        `).
		Order("current_risk_rating_value DESC")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("risk_register.reported_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	query = query.Limit(20)

	if err := query.Find(&topRiskRegisters).Error; err != nil {
		return nil, err
	}

	return topRiskRegisters, nil
}
