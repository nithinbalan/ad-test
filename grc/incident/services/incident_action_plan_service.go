package services

import (
	"errors"
	"fmt"
	"os"
	"time"

	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	incidentDTOs "jethur/grc/incident/dtos"
	incidentRepositories "jethur/grc/incident/repositories"
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

type IncidentActionPlanServiceInterface interface {
	GetAll(incidentID uint, params *incidentDTOs.IncidentActionPlanListQueryParamsDTO) (*incidentDTOs.IncidentActionPlanListResponseDTO, error)
	GetByID(incidentID uint, id uint) (*models.IncidentActionPlan, error)
	GetFullDetailsByID(incidentID uint, id uint) (*incidentDTOs.IncidentActionPlanDetailDTO, error)
	Create(incidentID uint, currentUserID uint, IncidentActionPlanFormDTO *incidentDTOs.IncidentActionPlanFormDTO) (uint, error)
	Update(incidentID uint, id uint, currentUserID uint, incident *incidentDTOs.IncidentActionPlanFormDTO) error
	Delete(incidentID uint, ids []uint, currentUserID uint) ([]uint, error)
	StatusUpdate(incidentID uint, id uint, currentUserID uint, incidentActionPlanStatusUpdateFormDTO *incidentDTOs.IncidentActionPlanStatusUpdateFormDTO, tx utility.DatabaseTransaction) error

	ExportIncidentActionPlan(data []*incidentDTOs.IncidentActionPlanListDTO, currentUserID uint) (*os.File, error)
}

// Service and implementations
type IncidentActionPlanService struct {
	repo                incidentRepositories.IncidentActionPlanRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	incidentRegisterSvc IncidentRegisterServiceInterface
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface
	severityLevelSvc    masterServices.SeverityLevelServiceInterface
	incidentTimelineSvc IncidentTimelineServiceInterface
	userSettingSvc      organizationServices.UserSettingServiceInterface
	taskSvc             taskServices.TaskServiceInterface
	userSvc             organizationServices.UserServiceInterface
	notificationSvc     commonServices.NotificationServiceInterface
}

func NewIncidentActionPlanService(
	repo incidentRepositories.IncidentActionPlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	incidentRegisterSvc IncidentRegisterServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	severityLevelSvc masterRepositories.SeverityLevelRepositoryInterface,
	incidentTimelineSvc IncidentTimelineServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
	taskSvc taskServices.TaskServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) IncidentActionPlanServiceInterface {
	return &IncidentActionPlanService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		incidentRegisterSvc: incidentRegisterSvc,
		actionPlanStatusSvc: actionPlanStatusSvc,
		severityLevelSvc:    severityLevelSvc,
		incidentTimelineSvc: incidentTimelineSvc,
		userSettingSvc:      userSettingSvc,
		taskSvc:             taskSvc,
		userSvc:             userSvc,
		notificationSvc:     notificationSvc,
	}
}

func (iaps *IncidentActionPlanService) GetAll(incidentID uint, params *incidentDTOs.IncidentActionPlanListQueryParamsDTO) (*incidentDTOs.IncidentActionPlanListResponseDTO, error) {
	if incidentID != 0 {
		if _, err := iaps.incidentRegisterSvc.GetByID(incidentID); err != nil {
			return nil, err
		}
	}
	return iaps.repo.GetAll(incidentID, params)
}

func (iaps IncidentActionPlanService) GetByID(incidentID uint, id uint) (*models.IncidentActionPlan, error) {
	if _, err := iaps.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return nil, err
	}

	incidentActionPlan, err := iaps.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return incidentActionPlan, nil
}

