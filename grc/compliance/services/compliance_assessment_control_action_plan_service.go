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
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	settingService "jethur/grc/settings/services"
	taskServices "jethur/grc/task/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ComplianceAssessmentControlActionPlanServiceInterface interface {
	GetAll(complianceAssessmentID uint, controlID uint, params *complianceDTOs.ComplianceActionPlanListQueryParamsDTO) (*complianceDTOs.ComplianceActionPlanListResponseDTO, error)
	GetByID(complianceAssessmentID uint, controlID uint, id uint) (*models.ComplianceActionPlan, error)
	GetFullDetailsByID(complianceAssessmentID uint, controlID uint, id uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceActionPlanDetailDTO, error)
	Create(complianceAssessmentID uint, controlID uint, currentUserID uint, ComplianceActionPlanFormDTO *complianceDTOs.ComplianceActionPlanFormDTO) (uint, error)
	Update(complianceAssessmentID uint, controlID uint, id uint, currentUserID uint, compliance *complianceDTOs.ComplianceActionPlanFormDTO) error
	Delete(complianceAssessmentID uint, controlID uint, ids []uint, currentUserID uint) ([]uint, error)

	StatusUpdate(complianceAssessmentID uint, controlID uint, id uint, currentUserID uint, complianceActionPlanStatusUpdateFormDTO *complianceDTOs.ComplianceActionPlanStatusUpdateFormDTO, tx utility.DatabaseTransaction) error
}

// Service and implementations
type ComplianceAssessmentControlActionPlanService struct {
	repo                           complianceRepositories.ComplianceActionPlanRepositoryInterface
	activityLogSvc                 logServices.ActivityLogServiceInterface
	complianceAssessmentSvc        ComplianceAssessmentServiceInterface
	complianceAssessmentControlSvc ComplianceAssessmentControlServiceInterface
	actionPlanStatusSvc            masterServices.ActionPlanStatusServiceInterface
	userSvc                        organizationServices.UserServiceInterface
	organizationUnitSvc            settingService.OrganizationUnitServiceInterface
	taskSvc                        taskServices.TaskServiceInterface
	notificationSvc                commonServices.NotificationServiceInterface
	delegateeUserSvc               commonServices.DelegateeUserserviceInterface
}

func NewComplianceAssessmentControlActionPlanService(
	repo complianceRepositories.ComplianceActionPlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	complianceAssessmentSvc ComplianceAssessmentServiceInterface,
	complianceAssessmentControlSvc ComplianceAssessmentControlServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	organizationUnitSvc settingService.OrganizationUnitServiceInterface,
	taskSvc taskServices.TaskServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	delegateeUserSvc commonServices.DelegateeUserserviceInterface,
) ComplianceAssessmentControlActionPlanServiceInterface {
	return &ComplianceAssessmentControlActionPlanService{
		repo:                           repo,
		activityLogSvc:                 activityLogSvc,
		complianceAssessmentSvc:        complianceAssessmentSvc,
		complianceAssessmentControlSvc: complianceAssessmentControlSvc,
		actionPlanStatusSvc:            actionPlanStatusSvc,
		userSvc:                        userSvc,
		organizationUnitSvc:            organizationUnitSvc,
		taskSvc:                        taskSvc,
		notificationSvc:                notificationSvc,
		delegateeUserSvc:               delegateeUserSvc,
	}
}

func (cacaps *ComplianceAssessmentControlActionPlanService) GetAll(complianceAssessmentID uint, controlID uint, params *complianceDTOs.ComplianceActionPlanListQueryParamsDTO) (*complianceDTOs.ComplianceActionPlanListResponseDTO, error) {
	if controlID != 0 {
		if _, err := cacaps.complianceAssessmentControlSvc.GetByID(controlID); err != nil {
			return nil, err
		}
	}
	return cacaps.repo.GetAll(uint(0), params, uint(0), &controlID)
}

func (cacaps ComplianceAssessmentControlActionPlanService) GetByID(complianceAssessmentID uint, controlID uint, id uint) (*models.ComplianceActionPlan, error) {
	if _, err := cacaps.complianceAssessmentControlSvc.GetByID(controlID); err != nil {
		return nil, err
	}

	complianceAssessmentControlActionPlan, err := cacaps.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return complianceAssessmentControlActionPlan, nil
}

