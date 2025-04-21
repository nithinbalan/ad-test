package services

import (
	"errors"
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	strategyDTOs "jethur/grc/strategy/dtos"
	StrategyProfileActionPlanRepositories "jethur/grc/strategy/repositories"
	taskServices "jethur/grc/task/services"
	"jethur/utility"
	"time"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type StrategyProfileActionPlanServiceInterface interface {
	GetAll(strategyProfileId uint, params *strategyDTOs.StrategyProfileActionPlanListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileActionPlanListDTO, error)
	GetByID(strategyProfileID uint, id uint) (*models.StrategyProfileActionPlan, error)
	GetFullDetailsByID(strategyProfileID uint, id uint, tx ...utility.DatabaseTransaction) (*strategyDTOs.StrategyProfileActionPlanDetailDTO, error)
	Create(strategyProfileID uint, currentUserID uint, dto *strategyDTOs.StrategyProfileActionPlanFormDTO) (uint, error)
	Update(strategyProfileID uint, id uint, currentUserID uint, dto *strategyDTOs.StrategyProfileActionPlanFormDTO) error
	Delete(strategyProfileID uint, ids []uint, currentUserID uint) ([]uint, error)
	StatusUpdate(strategyProfileId uint, id uint, currentUserID uint, strategyProfileActionPlanStatusUpdateFormDTO *strategyDTOs.StrategyProfileActionPlanStatusUpdateFormDTO, tx utility.DatabaseTransaction) error
}

type StrategyProfileActionPlanService struct {
	repo                StrategyProfileActionPlanRepositories.StrategyProfileActionPlanRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	userSvc             organizationServices.UserServiceInterface
	strategyProfileSvc  StrategyProfileServiceInterface
	ActionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface
	taskSvc             taskServices.TaskServiceInterface
	notificationSvc     commonServices.NotificationServiceInterface
	userSettingSvc      organizationServices.UserSettingServiceInterface
}

func NewStrategyProfileActionPlanService(
	repo StrategyProfileActionPlanRepositories.StrategyProfileActionPlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	strategyProfileSvc StrategyProfileServiceInterface,
	ActionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	taskSvc taskServices.TaskServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
) StrategyProfileActionPlanServiceInterface {
	return &StrategyProfileActionPlanService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		userSvc:             userSvc,
		strategyProfileSvc:  strategyProfileSvc,
		ActionPlanStatusSvc: ActionPlanStatusSvc,
		taskSvc:             taskSvc,
		notificationSvc:     notificationSvc,
		userSettingSvc:      userSettingSvc,
	}
}

func (spaps *StrategyProfileActionPlanService) GetAll(strategyProfileId uint, params *strategyDTOs.StrategyProfileActionPlanListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileActionPlanListDTO, error) {
	if _, err := spaps.strategyProfileSvc.GetByID(strategyProfileId); err != nil {
		return nil, err
	}
	return spaps.repo.GetAll(strategyProfileId, params)
}

func (spaps *StrategyProfileActionPlanService) GetByID(strategyProfileID uint, id uint) (*models.StrategyProfileActionPlan, error) {
	actionPlan, err := spaps.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return actionPlan, nil
}

func (spaps *StrategyProfileActionPlanService) GetFullDetailsByID(strategyProfileID uint, id uint, tx ...utility.DatabaseTransaction) (*strategyDTOs.StrategyProfileActionPlanDetailDTO, error) {
	if _, err := spaps.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return nil, err
	}
	var strategyProfileActionPlan strategyDTOs.StrategyProfileActionPlanDetailDTO

	actionPlan, err := spaps.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(actionPlan, &strategyProfileActionPlan)
	if err != nil {
		return nil, err
	}

	if actionPlan.TargetDate != nil {
		targetDate := actionPlan.TargetDate.Format("2006-01-02")
		strategyProfileActionPlan.TargetDateString = targetDate
	}

	strategyProfileActionPlan.CreatedAt = actionPlan.CreatedAt
	strategyProfileActionPlan.UpdatedAt = actionPlan.UpdatedAt

	return &strategyProfileActionPlan, nil
}

