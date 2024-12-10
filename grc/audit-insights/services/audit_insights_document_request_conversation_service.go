package services

import (
	"jethur/grc/audit/dtos"
	auditServices "jethur/grc/audit/services"
	logServices "jethur/grc/logs/services"
	"jethur/utility"
)

type AuditInsightsDocumentRequestConversationServiceInterface interface {
	GetAll(auditDocumentRequestID uint, queryParams *dtos.AuditDocumentRequestConversationListQueryParamsDTO, tx ...utility.DatabaseTransaction) ([]*dtos.AuditDocumentRequestConversationListDTO, error)
	Create(auditDocumentRequestID uint, currentUserID uint, formDTO *dtos.AuditDocumentRequestConversationFormDTO, tx ...utility.DatabaseTransaction) (uint, error)
	GetDocumentPath(auditDocumentRequestID uint, conversationID uint, token string, tx ...utility.DatabaseTransaction) (string, error)
}

type AuditInsightsDocumentRequestConversationService struct {
	activityLogSvc                      logServices.ActivityLogServiceInterface
	auditDocumentRequestConversationSvc auditServices.AuditDocumentRequestConversationServiceInterface
	auditDocumentRequestSvc             auditServices.AuditDocumentRequestServiceInterface
}

func NewAuditInsightsDocumentRequestConversationService(
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditDocumentRequestConversationSvc auditServices.AuditDocumentRequestConversationServiceInterface,
	auditDocumentRequestSvc auditServices.AuditDocumentRequestServiceInterface,
) AuditInsightsDocumentRequestConversationServiceInterface {
	return &AuditInsightsDocumentRequestConversationService{
		activityLogSvc:                      activityLogSvc,
		auditDocumentRequestConversationSvc: auditDocumentRequestConversationSvc,
		auditDocumentRequestSvc:             auditDocumentRequestSvc,
	}
}

func (aidrcs *AuditInsightsDocumentRequestConversationService) GetAll(auditDocumentRequestID uint, queryParams *dtos.AuditDocumentRequestConversationListQueryParamsDTO, tx ...utility.DatabaseTransaction) ([]*dtos.AuditDocumentRequestConversationListDTO, error) {
	documentRequest, err := aidrcs.auditDocumentRequestSvc.GetFullDetailsByID(auditDocumentRequestID, tx...)
	if err != nil {
		return nil, err
	}
	auditID := documentRequest.Audit.ID
	auditDocumentRequestConversation, err := aidrcs.auditDocumentRequestConversationSvc.GetAll(auditID, auditDocumentRequestID, queryParams)
	if err != nil {
		return nil, err
	}
	return auditDocumentRequestConversation, nil
}

func (aidrcs *AuditInsightsDocumentRequestConversationService) Create(auditDocumentRequestID uint, currentUserID uint, formDTO *dtos.AuditDocumentRequestConversationFormDTO, tx ...utility.DatabaseTransaction) (uint, error) {
	documentRequest, err := aidrcs.auditDocumentRequestSvc.GetFullDetailsByID(auditDocumentRequestID, tx...)
	if err != nil {
		return 0, err
	}
	auditID := documentRequest.Audit.ID

	auditDocumentRequest, err := aidrcs.auditDocumentRequestConversationSvc.Create(auditID, auditDocumentRequestID, currentUserID, formDTO, "audit-insight")
	if err != nil {
		return 0, err
	}

	return auditDocumentRequest, nil
}

func (aidrcs *AuditInsightsDocumentRequestConversationService) GetDocumentPath(auditDocumentRequestID uint, conversationID uint, token string, tx ...utility.DatabaseTransaction) (string, error) {
	documentRequest, err := aidrcs.auditDocumentRequestSvc.GetFullDetailsByID(auditDocumentRequestID, tx...)
	if err != nil {
		return "", err
	}
	auditID := documentRequest.Audit.ID

	auditDocumentRequest, err := aidrcs.auditDocumentRequestConversationSvc.GetDocumentPath(auditID, auditDocumentRequestID, conversationID, token)
	if err != nil {
		return "", err
	}

	return auditDocumentRequest, nil
}
