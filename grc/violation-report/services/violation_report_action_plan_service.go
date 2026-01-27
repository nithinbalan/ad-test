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
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	taskServices "jethur/grc/task/services"
	violationReportDTOs "jethur/grc/violation-report/dtos"
	violationReportRepositories "jethur/grc/violation-report/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ViolationReportActionPlanServiceInterface interface {
	GetAll(violationReportID uint, params *violationReportDTOs.ViolationReportActionPlanListQueryParamsDTO) (*violationReportDTOs.ViolationReportActionPlanListResponseDTO, error)
	GetByID(violationReportID uint, id uint) (*models.ViolationReportActionPlan, error)
	GetFullDetailsByID(violationReportID uint, id uint) (*violationReportDTOs.ViolationReportActionPlanDetailDTO, error)
	//GetUsedActionPlanMasterIDs(violationReportID uint) ([]uint, error)
	Create(violationReportID uint, currentUserID uint, violationReportActionPlanFormDTO *violationReportDTOs.ViolationReportActionPlanFormDTO) (uint, error)
	Update(violationReportID uint, id uint, currentUserID uint, violationReportActionPlanFormDTO *violationReportDTOs.ViolationReportActionPlanFormDTO) error
	Delete(violationReportID uint, ids []uint, currentUserID uint) ([]uint, error)
	ForceDelete(ids []uint, currentUserID uint) ([]uint, error)
	StatusUpdate(violationReportID uint, id uint, currentUserID uint, violationReportActionPlanStatusUpdateFormDTO *violationReportDTOs.ViolationReportActionPlanStatusUpdateFormDTO, tx utility.DatabaseTransaction) error
	GetAllByViolationReportID(violationReportID uint, tx utility.DatabaseTransaction) ([]*models.ViolationReportActionPlan, error)

	ExportViolationReportActionPlan(data []*violationReportDTOs.ViolationReportActionPlanListDTO, currentUserID uint) (*os.File, error)
}

// Service and implementations
type ViolationReportActionPlanService struct {
	repo                violationReportRepositories.ViolationReportActionPlanRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	violationReportSvc  ViolationReportServiceInterface
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface
	severityLevelSvc    masterServices.SeverityLevelServiceInterface
	userSettingSvc      organizationServices.UserSettingServiceInterface
	taskSvc             taskServices.TaskServiceInterface
	userSvc             organizationServices.UserServiceInterface
	notificationSvc     commonServices.NotificationServiceInterface
	delegateeUserSvc    commonServices.DelegateeUserserviceInterface
	actionPlanMasterSvc masterServices.ActionPlanMasterServiceInterface
}

func NewViolationReportActionPlanService(
	repo violationReportRepositories.ViolationReportActionPlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	violationReportSvc ViolationReportServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	severityLevelSvc masterServices.SeverityLevelServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
	taskSvc taskServices.TaskServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	delegateeUserSvc commonServices.DelegateeUserserviceInterface,
	actionPlanMasterSvc masterServices.ActionPlanMasterServiceInterface,
) ViolationReportActionPlanServiceInterface {
	return &ViolationReportActionPlanService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		violationReportSvc:  violationReportSvc,
		actionPlanStatusSvc: actionPlanStatusSvc,
		severityLevelSvc:    severityLevelSvc,
		userSettingSvc:      userSettingSvc,
		taskSvc:             taskSvc,
		userSvc:             userSvc,
		notificationSvc:     notificationSvc,
		delegateeUserSvc:    delegateeUserSvc,
		actionPlanMasterSvc: actionPlanMasterSvc,
	}
}

func (vraps *ViolationReportActionPlanService) GetAll(violationReportID uint, params *violationReportDTOs.ViolationReportActionPlanListQueryParamsDTO) (*violationReportDTOs.ViolationReportActionPlanListResponseDTO, error) {
	if violationReportID != 0 {
		if _, err := vraps.violationReportSvc.GetByID(violationReportID); err != nil {
			return nil, err
		}
	}
	return vraps.repo.GetAll(violationReportID, params)
}

func (vraps *ViolationReportActionPlanService) GetByID(violationReportID uint, id uint) (*models.ViolationReportActionPlan, error) {
	if _, err := vraps.violationReportSvc.GetByID(violationReportID); err != nil {
		return nil, err
	}

	violationReportActionPlan, err := vraps.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return violationReportActionPlan, nil
}

