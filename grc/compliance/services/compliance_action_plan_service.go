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
	settingService "jethur/grc/settings/services"
	taskServices "jethur/grc/task/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ComplianceActionPlanServiceInterface interface {
	GetAll(complianceID uint, params *complianceDTOs.ComplianceActionPlanListQueryParamsDTO) (*complianceDTOs.ComplianceActionPlanListResponseDTO, error)
	GetByID(complianceID uint, id uint) (*models.ComplianceActionPlan, error)
	GetFullDetailsByID(complianceID uint, id uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceActionPlanDetailDTO, error)
	Create(complianceID uint, currentUserID uint, ComplianceActionPlanFormDTO *complianceDTOs.ComplianceActionPlanFormDTO) (uint, error)
	Update(complianceID uint, id uint, currentUserID uint, compliance *complianceDTOs.ComplianceActionPlanFormDTO) error
	Delete(complianceID uint, ids []uint, currentUserID uint) ([]uint, error)
	StatusUpdate(complianceID uint, id uint, currentUserID uint, complianceActionPlanStatusUpdateFormDTO *complianceDTOs.ComplianceActionPlanStatusUpdateFormDTO, tx utility.DatabaseTransaction) error
}

// Service and implementations
type ComplianceActionPlanService struct {
	repo                     complianceRepositories.ComplianceActionPlanRepositoryInterface
	activityLogSvc           logServices.ActivityLogServiceInterface
	complianceRegisterSvc    ComplianceRegisterServiceInterface
	complianceRequirementSvc ComplianceRequirementServiceInterface
	actionPlanStatusSvc      masterServices.ActionPlanStatusServiceInterface
	severityLevelSvc         masterServices.SeverityLevelServiceInterface
	taskSvc                  taskServices.TaskServiceInterface
	moduleSvc                settingService.ModuleServiceInterface
	userSvc                  organizationServices.UserServiceInterface
	notificationSvc          commonServices.NotificationServiceInterface
}

func NewComplianceActionPlanService(
	repo complianceRepositories.ComplianceActionPlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	complianceRegisterSvc ComplianceRegisterServiceInterface,
	complianceRequirementSvc ComplianceRequirementServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	severityLevelSvc masterRepositories.SeverityLevelRepositoryInterface,
	taskSvc taskServices.TaskServiceInterface,
	moduleSvc settingService.ModuleServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) ComplianceActionPlanServiceInterface {
	return &ComplianceActionPlanService{
		repo:                     repo,
		activityLogSvc:           activityLogSvc,
		complianceRegisterSvc:    complianceRegisterSvc,
		complianceRequirementSvc: complianceRequirementSvc,
		actionPlanStatusSvc:      actionPlanStatusSvc,
		severityLevelSvc:         severityLevelSvc,
		taskSvc:                  taskSvc,
		moduleSvc:                moduleSvc,
		userSvc:                  userSvc,
		notificationSvc:          notificationSvc,
	}
}

func (caps *ComplianceActionPlanService) GetAll(complianceID uint, params *complianceDTOs.ComplianceActionPlanListQueryParamsDTO) (*complianceDTOs.ComplianceActionPlanListResponseDTO, error) {
	if complianceID != 0 {
		if _, err := caps.complianceRegisterSvc.GetByID(complianceID); err != nil {
			return nil, err
		}
	}
	return caps.repo.GetAll(complianceID, params)
}

func (caps ComplianceActionPlanService) GetByID(complianceID uint, id uint) (*models.ComplianceActionPlan, error) {
	if _, err := caps.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}

	complianceActionPlan, err := caps.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return complianceActionPlan, nil
}

func (caps ComplianceActionPlanService) GetFullDetailsByID(complianceID uint, id uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceActionPlanDetailDTO, error) {
	var complianceActionPlanDTO complianceDTOs.ComplianceActionPlanDetailDTO

	// Check if the compliance register exists
	if _, err := caps.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}

	// Get full details of the compliance action plan
	complianceActionPlan, err := caps.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Decode the compliance action plan into the DTO
	err = mapstructure.Decode(complianceActionPlan, &complianceActionPlanDTO)
	if err != nil {
		return nil, err
	}
	if complianceActionPlan.TargetDate != nil {
		targetDate := complianceActionPlan.TargetDate.Format("2006-01-02")
		complianceActionPlanDTO.TargetDateString = &targetDate
	}

	complianceActionPlanDTO.CreatedAt = complianceActionPlan.CreatedAt
	complianceActionPlanDTO.UpdatedAt = complianceActionPlan.UpdatedAt

	// TODO: Add any additional fields required for the details screen

	return &complianceActionPlanDTO, nil
}

