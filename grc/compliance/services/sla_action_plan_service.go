package services

import (
	"errors"
	"fmt"
	"time"

	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	complianceDTOs "jethur/grc/compliance/dtos"
	complianceRepositories "jethur/grc/compliance/repositories"
	logServices "jethur/grc/logs/services"
	masterRepositories "jethur/grc/masters/repositories"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	taskServices "jethur/grc/task/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type SLAActionPlanServiceInterface interface {
	GetAll(slaID uint, params *complianceDTOs.SLAActionPlanListQueryParamsDTO) (*complianceDTOs.SLAActionPlanListResponseDTO, error)
	GetByID(slaID uint, id uint) (*models.SLAActionPlan, error)
	GetFullDetailsByID(slaID uint, id uint) (*complianceDTOs.SLAActionPlanDetailDTO, error)
	Create(slaID uint, currentUserID uint, slaActionPlanFormDTO *complianceDTOs.SLAActionPlanFormDTO) (uint, error)
	Update(slaID uint, id uint, currentUserID uint, slaActionPlanFormDTO *complianceDTOs.SLAActionPlanFormDTO) error
	Delete(slaID uint, ids []uint, currentUserID uint) ([]uint, error)
	StatusUpdate(slaID uint, id uint, currentUserID uint, slaActionPlanStatusUpdateFormDTO *complianceDTOs.SLAActionPlanStatusUpdateFormDTO, tx utility.DatabaseTransaction) error
	GetSLAActionPlanCountByStatus(slaID uint) ([]complianceDTOs.SLAActionPlanStatusCountResponseDTO, error)
}

// Service and implementations
type SLAActionPlanService struct {
	repo                complianceRepositories.SLAActionPlanRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	slaRegisterSvc      SLARegisterServiceInterface
	slaCheckQuestionSvc SLACheckQuestionServiceInterface
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface
	severityLevelSvc    masterServices.SeverityLevelServiceInterface
	taskSvc             taskServices.TaskServiceInterface
	userSvc             organizationServices.UserServiceInterface
	notificationSvc     commonServices.NotificationServiceInterface
}

func NewSLAActionPlanService(
	repo complianceRepositories.SLAActionPlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	slaRegisterSvc SLARegisterServiceInterface,
	slaCheckQuestionSvc SLACheckQuestionServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	severityLevelSvc masterRepositories.SeverityLevelRepositoryInterface,
	taskSvc taskServices.TaskServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) SLAActionPlanServiceInterface {
	return &SLAActionPlanService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		slaRegisterSvc:      slaRegisterSvc,
		slaCheckQuestionSvc: slaCheckQuestionSvc,
		actionPlanStatusSvc: actionPlanStatusSvc,
		severityLevelSvc:    severityLevelSvc,
		taskSvc:             taskSvc,
		userSvc:             userSvc,
		notificationSvc:     notificationSvc,
	}
}

func (saps *SLAActionPlanService) GetAll(slaID uint, params *complianceDTOs.SLAActionPlanListQueryParamsDTO) (*complianceDTOs.SLAActionPlanListResponseDTO, error) {
	if _, err := saps.slaRegisterSvc.GetByID(slaID); err != nil {
		return nil, err
	}
	return saps.repo.GetAll(slaID, params)
}

func (saps *SLAActionPlanService) GetByID(slaID uint, id uint) (*models.SLAActionPlan, error) {
	if _, err := saps.slaRegisterSvc.GetByID(slaID); err != nil {
		return nil, err
	}

	slaActionPlan, err := saps.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return slaActionPlan, nil
}

func (saps *SLAActionPlanService) GetFullDetailsByID(slaID uint, id uint) (*complianceDTOs.SLAActionPlanDetailDTO, error) {
	if _, err := saps.slaRegisterSvc.GetByID(slaID); err != nil {
		return nil, err
	}

	var slaActionPlanDTO complianceDTOs.SLAActionPlanDetailDTO

	slaActionPlan, err := saps.repo.GetFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(slaActionPlan, &slaActionPlanDTO)
	if err != nil {
		return nil, err
	}

	if slaActionPlan.TargetDate != nil {
		targetDate := slaActionPlan.TargetDate.Format("2006-01-02")
		slaActionPlanDTO.TargetDateString = &targetDate
	}

	return &slaActionPlanDTO, nil
}

