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
	GetAll(sort, sortBy, keyword, status string) ([]*models.AssetCriticalityMatrixRatingScale, error)
	GetByID(id uint) (*models.AssetCriticalityMatrixRatingScale, error)
	Create(currentUserID uint, scale *models.AssetCriticalityMatrixRatingScale) (*models.AssetCriticalityMatrixRatingScale, error)
	Update(id, currentUserID uint, scale *models.AssetCriticalityMatrixRatingScale) (*models.AssetCriticalityMatrixRatingScale, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

type AssetCriticalityMatrixRatingScaleService struct {
	repo           repositories.AssetCriticalityMatrixRatingScaleRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewAssetCriticalityMatrixRatingScaleService(repo repositories.AssetCriticalityMatrixRatingScaleRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) AssetCriticalityMatrixRatingScaleServiceInterface {
	return &AssetCriticalityMatrixRatingScaleService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (s *AssetCriticalityMatrixRatingScaleService) GetAll(sort, sortBy, keyword, status string) ([]*models.AssetCriticalityMatrixRatingScale, error) {
	return s.repo.GetAll(sort, sortBy, keyword, status)
}

func (s *AssetCriticalityMatrixRatingScaleService) GetByID(id uint) (*models.AssetCriticalityMatrixRatingScale, error) {
	scale, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return scale, nil
}

func (s *AssetCriticalityMatrixRatingScaleService) Create(currentUserID uint, scale *models.AssetCriticalityMatrixRatingScale) (*models.AssetCriticalityMatrixRatingScale, error) {
	createdScale, err := s.repo.Create(scale)
	if err != nil {
		return nil, err
	}

	go func() {
		err := s.activityLogSvc.LogActivity(
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

func (s *AssetCriticalityMatrixRatingScaleService) Update(id, currentUserID uint, updatedScale *models.AssetCriticalityMatrixRatingScale) (*models.AssetCriticalityMatrixRatingScale, error) {
	existingScale, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	result, err := s.repo.Update(id, updatedScale)
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
		err := s.activityLogSvc.LogActivity(
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

func (s *AssetCriticalityMatrixRatingScaleService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {
		_, err := s.repo.GetByID(id)
		if err != nil {
			fmt.Print(err)
			continue
		}

		err = s.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {
		go func() {
			err := s.activityLogSvc.LogBulkActivities(
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
