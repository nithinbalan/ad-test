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

type AssetCriticalityMatrixRatingScaleServiceInterface interface {
	GetAll(matrixID uint, sort, sortBy, keyword, status string) ([]*models.AssetCriticalityMatrixRatingScale, error)
	GetByID(matrixID uint, id uint) (*models.AssetCriticalityMatrixRatingScale, error)
	Create(matrixID uint, currentUserID uint, scale *models.AssetCriticalityMatrixRatingScale) (*models.AssetCriticalityMatrixRatingScale, error)
	Update(matrixID uint, id, currentUserID uint, scale *models.AssetCriticalityMatrixRatingScale) (*models.AssetCriticalityMatrixRatingScale, error)
	Delete(matrixID uint, ids []uint, currentUserID uint) ([]uint, error)

	GetRatingScaleByScore(matrixID uint, score float64) (*models.AssetCriticalityMatrixRatingScale, error)
}

type AssetCriticalityMatrixRatingScaleService struct {
	repo                      repositories.AssetCriticalityMatrixRatingScaleRepositoryInterface
	activityLogSvc            services.ActivityLogServiceInterface
	assetCriticalityMatrixSvc AssetCriticalityMatrixServiceInterface
}

func NewAssetCriticalityMatrixRatingScaleService(repo repositories.AssetCriticalityMatrixRatingScaleRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface, assetCriticalityMatrixSvc AssetCriticalityMatrixServiceInterface) AssetCriticalityMatrixRatingScaleServiceInterface {
	return &AssetCriticalityMatrixRatingScaleService{
		repo:                      repo,
		activityLogSvc:            activityLogSvc,
		assetCriticalityMatrixSvc: assetCriticalityMatrixSvc,
	}
}

func (acmrss *AssetCriticalityMatrixRatingScaleService) GetAll(matrixID uint, sort, sortBy, keyword, status string) ([]*models.AssetCriticalityMatrixRatingScale, error) {
	return acmrss.repo.GetAll(matrixID, sort, sortBy, keyword, status)
}

func (acmrss *AssetCriticalityMatrixRatingScaleService) GetByID(matrixID uint, id uint) (*models.AssetCriticalityMatrixRatingScale, error) {
	scale, err := acmrss.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return scale, nil
}

func (acmrss *AssetCriticalityMatrixRatingScaleService) Create(matrixID uint, currentUserID uint, scale *models.AssetCriticalityMatrixRatingScale) (*models.AssetCriticalityMatrixRatingScale, error) {
	scale.AssetCriticalityMatrixID = matrixID

	exists, err := acmrss.repo.ExistsByTitle(scale.TitleEn, scale.TitleAr)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}
	createdScale, err := acmrss.repo.Create(scale)
	if err != nil {
		return nil, err
	}

	go func() {
		err := acmrss.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AssetCriticalityMatrixRatingScale{}.TableName(),
			createdScale.ID,
			createdScale,
			currentUserID,
			nil,
		)
		if err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdScale, nil
}

func (acmrss *AssetCriticalityMatrixRatingScaleService) Update(matrixID uint, id, currentUserID uint, updatedScale *models.AssetCriticalityMatrixRatingScale) (*models.AssetCriticalityMatrixRatingScale, error) {
	existingScale, err := acmrss.GetByID(matrixID, id)
	if err != nil {
		return nil, err
	}

	if updatedScale.TitleEn != existingScale.TitleEn || updatedScale.TitleAr != existingScale.TitleAr {
		exists, err := acmrss.repo.ExistsByTitleExcludingID(id, updatedScale.TitleEn, updatedScale.TitleAr)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	result, err := acmrss.repo.Update(id, updatedScale)
	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "ScoreFrom", "ScoreTo", "Type", "ColorCode"}
	changes, err := utility.CompareFields(existingScale, updatedScale, fieldsToCompare)
	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := acmrss.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.AssetCriticalityMatrixRatingScale{}.TableName(),
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

func (acmrss *AssetCriticalityMatrixRatingScaleService) Delete(matrixID uint, ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {
		_, err := acmrss.repo.GetByID(id)
		if err != nil {
			fmt.Print(err)
			continue
		}

		err = acmrss.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {
		go func() {
			err := acmrss.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.AssetCriticalityMatrixRatingScale{}.TableName(),
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

func (acmrss *AssetCriticalityMatrixRatingScaleService) GetRatingScaleByScore(matrixID uint, score float64) (*models.AssetCriticalityMatrixRatingScale, error) {
	return acmrss.repo.GetRatingScaleByScore(matrixID, score)
}
