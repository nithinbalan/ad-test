package services

import (
	"fmt"
	"jethur/grc/kpi-register/dtos"
	"jethur/grc/kpi-register/repositories"
	"math"
)

type KPIDashboardServiceInterface interface {
	GetKPISummary(params *dtos.KPIDashboardQueryParamsDTO) (*dtos.KPISummaryDTO, error)
	GetKPIStatusSummary(params *dtos.KPIDashboardQueryParamsDTO) ([]*dtos.KPIStatusCountDTO, error)
	GetKPIPassFailSummary(params *dtos.KPIDashboardQueryParamsDTO) (*dtos.KPIPassFailSummaryResponseDTO, error)
	GetTopPerformingKPIs(params *dtos.KPIDashboardQueryParamsDTO) ([]*dtos.PerformanceSummaryDTO, error)
	GetBottomPerformingKPIs(params *dtos.KPIDashboardQueryParamsDTO) ([]*dtos.PerformanceSummaryDTO, error)
	GetCategorywiseKPIStats(params *dtos.KPIDashboardQueryParamsDTO) ([]*dtos.KPIPerformanceCategoryStatsDTO, error)
	GetDepartmentwiseKPIStats(params *dtos.KPIDashboardQueryParamsDTO) ([]*dtos.KPIDepartmentStatsDTO, error)
	GetOverdueForAchievementDetails(params *dtos.KPIDashboardQueryParamsDTO) (*dtos.KPIOverdueAchievementResponseDTO, error)
	GetOverdueForReviewDetails(params *dtos.KPIDashboardQueryParamsDTO) (*dtos.KPIOverdueReviewResponseDTO, error)

	GetKPIStatusOverview(params *dtos.KPIDashboardQueryParamsDTO) (*dtos.KPIStatusOverviewResponseDTO, error)
	GetAchievementSummary(params *dtos.KPIDashboardQueryParamsDTO) (*dtos.KPIAchievementSummaryDTO, error)
	GetKPIByReviewStatus(params *dtos.KPIDashboardQueryParamsDTO) (*dtos.KPIReviewStatusDistributionResponseDTO, error)
	GetKPIAchievementByDomain(params *dtos.KPIDashboardQueryParamsDTO) ([]*dtos.KPIAchievementByDomainDTO, error)
	GetAllKPIs(params *dtos.KPIDashboardQueryParamsDTO) (*dtos.KPIRegisterListResponseDTO, error)
	GetTopKPIs(params *dtos.KPIDashboardQueryParamsDTO) ([]dtos.KPIRegisterListDTO, error)
	GetBottomKPIs(params *dtos.KPIDashboardQueryParamsDTO) ([]dtos.KPIRegisterListDTO, error)
}

type KPIDashboardService struct {
	repo         repositories.KPIDashboardRepositoryInterface
	kpiStatusSvc KPIStatusServiceInterface
}

func NewKPIDashboardService(
	repo repositories.KPIDashboardRepositoryInterface,
	kpiStatusSvc KPIStatusServiceInterface,
) KPIDashboardServiceInterface {
	return &KPIDashboardService{
		repo:         repo,
		kpiStatusSvc: kpiStatusSvc,
	}
}

func (kds *KPIDashboardService) GetKPISummary(params *dtos.KPIDashboardQueryParamsDTO) (*dtos.KPISummaryDTO, error) {

	result, err := kds.repo.GetKPISummary(params)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return &dtos.KPISummaryDTO{}, nil
	}

	return result, nil
}

func (kds *KPIDashboardService) GetKPIStatusSummary(params *dtos.KPIDashboardQueryParamsDTO) ([]*dtos.KPIStatusCountDTO, error) {
	results, err := kds.repo.GetKPIStatusSummary(params)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (kds *KPIDashboardService) GetKPIPassFailSummary(params *dtos.KPIDashboardQueryParamsDTO) (*dtos.KPIPassFailSummaryResponseDTO, error) {
	return kds.repo.GetKPIPassFailSummary(params)
}

func (kds *KPIDashboardService) GetTopPerformingKPIs(params *dtos.KPIDashboardQueryParamsDTO) ([]*dtos.PerformanceSummaryDTO, error) {
	result, err := kds.repo.GetTopPerformingKPIs(params)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return []*dtos.PerformanceSummaryDTO{}, nil
	}
	return result, nil
}

func (kds *KPIDashboardService) GetBottomPerformingKPIs(params *dtos.KPIDashboardQueryParamsDTO) ([]*dtos.PerformanceSummaryDTO, error) {
	result, err := kds.repo.GetBottomPerformingKPIs(params)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return []*dtos.PerformanceSummaryDTO{}, nil
	}
	return result, nil
}

func (kds *KPIDashboardService) GetCategorywiseKPIStats(params *dtos.KPIDashboardQueryParamsDTO) ([]*dtos.KPIPerformanceCategoryStatsDTO, error) {
	return kds.repo.GetCategorywiseKPIStats(params)
}

func (kds *KPIDashboardService) GetDepartmentwiseKPIStats(params *dtos.KPIDashboardQueryParamsDTO) ([]*dtos.KPIDepartmentStatsDTO, error) {
	return kds.repo.GetDepartmentwiseKPIStats(params)
}

