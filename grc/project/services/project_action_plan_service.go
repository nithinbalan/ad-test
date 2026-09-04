package services

import (
	"errors"
	"fmt"
	"time"

	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	projectDTOs "jethur/grc/project/dtos"
	projectRepositories "jethur/grc/project/repositories"
	taskServices "jethur/grc/task/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ProjectActionPlanServiceInterface interface {
	GetAll(projectID uint, params *projectDTOs.ProjectActionPlanListQueryParamsDTO) (*projectDTOs.ProjectActionPlanListResponseDTO, error)
	GetByID(projectID uint, id uint) (*models.ProjectActionPlan, error)
	GetFullDetailsByID(projectID uint, id uint, tx ...utility.DatabaseTransaction) (*projectDTOs.ProjectActionPlanDetailDTO, error)
	Create(projectID uint, currentUserID uint, ProjectActionPlanFormDTO *projectDTOs.ProjectActionPlanFormDTO) (uint, error)
	Update(projectID uint, id uint, currentUserID uint, project *projectDTOs.ProjectActionPlanFormDTO) error
	Delete(projectID uint, ids []uint, currentUserID uint) ([]uint, error)
	StatusUpdate(projectID uint, id uint, currentUserID uint, projectActionPlanStatusUpdateFormDTO *projectDTOs.ProjectActionPlanStatusUpdateFormDTO, tx utility.DatabaseTransaction) error
}

// Service and implementations
type ProjectActionPlanService struct {
	repo                projectRepositories.ProjectActionPlanRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface
	userSvc             organizationServices.UserServiceInterface
	projectSvc          ProjectServiceInterface
	projectMilestoneSvc ProjectMilestoneServiceInterface
	severityLevelSvc    masterServices.SeverityLevelServiceInterface
	delegateeUserSvc    commonServices.DelegateeUserserviceInterface
	taskSvc             taskServices.TaskServiceInterface
	notificationSvc     commonServices.NotificationServiceInterface
}

func NewProjectActionPlanService(
	repo projectRepositories.ProjectActionPlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	projectSvc ProjectServiceInterface,
	projectMilestoneSvc ProjectMilestoneServiceInterface,
	severityLevelSvc masterServices.SeverityLevelServiceInterface,
	delegateeUserSvc commonServices.DelegateeUserserviceInterface,
	taskSvc taskServices.TaskServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) ProjectActionPlanServiceInterface {
	return &ProjectActionPlanService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		actionPlanStatusSvc: actionPlanStatusSvc,
		userSvc:             userSvc,
		projectSvc:          projectSvc,
		projectMilestoneSvc: projectMilestoneSvc,
		severityLevelSvc:    severityLevelSvc,
		delegateeUserSvc:    delegateeUserSvc,
		taskSvc:             taskSvc,
		notificationSvc:     notificationSvc,
	}
}

func (paps *ProjectActionPlanService) GetAll(projectID uint, params *projectDTOs.ProjectActionPlanListQueryParamsDTO) (*projectDTOs.ProjectActionPlanListResponseDTO, error) {
	if projectID != 0 {
		if _, err := paps.projectSvc.GetByID(projectID); err != nil {
			return nil, err
		}
	}
	return paps.repo.GetAll(projectID, params)
}

func (paps ProjectActionPlanService) GetByID(projectID uint, id uint) (*models.ProjectActionPlan, error) {
	if _, err := paps.projectSvc.GetByID(projectID); err != nil {
		return nil, err
	}

	projectActionPlan, err := paps.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return projectActionPlan, nil
}

