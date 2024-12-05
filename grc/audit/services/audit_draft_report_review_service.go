package services

import (
	"encoding/json"
	"fmt"
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	taskServices "jethur/grc/task/services"
	"jethur/utility"
	"time"
)

type AuditDraftReportReviewServiceInterface interface {
	GetAll(auditID uint, auditDraftReportID uint, params *auditDTOs.AuditReportReviewListQueryParamsDTO) ([]*auditDTOs.AuditReportReviewListDTO, error)
	Submit(auditID uint, auditDraftReportID uint, currentUserID uint, auditReportWorkflowFormDTO *auditDTOs.AuditReportWorkflowSubmitFormDTO) error
	ReSubmit(auditID uint, auditDraftReportID uint, currentUserID uint, auditReportWorkflowFormDTO *auditDTOs.AuditReportWorkflowSubmitFormDTO) error
	Approve(auditID uint, auditDraftReportID uint, currentUserID uint, auditReportWorkflowFormDTO *auditDTOs.AuditReportWorkflowApproveFormDTO) error
	Revert(auditID uint, auditDraftReportID uint, currentUserID uint, auditReportWorkflowFormDTO *auditDTOs.AuditReportWorkflowApproveFormDTO) error
}

type AuditDraftReportReviewService struct {
	repo                      repositories.AuditDraftReportReviewRepositoryInterface
	activityLogSvc            logServices.ActivityLogServiceInterface
	auditSvc                  AuditServiceInterface
	auditDraftReportSvc       AuditDraftReportServiceInterface
	workflowStatusSvc         masterServices.WorkflowStatusServiceInterface
	userSvc                   organizationServices.UserServiceInterface
	auditDraftReportStatusSvc AuditDraftReportStatusServiceInterface
	notificationSvc           commonServices.NotificationServiceInterface
	taskSvc                   taskServices.TaskServiceInterface
}

func NewAuditDraftReportReviewService(
	repo repositories.AuditDraftReportReviewRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditSvc AuditServiceInterface,
	auditDraftReportSvc AuditDraftReportServiceInterface,
	workflowStatusSvc masterServices.WorkflowStatusServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	auditDraftReportStatusSvc AuditDraftReportStatusServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	taskSvc taskServices.TaskServiceInterface,
) AuditDraftReportReviewServiceInterface {
	return &AuditDraftReportReviewService{
		repo:                      repo,
		activityLogSvc:            activityLogSvc,
		auditSvc:                  auditSvc,
		auditDraftReportSvc:       auditDraftReportSvc,
		workflowStatusSvc:         workflowStatusSvc,
		userSvc:                   userSvc,
		auditDraftReportStatusSvc: auditDraftReportStatusSvc,
		notificationSvc:           notificationSvc,
		taskSvc:                   taskSvc,
	}
}
func (adrcs *AuditDraftReportReviewService) GetAll(auditID uint, auditDraftReportID uint, params *auditDTOs.AuditReportReviewListQueryParamsDTO) ([]*auditDTOs.AuditReportReviewListDTO, error) {
	if _, err := adrcs.auditSvc.GetByID(auditID); err != nil {
		return nil, err
	}
	if _, err := adrcs.auditDraftReportSvc.GetByID(auditDraftReportID); err != nil {
		return nil, err
	}

	return adrcs.repo.GetAll(auditDraftReportID, params)
}

