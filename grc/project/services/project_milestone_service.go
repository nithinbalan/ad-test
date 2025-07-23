package services

import (
	"errors"
	"fmt"
	"time"

	masterServices "jethur/grc/masters/services"

	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/grc/project/dtos"
	"jethur/grc/project/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ProjectMilestoneServiceInterface interface {
	GetAll(projectID uint, params *dtos.ProjectMilestoneListQueryParamsDTO) (*dtos.ProjectMilestoneListResponseDTO, error)
	GetByID(projectID uint, id uint) (*models.ProjectMilestone, error)
	GetFullDetailsByID(projectID uint, id uint) (*dtos.ProjectMilestoneDetailDTO, error)
	Create(projectID uint, currentUserID uint, form *dtos.ProjectMilestoneFormDTO) (uint, error)
	Update(projectID uint, id uint, currentUserID uint, form *dtos.ProjectMilestoneFormDTO) error
	Delete(projectID uint, ids []uint, currentUserID uint) ([]uint, error)
}

type ProjectMilestoneService struct {
	repo               repositories.ProjectMilestoneRepositoryInterface
	activityLog        logServices.ActivityLogServiceInterface
	userSvc            organizationServices.UserServiceInterface
	projectSvc         ProjectServiceInterface
	milestoneStatusSvc masterServices.ProjectMilestoneStatusServiceInterface
}

func NewProjectMilestoneService(
	repo repositories.ProjectMilestoneRepositoryInterface,
	activityLog logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	projectSvc ProjectServiceInterface,
	milestoneStatusSvc masterServices.ProjectMilestoneStatusServiceInterface,
) ProjectMilestoneServiceInterface {
	return &ProjectMilestoneService{
		repo:               repo,
		activityLog:        activityLog,
		userSvc:            userSvc,
		projectSvc:         projectSvc,
		milestoneStatusSvc: milestoneStatusSvc,
	}
}

func (pms *ProjectMilestoneService) GetAll(projectID uint, params *dtos.ProjectMilestoneListQueryParamsDTO) (*dtos.ProjectMilestoneListResponseDTO, error) {
	return pms.repo.GetAll(projectID, params)
}

func (pms *ProjectMilestoneService) GetByID(projectID uint, id uint) (*models.ProjectMilestone, error) {
	if _, err := pms.projectSvc.GetByID(projectID); err != nil {
		return nil, err
	}
	milestone, err := pms.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return milestone, nil
}

func (pms *ProjectMilestoneService) GetFullDetailsByID(projectID uint, id uint) (*dtos.ProjectMilestoneDetailDTO, error) {
	var milestoneDTO dtos.ProjectMilestoneDetailDTO

	// Check if the compliance register exists
	if _, err := pms.projectSvc.GetByID(projectID); err != nil {
		return nil, err
	}

	milestone, err := pms.repo.GetFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &utility.ResourceNotFoundError{}
		}
		return nil, err
	}

	err = mapstructure.Decode(milestone, &milestoneDTO)
	if err != nil {
		return nil, err
	}

	// Format date fields for nullable *time.Time
	if milestone.PlannedDate != nil {
		str := milestone.PlannedDate.Format("2006-01-02")
		milestoneDTO.PlannedDateString = &str
	}
	if milestone.ActualDate != nil {
		str := milestone.ActualDate.Format("2006-01-02")
		milestoneDTO.ActualDateString = &str
	}

	milestoneDTO.CreatedAt = milestone.CreatedAt
	milestoneDTO.UpdatedAt = milestone.UpdatedAt

	return &milestoneDTO, nil
}

func (pms *ProjectMilestoneService) Create(projectID uint, currentUserID uint, form *dtos.ProjectMilestoneFormDTO) (uint, error) {
	project, err := pms.projectSvc.GetByID(projectID)
	if err != nil {
		return 0, err
	}
	plannedDate, actualDate, err := parseMilestoneDates(form)
	if err != nil {
		return 0, err
	}

	milestone := &models.ProjectMilestone{
		ProjectID:                project.ID,
		TitleEn:                  form.TitleEn,
		TitleAr:                  form.TitleAr,
		DescriptionEn:            form.DescriptionEn,
		DescriptionAr:            form.DescriptionAr,
		PlannedDate:              plannedDate,
		ActualDate:               actualDate,
		ProjectMilestoneStatusID: form.ProjectMilestoneStatusID,
		CreatedBYID:              currentUserID,
		CreatedAt:                time.Now(),
	}

	id, err := pms.repo.Create(milestone)
	if err != nil {
		return 0, err
	}

	go func() {
		activityData := map[string]interface{}{
			"project_id":                  projectID,
			"title_en":                    form.TitleEn,
			"title_ar":                    form.TitleAr,
			"description_en":              form.DescriptionEn,
			"description_ar":              form.DescriptionAr,
			"planned_date":                form.PlannedDate,
			"actual_date":                 form.ActualDate,
			"project_milestone_status_id": form.ProjectMilestoneStatusID,
		}
		err := pms.activityLog.LogActivity(models.CreateAction, models.ProjectMilestone{}.TableName(), id, activityData, currentUserID, nil)
		if err != nil {
			utility.LogError("Failed to log create milestone activity", map[string]interface{}{
				"recordID": id,
				"error":    err.Error(),
			})
		}
	}()

	return id, nil
}

