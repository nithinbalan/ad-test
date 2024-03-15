package services

import (
	"fmt"
	activityLogModels "jethur/grc/logs/models"
	"jethur/grc/logs/services"
	"jethur/grc/masters/models"
	"jethur/grc/masters/repositories"
	"jethur/utility"
	"log"
)

type ChecklistQuestionServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string) ([]*models.ChecklistQuestionWithChecklistTitle, error)
	GetByID(id uint) (*models.ChecklistQuestion, error)
	Create(currentUserID uint, checklistQuestion *models.ChecklistQuestion) (*models.ChecklistQuestion, error)
	Update(id uint, currentUserID uint, checklistQuestion *models.ChecklistQuestion) (*models.ChecklistQuestion, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type ChecklistQuestionService struct {
	repo           repositories.ChecklistQuestionRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewChecklistQuestionService(repo repositories.ChecklistQuestionRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) ChecklistQuestionServiceInterface {
	return &ChecklistQuestionService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (cqs *ChecklistQuestionService) GetAll(sort string, sortBy string, keyword string) ([]*models.ChecklistQuestionWithChecklistTitle, error) {
	return cqs.repo.GetAll(sort, sortBy, keyword)
}

func (cqs ChecklistQuestionService) GetByID(id uint) (*models.ChecklistQuestion, error) {
	return cqs.repo.GetByID(id)
}

func (cqs *ChecklistQuestionService) Create(currentUserID uint, checklistQuestion *models.ChecklistQuestion) (*models.ChecklistQuestion, error) {

	if exists, err := cqs.repo.ExistsByTitle(checklistQuestion.QuestionEn, checklistQuestion.QuestionAr, checklistQuestion.ChecklistID); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the question must be unique"), "unique", "question")
	}

	if existsInLibrary, err := cqs.repo.ExistsQuestionLibrary(checklistQuestion.QuestionLibraryID); err != nil {
		return nil, err
	} else if existsInLibrary {
		return nil, utility.NewValidationError(fmt.Errorf("the question already exists in this checklist"), "unique", "question")
	}

	createdChecklistQuestion, err := cqs.repo.Create(checklistQuestion)
	if err != nil {
		return nil, err
	}

	go func() {
		err := cqs.activityLogSvc.LogActivity(
			activityLogModels.CreateAction,
			models.ChecklistQuestion{}.TableName(),
			checklistQuestion.ID,
			createdChecklistQuestion,
			currentUserID,
			nil,
		)

		if err != nil {
			// @todo - save to log file
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdChecklistQuestion, nil
}

func (cqs *ChecklistQuestionService) Update(id uint, currentUserID uint, updatedChecklistQuestion *models.ChecklistQuestion) (*models.ChecklistQuestion, error) {
	existingChecklistQuestion, err := cqs.repo.GetByID(id)

	if err != nil {
		return nil, err
	}

	if updatedChecklistQuestion.QuestionEn != existingChecklistQuestion.QuestionEn || updatedChecklistQuestion.QuestionAr != existingChecklistQuestion.QuestionAr || updatedChecklistQuestion.ChecklistID != existingChecklistQuestion.ChecklistID {
		if exists, err := cqs.repo.ExistsByTitleExcludingID(id, updatedChecklistQuestion.QuestionEn, updatedChecklistQuestion.QuestionAr, updatedChecklistQuestion.ChecklistID); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the question must be unique"), "unique", "question")
		}
	}

	if existsInLibrary, err := cqs.repo.ExistsQuestionLibraryExcludingID(id, updatedChecklistQuestion.QuestionLibraryID); err != nil {
		return nil, err
	} else if existsInLibrary {
		return nil, utility.NewValidationError(fmt.Errorf("the question already exists in this checklist"), "unique", "question")
	}

	result, err := cqs.repo.Update(id, updatedChecklistQuestion)

	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"QuestionEn", "QuestionAr"}
	changes, err := utility.CompareFields(existingChecklistQuestion, updatedChecklistQuestion, fieldsToCompare)

	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := cqs.activityLogSvc.LogActivity(
			activityLogModels.UpdateAction,
			models.ChecklistQuestion{}.TableName(),
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

func (cqs *ChecklistQuestionService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := cqs.repo.GetByID(id)

		if err != nil {
			continue
		}

		err = cqs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := cqs.activityLogSvc.LogBulkActivities(
				activityLogModels.DeleteAction, // Make sure this action constant is defined appropriately
				models.ChecklistQuestion{}.TableName(),
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