func (paps ProjectActionPlanService) GetFullDetailsByID(projectID uint, id uint, tx ...utility.DatabaseTransaction) (*projectDTOs.ProjectActionPlanDetailDTO, error) {
	var projectActionPlanDTO projectDTOs.ProjectActionPlanDetailDTO

	// Check if the PROJECT register exists
	if _, err := paps.projectSvc.GetByID(projectID); err != nil {
		return nil, err
	}

	// Get full details of the PROJECT action plan
	projectActionPlan, err := paps.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Decode the PROJECT action plan into the DTO
	err = mapstructure.Decode(projectActionPlan, &projectActionPlanDTO)
	if err != nil {
		return nil, err
	}
	if projectActionPlan.TargetDate != nil {
		targetDate := projectActionPlan.TargetDate.Format("2006-01-02")
		projectActionPlanDTO.TargetDateString = &targetDate
	}

	projectActionPlanDTO.CreatedAt = projectActionPlan.CreatedAt
	projectActionPlanDTO.UpdatedAt = projectActionPlan.UpdatedAt

	// TODO: Add any additional fields required for the details screen

	return &projectActionPlanDTO, nil
}

func (paps *ProjectActionPlanService) Create(projectID uint, currentUserID uint, ProjectActionPlanFormDTO *projectDTOs.ProjectActionPlanFormDTO) (uint, error) {
	tx := paps.repo.BeginTransaction()

	project, err := paps.projectSvc.GetByID(projectID)
	if err != nil {
		return 0, err
	}
	dates, err := paps.parseProjectActionPlanDates(ProjectActionPlanFormDTO)
	if err != nil {
		return 0, err
	}

	actionPlanStatus, err := paps.actionPlanStatusSvc.GetByType("new")
	if err != nil {
		return 0, err
	}

	projectActionPlan := &models.ProjectActionPlan{
		ProjectID:          &project.ID,
		ProjectMilestoneID: ProjectActionPlanFormDTO.ProjectMilestoneID,
		Title:              ProjectActionPlanFormDTO.Title,
		Description:        ProjectActionPlanFormDTO.Description,
		ResponsibleUserID:  ProjectActionPlanFormDTO.ResponsibleUserID,
		SeverityLevelID:    ProjectActionPlanFormDTO.SeverityLevelID,
		ReviewUserID:       ProjectActionPlanFormDTO.ReviewUserID,
		Comments:           ProjectActionPlanFormDTO.Comments,
		TargetDate:         dates.TargetDate,
		Percentage:         &ProjectActionPlanFormDTO.Percentage,
		ActionPlanStatusID: actionPlanStatus.ID,
		CreatedBYID:        currentUserID,
	}

	if ProjectActionPlanFormDTO.ResponsibleUserID != 0 {
		delegateeID, err := paps.delegateeUserSvc.AssignDelegateeIfPermitted(ProjectActionPlanFormDTO.ResponsibleUserID, "project_register")
		if err != nil {
			return 0, err
		}
		if delegateeID != nil {
			projectActionPlan.DelegateeUserID = delegateeID
		}
	}

	id, err := paps.repo.Create(projectActionPlan, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	urlParam := fmt.Sprintf(`{"project_id":%d, "action_plan_id":%d}`, projectID, id)
	taskDTO := commonDTOs.TaskFormDTO{
		ModuleIdentifier:    "organization",
		ProjectActionPlanID: &id,
		Title:               ProjectActionPlanFormDTO.Title,
		Description:         ProjectActionPlanFormDTO.Description,
		TargetDate:          dates.TargetDate,
		URLParam:            &urlParam,
		QueryParam:          nil,
		SeverityLevelID:     ProjectActionPlanFormDTO.SeverityLevelID,
		ResponsibleUserID:   ProjectActionPlanFormDTO.ResponsibleUserID,
	}

	if _, err := paps.taskSvc.Create(taskDTO, currentUserID, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	go func() {
		activityData := map[string]interface{}{
			"project_id":           projectID,
			"title":                ProjectActionPlanFormDTO.Title,
			"description":          ProjectActionPlanFormDTO.Description,
			"responsible_user_id":  ProjectActionPlanFormDTO.ResponsibleUserID,
			"severity_level_id":    ProjectActionPlanFormDTO.SeverityLevelID,
			"review_user_id":       ProjectActionPlanFormDTO.ReviewUserID,
			"target_date":          dates.TargetDate,
			"comments":             ProjectActionPlanFormDTO.Comments,
			"percentage":           ProjectActionPlanFormDTO.Percentage,
			"project_milestone_id": ProjectActionPlanFormDTO.ProjectMilestoneID,
		}
		err := paps.activityLogSvc.LogActivity(models.CreateAction, models.ProjectActionPlan{}.TableName(), id, activityData, currentUserID, "project", "project_register", "action_plan", nil)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ProjectActionPlan{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}

		// prepare notification
		urlParam := fmt.Sprintf(`{"project_id":%d, "action_plan_id":%d}`, projectID, id)
		queryParam := ""
		entity := "project-action-plan"
		action := "responsible-user-assign"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{ProjectActionPlanFormDTO.ResponsibleUserID},
			ModuleIdentifier:     "project",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "Project Action Plan Responsibility Assigned",
			HeadingAr:            "تم تعيين مسؤولية خطة عمل المشروع",
			TitleEn:              ProjectActionPlanFormDTO.Title,
			TitleAr:              ProjectActionPlanFormDTO.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			QueryParam:           queryParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		responsibleUser, err := paps.userSvc.GetByID(ProjectActionPlanFormDTO.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch responsible user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		assignedBy, err := paps.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch assigning user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("Project Action Plan/%s: Responsibility Assigned", ProjectActionPlanFormDTO.Title)
		appURL := utility.GetAppURL()
		targetDate := ""
		if dates.TargetDate != nil {
			targetDate = dates.TargetDate.Format("2006-01-02")
		}
		emailData := map[string]interface{}{
			"FirstNameEn":           responsibleUser.FirstNameEn,
			"LastNameEn":            responsibleUser.LastNameEn,
			"ProjectTitle":          project.TitleEn,
			"ProjectTitleAr":        project.TitleAr,
			"ActionPlanTitle":       ProjectActionPlanFormDTO.Title,
			"ActionPlanDescription": ProjectActionPlanFormDTO.Description,
			"TargetDate":            targetDate,
			"AssignedByFirstNameEn": assignedBy.FirstNameEn,
			"AssignedByLastNameEn":  assignedBy.LastNameEn,
			"AssignedOn":            time.Now().Format("2006-01-02"),
			"RedirectURL":           fmt.Sprintf("%s/project/%d/action-plans?action_plan_id=%d", appURL, projectID, id),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "project/templates/email/project_action_plan_create_responsible_user_mail.html",
			EmailData:            emailData,
		}

		// send notification and mail
		err = paps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "project_action_plan")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	return id, nil
}

func (paps *ProjectActionPlanService) Update(projectID uint, id uint, currentUserID uint, ProjectActionPlanFormDTO *projectDTOs.ProjectActionPlanFormDTO) error {
	tx := paps.repo.BeginTransaction()

	// Validate project milestone exists
	if ProjectActionPlanFormDTO.ProjectMilestoneID != nil {
		if _, err := paps.projectMilestoneSvc.GetByID(projectID, *ProjectActionPlanFormDTO.ProjectMilestoneID); err != nil {
			return err
		}
	}

	existingProjectActionPlan, err := paps.GetByID(projectID, id)
	if err != nil {
		return err
	}

	oldProjectActionPlan, err := paps.GetFullDetailsByID(projectID, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	oldResponsibleUserID := existingProjectActionPlan.ResponsibleUserID

	project, err := paps.projectSvc.GetByID(projectID)
	if err != nil {
		return err
	}

	if err := paps.validateProjectActionPlanFormDTO(&project.ID, ProjectActionPlanFormDTO); err != nil {
		return err
	}

	dates, err := paps.parseProjectActionPlanDates(ProjectActionPlanFormDTO)
	if err != nil {
		return err
	}

	// Update Project ActionPlan
	existingProjectActionPlan.Title = ProjectActionPlanFormDTO.Title
	existingProjectActionPlan.Description = ProjectActionPlanFormDTO.Description
	existingProjectActionPlan.ProjectMilestoneID = ProjectActionPlanFormDTO.ProjectMilestoneID
	existingProjectActionPlan.Comments = ProjectActionPlanFormDTO.Comments
	existingProjectActionPlan.ResponsibleUserID = ProjectActionPlanFormDTO.ResponsibleUserID
	existingProjectActionPlan.SeverityLevelID = ProjectActionPlanFormDTO.SeverityLevelID
	existingProjectActionPlan.ReviewUserID = ProjectActionPlanFormDTO.ReviewUserID
	existingProjectActionPlan.TargetDate = dates.TargetDate
	existingProjectActionPlan.UpdatedBYID = &currentUserID

	if ProjectActionPlanFormDTO.ResponsibleUserID != 0 && oldResponsibleUserID != ProjectActionPlanFormDTO.ResponsibleUserID {
		existingProjectActionPlan.DelegateeUserID = nil
		delegateeID, err := paps.delegateeUserSvc.AssignDelegateeIfPermitted(ProjectActionPlanFormDTO.ResponsibleUserID, "project_register")
		if err != nil {
			return err
		}
		if delegateeID != nil {
			existingProjectActionPlan.DelegateeUserID = delegateeID
		}
	}

	err = paps.repo.Update(existingProjectActionPlan, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// //update task
	urlParam := fmt.Sprintf(`{"project_id":%d, "action_plan_id":%d}`, projectID, id)

	taskUpdateDTO := commonDTOs.TaskFormDTO{
		Title:             ProjectActionPlanFormDTO.Title,
		Description:       ProjectActionPlanFormDTO.Description,
		TargetDate:        dates.TargetDate,
		URLParam:          &urlParam,
		SeverityLevelID:   ProjectActionPlanFormDTO.SeverityLevelID,
		ResponsibleUserID: ProjectActionPlanFormDTO.ResponsibleUserID,
	}

	err = paps.taskSvc.Update(currentUserID, taskUpdateDTO, "project_action_plan_id", id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	go func() {
		existingActionPlan, err := paps.GetFullDetailsByID(projectID, id)
		if err != nil {
			return
		}

		// log changes
		changedData := make(map[string]interface{})
		oldData := map[string]interface{}{
			"title":              oldProjectActionPlan.Title,
			"description":        oldProjectActionPlan.Description,
			"severity_level":     oldProjectActionPlan.SeverityLevel.TitleEn,
			"responsible_user":   fmt.Sprintf("%s %s", oldProjectActionPlan.ResponsibleUser.FirstNameEN, oldProjectActionPlan.ResponsibleUser.LastNameEN),
			"review_user":        utility.GetFieldValue(oldProjectActionPlan.ReviewUser),
			"target_date":        oldProjectActionPlan.TargetDateString,
			"percentage":         oldProjectActionPlan.Percentage,
			"comment":            oldProjectActionPlan.Comments,
			"action_plan_status": oldProjectActionPlan.ActionPlanStatus.TitleEn,
		}

		newData := map[string]interface{}{
			"title":              existingActionPlan.Title,
			"description":        existingActionPlan.Description,
			"comments":           existingActionPlan.Comments,
			"responsible_user":   fmt.Sprintf("%s %s", existingActionPlan.ResponsibleUser.FirstNameEN, existingActionPlan.ResponsibleUser.LastNameEN),
			"review_user":        utility.GetFieldValue(existingActionPlan.ReviewUser),
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
			err := paps.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.ProjectActionPlan{}.TableName(),
				id,
				changedData,
				currentUserID,
				"project", "project_register", "action_plan", nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.ProjectActionPlan{}.TableName(),
					"recordID":  id,
					"userID":    currentUserID,
					"error":     err.Error(),
				}
				utility.LogError("Failed to log activity", fields)
			}
		}

		// sends notification if responsible user has changed
		if oldResponsibleUserID != existingProjectActionPlan.ResponsibleUserID {
			urlParam := fmt.Sprintf(`{"project_id":%d, "action_plan_id":%d}`, projectID, id)
			queryParam := ""
			entity := "project-action-plan"

			oldResponsibleUserAction := "responsible-user-remove"
			oldResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{oldResponsibleUserID},
				ModuleIdentifier:     "project",
				NotificationType:     string(commonDTOs.InformationType),
				HeadingEn:            "Project Action Plan Responsibility Removed",
				HeadingAr:            "تمت إزالة مسؤولية خطة العمل",
				TitleEn:              existingProjectActionPlan.Title,
				TitleAr:              existingProjectActionPlan.Title,
				Entity:               &entity,
				Action:               &oldResponsibleUserAction,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			newResponsibleUserAction := "responsible-user-assign"
			newResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{ProjectActionPlanFormDTO.ResponsibleUserID},
				ModuleIdentifier:     "project",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Project Action Plan Responsibility Assigned",
				HeadingAr:            "تم تعيين مسؤولية خطة العمل",
				TitleEn:              existingProjectActionPlan.Title,
				TitleAr:              existingProjectActionPlan.Title,
				Entity:               &entity,
				Action:               &newResponsibleUserAction,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			oldResponsibleUser, err := paps.userSvc.GetByID(oldResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch old responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			newResponsibleUser, err := paps.userSvc.GetByID(ProjectActionPlanFormDTO.ResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch new responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			assignedBy, err := paps.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch assigning user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			appURL := utility.GetAppURL()

			oldResponsibleUserEmailTo := []string{oldResponsibleUser.Email}
			oldResponsibleUserSubject := fmt.Sprintf("Project Action Plan/%s: Responsibility Removed", existingProjectActionPlan.Title)
			oldResponsibleUserTemplateRelativePath := "project/templates/email/project_action_plan_remove_responsible_user_mail.html"
			oldResponsibleUserEmailData := map[string]interface{}{
				"FirstNameEn":           oldResponsibleUser.FirstNameEn,
				"LastNameEn":            oldResponsibleUser.LastNameEn,
				"ProjectTitle":          project.TitleEn,
				"ProjectTitleAr":        project.TitleAr,
				"ActionPlanTitle":       existingProjectActionPlan.Title,
				"ActionPlanDescription": existingProjectActionPlan.Description,
				"RedirectURL":           fmt.Sprintf("%s/project/%d/action-plans?action_plan_id=%d", appURL, projectID, id),
			}

			oldResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              oldResponsibleUserEmailTo,
				Subject:              oldResponsibleUserSubject,
				TemplateRelativePath: oldResponsibleUserTemplateRelativePath,
				EmailData:            oldResponsibleUserEmailData,
			}

			newResponsibleUserEmailTo := []string{newResponsibleUser.Email}
			newResponsibleUserSubject := fmt.Sprintf("Project Action Plan/%s: Responsibility Assigned", existingProjectActionPlan.Title)
			newResponsibleUserTemplateRelativePath := "project/templates/email/project_action_plan_change_responsible_user_mail.html"
			targetDate := ""
			if existingProjectActionPlan.TargetDate != nil {
				targetDate = existingProjectActionPlan.TargetDate.Format("2006-01-02")
			}
			newResponsibleUserEmailData := map[string]interface{}{
				"FirstNameEn":           newResponsibleUser.FirstNameEn,
				"LastNameEn":            newResponsibleUser.LastNameEn,
				"ProjectTitle":          project.TitleEn,
				"ProjectTitleAr":        project.TitleAr,
				"ActionPlanTitle":       existingProjectActionPlan.Title,
				"ActionPlanDescription": existingProjectActionPlan.Description,
				"TargetDate":            targetDate,
				"AssignedByFirstNameEn": assignedBy.FirstNameEn,
				"AssignedByLastNameEn":  assignedBy.LastNameEn,
				"AssignedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           fmt.Sprintf("%s/project/%d/action-plans?action_plan_id=%d", appURL, projectID, id),
			}

			newResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              newResponsibleUserEmailTo,
				Subject:              newResponsibleUserSubject,
				TemplateRelativePath: newResponsibleUserTemplateRelativePath,
				EmailData:            newResponsibleUserEmailData,
			}

			err = paps.notificationSvc.SendNotification(
				newResponsibleUserNotificationDTO, newResponsibleUserEmailNotificationDTO, "project_action_plan",
			)
			if err != nil {
				utility.LogError("Failed to send new responsible user notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}

			err = paps.notificationSvc.SendNotification(
				oldResponsibleUserNotificationDTO, oldResponsibleUserEmailNotificationDTO, "project_action_plan",
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

func (paps *ProjectActionPlanService) Delete(projectID uint, ids []uint, currentUserID uint) ([]uint, error) {
	if _, err := paps.projectSvc.GetByID(projectID); err != nil {
		return nil, err
	}

	var failedIds []uint
	var successfulIds []uint
	var deletedProjectActionPlans []map[string]interface{}

	for _, id := range ids {
		projectActionPlan, err := paps.repo.GetByID(id)
		if err != nil {
			fmt.Print(err)
			failedIds = append(failedIds, id)
			continue
		}

		err = paps.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedProjectActionPlans = append(deletedProjectActionPlans, map[string]interface{}{
			"id":    id,
			"title": projectActionPlan.Title,
		})
	}

	if len(successfulIds) > 0 {
		go func() {
			err := paps.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.ProjectActionPlan{}.TableName(),
				successfulIds,
				deletedProjectActionPlans,
				currentUserID,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to log bulk delete project action plan activity", map[string]interface{}{
					"recordIDs": successfulIds,
					"error":     err.Error(),
				})
			}
		}()
	}

	return failedIds, nil
}

func (paps *ProjectActionPlanService) validateProjectActionPlanFormDTO(projectID *uint, ProjectActionPlanFormDTO *projectDTOs.ProjectActionPlanFormDTO) error {
	// Validate responsible user exists
	if _, err := paps.userSvc.GetByID(ProjectActionPlanFormDTO.ResponsibleUserID); err != nil {
		return utility.NewValidationError(fmt.Errorf("responsible user does not exist"), "not_exist", "responsible user")
	}

	if ProjectActionPlanFormDTO.SeverityLevelID != nil {
		if _, err := paps.severityLevelSvc.GetByID(*ProjectActionPlanFormDTO.SeverityLevelID); err != nil {
			return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
		}
	}

	if ProjectActionPlanFormDTO.ReviewUserID != nil {
		if ProjectActionPlanFormDTO.ResponsibleUserID == *ProjectActionPlanFormDTO.ReviewUserID {
			return utility.NewValidationError(fmt.Errorf("review user and responsible user must be distinct"), "distinct_review_user", "project action plan distinct review user")
		}
	}

	return nil
}

func (paps *ProjectActionPlanService) parseProjectActionPlanDates(dto *projectDTOs.ProjectActionPlanFormDTO) (*projectDTOs.ProjectActionPlanParsedDates, error) {
	parsed := &projectDTOs.ProjectActionPlanParsedDates{}

	var err error
	if dto.TargetDate != nil {
		parsed.TargetDate, err = utility.ParseDate(*dto.TargetDate)
		if err != nil {
			return nil, utility.NewValidationError(err, "invalid_date", "target date")
		}
	}

	return parsed, nil
}

func (paps *ProjectActionPlanService) StatusUpdate(projectID uint, id uint, currentUserID uint, projectActionPlanStatusUpdateFormDTO *projectDTOs.ProjectActionPlanStatusUpdateFormDTO, tx utility.DatabaseTransaction) error {
	existingProjectActionPlan, err := paps.GetByID(projectID, id)
	if err != nil {
		return err
	}

	existingProjectActionPlan.ActionPlanStatusID = projectActionPlanStatusUpdateFormDTO.ActionPlanStatusID
	existingProjectActionPlan.Comments = &projectActionPlanStatusUpdateFormDTO.Comment
	existingProjectActionPlan.Percentage = projectActionPlanStatusUpdateFormDTO.Percentage
	existingProjectActionPlan.UpdatedBYID = &currentUserID

	err = paps.repo.Update(existingProjectActionPlan, tx)
	if err != nil {
		return err
	}

	return nil
}
