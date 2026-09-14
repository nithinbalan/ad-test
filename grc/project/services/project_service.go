package services

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	incidentServices "jethur/grc/incident/services"
	logServices "jethur/grc/logs/services"
	masterDTOs "jethur/grc/masters/dtos"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationServices "jethur/grc/organization/services"
	"jethur/grc/project/dtos"
	"jethur/grc/project/repositories"
	"jethur/utility"

	"golang.org/x/net/html"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ProjectServiceInterface interface {
	GetAll(params *dtos.ProjectListQueryParamsDTO) (*dtos.ProjectListResponseDTO, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.Project, error)
	GetProjectByID(id uint, tx ...utility.DatabaseTransaction) (*models.Project, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*dtos.ProjectDetailDTO, error)
	Create(currentUserID uint, form *dtos.ProjectFormDTO) (uint, error)
	Update(id uint, currentUserID uint, form *dtos.ProjectFormDTO) error
	UpdateProject(project *models.Project, tx ...utility.DatabaseTransaction) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	UpdateStatus(id uint, currentUserID uint, form *dtos.ProjectStatusUpdateDTO) error
	ExportProjects(data []*dtos.ProjectListDTO, currentUserID uint) (*os.File, error)
}

type ProjectService struct {
	repo                  repositories.ProjectRepositoryInterface
	activityLog           logServices.ActivityLogServiceInterface
	projectTypeSvc        ProjectTypeServiceInterface
	projectStatusSvc      ProjectStatusServiceInterface
	userSvc               organizationServices.UserServiceInterface
	userSettingSvc        organizationServices.UserSettingServiceInterface
	notificationSvc       commonServices.NotificationServiceInterface
	incidentActionPlanSvc incidentServices.IncidentActionPlanServiceInterface
}

func NewProjectService(
	repo repositories.ProjectRepositoryInterface,
	activityLog logServices.ActivityLogServiceInterface,
	projectTypeSvc ProjectTypeServiceInterface,
	projectStatusSvc ProjectStatusServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	incidentActionPlanSvc incidentServices.IncidentActionPlanServiceInterface,
) ProjectServiceInterface {
	return &ProjectService{
		repo:                  repo,
		activityLog:           activityLog,
		projectTypeSvc:        projectTypeSvc,
		projectStatusSvc:      projectStatusSvc,
		userSvc:               userSvc,
		userSettingSvc:        userSettingSvc,
		notificationSvc:       notificationSvc,
		incidentActionPlanSvc: incidentActionPlanSvc,
	}
}

func (s *ProjectService) GetAll(params *dtos.ProjectListQueryParamsDTO) (*dtos.ProjectListResponseDTO, error) {
	return s.repo.GetAll(params)
}

func (s *ProjectService) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.Project, error) {
	project, err := s.repo.GetByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return project, nil
}

func (s *ProjectService) GetProjectByID(id uint, tx ...utility.DatabaseTransaction) (*models.Project, error) {
	return s.GetByID(id, tx...)
}

func (s *ProjectService) UpdateProject(project *models.Project, tx ...utility.DatabaseTransaction) error {
	return s.repo.Update(project, tx...)
}

