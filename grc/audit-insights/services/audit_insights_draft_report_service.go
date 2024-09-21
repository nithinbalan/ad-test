package services

import (
	"jethur/grc/audit-insights/repositories"
	auditDTOs "jethur/grc/audit/dtos"
	auditServices "jethur/grc/audit/services"
	logServices "jethur/grc/logs/services"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"
)

type AuditInsightsDraftReportServiceInterface interface {
	GetAll(params *auditDTOs.AuditDraftReportListQueryParamsDTO, currentUserID uint) (*auditDTOs.AuditDraftReportListResponseDTO, error)
	GetFullDetailsByID(auditID uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditDraftReportDetailDTO, error)
}

type AuditInsightsDraftReportService struct {
	repo                repositories.AuditInsightsDraftReportRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	userSvc             organizationServices.UserServiceInterface
	auditSvc            auditServices.AuditServiceInterface
	auditDraftReportSvc auditServices.AuditDraftReportServiceInterface
}

func NewAuditInsightsDraftReportService(
	repo repositories.AuditInsightsDraftReportRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	auditSvc auditServices.AuditServiceInterface,
	auditDraftReportSvc auditServices.AuditDraftReportServiceInterface,
) AuditInsightsDraftReportServiceInterface {
	return &AuditInsightsDraftReportService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		userSvc:             userSvc,
		auditSvc:            auditSvc,
		auditDraftReportSvc: auditDraftReportSvc,
	}
}

func (aidrs *AuditInsightsDraftReportService) GetAll(params *auditDTOs.AuditDraftReportListQueryParamsDTO, currentUserID uint) (*auditDTOs.AuditDraftReportListResponseDTO, error) {
	orgUnitID, err := aidrs.userSvc.GetUserPrimaryUnitID(currentUserID)
	if err != nil {
		return nil, err
	}
	return aidrs.repo.GetAll(params, orgUnitID)
}

func (aidrs *AuditInsightsDraftReportService) GetFullDetailsByID(auditID uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditDraftReportDetailDTO, error) {
	audit, err := aidrs.auditDraftReportSvc.GetFullDetailsByID(auditID, tx...)
	if err != nil {
		return nil, err
	}
	return audit, nil
}
