package services

import (
	"errors"
	"fmt"
	"time"

	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	masterRepositories "jethur/grc/masters/repositories"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	organizationalContextDTOs "jethur/grc/organizational-context/dtos"
	stakeholderRepositories "jethur/grc/organizational-context/repositories"
	settingService "jethur/grc/settings/services"
	taskServices "jethur/grc/task/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type StakeholderActionPlanServiceInterface interface {
	GetAll(stakeholderID uint, params *organizationalContextDTOs.StakeholderActionPlanListQueryParamsDTO) (*organizationalContextDTOs.StakeholderActionPlanListResponseDTO, error)
	GetByID(stakeholderID uint, id uint) (*models.StakeholderActionPlan, error)
	GetFullDetailsByID(stakeholderID uint, id uint, tx ...utility.DatabaseTransaction) (*organizationalContextDTOs.StakeholderActionPlanDetailDTO, error)
	Create(stakeholderID uint, currentUserID uint, StakeholderActionPlanFormDTO *organizationalContextDTOs.StakeholderActionPlanFormDTO) (uint, error)
	Update(stakeholderID uint, id uint, currentUserID uint, stakeholder *organizationalContextDTOs.StakeholderActionPlanFormDTO) error
	Delete(stakeholderID uint, ids []uint, currentUserID uint) ([]uint, error)
	StatusUpdate(stakeholderID uint, id uint, currentUserID uint, stakeholderActionPlanStatusUpdateFormDTO *organizationalContextDTOs.StakeholderActionPlanStatusUpdateFormDTO, tx utility.DatabaseTransaction) error
}

// Service and implementations
type StakeholderActionPlanService struct {
	repo                stakeholderRepositories.StakeholderActionPlanRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	stakeholderSvc      StakeholderServiceInterface
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface
	severityLevelSvc    masterServices.SeverityLevelServiceInterface
	taskSvc             taskServices.TaskServiceInterface
	moduleSvc           settingService.ModuleServiceInterface
	userSvc             organizationServices.UserServiceInterface
	notificationSvc     commonServices.NotificationServiceInterface
	userSettingSvc      organizationServices.UserSettingServiceInterface
	organizationUnitSvc settingService.OrganizationUnitServiceInterface
}

func NewStakeholderActionPlanService(
	repo stakeholderRepositories.StakeholderActionPlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	stakeholderSvc StakeholderServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	severityLevelSvc masterRepositories.SeverityLevelRepositoryInterface,
	taskSvc taskServices.TaskServiceInterface,
	moduleSvc settingService.ModuleServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
	organizationUnitSvc settingService.OrganizationUnitServiceInterface,
) StakeholderActionPlanServiceInterface {
	return &StakeholderActionPlanService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		stakeholderSvc:      stakeholderSvc,
		actionPlanStatusSvc: actionPlanStatusSvc,
		severityLevelSvc:    severityLevelSvc,
		taskSvc:             taskSvc,
		moduleSvc:           moduleSvc,
		userSvc:             userSvc,
		notificationSvc:     notificationSvc,
		userSettingSvc:      userSettingSvc,
		organizationUnitSvc: organizationUnitSvc,
	}
}

func (caps *StakeholderActionPlanService) GetAll(stakeholderID uint, params *organizationalContextDTOs.StakeholderActionPlanListQueryParamsDTO) (*organizationalContextDTOs.StakeholderActionPlanListResponseDTO, error) {
	if stakeholderID != 0 {
		if _, err := caps.stakeholderSvc.GetByID(stakeholderID); err != nil {
			return nil, err
		}
	}
	return caps.repo.GetAll(stakeholderID, params, uint(0))
}

func (caps StakeholderActionPlanService) GetByID(stakeholderID uint, id uint) (*models.StakeholderActionPlan, error) {
	if _, err := caps.stakeholderSvc.GetByID(stakeholderID); err != nil {
		return nil, err
	}

	stakeholderActionPlan, err := caps.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return stakeholderActionPlan, nil
}

