package services

import (
	"errors"
	"fmt"
	"os"
	"time"

	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	"jethur/grc/identity-access-review/dtos"
	"jethur/grc/identity-access-review/repositories"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	taskServices "jethur/grc/task/services"
	"jethur/utility"

	"gorm.io/gorm"
)

type IdentityReviewActionPlanServiceInterface interface {
	GetAll(identityReviewUserID uint, params *dtos.IdentityReviewActionPlanListQueryParamsDTO) (*dtos.IdentityReviewActionPlanListResponseDTO, error)
	GetByID(identityReviewUserID uint, id uint) (*models.IdentityReviewActionPlan, error)
	GetFullDetailsByID(identityReviewUserID uint, id uint) (*models.IdentityReviewActionPlan, error)
	Create(identityReviewUserID uint, currentUserID uint, formDTO *dtos.IdentityReviewActionPlanFormDTO) (uint, error)
	Update(identityReviewUserID uint, id uint, currentUserID uint, formDTO *dtos.IdentityReviewActionPlanFormDTO) error
	Delete(identityReviewUserID uint, ids []uint, currentUserID uint) ([]uint, error)
	ForceDelete(ids []uint, currentUserID uint) ([]uint, error)
	GetSummary(identityReviewUserID uint) (*dtos.IdentityReviewActionPlanSummaryDTO, error)
	Export(items []*dtos.IdentityReviewActionPlanListDTO, currentUserID uint) (*os.File, error)
	StatusUpdate(identityReviewUserID uint, id uint, currentUserID uint, formDTO *dtos.IdentityReviewActionPlanStatusUpdateFormDTO, tx ...utility.DatabaseTransaction) error
}

type IdentityReviewActionPlanService struct {
	repo                  repositories.IdentityReviewActionPlanRepositoryInterface
	activityLogSvc        logServices.ActivityLogServiceInterface
	identityReviewUserSvc IdentityReviewUserServiceInterface
	severityLevelSvc      masterServices.SeverityLevelServiceInterface
	actionPlanStatusSvc   masterServices.ActionPlanStatusServiceInterface
	userSvc               organizationServices.UserServiceInterface
	userSettingSvc        organizationServices.UserSettingServiceInterface
	taskSvc               taskServices.TaskServiceInterface
	delegateeUserSvc      commonServices.DelegateeUserserviceInterface
	accountDecisionSvc    AccountDecisionServiceInterface
	privilegedSvc         PrivilegedAccountServiceInterface
	userCommentSvc        IdentityReviewUserCommentServiceInterface
}

func NewIdentityReviewActionPlanService(
	repo repositories.IdentityReviewActionPlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	identityReviewUserSvc IdentityReviewUserServiceInterface,
	severityLevelSvc masterServices.SeverityLevelServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
	taskSvc taskServices.TaskServiceInterface,
	delegateeUserSvc commonServices.DelegateeUserserviceInterface,
	accountDecisionSvc AccountDecisionServiceInterface,
	privilegedSvc PrivilegedAccountServiceInterface,
	userCommentSvc IdentityReviewUserCommentServiceInterface,
) IdentityReviewActionPlanServiceInterface {
	return &IdentityReviewActionPlanService{
		repo:                  repo,
		activityLogSvc:        activityLogSvc,
		identityReviewUserSvc: identityReviewUserSvc,
		severityLevelSvc:      severityLevelSvc,
		actionPlanStatusSvc:   actionPlanStatusSvc,
		userSvc:               userSvc,
		userSettingSvc:        userSettingSvc,
		taskSvc:               taskSvc,
		delegateeUserSvc:      delegateeUserSvc,
		accountDecisionSvc:    accountDecisionSvc,
		privilegedSvc:         privilegedSvc,
		userCommentSvc:        userCommentSvc,
	}
}

func (iraps *IdentityReviewActionPlanService) GetAll(identityReviewUserID uint, params *dtos.IdentityReviewActionPlanListQueryParamsDTO) (*dtos.IdentityReviewActionPlanListResponseDTO, error) {
	if identityReviewUserID != 0 {
		if _, err := iraps.identityReviewUserSvc.GetByID(identityReviewUserID); err != nil {
			return nil, err
		}
	}
	return iraps.repo.GetAll(identityReviewUserID, params)
}

