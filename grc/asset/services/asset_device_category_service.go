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

type AssetDeviceCategoryServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AssetDeviceCategory, error)
	GetByID(id uint) (*models.AssetDeviceCategory, error)
	Create(currentUserID uint, assetDeviceCategoryForm *dtos.AssetDeviceCategoryFormDTO) (*models.AssetDeviceCategory, error)
	Update(id uint, currentUserID uint, assetDeviceCategoryForm *dtos.AssetDeviceCategoryFormDTO) (*models.AssetDeviceCategory, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

type AssetDeviceCategoryService struct {
	repo           repositories.AssetDeviceCategoryRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewAssetDeviceCategoryService(
	repo repositories.AssetDeviceCategoryRepositoryInterface,
	activityLogSvc services.ActivityLogServiceInterface,
) AssetDeviceCategoryServiceInterface {
	return &AssetDeviceCategoryService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (acs *AssetDeviceCategoryService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AssetDeviceCategory, error) {
	return acs.repo.GetAll(sort, sortBy, keyword, status)
}

func (acs AssetDeviceCategoryService) GetByID(id uint) (*models.AssetDeviceCategory, error) {
	assetDeviceCategory, err := acs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return assetDeviceCategory, nil
}

func (acs *AssetDeviceCategoryService) Create(currentUserID uint, assetDeviceCategoryForm *dtos.AssetDeviceCategoryFormDTO) (*models.AssetDeviceCategory, error) {
	if exists, err := acs.repo.ExistsByTitle(assetDeviceCategoryForm.TitleEn, assetDeviceCategoryForm.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	assetDeviceCategory := &models.AssetDeviceCategory{
		TitleEn:  assetDeviceCategoryForm.TitleEn,
		TitleAr:  assetDeviceCategoryForm.TitleAr,
		IsSystem: false,
		Status:   assetDeviceCategoryForm.Status,
	}

	createdCategory, err := acs.repo.Create(assetDeviceCategory)
	if err != nil {
		return nil, err
	}

	go func() {
		err := acs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AssetDeviceCategory{}.TableName(),
			createdCategory.ID,
			createdCategory,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdCategory, nil
}

func (acs *AssetDeviceCategoryService) Update(id uint, currentUserID uint, assetDeviceCategoryForm *dtos.AssetDeviceCategoryFormDTO) (*models.AssetDeviceCategory, error) {
	existingCategory, err := acs.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if existingCategory.IsSystem {
		return nil, utility.NewValidationError(fmt.Errorf("system categories cannot be modified"), "immutable", "is_system")
	}

	if assetDeviceCategoryForm.TitleEn != existingCategory.TitleEn || assetDeviceCategoryForm.TitleAr != existingCategory.TitleAr {
		if exists, err := acs.repo.ExistsByTitleExcludingID(id, assetDeviceCategoryForm.TitleEn, assetDeviceCategoryForm.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	updatedCategory := &models.AssetDeviceCategory{
		TitleEn: assetDeviceCategoryForm.TitleEn,
		TitleAr: assetDeviceCategoryForm.TitleAr,
		Status:  assetDeviceCategoryForm.Status,
	}

	result, err := acs.repo.Update(id, updatedCategory)
	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status"}
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
			models.AssetDeviceCategory{}.TableName(),
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

func (acs *AssetDeviceCategoryService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedAssetDeviceCategory []map[string]interface{}

	for _, id := range ids {
		assetDeviceCategory, err := acs.repo.GetByID(id)
		if err != nil {
			fmt.Print(err)
			continue
		}

		if assetDeviceCategory.IsSystem {
			failedIds = append(failedIds, id)
			continue
		}

		err = acs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedAssetDeviceCategory = append(deletedAssetDeviceCategory, map[string]interface{}{
			"id":       id,
			"title_en": assetDeviceCategory.TitleEn,
			"title_ar": assetDeviceCategory.TitleAr,
		})
	}

	if len(successfulIds) > 0 {
		go func() {
			err := acs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.AssetDeviceCategory{}.TableName(),
				successfulIds,
				deletedAssetDeviceCategory,
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

func (acs *AssetDeviceCategoryService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := acs.repo.Activate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := acs.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.AssetDeviceCategory{}.TableName(),
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

func (acs *AssetDeviceCategoryService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := acs.repo.Deactivate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := acs.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.AssetDeviceCategory{}.TableName(),
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