package services

import (
	"errors"
	"fmt"
	incidentDTOs "jethur/grc/incident/dtos"
	"jethur/grc/incident/repositories"
	activityLogServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"
	"log"

	"gorm.io/gorm"
)

type IncidentSubCategoryServiceInterface interface {
	GetAll(params *incidentDTOs.IncidentSubCategoryListQueryParamsDTO) ([]*incidentDTOs.IncidentSubCategoryDTO, error)
	GetByID(id uint) (*models.IncidentSubCategory, error)
	Create(currentUserID uint, incidentSubCategoryFormDTO *incidentDTOs.IncidentSubCategoryFormDTO) (uint, error)
	Update(id uint, currentUserID uint, incidentSubCategoryFormDTO *incidentDTOs.IncidentSubCategoryFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type IncidentSubCategoryService struct {
	repo                repositories.IncidentSubCategoryRepositoryInterface
	activityLogSvc      activityLogServices.ActivityLogServiceInterface
	incidentCategorySvc IncidentCategoryServiceInterface
}

func NewIncidentSubCategoryService(repo repositories.IncidentSubCategoryRepositoryInterface, activityLogSvc activityLogServices.ActivityLogServiceInterface, incidentCategorySvc IncidentCategoryServiceInterface) IncidentSubCategoryServiceInterface {
	return &IncidentSubCategoryService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		incidentCategorySvc: incidentCategorySvc,
	}
}

func (iscs *IncidentSubCategoryService) GetAll(params *incidentDTOs.IncidentSubCategoryListQueryParamsDTO) ([]*incidentDTOs.IncidentSubCategoryDTO, error) {
	return iscs.repo.GetAll(params)
}

func (iscs IncidentSubCategoryService) GetByID(id uint) (*models.IncidentSubCategory, error) {
	incidentSubCategory, err := iscs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return incidentSubCategory, nil
}

func (iscs *IncidentSubCategoryService) Create(currentUserID uint, incidentSubCategoryFormDTO *incidentDTOs.IncidentSubCategoryFormDTO) (uint, error) {
	if err := iscs.validateIncidentSubCategoryFormDTO(incidentSubCategoryFormDTO, nil); err != nil {
		return 0, err
	}

	incidentSubCategory := &models.IncidentSubCategory{
		IncidentCategoryID: incidentSubCategoryFormDTO.IncidentCategoryID,
		TitleEn:            incidentSubCategoryFormDTO.TitleEn,
		TitleAr:            incidentSubCategoryFormDTO.TitleAr,
		Status:             incidentSubCategoryFormDTO.Status,
	}

	id, err := iscs.repo.Create(incidentSubCategory)
	if err != nil {
		return 0, err
	}

	go func() {
		createdSubCategory, err := iscs.repo.GetByID(id)
		if err != nil {
			return
		}
		err = iscs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.IncidentSubCategory{}.TableName(),
			id,
			createdSubCategory,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.IncidentSubCategory{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil

}

func (iscs *IncidentSubCategoryService) Update(id uint, currentUserID uint, incidentSubCategoryFormDTO *incidentDTOs.IncidentSubCategoryFormDTO) error {
	existingIncidentSubCategory, err := iscs.GetByID(id)
	if err != nil {
		return err
	}

	incidentSubCategory, err := iscs.GetByID(id)
	if err != nil {
		return err
	}

	if err := iscs.validateIncidentSubCategoryFormDTO(incidentSubCategoryFormDTO, &id); err != nil {
		return err
	}

	existingIncidentSubCategory.IncidentCategoryID = incidentSubCategoryFormDTO.IncidentCategoryID
	existingIncidentSubCategory.TitleEn = incidentSubCategoryFormDTO.TitleEn
	existingIncidentSubCategory.TitleAr = incidentSubCategoryFormDTO.TitleAr
	existingIncidentSubCategory.Status = incidentSubCategoryFormDTO.Status

	err = iscs.repo.Update(existingIncidentSubCategory)

	if err != nil {
		return err
	}

	fieldsToCompare := []string{"IncidentCategoryID", "TitleEn", "TitleAr", "Status"}
	changes, err := utility.CompareFields(incidentSubCategory, existingIncidentSubCategory, fieldsToCompare)

	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}
	go func() {
		err := iscs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.IncidentSubCategory{}.TableName(),
			id,
			changedDataJSON,
			currentUserID,
			nil,
		)
		if err != nil {
			// @todo - save to log file
			fmt.Printf("Failed to log bulk delete activity: %v\n", err)
		}
	}()

	return nil

}

func (iscs *IncidentSubCategoryService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedIncidentSubCategory []map[string]interface{}
	for _, id := range ids {

		subCategory, err := iscs.repo.GetByID(id)

		if err != nil {
			fmt.Print(err)
			continue
		}

		err = iscs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
		deletedIncidentSubCategory = append(deletedIncidentSubCategory, map[string]interface{}{
			"id":       id,
			"title_en": subCategory.TitleEn,
			"title_ar": subCategory.TitleAr,
		})
	}

	if len(successfulIds) > 0 {

		go func() {
			err := iscs.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.IncidentSubCategory{}.TableName(),
				successfulIds,
				deletedIncidentSubCategory,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (iscs *IncidentSubCategoryService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := iscs.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := iscs.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.IncidentSubCategory{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return err
}

func (iscs *IncidentSubCategoryService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := iscs.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := iscs.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.IncidentSubCategory{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				// You might want to log this error or handle it appropriately
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return err
}

func (iscs *IncidentSubCategoryService) validateIncidentSubCategoryFormDTO(incidentSubCategoryFormDTO *incidentDTOs.IncidentSubCategoryFormDTO, id *uint) error {
	if _, err := iscs.incidentCategorySvc.GetByID(incidentSubCategoryFormDTO.IncidentCategoryID); err != nil {
		return utility.NewValidationError(fmt.Errorf("incident category does not exist"), "not_exist", "incident category")
	}

	if id == nil {
		if exists, err := iscs.repo.ExistsByTitle(incidentSubCategoryFormDTO.TitleEn, incidentSubCategoryFormDTO.TitleAr, incidentSubCategoryFormDTO.IncidentCategoryID); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	} else {
		if exists, err := iscs.repo.ExistsByTitleExcludingID(*id, incidentSubCategoryFormDTO.TitleEn, incidentSubCategoryFormDTO.TitleAr, incidentSubCategoryFormDTO.IncidentCategoryID); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	return nil
}