func (caps *ComplianceActionPlanService) Create(complianceID uint, currentUserID uint, ComplianceActionPlanFormDTO *complianceDTOs.ComplianceActionPlanFormDTO) (uint, error) {
	tx := caps.repo.BeginTransaction()

	compliance, err := caps.complianceRegisterSvc.GetByID(complianceID)
	if err != nil {
		return 0, err
	}

	if err := caps.validateComplianceActionPlanFormDTO(complianceID, ComplianceActionPlanFormDTO); err != nil {
		return 0, err
	}

	dates, err := caps.parseComplianceActionPlanDates(ComplianceActionPlanFormDTO)
	if err != nil {
		return 0, err
	}

	complianceActionPlan := &models.ComplianceActionPlan{
		ComplianceRegisterID:    complianceID,
		ComplianceRequirementID: ComplianceActionPlanFormDTO.ComplianceRequirementID,
		Title:                   ComplianceActionPlanFormDTO.Title,
		Description:             ComplianceActionPlanFormDTO.Description,
		SeverityLevelID:         ComplianceActionPlanFormDTO.SeverityLevelID,
		ResponsibleUserID:       ComplianceActionPlanFormDTO.ResponsibleUserID,
		ReviewUserID:            ComplianceActionPlanFormDTO.ReviewUserID,
		TargetDate:              dates.TargetDate,
		Percentage:              ComplianceActionPlanFormDTO.Percentage,
		Comments:                ComplianceActionPlanFormDTO.Comments,
		ActionPlanStatusID:      ComplianceActionPlanFormDTO.ActionPlanStatusID,
		CreatedByID:             currentUserID,
	}

	id, err := caps.repo.Create(complianceActionPlan, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	actionPlanStatusID := ComplianceActionPlanFormDTO.ActionPlanStatusID
	complianceActionPlanUpdate := &models.ComplianceActionPlanStatusUpdate{
		ComplianceActionPlanID: id,
		ActionPlanStatusID:     actionPlanStatusID,
		Percentage:             ComplianceActionPlanFormDTO.Percentage,
		Comment:                "Initial action plan status",
		CreatedByID:            currentUserID,
	}

	_, err = caps.repo.CreateInitialComplianceActionPlanStatusUpdate(complianceActionPlanUpdate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	urlParam := fmt.Sprintf(`{"compliance_register_id":%d, "action_plan_id":%d}`, complianceID, id)
	taskDTO := commonDTOs.TaskFormDTO{
		ModuleIdentifier:       "compliance",
		ComplianceActionPlanID: &id,
		Title:                  ComplianceActionPlanFormDTO.Title,
		Description:            ComplianceActionPlanFormDTO.Description,
		TargetDate:             dates.TargetDate,
		URLParam:               &urlParam,
		QueryParam:             nil,
		SeverityLevelID:        ComplianceActionPlanFormDTO.SeverityLevelID,
		ResponsibleUserID:      ComplianceActionPlanFormDTO.ResponsibleUserID,
	}

	if _, err := caps.taskSvc.Create(taskDTO, currentUserID, tx); err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()

	go func() {
		existingComplianceActionPlan, err := caps.GetFullDetailsByID(complianceID, id)
		if err != nil {
			utility.LogError("Failed to fetch compliance action plan", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = caps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ComplianceActionPlan{}.TableName(),
			id,
			existingComplianceActionPlan,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ComplianceActionPlan{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}

		// prepare notification
		urlParam := fmt.Sprintf(`{"compliance_register_id":%d, "action_plan_id":%d}`, complianceID, id)
		queryParam := ""
		entity := "compliance-action-plan"
		action := "responsible-user-assign"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserID:               ComplianceActionPlanFormDTO.ResponsibleUserID,
			ModuleIdentifier:     "compliance",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "Compliance Action Plan Responsibility Assigned",
			HeadingAr:            "تم تعيين مسؤولية خطة عمل الامتثال",
			TitleEn:              ComplianceActionPlanFormDTO.Title,
			TitleAr:              ComplianceActionPlanFormDTO.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			QueryParam:           queryParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		responsibleUser, err := caps.userSvc.GetByID(ComplianceActionPlanFormDTO.ResponsibleUserID)
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
		subject := fmt.Sprintf("compliance Action Plan/%s: Responsibility Assigned", ComplianceActionPlanFormDTO.Title)
		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstNameEn":               responsibleUser.FirstNameEn,
			"LastNameEn":                responsibleUser.LastNameEn,
			"ComplianceReferenceNumber": compliance.ReferenceNumber,
			"ComplianceTitle":           compliance.TitleEn,
			"ActionPlanTitle":           ComplianceActionPlanFormDTO.Title,
			"ActionPlanDescription":     ComplianceActionPlanFormDTO.Description,
			"TargetDate":                dates.TargetDate.Format("2006-01-02"),
			"AssignedByFirstNameEn":     assignedBy.FirstNameEn,
			"AssignedByLastNameEn":      assignedBy.LastNameEn,
			"AssignedOn":                time.Now().Format("2006-01-02"),
			"RedirectURL":               fmt.Sprintf("%scompliance/register/%d/action-plans?action_plan_id=%d", appURL, complianceID, id),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "compliance/templates/email/compliance_action_plan_create_responsible_user_mail.html",
			EmailData:            emailData,
		}

		// send notification and mail
		err = caps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "compliance_action_plan")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	return id, nil
}

func (caps *ComplianceActionPlanService) Update(complianceID uint, id uint, currentUserID uint, ComplianceActionPlanFormDTO *complianceDTOs.ComplianceActionPlanFormDTO) error {
	tx := caps.repo.BeginTransaction()

	compliance, err := caps.complianceRegisterSvc.GetByID(complianceID)
	if err != nil {
		return err
	}

	existingComplianceActionPlan, err := caps.GetByID(complianceID, id)
	if err != nil {
		return err
	}

	oldComplianceActionPlan, err := caps.GetFullDetailsByID(complianceID, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	oldResponsibleUserID := existingComplianceActionPlan.ResponsibleUserID

	if err := caps.validateComplianceActionPlanFormDTO(complianceID, ComplianceActionPlanFormDTO); err != nil {
		return err
	}

	dates, err := caps.parseComplianceActionPlanDates(ComplianceActionPlanFormDTO)
	if err != nil {
		return err
	}

	// Update Compliance ActionPlan
	existingComplianceActionPlan.Title = ComplianceActionPlanFormDTO.Title
	existingComplianceActionPlan.Description = ComplianceActionPlanFormDTO.Description
	existingComplianceActionPlan.SeverityLevelID = ComplianceActionPlanFormDTO.SeverityLevelID
	existingComplianceActionPlan.ResponsibleUserID = ComplianceActionPlanFormDTO.ResponsibleUserID
	existingComplianceActionPlan.ReviewUserID = ComplianceActionPlanFormDTO.ReviewUserID
	existingComplianceActionPlan.TargetDate = dates.TargetDate
	existingComplianceActionPlan.Percentage = ComplianceActionPlanFormDTO.Percentage
	existingComplianceActionPlan.Comments = ComplianceActionPlanFormDTO.Comments
	existingComplianceActionPlan.ActionPlanStatusID = ComplianceActionPlanFormDTO.ActionPlanStatusID
	existingComplianceActionPlan.UpdatedByID = &currentUserID

	err = caps.repo.Update(existingComplianceActionPlan, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	// //update task
	urlParam := fmt.Sprintf(`{"compliance_register_id":%d, "action_plan_id":%d}`, complianceID, id)

	taskUpdateDTO := commonDTOs.TaskFormDTO{
		Title:             ComplianceActionPlanFormDTO.Title,
		Description:       ComplianceActionPlanFormDTO.Description,
		TargetDate:        dates.TargetDate,
		URLParam:          &urlParam,
		SeverityLevelID:   ComplianceActionPlanFormDTO.SeverityLevelID,
		ResponsibleUserID: ComplianceActionPlanFormDTO.ResponsibleUserID,
	}

	err = caps.taskSvc.Update(currentUserID, taskUpdateDTO, "compliance_action_plan_id", id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	go func() {
		existingActionPlan, err := caps.GetFullDetailsByID(complianceID, id)
		if err != nil {
			return
		}
		//log changes
		changedData := make(map[string]interface{})
		oldData := map[string]interface{}{
			"title":              oldComplianceActionPlan.Title,
			"description":        oldComplianceActionPlan.Description,
			"severity_level":     oldComplianceActionPlan.SeverityLevel.TitleEn,
			"responsible_user":   fmt.Sprintf("%s %s", oldComplianceActionPlan.ResponsibleUser.FirstNameEN, oldComplianceActionPlan.ResponsibleUser.LastNameEN),
			"review_user":        fmt.Sprintf("%s %s", oldComplianceActionPlan.ReviewUser.FirstNameEN, oldComplianceActionPlan.ReviewUser.LastNameEN),
			"target_date":        oldComplianceActionPlan.TargetDateString,
			"percentage":         oldComplianceActionPlan.Percentage,
			"comment":            oldComplianceActionPlan.Comments,
			"action_plan_status": oldComplianceActionPlan.ActionPlanStatus.TitleEn,
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
			err := caps.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.ComplianceActionPlan{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.ComplianceActionPlan{}.TableName(),
					"recordID":  id,
					"userID":    currentUserID,
					"error":     err.Error(),
				}

				utility.LogError("Failed to log activity", fields)
			}
		}
		// sends notification if responsible user has changed
		if oldResponsibleUserID != existingComplianceActionPlan.ResponsibleUserID {
			urlParam := fmt.Sprintf(`{"compliance_register_id":%d, "action_plan_id":%d}`, complianceID, id)
			queryParam := ""
			entity := "compliance-action-plan"

			oldResponsibleUserAction := "responsible-user-remove"
			oldResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserID:               oldResponsibleUserID,
				ModuleIdentifier:     "compliance",
				NotificationType:     string(commonDTOs.InformationType),
				HeadingEn:            "Compliance Action Plan Responsibility Removed",
				HeadingAr:            "تمت إزالة مسؤولية خطة عمل الامتثال",
				TitleEn:              existingComplianceActionPlan.Title,
				TitleAr:              existingComplianceActionPlan.Title,
				Entity:               &entity,
				Action:               &oldResponsibleUserAction,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			newResponsibleUserAction := "responsible-user-assign"
			newResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserID:               ComplianceActionPlanFormDTO.ResponsibleUserID,
				ModuleIdentifier:     "compliance",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Compliance Action Plan Responsibility Assigned",
				HeadingAr:            "تم تعيين مسؤولية خطة عمل الامتثال",
				TitleEn:              existingComplianceActionPlan.Title,
				TitleAr:              existingComplianceActionPlan.Title,
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

			newResponsibleUser, err := caps.userSvc.GetByID(ComplianceActionPlanFormDTO.ResponsibleUserID)
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
			oldResponsibleUserSubject := fmt.Sprintf("Compliance Action Plan/%s: Responsibility Removed", existingComplianceActionPlan.Title)
			oldResponsibleUserTemplateRelativePath := "compliance/templates/email/compliance_action_plan_remove_responsible_user_mail.html"
			oldResponsibleUserEmailData := map[string]interface{}{
				"FirstNameEn":               oldResponsibleUser.FirstNameEn,
				"LastNameEn":                oldResponsibleUser.LastNameEn,
				"ComplianceReferenceNumber": compliance.ReferenceNumber,
				"ComplianceTitle":           compliance.TitleEn,
				"ActionPlanTitle":           existingComplianceActionPlan.Title,
				"ActionPlanDescription":     existingComplianceActionPlan.Description,
				"RedirectURL":               fmt.Sprintf("%scompliance/register/%d/action-plans?action_plan_id=%d", appURL, complianceID, id),
			}

			oldResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              oldResponsibleUserEmailTo,
				Subject:              oldResponsibleUserSubject,
				TemplateRelativePath: oldResponsibleUserTemplateRelativePath,
				EmailData:            oldResponsibleUserEmailData,
			}

			newResponsibleUserEmailTo := []string{newResponsibleUser.Email}
			newResponsibleUserSubject := fmt.Sprintf("Compliance Action Plan/%s: Responsibility Assigned", existingComplianceActionPlan.Title)
			newResponsibleUserTemplateRelativePath := "compliance/templates/email/compliance_action_plan_change_responsible_user_mail.html"
			newResponsibleUserEmailData := map[string]interface{}{
				"FirstNameEn":               newResponsibleUser.FirstNameEn,
				"LastNameEn":                newResponsibleUser.LastNameEn,
				"ComplianceReferenceNumber": compliance.ReferenceNumber,
				"ComplianceTitle":           compliance.TitleEn,
				"ActionPlanTitle":           existingComplianceActionPlan.Title,
				"ActionPlanDescription":     existingComplianceActionPlan.Description,
				"TargetDate":                existingComplianceActionPlan.TargetDate.Format("2006-01-02"),
				"AssignedByFirstNameEn":     assignedBy.FirstNameEn,
				"AssignedByLastNameEn":      assignedBy.LastNameEn,
				"AssignedOn":                time.Now().Format("2006-01-02"),
				"RedirectURL":               fmt.Sprintf("%scompliance/register/%d/action-plans?action_plan_id=%d", appURL, complianceID, id),
			}

			newResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              newResponsibleUserEmailTo,
				Subject:              newResponsibleUserSubject,
				TemplateRelativePath: newResponsibleUserTemplateRelativePath,
				EmailData:            newResponsibleUserEmailData,
			}

			err = caps.notificationSvc.SendNotification(
				newResponsibleUserNotificationDTO, newResponsibleUserEmailNotificationDTO, "compliance_action_plan",
			)
			if err != nil {
				utility.LogError("Failed to send new responsible user notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}

			err = caps.notificationSvc.SendNotification(
				oldResponsibleUserNotificationDTO, oldResponsibleUserEmailNotificationDTO, "compliance_action_plan",
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

func (caps *ComplianceActionPlanService) Delete(complianceID uint, ids []uint, currentUserID uint) ([]uint, error) {
	if _, err := caps.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}

	var failedIds []uint
	var successfulIds []uint
	var deletedComplianceActionPlan []map[string]interface{}
	var notifications []struct {
		OwnerID         *uint
		TitleEn         string
		TitleAr         string
		Description     string
		ReferenceNumber string
		ComplianceTitle string
	}
	for _, id := range ids {

		complianceActionPlan, err := caps.GetFullDetailsByID(complianceID, id)
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
			ComplianceTitle string
		}{
			OwnerID:         &complianceActionPlan.ResponsibleUserID,
			TitleEn:         complianceActionPlan.Title,
			TitleAr:         complianceActionPlan.Title,
			Description:     complianceActionPlan.Description,
			ReferenceNumber: complianceActionPlan.ComplianceRegister.ReferenceNumber,
			ComplianceTitle: complianceActionPlan.ComplianceRegister.TitleEn,
		})

		if len(complianceActionPlan.ActionPlanStatusUpdates) > 0 {
			var deleteStatusUpdateIds []*uint
			for _, update := range complianceActionPlan.ActionPlanStatusUpdates {
				deleteStatusUpdateIds = append(deleteStatusUpdateIds, &update.ID)
			}
			if err := caps.repo.DeleteComplianceActionPlanStatusUpdates(deleteStatusUpdateIds); err != nil {
				failedIds = append(failedIds, id)
				continue
			}
		}
		existingTask, err := caps.taskSvc.GetByFieldID("compliance_action_plan_id", id)
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

		deletedComplianceActionPlan = append(deletedComplianceActionPlan, map[string]interface{}{
			"id":                     id,
			"compliance_register_id": complianceID,
			"title":                  complianceActionPlan.Title,
		})

	}

	if len(failedIds) > 0 {
		return failedIds, nil
	}

	if len(successfulIds) > 0 {
		go func() {
			err := caps.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.ComplianceActionPlan{}.TableName(),
				successfulIds,
				deletedComplianceActionPlan,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}

			for _, notification := range notifications {
				urlParam := ""
				entity := "compliance-action-plan"
				action := "delete"

				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserID:               *notification.OwnerID,
					ModuleIdentifier:     "compliance",
					NotificationType:     string(commonDTOs.DeletionType),
					HeadingEn:            "Compliance Action Plan Deleted",
					HeadingAr:            "تم حذف خطة عمل الامتثال",
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
				subject := fmt.Sprintf("Compliance Action Plan/%s: Deleted", notification.TitleEn)
				templateRelativePath := "compliance/templates/email/compliance_action_plan_delete_mail.html"

				emailData := map[string]interface{}{
					"FirstNameEn":               ownerUser.FirstNameEn,
					"LastNameEn":                ownerUser.LastNameEn,
					"ComplianceReferenceNumber": notification.ReferenceNumber,
					"ComplianceTitle":           notification.ComplianceTitle,
					"ActionPlanTitle":           notification.TitleEn,
					"ActionPlanDescription":     notification.Description,
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              ownerEmailTo,
					Subject:              subject,
					TemplateRelativePath: templateRelativePath,
					EmailData:            emailData,
				}

				err = caps.notificationSvc.SendNotification(
					dbNotificationDTO, emailNotificationDTO, "compliance_action_plan")
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

func (caps *ComplianceActionPlanService) StatusUpdate(complianceID uint, id uint, currentUserID uint, complianceActionPlanStatusUpdateFormDTO *complianceDTOs.ComplianceActionPlanStatusUpdateFormDTO, tx utility.DatabaseTransaction) error {
	existingComplianceActionPlan, err := caps.GetByID(complianceID, id)
	if err != nil {
		return err
	}

	existingComplianceActionPlan.ActionPlanStatusID = complianceActionPlanStatusUpdateFormDTO.ActionPlanStatusID
	existingComplianceActionPlan.Comments = &complianceActionPlanStatusUpdateFormDTO.Comment
	existingComplianceActionPlan.Percentage = complianceActionPlanStatusUpdateFormDTO.Percentage
	existingComplianceActionPlan.UpdatedByID = &currentUserID

	err = caps.repo.Update(existingComplianceActionPlan, tx)
	if err != nil {
		return err
	}

	isReviewUserNull := existingComplianceActionPlan.ReviewUserID == nil

	// Fetch the corresponding task and update its status
	existingTask, err := caps.taskSvc.GetByFieldID("compliance_action_plan_id", id)
	if err != nil {
		return err
	}
	err = caps.taskSvc.UpdateStatus(currentUserID, existingTask.ID, complianceActionPlanStatusUpdateFormDTO.ActionPlanStatusID, isReviewUserNull, tx)
	if err != nil {
		return err
	}

	return nil
}

func (caps *ComplianceActionPlanService) validateComplianceActionPlanFormDTO(complianceID uint, ComplianceActionPlanFormDTO *complianceDTOs.ComplianceActionPlanFormDTO) error {
	if ComplianceActionPlanFormDTO.SeverityLevelID != nil {
		if _, err := caps.severityLevelSvc.GetByID(*ComplianceActionPlanFormDTO.SeverityLevelID); err != nil {
			return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
		}
	}

	if ComplianceActionPlanFormDTO.ComplianceRequirementID != nil {
		if _, err := caps.complianceRequirementSvc.GetByID(complianceID, *ComplianceActionPlanFormDTO.ComplianceRequirementID); err != nil {
			return utility.NewValidationError(fmt.Errorf("compliance requirement does not exist"), "not_exist", "compliance requirement")
		}
	}

	if _, err := caps.actionPlanStatusSvc.GetByID(ComplianceActionPlanFormDTO.ActionPlanStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("action plan status does not exist"), "not_exist", "action plan status")
	}

	if ComplianceActionPlanFormDTO.Percentage == nil {
		defaultValue := uint(0)
		ComplianceActionPlanFormDTO.Percentage = &defaultValue
	}

	return nil
}

func (caps *ComplianceActionPlanService) parseComplianceActionPlanDates(dto *complianceDTOs.ComplianceActionPlanFormDTO) (*complianceDTOs.ComplianceActionPlanParsedDates, error) {
	parsed := &complianceDTOs.ComplianceActionPlanParsedDates{}

	var err error
	if *dto.TargetDate != "" {
		parsed.TargetDate, err = utility.ParseDate(*dto.TargetDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid target date"), "invalid_date", "target date")
		}
	}

	return parsed, nil
}
