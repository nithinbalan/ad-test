package services

import (
	"errors"
	"fmt"
	"time"

	assetDTOs "jethur/grc/asset/dtos"
	"jethur/grc/asset/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type AssetRiskResponseServiceInterface interface {
	GetFullDetailsByID(riskID uint) (*assetDTOs.AssetRiskResponseDTO, error)
	GetByID(riskID uint, id uint) (*models.RiskResponse, error)
	GetByRiskRegisterID(riskID uint, tx ...utility.DatabaseTransaction) (*models.RiskResponse, error)
	Create(riskID uint, currentUserID uint, formDTO *assetDTOs.AssetRiskResponseFormDTO) (uint, error)
	Update(riskID uint, id uint, currentUserID uint, formDTO *assetDTOs.AssetRiskResponseFormDTO) error
	Delete(riskID uint, id uint, currentUserID uint) (uint, error)
}

type AssetRiskResponseService struct {
	repo                 repositories.AssetRiskResponseRepositoryInterface
	assetRiskRegisterSvc AssetRiskRegisterServiceInterface
	activityLogSvc       logServices.ActivityLogServiceInterface
}

func NewAssetRiskResponseService(
	repo repositories.AssetRiskResponseRepositoryInterface,
	assetRiskRegisterSvc AssetRiskRegisterServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
) AssetRiskResponseServiceInterface {
	return &AssetRiskResponseService{
		repo:                 repo,
		assetRiskRegisterSvc: assetRiskRegisterSvc,
		activityLogSvc:       activityLogSvc,
	}
}

func (arrs *AssetRiskResponseService) GetFullDetailsByID(riskID uint) (*assetDTOs.AssetRiskResponseDTO, error) {
	riskResponse, err := arrs.repo.GetFullDetailsByRiskRegisterID(riskID)
	if err != nil {
		return nil, utility.NewResourceNotFoundError(err, riskID)
	}

	return toAssetRiskResponseDTO(riskResponse)
}

func (arrs *AssetRiskResponseService) GetByRiskRegisterID(riskID uint, tx ...utility.DatabaseTransaction) (*models.RiskResponse, error) {
	riskResponse, err := arrs.repo.GetByRiskRegisterID(riskID, tx...)
	if err != nil {
		return nil, utility.NewResourceNotFoundError(err, riskID)
	}

	return riskResponse, nil
}

func (arrs *AssetRiskResponseService) GetByID(riskID uint, id uint) (*models.RiskResponse, error) {
	if _, err := arrs.assetRiskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}

	riskResponse, err := arrs.repo.GetByID(id)
	if err != nil {
		return nil, utility.NewResourceNotFoundError(err, id)
	}

	return riskResponse, nil
}

func (arrs *AssetRiskResponseService) Create(riskID uint, currentUserID uint, formDTO *assetDTOs.AssetRiskResponseFormDTO) (id uint, err error) {
	tx := arrs.repo.BeginTransaction()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil {
			tx.Rollback()
			return
		}

		err = tx.Commit()
	}()

	riskRegister, err := arrs.assetRiskRegisterSvc.GetByID(riskID, tx)
	if err != nil {
		return 0, err
	}

	if _, existingErr := arrs.repo.GetByRiskRegisterID(riskID, tx); existingErr == nil {
		err = utility.NewValidationError(fmt.Errorf("a risk response already exists for this risk"), "unique", "risk_id")
		return 0, err
	} else if !errors.Is(existingErr, gorm.ErrRecordNotFound) {
		err = existingErr
		return 0, err
	}

	riskResponse := &models.RiskResponse{
		RiskRegisterID: riskID,
		CreatedBYID:    &currentUserID,
		CreatedAt:      time.Now(),
	}

	id, err = arrs.repo.Create(riskResponse, tx)
	if err != nil {
		return 0, err
	}

	now := time.Now()
	riskRegister.RiskResponseStrategyID = &formDTO.RiskResponseStrategyID
	riskRegister.UpdatedBYID = &currentUserID
	riskRegister.UpdatedAt = &now

	if err = arrs.assetRiskRegisterSvc.UpdateRiskRegister(riskRegister, tx); err != nil {
		return 0, err
	}

	timelineType, err := arrs.assetRiskRegisterSvc.GetTimelineTypeByType("risk_response", tx)
	if err != nil {
		return 0, err
	}
	timelineStatus, err := arrs.assetRiskRegisterSvc.GetTimelineStatusByType("completed", tx)
	if err != nil {
		return 0, err
	}
	riskTimeline, err := arrs.assetRiskRegisterSvc.GetRiskTimelineByRiskAndType(riskID, timelineType.ID, tx)
	if err != nil {
		return 0, err
	}
	riskTimeline.TimelineStatusID = timelineStatus.ID
	riskTimeline.PerformedByID = &currentUserID
	riskTimeline.PerformedAt = &now

	if err = arrs.assetRiskRegisterSvc.UpdateRiskTimeline(riskTimeline, tx); err != nil {
		return 0, err
	}

	go func() {
		newResponse, logErr := arrs.GetFullDetailsByID(riskID)
		if logErr != nil {
			utility.LogError("Failed to fetch asset risk response", map[string]interface{}{
				"error": logErr.Error(),
			})
			return
		}

		logErr = arrs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.RiskResponse{}.TableName(),
			id,
			newResponse,
			currentUserID,
			"asset",
			"asset_risk",
			"risk",
			nil,
		)

		if logErr != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":   models.CreateAction,
				"recordID": id,
				"error":    logErr.Error(),
			})
		}
	}()

	return id, nil
}

