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
	GetAll() ([]*dtos.AssetCriticalityMatrixResponseDTO, error)
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

func (acms *AssetCriticalityMatrixService) GetAll() ([]*dtos.AssetCriticalityMatrixResponseDTO, error) {
	matrices, err := acms.repo.GetAll()
	if err != nil {
		return nil, err
	}

	// Transform models to DTOs using mapstructure
	var response []*dtos.AssetCriticalityMatrixResponseDTO
	for _, matrix := range matrices {
		var dto dtos.AssetCriticalityMatrixResponseDTO
		if err := mapstructure.Decode(matrix, &dto); err != nil {
			return nil, fmt.Errorf("failed to map matrix to DTO: %w", err)
		}
		response = append(response, &dto)
	}

	return response, nil
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

	if err := acms.validateAssetCriticalityMatrixCategories(matrixFormDTO.AssetCategoryIDs, 0, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

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
		existingCriticalityMatrix, err := acms.GetFullDetailsByID(matrixID)
		if err != nil {
			utility.LogError("Failed to fetch criticality matrix", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = acms.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AssetCriticalityMatrix{}.TableName(),
			matrixID,
			existingCriticalityMatrix,
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

	oldCriticalityMatrix, err := acms.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := acms.validateAssetCriticalityMatrixCategories(matrixFormDTO.AssetCategoryIDs, id, tx); err != nil {
		tx.Rollback()
		return err
	}

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
	go func() {
		existingCriticalityMatrix, err := acms.GetFullDetailsByID(id)
		if err != nil {
			return
		}
		// Determine control and document changes
		categoriesAdded := []uint{}
		categoriesRemoved := []uint{}
		for id := range newCategoryIDs {
			if !existingCategoryIDs[id] {
				categoriesAdded = append(categoriesAdded, id)
			}
		}
		for id := range existingCategoryIDs {
			if !newCategoryIDs[id] {
				categoriesRemoved = append(categoriesRemoved, id)
			}
		}

		// Prepare log data
		changedData := make(map[string]interface{})
		associatedItems := make(map[string]interface{})

		if len(categoriesAdded) > 0 || len(categoriesRemoved) > 0 {
			associatedItems["control_ids"] = map[string]interface{}{
				"added":   categoriesAdded,
				"removed": categoriesRemoved,
			}
		}

		if len(associatedItems) > 0 {
			changedData["associated_items"] = associatedItems
		}
		//log changes
		oldData := map[string]interface{}{
			"title_en": oldCriticalityMatrix.TitleEn,
			"title_ar": oldCriticalityMatrix.TitleAr,
			"status":   oldCriticalityMatrix.Status,
		}

		newData := map[string]interface{}{
			"title_en": existingCriticalityMatrix.TitleEn,
			"title_ar": existingCriticalityMatrix.TitleAr,
			"status":   existingCriticalityMatrix.Status,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := acms.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.AssetCriticalityMatrix{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.AssetCriticalityMatrix{}.TableName(),
					"recordID":  id,
					"userID":    currentUserID,
					"error":     err.Error(),
				}

				utility.LogError("Failed to log activity", fields)
			}
		}
	}()
	return nil
}

func (acms *AssetCriticalityMatrixService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedCriticalityMatrix []map[string]interface{}

	tx := acms.repo.BeginTransaction()

	for _, id := range ids {

		criticalityMatrix, err := acms.repo.GetByID(id)
		if err != nil {
			fmt.Print(err)
			continue
		}

		// Delete associated categories within the transaction
		if err := acms.repo.DeleteCategoriesByMatrixID(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		// Delete the matrix itself within the transaction
		err = acms.repo.Delete(id, tx)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedCriticalityMatrix = append(deletedCriticalityMatrix, map[string]interface{}{
			"id":       id,
			"title_en": criticalityMatrix.TitleEn,
			"title_ar": criticalityMatrix.TitleAr,
		})
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
				deletedCriticalityMatrix,
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

func (acms *AssetCriticalityMatrixService) validateAssetCriticalityMatrixCategories(categoryIDs []uint, excludeMatrixID uint, tx ...utility.DatabaseTransaction) error {
	for _, categoryID := range categoryIDs {
		exists, err := acms.repo.ExistsByCategoryID(categoryID, excludeMatrixID, tx...)
		if err != nil {
			return fmt.Errorf("failed to validate category ID %d: %w", categoryID, err)
		}

		// If a matrix exists for the category (excluding the current matrix being updated)
		if exists {
			return utility.NewValidationError(fmt.Errorf("a criticality matrix already exists for category"), "matrix_exists", "Asset Category")
		}
	}
	return nil
}

func (acms *AssetCriticalityMatrixService) GetCriticalityMatrixByCategoryID(assetCategoryID uint) (uint, error) {
	return acms.repo.GetCriticalityMatrixByCategoryID(assetCategoryID)
}
