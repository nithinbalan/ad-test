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

type SwotPestelQuestionServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string, isSwot string, isPestel string) ([]*models.SwotPestelQuestion, error)
	GetByID(id uint) (*models.SwotPestelQuestion, error)
	Create(currentUserID uint, swotPestelQuestion *models.SwotPestelQuestion) (*models.SwotPestelQuestion, error)
	Update(id uint, currentUserID uint, swotPestelQuestion *models.SwotPestelQuestion) (*models.SwotPestelQuestion, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type SwotPestelQuestionService struct {
	repo           repositories.SwotPestelQuestionRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewSwotPestelQuestionService(repo repositories.SwotPestelQuestionRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) SwotPestelQuestionServiceInterface {
	return &SwotPestelQuestionService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (sqs *SwotPestelQuestionService) GetAll(sort string, sortBy string, keyword string, status string, isSwot string, isPestel string) ([]*models.SwotPestelQuestion, error) {
	return sqs.repo.GetAll(sort, sortBy, keyword, status, isSwot, isPestel)
}

func (sqs SwotPestelQuestionService) GetByID(id uint) (*models.SwotPestelQuestion, error) {
	swotPestelQuestion, err := sqs.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return swotPestelQuestion, nil
}

func (sqs *SwotPestelQuestionService) Create(currentUserID uint, swotPestelQuestion *models.SwotPestelQuestion) (*models.SwotPestelQuestion, error) {
	if exists, err := sqs.repo.ExistsByTitle(swotPestelQuestion.SwotPestelQuestionGroupID, swotPestelQuestion.TitleEn, swotPestelQuestion.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	existingQuestions, err := sqs.repo.GetAllOrdered(swotPestelQuestion.SwotPestelQuestionGroupID)
	if err != nil {
		return nil, err
	}

	newOrder := swotPestelQuestion.Order
	for _, question := range existingQuestions {
		if question.Order >= newOrder {
			question.Order++ // Increment order for existing records
			if _, err := sqs.repo.Update(question.ID, question); err != nil {
				return nil, err
			}
		}
	}

	createdSwotPestelQuestion, err := sqs.repo.Create(swotPestelQuestion)
	if err != nil {
		return nil, err
	}

	go func() {
		err := sqs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.SwotPestelQuestion{}.TableName(),
			createdSwotPestelQuestion.ID,
			createdSwotPestelQuestion,
			currentUserID,
			nil,
		)

		if err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdSwotPestelQuestion, nil
}

func (sqs *SwotPestelQuestionService) Update(id uint, currentUserID uint, updatedSwotPestelQuestion *models.SwotPestelQuestion) (*models.SwotPestelQuestion, error) {
	existingSwotPestelQuestion, err := sqs.GetByID(id)

	if err != nil {
		return nil, err
	}

	if updatedSwotPestelQuestion.TitleEn != existingSwotPestelQuestion.TitleEn || updatedSwotPestelQuestion.TitleAr != existingSwotPestelQuestion.TitleAr {
		if exists, err := sqs.repo.ExistsByTitleExcludingID(updatedSwotPestelQuestion.SwotPestelQuestionGroupID, id, updatedSwotPestelQuestion.TitleEn, updatedSwotPestelQuestion.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	}

	result, err := sqs.repo.Update(id, updatedSwotPestelQuestion)

	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "order", "Status"}
	changes, err := utility.CompareFields(existingSwotPestelQuestion, updatedSwotPestelQuestion, fieldsToCompare)

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
			models.SwotPestelQuestion{}.TableName(),
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

func (sqs *SwotPestelQuestionService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint

	groupedQuestions := make(map[uint][]*models.SwotPestelQuestion)

	for _, id := range ids {
		swotPestelQuestion, err := sqs.repo.GetByID(id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		groupedQuestions[swotPestelQuestion.SwotPestelQuestionGroupID] = append(groupedQuestions[swotPestelQuestion.SwotPestelQuestionGroupID], swotPestelQuestion)
	}

	for groupID, questionsToDelete := range groupedQuestions {
		swotPestelQuestions, err := sqs.repo.GetAllOrdered(groupID)
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
		for _, q := range swotPestelQuestions {
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
				models.SwotPestelQuestion{}.TableName(),
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

func (sqs *SwotPestelQuestionService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := sqs.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := sqs.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.SwotPestelQuestion{}.TableName(),
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

func (sqs *SwotPestelQuestionService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := sqs.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := sqs.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.SwotPestelQuestion{}.TableName(),
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
