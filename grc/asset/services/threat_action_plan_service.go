package services

import (
	"errors"
	"fmt"
	"os"
	"time"

	assetDTOs "jethur/grc/asset/dtos"
	assetRepositories "jethur/grc/asset/repositories"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	masterRepositories "jethur/grc/masters/repositories"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ThreatActionPlanServiceInterface interface {
	GetAll(threatID uint, params *assetDTOs.ThreatActionPlanListQueryParamsDTO) (*assetDTOs.ThreatActionPlanListResponseDTO, error)
	GetByID(threatID uint, id uint) (*models.ThreatActionPlan, error)
	GetFullDetailsByID(threatID uint, id uint) (*assetDTOs.ThreatActionPlanDetailDTO, error)
	Create(threatID uint, currentUserID uint, formDTO *assetDTOs.ThreatActionPlanFormDTO) (uint, error)
	Update(threatID uint, id uint, currentUserID uint, formDTO *assetDTOs.ThreatActionPlanFormDTO) error
	Delete(threatID uint, ids []uint, currentUserID uint) ([]uint, error)
	StatusUpdate(threatID uint, id uint, currentUserID uint, statusUpdateFormDTO *assetDTOs.ThreatActionPlanStatusUpdateFormDTO, tx utility.DatabaseTransaction) error
	ExportThreatActionPlan(data []*assetDTOs.ThreatActionPlanListDTO, currentUserID uint) (*os.File, error)

	GetThreatActionPlanSummary(threatID uint) (*assetDTOs.ThreatActionPlanSummaryDTO, error)
}

type ThreatActionPlanService struct {
	repo                assetRepositories.ThreatActionPlanRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	threatSvc           ThreatServiceInterface
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface
	severityLevelSvc    masterServices.SeverityLevelServiceInterface
	userSettingSvc      organizationServices.UserSettingServiceInterface
	userSvc             organizationServices.UserServiceInterface
	notificationSvc     commonServices.NotificationServiceInterface
}

func NewThreatActionPlanService(
	repo assetRepositories.ThreatActionPlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	threatSvc ThreatServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	severityLevelSvc masterRepositories.SeverityLevelRepositoryInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) ThreatActionPlanServiceInterface {
	return &ThreatActionPlanService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		threatSvc:           threatSvc,
		actionPlanStatusSvc: actionPlanStatusSvc,
		severityLevelSvc:    severityLevelSvc,
		userSettingSvc:      userSettingSvc,
		userSvc:             userSvc,
		notificationSvc:     notificationSvc,
	}
}

func (taps *ThreatActionPlanService) GetAll(threatID uint, params *assetDTOs.ThreatActionPlanListQueryParamsDTO) (*assetDTOs.ThreatActionPlanListResponseDTO, error) {
	if threatID != 0 {
		if _, err := taps.threatSvc.GetByID(threatID); err != nil {
			return nil, err
		}
	}
	return taps.repo.GetAll(threatID, params)
}

func (taps *ThreatActionPlanService) GetByID(threatID uint, id uint) (*models.ThreatActionPlan, error) {
	if _, err := taps.threatSvc.GetByID(threatID); err != nil {
		return nil, err
	}

	actionPlan, err := taps.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return actionPlan, nil
}

func (taps *ThreatActionPlanService) GetFullDetailsByID(threatID uint, id uint) (*assetDTOs.ThreatActionPlanDetailDTO, error) {
	var actionPlanDTO assetDTOs.ThreatActionPlanDetailDTO

	if _, err := taps.threatSvc.GetByID(threatID); err != nil {
		return nil, err
	}

	actionPlan, err := taps.repo.GetFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(actionPlan, &actionPlanDTO)
	if err != nil {
		return nil, err
	}

	if actionPlan.TargetDate != nil {
		targetDate := actionPlan.TargetDate.Format("2006-01-02")
		actionPlanDTO.TargetDateString = &targetDate
	}

	return &actionPlanDTO, nil
}

