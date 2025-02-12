package services

import (
	"errors"
	"fmt"
	assetDTOs "jethur/grc/asset/dtos"
	"jethur/grc/asset/repositories"
	activityLogServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"
	"log"

	"gorm.io/gorm"
)

type AssetCategoryServiceInterface interface {
	GetAll(params *assetDTOs.AssetCategoryListQueryParamsDTO) ([]*assetDTOs.AssetCategoryDTO, error)
	GetByID(id uint) (*models.AssetCategory, error)
	Create(currentUserID uint, assetCategoryFormDTO *assetDTOs.AssetCategoryFormDTO) (uint, error)
	Update(id uint, currentUserID uint, assetCategoryFormDTO *assetDTOs.AssetCategoryFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type AssetCategoryService struct {
	repo           repositories.AssetCategoryRepositoryInterface
	activityLogSvc activityLogServices.ActivityLogServiceInterface
}

func NewAssetCategoryService(repo repositories.AssetCategoryRepositoryInterface, activityLogSvc activityLogServices.ActivityLogServiceInterface) AssetCategoryServiceInterface {
	return &AssetCategoryService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (acs *AssetCategoryService) GetAll(params *assetDTOs.AssetCategoryListQueryParamsDTO) ([]*assetDTOs.AssetCategoryDTO, error) {
	return acs.repo.GetAll(params)
}

func (acs AssetCategoryService) GetByID(id uint) (*models.AssetCategory, error) {
	assetCategory, err := acs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return assetCategory, nil
}

func (acs *AssetCategoryService) Create(currentUserID uint, assetCategoryFormDTO *assetDTOs.AssetCategoryFormDTO) (uint, error) {
	if err := acs.validateAssetCategoryFormDTO(assetCategoryFormDTO, nil); err != nil {
		return 0, err
	}

	assetCategory := &models.AssetCategory{
		TitleEn: assetCategoryFormDTO.TitleEn,
		TitleAr: assetCategoryFormDTO.TitleAr,
		Status:  assetCategoryFormDTO.Status,
	}

	id, err := acs.repo.Create(assetCategory)
	if err != nil {
		return 0, err
	}

	go func() {
		createdAssetCategory, err := acs.GetByID(id)
		if err != nil {
			return
		}

		err = acs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AssetCategory{}.TableName(),
			id,
			createdAssetCategory,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AssetCategory{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil

}

func (acs *AssetCategoryService) Update(id uint, currentUserID uint, assetCategoryFormDTO *assetDTOs.AssetCategoryFormDTO) error {
	oldAssetCategory, err := acs.GetByID(id)
	if err != nil {
		return err
	}

	existingAssetCategory, err := acs.GetByID(id)
	if err != nil {
		return err
	}

	if err := acs.validateAssetCategoryFormDTO(assetCategoryFormDTO, &id); err != nil {
		return err
	}

	existingAssetCategory.TitleEn = assetCategoryFormDTO.TitleEn
	existingAssetCategory.TitleAr = assetCategoryFormDTO.TitleAr
	existingAssetCategory.Status = assetCategoryFormDTO.Status

	err = acs.repo.Update(existingAssetCategory)
	if err != nil {
		return err
	}

	fieldsToCompare := []string{"AssetTypeID", "TitleEn", "TitleAr", "Status"}
	changes, err := utility.CompareFields(oldAssetCategory, existingAssetCategory, fieldsToCompare)

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

func (acs *AssetCategoryService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedAssetCategory []map[string]interface{}
	for _, id := range ids {
		assetCategory, err := acs.repo.GetByID(id)
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
		deletedAssetCategory = append(deletedAssetCategory, map[string]interface{}{
			"id":       id,
			"title_en": assetCategory.TitleEn,
			"title_ar": assetCategory.TitleAr,
		})
	}

	if len(successfulIds) > 0 {
		go func() {
			err := acs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.AssetCategory{}.TableName(),
				successfulIds,
				deletedAssetCategory,
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

func (acs *AssetCategoryService) validateAssetCategoryFormDTO(assetCategoryFormDTO *assetDTOs.AssetCategoryFormDTO, id *uint) error {
	if id == nil {
		if exists, err := acs.repo.ExistsByTitle(assetCategoryFormDTO.TitleEn, assetCategoryFormDTO.TitleAr); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	} else {
		if exists, err := acs.repo.ExistsByTitleExcludingID(*id, assetCategoryFormDTO.TitleEn, assetCategoryFormDTO.TitleAr); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	return nil
}
