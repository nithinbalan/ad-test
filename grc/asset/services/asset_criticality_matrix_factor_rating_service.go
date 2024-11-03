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
	repo           repositories.AssetCriticalityMatrixFactorRatingRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewAssetCriticalityMatrixFactorRatingService(repo repositories.AssetCriticalityMatrixFactorRatingRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) AssetCriticalityMatrixFactorRatingServiceInterface {
	return &AssetCriticalityMatrixFactorRatingService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (s *AssetCriticalityMatrixFactorRatingService) GetAll(sort, sortBy, keyword, status string) ([]*models.AssetCriticalityMatrixFactorRating, error) {
	return s.repo.GetAll(sort, sortBy, keyword, status)
}

func (s *AssetCriticalityMatrixFactorRatingService) GetByID(id uint) (*models.AssetCriticalityMatrixFactorRating, error) {
	rating, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return rating, nil
}

func (s *AssetCriticalityMatrixFactorRatingService) Create(currentUserID uint, rating *models.AssetCriticalityMatrixFactorRating) (*models.AssetCriticalityMatrixFactorRating, error) {
	createdRating, err := s.repo.Create(rating)
	if err != nil {
		return nil, err
	}

	go func() {
		err := s.activityLogSvc.LogActivity(
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

func (s *AssetCriticalityMatrixFactorRatingService) Update(id, currentUserID uint, updatedRating *models.AssetCriticalityMatrixFactorRating) (*models.AssetCriticalityMatrixFactorRating, error) {
	existingRating, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	result, err := s.repo.Update(id, updatedRating)
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
		err := s.activityLogSvc.LogActivity(
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

func (s *AssetCriticalityMatrixFactorRatingService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
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
