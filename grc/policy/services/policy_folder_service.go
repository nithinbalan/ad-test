package services

import (
	"fmt"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/policy/dtos"
	policyDTOs "jethur/grc/policy/dtos"
	"jethur/grc/policy/repositories"
	"jethur/utility"
	"log"
)

type PolicyFolderServiceInterface interface {
	GetAll(params *policyDTOs.PolicyRegisterListQueryParamsDTO) (*policyDTOs.PolicyRegisterListResponseDTO, error)
	GetByID(id uint) (*models.Document, error)
	Create(currentUserID uint, formDTO *dtos.PolicyFolderFormDTO) (*models.Document, error)
	Update(id uint, currentUserID uint, formDTO *dtos.PolicyFolderFormDTO) error
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

func (cfs *PolicyFolderService) GetAll(params *policyDTOs.PolicyRegisterListQueryParamsDTO) (*policyDTOs.PolicyRegisterListResponseDTO, error) {
	return cfs.repo.GetAll(params)
}

func (cfs PolicyFolderService) GetByID(id uint) (*models.Document, error) {
	return cfs.repo.GetByID(id)
}

func (cfs *PolicyFolderService) Create(currentUserID uint, formDTO *dtos.PolicyFolderFormDTO) (*models.Document, error) {
	if exists, err := cfs.repo.ExistsByTitle(formDTO.TitleEn, formDTO.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	folder := &models.Document{
		TitleEn:              formDTO.TitleEn,
		TitleAr:              formDTO.TitleEn,
		IsFolder:             true,
		ParentDocumentID:     formDTO.ParentDocumentID,
		DocumentAccessTypeID: formDTO.DocumentAccessTypeID,
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
			models.Document{}.TableName(),
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

func (cfs *PolicyFolderService) Update(id uint, currentUserID uint, formDTO *dtos.PolicyFolderFormDTO) error {
	existingFolder, err := cfs.repo.GetByID(id)

	if err != nil {
		return err
	}

	if formDTO.TitleEn != existingFolder.TitleEn || formDTO.TitleAr != existingFolder.TitleAr {
		if exists, err := cfs.repo.ExistsByTitleExcludingID(id, formDTO.TitleEn, formDTO.TitleAr); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	existingFolder.UpdatedBYID = &currentUserID

	result, err := cfs.repo.Update(id, existingFolder)

	if err != nil {
		return err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status"}
	changes, err := utility.CompareFields(existingFolder, existingFolder, fieldsToCompare)

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
			models.Document{}.TableName(),
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
	return err
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
				models.Document{}.TableName(),
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
				models.Document{}.TableName(),
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
				models.Document{}.TableName(),
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
