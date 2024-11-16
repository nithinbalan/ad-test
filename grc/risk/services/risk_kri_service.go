package services

import (
	"errors"
	"fmt"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"
	riskRepositories "jethur/grc/risk/repositories"
	"jethur/utility"
	"time"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type RiskKriServiceInterface interface {
	GetAll(riskID uint, params *riskDTOs.RiskKriListQueryParamsDTO) (*riskDTOs.RiskKriListResponseDTO, error)
	GetByID(riskID uint, id uint) (*models.RiskKri, error)
	GetFullDetailsByID(riskID uint, id uint) (*riskDTOs.RiskKriDetailDTO, error)
	Create(riskID uint, currentUserID uint, RiskKriFormDTO *riskDTOs.RiskKriFormDTO) ([]uint, error)
	Update(riskID uint, currentUserID uint, risk *riskDTOs.RiskKriFormDTO) error
	Delete(riskID uint, ids []uint, currentUserID uint) ([]uint, error)
	GetCountByRiskID(riskID uint) (int64, error)
}

// Service and implementations
type RiskKriService struct {
	repo            riskRepositories.RiskKriRepositoryInterface
	activityLogSvc  logServices.ActivityLogServiceInterface
	riskRegisterSvc RiskRegisterServiceInterface
	kriRegisterSvc  KRIRegisterServiceInterface
	riskTimelineSvc RiskTimelineServiceInterface
}

func NewRiskKriService(
	repo riskRepositories.RiskKriRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	riskRegisterSvc RiskRegisterServiceInterface,
	kriRegisterSvc KRIRegisterServiceInterface,
	riskTimelineSvc RiskTimelineServiceInterface,
) RiskKriServiceInterface {
	return &RiskKriService{
		repo:            repo,
		activityLogSvc:  activityLogSvc,
		riskRegisterSvc: riskRegisterSvc,
		kriRegisterSvc:  kriRegisterSvc,
		riskTimelineSvc: riskTimelineSvc,
	}
}

func (rks *RiskKriService) GetAll(riskID uint, params *riskDTOs.RiskKriListQueryParamsDTO) (*riskDTOs.RiskKriListResponseDTO, error) {
	if riskID != 0 {
		if _, err := rks.riskRegisterSvc.GetByID(riskID); err != nil {
			return nil, err
		}
	}
	return rks.repo.GetAll(riskID, params)
}

func (rks RiskKriService) GetByID(riskID uint, id uint) (*models.RiskKri, error) {
	if _, err := rks.riskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}

	riskKri, err := rks.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return riskKri, nil
}

func (rks RiskKriService) GetFullDetailsByID(riskID uint, id uint) (*riskDTOs.RiskKriDetailDTO, error) {
	var riskKriDTO riskDTOs.RiskKriDetailDTO

	if _, err := rks.riskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}

	riskKri, err := rks.repo.GetFullDetailsByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(riskKri, &riskKriDTO)
	if err != nil {
		return nil, err
	}

	return &riskKriDTO, nil
}

