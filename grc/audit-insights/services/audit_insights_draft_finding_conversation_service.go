package services

import (
	"jethur/grc/audit-insights/repositories"
	"jethur/grc/audit/dtos"
	auditServices "jethur/grc/audit/services"
	logServices "jethur/grc/logs/services"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"
)

type AuditInsightsDraftFindingConversationServiceInterface interface {
	GetAll(auditDraftFindingID uint, currentUserID uint, queryParams *dtos.AuditDraftFindingConversationListQueryParamsDTO, tx ...utility.DatabaseTransaction) ([]*dtos.AuditDraftFindingConversationListDTO, error)
	Create(auditDraftFindingID uint, currentUserID uint, formDTO *dtos.AuditDraftFindingConversationFormDTO, tx ...utility.DatabaseTransaction) (uint, error)
	GetDocumentPath(auditDraftFindingID uint, id uint, token string) (string, error)
}

type AuditInsightsDraftFindingConversationService struct {
	repo                             repositories.AuditInsightsDraftFindingConversationRepositoryInterface
	activityLogSvc                   logServices.ActivityLogServiceInterface
	auditDraftFindingStatusSvc       auditServices.AuditDraftFindingStatusServiceInterface
	auditDraftFindingConversationSvc auditServices.AuditDraftFindingConversationServiceInterface
	auditDraftFindingSvc             auditServices.AuditDraftFindingServiceInterface
	auditSvc                         auditServices.AuditServiceInterface
	userSvc                          organizationServices.UserServiceInterface
}

func NewAuditInsightsDraftFindingConversationService(
	repo repositories.AuditInsightsDraftFindingConversationRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditDraftFindingStatusSvc auditServices.AuditDraftFindingStatusServiceInterface,
	auditDraftFindingConversationSvc auditServices.AuditDraftFindingConversationServiceInterface,
	auditDraftFindingSvc auditServices.AuditDraftFindingServiceInterface,
	auditSvc auditServices.AuditServiceInterface,
	userSvc organizationServices.UserServiceInterface,
) AuditInsightsDraftFindingConversationServiceInterface {
	return &AuditInsightsDraftFindingConversationService{
		repo:                             repo,
		activityLogSvc:                   activityLogSvc,
		auditDraftFindingStatusSvc:       auditDraftFindingStatusSvc,
		auditDraftFindingConversationSvc: auditDraftFindingConversationSvc,
		auditDraftFindingSvc:             auditDraftFindingSvc,
		auditSvc:                         auditSvc,
		userSvc:                          userSvc,
	}
}

func (aidfcs *AuditInsightsDraftFindingConversationService) GetAll(auditDraftFindingID uint, currentUserID uint, queryParams *dtos.AuditDraftFindingConversationListQueryParamsDTO, tx ...utility.DatabaseTransaction) ([]*dtos.AuditDraftFindingConversationListDTO, error) {
	_, err := aidfcs.auditDraftFindingSvc.GetFullDetailsByID(auditDraftFindingID, tx...)
	if err != nil {
		return nil, err
	}
	orgUnitID, err := aidfcs.userSvc.GetUserPrimaryUnitID(currentUserID)
	if err != nil {
		return nil, err
	}
	auditDraftFindingConversation, err := aidfcs.repo.GetAll(orgUnitID, auditDraftFindingID, queryParams)
	if err != nil {
		return nil, err
	}
	return auditDraftFindingConversation, nil
}

func (aidfcs *AuditInsightsDraftFindingConversationService) Create(auditDraftFindingID uint, currentUserID uint, formDTO *dtos.AuditDraftFindingConversationFormDTO, tx ...utility.DatabaseTransaction) (uint, error) {

	auditDocumentRequest, err := aidfcs.auditDraftFindingConversationSvc.Create(auditDraftFindingID, currentUserID, formDTO)
	if err != nil {
		return 0, err
	}

	return auditDocumentRequest, nil
}

func (aidfcs *AuditInsightsDraftFindingConversationService) GetDocumentPath(auditDraftFindingID uint, id uint, token string) (string, error) {

	auditDocumentRequest, err := aidfcs.auditDraftFindingConversationSvc.GetDocumentPath(auditDraftFindingID, id, token)
	if err != nil {
		return "", err
	}

	return auditDocumentRequest, nil
}
