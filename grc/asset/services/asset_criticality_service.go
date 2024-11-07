package services

import (
	"jethur/grc/asset/dtos"
	"jethur/grc/asset/models"
	"jethur/grc/asset/repositories"
	activityLogServices "jethur/grc/logs/services"
	"jethur/utility"
)

type AssetCriticalityServiceInterface interface {
	CreateAssetCriticality(assetID uint, currentUserID uint, assetCriticalityFormDTO *dtos.AssetCriticalityFormDTO) (uint, error)
}

type AssetCriticalityService struct {
	repo           repositories.AssetCriticalityRepositoryInterface
	activityLogSvc activityLogServices.ActivityLogServiceInterface
}

func NewAssetCriticalityService(repo repositories.AssetCriticalityRepositoryInterface, contributionRepo repositories.AssetCriticalityRepositoryInterface, activityLogSvc activityLogServices.ActivityLogServiceInterface) *AssetCriticalityService {
	return &AssetCriticalityService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (acs *AssetCriticalityService) CreateAssetCriticality(assetID uint, currentUserID uint, assetCriticalityFormDTO *dtos.AssetCriticalityFormDTO) (uint, error) {
	// Create the main AssetCriticality entry
	assetCriticality := &models.AssetCriticality{
		AssetID: assetID,
	}

	id, err := acs.repo.Create(assetCriticality)
	if err != nil {
		return 0, err
	}

	// Create entries in asset_criticality_contributions for each factor
	for _, contribution := range assetCriticalityFormDTO.s {
		contributionEntry := &models.AssetCriticality{
			AssetCriticalityID:                   id,
			AssetCriticalityMatrixFactorID:       contribution.AssetCriticalityMatrixFactorID,
			AssetCriticalityMatrixFactorRatingID: contribution.AssetCriticalityMatrixFactorRatingID,
			Score:                                contribution.Score,
			CreatedByID:                          currentUserID,
		}
		if _, err := acs.contributionRepo.Create(contributionEntry); err != nil {
			return 0, err
		}
	}

	// Log activity
	go func() {
		err := acs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AssetCriticality{}.TableName(),
			id,
			assetCriticalityFormDTO,
			currentUserID,
			nil,
		)
		if err != nil {
			utility.LogError("Failed to log asset criticality creation", map[string]interface{}{
				"assetID": assetID,
				"userID":  currentUserID,
				"error":   err.Error(),
			})
		}
	}()

	return id, nil
}
