package services

import (
	reportDTOs "jethur/grc/reports/dtos"
	"jethur/grc/reports/repositories"
)

type ProjectSecurityGateReportServiceInterface interface {
	GetProjectSecurityGateByStage(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectSecurityGateByStageDTO, error)
	GetProjectSecurityGateByReadinessStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectSecurityGateByReadinessStatusDTO, error)
	GetProjectSecurityGateByClassification(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectSecurityGateByClassificationDTO, error)
}

type ProjectSecurityGateReportService struct {
	repo repositories.ProjectSecurityGateReportRepositoryInterface
}

func NewProjectSecurityGateReportService(repo repositories.ProjectSecurityGateReportRepositoryInterface) ProjectSecurityGateReportServiceInterface {
	return &ProjectSecurityGateReportService{repo: repo}
}

func (psgrs *ProjectSecurityGateReportService) GetProjectSecurityGateByStage(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectSecurityGateByStageDTO, error) {
	return psgrs.repo.GetProjectSecurityGateByStage(params)
}

func (psgrs *ProjectSecurityGateReportService) GetProjectSecurityGateByReadinessStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectSecurityGateByReadinessStatusDTO, error) {
	return psgrs.repo.GetProjectSecurityGateByReadinessStatus(params)
}

func (psgrs *ProjectSecurityGateReportService) GetProjectSecurityGateByClassification(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ProjectSecurityGateByClassificationDTO, error) {
	return psgrs.repo.GetProjectSecurityGateByClassification(params)
}