func (iaps IncidentActionPlanService) GetFullDetailsByID(incidentID uint, id uint) (*incidentDTOs.IncidentActionPlanDetailDTO, error) {
	var incidentActionPlanDTO incidentDTOs.IncidentActionPlanDetailDTO

	if _, err := iaps.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return nil, err
	}

	incidentActionPlan, err := iaps.repo.GetFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(incidentActionPlan, &incidentActionPlanDTO)
	if err != nil {
		return nil, err
	}

	// todo :- also add extra fields which required in details screen
	if incidentActionPlan.TargetDate != nil {
		targetDate := incidentActionPlan.TargetDate.Format("2006-01-02")
		incidentActionPlanDTO.TargetDateString = &targetDate
	}

	incidentActionPlanDTO.CreatedAt = incidentActionPlan.CreatedAt
	incidentActionPlanDTO.UpdatedAt = incidentActionPlan.UpdatedAt

	return &incidentActionPlanDTO, nil
}

func (iaps *IncidentActionPlanService) Create(incidentID uint, currentUserID uint, IncidentActionPlanFormDTO *incidentDTOs.IncidentActionPlanFormDTO) (uint, error) {
	tx := iaps.repo.BeginTransaction()

	incident, err := iaps.incidentRegisterSvc.GetByID(incidentID)
	if err != nil {
		return 0, err
	}

	if err := iaps.validateIncidentActionPlanFormDTO(IncidentActionPlanFormDTO); err != nil {
		return 0, err
	}

	dates, err := iaps.parseIncidentActionPlanDates(IncidentActionPlanFormDTO)
	if err != nil {
		return 0, err
	}

	actionPlanStatus, err := iaps.actionPlanStatusSvc.GetByType("new")
	if err != nil {
		return 0, err
	}

	incidentActionPlan := &models.IncidentActionPlan{
		IncidentRegisterID: incidentID,
		Title:              IncidentActionPlanFormDTO.Title,
		Description:        IncidentActionPlanFormDTO.Description,
		SeverityLevelID:    IncidentActionPlanFormDTO.SeverityLevelID,
		ResponsibleUserID:  IncidentActionPlanFormDTO.ResponsibleUserID,
		ReviewUserID:       IncidentActionPlanFormDTO.ReviewUserID,
		TargetDate:         dates.TargetDate,
		Percentage:         IncidentActionPlanFormDTO.Percentage,
		Comments:           IncidentActionPlanFormDTO.Comments,
		ActionPlanStatusID: actionPlanStatus.ID,
		CreatedBYID:        currentUserID,
	}

	id, err := iaps.repo.Create(incidentActionPlan, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	actionPlanStatusID := actionPlanStatus.ID
	incidentActionPlanUpdate := &models.IncidentActionPlanStatusUpdate{
		IncidentActionPlanID: id,
		ActionPlanStatusID:   actionPlanStatusID,
		Percentage:           IncidentActionPlanFormDTO.Percentage,
		Comment:              "Initial action plan status",
		CreatedByID:          currentUserID,
	}
	_, err = iaps.repo.CreateInitialIncidentActionPlanStatusUpdate(incidentActionPlanUpdate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	// Update timeline
	// todo
	actinPlanCount, err := iaps.repo.GetCountByRiskID(incidentID, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if actinPlanCount == 1 {
		timelineType, err := iaps.incidentTimelineSvc.GetByType("action-plans")
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		timelineStatus, err := iaps.incidentTimelineSvc.GetByStatusType("completed")
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		currentTime := time.Now()
		timelineFormDTO := &incidentDTOs.IncidentTimelineFormDTO{
			PerformedByID:    &currentUserID,
			TimelineStatusID: timelineStatus.ID,
			PerformedAt:      &currentTime,
		}
		if err := iaps.incidentTimelineSvc.Update(incidentID, timelineFormDTO, timelineType.ID, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	urlParam := fmt.Sprintf(`{"incident_register_id":%d, "action_plan_id":%d}`, incidentID, id)
	taskDTO := commonDTOs.TaskFormDTO{
		ModuleIdentifier:     "incident",
		IncidentActionPlanID: &id,
		Title:                IncidentActionPlanFormDTO.Title,
		Description:          IncidentActionPlanFormDTO.Description,
		TargetDate:           dates.TargetDate,
		URLParam:             &urlParam,
		QueryParam:           nil,
		SeverityLevelID:      IncidentActionPlanFormDTO.SeverityLevelID,
		ResponsibleUserID:    IncidentActionPlanFormDTO.ResponsibleUserID,
	}

	if _, err := iaps.taskSvc.Create(taskDTO, currentUserID, tx); err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()

	go func() {

		existingIncidentActionPlan, err := iaps.GetFullDetailsByID(incidentID, id)
		if err != nil {
			utility.LogError("Failed to fetch incident action plan", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = iaps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.IncidentActionPlan{}.TableName(),
			id,
			existingIncidentActionPlan,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.IncidentActionPlan{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}

		// prepare notification
		urlParam := fmt.Sprintf(`{"incident_register_id":%d, "action_plan_id":%d}`, incidentID, id)
		queryParam := ""
		entity := "incident-action-plan"
		action := "responsible-user-assign"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserID:               IncidentActionPlanFormDTO.ResponsibleUserID,
			ModuleIdentifier:     "incident",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "Incident Action Plan Responsibility Assigned",
			HeadingAr:            "تم تعيين مسؤولية خطة عمل الحادث",
			TitleEn:              IncidentActionPlanFormDTO.Title,
			TitleAr:              IncidentActionPlanFormDTO.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			QueryParam:           queryParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		responsibleUser, err := iaps.userSvc.GetByID(IncidentActionPlanFormDTO.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch responsible user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		assignedBy, err := iaps.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch assigning user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("Incident Action Plan/%s: Responsibility Assigned", IncidentActionPlanFormDTO.Title)
		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstNameEn":             responsibleUser.FirstNameEn,
			"LastNameEn":              responsibleUser.LastNameEn,
			"IncidentReferenceNumber": incident.ReferenceNumber,
			"IncidentTitle":           incident.Title,
			"ActionPlanTitle":         IncidentActionPlanFormDTO.Title,
			"ActionPlanDescription":   IncidentActionPlanFormDTO.Description,
			"TargetDate":              dates.TargetDate.Format("2006-01-02"),
			"AssignedByFirstNameEn":   assignedBy.FirstNameEn,
			"AssignedByLastNameEn":    assignedBy.LastNameEn,
			"AssignedOn":              time.Now().Format("2006-01-02"),
			"RedirectURL":             fmt.Sprintf("%sincident/action-plan?action_plan_id=%d", appURL, id),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "incident/templates/email/incident_action_plan_create_responsible_user_mail.html",
			EmailData:            emailData,
		}

		// send notification and mail
		err = iaps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "incident_action_plan")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	return id, nil
}

func (iaps *IncidentActionPlanService) Update(incidentID uint, id uint, currentUserID uint, IncidentActionPlanFormDTO *incidentDTOs.IncidentActionPlanFormDTO) error {
	tx := iaps.repo.BeginTransaction()

	incident, err := iaps.incidentRegisterSvc.GetByID(incidentID)
	if err != nil {
		return err
	}

	existingIncidentActionPlan, err := iaps.GetByID(incidentID, id)
	if err != nil {
		return err
	}

	oldIncidentActionPlan, err := iaps.GetFullDetailsByID(incidentID, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	oldResponsibleUserID := existingIncidentActionPlan.ResponsibleUserID

	if err := iaps.validateIncidentActionPlanFormDTO(IncidentActionPlanFormDTO); err != nil {
		return err
	}

	dates, err := iaps.parseIncidentActionPlanDates(IncidentActionPlanFormDTO)
	if err != nil {
		return err
	}

	// Update Incident ActionPlan
	existingIncidentActionPlan.Title = IncidentActionPlanFormDTO.Title
	existingIncidentActionPlan.Description = IncidentActionPlanFormDTO.Description
	existingIncidentActionPlan.SeverityLevelID = IncidentActionPlanFormDTO.SeverityLevelID
	existingIncidentActionPlan.ResponsibleUserID = IncidentActionPlanFormDTO.ResponsibleUserID
	existingIncidentActionPlan.ReviewUserID = IncidentActionPlanFormDTO.ReviewUserID
	existingIncidentActionPlan.TargetDate = dates.TargetDate
	existingIncidentActionPlan.Percentage = IncidentActionPlanFormDTO.Percentage
	existingIncidentActionPlan.Comments = IncidentActionPlanFormDTO.Comments
	existingIncidentActionPlan.UpdatedBYID = &currentUserID

	err = iaps.repo.Update(existingIncidentActionPlan, tx)
	if err != nil {
		return err
	}

	urlParam := fmt.Sprintf(`{"incident_register_id":%d, "action_plan_id":%d}`, incidentID, id)

	taskUpdateDTO := commonDTOs.TaskFormDTO{
		Title:             IncidentActionPlanFormDTO.Title,
		Description:       IncidentActionPlanFormDTO.Description,
		TargetDate:        dates.TargetDate,
		URLParam:          &urlParam,
		SeverityLevelID:   IncidentActionPlanFormDTO.SeverityLevelID,
		ResponsibleUserID: IncidentActionPlanFormDTO.ResponsibleUserID,
	}

	err = iaps.taskSvc.Update(currentUserID, taskUpdateDTO, "incident_action_plan_id", id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		existingActionPlan, err := iaps.GetFullDetailsByID(incidentID, id)
		if err != nil {
			return
		}
		//log changes
		changedData := make(map[string]interface{})
		oldData := map[string]interface{}{
			"title":              oldIncidentActionPlan.Title,
			"description":        oldIncidentActionPlan.Description,
			"severity_level":     oldIncidentActionPlan.SeverityLevel.TitleEn,
			"responsible_user":   fmt.Sprintf("%s %s", oldIncidentActionPlan.ResponsibleUser.FirstNameEN, oldIncidentActionPlan.ResponsibleUser.LastNameEN),
			"review_user":        fmt.Sprintf("%s %s", oldIncidentActionPlan.ReviewUser.FirstNameEN, oldIncidentActionPlan.ReviewUser.LastNameEN),
			"target_date":        oldIncidentActionPlan.TargetDateString,
			"percentage":         oldIncidentActionPlan.Percentage,
			"comment":            oldIncidentActionPlan.Comments,
			"action_plan_status": oldIncidentActionPlan.ActionPlanStatus.TitleEn,
		}

		newData := map[string]interface{}{
			"title":              existingActionPlan.Title,
			"description":        existingActionPlan.Description,
			"severity_level":     existingActionPlan.SeverityLevel.TitleEn,
			"responsible_user":   fmt.Sprintf("%s %s", existingActionPlan.ResponsibleUser.FirstNameEN, existingActionPlan.ResponsibleUser.LastNameEN),
			"review_user":        fmt.Sprintf("%s %s", existingActionPlan.ReviewUser.FirstNameEN, existingActionPlan.ReviewUser.LastNameEN),
			"target_date":        existingActionPlan.TargetDateString,
			"percentage":         existingActionPlan.Percentage,
			"comment":            existingActionPlan.Comments,
			"action_plan_status": existingActionPlan.ActionPlanStatus.TitleEn,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := iaps.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.IncidentActionPlan{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.IncidentActionPlan{}.TableName(),
					"recordID":  id,
					"userID":    currentUserID,
					"error":     err.Error(),
				}

				utility.LogError("Failed to log activity", fields)
			}
		}
		// sends notification if responsible user has changed
		if oldResponsibleUserID != existingIncidentActionPlan.ResponsibleUserID {
			urlParam := fmt.Sprintf(`{"incident_register_id":%d, "action_plan_id":%d}`, incidentID, id)
			queryParam := ""
			entity := "incident-action-plan"

			oldResponsibleUserAction := "responsible-user-remove"
			oldResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserID:               oldResponsibleUserID,
				ModuleIdentifier:     "incident",
				NotificationType:     string(commonDTOs.InformationType),
				HeadingEn:            "Incident Action Plan Responsibility Removed",
				HeadingAr:            "تمت إزالة مسؤولية خطة عمل الحادث",
				TitleEn:              existingIncidentActionPlan.Title,
				TitleAr:              existingIncidentActionPlan.Title,
				Entity:               &entity,
				Action:               &oldResponsibleUserAction,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			newResponsibleUserAction := "responsible-user-assign"
			newResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserID:               IncidentActionPlanFormDTO.ResponsibleUserID,
				ModuleIdentifier:     "incident",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Incident Action Plan Responsibility Assigned",
				HeadingAr:            "تم تعيين مسؤولية خطة عمل الحادث",
				TitleEn:              existingIncidentActionPlan.Title,
				TitleAr:              existingIncidentActionPlan.Title,
				Entity:               &entity,
				Action:               &newResponsibleUserAction,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			oldResponsibleUser, err := iaps.userSvc.GetByID(oldResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch old responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			newResponsibleUser, err := iaps.userSvc.GetByID(IncidentActionPlanFormDTO.ResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch new responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			assignedBy, err := iaps.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch assigning user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			appURL := utility.GetAppURL()

			oldResponsibleUserEmailTo := []string{oldResponsibleUser.Email}
			oldResponsibleUserSubject := fmt.Sprintf("Incident Action Plan/%s: Responsibility Removed", existingIncidentActionPlan.Title)
			oldResponsibleUserTemplateRelativePath := "incident/templates/email/incident_action_plan_remove_responsible_user_mail.html"
			oldResponsibleUserEmailData := map[string]interface{}{
				"FirstNameEn":             oldResponsibleUser.FirstNameEn,
				"LastNameEn":              oldResponsibleUser.LastNameEn,
				"IncidentReferenceNumber": incident.ReferenceNumber,
				"IncidentTitle":           incident.Title,
				"ActionPlanTitle":         existingIncidentActionPlan.Title,
				"ActionPlanDescription":   existingIncidentActionPlan.Description,
				"RedirectURL":             fmt.Sprintf("%sincident/action-plan?action_plan_id=%d", appURL, id),
			}

			oldResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              oldResponsibleUserEmailTo,
				Subject:              oldResponsibleUserSubject,
				TemplateRelativePath: oldResponsibleUserTemplateRelativePath,
				EmailData:            oldResponsibleUserEmailData,
			}

			newResponsibleUserEmailTo := []string{newResponsibleUser.Email}
			newResponsibleUserSubject := fmt.Sprintf("Incident Action Plan/%s: Responsibility Assigned", existingIncidentActionPlan.Title)
			newResponsibleUserTemplateRelativePath := "incident/templates/email/incident_action_plan_change_responsible_user_mail.html"
			newResponsibleUserEmailData := map[string]interface{}{
				"FirstNameEn":             newResponsibleUser.FirstNameEn,
				"LastNameEn":              newResponsibleUser.LastNameEn,
				"IncidentReferenceNumber": incident.ReferenceNumber,
				"IncidentTitle":           incident.Title,
				"ActionPlanTitle":         existingIncidentActionPlan.Title,
				"ActionPlanDescription":   existingIncidentActionPlan.Description,
				"TargetDate":              existingIncidentActionPlan.TargetDate.Format("2006-01-02"),
				"AssignedByFirstNameEn":   assignedBy.FirstNameEn,
				"AssignedByLastNameEn":    assignedBy.LastNameEn,
				"AssignedOn":              time.Now().Format("2006-01-02"),
				"RedirectURL":             fmt.Sprintf("%sincident/action-plan?action_plan_id=%d", appURL, id),
			}

			newResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              newResponsibleUserEmailTo,
				Subject:              newResponsibleUserSubject,
				TemplateRelativePath: newResponsibleUserTemplateRelativePath,
				EmailData:            newResponsibleUserEmailData,
			}

			err = iaps.notificationSvc.SendNotification(
				newResponsibleUserNotificationDTO, newResponsibleUserEmailNotificationDTO, "incident_action_plan",
			)
			if err != nil {
				utility.LogError("Failed to send new responsible user notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}

			err = iaps.notificationSvc.SendNotification(
				oldResponsibleUserNotificationDTO, oldResponsibleUserEmailNotificationDTO, "incident_action_plan",
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

func (iaps *IncidentActionPlanService) Delete(incidentID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := iaps.repo.BeginTransaction()

	if _, err := iaps.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return nil, err
	}

	var failedIds []uint
	var successfulIds []uint
	var deletedIncidentActionPlan []map[string]interface{}
	var notifications []struct {
		OwnerID         *uint
		TitleEn         string
		TitleAr         string
		Description     string
		ReferenceNumber string
		IncidentTitle   string
	}
	for _, id := range ids {

		incidentActionPlan, err := iaps.repo.GetFullDetailsByID(id)
		if err != nil {
			continue
		}

		deletedIncidentActionPlan = append(deletedIncidentActionPlan, map[string]interface{}{
			"id":    id,
			"title": incidentActionPlan.Title,
		})

		// Prepare notification data before deletion
		notifications = append(notifications, struct {
			OwnerID         *uint
			TitleEn         string
			TitleAr         string
			Description     string
			ReferenceNumber string
			IncidentTitle   string
		}{
			OwnerID:         &incidentActionPlan.ResponsibleUserID,
			TitleEn:         incidentActionPlan.Title,
			TitleAr:         incidentActionPlan.Title,
			Description:     *incidentActionPlan.Description,
			ReferenceNumber: incidentActionPlan.IncidentRegister.ReferenceNumber,
			IncidentTitle:   incidentActionPlan.IncidentRegister.Title,
		})

		existingTask, err := iaps.taskSvc.GetByFieldID("incident_action_plan_id", id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		taskID, err := iaps.taskSvc.Delete(existingTask.ID, currentUserID)
		if err != nil {
			failedIds = append(failedIds, taskID)
		}

		if len(incidentActionPlan.ActionPlanStatusUpdates) > 0 {
			var deleteStatusUpdateIds []*uint
			for _, update := range incidentActionPlan.ActionPlanStatusUpdates {
				deleteStatusUpdateIds = append(deleteStatusUpdateIds, &update.ID)
			}
			if err := iaps.repo.DeleteIncidentActionPlanStatusUpdates(deleteStatusUpdateIds); err != nil {
				failedIds = append(failedIds, id)
				continue
			}
		}
		// Delete-> update timeline performedAt as null
		// todo
		actinPlanCount, err := iaps.repo.GetCountByRiskID(incidentID, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		if actinPlanCount == 1 {
			timelineType, err := iaps.incidentTimelineSvc.GetByType("action-plans")
			if err != nil {
				tx.Rollback()
				return failedIds, err
			}
			timelineStatus, err := iaps.incidentTimelineSvc.GetByStatusType("pending")
			if err != nil {
				tx.Rollback()
				return failedIds, err
			}
			timelineFormDTO := &incidentDTOs.IncidentTimelineFormDTO{
				PerformedByID:    nil,
				TimelineStatusID: timelineStatus.ID,
				PerformedAt:      nil,
			}
			if err := iaps.incidentTimelineSvc.Update(incidentID, timelineFormDTO, timelineType.ID, tx); err != nil {
				tx.Rollback()
				return failedIds, err
			}
		}
		err = iaps.repo.Delete(id, tx)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}
	tx.Commit()

	if len(successfulIds) > 0 {
		go func() {
			err := iaps.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.IncidentRegister{}.TableName(),
				successfulIds,
				deletedIncidentActionPlan,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}

			for _, notification := range notifications {
				urlParam := ""
				entity := "incident-action-plan"
				action := "delete"

				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserID:               *notification.OwnerID,
					ModuleIdentifier:     "incident",
					NotificationType:     string(commonDTOs.DeletionType),
					HeadingEn:            "Incident Action Plan Deleted",
					HeadingAr:            "تم حذف خطة عمل الحادث",
					TitleEn:              notification.TitleEn,
					TitleAr:              notification.TitleAr,
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				ownerUser, err := iaps.userSvc.GetByID(*notification.OwnerID)
				if err != nil {
					utility.LogError("Failed to fetch owner details for notification", map[string]interface{}{
						"error": err.Error(),
					})
					continue
				}

				ownerEmailTo := []string{ownerUser.Email}
				subject := fmt.Sprintf("incident Action Plan/%s: Deleted", notification.TitleEn)
				templateRelativePath := "incident/templates/email/incident_action_plan_delete_mail.html"

				emailData := map[string]interface{}{
					"FirstNameEn":             ownerUser.FirstNameEn,
					"LastNameEn":              ownerUser.LastNameEn,
					"IncidentReferenceNumber": notification.ReferenceNumber,
					"IncidentTitle":           notification.IncidentTitle,
					"ActionPlanTitle":         notification.TitleEn,
					"ActionPlanDescription":   notification.Description,
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              ownerEmailTo,
					Subject:              subject,
					TemplateRelativePath: templateRelativePath,
					EmailData:            emailData,
				}

				err = iaps.notificationSvc.SendNotification(
					dbNotificationDTO, emailNotificationDTO, "incident_action_plan")
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

func (iaps *IncidentActionPlanService) StatusUpdate(incidentID uint, id uint, currentUserID uint, incidentActionPlanStatusUpdateFormDTO *incidentDTOs.IncidentActionPlanStatusUpdateFormDTO, tx utility.DatabaseTransaction) error {
	existingIncidentActionPlan, err := iaps.GetByID(incidentID, id)
	if err != nil {
		return err
	}

	existingIncidentActionPlan.ActionPlanStatusID = incidentActionPlanStatusUpdateFormDTO.ActionPlanStatusID
	existingIncidentActionPlan.Percentage = incidentActionPlanStatusUpdateFormDTO.Percentage
	existingIncidentActionPlan.UpdatedBYID = &currentUserID

	err = iaps.repo.Update(existingIncidentActionPlan, tx)
	if err != nil {
		return err
	}

	isReviewUserNull := existingIncidentActionPlan.ReviewUserID == nil

	// Fetch the corresponding task and update its status
	existingTask, err := iaps.taskSvc.GetByFieldID("incident_action_plan_id", id)
	if err != nil {
		return err
	}
	err = iaps.taskSvc.UpdateStatus(currentUserID, existingTask.ID, incidentActionPlanStatusUpdateFormDTO.ActionPlanStatusID, isReviewUserNull, tx)
	if err != nil {
		return err
	}

	return nil
}

func (iaps *IncidentActionPlanService) validateIncidentActionPlanFormDTO(IncidentActionPlanFormDTO *incidentDTOs.IncidentActionPlanFormDTO) error {
	if IncidentActionPlanFormDTO.SeverityLevelID != nil {
		if _, err := iaps.severityLevelSvc.GetByID(*IncidentActionPlanFormDTO.SeverityLevelID); err != nil {
			return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
		}
	}

	return nil
}

func (iaps *IncidentActionPlanService) parseIncidentActionPlanDates(dto *incidentDTOs.IncidentActionPlanFormDTO) (*incidentDTOs.IncidentActionPlanParsedDates, error) {
	parsed := &incidentDTOs.IncidentActionPlanParsedDates{}

	var err error
	if *dto.TargetDate != "" {
		parsed.TargetDate, err = utility.ParseDate(*dto.TargetDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid target date"), "invalid_date", "target date")
		}
	}

	return parsed, nil
}

func (iaps *IncidentActionPlanService) ExportIncidentActionPlan(data []*incidentDTOs.IncidentActionPlanListDTO, currentUserID uint) (*os.File, error) {
	// Fetch user setting for language preference
	userSetting, err := iaps.userSettingSvc.GetUserSettingsByUserID(currentUserID)
	if err != nil {
		return nil, err
	}

	// Define headers based on user language setting
	var headers []string
	if userSetting.Language == "ar" {
		headers = []string{
			"SL No.",
			"Action Plan",
			"Incident",
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
			"Incident",
			"Description",
			"Severity",
			"Responsible User",
			"Target Date",
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
			actionPlan.IncidentTitle,
			utility.GetStringValue(actionPlan.Description),
			severity,
			responsibleUser,
			utility.GetStringValue(actionPlan.TargetDate),
			status,
		}
	}

	// Generate the Excel file
	return utility.GenerateExcelFile(headers, rows)
}
