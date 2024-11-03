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
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AssetCriticalityMatrixFactor, error)
	GetByID(id uint) (*models.AssetCriticalityMatrixFactor, error)
	Create(currentUserID uint, factor *models.AssetCriticalityMatrixFactor) (*models.AssetCriticalityMatrixFactor, error)
	Update(id uint, currentUserID uint, factor *models.AssetCriticalityMatrixFactor) (*models.AssetCriticalityMatrixFactor, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

type AssetCriticalityMatrixFactorService struct {
	repo           repositories.AssetCriticalityMatrixFactorRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewAssetCriticalityMatrixFactorService(repo repositories.AssetCriticalityMatrixFactorRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) AssetCriticalityMatrixFactorServiceInterface {
	return &AssetCriticalityMatrixFactorService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (s *AssetCriticalityMatrixFactorService) GetAll(sort, sortBy, keyword, status string) ([]*models.AssetCriticalityMatrixFactor, error) {
	return s.repo.GetAll(sort, sortBy, keyword, status)
}

func (s *AssetCriticalityMatrixFactorService) GetByID(id uint) (*models.AssetCriticalityMatrixFactor, error) {
	factor, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return factor, nil
}

func (s *AssetCriticalityMatrixFactorService) Create(currentUserID uint, factor *models.AssetCriticalityMatrixFactor) (*models.AssetCriticalityMatrixFactor, error) {
	createdFactor, err := s.repo.Create(factor)
	if err != nil {
		return nil, err
	}

	go func() {
		err := s.activityLogSvc.LogActivity(
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

func (s *AssetCriticalityMatrixFactorService) Update(id uint, currentUserID uint, updatedFactor *models.AssetCriticalityMatrixFactor) (*models.AssetCriticalityMatrixFactor, error) {
	existingFactor, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	result, err := s.repo.Update(id, updatedFactor)
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
		err := s.activityLogSvc.LogActivity(
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

func (s *AssetCriticalityMatrixFactorService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
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