func (cacaps ComplianceAssessmentControlActionPlanService) GetFullDetailsByID(complianceAssessmentID uint, controlID uint, id uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceActionPlanDetailDTO, error) {
	var complianceAssessmentControlActionPlanDTO complianceDTOs.ComplianceActionPlanDetailDTO

	// Check if the compliance register exists
	if _, err := cacaps.complianceAssessmentControlSvc.GetByID(controlID); err != nil {
		return nil, err
	}

	// Get full details of the compliance action plan
	complianceAssessmentControlActionPlan, err := cacaps.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Decode the compliance action plan into the DTO
	err = mapstructure.Decode(complianceAssessmentControlActionPlan, &complianceAssessmentControlActionPlanDTO)
	if err != nil {
		return nil, err
	}
	if complianceAssessmentControlActionPlan.TargetDate != nil {
		targetDate := complianceAssessmentControlActionPlan.TargetDate.Format("2006-01-02")
		complianceAssessmentControlActionPlanDTO.TargetDateString = &targetDate
	}

	responsibleUnitID, err := cacaps.userSvc.GetUserPrimaryUnitID(complianceAssessmentControlActionPlan.ResponsibleUserID)
	if err != nil {
		return nil, err
	}

	responsibleUnit, err := cacaps.organizationUnitSvc.GetFullDetailsByID(responsibleUnitID)
	if err != nil {
		return nil, err
	}
	complianceAssessmentControlActionPlanDTO.ResponsibleUserUnitID = &responsibleUnitID
	complianceAssessmentControlActionPlanDTO.ResponsibleUserUnit = responsibleUnit

	if complianceAssessmentControlActionPlan.ReviewUserID != nil {
		reviewUserUnitID, err := cacaps.userSvc.GetUserPrimaryUnitID(*complianceAssessmentControlActionPlan.ReviewUserID)
		if err != nil {
			return nil, err
		}

		reviewUserUnit, err := cacaps.organizationUnitSvc.GetFullDetailsByID(reviewUserUnitID)
		if err != nil {
			return nil, err
		}
		complianceAssessmentControlActionPlanDTO.ReviewUserUnitID = &reviewUserUnitID
		complianceAssessmentControlActionPlanDTO.ReviewUserUnit = reviewUserUnit
	}

	complianceAssessmentControlActionPlanDTO.CreatedAt = complianceAssessmentControlActionPlan.CreatedAt
	complianceAssessmentControlActionPlanDTO.UpdatedAt = complianceAssessmentControlActionPlan.UpdatedAt

	// TODO: Add any additional fields required for the details screen

	return &complianceAssessmentControlActionPlanDTO, nil
}

