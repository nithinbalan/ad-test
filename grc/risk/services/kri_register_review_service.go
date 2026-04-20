package services

import (
	"fmt"
	"time"

	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationSvc "jethur/grc/organization/services"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/grc/risk/repositories"
	taskServices "jethur/grc/task/services"
	"jethur/utility"
)

type KRIRegisterReviewServiceInterface interface {
	Submit(kriRegisterID uint, currentUserID uint) error
	Approve(kriRegisterID uint, currentUserID uint) error
	Reject(kriRegisterID uint, currentUserID uint) error
}

type KRIRegisterReviewService struct {
	repo                   repositories.KRIRegisterReviewRepositoryInterface
	logService             logServices.ActivityLogServiceInterface
	kriRegisterSvc         KRIRegisterServiceInterface
	kriRegisterWorkflowSvc KRIRegisterWorkflowServiceInterface
	workflowStatusSvc      masterServices.WorkflowStatusServiceInterface
	userSvc                organizationSvc.UserServiceInterface
	notificationSvc        commonServices.NotificationServiceInterface
	taskSvc                taskServices.TaskServiceInterface
}

func NewKRIRegisterReviewService(
	repo repositories.KRIRegisterReviewRepositoryInterface,
	logService logServices.ActivityLogServiceInterface,
	kriRegisterSvc KRIRegisterServiceInterface,
	kriRegisterWorkflowSvc KRIRegisterWorkflowServiceInterface,
	workflowStatusSvc masterServices.WorkflowStatusServiceInterface,
	userSvc organizationSvc.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	taskSvc taskServices.TaskServiceInterface,
) KRIRegisterReviewServiceInterface {
	return &KRIRegisterReviewService{
		repo:                   repo,
		logService:             logService,
		kriRegisterSvc:         kriRegisterSvc,
		kriRegisterWorkflowSvc: kriRegisterWorkflowSvc,
		workflowStatusSvc:      workflowStatusSvc,
		userSvc:                userSvc,
		notificationSvc:        notificationSvc,
		taskSvc:                taskSvc,
	}
}

