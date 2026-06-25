package services

import (
	"fmt"
	"time"

	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/repositories"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationSvc "jethur/grc/organization/services"
	taskServices "jethur/grc/task/services"
	"jethur/utility"
)

type ComplianceRegisterReviewServiceInterface interface {
	Submit(complianceRegisterID uint, currentUserID uint) error
	Approve(complianceRegisterID uint, currentUserID uint) error
	Reject(complianceRegisterID uint, currentUserID uint) error
	Reopen(complianceRegisterID uint, currentUserID uint) error
}

type ComplianceRegisterReviewService struct {
	repo                          repositories.ComplianceRegisterReviewRepositoryInterface
	logService                    logServices.ActivityLogServiceInterface
	complianceRegisterSvc         ComplianceRegisterServiceInterface
	complianceRegisterWorkflowSvc ComplianceRegisterWorkflowServiceInterface
	workflowStatusSvc             masterServices.WorkflowStatusServiceInterface
	userSvc                       organizationSvc.UserServiceInterface
	notificationSvc               commonServices.NotificationServiceInterface
	taskSvc                       taskServices.TaskServiceInterface
}

func NewComplianceRegisterReviewService(
	repo repositories.ComplianceRegisterReviewRepositoryInterface,
	logService logServices.ActivityLogServiceInterface,
	complianceRegisterSvc ComplianceRegisterServiceInterface,
	complianceRegisterWorkflowSvc ComplianceRegisterWorkflowServiceInterface,
	workflowStatusSvc masterServices.WorkflowStatusServiceInterface,
	userSvc organizationSvc.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	taskSvc taskServices.TaskServiceInterface,
) ComplianceRegisterReviewServiceInterface {
	return &ComplianceRegisterReviewService{
		repo:                          repo,
		logService:                    logService,
		complianceRegisterSvc:         complianceRegisterSvc,
		complianceRegisterWorkflowSvc: complianceRegisterWorkflowSvc,
		workflowStatusSvc:             workflowStatusSvc,
		userSvc:                       userSvc,
		notificationSvc:               notificationSvc,
		taskSvc:                       taskSvc,
	}
}

