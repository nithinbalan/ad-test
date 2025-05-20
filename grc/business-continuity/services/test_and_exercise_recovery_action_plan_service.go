package services

import (
	"errors"
	"fmt"
	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	TestAndExerciseRecoveryActionPlanRepositories "jethur/grc/business-continuity/repositories"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	taskServices "jethur/grc/task/services"
	"jethur/utility"
	"time"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type TestAndExerciseRecoveryActionPlanServiceInterface interface {
	GetAll(testAndExerciseId uint, params *businessContinuityDTOs.TestAndExerciseRecoveryActionPlanQueryParamsDTO) ([]*businessContinuityDTOs.TestAndExerciseRecoveryActionPlanListDTO, error)
	GetByID(testAndExerciseID uint, id uint) (*models.TestAndExerciseRecoveryActionPlan, error)
	GetFullDetailsByID(testAndExerciseID uint, id uint, tx ...utility.DatabaseTransaction) (*businessContinuityDTOs.TestAndExerciseRecoveryActionPlanDetailDTO, error)
	Create(testAndExerciseID uint, currentUserID uint, dto *businessContinuityDTOs.TestAndExerciseRecoveryActionPlanFormDTO) (uint, error)
	Update(testAndExerciseID uint, id uint, currentUserID uint, dto *businessContinuityDTOs.TestAndExerciseRecoveryActionPlanFormDTO) error
	Delete(testAndExerciseID uint, ids []uint, currentUserID uint) ([]uint, error)
	StatusUpdate(testAndExerciseId uint, id uint, currentUserID uint, testAndExerciseRecoveryActionPlanStatusUpdateFormDTO *businessContinuityDTOs.TaeRecoveryActionPlanStatusUpdateFormDTO, tx utility.DatabaseTransaction) error
}

type TestAndExerciseRecoveryActionPlanService struct {
	repo                TestAndExerciseRecoveryActionPlanRepositories.TestAndExerciseRecoveryActionPlanRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	userSvc             organizationServices.UserServiceInterface
	testAndExerciseSvc  TestAndExerciseServiceInterface
	ActionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface
	taskSvc             taskServices.TaskServiceInterface
	notificationSvc     commonServices.NotificationServiceInterface
	userSettingSvc      organizationServices.UserSettingServiceInterface
}

func NewTestAndExerciseRecoveryActionPlanService(
	repo TestAndExerciseRecoveryActionPlanRepositories.TestAndExerciseRecoveryActionPlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	testAndExerciseSvc TestAndExerciseServiceInterface,
	ActionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	taskSvc taskServices.TaskServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
) TestAndExerciseRecoveryActionPlanServiceInterface {
	return &TestAndExerciseRecoveryActionPlanService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		userSvc:             userSvc,
		testAndExerciseSvc:  testAndExerciseSvc,
		ActionPlanStatusSvc: ActionPlanStatusSvc,
		taskSvc:             taskSvc,
		notificationSvc:     notificationSvc,
		userSettingSvc:      userSettingSvc,
	}
}

func (teraps *TestAndExerciseRecoveryActionPlanService) GetAll(testAndExerciseId uint, params *businessContinuityDTOs.TestAndExerciseRecoveryActionPlanQueryParamsDTO) ([]*businessContinuityDTOs.TestAndExerciseRecoveryActionPlanListDTO, error) {
	if _, err := teraps.testAndExerciseSvc.GetByID(testAndExerciseId); err != nil {
		return nil, err
	}
	return teraps.repo.GetAll(testAndExerciseId, params)
}

func (teraps *TestAndExerciseRecoveryActionPlanService) GetByID(testAndExerciseID uint, id uint) (*models.TestAndExerciseRecoveryActionPlan, error) {
	actionPlan, err := teraps.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return actionPlan, nil
}

func (teraps *TestAndExerciseRecoveryActionPlanService) GetFullDetailsByID(testAndExerciseID uint, id uint, tx ...utility.DatabaseTransaction) (*businessContinuityDTOs.TestAndExerciseRecoveryActionPlanDetailDTO, error) {
	if _, err := teraps.testAndExerciseSvc.GetByID(testAndExerciseID); err != nil {
		return nil, err
	}
	var recoveryActionPlan businessContinuityDTOs.TestAndExerciseRecoveryActionPlanDetailDTO

	actionPlan, err := teraps.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(actionPlan, &recoveryActionPlan)
	if err != nil {
		return nil, err
	}

	if actionPlan.TargetDate != nil {
		targetDate := actionPlan.TargetDate.Format("2006-01-02")
		recoveryActionPlan.TargetDateString = &targetDate
	}

	recoveryActionPlan.CreatedAt = actionPlan.CreatedAt
	recoveryActionPlan.UpdatedAt = actionPlan.UpdatedAt

	return &recoveryActionPlan, nil
}

