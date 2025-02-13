package services

import (
	"fmt"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/repositories"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"
	"log"
)

type ComplianceFolderServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*complianceDTOs.ComplianceFolderListDTO, error)
	GetByID(id uint) (*models.ComplianceFolder, error)
	Create(currentUserID uint, folder *models.ComplianceFolder) (*models.ComplianceFolder, error)
	Update(id uint, currentUserID uint, folder *models.ComplianceFolder) (*models.ComplianceFolder, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type ComplianceFolderService struct {
	repo           repositories.ComplianceFolderRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewComplianceFolderService(repo repositories.ComplianceFolderRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) ComplianceFolderServiceInterface {
	return &ComplianceFolderService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (cfs *ComplianceFolderService) GetAll(sort string, sortBy string, keyword string, status string) ([]*complianceDTOs.ComplianceFolderListDTO, error) {
	return cfs.repo.GetAll(sort, sortBy, keyword, status)
}

func (cfs ComplianceFolderService) GetByID(id uint) (*models.ComplianceFolder, error) {
	return cfs.repo.GetByID(id)
}

func (cfs *ComplianceFolderService) Create(currentUserID uint, folder *models.ComplianceFolder) (*models.ComplianceFolder, error) {
	if exists, err := cfs.repo.ExistsByTitle(folder.TitleEn, folder.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	createdFolder, err := cfs.repo.Create(folder)
	if err != nil {
		return nil, err
	}

	go func() {
		err := cfs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ComplianceFolder{}.TableName(),
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

func (cfs *ComplianceFolderService) Update(id uint, currentUserID uint, updatedFolder *models.ComplianceFolder) (*models.ComplianceFolder, error) {
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
			models.ComplianceFolder{}.TableName(),
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

func (cfs *ComplianceFolderService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
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
				models.ComplianceFolder{}.TableName(),
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

func (cfs *ComplianceFolderService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := cfs.repo.Activate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := cfs.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.ComplianceFolder{}.TableName(),
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

func (cfs *ComplianceFolderService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := cfs.repo.Deactivate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := cfs.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.ComplianceFolder{}.TableName(),
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
