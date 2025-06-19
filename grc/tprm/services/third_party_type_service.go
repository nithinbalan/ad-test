package services

import (
	"errors"
	"fmt"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/tprm/repositories"
	"jethur/utility"
	"log"

	"gorm.io/gorm"
)

type ThirdPartyTypeServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ThirdPartyType, error)
	GetByID(id uint) (*models.ThirdPartyType, error)
	Create(currentUserID uint, thirdPartyType *models.ThirdPartyType) (*models.ThirdPartyType, error)
	Update(id uint, currentUserID uint, thirdPartyType *models.ThirdPartyType) (*models.ThirdPartyType, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type ThirdPartyTypeService struct {
	repo           repositories.ThirdPartyTypeRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewThirdPartyTypeService(repo repositories.ThirdPartyTypeRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) ThirdPartyTypeServiceInterface {
	return &ThirdPartyTypeService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (tpts *ThirdPartyTypeService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ThirdPartyType, error) {
	return tpts.repo.GetAll(sort, sortBy, keyword, status)
}

func (tpts ThirdPartyTypeService) GetByID(id uint) (*models.ThirdPartyType, error) {
	thirdPartyType, err := tpts.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return thirdPartyType, nil

}

func (tpts *ThirdPartyTypeService) Create(currentUserID uint, thirdPartyType *models.ThirdPartyType) (*models.ThirdPartyType, error) {

	if exists, err := tpts.repo.ExistsByTitle(thirdPartyType.TitleEn, thirdPartyType.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	createdType, err := tpts.repo.Create(thirdPartyType)
	if err != nil {
		return nil, err
	}

	go func() {
		err := tpts.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ThirdPartyType{}.TableName(),
			createdType.ID,
			createdType,
			currentUserID,
			nil,
		)

		if err != nil {
			// @todo - save to log file
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdType, nil
}

func (tpts *ThirdPartyTypeService) Update(id uint, currentUserID uint, updatedType *models.ThirdPartyType) (*models.ThirdPartyType, error) {
	existingType, err := tpts.repo.GetByID(id)

	if err != nil {
		return nil, err
	}

	if updatedType.TitleEn != existingType.TitleEn || updatedType.TitleAr != existingType.TitleAr {
		if exists, err := tpts.repo.ExistsByTitleExcludingID(id, updatedType.TitleEn, updatedType.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	}

	result, err := tpts.repo.Update(id, updatedType)

	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status"}
	changes, err := utility.CompareFields(existingType, updatedType, fieldsToCompare)

	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := tpts.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.ThirdPartyType{}.TableName(),
			result.ID,
			changedDataJSON,
			currentUserID,
			nil,
		)
		if err != nil {
			// @todo - save to log file
			fmt.Printf("Failed to log bulk delete activity: %v\n", err)
		}
	}()
	return result, err
}

func (tpts *ThirdPartyTypeService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedThirdPartyType []map[string]interface{}
	for _, id := range ids {

		thirdPartyType, err := tpts.repo.GetByID(id)

		if err != nil {
			fmt.Print(err)
			continue
		}

		err = tpts.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
		deletedThirdPartyType = append(deletedThirdPartyType, map[string]interface{}{
			"id":       id,
			"title_en": thirdPartyType.TitleEn,
			"title_ar": thirdPartyType.TitleAr,
		})
	}

	if len(successfulIds) > 0 {

		go func() {
			err := tpts.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.ThirdPartyType{}.TableName(),
				successfulIds,
				deletedThirdPartyType,
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

func (tpts *ThirdPartyTypeService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := tpts.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := tpts.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.ThirdPartyType{}.TableName(),
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

func (tpts *ThirdPartyTypeService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := tpts.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := tpts.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.ThirdPartyType{}.TableName(),
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
