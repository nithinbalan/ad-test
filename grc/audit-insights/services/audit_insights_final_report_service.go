package services

import (
	"jethur/grc/audit-insights/repositories"
	auditDTOs "jethur/grc/audit/dtos"
	auditServices "jethur/grc/audit/services"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"
)

type AuditInsightsFinalReportServiceInterface interface {
	GetAll(params *auditDTOs.AuditFinalReportListQueryParamsDTO, currentUserID uint) (*auditDTOs.AuditFinalReportListResponseDTO, error)
	GetFullDetailsByID(id uint, currentUserID uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditFinalReportDetailDTO, error)
}

type AuditInsightsFinalReportService struct {
	repo                repositories.AuditInsightsFinalReportRepositoryInterface
	userSvc             organizationServices.UserServiceInterface
	auditSvc            auditServices.AuditServiceInterface
	auditFinalReportSvc auditServices.AuditFinalReportServiceInterface
}

func NewAuditInsightsFinalReportService(
	repo repositories.AuditInsightsFinalReportRepositoryInterface,
	userSvc organizationServices.UserServiceInterface,
	auditSvc auditServices.AuditServiceInterface,
	auditFinalReportSvc auditServices.AuditFinalReportServiceInterface,
) AuditInsightsFinalReportServiceInterface {
	return &AuditInsightsFinalReportService{
		repo:                repo,
		userSvc:             userSvc,
		auditSvc:            auditSvc,
		auditFinalReportSvc: auditFinalReportSvc,
	}
}

func (aidrs *AuditInsightsFinalReportService) GetAll(params *auditDTOs.AuditFinalReportListQueryParamsDTO, currentUserID uint) (*auditDTOs.AuditFinalReportListResponseDTO, error) {
	orgUnitID, err := aidrs.userSvc.GetUserPrimaryUnitID(currentUserID)
	if err != nil {
		return nil, err
	}
	return aidrs.repo.GetAll(params, orgUnitID)
}

func (aidrs *AuditInsightsFinalReportService) GetFullDetailsByID(id uint, currentUserID uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditFinalReportDetailDTO, error) {
	auditFinalReport, err := aidrs.auditFinalReportSvc.GetByID(id)
	if err != nil {
		return nil, err
	}
	auditID := auditFinalReport.AuditID
	audit, err := aidrs.auditFinalReportSvc.GetFullDetailsByID(auditID, currentUserID, "audit-insight", tx...)
	if err != nil {
		return nil, err
	}
	return audit, nil
}
