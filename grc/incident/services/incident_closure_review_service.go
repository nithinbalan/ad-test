package services

import (
	"fmt"
	"time"

	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	incidentRepositories "jethur/grc/incident/repositories"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationSvc "jethur/grc/organization/services"
	"jethur/utility"
)

type IncidentClosureReviewServiceInterface interface {
	Submit(closureID uint, currentUserID uint) error
	Approve(closureID uint, currentUserID uint) error
	Revert(closureID uint, currentUserID uint) error
}

type IncidentClosureReviewService struct {
	repo                         incidentRepositories.IncidentClosureReviewRepositoryInterface
	logService                   logServices.ActivityLogServiceInterface
	incidentClosureSvc           IncidentClosureServiceInterface
	incidentClosureWorkflowSvc   IncidentClosureWorkflowServiceInterface
	workflowStatusSvc            masterServices.WorkflowStatusServiceInterface
	userSvc                      organizationSvc.UserServiceInterface
	notificationSvc              commonServices.NotificationServiceInterface
	incidentRegisterSvc          IncidentRegisterServiceInterface
	incidentStatusSvc            IncidentStatusServiceInterface
	incidentActionPlanSvc        IncidentActionPlanServiceInterface
	incidentRootCauseAnalysisSvc IncidentRootCauseAnalysisServiceInterface
}

func NewIncidentClosureReviewService(
	repo incidentRepositories.IncidentClosureReviewRepositoryInterface,
	logService logServices.ActivityLogServiceInterface,
	incidentClosureSvc IncidentClosureServiceInterface,
	incidentClosureWorkflowSvc IncidentClosureWorkflowServiceInterface,
	workflowStatusSvc masterServices.WorkflowStatusServiceInterface,
	userSvc organizationSvc.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	incidentRegisterSvc IncidentRegisterServiceInterface,
	incidentStatusSvc IncidentStatusServiceInterface,
	incidentActionPlanSvc IncidentActionPlanServiceInterface,
	incidentRootCauseAnalysisSvc IncidentRootCauseAnalysisServiceInterface,
) IncidentClosureReviewServiceInterface {
	return &IncidentClosureReviewService{
		repo:                         repo,
		logService:                   logService,
		incidentClosureSvc:           incidentClosureSvc,
		incidentClosureWorkflowSvc:   incidentClosureWorkflowSvc,
		workflowStatusSvc:            workflowStatusSvc,
		userSvc:                      userSvc,
		notificationSvc:              notificationSvc,
		incidentRegisterSvc:          incidentRegisterSvc,
		incidentStatusSvc:            incidentStatusSvc,
		incidentActionPlanSvc:        incidentActionPlanSvc,
		incidentRootCauseAnalysisSvc: incidentRootCauseAnalysisSvc,
	}
}