func (cacaps *ComplianceAssessmentControlActionPlanService) Create(complianceAssessmentID uint, controlID uint, currentUserID uint, ComplianceActionPlanFormDTO *complianceDTOs.ComplianceActionPlanFormDTO) (uint, error) {
	tx := cacaps.repo.BeginTransaction()

	if err := cacaps.validateComplianceActionPlanFormDTO(ComplianceActionPlanFormDTO); err != nil {
		return 0, err
	}

	dates, err := cacaps.parseComplianceAssessmentControlActionPlanDates(ComplianceActionPlanFormDTO)
	if err != nil {
		return 0, err
	}

	complianceAssessment, err := cacaps.complianceAssessmentSvc.GetFullDetailsByID(complianceAssessmentID)
	if err != nil {
		return 0, err
	}
	actionPlanStatus, err := cacaps.actionPlanStatusSvc.GetByType("new")
	if err != nil {
		return 0, err
	}

	complianceAssessmentControlActionPlan := &models.ComplianceActionPlan{
		ComplianceRegisterID:          nil,
		ComplianceAssessmentControlID: &controlID,
		Title:                         ComplianceActionPlanFormDTO.Title,
		Description:                   ComplianceActionPlanFormDTO.Description,
		ResponsibleUserID:             ComplianceActionPlanFormDTO.ResponsibleUserID,
		TargetDate:                    dates.TargetDate,
		Percentage:                    ComplianceActionPlanFormDTO.Percentage,
		ActionPlanStatusID:            actionPlanStatus.ID,
		SeverityLevelID:               ComplianceActionPlanFormDTO.SeverityLevelID,
		ReviewUserID:                  ComplianceActionPlanFormDTO.ReviewUserID,
		Comments:                      ComplianceActionPlanFormDTO.Comments,
		CreatedByID:                   currentUserID,
	}

	if ComplianceActionPlanFormDTO.ResponsibleUserID != 0 {
		delegateeID, err := cacaps.delegateeUserSvc.AssignDelegateeIfPermitted(ComplianceActionPlanFormDTO.ResponsibleUserID, "compliance_assessment")
		if err != nil {
			return 0, err
		}
		if delegateeID != nil {
			complianceAssessmentControlActionPlan.DelegateeUserID = delegateeID
		}
	}

	id, err := cacaps.repo.Create(complianceAssessmentControlActionPlan, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	actionPlanStatusID := actionPlanStatus.ID
	complianceActionPlanUpdate := &models.ComplianceActionPlanStatusUpdate{
		ComplianceActionPlanID: id,
		ActionPlanStatusID:     actionPlanStatusID,
		Percentage:             ComplianceActionPlanFormDTO.Percentage,
		Comment:                "Initial action plan status",
		CreatedByID:            currentUserID,
	}

	_, err = cacaps.repo.CreateInitialComplianceActionPlanStatusUpdate(complianceActionPlanUpdate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	urlParam := fmt.Sprintf(`{"compliance_assessment_id":%d,"control_id":%d, "action_plan_id":%d}`, complianceAssessmentID, controlID, id)
	taskDTO := commonDTOs.TaskFormDTO{
		ModuleIdentifier:       "compliance_assessment",
		ComplianceActionPlanID: &id,
		Title:                  ComplianceActionPlanFormDTO.Title,
		Description:            ComplianceActionPlanFormDTO.Description,
		TargetDate:             dates.TargetDate,
		URLParam:               &urlParam,
		QueryParam:             nil,
		SeverityLevelID:        ComplianceActionPlanFormDTO.SeverityLevelID,
		ResponsibleUserID:      ComplianceActionPlanFormDTO.ResponsibleUserID,
	}

	if _, err := cacaps.taskSvc.Create(taskDTO, currentUserID, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	go func() {
		existingComplianceAssessmentControlActionPlan, err := cacaps.GetFullDetailsByID(complianceAssessmentID, controlID, id)
		if err != nil {
			utility.LogError("Failed to fetch compliance action plan", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = cacaps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ComplianceActionPlan{}.TableName(),
			id,
			existingComplianceAssessmentControlActionPlan,
			currentUserID,
			"compliance",
			"compliance_assessment",
			"action_plan",
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
		urlParam := fmt.Sprintf(`{"compliance_assessment_id":%d,"control_id":%d, "action_plan_id":%d}`, complianceAssessmentID, controlID, id)
		queryParam := ""
		entity := "compliance-assessment-action-plan"
		action := "responsible-user-assign"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{ComplianceActionPlanFormDTO.ResponsibleUserID},
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

		responsibleUser, err := cacaps.userSvc.GetByID(ComplianceActionPlanFormDTO.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch responsible user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		assignedBy, err := cacaps.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch assigning user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("Compliance Assessment Requirement Action Plan%s: Responsibility Assigned", ComplianceActionPlanFormDTO.Title)
		appURL := utility.GetAppURL()
		// Get assessment title
		var assessmentTitle string
		if len(complianceAssessment.Compliances) > 0 && complianceAssessment.Compliances[0].ComplianceRegister != nil {
			assessmentTitle = complianceAssessment.Compliances[0].ComplianceRegister.TitleEn
		} else if complianceAssessment.ResponsibleOrganizationUnit != nil {
			assessmentTitle = complianceAssessment.ResponsibleOrganizationUnit.TitleEn + "'s Assessment"
		}

		emailData := map[string]interface{}{
			"FirstNameEn":           responsibleUser.FirstNameEn,
			"LastNameEn":            responsibleUser.LastNameEn,
			"ComplianceAssessment":  assessmentTitle,
			"ActionPlanTitle":       ComplianceActionPlanFormDTO.Title,
			"ActionPlanDescription": ComplianceActionPlanFormDTO.Description,
			"AssignedByFirstNameEn": assignedBy.FirstNameEn,
			"AssignedByLastNameEn":  assignedBy.LastNameEn,
			"AssignedOn":            time.Now().Format("2006-01-02"),
			"RedirectURL":           fmt.Sprintf("%s/compliance/assessments/%d/overview", appURL, complianceAssessmentID),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "compliance/templates/email/compliance_assessment_action_plan_create_responsible_user_mail.html",
			EmailData:            emailData,
		}

		// send notification and mail
		err = cacaps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "compliance_action_plan")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	return id, nil
}

func (cacaps *ComplianceAssessmentControlActionPlanService) Update(complianceAssessmentID uint, controlID uint, id uint, currentUserID uint, ComplianceActionPlanFormDTO *complianceDTOs.ComplianceActionPlanFormDTO) error {
	tx := cacaps.repo.BeginTransaction()

	existingComplianceAssessmentControlActionPlan, err := cacaps.GetByID(complianceAssessmentID, controlID, id)
	if err != nil {
		return err
	}

	oldComplianceAssessmentControlActionPlan, err := cacaps.GetFullDetailsByID(complianceAssessmentID, controlID, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	complianceAssessment, err := cacaps.complianceAssessmentSvc.GetFullDetailsByID(complianceAssessmentID)
	if err != nil {
		return err
	}

	oldResponsibleUserID := existingComplianceAssessmentControlActionPlan.ResponsibleUserID

	if err := cacaps.validateComplianceActionPlanFormDTO(ComplianceActionPlanFormDTO); err != nil {
		return err
	}

	dates, err := cacaps.parseComplianceAssessmentControlActionPlanDates(ComplianceActionPlanFormDTO)
	if err != nil {
		return err
	}

	// Update Compliance ActionPlan
	existingComplianceAssessmentControlActionPlan.Title = ComplianceActionPlanFormDTO.Title
	existingComplianceAssessmentControlActionPlan.Description = ComplianceActionPlanFormDTO.Description
	existingComplianceAssessmentControlActionPlan.ResponsibleUserID = ComplianceActionPlanFormDTO.ResponsibleUserID
	existingComplianceAssessmentControlActionPlan.TargetDate = dates.TargetDate
	existingComplianceAssessmentControlActionPlan.SeverityLevelID = ComplianceActionPlanFormDTO.SeverityLevelID
	existingComplianceAssessmentControlActionPlan.Comments = ComplianceActionPlanFormDTO.Comments
	existingComplianceAssessmentControlActionPlan.ReviewUserID = ComplianceActionPlanFormDTO.ReviewUserID

	existingComplianceAssessmentControlActionPlan.UpdatedByID = &currentUserID

	if ComplianceActionPlanFormDTO.ResponsibleUserID != 0 && oldResponsibleUserID != ComplianceActionPlanFormDTO.ResponsibleUserID {
		existingComplianceAssessmentControlActionPlan.DelegateeUserID = nil
		delegateeID, err := cacaps.delegateeUserSvc.AssignDelegateeIfPermitted(ComplianceActionPlanFormDTO.ResponsibleUserID, "compliance_assessment")
		if err != nil {
			return err
		}
		if delegateeID != nil {
			existingComplianceAssessmentControlActionPlan.DelegateeUserID = delegateeID
		}
	}

	err = cacaps.repo.Update(existingComplianceAssessmentControlActionPlan, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// //update task
	urlParam := fmt.Sprintf(`{"compliance_assessment_id":%d,"control_id":%d, "action_plan_id":%d}`, complianceAssessmentID, controlID, id)

	taskUpdateDTO := commonDTOs.TaskFormDTO{
		Title:             ComplianceActionPlanFormDTO.Title,
		Description:       ComplianceActionPlanFormDTO.Description,
		URLParam:          &urlParam,
		SeverityLevelID:   ComplianceActionPlanFormDTO.SeverityLevelID,
		ResponsibleUserID: ComplianceActionPlanFormDTO.ResponsibleUserID,
	}

	err = cacaps.taskSvc.Update(currentUserID, taskUpdateDTO, "compliance_action_plan_id", id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		existingActionPlan, err := cacaps.GetFullDetailsByID(complianceAssessmentID, controlID, id)
		if err != nil {
			return
		}
		// log changes
		changedData := make(map[string]interface{})
		oldData := map[string]interface{}{
			"title":              oldComplianceAssessmentControlActionPlan.Title,
			"description":        oldComplianceAssessmentControlActionPlan.Description,
			"responsible_user":   fmt.Sprintf("%s %s", oldComplianceAssessmentControlActionPlan.ResponsibleUser.FirstNameEN, oldComplianceAssessmentControlActionPlan.ResponsibleUser.LastNameEN),
			"target_date":        oldComplianceAssessmentControlActionPlan.TargetDateString,
			"percentage":         oldComplianceAssessmentControlActionPlan.Percentage,
			"action_plan_status": oldComplianceAssessmentControlActionPlan.ActionPlanStatus.TitleEn,
		}

		newData := map[string]interface{}{
			"title":              existingActionPlan.Title,
			"description":        existingActionPlan.Description,
			"responsible_user":   fmt.Sprintf("%s %s", existingActionPlan.ResponsibleUser.FirstNameEN, existingActionPlan.ResponsibleUser.LastNameEN),
			"target_date":        existingActionPlan.TargetDateString,
			"percentage":         existingActionPlan.Percentage,
			"action_plan_status": existingActionPlan.ActionPlanStatus.TitleEn,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := cacaps.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.ComplianceActionPlan{}.TableName(),
				id,
				changedData,
				currentUserID,
				"compliance",
				"compliance_assessment",
				"action_plan",
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
		if oldResponsibleUserID != existingComplianceAssessmentControlActionPlan.ResponsibleUserID {
			urlParam := fmt.Sprintf(`{"compliance_assessment_id":%d,"control_id":%d, "action_plan_id":%d}`, complianceAssessmentID, controlID, id)
			queryParam := ""
			entity := "compliance-assessment-action-plan"

			oldResponsibleUserAction := "responsible-user-remove"
			oldResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{oldResponsibleUserID},
				ModuleIdentifier:     "compliance",
				NotificationType:     string(commonDTOs.InformationType),
				HeadingEn:            "Compliance Action Plan Responsibility Removed",
				HeadingAr:            "تمت إزالة مسؤولية خطة عمل الامتثال",
				TitleEn:              existingComplianceAssessmentControlActionPlan.Title,
				TitleAr:              existingComplianceAssessmentControlActionPlan.Title,
				Entity:               &entity,
				Action:               &oldResponsibleUserAction,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			newResponsibleUserAction := "responsible-user-assign"
			newResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{ComplianceActionPlanFormDTO.ResponsibleUserID},
				ModuleIdentifier:     "compliance",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Compliance Action Plan Responsibility Assigned",
				HeadingAr:            "تم تعيين مسؤولية خطة عمل الامتثال",
				TitleEn:              existingComplianceAssessmentControlActionPlan.Title,
				TitleAr:              existingComplianceAssessmentControlActionPlan.Title,
				Entity:               &entity,
				Action:               &newResponsibleUserAction,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			oldResponsibleUser, err := cacaps.userSvc.GetByID(oldResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch old responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			newResponsibleUser, err := cacaps.userSvc.GetByID(ComplianceActionPlanFormDTO.ResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch new responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			assignedBy, err := cacaps.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch assigning user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			appURL := utility.GetAppURL()

			// Get assessment title for emails
			var assessmentTitle string
			if len(complianceAssessment.Compliances) > 0 && complianceAssessment.Compliances[0].ComplianceRegister != nil {
				assessmentTitle = complianceAssessment.Compliances[0].ComplianceRegister.TitleEn
			} else if complianceAssessment.ResponsibleOrganizationUnit != nil {
				assessmentTitle = complianceAssessment.ResponsibleOrganizationUnit.TitleEn + "'s Assessment"
			}

			oldResponsibleUserEmailTo := []string{oldResponsibleUser.Email}
			oldResponsibleUserSubject := fmt.Sprintf("Compliance Assessment Requirement Action Plan%s: Responsibility Removed", existingComplianceAssessmentControlActionPlan.Title)
			oldResponsibleUserTemplateRelativePath := "compliance/templates/email/compliance_assessment_action_plan_remove_responsible_user_mail.html"
			oldResponsibleUserEmailData := map[string]interface{}{
				"FirstNameEn":           oldResponsibleUser.FirstNameEn,
				"LastNameEn":            oldResponsibleUser.LastNameEn,
				"ComplianceAssessment":  assessmentTitle,
				"ActionPlanTitle":       existingComplianceAssessmentControlActionPlan.Title,
				"ActionPlanDescription": existingComplianceAssessmentControlActionPlan.Description,
				"RedirectURL":           fmt.Sprintf("%s/compliance/assessments/%d/overview", appURL, complianceAssessmentID),
			}

			oldResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              oldResponsibleUserEmailTo,
				Subject:              oldResponsibleUserSubject,
				TemplateRelativePath: oldResponsibleUserTemplateRelativePath,
				EmailData:            oldResponsibleUserEmailData,
			}

			newResponsibleUserEmailTo := []string{newResponsibleUser.Email}
			newResponsibleUserSubject := fmt.Sprintf("Compliance Assessment Requirement Action Plan%s: Responsibility Assigned", existingComplianceAssessmentControlActionPlan.Title)
			newResponsibleUserTemplateRelativePath := "compliance/templates/email/compliance_assessment_action_plan_change_responsible_user_mail.html"
			newResponsibleUserEmailData := map[string]interface{}{
				"FirstNameEn":           newResponsibleUser.FirstNameEn,
				"LastNameEn":            newResponsibleUser.LastNameEn,
				"ComplianceAssessment":  assessmentTitle,
				"ActionPlanTitle":       existingComplianceAssessmentControlActionPlan.Title,
				"ActionPlanDescription": existingComplianceAssessmentControlActionPlan.Description,
				"AssignedByFirstNameEn": assignedBy.FirstNameEn,
				"AssignedByLastNameEn":  assignedBy.LastNameEn,
				"AssignedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           fmt.Sprintf("%s/compliance/assessments/%d/overview", appURL, complianceAssessmentID),
			}

			newResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              newResponsibleUserEmailTo,
				Subject:              newResponsibleUserSubject,
				TemplateRelativePath: newResponsibleUserTemplateRelativePath,
				EmailData:            newResponsibleUserEmailData,
			}

			err = cacaps.notificationSvc.SendNotification(
				newResponsibleUserNotificationDTO, newResponsibleUserEmailNotificationDTO, "compliance_action_plan",
			)
			if err != nil {
				utility.LogError("Failed to send new responsible user notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}

			err = cacaps.notificationSvc.SendNotification(
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

func (cacaps *ComplianceAssessmentControlActionPlanService) Delete(complianceAssessmentID uint, controlID uint, ids []uint, currentUserID uint) ([]uint, error) {
	if _, err := cacaps.complianceAssessmentControlSvc.GetByID(controlID); err != nil {
		return nil, err
	}

	var failedIds []uint
	var successfulIds []uint
	var deletedComplianceAssessmentControlActionPlan []map[string]interface{}
	for _, id := range ids {

		complianceAssessmentControlActionPlan, err := cacaps.GetFullDetailsByID(complianceAssessmentID, controlID, id)
		if err != nil {
			continue
		}

		if len(complianceAssessmentControlActionPlan.ActionPlanStatusUpdates) > 0 {
			var deleteStatusUpdateIds []*uint
			var deleteDocumentIds []*uint
			for _, update := range complianceAssessmentControlActionPlan.ActionPlanStatusUpdates {
				deleteStatusUpdateIds = append(deleteStatusUpdateIds, &update.ID)
				for _, document := range update.Documents {
					deleteDocumentIds = append(deleteDocumentIds, &document.ID)
				}
			}
			if err := cacaps.repo.DeleteComplianceActionPlanStatusUpdateDocuments(deleteDocumentIds); err != nil {
				failedIds = append(failedIds, id)
				continue
			}
			if err := cacaps.repo.DeleteComplianceActionPlanStatusUpdates(deleteStatusUpdateIds); err != nil {
				failedIds = append(failedIds, id)
				continue
			}
		}

		if len(complianceAssessmentControlActionPlan.ActionPlanComments) > 0 {
			var deleteCommentIds []*uint
			var deleteCommentDocumentIds []*uint
			for _, update := range complianceAssessmentControlActionPlan.ActionPlanComments {
				deleteCommentIds = append(deleteCommentIds, &update.ID)
				for _, document := range update.Documents {
					deleteCommentDocumentIds = append(deleteCommentDocumentIds, &document.ID)
				}
			}
			if err := cacaps.repo.DeleteComplianceActionPlanCommentDocuments(deleteCommentDocumentIds); err != nil {
				failedIds = append(failedIds, id)
				continue
			}
			if err := cacaps.repo.DeleteComplianceActionPlanComments(deleteCommentIds); err != nil {
				failedIds = append(failedIds, id)
				continue
			}
		}

		existingTask, err := cacaps.taskSvc.GetByFieldID("compliance_action_plan_id", id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		taskID, err := cacaps.taskSvc.Delete(existingTask.ID, currentUserID)
		if err != nil {
			failedIds = append(failedIds, taskID)
		}

		err = cacaps.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedComplianceAssessmentControlActionPlan = append(deletedComplianceAssessmentControlActionPlan, map[string]interface{}{
			"id":                     id,
			"compliance_register_id": controlID,
			"title":                  complianceAssessmentControlActionPlan.Title,
		})

	}

	if len(failedIds) > 0 {
		return failedIds, nil
	}

	if len(successfulIds) > 0 {
		go func() {
			err := cacaps.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.ComplianceActionPlan{}.TableName(),
				successfulIds,
				deletedComplianceAssessmentControlActionPlan,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (cacaps *ComplianceAssessmentControlActionPlanService) StatusUpdate(assessmentID uint, controlID uint, id uint, currentUserID uint, complianceActionPlanStatusUpdateFormDTO *complianceDTOs.ComplianceActionPlanStatusUpdateFormDTO, tx utility.DatabaseTransaction) error {
	existingComplianceActionPlan, err := cacaps.GetByID(assessmentID, controlID, id)
	if err != nil {
		return err
	}

	existingComplianceActionPlan.ActionPlanStatusID = complianceActionPlanStatusUpdateFormDTO.ActionPlanStatusID
	existingComplianceActionPlan.Comments = &complianceActionPlanStatusUpdateFormDTO.Comment
	existingComplianceActionPlan.Percentage = complianceActionPlanStatusUpdateFormDTO.Percentage
	existingComplianceActionPlan.UpdatedByID = &currentUserID

	err = cacaps.repo.Update(existingComplianceActionPlan, tx)
	if err != nil {
		return err
	}

	isReviewUserNull := existingComplianceActionPlan.ReviewUserID == nil

	// Fetch the corresponding task and update its status
	existingTask, err := cacaps.taskSvc.GetByFieldID("compliance_action_plan_id", id)
	if err != nil {
		return err
	}
	err = cacaps.taskSvc.UpdateStatus(currentUserID, existingTask.ID, complianceActionPlanStatusUpdateFormDTO.ActionPlanStatusID, isReviewUserNull, tx)
	if err != nil {
		return err
	}

	return nil
}

func (cacaps *ComplianceAssessmentControlActionPlanService) validateComplianceActionPlanFormDTO(ComplianceActionPlanFormDTO *complianceDTOs.ComplianceActionPlanFormDTO) error {
	if ComplianceActionPlanFormDTO.Percentage == nil {
		defaultValue := int(0)
		ComplianceActionPlanFormDTO.Percentage = &defaultValue
	}

	return nil
}

func (cacaps *ComplianceAssessmentControlActionPlanService) parseComplianceAssessmentControlActionPlanDates(dto *complianceDTOs.ComplianceActionPlanFormDTO) (*complianceDTOs.ComplianceActionPlanParsedDates, error) {
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
