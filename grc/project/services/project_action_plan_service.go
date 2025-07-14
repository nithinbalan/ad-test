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
	GetAll(id uint, params *projectDTOs.ProjectActionPlanListQueryParamsDTO) (*projectDTOs.ProjectActionPlanListResponseDTO, error)
	GetByID(id uint) (*models.ProjectActionPlan, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*projectDTOs.ProjectActionPlanDetailDTO, error)
	Create(currentUserID uint, ProjectActionPlanFormDTO *projectDTOs.ProjectActionPlanFormDTO) (uint, error)
	Update(id uint, currentUserID uint, project *projectDTOs.ProjectActionPlanFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type ProjectActionPlanService struct {
	repo                projectRepositories.ProjectActionPlanRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface
	userSvc             organizationServices.UserServiceInterface
	projectMilestoneSvc ProjectMilestoneServiceInterface
}

func NewProjectActionPlanService(
	repo projectRepositories.ProjectActionPlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	projectMilestoneSvc ProjectMilestoneServiceInterface,
) ProjectActionPlanServiceInterface {
	return &ProjectActionPlanService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		actionPlanStatusSvc: actionPlanStatusSvc,
		userSvc:             userSvc,
		projectMilestoneSvc: projectMilestoneSvc,
	}
}

func (paps *ProjectActionPlanService) GetAll(id uint, params *projectDTOs.ProjectActionPlanListQueryParamsDTO) (*projectDTOs.ProjectActionPlanListResponseDTO, error) {
	if id != 0 {
		if _, err := paps.projectMilestoneSvc.GetByID(id); err != nil {
			return nil, err
		}
	}
	return paps.repo.GetAll(id, params)
}

func (paps ProjectActionPlanService) GetByID(id uint) (*models.ProjectActionPlan, error) {
	if _, err := paps.projectMilestoneSvc.GetByID(id); err != nil {
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

func (paps ProjectActionPlanService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*projectDTOs.ProjectActionPlanDetailDTO, error) {
	var projectActionPlanDTO projectDTOs.ProjectActionPlanDetailDTO

	// Check if the compliance register exists
	if _, err := paps.projectMilestoneSvc.GetByID(id); err != nil {
		return nil, err
	}

	// Get full details of the compliance action plan
	projectActionPlan, err := paps.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Decode the compliance action plan into the DTO
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

func (paps *ProjectActionPlanService) Create(currentUserID uint, ProjectActionPlanFormDTO *projectDTOs.ProjectActionPlanFormDTO) (uint, error) {
	tx := paps.repo.BeginTransaction()

	if err := paps.validateProjectActionPlanFormDTO(ProjectActionPlanFormDTO.ProjectMilestoneID, ProjectActionPlanFormDTO); err != nil {
		return 0, err
	}

	dates, err := paps.parseProjectActionPlanDates(ProjectActionPlanFormDTO)
	if err != nil {
		return 0, err
	}

	actionPlanStatus, err := paps.actionPlanStatusSvc.GetByID(ProjectActionPlanFormDTO.ActionPlanStatusID)
	if err != nil {
		return 0, err
	}

	projectActionPlan := &models.ProjectActionPlan{
		ProjectMilestoneID: &ProjectActionPlanFormDTO.ProjectMilestoneID,
		Title:              ProjectActionPlanFormDTO.Title,
		Description:        ProjectActionPlanFormDTO.Description,
		ResponsibleUserID:  ProjectActionPlanFormDTO.ResponsibleUserID,
		Comments:           ProjectActionPlanFormDTO.Comments,
		TargetDate:         dates.TargetDate,
		Percentage:         &ProjectActionPlanFormDTO.Percentage,
		ActionPlanStatusID: actionPlanStatus.ID,
		CreatedByID:        currentUserID,
	}

	id, err := paps.repo.Create(projectActionPlan, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	go func() {
		activityData := map[string]interface{}{
			"title":                ProjectActionPlanFormDTO.Title,
			"description":          ProjectActionPlanFormDTO.Description,
			"responsible_user_id":  ProjectActionPlanFormDTO.ResponsibleUserID,
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

func (paps *ProjectActionPlanService) Update(id uint, currentUserID uint, ProjectActionPlanFormDTO *projectDTOs.ProjectActionPlanFormDTO) error {
	tx := paps.repo.BeginTransaction()

	// Validate project milestone exists
	if _, err := paps.projectMilestoneSvc.GetByID(ProjectActionPlanFormDTO.ProjectMilestoneID); err != nil {
		return err
	}

	existingProjectActionPlan, err := paps.GetByID(id)
	if err != nil {
		return err
	}

	if err := paps.validateProjectActionPlanFormDTO(ProjectActionPlanFormDTO.ProjectMilestoneID, ProjectActionPlanFormDTO); err != nil {
		return err
	}

	dates, err := paps.parseProjectActionPlanDates(ProjectActionPlanFormDTO)
	if err != nil {
		return err
	}

	// Update Project ActionPlan
	existingProjectActionPlan.Title = ProjectActionPlanFormDTO.Title
	existingProjectActionPlan.Description = ProjectActionPlanFormDTO.Description
	existingProjectActionPlan.Comments = ProjectActionPlanFormDTO.Comments
	existingProjectActionPlan.ResponsibleUserID = ProjectActionPlanFormDTO.ResponsibleUserID
	existingProjectActionPlan.TargetDate = dates.TargetDate
	existingProjectActionPlan.Percentage = &ProjectActionPlanFormDTO.Percentage
	existingProjectActionPlan.ActionPlanStatusID = ProjectActionPlanFormDTO.ActionPlanStatusID
	existingProjectActionPlan.UpdatedByID = &currentUserID

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
			"target_date":           dates.TargetDate,
			"percentage":            ProjectActionPlanFormDTO.Percentage,
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

func (paps *ProjectActionPlanService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
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

func (paps *ProjectActionPlanService) validateProjectActionPlanFormDTO(projectMilestoneID uint, ProjectActionPlanFormDTO *projectDTOs.ProjectActionPlanFormDTO) error {
	// Validate responsible user exists
	if _, err := paps.userSvc.GetByID(ProjectActionPlanFormDTO.ResponsibleUserID); err != nil {
		return utility.NewValidationError(fmt.Errorf("responsible user does not exist"), "not_exist", "responsible user")
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
