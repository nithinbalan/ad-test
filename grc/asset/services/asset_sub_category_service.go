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

type AssetSubCategoryServiceInterface interface {
	GetAll(params *assetDTOs.AssetSubCategoryListQueryParamsDTO) ([]*assetDTOs.AssetSubCategoryDTO, error)
	GetByID(id uint) (*models.AssetSubCategory, error)
	Create(currentUserID uint, assetSubCategoryFormDTO *assetDTOs.AssetSubCategoryFormDTO) (uint, error)
	Update(id uint, currentUserID uint, assetSubCategoryFormDTO *assetDTOs.AssetSubCategoryFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type AssetSubCategoryService struct {
	repo             repositories.AssetSubCategoryRepositoryInterface
	assetCategorySvc AssetCategoryServiceInterface
	activityLogSvc   activityLogServices.ActivityLogServiceInterface
}

func NewAssetSubCategoryService(repo repositories.AssetSubCategoryRepositoryInterface, assetCategorySvc AssetCategoryServiceInterface, activityLogSvc activityLogServices.ActivityLogServiceInterface) AssetSubCategoryServiceInterface {
	return &AssetSubCategoryService{
		repo:             repo,
		assetCategorySvc: assetCategorySvc,
		activityLogSvc:   activityLogSvc,
	}
}

func (ascs *AssetSubCategoryService) GetAll(params *assetDTOs.AssetSubCategoryListQueryParamsDTO) ([]*assetDTOs.AssetSubCategoryDTO, error) {
	return ascs.repo.GetAll(params)
}

func (ascs AssetSubCategoryService) GetByID(id uint) (*models.AssetSubCategory, error) {
	assetSubCategory, err := ascs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return assetSubCategory, nil
}

func (ascs *AssetSubCategoryService) Create(currentUserID uint, assetSubCategoryFormDTO *assetDTOs.AssetSubCategoryFormDTO) (uint, error) {
	if err := ascs.validateAssetSubCategoryFormDTO(assetSubCategoryFormDTO, nil); err != nil {
		return 0, err
	}

	assetSubCategory := &models.AssetSubCategory{
		AssetCategoryID: assetSubCategoryFormDTO.AssetCategoryID,
		TitleEn:         assetSubCategoryFormDTO.TitleEn,
		TitleAr:         assetSubCategoryFormDTO.TitleAr,
		Status:          assetSubCategoryFormDTO.Status,
	}

	id, err := ascs.repo.Create(assetSubCategory)
	if err != nil {
		return 0, err
	}

	go func() {
		createdAssetSubCategory, err := ascs.GetByID(id)
		if err != nil {
			return
		}

		err = ascs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AssetSubCategory{}.TableName(),
			id,
			createdAssetSubCategory,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AssetSubCategory{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil

}

func (ascs *AssetSubCategoryService) Update(id uint, currentUserID uint, assetSubCategoryFormDTO *assetDTOs.AssetSubCategoryFormDTO) error {
	oldAssetSubCategory, err := ascs.GetByID(id)
	if err != nil {
		return err
	}

	existingAssetSubCategory, err := ascs.GetByID(id)
	if err != nil {
		return err
	}

	if err := ascs.validateAssetSubCategoryFormDTO(assetSubCategoryFormDTO, &id); err != nil {
		return err
	}

	existingAssetSubCategory.AssetCategoryID = assetSubCategoryFormDTO.AssetCategoryID
	existingAssetSubCategory.TitleEn = assetSubCategoryFormDTO.TitleEn
	existingAssetSubCategory.TitleAr = assetSubCategoryFormDTO.TitleAr
	existingAssetSubCategory.Status = assetSubCategoryFormDTO.Status

	err = ascs.repo.Update(existingAssetSubCategory)
	if err != nil {
		return err
	}

	fieldsToCompare := []string{"AssetCategoryID", "TitleEn", "TitleAr", "Status"}
	changes, err := utility.CompareFields(oldAssetSubCategory, existingAssetSubCategory, fieldsToCompare)

	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := ascs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.AssetSubCategory{}.TableName(),
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

func (ascs *AssetSubCategoryService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedAssetSubCategory []map[string]interface{}
	for _, id := range ids {
		assetSubCategory, err := ascs.repo.GetByID(id)
		if err != nil {
			fmt.Print(err)
			continue
		}

		err = ascs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
		deletedAssetSubCategory = append(deletedAssetSubCategory, map[string]interface{}{
			"id":       id,
			"title_en": assetSubCategory.TitleEn,
			"title_ar": assetSubCategory.TitleAr,
		})
	}

	if len(successfulIds) > 0 {
		go func() {
			err := ascs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.AssetSubCategory{}.TableName(),
				successfulIds,
				deletedAssetSubCategory,
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

func (ascs *AssetSubCategoryService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := ascs.repo.Activate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := ascs.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.AssetSubCategory{}.TableName(),
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

func (ascs *AssetSubCategoryService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := ascs.repo.Deactivate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := ascs.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.AssetSubCategory{}.TableName(),
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

func (ascs *AssetSubCategoryService) validateAssetSubCategoryFormDTO(assetSubCategoryFormDTO *assetDTOs.AssetSubCategoryFormDTO, id *uint) error {
	if assetSubCategoryFormDTO.AssetCategoryID != nil {
		if _, err := ascs.assetCategorySvc.GetByID(*assetSubCategoryFormDTO.AssetCategoryID); err != nil {
			return utility.NewValidationError(fmt.Errorf("asset type does not exist"), "not_exist", "asset type")
		}
	}

	if id == nil {
		if exists, err := ascs.repo.ExistsByTitle(assetSubCategoryFormDTO.TitleEn, assetSubCategoryFormDTO.TitleAr); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	} else {
		if exists, err := ascs.repo.ExistsByTitleExcludingID(*id, assetSubCategoryFormDTO.TitleEn, assetSubCategoryFormDTO.TitleAr); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	return nil
}
