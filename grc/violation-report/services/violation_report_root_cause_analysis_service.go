package services

import (
	"errors"
	"fmt"

	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	violationReportDTOs "jethur/grc/violation-report/dtos"
	violationReportRepositories "jethur/grc/violation-report/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type ViolationReportRootCauseAnalysisServiceInterface interface {
	GetAll(violationReportID uint, params *violationReportDTOs.ViolationReportRootCauseAnalysisListQueryParamsDTO) ([]*violationReportDTOs.ViolationReportRootCauseAnalysisListDTO, error)
	GetByID(violationReportID uint, id uint) (*models.ViolationReportRootCauseAnalysis, error)
	GetAllByViolationReportID(violationReportID uint, tx utility.DatabaseTransaction) ([]*models.ViolationReportRootCauseAnalysis, error)
	Create(violationReportID uint, currentUserID uint, rootCauseAnalysisFormDTO *violationReportDTOs.ViolationReportRootCauseAnalysisFormDTO) (uint, error)
	Update(violationReportID uint, id uint, currentUserID uint, rootCauseAnalysisFormDTO *violationReportDTOs.ViolationReportRootCauseAnalysisFormDTO) error
	Delete(violationReportID uint, id uint, currentUserID uint) (uint, error)
}

// Service and implementations
type ViolationReportRootCauseAnalysisService struct {
	repo               violationReportRepositories.ViolationReportRootCauseAnalysisRepositoryInterface
	activityLogSvc     logServices.ActivityLogServiceInterface
	violationReportSvc ViolationReportServiceInterface
}

func NewViolationReportRootCauseAnalysisService(
	repo violationReportRepositories.ViolationReportRootCauseAnalysisRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	violationReportSvc ViolationReportServiceInterface,
) ViolationReportRootCauseAnalysisServiceInterface {
	return &ViolationReportRootCauseAnalysisService{
		repo:               repo,
		activityLogSvc:     activityLogSvc,
		violationReportSvc: violationReportSvc,
	}
}

func (ircas *ViolationReportRootCauseAnalysisService) GetAll(violationReportID uint, params *violationReportDTOs.ViolationReportRootCauseAnalysisListQueryParamsDTO) ([]*violationReportDTOs.ViolationReportRootCauseAnalysisListDTO, error) {
	if _, err := ircas.violationReportSvc.GetByID(violationReportID); err != nil {
		return nil, err
	}
	return ircas.repo.GetAll(violationReportID, params)
}

func (ircas *ViolationReportRootCauseAnalysisService) GetByID(violationReportID uint, id uint) (*models.ViolationReportRootCauseAnalysis, error) {
	if _, err := ircas.violationReportSvc.GetByID(violationReportID); err != nil {
		return nil, err
	}

	violationReportRootCauseAnalysis, err := ircas.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return violationReportRootCauseAnalysis, nil
}

func (ircas *ViolationReportRootCauseAnalysisService) GetAllByViolationReportID(violationReportID uint, tx utility.DatabaseTransaction) ([]*models.ViolationReportRootCauseAnalysis, error) {
	return ircas.repo.GetAllByViolationReportID(violationReportID, tx)
}