func (s *ProjectService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*dtos.ProjectDetailDTO, error) {
	var projectDTO dtos.ProjectDetailDTO

	project, err := s.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &utility.ResourceNotFoundError{}
		}
		return nil, err
	}

	err = mapstructure.Decode(project, &projectDTO)
	if err != nil {
		return nil, err
	}

	// Format date fields
	if project.StartDate != (time.Time{}) {
		str := project.StartDate.Format("2006-01-02")
		projectDTO.StartDateString = str
	}
	if project.EndDate != (time.Time{}) {
		str := project.EndDate.Format("2006-01-02")
		projectDTO.EndDateString = str
	}

	projectDTO.CreatedAt = project.CreatedAt
	projectDTO.UpdatedAt = project.UpdatedAt

	projectDTO.IncidentActionPlanID = project.IncidentActionPlanID
	if project.IncidentActionPlan != nil {
		actionPlanDTO := &dtos.ProjectIncidentActionPlanDTO{
			ID:                 project.IncidentActionPlan.ID,
			IncidentRegisterID: project.IncidentActionPlan.IncidentRegisterID,
			Title:              project.IncidentActionPlan.Title,
			Description:        project.IncidentActionPlan.Description,
			Percentage:         project.IncidentActionPlan.Percentage,
			ActionPlanStatusID: project.IncidentActionPlan.ActionPlanStatusID,
			ResponsibleUserID:  project.IncidentActionPlan.ResponsibleUserID,
		}
		if project.IncidentActionPlan.IncidentRegister != nil {
			actionPlanDTO.IncidentReferenceNumber = &project.IncidentActionPlan.IncidentRegister.ReferenceNumber
			actionPlanDTO.IncidentTitle = &project.IncidentActionPlan.IncidentRegister.Title
		}
		if project.IncidentActionPlan.TargetDate != nil {
			targetDateStr := project.IncidentActionPlan.TargetDate.Format("2006-01-02")
			actionPlanDTO.TargetDate = &targetDateStr
		}
		if project.IncidentActionPlan.ActionPlanStatus != nil {
			actionPlanDTO.ActionPlanStatusTitleEn = &project.IncidentActionPlan.ActionPlanStatus.TitleEn
			actionPlanDTO.ActionPlanStatusTitleAr = &project.IncidentActionPlan.ActionPlanStatus.TitleAr
			actionPlanDTO.ActionPlanStatusColor = &project.IncidentActionPlan.ActionPlanStatus.ColorCode
		}
		if project.IncidentActionPlan.ResponsibleUser != nil {
			actionPlanDTO.ResponsibleUser = &organizationDTOs.UserMinimalDTO{
				ID:               int(project.IncidentActionPlan.ResponsibleUser.ID),
				FirstNameEN:      project.IncidentActionPlan.ResponsibleUser.FirstNameEn,
				LastNameEN:       project.IncidentActionPlan.ResponsibleUser.LastNameEn,
				FirstNameAR:      project.IncidentActionPlan.ResponsibleUser.FirstNameAr,
				LastNameAR:       project.IncidentActionPlan.ResponsibleUser.LastNameAr,
				ProfileImagePath: project.IncidentActionPlan.ResponsibleUser.ProfileImagePath,
				Email:            project.IncidentActionPlan.ResponsibleUser.Email,
				Initial:          project.IncidentActionPlan.ResponsibleUser.Initial,
				Color:            project.IncidentActionPlan.ResponsibleUser.Color,
			}
		}
		projectDTO.IncidentActionPlan = actionPlanDTO
	}

	if project.WorkflowStatus != nil {
		projectDTO.WorkflowStatus = &masterDTOs.WorkflowStatusDTO{
			ID:        project.WorkflowStatus.ID,
			TitleEn:   project.WorkflowStatus.TitleEn,
			TitleAr:   project.WorkflowStatus.TitleAr,
			ColorCode: project.WorkflowStatus.ColorCode,
			Type:      project.WorkflowStatus.Type,
			Status:    project.WorkflowStatus.Status,
		}
	}
	if project.SubmittedByUser != nil {
		projectDTO.SubmittedByUser = &organizationDTOs.UserMinimalDTO{
			ID:               int(project.SubmittedByUser.ID),
			FirstNameEN:      project.SubmittedByUser.FirstNameEn,
			LastNameEN:       project.SubmittedByUser.LastNameEn,
			FirstNameAR:      project.SubmittedByUser.FirstNameAr,
			LastNameAR:       project.SubmittedByUser.LastNameAr,
			ProfileImagePath: project.SubmittedByUser.ProfileImagePath,
			Email:            project.SubmittedByUser.Email,
			Initial:          project.SubmittedByUser.Initial,
			Color:            project.SubmittedByUser.Color,
		}
	}
	if project.ApprovedByUser != nil {
		projectDTO.ApprovedByUser = &organizationDTOs.UserMinimalDTO{
			ID:               int(project.ApprovedByUser.ID),
			FirstNameEN:      project.ApprovedByUser.FirstNameEn,
			LastNameEN:       project.ApprovedByUser.LastNameEn,
			FirstNameAR:      project.ApprovedByUser.FirstNameAr,
			LastNameAR:       project.ApprovedByUser.LastNameAr,
			ProfileImagePath: project.ApprovedByUser.ProfileImagePath,
			Email:            project.ApprovedByUser.Email,
			Initial:          project.ApprovedByUser.Initial,
			Color:            project.ApprovedByUser.Color,
		}
	}
	if project.RejectedByUser != nil {
		projectDTO.RejectedByUser = &organizationDTOs.UserMinimalDTO{
			ID:               int(project.RejectedByUser.ID),
			FirstNameEN:      project.RejectedByUser.FirstNameEn,
			LastNameEN:       project.RejectedByUser.LastNameEn,
			FirstNameAR:      project.RejectedByUser.FirstNameAr,
			LastNameAR:       project.RejectedByUser.LastNameAr,
			ProfileImagePath: project.RejectedByUser.ProfileImagePath,
			Email:            project.RejectedByUser.Email,
			Initial:          project.RejectedByUser.Initial,
			Color:            project.RejectedByUser.Color,
		}
	}

	return &projectDTO, nil
}

