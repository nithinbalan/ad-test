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

type QuestionGroupServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.QuestionGroup, error)
	GetByID(id uint) (*models.QuestionGroup, error)
	Create(currentUserID uint, questionGroup *models.QuestionGroup) (*models.QuestionGroup, error)
	Update(id uint, currentUserID uint, questionGroup *models.QuestionGroup) (*models.QuestionGroup, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type QuestionGroupService struct {
	repo           repositories.QuestionGroupRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewQuestionGroupService(repo repositories.QuestionGroupRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) QuestionGroupServiceInterface {
	return &QuestionGroupService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (qgs *QuestionGroupService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.QuestionGroup, error) {
	return qgs.repo.GetAll(sort, sortBy, keyword, status)
}

func (qgs QuestionGroupService) GetByID(id uint) (*models.QuestionGroup, error) {
	return qgs.repo.GetByID(id)
}

func (qgs *QuestionGroupService) Create(currentUserID uint, questionGroup *models.QuestionGroup) (*models.QuestionGroup, error) {

	if exists, err := qgs.repo.ExistsByTitle(questionGroup.TitleEn, questionGroup.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	createdQuestionGroup, err := qgs.repo.Create(questionGroup)
	if err != nil {
		return nil, err
	}

	go func() {
		err := qgs.activityLogSvc.LogActivity(
			activityLogModels.CreateAction,
			models.QuestionGroup{}.TableName(),
			questionGroup.ID,
			createdQuestionGroup,
			currentUserID,
			nil,
		)

		if err != nil {
			// @todo - save to log file
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdQuestionGroup, nil
}

func (qgs *QuestionGroupService) Update(id uint, currentUserID uint, updatedQuestionGroup *models.QuestionGroup) (*models.QuestionGroup, error) {
	existingQuestionGroup, err := qgs.repo.GetByID(id)

	if err != nil {
		return nil, err
	}

	if updatedQuestionGroup.TitleEn != existingQuestionGroup.TitleEn || updatedQuestionGroup.TitleAr != existingQuestionGroup.TitleAr {
		if exists, err := qgs.repo.ExistsByTitleExcludingID(id, updatedQuestionGroup.TitleEn, updatedQuestionGroup.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	}

	result, err := qgs.repo.Update(id, updatedQuestionGroup)

	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status"}
	changes, err := utility.CompareFields(existingQuestionGroup, updatedQuestionGroup, fieldsToCompare)

	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := qgs.activityLogSvc.LogActivity(
			activityLogModels.UpdateAction,
			models.QuestionGroup{}.TableName(),
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

func (qgs *QuestionGroupService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := qgs.repo.GetByID(id)

		if err != nil {
			continue
		}

		err = qgs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := qgs.activityLogSvc.LogBulkActivities(
				activityLogModels.DeleteAction, // Make sure this action constant is defined appropriately
				models.QuestionGroup{}.TableName(),
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

func (qgs *QuestionGroupService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := qgs.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := qgs.activityLogSvc.LogBulkActivities(
				activityLogModels.ActivateAction,
				models.QuestionGroup{}.TableName(),
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

func (qgs *QuestionGroupService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := qgs.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := qgs.activityLogSvc.LogBulkActivities(
				activityLogModels.DeactivateAction,
				models.QuestionGroup{}.TableName(),
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
