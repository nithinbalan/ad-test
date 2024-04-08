package services

import (
	"fmt"
	activityLogModels "jethur/grc/logs/models"
	"jethur/grc/logs/services"
	masterDTOs "jethur/grc/masters/dtos"
	masterModels "jethur/grc/masters/models"
	"jethur/grc/masters/repositories"
	"jethur/utility"
	"log"
)

type QuestionLibraryServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string, page int, pageSize int) (*masterDTOs.QuestionLibraryListResponseDTO, error)
	GetByID(id uint) (*masterModels.QuestionLibrary, error)
	Create(currentUserID uint, questionLibraryForm *masterDTOs.QuestionLibraryFormDTO) (uint, error)
	Update(id uint, currentUserID uint, questionLibraryForm *masterDTOs.QuestionLibraryFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type QuestionLibraryService struct {
	repo           repositories.QuestionLibraryRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewQuestionLibraryService(repo repositories.QuestionLibraryRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) QuestionLibraryServiceInterface {
	return &QuestionLibraryService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (qls *QuestionLibraryService) GetAll(sort string, sortBy string, keyword string, status string, page int, pageSize int) (*masterDTOs.QuestionLibraryListResponseDTO, error) {
	return qls.repo.GetAll(sort, sortBy, keyword, status, page, pageSize)
}

func (qls QuestionLibraryService) GetByID(id uint) (*masterModels.QuestionLibrary, error) {
	return qls.repo.GetByID(id)
}

func (qls *QuestionLibraryService) Create(currentUserID uint, questionLibraryForm *masterDTOs.QuestionLibraryFormDTO) (uint, error) {

	if exists, err := qls.repo.ExistsByTitle(questionLibraryForm.QuestionEn, questionLibraryForm.QuestionAr, questionLibraryForm.QuestionGroupID); err != nil {
		return 0, err
	} else if exists {
		return 0, utility.NewValidationError(fmt.Errorf("the question must be unique"), "unique", "question")
	}

	questionLibrary := &masterModels.QuestionLibrary{
		QuestionGroupID: questionLibraryForm.QuestionGroupID,
		QuestionEn:      questionLibraryForm.QuestionEn,
		QuestionAr:      questionLibraryForm.QuestionAr,
		Status:          questionLibraryForm.Status,
	}

	id, err := qls.repo.Create(questionLibrary)
	if err != nil {
		return 0, err
	}

	go func() {
		err := qls.activityLogSvc.LogActivity(
			activityLogModels.CreateAction,
			masterModels.QuestionLibrary{}.TableName(),
			id,
			questionLibraryForm,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    activityLogModels.CreateAction,
				"tableName": masterModels.QuestionLibrary{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil

}

func (qls *QuestionLibraryService) Update(id uint, currentUserID uint, updatedQuestionLibraryForm *masterDTOs.QuestionLibraryFormDTO) error {
	existingQuestionLibrary, err := qls.repo.GetByID(id)

	if err != nil {
		return err
	}

	if updatedQuestionLibraryForm.QuestionEn != existingQuestionLibrary.QuestionEn || updatedQuestionLibraryForm.QuestionAr != existingQuestionLibrary.QuestionAr || updatedQuestionLibraryForm.QuestionGroupID != existingQuestionLibrary.QuestionGroupID {
		if exists, err := qls.repo.ExistsByTitleExcludingID(id, updatedQuestionLibraryForm.QuestionEn, updatedQuestionLibraryForm.QuestionAr, updatedQuestionLibraryForm.QuestionGroupID); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the question must be unique"), "unique", "question")
		}

	}

	existingQuestionLibrary.QuestionGroupID = updatedQuestionLibraryForm.QuestionGroupID
	existingQuestionLibrary.QuestionEn = updatedQuestionLibraryForm.QuestionEn
	existingQuestionLibrary.QuestionAr = updatedQuestionLibraryForm.QuestionAr
	existingQuestionLibrary.Status = updatedQuestionLibraryForm.Status

	result, err := qls.repo.Update(id, existingQuestionLibrary)

	if err != nil {
		return err
	}

	fieldsToCompare := []string{"QuestionEn", "QuestionAr", "QuestionGroupID"}
	changes, err := utility.CompareFields(existingQuestionLibrary, updatedQuestionLibraryForm, fieldsToCompare)

	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := qls.activityLogSvc.LogActivity(
			activityLogModels.UpdateAction,
			masterModels.QuestionLibrary{}.TableName(),
			result.ID,
			changedDataJSON,
			currentUserID,
			nil,
		)
		if err != nil {
			// @todo - save to log file
			fmt.Printf("Failed to log activity: %v\n", err)
		}
	}()
	return err
}

func (qls *QuestionLibraryService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := qls.repo.GetByID(id)

		if err != nil {
			continue
		}

		err = qls.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := qls.activityLogSvc.LogBulkActivities(
				activityLogModels.DeleteAction, // Make sure this action constant is defined appropriately
				masterModels.QuestionLibrary{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (qls *QuestionLibraryService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := qls.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := qls.activityLogSvc.LogBulkActivities(
				activityLogModels.ActivateAction,
				masterModels.QuestionLibrary{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log activity: %v\n", err)
			}
		}()
	}

	return err
}

func (qls *QuestionLibraryService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := qls.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := qls.activityLogSvc.LogBulkActivities(
				activityLogModels.DeactivateAction,
				masterModels.QuestionLibrary{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				// You might want to log this error or handle it appropriately
				fmt.Printf("Failed to log activity: %v\n", err)
			}
		}()
	}

	return err
}
