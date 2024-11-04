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
	GetAll(sort, sortBy, keyword, status string) ([]*models.AssetCriticalityMatrixFactorRating, error)
	GetByID(id uint) (*models.AssetCriticalityMatrixFactorRating, error)
	Create(currentUserID uint, rating *models.AssetCriticalityMatrixFactorRating) (*models.AssetCriticalityMatrixFactorRating, error)
	Update(id, currentUserID uint, rating *models.AssetCriticalityMatrixFactorRating) (*models.AssetCriticalityMatrixFactorRating, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
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

func (acmfrs *AssetCriticalityMatrixFactorRatingService) GetAll(sort, sortBy, keyword, status string) ([]*models.AssetCriticalityMatrixFactorRating, error) {
	return acmfrs.repo.GetAll(sort, sortBy, keyword, status)
}

func (acmfrs *AssetCriticalityMatrixFactorRatingService) GetByID(id uint) (*models.AssetCriticalityMatrixFactorRating, error) {
	rating, err := acmfrs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return rating, nil
}

func (acmfrs *AssetCriticalityMatrixFactorRatingService) Create(currentUserID uint, rating *models.AssetCriticalityMatrixFactorRating) (*models.AssetCriticalityMatrixFactorRating, error) {
	matrixID, err := acmfrs.assetCriticalityMatrixSvc.GetFirstMatrixID()
	if err != nil {
		return nil, err
	}
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

func (acmfrs *AssetCriticalityMatrixFactorRatingService) Update(id, currentUserID uint, updatedRating *models.AssetCriticalityMatrixFactorRating) (*models.AssetCriticalityMatrixFactorRating, error) {
	existingRating, err := acmfrs.GetByID(id)
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

func (acmfrs *AssetCriticalityMatrixFactorRatingService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {
		_, err := acmfrs.repo.GetByID(id)
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
	}

	if len(successfulIds) > 0 {
		go func() {
			err := acmfrs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.AssetCriticalityMatrixFactorRating{}.TableName(),
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