func (arrs *AssetRiskResponseService) Update(riskID uint, id uint, currentUserID uint, formDTO *assetDTOs.AssetRiskResponseFormDTO) (err error) {
	tx := arrs.repo.BeginTransaction()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil {
			tx.Rollback()
			return
		}

		err = tx.Commit()
	}()

	riskRegister, err := arrs.assetRiskRegisterSvc.GetByID(riskID, tx)
	if err != nil {
		return err
	}

	existingResponse, err := arrs.repo.GetByID(id, tx)
	if err != nil {
		return utility.NewResourceNotFoundError(err, id)
	}

	oldResponse, err := arrs.GetFullDetailsByID(riskID)
	if err != nil {
		return err
	}

	now := time.Now()
	existingResponse.UpdatedBYID = &currentUserID
	existingResponse.UpdatedAt = &now

	if err = arrs.repo.Update(existingResponse, tx); err != nil {
		return err
	}

	riskRegister.RiskResponseStrategyID = &formDTO.RiskResponseStrategyID
	riskRegister.UpdatedBYID = &currentUserID
	riskRegister.UpdatedAt = &now

	if err = arrs.assetRiskRegisterSvc.UpdateRiskRegister(riskRegister, tx); err != nil {
		return err
	}

	go func() {
		newResponse, logErr := arrs.GetFullDetailsByID(riskID)
		if logErr != nil {
			utility.LogError("Failed to fetch asset risk response", map[string]interface{}{
				"error": logErr.Error(),
			})
			return
		}

		oldData := map[string]interface{}{}
		if oldResponse.RiskResponseStrategy != nil {
			oldData["risk_response_strategy"] = oldResponse.RiskResponseStrategy.TitleEn
		}
		newData := map[string]interface{}{}
		if newResponse.RiskResponseStrategy != nil {
			newData["risk_response_strategy"] = newResponse.RiskResponseStrategy.TitleEn
		}

		changes := utility.CompareFieldChanges(oldData, newData)
		if len(changes) == 0 {
			return
		}

		logErr = arrs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.RiskResponse{}.TableName(),
			id,
			changes,
			currentUserID,
			"asset",
			"asset_risk",
			"risk",
			nil,
		)

		if logErr != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":   models.UpdateAction,
				"recordID": id,
				"error":    logErr.Error(),
			})
		}
	}()

	return nil
}

func (arrs *AssetRiskResponseService) Delete(riskID uint, id uint, currentUserID uint) (failedID uint, err error) {
	tx := arrs.repo.BeginTransaction()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil {
			tx.Rollback()
			return
		}

		err = tx.Commit()
	}()

	riskRegister, err := arrs.assetRiskRegisterSvc.GetByID(riskID, tx)
	if err != nil {
		return 0, err
	}

	if _, err = arrs.repo.GetByID(id, tx); err != nil {
		return 0, utility.NewResourceNotFoundError(err, id)
	}

	if err = arrs.repo.Delete(id, tx); err != nil {
		return id, err
	}

	now := time.Now()
	riskRegister.RiskResponseStrategyID = nil
	riskRegister.UpdatedBYID = &currentUserID
	riskRegister.UpdatedAt = &now

	if err = arrs.assetRiskRegisterSvc.UpdateRiskRegister(riskRegister, tx); err != nil {
		return id, err
	}

	timelineType, err := arrs.assetRiskRegisterSvc.GetTimelineTypeByType("risk_response", tx)
	if err != nil {
		return id, err
	}
	timelineStatus, err := arrs.assetRiskRegisterSvc.GetTimelineStatusByType("pending", tx)
	if err != nil {
		return id, err
	}
	riskTimeline, err := arrs.assetRiskRegisterSvc.GetRiskTimelineByRiskAndType(riskID, timelineType.ID, tx)
	if err != nil {
		return id, err
	}
	riskTimeline.TimelineStatusID = timelineStatus.ID
	riskTimeline.PerformedByID = nil
	riskTimeline.PerformedAt = nil

	if err = arrs.assetRiskRegisterSvc.UpdateRiskTimeline(riskTimeline, tx); err != nil {
		return id, err
	}

	go func() {
		deletedResponse := map[string]interface{}{
			"id":               id,
			"risk_register_id": riskID,
		}

		logErr := arrs.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.RiskResponse{}.TableName(),
			id,
			deletedResponse,
			currentUserID,
			"asset",
			"asset_risk",
			"risk",
			nil,
		)

		if logErr != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":   models.DeleteAction,
				"recordID": id,
				"error":    logErr.Error(),
			})
		}
	}()

	return 0, nil
}

func toAssetRiskResponseDTO(riskResponse *models.RiskResponse) (*assetDTOs.AssetRiskResponseDTO, error) {
	var dto assetDTOs.AssetRiskResponseDTO

	if err := mapstructure.Decode(riskResponse, &dto); err != nil {
		return nil, err
	}

	dto.CreatedAt = riskResponse.CreatedAt
	dto.UpdatedAt = riskResponse.UpdatedAt

	if riskResponse.RiskRegister != nil && riskResponse.RiskRegister.RiskResponseStrategy != nil {
		strategy := riskResponse.RiskRegister.RiskResponseStrategy
		dto.RiskResponseStrategy = &assetDTOs.AssetRiskResponseStrategyDTO{
			ID:            strategy.ID,
			TitleEn:       strategy.TitleEn,
			TitleAr:       strategy.TitleAr,
			DescriptionEn: strategy.DescriptionEn,
			DescriptionAr: strategy.DescriptionAr,
			Status:        strategy.Status,
		}
	}

	return &dto, nil
}
