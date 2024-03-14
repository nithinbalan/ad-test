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
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ChecklistQuestion, error)
	GetByID(id uint) (*models.ChecklistQuestion, error)
	Create(currentUserID uint, checklistQuestion *models.ChecklistQuestion) (*models.ChecklistQuestion, error)
	Update(id uint, currentUserID uint, checklistQuestion *models.ChecklistQuestion) (*models.ChecklistQuestion, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
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

func (cqs *ChecklistQuestionService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ChecklistQuestion, error) {
	return cqs.repo.GetAll(sort, sortBy, keyword, status)
}

func (cqs ChecklistQuestionService) GetByID(id uint) (*models.ChecklistQuestion, error) {
	return cqs.repo.GetByID(id)
}

func (cqs *ChecklistQuestionService) Create(currentUserID uint, checklistQuestion *models.ChecklistQuestion) (*models.ChecklistQuestion, error) {

	if exists, err := cqs.repo.ExistsByTitle(checklistQuestion.TitleEn, checklistQuestion.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
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

	if updatedChecklistQuestion.TitleEn != existingChecklistQuestion.TitleEn || updatedChecklistQuestion.TitleAr != existingChecklistQuestion.TitleAr {
		if exists, err := cqs.repo.ExistsByTitleExcludingID(id, updatedChecklistQuestion.TitleEn, updatedChecklistQuestion.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	}

	result, err := cqs.repo.Update(id, updatedChecklistQuestion)

	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status"}
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

func (cqs *ChecklistQuestionService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := cqs.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := cqs.activityLogSvc.LogBulkActivities(
				activityLogModels.ActivateAction,
				models.ChecklistQuestion{}.TableName(),
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

func (cqs *ChecklistQuestionService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := cqs.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := cqs.activityLogSvc.LogBulkActivities(
				activityLogModels.DeactivateAction,
				models.ChecklistQuestion{}.TableName(),
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