func (ircas *ViolationReportRootCauseAnalysisService) Create(violationReportID uint, currentUserID uint, rootCauseAnalysisFormDTO *violationReportDTOs.ViolationReportRootCauseAnalysisFormDTO) (uint, error) {
	tx := ircas.repo.BeginTransaction()

	if _, err := ircas.violationReportSvc.GetByID(violationReportID); err != nil {
		return 0, err
	}
	// Increment the "Why" field by 1
	latestWhy, err := ircas.repo.GetLatestWhyByViolationReportID(violationReportID)
	if err != nil {
		return 0, err
	}
	why := latestWhy + 1

	violationReportRootCauseAnalysis := &models.ViolationReportRootCauseAnalysis{
		ViolationReportID: violationReportID,
		Why:               why,
		Question:          rootCauseAnalysisFormDTO.Question,
		Answer:            rootCauseAnalysisFormDTO.Answer,
		CreatedBYID:       currentUserID,
	}

	id, err := ircas.repo.Create(violationReportRootCauseAnalysis, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	go func() {
		existingViolationReportRootCauseAnalysis, err := ircas.GetByID(violationReportID, id)
		if err != nil {
			return
		}
		err = ircas.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ViolationReportRootCauseAnalysis{}.TableName(),
			id,
			existingViolationReportRootCauseAnalysis,
			currentUserID,
			"violation_report", "violation_report_register", "root_cause_analysis", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ViolationReportRootCauseAnalysis{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (ircas *ViolationReportRootCauseAnalysisService) Update(violationReportID uint, id uint, currentUserID uint, rootCauseAnalysisFormDTO *violationReportDTOs.ViolationReportRootCauseAnalysisFormDTO) error {
	tx := ircas.repo.BeginTransaction()

	if _, err := ircas.violationReportSvc.GetByID(violationReportID); err != nil {
		return err
	}
	oldViolationReportRootCauseAnalysis, err := ircas.GetByID(violationReportID, id)
	if err != nil {
		return err
	}

	existingViolationReportRootCauseAnalysis, err := ircas.GetByID(violationReportID, id)
	if err != nil {
		return err
	}

	// Update Violation Report RootCauseAnalysis
	existingViolationReportRootCauseAnalysis.Question = rootCauseAnalysisFormDTO.Question
	existingViolationReportRootCauseAnalysis.Answer = rootCauseAnalysisFormDTO.Answer
	existingViolationReportRootCauseAnalysis.UpdatedBYID = &currentUserID

	err = ircas.repo.Update(existingViolationReportRootCauseAnalysis, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		existingViolationReportRootCauseAnalysis, err := ircas.GetByID(violationReportID, id)
		if err != nil {
			return
		}
		// log changes
		changedData := make(map[string]interface{})
		oldData := map[string]interface{}{
			"question": oldViolationReportRootCauseAnalysis.Question,
			"answer":   oldViolationReportRootCauseAnalysis.Answer,
		}

		newData := map[string]interface{}{
			"question": existingViolationReportRootCauseAnalysis.Question,
			"answer":   existingViolationReportRootCauseAnalysis.Answer,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := ircas.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.ViolationReportRootCauseAnalysis{}.TableName(),
				id,
				changedData,
				currentUserID,
				"violation_report", "violation_report_register", "root_cause_analysis", nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.ViolationReportRootCauseAnalysis{}.TableName(),
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

func (ircas *ViolationReportRootCauseAnalysisService) Delete(violationReportID uint, id uint, currentUserID uint) (uint, error) {
	tx := ircas.repo.BeginTransaction()

	if _, err := ircas.violationReportSvc.GetByID(violationReportID); err != nil {
		return 0, err
	}

	var failedId uint
	var deletedViolationReportRootCauseAnalysis []map[string]interface{}

	entryToDelete, err := ircas.GetByID(violationReportID, id)
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
	err = ircas.repo.UpdateWhyValuesAfterDeletion(violationReportID, entryToDelete.Why, tx)
	if err != nil {
		tx.Rollback()
		fmt.Printf("Failed to update 'Why' values after deletion: %v\n", err)
		return 0, err
	}
	deletedViolationReportRootCauseAnalysis = append(deletedViolationReportRootCauseAnalysis, map[string]interface{}{
		"id":                  id,
		"violation_report_id": violationReportID,
	})
	tx.Commit()

	// Log successful delete activity
	go func() {
		err := ircas.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.ViolationReportRootCauseAnalysis{}.TableName(),
			id,
			deletedViolationReportRootCauseAnalysis,
			currentUserID,
			"violation_report", "violation_report_register", "root_cause_analysis", nil,
		)
		if err != nil {
			// Handle error if logging fails
			fmt.Printf("Failed to log delete activity: %v\n", err)
		}
	}()

	return failedId, nil
}
