package services

import (
	"errors"
	"fmt"
	"log"

	"jethur/grc/asset/dtos"
	"jethur/grc/asset/repositories"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type InstalledApplicationServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.InstalledApplication, error)
	GetByID(id uint) (*models.InstalledApplication, error)
	Create(currentUserID uint, installedApplicationForm *dtos.InstalledApplicationFormDTO) (*models.InstalledApplication, error)
	Update(id uint, currentUserID uint, installedApplicationForm *dtos.InstalledApplicationFormDTO) (*models.InstalledApplication, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

type InstalledApplicationService struct {
	repo           repositories.InstalledApplicationRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewInstalledApplicationService(
	repo repositories.InstalledApplicationRepositoryInterface,
	activityLogSvc services.ActivityLogServiceInterface,
) InstalledApplicationServiceInterface {
	return &InstalledApplicationService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (ias *InstalledApplicationService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.InstalledApplication, error) {
	return ias.repo.GetAll(sort, sortBy, keyword, status)
}

func (ias InstalledApplicationService) GetByID(id uint) (*models.InstalledApplication, error) {
	installedApplication, err := ias.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return installedApplication, nil
}

func (ias *InstalledApplicationService) Create(currentUserID uint, installedApplicationForm *dtos.InstalledApplicationFormDTO) (*models.InstalledApplication, error) {
	if exists, err := ias.repo.ExistsByTitle(installedApplicationForm.TitleEn, installedApplicationForm.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	installedApplication := &models.InstalledApplication{
		TitleEn: installedApplicationForm.TitleEn,
		TitleAr: installedApplicationForm.TitleAr,
		Status:  installedApplicationForm.Status,
	}

	createdApplication, err := ias.repo.Create(installedApplication)
	if err != nil {
		return nil, err
	}

	go func() {
		err := ias.activityLogSvc.LogActivity(
			models.CreateAction,
			models.InstalledApplication{}.TableName(),
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

func (ias *InstalledApplicationService) Update(id uint, currentUserID uint, installedApplicationForm *dtos.InstalledApplicationFormDTO) (*models.InstalledApplication, error) {
	existingApplication, err := ias.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if installedApplicationForm.TitleEn != existingApplication.TitleEn || installedApplicationForm.TitleAr != existingApplication.TitleAr {
		if exists, err := ias.repo.ExistsByTitleExcludingID(id, installedApplicationForm.TitleEn, installedApplicationForm.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	updatedApplication := &models.InstalledApplication{
		TitleEn: installedApplicationForm.TitleEn,
		TitleAr: installedApplicationForm.TitleAr,
		Status:  installedApplicationForm.Status,
	}

	result, err := ias.repo.Update(id, updatedApplication)
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
		err := ias.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.InstalledApplication{}.TableName(),
			result.ID,
			changedDataJSON,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			fmt.Printf("Failed to log bulk delete activity: %v\n", err)
		}
	}()
	return result, nil
}

func (ias *InstalledApplicationService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedInstalledApplication []map[string]interface{}

	for _, id := range ids {
		installedApplication, err := ias.repo.GetByID(id)
		if err != nil {
			fmt.Print(err)
			continue
		}

		err = ias.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedInstalledApplication = append(deletedInstalledApplication, map[string]interface{}{
			"id":       id,
			"title_en": installedApplication.TitleEn,
			"title_ar": installedApplication.TitleAr,
		})
	}

	if len(successfulIds) > 0 {
		go func() {
			err := ias.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.InstalledApplication{}.TableName(),
				successfulIds,
				deletedInstalledApplication,
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

func (ias *InstalledApplicationService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := ias.repo.Activate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := ias.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.InstalledApplication{}.TableName(),
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

func (ias *InstalledApplicationService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := ias.repo.Deactivate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := ias.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.InstalledApplication{}.TableName(),
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