func (caps StakeholderActionPlanService) GetFullDetailsByID(stakeholderID uint, id uint, tx ...utility.DatabaseTransaction) (*organizationalContextDTOs.StakeholderActionPlanDetailDTO, error) {
	var stakeholderActionPlanDTO organizationalContextDTOs.StakeholderActionPlanDetailDTO

	// Check if the stakeholder register exists
	if _, err := caps.stakeholderSvc.GetByID(stakeholderID); err != nil {
		return nil, err
	}

	// Get full details of the stakeholder action plan
	stakeholderActionPlan, err := caps.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Decode the stakeholder action plan into the DTO
	err = mapstructure.Decode(stakeholderActionPlan, &stakeholderActionPlanDTO)
	if err != nil {
		return nil, err
	}
	if stakeholderActionPlan.TargetDate != nil {
		targetDate := stakeholderActionPlan.TargetDate.Format("2006-01-02")
		stakeholderActionPlanDTO.TargetDateString = &targetDate
	}

	responsibleUnitID, err := caps.userSvc.GetUserPrimaryUnitID(stakeholderActionPlan.ResponsibleUserID)
	if err != nil {
		return nil, err
	}

	responsibleUnit, err := caps.organizationUnitSvc.GetFullDetailsByID(responsibleUnitID)
	if err != nil {
		return nil, err
	}
	stakeholderActionPlanDTO.ResponsibleUnitID = &responsibleUnitID
	stakeholderActionPlanDTO.ResponsibleUnit = responsibleUnit

	// if stakeholderActionPlan.ReviewUserID != nil {
	// 	reviewUserUnitID, err := caps.userSvc.GetUserPrimaryUnitID(*stakeholderActionPlan.ReviewUserID)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	reviewUserUnit, err := caps.organizationUnitSvc.GetFullDetailsByID(reviewUserUnitID)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	stakeholderActionPlanDTO.ReviewUserUnitID = &reviewUserUnitID
	// 	stakeholderActionPlanDTO.ReviewUserUnit = reviewUserUnit
	// }

	stakeholderActionPlanDTO.CreatedAt = stakeholderActionPlan.CreatedAt
	stakeholderActionPlanDTO.UpdatedAt = stakeholderActionPlan.UpdatedAt

	// TODO: Add any additional fields required for the details screen

	return &stakeholderActionPlanDTO, nil
}

