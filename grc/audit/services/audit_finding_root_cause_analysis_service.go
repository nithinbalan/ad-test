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
	GetAll(auditID uint, findingID uint, params *auditDTOs.AuditFindingRootCauseAnalysisListQueryParamsDTO) ([]*auditDTOs.AuditFindingRootCauseAnalysisListDTO, error)
	GetByID(auditID uint, findingID uint, id uint) (*models.AuditFindingRootCauseAnalysis, error)
	Create(auditID uint, findingID uint, currentUserID uint, rootCauseAnalysisFormDTO *auditDTOs.AuditFindingRootCauseAnalysisFormDTO) (uint, error)
	Update(auditID uint, findingID uint, id uint, currentUserID uint, rootCauseAnalysisFormDTO *auditDTOs.AuditFindingRootCauseAnalysisFormDTO) error
	Delete(auditID uint, findingID uint, id uint, currentUserID uint) (uint, error)
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

func (afrcas *AuditFindingRootCauseAnalysisService) GetAll(auditID uint, findingID uint, params *auditDTOs.AuditFindingRootCauseAnalysisListQueryParamsDTO) ([]*auditDTOs.AuditFindingRootCauseAnalysisListDTO, error) {
	if _, err := afrcas.auditFindingSvc.GetByID(findingID); err != nil {
		return nil, err
	}
	return afrcas.repo.GetAll(findingID, params)
}

func (afrcas *AuditFindingRootCauseAnalysisService) GetByID(auditID uint, findingID uint, id uint) (*models.AuditFindingRootCauseAnalysis, error) {
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

func (afrcas *AuditFindingRootCauseAnalysisService) Create(auditID uint, findingID uint, currentUserID uint, rootCauseAnalysisFormDTO *auditDTOs.AuditFindingRootCauseAnalysisFormDTO) (uint, error) {
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
		err := afrcas.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditFindingRootCauseAnalysis{}.TableName(),
			id,
			rootCauseAnalysisFormDTO,
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

func (afrcas *AuditFindingRootCauseAnalysisService) Update(auditID uint, findingID uint, id uint, currentUserID uint, rootCauseAnalysisFormDTO *auditDTOs.AuditFindingRootCauseAnalysisFormDTO) error {
	tx := afrcas.repo.BeginTransaction()

	if _, err := afrcas.auditFindingSvc.GetByID(findingID); err != nil {
		return err
	}

	existingAuditFindingRootCauseAnalysis, err := afrcas.GetByID(auditID, findingID, id)

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

	return nil
}

func (afrcas *AuditFindingRootCauseAnalysisService) Delete(auditID uint, findingID uint, id uint, currentUserID uint) (uint, error) {
	tx := afrcas.repo.BeginTransaction()

	if _, err := afrcas.auditFindingSvc.GetByID(findingID); err != nil {
		return 0, err
	}

	var failedId uint

	entryToDelete, err := afrcas.GetByID(auditID, findingID, id)
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
	tx.Commit()

	// Log successful delete activity
	go func() {
		err := afrcas.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.AuditFindingRootCauseAnalysis{}.TableName(),
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
