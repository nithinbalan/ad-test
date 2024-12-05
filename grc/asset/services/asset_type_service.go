package services

import (
	"errors"
	"fmt"
	"jethur/grc/asset/repositories"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"
	"log"

	"gorm.io/gorm"
)

type AssetTypeServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AssetType, error)
	GetByID(id uint) (*models.AssetType, error)
	Create(currentUserID uint, assetType *models.AssetType) (*models.AssetType, error)
	Update(id uint, currentUserID uint, assetType *models.AssetType) (*models.AssetType, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type AssetTypeService struct {
	repo           repositories.AssetTypeRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewAssetTypeService(repo repositories.AssetTypeRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) AssetTypeServiceInterface {
	return &AssetTypeService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (ats *AssetTypeService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AssetType, error) {
	return ats.repo.GetAll(sort, sortBy, keyword, status)
}

func (ats AssetTypeService) GetByID(id uint) (*models.AssetType, error) {
	assetType, err := ats.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return assetType, nil
}

func (ats *AssetTypeService) Create(currentUserID uint, assetType *models.AssetType) (*models.AssetType, error) {
	if exists, err := ats.repo.ExistsByTitle(assetType.TitleEn, assetType.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	createdType, err := ats.repo.Create(assetType)
	if err != nil {
		return nil, err
	}

	go func() {
		err := ats.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AssetType{}.TableName(),
			createdType.ID,
			createdType,
			currentUserID,
			nil,
		)

		if err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdType, nil
}

func (ats *AssetTypeService) Update(id uint, currentUserID uint, updatedType *models.AssetType) (*models.AssetType, error) {
	existingType, err := ats.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if updatedType.TitleEn != existingType.TitleEn || updatedType.TitleAr != existingType.TitleAr {
		if exists, err := ats.repo.ExistsByTitleExcludingID(id, updatedType.TitleEn, updatedType.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	result, err := ats.repo.Update(id, updatedType)
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
		err := ats.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.AssetType{}.TableName(),
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

func (ats *AssetTypeService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {
		_, err := ats.repo.GetByID(id)
		if err != nil {
			fmt.Print(err)
			continue
		}

		err = ats.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {
		go func() {
			err := ats.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.AssetType{}.TableName(),
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

func (ats *AssetTypeService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := ats.repo.Activate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := ats.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.AssetType{}.TableName(),
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

func (ats *AssetTypeService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := ats.repo.Deactivate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := ats.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.AssetType{}.TableName(),
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
