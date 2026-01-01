package services

import (
	"fmt"
	tprmDTOs "jethur/grc/tprm/dtos"
	"jethur/grc/tprm/repositories"
	"math"
	"time"
)

type TPRMDashboardServiceInterface interface {
	GetThirdPartyInsights(params *tprmDTOs.TPRMDashboardQueryParamsDTO) (*tprmDTOs.TPRMThirdPartyInsightsResponseDTO, error)
	GetAssessmentInsights(params *tprmDTOs.TPRMDashboardQueryParamsDTO) (*tprmDTOs.TPRMAssessmentInsightsResponseDTO, error)
	GetSystemStatistics(params *tprmDTOs.TPRMDashboardQueryParamsDTO) (*tprmDTOs.TPRMSystemStatisticsDTO, error)
	GetOngoingAssessments(params *tprmDTOs.TPRMDashboardQueryParamsDTO) ([]tprmDTOs.TPRMOngoingAssessmentDTO, error)
	GetThirdPartyServiceExpiry(params *tprmDTOs.TPRMDashboardQueryParamsDTO) ([]tprmDTOs.TPRMServiceExpiryDTO, error)
	GetAssessmentByTemplateType(params *tprmDTOs.TPRMDashboardQueryParamsDTO) ([]tprmDTOs.TPRMAssessmentByTemplateDTO, error)
	GetThirdPartyByType(params *tprmDTOs.TPRMDashboardQueryParamsDTO) ([]tprmDTOs.TPRMThirdPartyByTypeDTO, error)
	GetAssessmentRiskMonthlyTrend(params *tprmDTOs.TPRMDashboardQueryParamsDTO) ([]tprmDTOs.TPRMAssessmentRiskMonthlyTrendDTO, error)
}

type TPRMDashboardService struct {
	repo repositories.TPRMDashboardRepositoryInterface
}

func NewTPRMDashboardService(repo repositories.TPRMDashboardRepositoryInterface) TPRMDashboardServiceInterface {
	return &TPRMDashboardService{
		repo: repo,
	}
}

func (tds *TPRMDashboardService) GetThirdPartyInsights(params *tprmDTOs.TPRMDashboardQueryParamsDTO) (*tprmDTOs.TPRMThirdPartyInsightsResponseDTO, error) {

	progress, err := tds.repo.GetThirdPartyAssessmentProgress(params)
	if err != nil {
		return nil, err
	}

	riskLevels, err := tds.repo.GetThirdPartyRiskLevelDistribution(params)
	if err != nil {
		return nil, err
	}

	var riskTotal int64 = 0
	for _, r := range riskLevels {
		riskTotal += r.Count
	}

	if riskTotal > 0 {
		for i := range riskLevels {
			percent := (float64(riskLevels[i].Count) / float64(riskTotal)) * 100
			riskLevels[i].Percentage = math.Round(percent)
		}
	}

	return &tprmDTOs.TPRMThirdPartyInsightsResponseDTO{
		AssessmentProgress: *progress,
		RiskLevels:         riskLevels,
	}, nil
}

func (tds *TPRMDashboardService) GetAssessmentInsights(params *tprmDTOs.TPRMDashboardQueryParamsDTO) (*tprmDTOs.TPRMAssessmentInsightsResponseDTO, error) {

	insights, err := tds.repo.GetAssessmentInsights(params)
	if err != nil {
		return nil, err
	}

	if insights.TotalAssessments > 0 {
		insights.CompletionRate = math.Round(
			(float64(insights.Completed) / float64(insights.TotalAssessments)) * 100,
		)
	} else {
		insights.CompletionRate = 0
	}

	riskItems, err := tds.repo.GetAssessmentRiskLevelDistribution(params)
	if err != nil {
		return nil, err
	}

	var riskTotal int64 = 0
	for _, item := range riskItems {
		riskTotal += item.Count
	}

	if riskTotal > 0 {
		for i := range riskItems {
			riskItems[i].Percentage = math.Round(
				(float64(riskItems[i].Count) / float64(riskTotal)) * 100,
			)
		}
	} else {
		for i := range riskItems {
			riskItems[i].Percentage = 0
		}
	}

	resp := &tprmDTOs.TPRMAssessmentInsightsResponseDTO{
		Insights: *insights,
		RiskRatingDistribution: tprmDTOs.TPRMRiskRatingDistributionDTO{
			Total: riskTotal,
			Items: riskItems,
		},
	}

	return resp, nil
}

