package services

import (
	"fmt"
	"jethur/grc/libraries/models"
	"jethur/grc/libraries/repositories"
	activityLogModels "jethur/grc/logs/models"
	"jethur/grc/logs/services"
	"jethur/utility"
	"log"
)

type QuestionLibraryServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string) ([]*models.QuestionLibrary, error)
	GetByID(id uint) (*models.QuestionLibrary, error)
	Create(currentUserID uint, questionLibrary *models.QuestionLibrary) (*models.QuestionLibrary, error)
	Update(id uint, currentUserID uint, questionLibrary *models.QuestionLibrary) (*models.QuestionLibrary, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
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

func (qls *QuestionLibraryService) GetAll(sort string, sortBy string, keyword string) ([]*models.QuestionLibrary, error) {
	return qls.repo.GetAll(sort, sortBy, keyword)
}

func (qls QuestionLibraryService) GetByID(id uint) (*models.QuestionLibrary, error) {
	return qls.repo.GetByID(id)
}

func (qls *QuestionLibraryService) Create(currentUserID uint, questionLibrary *models.QuestionLibrary) (*models.QuestionLibrary, error) {

	if exists, err := qls.repo.ExistsByTitle(questionLibrary.QuestionEn, questionLibrary.QuestionAr, questionLibrary.QuestionGroupID); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the question must be unique"), "unique", "question")
	}

	createdQuestionLibrary, err := qls.repo.Create(questionLibrary)
	if err != nil {
		return nil, err
	}

	go func() {
		err := qls.activityLogSvc.LogActivity(
			activityLogModels.CreateAction,
			models.QuestionLibrary{}.TableName(),
			questionLibrary.ID,
			createdQuestionLibrary,
			currentUserID,
			nil,
		)

		if err != nil {
			// @todo - save to log file
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdQuestionLibrary, nil
}

func (qls *QuestionLibraryService) Update(id uint, currentUserID uint, updatedQuestionLibrary *models.QuestionLibrary) (*models.QuestionLibrary, error) {
	existingQuestionLibrary, err := qls.repo.GetByID(id)

	if err != nil {
		return nil, err
	}

	if updatedQuestionLibrary.QuestionEn != existingQuestionLibrary.QuestionEn || updatedQuestionLibrary.QuestionAr != existingQuestionLibrary.QuestionAr || updatedQuestionLibrary.QuestionGroupID != existingQuestionLibrary.QuestionGroupID {
		if exists, err := qls.repo.ExistsByTitleExcludingID(id, updatedQuestionLibrary.QuestionEn, updatedQuestionLibrary.QuestionAr, updatedQuestionLibrary.QuestionGroupID); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the question must be unique"), "unique", "question")
		}

	}

	result, err := qls.repo.Update(id, updatedQuestionLibrary)

	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"QuestionEn", "QuestionAr", "QuestionGroupID"}
	changes, err := utility.CompareFields(existingQuestionLibrary, updatedQuestionLibrary, fieldsToCompare)

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
	return result, err
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
