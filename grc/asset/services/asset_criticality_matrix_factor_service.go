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

type AssetCriticalityMatrixFactorServiceInterface interface {
	GetAll(matrixID uint, sort string, sortBy string, keyword string, status string) ([]*models.AssetCriticalityMatrixFactor, error)
	GetByID(matrixID uint, id uint) (*models.AssetCriticalityMatrixFactor, error)
	Create(matrixID uint, currentUserID uint, factor *models.AssetCriticalityMatrixFactor) (*models.AssetCriticalityMatrixFactor, error)
	Update(matrixID uint, id uint, currentUserID uint, factor *models.AssetCriticalityMatrixFactor) (*models.AssetCriticalityMatrixFactor, error)
	Delete(matrixID uint, ids []uint, currentUserID uint) ([]uint, error)
}

type AssetCriticalityMatrixFactorService struct {
	repo                      repositories.AssetCriticalityMatrixFactorRepositoryInterface
	activityLogSvc            services.ActivityLogServiceInterface
	assetCriticalityMatrixSvc AssetCriticalityMatrixServiceInterface
}

func NewAssetCriticalityMatrixFactorService(repo repositories.AssetCriticalityMatrixFactorRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface, assetCriticalityMatrixSvc AssetCriticalityMatrixServiceInterface) AssetCriticalityMatrixFactorServiceInterface {
	return &AssetCriticalityMatrixFactorService{
		repo:                      repo,
		activityLogSvc:            activityLogSvc,
		assetCriticalityMatrixSvc: assetCriticalityMatrixSvc,
	}
}

func (acmfs *AssetCriticalityMatrixFactorService) GetAll(matrixID uint, sort, sortBy, keyword, status string) ([]*models.AssetCriticalityMatrixFactor, error) {
	return acmfs.repo.GetAll(matrixID, sort, sortBy, keyword, status)
}

func (acmfs *AssetCriticalityMatrixFactorService) GetByID(matrixID uint, id uint) (*models.AssetCriticalityMatrixFactor, error) {
	factor, err := acmfs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return factor, nil
}

func (acmfs *AssetCriticalityMatrixFactorService) Create(matrixID uint, currentUserID uint, factor *models.AssetCriticalityMatrixFactor) (*models.AssetCriticalityMatrixFactor, error) {
	factor.AssetCriticalityMatrixID = matrixID

	// if exists, err := acmfs.repo.ExistsByTitle(factor.TitleEn, factor.TitleAr); err != nil {
	// 	return nil, err
	// } else if exists {
	// 	return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	// }

	createdFactor, err := acmfs.repo.Create(factor)
	if err != nil {
		return nil, err
	}

	go func() {
		err := acmfs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AssetCriticalityMatrixFactor{}.TableName(),
			createdFactor.ID,
			createdFactor,
			currentUserID,
			nil,
		)
		if err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdFactor, nil
}

func (acmfs *AssetCriticalityMatrixFactorService) Update(matrixID uint, id uint, currentUserID uint, updatedFactor *models.AssetCriticalityMatrixFactor) (*models.AssetCriticalityMatrixFactor, error) {
	existingFactor, err := acmfs.GetByID(matrixID, id)
	if err != nil {
		return nil, err
	}

	// if updatedFactor.TitleEn != existingFactor.TitleEn || updatedFactor.TitleAr != existingFactor.TitleAr {
	// 	if exists, err := acmfs.repo.ExistsByTitleExcludingID(id, updatedFactor.TitleEn, updatedFactor.TitleAr); err != nil {
	// 		return nil, err
	// 	} else if exists {
	// 		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	// 	}
	// }

	result, err := acmfs.repo.Update(id, updatedFactor)
	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Weightage"}
	changes, err := utility.CompareFields(existingFactor, updatedFactor, fieldsToCompare)
	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := acmfs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.AssetCriticalityMatrixFactor{}.TableName(),
			result.ID,
			changedDataJSON,
			currentUserID,
			nil,
		)
		if err != nil {
			fmt.Printf("Failed to log update activity: %v\n", err)
		}
	}()
	return result, err
}

func (acmfs *AssetCriticalityMatrixFactorService) Delete(matrixID uint, ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {
		_, err := acmfs.repo.GetByID(id)
		if err != nil {
			fmt.Print(err)
			continue
		}

		err = acmfs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {
		go func() {
			err := acmfs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.AssetCriticalityMatrixFactor{}.TableName(),
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
