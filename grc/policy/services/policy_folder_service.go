package services

import (
	"fmt"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	policyDTOs "jethur/grc/policy/dtos"
	"jethur/grc/policy/repositories"
	"jethur/utility"
	"log"
)

type PolicyFolderServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*policyDTOs.PolicyFolderListDTO, error)
	GetByID(id uint) (*models.PolicyFolder, error)
	Create(currentUserID uint, folder *models.PolicyFolder) (*models.PolicyFolder, error)
	Update(id uint, currentUserID uint, folder *models.PolicyFolder) (*models.PolicyFolder, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type PolicyFolderService struct {
	repo           repositories.PolicyFolderRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewPolicyFolderService(repo repositories.PolicyFolderRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) PolicyFolderServiceInterface {
	return &PolicyFolderService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (cfs *PolicyFolderService) GetAll(sort string, sortBy string, keyword string, status string) ([]*policyDTOs.PolicyFolderListDTO, error) {
	return cfs.repo.GetAll(sort, sortBy, keyword, status)
}

func (cfs PolicyFolderService) GetByID(id uint) (*models.PolicyFolder, error) {
	return cfs.repo.GetByID(id)
}

func (cfs *PolicyFolderService) Create(currentUserID uint, folder *models.PolicyFolder) (*models.PolicyFolder, error) {
	if exists, err := cfs.repo.ExistsByTitle(folder.TitleEn, folder.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	// Set the CreatedByID field
	folder.CreatedBYID = currentUserID

	createdFolder, err := cfs.repo.Create(folder)
	if err != nil {
		return nil, err
	}

	go func() {
		err := cfs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.PolicyFolder{}.TableName(),
			folder.ID,
			createdFolder,
			currentUserID,
			nil,
		)

		if err != nil {
			// @todo - save to log file
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdFolder, nil
}

func (cfs *PolicyFolderService) Update(id uint, currentUserID uint, updatedFolder *models.PolicyFolder) (*models.PolicyFolder, error) {
	existingFolder, err := cfs.repo.GetByID(id)

	if err != nil {
		return nil, err
	}

	if updatedFolder.TitleEn != existingFolder.TitleEn || updatedFolder.TitleAr != existingFolder.TitleAr {
		if exists, err := cfs.repo.ExistsByTitleExcludingID(id, updatedFolder.TitleEn, updatedFolder.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	updatedFolder.UpdatedBYID = &currentUserID

	result, err := cfs.repo.Update(id, updatedFolder)

	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status"}
	changes, err := utility.CompareFields(existingFolder, updatedFolder, fieldsToCompare)

	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := cfs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.PolicyFolder{}.TableName(),
			result.ID,
			changedDataJSON,
			currentUserID,
			nil,
		)
		if err != nil {
			// @todo - save to log file
			fmt.Printf("Failed to log activity: %v\n", err)
		}
	}()
	return result, err
}

func (cfs *PolicyFolderService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedFolders []map[string]interface{}

	for _, id := range ids {
		folder, err := cfs.repo.GetByID(id)
		if err != nil {
			continue
		}

		err = cfs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedFolders = append(deletedFolders, map[string]interface{}{
			"id":       id,
			"title_en": folder.TitleEn,
			"title_ar": folder.TitleAr,
		})
	}

	if len(successfulIds) > 0 {
		go func() {
			err := cfs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.PolicyFolder{}.TableName(),
				successfulIds,
				deletedFolders,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (cfs *PolicyFolderService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := cfs.repo.Activate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := cfs.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.PolicyFolder{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log activity: %v\n", err)
			}
		}()
	}

	return err
}

func (cfs *PolicyFolderService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := cfs.repo.Deactivate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := cfs.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.PolicyFolder{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log activity: %v\n", err)
			}
		}()
	}

	return err
}
