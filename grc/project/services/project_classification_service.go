package services

import (
	"errors"
	"fmt"
	"log"

	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/project/dtos"
	"jethur/grc/project/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type ProjectClassificationServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProjectClassification, error)
	GetByID(id uint) (*models.ProjectClassification, error)
	Create(currentUserID uint, projectClassificationForm *dtos.ProjectClassificationFormDTO) (*models.ProjectClassification, error)
	Update(id uint, currentUserID uint, projectClassificationForm *dtos.ProjectClassificationFormDTO) (*models.ProjectClassification, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

type ProjectClassificationService struct {
	repo           repositories.ProjectClassificationRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewProjectClassificationService(
	repo repositories.ProjectClassificationRepositoryInterface,
	activityLogSvc services.ActivityLogServiceInterface,
) ProjectClassificationServiceInterface {
	return &ProjectClassificationService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (pcs *ProjectClassificationService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProjectClassification, error) {
	return pcs.repo.GetAll(sort, sortBy, keyword, status)
}

func (pcs ProjectClassificationService) GetByID(id uint) (*models.ProjectClassification, error) {
	projectClassification, err := pcs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return projectClassification, nil
}

func (pcs *ProjectClassificationService) Create(currentUserID uint, projectClassificationForm *dtos.ProjectClassificationFormDTO) (*models.ProjectClassification, error) {
	if exists, err := pcs.repo.ExistsByTitle(projectClassificationForm.TitleEn, projectClassificationForm.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	projectClassification := &models.ProjectClassification{
		TitleEn: projectClassificationForm.TitleEn,
		TitleAr: projectClassificationForm.TitleAr,
		Status:  projectClassificationForm.Status,
	}

	createdClassification, err := pcs.repo.Create(projectClassification)
	if err != nil {
		return nil, err
	}

	go func() {
		err := pcs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ProjectClassification{}.TableName(),
			createdClassification.ID,
			createdClassification,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdClassification, nil
}

func (pcs *ProjectClassificationService) Update(id uint, currentUserID uint, projectClassificationForm *dtos.ProjectClassificationFormDTO) (*models.ProjectClassification, error) {
	existingClassification, err := pcs.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if projectClassificationForm.TitleEn != existingClassification.TitleEn || projectClassificationForm.TitleAr != existingClassification.TitleAr {
		if exists, err := pcs.repo.ExistsByTitleExcludingID(id, projectClassificationForm.TitleEn, projectClassificationForm.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	updatedClassification := &models.ProjectClassification{
		TitleEn: projectClassificationForm.TitleEn,
		TitleAr: projectClassificationForm.TitleAr,
		Status:  projectClassificationForm.Status,
	}

	result, err := pcs.repo.Update(id, updatedClassification)
	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status"}
	changes, err := utility.CompareFields(existingClassification, updatedClassification, fieldsToCompare)
	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := pcs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.ProjectClassification{}.TableName(),
			result.ID,
			changedDataJSON,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			fmt.Printf("Failed to log bulk delete activity: %v\n", err)
		}
	}()
	return result, err
}

func (pcs *ProjectClassificationService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedProjectClassification []map[string]interface{}

	for _, id := range ids {
		projectClassification, err := pcs.repo.GetByID(id)
		if err != nil {
			fmt.Print(err)
			continue
		}

		err = pcs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedProjectClassification = append(deletedProjectClassification, map[string]interface{}{
			"id":       id,
			"title_en": projectClassification.TitleEn,
			"title_ar": projectClassification.TitleAr,
		})
	}

	if len(successfulIds) > 0 {
		go func() {
			err := pcs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.ProjectClassification{}.TableName(),
				successfulIds,
				deletedProjectClassification,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (pcs *ProjectClassificationService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := pcs.repo.Activate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := pcs.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.ProjectClassification{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return err
}

func (pcs *ProjectClassificationService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := pcs.repo.Deactivate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := pcs.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.ProjectClassification{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return err
}