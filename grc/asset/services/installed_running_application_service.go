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

type InstalledRunningApplicationServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.InstalledRunningApplication, error)
	GetByID(id uint) (*models.InstalledRunningApplication, error)
	Create(currentUserID uint, installedRunningApplicationForm *dtos.InstalledRunningApplicationFormDTO) (*models.InstalledRunningApplication, error)
	Update(id uint, currentUserID uint, installedRunningApplicationForm *dtos.InstalledRunningApplicationFormDTO) (*models.InstalledRunningApplication, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

type InstalledRunningApplicationService struct {
	repo           repositories.InstalledRunningApplicationRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewInstalledRunningApplicationService(
	repo repositories.InstalledRunningApplicationRepositoryInterface,
	activityLogSvc services.ActivityLogServiceInterface,
) InstalledRunningApplicationServiceInterface {
	return &InstalledRunningApplicationService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (iras *InstalledRunningApplicationService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.InstalledRunningApplication, error) {
	return iras.repo.GetAll(sort, sortBy, keyword, status)
}

func (iras InstalledRunningApplicationService) GetByID(id uint) (*models.InstalledRunningApplication, error) {
	installedRunningApplication, err := iras.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return installedRunningApplication, nil
}

func (iras *InstalledRunningApplicationService) Create(currentUserID uint, installedRunningApplicationForm *dtos.InstalledRunningApplicationFormDTO) (*models.InstalledRunningApplication, error) {
	if exists, err := iras.repo.ExistsByTitle(installedRunningApplicationForm.TitleEn, installedRunningApplicationForm.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	installedRunningApplication := &models.InstalledRunningApplication{
		TitleEn: installedRunningApplicationForm.TitleEn,
		TitleAr: installedRunningApplicationForm.TitleAr,
		Status:  installedRunningApplicationForm.Status,
	}

	createdApplication, err := iras.repo.Create(installedRunningApplication)
	if err != nil {
		return nil, err
	}

	go func() {
		err := iras.activityLogSvc.LogActivity(
			models.CreateAction,
			models.InstalledRunningApplication{}.TableName(),
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

func (iras *InstalledRunningApplicationService) Update(id uint, currentUserID uint, installedRunningApplicationForm *dtos.InstalledRunningApplicationFormDTO) (*models.InstalledRunningApplication, error) {
	existingApplication, err := iras.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if installedRunningApplicationForm.TitleEn != existingApplication.TitleEn || installedRunningApplicationForm.TitleAr != existingApplication.TitleAr {
		if exists, err := iras.repo.ExistsByTitleExcludingID(id, installedRunningApplicationForm.TitleEn, installedRunningApplicationForm.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	updatedApplication := &models.InstalledRunningApplication{
		TitleEn: installedRunningApplicationForm.TitleEn,
		TitleAr: installedRunningApplicationForm.TitleAr,
		Status:  installedRunningApplicationForm.Status,
	}

	result, err := iras.repo.Update(id, updatedApplication)
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
		err := iras.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.InstalledRunningApplication{}.TableName(),
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

func (iras *InstalledRunningApplicationService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedInstalledRunningApplication []map[string]interface{}

	for _, id := range ids {
		installedRunningApplication, err := iras.repo.GetByID(id)
		if err != nil {
			fmt.Print(err)
			continue
		}

		err = iras.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedInstalledRunningApplication = append(deletedInstalledRunningApplication, map[string]interface{}{
			"id":       id,
			"title_en": installedRunningApplication.TitleEn,
			"title_ar": installedRunningApplication.TitleAr,
		})
	}

	if len(successfulIds) > 0 {
		go func() {
			err := iras.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.InstalledRunningApplication{}.TableName(),
				successfulIds,
				deletedInstalledRunningApplication,
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

func (iras *InstalledRunningApplicationService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := iras.repo.Activate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := iras.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.InstalledRunningApplication{}.TableName(),
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

func (iras *InstalledRunningApplicationService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := iras.repo.Deactivate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := iras.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.InstalledRunningApplication{}.TableName(),
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