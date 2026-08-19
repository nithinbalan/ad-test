package services

import (
	identityAccessReviewDashboardDtos "jethur/grc/kpi/dtos"
	identityAccessReviewDashboardRepositories "jethur/grc/kpi/repositories"
)

type IdentityAccessReviewDashboardServiceInterface interface {
	GetIdentityAccessReviewSummary() (*identityAccessReviewDashboardDtos.IdentityAccessReviewSummaryDTO, error)
	GetCampaignsByLifecycle() ([]*identityAccessReviewDashboardDtos.IdentityAccessReviewCampaignByLifecycleDTO, error)
	GetCompletionBySystem() ([]*identityAccessReviewDashboardDtos.IdentityAccessReviewCompletionBySystemDTO, error)
	GetActiveAndRecentCampaigns() ([]*identityAccessReviewDashboardDtos.IdentityAccessReviewRecentCampaignDTO, error)
}

type IdentityAccessReviewDashboardService struct {
	repo identityAccessReviewDashboardRepositories.IdentityAccessReviewDashboardRepositoryInterface
}

func NewIdentityAccessReviewDashboardService(repo identityAccessReviewDashboardRepositories.IdentityAccessReviewDashboardRepositoryInterface) IdentityAccessReviewDashboardServiceInterface {
	return &IdentityAccessReviewDashboardService{repo: repo}
}

func (service *IdentityAccessReviewDashboardService) GetIdentityAccessReviewSummary() (*identityAccessReviewDashboardDtos.IdentityAccessReviewSummaryDTO, error) {
	return service.repo.GetIdentityAccessReviewSummary()
}

func (service *IdentityAccessReviewDashboardService) GetCampaignsByLifecycle() ([]*identityAccessReviewDashboardDtos.IdentityAccessReviewCampaignByLifecycleDTO, error) {
	return service.repo.GetCampaignsByLifecycle()
}

func (service *IdentityAccessReviewDashboardService) GetCompletionBySystem() ([]*identityAccessReviewDashboardDtos.IdentityAccessReviewCompletionBySystemDTO, error) {
	return service.repo.GetCompletionBySystem()
}

func (service *IdentityAccessReviewDashboardService) GetActiveAndRecentCampaigns() ([]*identityAccessReviewDashboardDtos.IdentityAccessReviewRecentCampaignDTO, error) {
	return service.repo.GetActiveAndRecentCampaigns()
}
