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
	GetAuditFindingByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditFindingByStatusDTO, error)
	GetAuditFindingBySeverityLevel(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditFindingBySeverityLevelDTO, error)
	GetAuditFindingByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditFindingByResponsibleUnitDTO, error)
	GetAuditFindingByMonth(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditFindingByMonthDTO, error)
	GetAuditFindingByAuditor(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditFindingByAuditorDTO, error)

	GetAuditActionPlanByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditActionPlanByStatusDTO, error)
	GetAuditActionPlanBySeverityLevel(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditActionPlanBySeverityLevelDTO, error)
	GetAuditActionPlanByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditActionPlanByResponsibleUserDTO, error)
	GetAuditActionPlanByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditActionPlanByResponsibleUnitDTO, error)
	GetAuditActionPlanOnTimeVsDelayed(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditActionPlanOnTimeVsDelayedDTO, error)
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

func (ars *AuditReportService) GetAuditFindingByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditFindingByStatusDTO, error) {
	return ars.repo.GetAuditFindingByStatus(params)
}

func (ars *AuditReportService) GetAuditFindingBySeverityLevel(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditFindingBySeverityLevelDTO, error) {
	return ars.repo.GetAuditFindingBySeverityLevel(params)
}

func (ars *AuditReportService) GetAuditFindingByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditFindingByResponsibleUnitDTO, error) {
	return ars.repo.GetAuditFindingByResponsibleUnit(params)
}

func (ars *AuditReportService) GetAuditFindingByMonth(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditFindingByMonthDTO, error) {
	return ars.repo.GetAuditFindingByMonth(params)
}

func (ars *AuditReportService) GetAuditFindingByAuditor(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditFindingByAuditorDTO, error) {
	return ars.repo.GetAuditFindingByAuditor(params)
}

func (ars *AuditReportService) GetAuditActionPlanByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditActionPlanByStatusDTO, error) {
	return ars.repo.GetAuditActionPlanByStatus(params)
}

func (ars *AuditReportService) GetAuditActionPlanBySeverityLevel(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditActionPlanBySeverityLevelDTO, error) {
	return ars.repo.GetAuditActionPlanBySeverityLevel(params)
}

func (ars *AuditReportService) GetAuditActionPlanByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditActionPlanByResponsibleUserDTO, error) {
	return ars.repo.GetAuditActionPlanByResponsibleUser(params)
}

func (ars *AuditReportService) GetAuditActionPlanByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditActionPlanByResponsibleUnitDTO, error) {
	return ars.repo.GetAuditActionPlanByResponsibleUnit(params)
}

func (ars *AuditReportService) GetAuditActionPlanOnTimeVsDelayed(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.AuditActionPlanOnTimeVsDelayedDTO, error) {
	return ars.repo.GetAuditActionPlanOnTimeVsDelayed(params)
}
