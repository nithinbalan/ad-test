package services

import (
	reportDTOs "jethur/grc/reports/dtos"
	"jethur/grc/reports/repositories"
)

type AuditReportServiceInterface interface {
	GetAuditPlanByCategory(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditPlanByCategoryDTO, error)
	GetAuditPlanByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditPlanByStatusDTO, error)
	GetAuditPlanByTeam(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditPlanByTeamDTO, error)
	GetAuditPlanByTeamLead(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditPlanByTeamLeadDTO, error)
	GetAuditPlanByYear(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditPlanByYearDTO, error)

	GetAuditByCategory(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditByCategoryDTO, error)
	GetAuditByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditByStatusDTO, error)
	GetAuditByOpenVsClosed(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditByOpenVsClosedDTO, error)
	GetAuditByMonth(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditByMonthDTO, error)
	GetAuditByTeamLead(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditByTeamLeadDTO, error)

	GetAuditFindingByCategory(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditFindingByCategoryDTO, error)
}

type AuditReportService struct {
	repo repositories.AuditReportRepositoryInterface
}

func NewAuditReportService(repo repositories.AuditReportRepositoryInterface) AuditReportServiceInterface {
	return &AuditReportService{repo: repo}
}

func (ars *AuditReportService) GetAuditPlanByCategory(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditPlanByCategoryDTO, error) {
	return ars.repo.GetAuditPlanByCategory(params)
}

func (ars *AuditReportService) GetAuditPlanByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditPlanByStatusDTO, error) {
	return ars.repo.GetAuditPlanByStatus(params)
}

func (ars *AuditReportService) GetAuditPlanByTeam(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditPlanByTeamDTO, error) {
	return ars.repo.GetAuditPlanByTeam(params)
}

func (ars *AuditReportService) GetAuditPlanByTeamLead(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditPlanByTeamLeadDTO, error) {
	return ars.repo.GetAuditPlanByTeamLead(params)
}

func (ars *AuditReportService) GetAuditPlanByYear(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditPlanByYearDTO, error) {
	return ars.repo.GetAuditPlanByYear(params)
}

func (ars *AuditReportService) GetAuditByCategory(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditByCategoryDTO, error) {
	return ars.repo.GetAuditByCategory(params)
}

func (ars *AuditReportService) GetAuditByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditByStatusDTO, error) {
	return ars.repo.GetAuditByStatus(params)
}

func (ars *AuditReportService) GetAuditByOpenVsClosed(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditByOpenVsClosedDTO, error) {
	return ars.repo.GetAuditByOpenVsClosed(params)
}

func (ars *AuditReportService) GetAuditByMonth(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditByMonthDTO, error) {
	return ars.repo.GetAuditByMonth(params)
}

func (ars *AuditReportService) GetAuditByTeamLead(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditByTeamLeadDTO, error) {
	return ars.repo.GetAuditByTeamLead(params)
}

func (ars *AuditReportService) GetAuditFindingByCategory(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditFindingByCategoryDTO, error) {
	return ars.repo.GetAuditFindingByCategory(params)
}