func (teraps *TestAndExerciseRecoveryActionPlanService) Create(testAndExerciseID uint, currentUserID uint, recoveryActionPlanFormDTO *businessContinuityDTOs.TestAndExerciseRecoveryActionPlanFormDTO) (uint, error) {
	tx := teraps.repo.BeginTransaction()

	testAndExercise, err := teraps.testAndExerciseSvc.GetFullDetailsByID(testAndExerciseID)
	if err != nil {
		return 0, err
	}

	if _, err := teraps.testAndExerciseSvc.GetByID(testAndExerciseID); err != nil {
		return 0, err
	}

	responsibleUnitID, err := teraps.userSvc.GetUserPrimaryUnitID(recoveryActionPlanFormDTO.ResponsibleUserID)
	if err != nil {
		return 0, utility.NewValidationError(errors.New("organization unit not set for responsible user"), "not_exist", "responsible unit")
	}

	actionPlanStatus, err := teraps.ActionPlanStatusSvc.GetByType("new")
	if err != nil {
		return 0, err
	}

	dates, err := teraps.parseTestAndExerciseRecoveryActionPlanDates(recoveryActionPlanFormDTO)
	if err != nil {
		return 0, err
	}

	actionPlan := &models.TestAndExerciseRecoveryActionPlan{
		TestAndExerciseID:  testAndExerciseID,
		BcpPhaseID:         recoveryActionPlanFormDTO.BcpPhaseID,
		ProcessID:          recoveryActionPlanFormDTO.ProcessID,
		SeverityLevelID:    recoveryActionPlanFormDTO.SeverityLevelID,
		Title:              recoveryActionPlanFormDTO.Title,
		Description:        recoveryActionPlanFormDTO.Description,
		ResponsibleUserID:  recoveryActionPlanFormDTO.ResponsibleUserID,
		ResponsibleUnitID:  responsibleUnitID,
		ReviewUserID:       recoveryActionPlanFormDTO.ReviewUserID,
		TargetDate:         dates.TargetDate,
		ActionPlanStatusId: &actionPlanStatus.ID,
		IsRecovery:         recoveryActionPlanFormDTO.IsRecovery,
		Percentage:         recoveryActionPlanFormDTO.Percentage,
		CreatedBYID:        currentUserID,
	}

	id, err := teraps.repo.Create(actionPlan, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	actionPlanStatusID := actionPlanStatus.ID
	recoveryActionPlanUpdate := &models.TaeRecoveryActionPlanStatusUpdate{
		TestAndExerciseRecoveryActionPlanID: id,
		ActionPlanStatusID:                  actionPlanStatusID,
		Comment:                             "Initial action plan status",
		Percentage:                          recoveryActionPlanFormDTO.Percentage,
		CreatedByID:                         currentUserID,
	}

	_, err = teraps.repo.CreateInitialSLAActionPlanStatusUpdate(recoveryActionPlanUpdate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	urlParam := fmt.Sprintf(`{"test_and_exercise_id":%d, "action_plan_id":%d}`, testAndExerciseID, id)
	if recoveryActionPlanFormDTO.IsRecovery {
		urlParam = fmt.Sprintf(`{"test_and_exercise_id":%d, "recovery_action_plan_id":%d}`, testAndExerciseID, id)
	}
	taskDTO := commonDTOs.TaskFormDTO{
		ModuleIdentifier:     "business_continuity",
		RecoveryActionPlanID: &id,
		Title:                recoveryActionPlanFormDTO.Title,
		TargetDate:           dates.TargetDate,
		URLParam:             &urlParam,
		QueryParam:           nil,
		SeverityLevelID:      recoveryActionPlanFormDTO.SeverityLevelID,
		ResponsibleUserID:    recoveryActionPlanFormDTO.ResponsibleUserID,
	}

	if _, err := teraps.taskSvc.Create(taskDTO, currentUserID, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	go func() {
		existingTestAndExerciseRecoveryActionPlan, err := teraps.GetFullDetailsByID(testAndExerciseID, actionPlan.ID)
		if err != nil {
			utility.LogError("Failed to fetch test and exercise", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = teraps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.TestAndExerciseRecoveryActionPlan{}.TableName(),
			id,
			existingTestAndExerciseRecoveryActionPlan,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.TestAndExerciseRecoveryActionPlan{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
		// prepare notification
		appURL := utility.GetAppURL()
		urlParam := fmt.Sprintf(`{"test_and_exercise_id":%d, "action_plan_id":%d}`, testAndExerciseID, id)
		redirectUrl := fmt.Sprintf("%sbusiness-continuity/test-and-exercise/%d/action-plans?action_plan_id=%d", appURL, testAndExerciseID, id)
		if recoveryActionPlanFormDTO.IsRecovery {
			urlParam = fmt.Sprintf(`{"test_and_exercise_id":%d, "recovery_action_plan_id":%d}`, testAndExerciseID, id)
			redirectUrl = fmt.Sprintf("%sbusiness-continuity/test-and-exercise/%d/recovery-action-plans?action_plan_id=%d", appURL, testAndExerciseID, id)
		}
		queryParam := ""
		entity := "test-and-exercise-action-plan"
		action := "responsible-user-assign"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{recoveryActionPlanFormDTO.ResponsibleUserID},
			ModuleIdentifier:     "business_continuity",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "Action Plan Responsibility Assigned",
			HeadingAr:            "مسؤولية خطة العمل الموكلة",
			TitleEn:              recoveryActionPlanFormDTO.Title,
			TitleAr:              recoveryActionPlanFormDTO.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			QueryParam:           queryParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		responsibleUser, err := teraps.userSvc.GetByID(recoveryActionPlanFormDTO.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch responsible user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		assignedBy, err := teraps.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch assigning user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("Action Plan/%s: Responsibility Assigned", recoveryActionPlanFormDTO.Title)
		emailData := map[string]interface{}{
			"FirstNameEn":           responsibleUser.FirstNameEn,
			"LastNameEn":            responsibleUser.LastNameEn,
			"ReferenceNumber":       testAndExercise.ReferenceNumber,
			"TestAndExerciseTitle":  testAndExercise.BcpTestAndExercise.Title,
			"ActionPlanTitle":       recoveryActionPlanFormDTO.Title,
			"TargetDate":            dates.TargetDate.Format("2006-01-02"),
			"AssignedByFirstNameEn": assignedBy.FirstNameEn,
			"AssignedByLastNameEn":  assignedBy.LastNameEn,
			"AssignedOn":            time.Now().Format("2006-01-02"),
			"RedirectURL":           redirectUrl,
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "business-continuity/templates/email/action_plan_create_responsible_user_mail.html",
			EmailData:            emailData,
		}

		// send notification and mail
		err = teraps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "test_and_exercise")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()
	return id, nil
}

func (teraps *TestAndExerciseRecoveryActionPlanService) Update(testAndExerciseID uint, id uint, currentUserID uint, recoveryActionPlanFormDTO *businessContinuityDTOs.TestAndExerciseRecoveryActionPlanFormDTO) error {
	tx := teraps.repo.BeginTransaction()

	testAndExercise, err := teraps.testAndExerciseSvc.GetFullDetailsByID(testAndExerciseID)
	if err != nil {
		return err
	}

	existingRecoveryPlan, err := teraps.repo.GetByID(id)
	if err != nil {
		return utility.NewResourceNotFoundError(err, id)
	}

	oldRecoveryActionPlan, err := teraps.GetFullDetailsByID(testAndExerciseID, id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	oldResponsibleUserID := existingRecoveryPlan.ResponsibleUserID

	if existingRecoveryPlan.ResponsibleUserID != recoveryActionPlanFormDTO.ResponsibleUserID {
		responsibleUnitID, err := teraps.userSvc.GetUserPrimaryUnitID(existingRecoveryPlan.ResponsibleUserID)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "risk responsible unit")
		}
		existingRecoveryPlan.ResponsibleUnitID = responsibleUnitID
	}

	dates, err := teraps.parseTestAndExerciseRecoveryActionPlanDates(recoveryActionPlanFormDTO)
	if err != nil {
		return err
	}

	existingRecoveryPlan.BcpPhaseID = recoveryActionPlanFormDTO.BcpPhaseID
	existingRecoveryPlan.ProcessID = recoveryActionPlanFormDTO.ProcessID
	existingRecoveryPlan.SeverityLevelID = recoveryActionPlanFormDTO.SeverityLevelID
	existingRecoveryPlan.Title = recoveryActionPlanFormDTO.Title
	existingRecoveryPlan.Description = recoveryActionPlanFormDTO.Description
	existingRecoveryPlan.ResponsibleUserID = recoveryActionPlanFormDTO.ResponsibleUserID
	existingRecoveryPlan.TargetDate = dates.TargetDate
	existingRecoveryPlan.IsRecovery = recoveryActionPlanFormDTO.IsRecovery
	existingRecoveryPlan.Percentage = recoveryActionPlanFormDTO.Percentage
	existingRecoveryPlan.ReviewUserID = recoveryActionPlanFormDTO.ReviewUserID
	existingRecoveryPlan.UpdatedBYID = &currentUserID

	err = teraps.repo.Update(existingRecoveryPlan, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// //update task
	task, _ := teraps.taskSvc.GetByResponsibleUserID("recovery_action_plan_id", id, recoveryActionPlanFormDTO.ResponsibleUserID)

	if task != nil {
		urlParam := fmt.Sprintf(`{"test_and_exercise_id":%d, "action_plan_id":%d}`, testAndExerciseID, id)
		if recoveryActionPlanFormDTO.IsRecovery {
			urlParam = fmt.Sprintf(`{"test_and_exercise_id":%d, "recovery_action_plan_id":%d}`, testAndExerciseID, id)
		}
		taskUpdateDTO := commonDTOs.TaskFormDTO{
			Title:             recoveryActionPlanFormDTO.Title,
			TargetDate:        dates.TargetDate,
			URLParam:          &urlParam,
			SeverityLevelID:   recoveryActionPlanFormDTO.SeverityLevelID,
			ResponsibleUserID: recoveryActionPlanFormDTO.ResponsibleUserID,
		}

		err = teraps.taskSvc.Update(currentUserID, taskUpdateDTO, "recovery_action_plan_id", id, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	go func() {
		existingActionPlan, err := teraps.GetFullDetailsByID(testAndExerciseID, id)
		if err != nil {
			return
		}

		changedData := make(map[string]interface{})

		// log changes
		oldData := map[string]interface{}{
			"responsible_user": fmt.Sprintf("%s %s", oldRecoveryActionPlan.ResponsibleUser.FirstNameEN, oldRecoveryActionPlan.ResponsibleUser.LastNameEN),
			"bcp_phase":        oldRecoveryActionPlan.BcpPhase.TitleEn,
			"process":          oldRecoveryActionPlan.Process.TitleEn,
			"severity_level":   oldRecoveryActionPlan.SeverityLevel.TitleEn,
			"target_date":      oldRecoveryActionPlan.TargetDateString,
		}

		newData := map[string]interface{}{
			"responsible_user": fmt.Sprintf("%s %s", existingActionPlan.ResponsibleUser.FirstNameEN, existingActionPlan.ResponsibleUser.LastNameEN),
			"bcp_phase":        existingActionPlan.BcpPhase.TitleEn,
			"process":          existingActionPlan.Process.TitleEn,
			"severity_level":   existingActionPlan.SeverityLevel.TitleEn,
			"target_date":      existingActionPlan.TargetDateString,
		}

		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}

		if len(changedData) > 0 {
			err := teraps.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.TestAndExerciseRecoveryActionPlan{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.TestAndExerciseRecoveryActionPlan{}.TableName(),
					"recordID":  id,
					"userID":    currentUserID,
					"error":     err.Error(),
				}

				utility.LogError("Failed to log activity", fields)
			}
		}

		appURL := utility.GetAppURL()
		urlParam := fmt.Sprintf(`{"test_and_exercise_id":%d, "action_plan_id":%d}`, testAndExerciseID, id)
		redirectUrl := fmt.Sprintf("%sbusiness-continuity/test-and-exercise/%d/action-plans?action_plan_id=%d", appURL, testAndExerciseID, id)
		if recoveryActionPlanFormDTO.IsRecovery {
			urlParam = fmt.Sprintf(`{"test_and_exercise_id":%d, "recovery_action_plan_id":%d}`, testAndExerciseID, id)
			redirectUrl = fmt.Sprintf("%sbusiness-continuity/test-and-exercise/%d/recovery-action-plans?action_plan_id=%d", appURL, testAndExerciseID, id)
		}

		// sends notification if responsible user has changed
		if oldResponsibleUserID != existingActionPlan.ResponsibleUserID {
			queryParam := ""
			entity := "test-and-exercise-action-plan"

			oldResponsibleUserAction := "responsible-user-remove"
			oldResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{oldResponsibleUserID},
				ModuleIdentifier:     "business_continuity",
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
				UserIDs:              []uint{recoveryActionPlanFormDTO.ResponsibleUserID},
				ModuleIdentifier:     "business_continuity",
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

			oldResponsibleUser, err := teraps.userSvc.GetByID(oldResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch old responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			newResponsibleUser, err := teraps.userSvc.GetByID(recoveryActionPlanFormDTO.ResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch new responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			assignedBy, err := teraps.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch assigning user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			oldResponsibleUserEmailTo := []string{oldResponsibleUser.Email}
			oldResponsibleUserSubject := fmt.Sprintf("Action Plan/%s: Responsibility Removed", existingActionPlan.Title)
			oldResponsibleUserTemplateRelativePath := "business-continuity/templates/email/action_plan_remove_responsible_user_mail.html"
			oldResponsibleUserEmailData := map[string]interface{}{
				"FirstNameEn":          oldResponsibleUser.FirstNameEn,
				"LastNameEn":           oldResponsibleUser.LastNameEn,
				"ReferenceNumber":      testAndExercise.ReferenceNumber,
				"TestAndExerciseTitle": testAndExercise.BcpTestAndExercise.Title,
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
			newResponsibleUserTemplateRelativePath := "business-continuity/templates/email/action_plan_change_responsible_user_mail.html"
			newResponsibleUserEmailData := map[string]interface{}{
				"FirstNameEn":           newResponsibleUser.FirstNameEn,
				"LastNameEn":            newResponsibleUser.LastNameEn,
				"ReferenceNumber":       testAndExercise.ReferenceNumber,
				"TestAndExerciseTitle":  testAndExercise.BcpTestAndExercise.Title,
				"ActionPlanTitle":       existingRecoveryPlan.Title,
				"TargetDate":            existingRecoveryPlan.TargetDate.Format("2006-01-02"),
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

			err = teraps.notificationSvc.SendNotification(
				newResponsibleUserNotificationDTO, newResponsibleUserEmailNotificationDTO, "test_and_exercise",
			)
			if err != nil {
				utility.LogError("Failed to send new responsible user notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}

			err = teraps.notificationSvc.SendNotification(
				oldResponsibleUserNotificationDTO, oldResponsibleUserEmailNotificationDTO, "test_and_exercise",
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

func (teraps *TestAndExerciseRecoveryActionPlanService) Delete(testAndExerciseID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := teraps.repo.BeginTransaction()

	if _, err := teraps.testAndExerciseSvc.GetByID(testAndExerciseID); err != nil {
		return nil, err
	}

	testAndExercise, err := teraps.testAndExerciseSvc.GetFullDetailsByID(testAndExerciseID)
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
		TestAndExerciseTitle string
	}

	for _, id := range ids {
		actionPlan, err := teraps.GetFullDetailsByID(testAndExerciseID, id, tx)
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
			TestAndExerciseTitle string
		}{
			OwnerID:              &actionPlan.ResponsibleUserID,
			TitleEn:              actionPlan.Title,
			TitleAr:              actionPlan.Title,
			ReferenceNumber:      testAndExercise.ReferenceNumber,
			TestAndExerciseTitle: testAndExercise.BcpTestAndExercise.Title,
		})

		deletedActionPlans = append(deletedActionPlans, map[string]interface{}{
			"id":                id,
			"test_and_exercise": actionPlan.ID,
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
			if err := teraps.repo.DeleteActionPlanStatusUpdateDocuments(deleteDocumentIds); err != nil {
				failedIDs = append(failedIDs, id)
				continue
			}
			if err := teraps.repo.DeleteActionPlanStatusUpdates(deleteStatusUpdateIds); err != nil {
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
			if err := teraps.repo.DeleteActionPlanCommentDocuments(deleteCommentDocumentIds); err != nil {
				failedIDs = append(failedIDs, id)
				continue
			}
			if err := teraps.repo.DeleteActionPlanComments(deleteCommentIds); err != nil {
				failedIDs = append(failedIDs, id)
				continue
			}
		}

		existingTask, err := teraps.taskSvc.GetByFieldID("recovery_action_plan_id", id)
		if err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		taskID, err := teraps.taskSvc.Delete(existingTask.ID, currentUserID)
		if err != nil {
			failedIDs = append(failedIDs, taskID)
		}

		if err := teraps.repo.Delete(id, tx); err != nil {
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
			err := teraps.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.TestAndExerciseRecoveryActionPlan{}.TableName(),
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
			entity := "test-and-exercise-action-plan"
			action := "delete"

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{*notification.OwnerID},
				ModuleIdentifier:     "business_continuity",
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

			ownerUser, err := teraps.userSvc.GetByID(*notification.OwnerID)
			if err != nil {
				utility.LogError("Failed to fetch owner details for notification", map[string]interface{}{
					"error": err.Error(),
				})
				continue
			}

			ownerEmailTo := []string{ownerUser.Email}
			subject := fmt.Sprintf("Action Plan/%s: Deleted", notification.TitleEn)
			templateRelativePath := "business-continuity/templates/email/action_plan_delete_mail.html"

			emailData := map[string]interface{}{
				"FirstNameEn":           ownerUser.FirstNameEn,
				"LastNameEn":            ownerUser.LastNameEn,
				"ReferenceNumber":       notification.ReferenceNumber,
				"TestAndExerciseTitle":  notification.TestAndExerciseTitle,
				"ActionPlanTitle":       notification.TitleEn,
				"ActionPlanDescription": notification.Description,
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              ownerEmailTo,
				Subject:              subject,
				TemplateRelativePath: templateRelativePath,
				EmailData:            emailData,
			}

			err = teraps.notificationSvc.SendNotification(
				dbNotificationDTO, emailNotificationDTO, "test_and_exercise")
			if err != nil {
				utility.LogError("Failed to send deletion notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}

		}
	}

	return failedIDs, nil
}

func (teraps *TestAndExerciseRecoveryActionPlanService) parseTestAndExerciseRecoveryActionPlanDates(dto *businessContinuityDTOs.TestAndExerciseRecoveryActionPlanFormDTO) (*businessContinuityDTOs.TestAndExerciseRecoveryActionPlanParsedDates, error) {
	parsed := &businessContinuityDTOs.TestAndExerciseRecoveryActionPlanParsedDates{}

	var err error

	if &dto.TargetDate != nil {
		parsed.TargetDate, err = utility.ParseDate(dto.TargetDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid target date"), "invalid_date", "target date")
		}
	}

	return parsed, nil
}

func (teraps *TestAndExerciseRecoveryActionPlanService) StatusUpdate(actionPlanID uint, id uint, currentUserID uint, testAndExerciseRecoveryActionPlanStatusUpdateFormDTO *businessContinuityDTOs.TaeRecoveryActionPlanStatusUpdateFormDTO, tx utility.DatabaseTransaction) error {
	existingTestAndExerciseRecoveryActionPlan, err := teraps.GetByID(actionPlanID, id)
	if err != nil {
		return err
	}

	existingTestAndExerciseRecoveryActionPlan.ActionPlanStatusId = &testAndExerciseRecoveryActionPlanStatusUpdateFormDTO.ActionPlanStatusID
	existingTestAndExerciseRecoveryActionPlan.Percentage = testAndExerciseRecoveryActionPlanStatusUpdateFormDTO.Percentage
	existingTestAndExerciseRecoveryActionPlan.UpdatedBYID = &currentUserID

	err = teraps.repo.Update(existingTestAndExerciseRecoveryActionPlan, tx)
	if err != nil {
		return err
	}
	isReviewUserNull := existingTestAndExerciseRecoveryActionPlan.ReviewUserID == nil
	// Fetch the corresponding task and update its status

	task, _ := teraps.taskSvc.GetByResponsibleUserID("recovery_action_plan_id", actionPlanID, existingTestAndExerciseRecoveryActionPlan.ResponsibleUserID)

	if task != nil {

		existingTask, err := teraps.taskSvc.GetByFieldID("recovery_action_plan_id", id)
		if err != nil {
			return err
		}
		err = teraps.taskSvc.UpdateStatus(currentUserID, existingTask.ID, testAndExerciseRecoveryActionPlanStatusUpdateFormDTO.ActionPlanStatusID, isReviewUserNull, tx)
		if err != nil {
			return err
		}
	}

	return nil
}
