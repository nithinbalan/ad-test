package services

import (
	"fmt"
	"time"

	"jethur/grc/audit-insights/dtos"
	"jethur/grc/audit-insights/repositories"
	auditDTOs "jethur/grc/audit/dtos"
	auditServices "jethur/grc/audit/services"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"
)

type AuditInsightsDraftFindingServiceInterface interface {
	Accept(auditDraftFindingID uint, currentUserID uint, auditInsightsDraftFindingApproveFormDTO *dtos.AuditInsightsDraftFindingApproveFormDTO) error
	Challenge(auditDraftFindingID uint, currentUserID uint, auditInsightsDraftFindingChallengeFormDTO *dtos.AuditInsightsDraftFindingChallengeFormDTO) error
}

type AuditInsightsDraftFindingService struct {
	repo                             repositories.AuditInsightsDraftFindingRepositoryInterface
	activityLogSvc                   logServices.ActivityLogServiceInterface
	userSvc                          organizationServices.UserServiceInterface
	auditDraftReportSvc              auditServices.AuditDraftReportServiceInterface
	auditDraftReportStatusSvc        auditServices.AuditDraftReportStatusServiceInterface
	auditDraftFindingSvc             auditServices.AuditDraftFindingServiceInterface
	auditDraftFindingStatusSvc       auditServices.AuditDraftFindingStatusServiceInterface
	auditDraftFindingConversationSvc auditServices.AuditDraftFindingConversationServiceInterface
	auditDraftFindingActionPlanSvc   auditServices.AuditDraftFindingActionPlanServiceInterface
	notificationSvc                  commonServices.NotificationServiceInterface
}

func NewAuditInsightsDraftFindingService(
	repo repositories.AuditInsightsDraftFindingRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	auditDraftReportSvc auditServices.AuditDraftReportServiceInterface,
	auditDraftReportStatusSvc auditServices.AuditDraftReportStatusServiceInterface,
	auditDraftFindingSvc auditServices.AuditDraftFindingServiceInterface,
	auditDraftFindingStatusSvc auditServices.AuditDraftFindingStatusServiceInterface,
	auditDraftFindingConversationSvc auditServices.AuditDraftFindingConversationServiceInterface,
	auditDraftFindingActionPlanSvc auditServices.AuditDraftFindingActionPlanServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) AuditInsightsDraftFindingServiceInterface {
	return &AuditInsightsDraftFindingService{
		repo:                             repo,
		activityLogSvc:                   activityLogSvc,
		userSvc:                          userSvc,
		auditDraftReportSvc:              auditDraftReportSvc,
		auditDraftReportStatusSvc:        auditDraftReportStatusSvc,
		auditDraftFindingSvc:             auditDraftFindingSvc,
		auditDraftFindingStatusSvc:       auditDraftFindingStatusSvc,
		auditDraftFindingConversationSvc: auditDraftFindingConversationSvc,
		auditDraftFindingActionPlanSvc:   auditDraftFindingActionPlanSvc,
		notificationSvc:                  notificationSvc,
	}
}

