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

type AssetDocumentCategoryServiceInterface interface {
	GetAll(params *assetDTOs.AssetDocumentCategoryListQueryParamsDTO) ([]*assetDTOs.AssetDocumentCategoryDTO, error)
	GetByID(id uint) (*models.AssetDocumentCategory, error)
	Create(currentUserID uint, assetDocumentCategoryFormDTO *assetDTOs.AssetDocumentCategoryFormDTO) (uint, error)
	Update(id uint, currentUserID uint, assetDocumentCategoryFormDTO *assetDTOs.AssetDocumentCategoryFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type AssetDocumentCategoryService struct {
	repo           repositories.AssetDocumentCategoryRepositoryInterface
	assetTypeSvc   AssetTypeServiceInterface
	activityLogSvc activityLogServices.ActivityLogServiceInterface
}

func NewAssetDocumentCategoryService(repo repositories.AssetDocumentCategoryRepositoryInterface, assetTypeSvc AssetTypeServiceInterface, activityLogSvc activityLogServices.ActivityLogServiceInterface) AssetDocumentCategoryServiceInterface {
	return &AssetDocumentCategoryService{
		repo:           repo,
		assetTypeSvc:   assetTypeSvc,
		activityLogSvc: activityLogSvc,
	}
}

func (adcs *AssetDocumentCategoryService) GetAll(params *assetDTOs.AssetDocumentCategoryListQueryParamsDTO) ([]*assetDTOs.AssetDocumentCategoryDTO, error) {
	return adcs.repo.GetAll(params)
}

func (adcs AssetDocumentCategoryService) GetByID(id uint) (*models.AssetDocumentCategory, error) {
	assetDocumentCategory, err := adcs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return assetDocumentCategory, nil
}

func (adcs *AssetDocumentCategoryService) Create(currentUserID uint, assetDocumentCategoryFormDTO *assetDTOs.AssetDocumentCategoryFormDTO) (uint, error) {
	if err := adcs.validateAssetDocumentCategoryFormDTO(assetDocumentCategoryFormDTO, nil); err != nil {
		return 0, err
	}

	assetDocumentCategory := &models.AssetDocumentCategory{
		TitleEn:       assetDocumentCategoryFormDTO.TitleEn,
		TitleAr:       assetDocumentCategoryFormDTO.TitleAr,
		DescriptionEn: assetDocumentCategoryFormDTO.DescriptionEn,
		DescriptionAr: assetDocumentCategoryFormDTO.DescriptionAr,
		Status:        assetDocumentCategoryFormDTO.Status,
	}

	id, err := adcs.repo.Create(assetDocumentCategory)
	if err != nil {
		return 0, err
	}

	go func() {
		createdAssetDocumentCategory, err := adcs.GetByID(id)
		if err != nil {
			return
		}

		err = adcs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AssetDocumentCategory{}.TableName(),
			id,
			createdAssetDocumentCategory,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AssetDocumentCategory{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil

}

func (adcs *AssetDocumentCategoryService) Update(id uint, currentUserID uint, assetDocumentCategoryFormDTO *assetDTOs.AssetDocumentCategoryFormDTO) error {
	oldAssetDocumentCategory, err := adcs.GetByID(id)
	if err != nil {
		return err
	}

	existingAssetDocumentCategory, err := adcs.GetByID(id)
	if err != nil {
		return err
	}

	if err := adcs.validateAssetDocumentCategoryFormDTO(assetDocumentCategoryFormDTO, &id); err != nil {
		return err
	}

	existingAssetDocumentCategory.TitleEn = assetDocumentCategoryFormDTO.TitleEn
	existingAssetDocumentCategory.TitleAr = assetDocumentCategoryFormDTO.TitleAr
	existingAssetDocumentCategory.DescriptionEn = assetDocumentCategoryFormDTO.DescriptionEn
	existingAssetDocumentCategory.DescriptionAr = assetDocumentCategoryFormDTO.DescriptionAr
	existingAssetDocumentCategory.Status = assetDocumentCategoryFormDTO.Status

	err = adcs.repo.Update(existingAssetDocumentCategory)
	if err != nil {
		return err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status"}
	changes, err := utility.CompareFields(oldAssetDocumentCategory, existingAssetDocumentCategory, fieldsToCompare)

	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := adcs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.AssetDocumentCategory{}.TableName(),
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

func (adcs *AssetDocumentCategoryService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedAssetDocumentCategory []map[string]interface{}
	for _, id := range ids {
		assetDocumentCategory, err := adcs.repo.GetByID(id)
		if err != nil {
			fmt.Print(err)
			continue
		}

		err = adcs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
		deletedAssetDocumentCategory = append(deletedAssetDocumentCategory, map[string]interface{}{
			"id":       id,
			"title_en": assetDocumentCategory.TitleEn,
			"title_ar": assetDocumentCategory.TitleAr,
		})
	}

	if len(successfulIds) > 0 {
		go func() {
			err := adcs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.AssetDocumentCategory{}.TableName(),
				successfulIds,
				deletedAssetDocumentCategory,
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

func (adcs *AssetDocumentCategoryService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := adcs.repo.Activate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := adcs.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.AssetDocumentCategory{}.TableName(),
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

func (adcs *AssetDocumentCategoryService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := adcs.repo.Deactivate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := adcs.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.AssetDocumentCategory{}.TableName(),
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

func (adcs *AssetDocumentCategoryService) validateAssetDocumentCategoryFormDTO(assetDocumentCategoryFormDTO *assetDTOs.AssetDocumentCategoryFormDTO, id *uint) error {
	if id == nil {
		if exists, err := adcs.repo.ExistsByTitle(assetDocumentCategoryFormDTO.TitleEn, assetDocumentCategoryFormDTO.TitleAr); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	} else {
		if exists, err := adcs.repo.ExistsByTitleExcludingID(*id, assetDocumentCategoryFormDTO.TitleEn, assetDocumentCategoryFormDTO.TitleAr); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	return nil
}
