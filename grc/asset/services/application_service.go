package services

import (
	"errors"
	"fmt"
	"log"

	"jethur/grc/asset/repositories"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ApplicationServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.Application, error)
	GetByID(id uint) (*models.Application, error)
	Create(currentUserID uint, application *models.Application) (*models.Application, error)
	Update(id uint, currentUserID uint, application *models.Application) (*models.Application, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

type ApplicationService struct {
	repo           repositories.ApplicationRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewApplicationService(repo repositories.ApplicationRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) ApplicationServiceInterface {
	return &ApplicationService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (as *ApplicationService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.Application, error) {
	return as.repo.GetAll(sort, sortBy, keyword, status)
}

func (as *ApplicationService) GetByID(id uint) (*models.Application, error) {
	application, err := as.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return application, nil
}

func (as *ApplicationService) Create(currentUserID uint, application *models.Application) (*models.Application, error) {
	if exists, err := as.repo.ExistsByTitle(application.TitleEn, application.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	createdApplication, err := as.repo.Create(application)
	if err != nil {
		return nil, err
	}

	go func() {
		err := as.activityLogSvc.LogActivity(
			models.CreateAction,
			models.Application{}.TableName(),
			createdApplication.ID,
			createdApplication,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdApplication, nil
}

func (as *ApplicationService) Update(id uint, currentUserID uint, updatedApplication *models.Application) (*models.Application, error) {
	existingApplication, err := as.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if updatedApplication.TitleEn != existingApplication.TitleEn || updatedApplication.TitleAr != existingApplication.TitleAr {
		if exists, err := as.repo.ExistsByTitleExcludingID(id, updatedApplication.TitleEn, updatedApplication.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	result, err := as.repo.Update(id, updatedApplication)
	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status"}
	changes, err := utility.CompareFields(existingApplication, updatedApplication, fieldsToCompare)
	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := as.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.Application{}.TableName(),
			result.ID,
			changedDataJSON,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			fmt.Printf("Failed to log update activity: %v\n", err)
		}
	}()

	return result, err
}

func (as *ApplicationService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedApplications []map[string]interface{}

	for _, id := range ids {
		application, err := as.repo.GetByID(id)
		if err != nil {
			fmt.Print(err)
			continue
		}

		err = as.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
		deletedApplications = append(deletedApplications, map[string]interface{}{
			"id":       id,
			"title_en": application.TitleEn,
			"title_ar": application.TitleAr,
		})
	}

	if len(successfulIds) > 0 {
		go func() {
			err := as.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.Application{}.TableName(),
				successfulIds,
				deletedApplications,
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

func (as *ApplicationService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := as.repo.Activate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := as.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.Application{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk activate activity: %v\n", err)
			}
		}()
	}

	return err
}

func (as *ApplicationService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := as.repo.Deactivate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := as.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.Application{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk deactivate activity: %v\n", err)
			}
		}()
	}

	return err
}