func (aidfs *AuditInsightsDraftFindingService) Accept(auditDraftFindingID uint, currentUserID uint, auditInsightsDraftFindingApproveFormDTO *dtos.AuditInsightsDraftFindingApproveFormDTO) error {
	existingAuditDraftFinding, err := aidfs.auditDraftFindingSvc.GetFullDetailsByID(auditDraftFindingID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, auditDraftFindingID)
	}
	orgUnitID, err := aidfs.userSvc.GetUserPrimaryUnitID(currentUserID)
	if err != nil {
		return err
	}
	if existingAuditDraftFinding.AuditableItem.ResponsibleUnitID != orgUnitID {
		return utility.NewForbiddenError("You are not authorized to accept this draft finding. Only users from the responsible unit can perform this action.")
	}

	// Retrieve the AuditDraftReport using AuditID from AuditDraftFinding.Audit
	auditID := existingAuditDraftFinding.AuditID
	auditDraftReport, err := aidfs.auditDraftReportSvc.GetByAuditID(auditID)
	if err != nil {
		return err
	}

	// Check if AuditDraftReportStatus is 'InDepartmentReview'
	inReviewStatus, err := aidfs.auditDraftReportStatusSvc.GetByType("in-department-review")
	if err != nil {
		return err
	}
	if auditDraftReport.AuditDraftReportStatusID != inReviewStatus.ID {
		return utility.NewForbiddenError("Action not allowed: Only audit draft reports with an 'In Department Review' status can be accepted.")
	}

	// auditDraftFindingStatus, err := aidfs.auditDraftFindingStatusSvc.GetByType("open")
	// if err != nil {
	// 	return err
	// }
	// if existingAuditDraftFinding.AuditDraftFindingStatusID != auditDraftFindingStatus.ID {
	// 	return utility.NewForbiddenError("Action not allowed: Only audit draft findings with an 'open' status can be processed. Please ensure the audit draft finding is open before proceeding.")
	// }
	if existingAuditDraftFinding.AuditDraftFindingStatus.IsFinalized {
		return utility.NewForbiddenError("Action not allowed: Only audit draft findings with an 'open' status can be processed. Please ensure the audit draft finding is not finalized before proceeding.")
	}
	tx := aidfs.repo.BeginTransaction()

	conversation := &models.AuditDraftFindingConversation{
		AuditDraftFindingID: auditDraftFindingID,
		Comment:             auditInsightsDraftFindingApproveFormDTO.Comment,
		Type:                "DR",
		IsAccepted:          true,
		CreatedBYID:         currentUserID,
	}
	_, err = aidfs.repo.CreateConversation(conversation, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	auditDraftFindingStatusAccepted, err := aidfs.auditDraftFindingStatusSvc.GetByType("accepted")
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update Audit Draft Finding Status
	if err := aidfs.auditDraftFindingSvc.StatusUpdate(auditDraftFindingID, currentUserID, auditDraftFindingStatusAccepted.ID, tx); err != nil {
		tx.Rollback()
		return err
	}
	auditDraftFindingActionPlanFormDTO := auditDTOs.AuditDraftFindingActionPlanFormDTO{
		Title:             auditInsightsDraftFindingApproveFormDTO.Title,
		Description:       auditInsightsDraftFindingApproveFormDTO.Description,
		ResponsibleUserID: auditInsightsDraftFindingApproveFormDTO.ResponsibleUserID,
		TargetDate:        auditInsightsDraftFindingApproveFormDTO.TargetDate,
	}

	if _, err := aidfs.auditDraftFindingActionPlanSvc.Create(auditDraftFindingID, currentUserID, &auditDraftFindingActionPlanFormDTO, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		urlParam := fmt.Sprintf(`{"id":%d}`, auditDraftReport.ID)
		entity := "audit-draft-report"
		action := "approval"

		// Create a map to track unique recipients
		recipientsMap := make(map[uint]bool)

		// Add participants excluding the current user
		for _, participant := range auditDraftReport.Audit.AuditPlan.Members {
			if participant.UserID != currentUserID {
				recipientsMap[participant.UserID] = true
			}
		}

		// Add InternalReviewSubmittedBYID
		recipientsMap[*auditDraftReport.Audit.AuditPlan.TeamLeadID] = true

		// Extract unique recipients into a slice
		var uniqueRecipients []uint
		for userID := range recipientsMap {
			uniqueRecipients = append(uniqueRecipients, userID)
		}

		// Now `uniqueRecipients` contains all unique UserIDs

		// Iterate through recipients
		for _, userID := range uniqueRecipients {
			responsibleUser, err := aidfs.userSvc.GetByID(userID)
			if err != nil {
				utility.LogError("Failed to fetch review user", map[string]interface{}{
					"error": err.Error(),
				})
				continue
			}

			assignedBy, err := aidfs.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch assigning user", map[string]interface{}{
					"error": err.Error(),
				})
				continue
			}

			emailTo := []string{responsibleUser.Email}
			subject := fmt.Sprintf("Audit Draft Report /%s: Findings Accepted Notification", auditDraftReport.Audit.AuditPlan.Title)
			appURL := utility.GetAppURL()

			emailData := map[string]interface{}{
				"FirstNameEn":           responsibleUser.FirstNameEn,
				"LastNameEn":            responsibleUser.LastNameEn,
				"AuditTitle":            auditDraftReport.Audit.AuditPlan.Title,
				"FindingTitle":          existingAuditDraftFinding.Title,
				"FindingDescription":    existingAuditDraftFinding.Description,
				"FindingCategoryTitle":  existingAuditDraftFinding.AuditFindingCategory.TitleEn,
				"AcceptedByFirstNameEn": assignedBy.FirstNameEn,
				"AcceptedByLastNameEn":  assignedBy.LastNameEn,
				"AcceptedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           fmt.Sprintf("%s/audit/audits/%d/draft-report", appURL, auditID),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "audit/templates/email/audit_draft_report_finding_accept_user_mail.html",
				EmailData:            emailData,
			}

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{userID},
				ModuleIdentifier:     "audit_insights",
				NotificationType:     string(commonDTOs.ReviewType),
				HeadingEn:            "Approve Audit Draft Report Finding",
				HeadingAr:            "قبول نتائج مسودة تقرير التدقيق",
				TitleEn:              existingAuditDraftFinding.Title,
				TitleAr:              existingAuditDraftFinding.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			// Send Notification and Email
			err = aidfs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "audit")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
		err := aidfs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditDraftFindingConversation{}.TableName(),
			auditDraftFindingID,
			auditInsightsDraftFindingApproveFormDTO,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AuditDraftFindingConversation{}.TableName(),
				"recordID":  auditDraftFindingID,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (aidfs *AuditInsightsDraftFindingService) Challenge(auditDraftFindingID uint, currentUserID uint, auditInsightsDraftFindingChallengeFormDTO *dtos.AuditInsightsDraftFindingChallengeFormDTO) error {
	existingAuditDraftFinding, err := aidfs.auditDraftFindingSvc.GetFullDetailsByID(auditDraftFindingID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, auditDraftFindingID)
	}
	orgUnitID, err := aidfs.userSvc.GetUserPrimaryUnitID(currentUserID)
	if err != nil {
		return err
	}

	if existingAuditDraftFinding.AuditableItem.ResponsibleUnitID != orgUnitID {
		return utility.NewForbiddenError("You are not authorized to accept this draft finding. Only users from the responsible unit can perform this action.")
	}

	// Retrieve the AuditDraftReport using AuditID from AuditDraftFinding.Audit
	auditID := existingAuditDraftFinding.AuditID
	auditDraftReport, err := aidfs.auditDraftReportSvc.GetByAuditID(auditID)
	if err != nil {
		return err
	}

	// Check if AuditDraftReportStatus is 'InDepartmentReview'
	inReviewStatus, err := aidfs.auditDraftReportStatusSvc.GetByType("in-department-review")
	if err != nil {
		return err
	}
	if auditDraftReport.AuditDraftReportStatusID != inReviewStatus.ID {
		return utility.NewForbiddenError("Action not allowed: Only audit draft reports with an 'In Department Review' status can be challenged.")
	}

	// auditDraftFindingStatus, err := aidfs.auditDraftFindingStatusSvc.GetByType("open")
	// if err != nil {
	// 	return err
	// }
	// if existingAuditDraftFinding.AuditDraftFindingStatusID != auditDraftFindingStatus.ID {
	// 	return utility.NewForbiddenError("Action not allowed: Only audit draft findings with an 'open' status can be processed. Please ensure the audit draft finding is open before proceeding.")
	// }
	if existingAuditDraftFinding.AuditDraftFindingStatus.IsFinalized {
		return utility.NewForbiddenError("Action not allowed: Only audit draft findings with an 'open' status can be processed. Please ensure the audit draft finding is not finalized before proceeding.")
	}
	tx := aidfs.repo.BeginTransaction()

	conversation := &models.AuditDraftFindingConversation{
		AuditDraftFindingID: auditDraftFindingID,
		Comment:             auditInsightsDraftFindingChallengeFormDTO.Comment,
		Type:                "DR",
		CreatedBYID:         currentUserID,
	}
	_, err = aidfs.repo.CreateConversation(conversation, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	auditDraftFindingStatusChallenged, err := aidfs.auditDraftFindingStatusSvc.GetByType("challenged")
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update Audit Draft Finding Status
	if err := aidfs.auditDraftFindingSvc.StatusUpdate(auditDraftFindingID, currentUserID, auditDraftFindingStatusChallenged.ID); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	go func() {
		urlParam := fmt.Sprintf(`{"id":%d}`, auditDraftReport.ID)
		entity := "audit-draft-report"
		action := "challenge"

		// Create a map to track unique recipients
		recipientsMap := make(map[uint]bool)

		// Add participants excluding the current user
		for _, participant := range auditDraftReport.Audit.AuditPlan.Members {
			if participant.UserID != currentUserID {
				recipientsMap[participant.UserID] = true
			}
		}

		// Add InternalReviewSubmittedBYID
		recipientsMap[*auditDraftReport.Audit.AuditPlan.TeamLeadID] = true

		// Extract unique recipients into a slice
		var uniqueRecipients []uint
		for userID := range recipientsMap {
			uniqueRecipients = append(uniqueRecipients, userID)
		}

		// Now `uniqueRecipients` contains all unique UserIDs

		// Iterate through recipients
		for _, userID := range uniqueRecipients {
			responsibleUser, err := aidfs.userSvc.GetByID(userID)
			if err != nil {
				utility.LogError("Failed to fetch review user", map[string]interface{}{
					"error": err.Error(),
				})
				continue
			}

			assignedBy, err := aidfs.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch assigning user", map[string]interface{}{
					"error": err.Error(),
				})
				continue
			}

			emailTo := []string{responsibleUser.Email}
			subject := fmt.Sprintf("Audit Draft Report /%s: Findings Challenged Notification", auditDraftReport.Audit.AuditPlan.Title)
			appURL := utility.GetAppURL()

			emailData := map[string]interface{}{
				"FirstNameEn":             responsibleUser.FirstNameEn,
				"LastNameEn":              responsibleUser.LastNameEn,
				"AuditTitle":              auditDraftReport.Audit.AuditPlan.Title,
				"FindingTitle":            existingAuditDraftFinding.Title,
				"FindingDescription":      existingAuditDraftFinding.Description,
				"FindingCategoryTitle":    existingAuditDraftFinding.AuditFindingCategory.TitleEn,
				"ChallengedByFirstNameEn": assignedBy.FirstNameEn,
				"ChallengedByLastNameEn":  assignedBy.LastNameEn,
				"ChallengedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":             fmt.Sprintf("%s/audit/audits/%d/draft-report", appURL, auditID),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "audit/templates/email/audit_draft_report_finding_challenge_user_mail.html",
				EmailData:            emailData,
			}

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{userID},
				ModuleIdentifier:     "audit_insights",
				NotificationType:     string(commonDTOs.ReviewType),
				HeadingEn:            "Challenge Audit Draft Report Finding",
				HeadingAr:            "تحدي نتائج مسودة تقرير التدقيق",
				TitleEn:              existingAuditDraftFinding.Title,
				TitleAr:              existingAuditDraftFinding.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			// Send Notification and Email
			err = aidfs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "audit")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
		err := aidfs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditDraftFindingConversation{}.TableName(),
			auditDraftFindingID,
			auditInsightsDraftFindingChallengeFormDTO,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AuditDraftFindingConversation{}.TableName(),
				"recordID":  auditDraftFindingID,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}
