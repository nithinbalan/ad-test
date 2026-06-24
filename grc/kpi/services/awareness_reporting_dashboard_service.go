// awareness_plan_dashboard_service.go

package services

import (
	awarenessDTOs "jethur/grc/kpi/dtos"
	awarenessRepositories "jethur/grc/kpi/repositories"
)

type AwarenessDashboardServiceInterface interface {
	GetAwarenessPlanSummary(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) (*awarenessDTOs.AwarenessPlanSummaryDTO, error)
	GetCampaignsByType(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) ([]*awarenessDTOs.AwarenessCampaignsByTypeDTO, error)
	GetPriorityDistribution(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) ([]*awarenessDTOs.AwarenessPriorityDistributionDTO, error)
	GetCampaignsPerMonth(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) ([]*awarenessDTOs.AwarenessCampaignsPerMonthDTO, error)
	GetChannelMixPerMonth(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) ([]*awarenessDTOs.AwarenessChannelMixPerMonthDTO, error)
	GetCoverageByAudience(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) ([]*awarenessDTOs.AwarenessCoverageByAudienceDTO, error)

	GetAwarenessProgressSummary(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) (*awarenessDTOs.AwarenessProgressSummaryDTO, error)
	GetCompletionByChannel(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) ([]*awarenessDTOs.AwarenessCompletionByChannelDTO, error)
	GetCompletionPerCampaign(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) ([]*awarenessDTOs.AwarenessCompletionPerCampaignDTO, error)
	GetSurveyCompletionStatus(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) (*awarenessDTOs.AwarenessSurveyCompletionDTO, error)
	GetCompletionByDepartment(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) ([]*awarenessDTOs.AwarenessCompletionByDepartmentDTO, error)
	GetPhishingResults(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) (*awarenessDTOs.AwarenessPhishingResultsDTO, error)
}

type AwarenessDashboardService struct {
	repo awarenessRepositories.AwarenessDashboardRepositoryInterface
}

func NewAwarenessDashboardService(repo awarenessRepositories.AwarenessDashboardRepositoryInterface) AwarenessDashboardServiceInterface {
	return &AwarenessDashboardService{repo: repo}
}

func (s *AwarenessDashboardService) GetAwarenessPlanSummary(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) (*awarenessDTOs.AwarenessPlanSummaryDTO, error) {
	return s.repo.GetAwarenessPlanSummary(params)
}

func (s *AwarenessDashboardService) GetCampaignsByType(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) ([]*awarenessDTOs.AwarenessCampaignsByTypeDTO, error) {
	return s.repo.GetCampaignsByType(params)
}

func (s *AwarenessDashboardService) GetPriorityDistribution(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) ([]*awarenessDTOs.AwarenessPriorityDistributionDTO, error) {
	return s.repo.GetPriorityDistribution(params)
}

func (s *AwarenessDashboardService) GetCampaignsPerMonth(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) ([]*awarenessDTOs.AwarenessCampaignsPerMonthDTO, error) {
	return s.repo.GetCampaignsPerMonth(params)
}

func (s *AwarenessDashboardService) GetChannelMixPerMonth(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) ([]*awarenessDTOs.AwarenessChannelMixPerMonthDTO, error) {
	return s.repo.GetChannelMixPerMonth(params)
}

func (s *AwarenessDashboardService) GetCoverageByAudience(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) ([]*awarenessDTOs.AwarenessCoverageByAudienceDTO, error) {
	return s.repo.GetCoverageByAudience(params)
}

func (s *AwarenessDashboardService) GetAwarenessProgressSummary(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) (*awarenessDTOs.AwarenessProgressSummaryDTO, error) {
	return s.repo.GetAwarenessProgressSummary(params)
}

func (s *AwarenessDashboardService) GetCompletionByChannel(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) ([]*awarenessDTOs.AwarenessCompletionByChannelDTO, error) {
	return s.repo.GetCompletionByChannel(params)
}

func (s *AwarenessDashboardService) GetCompletionPerCampaign(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) ([]*awarenessDTOs.AwarenessCompletionPerCampaignDTO, error) {
	return s.repo.GetCompletionPerCampaign(params)
}

func (s *AwarenessDashboardService) GetSurveyCompletionStatus(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) (*awarenessDTOs.AwarenessSurveyCompletionDTO, error) {
	return s.repo.GetSurveyCompletionStatus(params)
}

func (s *AwarenessDashboardService) GetCompletionByDepartment(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) ([]*awarenessDTOs.AwarenessCompletionByDepartmentDTO, error) {
	return s.repo.GetCompletionByDepartment(params)
}

func (s *AwarenessDashboardService) GetPhishingResults(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) (*awarenessDTOs.AwarenessPhishingResultsDTO, error) {
	return s.repo.GetPhishingResults(params)
}