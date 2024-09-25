package services

import (
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"
)

type AuditDraftReportReviewServiceInterface interface {
	Submit(auditID uint, auditDraftReportID uint, currentUserID uint, auditReportWorkflowFormDTO *auditDTOs.AuditReportWorkflowFormDTO) error
}

type AuditDraftReportReviewService struct {
	repo                   repositories.AuditDraftReportReviewRepositoryInterface
	activityLogSvc         logServices.ActivityLogServiceInterface
	auditSvc               AuditServiceInterface
	auditDraftReportSvc    AuditDraftReportServiceInterface
	auditWorkflowStatusSvc AuditWorkflowStatusServiceInterface
	userSvc                organizationServices.UserServiceInterface
}

func NewAuditDraftReportReviewService(
	repo repositories.AuditDraftReportReviewRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditSvc AuditServiceInterface,
	auditDraftReportSvc AuditDraftReportServiceInterface,
	auditWorkflowStatusSvc AuditWorkflowStatusServiceInterface,
	userSvc organizationServices.UserServiceInterface,
) AuditDraftReportReviewServiceInterface {
	return &AuditDraftReportReviewService{
		repo:                   repo,
		activityLogSvc:         activityLogSvc,
		auditSvc:               auditSvc,
		auditDraftReportSvc:    auditDraftReportSvc,
		auditWorkflowStatusSvc: auditWorkflowStatusSvc,
		userSvc:                userSvc,
	}
}

func (adrrs *AuditDraftReportReviewService) Submit(auditID uint, auditDraftReportID uint, currentUserID uint, auditReportWorkflowFormDTO *auditDTOs.AuditReportWorkflowFormDTO) error {
	tx := adrrs.repo.BeginTransaction()
	if _, err := adrrs.auditSvc.GetByID(auditID); err != nil {
		return err
	}

	auditDraftReportStatus, err := adrrs.auditWorkflowStatusSvc.GetByType("pending")
	if err != nil {
		return err
	}
	approvalType := "anyone"
	workflowType := "draft-report"
	auditReportWorkflow := &models.AuditReportWorkflow{
		AuditDraftReportID:    &auditDraftReportID,
		AuditWorkflowStatusID: auditDraftReportStatus.ID,
		Level:                 1,
		ApprovalType:          approvalType,
		Type:                  workflowType,
		Comment:               auditReportWorkflowFormDTO.Comment,
		CreatedBYID:           currentUserID,
	}
	id, err := adrrs.repo.Create(auditReportWorkflow, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	var participants []*models.AuditReportWorkflowParticipant
	for _, userID := range auditReportWorkflowFormDTO.UserIDs {
		auditReportWorkflowParticipant := &models.AuditReportWorkflowParticipant{
			AuditReportWorkflowID: id,
			UserID:                *userID,
			CreatedBYID:           currentUserID,
		}
		participants = append(participants, auditReportWorkflowParticipant)
	}

	if len(participants) > 0 {
		if err := adrrs.repo.CreateAuditReportWorkflowParticipants(participants, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	//todo - notification to workflow user

	tx.Commit()

	go func() {
		err := adrrs.activityLogSvc.LogActivity(
			models.SubmitAction,
			models.AuditReportWorkflow{}.TableName(),
			id,
			auditReportWorkflowFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.SubmitAction,
				"tableName": models.AuditReportWorkflow{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}
