package services

import (
	"jethur/grc/kpi/dtos"
	"jethur/grc/kpi/repositories"
)

type ProjectSecurityGateDashboardServiceInterface interface {
	GetSummary(params *dtos.ProjectSecurityGateDashboardQueryParamsDTO) (*dtos.ProjectSecurityGateDashboardSummaryDTO, error)
	GetExceptionSummary(params *dtos.ProjectSecurityGateDashboardQueryParamsDTO) (*dtos.ProjectSecurityGateExceptionSummaryDTO, error)
	GetChecklistsByResponseStatus(params *dtos.ProjectSecurityGateDashboardQueryParamsDTO) (*dtos.ProjectSecurityGateChecklistsByResponseStatusDTO, error)
	GetByReadinessStatus(params *dtos.ProjectSecurityGateDashboardQueryParamsDTO) (*dtos.ProjectSecurityGatesByReadinessStatusDTO, error)
	GetByClassification(params *dtos.ProjectSecurityGateDashboardQueryParamsDTO) (*dtos.ProjectSecurityGatesByClassificationDTO, error)
}

type ProjectSecurityGateDashboardService struct {
	repo repositories.ProjectSecurityGateDashboardRepositoryInterface
}

func NewProjectSecurityGateDashboardService(repo repositories.ProjectSecurityGateDashboardRepositoryInterface) ProjectSecurityGateDashboardServiceInterface {
	return &ProjectSecurityGateDashboardService{repo: repo}
}

func (psgds *ProjectSecurityGateDashboardService) GetSummary(params *dtos.ProjectSecurityGateDashboardQueryParamsDTO) (*dtos.ProjectSecurityGateDashboardSummaryDTO, error) {
	return psgds.repo.GetSummary(params)
}

func (psgds *ProjectSecurityGateDashboardService) GetExceptionSummary(params *dtos.ProjectSecurityGateDashboardQueryParamsDTO) (*dtos.ProjectSecurityGateExceptionSummaryDTO, error) {
	return psgds.repo.GetExceptionSummary(params)
}

func (psgds *ProjectSecurityGateDashboardService) GetChecklistsByResponseStatus(params *dtos.ProjectSecurityGateDashboardQueryParamsDTO) (*dtos.ProjectSecurityGateChecklistsByResponseStatusDTO, error) {
	return psgds.repo.GetChecklistsByResponseStatus(params)
}

func (psgds *ProjectSecurityGateDashboardService) GetByReadinessStatus(params *dtos.ProjectSecurityGateDashboardQueryParamsDTO) (*dtos.ProjectSecurityGatesByReadinessStatusDTO, error) {
	return psgds.repo.GetByReadinessStatus(params)
}

func (psgds *ProjectSecurityGateDashboardService) GetByClassification(params *dtos.ProjectSecurityGateDashboardQueryParamsDTO) (*dtos.ProjectSecurityGatesByClassificationDTO, error) {
	return psgds.repo.GetByClassification(params)
}
