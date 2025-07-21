package services

import (
	"errors"
	"fmt"

	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	projectDTOs "jethur/grc/project/dtos"
	projectRepositories "jethur/grc/project/repositories"
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
}

func NewProjectActionPlanService(
	repo projectRepositories.ProjectActionPlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	projectSvc ProjectServiceInterface,
	projectMilestoneSvc ProjectMilestoneServiceInterface,
	severityLevelSvc masterServices.SeverityLevelServiceInterface,
) ProjectActionPlanServiceInterface {
	return &ProjectActionPlanService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		actionPlanStatusSvc: actionPlanStatusSvc,
		userSvc:             userSvc,
		projectSvc:          projectSvc,
		projectMilestoneSvc: projectMilestoneSvc,
		severityLevelSvc:    severityLevelSvc,
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

	var actionPlanStatusID *uint
	if ProjectActionPlanFormDTO.ActionPlanStatusID != nil {
		actionPlanStatus, err := paps.actionPlanStatusSvc.GetByID(*ProjectActionPlanFormDTO.ActionPlanStatusID)
		if err != nil {
			return 0, err
		}
		actionPlanStatusID = &actionPlanStatus.ID
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
		ActionPlanStatusID: actionPlanStatusID,
		CreatedBYID:        currentUserID,
	}

	id, err := paps.repo.Create(projectActionPlan, tx)
	if err != nil {
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
		err := paps.activityLogSvc.LogActivity(models.CreateAction, models.ProjectActionPlan{}.TableName(), id, activityData, currentUserID, nil)
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
	existingProjectActionPlan.Percentage = &ProjectActionPlanFormDTO.Percentage
	existingProjectActionPlan.ActionPlanStatusID = ProjectActionPlanFormDTO.ActionPlanStatusID
	existingProjectActionPlan.UpdatedBYID = &currentUserID

	err = paps.repo.Update(existingProjectActionPlan, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	go func() {
		activityData := map[string]interface{}{
			"title":                 ProjectActionPlanFormDTO.Title,
			"description":           ProjectActionPlanFormDTO.Description,
			"comments":              ProjectActionPlanFormDTO.Comments,
			"responsible_user_id":   ProjectActionPlanFormDTO.ResponsibleUserID,
			"severity_level_id":     ProjectActionPlanFormDTO.SeverityLevelID,
			"review_user_id":        ProjectActionPlanFormDTO.ReviewUserID,
			"target_date":           dates.TargetDate,
			"percentage":            ProjectActionPlanFormDTO.Percentage,
			"project_id":            projectID,
			"project_milestone_id":  ProjectActionPlanFormDTO.ProjectMilestoneID,
			"action_plan_status_id": ProjectActionPlanFormDTO.ActionPlanStatusID,
		}
		err := paps.activityLogSvc.LogActivity(models.UpdateAction, models.ProjectActionPlan{}.TableName(), id, activityData, currentUserID, nil)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.ProjectActionPlan{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
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

	existingProjectActionPlan.ActionPlanStatusID = &projectActionPlanStatusUpdateFormDTO.ActionPlanStatusID
	existingProjectActionPlan.Comments = &projectActionPlanStatusUpdateFormDTO.Comment
	existingProjectActionPlan.Percentage = projectActionPlanStatusUpdateFormDTO.Percentage
	existingProjectActionPlan.UpdatedBYID = &currentUserID

	err = paps.repo.Update(existingProjectActionPlan, tx)
	if err != nil {
		return err
	}

	return nil
}
