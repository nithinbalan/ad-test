package services

import (
	"errors"
	"fmt"
	"os"
	"time"

	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	meetingDTOs "jethur/grc/meeting/dtos"
	meetingRepositories "jethur/grc/meeting/repositories"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	taskServices "jethur/grc/task/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type MeetingActionPlanServiceInterface interface {
	GetAll(meetingID uint, params *meetingDTOs.MeetingActionPlanListQueryParamsDTO) (*meetingDTOs.MeetingActionPlanListResponseDTO, error)
	GetByID(meetingID uint, id uint) (*models.MeetingActionPlan, error)
	GetFullDetailsByID(meetingID uint, id uint) (*meetingDTOs.MeetingActionPlanDetailDTO, error)
	Create(meetingID uint, currentUserID uint, MeetingActionPlanFormDTO *meetingDTOs.MeetingActionPlanFormDTO) (uint, error)
	Update(meetingID uint, id uint, currentUserID uint, MeetingActionPlanFormDTO *meetingDTOs.MeetingActionPlanFormDTO) error
	Delete(meetingID uint, ids []uint, currentUserID uint) ([]uint, error)
	StatusUpdate(meetingID uint, id uint, currentUserID uint, meetingActionPlanStatusUpdateFormDTO *meetingDTOs.MeetingActionPlanStatusUpdateFormDTO, tx ...utility.DatabaseTransaction) error

	GetMeetingActionPlanSummary(meetingID uint) (*meetingDTOs.ActionPlanSummaryDTO, error)
	ExportMeetingActionPlan(data []*meetingDTOs.MeetingActionPlanListDTO, currentUserID uint) (*os.File, error)
}

// Service and implementations
type MeetingActionPlanService struct {
	repo                   meetingRepositories.MeetingActionPlanRepositoryInterface
	activityLogSvc         logServices.ActivityLogServiceInterface
	meetingSvc             MeetingServiceInterface
	severityLevelSvc       masterServices.SeverityLevelServiceInterface
	actionPlanStatusSvc    masterServices.ActionPlanStatusServiceInterface
	userSvc                organizationServices.UserServiceInterface
	notificationSvc        commonServices.NotificationServiceInterface
	meetingMinuteSvc       MeetingMinuteServiceInterface
	meetingMinuteStatusSvc MeetingMinuteStatusServiceInterface
	userSettingSvc         organizationServices.UserSettingServiceInterface
	taskSvc                taskServices.TaskServiceInterface
}

func NewMeetingActionPlanService(
	repo meetingRepositories.MeetingActionPlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	meetingSvc MeetingServiceInterface,
	severityLevelSvc masterServices.SeverityLevelServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	meetingMinuteSvc MeetingMinuteServiceInterface,
	meetingMinuteStatusSvc MeetingMinuteStatusServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
	taskSvc taskServices.TaskServiceInterface,
) MeetingActionPlanServiceInterface {
	return &MeetingActionPlanService{
		repo:                   repo,
		activityLogSvc:         activityLogSvc,
		meetingSvc:             meetingSvc,
		severityLevelSvc:       severityLevelSvc,
		actionPlanStatusSvc:    actionPlanStatusSvc,
		userSvc:                userSvc,
		notificationSvc:        notificationSvc,
		meetingMinuteSvc:       meetingMinuteSvc,
		meetingMinuteStatusSvc: meetingMinuteStatusSvc,
		userSettingSvc:         userSettingSvc,
		taskSvc:                taskSvc,
	}
}

func (maps *MeetingActionPlanService) GetAll(meetingID uint, params *meetingDTOs.MeetingActionPlanListQueryParamsDTO) (*meetingDTOs.MeetingActionPlanListResponseDTO, error) {
	if meetingID != 0 {
		if _, err := maps.meetingSvc.GetByID(meetingID); err != nil {
			return nil, err
		}
	}
	return maps.repo.GetAll(meetingID, params)
}

func (maps MeetingActionPlanService) GetByID(meetingID uint, id uint) (*models.MeetingActionPlan, error) {
	if _, err := maps.meetingSvc.GetByID(meetingID); err != nil {
		return nil, err
	}

	meetingActionPlan, err := maps.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return meetingActionPlan, nil
}

