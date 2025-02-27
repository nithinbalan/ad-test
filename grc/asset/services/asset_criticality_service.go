package services

import (
	"errors"
	"fmt"
	assetDTOs "jethur/grc/asset/dtos"
	"jethur/grc/asset/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AssetCriticalityServiceInterface interface {
	GetAll(assetID uint, params *assetDTOs.AssetCriticalityContributionListQueryParamsDTO) ([]*assetDTOs.AssetCriticalityContributionListDTO, error)
	Create(assetID uint, contributions []assetDTOs.AssetCriticalityFormDTO, currentUserID uint) error
	Update(assetID uint, id uint, contribution assetDTOs.AssetCriticalityFormDTO, currentUserID uint) error
	Delete(id uint, userID uint) error
}

type AssetCriticalityService struct {
	repo                                  repositories.AssetCriticalityRepositoryInterface
	activityLogSvc                        logServices.ActivityLogServiceInterface
	assetRegisterSvc                      AssetRegisterServiceInterface
	assetCriticalityMatrixFactorSvc       AssetCriticalityMatrixFactorServiceInterface
	assetCriticalityMatrixFactorRatingSvc AssetCriticalityMatrixFactorRatingServiceInterface
}

func NewAssetCriticalityService(repo repositories.AssetCriticalityRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	assetRegisterSvc AssetRegisterServiceInterface,
	assetCriticalityMatrixFactorSvc AssetCriticalityMatrixFactorServiceInterface,
	assetCriticalityMatrixFactorRatingSvc AssetCriticalityMatrixFactorRatingServiceInterface) AssetCriticalityServiceInterface {
	return &AssetCriticalityService{
		repo:                                  repo,
		activityLogSvc:                        activityLogSvc,
		assetRegisterSvc:                      assetRegisterSvc,
		assetCriticalityMatrixFactorSvc:       assetCriticalityMatrixFactorSvc,
		assetCriticalityMatrixFactorRatingSvc: assetCriticalityMatrixFactorRatingSvc,
	}
}

func (acs *AssetCriticalityService) GetAll(assetID uint, params *assetDTOs.AssetCriticalityContributionListQueryParamsDTO) ([]*assetDTOs.AssetCriticalityContributionListDTO, error) {
	if _, err := acs.assetRegisterSvc.GetByID(assetID); err != nil {
		return nil, err
	}
	return acs.repo.GetAll(assetID, params)
}

func (acs *AssetCriticalityService) GetByID(id uint) (*models.AssetCriticalityContribution, error) {
	asset, err := acs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return asset, nil
}