func (crrs *ComplianceRegisterReviewService) Submit(complianceRegisterID uint, currentUserID uint) error {
	// Get full details for validation and notifications
	complianceRegisterDetail, err := crrs.complianceRegisterSvc.GetFullDetailsByID(complianceRegisterID)
	if err != nil {
		return err
	}

	// Get the model for database operations
	complianceRegister, err := crrs.complianceRegisterSvc.GetByID(complianceRegisterID)
	if err != nil {
		return err
	}

	totalRequirements := 0
	for _, clause := range complianceRegisterDetail.ComplianceClauses {
		totalRequirements += len(clause.ComplianceRequirements)
	}
	if totalRequirements == 0 {
		return utility.NewValidationError(fmt.Errorf("compliance register has no requirements"), "required", "compliance requirements")
	}

	if complianceRegister.SubmittedBYID != nil && complianceRegister.WorkflowStatusID != nil {
		workflowStatus, err := crrs.workflowStatusSvc.GetByID(*complianceRegister.WorkflowStatusID)
		if err != nil {
			return err
		}

		if workflowStatus.Type != "revert" {
			return utility.NewForbiddenError("compliance_register_already_submitted")
		}
	}

	err = crrs.complianceRegisterWorkflowSvc.ValidateWorkflowBeforeCreate("compliance_register")
	if err != nil {
		return err
	}

	tx := crrs.repo.BeginTransaction()
	currentDate := time.Now()
	level := 1

	workflowStatus, err := crrs.workflowStatusSvc.GetByType("submit")
	if err != nil {
		tx.Rollback()
		return err
	}

	err = crrs.complianceRegisterWorkflowSvc.DeleteAllWorkflowsForComplianceRegister(complianceRegisterID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = crrs.complianceRegisterWorkflowSvc.CreateComplianceRegisterWorkflow(complianceRegisterID, "compliance_register", currentUserID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	complianceRegister.WorkflowStatusID = &workflowStatus.ID
	complianceRegister.SubmittedBYID = &currentUserID
	complianceRegister.SubmittedAt = &currentDate
	complianceRegister.NextReviewLevel = &level
	complianceRegister.UpdatedBYID = &currentUserID

	err = crrs.repo.UpdateComplianceRegister(complianceRegister, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	complianceRegisterWorkflowLevel, err := crrs.complianceRegisterWorkflowSvc.GetComplianceRegisterWorkflowByLevel(complianceRegisterID, level, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Create submit review in both current and history tables
	if err := crrs.createReviewWithHistory(complianceRegisterWorkflowLevel.ID, complianceRegisterID, level, workflowStatus.ID, currentUserID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	// Send notifications asynchronously
	go func() {
		urlParam := fmt.Sprintf(`{"id":%d}`, complianceRegisterDetail.ID)
		entity := "compliance-register-review"
		action := "review-user-assign"

		// Get compliance details from DTO for notifications
		referenceNumber := complianceRegisterDetail.ReferenceNumber
		complianceTitle := complianceRegister.TitleEn

		for _, participant := range complianceRegisterWorkflowLevel.ComplianceRegisterWorkflowParticipants {
			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{participant.UserID},
				ModuleIdentifier:     "compliance",
				NotificationType:     string(commonDTOs.ReviewType),
				ReferenceNumber:      referenceNumber,
				HeadingEn:            "Compliance Register Review",
				HeadingAr:            "مراجعة سجل الامتثال",
				TitleEn:              complianceTitle,
				TitleAr:              complianceTitle,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			responsibleUser, err := crrs.userSvc.GetByID(participant.UserID)
			if err != nil {
				utility.LogError("Failed to fetch review user", map[string]interface{}{
					"error": err.Error(),
				})
				continue
			}

			assignedBy, err := crrs.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch assigning user", map[string]interface{}{
					"error": err.Error(),
				})
				continue
			}

			emailTo := []string{responsibleUser.Email}
			subject := "Compliance Register: Review Assigned"
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":           responsibleUser.FirstNameEn,
				"LastNameEn":            responsibleUser.LastNameEn,
				"Title":                 complianceTitle,
				"ReferenceNumber":       referenceNumber,
				"AssignedByFirstNameEn": assignedBy.FirstNameEn,
				"AssignedByLastNameEn":  assignedBy.LastNameEn,
				"AssignedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           fmt.Sprintf("%s/compliance/register/%d/overview", appURL, complianceRegisterDetail.ID),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "compliance/templates/email/compliance_register_workflow_review_user_mail.html",
				EmailData:            emailData,
			}

			err = crrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "compliance_register")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()

	return nil
}

func (crrs *ComplianceRegisterReviewService) Approve(complianceRegisterID uint, currentUserID uint) error {
	tx := crrs.repo.BeginTransaction()

	complianceRegister, err := crrs.complianceRegisterSvc.GetByID(complianceRegisterID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if complianceRegister.NextReviewLevel == nil {
		tx.Rollback()
		return utility.NewForbiddenError("review level not found")
	}

	currentApprovalLevel := complianceRegister.NextReviewLevel
	currentWorkflow, err := crrs.complianceRegisterWorkflowSvc.GetComplianceRegisterWorkflowByLevel(complianceRegisterID, *currentApprovalLevel, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	workflowStatus, err := crrs.workflowStatusSvc.GetByType("approve")
	if err != nil {
		tx.Rollback()
		return err
	}

	switch currentWorkflow.WorkflowLevelType.Type {
	case "individual":
		if currentUserID != currentWorkflow.ComplianceRegisterWorkflowParticipants[0].UserID {
			tx.Rollback()
			return utility.NewForbiddenError("you have no permission to approve")
		}

		// Create approval review in both current and history tables
		if err := crrs.createReviewWithHistory(currentWorkflow.ID, complianceRegisterID, currentWorkflow.Level, workflowStatus.ID, currentUserID, tx); err != nil {
			tx.Rollback()
			return err
		}

		complianceRegisterWorkflow, err := crrs.complianceRegisterWorkflowSvc.GetComplianceRegisterWorkflowByID(currentWorkflow.ID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		complianceRegisterWorkflow.WorkflowStatusID = workflowStatus.ID
		if err := crrs.complianceRegisterWorkflowSvc.UpdateComplianceRegisterWorkflow(complianceRegisterWorkflow, tx); err != nil {
			tx.Rollback()
			return err
		}

		nextWorkflowLevel, err := crrs.complianceRegisterWorkflowSvc.GetNextLevel(complianceRegisterID, *currentApprovalLevel, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		existingComplianceRegister, err := crrs.complianceRegisterSvc.GetByID(complianceRegisterID)
		if err != nil {
			tx.Rollback()
			return err
		}

		if nextWorkflowLevel == nil {
			approvedAt := time.Now()
			existingComplianceRegister.WorkflowStatusID = &workflowStatus.ID
			existingComplianceRegister.ApprovedBYID = &currentUserID
			existingComplianceRegister.ApprovedAt = &approvedAt
			existingComplianceRegister.NextReviewLevel = nil
		} else {
			existingComplianceRegister.NextReviewLevel = &nextWorkflowLevel.Level
		}

		if err := crrs.repo.UpdateComplianceRegister(existingComplianceRegister, tx); err != nil {
			tx.Rollback()
			return err
		}

		// Commit transaction before sending notifications
		if err := tx.Commit(); err != nil {
			return err
		}

		// Send notifications after transaction commit
		if nextWorkflowLevel != nil {
			nextWorkflow, err := crrs.complianceRegisterWorkflowSvc.GetComplianceRegisterWorkflowByLevel(complianceRegisterID, nextWorkflowLevel.Level)
			if err != nil {
				return err
			}

			complianceDetails, err := crrs.complianceRegisterSvc.GetFullDetailsByID(complianceRegisterID)
			if err != nil {
				return err
			}
			crrs.notifyNextLevelUsers(complianceRegisterID, nextWorkflow, complianceDetails, currentUserID)
		}

		complianceDetails, err := crrs.complianceRegisterSvc.GetFullDetailsByID(complianceRegisterID)
		if err != nil {
			return err
		}
		crrs.sendApprovalNotification(complianceRegisterID, currentUserID, complianceDetails)

	case "multiple-users", "role":
		allowed := false
		for _, participant := range currentWorkflow.ComplianceRegisterWorkflowParticipants {
			if participant.UserID == currentUserID {
				allowed = true
				break
			}
		}
		if !allowed {
			tx.Rollback()
			return utility.NewForbiddenError("you have no permission to approve")
		}

		// Create approval review in both current and history tables
		if err := crrs.createReviewWithHistory(currentWorkflow.ID, complianceRegisterID, currentWorkflow.Level, workflowStatus.ID, currentUserID, tx); err != nil {
			tx.Rollback()
			return err
		}

		if currentWorkflow.ApprovalTypes == "all" {
			allApproved, err := crrs.complianceRegisterWorkflowSvc.CheckAllParticipantsApproved(currentWorkflow.ID, currentWorkflow.ComplianceRegisterWorkflowParticipants, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
			if !allApproved {
				return tx.Commit()
			}
		}

		complianceRegisterWorkflow, err := crrs.complianceRegisterWorkflowSvc.GetComplianceRegisterWorkflowByID(currentWorkflow.ID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
		complianceRegisterWorkflow.WorkflowStatusID = workflowStatus.ID
		if err := crrs.complianceRegisterWorkflowSvc.UpdateComplianceRegisterWorkflow(complianceRegisterWorkflow, tx); err != nil {
			tx.Rollback()
			return err
		}

		nextWorkflowLevel, err := crrs.complianceRegisterWorkflowSvc.GetNextLevel(complianceRegisterID, *currentApprovalLevel, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		existingComplianceRegister, err := crrs.complianceRegisterSvc.GetByID(complianceRegisterID)
		if err != nil {
			tx.Rollback()
			return err
		}

		if nextWorkflowLevel == nil {
			approvedAt := time.Now()
			existingComplianceRegister.WorkflowStatusID = &workflowStatus.ID
			existingComplianceRegister.ApprovedBYID = &currentUserID
			existingComplianceRegister.ApprovedAt = &approvedAt
			existingComplianceRegister.NextReviewLevel = nil
		} else {
			existingComplianceRegister.NextReviewLevel = &nextWorkflowLevel.Level
		}

		if err := crrs.repo.UpdateComplianceRegister(existingComplianceRegister, tx); err != nil {
			tx.Rollback()
			return err
		}

		// Commit transaction before sending notifications
		if err := tx.Commit(); err != nil {
			return err
		}

		// Send notifications after transaction commit
		if nextWorkflowLevel != nil {
			nextWorkflow, err := crrs.complianceRegisterWorkflowSvc.GetComplianceRegisterWorkflowByLevel(complianceRegisterID, nextWorkflowLevel.Level)
			if err != nil {
				return err
			}
			complianceDetails, err := crrs.complianceRegisterSvc.GetFullDetailsByID(complianceRegisterID)
			if err != nil {
				return err
			}
			crrs.notifyNextLevelUsers(complianceRegisterID, nextWorkflow, complianceDetails, currentUserID)
		}

		complianceDetails, err := crrs.complianceRegisterSvc.GetFullDetailsByID(complianceRegisterID)
		if err != nil {
			return err
		}
		crrs.sendApprovalNotification(complianceRegisterID, currentUserID, complianceDetails)
	}

	return nil
}

func (crrs *ComplianceRegisterReviewService) Reject(complianceRegisterID uint, currentUserID uint) error {
	// Get full details for notifications
	complianceRegisterDetail, err := crrs.complianceRegisterSvc.GetFullDetailsByID(complianceRegisterID)
	if err != nil {
		return err
	}

	workflowStatus, err := crrs.workflowStatusSvc.GetByType("revert")
	if err != nil {
		return err
	}

	tx := crrs.repo.BeginTransaction()

	complianceRegister, err := crrs.complianceRegisterSvc.GetByID(complianceRegisterID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if complianceRegister.NextReviewLevel == nil {
		tx.Rollback()
		return utility.NewForbiddenError("review level not found")
	}

	currentApprovalLevel := complianceRegister.NextReviewLevel
	currentWorkflow, err := crrs.complianceRegisterWorkflowSvc.GetComplianceRegisterWorkflowByLevel(complianceRegisterID, *currentApprovalLevel, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Create revert review in both current and history tables
	if err := crrs.createReviewWithHistory(currentWorkflow.ID, complianceRegisterID, currentWorkflow.Level, workflowStatus.ID, currentUserID, tx); err != nil {
		tx.Rollback()
		return err
	}

	complianceRegisterWorkflow, err := crrs.complianceRegisterWorkflowSvc.GetComplianceRegisterWorkflowByID(currentWorkflow.ID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	complianceRegisterWorkflow.WorkflowStatusID = workflowStatus.ID
	if err := crrs.complianceRegisterWorkflowSvc.UpdateComplianceRegisterWorkflow(complianceRegisterWorkflow, tx); err != nil {
		tx.Rollback()
		return err
	}

	existingComplianceRegister, err := crrs.complianceRegisterSvc.GetByID(complianceRegisterID)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingComplianceRegister.WorkflowStatusID = &workflowStatus.ID
	existingComplianceRegister.NextReviewLevel = nil
	existingComplianceRegister.SubmittedAt = nil
	existingComplianceRegister.SubmittedBYID = nil
	existingComplianceRegister.UpdatedBYID = &currentUserID
	if err := crrs.repo.UpdateComplianceRegister(existingComplianceRegister, tx); err != nil {
		tx.Rollback()
		return err
	}

	crrs.sendRejectedNotification(complianceRegisterID, currentUserID, complianceRegisterDetail)

	return tx.Commit()
}

func (crrs *ComplianceRegisterReviewService) Reopen(complianceRegisterID uint, currentUserID uint) error {
	// Get full details for notifications
	complianceRegisterDetail, err := crrs.complianceRegisterSvc.GetFullDetailsByID(complianceRegisterID)
	if err != nil {
		return err
	}

	// Get the model for database operations
	complianceRegister, err := crrs.complianceRegisterSvc.GetByID(complianceRegisterID)
	if err != nil {
		return err
	}

	if complianceRegister.SubmittedBYID != nil && complianceRegister.WorkflowStatusID != nil {
		workflowStatus, err := crrs.workflowStatusSvc.GetByID(*complianceRegister.WorkflowStatusID)
		if err != nil {
			return err
		}

		if workflowStatus.Type != "approve" {
			return utility.NewForbiddenError("compliance register already submitted")
		}
	}

	if err := crrs.complianceRegisterWorkflowSvc.ValidateWorkflowBeforeCreate("compliance_register"); err != nil {
		return err
	}

	tx := crrs.repo.BeginTransaction()

	workflowStatus, err := crrs.workflowStatusSvc.GetByType("reopen")
	if err != nil {
		tx.Rollback()
		return err
	}

	// Check if compliance register was submitted (has workflow history)
	if complianceRegister.ApprovedBYID != nil {
		// Get level 1 workflow to attach reopen history entry
		levelOneWorkflow, err := crrs.complianceRegisterWorkflowSvc.GetComplianceRegisterWorkflowByLevel(complianceRegisterID, 1, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		// Create reopen review in both current and history tables at level 1
		if err := crrs.createReviewWithHistory(levelOneWorkflow.ID, complianceRegisterID, 1, workflowStatus.ID, currentUserID, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	complianceRegister.WorkflowStatusID = &workflowStatus.ID
	complianceRegister.SubmittedBYID = nil
	complianceRegister.SubmittedAt = nil
	complianceRegister.NextReviewLevel = nil
	complianceRegister.ApprovedAt = nil
	complianceRegister.ApprovedBYID = nil
	complianceRegister.UpdatedBYID = &currentUserID

	if err := crrs.repo.UpdateComplianceRegister(complianceRegister, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	go func() {

		urlParam := fmt.Sprintf(`{"id":%d}`, complianceRegisterDetail.ID)
		entity := "compliance-register-review"
		action := "reopened"

		if complianceRegister.ResponsibleUserID == nil {
			utility.LogError("Failed to fetch compliance register responsible user", map[string]interface{}{
				"error": "ResponsibleUserID is nil",
			})
			return
		}

		responsibleUser, err := crrs.userSvc.GetByID(*complianceRegister.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch compliance register responsible user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		// Get compliance details from DTO for notifications
		var referenceNumber string
		var complianceTitle string
		if complianceRegisterDetail != nil {
			referenceNumber = complianceRegisterDetail.ReferenceNumber
			complianceTitle = complianceRegisterDetail.TitleEn
		}

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{*complianceRegister.ResponsibleUserID},
			ModuleIdentifier:     "compliance",
			NotificationType:     string(commonDTOs.RevertType),
			HeadingEn:            "Compliance Register Reopened",
			HeadingAr:            "تم إعادة فتح سجل الامتثال",
			ReferenceNumber:      referenceNumber,
			TitleEn:              complianceTitle,
			TitleAr:              complianceTitle,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		reopenedByUser, err := crrs.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch compliance register reopen user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		appURL := utility.GetAppURL()
		emailTo := []string{responsibleUser.Email}
		subject := "Compliance Register: Review Reopened"
		redirectURL := fmt.Sprintf("%s/compliance/register/%d/overview", appURL, complianceRegisterDetail.ID)

		emailData := map[string]interface{}{
			"FirstNameEn":           responsibleUser.FirstNameEn,
			"LastNameEn":            responsibleUser.LastNameEn,
			"Title":                 complianceTitle,
			"ReferenceNumber":       referenceNumber,
			"ReopenedByFirstNameEn": reopenedByUser.FirstNameEn,
			"ReopenedByLastNameEn":  reopenedByUser.LastNameEn,
			"ReopenedAt":            time.Now().Format("2006-01-02"),
			"RedirectURL":           redirectURL,
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "compliance/templates/email/compliance_register_workflow_reopen_user_mail.html",
			EmailData:            emailData,
		}

		err = crrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "compliance_register")
		if err != nil {
			utility.LogError("Failed to send compliance register reopened notification email", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	return nil
}

func (crrs *ComplianceRegisterReviewService) sendApprovalNotification(id uint, currentUserID uint, complianceDetail *complianceDTOs.ComplianceRegisterDetailDTO) {
	go func() {
		if complianceDetail.ResponsibleUser.ID == 0 {
			return
		}

		urlParam := fmt.Sprintf(`{"id":%d}`, complianceDetail.ID)
		entity := "compliance-register-review"
		action := "approval"

		// Get compliance details for notifications
		var referenceNumber string
		var complianceTitle string
		if complianceDetail != nil {
			referenceNumber = complianceDetail.ReferenceNumber
			complianceTitle = complianceDetail.TitleEn
		}

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{uint(complianceDetail.ResponsibleUser.ID)},
			ModuleIdentifier:     "compliance",
			NotificationType:     string(commonDTOs.ApprovalType),
			HeadingEn:            "Compliance Register Approval",
			HeadingAr:            "Compliance Register Approval",
			ReferenceNumber:      referenceNumber,
			TitleEn:              complianceTitle,
			TitleAr:              complianceTitle,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		approvedByUser, err := crrs.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch approving user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"Title":                 complianceTitle,
			"Description":           "",
			"FirstNameEn":           complianceDetail.ResponsibleUser.FirstNameEN,
			"LastNameEn":            complianceDetail.ResponsibleUser.LastNameEN,
			"ReferenceNumber":       referenceNumber,
			"ApprovedByFirstNameEn": approvedByUser.FirstNameEn,
			"ApprovedByLastNameEn":  approvedByUser.LastNameEn,
			"ApprovedAt":            time.Now().Format("2006-01-02"),
			"RedirectURL":           fmt.Sprintf("%s/compliance/register/%d/overview", appURL, complianceDetail.ID),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              []string{complianceDetail.ResponsibleUser.Email},
			Subject:              "Compliance Register Approval Notification",
			TemplateRelativePath: "compliance/templates/email/compliance_register_workflow_approve_user_mail.html",
			EmailData:            emailData,
		}

		err = crrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "compliance_register")
		if err != nil {
			utility.LogError("Failed to send approval notification email", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()
}

func (crrs *ComplianceRegisterReviewService) notifyNextLevelUsers(id uint, nextWorkflowLevel *models.ComplianceRegisterWorkflow, complianceDetail *complianceDTOs.ComplianceRegisterDetailDTO, currentUserID uint) {
	go func() {
		if complianceDetail.CreatedByUser.ID == 0 {
			return
		}

		urlParam := fmt.Sprintf(`{"id":%d}`, complianceDetail.ID)
		entity := "compliance-register-review"
		action := "review-user-assign"

		// Get compliance details for notifications
		var referenceNumber string
		var complianceTitle string
		if complianceDetail != nil {
			referenceNumber = complianceDetail.ReferenceNumber
			complianceTitle = complianceDetail.TitleEn
		}

		for _, participant := range nextWorkflowLevel.ComplianceRegisterWorkflowParticipants {
			responsibleUser, err := crrs.userSvc.GetByID(participant.UserID)
			if err != nil {
				utility.LogError("Failed to fetch review user", map[string]interface{}{
					"error": err.Error(),
				})
				continue
			}

			createdByID := uint(complianceDetail.CreatedByUser.ID)
			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{participant.UserID},
				ModuleIdentifier:     "compliance",
				NotificationType:     string(commonDTOs.ReviewType),
				HeadingEn:            "Compliance Register Review",
				HeadingAr:            "مراجعة سجل الامتثال",
				ReferenceNumber:      referenceNumber,
				TitleEn:              complianceTitle,
				TitleAr:              complianceTitle,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &createdByID,
			}

			assignedBy, err := crrs.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch assigning user", map[string]interface{}{
					"error": err.Error(),
				})
				continue
			}

			emailTo := []string{responsibleUser.Email}
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":           responsibleUser.FirstNameEn,
				"LastNameEn":            responsibleUser.LastNameEn,
				"Title":                 complianceTitle,
				"Description":           "",
				"ReferenceNumber":       referenceNumber,
				"ResponsibleUserName":   complianceDetail.CreatedByUser.FirstNameEN + " " + complianceDetail.CreatedByUser.LastNameEN,
				"ResponsibleUserEmail":  complianceDetail.CreatedByUser.Email,
				"AssignedByFirstNameEn": assignedBy.FirstNameEn,
				"AssignedByLastNameEn":  assignedBy.LastNameEn,
				"AssignedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           fmt.Sprintf("%s/compliance/register/%d/overview", appURL, complianceDetail.ID),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              "Compliance Register Review Assigned",
				TemplateRelativePath: "compliance/templates/email/compliance_register_workflow_review_user_mail.html",
				EmailData:            emailData,
			}

			err = crrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "compliance_register")
			if err != nil {
				utility.LogError("Failed to send next level review notification", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()
}

func (crrs *ComplianceRegisterReviewService) sendRejectedNotification(id uint, currentUserID uint, complianceRegisterDetail *complianceDTOs.ComplianceRegisterDetailDTO) {
	go func() {
		if complianceRegisterDetail.ResponsibleUser.ID == 0 {
			return
		}

		urlParam := fmt.Sprintf(`{"id":%d}`, complianceRegisterDetail.ID)
		entity := "compliance-register-review"
		action := "reverted"

		// Get compliance details from DTO for notifications
		var referenceNumber string
		var complianceTitle string
		if complianceRegisterDetail != nil {
			referenceNumber = complianceRegisterDetail.ReferenceNumber
			complianceTitle = complianceRegisterDetail.TitleEn
		}

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{complianceRegisterDetail.ResponsibleUserID},
			ModuleIdentifier:     "compliance",
			NotificationType:     string(commonDTOs.RejectType),
			HeadingEn:            "Compliance Register Reverted",
			HeadingAr:            "تم رفض سجل الامتثال",
			ReferenceNumber:      referenceNumber,
			TitleEn:              complianceTitle,
			TitleAr:              complianceTitle,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		rejectedByUser, err := crrs.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch review performed user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstNameEn":           complianceRegisterDetail.ResponsibleUser.FirstNameEN,
			"LastNameEn":            complianceRegisterDetail.ResponsibleUser.LastNameEN,
			"Title":                 complianceTitle,
			"ReferenceNumber":       referenceNumber,
			"RejectedByFirstNameEn": rejectedByUser.FirstNameEn,
			"RejectedByLastNameEn":  rejectedByUser.LastNameEn,
			"RejectedAt":            time.Now().Format("2006-01-02"),
			"RedirectURL":           fmt.Sprintf("%s/compliance/register/%d/overview", appURL, complianceRegisterDetail.ID),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              []string{complianceRegisterDetail.ResponsibleUser.Email},
			Subject:              "Compliance Register: Reverted Notification",
			TemplateRelativePath: "compliance/templates/email/compliance_register_workflow_reject_user_mail.html",
			EmailData:            emailData,
		}

		err = crrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "compliance_register")
		if err != nil {
			utility.LogError("Failed to send rejected notification email", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()
}

// Helper function to create review in both current and history tables
func (crrs *ComplianceRegisterReviewService) createReviewWithHistory(workflowID uint, complianceRegisterID uint, level int, statusID uint, userID uint, tx utility.DatabaseTransaction) error {
	// Create in current reviews table
	complianceRegisterReview := &models.ComplianceRegisterWorkflowReview{
		ComplianceRegisterWorkflowID: workflowID,
		WorkflowStatusID:             statusID,
		CreatedBYID:                  userID,
	}
	if err := crrs.repo.CreateComplianceRegisterReview(complianceRegisterReview, tx); err != nil {
		return err
	}

	// Create in history table for audit trail (no FK to workflow to avoid deletion conflicts)
	historyReview := &models.ComplianceRegisterWorkflowHistory{
		ComplianceRegisterID: complianceRegisterID,
		Level:                level,
		WorkflowStatusID:     statusID,
		CreatedBYID:          userID,
		CreatedAt:            time.Now(),
	}
	if err := crrs.repo.CreateComplianceRegisterWorkflowHistory(historyReview, tx); err != nil {
		return err
	}

	return nil
}
