package services

import (
	projectDTOs "jethur/grc/kpi/dtos"
	projectRepositories "jethur/grc/kpi/repositories"
)

type ProjectDashboardServiceInterface interface {
	GetProjectSummary(params *projectDTOs.ProjectDashboardQueryParamsDTO) (*projectDTOs.ProjectDashboardSummaryDTO, error)
	GetProjectsByPhase(params *projectDTOs.ProjectDashboardQueryParamsDTO) (*projectDTOs.ProjectsByPhaseDTO, error)
	GetProjectsByClassification(params *projectDTOs.ProjectDashboardQueryParamsDTO) (*projectDTOs.ProjectsByClassificationDTO, error)
	GetProjectsByStatus(params *projectDTOs.ProjectDashboardQueryParamsDTO) (*projectDTOs.ProjectsByStatusDTO, error)
	GetProjectCompletion(params *projectDTOs.ProjectDashboardQueryParamsDTO) ([]projectDTOs.ProjectCompletionListDTO, error)
	GetProjectsDueSoonAndOverdue(params *projectDTOs.ProjectDashboardQueryParamsDTO) ([]projectDTOs.ProjectDueSoonOverdueListDTO, error)
}

type ProjectDashboardService struct {
	repo projectRepositories.ProjectDashboardRepositoryInterface
}

func NewProjectDashboardService(repo projectRepositories.ProjectDashboardRepositoryInterface) ProjectDashboardServiceInterface {
	return &ProjectDashboardService{repo: repo}
}

func (pds *ProjectDashboardService) GetProjectSummary(params *projectDTOs.ProjectDashboardQueryParamsDTO) (*projectDTOs.ProjectDashboardSummaryDTO, error) {
	return pds.repo.GetProjectSummary(params)
}

func (pds *ProjectDashboardService) GetProjectsByPhase(params *projectDTOs.ProjectDashboardQueryParamsDTO) (*projectDTOs.ProjectsByPhaseDTO, error) {
	return pds.repo.GetProjectsByPhase(params)
}

func (pds *ProjectDashboardService) GetProjectsByClassification(params *projectDTOs.ProjectDashboardQueryParamsDTO) (*projectDTOs.ProjectsByClassificationDTO, error) {
	return pds.repo.GetProjectsByClassification(params)
}

func (pds *ProjectDashboardService) GetProjectsByStatus(params *projectDTOs.ProjectDashboardQueryParamsDTO) (*projectDTOs.ProjectsByStatusDTO, error) {
	return pds.repo.GetProjectsByStatus(params)
}

func (pds *ProjectDashboardService) GetProjectCompletion(params *projectDTOs.ProjectDashboardQueryParamsDTO) ([]projectDTOs.ProjectCompletionListDTO, error) {
	return pds.repo.GetProjectCompletion(params)
}

func (pds *ProjectDashboardService) GetProjectsDueSoonAndOverdue(params *projectDTOs.ProjectDashboardQueryParamsDTO) ([]projectDTOs.ProjectDueSoonOverdueListDTO, error) {
	return pds.repo.GetProjectsDueSoonAndOverdue(params)
}