func (s *ProjectService) Create(currentUserID uint, form *dtos.ProjectFormDTO) (uint, error) {
	if form.IncidentActionPlanID != nil && *form.IncidentActionPlanID != 0 {
		_, err := s.incidentActionPlanSvc.GetActionPlanByID(*form.IncidentActionPlanID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return 0, utility.NewValidationError(fmt.Errorf("incident action plan does not exist"), "not_exist", "incident action plan")
			}
			return 0, err
		}

		exists, err := s.repo.ExistsByIncidentActionPlanID(*form.IncidentActionPlanID)
		if err != nil {
			return 0, err
		}
		if exists {
			return 0, utility.NewValidationError(fmt.Errorf("incident action plan has already been used"), "already_used", "incident action plan")
		}
	}

	startDate, endDate, err := parseProjectDates(form)
	if err != nil {
		return 0, err
	}

	project := &models.Project{
		TitleEn:                 form.TitleEn,
		TitleAr:                 form.TitleAr,
		DescriptionEn:           form.DescriptionEn,
		DescriptionAr:           form.DescriptionAr,
		Budget:                  form.Budget,
		Scope:                   form.Scope,
		NumberOfResources:       form.NumberOfResources,
		ResponsibleUserID:       form.ResponsibleUserID,
		StartDate:               *startDate,
		EndDate:                 *endDate,
		ProjectTypeID:           form.ProjectTypeID,
		ProjectStatusID:         form.ProjectStatusID,
		ProjectPhaseID:          form.ProjectPhaseID,
		ProjectClassificationID: &form.ProjectClassificationID,
		IncidentActionPlanID:    form.IncidentActionPlanID,
		CreatedBYID:             currentUserID,
	}

	id, err := s.repo.Create(project)
	if err != nil {
		return 0, err
	}

	go func() {
		activityData := map[string]interface{}{
			"title_en":                  form.TitleEn,
			"title_ar":                  form.TitleAr,
			"description_en":            form.DescriptionEn,
			"description_ar":            form.DescriptionAr,
			"budget":                    form.Budget,
			"scope":                     form.Scope,
			"number_of_resources":       form.NumberOfResources,
			"responsible_user_id":       form.ResponsibleUserID,
			"start_date":                form.StartDate,
			"end_date":                  form.EndDate,
			"project_type_id":           form.ProjectTypeID,
			"project_status_id":         form.ProjectStatusID,
			"project_phase_id":          form.ProjectPhaseID,
			"project_classification_id": form.ProjectClassificationID,
			"incident_action_plan_id":   form.IncidentActionPlanID,
		}
		err := s.activityLog.LogActivity(models.CreateAction, models.Project{}.TableName(), id, activityData, currentUserID, "project", "project_register", "", nil)
		if err != nil {
			utility.LogError("Failed to log create activity", map[string]interface{}{
				"recordID": id,
				"error":    err.Error(),
			})
		}
		// Prepare notification and email
		urlParam := fmt.Sprintf(`{"id":%d}`, id)
		entity := "project"
		action := "responsible-user-assign"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{form.ResponsibleUserID},
			ModuleIdentifier:     "project",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "Project Responsibility Assigned",
			HeadingAr:            "تم تعيين مسؤولية المشروع",
			TitleEn:              form.TitleEn,
			TitleAr:              form.TitleAr,
			ReferenceNumber:      "",
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		responsibleUser, err := s.userSvc.GetByID(form.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch responsible user for project assignment notification", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		createdBy, err := s.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch creator user for project assignment notification", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("Project/%s: Responsibility Assigned", form.TitleEn)
		appURL := utility.GetAppURL()

		emailData := map[string]interface{}{
			"FirstNameEn":          responsibleUser.FirstNameEn,
			"LastNameEn":           responsibleUser.LastNameEn,
			"ProjectTitleEn":       form.TitleEn,
			"ProjectTitleAr":       form.TitleAr,
			"ProjectDescription":   stripHTML(form.DescriptionEn),
			"CreatedByFirstNameEn": createdBy.FirstNameEn,
			"CreatedByLastNameEn":  createdBy.LastNameEn,
			"AssignedOn":           time.Now().Format("2006-01-02"),
			"RedirectURL":          fmt.Sprintf("%s/project/%d/overview", appURL, id),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "project/templates/email/project_create_responsible_user_mail.html",
			EmailData:            emailData,
		}

		err = s.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "project_register")
		if err != nil {
			utility.LogError("Failed to send project assignment notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	return id, nil
}

func (s *ProjectService) Update(id uint, currentUserID uint, form *dtos.ProjectFormDTO) error {
	project, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &utility.ResourceNotFoundError{}
		}
		return err
	}

	oldResponsibleUserID := project.ResponsibleUserID

	oldData := map[string]interface{}{
		"title_en":                  project.TitleEn,
		"title_ar":                  project.TitleAr,
		"description_en":            project.DescriptionEn,
		"description_ar":            project.DescriptionAr,
		"budget":                    project.Budget,
		"scope":                     project.Scope,
		"number_of_resources":       project.NumberOfResources,
		"responsible_user_id":       project.ResponsibleUserID,
		"start_date":                project.StartDate,
		"end_date":                  project.EndDate,
		"project_type_id":           project.ProjectTypeID,
		"project_status_id":         project.ProjectStatusID,
		"project_phase_id":          project.ProjectPhaseID,
		"project_classification_id": project.ProjectClassificationID,
		"incident_action_plan_id":   project.IncidentActionPlanID,
	}

	if form.IncidentActionPlanID != nil {
		if *form.IncidentActionPlanID != 0 {
			_, err := s.incidentActionPlanSvc.GetActionPlanByID(*form.IncidentActionPlanID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return utility.NewValidationError(fmt.Errorf("incident action plan does not exist"), "not_exist", "incident action plan")
				}
				return err
			}

			exists, err := s.repo.ExistsByIncidentActionPlanIDExcludingID(id, *form.IncidentActionPlanID)
			if err != nil {
				return err
			}
			if exists {
				return utility.NewValidationError(fmt.Errorf("incident action plan has already been used"), "already_used", "incident action plan")
			}
			project.IncidentActionPlanID = form.IncidentActionPlanID
		} else {
			project.IncidentActionPlanID = nil
		}
	}

	startDate, endDate, err := parseProjectDates(form)
	if err != nil {
		return err
	}

	project.TitleEn = form.TitleEn
	project.TitleAr = form.TitleAr
	project.DescriptionEn = form.DescriptionEn
	project.DescriptionAr = form.DescriptionAr
	project.Budget = form.Budget
	project.Scope = form.Scope
	project.NumberOfResources = form.NumberOfResources
	project.ResponsibleUserID = form.ResponsibleUserID
	project.StartDate = *startDate
	project.EndDate = *endDate
	project.ProjectTypeID = form.ProjectTypeID
	project.ProjectStatusID = form.ProjectStatusID
	project.ProjectPhaseID = form.ProjectPhaseID
	project.ProjectClassificationID = &form.ProjectClassificationID
	project.UpdatedBYID = &currentUserID

	err = s.repo.Update(project)
	if err != nil {
		return err
	}

	go func() {
		newData := map[string]interface{}{
			"title_en":                  form.TitleEn,
			"title_ar":                  form.TitleAr,
			"description_en":            form.DescriptionEn,
			"description_ar":            form.DescriptionAr,
			"budget":                    form.Budget,
			"scope":                     form.Scope,
			"number_of_resources":       form.NumberOfResources,
			"responsible_user_id":       form.ResponsibleUserID,
			"start_date":                form.StartDate,
			"end_date":                  form.EndDate,
			"project_type_id":           form.ProjectTypeID,
			"project_status_id":         form.ProjectStatusID,
			"project_phase_id":          form.ProjectPhaseID,
			"project_classification_id": form.ProjectClassificationID,
			"incident_action_plan_id":   project.IncidentActionPlanID,
		}
		changes := utility.CompareFieldChanges(oldData, newData)
		if len(changes) > 0 {
			err := s.activityLog.LogActivity(models.UpdateAction, models.Project{}.TableName(), id, changes, currentUserID, "project", "project_register", "", nil)
			if err != nil {
				utility.LogError("Failed to log update activity", map[string]interface{}{
					"recordID": id,
					"error":    err.Error(),
				})
			}
		}

		// Send notifications to old and new responsible users if the responsible user has changed
		if oldResponsibleUserID != form.ResponsibleUserID {
			// URL param for notifications
			urlParam := fmt.Sprintf(`{"id":%d}`, id)
			entity := "project"

			// Old responsible user notification (Removal)
			oldResponsibleUserAction := "responsible-user-remove"
			oldResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{oldResponsibleUserID},
				ModuleIdentifier:     "project",
				NotificationType:     string(commonDTOs.InformationType),
				HeadingEn:            "Project Responsibility Removed",
				HeadingAr:            "تمت إزالة مسؤولية المشروع",
				TitleEn:              project.TitleEn,
				TitleAr:              project.TitleAr,
				ReferenceNumber:      "",
				Entity:               &entity,
				Action:               &oldResponsibleUserAction,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			// New responsible user notification (Assignment)
			newResponsibleUserAction := "responsible-user-assign"
			newResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{form.ResponsibleUserID},
				ModuleIdentifier:     "project",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Project Responsibility Assigned",
				HeadingAr:            "تم تعيين مسؤولية المشروع",
				TitleEn:              project.TitleEn,
				TitleAr:              project.TitleAr,
				ReferenceNumber:      "",
				Entity:               &entity,
				Action:               &newResponsibleUserAction,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			// Fetch the responsible users' data
			oldResponsibleUser, err := s.userSvc.GetByID(oldResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch old responsible user for project assignment change", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			newResponsibleUser, err := s.userSvc.GetByID(form.ResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch new responsible user for project assignment change", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			// Fetch details of the user making the assignment
			updatedByUser, err := s.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch assigning user for project assignment change", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			appURL := utility.GetAppURL()

			// Email notification for the old responsible user
			oldResponsibleUserEmailTo := []string{oldResponsibleUser.Email}
			oldResponsibleUserSubject := fmt.Sprintf("Project/%s: Responsibility Removed", project.TitleEn)
			oldResponsibleUserEmailData := map[string]interface{}{
				"FirstNameEn":          oldResponsibleUser.FirstNameEn,
				"LastNameEn":           oldResponsibleUser.LastNameEn,
				"ProjectTitle":         project.TitleEn,
				"ProjectTitleAr":       project.TitleAr,
				"ProjectDescription":   stripHTML(project.DescriptionEn),
				"UpdatedByFirstNameEn": updatedByUser.FirstNameEn,
				"UpdatedByLastNameEn":  updatedByUser.LastNameEn,
				"RedirectURL":          fmt.Sprintf("%s/project/%d/overview", appURL, id),
			}

			oldResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              oldResponsibleUserEmailTo,
				Subject:              oldResponsibleUserSubject,
				TemplateRelativePath: "project/templates/email/project_remove_responsible_user_mail.html",
				EmailData:            oldResponsibleUserEmailData,
			}

			// Email notification for the new responsible user
			newResponsibleUserEmailTo := []string{newResponsibleUser.Email}
			newResponsibleUserSubject := fmt.Sprintf("Project/%s: Responsibility Assigned", form.TitleEn)
			newResponsibleUserEmailData := map[string]interface{}{
				"FirstNameEn":           newResponsibleUser.FirstNameEn,
				"LastNameEn":            newResponsibleUser.LastNameEn,
				"ProjectTitle":          form.TitleEn,
				"ProjectTitleAr":        form.TitleAr,
				"ProjectDescription":    stripHTML(form.DescriptionEn),
				"AssignedByFirstNameEn": updatedByUser.FirstNameEn,
				"AssignedByLastNameEn":  updatedByUser.LastNameEn,
				"AssignedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           fmt.Sprintf("%s/project/%d/overview", appURL, id),
			}

			newResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              newResponsibleUserEmailTo,
				Subject:              newResponsibleUserSubject,
				TemplateRelativePath: "project/templates/email/project_change_responsible_user_mail.html",
				EmailData:            newResponsibleUserEmailData,
			}

			// Send notifications and emails
			err = s.notificationSvc.SendNotification(newResponsibleUserNotificationDTO, newResponsibleUserEmailNotificationDTO, "project_register")
			if err != nil {
				utility.LogError("Failed to send notifications to new responsible user for project", map[string]interface{}{
					"error": err.Error(),
				})
			}

			err = s.notificationSvc.SendNotification(oldResponsibleUserNotificationDTO, oldResponsibleUserEmailNotificationDTO, "project_register")
			if err != nil {
				utility.LogError("Failed to send notifications to old responsible user for project", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()

	return nil
}

func (s *ProjectService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedProjects []map[string]interface{}

	for _, id := range ids {
		project, err := s.repo.GetByID(id)
		if err != nil {
			fmt.Print(err)
			failedIds = append(failedIds, id)
			continue
		}
		err = s.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedProjects = append(deletedProjects, map[string]interface{}{
			"id":       id,
			"title_en": project.TitleEn,
			"title_ar": project.TitleAr,
		})
	}
	if len(successfulIds) > 0 {
		go func() {
			err := s.activityLog.LogBulkActivities(
				models.DeleteAction,
				models.Project{}.TableName(),
				successfulIds,
				deletedProjects,
				currentUserID,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to log bulk delete activity", map[string]interface{}{
					"recordIDs": successfulIds,
					"error":     err.Error(),
				})
			}
		}()
	}

	return failedIds, nil
}

