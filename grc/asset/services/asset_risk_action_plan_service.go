package services

import (
	"errors"
	"fmt"
	"os"
	"time"

	assetDTOs "jethur/grc/asset/dtos"
	"jethur/grc/asset/repositories"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	masterRepositories "jethur/grc/masters/repositories"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationServices "jethur/grc/organization/services"
	taskServices "jethur/grc/task/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type AssetRiskActionPlanServiceInterface interface {
	GetAll(riskID uint, currentUserID uint, params *assetDTOs.AssetRiskActionPlanListQueryParamsDTO) (*assetDTOs.AssetRiskActionPlanListResponseDTO, error)
	GetByID(riskID uint, id uint) (*models.RiskActionPlan, error)
	GetFullDetailsByID(riskID uint, id uint, currentUserID uint, isPassiveUser bool) (*assetDTOs.AssetRiskActionPlanDetailDTO, error)
	Create(riskID uint, currentUserID uint, formDTO *assetDTOs.AssetRiskActionPlanFormDTO) (uint, error)
	Update(riskID uint, id uint, currentUserID uint, formDTO *assetDTOs.AssetRiskActionPlanFormDTO) error
	Delete(riskID uint, ids []uint, currentUserID uint) ([]uint, error)
	StatusUpdate(riskID uint, id uint, currentUserID uint, statusUpdateFormDTO *assetDTOs.AssetRiskActionPlanStatusUpdateFormDTO, tx utility.DatabaseTransaction) error

	ExportAssetRiskActionPlan(data []*assetDTOs.AssetRiskActionPlanListDTO, currentUserID uint) (*os.File, error)
}

// Service and implementations
type AssetRiskActionPlanService struct {
	repo                           repositories.AssetRiskActionPlanRepositoryInterface
	activityLogSvc                 logServices.ActivityLogServiceInterface
	assetRiskRegisterSvc           AssetRiskRegisterServiceInterface
	actionPlanStatusSvc            masterServices.ActionPlanStatusServiceInterface
	severityLevelSvc               masterServices.SeverityLevelServiceInterface
	userSettingSvc                 organizationServices.UserSettingServiceInterface
	userSvc                        organizationServices.UserServiceInterface
	notificationSvc                commonServices.NotificationServiceInterface
	taskSvc                        taskServices.TaskServiceInterface
	controlSvc                     organizationServices.ControlServiceInterface
	controlImplementationStatusSvc organizationServices.ControlImplementationStatusServiceInterface
	controlEfficiencyStatusSvc     organizationServices.ControlEfficiencyStatusServiceInterface
	delegateeUserSvc               commonServices.DelegateeUserserviceInterface
}

func NewAssetRiskActionPlanService(
	repo repositories.AssetRiskActionPlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	assetRiskRegisterSvc AssetRiskRegisterServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	severityLevelSvc masterRepositories.SeverityLevelRepositoryInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	taskSvc taskServices.TaskServiceInterface,
	controlSvc organizationServices.ControlServiceInterface,
	controlImplementationStatusSvc organizationServices.ControlImplementationStatusServiceInterface,
	controlEfficiencyStatusSvc organizationServices.ControlEfficiencyStatusServiceInterface,
	delegateeUserSvc commonServices.DelegateeUserserviceInterface,
) AssetRiskActionPlanServiceInterface {
	return &AssetRiskActionPlanService{
		repo:                           repo,
		activityLogSvc:                 activityLogSvc,
		assetRiskRegisterSvc:           assetRiskRegisterSvc,
		actionPlanStatusSvc:            actionPlanStatusSvc,
		severityLevelSvc:               severityLevelSvc,
		userSettingSvc:                 userSettingSvc,
		userSvc:                        userSvc,
		notificationSvc:                notificationSvc,
		taskSvc:                        taskSvc,
		controlSvc:                     controlSvc,
		controlImplementationStatusSvc: controlImplementationStatusSvc,
		controlEfficiencyStatusSvc:     controlEfficiencyStatusSvc,
		delegateeUserSvc:               delegateeUserSvc,
	}
}

func (araps *AssetRiskActionPlanService) GetAll(riskID uint, currentUserID uint, params *assetDTOs.AssetRiskActionPlanListQueryParamsDTO) (*assetDTOs.AssetRiskActionPlanListResponseDTO, error) {
	if riskID != 0 {
		if _, err := araps.assetRiskRegisterSvc.GetByID(riskID); err != nil {
			return nil, err
		}
	}

	if currentUserID != 0 && riskID == 0 {
		params.IsAllAccess = false
		params.AccessUserID = currentUserID

		isSuperAdmin, err := araps.userSvc.IsSuperAdmin(currentUserID)
		if err != nil {
			return nil, err
		}

		if isSuperAdmin {
			params.IsAllAccess = true
		}
	}

	return araps.repo.GetAll(riskID, params)
}

