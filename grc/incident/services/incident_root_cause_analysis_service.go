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
	Update(incidentID uint, currentUserID uint, rootCauseAnalysisFormDTO *incidentDTOs.IncidentRootCauseAnalysisFormDTO) error
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
	tx := ircas.repo.BeginTransaction()

	if _, err := ircas.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return 0, err
	}

	var incidentRootCauseAnalyses []*models.IncidentRootCauseAnalysis

	for _, analysisData := range rootCauseAnalysisFormDTO.Analyses {
		incidentRootCauseAnalysis := &models.IncidentRootCauseAnalysis{
			IncidentRegisterID: incidentID,
			Why:                analysisData.Why,
			Question:           analysisData.Question,
			Answer:             analysisData.Answer,
			CreatedByID:        currentUserID,
		}
		incidentRootCauseAnalyses = append(incidentRootCauseAnalyses, incidentRootCauseAnalysis)

	}
	if len(incidentRootCauseAnalyses) > 0 {
		if err := ircas.repo.Create(incidentRootCauseAnalyses, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	tx.Commit()

	go func() {
		err := ircas.activityLogSvc.LogActivity(
			models.CreateAction,
			models.IncidentRootCauseAnalysis{}.TableName(),
			incidentID,
			rootCauseAnalysisFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.IncidentRootCauseAnalysis{}.TableName(),
				"recordID":  incidentID,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return incidentID, nil
}

func (ircas *IncidentRootCauseAnalysisService) Update(incidentID uint, currentUserID uint, rootCauseAnalysisFormDTO *incidentDTOs.IncidentRootCauseAnalysisFormDTO) error {
	tx := ircas.repo.BeginTransaction()

	if _, err := ircas.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return err
	}

	// Get all existing incident root cause analyses for the given incident ID
	existingAnalyses, err := ircas.repo.GetByIncidentID(incidentID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Create a map to store existing analyses by their ID for easier lookup
	existingAnalysisMap := make(map[uint]*models.IncidentRootCauseAnalysis)
	for _, analysis := range existingAnalyses {
		existingAnalysisMap[analysis.ID] = analysis
	}

	var updatedAnalyses []*models.IncidentRootCauseAnalysis
	var addedAnalyses []*models.IncidentRootCauseAnalysis

	// Update or delete existing analyses based on the data in the form DTO
	for _, analysisData := range rootCauseAnalysisFormDTO.Analyses {
		if analysisData.ID == nil {
			// Analysis ID is nil, indicating it's a new analysis
			newAnalysis := &models.IncidentRootCauseAnalysis{
				IncidentRegisterID: incidentID,
				Why:                analysisData.Why,
				Question:           analysisData.Question,
				Answer:             analysisData.Answer,
				CreatedByID:        currentUserID,
			}
			addedAnalyses = append(addedAnalyses, newAnalysis)
		} else if existingAnalysis, ok := existingAnalysisMap[*analysisData.ID]; ok {
			// Analysis exists, update it
			existingAnalysis.Why = analysisData.Why
			existingAnalysis.Question = analysisData.Question
			existingAnalysis.Answer = analysisData.Answer
			updatedAnalyses = append(updatedAnalyses, existingAnalysis)
			// Remove the analysis from the map since it has been updated
			delete(existingAnalysisMap, *analysisData.ID)
		}

	}
	// Delete any analyses that are no longer present in the form DTO
	var deletedAnalysisIDs []*uint
	for _, analysis := range existingAnalysisMap {
		deletedAnalysisIDs = append(deletedAnalysisIDs, &analysis.ID)
	}
	// Delete existing analyses
	if len(deletedAnalysisIDs) > 0 {
		if err := ircas.repo.DeleteAnalyses(deletedAnalysisIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	// Update existing analyses
	if len(updatedAnalyses) > 0 {
		if err := ircas.repo.Update(updatedAnalyses, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Add new analyses
	if len(addedAnalyses) > 0 {
		if err := ircas.repo.Create(addedAnalyses, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
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