func (acs *AssetCriticalityService) Create(assetID uint, assetCriticalityFormDTO []assetDTOs.AssetCriticalityFormDTO, currentUserID uint) error {
	tx := acs.repo.BeginTransaction()
	var createdIds []uint
	var addedAssetCriticality []map[string]interface{}
	asset, err := acs.assetRegisterSvc.GetByID(assetID)
	if err != nil {
		return err
	}

	if err := acs.validateAssetCriticalityFormDTO(asset.AssetCriticalityMatrixID, assetCriticalityFormDTO); err != nil {
		return err
	}

	// Delete all existing contributions for this asset
	if err := acs.repo.Delete(assetID, tx); err != nil {
		tx.Rollback()
		return err
	}

	for _, dto := range assetCriticalityFormDTO {
		factor, err := acs.assetCriticalityMatrixFactorSvc.GetByID(asset.AssetCriticalityMatrixID, dto.AssetCriticalityMatrixFactorID)
		if err != nil {
			return err
		}

		rating, err := acs.assetCriticalityMatrixFactorRatingSvc.GetByID(asset.AssetCriticalityMatrixID, dto.AssetCriticalityMatrixFactorRatingID)
		if err != nil {
			return err
		}

		// Calculate the score as weightage * score
		calculatedScore := factor.Weightage * float64(rating.Score)

		// Create the AssetCriticalityContribution with the calculated score
		contribution := &models.AssetCriticalityContribution{
			AssetRegisterID:                      assetID,
			AssetCriticalityMatrixFactorID:       dto.AssetCriticalityMatrixFactorID,
			AssetCriticalityMatrixFactorRatingID: dto.AssetCriticalityMatrixFactorRatingID,
			Comments:                             dto.Comments,
			Score:                                calculatedScore,
			CreatedByID:                          currentUserID,
		}
		id, err := acs.repo.Create(contribution, tx)

		if err != nil {
			tx.Rollback()
			return err
		}
		addedAssetCriticality = append(addedAssetCriticality, map[string]interface{}{
			"id":           id,
			"contribution": contribution,
		})

		createdIds = append(createdIds, id)
	}

	//update total score in asset register
	totalScore, err := acs.repo.CalculateTotalScore(assetID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := acs.assetRegisterSvc.UpdateTotalScore(assetID, totalScore, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		err := acs.activityLogSvc.LogBulkActivities(
			models.CreateAction,
			models.AssetCriticalityContribution{}.TableName(),
			createdIds,
			addedAssetCriticality,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AssetCriticalityContribution{}.TableName(),
				"recordID":  assetID,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log control activity", fields)
		}
	}()

	return nil
}

func (acs *AssetCriticalityService) Update(assetID uint, id uint, dto assetDTOs.AssetCriticalityFormDTO, userID uint) error {
	tx := acs.repo.BeginTransaction()

	asset, err := acs.assetRegisterSvc.GetByID(assetID)
	if err != nil {
		return err
	}

	oldContribution, err := acs.GetByID(id)
	if err != nil {
		return err
	}

	contribution, err := acs.GetByID(id)
	if err != nil {
		return err
	}

	factor, err := acs.assetCriticalityMatrixFactorSvc.GetByID(asset.AssetCriticalityMatrixID, dto.AssetCriticalityMatrixFactorID)
	if err != nil {
		return err
	}

	rating, err := acs.assetCriticalityMatrixFactorRatingSvc.GetByID(asset.AssetCriticalityMatrixID, dto.AssetCriticalityMatrixFactorRatingID)
	if err != nil {
		return err
	}

	// Calculate the score as weightage * score
	calculatedScore := factor.Weightage * float64(rating.Score)

	contribution.AssetCriticalityMatrixFactorID = dto.AssetCriticalityMatrixFactorID
	contribution.AssetCriticalityMatrixFactorRatingID = dto.AssetCriticalityMatrixFactorRatingID
	contribution.Score = calculatedScore
	contribution.Comments = dto.Comments
	contribution.UpdatedByID = &userID

	err = acs.repo.Update(contribution, tx)
	if err != nil {
		return err
	}

	totalScore, err := acs.repo.CalculateTotalScore(assetID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := acs.assetRegisterSvc.UpdateTotalScore(assetID, totalScore, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	go func() {
		existingContribution, err := acs.GetByID(id)
		if err != nil {
			return
		}
		// Prepare log data
		changedData := make(map[string]interface{})
		//log changes
		oldData := map[string]interface{}{
			"asset_criticality_matrix_factor":        oldContribution.AssetCriticalityMatrixFactor.TitleEn,
			"asset_criticality_matrix_factor_rating": oldContribution.AssetCriticalityMatrixFactorRating.TitleEn,
			"score":                                  oldContribution.Score,
			"comments":                               oldContribution.Comments,
		}

		newData := map[string]interface{}{
			"asset_criticality_matrix_factor":        existingContribution.AssetCriticalityMatrixFactor.TitleEn,
			"asset_criticality_matrix_factor_rating": existingContribution.AssetCriticalityMatrixFactorRating.TitleEn,
			"score":                                  existingContribution.Score,
			"comments":                               existingContribution.Comments,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := acs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.AssetCriticalityContribution{}.TableName(),
				id,
				changedData,
				userID,
				nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.AssetCriticalityContribution{}.TableName(),
					"recordID":  id,
					"userID":    userID,
					"error":     err.Error(),
				}
				fmt.Printf("Failed to log activity: %v\n", fields)
			}
		}
	}()
	return nil
}

func (acs *AssetCriticalityService) Delete(assetID uint, userID uint) error {
	tx := acs.repo.BeginTransaction()
	var createdIds []uint
	var deletedContribution []map[string]interface{}
	_, err := acs.assetRegisterSvc.GetByID(assetID)
	if err != nil {
		return err
	}

	contributions, err := acs.repo.GetAllContribution(assetID)
	if err != nil {
		return nil
	}

	if err := acs.repo.Delete(assetID, tx); err != nil {
		tx.Rollback()
		return err
	}

	err = acs.assetRegisterSvc.UpdateTotalScore(assetID, 0, tx)
	if err != nil {
		return err
	}
	for _, contribution := range contributions {
		createdIds = append(createdIds, contribution.ID)

		deletedContribution = append(deletedContribution, map[string]interface{}{
			"id":                                     contribution.ID,
			"asset_id":                               assetID,
			"asset_criticality_matrix_factor":        contribution.AssetCriticalityMatrixFactor.TitleEn,
			"asset_criticality_matrix_factor_rating": contribution.AssetCriticalityMatrixFactorRating.TitleEn,
		})
	}
	fmt.Println(deletedContribution)
	tx.Commit()
	go func() {
		err := acs.activityLogSvc.LogBulkActivities(
			models.DeleteAction,
			models.AssetCriticalityContribution{}.TableName(),
			createdIds,
			deletedContribution,
			userID,
			nil,
		)
		if err != nil {
			// Handle error if logging fails
			fmt.Printf("Failed to log delete activity: %v\n", err)
		}
	}()
	return nil
}

func (acs *AssetCriticalityService) validateAssetCriticalityFormDTO(matrixID uint, assetCriticalityFormDTO []assetDTOs.AssetCriticalityFormDTO) error {
	for _, dto := range assetCriticalityFormDTO {
		if _, err := acs.assetCriticalityMatrixFactorSvc.GetByID(matrixID, dto.AssetCriticalityMatrixFactorID); err != nil {
			return utility.NewValidationError(fmt.Errorf("asset criticality matrix factor does not exist"), "not_exist", "asset criticality matrix factor")
		}

		if _, err := acs.assetCriticalityMatrixFactorRatingSvc.GetByID(matrixID, dto.AssetCriticalityMatrixFactorRatingID); err != nil {
			return utility.NewValidationError(fmt.Errorf("asset criticality matrix factor rating does not exist"), "not_exist", "asset criticality matrix factor rating")
		}
	}
	return nil
}
