package services

import (
	"errors"
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationRepositories "jethur/grc/organization/repositories"
	taskServices "jethur/grc/task/services"
	"jethur/utility"
	"time"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ControlActionPlanServiceInterface interface {
	GetAll(controlID uint, params *organizationDTOs.ControlActionPlanListQueryParamsDTO) (*organizationDTOs.ControlActionPlanListResponseDTO, error)
	GetByID(controlID uint, id uint) (*models.ControlActionPlan, error)
	GetFullDetailsByID(controlID uint, id uint) (*organizationDTOs.ControlActionPlanDetailDTO, error)
	Create(controlID uint, currentUserID uint, ControlActionPlanFormDTO *organizationDTOs.ControlActionPlanFormDTO) (uint, error)
	Update(controlID uint, id uint, currentUserID uint, ControlActionPlanFormDTO *organizationDTOs.ControlActionPlanFormDTO) error
	Delete(controlID uint, ids []uint, currentUserID uint) ([]uint, error)
	StatusUpdate(controlID uint, id uint, currentUserID uint, actionPlanStatusID uint, percentage uint, tx ...utility.DatabaseTransaction) error
}

// Service and implementations
type ControlActionPlanService struct {
	repo                organizationRepositories.ControlActionPlanRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	controlSvc          ControlServiceInterface
	severityLevelSvc    masterServices.SeverityLevelServiceInterface
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface
	userSvc             UserServiceInterface
	notificationSvc     commonServices.NotificationServiceInterface
	taskSvc             taskServices.TaskServiceInterface
}

func NewControlActionPlanService(
	repo organizationRepositories.ControlActionPlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	controlSvc ControlServiceInterface,
	severityLevelSvc masterServices.SeverityLevelServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	userSvc UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	taskSvc taskServices.TaskServiceInterface) ControlActionPlanServiceInterface {
	return &ControlActionPlanService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		controlSvc:          controlSvc,
		severityLevelSvc:    severityLevelSvc,
		actionPlanStatusSvc: actionPlanStatusSvc,
		userSvc:             userSvc,
		notificationSvc:     notificationSvc,
		taskSvc:             taskSvc,
	}
}

func (caps *ControlActionPlanService) GetAll(controlID uint, params *organizationDTOs.ControlActionPlanListQueryParamsDTO) (*organizationDTOs.ControlActionPlanListResponseDTO, error) {
	if _, err := caps.controlSvc.GetByID(controlID); err != nil {
		return nil, err
	}
	return caps.repo.GetAll(controlID, params)
}

func (caps ControlActionPlanService) GetByID(controlID uint, id uint) (*models.ControlActionPlan, error) {
	if _, err := caps.controlSvc.GetByID(controlID); err != nil {
		return nil, err
	}

	controlActionPlan, err := caps.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return controlActionPlan, nil
}

func (caps ControlActionPlanService) GetFullDetailsByID(controlID uint, id uint) (*organizationDTOs.ControlActionPlanDetailDTO, error) {
	var controlActionPlanDTO organizationDTOs.ControlActionPlanDetailDTO

	if _, err := caps.controlSvc.GetByID(controlID); err != nil {
		return nil, err
	}

	controlActionPlan, err := caps.repo.GetFullDetailsByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(controlActionPlan, &controlActionPlanDTO)
	if err != nil {
		return nil, err
	}

	//todo :- also add extra fields which required in details screen
	if controlActionPlan.TargetDate != nil {
		targetDate := controlActionPlan.TargetDate.Format("2006-01-02")
		controlActionPlanDTO.TargetDateString = &targetDate
	}

	controlActionPlanDTO.CreatedAt = controlActionPlan.CreatedAt
	controlActionPlanDTO.UpdatedAt = controlActionPlan.UpdatedAt

	return &controlActionPlanDTO, nil
}