func (adrrs *AuditDraftReportReviewService) Submit(auditID uint, auditDraftReportID uint, currentUserID uint, auditReportWorkflowFormDTO *auditDTOs.AuditReportWorkflowSubmitFormDTO) error {
	tx := adrrs.repo.BeginTransaction()
	if _, err := adrrs.auditSvc.GetByID(auditID); err != nil {
		tx.Rollback()
		return err
	}
	existingAuditDraftReport, err := adrrs.auditDraftReportSvc.GetByID(auditDraftReportID)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, auditDraftReportID)
	}
	auditDraftReportStatusDraft, err := adrrs.auditDraftReportStatusSvc.GetByType("draft")
	if err != nil {
		tx.Rollback()
		return err
	}
	if existingAuditDraftReport.AuditDraftReportStatusID != auditDraftReportStatusDraft.ID {
		tx.Rollback()
		return utility.NewForbiddenError("The audit draft report has already been submitted.")
	}

	workflowSubmitStatus, err := adrrs.workflowStatusSvc.GetByType("submit")
	if err != nil {
		tx.Rollback()
		return err
	}
	// No existing workflow, create a new one
	approvalType := "anyone"
	workflowType := "draft"
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
		AuditDraftReportID:    &auditDraftReportID,
		AuditReportWorkflowID: auditReportWorkflowID,
		WorkflowStatusID:      workflowSubmitStatus.ID,
		UserID:                &currentUserID,
		Comment:               auditReportWorkflowFormDTO.Comment,
		CreatedBYID:           currentUserID,
	}
	id, err := adrrs.repo.Create(auditReportReview, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	//update status & AuditeeReviewSubmitted user details on draft report
	auditDraftReportInReviewStatus, err := adrrs.auditDraftReportStatusSvc.GetByType("in-review")
	if err != nil {
		tx.Rollback()
		return err
	}
	nextReviewLevel := 1

	err = adrrs.auditDraftReportSvc.StatusUpdate(auditDraftReportID, currentUserID, auditDraftReportInReviewStatus.ID, nextReviewLevel, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	existingAuditDraftReportFullDetails, err := adrrs.auditDraftReportSvc.GetByAuditID(auditID)
	if err != nil {
		tx.Rollback()
		return err
	}
	urlParam := fmt.Sprintf(`{"id":%d}`, auditID)
	var taskIDs []uint
	auditFinalReportID := &existingAuditDraftReport.ID
	for _, participant := range participants {
		taskDTO := commonDTOs.TaskFormDTO{
			ModuleIdentifier:      "audit",
			AuditFinalReportID:    auditFinalReportID,
			AuditReportWorkflowID: &auditReportWorkflowID,
			Title:                 existingAuditDraftReportFullDetails.Audit.AuditPlan.Title,
			Description:           existingAuditDraftReportFullDetails.Audit.AuditPlan.Description,
			URLParam:              &urlParam,
			QueryParam:            nil,
			ResponsibleUserID:     participant.User.ID,
		}

		taskID, err := adrrs.taskSvc.Create(taskDTO, currentUserID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
		taskIDs = append(taskIDs, taskID)

	}

	for _, taskID := range taskIDs {
		relatedTaskIDs := []uint{}
		for _, otherTaskID := range taskIDs {
			if otherTaskID != taskID {
				relatedTaskIDs = append(relatedTaskIDs, otherTaskID)
			}
		}

		relatedTaskIDsJSON, err := json.Marshal(relatedTaskIDs)
		if err != nil {
			tx.Rollback()
			return err
		}

		err = adrrs.taskSvc.UpdateTaskWithRelatedIDs(taskID, string(relatedTaskIDsJSON), tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	go func() {
		entity := "audit"
		action := "review-user-assign"

		for _, participant := range auditReportWorkflowFormDTO.UserIDs {
			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserID:               *participant,
				ModuleIdentifier:     "audit",
				NotificationType:     string(commonDTOs.ReviewType),
				HeadingEn:            "Audit Draft Report Review",
				HeadingAr:            "مراجعة مسودة تقرير التدقيق",
				TitleEn:              existingAuditDraftReportFullDetails.Audit.AuditPlan.Title,
				TitleAr:              existingAuditDraftReportFullDetails.Audit.AuditPlan.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			responsibleUser, err := adrrs.userSvc.GetByID(*participant)
			if err != nil {
				utility.LogError("Failed to fetch review user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			assignedBy, err := adrrs.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch assigning user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			emailTo := []string{responsibleUser.Email}
			subject := fmt.Sprintf("Audit Finale Report /%s: Review Assigned", existingAuditDraftReportFullDetails.Audit.AuditPlan.Title)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":           responsibleUser.FirstNameEn,
				"LastNameEn":            responsibleUser.LastNameEn,
				"AuditTitle":            existingAuditDraftReportFullDetails.Audit.AuditPlan.Title,
				"Reference Number":      existingAuditDraftReportFullDetails.Audit.AuditPlan.ReferenceNumber,
				"AssignedByFirstNameEn": assignedBy.FirstNameEn,
				"AssignedByLastNameEn":  assignedBy.LastNameEn,
				"AssignedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           fmt.Sprintf("%s/audit/audits/%d/draft-report", appURL, auditID),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "audit/templates/email/audit_draft_report_workflow_review_user_mail.html",
				EmailData:            emailData,
			}

			// Send Notification and Email
			err = adrrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "audit")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
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

func (adrrs *AuditDraftReportReviewService) Approve(auditID uint, auditDraftReportID uint, currentUserID uint, auditReportWorkflowFormDTO *auditDTOs.AuditReportWorkflowApproveFormDTO) error {
	tx := adrrs.repo.BeginTransaction()
	if _, err := adrrs.auditSvc.GetByID(auditID); err != nil {
		tx.Rollback()
		return err
	}
	existingAuditDraftReport, err := adrrs.auditDraftReportSvc.GetByID(auditDraftReportID)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, auditDraftReportID)
	}
	// Get Existing Audit Report Workflow
	existingAuditReportWorkflow, err := adrrs.repo.GetAuditReportWorkflowByDraftReportID(auditDraftReportID, tx)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, auditDraftReportID)
	}

	// Check if current user is a participant
	isParticipant := false
	for _, participant := range existingAuditReportWorkflow.AuditReportWorkflowParticipants {
		if participant.UserID == currentUserID {
			isParticipant = true
			break
		}
	}

	// If the current user is not a participant, return a forbidden error
	if !isParticipant {
		tx.Rollback()
		return utility.NewForbiddenError("The current user is not authorized to perform this action.")
	}

	//check level is 1
	if existingAuditDraftReport.NextReviewLevel != 1 {
		tx.Rollback()
		return utility.NewForbiddenError("The audit draft report is not at the required review level for this action.")
	}
	auditDraftReportStatusInReview, err := adrrs.auditDraftReportStatusSvc.GetByType("in-review")
	if err != nil {
		tx.Rollback()
		return err
	}
	if existingAuditDraftReport.AuditDraftReportStatusID != auditDraftReportStatusInReview.ID {
		tx.Rollback()
		return utility.NewForbiddenError("The audit draft report is either not in review or has not been submitted.")
	}
	workflowApproveStatus, err := adrrs.workflowStatusSvc.GetByType("approve")
	if err != nil {
		tx.Rollback()
		return err
	}
	// Update the WorkflowStatusID to workflowApproveStatus
	existingAuditReportWorkflow.WorkflowStatusID = workflowApproveStatus.ID

	// Save the updated AuditReportWorkflow
	if err := adrrs.repo.UpdateAuditReportWorkflow(existingAuditReportWorkflow, tx); err != nil {
		tx.Rollback()
		return err
	}
	//add review for history
	auditReportReview := &models.AuditReportReview{
		AuditDraftReportID:    &auditDraftReportID,
		AuditReportWorkflowID: existingAuditReportWorkflow.ID,
		WorkflowStatusID:      workflowApproveStatus.ID,
		UserID:                &currentUserID,
		Comment:               auditReportWorkflowFormDTO.Comment,
		CreatedBYID:           currentUserID,
	}
	id, err := adrrs.repo.Create(auditReportReview, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	//update status & AuditeeReviewSubmitted user details on draft report
	auditDraftReportApprovedStatus, err := adrrs.auditDraftReportStatusSvc.GetByType("approved")
	if err != nil {
		tx.Rollback()
		return err
	}
	nextReviewLevel := 2
	err = adrrs.auditDraftReportSvc.StatusUpdate(auditDraftReportID, currentUserID, auditDraftReportApprovedStatus.ID, nextReviewLevel, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = adrrs.taskSvc.UpdateAuditDraftReportTaskStatus(auditDraftReportID, currentUserID, existingAuditReportWorkflow.ID)
	if err != nil {
		return err
	}

	tx.Commit()

	go func() {
		urlParam := fmt.Sprintf(`{"id":%d}`, auditID)
		entity := "audit"
		action := "approve"

		// Create a map to track unique recipients
		recipientsMap := make(map[uint]bool)

		// Add participants excluding the current user
		for _, participant := range existingAuditReportWorkflow.AuditReportWorkflowParticipants {
			if participant.UserID != currentUserID {
				recipientsMap[participant.UserID] = true
			}
		}

		// Add InternalReviewSubmittedBYID
		recipientsMap[*existingAuditDraftReport.InternalReviewSubmittedBYID] = true

		// Iterate through recipients
		for userID := range recipientsMap {
			responsibleUser, err := adrrs.userSvc.GetByID(userID)
			if err != nil {
				utility.LogError("Failed to fetch review user", map[string]interface{}{
					"error": err.Error(),
				})
				continue
			}

			assignedBy, err := adrrs.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch assigning user", map[string]interface{}{
					"error": err.Error(),
				})
				continue
			}

			emailTo := []string{responsibleUser.Email}
			subject := fmt.Sprintf("Audit Final Report /%s: Approval Notification", existingAuditDraftReport.Audit.AuditPlan.Title)
			appURL := utility.GetAppURL()

			emailData := map[string]interface{}{
				"FirstNameEn":           responsibleUser.FirstNameEn,
				"LastNameEn":            responsibleUser.LastNameEn,
				"AuditTitle":            existingAuditDraftReport.Audit.AuditPlan.Title,
				"Reference Number":      existingAuditDraftReport.Audit.AuditPlan.ReferenceNumber,
				"AssignedByFirstNameEn": assignedBy.FirstNameEn,
				"AssignedByLastNameEn":  assignedBy.LastNameEn,
				"AssignedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           fmt.Sprintf("%s/audit/audits/%d/draft-report", appURL, auditID),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "audit/templates/email/audit_draft_report_workflow_approve_user_mail.html",
				EmailData:            emailData,
			}

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserID:               userID,
				ModuleIdentifier:     "audit",
				NotificationType:     string(commonDTOs.ReviewType),
				HeadingEn:            "Approve Audit Draft Report",
				HeadingAr:            "مراجعة مسودة تقرير التدقيق",
				TitleEn:              existingAuditDraftReport.Audit.AuditPlan.Title,
				TitleAr:              existingAuditDraftReport.Audit.AuditPlan.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			// Send Notification and Email
			err = adrrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "audit")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}

		// Log activity
		err := adrrs.activityLogSvc.LogActivity(
			models.ApproveAction,
			models.AuditReportReview{}.TableName(),
			id,
			auditReportWorkflowFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.ApproveAction,
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

func (adrrs *AuditDraftReportReviewService) Revert(auditID uint, auditDraftReportID uint, currentUserID uint, auditReportWorkflowFormDTO *auditDTOs.AuditReportWorkflowApproveFormDTO) error {
	tx := adrrs.repo.BeginTransaction()
	if _, err := adrrs.auditSvc.GetByID(auditID); err != nil {
		tx.Rollback()
		return err
	}
	existingAuditDraftReport, err := adrrs.auditDraftReportSvc.GetByID(auditDraftReportID)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, auditDraftReportID)
	}

	existingAuditReportWorkflow, err := adrrs.repo.GetAuditReportWorkflowByDraftReportID(auditDraftReportID, tx)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, auditDraftReportID)
	}

	// Check if current user is a participant
	isParticipant := false
	for _, participant := range existingAuditReportWorkflow.AuditReportWorkflowParticipants {
		if participant.UserID == currentUserID {
			isParticipant = true
			break
		}
	}

	// If the current user is not a participant, return a forbidden error
	if !isParticipant {
		tx.Rollback()
		return utility.NewForbiddenError("The current user is not authorized to perform this action.")
	}

	//check level is = 1
	if existingAuditDraftReport.NextReviewLevel != 1 {
		tx.Rollback()
		return utility.NewForbiddenError("The audit draft report is not at the required review level for this action.")
	}
	auditDraftReportStatusInReview, err := adrrs.auditDraftReportStatusSvc.GetByType("in-review")
	if err != nil {
		tx.Rollback()
		return err
	}
	if existingAuditDraftReport.AuditDraftReportStatusID != auditDraftReportStatusInReview.ID {
		tx.Rollback()
		return utility.NewForbiddenError("The audit draft report is either not in review or has not been submitted.")
	}
	workflowRevertStatus, err := adrrs.workflowStatusSvc.GetByType("revert")
	if err != nil {
		tx.Rollback()
		return err
	}
	// Update the WorkflowStatusID to workflowRevertStatus
	existingAuditReportWorkflow.WorkflowStatusID = workflowRevertStatus.ID

	// Save the updated AuditReportWorkflow
	if err := adrrs.repo.UpdateAuditReportWorkflow(existingAuditReportWorkflow, tx); err != nil {
		tx.Rollback()
		return err
	}
	//add review for history
	auditReportReview := &models.AuditReportReview{
		AuditDraftReportID:    &auditDraftReportID,
		AuditReportWorkflowID: existingAuditReportWorkflow.ID,
		WorkflowStatusID:      workflowRevertStatus.ID,
		UserID:                &currentUserID,
		Comment:               auditReportWorkflowFormDTO.Comment,
		CreatedBYID:           currentUserID,
	}
	id, err := adrrs.repo.Create(auditReportReview, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	//update status & AuditeeReviewSubmitted user details on draft report
	auditDraftReportDraftStatus, err := adrrs.auditDraftReportStatusSvc.GetByType("revert")
	if err != nil {
		tx.Rollback()
		return err
	}
	nextReviewLevel := 0
	err = adrrs.auditDraftReportSvc.StatusUpdate(auditDraftReportID, currentUserID, auditDraftReportDraftStatus.ID, nextReviewLevel, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	//todo - notification to workflow user

	tx.Commit()

	go func() {

		urlParam := fmt.Sprintf(`{"id":%d}`, auditID)
		entity := "audit"
		action := "revert"
		// Create a map to track unique recipients
		recipientsMap := make(map[uint]bool)

		// Add participants excluding the current user
		for _, participant := range existingAuditReportWorkflow.AuditReportWorkflowParticipants {
			if participant.UserID != currentUserID {
				recipientsMap[participant.UserID] = true
			}
		}

		// Add InternalReviewSubmittedBYID
		recipientsMap[*existingAuditDraftReport.InternalReviewSubmittedBYID] = true

		for userID := range recipientsMap {
			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserID:               userID,
				ModuleIdentifier:     "audit",
				NotificationType:     string(commonDTOs.ReviewType),
				HeadingEn:            "Revert Audit Draft Report",
				HeadingAr:            "مراجعة مسودة تقرير التدقيق",
				TitleEn:              existingAuditDraftReport.Audit.AuditPlan.Title,
				TitleAr:              existingAuditDraftReport.Audit.AuditPlan.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			responsibleUser, err := adrrs.userSvc.GetByID(userID)
			if err != nil {
				utility.LogError("Failed to fetch review user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			assignedBy, err := adrrs.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch assigning user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			emailTo := []string{responsibleUser.Email}
			subject := fmt.Sprintf("Audit Draft Report /%s: Revert Notification", existingAuditDraftReport.Audit.AuditPlan.Title)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":           responsibleUser.FirstNameEn,
				"LastNameEn":            responsibleUser.LastNameEn,
				"AuditTitle":            existingAuditDraftReport.Audit.AuditPlan.Title,
				"Reference Number":      existingAuditDraftReport.Audit.AuditPlan.ReferenceNumber,
				"RevertedByFirstNameEn": assignedBy.FirstNameEn,
				"RevertedByLastNameEn":  assignedBy.LastNameEn,
				"RevertedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           fmt.Sprintf("%s/audit/audits/%d/draft-report", appURL, auditID),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "audit/templates/email/audit_draft_report_workflow_revert_user_mail.html",
				EmailData:            emailData,
			}

			// Send Notification and Email
			err = adrrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "audit")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
		err := adrrs.activityLogSvc.LogActivity(
			models.RejectAction,
			models.AuditReportReview{}.TableName(),
			id,
			auditReportWorkflowFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.RejectAction,
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

func (adrrs *AuditDraftReportReviewService) ReSubmit(auditID uint, auditDraftReportID uint, currentUserID uint, auditReportWorkflowFormDTO *auditDTOs.AuditReportWorkflowSubmitFormDTO) error {
	tx := adrrs.repo.BeginTransaction()
	if _, err := adrrs.auditSvc.GetByID(auditID); err != nil {
		return err
	}
	existingAuditDraftReport, err := adrrs.auditDraftReportSvc.GetByID(auditDraftReportID)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, auditDraftReportID)
	}
	auditDraftReportStatusDraft, err := adrrs.auditDraftReportStatusSvc.GetByType("revert")
	if err != nil {
		return err
	}
	if existingAuditDraftReport.AuditDraftReportStatusID != auditDraftReportStatusDraft.ID {
		tx.Rollback()
		return utility.NewForbiddenError("The audit draft report is not in the 'reverted' status.")
	}
	// Check if there is an existing AuditReportWorkflow for the draft report
	existingAuditReportWorkflow, err := adrrs.repo.GetAuditReportWorkflowByDraftReportID(auditDraftReportID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	workflowSubmitStatus, err := adrrs.workflowStatusSvc.GetByType("submit")
	if err != nil {
		return err
	}
	// If the workflow already exists, we will not create a new one
	// In the case of re-submit

	auditReportWorkflowID := existingAuditReportWorkflow.ID
	// Create a map for existing participant user IDs
	existingUserIDs := make(map[uint]bool)
	for _, participant := range existingAuditReportWorkflow.AuditReportWorkflowParticipants {
		existingUserIDs[participant.UserID] = true
	}

	// Create a map for new user IDs from the form DTO
	newUserIDs := make(map[uint]bool)
	for _, userID := range auditReportWorkflowFormDTO.UserIDs {
		newUserIDs[*userID] = true
	}

	// Prepare a slice to hold users to add
	var participantsToAdd []*models.AuditReportWorkflowParticipant
	for _, userID := range auditReportWorkflowFormDTO.UserIDs {
		if _, exists := existingUserIDs[*userID]; !exists {
			participant := &models.AuditReportWorkflowParticipant{
				AuditReportWorkflowID: auditReportWorkflowID,
				UserID:                *userID,
				CreatedBYID:           currentUserID,
			}
			participantsToAdd = append(participantsToAdd, participant)
		}
	}

	// Prepare a slice to hold users to remove
	var participantsToRemove []*models.AuditReportWorkflowParticipant
	for userID := range existingUserIDs {
		if _, exists := newUserIDs[userID]; !exists {
			participantsToRemove = append(participantsToRemove, &models.AuditReportWorkflowParticipant{
				AuditReportWorkflowID: auditReportWorkflowID,
				UserID:                userID,
			})
		}
	}
	// Extract participant IDs to remove
	var participantIDsToRemove []*uint
	for _, user := range participantsToRemove {
		userID := user.UserID
		participantIDsToRemove = append(participantIDsToRemove, &userID)
	}

	// Add new participants
	if len(participantsToAdd) > 0 {
		if err := adrrs.repo.CreateAuditReportWorkflowParticipants(participantsToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Remove old participants
	if len(participantIDsToRemove) > 0 {
		if err := adrrs.repo.DeleteAuditReportWorkflowParticipants(participantIDsToRemove, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	//add review for history
	auditReportReview := &models.AuditReportReview{
		AuditDraftReportID:    &auditDraftReportID,
		AuditReportWorkflowID: auditReportWorkflowID,
		WorkflowStatusID:      workflowSubmitStatus.ID,
		UserID:                &currentUserID,
		Comment:               auditReportWorkflowFormDTO.Comment,
		CreatedBYID:           currentUserID,
	}
	id, err := adrrs.repo.Create(auditReportReview, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	//update status & AuditeeReviewSubmitted user details on draft report
	auditDraftReportInReviewStatus, err := adrrs.auditDraftReportStatusSvc.GetByType("in-review")
	if err != nil {
		return err
	}
	nextReviewLevel := 1

	err = adrrs.auditDraftReportSvc.StatusUpdate(auditDraftReportID, currentUserID, auditDraftReportInReviewStatus.ID, nextReviewLevel, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	//todo - notification to workflow user

	tx.Commit()

	go func() {
		urlParam := fmt.Sprintf(`{"id":%d}`, auditID)
		entity := "audit"
		action := "review-user-assign"

		for _, participant := range auditReportWorkflowFormDTO.UserIDs {
			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserID:               *participant,
				ModuleIdentifier:     "audit",
				NotificationType:     string(commonDTOs.ReviewType),
				HeadingEn:            "Audit Draft Report Review",
				HeadingAr:            "مراجعة مسودة تقرير التدقيق",
				TitleEn:              existingAuditDraftReport.Audit.AuditPlan.Title,
				TitleAr:              existingAuditDraftReport.Audit.AuditPlan.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			responsibleUser, err := adrrs.userSvc.GetByID(*participant)
			if err != nil {
				utility.LogError("Failed to fetch review user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			assignedBy, err := adrrs.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch assigning user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			emailTo := []string{responsibleUser.Email}
			subject := fmt.Sprintf("Audit Finale Report /%s: Review Assigned", existingAuditDraftReport.Audit.AuditPlan.Title)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":           responsibleUser.FirstNameEn,
				"LastNameEn":            responsibleUser.LastNameEn,
				"AuditTitle":            existingAuditDraftReport.Audit.AuditPlan.Title,
				"Reference Number":      existingAuditDraftReport.Audit.AuditPlan.ReferenceNumber,
				"AssignedByFirstNameEn": assignedBy.FirstNameEn,
				"AssignedByLastNameEn":  assignedBy.LastNameEn,
				"AssignedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           fmt.Sprintf("%s/audit/audits/%d/draft-report", appURL, auditID),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "audit/templates/email/audit_draft_report_workflow_re_submit_user_mail.html",
				EmailData:            emailData,
			}

			// Send Notification and Email
			err = adrrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "audit")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
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
