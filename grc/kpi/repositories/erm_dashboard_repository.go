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
		Joins("JOIN risk_matrix_risk_ratings ON risk_register.current_risk_rating_id = risk_matrix_risk_ratings.id")

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

		// Build base query for this specific rating
		baseQuery := repo.db.Model(&models.RiskRegister{}).
			Where("current_risk_rating_id = ?", rating.ID)

		if params.DateFrom != "" && params.DateTo != "" {
			baseQuery = baseQuery.Where("reported_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
		}

		// Total count for this risk rating
		err := baseQuery.Count(&count).Error
		if err != nil {
			return nil, err
		}

		// Open count for this risk rating
		openQuery := baseQuery.
			Joins("JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id").
			Where("risk_statuses.type = ?", "open")

		err = openQuery.Count(&openCount).Error
		if err != nil {
			return nil, err
		}

		// Calculate percentage
		var percentage float64
		if totalRisksRegistered > 0 {
			percentage = math.Round(float64(count) / float64(totalRisksRegistered) * 100)
		}

		// Populate risk rating details
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

	// Base query
	baseQuery := repo.db.Model(&models.RiskRegister{}).
		Select("risk_matrix_risk_ratings.type AS risk_level, COUNT(*) AS count").
		Joins("JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id").
		Joins("JOIN risk_matrix_risk_ratings ON risk_matrix_risk_ratings.id = risk_register.current_risk_rating_id").
		Where("next_review_date < ?", currentTime).
		Where("next_review_date IS NOT NULL").
		Where("risk_statuses.type = ?", "open").
		Group("risk_matrix_risk_ratings.type")

	if params.DateFrom != "" && params.DateTo != "" {
		baseQuery = baseQuery.Where(
			"reported_at BETWEEN ? AND ?",
			params.DateFrom+" 00:00:00",
			params.DateTo+" 23:59:59",
		)
	}

	type result struct {
		RiskLevel string
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
		var level string
		switch r.RiskLevel {
		case "high", "very_high":
			level = "High"
		case "medium":
			level = "Moderate"
		case "low", "very_low":
			level = "Low"
		default:
			continue
		}

		risksOverdueDetails = append(risksOverdueDetails, kpiDTOs.ERMOverdueRiskDetailDTO{
			RiskLevel: level,
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
	// Build the base query with all necessary joins
	baseQuery := repo.db.Model(&models.RiskRegister{}).
		Joins("JOIN risk_evaluations ON risk_evaluations.risk_register_id = risk_register.id").
		Joins("JOIN risk_appetite_alignments ON risk_appetite_alignments.risk_evaluation_id = risk_evaluations.id").
		Joins("JOIN risk_appetite_tolerance_levels ON risk_appetite_tolerance_levels.id = risk_appetite_alignments.risk_appetite_tolerance_level_id").
		Joins("LEFT JOIN risk_matrix_risk_ratings ON risk_matrix_risk_ratings.id = risk_register.current_risk_rating_id").
		Where("risk_appetite_tolerance_levels.type = ?", "above")

	// Apply date filter inline if provided
	if params.DateFrom != "" && params.DateTo != "" {
		baseQuery = baseQuery.Where("risk_register.reported_at BETWEEN ? AND ?",
			params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	var result struct {
		Total    int64 `gorm:"column:total"`
		High     int64 `gorm:"column:high"`
		Moderate int64 `gorm:"column:moderate"`
		Low      int64 `gorm:"column:low"`
	}

	err := baseQuery.Select(`
        COUNT(*) as total,
        COUNT(CASE WHEN risk_matrix_risk_ratings.type IN ('high', 'very_high') THEN 1 END) as high,
        COUNT(CASE WHEN risk_matrix_risk_ratings.type = 'medium' THEN 1 END) as moderate,
        COUNT(CASE WHEN risk_matrix_risk_ratings.type IN ('low', 'very_low') THEN 1 END) as low
    `).Scan(&result).Error

	if err != nil {
		return nil, err
	}

	// Build risks by level array
	risksByLevels := []kpiDTOs.ERMRisksAboveAppetiteLevelDTO{}

	if result.High > 0 {
		risksByLevels = append(risksByLevels, kpiDTOs.ERMRisksAboveAppetiteLevelDTO{
			RiskLevel: "High",
			Count:     result.High,
		})
	}

	if result.Moderate > 0 {
		risksByLevels = append(risksByLevels, kpiDTOs.ERMRisksAboveAppetiteLevelDTO{
			RiskLevel: "Moderate",
			Count:     result.Moderate,
		})
	}

	if result.Low > 0 {
		risksByLevels = append(risksByLevels, kpiDTOs.ERMRisksAboveAppetiteLevelDTO{
			RiskLevel: "Low",
			Count:     result.Low,
		})
	}

	return &kpiDTOs.ERMRisksAboveAppetiteDTO{
		TotalRisksAboveAppetite:    result.Total,
		AboveAppetiteRisksByLevels: risksByLevels,
	}, nil
}

func (repo *ERMDashboardRepository) GetERMLeadershipAttentionMetrics(params *kpiDTOs.ERMDashboardQueryParamsDTO) (*kpiDTOs.ERMLeadershipAttentionDTO, error) {
	// Strategic objectives at risk
	var strategicObjectivesAtRisk int64
	strategicQuery := repo.db.Model(&models.RiskRegister{}).
		Joins("JOIN risk_strategic_objectives ON risk_strategic_objectives.risk_register_id = risk_register.id").
		Joins("JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id").
		Where("risk_statuses.type = ?", "open")

	if params.DateFrom != "" && params.DateTo != "" {
		strategicQuery = strategicQuery.Where("risk_register.reported_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	if err := strategicQuery.Count(&strategicObjectivesAtRisk).Error; err != nil {
		return nil, err
	}

	// Critical risks uncontrolled
	var criticalRisksUncontrolled int64
	criticalQuery := repo.db.Model(&models.RiskRegister{}).
		Joins("JOIN risk_matrix_risk_ratings ON risk_matrix_risk_ratings.id = risk_register.current_risk_rating_id").
		Joins("JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id").
		Where("risk_matrix_risk_ratings.type IN ?", []string{"high", "very_high"}).
		Where("risk_statuses.type = ?", "open").
		Where(`NOT EXISTS (
            SELECT 1 FROM risk_controls
            JOIN controls ON controls.id = risk_controls.control_id
            JOIN control_efficiency_statuses control_efficiency_statuses ON control_efficiency_statuses.id = controls.control_efficiency_status_id
            WHERE risk_controls.risk_register_id = risk_register.id
            AND control_efficiency_statuses.type = 'efficient'
        )`)

	if params.DateFrom != "" && params.DateTo != "" {
		criticalQuery = criticalQuery.Where("risk_register.reported_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	if err := criticalQuery.Count(&criticalRisksUncontrolled).Error; err != nil {
		return nil, err
	}

	// Overdue high-risk actions
	var overdueHighRiskActions int64
	currentTime := time.Now()

	actionQuery := repo.db.Model(&models.RiskActionPlan{}).
		Joins("JOIN risk_register ON risk_register.id = risk_action_plans.risk_register_id").
		Joins("JOIN risk_matrix_risk_ratings ON risk_matrix_risk_ratings.id = risk_register.current_risk_rating_id").
		Joins("JOIN action_plan_statuses ON action_plan_statuses.id = risk_action_plans.action_plan_status_id").
		Where("risk_matrix_risk_ratings.type IN ?", []string{"high", "very_high"}).
		Where("action_plan_statuses.type != ?", "completed").
		Where("risk_action_plans.target_date < ?", currentTime)

	if params.DateFrom != "" && params.DateTo != "" {
		actionQuery = actionQuery.Where("risk_register.reported_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	if err := actionQuery.Count(&overdueHighRiskActions).Error; err != nil {
		return nil, err
	}

	// Compliance-linked high risks
	var complianceLinkedHighRisks int64
	complianceQuery := repo.db.Model(&models.RiskRegister{}).
		Joins("JOIN risk_compliance_requirements ON risk_compliance_requirements.risk_register_id = risk_register.id").
		Joins("JOIN risk_matrix_risk_ratings ON risk_matrix_risk_ratings.id = risk_register.current_risk_rating_id").
		Joins("JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id").
		Where("risk_matrix_risk_ratings.type IN ?", []string{"high", "very_high"}).
		Where("risk_statuses.type = ?", "open")

	if params.DateFrom != "" && params.DateTo != "" {
		complianceQuery = complianceQuery.Where("risk_register.reported_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	if err := complianceQuery.Count(&complianceLinkedHighRisks).Error; err != nil {
		return nil, err
	}

	return &kpiDTOs.ERMLeadershipAttentionDTO{
		StrategicObjectivesAtRisk: strategicObjectivesAtRisk,
		CriticalRisksUncontrolled: criticalRisksUncontrolled,
		OverdueHighRiskActions:    overdueHighRiskActions,
		ComplianceLinkedHighRisks: complianceLinkedHighRisks,
	}, nil
}