func (taps *ThreatActionPlanService) Create(threatID uint, currentUserID uint, formDTO *assetDTOs.ThreatActionPlanFormDTO) (uint, error) {
	tx := taps.repo.BeginTransaction()

	threat, err := taps.threatSvc.GetByID(threatID)
	if err != nil {
		return 0, err
	}

	if err := taps.validateActionPlanFormDTO(formDTO); err != nil {
		return 0, err
	}

	dates, err := taps.parseActionPlanDates(formDTO)
	if err != nil {
		return 0, err
	}

	actionPlanStatus, err := taps.actionPlanStatusSvc.GetByType("new")
	if err != nil {
		return 0, err
	}

	actionPlan := &models.ThreatActionPlan{
		ThreatID:           threatID,
		Title:              formDTO.Title,
		Description:        formDTO.Description,
		SeverityLevelID:    formDTO.SeverityLevelID,
		ResponsibleUserID:  formDTO.ResponsibleUserID,
		ReviewUserID:       formDTO.ReviewUserID,
		TargetDate:         dates.TargetDate,
		Percentage:         formDTO.Percentage,
		Comments:           formDTO.Comments,
		ActionPlanStatusID: actionPlanStatus.ID,
		CreatedByID:        currentUserID,
	}

	id, err := taps.repo.Create(actionPlan)
	if err != nil {
		return 0, err
	}

	actionPlanUpdate := &models.ThreatActionPlanStatusUpdate{
		ThreatActionPlanID: id,
		ActionPlanStatusID: actionPlanStatus.ID,
		Percentage:         formDTO.Percentage,
		Comment:            "Initial action plan status",
		CreatedByID:        currentUserID,
	}
	_, err = taps.repo.CreateInitialThreatActionPlanStatusUpdate(actionPlanUpdate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	go func() {
		existingActionPlan, err := taps.GetFullDetailsByID(threatID, id)
		if err != nil {
			utility.LogError("Failed to fetch threat action plan", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		err = taps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ThreatActionPlan{}.TableName(),
			id,
			existingActionPlan,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ThreatActionPlan{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}

		responsibleUser, err := taps.userSvc.GetByID(formDTO.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch responsible user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		assignedBy, err := taps.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch assigning user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		urlParam := fmt.Sprintf(`{"threat_id":%d, "action_plan_id":%d}`, threatID, id)
		queryParam := ""
		entity := "threat-action-plan"
		action := "responsible-user-assign"

		// Prepare DB Notification
		dbNotification := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{formDTO.ResponsibleUserID},
			ModuleIdentifier:     "threat",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "Threat Action Plan Responsibility Assigned",
			HeadingAr:            "تم تعيين مسؤولية خطة عمل التهديد",
			TitleEn:              formDTO.Title,
			TitleAr:              formDTO.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			QueryParam:           queryParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		// Prepare Email Notification
		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("Threat Action Plan/%s: Responsibility Assigned", formDTO.Title)
		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstNameEn":           responsibleUser.FirstNameEn,
			"LastNameEn":            responsibleUser.LastNameEn,
			"ThreatTitle":           threat.TitleEn,
			"ActionPlanTitle":       formDTO.Title,
			"ActionPlanDescription": formDTO.Description,
			"TargetDate":            dates.TargetDate.Format("2006-01-02"),
			"AssignedByFirstNameEn": assignedBy.FirstNameEn,
			"AssignedByLastNameEn":  assignedBy.LastNameEn,
			"AssignedOn":            time.Now().Format("2006-01-02"),
			"RedirectURL":           fmt.Sprintf("%s/asset/threats/%d/action-plans", appURL, id),
		}
		emailNotification := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "asset/templates/email/threat_action_plan_create_responsible_user_mail.html",
			EmailData:            emailData,
		}

		// Send Notification and Email
		err = taps.notificationSvc.SendNotification(dbNotification, emailNotification, "threat_action_plan")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	return id, nil
}

func (taps *ThreatActionPlanService) Update(threatID uint, id uint, currentUserID uint, formDTO *assetDTOs.ThreatActionPlanFormDTO) error {
	tx := taps.repo.BeginTransaction()

	oldActionPlan, err := taps.GetFullDetailsByID(threatID, id)
	if err != nil {
		return err
	}

	threat, err := taps.threatSvc.GetByID(threatID)
	if err != nil {
		return err
	}

	existingActionPlan, err := taps.GetByID(threatID, id)
	if err != nil {
		return err
	}

	oldResponsibleUserID := existingActionPlan.ResponsibleUserID

	if err := taps.validateActionPlanFormDTO(formDTO); err != nil {
		return err
	}

	dates, err := taps.parseActionPlanDates(formDTO)
	if err != nil {
		return err
	}

	existingActionPlan.Title = formDTO.Title
	existingActionPlan.Description = formDTO.Description
	existingActionPlan.SeverityLevelID = formDTO.SeverityLevelID
	existingActionPlan.ResponsibleUserID = formDTO.ResponsibleUserID
	existingActionPlan.ReviewUserID = formDTO.ReviewUserID
	existingActionPlan.TargetDate = dates.TargetDate
	existingActionPlan.Percentage = formDTO.Percentage
	existingActionPlan.Comments = formDTO.Comments
	existingActionPlan.UpdatedByID = &currentUserID

	err = taps.repo.Update(existingActionPlan, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		existingThreatActionPlan, err := taps.GetFullDetailsByID(threatID, id)
		if err != nil {
			return
		}
		// Prepare log data
		changedData := make(map[string]interface{})

		// log changes
		oldData := map[string]interface{}{
			"title":              oldActionPlan.Title,
			"description":        oldActionPlan.Description,
			"severity_level":     oldActionPlan.SeverityLevel.TitleEn,
			"responsible_user":   fmt.Sprintf("%s %s", oldActionPlan.ResponsibleUser.FirstNameEN, oldActionPlan.ResponsibleUser.LastNameEN),
			"review_user":        fmt.Sprintf("%s %s", oldActionPlan.ReviewUser.FirstNameEN, oldActionPlan.ReviewUser.LastNameEN),
			"target_date":        oldActionPlan.TargetDateString,
			"percentage":         oldActionPlan.Percentage,
			"comment":            oldActionPlan.Comments,
			"action_plan_status": oldActionPlan.ActionPlanStatus.TitleEn,
		}

		newData := map[string]interface{}{
			"title":              existingThreatActionPlan.Title,
			"description":        existingThreatActionPlan.Description,
			"severity_level":     existingThreatActionPlan.SeverityLevel.TitleEn,
			"responsible_user":   fmt.Sprintf("%s %s", existingThreatActionPlan.ResponsibleUser.FirstNameEN, existingThreatActionPlan.ResponsibleUser.LastNameEN),
			"review_user":        fmt.Sprintf("%s %s", existingThreatActionPlan.ReviewUser.FirstNameEN, existingThreatActionPlan.ReviewUser.LastNameEN),
			"target_date":        existingThreatActionPlan.TargetDateString,
			"percentage":         existingThreatActionPlan.Percentage,
			"comment":            existingThreatActionPlan.Comments,
			"action_plan_status": existingThreatActionPlan.ActionPlanStatus.TitleEn,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := taps.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.ThreatActionPlan{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.ThreatActionPlan{}.TableName(),
					"recordID":  id,
					"userID":    currentUserID,
					"error":     err.Error(),
				}

				utility.LogError("Failed to log activity", fields)
			}
		}

		// send notification if responsible user changed
		if oldResponsibleUserID != formDTO.ResponsibleUserID {
			urlParam := fmt.Sprintf(`{"threat_id":%d, "action_plan_id":%d}`, threatID, id)
			queryParam := ""
			entity := "threat-action-plan"

			oldResponsibleUserAction := "responsible-user-remove"
			oldResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{oldResponsibleUserID},
				ModuleIdentifier:     "threat",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Threat Action Plan Responsibility Removed",
				HeadingAr:            "تمت إزالة مسؤولية خطة عمل التهديد",
				TitleEn:              existingActionPlan.Title,
				TitleAr:              existingActionPlan.Title,
				Entity:               &entity,
				Action:               &oldResponsibleUserAction,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			newResponsibleUserAction := "responsible-user-assign"
			newResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{formDTO.ResponsibleUserID},
				ModuleIdentifier:     "threat",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Threat Action Plan Responsibility Assigned",
				HeadingAr:            "تم تعيين مسؤولية خطة عمل التهديد",
				TitleEn:              existingActionPlan.Title,
				TitleAr:              existingActionPlan.Title,
				Entity:               &entity,
				Action:               &newResponsibleUserAction,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			oldResponsibleUser, err := taps.userSvc.GetByID(oldResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch old responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			newResponsibleUser, err := taps.userSvc.GetByID(formDTO.ResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch new responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			assignedBy, err := taps.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch assigning user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			appURL := utility.GetAppURL()

			oldResponsibleUserEmailTo := []string{oldResponsibleUser.Email}
			oldResponsibleUserSubject := fmt.Sprintf("Threat Action Plan/%s: Responsibility Removed", existingActionPlan.Title)
			oldResponsibleUserTemplateRelativePath := "asset/templates/email/threat_action_plan_remove_responsible_user_mail.html"
			oldResponsibleUserEmailData := map[string]interface{}{
				"FirstNameEn":           oldResponsibleUser.FirstNameEn,
				"LastNameEn":            oldResponsibleUser.LastNameEn,
				"ThreatTitle":           threat.TitleEn,
				"ActionPlanTitle":       existingActionPlan.Title,
				"ActionPlanDescription": existingActionPlan.Description,
				"AssignedByFirstNameEn": assignedBy.FirstNameEn,
				"AssignedByLastNameEn":  assignedBy.LastNameEn,
				"AssignedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           fmt.Sprintf("%s/asset/threats/%d/action-plans", appURL, id),
			}

			oldResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              oldResponsibleUserEmailTo,
				Subject:              oldResponsibleUserSubject,
				TemplateRelativePath: oldResponsibleUserTemplateRelativePath,
				EmailData:            oldResponsibleUserEmailData,
			}

			newResponsibleUserEmailTo := []string{newResponsibleUser.Email}
			newResponsibleUserSubject := fmt.Sprintf("Threat Action Plan/%s: Responsibility Assigned", existingActionPlan.Title)
			newResponsibleUserTemplateRelativePath := "asset/templates/email/threat_action_plan_change_responsible_user_mail.html"
			newResponsibleUserEmailData := map[string]interface{}{
				"FirstNameEn":           newResponsibleUser.FirstNameEn,
				"LastNameEn":            newResponsibleUser.LastNameEn,
				"ThreatTitle":           threat.TitleEn,
				"ActionPlanTitle":       existingActionPlan.Title,
				"ActionPlanDescription": existingActionPlan.Description,
				"TargetDate":            existingActionPlan.TargetDate.Format("2006-01-02"),
				"AssignedByFirstNameEn": assignedBy.FirstNameEn,
				"AssignedByLastNameEn":  assignedBy.LastNameEn,
				"AssignedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           fmt.Sprintf("%s/asset/threats/%d/action-plans", appURL, id),
			}

			newResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              newResponsibleUserEmailTo,
				Subject:              newResponsibleUserSubject,
				TemplateRelativePath: newResponsibleUserTemplateRelativePath,
				EmailData:            newResponsibleUserEmailData,
			}

			err = taps.notificationSvc.SendNotification(
				newResponsibleUserNotificationDTO, newResponsibleUserEmailNotificationDTO, "threat_action_plan",
			)
			if err != nil {
				utility.LogError("Failed to send new responsible user notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}

			err = taps.notificationSvc.SendNotification(
				oldResponsibleUserNotificationDTO, oldResponsibleUserEmailNotificationDTO, "threat_action_plan",
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

func (taps *ThreatActionPlanService) Delete(threatID uint, ids []uint, currentUserID uint) ([]uint, error) {
	threat, err := taps.threatSvc.GetByID(threatID)
	if err != nil {
		return nil, err
	}

	var failedIds []uint
	var successfulIds []uint
	var deletedActionPlan []map[string]interface{}
	var notifications []struct {
		UserID      uint
		TitleEn     string
		TitleAr     string
		Description string
		ThreatTitle string
	}

	for _, id := range ids {
		actionPlan, err := taps.repo.GetFullDetailsByID(id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		// Prepare notification data before deletion
		notifications = append(notifications, struct {
			UserID      uint
			TitleEn     string
			TitleAr     string
			Description string
			ThreatTitle string
		}{
			UserID:      actionPlan.ResponsibleUserID,
			TitleEn:     actionPlan.Title,
			TitleAr:     actionPlan.Title,
			Description: *actionPlan.Description,
			ThreatTitle: threat.TitleEn,
		})

		if len(actionPlan.ActionPlanStatusUpdates) > 0 {
			var deleteStatusUpdateIds []*uint
			for _, update := range actionPlan.ActionPlanStatusUpdates {
				deleteStatusUpdateIds = append(deleteStatusUpdateIds, &update.ID)
			}
			if err := taps.repo.DeleteThreatActionPlanStatusUpdates(deleteStatusUpdateIds); err != nil {
				failedIds = append(failedIds, id)
				continue
			}
		}

		err = taps.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedActionPlan = append(deletedActionPlan, map[string]interface{}{
			"id":    id,
			"title": actionPlan.Title,
		})
	}

	if len(successfulIds) > 0 {
		go func() {
			err := taps.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.ThreatActionPlan{}.TableName(),
				successfulIds,
				deletedActionPlan,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}

			for _, notification := range notifications {
				urlParam := ""
				entity := "threat-action-plan"
				action := "delete"

				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserIDs:              []uint{notification.UserID},
					ModuleIdentifier:     "threat",
					NotificationType:     string(commonDTOs.DeletionType),
					HeadingEn:            "Threat Action Plan Deleted",
					HeadingAr:            "تم حذف خطة عمل التهديد",
					TitleEn:              notification.TitleEn,
					TitleAr:              notification.TitleAr,
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				ownerUser, err := taps.userSvc.GetByID(notification.UserID)
				if err != nil {
					utility.LogError("Failed to fetch owner details for notification", map[string]interface{}{
						"error": err.Error(),
					})
					continue
				}

				ownerEmailTo := []string{ownerUser.Email}
				subject := fmt.Sprintf("Threat Action Plan/%s: Deleted", notification.TitleEn)
				templateRelativePath := "asset/templates/email/threat_action_plan_delete_mail.html"

				emailData := map[string]interface{}{
					"FirstNameEn":           ownerUser.FirstNameEn,
					"LastNameEn":            ownerUser.LastNameEn,
					"ThreatTitle":           notification.ThreatTitle,
					"ActionPlanTitle":       notification.TitleEn,
					"ActionPlanDescription": notification.Description,
					"AssignedOn":            time.Now().Format("2006-01-02"),
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              ownerEmailTo,
					Subject:              subject,
					TemplateRelativePath: templateRelativePath,
					EmailData:            emailData,
				}

				err = taps.notificationSvc.SendNotification(
					dbNotificationDTO, emailNotificationDTO, "threat_action_plan")
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

func (taps *ThreatActionPlanService) StatusUpdate(threatID uint, id uint, currentUserID uint, statusUpdateFormDTO *assetDTOs.ThreatActionPlanStatusUpdateFormDTO, tx utility.DatabaseTransaction) error {
	existingActionPlan, err := taps.GetByID(threatID, id)
	if err != nil {
		return err
	}

	existingActionPlan.ActionPlanStatusID = statusUpdateFormDTO.ActionPlanStatusID
	existingActionPlan.Percentage = statusUpdateFormDTO.Percentage
	existingActionPlan.UpdatedByID = &currentUserID

	err = taps.repo.Update(existingActionPlan, tx)
	if err != nil {
		return err
	}

	return nil
}

func (taps *ThreatActionPlanService) validateActionPlanFormDTO(formDTO *assetDTOs.ThreatActionPlanFormDTO) error {
	if formDTO.SeverityLevelID != nil {
		if _, err := taps.severityLevelSvc.GetByID(*formDTO.SeverityLevelID); err != nil {
			return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
		}
	}

	return nil
}

func (taps *ThreatActionPlanService) parseActionPlanDates(dto *assetDTOs.ThreatActionPlanFormDTO) (*assetDTOs.ThreatActionPlanParsedDates, error) {
	parsed := &assetDTOs.ThreatActionPlanParsedDates{}

	var err error
	if *dto.TargetDate != "" {
		parsed.TargetDate, err = utility.ParseDate(*dto.TargetDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid target date"), "invalid_date", "target date")
		}
	}

	return parsed, nil
}

func (taps *ThreatActionPlanService) GetThreatActionPlanSummary(threatID uint) (*assetDTOs.ThreatActionPlanSummaryDTO, error) {
	openCount, closedCount, err := taps.repo.GetThreatActionPlanSummary(threatID)
	if err != nil {
		return nil, err
	}

	summary := &assetDTOs.ThreatActionPlanSummaryDTO{
		OpenCount:   openCount,
		ClosedCount: closedCount,
	}

	return summary, nil
}

func (taps *ThreatActionPlanService) ExportThreatActionPlan(data []*assetDTOs.ThreatActionPlanListDTO, currentUserID uint) (*os.File, error) {
	userSetting, err := taps.userSettingSvc.GetUserSettingsByUserID(currentUserID)
	if err != nil {
		return nil, err
	}

	var headers []string
	if userSetting.Language == "ar" {
		headers = []string{
			"SL No.",
			"Action Plan",
			"Threat",
			"Description",
			"Severity",
			"Responsible User",
			"Target Date",
			"Status",
		}
	} else {
		headers = []string{
			"SL No.",
			"Action Plan",
			"Threat",
			"Description",
			"Severity",
			"Responsible User",
			"Target Date",
			"Status",
		}
	}

	rows := make([][]interface{}, len(data))
	for i, actionPlan := range data {
		var responsibleUser, severity, status, threatTitle string

		if userSetting.Language == "ar" {
			responsibleUser = actionPlan.ResponsibleUserFirstNameAr
			severity = utility.GetStringValue(actionPlan.SeverityLevelTitleAr)
			status = actionPlan.ActionPlanStatusTitleAr
			threatTitle = actionPlan.ThreatTitleAr
		} else {
			responsibleUser = actionPlan.ResponsibleUserFirstNameEn
			severity = utility.GetStringValue(actionPlan.SeverityLevelTitleEn)
			status = actionPlan.ActionPlanStatusTitleEn
			threatTitle = actionPlan.ThreatTitleEn
		}

		rows[i] = []interface{}{
			i + 1,
			actionPlan.Title,
			threatTitle,
			utility.GetStringValue(actionPlan.Description),
			severity,
			responsibleUser,
			utility.GetStringValue(actionPlan.TargetDate),
			status,
		}
	}

	return utility.GenerateExcelFile(headers, rows)
}