func (saps *SLAActionPlanService) Create(slaID uint, currentUserID uint, slaActionPlanFormDTO *complianceDTOs.SLAActionPlanFormDTO) (uint, error) {
	tx := saps.repo.BeginTransaction()
	sla, err := saps.slaRegisterSvc.GetByID(slaID)
	if err != nil {
		return 0, err
	}

	if err := saps.validateSLAActionPlanFormDTO(slaID, slaActionPlanFormDTO); err != nil {
		return 0, err
	}

	dates, err := saps.parseSLAActionPlanDates(slaActionPlanFormDTO)
	if err != nil {
		return 0, err
	}

	actionPlanStatus, err := saps.actionPlanStatusSvc.GetByType("new")
	if err != nil {
		return 0, err
	}

	slaActionPlan := &models.SLAActionPlan{
		SLARegisterID:      slaID,
		SLACheckQuestionID: slaActionPlanFormDTO.SLACheckQuestionID,
		Title:              slaActionPlanFormDTO.Title,
		Description:        slaActionPlanFormDTO.Description,
		SeverityLevelID:    slaActionPlanFormDTO.SeverityLevelID,
		ResponsibleUserID:  slaActionPlanFormDTO.ResponsibleUserID,
		ReviewUserID:       slaActionPlanFormDTO.ReviewUserID,
		TargetDate:         dates.TargetDate,
		Percentage:         slaActionPlanFormDTO.Percentage,
		Comments:           slaActionPlanFormDTO.Comments,
		ActionPlanStatusID: actionPlanStatus.ID,
		CreatedByID:        currentUserID,
	}

	id, err := saps.repo.Create(slaActionPlan)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	actionPlanStatusID := actionPlanStatus.ID
	slaActionPlanUpdate := &models.SLAActionPlanStatusUpdate{
		SLAActionPlanID:    id,
		ActionPlanStatusID: actionPlanStatusID,
		Comment:            "Initial action plan status",
		Percentage:         slaActionPlanFormDTO.Percentage,
		CreatedByID:        currentUserID,
	}

	_, err = saps.repo.CreateInitialSLAActionPlanStatusUpdate(slaActionPlanUpdate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	urlParam := fmt.Sprintf(`{"sla_register_id":%d, "action_plan_id":%d}`, slaID, id)
	taskDTO := commonDTOs.TaskFormDTO{
		ModuleIdentifier:  "compliance",
		SLAActionPlanID:   &id,
		Title:             slaActionPlanFormDTO.Title,
		Description:       slaActionPlanFormDTO.Description,
		TargetDate:        dates.TargetDate,
		URLParam:          &urlParam,
		QueryParam:        nil,
		SeverityLevelID:   slaActionPlanFormDTO.SeverityLevelID,
		ResponsibleUserID: slaActionPlanFormDTO.ResponsibleUserID,
	}

	if _, err := saps.taskSvc.Create(taskDTO, currentUserID, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	go func() {
		existingSlaActionPlan, err := saps.GetFullDetailsByID(slaID, id)
		if err != nil {
			utility.LogError("Failed to fetch compliance action plan", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = saps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.SLAActionPlan{}.TableName(),
			id,
			existingSlaActionPlan,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.SLAActionPlan{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}

		// prepare notification
		urlParam := fmt.Sprintf(`{"sla_register_id":%d, "action_plan_id":%d}`, slaID, id)
		queryParam := ""
		entity := "sla-action-plan"
		action := "responsible-user-assign"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserID:               slaActionPlanFormDTO.ResponsibleUserID,
			ModuleIdentifier:     "compliance",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "Sla Action Plan Responsiblity Assigned",
			HeadingAr:            "تم تعيين مسؤولية خطة عمل Sla",
			TitleEn:              slaActionPlanFormDTO.Title,
			TitleAr:              slaActionPlanFormDTO.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			QueryParam:           queryParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		responsibleUser, err := saps.userSvc.GetByID(slaActionPlanFormDTO.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch responsible user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		assignedBy, err := saps.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch assigning user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("Sla Action Plan/%s: Responsibility Assigned", slaActionPlanFormDTO.Title)
		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstNameEn":           responsibleUser.FirstNameEn,
			"LastNameEn":            responsibleUser.LastNameEn,
			"SlaReferenceNumber":    sla.ReferenceNumber,
			"SlaTitle":              sla.TitleEn,
			"ActionPlanTitle":       slaActionPlanFormDTO.Title,
			"ActionPlanDescription": slaActionPlanFormDTO.Description,
			"TargetDate":            dates.TargetDate.Format("2006-01-02"),
			"AssignedByFirstNameEn": assignedBy.FirstNameEn,
			"AssignedByLastNameEn":  assignedBy.LastNameEn,
			"AssignedOn":            time.Now().Format("2006-01-02"),
			"RedirectURL":           fmt.Sprintf("%scompliance/sla-contracts/%d/action-plans?action_plan_id=%d", appURL, slaID, id),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "compliance/templates/email/sla_action_plan_create_responsible_user_mail.html",
			EmailData:            emailData,
		}

		// send notification and mail
		err = saps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "sla_register")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	return id, nil
}

func (saps *SLAActionPlanService) Update(slaID uint, id uint, currentUserID uint, slaActionPlanFormDTO *complianceDTOs.SLAActionPlanFormDTO) error {
	tx := saps.repo.BeginTransaction()

	sla, err := saps.slaRegisterSvc.GetByID(slaID)
	if err != nil {
		return err
	}

	existingSLAActionPlan, err := saps.GetByID(slaID, id)
	if err != nil {
		return err
	}

	oldSLAActionPlan, err := saps.GetFullDetailsByID(slaID, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	oldResponsibleUserID := existingSLAActionPlan.ResponsibleUserID

	if err := saps.validateSLAActionPlanFormDTO(slaID, slaActionPlanFormDTO); err != nil {
		return err
	}

	dates, err := saps.parseSLAActionPlanDates(slaActionPlanFormDTO)
	if err != nil {
		return err
	}

	// Update SLA ActionPlan
	existingSLAActionPlan.Title = slaActionPlanFormDTO.Title
	existingSLAActionPlan.SLACheckQuestionID = slaActionPlanFormDTO.SLACheckQuestionID
	existingSLAActionPlan.Description = slaActionPlanFormDTO.Description
	existingSLAActionPlan.SeverityLevelID = slaActionPlanFormDTO.SeverityLevelID
	existingSLAActionPlan.ResponsibleUserID = slaActionPlanFormDTO.ResponsibleUserID
	existingSLAActionPlan.ReviewUserID = slaActionPlanFormDTO.ReviewUserID
	existingSLAActionPlan.TargetDate = dates.TargetDate
	existingSLAActionPlan.Percentage = slaActionPlanFormDTO.Percentage
	existingSLAActionPlan.Comments = slaActionPlanFormDTO.Comments
	existingSLAActionPlan.UpdatedByID = &currentUserID

	err = saps.repo.Update(existingSLAActionPlan, tx)
	if err != nil {
		return err
	}

	// //update task
	urlParam := fmt.Sprintf(`{"sla_register_id":%d, "action_plan_id":%d}`, slaID, id)

	taskUpdateDTO := commonDTOs.TaskFormDTO{
		Title:             slaActionPlanFormDTO.Title,
		Description:       slaActionPlanFormDTO.Description,
		TargetDate:        dates.TargetDate,
		URLParam:          &urlParam,
		SeverityLevelID:   slaActionPlanFormDTO.SeverityLevelID,
		ResponsibleUserID: slaActionPlanFormDTO.ResponsibleUserID,
	}

	err = saps.taskSvc.Update(currentUserID, taskUpdateDTO, "sla_action_plan_id", id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	go func() {
		existingActionPlan, err := saps.GetFullDetailsByID(slaID, id)
		if err != nil {
			return
		}
		//log changes
		changedData := make(map[string]interface{})
		oldData := map[string]interface{}{
			"title":              oldSLAActionPlan.Title,
			"sla_check_question": utility.GetFieldValue(oldSLAActionPlan.SLACheckQuestion).(string),
			"description":        oldSLAActionPlan.Description,
			"severity_level":     oldSLAActionPlan.SeverityLevel.TitleEn,
			"responsible_user":   fmt.Sprintf("%s %s", oldSLAActionPlan.ResponsibleUser.FirstNameEN, oldSLAActionPlan.ResponsibleUser.LastNameEN),
			"review_user":        fmt.Sprintf("%s %s", oldSLAActionPlan.ReviewUser.FirstNameEN, oldSLAActionPlan.ReviewUser.LastNameEN),
			"target_date":        oldSLAActionPlan.TargetDateString,
			"percentage":         oldSLAActionPlan.Percentage,
			"comment":            oldSLAActionPlan.Comments,
			"action_plan_status": oldSLAActionPlan.ActionPlanStatus.TitleEn,
		}

		fmt.Println(oldData)

		newData := map[string]interface{}{
			"title":              existingActionPlan.Title,
			"sla_check_question": utility.GetFieldValue(existingActionPlan.SLACheckQuestion).(string),
			"description":        existingActionPlan.Description,
			"severity_level":     existingActionPlan.SeverityLevel.TitleEn,
			"responsible_user":   fmt.Sprintf("%s %s", existingActionPlan.ResponsibleUser.FirstNameEN, existingActionPlan.ResponsibleUser.LastNameEN),
			"review_user":        fmt.Sprintf("%s %s", existingActionPlan.ReviewUser.FirstNameEN, existingActionPlan.ReviewUser.LastNameEN),
			"target_date":        existingActionPlan.TargetDateString,
			"percentage":         existingActionPlan.Percentage,
			"comment":            existingActionPlan.Comments,
			"action_plan_status": existingActionPlan.ActionPlanStatus.TitleEn,
		}

		fmt.Println(newData)
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := saps.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.SLAActionPlan{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.SLAActionPlan{}.TableName(),
					"recordID":  id,
					"userID":    currentUserID,
					"error":     err.Error(),
				}

				utility.LogError("Failed to log activity", fields)
			}
		}
		// sends notification if responsible user has changed
		if oldResponsibleUserID != existingSLAActionPlan.ResponsibleUserID {
			urlParam := fmt.Sprintf(`{"sla_register_id":%d, "action_plan_id":%d}`, slaID, id)
			queryParam := ""
			entity := "sla-action-plan"

			oldResponsibleUserAction := "responsible-user-remove"
			oldResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserID:               oldResponsibleUserID,
				ModuleIdentifier:     "compliance",
				NotificationType:     string(commonDTOs.InformationType),
				HeadingEn:            "Sla Action Plan Responsibility Removed",
				HeadingAr:            "تمت إزالة مسؤولية خطة عمل Sla",
				TitleEn:              existingSLAActionPlan.Title,
				TitleAr:              existingSLAActionPlan.Title,
				Entity:               &entity,
				Action:               &oldResponsibleUserAction,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			newResponsibleUserAction := "responsible-user-assign"
			newResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserID:               existingSLAActionPlan.ResponsibleUserID,
				ModuleIdentifier:     "compliance",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Sla Action Plan Responsibility Assigned",
				HeadingAr:            "تم تعيين مسؤولية خطة عمل Sla",
				TitleEn:              existingSLAActionPlan.Title,
				TitleAr:              existingSLAActionPlan.Title,
				Entity:               &entity,
				Action:               &newResponsibleUserAction,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			oldResponsibleUser, err := saps.userSvc.GetByID(oldResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch old responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			newResponsibleUser, err := saps.userSvc.GetByID(existingSLAActionPlan.ResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch new responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			assignedBy, err := saps.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch assigning user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			appURL := utility.GetAppURL()

			oldResponsibleUserEmailTo := []string{oldResponsibleUser.Email}
			oldResponsibleUserSubject := fmt.Sprintf("Sla Action Plan/%s: Responsibility Removed", existingSLAActionPlan.Title)
			oldResponsibleUserTemplateRelativePath := "compliance/templates/email/sla_action_plan_remove_responsible_user_mail.html"
			oldResponsibleUserEmailData := map[string]interface{}{
				"FirstNameEn":           oldResponsibleUser.FirstNameEn,
				"LastNameEn":            oldResponsibleUser.LastNameEn,
				"SlaReferenceNumber":    sla.ReferenceNumber,
				"SlaTitle":              sla.TitleEn,
				"ActionPlanTitle":       existingSLAActionPlan.Title,
				"ActionPlanDescription": existingSLAActionPlan.Description,
				"RedirectURL":           fmt.Sprintf("%scompliance/sla-contracts/%d/action-plans?action_plan_id=%d", appURL, slaID, id),
			}

			oldResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              oldResponsibleUserEmailTo,
				Subject:              oldResponsibleUserSubject,
				TemplateRelativePath: oldResponsibleUserTemplateRelativePath,
				EmailData:            oldResponsibleUserEmailData,
			}

			newResponsibleUserEmailTo := []string{newResponsibleUser.Email}
			newResponsibleUserSubject := fmt.Sprintf("sla Action Plan/%s: Responsibility Assigned", existingSLAActionPlan.Title)
			newResponsibleUserTemplateRelativePath := "compliance/templates/email/sla_action_plan_change_responsible_user_mail.html"
			newResponsibleUserEmailData := map[string]interface{}{
				"FirstNameEn":           newResponsibleUser.FirstNameEn,
				"LastNameEn":            newResponsibleUser.LastNameEn,
				"SlaReferenceNumber":    sla.ReferenceNumber,
				"SlaTitle":              sla.TitleEn,
				"ActionPlanTitle":       existingSLAActionPlan.Title,
				"ActionPlanDescription": existingSLAActionPlan.Description,
				"TargetDate":            existingSLAActionPlan.TargetDate.Format("2006-01-02"),
				"AssignedByFirstNameEn": assignedBy.FirstNameEn,
				"AssignedByLastNameEn":  assignedBy.LastNameEn,
				"AssignedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           fmt.Sprintf("%scompliance/sla-contracts/%d/action-plans?action_plan_id=%d", appURL, slaID, id),
			}

			newResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              newResponsibleUserEmailTo,
				Subject:              newResponsibleUserSubject,
				TemplateRelativePath: newResponsibleUserTemplateRelativePath,
				EmailData:            newResponsibleUserEmailData,
			}

			err = saps.notificationSvc.SendNotification(
				newResponsibleUserNotificationDTO, newResponsibleUserEmailNotificationDTO, "sla_register_action_plan",
			)
			if err != nil {
				utility.LogError("Failed to send new responsible user notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}

			err = saps.notificationSvc.SendNotification(
				oldResponsibleUserNotificationDTO, oldResponsibleUserEmailNotificationDTO, "sla_register_action_plan",
			)
			if err != nil {
				utility.LogError("Failed to send old responsible user notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()

	return nil
}

func (saps *SLAActionPlanService) Delete(slaID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := saps.repo.BeginTransaction()

	sla, err := saps.slaRegisterSvc.GetByID(slaID)
	if err != nil {
		return nil, err
	}

	var failedIds []uint
	var successfulIds []uint
	var deletedSlaActionPlan []map[string]interface{}
	var notifications []struct {
		OwnerID         *uint
		TitleEn         string
		TitleAr         string
		Description     string
		ReferenceNumber string
		SlaTitle        string
	}
	for _, id := range ids {

		slaActionPlan, err := saps.GetFullDetailsByID(slaID, id)
		if err != nil {
			continue
		}

		// Prepare notification data before deletion
		notifications = append(notifications, struct {
			OwnerID         *uint
			TitleEn         string
			TitleAr         string
			Description     string
			ReferenceNumber string
			SlaTitle        string
		}{
			OwnerID:         &slaActionPlan.ResponsibleUserID,
			TitleEn:         slaActionPlan.Title,
			TitleAr:         slaActionPlan.Title,
			Description:     slaActionPlan.Description,
			ReferenceNumber: sla.ReferenceNumber,
			SlaTitle:        slaActionPlan.Title,
		})

		if len(slaActionPlan.ActionPlanStatusUpdates) > 0 {
			var deleteStatusUpdateIds []*uint
			for _, update := range slaActionPlan.ActionPlanStatusUpdates {
				deleteStatusUpdateIds = append(deleteStatusUpdateIds, &update.ID)
			}
			if err := saps.repo.DeleteSLAActionPlanStatusUpdates(deleteStatusUpdateIds, tx); err != nil {
				failedIds = append(failedIds, id)
				continue
			}
		}

		existingTask, err := saps.taskSvc.GetByFieldID("sla_action_plan_id", id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		taskID, err := saps.taskSvc.Delete(existingTask.ID, currentUserID)
		if err != nil {
			failedIds = append(failedIds, taskID)
		}

		err = saps.repo.Delete(id, tx)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
		deletedSlaActionPlan = append(deletedSlaActionPlan, map[string]interface{}{
			"id":              id,
			"sla_register_id": slaID,
			"title":           slaActionPlan.Title,
		})
	}

	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}

	tx.Commit()

	if len(successfulIds) > 0 {
		go func() {
			err := saps.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.SLARegister{}.TableName(),
				successfulIds,
				deletedSlaActionPlan,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}

			for _, notification := range notifications {
				urlParam := ""
				entity := "sla-action-plan"
				action := "delete"

				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserID:               *notification.OwnerID,
					ModuleIdentifier:     "compliance",
					NotificationType:     string(commonDTOs.DeletionType),
					HeadingEn:            "Sla Action Plan Deleted",
					HeadingAr:            "تم حذف خطة عمل Sla",
					TitleEn:              notification.TitleEn,
					TitleAr:              notification.TitleAr,
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				ownerUser, err := saps.userSvc.GetByID(*notification.OwnerID)
				if err != nil {
					utility.LogError("Failed to fetch owner details for notification", map[string]interface{}{
						"error": err.Error(),
					})
					continue
				}

				ownerEmailTo := []string{ownerUser.Email}
				subject := fmt.Sprintf("Sla Action Plan/%s: Deleted", notification.TitleEn)
				templateRelativePath := "compliance/templates/email/sla_action_plan_delete_mail.html"

				emailData := map[string]interface{}{
					"FirstNameEn":           ownerUser.FirstNameEn,
					"LastNameEn":            ownerUser.LastNameEn,
					"SlaReferenceNumber":    notification.ReferenceNumber,
					"SlaTitle":              notification.SlaTitle,
					"ActionPlanTitle":       notification.TitleEn,
					"ActionPlanDescription": notification.Description,
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              ownerEmailTo,
					Subject:              subject,
					TemplateRelativePath: templateRelativePath,
					EmailData:            emailData,
				}

				err = saps.notificationSvc.SendNotification(
					dbNotificationDTO, emailNotificationDTO, "sla_register_action_plan")
				if err != nil {
					utility.LogError("Failed to send deletion notifications", map[string]interface{}{
						"error": err.Error(),
					})
				}

			}
		}()
	}

	return failedIds, nil
}

func (saps *SLAActionPlanService) StatusUpdate(slaID uint, id uint, currentUserID uint, slaActionPlanStatusUpdateFormDTO *complianceDTOs.SLAActionPlanStatusUpdateFormDTO, tx utility.DatabaseTransaction) error {
	existingSLAActionPlan, err := saps.GetByID(slaID, id)
	if err != nil {
		return err
	}

	existingSLAActionPlan.ActionPlanStatusID = slaActionPlanStatusUpdateFormDTO.ActionPlanStatusID
	existingSLAActionPlan.Percentage = slaActionPlanStatusUpdateFormDTO.Percentage
	existingSLAActionPlan.UpdatedByID = &currentUserID

	err = saps.repo.Update(existingSLAActionPlan, tx)
	if err != nil {
		return err
	}

	isReviewUserNull := existingSLAActionPlan.ReviewUserID == nil

	// Fetch the corresponding task and update its status
	existingTask, err := saps.taskSvc.GetByFieldID("sla_action_plan_id", id)
	if err != nil {
		return err
	}
	err = saps.taskSvc.UpdateStatus(currentUserID, existingTask.ID, slaActionPlanStatusUpdateFormDTO.ActionPlanStatusID, isReviewUserNull, tx)
	if err != nil {
		return err
	}

	return nil
}

func (saps *SLAActionPlanService) validateSLAActionPlanFormDTO(slaID uint, slaActionPlanFormDTO *complianceDTOs.SLAActionPlanFormDTO) error {
	if slaActionPlanFormDTO.SeverityLevelID != nil {
		if _, err := saps.severityLevelSvc.GetByID(*slaActionPlanFormDTO.SeverityLevelID); err != nil {
			return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
		}
	}

	if slaActionPlanFormDTO.SLACheckQuestionID != nil {
		if _, err := saps.slaCheckQuestionSvc.GetQuestionByID(slaID, *slaActionPlanFormDTO.SLACheckQuestionID); err != nil {
			return utility.NewValidationError(fmt.Errorf("SLA requirement does not exist"), "not_exist", "SLA requirement")
		}
	}

	return nil
}

func (saps *SLAActionPlanService) parseSLAActionPlanDates(dto *complianceDTOs.SLAActionPlanFormDTO) (*complianceDTOs.SLAActionPlanParsedDates, error) {
	parsed := &complianceDTOs.SLAActionPlanParsedDates{}

	var err error
	if *dto.TargetDate != "" {
		parsed.TargetDate, err = utility.ParseDate(*dto.TargetDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid target date"), "invalid_date", "target date")
		}
	}

	return parsed, nil
}

func (saps *SLAActionPlanService) GetSLAActionPlanCountByStatus(slaID uint) ([]complianceDTOs.SLAActionPlanStatusCountResponseDTO, error) {
	if _, err := saps.slaRegisterSvc.GetByID(slaID); err != nil {
		return nil, err
	}

	response, err := saps.repo.GetSLAActionPlanCountByStatus(slaID)
	if err != nil {
		return nil, err
	}

	return response, nil
}
