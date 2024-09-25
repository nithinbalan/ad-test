package services

import (
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"
	"time"
)

type AuditDraftReportReviewServiceInterface interface {
	Submit(auditID uint, auditDraftReportID uint, currentUserID uint, auditReportWorkflowFormDTO *auditDTOs.AuditReportWorkflowFormDTO) error
}

type AuditDraftReportReviewService struct {
	repo                      repositories.AuditDraftReportReviewRepositoryInterface
	activityLogSvc            logServices.ActivityLogServiceInterface
	auditSvc                  AuditServiceInterface
	auditDraftReportSvc       AuditDraftReportServiceInterface
	workflowStatusSvc         masterServices.WorkflowStatusServiceInterface
	userSvc                   organizationServices.UserServiceInterface
	auditDraftReportStatusSvc AuditDraftReportStatusServiceInterface
}

func NewAuditDraftReportReviewService(
	repo repositories.AuditDraftReportReviewRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditSvc AuditServiceInterface,
	auditDraftReportSvc AuditDraftReportServiceInterface,
	workflowStatusSvc masterServices.WorkflowStatusServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	auditDraftReportStatusSvc AuditDraftReportStatusServiceInterface,
) AuditDraftReportReviewServiceInterface {
	return &AuditDraftReportReviewService{
		repo:                      repo,
		activityLogSvc:            activityLogSvc,
		auditSvc:                  auditSvc,
		auditDraftReportSvc:       auditDraftReportSvc,
		workflowStatusSvc:         workflowStatusSvc,
		userSvc:                   userSvc,
		auditDraftReportStatusSvc: auditDraftReportStatusSvc,
	}
}

func (adrrs *AuditDraftReportReviewService) Submit(auditID uint, auditDraftReportID uint, currentUserID uint, auditReportWorkflowFormDTO *auditDTOs.AuditReportWorkflowFormDTO) error {
	tx := adrrs.repo.BeginTransaction()
	if _, err := adrrs.auditSvc.GetByID(auditID); err != nil {
		return err
	}
	existingAuditDraftReport, err := adrrs.auditDraftReportSvc.GetByID(auditDraftReportID)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, auditDraftReportID)
	}
	auditDraftReportStatusDraft, err := adrrs.auditDraftReportStatusSvc.GetByType("draft")
	if err != nil {
		return err
	}
	if existingAuditDraftReport.AuditDraftReportStatusID != auditDraftReportStatusDraft.ID {
		tx.Rollback()
		return utility.NewForbiddenError("The audit draft report has already been submitted.")
	}
	workflowSubmitStatus, err := adrrs.workflowStatusSvc.GetByType("submit")
	if err != nil {
		return err
	}
	approvalType := "anyone"
	workflowType := "draft-report"
	auditReportWorkflow := &models.AuditReportWorkflow{
		AuditDraftReportID: &auditDraftReportID,
		WorkflowStatusID:   workflowSubmitStatus.ID,
		Level:              1,
		ApprovalType:       approvalType,
		Type:               workflowType,
		CreatedBYID:        currentUserID,
	}
	auditReportWorkflowID, err := adrrs.repo.CreateAuditReportWorkflow(auditReportWorkflow, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	var participants []*models.AuditReportWorkflowParticipant
	for _, userID := range auditReportWorkflowFormDTO.UserIDs {
		auditReportWorkflowParticipant := &models.AuditReportWorkflowParticipant{
			AuditReportWorkflowID: auditReportWorkflowID,
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
	//add review for history
	auditReportReview := &models.AuditReportReview{
		AuditDraftReportID:    auditDraftReportID,
		AuditReportWorkflowID: auditReportWorkflowID,
		WorkflowStatusID:      workflowSubmitStatus.ID, // Initially set to 'pending'
		UserID:                currentUserID,
		CreatedByID:           currentUserID,
	}
	id, err := adrrs.repo.Create(auditReportReview, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	auditDraftReportInReviewStatus, err := adrrs.auditDraftReportStatusSvc.GetByType("in-review")
	if err != nil {
		return err
	}
	now := time.Now()
	existingAuditDraftReport.SubmittedAt = &now
	existingAuditDraftReport.AuditDraftReportStatusID = auditDraftReportInReviewStatus.ID
	existingAuditDraftReport.SubmittedBYID = &currentUserID

	// err = adrrs.auditDraftReportSvc.Update(existingAuditDraftReport, tx)
	// if err != nil {
	// 	tx.Rollback()
	// 	return err
	// }
	//todo - notification to workflow user

	tx.Commit()

	go func() {
		err := adrrs.activityLogSvc.LogActivity(
			models.SubmitAction,
			models.AuditReportReview{}.TableName(),
			id,
			auditReportWorkflowFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.SubmitAction,
				"tableName": models.AuditReportReview{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}