func (araps AssetRiskActionPlanService) GetByID(riskID uint, id uint) (*models.RiskActionPlan, error) {
	if riskID != 0 {
		if _, err := araps.assetRiskRegisterSvc.GetByID(riskID); err != nil {
			return nil, err
		}
	}

	riskActionPlan, err := araps.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return riskActionPlan, nil
}

func (araps AssetRiskActionPlanService) GetFullDetailsByID(riskID uint, id uint, currentUserID uint, isPassiveUser bool) (*assetDTOs.AssetRiskActionPlanDetailDTO, error) {
	var riskActionPlanDTO assetDTOs.AssetRiskActionPlanDetailDTO

	if riskID != 0 {
		if _, err := araps.assetRiskRegisterSvc.GetByID(riskID); err != nil {
			return nil, err
		}
	}

	riskActionPlan, err := araps.repo.GetFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	if riskActionPlan.RiskRegister == nil || !riskActionPlan.RiskRegister.IsCyberRisk {
		return nil, utility.NewResourceNotFoundError(gorm.ErrRecordNotFound, id)
	}

	if isPassiveUser {
		if riskActionPlan.ResponsibleUserID != currentUserID {
			return nil, utility.NewForbiddenError("Permission denied.")
		}

		if riskActionPlan.DelegateeUserID != nil && *riskActionPlan.DelegateeUserID != currentUserID {
			return nil, utility.NewForbiddenError("Permission denied.")
		}
	}

	err = mapstructure.Decode(riskActionPlan, &riskActionPlanDTO)
	if err != nil {
		return nil, err
	}

	if riskActionPlan.StartDate != nil {
		startDate := riskActionPlan.StartDate.Format("2006-01-02")
		riskActionPlanDTO.StartDateString = &startDate
	}

	if riskActionPlan.TargetDate != nil {
		targetDate := riskActionPlan.TargetDate.Format("2006-01-02")
		riskActionPlanDTO.TargetDateString = &targetDate
	}

	return &riskActionPlanDTO, nil
}

