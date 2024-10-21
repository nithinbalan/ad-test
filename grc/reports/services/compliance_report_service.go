package services

import (
	reportDTOs "jethur/grc/reports/dtos"
	"jethur/grc/reports/repositories"
)

type ComplianceReportServiceInterface interface {
	GetComplianceByCategory(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ComplianceByCategoryDTO, error)
}

type ComplianceReportService struct {
	repo repositories.ComplianceReportRepositoryInterface
}

func NewComplianceReportService(repo repositories.ComplianceReportRepositoryInterface) ComplianceReportServiceInterface {
	return &ComplianceReportService{repo: repo}
}

func (crs *ComplianceReportService) GetComplianceByCategory(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ComplianceByCategoryDTO, error) {
	return crs.repo.GetComplianceByCategory(params)
}