func (iraps *IdentityReviewActionPlanService) GetByID(identityReviewUserID uint, id uint) (*models.IdentityReviewActionPlan, error) {
	actionPlan, err := iraps.repo.GetByID(identityReviewUserID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return actionPlan, nil
}

func (iraps *IdentityReviewActionPlanService) GetFullDetailsByID(identityReviewUserID uint, id uint) (*models.IdentityReviewActionPlan, error) {
	actionPlan, err := iraps.repo.GetFullDetailsByID(identityReviewUserID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return actionPlan, nil
}

func (iraps *IdentityReviewActionPlanService) Create(identityReviewUserID uint, currentUserID uint, formDTO *dtos.IdentityReviewActionPlanFormDTO) (uint, error) {
	if _, err := iraps.identityReviewUserSvc.GetByID(identityReviewUserID); err != nil {
		return 0, err
	}

	accountDecision, err := iraps.accountDecisionSvc.GetByID(formDTO.AccountDecisionID)
	if err != nil {
		return 0, utility.NewValidationError(fmt.Errorf("account decision does not exist"), "not_exist", "account_decision_id")
	}

	tx := iraps.repo.BeginTransaction()

	if _, err := iraps.identityReviewUserSvc.UpdateAccountDecision(identityReviewUserID, currentUserID, &dtos.IdentityReviewUserAccountDecisionUpdateFormDTO{
		AccountDecisionID: formDTO.AccountDecisionID,
	}, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	if accountDecision.Type == "keep" {
		if err := tx.Commit(); err != nil {
			return 0, err
		}
		return 0, nil
	}

	if err := iraps.validateForm(formDTO); err != nil {
		tx.Rollback()
		return 0, err
	}
	if accountDecision.Type == "modify" {
		if _, err := iraps.identityReviewUserSvc.UpdatePrivilege(
			identityReviewUserID,
			currentUserID,
			&dtos.IdentityReviewUserPrivilegeUpdateFormDTO{PrivilegeID: formDTO.PrivilegedID},
			tx,
		); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	dates, err := iraps.parseDates(formDTO)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	actionPlanStatus, err := iraps.actionPlanStatusSvc.GetByType("new")
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	actionPlan := &models.IdentityReviewActionPlan{
		IdentityReviewUserID: identityReviewUserID,
		Title:                formDTO.Title,
		Description:          formDTO.Description,
		SeverityLevelID:      formDTO.SeverityLevelID,
		TargetDate:           dates.TargetDate,
		ResponsibleUserID:    formDTO.ResponsibleUserID,
		ReviewUserID:         formDTO.ReviewUserID,
		ActionPlanStatusID:   actionPlanStatus.ID,
		Percentage:           formDTO.Percentage,
		Comments:             formDTO.Comments,
		CreatedByID:          currentUserID,
	}

	if formDTO.ResponsibleUserID != 0 {
		delegateeID, err := iraps.delegateeUserSvc.AssignDelegateeIfPermitted(
			formDTO.ResponsibleUserID,
			"identity_access_review",
		)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		if delegateeID != nil {
			actionPlan.DelegateeUserID = delegateeID
		}
	}

	id, err := iraps.repo.Create(actionPlan, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if formDTO.Description != nil {
		if _, err := iraps.userCommentSvc.CreateInitial(
			identityReviewUserID,
			formDTO.AccountDecisionID,
			currentUserID,
			*formDTO.Description,
			tx,
		); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	percentage := formDTO.Percentage
	if err := iraps.repo.CreateInitialStatusUpdate(&models.IdentityReviewActionPlanStatusUpdate{
		IdentityReviewActionPlanID: id,
		ActionPlanStatusID:         actionPlanStatus.ID,
		Comment:                    "Initial action plan status",
		Percentage:                 &percentage,
		CreatedBYID:                currentUserID,
	}, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	urlParam := fmt.Sprintf(`{"identity_review_user_id":%d,"action_plan_id":%d}`, identityReviewUserID, id)
	if _, err := iraps.taskSvc.Create(commonDTOs.TaskFormDTO{
		ModuleIdentifier:           "identity_access_review",
		IdentityReviewActionPlanID: &id,
		Title:                      formDTO.Title,
		Description:                formDTO.Description,
		TargetDate:                 dates.TargetDate,
		URLParam:                   &urlParam,
		SeverityLevelID:            formDTO.SeverityLevelID,
		ResponsibleUserID:          formDTO.ResponsibleUserID,
	}, currentUserID, tx); err != nil {
		tx.Rollback()
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	go func() {
		existingActionPlan, err := iraps.GetFullDetailsByID(identityReviewUserID, id)
		if err != nil {
			utility.LogError("Failed to fetch identity review action plan", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = iraps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.IdentityReviewActionPlan{}.TableName(),
			id,
			existingActionPlan,
			currentUserID,
			"identity_access_review", "identity_review_user", "action_plan", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.IdentityReviewActionPlan{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (iraps *IdentityReviewActionPlanService) Update(identityReviewUserID uint, id uint, currentUserID uint, formDTO *dtos.IdentityReviewActionPlanFormDTO) error {
	actionPlan, err := iraps.GetByID(identityReviewUserID, id)
	if err != nil {
		return err
	}

	oldActionPlan, err := iraps.GetFullDetailsByID(identityReviewUserID, id)
	if err != nil {
		return err
	}

	oldResponsibleUserID := actionPlan.ResponsibleUserID

	if err := iraps.validateForm(formDTO); err != nil {
		return err
	}

	dates, err := iraps.parseDates(formDTO)
	if err != nil {
		return err
	}

	actionPlan.Title = formDTO.Title
	actionPlan.Description = formDTO.Description
	actionPlan.SeverityLevelID = formDTO.SeverityLevelID
	actionPlan.TargetDate = dates.TargetDate
	actionPlan.ResponsibleUserID = formDTO.ResponsibleUserID
	actionPlan.ReviewUserID = formDTO.ReviewUserID
	actionPlan.Percentage = formDTO.Percentage
	actionPlan.Comments = formDTO.Comments
	actionPlan.UpdatedByID = &currentUserID

	if oldResponsibleUserID != formDTO.ResponsibleUserID {
		actionPlan.DelegateeUserID = nil

		if formDTO.ResponsibleUserID != 0 {
			delegateeID, err := iraps.delegateeUserSvc.AssignDelegateeIfPermitted(
				formDTO.ResponsibleUserID,
				"identity_access_review",
			)
			if err != nil {
				return err
			}

			if delegateeID != nil {
				actionPlan.DelegateeUserID = delegateeID
			}
		}
	}

	tx := iraps.repo.BeginTransaction()

	if err := iraps.repo.Update(actionPlan, tx); err != nil {
		tx.Rollback()
		return err
	}

	urlParam := fmt.Sprintf(
		`{"identity_review_user_id":%d,"action_plan_id":%d}`,
		identityReviewUserID,
		id,
	)

	taskUpdateDTO := commonDTOs.TaskFormDTO{
		Title:             formDTO.Title,
		Description:       formDTO.Description,
		TargetDate:        dates.TargetDate,
		URLParam:          &urlParam,
		SeverityLevelID:   formDTO.SeverityLevelID,
		ResponsibleUserID: formDTO.ResponsibleUserID,
	}

	if err := iraps.taskSvc.Update(currentUserID, taskUpdateDTO, "identity_review_action_plan_id", id, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	go func() {
		existingActionPlan, err := iraps.GetFullDetailsByID(identityReviewUserID, id)
		if err != nil {
			utility.LogError(
				"Failed to fetch identity review action plan",
				map[string]interface{}{
					"error": err.Error(),
				},
			)
			return
		}

		oldData := map[string]interface{}{
			"title":             oldActionPlan.Title,
			"description":       oldActionPlan.Description,
			"severity_level_id": oldActionPlan.SeverityLevelID,
			"responsible_user":  oldActionPlan.ResponsibleUserID,
			"review_user":       oldActionPlan.ReviewUserID,
			"target_date":       oldActionPlan.TargetDate,
			"percentage":        oldActionPlan.Percentage,
			"comment":           oldActionPlan.Comments,
		}

		newData := map[string]interface{}{
			"title":             existingActionPlan.Title,
			"description":       existingActionPlan.Description,
			"severity_level_id": existingActionPlan.SeverityLevelID,
			"responsible_user":  existingActionPlan.ResponsibleUserID,
			"review_user":       existingActionPlan.ReviewUserID,
			"target_date":       existingActionPlan.TargetDate,
			"percentage":        existingActionPlan.Percentage,
			"comment":           existingActionPlan.Comments,
		}

		changedData := utility.CompareFieldChanges(oldData, newData)

		if len(changedData) > 0 {
			if err := iraps.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.IdentityReviewActionPlan{}.TableName(),
				id,
				changedData,
				currentUserID,
				"identity_access_review",
				"identity_review_user",
				"action_plan",
				nil,
			); err != nil {
				utility.LogError(
					"Failed to log activity",
					map[string]interface{}{
						"action":    models.UpdateAction,
						"tableName": models.IdentityReviewActionPlan{}.TableName(),
						"recordID":  id,
						"userID":    currentUserID,
						"error":     err.Error(),
					},
				)
			}
		}
	}()

	return nil
}

func (iraps *IdentityReviewActionPlanService) Delete(identityReviewUserID uint, ids []uint, currentUserID uint) ([]uint, error) {
	if _, err := iraps.identityReviewUserSvc.GetByID(identityReviewUserID); err != nil {
		return nil, err
	}

	tx := iraps.repo.BeginTransaction()

	var failedIDs []uint
	var successfulIDs []uint
	var deletedActionPlans []map[string]interface{}

	for _, id := range ids {
		actionPlan, err := iraps.GetFullDetailsByID(identityReviewUserID, id)
		if err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		existingTask, err := iraps.taskSvc.GetByFieldID("identity_review_action_plan_id", id, tx)
		if err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		if _, err := iraps.taskSvc.Delete(existingTask.ID, currentUserID); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		deletedActionPlans = append(
			deletedActionPlans,
			map[string]interface{}{
				"id":    id,
				"title": actionPlan.Title,
			},
		)

		if len(actionPlan.ActionPlanStatusUpdates) > 0 {
			statusUpdateIDs := make([]*uint, 0, len(actionPlan.ActionPlanStatusUpdates))
			for _, update := range actionPlan.ActionPlanStatusUpdates {
				statusUpdateIDs = append(statusUpdateIDs, &update.ID)
			}
			if err := iraps.repo.DeleteIdentityReviewActionPlanStatusUpdates(statusUpdateIDs, tx); err != nil {
				failedIDs = append(failedIDs, id)
				continue
			}
		}

		if err := iraps.repo.Delete(id, tx); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		successfulIDs = append(successfulIDs, id)
	}

	if len(failedIDs) > 0 {
		tx.Rollback()
		return failedIDs, nil
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	if len(successfulIDs) > 0 {
		go func() {
			if err := iraps.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.IdentityReviewActionPlan{}.TableName(),
				successfulIDs,
				deletedActionPlans,
				currentUserID,
				nil,
			); err != nil {
				utility.LogError(
					"Failed to log bulk delete activity",
					map[string]interface{}{
						"error": err.Error(),
					},
				)
			}
		}()
	}

	return failedIDs, nil
}

func (iraps *IdentityReviewActionPlanService) ForceDelete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := iraps.repo.BeginTransaction()

	var failedIDs []uint
	var deletedData []map[string]interface{}

	for _, id := range ids {
		if err := iraps.repo.ForceDelete(id, tx); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		deletedData = append(deletedData, map[string]interface{}{
			"id": id,
		})
	}

	if len(failedIDs) > 0 {
		tx.Rollback()
		return failedIDs, nil
	}

	if err := tx.Commit(); err != nil {
		return ids, err
	}

	go func() {
		if err := iraps.activityLogSvc.LogBulkActivities(
			models.DeleteAction,
			models.IdentityReviewActionPlan{}.TableName(),
			ids,
			deletedData,
			currentUserID,
			nil,
		); err != nil {
			utility.LogError(
				"Failed to log bulk force-delete activity",
				map[string]interface{}{
					"error": err.Error(),
				},
			)
		}
	}()

	return nil, nil
}

func (iraps *IdentityReviewActionPlanService) GetSummary(identityReviewUserID uint) (*dtos.IdentityReviewActionPlanSummaryDTO, error) {
	if identityReviewUserID != 0 {
		if _, err := iraps.identityReviewUserSvc.GetByID(identityReviewUserID); err != nil {
			return nil, err
		}
	}
	openCount, closedCount, err := iraps.repo.GetSummary(identityReviewUserID)
	if err != nil {
		return nil, err
	}
	return &dtos.IdentityReviewActionPlanSummaryDTO{OpenCount: openCount, ClosedCount: closedCount}, nil
}

func (iraps *IdentityReviewActionPlanService) Export(items []*dtos.IdentityReviewActionPlanListDTO, currentUserID uint) (*os.File, error) {
	userSetting, err := iraps.userSettingSvc.GetUserSettingsByUserID(currentUserID)
	if err != nil {
		return nil, err
	}

	headers := []string{
		"SL No.",
		"Action Plan",
		"Identity Review User",
		"Description",
		"Severity",
		"Responsible User",
		"Target Date",
		"Percentage",
		"Status",
	}
	if userSetting.Language == "ar" {
		headers = []string{
			"الرقم التسلسلي",
			"خطة العمل",
			"مستخدم مراجعة الهوية",
			"الوصف",
			"درجة الخطورة",
			"المستخدم المسؤول",
			"التاريخ المستهدف",
			"النسبة المئوية",
			"الحالة",
		}
	}

	rows := make([][]interface{}, len(items))
	for i, actionPlan := range items {
		responsibleUser := actionPlan.ResponsibleUserFirstNameEn + " " + actionPlan.ResponsibleUserLastNameEn
		severity := utility.GetStringValue(actionPlan.SeverityLevelTitleEn)
		status := actionPlan.ActionPlanStatusTitleEn

		if userSetting.Language == "ar" {
			responsibleUser = actionPlan.ResponsibleUserFirstNameAr + " " + actionPlan.ResponsibleUserLastNameAr
			severity = utility.GetStringValue(actionPlan.SeverityLevelTitleAr)
			status = actionPlan.ActionPlanStatusTitleAr
		}

		rows[i] = []interface{}{
			i + 1,
			actionPlan.Title,
			actionPlan.IdentityReviewUsername,
			utility.GetStringValue(actionPlan.Description),
			severity,
			responsibleUser,
			utility.GetStringValue(actionPlan.TargetDate),
			actionPlan.Percentage,
			status,
		}
	}

	return utility.GenerateExcelFile(headers, rows)
}

func (iraps *IdentityReviewActionPlanService) StatusUpdate(identityReviewUserID uint, id uint, currentUserID uint, formDTO *dtos.IdentityReviewActionPlanStatusUpdateFormDTO, tx ...utility.DatabaseTransaction) error {
	actionPlan, err := iraps.GetByID(identityReviewUserID, id)
	if err != nil {
		return err
	}
	actionPlan.ActionPlanStatusID = formDTO.ActionPlanStatusID
	actionPlan.Percentage = *formDTO.Percentage
	actionPlan.UpdatedByID = &currentUserID
	if err := iraps.repo.Update(actionPlan, tx...); err != nil {
		return err
	}
	task, err := iraps.taskSvc.GetByFieldID("identity_review_action_plan_id", id, tx...)
	if err != nil {
		return err
	}
	return iraps.taskSvc.UpdateStatus(currentUserID, task.ID, formDTO.ActionPlanStatusID, actionPlan.ReviewUserID == nil, tx...)
}

func (iraps *IdentityReviewActionPlanService) validateForm(formDTO *dtos.IdentityReviewActionPlanFormDTO) error {
	if formDTO.SeverityLevelID != nil {
		if _, err := iraps.severityLevelSvc.GetByID(*formDTO.SeverityLevelID); err != nil {
			return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity_level_id")
		}
	}
	if _, err := iraps.userSvc.GetByID(formDTO.ResponsibleUserID); err != nil {
		return utility.NewValidationError(fmt.Errorf("responsible user does not exist"), "not_exist", "responsible_user_id")
	}
	if formDTO.ReviewUserID != nil {
		if _, err := iraps.userSvc.GetByID(*formDTO.ReviewUserID); err != nil {
			return utility.NewValidationError(fmt.Errorf("review user does not exist"), "not_exist", "review_user_id")
		}
	}
	return nil
}

func (iraps *IdentityReviewActionPlanService) validateModifyDecision(formDTO *dtos.IdentityReviewActionPlanFormDTO) error {
	if formDTO.PrivilegedID == nil {
		return utility.NewValidationError(fmt.Errorf("privileged_id is required when account decision is modify"), "required", "privileged_id")
	}
	if _, err := iraps.privilegedSvc.GetByID(*formDTO.PrivilegedID); err != nil {
		return utility.NewValidationError(fmt.Errorf("privileged access does not exist"), "not_exist", "privileged_id")
	}
	return nil
}

func (iraps *IdentityReviewActionPlanService) parseDates(formDTO *dtos.IdentityReviewActionPlanFormDTO) (*dtos.IdentityReviewActionPlanParsedDates, error) {
	parsed := &dtos.IdentityReviewActionPlanParsedDates{}
	if formDTO.TargetDate == nil || *formDTO.TargetDate == "" {
		return parsed, nil
	}
	targetDate, err := time.Parse("2006-01-02", *formDTO.TargetDate)
	if err != nil {
		return nil, utility.NewValidationError(fmt.Errorf("invalid target date"), "invalid_date", "target_date")
	}
	parsed.TargetDate = &targetDate
	return parsed, nil
}