func (vraps *ViolationReportActionPlanService) GetFullDetailsByID(violationReportID uint, id uint) (*violationReportDTOs.ViolationReportActionPlanDetailDTO, error) {
	var violationReportActionPlanDTO violationReportDTOs.ViolationReportActionPlanDetailDTO

	if _, err := vraps.violationReportSvc.GetByID(violationReportID); err != nil {
		return nil, err
	}

	violationReportActionPlan, err := vraps.repo.GetFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(violationReportActionPlan, &violationReportActionPlanDTO)
	if err != nil {
		return nil, err
	}

	// todo :- also add extra fields which required in details screen
	if violationReportActionPlan.TargetDate != nil {
		targetDate := violationReportActionPlan.TargetDate.Format("2006-01-02")
		violationReportActionPlanDTO.TargetDateString = &targetDate
	}

	violationReportActionPlanDTO.CreatedAt = violationReportActionPlan.CreatedAt
	violationReportActionPlanDTO.UpdatedAt = violationReportActionPlan.UpdatedAt

	return &violationReportActionPlanDTO, nil
}

// func (vraps *ViolationReportActionPlanService) GetUsedActionPlanMasterIDs(violationReportID uint) ([]uint, error) {
// 	if _, err := vraps.violationReportSvc.GetByID(violationReportID); err != nil {
// 		return nil, err
// 	}

// 	return vraps.repo.GetUsedActionPlanMasterIDsByViolationReportID(violationReportID)
// }

func (vraps *ViolationReportActionPlanService) GetAllByViolationReportID(violationReportID uint, tx utility.DatabaseTransaction) ([]*models.ViolationReportActionPlan, error) {
	return vraps.repo.GetAllByViolationReportID(violationReportID, tx)
}

