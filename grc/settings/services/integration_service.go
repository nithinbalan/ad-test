package services

import (
	activityLogServices "jethur/grc/logs/services"
	integrationCategoryServices "jethur/grc/masters/services"
	integrationDTOs "jethur/grc/settings/dtos"
	"jethur/grc/settings/repositories"
)

type IntegrationServiceInterface interface {
	GetAll(params *integrationDTOs.IntegrationListQueryParamsDTO) ([]*integrationDTOs.IntegrationDTO, error)
}

// Service and implementations
type IntegrationService struct {
	repo                   repositories.IntegrationRepositoryInterface
	activityLogSvc         activityLogServices.ActivityLogServiceInterface
	integrationCategorySvc integrationCategoryServices.IntegrationCategoryServiceInterface
}

func NewIntegrationService(repo repositories.IntegrationRepositoryInterface, activityLogSvc activityLogServices.ActivityLogServiceInterface, integrationCategorySvc integrationCategoryServices.IntegrationCategoryServiceInterface) IntegrationServiceInterface {
	return &IntegrationService{
		repo:                   repo,
		activityLogSvc:         activityLogSvc,
		integrationCategorySvc: integrationCategorySvc,
	}
}

func (is *IntegrationService) GetAll(params *integrationDTOs.IntegrationListQueryParamsDTO) ([]*integrationDTOs.IntegrationDTO, error) {
	return is.repo.GetAll(params)
}
