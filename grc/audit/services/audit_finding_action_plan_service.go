package services

import (
	"errors"
	"fmt"
	"os"
	"time"

	auditFindingDTOs "jethur/grc/audit/dtos"
	auditFindingRepositories "jethur/grc/audit/repositories"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
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

type AuditFindingActionPlanServiceInterface interface {
	GetAll(findingID uint, params *auditFindingDTOs.AuditFindingActionPlanListQueryParamsDTO) (*auditFindingDTOs.AuditFindingActionPlanListResponseDTO, error)
	GetByID(findingID uint, id uint) (*models.AuditFindingActionPlan, error)
	GetFullDetailsByID(findingID uint, id uint) (*auditFindingDTOs.AuditFindingActionPlanDetailDTO, error)
	Create(findingID uint, currentUserID uint, AuditFindingActionPlanFormDTO *auditFindingDTOs.AuditFindingActionPlanFormDTO) (uint, error)
	Update(findingID uint, id uint, currentUserID uint, auditFinding *auditFindingDTOs.AuditFindingActionPlanFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	StatusUpdate(findingID uint, id uint, currentUserID uint, auditFindingActionPlanStatusUpdateFormDTO *auditFindingDTOs.AuditFindingActionPlanStatusUpdateFormDTO, tx utility.DatabaseTransaction) error

	GetAuditFindingActionPlanCountByStatus(findingId uint) ([]auditFindingDTOs.AuditFindingActionPlanStatusCountResponseDTO, error)
	ExportAuditFindingActionPlans(data []*auditFindingDTOs.AuditFindingActionPlanListDTO, currentUserID uint) (*os.File, error)
}

// Service and implementations
type AuditFindingActionPlanService struct {
	repo                auditFindingRepositories.AuditFindingActionPlanRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	auditFindingSvc     AuditFindingServiceInterface
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface
	severityLevelSvc    masterServices.SeverityLevelServiceInterface
	auditTimelineSvc    AuditTimelineServiceInterface
	taskSvc             taskServices.TaskServiceInterface
	userSvc             organizationServices.UserServiceInterface
	notificationSvc     commonServices.NotificationServiceInterface
	userSettingSvc      organizationServices.UserSettingServiceInterface
}

func NewAuditFindingActionPlanService(
	repo auditFindingRepositories.AuditFindingActionPlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditFindingSvc AuditFindingServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	severityLevelSvc masterRepositories.SeverityLevelRepositoryInterface,
	auditTimelineSvc AuditTimelineServiceInterface,
	taskSvc taskServices.TaskServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
) AuditFindingActionPlanServiceInterface {
	return &AuditFindingActionPlanService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		auditFindingSvc:     auditFindingSvc,
		actionPlanStatusSvc: actionPlanStatusSvc,
		severityLevelSvc:    severityLevelSvc,
		auditTimelineSvc:    auditTimelineSvc,
		taskSvc:             taskSvc,
		userSvc:             userSvc,
		notificationSvc:     notificationSvc,
		userSettingSvc:      userSettingSvc,
	}
}

func (afaps *AuditFindingActionPlanService) GetAll(findingID uint, params *auditFindingDTOs.AuditFindingActionPlanListQueryParamsDTO) (*auditFindingDTOs.AuditFindingActionPlanListResponseDTO, error) {
	return afaps.repo.GetAll(findingID, params)
}

func (afaps AuditFindingActionPlanService) GetByID(findingID uint, id uint) (*models.AuditFindingActionPlan, error) {
	if _, err := afaps.auditFindingSvc.GetByID(findingID); err != nil {
		return nil, err
	}

	auditFindingActionPlan, err := afaps.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return auditFindingActionPlan, nil
}

func (afaps AuditFindingActionPlanService) GetFullDetailsByID(findingID uint, id uint) (*auditFindingDTOs.AuditFindingActionPlanDetailDTO, error) {
	var auditFindingActionPlanDTO auditFindingDTOs.AuditFindingActionPlanDetailDTO

	if _, err := afaps.auditFindingSvc.GetByID(findingID); err != nil {
		return nil, err
	}

	auditFindingActionPlan, err := afaps.repo.GetFullDetailsByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(auditFindingActionPlan, &auditFindingActionPlanDTO)
	if err != nil {
		return nil, err
	}

	// todo :- also add extra fields which required in details screen
	if auditFindingActionPlan.TargetDate != nil {
		targetDate := auditFindingActionPlan.TargetDate.Format("2006-01-02")
		auditFindingActionPlanDTO.TargetDateString = &targetDate
	}

	auditFindingActionPlanDTO.CreatedAt = auditFindingActionPlan.CreatedAt
	auditFindingActionPlanDTO.UpdatedAt = auditFindingActionPlan.UpdatedAt

	return &auditFindingActionPlanDTO, nil
}

func (afaps *AuditFindingActionPlanService) Create(findingID uint, currentUserID uint, AuditFindingActionPlanFormDTO *auditFindingDTOs.AuditFindingActionPlanFormDTO) (uint, error) {
	tx := afaps.repo.BeginTransaction()
	audit, err := afaps.auditFindingSvc.GetByID(findingID)
	if err != nil {
		return 0, err
	}

	if err := afaps.validateAuditFindingActionPlanFormDTO(AuditFindingActionPlanFormDTO); err != nil {
		return 0, err
	}

	dates, err := afaps.parseAuditFindingActionPlanDates(AuditFindingActionPlanFormDTO)
	if err != nil {
		return 0, err
	}

	actionPlanStatus, err := afaps.actionPlanStatusSvc.GetByType("new")
	if err != nil {
		return 0, err
	}

	auditFindingActionPlan := &models.AuditFindingActionPlan{
		AuditFindingID:     findingID,
		Title:              AuditFindingActionPlanFormDTO.Title,
		Description:        AuditFindingActionPlanFormDTO.Description,
		SeverityLevelID:    AuditFindingActionPlanFormDTO.SeverityLevelID,
		ResponsibleUserID:  AuditFindingActionPlanFormDTO.ResponsibleUserID,
		ReviewUserID:       AuditFindingActionPlanFormDTO.ReviewUserID,
		TargetDate:         dates.TargetDate,
		Percentage:         AuditFindingActionPlanFormDTO.Percentage,
		Comments:           AuditFindingActionPlanFormDTO.Comments,
		ActionPlanStatusID: actionPlanStatus.ID,
		CreatedBYID:        currentUserID,
	}

	id, err := afaps.repo.Create(auditFindingActionPlan, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	actionPlanStatusID := actionPlanStatus.ID
	auditFindingActionPlanUpdate := &models.AuditFindingActionPlanStatusUpdate{
		AuditFindingActionPlanID: id,
		ActionPlanStatusID:       actionPlanStatusID,
		Percentage:               AuditFindingActionPlanFormDTO.Percentage,
		Comment:                  "Initial action plan status",
		CreatedByID:              currentUserID,
	}

	_, err = afaps.repo.CreateInitialAuditFindingActionPlanStatusUpdate(auditFindingActionPlanUpdate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	urlParam := fmt.Sprintf(`{"audit_id":%d,"audit_finding_id":%d, "action_plan_id":%d}`, audit.AuditID, findingID, id)
	taskDTO := commonDTOs.TaskFormDTO{
		ModuleIdentifier:         "audit",
		AuditFindingActionPlanID: &id,
		Title:                    AuditFindingActionPlanFormDTO.Title,
		Description:              AuditFindingActionPlanFormDTO.Description,
		TargetDate:               dates.TargetDate,
		URLParam:                 &urlParam,
		QueryParam:               nil,
		SeverityLevelID:          AuditFindingActionPlanFormDTO.SeverityLevelID,
		ResponsibleUserID:        AuditFindingActionPlanFormDTO.ResponsibleUserID,
	}

	if _, err := afaps.taskSvc.Create(taskDTO, currentUserID, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	go func() {
		existingAuditFindingActionPlan, err := afaps.GetFullDetailsByID(findingID, id)
		if err != nil {
			utility.LogError("Failed to fetch audit finding action plan", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		err = afaps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditFindingActionPlan{}.TableName(),
			id,
			existingAuditFindingActionPlan,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AuditFindingActionPlan{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}

		// prepare notification
		urlParam := fmt.Sprintf(`{"audit_id":%d,"audit_finding_id":%d, "action_plan_id":%d}`, audit.AuditID, findingID, id)
		queryParam := ""
		entity := "audit-finding-action-plan"
		action := "responsible-user-assign"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserID:               AuditFindingActionPlanFormDTO.ResponsibleUserID,
			ModuleIdentifier:     "audit",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "Audit Finding Action Plan Responsiblity Assigned",
			HeadingAr:            "تم تعيين مسؤولية خطة العمل الخاصة بالتدقيق",
			TitleEn:              AuditFindingActionPlanFormDTO.Title,
			TitleAr:              AuditFindingActionPlanFormDTO.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			QueryParam:           queryParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		responsibleUser, err := afaps.userSvc.GetByID(AuditFindingActionPlanFormDTO.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch responsible user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		assignedBy, err := afaps.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch assigning user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("Audit Finding Action Plan/%s: Responsibility Assigned", AuditFindingActionPlanFormDTO.Title)
		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstNameEn":                 responsibleUser.FirstNameEn,
			"LastNameEn":                  responsibleUser.LastNameEn,
			"AuditFindingReferenceNumber": audit.ReferenceNumber,
			"AuditFindingTitle":           audit.Title,
			"ActionPlanTitle":             AuditFindingActionPlanFormDTO.Title,
			"ActionPlanDescription":       AuditFindingActionPlanFormDTO.Description,
			"TargetDate":                  dates.TargetDate.Format("2006-01-02"),
			"AssignedByFirstNameEn":       assignedBy.FirstNameEn,
			"AssignedByLastNameEn":        assignedBy.LastNameEn,
			"AssignedOn":                  time.Now().Format("2006-01-02"),
			"RedirectURL":                 fmt.Sprintf("%saudit/audit-findings/%d/action-plans", appURL, findingID),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "audit/templates/email/audit_finding_action_plan_create_responsible_user_mail.html",
			EmailData:            emailData,
		}

		// send notification and mail
		err = afaps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "audit_finding_action_plan")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	return id, nil
}

func (afaps *AuditFindingActionPlanService) Update(findingID uint, id uint, currentUserID uint, AuditFindingActionPlanFormDTO *auditFindingDTOs.AuditFindingActionPlanFormDTO) error {
	tx := afaps.repo.BeginTransaction()

	audit, err := afaps.auditFindingSvc.GetByID(findingID)
	if err != nil {
		return err
	}

	oldAuditFindingActionPlan, err := afaps.repo.GetFullDetailsByID(id)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}

	existingAuditFindingActionPlan, err := afaps.GetByID(findingID, id)

	if err != nil {
		return err
	}

	oldResponsibleUserID := existingAuditFindingActionPlan.ResponsibleUserID

	if err := afaps.validateAuditFindingActionPlanFormDTO(AuditFindingActionPlanFormDTO); err != nil {
		return err
	}

	dates, err := afaps.parseAuditFindingActionPlanDates(AuditFindingActionPlanFormDTO)
	if err != nil {
		return err
	}

	// Update Audit Finding ActionPlan
	existingAuditFindingActionPlan.Title = AuditFindingActionPlanFormDTO.Title
	existingAuditFindingActionPlan.Description = AuditFindingActionPlanFormDTO.Description
	existingAuditFindingActionPlan.SeverityLevelID = AuditFindingActionPlanFormDTO.SeverityLevelID
	existingAuditFindingActionPlan.ResponsibleUserID = AuditFindingActionPlanFormDTO.ResponsibleUserID
	existingAuditFindingActionPlan.ReviewUserID = AuditFindingActionPlanFormDTO.ReviewUserID
	existingAuditFindingActionPlan.TargetDate = dates.TargetDate
	existingAuditFindingActionPlan.Percentage = AuditFindingActionPlanFormDTO.Percentage
	existingAuditFindingActionPlan.Comments = AuditFindingActionPlanFormDTO.Comments
	existingAuditFindingActionPlan.ActionPlanStatusID = AuditFindingActionPlanFormDTO.ActionPlanStatusID
	existingAuditFindingActionPlan.UpdatedBYID = &currentUserID

	err = afaps.repo.Update(existingAuditFindingActionPlan, tx)

	if err != nil {
		return err
	}

	// //update task
	urlParam := fmt.Sprintf(`{"audit_id":%d,"audit_finding_id":%d, "action_plan_id":%d}`, audit.AuditID, findingID, id)

	taskUpdateDTO := commonDTOs.TaskFormDTO{
		Title:             AuditFindingActionPlanFormDTO.Title,
		Description:       AuditFindingActionPlanFormDTO.Description,
		TargetDate:        dates.TargetDate,
		URLParam:          &urlParam,
		SeverityLevelID:   AuditFindingActionPlanFormDTO.SeverityLevelID,
		ResponsibleUserID: AuditFindingActionPlanFormDTO.ResponsibleUserID,
	}

	err = afaps.taskSvc.Update(currentUserID, taskUpdateDTO, "audit_finding_action_plan_id", id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		existingFindingActionPlan, err := afaps.repo.GetFullDetailsByID(id)
		if err != nil {
			return
		}
		// Prepare log data
		changedData := make(map[string]interface{})
		//log changes
		oldData := map[string]interface{}{
			"title":              oldAuditFindingActionPlan.Title,
			"description":        oldAuditFindingActionPlan.Description,
			"severity_level":     oldAuditFindingActionPlan.SeverityLevel.TitleEn,
			"responsible_user":   fmt.Sprintf("%s %s", oldAuditFindingActionPlan.ResponsibleUser.FirstNameEn, oldAuditFindingActionPlan.ResponsibleUser.LastNameEn),
			"review_user":        utility.GetFieldValue(oldAuditFindingActionPlan.ReviewUser),
			"target_date":        oldAuditFindingActionPlan.TargetDate,
			"percentage":         oldAuditFindingActionPlan.Percentage,
			"comments":           oldAuditFindingActionPlan.Comments,
			"action_plan_status": oldAuditFindingActionPlan.ActionPlanStatus.TitleEn,
		}

		newData := map[string]interface{}{
			"title":              existingFindingActionPlan.Title,
			"description":        existingFindingActionPlan.Description,
			"severity_level":     existingFindingActionPlan.SeverityLevel.TitleEn,
			"responsible_user":   fmt.Sprintf("%s %s", existingFindingActionPlan.ResponsibleUser.FirstNameEn, existingFindingActionPlan.ResponsibleUser.LastNameEn),
			"review_user":        utility.GetFieldValue(existingFindingActionPlan.ReviewUser),
			"target_date":        existingFindingActionPlan.TargetDate,
			"percentage":         existingFindingActionPlan.Percentage,
			"comments":           existingFindingActionPlan.Comments,
			"action_plan_status": existingFindingActionPlan.ActionPlanStatus.TitleEn,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := afaps.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.AuditFindingActionPlan{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.AuditFindingActionPlan{}.TableName(),
					"recordID":  id,
					"userID":    currentUserID,
					"error":     err.Error(),
				}

				utility.LogError("Failed to log activity", fields)
			}
		}
		// sends notification if responsible user has changed
		if oldResponsibleUserID != existingAuditFindingActionPlan.ResponsibleUserID {
			urlParam := fmt.Sprintf(`{"audit_id":%d,"audit_finding_id":%d, "action_plan_id":%d}`, audit.AuditID, findingID, id)
			queryParam := ""
			entity := "audit-finding-action-plan"

			oldResponsibleUserAction := "responsible-user-remove"
			oldResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserID:               oldResponsibleUserID,
				ModuleIdentifier:     "audit",
				NotificationType:     string(commonDTOs.InformationType),
				HeadingEn:            "Audit Finding Action Plan Responsibility Removed",
				HeadingAr:            "تمت إزالة مسؤولية خطة عمل العثور على التدقيق",
				TitleEn:              existingAuditFindingActionPlan.Title,
				TitleAr:              existingAuditFindingActionPlan.Title,
				Entity:               &entity,
				Action:               &oldResponsibleUserAction,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			newResponsibleUserAction := "responsible-user-assign"
			newResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserID:               AuditFindingActionPlanFormDTO.ResponsibleUserID,
				ModuleIdentifier:     "audit",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Audit Finding Action Plan Responsibility Assigned",
				HeadingAr:            "تم تعيين مسؤولية خطة العمل الخاصة بالتدقيق",
				TitleEn:              existingAuditFindingActionPlan.Title,
				TitleAr:              existingAuditFindingActionPlan.Title,
				Entity:               &entity,
				Action:               &newResponsibleUserAction,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			oldResponsibleUser, err := afaps.userSvc.GetByID(oldResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch old responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			newResponsibleUser, err := afaps.userSvc.GetByID(AuditFindingActionPlanFormDTO.ResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch new responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			assignedBy, err := afaps.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch assigning user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			appURL := utility.GetAppURL()

			oldResponsibleUserEmailTo := []string{oldResponsibleUser.Email}
			oldResponsibleUserSubject := fmt.Sprintf("Audit Finding Action Plan/%s: Responsibility Removed", existingAuditFindingActionPlan.Title)
			oldResponsibleUserTemplateRelativePath := "audit/templates/email/audit_finding_action_plan_remove_responsible_user_mail.html"
			oldResponsibleUserEmailData := map[string]interface{}{
				"FirstNameEn":                 oldResponsibleUser.FirstNameEn,
				"LastNameEn":                  oldResponsibleUser.LastNameEn,
				"AuditFindingReferenceNumber": audit.ReferenceNumber,
				"AuditFindingTitle":           audit.Title,
				"ActionPlanTitle":             existingAuditFindingActionPlan.Title,
				"ActionPlanDescription":       existingAuditFindingActionPlan.Description,
				"RedirectURL":                 fmt.Sprintf("%saudit/audit-findings/%d/action-plans", appURL, findingID),
			}

			oldResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              oldResponsibleUserEmailTo,
				Subject:              oldResponsibleUserSubject,
				TemplateRelativePath: oldResponsibleUserTemplateRelativePath,
				EmailData:            oldResponsibleUserEmailData,
			}

			newResponsibleUserEmailTo := []string{newResponsibleUser.Email}
			newResponsibleUserSubject := fmt.Sprintf("Audit Finding Action Plan/%s: Responsibility Assigned", existingAuditFindingActionPlan.Title)
			newResponsibleUserTemplateRelativePath := "audit/templates/email/audit_finding_action_plan_change_responsible_user_mail.html"
			newResponsibleUserEmailData := map[string]interface{}{
				"FirstNameEn":                 newResponsibleUser.FirstNameEn,
				"LastNameEn":                  newResponsibleUser.LastNameEn,
				"AuditFindingReferenceNumber": audit.ReferenceNumber,
				"AuditFindingTitle":           audit.Title,
				"ActionPlanTitle":             existingAuditFindingActionPlan.Title,
				"ActionPlanDescription":       existingAuditFindingActionPlan.Description,
				"TargetDate":                  existingAuditFindingActionPlan.TargetDate.Format("2006-01-02"),
				"AssignedByFirstNameEn":       assignedBy.FirstNameEn,
				"AssignedByLastNameEn":        assignedBy.LastNameEn,
				"AssignedOn":                  time.Now().Format("2006-01-02"),
				"RedirectURL":                 fmt.Sprintf("%saudit/audit-findings/%d/action-plans", appURL, findingID),
			}

			newResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              newResponsibleUserEmailTo,
				Subject:              newResponsibleUserSubject,
				TemplateRelativePath: newResponsibleUserTemplateRelativePath,
				EmailData:            newResponsibleUserEmailData,
			}

			err = afaps.notificationSvc.SendNotification(
				newResponsibleUserNotificationDTO, newResponsibleUserEmailNotificationDTO, "audit_finding_action_plan",
			)
			if err != nil {
				utility.LogError("Failed to send new responsible user notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}

			err = afaps.notificationSvc.SendNotification(
				oldResponsibleUserNotificationDTO, oldResponsibleUserEmailNotificationDTO, "audit_finding_action_plan",
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

func (afaps *AuditFindingActionPlanService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := afaps.repo.BeginTransaction()

	var failedIds []uint
	var successfulIds []uint
	var deletedAuditFindingActionPlan []map[string]interface{}
	var notifications []struct {
		OwnerID           *uint
		TitleEn           string
		TitleAr           string
		Description       string
		ReferenceNumber   string
		AuditFindingTitle string
	}
	for _, id := range ids {

		auditFindingActionPlan, err := afaps.repo.GetFullDetailsByID(id)

		if err != nil {
			continue
		}

		// Prepare notification data before deletion
		notifications = append(notifications, struct {
			OwnerID           *uint
			TitleEn           string
			TitleAr           string
			Description       string
			ReferenceNumber   string
			AuditFindingTitle string
		}{
			OwnerID:           &auditFindingActionPlan.ResponsibleUserID,
			TitleEn:           auditFindingActionPlan.Title,
			TitleAr:           auditFindingActionPlan.Title,
			Description:       *auditFindingActionPlan.Description,
			ReferenceNumber:   auditFindingActionPlan.AuditFinding.ReferenceNumber,
			AuditFindingTitle: auditFindingActionPlan.AuditFinding.Title,
		})

		if len(auditFindingActionPlan.ActionPlanStatusUpdates) > 0 {
			var deleteStatusUpdateIds []*uint
			for _, update := range auditFindingActionPlan.ActionPlanStatusUpdates {
				deleteStatusUpdateIds = append(deleteStatusUpdateIds, &update.ID)
			}
			if err := afaps.repo.DeleteAuditFindingActionPlanStatusUpdates(deleteStatusUpdateIds); err != nil {
				failedIds = append(failedIds, id)
				continue
			}
		}

		existingTask, err := afaps.taskSvc.GetByFieldID("audit_finding_action_plan_id", id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		taskID, err := afaps.taskSvc.Delete(existingTask.ID, currentUserID)
		if err != nil {
			failedIds = append(failedIds, taskID)
		}

		err = afaps.repo.Delete(id, tx)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedAuditFindingActionPlan = append(deletedAuditFindingActionPlan, map[string]interface{}{
			"id":            id,
			"title":         auditFindingActionPlan.Title,
			"audit_finding": auditFindingActionPlan.AuditFinding.Title,
		})
	}
	tx.Commit()

	if len(successfulIds) > 0 {

		go func() {
			err := afaps.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.AuditFinding{}.TableName(),
				successfulIds,
				deletedAuditFindingActionPlan,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}

			for _, notification := range notifications {
				urlParam := ""
				entity := "audit-finding-action-plan"
				action := "delete"

				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserID:               *notification.OwnerID,
					ModuleIdentifier:     "audit",
					NotificationType:     string(commonDTOs.DeletionType),
					HeadingEn:            "Audit Finding Action Plan Deleted",
					HeadingAr:            "تم حذف خطة عمل نتائج التدقيق",
					TitleEn:              notification.TitleEn,
					TitleAr:              notification.TitleAr,
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				ownerUser, err := afaps.userSvc.GetByID(*notification.OwnerID)
				if err != nil {
					utility.LogError("Failed to fetch owner details for notification", map[string]interface{}{
						"error": err.Error(),
					})
					continue
				}

				ownerEmailTo := []string{ownerUser.Email}
				subject := fmt.Sprintf("Audit Finding Action Plan/%s: Deleted", notification.TitleEn)
				templateRelativePath := "audit/templates/email/audit_finding_action_plan_delete_mail.html"

				emailData := map[string]interface{}{
					"FirstNameEn":                 ownerUser.FirstNameEn,
					"LastNameEn":                  ownerUser.LastNameEn,
					"AuditFindingReferenceNumber": notification.ReferenceNumber,
					"AuditFindingTitle":           notification.AuditFindingTitle,
					"ActionPlanTitle":             notification.TitleEn,
					"ActionPlanDescription":       notification.Description,
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              ownerEmailTo,
					Subject:              subject,
					TemplateRelativePath: templateRelativePath,
					EmailData:            emailData,
				}

				err = afaps.notificationSvc.SendNotification(
					dbNotificationDTO, emailNotificationDTO, "audit_finding_action_plan")
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

func (afaps *AuditFindingActionPlanService) StatusUpdate(findingID uint, id uint, currentUserID uint, auditFindingActionPlanStatusUpdateFormDTO *auditFindingDTOs.AuditFindingActionPlanStatusUpdateFormDTO, tx utility.DatabaseTransaction) error {

	existingAuditFindingActionPlan, err := afaps.GetByID(findingID, id)

	if err != nil {
		return err
	}

	existingAuditFindingActionPlan.ActionPlanStatusID = auditFindingActionPlanStatusUpdateFormDTO.ActionPlanStatusID
	existingAuditFindingActionPlan.Percentage = auditFindingActionPlanStatusUpdateFormDTO.Percentage
	existingAuditFindingActionPlan.UpdatedBYID = &currentUserID

	err = afaps.repo.Update(existingAuditFindingActionPlan, tx)

	if err != nil {
		return err
	}

	isReviewUserNull := existingAuditFindingActionPlan.ReviewUserID == nil

	// Fetch the corresponding task and update its status
	existingTask, err := afaps.taskSvc.GetByFieldID("audit_finding_action_plan_id", id)
	if err != nil {
		return err
	}
	err = afaps.taskSvc.UpdateStatus(currentUserID, existingTask.ID, auditFindingActionPlanStatusUpdateFormDTO.ActionPlanStatusID, isReviewUserNull, tx)
	if err != nil {
		return err
	}

	return nil
}

func (afaps *AuditFindingActionPlanService) validateAuditFindingActionPlanFormDTO(AuditFindingActionPlanFormDTO *auditFindingDTOs.AuditFindingActionPlanFormDTO) error {
	if AuditFindingActionPlanFormDTO.SeverityLevelID != nil {
		if _, err := afaps.severityLevelSvc.GetByID(*AuditFindingActionPlanFormDTO.SeverityLevelID); err != nil {
			return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
		}
	}

	return nil
}

func (afaps *AuditFindingActionPlanService) parseAuditFindingActionPlanDates(dto *auditFindingDTOs.AuditFindingActionPlanFormDTO) (*auditFindingDTOs.AuditFindingActionPlanParsedDates, error) {
	parsed := &auditFindingDTOs.AuditFindingActionPlanParsedDates{}

	var err error
	if *dto.TargetDate != "" {
		parsed.TargetDate, err = utility.ParseDate(*dto.TargetDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid target date"), "invalid_date", "target date")
		}
	}

	return parsed, nil
}

func (afaps *AuditFindingActionPlanService) GetAuditFindingActionPlanCountByStatus(findingId uint) ([]auditFindingDTOs.AuditFindingActionPlanStatusCountResponseDTO, error) {
	if _, err := afaps.auditFindingSvc.GetByID(findingId); err != nil {
		return nil, err
	}

	response, err := afaps.repo.GetAuditFindingActionPlanCountByStatus(findingId)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (afs *AuditFindingActionPlanService) ExportAuditFindingActionPlans(data []*auditFindingDTOs.AuditFindingActionPlanListDTO, currentUserID uint) (*os.File, error) {
	// Retrieve user settings for language preference
	userSetting, err := afs.userSettingSvc.GetUserSetting(currentUserID)
	if err != nil {
		return nil, err
	}

	var headers []string
	// Define headers based on the user's language setting
	if userSetting.Language == "ar" {
		headers = []string{
			"SL No.",
			"Reference Number",
			"Finding Title (AR)",
			"Finding Description (AR)",
			"Action Plan Title (AR)",
			"Action Plan Description (AR)",
			"Severity Level (AR)",
			"Severity Color",
			"Responsible User (AR)",
			"Review User (AR)",
			"Target Date",
			"Percentage",
			"Status (AR)",
			"Status Color",
			"Created At",
		}
	} else {
		headers = []string{
			"SL No.",
			"Reference Number",
			"Finding Title (EN)",
			"Finding Description (EN)",
			"Action Plan Title (EN)",
			"Action Plan Description (EN)",
			"Severity Level (EN)",
			"Severity Color",
			"Responsible User (EN)",
			"Review User (EN)",
			"Target Date",
			"Percentage",
			"Status (EN)",
			"Status Color",
			"Created At",
		}
	}

	// Convert data to [][]interface{} based on language preference
	rows := make([][]interface{}, len(data))
	for i, item := range data {
		if userSetting.Language == "ar" {
			rows[i] = []interface{}{
				i + 1,
				item.AuditFindingReferenceNumber,
				item.AuditFindingTitle,
				safeString(&item.AuditFindingDescription),
				item.Title,
				item.Description,
				item.SeverityLevelTitleAr,
				item.SeverityLevelColorCode,
				item.ResponsibleUserFirstNameAr + " " + item.ResponsibleUserLastNameAr,
				safeString(item.ReviewUserFirstNameAr) + " " + safeString(item.ReviewUserLastNameAr),
				item.TargetDate,
				item.Percentage,
				item.ActionPlanStatusTitleAr,
				item.ActionPlanStatusColorCode,
				item.CreatedAt,
			}
		} else {
			rows[i] = []interface{}{
				i + 1,
				item.AuditFindingReferenceNumber,
				item.AuditFindingTitle,
				safeString(&item.AuditFindingDescription),
				item.Title,
				item.Description,
				item.SeverityLevelTitleEn,
				item.SeverityLevelColorCode,
				item.ResponsibleUserFirstNameEn + " " + item.ResponsibleUserLastNameEn,
				safeString(&item.ReviewUserFirstNameEn) + " " + safeString(item.ReviewUserLastNameEn),
				item.TargetDate,
				item.Percentage,
				item.ActionPlanStatusTitleEn,
				item.ActionPlanStatusColorCode,
				item.CreatedAt,
			}
		}
	}

	// Generate the Excel file
	return utility.GenerateExcelFile(headers, rows)
}
