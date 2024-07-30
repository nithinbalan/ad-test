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

	"gorm.io/gorm"
)

type RiskRootCauseAnalysisServiceInterface interface {
	GetAll(riskID uint, params *riskDTOs.RiskRootCauseAnalysisListQueryParamsDTO) ([]*riskDTOs.RiskRootCauseAnalysisListDTO, error)
	GetByID(riskID uint, id uint) (*models.RiskRootCauseAnalysis, error)
	Create(riskID uint, currentUserID uint, rootCauseAnalysisFormDTO *riskDTOs.RiskRootCauseAnalysisFormDTO) (uint, error)
	Update(riskID uint, id uint, currentUserID uint, rootCauseAnalysisFormDTO *riskDTOs.RiskRootCauseAnalysisFormDTO) error
	Delete(riskID uint, id uint, currentUserID uint) (uint, error)
	GetCountByRiskID(riskID uint) (int64, error)
}

// Service and implementations
type RiskRootCauseAnalysisService struct {
	repo            riskRepositories.RiskRootCauseAnalysisRepositoryInterface
	activityLogSvc  logServices.ActivityLogServiceInterface
	riskRegisterSvc RiskRegisterServiceInterface
	riskTimelineSvc RiskTimelineServiceInterface
}

func NewRiskRootCauseAnalysisService(
	repo riskRepositories.RiskRootCauseAnalysisRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	riskRegisterSvc RiskRegisterServiceInterface,
	riskTimelineSvc RiskTimelineServiceInterface,
) RiskRootCauseAnalysisServiceInterface {
	return &RiskRootCauseAnalysisService{
		repo:            repo,
		activityLogSvc:  activityLogSvc,
		riskRegisterSvc: riskRegisterSvc,
		riskTimelineSvc: riskTimelineSvc,
	}
}

func (rrcas *RiskRootCauseAnalysisService) GetAll(riskID uint, params *riskDTOs.RiskRootCauseAnalysisListQueryParamsDTO) ([]*riskDTOs.RiskRootCauseAnalysisListDTO, error) {
	if _, err := rrcas.riskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}
	return rrcas.repo.GetAll(riskID, params)
}

func (rrcas *RiskRootCauseAnalysisService) GetByID(riskID uint, id uint) (*models.RiskRootCauseAnalysis, error) {
	if _, err := rrcas.riskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}

	riskRootCauseAnalysis, err := rrcas.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return riskRootCauseAnalysis, nil
}

func (rrcas *RiskRootCauseAnalysisService) Create(riskID uint, currentUserID uint, rootCauseAnalysisFormDTO *riskDTOs.RiskRootCauseAnalysisFormDTO) (uint, error) {
	tx := rrcas.repo.BeginTransaction()

	if _, err := rrcas.riskRegisterSvc.GetByID(riskID); err != nil {
		return 0, err
	}

	// Increment the "Why" field by 1
	latestWhy, err := rrcas.repo.GetLatestWhyByRiskID(riskID)
	if err != nil {
		return 0, err
	}
	why := latestWhy + 1

	riskRootCauseAnalysis := &models.RiskRootCauseAnalysis{
		RiskRegisterID: riskID,
		Why:            why,
		Question:       rootCauseAnalysisFormDTO.Question,
		Answer:         rootCauseAnalysisFormDTO.Answer,
		CreatedByID:    currentUserID,
	}

	id, err := rrcas.repo.Create(riskRootCauseAnalysis)
	if err != nil {
		return 0, err
	}
	//Update timeline
	//todo
	rcaCount, err := rrcas.repo.GetCountByRiskID(riskID)
	if err != nil {
		return 0, err
	}
	if rcaCount == 1 {
		timelineType, err := rrcas.riskTimelineSvc.GetByType("rca")
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		timelineStatus, err := rrcas.riskTimelineSvc.GetByStatusType("completed")
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		currentTime := time.Now()
		timelineFormDTO := &riskDTOs.RiskTimelineFormDTO{
			PerformedByID:    &currentUserID,
			TimelineStatusID: timelineStatus.ID,
			PerformedAt:      &currentTime,
		}
		if err := rrcas.riskTimelineSvc.Update(riskID, timelineFormDTO, timelineType.ID, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	tx.Commit()

	go func() {
		err := rrcas.activityLogSvc.LogActivity(
			models.CreateAction,
			models.RiskRootCauseAnalysis{}.TableName(),
			id,
			rootCauseAnalysisFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.RiskRootCauseAnalysis{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (rrcas *RiskRootCauseAnalysisService) Update(riskID uint, id uint, currentUserID uint, rootCauseAnalysisFormDTO *riskDTOs.RiskRootCauseAnalysisFormDTO) error {

	if _, err := rrcas.riskRegisterSvc.GetByID(riskID); err != nil {
		return err
	}

	existingRiskRootCauseAnalysis, err := rrcas.GetByID(riskID, id)

	if err != nil {
		return err
	}

	// Update Risk RootCauseAnalysis
	existingRiskRootCauseAnalysis.Question = rootCauseAnalysisFormDTO.Question
	existingRiskRootCauseAnalysis.Answer = rootCauseAnalysisFormDTO.Answer
	existingRiskRootCauseAnalysis.UpdatedByID = &currentUserID

	err = rrcas.repo.Update(existingRiskRootCauseAnalysis)

	if err != nil {
		return err
	}

	return nil
}

func (rrcas *RiskRootCauseAnalysisService) Delete(riskID uint, id uint, currentUserID uint) (uint, error) {
	tx := rrcas.repo.BeginTransaction()

	if _, err := rrcas.riskRegisterSvc.GetByID(riskID); err != nil {
		return 0, err
	}

	var failedId uint

	entryToDelete, err := rrcas.GetByID(riskID, id)
	if err != nil {
		return 0, err
	}
	//Delete-> update timeline performedAt as null
	//todo
	rcaCount, err := rrcas.repo.GetCountByRiskID(riskID)
	if err != nil {
		return 0, err
	}
	if rcaCount == 1 {
		timelineType, err := rrcas.riskTimelineSvc.GetByType("rca")
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		timelineStatus, err := rrcas.riskTimelineSvc.GetByStatusType("pending")
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		timelineFormDTO := &riskDTOs.RiskTimelineFormDTO{
			PerformedByID:    nil,
			TimelineStatusID: timelineStatus.ID,
			PerformedAt:      nil,
		}
		if err := rrcas.riskTimelineSvc.Update(riskID, timelineFormDTO, timelineType.ID, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	// Delete the entry
	err = rrcas.repo.Delete(id, tx)
	if err != nil {
		tx.Rollback()
		failedId = id
		return failedId, err
	}

	// Reassign "Why" values for subsequent entries
	err = rrcas.repo.UpdateWhyValuesAfterDeletion(riskID, entryToDelete.Why, tx)
	if err != nil {
		tx.Rollback()
		fmt.Printf("Failed to update 'Why' values after deletion: %v\n", err)
		return 0, err
	}

	tx.Commit()

	// Log successful delete activity
	go func() {
		err := rrcas.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.RiskRootCauseAnalysis{}.TableName(),
			id,
			nil,
			currentUserID,
			nil,
		)
		if err != nil {
			// Handle error if logging fails
			fmt.Printf("Failed to log delete activity: %v\n", err)
		}
	}()

	return failedId, nil
}

func (rrcas *RiskRootCauseAnalysisService) GetCountByRiskID(riskID uint) (int64, error) {
	if _, err := rrcas.riskRegisterSvc.GetByID(riskID); err != nil {
		return 0, err
	}

	return rrcas.repo.GetCountByRiskID(riskID)
}
