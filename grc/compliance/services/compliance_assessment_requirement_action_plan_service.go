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

type ComplianceAssessmentRequirementActionPlanServiceInterface interface {
	GetAll(complianceAssessmentID uint, requirementID uint, params *complianceDTOs.ComplianceActionPlanListQueryParamsDTO) (*complianceDTOs.ComplianceActionPlanListResponseDTO, error)
	GetByID(complianceAssessmentID uint, requirementID uint, id uint) (*models.ComplianceActionPlan, error)
	GetFullDetailsByID(complianceAssessmentID uint, requirementID uint, id uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceActionPlanDetailDTO, error)
	Create(complianceAssessmentID uint, requirementID uint, currentUserID uint, ComplianceActionPlanFormDTO *complianceDTOs.ComplianceActionPlanFormDTO) (uint, error)
	Update(complianceAssessmentID uint, requirementID uint, id uint, currentUserID uint, compliance *complianceDTOs.ComplianceActionPlanFormDTO) error
	Delete(complianceAssessmentID uint, requirementID uint, ids []uint, currentUserID uint) ([]uint, error)

	StatusUpdate(complianceAssessmentID uint, requirementID uint, id uint, currentUserID uint, complianceActionPlanStatusUpdateFormDTO *complianceDTOs.ComplianceActionPlanStatusUpdateFormDTO, tx utility.DatabaseTransaction) error
}

// Service and implementations
type ComplianceAssessmentRequirementActionPlanService struct {
	repo                               complianceRepositories.ComplianceActionPlanRepositoryInterface
	activityLogSvc                     logServices.ActivityLogServiceInterface
	complianceAssessmentSvc            ComplianceAssessmentServiceInterface
	complianceAssessmentRequirementSvc ComplianceAssessmentRequirementServiceInterface
	actionPlanStatusSvc                masterServices.ActionPlanStatusServiceInterface
	userSvc                            organizationServices.UserServiceInterface
	organizationUnitSvc                settingService.OrganizationUnitServiceInterface
	taskSvc                            taskServices.TaskServiceInterface
	notificationSvc                    commonServices.NotificationServiceInterface
}

func NewComplianceAssessmentRequirementActionPlanService(
	repo complianceRepositories.ComplianceActionPlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	complianceAssessmentSvc ComplianceAssessmentServiceInterface,
	complianceAssessmentRequirementSvc ComplianceAssessmentRequirementServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	organizationUnitSvc settingService.OrganizationUnitServiceInterface,
	taskSvc taskServices.TaskServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) ComplianceAssessmentRequirementActionPlanServiceInterface {
	return &ComplianceAssessmentRequirementActionPlanService{
		repo:                               repo,
		activityLogSvc:                     activityLogSvc,
		complianceAssessmentSvc:            complianceAssessmentSvc,
		complianceAssessmentRequirementSvc: complianceAssessmentRequirementSvc,
		actionPlanStatusSvc:                actionPlanStatusSvc,
		userSvc:                            userSvc,
		organizationUnitSvc:                organizationUnitSvc,
		taskSvc:                            taskSvc,
		notificationSvc:                    notificationSvc,
	}
}

func (caraps *ComplianceAssessmentRequirementActionPlanService) GetAll(complianceAssessmentID uint, requirementID uint, params *complianceDTOs.ComplianceActionPlanListQueryParamsDTO) (*complianceDTOs.ComplianceActionPlanListResponseDTO, error) {
	if requirementID != 0 {
		if _, err := caraps.complianceAssessmentRequirementSvc.GetByID(requirementID); err != nil {
			return nil, err
		}
	}
	return caraps.repo.GetAll(uint(0), params, requirementID)
}

func (caraps ComplianceAssessmentRequirementActionPlanService) GetByID(complianceAssessmentID uint, requirementID uint, id uint) (*models.ComplianceActionPlan, error) {
	if _, err := caraps.complianceAssessmentRequirementSvc.GetByID(requirementID); err != nil {
		return nil, err
	}

	complianceAssessmentRequirementActionPlan, err := caraps.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return complianceAssessmentRequirementActionPlan, nil
}