func (pms *ProjectMilestoneService) Update(projectID uint, id uint, currentUserID uint, form *dtos.ProjectMilestoneFormDTO) error {

	milestone, err := pms.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &utility.ResourceNotFoundError{}
		}
		return err
	}

	oldData := map[string]interface{}{
		"title_en":                    milestone.TitleEn,
		"title_ar":                    milestone.TitleAr,
		"description_en":              milestone.DescriptionEn,
		"description_ar":              milestone.DescriptionAr,
		"planned_date":                milestone.PlannedDate,
		"actual_date":                 milestone.ActualDate,
		"project_milestone_status_id": milestone.ProjectMilestoneStatusID,
	}

	plannedDate, actualDate, err := parseMilestoneDates(form)
	if err != nil {
		return err
	}

	milestone.TitleEn = form.TitleEn
	milestone.TitleAr = form.TitleAr
	milestone.DescriptionEn = form.DescriptionEn
	milestone.DescriptionAr = form.DescriptionAr
	milestone.PlannedDate = plannedDate
	milestone.ActualDate = actualDate
	milestone.ProjectMilestoneStatusID = form.ProjectMilestoneStatusID
	milestone.UpdatedBYID = &currentUserID
	milestone.UpdatedAt = timePtr(time.Now())

	err = pms.repo.Update(milestone)
	if err != nil {
		return err
	}

	go func() {
		newData := map[string]interface{}{
			"title_en":                    form.TitleEn,
			"title_ar":                    form.TitleAr,
			"description_en":              form.DescriptionEn,
			"description_ar":              form.DescriptionAr,
			"planned_date":                form.PlannedDate,
			"actual_date":                 form.ActualDate,
			"project_milestone_status_id": form.ProjectMilestoneStatusID,
		}
		changes := utility.CompareFieldChanges(oldData, newData)
		if len(changes) > 0 {
			err := pms.activityLog.LogActivity(models.UpdateAction, models.ProjectMilestone{}.TableName(), id, changes, currentUserID, nil)
			if err != nil {
				utility.LogError("Failed to log update milestone activity", map[string]interface{}{
					"recordID": id,
					"error":    err.Error(),
				})
			}
		}
	}()

	return nil
}

func (pms *ProjectMilestoneService) Delete(projectID uint, ids []uint, currentUserID uint) ([]uint, error) {
	if _, err := pms.projectSvc.GetByID(projectID); err != nil {
		return nil, err
	}

	var failedIds []uint
	var successfulIds []uint
	var deletedMilestones []map[string]interface{}

	for _, id := range ids {
		milestone, err := pms.repo.GetByID(id)
		if err != nil {
			fmt.Print(err)
			failedIds = append(failedIds, id)
			continue
		}

		err = pms.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedMilestones = append(deletedMilestones, map[string]interface{}{
			"id":       id,
			"title_en": milestone.TitleEn,
			"title_ar": milestone.TitleAr,
		})
	}

	if len(successfulIds) > 0 {
		go func() {
			err := pms.activityLog.LogBulkActivities(
				models.DeleteAction,
				models.ProjectMilestone{}.TableName(),
				successfulIds,
				deletedMilestones,
				currentUserID,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to log bulk delete milestone activity", map[string]interface{}{
					"recordIDs": successfulIds,
					"error":     err.Error(),
				})
			}
		}()
	}

	return failedIds, nil
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func parseMilestoneDates(form *dtos.ProjectMilestoneFormDTO) (plannedDate, actualDate *time.Time, err error) {
	if form.PlannedDate != nil {
		pd, err := utility.ParseDate(*form.PlannedDate)
		if err != nil {
			return nil, nil, utility.NewValidationError(err, "invalid_date", "planned_date")
		}
		plannedDate = pd
	}
	if form.ActualDate != nil {
		ad, err := utility.ParseDate(*form.ActualDate)
		if err != nil {
			return nil, nil, utility.NewValidationError(err, "invalid_date", "actual_date")
		}
		actualDate = ad
	}
	return plannedDate, actualDate, nil
}
