package services

import (
	"errors"
	"fmt"
	incidentDTOs "jethur/grc/incident/dtos"
	incidentRepositories "jethur/grc/incident/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type IncidentRootCauseAnalysisServiceInterface interface {
	GetAll(incidentID uint, params *incidentDTOs.IncidentRootCauseAnalysisListQueryParamsDTO) ([]*incidentDTOs.IncidentRootCauseAnalysisListDTO, error)
	GetByID(incidentID uint, id uint) (*models.IncidentRootCauseAnalysis, error)
	Create(incidentID uint, currentUserID uint, rootCauseAnalysisFormDTO *incidentDTOs.IncidentRootCauseAnalysisFormDTO) (uint, error)
	Update(incidentID uint, id uint, currentUserID uint, rootCauseAnalysisFormDTO *incidentDTOs.IncidentRootCauseAnalysisFormDTO) error
	Delete(incidentID uint, id uint, currentUserID uint) (uint, error)
}

// Service and implementations
type IncidentRootCauseAnalysisService struct {
	repo                incidentRepositories.IncidentRootCauseAnalysisRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	incidentRegisterSvc IncidentRegisterServiceInterface
}

func NewIncidentRootCauseAnalysisService(
	repo incidentRepositories.IncidentRootCauseAnalysisRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	incidentRegisterSvc IncidentRegisterServiceInterface,
) IncidentRootCauseAnalysisServiceInterface {
	return &IncidentRootCauseAnalysisService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		incidentRegisterSvc: incidentRegisterSvc,
	}
}

func (ircas *IncidentRootCauseAnalysisService) GetAll(incidentID uint, params *incidentDTOs.IncidentRootCauseAnalysisListQueryParamsDTO) ([]*incidentDTOs.IncidentRootCauseAnalysisListDTO, error) {
	if _, err := ircas.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return nil, err
	}
	return ircas.repo.GetAll(incidentID, params)
}

func (ircas *IncidentRootCauseAnalysisService) GetByID(incidentID uint, id uint) (*models.IncidentRootCauseAnalysis, error) {
	if _, err := ircas.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return nil, err
	}

	incidentRootCauseAnalysis, err := ircas.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return incidentRootCauseAnalysis, nil
}

func (ircas *IncidentRootCauseAnalysisService) Create(incidentID uint, currentUserID uint, rootCauseAnalysisFormDTO *incidentDTOs.IncidentRootCauseAnalysisFormDTO) (uint, error) {
	if _, err := ircas.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return 0, err
	}
	// Increment the "Why" field by 1
	latestWhy, err := ircas.repo.GetLatestWhyByIncidentID(incidentID)
	if err != nil {
		return 0, err
	}
	why := latestWhy + 1

	incidentRootCauseAnalysis := &models.IncidentRootCauseAnalysis{
		IncidentRegisterID: incidentID,
		Why:                why,
		Question:           rootCauseAnalysisFormDTO.Question,
		Answer:             rootCauseAnalysisFormDTO.Answer,
		CreatedByID:        currentUserID,
	}

	id, err := ircas.repo.Create(incidentRootCauseAnalysis)
	if err != nil {
		return 0, err
	}

	go func() {
		err := ircas.activityLogSvc.LogActivity(
			models.CreateAction,
			models.IncidentRootCauseAnalysis{}.TableName(),
			id,
			rootCauseAnalysisFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.IncidentRootCauseAnalysis{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (ircas *IncidentRootCauseAnalysisService) Update(incidentID uint, id uint, currentUserID uint, rootCauseAnalysisFormDTO *incidentDTOs.IncidentRootCauseAnalysisFormDTO) error {

	if _, err := ircas.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return err
	}

	existingIncidentRootCauseAnalysis, err := ircas.GetByID(incidentID, id)

	if err != nil {
		return err
	}

	// Update Incident RootCauseAnalysis
	existingIncidentRootCauseAnalysis.Question = rootCauseAnalysisFormDTO.Question
	existingIncidentRootCauseAnalysis.Answer = rootCauseAnalysisFormDTO.Answer
	existingIncidentRootCauseAnalysis.UpdatedByID = &currentUserID

	err = ircas.repo.Update(existingIncidentRootCauseAnalysis)

	if err != nil {
		return err
	}

	return nil
}

func (ircas *IncidentRootCauseAnalysisService) Delete(incidentID uint, id uint, currentUserID uint) (uint, error) {
	tx := ircas.repo.BeginTransaction()

	if _, err := ircas.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return 0, err
	}

	var failedId uint

	entryToDelete, err := ircas.GetByID(incidentID, id)
	if err != nil {
		return 0, err
	}

	// Delete the entry
	if err := ircas.repo.Delete(id, tx); err != nil {
		tx.Rollback()
		failedId = id
		return failedId, err
	}
	// Reassign "Why" values for subsequent entries
	err = ircas.repo.UpdateWhyValuesAfterDeletion(incidentID, entryToDelete.Why, tx)
	if err != nil {
		tx.Rollback()
		fmt.Printf("Failed to update 'Why' values after deletion: %v\n", err)
		return 0, err
	}
	tx.Commit()

	// Log successful delete activity
	go func() {
		err := ircas.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.IncidentRootCauseAnalysis{}.TableName(),
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