func (rks *RiskKriService) Create(riskID uint, currentUserID uint, riskKriFormDTO *riskDTOs.RiskKriFormDTO) ([]uint, error) {
	tx := rks.repo.BeginTransaction()

	if _, err := rks.riskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}

	if err := rks.validateRiskKriFormDTO(riskID, riskKriFormDTO); err != nil {
		return nil, err
	}

	// Create KRIs
	var kris []*models.RiskKri
	for _, kriID := range riskKriFormDTO.KriIDs {
		kri := &models.RiskKri{
			RiskRegisterID: riskID,
			KriRegisterID:  *kriID,
			CreatedByID:    currentUserID,
		}
		kris = append(kris, kri)
	}

	ids, err := rks.repo.Create(kris, tx)

	if err != nil {
		tx.Rollback()
		return nil, err
	}
	//update timeline
	//todo
	kriCount, err := rks.repo.GetCountByRiskID(riskID, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if kriCount == 1 {
		timelineType, err := rks.riskTimelineSvc.GetByType("monitoring")
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		timelineStatus, err := rks.riskTimelineSvc.GetByStatusType("completed")
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		currentTime := time.Now()
		timelineFormDTO := &riskDTOs.RiskTimelineFormDTO{
			PerformedByID:    &currentUserID,
			TimelineStatusID: timelineStatus.ID,
			PerformedAt:      &currentTime,
		}
		if err := rks.riskTimelineSvc.Update(riskID, timelineFormDTO, timelineType.ID, tx); err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	tx.Commit()

	go func() {
		err := rks.activityLogSvc.LogBulkActivities(
			models.CreateAction,
			models.RiskKri{}.TableName(),
			ids,
			kris,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.RiskKri{}.TableName(),
				"recordID":  ids,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return ids, nil
}

func (rks *RiskKriService) Update(riskID uint, currentUserID uint, riskKriFormDTO *riskDTOs.RiskKriFormDTO) error {
	tx := rks.repo.BeginTransaction()
	existingRisk, err := rks.riskRegisterSvc.GetFullDetailsByID(riskID)
	if err != nil {
		return err
	}

	if err := rks.validateRiskKriFormDTO(riskID, riskKriFormDTO); err != nil {
		return err
	}

	existingKRIIDs := make(map[uint]bool)
	for _, kri := range existingRisk.Kris {
		existingKRIIDs[kri.KriID] = true
	}

	newKRIIDs := make(map[uint]bool)
	for _, kri := range riskKriFormDTO.KriIDs {
		newKRIIDs[*kri] = true
	}
	var krisToAdd []*models.RiskKri
	for _, newKRIID := range riskKriFormDTO.KriIDs {
		if _, exists := existingKRIIDs[*newKRIID]; !exists {
			kri := &models.RiskKri{
				RiskRegisterID: riskID,
				KriRegisterID:  *newKRIID,
				CreatedByID:    currentUserID,
			}
			krisToAdd = append(krisToAdd, kri)
		}
	}
	if len(krisToAdd) > 0 {
		_, err := rks.repo.Create(krisToAdd, tx)
		if err != nil {
			return err
		}
	}

	var deleteKRIIDs []*uint

	for _, existingKRI := range existingRisk.Kris {
		if _, exists := newKRIIDs[existingKRI.KriID]; !exists {
			deleteKRIIDs = append(deleteKRIIDs, existingKRI.ID)
		}
	}

	if len(deleteKRIIDs) > 0 {
		if err := rks.repo.DeleteRiskKRIs(deleteKRIIDs); err != nil {
			return err
		}
	}

	tx.Commit()
	go func() {
		krisAdded := []uint{}
		krisRemoved := []uint{}
		for id := range newKRIIDs {
			if !existingKRIIDs[id] {
				krisAdded = append(krisAdded, id)
			}
		}
		for id := range existingKRIIDs {
			if !newKRIIDs[id] {
				krisRemoved = append(krisRemoved, id)
			}
		}
		// Prepare log data
		changedData := make(map[string]interface{})
		associatedItems := make(map[string]interface{})

		if len(krisAdded) > 0 || len(krisRemoved) > 0 {
			associatedItems["kri_ids"] = map[string]interface{}{
				"added":   krisAdded,
				"removed": krisRemoved,
			}
		}
		if len(associatedItems) > 0 {
			changedData["associated_items"] = associatedItems
		}

		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := rks.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.RiskKri{}.TableName(),
				riskID,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to log activity", map[string]interface{}{
					"action":   models.UpdateAction,
					"error":    err.Error(),
					"recordID": riskID,
				})
			}
		}
	}()
	return nil
}

func (rks *RiskKriService) Delete(riskID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := rks.repo.BeginTransaction()

	if _, err := rks.riskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}

	var failedIds []uint
	var successfulIds []uint

	var deletedRiskKri []map[string]interface{}

	for _, id := range ids {

		_, err := rks.repo.GetByID(id)

		if err != nil {
			continue
		}
		//Delete-> update timeline performedAt as null
		//todo
		kriCount, err := rks.repo.GetCountByRiskID(riskID, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		if kriCount == 1 {
			timelineType, err := rks.riskTimelineSvc.GetByType("monitoring")
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			timelineStatus, err := rks.riskTimelineSvc.GetByStatusType("pending")
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			timelineFormDTO := &riskDTOs.RiskTimelineFormDTO{
				PerformedByID:    nil,
				TimelineStatusID: timelineStatus.ID,
				PerformedAt:      nil,
			}
			if err := rks.riskTimelineSvc.Update(riskID, timelineFormDTO, timelineType.ID, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		err = rks.repo.Delete(&id, tx)
		if err != nil {
			tx.Rollback()
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedRiskKri = append(deletedRiskKri, map[string]interface{}{
			"ID": id,
		})
	}
	tx.Commit()

	if len(successfulIds) > 0 {

		go func() {
			err := rks.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.RiskKri{}.TableName(),
				successfulIds,
				deletedRiskKri,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (rks *RiskKriService) validateRiskKriFormDTO(riskID uint, riskKriFormDTO *riskDTOs.RiskKriFormDTO) error {
	kriIds := riskKriFormDTO.KriIDs

	if len(kriIds) == 0 {
		return utility.NewValidationError(fmt.Errorf("no KRI IDs provided"), "empty", "kri")
	}

	kriIDs := make(map[uint]bool)

	for _, kriId := range kriIds {
		if kriIDs[*kriId] {
			return utility.NewValidationError(fmt.Errorf("duplicate kri ID: %d", kriId), "duplicate", "kri")
		}

		if _, err := rks.kriRegisterSvc.GetByID(*kriId); err != nil {
			return utility.NewValidationError(fmt.Errorf("kri does not exist for ID: %d", kriId), "not_exist", "kri")
		}

		// Check if the KRI is already associated with the given risk
		existingKRIs, err := rks.repo.GetKrisByRiskID(riskID)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("failed to fetch existing KRIs for the risk"), "fetch_error", "kri")
		}

		for _, existingKri := range existingKRIs {
			if existingKri.KriRegisterID == *kriId {
				return utility.NewValidationError(fmt.Errorf("KRI ID is already associated with this risk"), "already_associated", "kri")
			}
		}

		kriIDs[*kriId] = true
	}
	return nil
}

func (rks *RiskKriService) GetCountByRiskID(riskID uint) (int64, error) {
	if _, err := rks.riskRegisterSvc.GetByID(riskID); err != nil {
		return 0, err
	}

	return rks.repo.GetCountByRiskID(riskID)
}
