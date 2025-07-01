package services

import (
	"errors"
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
	GetByID(id uint) (*dtos.ProjectDetailDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*dtos.ProjectDetailDTO, error)
	Create(currentUserID uint, form *dtos.ProjectFormDTO) (uint, error)
	Update(id uint, currentUserID uint, form *dtos.ProjectFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
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
	return s.repo.GetAll(
		params.Sort,
		params.SortBy,
		params.Keyword,
		params.Status,
		params.Page,
		params.PageSize,
		params.Filters,
		params.QuickFilter,
	)
}

func (s *ProjectService) GetByID(id uint) (*dtos.ProjectDetailDTO, error) {
	var projectDTO dtos.ProjectDetailDTO

	project, err := s.repo.GetFullDetailsByID(id)
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

	// Assign date fields directly since they are already strings in the model
	projectDTO.StartDateString = project.StartDate
	projectDTO.EndDateString = project.EndDate

	projectDTO.CreatedAt = project.CreatedAt
	projectDTO.UpdatedAt = project.UpdatedAt

	return &projectDTO, nil
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

	// Assign date fields directly since they are already strings in the model
	projectDTO.StartDateString = project.StartDate
	projectDTO.EndDateString = project.EndDate

	projectDTO.CreatedAt = project.CreatedAt
	projectDTO.UpdatedAt = project.UpdatedAt

	return &projectDTO, nil
}

func (s *ProjectService) Create(currentUserID uint, form *dtos.ProjectFormDTO) (uint, error) {
	project := &models.Project{
		TitleEn:           form.TitleEn,
		TitleAr:           form.TitleAr,
		DescriptionEn:     form.DescriptionEn,
		DescriptionAr:     form.DescriptionAr,
		ResponsibleUserID: form.ResponsibleUserID,
		StartDate:         &form.StartDate,
		EndDate:           &form.EndDate,
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
		err := s.activityLog.LogActivity(models.CreateAction, models.Project{}.TableName(), id, activityData, currentUserID, nil)
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

	project.TitleEn = form.TitleEn
	project.TitleAr = form.TitleAr
	project.DescriptionEn = form.DescriptionEn
	project.DescriptionAr = form.DescriptionAr
	project.ResponsibleUserID = form.ResponsibleUserID
	project.StartDate = &form.StartDate
	project.EndDate = &form.EndDate
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
			err := s.activityLog.LogActivity(models.UpdateAction, models.Project{}.TableName(), id, changes, currentUserID, nil)
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
	deletedIDs, err := s.repo.Delete(ids)
	if err != nil {
		return nil, err
	}

	go func() {
		if len(deletedIDs) > 0 {
			err := s.activityLog.LogBulkActivities(models.DeleteAction, models.Project{}.TableName(), deletedIDs, nil, currentUserID, nil)
			if err != nil {
				utility.LogError("Failed to log bulk delete activity", map[string]interface{}{
					"recordIDs": deletedIDs,
					"error":     err.Error(),
				})
			}
		}
	}()

	return deletedIDs, nil
}