func (icrs *IncidentClosureReviewService) Submit(closureID uint, currentUserID uint) error {
	incidentClosure, err := icrs.incidentClosureSvc.GetByClosureID(closureID)
	if err != nil {
		return err
	}

	if incidentClosure.SubmittedBYID != nil && incidentClosure.WorkflowStatusID != nil {
		workflowStatus, err := icrs.workflowStatusSvc.GetByID(*incidentClosure.WorkflowStatusID)
		if err != nil {
			return err
		}

		if workflowStatus.Type != "revert" {
			return utility.NewForbiddenError("incident_closure_already_submitted")
		}
	}

	err = icrs.incidentClosureWorkflowSvc.ValidateWorkflowBeforeCreate("incident_closure")
	if err != nil {
		return err
	}

	tx := icrs.repo.BeginTransaction()
	currentDate := time.Now()
	level := 1

	rootCauseAnalyses, err := icrs.incidentRootCauseAnalysisSvc.GetAllByIncidentID(incidentClosure.IncidentRegisterID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	if len(rootCauseAnalyses) == 0 {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("at least one root cause analysis is required before submitting the incident closure"), "incident_closure_submission", "root cause analysis")
	}

	actionPlans, err := icrs.incidentActionPlanSvc.GetAllByIncidentID(incidentClosure.IncidentRegisterID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	if len(actionPlans) == 0 {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("at least one action plan is required before submitting the incident closure"), "incident_closure_submission", "action plan")
	}

	workflowStatus, err := icrs.workflowStatusSvc.GetByType("submit")
	if err != nil {
		tx.Rollback()
		return err
	}

	err = icrs.incidentClosureWorkflowSvc.DeleteAllWorkflowsForIncidentClosure(uint(closureID), tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = icrs.incidentClosureWorkflowSvc.CreateIncidentClosureWorkflow(uint(closureID), "incident_closure", currentUserID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	incidentClosure.WorkflowStatusID = &workflowStatus.ID
	incidentClosure.SubmittedBYID = &currentUserID
	incidentClosure.SubmittedAt = &currentDate
	incidentClosure.NextReviewLevel = &level
	incidentClosure.UpdatedBYID = &currentUserID

	err = icrs.repo.UpdateIncidentClosure(incidentClosure, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// cybersecurityApprovalStatus, err := icrs.incidentStatusSvc.GetByType("cybersecurity_approval")
	// if err != nil {
	// 	tx.Rollback()
	// 	return err
	// }

	// if err := icrs.incidentRegisterSvc.StatusUpdate(incidentClosure.IncidentRegisterID, currentUserID, cybersecurityApprovalStatus.ID, tx); err != nil {
	// 	tx.Rollback()
	// 	return err
	// }

	incidentClosureWorkflowLevel, err := icrs.incidentClosureWorkflowSvc.GetIncidentClosureWorkflowByLevel(incidentClosure.ID, level, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := icrs.createReviewWithHistory(incidentClosureWorkflowLevel.ID, uint(closureID), level, workflowStatus.ID, currentUserID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	go func() {
		closure, err := icrs.incidentClosureSvc.GetFullDetailsByClosureID(closureID)
		if err != nil {
			utility.LogError("Failed to fetch incident closure", map[string]interface{}{"error": err.Error()})
			return
		}

		urlParam := fmt.Sprintf(`{"incidentID":%d}`, closure.IncidentRegister.ID)
		entity := "incident-closure-review"
		action := "review-user-assign"
		incidentRef := fmt.Sprintf("%s--%s", closure.IncidentRegister.ReferenceNumber, closure.IncidentRegister.Title)

		for _, participant := range incidentClosureWorkflowLevel.IncidentClosureWorkflowParticipants {
			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{participant.UserID},
				ModuleIdentifier:     "incident",
				NotificationType:     string(commonDTOs.ReviewType),
				HeadingEn:            "Incident Closure Review",
				HeadingAr:            "مراجعة إغلاق الحادث",
				TitleEn:              fmt.Sprintf("Incident Closure Review Assigned - %s", incidentRef),
				TitleAr:              fmt.Sprintf("تعيين مراجعة إغلاق الحادث - %s", incidentRef),
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			responsibleUser, err := icrs.userSvc.GetByID(participant.UserID)
			if err != nil {
				utility.LogError("Failed to fetch review user", map[string]interface{}{"error": err.Error()})
				continue
			}

			assignedBy, err := icrs.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch assigning user", map[string]interface{}{"error": err.Error()})
				continue
			}

			emailTo := []string{responsibleUser.Email}
			subject := "Incident Closure: Review Assigned"
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":           responsibleUser.FirstNameEn,
				"LastNameEn":            responsibleUser.LastNameEn,
				"Title":                 "Incident Closure",
				"IncidentTitle":         closure.IncidentRegister.Title,
				"IncidentReferenceNo":   closure.IncidentRegister.ReferenceNumber,
				"CompletionDate":        closure.CompletionDateString,
				"ResolutionSummary":     closure.ResolutionSummary,
				"AssignedByFirstNameEn": assignedBy.FirstNameEn,
				"AssignedByLastNameEn":  assignedBy.LastNameEn,
				"AssignedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           fmt.Sprintf("%s/incident/register/%d/incident-closure", appURL, closure.IncidentRegister.ID),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "incident/templates/email/incident_closure_workflow_review_user_mail.html",
				EmailData:            emailData,
			}

			if err := icrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "incident_closure"); err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{"error": err.Error()})
			}
		}
	}()

	return nil
}