func (caraps ComplianceAssessmentRequirementActionPlanService) GetFullDetailsByID(complianceAssessmentID uint, requirementID uint, id uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceActionPlanDetailDTO, error) {
	var complianceAssessmentRequirementActionPlanDTO complianceDTOs.ComplianceActionPlanDetailDTO

	// Check if the compliance register exists
	if _, err := caraps.complianceAssessmentRequirementSvc.GetByID(requirementID); err != nil {
		return nil, err
	}

	// Get full details of the compliance action plan
	complianceAssessmentRequirementActionPlan, err := caraps.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Decode the compliance action plan into the DTO
	err = mapstructure.Decode(complianceAssessmentRequirementActionPlan, &complianceAssessmentRequirementActionPlanDTO)
	if err != nil {
		return nil, err
	}
	if complianceAssessmentRequirementActionPlan.TargetDate != nil {
		targetDate := complianceAssessmentRequirementActionPlan.TargetDate.Format("2006-01-02")
		complianceAssessmentRequirementActionPlanDTO.TargetDateString = &targetDate
	}

	responsibleUnitID, err := caraps.userSvc.GetUserPrimaryUnitID(complianceAssessmentRequirementActionPlan.ResponsibleUserID)
	if err != nil {
		return nil, err
	}

	responsibleUnit, err := caraps.organizationUnitSvc.GetFullDetailsByID(responsibleUnitID)
	if err != nil {
		return nil, err
	}
	complianceAssessmentRequirementActionPlanDTO.ResponsibleUserUnitID = &responsibleUnitID
	complianceAssessmentRequirementActionPlanDTO.ResponsibleUserUnit = responsibleUnit

	if complianceAssessmentRequirementActionPlan.ReviewUserID != nil {
		reviewUserUnitID, err := caraps.userSvc.GetUserPrimaryUnitID(*complianceAssessmentRequirementActionPlan.ReviewUserID)
		if err != nil {
			return nil, err
		}

		reviewUserUnit, err := caraps.organizationUnitSvc.GetFullDetailsByID(reviewUserUnitID)
		if err != nil {
			return nil, err
		}
		complianceAssessmentRequirementActionPlanDTO.ReviewUserUnitID = &reviewUserUnitID
		complianceAssessmentRequirementActionPlanDTO.ReviewUserUnit = reviewUserUnit
	}

	complianceAssessmentRequirementActionPlanDTO.CreatedAt = complianceAssessmentRequirementActionPlan.CreatedAt
	complianceAssessmentRequirementActionPlanDTO.UpdatedAt = complianceAssessmentRequirementActionPlan.UpdatedAt

	// TODO: Add any additional fields required for the details screen

	return &complianceAssessmentRequirementActionPlanDTO, nil
}