func (kds *KPIDashboardService) GetOverdueForAchievementDetails(params *dtos.KPIDashboardQueryParamsDTO) (*dtos.KPIOverdueAchievementResponseDTO, error) {

	summary, err := kds.repo.GetOverdueForAchievementSummary(params)
	if err != nil {
		return nil, err
	}

	list, err := kds.repo.GetTopOverdueForAchievementKPIs(params)
	if err != nil {
		return nil, err
	}

	if summary == nil {
		summary = &dtos.KPIOverdueSummaryDTO{}
	}

	return &dtos.KPIOverdueAchievementResponseDTO{
		TotalCount:         summary.TotalCount,
		AverageDaysOverdue: summary.AverageDaysOverdue,
		TopKPIs:            list,
	}, nil
}

func (kds *KPIDashboardService) GetOverdueForReviewDetails(params *dtos.KPIDashboardQueryParamsDTO) (*dtos.KPIOverdueReviewResponseDTO, error) {

	summary, err := kds.repo.GetOverdueReviewSummary(params)
	if err != nil {
		return nil, err
	}

	list, err := kds.repo.GetTopOverdueReviewKPIs(params)
	if err != nil {
		return nil, err
	}

	if summary == nil {
		summary = &dtos.KPIOverdueSummaryDTO{}
	}

	return &dtos.KPIOverdueReviewResponseDTO{
		TotalCount:         summary.TotalCount,
		AverageDaysOverdue: summary.AverageDaysOverdue,
		TopKPIs:            list,
	}, nil
}

func (kds *KPIDashboardService) GetKPIStatusOverview(params *dtos.KPIDashboardQueryParamsDTO) (*dtos.KPIStatusOverviewResponseDTO, error) {

	results, err := kds.repo.GetKPIStatusOverview(params)
	if err != nil {
		return nil, err
	}

	var total int64

	for _, r := range results {
		total += r.Count

		r.AverageAchievement = math.Round(r.AverageAchievement*100) / 100
		r.AverageTrend = math.Round(r.AverageTrend*100) / 100
	}

	return &dtos.KPIStatusOverviewResponseDTO{
		Total: total,
		Items: results,
	}, nil
}

func (kds *KPIDashboardService) GetAchievementSummary(params *dtos.KPIDashboardQueryParamsDTO) (*dtos.KPIAchievementSummaryDTO, error) {

	result, err := kds.repo.GetAchievementSummary(params)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (kds *KPIDashboardService) GetKPIByReviewStatus(params *dtos.KPIDashboardQueryParamsDTO) (*dtos.KPIReviewStatusDistributionResponseDTO, error) {

	results, err := kds.repo.GetKPIByReviewStatus(params)
	if err != nil {
		return nil, err
	}

	var total int64
	for _, r := range results {
		total += r.Count
	}

	return &dtos.KPIReviewStatusDistributionResponseDTO{
		Total: total,
		Items: results,
	}, nil
}

func (kds *KPIDashboardService) GetKPIAchievementByDomain(params *dtos.KPIDashboardQueryParamsDTO) ([]*dtos.KPIAchievementByDomainDTO, error) {

	results, err := kds.repo.GetKPIAchievementByDomain(params)
	if err != nil {
		return nil, err
	}

	statusList, err := kds.kpiStatusSvc.GetAll("", "", "", "active")
	if err != nil {
		return nil, err
	}

	statusMap := make(map[string]*dtos.KPIStatusDTO)

	for _, s := range statusList {
		statusMap[s.Type] = &dtos.KPIStatusDTO{
			ID:            s.ID,
			TitleEn:       s.TitleEn,
			TitleAr:       s.TitleAr,
			Type:          s.Type,
			ColorCode:     s.ColorCode,
			Status:        s.Status,
			Label:         s.Label,
			Order:         s.Order,
			IsNotAchieved: s.IsNotAchieved,
		}
	}

	for _, r := range results {

		achievement := r.AverageScore

		var statusType string

		switch {
		case achievement > 100:
			statusType = "exceeded"

		case math.Abs(achievement-100) < 0.0001:
			statusType = "achieved"

		case achievement >= 90:
			statusType = "near-target"

		case achievement >= 75:
			statusType = "needs-improvement"

		default:
			statusType = "low"
		}

		status, ok := statusMap[statusType]
		if !ok {
			return nil, fmt.Errorf("missing KPI status for type: %s", statusType)
		}

		r.KPIStatusID = status.ID
		r.KPIStatusTitleEn = status.TitleEn
		r.KPIStatusTitleAr = status.TitleAr
		r.KPIStatusType = status.Type
		r.KPIStatusColorCode = status.ColorCode
		r.KPIStatusLabel = status.Label
		r.KPIStatusOrder = status.Order
		r.KPIStatusIsNotAchieved = status.IsNotAchieved

		r.AverageScore = math.Round(achievement*100) / 100
	}

	return results, nil
}
func (kds *KPIDashboardService) GetAllKPIs(params *dtos.KPIDashboardQueryParamsDTO) (*dtos.KPIRegisterListResponseDTO, error) {
	return kds.repo.GetAllKPIs(params)
}

func (kds *KPIDashboardService) GetTopKPIs(params *dtos.KPIDashboardQueryParamsDTO) ([]dtos.KPIRegisterListDTO, error) {
	return kds.repo.GetTopKPIs(params)
}

func (kds *KPIDashboardService) GetBottomKPIs(params *dtos.KPIDashboardQueryParamsDTO) ([]dtos.KPIRegisterListDTO, error) {
	return kds.repo.GetBottomKPIs(params)
}