func (vraps *ViolationReportActionPlanService) Create(violationReportID uint, currentUserID uint, violationReportActionPlanFormDTO *violationReportDTOs.ViolationReportActionPlanFormDTO) (uint, error) {

	violationReport, err := vraps.violationReportSvc.GetByID(violationReportID)
	if err != nil {
		return 0, err
	}

	if err := vraps.validateViolationReportActionPlanFormDTO(violationReportActionPlanFormDTO, violationReportID); err != nil {
		return 0, err
	}

	dates, err := vraps.parseViolationReportActionPlanDates(violationReportActionPlanFormDTO)
	if err != nil {
		return 0, err
	}

	actionPlanStatus, err := vraps.actionPlanStatusSvc.GetByType("new")
	if err != nil {
		return 0, err
	}

	violationReportActionPlan := &models.ViolationReportActionPlan{
		ViolationReportID:  violationReportID,
		Title:              violationReportActionPlanFormDTO.Title,
		Description:        violationReportActionPlanFormDTO.Description,
		SeverityLevelID:    violationReportActionPlanFormDTO.SeverityLevelID,
		ResponsibleUserID:  violationReportActionPlanFormDTO.ResponsibleUserID,
		ReviewUserID:       violationReportActionPlanFormDTO.ReviewUserID,
		TargetDate:         dates.TargetDate,
		Percentage:         violationReportActionPlanFormDTO.Percentage,
		Comments:           violationReportActionPlanFormDTO.Comments,
		ActionPlanStatusID: actionPlanStatus.ID,
		CreatedBYID:        currentUserID,
	}

	if violationReportActionPlanFormDTO.ResponsibleUserID != 0 {
		delegateeID, err := vraps.delegateeUserSvc.AssignDelegateeIfPermitted(violationReportActionPlanFormDTO.ResponsibleUserID, "violation_report")
		if err != nil {
			return 0, err
		}
		if delegateeID != nil {
			violationReportActionPlan.DelegateeUserID = delegateeID
		}
	}
	tx := vraps.repo.BeginTransaction()

	id, err := vraps.repo.Create(violationReportActionPlan, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	actionPlanStatusID := actionPlanStatus.ID
	violationReportActionPlanUpdate := &models.ViolationReportActionPlanStatusUpdate{
		ViolationReportActionPlanID: id,
		ActionPlanStatusID:          actionPlanStatusID,
		Percentage:                  violationReportActionPlanFormDTO.Percentage,
		Comment:                     "Initial action plan status",
		CreatedByID:                 currentUserID,
	}
	_, err = vraps.repo.CreateInitialViolationReportActionPlanStatusUpdate(violationReportActionPlanUpdate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	urlParam := fmt.Sprintf(`{"violation_report_id":%d, "action_plan_id":%d}`, violationReportID, id)
	taskDTO := commonDTOs.TaskFormDTO{
		ModuleIdentifier:            "violation_report_register",
		ViolationReportActionPlanID: &id,
		Title:                       violationReportActionPlanFormDTO.Title,
		Description:                 violationReportActionPlanFormDTO.Description,
		TargetDate:                  dates.TargetDate,
		URLParam:                    &urlParam,
		QueryParam:                  nil,
		SeverityLevelID:             violationReportActionPlanFormDTO.SeverityLevelID,
		ResponsibleUserID:           violationReportActionPlanFormDTO.ResponsibleUserID,
	}

	if _, err := vraps.taskSvc.Create(taskDTO, currentUserID, tx); err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()

	go func() {
		existingViolationReportActionPlan, err := vraps.GetFullDetailsByID(violationReportID, id)
		if err != nil {
			utility.LogError("Failed to fetch violation report action plan", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = vraps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ViolationReportActionPlan{}.TableName(),
			id,
			existingViolationReportActionPlan,
			currentUserID,
			"violation_report", "violation_report_register", "action_plan", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ViolationReportActionPlan{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}

		// prepare notification
		urlParam := fmt.Sprintf(`{"violation_report_id":%d, "action_plan_id":%d}`, violationReportID, id)
		queryParam := ""
		entity := "violation-report-action-plan"
		action := "responsible-user-assign"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{violationReportActionPlanFormDTO.ResponsibleUserID},
			ModuleIdentifier:     "violation_report",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "Violation Report Action Plan Responsibility Assigned",
			HeadingAr:            "تم تعيين مسؤولية خطة عمل التقرير",
			TitleEn:              violationReportActionPlanFormDTO.Title,
			TitleAr:              violationReportActionPlanFormDTO.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			QueryParam:           queryParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		responsibleUser, err := vraps.userSvc.GetByID(violationReportActionPlanFormDTO.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch responsible user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		assignedBy, err := vraps.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch assigning user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("Violation Report Action Plan/%s: Responsibility Assigned", violationReportActionPlanFormDTO.Title)
		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstNameEn":                    responsibleUser.FirstNameEn,
			"LastNameEn":                     responsibleUser.LastNameEn,
			"ViolationReportReferenceNumber": violationReport.ReferenceNumber,
			"ViolationReportTitle":           violationReport.TitleEn,
			"ActionPlanTitle":                violationReportActionPlanFormDTO.Title,
			"ActionPlanDescription":          violationReportActionPlanFormDTO.Description,
			"TargetDate":                     dates.TargetDate.Format("2006-01-02"),
			"AssignedByFirstNameEn":          assignedBy.FirstNameEn,
			"AssignedByLastNameEn":           assignedBy.LastNameEn,
			"AssignedOn":                     time.Now().Format("2006-01-02"),
			"RedirectURL":                    fmt.Sprintf("%s/violation-report/action-plan?action_plan_id=%d", appURL, id),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "violation-report/templates/email/violation_report_action_plan_create_responsible_user_mail.html",
			EmailData:            emailData,
		}

		// send notification and mail
		err = vraps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "violation_report_action_plan")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	return id, nil
}

func (vraps *ViolationReportActionPlanService) Update(violationReportID uint, id uint, currentUserID uint, violationReportActionPlanFormDTO *violationReportDTOs.ViolationReportActionPlanFormDTO) error {

	violationReport, err := vraps.violationReportSvc.GetByID(violationReportID)
	if err != nil {
		return err
	}

	existingActionPlan, err := vraps.GetByID(violationReportID, id)
	if err != nil {
		return err
	}

	oldActionPlan, err := vraps.GetFullDetailsByID(violationReportID, id)
	if err != nil {
		return err
	}

	oldResponsibleUserID := existingActionPlan.ResponsibleUserID

	if err := vraps.validateViolationReportActionPlanFormDTO(violationReportActionPlanFormDTO, violationReportID); err != nil {
		return err
	}

	dates, err := vraps.parseViolationReportActionPlanDates(violationReportActionPlanFormDTO)
	if err != nil {
		return err
	}

	// Update Violation Report Action Plan
	existingActionPlan.Title = violationReportActionPlanFormDTO.Title
	existingActionPlan.Description = violationReportActionPlanFormDTO.Description
	existingActionPlan.SeverityLevelID = violationReportActionPlanFormDTO.SeverityLevelID
	existingActionPlan.ResponsibleUserID = violationReportActionPlanFormDTO.ResponsibleUserID
	existingActionPlan.ReviewUserID = violationReportActionPlanFormDTO.ReviewUserID
	existingActionPlan.TargetDate = dates.TargetDate
	existingActionPlan.Comments = violationReportActionPlanFormDTO.Comments
	existingActionPlan.UpdatedBYID = &currentUserID

	if violationReportActionPlanFormDTO.ResponsibleUserID != 0 &&
		oldResponsibleUserID != violationReportActionPlanFormDTO.ResponsibleUserID {

		existingActionPlan.DelegateeUserID = nil

		delegateeID, err := vraps.delegateeUserSvc.AssignDelegateeIfPermitted(
			violationReportActionPlanFormDTO.ResponsibleUserID,
			"violation_report",
		)
		if err != nil {
			return err
		}
		if delegateeID != nil {
			existingActionPlan.DelegateeUserID = delegateeID
		}
	}

	tx := vraps.repo.BeginTransaction()

	if err := vraps.repo.Update(existingActionPlan, tx); err != nil {
		tx.Rollback()
		return err
	}

	urlParam := fmt.Sprintf(`{"violation_report_id":%d, "action_plan_id":%d}`, violationReportID, id)

	taskUpdateDTO := commonDTOs.TaskFormDTO{
		Title:             violationReportActionPlanFormDTO.Title,
		Description:       violationReportActionPlanFormDTO.Description,
		TargetDate:        dates.TargetDate,
		URLParam:          &urlParam,
		SeverityLevelID:   violationReportActionPlanFormDTO.SeverityLevelID,
		ResponsibleUserID: violationReportActionPlanFormDTO.ResponsibleUserID,
	}

	if err := vraps.taskSvc.Update(currentUserID, taskUpdateDTO, "violation_report_action_plan_id", id, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {

		existingUpdatedActionPlan, err := vraps.GetFullDetailsByID(violationReportID, id)
		if err != nil {
			return
		}

		// Activity log
		oldData := map[string]interface{}{
			"title":              oldActionPlan.Title,
			"description":        oldActionPlan.Description,
			"severity_level":     oldActionPlan.SeverityLevel.TitleEn,
			"responsible_user":   fmt.Sprintf("%s %s", oldActionPlan.ResponsibleUser.FirstNameEN, oldActionPlan.ResponsibleUser.LastNameEN),
			"target_date":        oldActionPlan.TargetDateString,
			"percentage":         oldActionPlan.Percentage,
			"comment":            oldActionPlan.Comments,
			"action_plan_status": oldActionPlan.ActionPlanStatus.TitleEn,
		}

		if oldActionPlan.ReviewUser != nil {
			oldData["review_user"] = fmt.Sprintf("%s %s", oldActionPlan.ReviewUser.FirstNameEN, oldActionPlan.ReviewUser.LastNameEN)
		}

		newData := map[string]interface{}{
			"title":              existingUpdatedActionPlan.Title,
			"description":        existingUpdatedActionPlan.Description,
			"severity_level":     existingUpdatedActionPlan.SeverityLevel.TitleEn,
			"responsible_user":   fmt.Sprintf("%s %s", existingUpdatedActionPlan.ResponsibleUser.FirstNameEN, existingUpdatedActionPlan.ResponsibleUser.LastNameEN),
			"target_date":        existingUpdatedActionPlan.TargetDateString,
			"percentage":         existingUpdatedActionPlan.Percentage,
			"comment":            existingUpdatedActionPlan.Comments,
			"action_plan_status": existingUpdatedActionPlan.ActionPlanStatus.TitleEn,
		}

		if existingUpdatedActionPlan.ReviewUser != nil {
			newData["review_user"] = fmt.Sprintf("%s %s", existingUpdatedActionPlan.ReviewUser.FirstNameEN, existingUpdatedActionPlan.ReviewUser.LastNameEN)
		}

		changedData := utility.CompareFieldChanges(oldData, newData)

		if len(changedData) > 0 {
			_ = vraps.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.ViolationReportActionPlan{}.TableName(),
				id,
				changedData,
				currentUserID,
				"violation_report",
				"violation_report_register",
				"action_plan",
				nil,
			)
		}

		// Responsible user changed → notifications + email
		if oldResponsibleUserID != existingActionPlan.ResponsibleUserID {

			entity := "violation-report-action-plan"

			oldUserAction := "responsible-user-remove"
			newUserAction := "responsible-user-assign"

			oldNotification := &commonDTOs.DBNotificationDTO{
				UserIDs:          []uint{oldResponsibleUserID},
				ModuleIdentifier: "violation_report",
				NotificationType: string(commonDTOs.InformationType),
				HeadingEn:        "Violation Report Action Plan Responsibility Removed",
				HeadingAr:        "تمت إزالة مسؤولية خطة عمل التقرير",
				TitleEn:          existingActionPlan.Title,
				TitleAr:          existingActionPlan.Title,
				Entity:           &entity,
				Action:           &oldUserAction,
				URLParam:         urlParam,
				CreatedByID:      &currentUserID,
			}

			newNotification := &commonDTOs.DBNotificationDTO{
				UserIDs:          []uint{violationReportActionPlanFormDTO.ResponsibleUserID},
				ModuleIdentifier: "violation_report",
				NotificationType: string(commonDTOs.AssignmentType),
				HeadingEn:        "Violation Report Action Plan Responsibility Assigned",
				HeadingAr:        "تم تعيين مسؤولية خطة عمل التقرير",
				TitleEn:          existingActionPlan.Title,
				TitleAr:          existingActionPlan.Title,
				Entity:           &entity,
				Action:           &newUserAction,
				URLParam:         urlParam,
				CreatedByID:      &currentUserID,
			}

			oldResponsibleUser, err := vraps.userSvc.GetByID(oldResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch old responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			newResponsibleUser, err := vraps.userSvc.GetByID(violationReportActionPlanFormDTO.ResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch new responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			assignedBy, err := vraps.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch assigning user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			appURL := utility.GetAppURL()

			oldEmail := &commonDTOs.EmailNotificationDTO{
				EmailTo:              []string{oldResponsibleUser.Email},
				Subject:              fmt.Sprintf("Violation Report Action Plan/%s: Responsibility Removed", existingActionPlan.Title),
				TemplateRelativePath: "violation-report/templates/email/violation_report_action_plan_remove_responsible_user_mail.html",
				EmailData: map[string]interface{}{
					"FirstNameEn":                    oldResponsibleUser.FirstNameEn,
					"LastNameEn":                     oldResponsibleUser.LastNameEn,
					"ViolationReportReferenceNumber": violationReport.ReferenceNumber,
					"ViolationReportTitle":           violationReport.TitleEn,
					"ActionPlanTitle":                existingActionPlan.Title,
					"ActionPlanDescription":          existingActionPlan.Description,
					"RedirectURL":                    fmt.Sprintf("%s/violation-report/action-plan?action_plan_id=%d", appURL, id),
				},
			}

			newEmail := &commonDTOs.EmailNotificationDTO{
				EmailTo:              []string{newResponsibleUser.Email},
				Subject:              fmt.Sprintf("Violation Report Action Plan/%s: Responsibility Assigned", existingActionPlan.Title),
				TemplateRelativePath: "violation-report/templates/email/violation_report_action_plan_change_responsible_user_mail.html",
				EmailData: map[string]interface{}{
					"FirstNameEn":                    newResponsibleUser.FirstNameEn,
					"LastNameEn":                     newResponsibleUser.LastNameEn,
					"ViolationReportReferenceNumber": violationReport.ReferenceNumber,
					"ViolationReportTitle":           violationReport.TitleEn,
					"ActionPlanTitle":                existingActionPlan.Title,
					"ActionPlanDescription":          existingActionPlan.Description,
					"TargetDate":                     existingActionPlan.TargetDate.Format("2006-01-02"),
					"AssignedByFirstNameEn":          assignedBy.FirstNameEn,
					"AssignedByLastNameEn":           assignedBy.LastNameEn,
					"AssignedOn":                     time.Now().Format("2006-01-02"),
					"RedirectURL":                    fmt.Sprintf("%s/violation-report/action-plan?action_plan_id=%d", appURL, id),
				},
			}

			_ = vraps.notificationSvc.SendNotification(
				newNotification, newEmail, "violation_report_action_plan",
			)
			_ = vraps.notificationSvc.SendNotification(
				oldNotification, oldEmail, "violation_report_action_plan",
			)
		}
	}()

	return nil
}

func (vraps *ViolationReportActionPlanService) Delete(violationReportID uint, ids []uint, currentUserID uint) ([]uint, error) {

	if _, err := vraps.violationReportSvc.GetByID(violationReportID); err != nil {
		return nil, err
	}

	tx := vraps.repo.BeginTransaction()

	var failedIds []uint
	var successfulIds []uint
	var deletedActionPlans []map[string]interface{}
	var notifications []struct {
		OwnerID         *uint
		TitleEn         string
		TitleAr         string
		Description     *string
		ReferenceNumber string
		ViolationTitle  string
	}

	for _, id := range ids {

		actionPlan, err := vraps.repo.GetFullDetailsByID(id)
		if err != nil {
			continue
		}

		description := ""
		if actionPlan.Description != nil {
			description = *actionPlan.Description
		}

		deletedActionPlans = append(deletedActionPlans, map[string]interface{}{
			"id":    id,
			"title": actionPlan.Title,
		})

		// Prepare notification data before deletion
		notifications = append(notifications, struct {
			OwnerID         *uint
			TitleEn         string
			TitleAr         string
			Description     *string
			ReferenceNumber string
			ViolationTitle  string
		}{
			OwnerID:         &actionPlan.ResponsibleUserID,
			TitleEn:         actionPlan.Title,
			TitleAr:         actionPlan.Title,
			Description:     &description,
			ReferenceNumber: actionPlan.ViolationReport.ReferenceNumber,
			ViolationTitle:  actionPlan.ViolationReport.TitleEn,
		})

		existingTask, err := vraps.taskSvc.GetByFieldID(
			"violation_report_action_plan_id",
			id,
		)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		taskID, err := vraps.taskSvc.Delete(existingTask.ID, currentUserID)
		if err != nil {
			failedIds = append(failedIds, taskID)
		}

		if len(actionPlan.ActionPlanStatusUpdates) > 0 {
			var deleteStatusUpdateIds []*uint
			for _, update := range actionPlan.ActionPlanStatusUpdates {
				deleteStatusUpdateIds = append(deleteStatusUpdateIds, &update.ID)
			}
			if err := vraps.repo.DeleteViolationReportActionPlanStatusUpdates(
				deleteStatusUpdateIds,
			); err != nil {
				failedIds = append(failedIds, id)
				continue
			}
		}

		if err := vraps.repo.Delete(id, tx); err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}

	tx.Commit()

	if len(successfulIds) > 0 {
		go func() {

			err := vraps.activityLogSvc.LogBulkActivities(models.DeleteAction, models.ViolationReportActionPlan{}.TableName(), successfulIds, deletedActionPlans, currentUserID, nil)
			if err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}

			for _, notification := range notifications {

				entity := "violation-report-action-plan"
				action := "delete"

				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserIDs:              []uint{*notification.OwnerID},
					ModuleIdentifier:     "violation_report",
					NotificationType:     string(commonDTOs.DeletionType),
					HeadingEn:            "Violation Report Action Plan Deleted",
					HeadingAr:            "تم حذف خطة عمل التقرير",
					TitleEn:              notification.TitleEn,
					TitleAr:              notification.TitleAr,
					Entity:               &entity,
					Action:               &action,
					URLParam:             "",
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				ownerUser, err := vraps.userSvc.GetByID(*notification.OwnerID)
				if err != nil {
					utility.LogError("Failed to fetch owner details for notification", map[string]interface{}{
						"error": err.Error(),
					})
					continue
				}

				subject := fmt.Sprintf(
					"Violation Report Action Plan/%s: Deleted",
					notification.TitleEn,
				)

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              []string{ownerUser.Email},
					Subject:              subject,
					TemplateRelativePath: "violation-report/templates/email/violation_report_action_plan_delete_mail.html",
					EmailData: map[string]interface{}{
						"FirstNameEn":                    ownerUser.FirstNameEn,
						"LastNameEn":                     ownerUser.LastNameEn,
						"ViolationReportReferenceNumber": notification.ReferenceNumber,
						"ViolationReportTitle":           notification.ViolationTitle,
						"ActionPlanTitle":                notification.TitleEn,
						"ActionPlanDescription":          notification.Description,
					},
				}

				err = vraps.notificationSvc.SendNotification(
					dbNotificationDTO,
					emailNotificationDTO,
					"violation_report_action_plan",
				)
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

func (vraps *ViolationReportActionPlanService) StatusUpdate(violationReportID uint, id uint, currentUserID uint, violationReportActionPlanStatusUpdateFormDTO *violationReportDTOs.ViolationReportActionPlanStatusUpdateFormDTO, tx utility.DatabaseTransaction) error {
	existingViolationReportActionPlan, err := vraps.GetByID(violationReportID, id)
	if err != nil {
		return err
	}

	existingViolationReportActionPlan.ActionPlanStatusID = violationReportActionPlanStatusUpdateFormDTO.ActionPlanStatusID
	existingViolationReportActionPlan.Percentage = violationReportActionPlanStatusUpdateFormDTO.Percentage
	existingViolationReportActionPlan.UpdatedBYID = &currentUserID

	err = vraps.repo.Update(existingViolationReportActionPlan, tx)
	if err != nil {
		return err
	}

	isReviewUserNull := existingViolationReportActionPlan.ReviewUserID == nil

	// Fetch the corresponding task and update its status
	existingTask, err := vraps.taskSvc.GetByFieldID("violation_report_action_plan_id", id)
	if err != nil {
		return err
	}
	err = vraps.taskSvc.UpdateStatus(currentUserID, existingTask.ID, violationReportActionPlanStatusUpdateFormDTO.ActionPlanStatusID, isReviewUserNull, tx)
	if err != nil {
		return err
	}

	return nil
}

func (vraps *ViolationReportActionPlanService) validateViolationReportActionPlanFormDTO(violationReportActionPlanFormDTO *violationReportDTOs.ViolationReportActionPlanFormDTO, violationReportID uint) error {
	if violationReportActionPlanFormDTO.SeverityLevelID != nil {
		if _, err := vraps.severityLevelSvc.GetByID(*violationReportActionPlanFormDTO.SeverityLevelID); err != nil {
			return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
		}
	}

	return nil
}

func (vraps *ViolationReportActionPlanService) parseViolationReportActionPlanDates(dto *violationReportDTOs.ViolationReportActionPlanFormDTO) (*violationReportDTOs.ViolationReportActionPlanParsedDates, error) {
	parsed := &violationReportDTOs.ViolationReportActionPlanParsedDates{}

	var err error
	if *dto.TargetDate != "" {
		parsed.TargetDate, err = utility.ParseDate(*dto.TargetDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid target date"), "invalid_date", "target date")
		}
	}

	return parsed, nil
}

func (vraps *ViolationReportActionPlanService) ExportViolationReportActionPlan(data []*violationReportDTOs.ViolationReportActionPlanListDTO, currentUserID uint) (*os.File, error) {
	// Fetch user setting for language preference
	userSetting, err := vraps.userSettingSvc.GetUserSettingsByUserID(currentUserID)
	if err != nil {
		return nil, err
	}

	// Define headers based on user language setting
	var headers []string
	if userSetting.Language == "ar" {
		headers = []string{
			"SL No.",
			"Action Plan",
			"Violation",
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
			"Violation",
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
			actionPlan.ViolationReportTitleEn,
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

func (vraps *ViolationReportActionPlanService) ForceDelete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := vraps.repo.BeginTransaction()
	var failedIds []uint
	var deletedData []map[string]interface{}

	for _, id := range ids {
		if err := vraps.repo.ForceDelete(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		deletedData = append(deletedData, map[string]interface{}{
			"id": id,
		})
	}

	if len(failedIds) > 0 {
		return failedIds, nil
	}

	if err := tx.Commit(); err != nil {
		return ids, err
	}

	go func() {
		if err := vraps.activityLogSvc.LogBulkActivities(
			models.DeleteAction,
			models.ViolationReportActionPlan{}.TableName(),
			ids,
			deletedData,
			currentUserID,
			nil,
		); err != nil {
			fmt.Printf("Failed to log bulk delete activity: %v\n", err)
		}
	}()

	return nil, nil
}
