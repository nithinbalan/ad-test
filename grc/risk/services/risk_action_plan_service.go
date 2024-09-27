package services

import (
	"errors"
	"fmt"
	"os"
	"time"

	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	masterRepositories "jethur/grc/masters/repositories"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	riskDTOs "jethur/grc/risk/dtos"
	riskRepositories "jethur/grc/risk/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type RiskActionPlanServiceInterface interface {
	GetAll(riskID uint, params *riskDTOs.RiskActionPlanListQueryParamsDTO) (*riskDTOs.RiskActionPlanListResponseDTO, error)
	GetByID(riskID uint, id uint) (*models.RiskActionPlan, error)
	GetFullDetailsByID(riskID uint, id uint) (*riskDTOs.RiskActionPlanDetailDTO, error)
	Create(riskID uint, currentUserID uint, RiskActionPlanFormDTO *riskDTOs.RiskActionPlanFormDTO) (uint, error)
	Update(riskID uint, id uint, currentUserID uint, risk *riskDTOs.RiskActionPlanFormDTO) error
	Delete(riskID uint, ids []uint, currentUserID uint) ([]uint, error)
	StatusUpdate(riskID uint, id uint, currentUserID uint, ActionPlanStatusID uint, tx utility.DatabaseTransaction) error

	ExportRiskActionPlan(data []*riskDTOs.RiskActionPlanListDTO, currentUserID uint) (*os.File, error)
}

// Service and implementations
type RiskActionPlanService struct {
	repo                riskRepositories.RiskActionPlanRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	riskRegisterSvc     RiskRegisterServiceInterface
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface
	severityLevelSvc    masterServices.SeverityLevelServiceInterface
	userSettingSvc      organizationServices.UserSettingServiceInterface
	userSvc             organizationServices.UserServiceInterface
	notificationSvc     commonServices.NotificationServiceInterface
}

func NewRiskActionPlanService(
	repo riskRepositories.RiskActionPlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	riskRegisterSvc RiskRegisterServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	severityLevelSvc masterRepositories.SeverityLevelRepositoryInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) RiskActionPlanServiceInterface {
	return &RiskActionPlanService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		riskRegisterSvc:     riskRegisterSvc,
		actionPlanStatusSvc: actionPlanStatusSvc,
		severityLevelSvc:    severityLevelSvc,
		userSettingSvc:      userSettingSvc,
		userSvc:             userSvc,
		notificationSvc:     notificationSvc,
	}
}

func (raps *RiskActionPlanService) GetAll(riskID uint, params *riskDTOs.RiskActionPlanListQueryParamsDTO) (*riskDTOs.RiskActionPlanListResponseDTO, error) {
	if riskID != 0 {
		if _, err := raps.riskRegisterSvc.GetByID(riskID); err != nil {
			return nil, err
		}
	}
	return raps.repo.GetAll(riskID, params)
}

func (raps RiskActionPlanService) GetByID(riskID uint, id uint) (*models.RiskActionPlan, error) {
	if _, err := raps.riskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}

	riskActionPlan, err := raps.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return riskActionPlan, nil
}

func (raps RiskActionPlanService) GetFullDetailsByID(riskID uint, id uint) (*riskDTOs.RiskActionPlanDetailDTO, error) {
	var riskActionPlanDTO riskDTOs.RiskActionPlanDetailDTO

	if _, err := raps.riskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}

	riskActionPlan, err := raps.repo.GetFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(riskActionPlan, &riskActionPlanDTO)
	if err != nil {
		return nil, err
	}

	// todo :- also add extra fields which required in details screen

	return &riskActionPlanDTO, nil
}

