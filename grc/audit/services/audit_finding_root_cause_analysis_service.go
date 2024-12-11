package services

import (
	"errors"
	"fmt"
	auditDTOs "jethur/grc/audit/dtos"
	auditRepositories "jethur/grc/audit/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditFindingRootCauseAnalysisServiceInterface interface {
	GetAll(findingID uint, params *auditDTOs.AuditFindingRootCauseAnalysisListQueryParamsDTO) ([]*auditDTOs.AuditFindingRootCauseAnalysisListDTO, error)
	GetByID(findingID uint, id uint) (*models.AuditFindingRootCauseAnalysis, error)
	Create(findingID uint, currentUserID uint, rootCauseAnalysisFormDTO *auditDTOs.AuditFindingRootCauseAnalysisFormDTO) (uint, error)
	Update(findingID uint, id uint, currentUserID uint, rootCauseAnalysisFormDTO *auditDTOs.AuditFindingRootCauseAnalysisFormDTO) error
	Delete(findingID uint, id uint, currentUserID uint) (uint, error)
}

// Service and implementations
type AuditFindingRootCauseAnalysisService struct {
	repo             auditRepositories.AuditFindingRootCauseAnalysisRepositoryInterface
	activityLogSvc   logServices.ActivityLogServiceInterface
	auditFindingSvc  AuditFindingServiceInterface
	auditTimelineSvc AuditTimelineServiceInterface
}

func NewAuditFindingRootCauseAnalysisService(
	repo auditRepositories.AuditFindingRootCauseAnalysisRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditFindingSvc AuditFindingServiceInterface,
	auditTimelineSvc AuditTimelineServiceInterface,
) AuditFindingRootCauseAnalysisServiceInterface {
	return &AuditFindingRootCauseAnalysisService{
		repo:             repo,
		activityLogSvc:   activityLogSvc,
		auditFindingSvc:  auditFindingSvc,
		auditTimelineSvc: auditTimelineSvc,
	}
}

func (afrcas *AuditFindingRootCauseAnalysisService) GetAll(findingID uint, params *auditDTOs.AuditFindingRootCauseAnalysisListQueryParamsDTO) ([]*auditDTOs.AuditFindingRootCauseAnalysisListDTO, error) {
	if _, err := afrcas.auditFindingSvc.GetByID(findingID); err != nil {
		return nil, err
	}
	return afrcas.repo.GetAll(findingID, params)
}

func (afrcas *AuditFindingRootCauseAnalysisService) GetByID(findingID uint, id uint) (*models.AuditFindingRootCauseAnalysis, error) {
	if _, err := afrcas.auditFindingSvc.GetByID(findingID); err != nil {
		return nil, err
	}

	auditFindingRootCauseAnalysis, err := afrcas.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return auditFindingRootCauseAnalysis, nil
}

func (afrcas *AuditFindingRootCauseAnalysisService) Create(findingID uint, currentUserID uint, rootCauseAnalysisFormDTO *auditDTOs.AuditFindingRootCauseAnalysisFormDTO) (uint, error) {
	tx := afrcas.repo.BeginTransaction()

	if _, err := afrcas.auditFindingSvc.GetByID(findingID); err != nil {
		return 0, err
	}
	// Increment the "Why" field by 1
	latestWhy, err := afrcas.repo.GetLatestWhyByFindingID(findingID)
	if err != nil {
		return 0, err
	}
	why := latestWhy + 1

	auditFindingRootCauseAnalysis := &models.AuditFindingRootCauseAnalysis{
		AuditFindingID: findingID,
		Why:            why,
		Question:       rootCauseAnalysisFormDTO.Question,
		Answer:         rootCauseAnalysisFormDTO.Answer,
		CreatedBYID:    currentUserID,
	}

	id, err := afrcas.repo.Create(auditFindingRootCauseAnalysis, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	go func() {
		existingAuditFindingRootCauseAnalysis, err := afrcas.GetByID(findingID, id)
		if err != nil {
			utility.LogError("Failed to fetch audit finding root cause", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		err = afrcas.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditFindingRootCauseAnalysis{}.TableName(),
			id,
			existingAuditFindingRootCauseAnalysis,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AuditFindingRootCauseAnalysis{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (afrcas *AuditFindingRootCauseAnalysisService) Update(findingID uint, id uint, currentUserID uint, rootCauseAnalysisFormDTO *auditDTOs.AuditFindingRootCauseAnalysisFormDTO) error {
	tx := afrcas.repo.BeginTransaction()

	oldAuditFindingRootCauseAnalysis, err := afrcas.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}

	if _, err := afrcas.auditFindingSvc.GetByID(findingID); err != nil {
		return err
	}

	existingAuditFindingRootCauseAnalysis, err := afrcas.GetByID(findingID, id)

	if err != nil {
		return err
	}

	// Update Audit Finding RootCauseAnalysis
	existingAuditFindingRootCauseAnalysis.Question = rootCauseAnalysisFormDTO.Question
	existingAuditFindingRootCauseAnalysis.Answer = rootCauseAnalysisFormDTO.Answer
	existingAuditFindingRootCauseAnalysis.UpdatedBYID = &currentUserID

	err = afrcas.repo.Update(existingAuditFindingRootCauseAnalysis, tx)

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	go func() {
		existingFindingRootCauseAnalysis, err := afrcas.repo.GetByID(id)
		if err != nil {
			return
		}
		// Prepare log data
		changedData := make(map[string]interface{})
		//log changes
		oldData := map[string]interface{}{
			"question": oldAuditFindingRootCauseAnalysis.Question,
			"answer":   oldAuditFindingRootCauseAnalysis.Answer,
		}

		newData := map[string]interface{}{
			"question": existingFindingRootCauseAnalysis.Question,
			"answer":   existingFindingRootCauseAnalysis.Answer,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {

			err := afrcas.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.AuditFindingRootCauseAnalysis{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.AuditFindingRootCauseAnalysis{}.TableName(),
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

func (afrcas *AuditFindingRootCauseAnalysisService) Delete(findingID uint, id uint, currentUserID uint) (uint, error) {
	tx := afrcas.repo.BeginTransaction()

	if _, err := afrcas.auditFindingSvc.GetByID(findingID); err != nil {
		return 0, err
	}

	var failedId uint
	var deletedFindingRootCauseAnalysis []map[string]interface{}
	entryToDelete, err := afrcas.GetByID(findingID, id)
	if err != nil {
		return 0, err
	}

	// Delete the entry
	if err := afrcas.repo.Delete(id, tx); err != nil {
		tx.Rollback()
		failedId = id
		return failedId, err
	}
	// Reassign "Why" values for subsequent entries
	err = afrcas.repo.UpdateWhyValuesAfterDeletion(findingID, entryToDelete.Why, tx)
	if err != nil {
		tx.Rollback()
		fmt.Printf("Failed to update 'Why' values after deletion: %v\n", err)
		return 0, err
	}

	deletedFindingRootCauseAnalysis = append(deletedFindingRootCauseAnalysis, map[string]interface{}{
		"id":            id,
		"audit_finding": entryToDelete.AuditFinding.Title,
		"question":      entryToDelete.Question,
	})

	tx.Commit()

	// Log successful delete activity
	go func() {
		err := afrcas.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.AuditFindingRootCauseAnalysis{}.TableName(),
			id,
			deletedFindingRootCauseAnalysis,
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