func (spaps *StrategyProfileActionPlanService) Create(strategyProfileID uint, currentUserID uint, strategyProfileActionPlanFormDTO *strategyDTOs.StrategyProfileActionPlanFormDTO) (uint, error) {
	tx := spaps.repo.BeginTransaction()

	strategyProfile, err := spaps.strategyProfileSvc.GetFullDetailsByID(strategyProfileID)
	if err != nil {
		return 0, err
	}

	if _, err := spaps.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return 0, err
	}

	responsibleUnitID, err := spaps.userSvc.GetUserPrimaryUnitID(*strategyProfileActionPlanFormDTO.ResponsibleUserID)
	if err != nil {
		return 0, utility.NewValidationError(errors.New("organization unit not set for responsible user"), "not_exist", "responsible unit")
	}

	actionPlanStatus, err := spaps.ActionPlanStatusSvc.GetByType("new")
	if err != nil {
		return 0, err
	}

	dates, err := spaps.parseStrategyProfileActionPlanDates(strategyProfileActionPlanFormDTO)
	if err != nil {
		return 0, err
	}

	actionPlan := &models.StrategyProfileActionPlan{
		StrategyProfileID:           strategyProfileID,
		StrategyProfileFocusAreaID:  strategyProfileActionPlanFormDTO.StrategyProfileFocusAreaID,
		StrategyProfileGoalID:       strategyProfileActionPlanFormDTO.StrategyProfileGoalID,
		StrategyProfileObjectiveID:  strategyProfileActionPlanFormDTO.StrategyProfileObjectiveID,
		StrategyProfileInitiativeID: strategyProfileActionPlanFormDTO.StrategyProfileInitiativeID,
		StrategyProfileProjectID:    strategyProfileActionPlanFormDTO.StrategyProfileProjectID,
		StrategyProfileMilestoneID:  strategyProfileActionPlanFormDTO.StrategyProfileMilestoneID,
		SeverityLevelID:             strategyProfileActionPlanFormDTO.SeverityLevelID,
		Title:                       strategyProfileActionPlanFormDTO.Title,
		Description:                 strategyProfileActionPlanFormDTO.Description,
		ResponsibleUserID:           strategyProfileActionPlanFormDTO.ResponsibleUserID,
		ReviewUserID:                strategyProfileActionPlanFormDTO.ReviewUserID,
		ResponsibleUnitID:           &responsibleUnitID,
		TargetDate:                  dates.TargetDate,
		ActionPlanStatusId:          &actionPlanStatus.ID,
		CreatedBYID:                 currentUserID,
	}

	id, err := spaps.repo.Create(actionPlan, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	actionPlanStatusID := actionPlanStatus.ID
	strategyProfileActionPlanUpdate := &models.StrategyProfileActionPlanStatusUpdate{
		StrategyProfileActionPlanID: id,
		ActionPlanStatusID:          actionPlanStatusID,
		Comment:                     "Initial action plan status",
		CreatedByID:                 currentUserID,
	}

	_, err = spaps.repo.CreateInitialStrategyProfileActionPlanStatusUpdate(strategyProfileActionPlanUpdate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	urlParam := fmt.Sprintf(`{"strategy_profile_id":%d, "action_plan_id":%d}`, strategyProfileID, id)
	taskDTO := commonDTOs.TaskFormDTO{
		ModuleIdentifier:            "strategy",
		StrategyProfileActionPlanID: &id,
		Title:                       strategyProfileActionPlanFormDTO.Title,
		TargetDate:                  dates.TargetDate,
		URLParam:                    &urlParam,
		QueryParam:                  nil,
		SeverityLevelID:             strategyProfileActionPlanFormDTO.SeverityLevelID,
		ResponsibleUserID:           *strategyProfileActionPlanFormDTO.ResponsibleUserID,
	}

	if _, err := spaps.taskSvc.Create(taskDTO, currentUserID, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	go func() {
		existingStrategyProfileActionPlan, err := spaps.GetFullDetailsByID(strategyProfileID, actionPlan.ID)
		if err != nil {
			utility.LogError("Failed to fetch test and exercise", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = spaps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.StrategyProfileActionPlan{}.TableName(),
			id,
			existingStrategyProfileActionPlan,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.StrategyProfileActionPlan{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
		// prepare notification
		appURL := utility.GetAppURL()
		urlParam := fmt.Sprintf(`{"strategy_profile_id":%d, "action_plan_id":%d}`, strategyProfileID, id)
		redirectUrl := fmt.Sprintf("%sstrategy/strategy-profile/%d/action-plans?action_plan_id=%d", appURL, strategyProfileID, id)
		queryParam := ""
		entity := "strategy-profile-action-plan"
		action := "responsible-user-assign"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{*strategyProfileActionPlanFormDTO.ResponsibleUserID},
			ModuleIdentifier:     "strategy",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "Action Plan Responsibility Assigned",
			HeadingAr:            "مسؤولية خطة العمل الموكلة",
			TitleEn:              strategyProfileActionPlanFormDTO.Title,
			TitleAr:              strategyProfileActionPlanFormDTO.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			QueryParam:           queryParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		responsibleUser, err := spaps.userSvc.GetByID(*strategyProfileActionPlanFormDTO.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch responsible user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		assignedBy, err := spaps.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch assigning user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("Action Plan/%s: Responsibility Assigned", strategyProfileActionPlanFormDTO.Title)
		emailData := map[string]interface{}{
			"FirstNameEn":           responsibleUser.FirstNameEn,
			"LastNameEn":            responsibleUser.LastNameEn,
			"StrategyProfileTitle":  strategyProfile.Title,
			"ActionPlanTitle":       strategyProfileActionPlanFormDTO.Title,
			"TargetDate":            dates.TargetDate.Format("2006-01-02"),
			"AssignedByFirstNameEn": assignedBy.FirstNameEn,
			"AssignedByLastNameEn":  assignedBy.LastNameEn,
			"AssignedOn":            time.Now().Format("2006-01-02"),
			"RedirectURL":           redirectUrl,
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "strategy/templates/email/action_plan_create_responsible_user_mail.html",
			EmailData:            emailData,
		}

		// send notification and mail
		err = spaps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "strategy_profile_action_plan")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()
	return id, nil
}

func (spaps *StrategyProfileActionPlanService) Update(strategyProfileID uint, id uint, currentUserID uint, strategyProfileActionPlanFormDTO *strategyDTOs.StrategyProfileActionPlanFormDTO) error {
	tx := spaps.repo.BeginTransaction()

	strategyProfile, err := spaps.strategyProfileSvc.GetFullDetailsByID(strategyProfileID)
	if err != nil {
		return err
	}

	existingStrategyProfilePlan, err := spaps.repo.GetByID(id)
	if err != nil {
		return utility.NewResourceNotFoundError(err, id)
	}

	oldStrategyProfileActionPlan, err := spaps.GetFullDetailsByID(strategyProfileID, id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	oldResponsibleUserID := existingStrategyProfilePlan.ResponsibleUserID

	if existingStrategyProfilePlan.ResponsibleUserID != strategyProfileActionPlanFormDTO.ResponsibleUserID {
		responsibleUnitID, err := spaps.userSvc.GetUserPrimaryUnitID(*existingStrategyProfilePlan.ResponsibleUserID)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "risk responsible unit")
		}
		existingStrategyProfilePlan.ResponsibleUnitID = &responsibleUnitID
	}

	dates, err := spaps.parseStrategyProfileActionPlanDates(strategyProfileActionPlanFormDTO)
	if err != nil {
		return err
	}

	existingStrategyProfilePlan.SeverityLevelID = strategyProfileActionPlanFormDTO.SeverityLevelID
	existingStrategyProfilePlan.Title = strategyProfileActionPlanFormDTO.Title
	existingStrategyProfilePlan.Description = strategyProfileActionPlanFormDTO.Description
	existingStrategyProfilePlan.ResponsibleUserID = strategyProfileActionPlanFormDTO.ResponsibleUserID
	existingStrategyProfilePlan.TargetDate = dates.TargetDate
	existingStrategyProfilePlan.ReviewUserID = strategyProfileActionPlanFormDTO.ReviewUserID
	existingStrategyProfilePlan.UpdatedBYID = &currentUserID

	err = spaps.repo.Update(existingStrategyProfilePlan, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// //update task
	urlParam := fmt.Sprintf(`{"strategy_profile_id":%d, "action_plan_id":%d}`, strategyProfileID, id)
	taskUpdateDTO := commonDTOs.TaskFormDTO{
		Title:             strategyProfileActionPlanFormDTO.Title,
		TargetDate:        dates.TargetDate,
		URLParam:          &urlParam,
		SeverityLevelID:   strategyProfileActionPlanFormDTO.SeverityLevelID,
		ResponsibleUserID: *strategyProfileActionPlanFormDTO.ResponsibleUserID,
	}

	err = spaps.taskSvc.Update(currentUserID, taskUpdateDTO, "strategy_profile_action_plan_id", id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	go func() {
		existingActionPlan, err := spaps.GetFullDetailsByID(strategyProfileID, id)
		if err != nil {
			return
		}

		changedData := make(map[string]interface{})

		// log changes
		oldData := map[string]interface{}{
			"responsible_user": fmt.Sprintf("%s %s", oldStrategyProfileActionPlan.ResponsibleUser.FirstNameEN, oldStrategyProfileActionPlan.ResponsibleUser.LastNameEN),
			"severity_level":   oldStrategyProfileActionPlan.SeverityLevel.TitleEn,
			"target_date":      oldStrategyProfileActionPlan.TargetDateString,
		}

		newData := map[string]interface{}{
			"responsible_user": fmt.Sprintf("%s %s", existingActionPlan.ResponsibleUser.FirstNameEN, existingActionPlan.ResponsibleUser.LastNameEN),
			"severity_level":   existingActionPlan.SeverityLevel.TitleEn,
			"target_date":      existingActionPlan.TargetDateString,
		}

		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}

		if len(changedData) > 0 {
			err := spaps.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.StrategyProfileActionPlan{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.StrategyProfileActionPlan{}.TableName(),
					"recordID":  id,
					"userID":    currentUserID,
					"error":     err.Error(),
				}

				utility.LogError("Failed to log activity", fields)
			}
		}

		appURL := utility.GetAppURL()
		urlParam := fmt.Sprintf(`{"strategy_profile_id":%d, "action_plan_id":%d}`, strategyProfileID, id)
		redirectUrl := fmt.Sprintf("%sstrategy/strategy-profile/%d/action-plans?action_plan_id=%d", appURL, strategyProfileID, id)

		// sends notification if responsible user has changed
		if oldResponsibleUserID != existingActionPlan.ResponsibleUserID {
			queryParam := ""
			entity := "strategy-profile-action-plan"

			oldResponsibleUserAction := "responsible-user-remove"
			oldResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{*oldResponsibleUserID},
				ModuleIdentifier:     "strategy",
				NotificationType:     string(commonDTOs.InformationType),
				HeadingEn:            "Action Plan Responsibility Removed",
				HeadingAr:            "تمت إزالة مسؤولية خطة العمل",
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
				UserIDs:              []uint{*strategyProfileActionPlanFormDTO.ResponsibleUserID},
				ModuleIdentifier:     "strategy",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Action Plan Responsibility Assigned",
				HeadingAr:            "مسؤولية خطة العمل الموكلة",
				TitleEn:              existingActionPlan.Title,
				TitleAr:              existingActionPlan.Title,
				Entity:               &entity,
				Action:               &newResponsibleUserAction,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			oldResponsibleUser, err := spaps.userSvc.GetByID(*oldResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch old responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			newResponsibleUser, err := spaps.userSvc.GetByID(*strategyProfileActionPlanFormDTO.ResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch new responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			assignedBy, err := spaps.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch assigning user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			oldResponsibleUserEmailTo := []string{oldResponsibleUser.Email}
			oldResponsibleUserSubject := fmt.Sprintf("Action Plan/%s: Responsibility Removed", existingActionPlan.Title)
			oldResponsibleUserTemplateRelativePath := "strategy/templates/email/action_plan_remove_responsible_user_mail.html"
			oldResponsibleUserEmailData := map[string]interface{}{
				"FirstNameEn":          oldResponsibleUser.FirstNameEn,
				"LastNameEn":           oldResponsibleUser.LastNameEn,
				"StrategyProfileTitle": strategyProfile.Title,
				"ActionPlanTitle":      existingActionPlan.Title,
				"RedirectURL":          redirectUrl,
			}

			oldResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              oldResponsibleUserEmailTo,
				Subject:              oldResponsibleUserSubject,
				TemplateRelativePath: oldResponsibleUserTemplateRelativePath,
				EmailData:            oldResponsibleUserEmailData,
			}

			newResponsibleUserEmailTo := []string{newResponsibleUser.Email}
			newResponsibleUserSubject := fmt.Sprintf("Action Plan/%s: Responsibility Assigned", existingActionPlan.Title)
			newResponsibleUserTemplateRelativePath := "strategy/templates/email/action_plan_change_responsible_user_mail.html"
			newResponsibleUserEmailData := map[string]interface{}{
				"FirstNameEn":           newResponsibleUser.FirstNameEn,
				"LastNameEn":            newResponsibleUser.LastNameEn,
				"StrategyProfileTitle":  strategyProfile.Title,
				"ActionPlanTitle":       existingStrategyProfilePlan.Title,
				"TargetDate":            existingStrategyProfilePlan.TargetDate.Format("2006-01-02"),
				"AssignedByFirstNameEn": assignedBy.FirstNameEn,
				"AssignedByLastNameEn":  assignedBy.LastNameEn,
				"AssignedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           redirectUrl,
			}

			newResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              newResponsibleUserEmailTo,
				Subject:              newResponsibleUserSubject,
				TemplateRelativePath: newResponsibleUserTemplateRelativePath,
				EmailData:            newResponsibleUserEmailData,
			}

			err = spaps.notificationSvc.SendNotification(
				newResponsibleUserNotificationDTO, newResponsibleUserEmailNotificationDTO, "strategy_profile_action_plan",
			)
			if err != nil {
				utility.LogError("Failed to send new responsible user notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}

			err = spaps.notificationSvc.SendNotification(
				oldResponsibleUserNotificationDTO, oldResponsibleUserEmailNotificationDTO, "strategy_profile_action_plan",
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

func (spaps *StrategyProfileActionPlanService) Delete(strategyProfileID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := spaps.repo.BeginTransaction()

	if _, err := spaps.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return nil, err
	}

	strategyProfile, err := spaps.strategyProfileSvc.GetFullDetailsByID(strategyProfileID)
	if err != nil {
		return nil, err
	}

	var failedIDs []uint
	var successfulIDs []uint
	var deletedActionPlans []map[string]interface{}
	var notifications []struct {
		OwnerID              *uint
		TitleEn              string
		TitleAr              string
		Description          string
		ReferenceNumber      string
		StrategyProfileTitle string
	}

	for _, id := range ids {
		actionPlan, err := spaps.GetFullDetailsByID(strategyProfileID, id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		// Prepare notification data before deletion
		notifications = append(notifications, struct {
			OwnerID              *uint
			TitleEn              string
			TitleAr              string
			Description          string
			ReferenceNumber      string
			StrategyProfileTitle string
		}{
			OwnerID:              actionPlan.ResponsibleUserID,
			TitleEn:              actionPlan.Title,
			TitleAr:              actionPlan.Title,
			StrategyProfileTitle: strategyProfile.Title,
		})

		deletedActionPlans = append(deletedActionPlans, map[string]interface{}{
			"id":                           id,
			"strategy_profile_action_plan": actionPlan.ID,
		})

		if len(actionPlan.ActionPlanStatusUpdates) > 0 {
			var deleteStatusUpdateIds []*uint
			var deleteDocumentIds []*uint
			for _, update := range actionPlan.ActionPlanStatusUpdates {
				deleteStatusUpdateIds = append(deleteStatusUpdateIds, &update.ID)
				for _, document := range update.Documents {
					deleteDocumentIds = append(deleteDocumentIds, &document.ID)
				}
			}
			if err := spaps.repo.DeleteActionPlanStatusUpdateDocuments(deleteDocumentIds); err != nil {
				failedIDs = append(failedIDs, id)
				continue
			}
			if err := spaps.repo.DeleteActionPlanStatusUpdates(deleteStatusUpdateIds); err != nil {
				failedIDs = append(failedIDs, id)
				continue
			}
		}

		if len(actionPlan.ActionPlanComments) > 0 {
			var deleteCommentIds []*uint
			var deleteCommentDocumentIds []*uint
			for _, update := range actionPlan.ActionPlanComments {
				deleteCommentIds = append(deleteCommentIds, &update.ID)
				for _, document := range update.Documents {
					deleteCommentDocumentIds = append(deleteCommentDocumentIds, &document.ID)
				}
			}
			if err := spaps.repo.DeleteActionPlanCommentDocuments(deleteCommentDocumentIds); err != nil {
				failedIDs = append(failedIDs, id)
				continue
			}
			if err := spaps.repo.DeleteActionPlanComments(deleteCommentIds); err != nil {
				failedIDs = append(failedIDs, id)
				continue
			}
		}

		existingTask, err := spaps.taskSvc.GetByFieldID("strategy_profile_action_plan_id", id)
		if err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		taskID, err := spaps.taskSvc.Delete(existingTask.ID, currentUserID)
		if err != nil {
			failedIDs = append(failedIDs, taskID)
		}

		if err := spaps.repo.Delete(id, tx); err != nil {
			failedIDs = append(failedIDs, id)
		} else {
			successfulIDs = append(successfulIDs, id)
		}
	}

	if len(failedIDs) > 0 {
		tx.Rollback()
		return failedIDs, nil
	}

	tx.Commit()

	if len(successfulIDs) > 0 {
		go func() {
			err := spaps.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.StrategyProfileActionPlan{}.TableName(),
				successfulIDs,
				deletedActionPlans,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
		for _, notification := range notifications {
			urlParam := ""
			entity := "strategy-profile-action-plan"
			action := "delete"

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{*notification.OwnerID},
				ModuleIdentifier:     "strategy",
				NotificationType:     string(commonDTOs.DeletionType),
				HeadingEn:            "Action Plan Deleted",
				HeadingAr:            "تم حذف خطة العمل",
				TitleEn:              notification.TitleEn,
				TitleAr:              notification.TitleAr,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			ownerUser, err := spaps.userSvc.GetByID(*notification.OwnerID)
			if err != nil {
				utility.LogError("Failed to fetch owner details for notification", map[string]interface{}{
					"error": err.Error(),
				})
				continue
			}

			ownerEmailTo := []string{ownerUser.Email}
			subject := fmt.Sprintf("Action Plan/%s: Deleted", notification.TitleEn)
			templateRelativePath := "strategy/templates/email/action_plan_delete_mail.html"

			emailData := map[string]interface{}{
				"FirstNameEn":           ownerUser.FirstNameEn,
				"LastNameEn":            ownerUser.LastNameEn,
				"ReferenceNumber":       notification.ReferenceNumber,
				"StrategyProfileTitle":  notification.StrategyProfileTitle,
				"ActionPlanTitle":       notification.TitleEn,
				"ActionPlanDescription": notification.Description,
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              ownerEmailTo,
				Subject:              subject,
				TemplateRelativePath: templateRelativePath,
				EmailData:            emailData,
			}

			err = spaps.notificationSvc.SendNotification(
				dbNotificationDTO, emailNotificationDTO, "strategy_profile_action_plan")
			if err != nil {
				utility.LogError("Failed to send deletion notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}

		}
	}

	return failedIDs, nil
}

func (spaps *StrategyProfileActionPlanService) parseStrategyProfileActionPlanDates(dto *strategyDTOs.StrategyProfileActionPlanFormDTO) (*strategyDTOs.StrategyProfileActionPlanParsedDates, error) {
	parsed := &strategyDTOs.StrategyProfileActionPlanParsedDates{}

	var err error

	if &dto.TargetDate != nil {
		parsed.TargetDate, err = utility.ParseDate(dto.TargetDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid target date"), "invalid_date", "target date")
		}
	}

	return parsed, nil
}

func (spaps *StrategyProfileActionPlanService) StatusUpdate(slaID uint, id uint, currentUserID uint, strategyProfileActionPlanStatusUpdateFormDTO *strategyDTOs.StrategyProfileActionPlanStatusUpdateFormDTO, tx utility.DatabaseTransaction) error {
	existingStrategyProfileActionPlan, err := spaps.GetByID(slaID, id)
	if err != nil {
		return err
	}

	existingStrategyProfileActionPlan.ActionPlanStatusId = &strategyProfileActionPlanStatusUpdateFormDTO.ActionPlanStatusID
	existingStrategyProfileActionPlan.Percentage = strategyProfileActionPlanStatusUpdateFormDTO.Percentage
	existingStrategyProfileActionPlan.UpdatedBYID = &currentUserID

	err = spaps.repo.Update(existingStrategyProfileActionPlan, tx)
	if err != nil {
		return err
	}
	isReviewUserNull := existingStrategyProfileActionPlan.ReviewUserID == nil
	// Fetch the corresponding task and update its status
	existingTask, err := spaps.taskSvc.GetByFieldID("strategy_profile_action_plan_id", id)
	if err != nil {
		return err
	}
	err = spaps.taskSvc.UpdateStatus(currentUserID, existingTask.ID, strategyProfileActionPlanStatusUpdateFormDTO.ActionPlanStatusID, isReviewUserNull, tx)
	if err != nil {
		return err
	}

	return nil
}
