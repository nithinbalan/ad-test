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

	totalRiskQuery := repo.db.Model(&models.RiskRegister{}).
		Joins("JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id").
		Where("risk_statuses.type IN (?)", []string{"open", "closed"})

	if params.DateFrom != "" && params.DateTo != "" {
		totalRiskQuery = totalRiskQuery.Where("reported_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	// Total risks
	err := totalRiskQuery.Count(&totalRisks).Error
	if err != nil {
		return nil, err
	}

	// Open risks
	openRiskQuery := repo.db.Model(&models.RiskRegister{}).
		Joins("JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id").
		Where("risk_statuses.type = ?", "open")
	if params.DateFrom != "" && params.DateTo != "" {
		openRiskQuery = openRiskQuery.Where("risk_register.reported_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}
	err = openRiskQuery.Count(&openRisks).Error
	if err != nil {
		return nil, err
	}

	// Closed risks
	closedRiskQuery := repo.db.Model(&models.RiskRegister{}).
		Joins("JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id").
		Where("risk_statuses.type = ?", "closed")
	if params.DateFrom != "" && params.DateTo != "" {
		closedRiskQuery = closedRiskQuery.Where("risk_register.reported_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}
	err = closedRiskQuery.Count(&closedRisks).Error
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
		Where("risk_statuses.type != ?", "closed")
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

	// Base query for overdue risks
	overdueQuery := repo.db.Model(&models.RiskRegister{}).
		Where("next_review_date < ?", currentTime).
		Where("next_review_date IS NOT NULL")

	if params.DateFrom != "" && params.DateTo != "" {
		overdueQuery = overdueQuery.Where("reported_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	// Total overdue count
	var totalOverdue int64
	err := overdueQuery.Count(&totalOverdue).Error
	if err != nil {
		return nil, err
	}

	// Count overdue risks by risk level
	var overdueDetails []kpiDTOs.ERMOverdueRiskDetailDTO
	err = overdueQuery.
		Joins("JOIN risk_matrix_risk_ratings ON risk_matrix_risk_ratings.id = risk_register.current_risk_rating_id").
		Select(`
            CASE 
                WHEN risk_matrix_risk_ratings.type IN ('high', 'very_high') THEN 'High'
                WHEN risk_matrix_risk_ratings.type = 'medium' THEN 'Moderate'
                WHEN risk_matrix_risk_ratings.type IN ('low', 'very_low') THEN 'Low'
                ELSE 'Other'
            END as risk_level,
            COUNT(*) as count
        `).
		Group("risk_level").
		Scan(&overdueDetails).Error
	if err != nil {
		return nil, err
	}

	return &kpiDTOs.ERMOverdueRisksDTO{
		TotalOverdue:  totalOverdue,
		OverdueByRisk: overdueDetails,
	}, nil
}
