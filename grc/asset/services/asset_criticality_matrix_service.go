package services

import (
	"errors"
	"fmt"
	"jethur/grc/asset/dtos"
	"jethur/grc/asset/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type AssetCriticalityMatrixServiceInterface interface {
	GetAll(params *dtos.AssetCriticalityMatrixListQueryParamsDTO) (*dtos.AssetCriticalityMatrixListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*dtos.AssetCriticalityMatrixDetailDTO, error)
	Create(matrixFormDTO *dtos.AssetCriticalityMatrixFormDTO, currentUserID uint) (uint, error)
	Update(id uint, matrixFormDTO *dtos.AssetCriticalityMatrixFormDTO, currentUserID uint) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	GetCriticalityMatrixByCategoryID(assetCategoryID uint) (uint, error)
}

type AssetCriticalityMatrixService struct {
	repo           repositories.AssetCriticalityMatrixRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
}

func NewAssetCriticalityMatrixService(
	repo repositories.AssetCriticalityMatrixRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
) AssetCriticalityMatrixServiceInterface {
	return &AssetCriticalityMatrixService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (acms *AssetCriticalityMatrixService) GetAll(params *dtos.AssetCriticalityMatrixListQueryParamsDTO) (*dtos.AssetCriticalityMatrixListResponseDTO, error) {
	return acms.repo.GetAll(params)
}

func (acms *AssetCriticalityMatrixService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*dtos.AssetCriticalityMatrixDetailDTO, error) {
	matrix, err := acms.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	var matrixDetailDTO dtos.AssetCriticalityMatrixDetailDTO
	if err := mapstructure.Decode(matrix, &matrixDetailDTO); err != nil {
		return nil, err
	}

	return &matrixDetailDTO, nil
}

func (acms *AssetCriticalityMatrixService) Create(matrixFormDTO *dtos.AssetCriticalityMatrixFormDTO, currentUserID uint) (uint, error) {
	tx := acms.repo.BeginTransaction()

	matrix := &models.AssetCriticalityMatrix{
		TitleEn: matrixFormDTO.TitleEn,
		TitleAr: matrixFormDTO.TitleAr,
		Status:  matrixFormDTO.Status,
	}

	matrixID, err := acms.repo.Create(matrix, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var categories []*models.AssetCriticalityMatrixAssetCategory
	for _, categoryID := range matrixFormDTO.AssetCategoryIDs {
		category := &models.AssetCriticalityMatrixAssetCategory{
			AssetCriticalityMatrixID: matrixID,
			AssetCategoryID:          categoryID,
			CreatedBYID:              currentUserID,
		}
		categories = append(categories, category)
	}

	if len(categories) > 0 {
		if err := acms.repo.CreateCategories(categories, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()

	go func() {
		err := acms.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AssetCriticalityMatrix{}.TableName(),
			matrixID,
			matrixFormDTO,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AssetRegister{}.TableName(),
				"recordID":  matrixID,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log control activity", fields)
		}
	}()

	return matrixID, nil
}

func (acms *AssetCriticalityMatrixService) Update(id uint, matrixFormDTO *dtos.AssetCriticalityMatrixFormDTO, currentUserID uint) error {
	tx := acms.repo.BeginTransaction()

	matrix, err := acms.repo.GetByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	matrix.TitleEn = matrixFormDTO.TitleEn
	matrix.TitleAr = matrixFormDTO.TitleAr
	matrix.Status = matrixFormDTO.Status

	if err := acms.repo.Update(matrix, tx); err != nil {
		tx.Rollback()
		return err
	}

	existingMatrix, err := acms.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingCategoryIDs := make(map[uint]bool)
	for _, category := range existingMatrix.AssetCategories {
		existingCategoryIDs[category.AssetCategoryID] = true
	}

	newCategoryIDs := make(map[uint]bool)
	for _, categoryID := range matrixFormDTO.AssetCategoryIDs {
		newCategoryIDs[categoryID] = true
	}

	var categoriesToDelete []uint
	for categoryID := range existingCategoryIDs {
		if !newCategoryIDs[categoryID] {
			categoriesToDelete = append(categoriesToDelete, categoryID)
		}
	}

	if len(categoriesToDelete) > 0 {
		if err := acms.repo.DeleteCategories(id, categoriesToDelete, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	var categoriesToAdd []*models.AssetCriticalityMatrixAssetCategory
	for categoryID := range newCategoryIDs {
		if !existingCategoryIDs[categoryID] {
			categoriesToAdd = append(categoriesToAdd, &models.AssetCriticalityMatrixAssetCategory{
				AssetCriticalityMatrixID: id,
				AssetCategoryID:          categoryID,
				CreatedBYID:              currentUserID,
			})
		}
	}

	if len(categoriesToAdd) > 0 {
		if err := acms.repo.CreateCategories(categoriesToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

func (acms *AssetCriticalityMatrixService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint

	tx := acms.repo.BeginTransaction()

	for _, id := range ids {
		// Delete associated categories within the transaction
		if err := acms.repo.DeleteCategoriesByMatrixID(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		// Delete the matrix itself within the transaction
		err := acms.repo.Delete(id, tx)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(failedIds) > 0 {
		tx.Rollback()
	} else {
		tx.Commit()
	}

	// Log successful deletions asynchronously
	if len(successfulIds) > 0 {
		go func() {
			err := acms.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.AssetCriticalityMatrix{}.TableName(),
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

func (acms *AssetCriticalityMatrixService) GetCriticalityMatrixByCategoryID(assetCategoryID uint) (uint, error) {
	return acms.repo.GetCriticalityMatrixByCategoryID(assetCategoryID)
}
