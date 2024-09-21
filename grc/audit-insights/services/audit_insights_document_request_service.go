package services

import (
	"jethur/grc/audit-insights/repositories"
	auditDTOs "jethur/grc/audit/dtos"
	auditServices "jethur/grc/audit/services"
	logServices "jethur/grc/logs/services"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"
)

type AuditInsightsDocumentRequestServiceInterface interface {
	GetAll(params *auditDTOs.AuditDocumentRequestListQueryParamsDTO, currentUserID uint) (*auditDTOs.AuditDocumentRequestListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditDocumentRequestDetailDTO, error)
}

type AuditInsightsDocumentRequestService struct {
	repo                    repositories.AuditInsightsDocumentRequestRepositoryInterface
	activityLogSvc          logServices.ActivityLogServiceInterface
	userSvc                 organizationServices.UserServiceInterface
	auditDocumentRequestSvc auditServices.AuditDocumentRequestServiceInterface
}

func NewAuditInsightsDocumentRequestService(
	repo repositories.AuditInsightsDocumentRequestRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	auditDocumentRequestSvc auditServices.AuditDocumentRequestServiceInterface,

) AuditInsightsDocumentRequestServiceInterface {
	return &AuditInsightsDocumentRequestService{
		repo:                    repo,
		activityLogSvc:          activityLogSvc,
		userSvc:                 userSvc,
		auditDocumentRequestSvc: auditDocumentRequestSvc,
	}
}

func (aidrs *AuditInsightsDocumentRequestService) GetAll(params *auditDTOs.AuditDocumentRequestListQueryParamsDTO, currentUserID uint) (*auditDTOs.AuditDocumentRequestListResponseDTO, error) {
	orgUnitID, err := aidrs.userSvc.GetUserPrimaryUnitID(currentUserID)
	if err != nil {
		return nil, err
	}
	return aidrs.repo.GetAll(params, orgUnitID)
}

func (aidrs *AuditInsightsDocumentRequestService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditDocumentRequestDetailDTO, error) {
	documentRequest, err := aidrs.auditDocumentRequestSvc.GetFullDetailsByID(id, tx...)
	if err != nil {
		return nil, err
	}
	return documentRequest, nil
}
