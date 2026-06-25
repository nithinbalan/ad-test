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

func (ads *AwarenessDashboardService) GetAwarenessPlanSummary(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) (*awarenessDTOs.AwarenessPlanSummaryDTO, error) {
	return ads.repo.GetAwarenessPlanSummary(params)
}

func (ads *AwarenessDashboardService) GetCampaignsByType(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) ([]*awarenessDTOs.AwarenessCampaignsByTypeDTO, error) {
	return ads.repo.GetCampaignsByType(params)
}

func (ads *AwarenessDashboardService) GetPriorityDistribution(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) ([]*awarenessDTOs.AwarenessPriorityDistributionDTO, error) {
	return ads.repo.GetPriorityDistribution(params)
}

func (ads *AwarenessDashboardService) GetCampaignsPerMonth(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) ([]*awarenessDTOs.AwarenessCampaignsPerMonthDTO, error) {
	return ads.repo.GetCampaignsPerMonth(params)
}

func (ads *AwarenessDashboardService) GetChannelMixPerMonth(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) ([]*awarenessDTOs.AwarenessChannelMixPerMonthDTO, error) {
	return ads.repo.GetChannelMixPerMonth(params)
}

func (ads *AwarenessDashboardService) GetCoverageByAudience(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) ([]*awarenessDTOs.AwarenessCoverageByAudienceDTO, error) {
	return ads.repo.GetCoverageByAudience(params)
}

func (ads *AwarenessDashboardService) GetAwarenessProgressSummary(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) (*awarenessDTOs.AwarenessProgressSummaryDTO, error) {
	return ads.repo.GetAwarenessProgressSummary(params)
}

func (ads *AwarenessDashboardService) GetCompletionByChannel(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) ([]*awarenessDTOs.AwarenessCompletionByChannelDTO, error) {
	return ads.repo.GetCompletionByChannel(params)
}

func (ads *AwarenessDashboardService) GetCompletionPerCampaign(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) ([]*awarenessDTOs.AwarenessCompletionPerCampaignDTO, error) {
	return ads.repo.GetCompletionPerCampaign(params)
}

func (ads *AwarenessDashboardService) GetSurveyCompletionStatus(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) (*awarenessDTOs.AwarenessSurveyCompletionDTO, error) {
	return ads.repo.GetSurveyCompletionStatus(params)
}

func (ads *AwarenessDashboardService) GetCompletionByDepartment(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) ([]*awarenessDTOs.AwarenessCompletionByDepartmentDTO, error) {
	return ads.repo.GetCompletionByDepartment(params)
}

func (ads *AwarenessDashboardService) GetPhishingResults(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) (*awarenessDTOs.AwarenessPhishingResultsDTO, error) {
	return ads.repo.GetPhishingResults(params)
}