func (raps *RiskActionPlanService) Create(riskID uint, currentUserID uint, RiskActionPlanFormDTO *riskDTOs.RiskActionPlanFormDTO) (uint, error) {
	tx := raps.repo.BeginTransaction()

	risk, err := raps.riskRegisterSvc.GetByID(riskID)
	if err != nil {
		return 0, err
	}

	if err := raps.validateRiskActionPlanFormDTO(RiskActionPlanFormDTO); err != nil {
		return 0, err
	}

	dates, err := raps.parseRiskActionPlanDates(RiskActionPlanFormDTO)
	if err != nil {
		return 0, err
	}

	riskActionPlan := &models.RiskActionPlan{
		RiskRegisterID:     riskID,
		Title:              RiskActionPlanFormDTO.Title,
		Description:        RiskActionPlanFormDTO.Description,
		SeverityLevelID:    RiskActionPlanFormDTO.SeverityLevelID,
		ResponsibleUserID:  RiskActionPlanFormDTO.ResponsibleUserID,
		TargetDate:         dates.TargetDate,
		Percentage:         RiskActionPlanFormDTO.Percentage,
		Comments:           RiskActionPlanFormDTO.Comments,
		ActionPlanStatusID: RiskActionPlanFormDTO.ActionPlanStatusID,
		CreatedBYID:        currentUserID,
	}

	id, err := raps.repo.Create(riskActionPlan)
	if err != nil {
		return 0, err
	}
	actionPlanStatusID := RiskActionPlanFormDTO.ActionPlanStatusID
	riskActionPlanUpdate := &models.RiskActionPlanStatusUpdate{
		RiskActionPlanID:   id,
		ActionPlanStatusID: actionPlanStatusID,
		Comment:            "Initial action plan status",
		CreatedByID:        currentUserID,
	}
	_, err = raps.repo.CreateInitialRiskActionPlanStatusUpdate(riskActionPlanUpdate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()

	go func() {
		err := raps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.RiskActionPlan{}.TableName(),
			id,
			RiskActionPlanFormDTO,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.RiskActionPlan{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
		// prepare notification
		urlParam := fmt.Sprintf(`{"id":%d}`, id)
		entity := "risk-action-plan"
		action := "responsible-user-assign"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserID:               RiskActionPlanFormDTO.ResponsibleUserID,
			ModuleIdentifier:     "risk",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "Risk Action Plan Assignment",
			HeadingAr:            "تعيين خطة عمل المخاطر",
			TitleEn:              RiskActionPlanFormDTO.Title,
			TitleAr:              RiskActionPlanFormDTO.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		responsibleUser, err := raps.userSvc.GetByID(RiskActionPlanFormDTO.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch responsible user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		assignedBy, err := raps.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch assigning user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("Risk Action Plan/%s: Responsibility Assigned", RiskActionPlanFormDTO.Title)
		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstNameEn":           responsibleUser.FirstNameEn,
			"LastNameEn":            responsibleUser.LastNameEn,
			"RiskReferenceNumber":   risk.ReferenceNumber,
			"RiskTitle":             risk.Title,
			"ActionPlanTitle":       RiskActionPlanFormDTO.Title,
			"ActionPlanDescription": RiskActionPlanFormDTO.Description,
			"TargetDate":            dates.TargetDate.Format("2006-01-02"),
			"AssignedByFirstNameEn": assignedBy.FirstNameEn,
			"AssignedByLastNameEn":  assignedBy.LastNameEn,
			"AssignedOn":            time.Now().Format("2006-01-02"),
			"RedirectURL":           fmt.Sprintf("%s/risk-action-plan/%d", appURL, id),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "risk/templates/email/risk_action_plan_create_responsible_user_mail.html",
			EmailData:            emailData,
		}

		// send notification and mail
		err = raps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "risk_action_plan")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	return id, nil
}

func (raps *RiskActionPlanService) Update(riskID uint, id uint, currentUserID uint, RiskActionPlanFormDTO *riskDTOs.RiskActionPlanFormDTO) error {
	tx := raps.repo.BeginTransaction()

	risk, err := raps.riskRegisterSvc.GetByID(riskID)
	if err != nil {
		return err
	}

	existingRiskActionPlan, err := raps.GetByID(riskID, id)
	if err != nil {
		return err
	}

	oldResponsibleUserID := existingRiskActionPlan.ResponsibleUserID

	if err := raps.validateRiskActionPlanFormDTO(RiskActionPlanFormDTO); err != nil {
		return err
	}

	dates, err := raps.parseRiskActionPlanDates(RiskActionPlanFormDTO)
	if err != nil {
		return err
	}

	// Update Risk ActionPlan
	existingRiskActionPlan.Title = RiskActionPlanFormDTO.Title
	existingRiskActionPlan.Description = RiskActionPlanFormDTO.Description
	existingRiskActionPlan.SeverityLevelID = RiskActionPlanFormDTO.SeverityLevelID
	existingRiskActionPlan.ResponsibleUserID = RiskActionPlanFormDTO.ResponsibleUserID
	existingRiskActionPlan.TargetDate = dates.TargetDate
	existingRiskActionPlan.Percentage = RiskActionPlanFormDTO.Percentage
	existingRiskActionPlan.Comments = RiskActionPlanFormDTO.Comments
	existingRiskActionPlan.ActionPlanStatusID = RiskActionPlanFormDTO.ActionPlanStatusID
	existingRiskActionPlan.UpdatedBYID = &currentUserID

	err = raps.repo.Update(existingRiskActionPlan, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	go func() {
		// sends notification if responsible user has changed
		if oldResponsibleUserID != existingRiskActionPlan.ResponsibleUserID {
			urlParam := fmt.Sprintf(`{"id:%d"}`, id)
			entity := "risk-action-plan"

			oldResponsibleUserAction := "responsible-user-remove"
			oldResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserID:               oldResponsibleUserID,
				ModuleIdentifier:     "risk",
				NotificationType:     string(commonDTOs.InformationType),
				HeadingEn:            "Risk Action Plan Removal",
				HeadingAr:            " إزالة خطة عمل المخاطر",
				TitleEn:              existingRiskActionPlan.Title,
				TitleAr:              existingRiskActionPlan.Title,
				Entity:               &entity,
				Action:               &oldResponsibleUserAction,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			newResponsibleUserAction := "responsible-user-assign"
			newResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserID:               RiskActionPlanFormDTO.ResponsibleUserID,
				ModuleIdentifier:     "risk",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Risk Action Plan Assignment",
				HeadingAr:            "تعيين خطة عمل المخاطر",
				TitleEn:              existingRiskActionPlan.Title,
				Entity:               &entity,
				Action:               &newResponsibleUserAction,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			oldResponsibleUser, err := raps.userSvc.GetByID(oldResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch old responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			newResponsibleUser, err := raps.userSvc.GetByID(RiskActionPlanFormDTO.ResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch new responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			assignedBy, err := raps.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch assigning user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			appURL := utility.GetAppURL()

			oldResponsibleUserEmailTo := []string{oldResponsibleUser.Email}
			oldResponsibleUserSubject := fmt.Sprintf("Risk Action Plan/%s: Responsibility Removed", existingRiskActionPlan.Title)
			oldResponsibleUserTemplateRelativePath := "risk/templates/email/risk_action_plan_remove_responsible_user_mail.html"
			oldResponsibleUserEmailData := map[string]interface{}{
				"FirstNameEn":           oldResponsibleUser.FirstNameEn,
				"LastNameEn":            oldResponsibleUser.LastNameEn,
				"RiskReferenceNumber":       risk.ReferenceNumber,
				"RiskTitle":             risk.Title,
				"ActionPlanTitle":       existingRiskActionPlan.Title,
				"ActionPlanDescription": existingRiskActionPlan.Description,
				"RedirectURL":           fmt.Sprintf("%s/risk-action-plan/%d", appURL, id),
			}

			oldResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              oldResponsibleUserEmailTo,
				Subject:              oldResponsibleUserSubject,
				TemplateRelativePath: oldResponsibleUserTemplateRelativePath,
				EmailData:            oldResponsibleUserEmailData,
			}

			newResponsibleUserEmailTo := []string{newResponsibleUser.Email}
			newResponsibleUserSubject := fmt.Sprintf("Risk Action Plan/%s: Responsibility Assigned", existingRiskActionPlan.Title)
			newResponsibleUserTemplateRelativePath := "risk/templates/email/risk_action_plan_change_responsible_user_mail.html"
			newResponsibleUserEmailData := map[string]interface{}{
				"FirstNameEn":           newResponsibleUser.FirstNameEn,
				"LastNameEn":            newResponsibleUser.LastNameEn,
				"RiskReferenceNumber":   risk.ReferenceNumber,
				"RiskTitle":             risk.Title,
				"ActionPlanTitle":       existingRiskActionPlan.Description,
				"ActionPlanDescription": existingRiskActionPlan.Description,
				"TargetDate":            existingRiskActionPlan.TargetDate.Format("2006-01-02"),
				"AssignedByFirstNameEn": assignedBy.FirstNameEn,
				"AssignedByLastNameEn":  assignedBy.LastNameEn,
				"AssignedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           fmt.Sprintf("%s/risk-action-plan/%d", appURL, id),
			}

			newResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              newResponsibleUserEmailTo,
				Subject:              newResponsibleUserSubject,
				TemplateRelativePath: newResponsibleUserTemplateRelativePath,
				EmailData:            newResponsibleUserEmailData,
			}

			err = raps.notificationSvc.SendNotification(
				newResponsibleUserNotificationDTO, newResponsibleUserEmailNotificationDTO, "risk_action_plan",
			)
			if err != nil {
				utility.LogError("Failed to send new responsible user notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}

			err = raps.notificationSvc.SendNotification(
				oldResponsibleUserNotificationDTO, oldResponsibleUserEmailNotificationDTO, "risk_action_plan",
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

func (raps *RiskActionPlanService) Delete(riskID uint, ids []uint, currentUserID uint) ([]uint, error) {
	if _, err := raps.riskRegisterSvc.GetByID(riskID); err != nil {
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
		RiskTitle       string
	}

	for _, id := range ids {

		riskActionPlan, err := raps.repo.GetFullDetailsByID(id)
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
			RiskTitle       string
		}{
			OwnerID:         &riskActionPlan.ResponsibleUserID,
			TitleEn:         riskActionPlan.Title,
			TitleAr:         riskActionPlan.Title,
			Description:     *riskActionPlan.Description,
			ReferenceNumber: riskActionPlan.RiskRegister.ReferenceNumber,
			RiskTitle:       riskActionPlan.RiskRegister.Title,
		})

		if len(riskActionPlan.ActionPlanStatusUpdates) > 0 {
			var deleteStatusUpdateIds []*uint
			for _, update := range riskActionPlan.ActionPlanStatusUpdates {
				deleteStatusUpdateIds = append(deleteStatusUpdateIds, &update.ID)
			}
			if err := raps.repo.DeleteRiskActionPlanStatusUpdates(deleteStatusUpdateIds); err != nil {
				failedIds = append(failedIds, id)
				continue
			}
		}

		err = raps.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {
		go func() {
			err := raps.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.RiskRegister{}.TableName(),
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
				entity := "risk-action-plan"
				action := "delete"

				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserID:               *notification.OwnerID,
					ModuleIdentifier:     "risk",
					NotificationType:     string(commonDTOs.DeletionType),
					HeadingEn:            "Risk Action Plan Deletion",
					HeadingAr:            "حذف خطة عمل المخاطر",
					TitleEn:              notification.TitleEn,
					TitleAr:              notification.TitleAr,
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				ownerUser, err := raps.userSvc.GetByID(*notification.OwnerID)
				if err != nil {
					utility.LogError("Failed to fetch owner details for notification", map[string]interface{}{
						"error": err.Error(),
					})
					continue
				}

				ownerEmailTo := []string{ownerUser.Email}
				subject := fmt.Sprintf("Risk Action Plan/%s: Deleted", notification.TitleEn)
				templateRelativePath := "risk/templates/email/risk_action_plan_delete_mail.html"

				emailData := map[string]interface{}{
					"FirstNameEn":           ownerUser.FirstNameEn,
					"LastNameEn":            ownerUser.LastNameEn,
					"RiskReferenceNumber":   notification.ReferenceNumber,
					"RiskTitle":             notification.RiskTitle,
					"ActionPlanTitle":       notification.TitleEn,
					"ActionPlanDescription": notification.Description,
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              ownerEmailTo,
					Subject:              subject,
					TemplateRelativePath: templateRelativePath,
					EmailData:            emailData,
				}

				err = raps.notificationSvc.SendNotification(
					dbNotificationDTO, emailNotificationDTO, "risk_action_plan")
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

func (raps *RiskActionPlanService) StatusUpdate(riskID uint, id uint, currentUserID uint, ActionPlanStatusID uint, tx utility.DatabaseTransaction) error {
	existingRiskActionPlan, err := raps.GetByID(riskID, id)
	if err != nil {
		return err
	}

	existingRiskActionPlan.ActionPlanStatusID = ActionPlanStatusID
	existingRiskActionPlan.UpdatedBYID = &currentUserID

	err = raps.repo.Update(existingRiskActionPlan, tx)
	if err != nil {
		return err
	}

	return nil
}

func (raps *RiskActionPlanService) validateRiskActionPlanFormDTO(RiskActionPlanFormDTO *riskDTOs.RiskActionPlanFormDTO) error {
	if RiskActionPlanFormDTO.SeverityLevelID != nil {
		if _, err := raps.severityLevelSvc.GetByID(*RiskActionPlanFormDTO.SeverityLevelID); err != nil {
			return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
		}
	}

	if _, err := raps.actionPlanStatusSvc.GetByID(RiskActionPlanFormDTO.ActionPlanStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("action plan status does not exist"), "not_exist", "action plan status")
	}

	return nil
}

func (raps *RiskActionPlanService) parseRiskActionPlanDates(dto *riskDTOs.RiskActionPlanFormDTO) (*riskDTOs.RiskActionPlanParsedDates, error) {
	parsed := &riskDTOs.RiskActionPlanParsedDates{}

	var err error
	if *dto.TargetDate != "" {
		parsed.TargetDate, err = utility.ParseDate(*dto.TargetDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid target date"), "invalid_date", "target date")
		}
	}

	return parsed, nil
}

func (raps *RiskActionPlanService) ExportRiskActionPlan(data []*riskDTOs.RiskActionPlanListDTO, currentUserID uint) (*os.File, error) {
	// Fetch user setting for language preference
	userSetting, err := raps.userSettingSvc.GetUserSettingsByUserID(currentUserID)
	if err != nil {
		return nil, err
	}

	// Define headers based on user language setting
	var headers []string
	if userSetting.Language == "ar" {
		headers = []string{
			"SL No.",
			"Action Plan",
			"Risk",
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
			actionPlan.RiskTitle,
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
