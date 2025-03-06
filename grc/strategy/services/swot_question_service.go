package services

import (
	"errors"
	"fmt"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/strategy/repositories"
	"jethur/utility"
	"log"

	"gorm.io/gorm"
)

type SwotQuestionServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.SwotQuestion, error)
	GetByID(id uint) (*models.SwotQuestion, error)
	Create(currentUserID uint, swotQuestion *models.SwotQuestion) (*models.SwotQuestion, error)
	Update(id uint, currentUserID uint, swotQuestion *models.SwotQuestion) (*models.SwotQuestion, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type SwotQuestionService struct {
	repo           repositories.SwotQuestionRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewSwotQuestionService(repo repositories.SwotQuestionRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) SwotQuestionServiceInterface {
	return &SwotQuestionService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (sqs *SwotQuestionService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.SwotQuestion, error) {
	return sqs.repo.GetAll(sort, sortBy, keyword, status)
}

func (sqs SwotQuestionService) GetByID(id uint) (*models.SwotQuestion, error) {
	swotQuestion, err := sqs.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return swotQuestion, nil
}

func (sqs *SwotQuestionService) Create(currentUserID uint, swotQuestion *models.SwotQuestion) (*models.SwotQuestion, error) {
	if exists, err := sqs.repo.ExistsByTitle(swotQuestion.SwotQuestionGroupID, swotQuestion.TitleEn, swotQuestion.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	existingQuestions, err := sqs.repo.GetAllOrdered(swotQuestion.SwotQuestionGroupID)
	if err != nil {
		return nil, err
	}

	newOrder := swotQuestion.Order
	for _, question := range existingQuestions {
		if question.Order >= newOrder {
			question.Order++ // Increment order for existing records
			if _, err := sqs.repo.Update(question.ID, question); err != nil {
				return nil, err
			}
		}
	}

	createdSwotQuestion, err := sqs.repo.Create(swotQuestion)
	if err != nil {
		return nil, err
	}

	go func() {
		err := sqs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.SwotQuestion{}.TableName(),
			createdSwotQuestion.ID,
			createdSwotQuestion,
			currentUserID,
			nil,
		)

		if err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdSwotQuestion, nil
}

func (sqs *SwotQuestionService) Update(id uint, currentUserID uint, updatedSwotQuestion *models.SwotQuestion) (*models.SwotQuestion, error) {
	existingSwotQuestion, err := sqs.GetByID(id)

	if err != nil {
		return nil, err
	}

	if updatedSwotQuestion.TitleEn != existingSwotQuestion.TitleEn || updatedSwotQuestion.TitleAr != existingSwotQuestion.TitleAr {
		if exists, err := sqs.repo.ExistsByTitleExcludingID(updatedSwotQuestion.SwotQuestionGroupID, id, updatedSwotQuestion.TitleEn, updatedSwotQuestion.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	}

	swotQuestions, err := sqs.repo.GetAllOrdered(updatedSwotQuestion.SwotQuestionGroupID)
	if err != nil {
		return nil, err
	}

	if updatedSwotQuestion.Order != existingSwotQuestion.Order {
		updatedOrder := updatedSwotQuestion.Order
		oldOrder := existingSwotQuestion.Order

		for _, question := range swotQuestions {
			if question.ID == id {
				continue
			}
			if oldOrder < updatedOrder && question.Order > oldOrder && question.Order <= updatedOrder {
				question.Order--
			} else if oldOrder > updatedOrder && question.Order >= updatedOrder && question.Order < oldOrder {
				question.Order++
			}

			if _, err := sqs.repo.Update(question.ID, question); err != nil {
				return nil, err
			}
		}
	}

	result, err := sqs.repo.Update(id, updatedSwotQuestion)

	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "order", "Status"}
	changes, err := utility.CompareFields(existingSwotQuestion, updatedSwotQuestion, fieldsToCompare)

	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := sqs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.SwotQuestion{}.TableName(),
			result.ID,
			changedDataJSON,
			currentUserID,
			nil,
		)
		if err != nil {
			// @todo - save to log file
			fmt.Printf("Failed to log bulk delete activity: %v\n", err)
		}
	}()
	return result, err
}

func (sqs *SwotQuestionService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint

	groupedQuestions := make(map[uint][]*models.SwotQuestion)

	for _, id := range ids {
		swotQuestion, err := sqs.repo.GetByID(id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		groupedQuestions[swotQuestion.SwotQuestionGroupID] = append(groupedQuestions[swotQuestion.SwotQuestionGroupID], swotQuestion)
	}

	for groupID, questionsToDelete := range groupedQuestions {
		swotQuestions, err := sqs.repo.GetAllOrdered(groupID)
		if err != nil {
			return nil, err
		}

		deletingMap := make(map[uint]bool)
		for _, q := range questionsToDelete {
			deletingMap[q.ID] = true
		}

		for _, q := range questionsToDelete {
			if err := sqs.repo.Delete(q.ID); err != nil {
				failedIds = append(failedIds, q.ID)
			} else {
				successfulIds = append(successfulIds, q.ID)
			}
		}

		newOrder := 1
		for _, q := range swotQuestions {
			if !deletingMap[q.ID] {
				q.Order = newOrder
				_, err := sqs.repo.Update(q.ID, q)
				if err != nil {
					failedIds = append(failedIds, q.ID)
				}
				newOrder++
			}
		}
	}

	if len(successfulIds) > 0 {
		go func() {
			err := sqs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.SwotQuestion{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log delete activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (sqs *SwotQuestionService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := sqs.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := sqs.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.SwotQuestion{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return err
}

func (sqs *SwotQuestionService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := sqs.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := sqs.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.SwotQuestion{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				// You might want to log this error or handle it appropriately
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return err
}