func (icrs *IncidentClosureReviewService) Approve(closureID uint, currentUserID uint) error {
	incidentClosure, err := icrs.incidentClosureSvc.GetByClosureID(closureID)
	if err != nil {
		return err
	}

	tx := icrs.repo.BeginTransaction()

	if incidentClosure.NextReviewLevel == nil {
		tx.Rollback()
		return utility.NewForbiddenError("review level not found")
	}

	currentApprovalLevel := incidentClosure.NextReviewLevel
	currentWorkflow, err := icrs.incidentClosureWorkflowSvc.GetIncidentClosureWorkflowByLevel(incidentClosure.ID, *currentApprovalLevel, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	workflowStatus, err := icrs.workflowStatusSvc.GetByType("approve")
	if err != nil {
		tx.Rollback()
		return err
	}

	switch currentWorkflow.WorkflowLevelType.Type {
	case "individual", "hod":
		if len(currentWorkflow.IncidentClosureWorkflowParticipants) == 0 ||
			currentUserID != currentWorkflow.IncidentClosureWorkflowParticipants[0].UserID {
			tx.Rollback()
			return utility.NewForbiddenError("you have no permission to approve")
		}

		if err := icrs.createReviewWithHistory(currentWorkflow.ID, incidentClosure.ID, currentWorkflow.Level, workflowStatus.ID, currentUserID, tx); err != nil {
			tx.Rollback()
			return err
		}

		for _, participant := range currentWorkflow.IncidentClosureWorkflowParticipants {
			if participant.UserID == currentUserID {
				participant.WorkflowStatusID = workflowStatus.ID
				if err := icrs.repo.UpdateIncidentClosureWorkflowParticipant(participant, tx); err != nil {
					tx.Rollback()
					return err
				}
				break
			}
		}

		incidentClosureWorkflow, err := icrs.incidentClosureWorkflowSvc.GetIncidentClosureWorkflowByID(currentWorkflow.ID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		incidentClosureWorkflow.WorkflowStatusID = workflowStatus.ID
		if err := icrs.incidentClosureWorkflowSvc.UpdateIncidentClosureWorkflow(incidentClosureWorkflow, tx); err != nil {
			tx.Rollback()
			return err
		}

		nextWorkflowLevel, err := icrs.incidentClosureWorkflowSvc.GetNextLevel(incidentClosure.ID, *currentApprovalLevel, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		if nextWorkflowLevel == nil {
			approvedAt := time.Now()
			incidentClosure.WorkflowStatusID = &workflowStatus.ID
			incidentClosure.ApprovedBYID = &currentUserID
			incidentClosure.ApprovedAt = &approvedAt
			incidentClosure.NextReviewLevel = nil

			closedStatus, err := icrs.incidentStatusSvc.GetByType("closed")
			if err != nil {
				tx.Rollback()
				return err
			}

			if err := icrs.incidentRegisterSvc.StatusUpdate(incidentClosure.IncidentRegisterID, currentUserID, closedStatus.ID, tx); err != nil {
				tx.Rollback()
				return err
			}
		} else {
			incidentClosure.NextReviewLevel = &nextWorkflowLevel.Level
		}

		if err := icrs.repo.UpdateIncidentClosure(incidentClosure, tx); err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}

		if nextWorkflowLevel != nil {
			nextWorkflow, err := icrs.incidentClosureWorkflowSvc.GetIncidentClosureWorkflowByLevel(incidentClosure.ID, nextWorkflowLevel.Level)
			if err != nil {
				return err
			}
			icrs.notifyNextLevelUsers(incidentClosure.ID, nextWorkflow, currentUserID)
		}

		icrs.sendApprovalNotification(incidentClosure.ID, currentUserID)

	case "multiple-users", "role":
		allowed := false
		for _, participant := range currentWorkflow.IncidentClosureWorkflowParticipants {
			if participant.UserID == currentUserID {
				allowed = true
				break
			}
		}
		if !allowed {
			tx.Rollback()
			return utility.NewForbiddenError("you have no permission to approve")
		}

		if err := icrs.createReviewWithHistory(currentWorkflow.ID, incidentClosure.ID, currentWorkflow.Level, workflowStatus.ID, currentUserID, tx); err != nil {
			tx.Rollback()
			return err
		}

		for _, participant := range currentWorkflow.IncidentClosureWorkflowParticipants {
			if participant.UserID == currentUserID {
				participant.WorkflowStatusID = workflowStatus.ID
				if err := icrs.repo.UpdateIncidentClosureWorkflowParticipant(participant, tx); err != nil {
					tx.Rollback()
					return err
				}
				break
			}
		}

		if currentWorkflow.ApprovalTypes == "all" {
			allApproved, err := icrs.incidentClosureWorkflowSvc.CheckAllParticipantsApproved(currentWorkflow.ID, currentWorkflow.IncidentClosureWorkflowParticipants, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
			if !allApproved {
				return tx.Commit()
			}
		}

		incidentClosureWorkflow, err := icrs.incidentClosureWorkflowSvc.GetIncidentClosureWorkflowByID(currentWorkflow.ID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
		incidentClosureWorkflow.WorkflowStatusID = workflowStatus.ID
		if err := icrs.incidentClosureWorkflowSvc.UpdateIncidentClosureWorkflow(incidentClosureWorkflow, tx); err != nil {
			tx.Rollback()
			return err
		}

		nextWorkflowLevel, err := icrs.incidentClosureWorkflowSvc.GetNextLevel(incidentClosure.ID, *currentApprovalLevel, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		if nextWorkflowLevel == nil {
			approvedAt := time.Now()
			incidentClosure.WorkflowStatusID = &workflowStatus.ID
			incidentClosure.ApprovedBYID = &currentUserID
			incidentClosure.ApprovedAt = &approvedAt
			incidentClosure.NextReviewLevel = nil

			// closedStatus, err := icrs.incidentStatusSvc.GetByType("closed")
			// if err != nil {
			// 	tx.Rollback()
			// 	return err
			// }

			// if err := icrs.incidentRegisterSvc.StatusUpdate(incidentClosure.IncidentRegisterID, currentUserID, closedStatus.ID, tx); err != nil {
			// 	tx.Rollback()
			// 	return err
			// }
		} else {
			incidentClosure.NextReviewLevel = &nextWorkflowLevel.Level
		}

		if err := icrs.repo.UpdateIncidentClosure(incidentClosure, tx); err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}

		if nextWorkflowLevel != nil {
			nextWorkflow, err := icrs.incidentClosureWorkflowSvc.GetIncidentClosureWorkflowByLevel(incidentClosure.ID, nextWorkflowLevel.Level)
			if err != nil {
				return err
			}
			icrs.notifyNextLevelUsers(incidentClosure.ID, nextWorkflow, currentUserID)
		}

		icrs.sendApprovalNotification(incidentClosure.ID, currentUserID)
	}

	return nil
}

func (icrs *IncidentClosureReviewService) Revert(closureID uint, currentUserID uint) error {
	workflowStatus, err := icrs.workflowStatusSvc.GetByType("revert")
	if err != nil {
		return err
	}

	tx := icrs.repo.BeginTransaction()

	incidentClosure, err := icrs.incidentClosureSvc.GetByClosureID(closureID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if incidentClosure.NextReviewLevel == nil {
		tx.Rollback()
		return utility.NewForbiddenError("review level not found")
	}

	currentApprovalLevel := incidentClosure.NextReviewLevel
	currentWorkflow, err := icrs.incidentClosureWorkflowSvc.GetIncidentClosureWorkflowByLevel(incidentClosure.ID, *currentApprovalLevel, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := icrs.createReviewWithHistory(currentWorkflow.ID, incidentClosure.ID, currentWorkflow.Level, workflowStatus.ID, currentUserID, tx); err != nil {
		tx.Rollback()
		return err
	}

	incidentClosureWorkflow, err := icrs.incidentClosureWorkflowSvc.GetIncidentClosureWorkflowByID(currentWorkflow.ID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	incidentClosureWorkflow.WorkflowStatusID = workflowStatus.ID
	if err := icrs.incidentClosureWorkflowSvc.UpdateIncidentClosureWorkflow(incidentClosureWorkflow, tx); err != nil {
		tx.Rollback()
		return err
	}

	incidentClosure.WorkflowStatusID = &workflowStatus.ID
	incidentClosure.NextReviewLevel = nil
	incidentClosure.SubmittedAt = nil
	incidentClosure.SubmittedBYID = nil
	incidentClosure.UpdatedBYID = &currentUserID
	if err := icrs.repo.UpdateIncidentClosure(incidentClosure, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	icrs.sendRevertedNotification(incidentClosure.ID, currentUserID)

	return nil
}

func (icrs *IncidentClosureReviewService) createReviewWithHistory(workflowID uint, incidentClosureID uint, level int, statusID uint, userID uint, tx utility.DatabaseTransaction) error {
	workflowReview := &models.IncidentClosureWorkflowReview{
		IncidentClosureWorkflowID: workflowID,
		WorkflowStatusID:          statusID,
		CreatedBYID:               userID,
	}
	if err := icrs.repo.CreateIncidentClosureReview(workflowReview, tx); err != nil {
		return err
	}

	historyReview := &models.IncidentClosureWorkflowHistory{
		IncidentClosureID: incidentClosureID,
		Level:             level,
		WorkflowStatusID:  statusID,
		CreatedBYID:       userID,
		CreatedAt:         time.Now(),
	}
	if err := icrs.repo.CreateIncidentClosureWorkflowHistory(historyReview, tx); err != nil {
		return err
	}

	return nil
}

func (icrs *IncidentClosureReviewService) sendApprovalNotification(closureID uint, currentUserID uint) {
	go func() {
		incidentClosure, err := icrs.incidentClosureSvc.GetFullDetailsByClosureID(closureID)
		if err != nil {
			utility.LogError("Failed to fetch incident closure", map[string]interface{}{"error": err.Error()})
			return
		}

		urlParam := fmt.Sprintf(`{"incidentID":%d}`, incidentClosure.IncidentRegister.ID)
		entity := "incident-closure-review"
		action := "approval"
		incidentRef := fmt.Sprintf("%s--%s", incidentClosure.IncidentRegister.ReferenceNumber, incidentClosure.IncidentRegister.Title)

		var createdBYID uint
		if incidentClosure.CreatedBYID != nil {
			createdBYID = *incidentClosure.CreatedBYID
		}

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{createdBYID},
			ModuleIdentifier:     "incident",
			NotificationType:     string(commonDTOs.ApprovalType),
			HeadingEn:            "Incident Closure Approval",
			HeadingAr:            "تمت الموافقة على إغلاق الحادث",
			TitleEn:              fmt.Sprintf("Incident Closure Approval - %s", incidentRef),
			TitleAr:              fmt.Sprintf("تمت الموافقة على إغلاق الحادث - %s", incidentRef),
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		createdByUser, err := icrs.userSvc.GetByID(createdBYID)
		if err != nil {
			utility.LogError("Failed to fetch closure creator", map[string]interface{}{"error": err.Error()})
			return
		}

		approvedByUser, err := icrs.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch approving user", map[string]interface{}{"error": err.Error()})
			return
		}

		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstNameEn":           createdByUser.FirstNameEn,
			"LastNameEn":            createdByUser.LastNameEn,
			"IncidentTitle":         incidentClosure.IncidentRegister.Title,
			"IncidentReferenceNo":   incidentClosure.IncidentRegister.ReferenceNumber,
			"CompletionDate":        incidentClosure.CompletionDateString,
			"ResolutionSummary":     incidentClosure.ResolutionSummary,
			"ApprovedByFirstNameEn": approvedByUser.FirstNameEn,
			"ApprovedByLastNameEn":  approvedByUser.LastNameEn,
			"ApprovedAt":            time.Now().Format("2006-01-02"),
			"RedirectURL":           fmt.Sprintf("%s/incident/register/%d/incident-closure", appURL, incidentClosure.IncidentRegister.ID),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              []string{createdByUser.Email},
			Subject:              "Incident Closure Approval Notification",
			TemplateRelativePath: "incident/templates/email/incident_closure_workflow_approve_user_mail.html",
			EmailData:            emailData,
		}

		if err := icrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "incident_closure"); err != nil {
			utility.LogError("Failed to send approval notification email", map[string]interface{}{"error": err.Error()})
		}
	}()
}

func (icrs *IncidentClosureReviewService) notifyNextLevelUsers(closureID uint, nextWorkflowLevel *models.IncidentClosureWorkflow, currentUserID uint) {
	go func() {
		incidentClosure, err := icrs.incidentClosureSvc.GetFullDetailsByClosureID(closureID)
		if err != nil {
			utility.LogError("Failed to fetch incident closure", map[string]interface{}{"error": err.Error()})
			return
		}

		var createdBYID uint
		if incidentClosure.CreatedBYID != nil {
			createdBYID = *incidentClosure.CreatedBYID
		}

		urlParam := fmt.Sprintf(`{"incidentID":%d}`, incidentClosure.IncidentRegister.ID)
		entity := "incident-closure-review"
		action := "review-user-assign"
		incidentRef := fmt.Sprintf("%s--%s", incidentClosure.IncidentRegister.ReferenceNumber, incidentClosure.IncidentRegister.Title)

		for _, participant := range nextWorkflowLevel.IncidentClosureWorkflowParticipants {
			responsibleUser, err := icrs.userSvc.GetByID(participant.UserID)
			if err != nil {
				utility.LogError("Failed to fetch review user", map[string]interface{}{"error": err.Error()})
				continue
			}

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{participant.UserID},
				ModuleIdentifier:     "incident",
				NotificationType:     string(commonDTOs.ReviewType),
				HeadingEn:            "Incident Closure Review",
				HeadingAr:            "مراجعة إغلاق الحادث",
				TitleEn:              fmt.Sprintf("Incident Closure Review Assigned - %s", incidentRef),
				TitleAr:              fmt.Sprintf("تعيين مراجعة إغلاق الحادث - %s", incidentRef),
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &createdBYID,
			}

			assignedBy, err := icrs.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch assigning user", map[string]interface{}{"error": err.Error()})
				continue
			}

			createdByUser, err := icrs.userSvc.GetByID(createdBYID)
			if err != nil {
				utility.LogError("Failed to fetch closure creator", map[string]interface{}{"error": err.Error()})
				continue
			}

			emailTo := []string{responsibleUser.Email}
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":           responsibleUser.FirstNameEn,
				"LastNameEn":            responsibleUser.LastNameEn,
				"Title":                 "Incident Closure",
				"IncidentTitle":         incidentClosure.IncidentRegister.Title,
				"IncidentReferenceNo":   incidentClosure.IncidentRegister.ReferenceNumber,
				"CompletionDate":        incidentClosure.CompletionDateString,
				"ResolutionSummary":     incidentClosure.ResolutionSummary,
				"ResponsibleUserName":   createdByUser.FirstNameEn + " " + createdByUser.LastNameEn,
				"ResponsibleUserEmail":  createdByUser.Email,
				"AssignedByFirstNameEn": assignedBy.FirstNameEn,
				"AssignedByLastNameEn":  assignedBy.LastNameEn,
				"AssignedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           fmt.Sprintf("%s/incident/register/%d/incident-closure", appURL, incidentClosure.IncidentRegister.ID),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              "Incident Closure: Review Assigned",
				TemplateRelativePath: "incident/templates/email/incident_closure_workflow_review_user_mail.html",
				EmailData:            emailData,
			}

			if err := icrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "incident_closure"); err != nil {
				utility.LogError("Failed to send next level review notification", map[string]interface{}{"error": err.Error()})
			}
		}
	}()
}

func (icrs *IncidentClosureReviewService) sendRevertedNotification(closureID uint, currentUserID uint) {
	go func() {
		incidentClosure, err := icrs.incidentClosureSvc.GetFullDetailsByClosureID(closureID)
		if err != nil {
			utility.LogError("Failed to fetch incident closure", map[string]interface{}{"error": err.Error()})
			return
		}

		var createdBYID uint
		if incidentClosure.CreatedBYID != nil {
			createdBYID = *incidentClosure.CreatedBYID
		}

		createdByUser, err := icrs.userSvc.GetByID(createdBYID)
		if err != nil {
			utility.LogError("Failed to fetch closure creator", map[string]interface{}{"error": err.Error()})
			return
		}

		urlParam := fmt.Sprintf(`{"incidentID":%d}`, incidentClosure.IncidentRegister.ID)
		entity := "incident-closure-review"
		action := "reverted"
		incidentRef := fmt.Sprintf("%s--%s", incidentClosure.IncidentRegister.ReferenceNumber, incidentClosure.IncidentRegister.Title)

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{createdBYID},
			ModuleIdentifier:     "incident",
			NotificationType:     string(commonDTOs.RejectType),
			HeadingEn:            "Incident Closure Reverted",
			HeadingAr:            "تم رفض إغلاق الحادث",
			TitleEn:              fmt.Sprintf("Incident Closure Reverted - %s", incidentRef),
			TitleAr:              fmt.Sprintf("تم رفض إغلاق الحادث - %s", incidentRef),
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		rejectedByUser, err := icrs.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch review performed user", map[string]interface{}{"error": err.Error()})
			return
		}

		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstNameEn":           createdByUser.FirstNameEn,
			"LastNameEn":            createdByUser.LastNameEn,
			"Title":                 "Incident Closure",
			"IncidentTitle":         incidentClosure.IncidentRegister.Title,
			"IncidentReferenceNo":   incidentClosure.IncidentRegister.ReferenceNumber,
			"CompletionDate":        incidentClosure.CompletionDateString,
			"ResolutionSummary":     incidentClosure.ResolutionSummary,
			"RejectedByFirstNameEn": rejectedByUser.FirstNameEn,
			"RejectedByLastNameEn":  rejectedByUser.LastNameEn,
			"RejectedAt":            time.Now().Format("2006-01-02"),
			"RedirectURL":           fmt.Sprintf("%s/incident/register/%d/incident-closure", appURL, incidentClosure.IncidentRegister.ID),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              []string{createdByUser.Email},
			Subject:              "Incident Closure: Reverted Notification",
			TemplateRelativePath: "incident/templates/email/incident_closure_workflow_revert_user_mail.html",
			EmailData:            emailData,
		}

		if err := icrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "incident_closure"); err != nil {
			utility.LogError("Failed to send reverted notification email", map[string]interface{}{"error": err.Error()})
		}
	}()
}