func (tds *TPRMDashboardService) GetSystemStatistics(params *tprmDTOs.TPRMDashboardQueryParamsDTO) (*tprmDTOs.TPRMSystemStatisticsDTO, error) {
	return tds.repo.GetSystemStatistics(params)
}

func (tds *TPRMDashboardService) GetOngoingAssessments(params *tprmDTOs.TPRMDashboardQueryParamsDTO) ([]tprmDTOs.TPRMOngoingAssessmentDTO, error) {

	items, err := tds.repo.GetOngoingAssessments(params)
	if err != nil {
		return nil, err
	}

	today := time.Now().Truncate(24 * time.Hour)

	for i := range items {

		due := items[i].DueDate.Truncate(24 * time.Hour)

		diffDays := int(due.Sub(today).Hours() / 24)

		switch {
		case diffDays < 0:
			items[i].DaysLeft = fmt.Sprintf("Overdue (%d days)", -diffDays)
		case diffDays == 0:
			items[i].DaysLeft = "Due Today"
		case diffDays == 1:
			items[i].DaysLeft = "1 day left"
		default:
			items[i].DaysLeft = fmt.Sprintf("%d days left", diffDays)
		}
	}

	return items, nil
}

func (tds *TPRMDashboardService) GetThirdPartyServiceExpiry(params *tprmDTOs.TPRMDashboardQueryParamsDTO) ([]tprmDTOs.TPRMServiceExpiryDTO, error) {

	counts, err :=
		tds.repo.GetThirdPartyServiceExpiryCounts(params)
	if err != nil {
		return nil, err
	}

	resp := []tprmDTOs.TPRMServiceExpiryDTO{
		{
			Label: "Due within 7 d",
			Count: counts.DueWithin7d,
		},
		{
			Label: "Overdue (less than 30d)",
			Count: counts.OverdueLt30d,
		},
		{
			Label: "Overdue (30+ d)",
			Count: counts.OverdueGt30d,
		},
	}

	return resp, nil
}

func (tds *TPRMDashboardService) GetAssessmentByTemplateType(params *tprmDTOs.TPRMDashboardQueryParamsDTO) ([]tprmDTOs.TPRMAssessmentByTemplateDTO, error) {
	return tds.repo.GetAssessmentByTemplateType(params)
}

func (tds *TPRMDashboardService) GetThirdPartyByType(params *tprmDTOs.TPRMDashboardQueryParamsDTO) ([]tprmDTOs.TPRMThirdPartyByTypeDTO, error) {
	return tds.repo.GetThirdPartyByType(params)
}

func (tds *TPRMDashboardService) GetAssessmentRiskMonthlyTrend(params *tprmDTOs.TPRMDashboardQueryParamsDTO) ([]tprmDTOs.TPRMAssessmentRiskMonthlyTrendDTO, error) {

	var finalMonth time.Time
	if params != nil && params.DateTo != "" {
		parsed, err := time.Parse("2006-01-02", params.DateTo)
		if err != nil {
			return nil, err
		}
		finalMonth = time.Date(parsed.Year(), parsed.Month(), 1, 0, 0, 0, 0, time.UTC)
	} else {
		now := time.Now()
		finalMonth = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	}

	startMonth := finalMonth.AddDate(0, -11, 0)
	endMonth := finalMonth.AddDate(0, 1, 0)

	rows, err := tds.repo.GetAssessmentRiskMonthlyCounts(
		startMonth,
		endMonth,
		params,
	)
	if err != nil {
		return nil, err
	}

	monthMap := map[string]*tprmDTOs.TPRMAssessmentRiskMonthlyTrendDTO{}

	for i := 0; i < 12; i++ {
		monthTime := startMonth.AddDate(0, i, 0)
		monthKey := monthTime.Format("2006-01")
		monthLabel := monthTime.Format("Jan 2006")

		monthMap[monthKey] = &tprmDTOs.TPRMAssessmentRiskMonthlyTrendDTO{
			Month:  monthLabel,
			High:   0,
			Medium: 0,
			Low:    0,
		}
	}

	for _, r := range rows {
		if entry, ok := monthMap[r.Month]; ok {
			switch r.RiskType {
			case "high":
				entry.High = r.Count
			case "medium":
				entry.Medium = r.Count
			case "low":
				entry.Low = r.Count
			}
		}
	}

	results := make([]tprmDTOs.TPRMAssessmentRiskMonthlyTrendDTO, 0, 12)

	for i := 0; i < 12; i++ {
		monthKey := startMonth.AddDate(0, i, 0).Format("2006-01")
		results = append(results, *monthMap[monthKey])
	}

	return results, nil
}
