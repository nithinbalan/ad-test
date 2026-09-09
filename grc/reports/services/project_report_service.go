package services

import (
	reportDTOs "jethur/grc/reports/dtos"
	"jethur/grc/reports/repositories"
)

type ProjectReportServiceInterface interface {
	GetProjectByType(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectByTypeDTO, error)
	GetProjectByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectByStatusDTO, error)
	GetProjectByPhase(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectByPhaseDTO, error)
	GetProjectByClassification(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectByClassificationDTO, error)
	GetProjectByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectByResponsibleUnitDTO, error)
	GetProjectByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectByResponsibleUserDTO, error)
	GetProjectByYear(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectByYearDTO, error)
	GetProjectActionPlanByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectActionPlanByStatusDTO, error)
	GetProjectActionPlanByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectActionPlanByResponsibleUserDTO, error)
	GetProjectActionPlanByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectActionPlanByResponsibleUnitDTO, error)
}

type ProjectReportService struct {
	repo repositories.ProjectReportRepositoryInterface
}

func NewProjectReportService(repo repositories.ProjectReportRepositoryInterface) ProjectReportServiceInterface {
	return &ProjectReportService{repo: repo}
}

func (prs *ProjectReportService) GetProjectByType(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectByTypeDTO, error) {
	return prs.repo.GetProjectByType(params)
}

func (prs *ProjectReportService) GetProjectByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectByStatusDTO, error) {
	return prs.repo.GetProjectByStatus(params)
}

func (prs *ProjectReportService) GetProjectByPhase(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectByPhaseDTO, error) {
	return prs.repo.GetProjectByPhase(params)
}

func (prs *ProjectReportService) GetProjectByClassification(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectByClassificationDTO, error) {
	return prs.repo.GetProjectByClassification(params)
}

func (prs *ProjectReportService) GetProjectByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectByResponsibleUnitDTO, error) {
	return prs.repo.GetProjectByResponsibleUnit(params)
}

func (prs *ProjectReportService) GetProjectByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectByResponsibleUserDTO, error) {
	return prs.repo.GetProjectByResponsibleUser(params)
}

func (prs *ProjectReportService) GetProjectByYear(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectByYearDTO, error) {
	return prs.repo.GetProjectByYear(params)
}

func (prs *ProjectReportService) GetProjectActionPlanByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectActionPlanByStatusDTO, error) {
	return prs.repo.GetProjectActionPlanByStatus(params)
}

func (prs *ProjectReportService) GetProjectActionPlanByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectActionPlanByResponsibleUserDTO, error) {
	return prs.repo.GetProjectActionPlanByResponsibleUser(params)
}

func (prs *ProjectReportService) GetProjectActionPlanByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectActionPlanByResponsibleUnitDTO, error) {
	return prs.repo.GetProjectActionPlanByResponsibleUnit(params)
}