func parseProjectDates(form *dtos.ProjectFormDTO) (startDate, endDate *time.Time, err error) {
	if form.StartDate != "" {
		sd, err := utility.ParseDate(form.StartDate)
		if err != nil {
			return nil, nil, utility.NewValidationError(fmt.Errorf("invalid start date"), "invalid_date", "start_date")
		}
		startDate = sd
	}
	if form.EndDate != "" {
		ed, err := utility.ParseDate(form.EndDate)
		if err != nil {
			return nil, nil, utility.NewValidationError(fmt.Errorf("invalid end date"), "invalid_date", "end_date")
		}
		endDate = ed
	}
	return startDate, endDate, nil
}

func (s *ProjectService) UpdateStatus(id uint, currentUserID uint, projectStatusUpdateDTO *dtos.ProjectStatusUpdateDTO) error {
	tx := s.repo.BeginTransaction()

	project, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &utility.ResourceNotFoundError{}
		}
		return err
	}

	oldData := map[string]interface{}{
		"project_status_id": project.ProjectStatusID,
	}

	// Update the weightage
	project.ProjectStatusID = projectStatusUpdateDTO.ProjectStatusID
	project.UpdatedBYID = &currentUserID

	// Save the changes
	if err := s.repo.Update(project, tx); err != nil {
		tx.Rollback()
		return err
	}

	newData := map[string]interface{}{
		"project_status_id": project.ProjectStatusID,
	}

	// Commit first before logging
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	// Log the activity asynchronously
	go func() {
		err := s.activityLog.LogActivity(
			models.UpdateAction,
			models.Project{}.TableName(),
			id,
			map[string]interface{}{
				"old": oldData,
				"new": newData,
			},
			currentUserID,
			"project", "project_register", "", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.Project{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (s *ProjectService) ExportProjects(data []*dtos.ProjectListDTO, currentUserID uint) (*os.File, error) {
	userSetting, err := s.userSettingSvc.GetUserSetting(currentUserID)
	if err != nil {
		return nil, err
	}

	headers := []string{
		"Title",
		"Description",
		"Project Type",
		"Project Phase",
		"Classification",
		"Budget",
		"Scope",
		"Start Date",
		"End Date",
		"Responsible User",
		"Project Status",
	}

	rows := make([][]interface{}, len(data))
	for i, project := range data {
		if userSetting.Language == "ar" {
			rows[i] = []interface{}{
				project.TitleAr,
				stripHTML(project.DescriptionAr),
				project.ProjectTypeTitleAr,
				project.ProjectPhaseTitleAr,
				project.ProjectClassificationTitleAr,
				project.Budget,
				stripHTML(project.Scope),
				project.StartDate,
				project.EndDate,
				fmt.Sprintf("%s %s", project.ResponsibleUserFirstNameAr, project.ResponsibleUserLastNameAr),
				project.ProjectStatusTitleAr,
			}
		} else {
			rows[i] = []interface{}{
				project.TitleEn,
				stripHTML(project.DescriptionEn),
				project.ProjectTypeTitleEn,
				project.ProjectPhaseTitleEn,
				project.ProjectClassificationTitleEn,
				project.Budget,
				stripHTML(project.Scope),
				project.StartDate,
				project.EndDate,
				fmt.Sprintf("%s %s", project.ResponsibleUserFirstNameEn, project.ResponsibleUserLastNameEn),
				project.ProjectStatusTitleEn,
			}
		}
	}

	return utility.GenerateExcelFile(headers, rows)
}

func stripHTML(content *string) string {
	if content == nil {
		return ""
	}

	document, err := html.Parse(strings.NewReader(*content))
	if err != nil {
		return *content
	}

	var text strings.Builder
	var collectText func(*html.Node)
	collectText = func(node *html.Node) {
		if node.Type == html.TextNode {
			text.WriteString(node.Data)
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			collectText(child)
		}
	}

	collectText(document)
	return strings.TrimSpace(text.String())
}
