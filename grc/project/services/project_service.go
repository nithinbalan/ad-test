package services

import (
	"errors"
	"fmt"
	"time"

	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/grc/project/dtos"
	"jethur/grc/project/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ProjectServiceInterface interface {
	GetAll(params *dtos.ProjectListQueryParamsDTO) (*dtos.ProjectListResponseDTO, error)
	GetByID(id uint) (*models.Project, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*dtos.ProjectDetailDTO, error)
	Create(currentUserID uint, form *dtos.ProjectFormDTO) (uint, error)
	Update(id uint, currentUserID uint, form *dtos.ProjectFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	UpdateStatus(id uint, currentUserID uint, form *dtos.ProjectStatusUpdateDTO) error
}

type ProjectService struct {
	repo             repositories.ProjectRepositoryInterface
	activityLog      logServices.ActivityLogServiceInterface
	projectTypeSvc   ProjectTypeServiceInterface
	projectStatusSvc ProjectStatusServiceInterface
	userSvc          organizationServices.UserServiceInterface
}

func NewProjectService(
	repo repositories.ProjectRepositoryInterface,
	activityLog logServices.ActivityLogServiceInterface,
	projectTypeSvc ProjectTypeServiceInterface,
	projectStatusSvc ProjectStatusServiceInterface,
	userSvc organizationServices.UserServiceInterface,
) ProjectServiceInterface {
	return &ProjectService{
		repo:             repo,
		activityLog:      activityLog,
		projectTypeSvc:   projectTypeSvc,
		projectStatusSvc: projectStatusSvc,
		userSvc:          userSvc,
	}
}

func (s *ProjectService) GetAll(params *dtos.ProjectListQueryParamsDTO) (*dtos.ProjectListResponseDTO, error) {
	return s.repo.GetAll(params)
}

func (s *ProjectService) GetByID(id uint) (*models.Project, error) {
	project, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return project, nil
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

	return &projectDTO, nil
}

func (s *ProjectService) Create(currentUserID uint, form *dtos.ProjectFormDTO) (uint, error) {
	startDate, endDate, err := parseProjectDates(form)
	if err != nil {
		return 0, err
	}

	project := &models.Project{
		TitleEn:           form.TitleEn,
		TitleAr:           form.TitleAr,
		DescriptionEn:     form.DescriptionEn,
		DescriptionAr:     form.DescriptionAr,
		ResponsibleUserID: form.ResponsibleUserID,
		StartDate:         *startDate,
		EndDate:           *endDate,
		ProjectTypeID:     form.ProjectTypeID,
		ProjectStatusID:   form.ProjectStatusID,
		CreatedBYID:       currentUserID,
	}

	id, err := s.repo.Create(project)
	if err != nil {
		return 0, err
	}

	go func() {
		activityData := map[string]interface{}{
			"title_en":            form.TitleEn,
			"title_ar":            form.TitleAr,
			"description_en":      form.DescriptionEn,
			"description_ar":      form.DescriptionAr,
			"responsible_user_id": form.ResponsibleUserID,
			"start_date":          form.StartDate,
			"end_date":            form.EndDate,
			"project_type_id":     form.ProjectTypeID,
			"project_status_id":   form.ProjectStatusID,
		}
		err := s.activityLog.LogActivity(models.CreateAction, models.Project{}.TableName(), id, activityData, currentUserID, "", "", "", nil)
		if err != nil {
			utility.LogError("Failed to log create activity", map[string]interface{}{
				"recordID": id,
				"error":    err.Error(),
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

	oldData := map[string]interface{}{
		"title_en":            project.TitleEn,
		"title_ar":            project.TitleAr,
		"description_en":      project.DescriptionEn,
		"description_ar":      project.DescriptionAr,
		"responsible_user_id": project.ResponsibleUserID,
		"start_date":          project.StartDate,
		"end_date":            project.EndDate,
		"project_type_id":     project.ProjectTypeID,
		"project_status_id":   project.ProjectStatusID,
	}

	startDate, endDate, err := parseProjectDates(form)
	if err != nil {
		return err
	}

	project.TitleEn = form.TitleEn
	project.TitleAr = form.TitleAr
	project.DescriptionEn = form.DescriptionEn
	project.DescriptionAr = form.DescriptionAr
	project.ResponsibleUserID = form.ResponsibleUserID
	project.StartDate = *startDate
	project.EndDate = *endDate
	project.ProjectTypeID = form.ProjectTypeID
	project.ProjectStatusID = form.ProjectStatusID
	project.UpdatedBYID = &currentUserID

	err = s.repo.Update(project)
	if err != nil {
		return err
	}

	go func() {
		newData := map[string]interface{}{
			"title_en":            form.TitleEn,
			"title_ar":            form.TitleAr,
			"description_en":      form.DescriptionEn,
			"description_ar":      form.DescriptionAr,
			"responsible_user_id": form.ResponsibleUserID,
			"start_date":          form.StartDate,
			"end_date":            form.EndDate,
			"project_type_id":     form.ProjectTypeID,
			"project_status_id":   form.ProjectStatusID,
		}
		changes := utility.CompareFieldChanges(oldData, newData)
		if len(changes) > 0 {
			err := s.activityLog.LogActivity(models.UpdateAction, models.Project{}.TableName(), id, changes, currentUserID, "", "", "", nil)
			if err != nil {
				utility.LogError("Failed to log update activity", map[string]interface{}{
					"recordID": id,
					"error":    err.Error(),
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
			"", "", "", nil,
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
