package services

import (
	organizationServices "jethur/grc/organization/services"
	reportDTOs "jethur/grc/reports/dtos"

	"jethur/grc/reports/repositories"
)

type AuditInsightReportServiceInterface interface {
	GetAuditByStatus(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.AuditInsightAuditByStatusDTO, error)
	GetAuditByYear(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.AuditInsightAuditByYearDTO, error)
	GetAuditByMonth(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.AuditInsightAuditByMonthDTO, error)
	GetAuditPlanByCategory(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.AuditInsightAuditPlanByCategoryDTO, error)
	GetAuditPlanByTeamLead(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.AuditInsightAuditPlanByTeamLeadDTO, error)
	GetAuditPlanByAuditableItem(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.AuditInsightAuditPlanByAuditableItemDTO, error)
	GetAuditFindingByCategory(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.AuditInsightAuditFindingByCategoryDTO, error)
	GetAuditFindingByStatus(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.AuditInsightAuditFindingByStatusDTO, error)
	GetAuditFindingBySeverityLevel(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.AuditInsightAuditFindingBySeverityLevelDTO, error)
	GetAuditFindingByMonth(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.AuditInsightAuditFindingByMonthDTO, error)
	GetAuditFindingByYear(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.AuditInsightAuditFindingByYearDTO, error)
	GetAuditActionPlanByStatus(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.AuditInsightAuditActionPlanByStatusDTO, error)
	GetAuditActionPlanByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.AuditInsightAuditActionPlanByResponsibleUserDTO, error)
	GetAuditDraftReportByStatus(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.AuditInsightAuditDraftReportByStatusDTO, error)
	GetAuditFinalReportByStatus(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.AuditInsightAuditFinalReportByStatusDTO, error)
}

type AuditInsightReportService struct {
	repo    repositories.AuditInsightReportRepositoryInterface
	userSvc organizationServices.UserServiceInterface
}

func NewAuditInsightReportService(repo repositories.AuditInsightReportRepositoryInterface,
	userSvc organizationServices.UserServiceInterface) AuditInsightReportServiceInterface {
	return &AuditInsightReportService{repo: repo, userSvc: userSvc}
}

func (airs *AuditInsightReportService) GetAuditByStatus(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.AuditInsightAuditByStatusDTO, error) {
	orgUnitID, err := airs.userSvc.GetUserPrimaryUnitID(currentUserID)
	if err != nil {
		return nil, err
	}
	return airs.repo.GetAuditByStatus(params, orgUnitID)
}

func (airs *AuditInsightReportService) GetAuditByYear(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.AuditInsightAuditByYearDTO, error) {
	orgUnitID, err := airs.userSvc.GetUserPrimaryUnitID(currentUserID)
	if err != nil {
		return nil, err
	}
	return airs.repo.GetAuditByYear(params, orgUnitID)
}

func (airs *AuditInsightReportService) GetAuditByMonth(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.AuditInsightAuditByMonthDTO, error) {
	orgUnitID, err := airs.userSvc.GetUserPrimaryUnitID(currentUserID)
	if err != nil {
		return nil, err
	}
	return airs.repo.GetAuditByMonth(params, orgUnitID)
}

func (airs *AuditInsightReportService) GetAuditPlanByCategory(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.AuditInsightAuditPlanByCategoryDTO, error) {
	orgUnitID, err := airs.userSvc.GetUserPrimaryUnitID(currentUserID)
	if err != nil {
		return nil, err
	}
	return airs.repo.GetAuditPlanByCategory(params, orgUnitID)
}

func (airs *AuditInsightReportService) GetAuditPlanByTeamLead(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.AuditInsightAuditPlanByTeamLeadDTO, error) {
	orgUnitID, err := airs.userSvc.GetUserPrimaryUnitID(currentUserID)
	if err != nil {
		return nil, err
	}
	return airs.repo.GetAuditPlanByTeamLead(params, orgUnitID)
}

func (airs *AuditInsightReportService) GetAuditPlanByAuditableItem(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.AuditInsightAuditPlanByAuditableItemDTO, error) {
	orgUnitID, err := airs.userSvc.GetUserPrimaryUnitID(currentUserID)
	if err != nil {
		return nil, err
	}
	return airs.repo.GetAuditPlanByAuditableItem(params, orgUnitID)
}

func (airs *AuditInsightReportService) GetAuditFindingByCategory(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.AuditInsightAuditFindingByCategoryDTO, error) {
	orgUnitID, err := airs.userSvc.GetUserPrimaryUnitID(currentUserID)
	if err != nil {
		return nil, err
	}
	return airs.repo.GetAuditFindingByCategory(params, orgUnitID)
}

func (airs *AuditInsightReportService) GetAuditFindingByStatus(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.AuditInsightAuditFindingByStatusDTO, error) {
	orgUnitID, err := airs.userSvc.GetUserPrimaryUnitID(currentUserID)
	if err != nil {
		return nil, err
	}
	return airs.repo.GetAuditFindingByStatus(params, orgUnitID)
}

func (airs *AuditInsightReportService) GetAuditFindingBySeverityLevel(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.AuditInsightAuditFindingBySeverityLevelDTO, error) {
	orgUnitID, err := airs.userSvc.GetUserPrimaryUnitID(currentUserID)
	if err != nil {
		return nil, err
	}
	return airs.repo.GetAuditFindingBySeverityLevel(params, orgUnitID)
}

func (airs *AuditInsightReportService) GetAuditFindingByMonth(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.AuditInsightAuditFindingByMonthDTO, error) {
	orgUnitID, err := airs.userSvc.GetUserPrimaryUnitID(currentUserID)
	if err != nil {
		return nil, err
	}
	return airs.repo.GetAuditFindingByMonth(params, orgUnitID)
}

func (airs *AuditInsightReportService) GetAuditFindingByYear(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.AuditInsightAuditFindingByYearDTO, error) {
	orgUnitID, err := airs.userSvc.GetUserPrimaryUnitID(currentUserID)
	if err != nil {
		return nil, err
	}
	return airs.repo.GetAuditFindingByYear(params, orgUnitID)
}

func (airs *AuditInsightReportService) GetAuditActionPlanByStatus(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.AuditInsightAuditActionPlanByStatusDTO, error) {
	orgUnitID, err := airs.userSvc.GetUserPrimaryUnitID(currentUserID)
	if err != nil {
		return nil, err
	}
	return airs.repo.GetAuditActionPlanByStatus(params, orgUnitID)
}

func (airs *AuditInsightReportService) GetAuditActionPlanByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.AuditInsightAuditActionPlanByResponsibleUserDTO, error) {
	orgUnitID, err := airs.userSvc.GetUserPrimaryUnitID(currentUserID)
	if err != nil {
		return nil, err
	}
	return airs.repo.GetAuditActionPlanByResponsibleUser(params, orgUnitID)
}

func (airs *AuditInsightReportService) GetAuditDraftReportByStatus(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.AuditInsightAuditDraftReportByStatusDTO, error) {
	orgUnitID, err := airs.userSvc.GetUserPrimaryUnitID(currentUserID)
	if err != nil {
		return nil, err
	}
	return airs.repo.GetAuditDraftReportByStatus(params, orgUnitID)
}

func (airs *AuditInsightReportService) GetAuditFinalReportByStatus(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.AuditInsightAuditFinalReportByStatusDTO, error) {
	orgUnitID, err := airs.userSvc.GetUserPrimaryUnitID(currentUserID)
	if err != nil {
		return nil, err
	}
	return airs.repo.GetAuditFinalReportByStatus(params, orgUnitID)
}