func (araps *AssetRiskActionPlanService) Create(riskID uint, currentUserID uint, formDTO *assetDTOs.AssetRiskActionPlanFormDTO) (id uint, err error) {
	tx := araps.repo.BeginTransaction()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil {
			tx.Rollback()
			return
		}

		err = tx.Commit()
	}()

	risk, err := araps.assetRiskRegisterSvc.GetByID(riskID, tx)
	if err != nil {
		return 0, err
	}

	if err := araps.validateAssetRiskActionPlanFormDTO(formDTO); err != nil {
		return 0, err
	}

	dates, err := araps.parseAssetRiskActionPlanDates(formDTO)
	if err != nil {
		return 0, err
	}

	actionPlanStatus, err := araps.actionPlanStatusSvc.GetByType("new")
	if err != nil {
		return 0, err
	}

	riskActionPlan := &models.RiskActionPlan{
		RiskRegisterID:     riskID,
		Title:              formDTO.Title,
		Description:        formDTO.Description,
		SeverityLevelID:    formDTO.SeverityLevelID,
		ResponsibleUnitID:  formDTO.ResponsibleUnitID,
		ResponsibleUserID:  formDTO.ResponsibleUserID,
		ReviewUserID:       formDTO.ReviewUserID,
		StartDate:          dates.StartDate,
		TargetDate:         dates.TargetDate,
		Cost:               formDTO.Cost,
		Percentage:         formDTO.Percentage,
		Comments:           formDTO.Comments,
		Notes:              formDTO.Notes,
		ActionPlanStatusID: actionPlanStatus.ID,
		CreatedBYID:        currentUserID,
	}

	if formDTO.ResponsibleUserID != 0 {
		delegateeID, err := araps.delegateeUserSvc.AssignDelegateeIfPermitted(formDTO.ResponsibleUserID, "asset_risk")
		if err != nil {
			return 0, err
		}
		if delegateeID != nil {
			riskActionPlan.DelegateeUserID = delegateeID
		}
	}
	id, err = araps.repo.Create(riskActionPlan, tx)
	if err != nil {
		return 0, err
	}

	actionPlanStatusID := actionPlanStatus.ID
	riskActionPlanUpdate := &models.RiskActionPlanStatusUpdate{
		RiskActionPlanID:   id,
		ActionPlanStatusID: actionPlanStatusID,
		Percentage:         formDTO.Percentage,
		Comment:            "Initial action plan status",
		CreatedByID:        currentUserID,
	}
	_, err = araps.repo.CreateInitialRiskActionPlanStatusUpdate(riskActionPlanUpdate, tx)
	if err != nil {
		return 0, err
	}

	urlParam := fmt.Sprintf(`{"risk_register_id":%d, "action_plan_id":%d}`, riskID, id)
	taskDTO := commonDTOs.TaskFormDTO{
		ModuleIdentifier:  "asset_risk",
		RiskActionPlanID:  &id,
		Title:             formDTO.Title,
		Description:       formDTO.Description,
		TargetDate:        dates.TargetDate,
		URLParam:          &urlParam,
		QueryParam:        nil,
		SeverityLevelID:   formDTO.SeverityLevelID,
		ResponsibleUserID: formDTO.ResponsibleUserID,
	}

	if _, err := araps.taskSvc.Create(taskDTO, currentUserID, tx); err != nil {
		return 0, err
	}

	defaultEfficiencyStatus, err := araps.controlEfficiencyStatusSvc.GetByType("inefficient")
	if err != nil {
		return 0, err
	}
	defaultImplementationStatus, err := araps.controlImplementationStatusSvc.GetByType("not-implemented")
	if err != nil {
		return 0, err
	}

	for _, control := range formDTO.Controls {

		// Create Control
		controlDTO := &organizationDTOs.ControlFormDTO{
			ControlTypeID:                 control.ControlTypeID,
			TitleEn:                       control.TitleEn,
			TitleAr:                       control.TitleAr,
			OwnerID:                       control.OwnerID,
			ControlEfficiencyStatusID:     &defaultEfficiencyStatus.ID,
			IsApplicable:                  true,
			ControlImplementationStatusID: &defaultImplementationStatus.ID,
		}

		controlId, err := araps.controlSvc.Create(currentUserID, controlDTO, tx)
		if err != nil {
			return 0, err
		}

		riskActionPlanControl := &models.RiskActionPlanControl{
			RiskActionPlanID: id,
			ControlID:        controlId,
			CreatedByID:      currentUserID,
		}

		if err := araps.repo.CreateRiskActionPlanControl(riskActionPlanControl, tx); err != nil {
			return 0, err
		}

		// Save risk control
		riskControl := &models.RiskControl{
			RiskRegisterID: riskID,
			ControlID:      controlId,
			CreatedByID:    currentUserID,
		}

		if err := araps.repo.CreateRiskControls(riskControl, tx); err != nil {
			return 0, err
		}

	}

	go func() {
		existingRiskActionPlan, err := araps.GetFullDetailsByID(riskActionPlan.RiskRegisterID, riskActionPlan.ID, 0, false)
		if err != nil {
			utility.LogError("Failed to fetch responsible user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = araps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.RiskActionPlan{}.TableName(),
			id,
			existingRiskActionPlan,
			currentUserID,
			"asset", "asset_risk", "action_plan", nil,
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
		urlParam := fmt.Sprintf(`{"risk_register_id":%d, "action_plan_id":%d}`, riskID, id)
		queryParam := ""
		entity := "asset-risk-action-plan"
		action := "responsible-user-assign"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{formDTO.ResponsibleUserID},
			ModuleIdentifier:     "asset_risk",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "Risk Action Plan Responsibility Assigned",
			HeadingAr:            "المسؤولية المعينة لخطة عمل المخاطر",
			TitleEn:              formDTO.Title,
			TitleAr:              formDTO.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			QueryParam:           queryParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		responsibleUser, err := araps.userSvc.GetByID(formDTO.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch responsible user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		assignedBy, err := araps.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch assigning user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		actionPlanDescription := "-"
		if formDTO.Description != nil && *formDTO.Description != "" {
			actionPlanDescription = *formDTO.Description
		}

		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("Asset Risk Action Plan: %s - Responsibility Assigned", formDTO.Title)
		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstNameEn":           responsibleUser.FirstNameEn,
			"LastNameEn":            responsibleUser.LastNameEn,
			"RiskReferenceNumber":   risk.ReferenceNumber,
			"RiskTitle":             risk.Title,
			"ActionPlanTitle":       formDTO.Title,
			"ActionPlanDescription": actionPlanDescription,
			"TargetDate":            dates.TargetDate.Format("2006-01-02"),
			"AssignedByFirstNameEn": assignedBy.FirstNameEn,
			"AssignedByLastNameEn":  assignedBy.LastNameEn,
			"AssignedOn":            time.Now().Format("2006-01-02"),
			"RedirectURL":           fmt.Sprintf("%s/asset/risk-mitigations?id=%d", appURL, id),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "risk/templates/email/risk_action_plan_create_responsible_user_mail.html",
			EmailData:            emailData,
		}

		// send notification and mail
		err = araps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "asset_risk_action_plan")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	return id, nil
}

func (araps *AssetRiskActionPlanService) Update(riskID uint, id uint, currentUserID uint, formDTO *assetDTOs.AssetRiskActionPlanFormDTO) (err error) {

	oldRiskActionPlan, err := araps.GetFullDetailsByID(riskID, id, 0, false)
	if err != nil {
		return err
	}

	risk, err := araps.assetRiskRegisterSvc.GetByID(riskID)
	if err != nil {
		return err
	}

	existingRiskActionPlan, err := araps.GetByID(riskID, id)
	if err != nil {
		return err
	}

	oldResponsibleUserID := existingRiskActionPlan.ResponsibleUserID

	if err := araps.validateAssetRiskActionPlanFormDTO(formDTO); err != nil {
		return err
	}

	dates, err := araps.parseAssetRiskActionPlanDates(formDTO)
	if err != nil {
		return err
	}

	tx := araps.repo.BeginTransaction()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil {
			tx.Rollback()
			return
		}

		err = tx.Commit()
	}()

	// Update Risk ActionPlan
	existingRiskActionPlan.Title = formDTO.Title
	existingRiskActionPlan.Description = formDTO.Description
	existingRiskActionPlan.SeverityLevelID = formDTO.SeverityLevelID
	existingRiskActionPlan.ResponsibleUnitID = formDTO.ResponsibleUnitID
	existingRiskActionPlan.ResponsibleUserID = formDTO.ResponsibleUserID
	existingRiskActionPlan.ReviewUserID = formDTO.ReviewUserID
	existingRiskActionPlan.StartDate = dates.StartDate
	existingRiskActionPlan.TargetDate = dates.TargetDate
	existingRiskActionPlan.Cost = formDTO.Cost
	existingRiskActionPlan.Comments = formDTO.Comments
	existingRiskActionPlan.Notes = formDTO.Notes
	existingRiskActionPlan.UpdatedBYID = &currentUserID

	if formDTO.ResponsibleUserID != 0 && oldResponsibleUserID != formDTO.ResponsibleUserID {
		existingRiskActionPlan.DelegateeUserID = nil
		delegateeID, err := araps.delegateeUserSvc.AssignDelegateeIfPermitted(formDTO.ResponsibleUserID, "asset_risk")
		if err != nil {
			return err
		}
		if delegateeID != nil {
			existingRiskActionPlan.DelegateeUserID = delegateeID
		}
	}

	err = araps.repo.Update(existingRiskActionPlan, tx)
	if err != nil {
		return err
	}

	urlParam := fmt.Sprintf(`{"risk_register_id":%d, "action_plan_id":%d}`, riskID, id)

	taskUpdateDTO := commonDTOs.TaskFormDTO{
		Title:             formDTO.Title,
		Description:       formDTO.Description,
		TargetDate:        dates.TargetDate,
		URLParam:          &urlParam,
		SeverityLevelID:   formDTO.SeverityLevelID,
		ResponsibleUserID: formDTO.ResponsibleUserID,
	}

	err = araps.taskSvc.Update(currentUserID, taskUpdateDTO, "risk_action_plan_id", id, tx)
	if err != nil {
		return err
	}

	existingControlIDs := make(map[uint]bool)
	for _, ctrl := range oldRiskActionPlan.Controls {
		existingControlIDs[ctrl.ControlID] = true
	}

	newControlIDs := make(map[uint]bool)
	for _, ctrl := range formDTO.Controls {
		if ctrl.ID != nil {
			newControlIDs[*ctrl.ID] = true
		}
	}

	var deleteControlIDs []*uint
	for _, existingActionPlan := range oldRiskActionPlan.Controls {
		if _, exists := newControlIDs[existingActionPlan.ControlID]; !exists {
			deleteControlIDs = append(deleteControlIDs, existingActionPlan.ID)
		}
	}

	if len(deleteControlIDs) > 0 {
		if err := araps.repo.DeleteRiskActionPlanControl(deleteControlIDs, tx); err != nil {
			return err
		}
	}

	var deletRiskControlIDs []*uint
	for _, existingRiskControl := range oldRiskActionPlan.RiskRegister.Controls {
		if _, exists := existingControlIDs[existingRiskControl.ControlID]; exists {
			if _, exists := newControlIDs[existingRiskControl.ControlID]; !exists {
				deletRiskControlIDs = append(deletRiskControlIDs, existingRiskControl.ID)
			}
		}
	}

	if len(deletRiskControlIDs) > 0 {
		if err := araps.repo.DeleteRiskControls(deletRiskControlIDs, tx); err != nil {
			return err
		}
	}

	defaultEfficiencyStatus, err := araps.controlEfficiencyStatusSvc.GetByType("inefficient")
	if err != nil {
		return err
	}
	defaultImplementationStatus, err := araps.controlImplementationStatusSvc.GetByType("not-implemented")
	if err != nil {
		return err
	}

	for _, control := range formDTO.Controls {
		if control.ID == nil {
			// Create Control
			controlDTO := &organizationDTOs.ControlFormDTO{
				ControlTypeID:                 control.ControlTypeID,
				TitleEn:                       control.TitleEn,
				TitleAr:                       control.TitleAr,
				OwnerID:                       control.OwnerID,
				ControlEfficiencyStatusID:     &defaultEfficiencyStatus.ID,
				IsApplicable:                  true,
				ControlImplementationStatusID: &defaultImplementationStatus.ID,
			}

			controlId, err := araps.controlSvc.Create(currentUserID, controlDTO, tx)
			if err != nil {
				return err
			}

			riskActionPlanControl := &models.RiskActionPlanControl{
				RiskActionPlanID: id,
				ControlID:        controlId,
				CreatedByID:      currentUserID,
			}

			if err := araps.repo.CreateRiskActionPlanControl(riskActionPlanControl, tx); err != nil {
				return err
			}

			// Save risk control
			riskControl := &models.RiskControl{
				RiskRegisterID: riskID,
				ControlID:      controlId,
				CreatedByID:    currentUserID,
			}

			if err := araps.repo.CreateRiskControls(riskControl, tx); err != nil {
				return err
			}
		}
	}

	go func() {
		existingActionPlan, err := araps.GetFullDetailsByID(riskID, id, 0, false)
		if err != nil {
			return
		}

		// Prepare log data
		changedData := make(map[string]interface{})

		// log changes
		oldData := map[string]interface{}{
			"title":              oldRiskActionPlan.Title,
			"description":        oldRiskActionPlan.Description,
			"severity_level":     oldRiskActionPlan.SeverityLevel.TitleEn,
			"responsible_user":   fmt.Sprintf("%s %s", oldRiskActionPlan.ResponsibleUser.FirstNameEN, oldRiskActionPlan.ResponsibleUser.LastNameEN),
			"review_user":        utility.GetFieldValue(oldRiskActionPlan.ReviewUser).(string),
			"target_date":        oldRiskActionPlan.TargetDateString,
			"percentage":         oldRiskActionPlan.Percentage,
			"comments":           oldRiskActionPlan.Comments,
			"action_plan_status": oldRiskActionPlan.ActionPlanStatus.TitleEn,
		}

		newData := map[string]interface{}{
			"title":              existingActionPlan.Title,
			"description":        existingActionPlan.Description,
			"severity_level":     existingActionPlan.SeverityLevel.TitleEn,
			"responsible_user":   fmt.Sprintf("%s %s", existingActionPlan.ResponsibleUser.FirstNameEN, existingActionPlan.ResponsibleUser.LastNameEN),
			"review_user":        utility.GetFieldValue(existingActionPlan.ReviewUser).(string),
			"target_date":        existingActionPlan.TargetDateString,
			"percentage":         existingActionPlan.Percentage,
			"comments":           existingActionPlan.Comments,
			"action_plan_status": existingActionPlan.ActionPlanStatus.TitleEn,
		}

		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}

		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := araps.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.RiskActionPlan{}.TableName(),
				id,
				changedData,
				currentUserID,
				"asset", "asset_risk", "action_plan", nil,
			)
			if err != nil {
				utility.LogError("Failed to log activity", map[string]interface{}{
					"action":   models.UpdateAction,
					"error":    err.Error(),
					"recordID": id,
				})
			}
		}

		// sends notification if responsible user has changed
		if oldResponsibleUserID != existingRiskActionPlan.ResponsibleUserID {
			urlParam := fmt.Sprintf(`{"risk_register_id":%d, "action_plan_id":%d}`, riskID, id)
			queryParam := ""
			entity := "asset-risk-action-plan"

			oldResponsibleUserAction := "responsible-user-remove"
			oldResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{oldResponsibleUserID},
				ModuleIdentifier:     "asset_risk",
				NotificationType:     string(commonDTOs.InformationType),
				HeadingEn:            "Risk Action Plan Responsibility Removed",
				HeadingAr:            "تمت إزالة مسؤولية خطة عمل المخاطر",
				TitleEn:              existingRiskActionPlan.Title,
				TitleAr:              existingRiskActionPlan.Title,
				Entity:               &entity,
				Action:               &oldResponsibleUserAction,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			newResponsibleUserAction := "responsible-user-assign"
			newResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{formDTO.ResponsibleUserID},
				ModuleIdentifier:     "asset_risk",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Risk Action Plan Responsibility Assigned",
				HeadingAr:            "المسؤولية المعينة لخطة عمل المخاطر",
				TitleEn:              existingRiskActionPlan.Title,
				TitleAr:              existingRiskActionPlan.Title,
				Entity:               &entity,
				Action:               &newResponsibleUserAction,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			oldResponsibleUser, err := araps.userSvc.GetByID(oldResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch old responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			newResponsibleUser, err := araps.userSvc.GetByID(formDTO.ResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch new responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			assignedBy, err := araps.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch assigning user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			appURL := utility.GetAppURL()

			actionPlanDescription := "-"
			if existingRiskActionPlan.Description != nil && *existingRiskActionPlan.Description != "" {
				actionPlanDescription = *existingRiskActionPlan.Description
			}

			oldResponsibleUserEmailTo := []string{oldResponsibleUser.Email}
			oldResponsibleUserSubject := fmt.Sprintf("Asset Risk Action Plan: %s - Responsibility Removed", existingRiskActionPlan.Title)
			oldResponsibleUserTemplateRelativePath := "risk/templates/email/risk_action_plan_remove_responsible_user_mail.html"
			oldResponsibleUserEmailData := map[string]interface{}{
				"FirstNameEn":           oldResponsibleUser.FirstNameEn,
				"LastNameEn":            oldResponsibleUser.LastNameEn,
				"RiskReferenceNumber":   risk.ReferenceNumber,
				"RiskTitle":             risk.Title,
				"ActionPlanTitle":       existingRiskActionPlan.Title,
				"ActionPlanDescription": actionPlanDescription,
				"RedirectURL":           fmt.Sprintf("%s/asset/risk-mitigations?id=%d", appURL, id),
			}

			oldResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              oldResponsibleUserEmailTo,
				Subject:              oldResponsibleUserSubject,
				TemplateRelativePath: oldResponsibleUserTemplateRelativePath,
				EmailData:            oldResponsibleUserEmailData,
			}

			newResponsibleUserEmailTo := []string{newResponsibleUser.Email}
			newResponsibleUserSubject := fmt.Sprintf("Asset Risk Action Plan: %s - Responsibility Assigned", existingRiskActionPlan.Title)
			newResponsibleUserTemplateRelativePath := "risk/templates/email/risk_action_plan_change_responsible_user_mail.html"
			newResponsibleUserEmailData := map[string]interface{}{
				"FirstNameEn":           newResponsibleUser.FirstNameEn,
				"LastNameEn":            newResponsibleUser.LastNameEn,
				"RiskReferenceNumber":   risk.ReferenceNumber,
				"RiskTitle":             risk.Title,
				"ActionPlanTitle":       existingRiskActionPlan.Title,
				"ActionPlanDescription": actionPlanDescription,
				"TargetDate":            existingRiskActionPlan.TargetDate.Format("2006-01-02"),
				"AssignedByFirstNameEn": assignedBy.FirstNameEn,
				"AssignedByLastNameEn":  assignedBy.LastNameEn,
				"AssignedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           fmt.Sprintf("%s/asset/risk-mitigations?id=%d", appURL, id),
			}

			newResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              newResponsibleUserEmailTo,
				Subject:              newResponsibleUserSubject,
				TemplateRelativePath: newResponsibleUserTemplateRelativePath,
				EmailData:            newResponsibleUserEmailData,
			}

			err = araps.notificationSvc.SendNotification(
				newResponsibleUserNotificationDTO, newResponsibleUserEmailNotificationDTO, "asset_risk_action_plan",
			)
			if err != nil {
				utility.LogError("Failed to send new responsible user notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}

			err = araps.notificationSvc.SendNotification(
				oldResponsibleUserNotificationDTO, oldResponsibleUserEmailNotificationDTO, "asset_risk_action_plan",
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

func (araps *AssetRiskActionPlanService) Delete(riskID uint, ids []uint, currentUserID uint) (failedIds []uint, err error) {
	if _, err := araps.assetRiskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}

	tx := araps.repo.BeginTransaction()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if len(failedIds) > 0 {
			tx.Rollback()
			return
		}

		err = tx.Commit()
	}()

	var successfulIds []uint
	var notifications []struct {
		OwnerID         *uint
		TitleEn         string
		TitleAr         string
		Description     *string
		ReferenceNumber string
		RiskTitle       string
	}

	var deletedRiskActionPlan []map[string]interface{}

	for _, id := range ids {

		riskActionPlan, err := araps.repo.GetFullDetailsByID(id, tx)
		if err != nil {
			continue
		}

		description := "-"
		if riskActionPlan.Description != nil && *riskActionPlan.Description != "" {
			description = *riskActionPlan.Description
		}

		referenceNumber := "-"
		riskTitle := "-"

		if riskActionPlan.RiskRegister != nil {
			if riskActionPlan.RiskRegister.ReferenceNumber != nil &&
				*riskActionPlan.RiskRegister.ReferenceNumber != "" {
				referenceNumber = *riskActionPlan.RiskRegister.ReferenceNumber
			}

			riskTitle = riskActionPlan.RiskRegister.Title
		}

		// Prepare notification data before deletion
		notifications = append(notifications, struct {
			OwnerID         *uint
			TitleEn         string
			TitleAr         string
			Description     *string
			ReferenceNumber string
			RiskTitle       string
		}{
			OwnerID:         &riskActionPlan.ResponsibleUserID,
			TitleEn:         riskTitle,
			TitleAr:         riskTitle,
			Description:     &description,
			ReferenceNumber: referenceNumber,
			RiskTitle:       riskActionPlan.RiskRegister.Title,
		})

		deletedRiskActionPlan = append(deletedRiskActionPlan, map[string]interface{}{
			"id":    id,
			"title": riskActionPlan.Title,
		})

		existingTask, err := araps.taskSvc.GetByFieldID("risk_action_plan_id", id, tx)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		taskID, err := araps.taskSvc.Delete(existingTask.ID, currentUserID)
		if err != nil {
			failedIds = append(failedIds, taskID)
		}

		existingRiskActionPlan, err := araps.GetFullDetailsByID(riskID, id, 0, false)
		if err != nil {
			continue
		}

		newControlIDs := make(map[uint]bool)
		for _, ctrl := range existingRiskActionPlan.Controls {
			if ctrl.ID != nil {
				newControlIDs[ctrl.ControlID] = true
			}
		}

		if len(existingRiskActionPlan.Controls) > 0 {
			var deleteControlIds []*uint
			for _, control := range existingRiskActionPlan.Controls {
				deleteControlIds = append(deleteControlIds, control.ID)
			}
			if err := araps.repo.DeleteRiskActionPlanControl(deleteControlIds, tx); err != nil {
				failedIds = append(failedIds, id)
				continue
			}
		}

		if len(existingRiskActionPlan.RiskRegister.Controls) > 0 {
			var deleteRiskControlIds []*uint
			for _, control := range existingRiskActionPlan.RiskRegister.Controls {
				if _, exists := newControlIDs[control.ControlID]; exists {
					deleteRiskControlIds = append(deleteRiskControlIds, control.ID)
				}
			}
			if err := araps.repo.DeleteRiskControls(deleteRiskControlIds, tx); err != nil {
				failedIds = append(failedIds, id)
				continue
			}
		}

		if len(riskActionPlan.ActionPlanStatusUpdates) > 0 {
			var deleteStatusUpdateIds []*uint
			for _, update := range riskActionPlan.ActionPlanStatusUpdates {
				deleteStatusUpdateIds = append(deleteStatusUpdateIds, &update.ID)
			}
			if err := araps.repo.DeleteRiskActionPlanStatusUpdates(deleteStatusUpdateIds, tx); err != nil {
				failedIds = append(failedIds, id)
				continue
			}
		}

		err = araps.repo.Delete(id, tx)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {
		go func() {
			err := araps.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.RiskActionPlan{}.TableName(),
				successfulIds,
				deletedRiskActionPlan,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}

			for _, notification := range notifications {
				urlParam := ""
				entity := "asset-risk-action-plan"
				action := "delete"

				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserIDs:              []uint{*notification.OwnerID},
					ModuleIdentifier:     "asset_risk",
					NotificationType:     string(commonDTOs.DeletionType),
					HeadingEn:            "Risk Action Plan Deleted",
					HeadingAr:            "تم حذف خطة عمل المخاطر",
					TitleEn:              notification.TitleEn,
					TitleAr:              notification.TitleAr,
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				ownerUser, err := araps.userSvc.GetByID(*notification.OwnerID)
				if err != nil {
					utility.LogError("Failed to fetch owner details for notification", map[string]interface{}{
						"error": err.Error(),
					})
					continue
				}

				ownerEmailTo := []string{ownerUser.Email}
				subject := fmt.Sprintf("Asset Risk Action Plan: %s - Deleted", notification.TitleEn)
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

				err = araps.notificationSvc.SendNotification(
					dbNotificationDTO, emailNotificationDTO, "asset_risk_action_plan")
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

func (araps *AssetRiskActionPlanService) StatusUpdate(riskID uint, id uint, currentUserID uint, statusUpdateFormDTO *assetDTOs.AssetRiskActionPlanStatusUpdateFormDTO, tx utility.DatabaseTransaction) error {
	existingRiskActionPlan, err := araps.GetByID(riskID, id)
	if err != nil {
		return err
	}

	existingRiskActionPlan.ActionPlanStatusID = statusUpdateFormDTO.ActionPlanStatusID
	existingRiskActionPlan.Percentage = statusUpdateFormDTO.Percentage
	existingRiskActionPlan.UpdatedBYID = &currentUserID

	err = araps.repo.Update(existingRiskActionPlan, tx)
	if err != nil {
		return err
	}

	isReviewUserNull := existingRiskActionPlan.ReviewUserID == nil

	// Fetch the corresponding task and update its status
	existingTask, err := araps.taskSvc.GetByFieldID("risk_action_plan_id", id)
	if err != nil {
		return err
	}
	err = araps.taskSvc.UpdateStatus(currentUserID, existingTask.ID, statusUpdateFormDTO.ActionPlanStatusID, isReviewUserNull, tx)
	if err != nil {
		return err
	}

	return nil
}

func (araps *AssetRiskActionPlanService) validateAssetRiskActionPlanFormDTO(formDTO *assetDTOs.AssetRiskActionPlanFormDTO) error {
	if formDTO.SeverityLevelID != nil {
		if _, err := araps.severityLevelSvc.GetByID(*formDTO.SeverityLevelID); err != nil {
			return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
		}
	}
	if formDTO.Percentage == nil {
		defaultValue := uint(0)
		formDTO.Percentage = &defaultValue
	}
	return nil
}

func (araps *AssetRiskActionPlanService) parseAssetRiskActionPlanDates(dto *assetDTOs.AssetRiskActionPlanFormDTO) (*assetDTOs.AssetRiskActionPlanParsedDates, error) {
	parsed := &assetDTOs.AssetRiskActionPlanParsedDates{}

	var err error
	if dto.StartDate != nil && *dto.StartDate != "" {
		parsed.StartDate, err = utility.ParseDate(*dto.StartDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid start date"), "invalid_date", "start date")
		}
	}

	if dto.TargetDate != nil && *dto.TargetDate != "" {
		parsed.TargetDate, err = utility.ParseDate(*dto.TargetDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid target date"), "invalid_date", "target date")
		}
	}

	return parsed, nil
}

func (araps *AssetRiskActionPlanService) ExportAssetRiskActionPlan(data []*assetDTOs.AssetRiskActionPlanListDTO, currentUserID uint) (*os.File, error) {
	// Fetch user setting for language preference
	userSetting, err := araps.userSettingSvc.GetUserSettingsByUserID(currentUserID)
	if err != nil {
		return nil, err
	}

	headers := []string{
		"SL No.",
		"Action Plan",
		"Risk",
		"Description",
		"Severity",
		"Responsible User",
		"Target Date",
		"Status",
	}

	// Prepare the data rows
	rows := make([][]interface{}, len(data))
	for i, actionPlan := range data {
		var responsibleUser, severity, status string

		if userSetting.Language == "ar" {
			responsibleUser = actionPlan.ResponsibleUserFirstNameAr + " " + actionPlan.ResponsibleUserLastNameAr
			severity = actionPlan.SeverityLevelTitleAr
			status = actionPlan.ActionPlanStatusTitleAr
		} else {
			responsibleUser = actionPlan.ResponsibleUserFirstNameEn + " " + actionPlan.ResponsibleUserLastNameEn
			severity = actionPlan.SeverityLevelTitleEn
			status = actionPlan.ActionPlanStatusTitleEn
		}

		rows[i] = []interface{}{
			i + 1,
			actionPlan.Title,
			actionPlan.RiskTitle,
			actionPlan.Description,
			severity,
			responsibleUser,
			actionPlan.TargetDate,
			status,
		}
	}

	// Generate the Excel file
	return utility.GenerateExcelFile(headers, rows)
}
