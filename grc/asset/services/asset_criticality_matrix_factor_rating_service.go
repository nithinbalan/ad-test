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

type AssetCriticalityMatrixFactorRatingServiceInterface interface {
	GetAll(matrixID uint, sort, sortBy, keyword, status string) ([]*models.AssetCriticalityMatrixFactorRating, error)
	GetByID(matrixID uint, id uint) (*models.AssetCriticalityMatrixFactorRating, error)
	Create(matrixID uint, currentUserID uint, rating *models.AssetCriticalityMatrixFactorRating) (*models.AssetCriticalityMatrixFactorRating, error)
	Update(matrixID uint, id, currentUserID uint, rating *models.AssetCriticalityMatrixFactorRating) (*models.AssetCriticalityMatrixFactorRating, error)
	Delete(matrixID uint, ids []uint, currentUserID uint) ([]uint, error)
}

type AssetCriticalityMatrixFactorRatingService struct {
	repo                      repositories.AssetCriticalityMatrixFactorRatingRepositoryInterface
	activityLogSvc            services.ActivityLogServiceInterface
	assetCriticalityMatrixSvc AssetCriticalityMatrixServiceInterface
}

func NewAssetCriticalityMatrixFactorRatingService(repo repositories.AssetCriticalityMatrixFactorRatingRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface, assetCriticalityMatrixSvc AssetCriticalityMatrixServiceInterface) AssetCriticalityMatrixFactorRatingServiceInterface {
	return &AssetCriticalityMatrixFactorRatingService{
		repo:                      repo,
		activityLogSvc:            activityLogSvc,
		assetCriticalityMatrixSvc: assetCriticalityMatrixSvc,
	}
}

func (acmfrs *AssetCriticalityMatrixFactorRatingService) GetAll(matrixID uint, sort, sortBy, keyword, status string) ([]*models.AssetCriticalityMatrixFactorRating, error) {
	return acmfrs.repo.GetAll(matrixID, sort, sortBy, keyword, status)
}

func (acmfrs *AssetCriticalityMatrixFactorRatingService) GetByID(matrixID uint, id uint) (*models.AssetCriticalityMatrixFactorRating, error) {
	rating, err := acmfrs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return rating, nil
}

func (acmfrs *AssetCriticalityMatrixFactorRatingService) Create(matrixID uint, currentUserID uint, rating *models.AssetCriticalityMatrixFactorRating) (*models.AssetCriticalityMatrixFactorRating, error) {

	rating.AssetCriticalityMatrixID = matrixID

	// Check if title exists
	if exists, err := acmfrs.repo.ExistsByTitle(rating.TitleEn, rating.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	createdRating, err := acmfrs.repo.Create(rating)
	if err != nil {
		return nil, err
	}

	go func() {
		err := acmfrs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AssetCriticalityMatrixFactorRating{}.TableName(),
			createdRating.ID,
			createdRating,
			currentUserID,
			nil,
		)
		if err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdRating, nil
}

func (acmfrs *AssetCriticalityMatrixFactorRatingService) Update(matrixID uint, id, currentUserID uint, updatedRating *models.AssetCriticalityMatrixFactorRating) (*models.AssetCriticalityMatrixFactorRating, error) {
	existingRating, err := acmfrs.GetByID(matrixID, id)
	if err != nil {
		return nil, err
	}

	// Validate unique title for update
	if updatedRating.TitleEn != existingRating.TitleEn || updatedRating.TitleAr != existingRating.TitleAr {
		if exists, err := acmfrs.repo.ExistsByTitleExcludingID(id, updatedRating.TitleEn, updatedRating.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	result, err := acmfrs.repo.Update(id, updatedRating)
	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Score", "Type", "ColorCode"}
	changes, err := utility.CompareFields(existingRating, updatedRating, fieldsToCompare)
	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := acmfrs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.AssetCriticalityMatrixFactorRating{}.TableName(),
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

func (acmfrs *AssetCriticalityMatrixFactorRatingService) Delete(matrixID uint, ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedFactorRating []map[string]interface{}
	for _, id := range ids {
		factorRating, err := acmfrs.repo.GetByID(id)
		if err != nil {
			fmt.Print(err)
			continue
		}

		err = acmfrs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedFactorRating = append(deletedFactorRating, map[string]interface{}{
			"id":       id,
			"title_en": factorRating.TitleEn,
			"title_ar": factorRating.TitleAr,
		})
	}

	if len(successfulIds) > 0 {
		go func() {
			err := acmfrs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.AssetCriticalityMatrixFactorRating{}.TableName(),
				successfulIds,
				deletedFactorRating,
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
