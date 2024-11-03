package services

import (
	"errors"
	"fmt"
	"jethur/grc/asset/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"
	"log"

	"gorm.io/gorm"
)

type AssetCategoryServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AssetCategory, error)
	GetByID(id uint) (*models.AssetCategory, error)
	Create(currentUserID uint, assetCategory *models.AssetCategory) (*models.AssetCategory, error)
	Update(id uint, currentUserID uint, assetCategory *models.AssetCategory) (*models.AssetCategory, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service implementation
type AssetCategoryService struct {
	repo           repositories.AssetCategoryRepositoryInterface
	assetTypeSvc   AssetTypeServiceInterface
	activityLogSvc logServices.ActivityLogServiceInterface
}

func NewAssetCategoryService(
	repo repositories.AssetCategoryRepositoryInterface,
	assetTypeSvc AssetTypeServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
) AssetCategoryServiceInterface {
	return &AssetCategoryService{
		repo:           repo,
		assetTypeSvc:   assetTypeSvc,
		activityLogSvc: activityLogSvc,
	}
}

func (acs *AssetCategoryService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AssetCategory, error) {
	return acs.repo.GetAll(sort, sortBy, keyword, status)
}

func (acs *AssetCategoryService) GetByID(id uint) (*models.AssetCategory, error) {
	assetCategory, err := acs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return assetCategory, nil
}

func (acs *AssetCategoryService) Create(currentUserID uint, assetCategory *models.AssetCategory) (*models.AssetCategory, error) {
	if assetCategory.AssetTypeID != nil {
		// Validate asset type exists
		_, err := acs.assetTypeSvc.GetByID(*assetCategory.AssetTypeID)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("asset type does not exist"), "not_exist", "asset type")
		}
	}

	if exists, err := acs.repo.ExistsByTitle(assetCategory.TitleEn, assetCategory.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	createdCategory, err := acs.repo.Create(assetCategory)
	if err != nil {
		return nil, err
	}

	go func() {
		err := acs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AssetCategory{}.TableName(),
			createdCategory.ID,
			createdCategory,
			currentUserID,
			nil,
		)
		if err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdCategory, nil
}

func (acs *AssetCategoryService) Update(id uint, currentUserID uint, updatedCategory *models.AssetCategory) (*models.AssetCategory, error) {
	existingCategory, err := acs.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if updatedCategory.AssetTypeID != nil && *updatedCategory.AssetTypeID != *existingCategory.AssetTypeID {
		// Validate asset type exists
		_, err := acs.assetTypeSvc.GetByID(*updatedCategory.AssetTypeID)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("asset type does not exist"), "not_exist", "asset type")
		}
	}

	if updatedCategory.TitleEn != existingCategory.TitleEn || updatedCategory.TitleAr != existingCategory.TitleAr {
		if exists, err := acs.repo.ExistsByTitleExcludingID(id, updatedCategory.TitleEn, updatedCategory.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	result, err := acs.repo.Update(id, updatedCategory)
	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status", "AssetTypeID"}
	changes, err := utility.CompareFields(existingCategory, updatedCategory, fieldsToCompare)
	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := acs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.AssetCategory{}.TableName(),
			result.ID,
			changedDataJSON,
			currentUserID,
			nil,
		)
		if err != nil {
			fmt.Printf("Failed to log bulk delete activity: %v\n", err)
		}
	}()
	return result, err
}

func (acs *AssetCategoryService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {
		_, err := acs.repo.GetByID(id)
		if err != nil {
			fmt.Print(err)
			continue
		}

		err = acs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {
		go func() {
			err := acs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.AssetCategory{}.TableName(),
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

	return failedIds, nil
}

func (acs *AssetCategoryService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := acs.repo.Activate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := acs.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.AssetCategory{}.TableName(),
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

func (acs *AssetCategoryService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := acs.repo.Deactivate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := acs.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.AssetCategory{}.TableName(),
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
