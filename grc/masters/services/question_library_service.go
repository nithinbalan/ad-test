package services

import (
	"fmt"
	"jethur/grc/logs/services"
	masterDTOs "jethur/grc/masters/dtos"
	"jethur/grc/masters/repositories"
	"jethur/grc/models"
	"jethur/utility"
	"log"
)

type QuestionLibraryServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string, page int, pageSize int) (*masterDTOs.QuestionLibraryListResponseDTO, error)
	GetByID(id uint) (*models.QuestionLibrary, error)
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

func (qls QuestionLibraryService) GetByID(id uint) (*models.QuestionLibrary, error) {
	return qls.repo.GetByID(id)
}

func (qls *QuestionLibraryService) Create(currentUserID uint, questionLibraryForm *masterDTOs.QuestionLibraryFormDTO) (uint, error) {

	if err := qls.validateQuestionLibraryForm(questionLibraryForm, nil); err != nil {
		return 0, err
	}

	questionLibrary := &models.QuestionLibrary{
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
			models.CreateAction,
			models.QuestionLibrary{}.TableName(),
			id,
			questionLibraryForm,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.QuestionLibrary{}.TableName(),
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

	if err := qls.validateQuestionLibraryForm(updatedQuestionLibraryForm, &id); err != nil {
		return err
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
			models.UpdateAction,
			models.QuestionLibrary{}.TableName(),
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
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.QuestionLibrary{}.TableName(),
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
				models.ActivateAction,
				models.QuestionLibrary{}.TableName(),
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
				models.DeactivateAction,
				models.QuestionLibrary{}.TableName(),
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

func (qls *QuestionLibraryService) validateQuestionLibraryForm(questionLibraryForm *masterDTOs.QuestionLibraryFormDTO, id *uint) error {

	if id == nil {

		if exists, err := qls.repo.ExistsByTitle(questionLibraryForm.QuestionEn, questionLibraryForm.QuestionAr, questionLibraryForm.QuestionGroupID); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the question must be unique"), "unique", "question")
		}

	} else {

		if exists, err := qls.repo.ExistsByTitleExcludingID(*id, questionLibraryForm.QuestionEn, questionLibraryForm.QuestionAr, questionLibraryForm.QuestionGroupID); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the question must be unique"), "unique", "question")
		}

	}

	return nil
}