func (caraps *ComplianceAssessmentRequirementActionPlanService) Create(complianceAssessmentID uint, requirementID uint, currentUserID uint, ComplianceActionPlanFormDTO *complianceDTOs.ComplianceActionPlanFormDTO) (uint, error) {
	tx := caraps.repo.BeginTransaction()

	if err := caraps.validateComplianceActionPlanFormDTO(ComplianceActionPlanFormDTO); err != nil {
		return 0, err
	}

	dates, err := caraps.parseComplianceAssessmentRequirementActionPlanDates(ComplianceActionPlanFormDTO)
	if err != nil {
		return 0, err
	}

	complianceAssessment, err := caraps.complianceAssessmentSvc.GetFullDetailsByID(complianceAssessmentID)
	if err != nil {
		return 0, err
	}

	actionPlanStatus, err := caraps.actionPlanStatusSvc.GetByType("new")
	if err != nil {
		return 0, err
	}

	complianceAssessmentRequirementActionPlan := &models.ComplianceActionPlan{
		ComplianceRegisterID:              nil,
		ComplianceAssessmentRequirementID: &requirementID,
		Title:                             ComplianceActionPlanFormDTO.Title,
		Description:                       ComplianceActionPlanFormDTO.Description,
		ResponsibleUserID:                 ComplianceActionPlanFormDTO.ResponsibleUserID,
		TargetDate:                        dates.TargetDate,
		Percentage:                        ComplianceActionPlanFormDTO.Percentage,
		ActionPlanStatusID:                actionPlanStatus.ID,
		SeverityLevelID:                   ComplianceActionPlanFormDTO.SeverityLevelID,
		ReviewUserID:                      ComplianceActionPlanFormDTO.ReviewUserID,
		Comments:                          ComplianceActionPlanFormDTO.Comments,
		CreatedByID:                       currentUserID,
	}

	id, err := caraps.repo.Create(complianceAssessmentRequirementActionPlan, tx)
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

	_, err = caraps.repo.CreateInitialComplianceActionPlanStatusUpdate(complianceActionPlanUpdate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	urlParam := fmt.Sprintf(`{"compliance_assessment_id":%d,"requirementt_id":%d, "action_plan_id":%d}`, complianceAssessmentID, requirementID, id)
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

	if _, err := caraps.taskSvc.Create(taskDTO, currentUserID, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	go func() {
		existingComplianceAssessmentRequirementActionPlan, err := caraps.GetFullDetailsByID(complianceAssessmentID, requirementID, id)
		if err != nil {
			utility.LogError("Failed to fetch compliance action plan", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = caraps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ComplianceActionPlan{}.TableName(),
			id,
			existingComplianceAssessmentRequirementActionPlan,
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
		urlParam := fmt.Sprintf(`{"compliance_assessment_id":%d,"requirement_id":%d, "action_plan_id":%d}`, complianceAssessmentID, requirementID, id)
		queryParam := ""
		entity := "compliance-action-plan"
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

		responsibleUser, err := caraps.userSvc.GetByID(ComplianceActionPlanFormDTO.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch responsible user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		assignedBy, err := caraps.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch assigning user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("Compliance Assessment Requirement Action Plan%s: Responsibility Assigned", ComplianceActionPlanFormDTO.Title)
		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstNameEn":           responsibleUser.FirstNameEn,
			"LastNameEn":            responsibleUser.LastNameEn,
			"ComplianceAssessment":  complianceAssessment.ComplianceRegister.TitleEn,
			"ActionPlanTitle":       ComplianceActionPlanFormDTO.Title,
			"ActionPlanDescription": ComplianceActionPlanFormDTO.Description,
			"AssignedByFirstNameEn": assignedBy.FirstNameEn,
			"AssignedByLastNameEn":  assignedBy.LastNameEn,
			"AssignedOn":            time.Now().Format("2006-01-02"),
			"RedirectURL":           fmt.Sprintf("%scompliance/assessments/%d/overview", appURL, complianceAssessmentID),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "compliance/templates/email/compliance_assessment_action_plan_create_responsible_user_mail.html",
			EmailData:            emailData,
		}

		// send notification and mail
		err = caraps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "compliance_action_plan")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	return id, nil
}

func (caraps *ComplianceAssessmentRequirementActionPlanService) Update(complianceAssessmentID uint, requirementID uint, id uint, currentUserID uint, ComplianceActionPlanFormDTO *complianceDTOs.ComplianceActionPlanFormDTO) error {
	tx := caraps.repo.BeginTransaction()

	existingComplianceAssessmentRequirementActionPlan, err := caraps.GetByID(complianceAssessmentID, requirementID, id)
	if err != nil {
		return err
	}

	oldComplianceAssessmentRequirementActionPlan, err := caraps.GetFullDetailsByID(complianceAssessmentID, requirementID, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	complianceAssessment, err := caraps.complianceAssessmentSvc.GetFullDetailsByID(complianceAssessmentID)
	if err != nil {
		return err
	}

	oldResponsibleUserID := existingComplianceAssessmentRequirementActionPlan.ResponsibleUserID

	if err := caraps.validateComplianceActionPlanFormDTO(ComplianceActionPlanFormDTO); err != nil {
		return err
	}

	dates, err := caraps.parseComplianceAssessmentRequirementActionPlanDates(ComplianceActionPlanFormDTO)
	if err != nil {
		return err
	}

	// Update Compliance ActionPlan
	existingComplianceAssessmentRequirementActionPlan.Title = ComplianceActionPlanFormDTO.Title
	existingComplianceAssessmentRequirementActionPlan.Description = ComplianceActionPlanFormDTO.Description
	existingComplianceAssessmentRequirementActionPlan.ResponsibleUserID = ComplianceActionPlanFormDTO.ResponsibleUserID
	existingComplianceAssessmentRequirementActionPlan.TargetDate = dates.TargetDate
	existingComplianceAssessmentRequirementActionPlan.Percentage = ComplianceActionPlanFormDTO.Percentage
	existingComplianceAssessmentRequirementActionPlan.ActionPlanStatusID = ComplianceActionPlanFormDTO.ActionPlanStatusID
	existingComplianceAssessmentRequirementActionPlan.SeverityLevelID = ComplianceActionPlanFormDTO.SeverityLevelID
	existingComplianceAssessmentRequirementActionPlan.Comments = ComplianceActionPlanFormDTO.Comments
	existingComplianceAssessmentRequirementActionPlan.ReviewUserID = ComplianceActionPlanFormDTO.ReviewUserID

	existingComplianceAssessmentRequirementActionPlan.UpdatedByID = &currentUserID

	err = caraps.repo.Update(existingComplianceAssessmentRequirementActionPlan, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// //update task
	urlParam := fmt.Sprintf(`{"compliance_assessment_id":%d,"requirement_id":%d, "action_plan_id":%d}`, complianceAssessmentID, requirementID, id)

	taskUpdateDTO := commonDTOs.TaskFormDTO{
		Title:             ComplianceActionPlanFormDTO.Title,
		Description:       ComplianceActionPlanFormDTO.Description,
		URLParam:          &urlParam,
		SeverityLevelID:   ComplianceActionPlanFormDTO.SeverityLevelID,
		ResponsibleUserID: ComplianceActionPlanFormDTO.ResponsibleUserID,
	}

	err = caraps.taskSvc.Update(currentUserID, taskUpdateDTO, "compliance_action_plan_id", id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		existingActionPlan, err := caraps.GetFullDetailsByID(complianceAssessmentID, requirementID, id)
		if err != nil {
			return
		}
		// log changes
		changedData := make(map[string]interface{})
		oldData := map[string]interface{}{
			"title":              oldComplianceAssessmentRequirementActionPlan.Title,
			"description":        oldComplianceAssessmentRequirementActionPlan.Description,
			"responsible_user":   fmt.Sprintf("%s %s", oldComplianceAssessmentRequirementActionPlan.ResponsibleUser.FirstNameEN, oldComplianceAssessmentRequirementActionPlan.ResponsibleUser.LastNameEN),
			"target_date":        oldComplianceAssessmentRequirementActionPlan.TargetDateString,
			"percentage":         oldComplianceAssessmentRequirementActionPlan.Percentage,
			"action_plan_status": oldComplianceAssessmentRequirementActionPlan.ActionPlanStatus.TitleEn,
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
			err := caraps.activityLogSvc.LogActivity(
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
		if oldResponsibleUserID != existingComplianceAssessmentRequirementActionPlan.ResponsibleUserID {
			urlParam := fmt.Sprintf(`{"compliance_assessment_id":%d,"requirement_id":%d, "action_plan_id":%d}`, complianceAssessmentID, requirementID, id)
			queryParam := ""
			entity := "compliance-action-plan"

			oldResponsibleUserAction := "responsible-user-remove"
			oldResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{oldResponsibleUserID},
				ModuleIdentifier:     "compliance",
				NotificationType:     string(commonDTOs.InformationType),
				HeadingEn:            "Compliance Action Plan Responsibility Removed",
				HeadingAr:            "تمت إزالة مسؤولية خطة عمل الامتثال",
				TitleEn:              existingComplianceAssessmentRequirementActionPlan.Title,
				TitleAr:              existingComplianceAssessmentRequirementActionPlan.Title,
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
				TitleEn:              existingComplianceAssessmentRequirementActionPlan.Title,
				TitleAr:              existingComplianceAssessmentRequirementActionPlan.Title,
				Entity:               &entity,
				Action:               &newResponsibleUserAction,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			oldResponsibleUser, err := caraps.userSvc.GetByID(oldResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch old responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			newResponsibleUser, err := caraps.userSvc.GetByID(ComplianceActionPlanFormDTO.ResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch new responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			assignedBy, err := caraps.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch assigning user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			appURL := utility.GetAppURL()

			oldResponsibleUserEmailTo := []string{oldResponsibleUser.Email}
			oldResponsibleUserSubject := fmt.Sprintf("Compliance Assessment Requirement Action Plan%s: Responsibility Removed", existingComplianceAssessmentRequirementActionPlan.Title)
			oldResponsibleUserTemplateRelativePath := "compliance/templates/email/compliance_assessment_action_plan_remove_responsible_user_mail.html"
			oldResponsibleUserEmailData := map[string]interface{}{
				"FirstNameEn":           oldResponsibleUser.FirstNameEn,
				"LastNameEn":            oldResponsibleUser.LastNameEn,
				"ActionPlanTitle":       existingComplianceAssessmentRequirementActionPlan.Title,
				"ActionPlanDescription": existingComplianceAssessmentRequirementActionPlan.Description,
				"RedirectURL":           fmt.Sprintf("%scompliance/assessments/%d/overview", appURL, complianceAssessmentID),
			}

			oldResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              oldResponsibleUserEmailTo,
				Subject:              oldResponsibleUserSubject,
				TemplateRelativePath: oldResponsibleUserTemplateRelativePath,
				EmailData:            oldResponsibleUserEmailData,
			}

			newResponsibleUserEmailTo := []string{newResponsibleUser.Email}
			newResponsibleUserSubject := fmt.Sprintf("Compliance Assessment Requirement Action Plan%s: Responsibility Assigned", existingComplianceAssessmentRequirementActionPlan.Title)
			newResponsibleUserTemplateRelativePath := "compliance/templates/email/compliance_assessment_action_plan_change_responsible_user_mail.html"
			newResponsibleUserEmailData := map[string]interface{}{
				"FirstNameEn":           newResponsibleUser.FirstNameEn,
				"LastNameEn":            newResponsibleUser.LastNameEn,
				"ComplianceAssessment":  complianceAssessment.ComplianceRegister.TitleEn,
				"ActionPlanTitle":       existingComplianceAssessmentRequirementActionPlan.Title,
				"ActionPlanDescription": existingComplianceAssessmentRequirementActionPlan.Description,
				"AssignedByFirstNameEn": assignedBy.FirstNameEn,
				"AssignedByLastNameEn":  assignedBy.LastNameEn,
				"AssignedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           fmt.Sprintf("%scompliance/assessments/%d/overview", appURL, complianceAssessmentID),
			}

			newResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              newResponsibleUserEmailTo,
				Subject:              newResponsibleUserSubject,
				TemplateRelativePath: newResponsibleUserTemplateRelativePath,
				EmailData:            newResponsibleUserEmailData,
			}

			err = caraps.notificationSvc.SendNotification(
				newResponsibleUserNotificationDTO, newResponsibleUserEmailNotificationDTO, "compliance_action_plan",
			)
			if err != nil {
				utility.LogError("Failed to send new responsible user notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}

			err = caraps.notificationSvc.SendNotification(
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

func (caraps *ComplianceAssessmentRequirementActionPlanService) Delete(complianceAssessmentID uint, requirementID uint, ids []uint, currentUserID uint) ([]uint, error) {
	if _, err := caraps.complianceAssessmentRequirementSvc.GetByID(requirementID); err != nil {
		return nil, err
	}

	var failedIds []uint
	var successfulIds []uint
	var deletedComplianceAssessmentRequirementActionPlan []map[string]interface{}
	for _, id := range ids {

		complianceAssessmentRequirementActionPlan, err := caraps.GetFullDetailsByID(complianceAssessmentID, requirementID, id)
		if err != nil {
			continue
		}

		err = caraps.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedComplianceAssessmentRequirementActionPlan = append(deletedComplianceAssessmentRequirementActionPlan, map[string]interface{}{
			"id":                     id,
			"compliance_register_id": requirementID,
			"title":                  complianceAssessmentRequirementActionPlan.Title,
		})

	}

	if len(failedIds) > 0 {
		return failedIds, nil
	}

	if len(successfulIds) > 0 {
		go func() {
			err := caraps.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.ComplianceActionPlan{}.TableName(),
				successfulIds,
				deletedComplianceAssessmentRequirementActionPlan,
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

func (caraps *ComplianceAssessmentRequirementActionPlanService) StatusUpdate(assessmentID uint, requirementID uint, id uint, currentUserID uint, complianceActionPlanStatusUpdateFormDTO *complianceDTOs.ComplianceActionPlanStatusUpdateFormDTO, tx utility.DatabaseTransaction) error {
	existingComplianceActionPlan, err := caraps.GetByID(assessmentID, requirementID, id)
	if err != nil {
		return err
	}

	existingComplianceActionPlan.ActionPlanStatusID = complianceActionPlanStatusUpdateFormDTO.ActionPlanStatusID
	existingComplianceActionPlan.Comments = &complianceActionPlanStatusUpdateFormDTO.Comment
	existingComplianceActionPlan.Percentage = complianceActionPlanStatusUpdateFormDTO.Percentage
	existingComplianceActionPlan.UpdatedByID = &currentUserID

	err = caraps.repo.Update(existingComplianceActionPlan, tx)
	if err != nil {
		return err
	}

	isReviewUserNull := existingComplianceActionPlan.ReviewUserID == nil

	// Fetch the corresponding task and update its status
	existingTask, err := caraps.taskSvc.GetByFieldID("compliance_action_plan_id", id)
	if err != nil {
		return err
	}
	err = caraps.taskSvc.UpdateStatus(currentUserID, existingTask.ID, complianceActionPlanStatusUpdateFormDTO.ActionPlanStatusID, isReviewUserNull, tx)
	if err != nil {
		return err
	}

	return nil
}

func (caraps *ComplianceAssessmentRequirementActionPlanService) validateComplianceActionPlanFormDTO(ComplianceActionPlanFormDTO *complianceDTOs.ComplianceActionPlanFormDTO) error {
	if ComplianceActionPlanFormDTO.Percentage == nil {
		defaultValue := int(0)
		ComplianceActionPlanFormDTO.Percentage = &defaultValue
	}

	return nil
}

func (caraps *ComplianceAssessmentRequirementActionPlanService) parseComplianceAssessmentRequirementActionPlanDates(dto *complianceDTOs.ComplianceActionPlanFormDTO) (*complianceDTOs.ComplianceActionPlanParsedDates, error) {
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