func (caps *StakeholderActionPlanService) Create(stakeholderID uint, currentUserID uint, StakeholderActionPlanFormDTO *organizationalContextDTOs.StakeholderActionPlanFormDTO) (uint, error) {
	tx := caps.repo.BeginTransaction()

	stakeholder, err := caps.stakeholderSvc.GetByID(stakeholderID)
	if err != nil {
		return 0, err
	}

	if err := caps.validateStakeholderActionPlanFormDTO(stakeholderID, StakeholderActionPlanFormDTO); err != nil {
		return 0, err
	}

	dates, err := caps.parseStakeholderActionPlanDates(StakeholderActionPlanFormDTO)
	if err != nil {
		return 0, err
	}

	actionPlanStatus, err := caps.actionPlanStatusSvc.GetByType("new")
	if err != nil {
		return 0, err
	}

	var responsibleUnitID uint
	unitID, err := caps.userSvc.GetUserPrimaryUnitID(StakeholderActionPlanFormDTO.ResponsibleUserID)
	if err != nil {
		return 0, utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "responsible unit")
	}
	responsibleUnitID = unitID

	stakeholderActionPlan := &models.StakeholderActionPlan{
		StakeholderID:      &stakeholderID,
		Title:              StakeholderActionPlanFormDTO.Title,
		Description:        StakeholderActionPlanFormDTO.Description,
		SeverityLevelID:    StakeholderActionPlanFormDTO.SeverityLevelID,
		ResponsibleUserID:  StakeholderActionPlanFormDTO.ResponsibleUserID,
		ResponsibleUnitID:  responsibleUnitID,
		ReviewUserID:       StakeholderActionPlanFormDTO.ReviewUserID,
		TargetDate:         dates.TargetDate,
		ActionPlanStatusID: actionPlanStatus.ID,
		CreatedByID:        currentUserID,
	}

	id, err := caps.repo.Create(stakeholderActionPlan, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	actionPlanStatusID := actionPlanStatus.ID
	stakeholderActionPlanUpdate := &models.StakeholderActionPlanStatusUpdate{
		StakeholderActionPlanID: id,
		ActionPlanStatusID:      actionPlanStatusID,
		Comment:                 "Initial action plan status",
		CreatedByID:             currentUserID,
	}

	_, err = caps.repo.CreateInitialStakeholderActionPlanStatusUpdate(stakeholderActionPlanUpdate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	urlParam := fmt.Sprintf(`{"stakeholder_id":%d, "action_plan_id":%d}`, stakeholderID, id)
	taskDTO := commonDTOs.TaskFormDTO{
		ModuleIdentifier:        "stakeholder_analysis",
		StakeholderActionPlanID: &id,
		Title:                   StakeholderActionPlanFormDTO.Title,
		Description:             StakeholderActionPlanFormDTO.Description,
		TargetDate:              dates.TargetDate,
		URLParam:                &urlParam,
		QueryParam:              nil,
		SeverityLevelID:         StakeholderActionPlanFormDTO.SeverityLevelID,
		ResponsibleUserID:       StakeholderActionPlanFormDTO.ResponsibleUserID,
	}

	if _, err := caps.taskSvc.Create(taskDTO, currentUserID, tx); err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()

	go func() {
		existingStakeholderActionPlan, err := caps.GetFullDetailsByID(stakeholderID, id)
		if err != nil {
			utility.LogError("Failed to fetch stakeholder action plan", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = caps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.StakeholderActionPlan{}.TableName(),
			id,
			existingStakeholderActionPlan,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.StakeholderActionPlan{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}

		// prepare notification
		urlParam := fmt.Sprintf(`{"stakeholder_id":%d, "action_plan_id":%d}`, stakeholderID, id)
		queryParam := ""
		entity := "stakeholder-action-plan"
		action := "responsible-user-assign"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{StakeholderActionPlanFormDTO.ResponsibleUserID},
			ModuleIdentifier:     "stakeholder_analysis",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "Stakeholder Action Plan Responsibility Assigned",
			HeadingAr:            "مسؤولية خطة عمل أصحاب المصلحة المعينة",
			TitleEn:              StakeholderActionPlanFormDTO.Title,
			TitleAr:              StakeholderActionPlanFormDTO.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			QueryParam:           queryParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		responsibleUser, err := caps.userSvc.GetByID(StakeholderActionPlanFormDTO.ResponsibleUserID)
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
		subject := fmt.Sprintf("stakeholder Action Plan/%s: Responsibility Assigned", StakeholderActionPlanFormDTO.Title)
		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstNameEn":           responsibleUser.FirstNameEn,
			"LastNameEn":            responsibleUser.LastNameEn,
			"StakeholderTitle":      stakeholder.TitleEn,
			"ActionPlanTitle":       StakeholderActionPlanFormDTO.Title,
			"ActionPlanDescription": StakeholderActionPlanFormDTO.Description,
			"TargetDate":            dates.TargetDate.Format("2006-01-02"),
			"AssignedByFirstNameEn": assignedBy.FirstNameEn,
			"AssignedByLastNameEn":  assignedBy.LastNameEn,
			"AssignedOn":            time.Now().Format("2006-01-02"),
			"RedirectURL":           fmt.Sprintf("%sstakeholder/%d/action-plans?action_plan_id=%d", appURL, stakeholderID, id),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "organizational-context/templates/email/stakeholder_action_plan_create_responsible_user_mail.html",
			EmailData:            emailData,
		}

		// send notification and mail
		err = caps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "stakeholder_action_plan")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	return id, nil
}

func (caps *StakeholderActionPlanService) Update(stakeholderID uint, id uint, currentUserID uint, StakeholderActionPlanFormDTO *organizationalContextDTOs.StakeholderActionPlanFormDTO) error {
	tx := caps.repo.BeginTransaction()

	stakeholder, err := caps.stakeholderSvc.GetByID(stakeholderID)
	if err != nil {
		return err
	}

	existingStakeholderActionPlan, err := caps.GetByID(stakeholderID, id)
	if err != nil {
		return err
	}

	oldStakeholderActionPlan, err := caps.GetFullDetailsByID(stakeholderID, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	oldResponsibleUserID := existingStakeholderActionPlan.ResponsibleUserID

	if err := caps.validateStakeholderActionPlanFormDTO(stakeholderID, StakeholderActionPlanFormDTO); err != nil {
		return err
	}

	dates, err := caps.parseStakeholderActionPlanDates(StakeholderActionPlanFormDTO)
	if err != nil {
		return err
	}

	// Update Stakeholder ActionPlan
	existingStakeholderActionPlan.Title = StakeholderActionPlanFormDTO.Title
	existingStakeholderActionPlan.Description = StakeholderActionPlanFormDTO.Description
	existingStakeholderActionPlan.SeverityLevelID = StakeholderActionPlanFormDTO.SeverityLevelID
	existingStakeholderActionPlan.ResponsibleUserID = StakeholderActionPlanFormDTO.ResponsibleUserID
	existingStakeholderActionPlan.ReviewUserID = StakeholderActionPlanFormDTO.ReviewUserID
	existingStakeholderActionPlan.TargetDate = dates.TargetDate
	existingStakeholderActionPlan.UpdatedByID = &currentUserID

	err = caps.repo.Update(existingStakeholderActionPlan, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	// //update task
	urlParam := fmt.Sprintf(`{"stakeholder_id":%d, "action_plan_id":%d}`, stakeholderID, id)

	taskUpdateDTO := commonDTOs.TaskFormDTO{
		Title:             StakeholderActionPlanFormDTO.Title,
		Description:       StakeholderActionPlanFormDTO.Description,
		TargetDate:        dates.TargetDate,
		URLParam:          &urlParam,
		SeverityLevelID:   StakeholderActionPlanFormDTO.SeverityLevelID,
		ResponsibleUserID: StakeholderActionPlanFormDTO.ResponsibleUserID,
	}

	err = caps.taskSvc.Update(currentUserID, taskUpdateDTO, "stakeholder_action_plan_id", id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	go func() {
		existingActionPlan, err := caps.GetFullDetailsByID(stakeholderID, id)
		if err != nil {
			return
		}
		// log changes
		changedData := make(map[string]interface{})
		oldData := map[string]interface{}{
			"title":              oldStakeholderActionPlan.Title,
			"description":        oldStakeholderActionPlan.Description,
			"severity_level":     oldStakeholderActionPlan.SeverityLevel.TitleEn,
			"responsible_user":   fmt.Sprintf("%s %s", oldStakeholderActionPlan.ResponsibleUser.FirstNameEN, oldStakeholderActionPlan.ResponsibleUser.LastNameEN),
			"review_user":        utility.GetFieldValue(oldStakeholderActionPlan.ReviewUser),
			"target_date":        oldStakeholderActionPlan.TargetDateString,
			"action_plan_status": oldStakeholderActionPlan.ActionPlanStatus.TitleEn,
		}

		newData := map[string]interface{}{
			"title":              existingActionPlan.Title,
			"description":        existingActionPlan.Description,
			"severity_level":     existingActionPlan.SeverityLevel.TitleEn,
			"responsible_user":   fmt.Sprintf("%s %s", existingActionPlan.ResponsibleUser.FirstNameEN, existingActionPlan.ResponsibleUser.LastNameEN),
			"review_user":        utility.GetFieldValue(existingActionPlan.ReviewUser),
			"target_date":        existingActionPlan.TargetDateString,
			"action_plan_status": existingActionPlan.ActionPlanStatus.TitleEn,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := caps.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.StakeholderActionPlan{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.StakeholderActionPlan{}.TableName(),
					"recordID":  id,
					"userID":    currentUserID,
					"error":     err.Error(),
				}

				utility.LogError("Failed to log activity", fields)
			}
		}
		// sends notification if responsible user has changed
		if oldResponsibleUserID != existingStakeholderActionPlan.ResponsibleUserID {
			urlParam := fmt.Sprintf(`{"stakeholder_id":%d, "action_plan_id":%d}`, stakeholderID, id)
			queryParam := ""
			entity := "stakeholder-action-plan"

			oldResponsibleUserAction := "responsible-user-remove"
			oldResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{oldResponsibleUserID},
				ModuleIdentifier:     "stakeholder_analysis",
				NotificationType:     string(commonDTOs.InformationType),
				HeadingEn:            "Stakeholder Action Plan Responsibility Removed",
				HeadingAr:            "تمت إزالة مسؤولية خطة عمل أصحاب المصلحة",
				TitleEn:              existingStakeholderActionPlan.Title,
				TitleAr:              existingStakeholderActionPlan.Title,
				Entity:               &entity,
				Action:               &oldResponsibleUserAction,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			newResponsibleUserAction := "responsible-user-assign"
			newResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{StakeholderActionPlanFormDTO.ResponsibleUserID},
				ModuleIdentifier:     "stakeholder_analysis",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Stakeholder Action Plan Responsibility Assigned",
				HeadingAr:            "مسؤولية خطة عمل أصحاب المصلحة المعينة",
				TitleEn:              existingStakeholderActionPlan.Title,
				TitleAr:              existingStakeholderActionPlan.Title,
				Entity:               &entity,
				Action:               &newResponsibleUserAction,
				URLParam:             urlParam,
				QueryParam:           queryParam,
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

			newResponsibleUser, err := caps.userSvc.GetByID(StakeholderActionPlanFormDTO.ResponsibleUserID)
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
			oldResponsibleUserSubject := fmt.Sprintf("Stakeholder Action Plan/%s: Responsibility Removed", existingStakeholderActionPlan.Title)
			oldResponsibleUserTemplateRelativePath := "organizational-context/templates/email/stakeholder_action_plan_remove_responsible_user_mail.html"
			oldResponsibleUserEmailData := map[string]interface{}{
				"FirstNameEn":           oldResponsibleUser.FirstNameEn,
				"LastNameEn":            oldResponsibleUser.LastNameEn,
				"StakeholderTitle":      stakeholder.TitleEn,
				"ActionPlanTitle":       existingStakeholderActionPlan.Title,
				"ActionPlanDescription": existingStakeholderActionPlan.Description,
				"RedirectURL":           fmt.Sprintf("%sstakeholder/%d/action-plans?action_plan_id=%d", appURL, stakeholderID, id),
			}

			oldResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              oldResponsibleUserEmailTo,
				Subject:              oldResponsibleUserSubject,
				TemplateRelativePath: oldResponsibleUserTemplateRelativePath,
				EmailData:            oldResponsibleUserEmailData,
			}

			newResponsibleUserEmailTo := []string{newResponsibleUser.Email}
			newResponsibleUserSubject := fmt.Sprintf("Stakeholder Action Plan/%s: Responsibility Assigned", existingStakeholderActionPlan.Title)
			newResponsibleUserTemplateRelativePath := "organizational-context/templates/email/stakeholder_action_plan_change_responsible_user_mail.html"
			newResponsibleUserEmailData := map[string]interface{}{
				"FirstNameEn":           newResponsibleUser.FirstNameEn,
				"LastNameEn":            newResponsibleUser.LastNameEn,
				"StakeholderTitle":      stakeholder.TitleEn,
				"ActionPlanTitle":       existingStakeholderActionPlan.Title,
				"ActionPlanDescription": existingStakeholderActionPlan.Description,
				"TargetDate":            existingStakeholderActionPlan.TargetDate.Format("2006-01-02"),
				"AssignedByFirstNameEn": assignedBy.FirstNameEn,
				"AssignedByLastNameEn":  assignedBy.LastNameEn,
				"AssignedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           fmt.Sprintf("%sstakeholder/%d/action-plans?action_plan_id=%d", appURL, stakeholderID, id),
			}

			newResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              newResponsibleUserEmailTo,
				Subject:              newResponsibleUserSubject,
				TemplateRelativePath: newResponsibleUserTemplateRelativePath,
				EmailData:            newResponsibleUserEmailData,
			}

			err = caps.notificationSvc.SendNotification(
				newResponsibleUserNotificationDTO, newResponsibleUserEmailNotificationDTO, "stakeholder_action_plan",
			)
			if err != nil {
				utility.LogError("Failed to send new responsible user notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}

			err = caps.notificationSvc.SendNotification(
				oldResponsibleUserNotificationDTO, oldResponsibleUserEmailNotificationDTO, "stakeholder_action_plan",
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

func (caps *StakeholderActionPlanService) Delete(stakeholderID uint, ids []uint, currentUserID uint) ([]uint, error) {
	if _, err := caps.stakeholderSvc.GetByID(stakeholderID); err != nil {
		return nil, err
	}

	var failedIds []uint
	var successfulIds []uint
	var deletedStakeholderActionPlan []map[string]interface{}
	var notifications []struct {
		OwnerID          *uint
		TitleEn          string
		TitleAr          string
		Description      string
		ReferenceNumber  string
		StakeholderTitle string
	}
	for _, id := range ids {

		stakeholderActionPlan, err := caps.GetFullDetailsByID(stakeholderID, id)
		if err != nil {
			continue
		}

		// Prepare notification data before deletion
		notifications = append(notifications, struct {
			OwnerID          *uint
			TitleEn          string
			TitleAr          string
			Description      string
			ReferenceNumber  string
			StakeholderTitle string
		}{
			OwnerID:          &stakeholderActionPlan.ResponsibleUserID,
			TitleEn:          stakeholderActionPlan.Title,
			TitleAr:          stakeholderActionPlan.Title,
			Description:      *stakeholderActionPlan.Description,
			StakeholderTitle: stakeholderActionPlan.Stakeholder.TitleEn,
		})

		if len(stakeholderActionPlan.ActionPlanStatusUpdates) > 0 {
			var deleteStatusUpdateIds []*uint
			for _, update := range stakeholderActionPlan.ActionPlanStatusUpdates {
				deleteStatusUpdateIds = append(deleteStatusUpdateIds, &update.ID)
			}
			if err := caps.repo.DeleteStakeholderActionPlanStatusUpdates(deleteStatusUpdateIds); err != nil {
				failedIds = append(failedIds, id)
				continue
			}
		}
		existingTask, err := caps.taskSvc.GetByFieldID("stakeholder_action_plan_id", id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		taskID, err := caps.taskSvc.Delete(existingTask.ID, currentUserID)
		if err != nil {
			failedIds = append(failedIds, taskID)
		}

		err = caps.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedStakeholderActionPlan = append(deletedStakeholderActionPlan, map[string]interface{}{
			"id":             id,
			"stakeholder_id": stakeholderID,
			"title":          stakeholderActionPlan.Title,
		})

	}

	if len(failedIds) > 0 {
		return failedIds, nil
	}

	if len(successfulIds) > 0 {
		go func() {
			err := caps.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.StakeholderActionPlan{}.TableName(),
				successfulIds,
				deletedStakeholderActionPlan,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}

			for _, notification := range notifications {
				urlParam := ""
				entity := "stakeholder-action-plan"
				action := "delete"

				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserIDs:              []uint{*notification.OwnerID},
					ModuleIdentifier:     "stakeholder_analysis",
					NotificationType:     string(commonDTOs.DeletionType),
					HeadingEn:            "Stakeholder Action Plan Deleted",
					HeadingAr:            "تم حذف خطة عمل أصحاب المصلحة",
					TitleEn:              notification.TitleEn,
					TitleAr:              notification.TitleAr,
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				ownerUser, err := caps.userSvc.GetByID(*notification.OwnerID)
				if err != nil {
					utility.LogError("Failed to fetch owner details for notification", map[string]interface{}{
						"error": err.Error(),
					})
					continue
				}

				ownerEmailTo := []string{ownerUser.Email}
				subject := fmt.Sprintf("Stakeholder Action Plan/%s: Deleted", notification.TitleEn)
				templateRelativePath := "organizational-context/templates/email/stakeholder_action_plan_delete_mail.html"

				emailData := map[string]interface{}{
					"FirstNameEn":                ownerUser.FirstNameEn,
					"LastNameEn":                 ownerUser.LastNameEn,
					"StakeholderReferenceNumber": notification.ReferenceNumber,
					"StakeholderTitle":           notification.StakeholderTitle,
					"ActionPlanTitle":            notification.TitleEn,
					"ActionPlanDescription":      notification.Description,
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              ownerEmailTo,
					Subject:              subject,
					TemplateRelativePath: templateRelativePath,
					EmailData:            emailData,
				}

				err = caps.notificationSvc.SendNotification(
					dbNotificationDTO, emailNotificationDTO, "stakeholder_action_plan")
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

func (caps *StakeholderActionPlanService) StatusUpdate(stakeholderID uint, id uint, currentUserID uint, stakeholderActionPlanStatusUpdateFormDTO *organizationalContextDTOs.StakeholderActionPlanStatusUpdateFormDTO, tx utility.DatabaseTransaction) error {
	existingStakeholderActionPlan, err := caps.GetByID(stakeholderID, id)
	if err != nil {
		return err
	}

	existingStakeholderActionPlan.ActionPlanStatusID = stakeholderActionPlanStatusUpdateFormDTO.ActionPlanStatusID
	existingStakeholderActionPlan.Comments = &stakeholderActionPlanStatusUpdateFormDTO.Comment
	existingStakeholderActionPlan.Percentage = stakeholderActionPlanStatusUpdateFormDTO.Percentage
	existingStakeholderActionPlan.UpdatedByID = &currentUserID

	err = caps.repo.Update(existingStakeholderActionPlan, tx)
	if err != nil {
		return err
	}

	isReviewUserNull := existingStakeholderActionPlan.ReviewUserID == nil

	// Fetch the corresponding task and update its status
	existingTask, err := caps.taskSvc.GetByFieldID("stakeholder_action_plan_id", id)
	if err != nil {
		return err
	}
	err = caps.taskSvc.UpdateStatus(currentUserID, existingTask.ID, stakeholderActionPlanStatusUpdateFormDTO.ActionPlanStatusID, isReviewUserNull, tx)
	if err != nil {
		return err
	}

	return nil
}

func (caps *StakeholderActionPlanService) validateStakeholderActionPlanFormDTO(stakeholderID uint, StakeholderActionPlanFormDTO *organizationalContextDTOs.StakeholderActionPlanFormDTO) error {
	if StakeholderActionPlanFormDTO.SeverityLevelID != nil {
		if _, err := caps.severityLevelSvc.GetByID(*StakeholderActionPlanFormDTO.SeverityLevelID); err != nil {
			return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
		}
	}

	if StakeholderActionPlanFormDTO.ReviewUserID != nil {
		if StakeholderActionPlanFormDTO.ResponsibleUserID == *StakeholderActionPlanFormDTO.ReviewUserID {
			return utility.NewValidationError(fmt.Errorf("review user and responsible user must be distinct"), "distinct_review_user", "stakeholder action plan distinct review user")
		}
	}

	return nil
}

func (caps *StakeholderActionPlanService) parseStakeholderActionPlanDates(dto *organizationalContextDTOs.StakeholderActionPlanFormDTO) (*organizationalContextDTOs.StakeholderActionPlanParsedDates, error) {
	parsed := &organizationalContextDTOs.StakeholderActionPlanParsedDates{}

	var err error
	if dto.TargetDate != "" {
		parsed.TargetDate, err = utility.ParseDate(dto.TargetDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid target date"), "invalid_date", "target date")
		}
	}

	return parsed, nil
}