func (maps MeetingActionPlanService) GetFullDetailsByID(meetingID uint, id uint) (*meetingDTOs.MeetingActionPlanDetailDTO, error) {
	var meetingActionPlanDTO meetingDTOs.MeetingActionPlanDetailDTO

	if _, err := maps.meetingSvc.GetByID(meetingID); err != nil {
		return nil, err
	}

	meetingActionPlan, err := maps.repo.GetFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(meetingActionPlan, &meetingActionPlanDTO)
	if err != nil {
		return nil, err
	}

	if meetingActionPlan.TargetDate != nil {
		targetDate := meetingActionPlan.TargetDate.Format("2006-01-02")
		meetingActionPlanDTO.TargetDateString = &targetDate
	}
	// todo :- also add extra fields which required in details screen

	return &meetingActionPlanDTO, nil
}

func (maps *MeetingActionPlanService) Create(meetingID uint, currentUserID uint, MeetingActionPlanFormDTO *meetingDTOs.MeetingActionPlanFormDTO) (uint, error) {
	tx := maps.repo.BeginTransaction()

	meeting, err := maps.meetingSvc.GetByID(meetingID)
	if err != nil {
		return 0, err
	}

	if err := maps.validateMeetingActionPlanFormDTO(MeetingActionPlanFormDTO); err != nil {
		return 0, err
	}

	dates, err := maps.parseMeetingActionPlanDates(MeetingActionPlanFormDTO)
	if err != nil {
		return 0, err
	}

	meetingActionPlan := &models.MeetingActionPlan{
		MeetingID:          meetingID,
		Title:              MeetingActionPlanFormDTO.Title,
		Description:        MeetingActionPlanFormDTO.Description,
		SeverityLevelID:    MeetingActionPlanFormDTO.SeverityLevelID,
		ResponsibleUserID:  MeetingActionPlanFormDTO.ResponsibleUserID,
		ReviewUserID:       MeetingActionPlanFormDTO.ReviewUserID,
		MeetingMinuteID:    MeetingActionPlanFormDTO.MeetingMinuteID,
		TargetDate:         dates.TargetDate,
		Percentage:         *MeetingActionPlanFormDTO.Percentage,
		ActionPlanStatusID: MeetingActionPlanFormDTO.ActionPlanStatusID,
		CreatedByID:        currentUserID,
	}

	id, err := maps.repo.Create(meetingActionPlan, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	meetingMinuteStatus, err := maps.meetingMinuteStatusSvc.GetByType("open")
	if err != nil {
		return 0, err
	}
	if err := maps.meetingMinuteSvc.StatusUpdate(*MeetingActionPlanFormDTO.MeetingMinuteID, currentUserID, uint(meetingMinuteStatus.ID), tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	actionPlanStatusID := MeetingActionPlanFormDTO.ActionPlanStatusID
	meetingActionPlanUpdate := &models.MeetingActionPlanStatusUpdate{
		MeetingActionPlanID: id,
		ActionPlanStatusID:  actionPlanStatusID,
		Comment:             "Initial action plan status",
		CreatedByID:         currentUserID,
	}

	_, err = maps.repo.CreateInitialMeetingActionPlanStatusUpdate(meetingActionPlanUpdate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	urlParam := fmt.Sprintf(`{"id":%d}`, id)
	taskDTO := commonDTOs.TaskFormDTO{
		ModuleIdentifier:    "meeting",
		MeetingActionPlanID: &id,
		Title:               MeetingActionPlanFormDTO.Title,
		Description:         MeetingActionPlanFormDTO.Description,
		TargetDate:          dates.TargetDate,
		URLParam:            &urlParam,
		QueryParam:          nil,
		SeverityLevelID:     MeetingActionPlanFormDTO.SeverityLevelID,
		ResponsibleUserID:   MeetingActionPlanFormDTO.ResponsibleUserID,
	}

	if err := maps.taskSvc.Create(taskDTO, currentUserID, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	go func() {
		err := maps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.MeetingActionPlan{}.TableName(),
			id,
			MeetingActionPlanFormDTO,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.MeetingActionPlan{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}

		// prepare notification
		urlParam := fmt.Sprintf(`{"id":%d}`, id)
		entity := "meeting-action-plan"
		action := "responsible-user-assign"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserID:               MeetingActionPlanFormDTO.ResponsibleUserID,
			ModuleIdentifier:     "meeting",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "Meeting Action Plan Assignment",
			HeadingAr:            "تعيين خطة عمل الاجتماع",
			TitleEn:              MeetingActionPlanFormDTO.Title,
			TitleAr:              MeetingActionPlanFormDTO.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		responsibleUser, err := maps.userSvc.GetByID(MeetingActionPlanFormDTO.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch responsible user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		assignedBy, err := maps.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch assigning user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("Meeting Action Plan/%s: Responsibility Assigned", MeetingActionPlanFormDTO.Title)
		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstNameEn":            responsibleUser.FirstNameEn,
			"LastNameEn":             responsibleUser.LastNameEn,
			"MeetingReferenceNumber": meeting.ReferenceNumber,
			"MeetingTitle":           meeting.Title,
			"ActionPlanTitle":        MeetingActionPlanFormDTO.Title,
			"ActionPlanDescription":  MeetingActionPlanFormDTO.Description,
			"TargetDate":             dates.TargetDate.Format("2006-01-02"),
			"AssignedByFirstNameEn":  assignedBy.FirstNameEn,
			"AssignedByLastNameEn":   assignedBy.LastNameEn,
			"AssignedOn":             time.Now().Format("2006-01-02"),
			"RedirectURL":            fmt.Sprintf("%sreview-meeting/action-plan-list?action_plan_id=%d", appURL, id),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "meeting/templates/email/meeting_action_plan_create_responsible_user_mail.html",
			EmailData:            emailData,
		}

		// send notification and mail
		err = maps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "meeting_action_plan")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	return id, nil
}

func (maps *MeetingActionPlanService) Update(meetingID uint, id uint, currentUserID uint, MeetingActionPlanFormDTO *meetingDTOs.MeetingActionPlanFormDTO) error {
	tx := maps.repo.BeginTransaction()

	meeting, err := maps.meetingSvc.GetByID(meetingID)
	if err != nil {
		return err
	}

	existingMeetingActionPlan, err := maps.GetByID(meetingID, id)
	if err != nil {
		return err
	}

	oldResponsibleUserID := existingMeetingActionPlan.ResponsibleUserID

	if err := maps.validateMeetingActionPlanFormDTO(MeetingActionPlanFormDTO); err != nil {
		return err
	}

	dates, err := maps.parseMeetingActionPlanDates(MeetingActionPlanFormDTO)
	if err != nil {
		return err
	}

	// Update Meeting ActionPlan
	existingMeetingActionPlan.Title = MeetingActionPlanFormDTO.Title
	existingMeetingActionPlan.Description = MeetingActionPlanFormDTO.Description
	existingMeetingActionPlan.SeverityLevelID = MeetingActionPlanFormDTO.SeverityLevelID
	existingMeetingActionPlan.ResponsibleUserID = MeetingActionPlanFormDTO.ResponsibleUserID
	existingMeetingActionPlan.ReviewUserID = MeetingActionPlanFormDTO.ReviewUserID
	existingMeetingActionPlan.TargetDate = dates.TargetDate
	existingMeetingActionPlan.Percentage = *MeetingActionPlanFormDTO.Percentage
	existingMeetingActionPlan.ActionPlanStatusID = MeetingActionPlanFormDTO.ActionPlanStatusID
	existingMeetingActionPlan.MeetingMinuteID = MeetingActionPlanFormDTO.MeetingMinuteID
	existingMeetingActionPlan.UpdatedByID = &currentUserID

	err = maps.repo.Update(existingMeetingActionPlan)
	if err != nil {
		return err
	}

	urlParam := fmt.Sprintf(`{"id":%d}`, id)
	taskUpdateDTO := commonDTOs.TaskFormDTO{
		Title:             MeetingActionPlanFormDTO.Title,
		Description:       MeetingActionPlanFormDTO.Description,
		TargetDate:        dates.TargetDate,
		URLParam:          &urlParam,
		SeverityLevelID:   MeetingActionPlanFormDTO.SeverityLevelID,
		ResponsibleUserID: MeetingActionPlanFormDTO.ResponsibleUserID,
	}

	err = maps.taskSvc.Update(currentUserID, taskUpdateDTO, "meeting_action_plan_id", id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		// send notification if responsible user changed
		if oldResponsibleUserID != existingMeetingActionPlan.ResponsibleUserID {
			urlParam := fmt.Sprintf(`{"id":%d}`, id)
			entity := "meeting-action-plan"

			oldResponsibleUserAction := "responsible-user-remove"
			oldResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserID:               oldResponsibleUserID,
				ModuleIdentifier:     "meeting",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Meeting Action Plan Assignment",
				HeadingAr:            "تعيين خطة عمل الاجتماع",
				TitleEn:              existingMeetingActionPlan.Title,
				TitleAr:              existingMeetingActionPlan.Title,
				Entity:               &entity,
				Action:               &oldResponsibleUserAction,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			newResponsibleUserAction := "responsible-user-assign"
			newResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserID:               MeetingActionPlanFormDTO.ResponsibleUserID,
				ModuleIdentifier:     "meeting",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Meeting Action Plan Assignment",
				HeadingAr:            "تعيين خطة عمل الاجتماع",
				TitleEn:              existingMeetingActionPlan.Title,
				TitleAr:              existingMeetingActionPlan.Title,
				Entity:               &entity,
				Action:               &newResponsibleUserAction,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			oldResponsibleUser, err := maps.userSvc.GetByID(oldResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch old responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			newResponsibleUser, err := maps.userSvc.GetByID(MeetingActionPlanFormDTO.ResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch new responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			assignedBy, err := maps.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch assigning user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			appURL := utility.GetAppURL()

			oldResponsibleUserEmailTo := []string{oldResponsibleUser.Email}
			oldResponsibleUserSubject := fmt.Sprintf("Meeting Action Plan/%s: Responsibility Removed", existingMeetingActionPlan.Title)
			oldResponsibleUserTemplateRelativePath := "meeting/templates/email/meeting_action_plan_remove_responsible_user_mail.html"
			oldResponsibleUserEmailData := map[string]interface{}{
				"FirstNameEn":            oldResponsibleUser.FirstNameEn,
				"LastNameEn":             oldResponsibleUser.LastNameEn,
				"MeetingReferenceNumber": meeting.ReferenceNumber,
				"MeetingTitle":           meeting.Title,
				"ActionPlanTitle":        existingMeetingActionPlan.Title,
				"ActionPlanDescription":  existingMeetingActionPlan.Description,
				"RedirectURL":            fmt.Sprintf("%sreview-meeting/action-plan-list?action_plan_id=%d", appURL, id),
			}

			oldResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              oldResponsibleUserEmailTo,
				Subject:              oldResponsibleUserSubject,
				TemplateRelativePath: oldResponsibleUserTemplateRelativePath,
				EmailData:            oldResponsibleUserEmailData,
			}

			newResponsibleUserEmailTo := []string{newResponsibleUser.Email}
			newResponsibleUserSubject := fmt.Sprintf("Meeting Action Plan/%s: Responsibility Assigned", existingMeetingActionPlan.Title)
			newResponsibleUserTemplateRelativePath := "meeting/templates/email/meeting_action_plan_change_responsible_user_mail.html"
			newResponsibleUserEmailData := map[string]interface{}{
				"FirstNameEn":            newResponsibleUser.FirstNameEn,
				"LastNameEn":             newResponsibleUser.LastNameEn,
				"MeetingReferenceNumber": meeting.ReferenceNumber,
				"MeetingTitle":           meeting.Title,
				"ActionPlanTitle":        existingMeetingActionPlan.Title,
				"ActionPlanDescription":  existingMeetingActionPlan.Description,
				"TargetDate":             existingMeetingActionPlan.TargetDate.Format("2006-01-02"),
				"AssignedByFirstNameEn":  assignedBy.FirstNameEn,
				"AssignedByLastNameEn":   assignedBy.LastNameEn,
				"AssignedOn":             time.Now().Format("2006-01-02"),
				"RedirectURL":            fmt.Sprintf("%sreview-meeting/action-plan-list?action_plan_id=%d", appURL, id),
			}

			newResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              newResponsibleUserEmailTo,
				Subject:              newResponsibleUserSubject,
				TemplateRelativePath: newResponsibleUserTemplateRelativePath,
				EmailData:            newResponsibleUserEmailData,
			}

			err = maps.notificationSvc.SendNotification(
				newResponsibleUserNotificationDTO, newResponsibleUserEmailNotificationDTO, "meeting_action_plan",
			)
			if err != nil {
				utility.LogError("Failed to send new responsible user notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}

			err = maps.notificationSvc.SendNotification(
				oldResponsibleUserNotificationDTO, oldResponsibleUserEmailNotificationDTO, "meeting_action_plan",
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

func (maps *MeetingActionPlanService) Delete(meetingID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := maps.repo.BeginTransaction()

	meeting, err := maps.meetingSvc.GetByID(meetingID)
	if err != nil {
		return nil, err
	}

	var failedIds []uint
	var successfulIds []uint
	var notifications []struct {
		OwnerID         *uint
		TitleEn         string
		TitleAr         string
		Description     string
		ReferenceNumber string
		MeetingTitle    string
	}

	for _, id := range ids {
		meetingActionPlan, err := maps.repo.GetFullDetailsByID(id)
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
			MeetingTitle    string
		}{
			OwnerID:         &meetingActionPlan.ResponsibleUserID,
			TitleEn:         meetingActionPlan.Title,
			TitleAr:         meetingActionPlan.Title,
			Description:     *meetingActionPlan.Description,
			ReferenceNumber: meeting.ReferenceNumber,
			MeetingTitle:    meeting.Title,
		})

		existingTask, err := maps.taskSvc.GetByFieldID("meeting_action_plan_id", id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		taskID, err := maps.taskSvc.Delete(existingTask.ID, currentUserID)
		if err != nil {
			failedIds = append(failedIds, taskID)
		}

		if len(meetingActionPlan.ActionPlanStatusUpdates) > 0 {
			var deleteStatusUpdateIds []*uint
			for _, update := range meetingActionPlan.ActionPlanStatusUpdates {
				deleteStatusUpdateIds = append(deleteStatusUpdateIds, &update.ID)
			}
			if err := maps.repo.DeleteMeetingActionPlanStatusUpdates(deleteStatusUpdateIds, tx); err != nil {
				failedIds = append(failedIds, id)
				continue
			}
		}

		err = maps.repo.Delete(id, tx)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {
		tx.Commit()

		go func() {
			err := maps.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.MeetingActionPlan{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}

			for _, notification := range notifications {
				urlParam := ""
				entity := "meeting-action-plan"
				action := "delete"

				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserID:               *notification.OwnerID,
					ModuleIdentifier:     "meeting",
					NotificationType:     string(commonDTOs.DeletionType),
					HeadingEn:            "Meeting Action Plan Deletion",
					HeadingAr:            "حذف خطة عمل الاجتماع",
					TitleEn:              notification.TitleEn,
					TitleAr:              notification.TitleAr,
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				ownerUser, err := maps.userSvc.GetByID(*notification.OwnerID)
				if err != nil {
					utility.LogError("Failed to fetch owner details for notification", map[string]interface{}{
						"error": err.Error(),
					})
					continue
				}

				ownerEmailTo := []string{ownerUser.Email}
				subject := fmt.Sprintf("Meeting Action Plan/%s: Deleted", notification.TitleEn)
				templateRelativePath := "meeting/templates/email/meeting_action_plan_delete_mail.html"

				emailData := map[string]interface{}{
					"FirstNameEn":            ownerUser.FirstNameEn,
					"LastNameEn":             ownerUser.LastNameEn,
					"MeetingReferenceNumber": notification.ReferenceNumber,
					"MeetingTitle":           notification.MeetingTitle,
					"ActionPlanTitle":        notification.TitleEn,
					"ActionPlanDescription":  notification.Description,
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              ownerEmailTo,
					Subject:              subject,
					TemplateRelativePath: templateRelativePath,
					EmailData:            emailData,
				}

				err = maps.notificationSvc.SendNotification(
					dbNotificationDTO, emailNotificationDTO, "meeting_action_plan")
				if err != nil {
					utility.LogError("Failed to send deletion notifications", map[string]interface{}{
						"error": err.Error(),
					})
				}

			}
		}()
	} else {
		tx.Rollback()
	}

	return failedIds, nil
}

func (maps *MeetingActionPlanService) StatusUpdate(meetingID uint, id uint, currentUserID uint, meetingActionPlanStatusUpdateFormDTO *meetingDTOs.MeetingActionPlanStatusUpdateFormDTO, tx ...utility.DatabaseTransaction) error {
	existingMeetingActionPlan, err := maps.GetByID(meetingID, id)
	if err != nil {
		return err
	}

	existingMeetingActionPlan.ActionPlanStatusID = meetingActionPlanStatusUpdateFormDTO.ActionPlanStatusID
	existingMeetingActionPlan.Percentage = meetingActionPlanStatusUpdateFormDTO.Percentage
	existingMeetingActionPlan.Comments = &meetingActionPlanStatusUpdateFormDTO.Comment

	existingMeetingActionPlan.UpdatedByID = &currentUserID

	err = maps.repo.Update(existingMeetingActionPlan, tx...)
	if err != nil {
		return err
	}

	isReviewUserNull := existingMeetingActionPlan.ReviewUserID == nil

	// Fetch the corresponding task and update its status
	existingTask, err := maps.taskSvc.GetByFieldID("meeting_action_plan_id", id)
	if err != nil {
		return err
	}
	err = maps.taskSvc.UpdateStatus(currentUserID, existingTask.ID, meetingActionPlanStatusUpdateFormDTO.ActionPlanStatusID, isReviewUserNull, tx...)
	if err != nil {
		return err
	}

	return nil
}

func (maps *MeetingActionPlanService) validateMeetingActionPlanFormDTO(meetingActionPlanFormDTO *meetingDTOs.MeetingActionPlanFormDTO) error {
	if meetingActionPlanFormDTO.SeverityLevelID != nil {
		if _, err := maps.severityLevelSvc.GetByID(*meetingActionPlanFormDTO.SeverityLevelID); err != nil {
			return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
		}
	}

	if _, err := maps.userSvc.GetByID(uint(meetingActionPlanFormDTO.ResponsibleUserID)); err != nil {
		return utility.NewValidationError(fmt.Errorf("responsible user does not exist"), "not_exist", "responsible user")
	}

	if _, err := maps.actionPlanStatusSvc.GetByID(meetingActionPlanFormDTO.ActionPlanStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("action plan status does not exist"), "not_exist", "action plan status")
	}
	if _, err := maps.meetingMinuteSvc.GetByID(*meetingActionPlanFormDTO.MeetingMinuteID); err != nil {
		return utility.NewValidationError(fmt.Errorf("meeting minute does not exist"), "not_exist", "meeting minute")
	}

	return nil
}

func (maps *MeetingActionPlanService) parseMeetingActionPlanDates(dto *meetingDTOs.MeetingActionPlanFormDTO) (*meetingDTOs.MeetingActionPlanParsedDates, error) {
	parsed := &meetingDTOs.MeetingActionPlanParsedDates{}

	var err error
	if *dto.TargetDate != "" {
		parsed.TargetDate, err = utility.ParseDate(*dto.TargetDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid target date"), "invalid_date", "target date")
		}
	}

	return parsed, nil
}

func (maps *MeetingActionPlanService) GetMeetingActionPlanSummary(meetingID uint) (*meetingDTOs.ActionPlanSummaryDTO, error) {
	openCount, closedCount, err := maps.repo.GetMeetingActionPlanSummary(meetingID)
	if err != nil {
		return nil, err
	}

	summary := &meetingDTOs.ActionPlanSummaryDTO{
		OpenCount:   openCount,
		ClosedCount: closedCount,
	}

	return summary, nil
}

func (maps *MeetingActionPlanService) ExportMeetingActionPlan(data []*meetingDTOs.MeetingActionPlanListDTO, currentUserID uint) (*os.File, error) {
	// Fetch user setting for language preference
	userSetting, err := maps.userSettingSvc.GetUserSettingsByUserID(currentUserID)
	if err != nil {
		return nil, err
	}

	// Define headers based on user language setting
	var headers []string
	if userSetting.Language == "ar" {
		headers = []string{
			"SL No.",
			"Action Plan",
			"Meeting",
			"Description",
			"Severity",
			"Responsible User",
			"Target Date",
			"Percentage",
			"Status",
		}
	} else {
		headers = []string{
			"SL No.",
			"Action Plan",
			"Meeting",
			"Description",
			"Severity",
			"Responsible User",
			"Target Date",
			"Percentage",
			"Status",
		}
	}
	// Prepare the data rows
	rows := make([][]interface{}, len(data))
	for i, actionPlan := range data {
		var responsibleUser, severity, status string

		if userSetting.Language == "ar" {
			responsibleUser = actionPlan.ResponsibleUserFirstNameAr
			severity = utility.GetStringValue(actionPlan.SeverityLevelTitleAr)
			status = actionPlan.ActionPlanStatusTitleAr
		} else {
			responsibleUser = actionPlan.ResponsibleUserFirstNameEn
			severity = utility.GetStringValue(actionPlan.SeverityLevelTitleEn)
			status = actionPlan.ActionPlanStatusTitleEn
		}

		rows[i] = []interface{}{
			i + 1,
			actionPlan.Title,
			actionPlan.MeetingTitle,
			utility.GetStringValue(actionPlan.Description),
			severity,
			responsibleUser,
			utility.GetStringValue(actionPlan.TargetDate),
			actionPlan.Percentage,
			status,
		}
	}

	// Generate the Excel file
	return utility.GenerateExcelFile(headers, rows)
}