func (krs *KRIRegisterReviewService) Submit(kriRegisterID uint, currentUserID uint) error {
	kriRegister, err := krs.kriRegisterSvc.GetByID(kriRegisterID)
	if err != nil {
		return err
	}

	tx := krs.repo.BeginTransaction()
	currentDate := time.Now()
	level := 1

	workflowStatus, err := krs.workflowStatusSvc.GetByType("submit")
	if err != nil {
		tx.Rollback()
		return err
	}

	kriRegister.WorkflowStatusID = &workflowStatus.ID
	kriRegister.SubmittedByID = &currentUserID
	kriRegister.SubmittedAt = &currentDate
	kriRegister.NextReviewLevel = &level
	kriRegister.UpdatedBYID = &currentUserID

	err = krs.repo.UpdateKRIRegister(kriRegister, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	kriWorkflowLevel, err := krs.kriRegisterWorkflowSvc.GetKRIRegisterWorkflowByLevel(kriRegister.ID, level, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	kriReview := &models.KRIRegisterWorkflowReview{
		KRIRegisterWorkflowID: kriWorkflowLevel.ID,
		WorkflowStatusID:      workflowStatus.ID,
		CreatedBYID:           currentUserID,
	}

	err = krs.repo.CreateKRIRegisterReview(kriReview, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	go func() {
		urlParam := fmt.Sprintf(`{"id":%d}`, kriRegisterID)
		entity := "kri-review"
		action := "review-user-assign"
		referenceNumber := ""
		if kriRegister.ReferenceNumber != nil {
			referenceNumber = *kriRegister.ReferenceNumber
		}

		for _, participant := range kriWorkflowLevel.KRIRegisterWorkflowParticipants {
			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{participant.UserID},
				ModuleIdentifier:     "risk",
				NotificationType:     string(commonDTOs.ReviewType),
				ReferenceNumber:      referenceNumber,
				HeadingEn:            "KRI Review",
				HeadingAr:            "KRI Review",
				TitleEn:              kriRegister.Title,
				TitleAr:              kriRegister.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			responsibleUser, err := krs.userSvc.GetByID(participant.UserID)
			if err != nil {
				utility.LogError("Failed to fetch review user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			assignedBy, err := krs.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch assigning user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":           responsibleUser.FirstNameEn,
				"LastNameEn":            responsibleUser.LastNameEn,
				"Title":                 kriRegister.Title,
				"ReferenceNumber":       referenceNumber,
				"AssignedByFirstNameEn": assignedBy.FirstNameEn,
				"AssignedByLastNameEn":  assignedBy.LastNameEn,
				"AssignedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           fmt.Sprintf("%srisk/kri-register/%d/overview", appURL, kriRegisterID),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              []string{responsibleUser.Email},
				Subject:              fmt.Sprintf("KRI %s: Review Assigned", kriRegister.Title),
				TemplateRelativePath: "risk/templates/email/kri_register_workflow_review_user_mail.html",
				EmailData:            emailData,
			}

			err = krs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "kri_register")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()

	return nil
}

func (krs *KRIRegisterReviewService) Approve(kriRegisterID uint, currentUserID uint) error {
	tx := krs.repo.BeginTransaction()

	kriRegister, err := krs.kriRegisterSvc.GetByID(kriRegisterID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if kriRegister.NextReviewLevel == nil {
		tx.Rollback()
		return utility.NewForbiddenError("review level not found")
	}

	currentApprovalLevel := kriRegister.NextReviewLevel
	currentWorkflow, err := krs.kriRegisterWorkflowSvc.GetKRIRegisterWorkflowByLevel(kriRegisterID, *currentApprovalLevel, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	workflowStatus, err := krs.workflowStatusSvc.GetByType("approve")
	if err != nil {
		tx.Rollback()
		return err
	}

	switch currentWorkflow.WorkflowLevelType.Type {
	case "individual":
		if currentUserID != currentWorkflow.KRIRegisterWorkflowParticipants[0].UserID {
			tx.Rollback()
			return utility.NewForbiddenError("you have no permission to approve")
		}

		kriReview := &models.KRIRegisterWorkflowReview{
			KRIRegisterWorkflowID: currentWorkflow.ID,
			WorkflowStatusID:      workflowStatus.ID,
			CreatedBYID:           currentUserID,
		}
		if err := krs.repo.CreateKRIRegisterReview(kriReview, tx); err != nil {
			tx.Rollback()
			return err
		}

		kriWorkflow, err := krs.kriRegisterWorkflowSvc.GetKRIRegisterWorkflowByID(currentWorkflow.ID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		kriWorkflow.WorkflowStatusID = workflowStatus.ID
		if err := krs.kriRegisterWorkflowSvc.UpdateKRIRegisterWorkflow(kriWorkflow, tx); err != nil {
			tx.Rollback()
			return err
		}

		nextWorkflowLevel, err := krs.kriRegisterWorkflowSvc.GetNextLevel(kriRegisterID, *currentApprovalLevel, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		existingKRI, err := krs.kriRegisterSvc.GetByID(kriRegisterID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		if nextWorkflowLevel == nil {
			approvedAt := time.Now()
			existingKRI.WorkflowStatusID = &workflowStatus.ID
			existingKRI.ApprovedByID = &currentUserID
			existingKRI.ApprovedAt = &approvedAt
			existingKRI.NextReviewLevel = nil
		} else {
			existingKRI.NextReviewLevel = &nextWorkflowLevel.Level
		}

		if err := krs.repo.UpdateKRIRegister(existingKRI, tx); err != nil {
			tx.Rollback()
			return err
		}

		if nextWorkflowLevel != nil {
			nextWorkflow, err := krs.kriRegisterWorkflowSvc.GetKRIRegisterWorkflowByLevel(kriRegisterID, nextWorkflowLevel.Level, tx)
			if err != nil {
				tx.Rollback()
				return err
			}

			kriDetails, err := krs.kriRegisterSvc.GetFullDetailsByID(kriRegisterID, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
			krs.notifyNextLevelUsers(kriRegisterID, nextWorkflow, kriDetails)
		}

		kriDetails, err := krs.kriRegisterSvc.GetFullDetailsByID(kriRegisterID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
		krs.sendApprovalNotification(kriRegisterID, currentUserID, kriDetails)

	case "multiple-users", "role":
		allowed := false
		for _, participant := range currentWorkflow.KRIRegisterWorkflowParticipants {
			if participant.UserID == currentUserID {
				allowed = true
				break
			}
		}
		if !allowed {
			tx.Rollback()
			return utility.NewForbiddenError("you have no permission to approve")
		}

		kriReview := &models.KRIRegisterWorkflowReview{
			KRIRegisterWorkflowID: currentWorkflow.ID,
			WorkflowStatusID:      workflowStatus.ID,
			CreatedBYID:           currentUserID,
		}
		if err := krs.repo.CreateKRIRegisterReview(kriReview, tx); err != nil {
			tx.Rollback()
			return err
		}

		if currentWorkflow.ApprovalTypes == "all" {
			allApproved, err := krs.kriRegisterWorkflowSvc.CheckAllParticipantsApproved(currentWorkflow.ID, currentWorkflow.KRIRegisterWorkflowParticipants, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
			if !allApproved {
				return tx.Commit()
			}
		}

		kriWorkflow, err := krs.kriRegisterWorkflowSvc.GetKRIRegisterWorkflowByID(currentWorkflow.ID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
		kriWorkflow.WorkflowStatusID = workflowStatus.ID
		if err := krs.kriRegisterWorkflowSvc.UpdateKRIRegisterWorkflow(kriWorkflow, tx); err != nil {
			tx.Rollback()
			return err
		}

		nextWorkflowLevel, err := krs.kriRegisterWorkflowSvc.GetNextLevel(kriRegisterID, *currentApprovalLevel, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		existingKRI, err := krs.kriRegisterSvc.GetByID(kriRegisterID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		if nextWorkflowLevel == nil {
			approvedAt := time.Now()
			existingKRI.WorkflowStatusID = &workflowStatus.ID
			existingKRI.ApprovedByID = &currentUserID
			existingKRI.ApprovedAt = &approvedAt
			existingKRI.NextReviewLevel = nil
		} else {
			existingKRI.NextReviewLevel = &nextWorkflowLevel.Level
		}

		if err := krs.repo.UpdateKRIRegister(existingKRI, tx); err != nil {
			tx.Rollback()
			return err
		}

		if nextWorkflowLevel != nil {
			nextWorkflow, err := krs.kriRegisterWorkflowSvc.GetKRIRegisterWorkflowByLevel(kriRegisterID, nextWorkflowLevel.Level, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
			kriDetails, err := krs.kriRegisterSvc.GetFullDetailsByID(kriRegisterID, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
			krs.notifyNextLevelUsers(kriRegisterID, nextWorkflow, kriDetails)
		}

		kriDetails, err := krs.kriRegisterSvc.GetFullDetailsByID(kriRegisterID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
		krs.sendApprovalNotification(kriRegisterID, currentUserID, kriDetails)
	}

	return tx.Commit()
}

func (krs *KRIRegisterReviewService) Reject(kriRegisterID uint, currentUserID uint) error {
	workflowStatus, err := krs.workflowStatusSvc.GetByType("reject")
	if err != nil {
		return err
	}

	tx := krs.repo.BeginTransaction()

	kriRegister, err := krs.kriRegisterSvc.GetByID(kriRegisterID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if kriRegister.NextReviewLevel == nil {
		tx.Rollback()
		return utility.NewForbiddenError("review level not found")
	}

	currentApprovalLevel := kriRegister.NextReviewLevel
	currentWorkflow, err := krs.kriRegisterWorkflowSvc.GetKRIRegisterWorkflowByLevel(kriRegisterID, *currentApprovalLevel, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	kriWorkflowReview := &models.KRIRegisterWorkflowReview{
		KRIRegisterWorkflowID: currentWorkflow.ID,
		WorkflowStatusID:      workflowStatus.ID,
		CreatedBYID:           currentUserID,
	}

	if err := krs.repo.CreateKRIRegisterReview(kriWorkflowReview, tx); err != nil {
		tx.Rollback()
		return err
	}

	kriWorkflow, err := krs.kriRegisterWorkflowSvc.GetKRIRegisterWorkflowByID(currentWorkflow.ID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	kriWorkflow.WorkflowStatusID = workflowStatus.ID
	if err := krs.kriRegisterWorkflowSvc.UpdateKRIRegisterWorkflow(kriWorkflow, tx); err != nil {
		tx.Rollback()
		return err
	}

	existingKRI, err := krs.kriRegisterSvc.GetByID(kriRegisterID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingKRI.WorkflowStatusID = &workflowStatus.ID
	existingKRI.UpdatedBYID = &currentUserID
	if err := krs.repo.UpdateKRIRegister(existingKRI, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	kriDetails, err := krs.kriRegisterSvc.GetFullDetailsByID(kriRegisterID)
	if err != nil {
		return err
	}
	krs.sendRejectedNotification(kriRegisterID, currentUserID, kriDetails)

	return nil
}

func (krs *KRIRegisterReviewService) sendApprovalNotification(id uint, currentUserID uint, kriDetail *riskDTOs.KRIRegisterDetailsDTO) {
	go func() {
		if kriDetail.ResponsibleUser == nil {
			return
		}

		urlParam := fmt.Sprintf(`{"id":%d}`, kriDetail.ID)
		entity := "kri-review"
		action := "approval"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{uint(kriDetail.ResponsibleUser.ID)},
			ModuleIdentifier:     "risk",
			NotificationType:     string(commonDTOs.ApprovalType),
			HeadingEn:            "KRI Approval",
			HeadingAr:            "KRI Approval",
			ReferenceNumber:      kriDetail.ReferenceNumber,
			TitleEn:              kriDetail.Title,
			TitleAr:              kriDetail.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		approvedByUser, err := krs.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch approving user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"Title":                 kriDetail.Title,
			"Description":           utility.GetStringValue(kriDetail.Description),
			"FirstNameEn":           kriDetail.ResponsibleUser.FirstNameEN,
			"LastNameEn":            kriDetail.ResponsibleUser.LastNameEN,
			"ApprovedByFirstNameEn": approvedByUser.FirstNameEn,
			"ApprovedByLastNameEn":  approvedByUser.LastNameEn,
			"ApprovedAt":            time.Now().Format("2006-01-02"),
			"RedirectURL":           fmt.Sprintf("%srisk/kri-register/%d/overview", appURL, id),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              []string{kriDetail.ResponsibleUser.Email},
			Subject:              fmt.Sprintf("KRI %s: Approval Notification", kriDetail.Title),
			TemplateRelativePath: "risk/templates/email/kri_register_workflow_approve_user_mail.html",
			EmailData:            emailData,
		}

		err = krs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "kri_register")
		if err != nil {
			utility.LogError("Failed to send approval notification email", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()
}

func (krs *KRIRegisterReviewService) notifyNextLevelUsers(id uint, nextWorkflowLevel *models.KRIRegisterWorkflow, kriDetail *riskDTOs.KRIRegisterDetailsDTO) {
	go func() {
		if kriDetail.ResponsibleUser == nil {
			return
		}

		urlParam := fmt.Sprintf(`{"id":%d}`, kriDetail.ID)
		entity := "kri-review"
		action := "review-user-assign"

		for _, participant := range nextWorkflowLevel.KRIRegisterWorkflowParticipants {
			responsibleUser, err := krs.userSvc.GetByID(participant.UserID)
			if err != nil {
				utility.LogError("Failed to fetch review user", map[string]interface{}{
					"error": err.Error(),
				})
				continue
			}

			createdByID := uint(kriDetail.ResponsibleUserID)
			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{participant.UserID},
				ModuleIdentifier:     "risk",
				NotificationType:     string(commonDTOs.ReviewType),
				HeadingEn:            "KRI Review",
				HeadingAr:            "KRI Review",
				ReferenceNumber:      kriDetail.ReferenceNumber,
				TitleEn:              kriDetail.Title,
				TitleAr:              kriDetail.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &createdByID,
			}

			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":           responsibleUser.FirstNameEn,
				"LastNameEn":            responsibleUser.LastNameEn,
				"Title":                 kriDetail.Title,
				"Description":           utility.GetStringValue(kriDetail.Description),
				"AssignedByFirstNameEn": kriDetail.ResponsibleUser.FirstNameEN,
				"AssignedByLastNameEn":  kriDetail.ResponsibleUser.LastNameEN,
				"AssignedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           fmt.Sprintf("%srisk/kri-register/%d/overview", appURL, id),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              []string{responsibleUser.Email},
				Subject:              fmt.Sprintf("KRI %s: Review Assigned", kriDetail.Title),
				TemplateRelativePath: "risk/templates/email/kri_register_workflow_review_user_mail.html",
				EmailData:            emailData,
			}

			err = krs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "kri_register")
			if err != nil {
				utility.LogError("Failed to send next level review notification", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()
}

func (krs *KRIRegisterReviewService) sendRejectedNotification(id uint, currentUserID uint, kri *riskDTOs.KRIRegisterDetailsDTO) {
	go func() {
		if kri.ResponsibleUser == nil {
			return
		}

		urlParam := fmt.Sprintf(`{"id":%d}`, id)
		entity := "kri-review"
		action := "rejected"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{uint(kri.ResponsibleUser.ID)},
			ModuleIdentifier:     "risk",
			NotificationType:     string(commonDTOs.RejectType),
			HeadingEn:            "KRI Rejected",
			HeadingAr:            "KRI Rejected",
			ReferenceNumber:      kri.ReferenceNumber,
			TitleEn:              kri.Title,
			TitleAr:              kri.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		rejectedByUser, err := krs.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch review performed user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstNameEn":           kri.ResponsibleUser.FirstNameEN,
			"LastNameEn":            kri.ResponsibleUser.LastNameEN,
			"Title":                 kri.Title,
			"RejectedByFirstNameEn": rejectedByUser.FirstNameEn,
			"RejectedByLastNameEn":  rejectedByUser.LastNameEn,
			"ApprovedAt":            time.Now().Format("2006-01-02"),
			"RedirectURL":           fmt.Sprintf("%srisk/kri-register/%d/overview", appURL, id),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              []string{kri.ResponsibleUser.Email},
			Subject:              fmt.Sprintf("KRI %s: Rejected Notification", kri.Title),
			TemplateRelativePath: "risk/templates/email/kri_register_workflow_reject_user_mail.html",
			EmailData:            emailData,
		}

		err = krs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "kri_register")
		if err != nil {
			utility.LogError("Failed to send rejected notification email", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()
}
