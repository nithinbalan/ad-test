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
			Score:                                calculatedScore,
			CreatedByID:                          currentUserID,
		}
		if err := acs.repo.Create(contribution, tx); err != nil {
			tx.Rollback()
			return err
		}
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
		err := acs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AssetCriticalityContribution{}.TableName(),
			assetID,
			assetCriticalityFormDTO,
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

	return nil
}

func (acs *AssetCriticalityService) Delete(assetID uint, userID uint) error {
	tx := acs.repo.BeginTransaction()

	_, err := acs.assetRegisterSvc.GetByID(assetID)
	if err != nil {
		return err
	}

	if err := acs.repo.Delete(assetID, tx); err != nil {
		tx.Rollback()
		return err
	}

	err = acs.assetRegisterSvc.UpdateTotalScore(assetID, 0, tx)
	if err != nil {
		return err
	}

	tx.Commit()
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