func (caps *ControlActionPlanService) Create(controlID uint, currentUserID uint, ControlActionPlanFormDTO *organizationDTOs.ControlActionPlanFormDTO) (uint, error) {
	tx := caps.repo.BeginTransaction()

	control, err := caps.controlSvc.GetByID(controlID)
	if err != nil {
		return 0, err
	}

	if err := caps.validateControlActionPlanFormDTO(ControlActionPlanFormDTO); err != nil {
		return 0, err
	}

	dates, err := caps.parseControlActionPlanDates(ControlActionPlanFormDTO)
	if err != nil {
		return 0, err
	}

	controlActionPlan := &models.ControlActionPlan{
		ControlID:          controlID,
		Title:              ControlActionPlanFormDTO.Title,
		Description:        ControlActionPlanFormDTO.Description,
		SeverityLevelID:    ControlActionPlanFormDTO.SeverityLevelID,
		ResponsibleUserID:  ControlActionPlanFormDTO.ResponsibleUserID,
		TargetDate:         dates.TargetDate,
		Percentage:         *ControlActionPlanFormDTO.Percentage,
		ActionPlanStatusID: ControlActionPlanFormDTO.ActionPlanStatusID,
		CreatedByID:        currentUserID,
	}

	id, err := caps.repo.Create(controlActionPlan, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	actionPlanStatusID := ControlActionPlanFormDTO.ActionPlanStatusID
	controlActionPlanUpdate := &models.ControlActionPlanStatusUpdate{
		ControlActionPlanID: id,
		ActionPlanStatusID:  actionPlanStatusID,
		Comment:             "Initial action plan status",
		CreatedByID:         currentUserID,
	}

	_, err = caps.repo.CreateInitialControlActionPlanStatusUpdate(controlActionPlanUpdate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	urlParam := fmt.Sprintf(`{"id":%d}`, id)

	taskDTO := commonDTOs.TaskFormDTO{
		ModuleIdentifier:    "organization",
		ControlActionPlanID: &id,
		Title:               ControlActionPlanFormDTO.Title,
		Description:         ControlActionPlanFormDTO.Description,
		TargetDate:          dates.TargetDate,
		URLParam:            &urlParam,
		QueryParam:          nil,
		SeverityLevelID:     ControlActionPlanFormDTO.SeverityLevelID,
		ResponsibleUserID:   ControlActionPlanFormDTO.ResponsibleUserID,
	}

	if err := caps.taskSvc.Create(taskDTO, currentUserID, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	go func() {
		// Log activity
		err := caps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ControlActionPlan{}.TableName(),
			id,
			ControlActionPlanFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ControlActionPlan{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}

		// Prepare notification
		urlParam := fmt.Sprintf(`{"id":%d}`, id)
		entity := "control-action-plan"
		action := "responsible-user-assign"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserID:               ControlActionPlanFormDTO.ResponsibleUserID,
			ModuleIdentifier:     "organization",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "Control Action Plan Assignment",
			HeadingAr:            "تعيين خطة عمل التحكم",
			TitleEn:              ControlActionPlanFormDTO.Title,
			TitleAr:              ControlActionPlanFormDTO.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		responsibleUser, err := caps.userSvc.GetByID(ControlActionPlanFormDTO.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch responsible user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		assignedBy, err := caps.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch assigning user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("Control Action Plan/%s: Responsibility Assigned", ControlActionPlanFormDTO.Title)
		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstNameEn":            responsibleUser.FirstNameEn,
			"LastNameEn":             responsibleUser.LastNameEn,
			"ControlReferenceNumber": control.ReferenceNumber,
			"ControlTitle":           control.TitleEn,
			"ActionPlanTitle":        ControlActionPlanFormDTO.Title,
			"ActionPlanDescription":  ControlActionPlanFormDTO.Description,
			"TargetDate":             dates.TargetDate.Format("2006-01-02"),
			"AssignedByFirstNameEn":  assignedBy.FirstNameEn,
			"AssignedByLastNameEn":   assignedBy.LastNameEn,
			"AssignedOn":             time.Now().Format("2006-01-02"),
			"RedirectURL":            fmt.Sprintf("%s/control-action-plan/%d", appURL, id),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "organization/templates/email/control_action_plan_create_responsible_user_mail.html",
			EmailData:            emailData,
		}

		// Send Notification and Email
		err = caps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "control_action_plan")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	return id, nil
}

func (caps *ControlActionPlanService) Update(controlID uint, id uint, currentUserID uint, ControlActionPlanFormDTO *organizationDTOs.ControlActionPlanFormDTO) error {
	control, err := caps.controlSvc.GetByID(controlID)
	if err != nil {
		return err
	}

	existingControlActionPlan, err := caps.GetByID(controlID, id)
	if err != nil {
		return err
	}

	oldResponsibleUserID := existingControlActionPlan.ResponsibleUserID

	// Validate the incoming data
	if err := caps.validateControlActionPlanFormDTO(ControlActionPlanFormDTO); err != nil {
		return err
	}

	dates, err := caps.parseControlActionPlanDates(ControlActionPlanFormDTO)
	if err != nil {
		return err
	}

	existingControlActionPlan.Title = ControlActionPlanFormDTO.Title
	existingControlActionPlan.Description = ControlActionPlanFormDTO.Description
	existingControlActionPlan.SeverityLevelID = ControlActionPlanFormDTO.SeverityLevelID
	existingControlActionPlan.ResponsibleUserID = ControlActionPlanFormDTO.ResponsibleUserID
	existingControlActionPlan.TargetDate = dates.TargetDate
	existingControlActionPlan.Percentage = *ControlActionPlanFormDTO.Percentage
	existingControlActionPlan.ActionPlanStatusID = ControlActionPlanFormDTO.ActionPlanStatusID
	existingControlActionPlan.UpdatedByID = &currentUserID

	err = caps.repo.Update(existingControlActionPlan)
	if err != nil {
		return err
	}

	go func() {
		if oldResponsibleUserID != ControlActionPlanFormDTO.ResponsibleUserID {

			urlParam := fmt.Sprintf(`{"id":%d}`, id)
			entity := "control-action-plan"

			oldResponsibleUserAction := "responsible-user-remove"
			oldResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserID:               oldResponsibleUserID,
				ModuleIdentifier:     "organization",
				NotificationType:     string(commonDTOs.InformationType),
				HeadingEn:            "Control Action Plan Removal",
				HeadingAr:            "إزالة خطة عمل التحكم",
				TitleEn:              existingControlActionPlan.Title,
				TitleAr:              existingControlActionPlan.Title,
				Entity:               &entity,
				Action:               &oldResponsibleUserAction,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			newResponsibleUserAction := "responsible-user-assign"
			newResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserID:               ControlActionPlanFormDTO.ResponsibleUserID,
				ModuleIdentifier:     "organization",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Control Action Plan Assignment",
				HeadingAr:            "تعيين خطة عمل التحكم",
				TitleEn:              existingControlActionPlan.Title,
				TitleAr:              existingControlActionPlan.Title,
				Entity:               &entity,
				Action:               &newResponsibleUserAction,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			oldResponsibleUser, err := caps.userSvc.GetByID(oldResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch old responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			newResponsibleUser, err := caps.userSvc.GetByID(ControlActionPlanFormDTO.ResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch new responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			assignedBy, err := caps.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch assigning user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			appURL := utility.GetAppURL()

			oldResponsibleUserEmailTo := []string{oldResponsibleUser.Email}
			oldResponsibleUserSubject := fmt.Sprintf("Control Action Plan/%s: Responsibility Removed", existingControlActionPlan.Title)
			oldResponsibleUserTemplateRelativePath := "organization/templates/email/control_action_plan_remove_responsible_user_mail.html"
			oldResponsibleUserEmailData := map[string]interface{}{
				"FirstNameEn":            oldResponsibleUser.FirstNameEn,
				"LastNameEn":             oldResponsibleUser.LastNameEn,
				"ControlReferenceNumber": control.ReferenceNumber,
				"ControlTitle":           existingControlActionPlan.Title,
				"ActionPlanTitle":        existingControlActionPlan.Title,
				"ActionPlanDescription":  existingControlActionPlan.Description,
				"RedirectURL":            fmt.Sprintf("%s/control-action-plan/%d", appURL, id),
			}

			oldResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              oldResponsibleUserEmailTo,
				Subject:              oldResponsibleUserSubject,
				TemplateRelativePath: oldResponsibleUserTemplateRelativePath,
				EmailData:            oldResponsibleUserEmailData,
			}

			newResponsibleUserEmailTo := []string{newResponsibleUser.Email}
			newResponsibleUserSubject := fmt.Sprintf("Control Action Plan/%s: Responsibility Assigned", existingControlActionPlan.Title)
			newResponsibleUserTemplateRelativePath := "organization/templates/email/control_action_plan_change_responsible_user_mail.html"
			newResponsibleUserEmailData := map[string]interface{}{
				"FirstNameEn":            newResponsibleUser.FirstNameEn,
				"LastNameEn":             newResponsibleUser.LastNameEn,
				"ControlReferenceNumber": control.ReferenceNumber,
				"ControlTitle":           control.TitleEn,
				"ActionPlanTitle":        existingControlActionPlan.Title,
				"ActionPlanDescription":  existingControlActionPlan.Description,
				"TargetDate":             existingControlActionPlan.TargetDate.Format("2006-01-02"),
				"AssignedByFirstNameEn":  assignedBy.FirstNameEn,
				"AssignedByLastNameEn":   assignedBy.LastNameEn,
				"AssignedOn":             time.Now().Format("2006-01-02"),
				"RedirectURL":            fmt.Sprintf("%s/control-action-plan/%d", appURL, id),
			}

			newResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              newResponsibleUserEmailTo,
				Subject:              newResponsibleUserSubject,
				TemplateRelativePath: newResponsibleUserTemplateRelativePath,
				EmailData:            newResponsibleUserEmailData,
			}

			// Send notifications and emails
			err = caps.notificationSvc.SendNotification(
				newResponsibleUserNotificationDTO, newResponsibleUserEmailNotificationDTO, "control_action_plan",
			)
			if err != nil {
				utility.LogError("Failed to send new responsible user notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}

			err = caps.notificationSvc.SendNotification(
				oldResponsibleUserNotificationDTO, oldResponsibleUserEmailNotificationDTO, "control_action_plan",
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

func (caps *ControlActionPlanService) Delete(controlID uint, ids []uint, currentUserID uint) ([]uint, error) {
	if _, err := caps.controlSvc.GetByID(controlID); err != nil {
		return nil, err
	}

	var failedIds []uint
	var successfulIds []uint
	var notifications []struct {
		ResponsibleUserID *uint
		TitleEn           string
		TitleAr           string
		ReferenceNumber   string
	}

	for _, id := range ids {
		controlActionPlan, err := caps.GetFullDetailsByID(controlID, id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		notifications = append(notifications, struct {
			ResponsibleUserID *uint
			TitleEn           string
			TitleAr           string
			ReferenceNumber   string
		}{
			ResponsibleUserID: &controlActionPlan.ResponsibleUserID,
			TitleEn:           controlActionPlan.Title,
			TitleAr:           controlActionPlan.Title,
		})

		// Delete associated ActionPlanStatusUpdates if any
		if len(controlActionPlan.ActionPlanStatusUpdates) > 0 {
			var deleteStatusUpdateIds []*uint
			for _, update := range controlActionPlan.ActionPlanStatusUpdates {
				deleteStatusUpdateIds = append(deleteStatusUpdateIds, &update.ID)
			}
			if err := caps.repo.DeleteControlActionPlanStatusUpdates(deleteStatusUpdateIds); err != nil {
				failedIds = append(failedIds, id)
				continue
			}
		}

		err = caps.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {
		go func() {
			err := caps.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.ControlActionPlan{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}

			// Send notifications to the responsible users of the deleted ControlActionPlans
			for _, notification := range notifications {
				urlParam := ""
				queryParam := ""
				entity := "control-action-plan"
				action := "delete"

				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserID:               *notification.ResponsibleUserID,
					ModuleIdentifier:     "organization",
					NotificationType:     string(commonDTOs.DeletionType),
					ReferenceNumber:      notification.ReferenceNumber,
					HeadingEn:            "Control Action Plan Deletion",
					HeadingAr:            "حذف خطة عمل التحكم",
					TitleEn:              notification.TitleEn,
					TitleAr:              notification.TitleAr,
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					QueryParam:           queryParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				responsibleUser, err := caps.userSvc.GetByID(*notification.ResponsibleUserID)
				if err != nil {
					utility.LogError("Failed to fetch responsible user details for notification", map[string]interface{}{
						"error": err.Error(),
					})
					continue
				}

				emailTo := []string{responsibleUser.Email}
				subject := fmt.Sprintf("Control Action Plan/%s: Deleted", notification.TitleEn)
				templateRelativePath := "organization/templates/email/control_action_plan_delete_mail.html"
				redirectURL := fmt.Sprintf("control-action-plan/deleted/%d", *notification.ResponsibleUserID)

				emailData := organizationDTOs.ControlActionPlanMailDTO{
					FirstNameEn: responsibleUser.FirstNameEn,
					LastNameEn:  responsibleUser.LastNameEn,
					MessageEn:   notification.TitleEn,
					RedirectURL: redirectURL,
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              emailTo,
					Subject:              subject,
					TemplateRelativePath: templateRelativePath,
					EmailData:            emailData,
				}

				err = caps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "control_action_plan")
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

func (caps *ControlActionPlanService) StatusUpdate(controlID uint, id uint, currentUserID uint, actionPlanStatusID uint, percentage uint, tx ...utility.DatabaseTransaction) error {
	existingControlActionPlan, err := caps.GetByID(controlID, id)

	if err != nil {
		return err
	}

	existingControlActionPlan.ActionPlanStatusID = actionPlanStatusID
	existingControlActionPlan.Percentage = percentage
	existingControlActionPlan.UpdatedByID = &currentUserID

	err = caps.repo.Update(existingControlActionPlan, tx...)

	if err != nil {
		return err
	}

	return nil
}

func (caps *ControlActionPlanService) validateControlActionPlanFormDTO(ControlActionPlanFormDTO *organizationDTOs.ControlActionPlanFormDTO) error {
	if ControlActionPlanFormDTO.SeverityLevelID != nil {
		if _, err := caps.severityLevelSvc.GetByID(*ControlActionPlanFormDTO.SeverityLevelID); err != nil {
			return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
		}
	}

	if _, err := caps.userSvc.GetByID(uint(ControlActionPlanFormDTO.ResponsibleUserID)); err != nil {
		return utility.NewValidationError(fmt.Errorf("responsible user does not exist"), "not_exist", "responsible user")
	}

	if _, err := caps.actionPlanStatusSvc.GetByID(ControlActionPlanFormDTO.ActionPlanStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("action plan status does not exist"), "not_exist", "action plan status")
	}

	if ControlActionPlanFormDTO.Percentage == nil {
		defaultValue := uint(0)
		ControlActionPlanFormDTO.Percentage = &defaultValue
	}

	return nil
}

func (caps *ControlActionPlanService) parseControlActionPlanDates(dto *organizationDTOs.ControlActionPlanFormDTO) (*organizationDTOs.ControlActionPlanParsedDates, error) {
	parsed := &organizationDTOs.ControlActionPlanParsedDates{}

	var err error
	if *dto.TargetDate != "" {
		parsed.TargetDate, err = utility.ParseDate(*dto.TargetDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid target date"), "invalid_date", "target date")
		}
	}

	return parsed, nil
}
