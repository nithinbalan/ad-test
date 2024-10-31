package services

import (
	reportDTOs "jethur/grc/reports/dtos"
	"jethur/grc/reports/repositories"
)

type ComplianceReportServiceInterface interface {
	GetComplianceByComplianceCategory(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ComplianceByCategoryDTO, error)
	GetComplianceBySeverityLevel(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ComplianceBySeverityLevelDTO, error)
	GetComplianceByComplianceStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ComplianceByStatusDTO, error)
	GetComplianceByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ComplianceByResponsibleUnitDTO, error)
	GetComplianceByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ComplianceByResponsibleUserDTO, error)

	GetComplianceRequirementByComplianceStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ComplianceRequirementByComplianceStatusDTO, error)
	GetComplianceRequirementBySeverity(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ComplianceRequirementBySeverityLevelDTO, error)
	GetComplianceRequirementByCategory(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ComplianceRequirementByCategoryDTO, error)
	GetComplianceRequirementByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ComplianceRequirementByResponsibleUserDTO, error)

	GetComplianceActionPlanByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ComplianceActionPlanByStatusDTO, error)
	GetComplianceActionPlanByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ComplianceActionPlanByResponsibleUserDTO, error)
}

type ComplianceReportService struct {
	repo repositories.ComplianceReportRepositoryInterface
}

func NewComplianceReportService(repo repositories.ComplianceReportRepositoryInterface) ComplianceReportServiceInterface {
	return &ComplianceReportService{repo: repo}
}

func (crs *ComplianceReportService) GetComplianceByComplianceCategory(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ComplianceByCategoryDTO, error) {
	return crs.repo.GetComplianceByComplianceCategory(params)
}

func (crs *ComplianceReportService) GetComplianceBySeverityLevel(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ComplianceBySeverityLevelDTO, error) {
	return crs.repo.GetComplianceBySeverityLevel(params)
}

func (crs *ComplianceReportService) GetComplianceByComplianceStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ComplianceByStatusDTO, error) {
	return crs.repo.GetComplianceByComplianceStatus(params)
}

func (crs *ComplianceReportService) GetComplianceByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ComplianceByResponsibleUnitDTO, error) {
	return crs.repo.GetComplianceByResponsibleUnit(params)
}

func (crs *ComplianceReportService) GetComplianceByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ComplianceByResponsibleUserDTO, error) {
	return crs.repo.GetComplianceByResponsibleUser(params)
}

func (crs *ComplianceReportService) GetComplianceRequirementByComplianceStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ComplianceRequirementByComplianceStatusDTO, error) {
	return crs.repo.GetComplianceRequirementByComplianceStatus(params)
}

func (crs *ComplianceReportService) GetComplianceRequirementBySeverity(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ComplianceRequirementBySeverityLevelDTO, error) {
	return crs.repo.GetComplianceRequirementBySeverity(params)
}

func (crs *ComplianceReportService) GetComplianceRequirementByCategory(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ComplianceRequirementByCategoryDTO, error) {
	return crs.repo.GetComplianceRequirementByCategory(params)
}

func (crs *ComplianceReportService) GetComplianceRequirementByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ComplianceRequirementByResponsibleUserDTO, error) {
	return crs.repo.GetComplianceRequirementByResponsibleUser(params)
}

func (crs *ComplianceReportService) GetComplianceActionPlanByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ComplianceActionPlanByStatusDTO, error) {
	return crs.repo.GetComplianceActionPlanByStatus(params)
}

func (crs *ComplianceReportService) GetComplianceActionPlanByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ComplianceActionPlanByResponsibleUserDTO, error) {
	return crs.repo.GetComplianceActionPlanByResponsibleUser(params)
}
