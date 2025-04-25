package services

import (
	"errors"
	"fmt"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	organizationalContextDTOs "jethur/grc/organizational-context/dtos"
	"jethur/grc/organizational-context/repositories"
	"jethur/utility"
	"log"
	"time"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type SwotPestelQuestionServiceInterface interface {
	GetAll(params *organizationalContextDTOs.SwotPestelQuestionListQueryParamsDTO) ([]*organizationalContextDTOs.SwotPestelQuestionListDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*organizationalContextDTOs.SwotPestelQuestionDTO, error)
	GetByID(id uint) (*models.SwotPestelQuestion, error)
	Create(currentUserID uint, swotPestelQuestion *models.SwotPestelQuestion) (*models.SwotPestelQuestion, error)
	Update(id uint, currentUserID uint, swotPestelQuestion *models.SwotPestelQuestion) (*models.SwotPestelQuestion, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error

	UpdateAnswer(id uint, currentUserID uint, answerFormDTO *organizationalContextDTOs.AnswerFormDTO) error
	CreateInvitees(id uint, currentUserID uint, inviteeFormDTO *organizationalContextDTOs.InviteeFormDTO) error
	DeleteInvitees(questionId uint, ids []*uint, currentUserID uint) ([]uint, error)
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

func (sqs *SwotPestelQuestionService) GetAll(params *organizationalContextDTOs.SwotPestelQuestionListQueryParamsDTO) ([]*organizationalContextDTOs.SwotPestelQuestionListDTO, error) {
	// Fetch raw models from the repository
	return sqs.repo.GetAll(params)
}

func (sqs SwotPestelQuestionService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*organizationalContextDTOs.SwotPestelQuestionDTO, error) {
	// Initialize your DTO
	var swotPestelQuestionDTO organizationalContextDTOs.SwotPestelQuestionDTO

	question, err := sqs.repo.GetFullDetailsByID(id, tx...)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Map the result to your DTO
	err = mapstructure.Decode(question, &swotPestelQuestionDTO)
	if err != nil {
		return nil, err
	}

	swotPestelQuestionDTO.AnsweredAt = question.AnsweredAt

	return &swotPestelQuestionDTO, nil
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

func (sqs *SwotPestelQuestionService) UpdateAnswer(id uint, currentUserID uint, answerFormDTO *organizationalContextDTOs.AnswerFormDTO) error {
	existingAnswer, err := sqs.GetByID(id)
	if err != nil {
		return err
	}

	oldAnswer, err := sqs.GetByID(id)
	if err != nil {
		return err
	}

	now := time.Now()
	existingAnswer.Answer = answerFormDTO.Answer
	existingAnswer.AnsweredBYID = &currentUserID
	existingAnswer.AnsweredAt = &now

	result, err := sqs.repo.Update(id, existingAnswer)
	if err != nil {
		return err
	}

	go func() {
		existingAnswer, err := sqs.GetByID(id)
		if err != nil {
			return
		}

		changedData := make(map[string]interface{})

		//log changes
		oldData := map[string]interface{}{
			"answer": oldAnswer.Answer,
		}

		newData := map[string]interface{}{
			"answer": existingAnswer.Answer,
		}

		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		if len(changedData) > 0 {
			err := sqs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.SwotPestelQuestion{}.TableName(),
				result.ID,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}
	}()
	return err
}

func (sqs *SwotPestelQuestionService) CreateInvitees(id uint, currentUserID uint, inviteeFormDTO *organizationalContextDTOs.InviteeFormDTO) error {
	question, err := sqs.repo.GetFullDetailsByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, id)
		}
		return err
	}

	var existingInvitees []*uint
	for _, user := range question.SwotPestelQuestionInvitees {
		userID := user.UserID
		existingInvitees = append(existingInvitees, &userID)
	}

	if err := sqs.repo.DeleteSwotPestelQuestionInvitees(id, existingInvitees); err != nil {
		return err
	}

	// Prepare invitee records
	var invitees []*models.SwotPestelQuestionInvitee
	for _, userID := range inviteeFormDTO.UserIds {
		if userID != nil {
			invitee := &models.SwotPestelQuestionInvitee{
				SwotPestelQuestionID: id,
				UserID:               *userID,
				CreatedBYID:          currentUserID,
			}
			invitees = append(invitees, invitee)
		}
	}

	if len(invitees) > 0 {
		if err := sqs.repo.CreateSwotPestelQuestionInvitees(invitees); err != nil {
			return err
		}
	}

	return nil

}

func (spqfs *SwotPestelQuestionService) DeleteInvitees(questionId uint, ids []*uint, currentUserID uint) ([]uint, error) {
	tx := spqfs.repo.BeginTransaction()
	var failedIds []uint
	var successfulIds []uint
	var deletedSwotPestelQuestionInvitee []map[string]interface{}

	for _, id := range ids {

		deletedSwotPestelQuestionInvitee = append(deletedSwotPestelQuestionInvitee, map[string]interface{}{
			"id": id,
		})

		// Delete the business
		if err := spqfs.repo.DeleteSwotPestelQuestionInvitees(questionId, ids); err != nil {
			failedIds = append(failedIds, *id)
			continue
		} else {
			successfulIds = append(successfulIds, *id)
		}
	}
	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}

	tx.Commit()

	// Log bulk delete activities for successful deletions
	if len(successfulIds) > 0 {
		go func() {
			err := spqfs.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.SwotPestelQuestionInvitee{}.TableName(),
				successfulIds,
				deletedSwotPestelQuestionInvitee,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}

		}()
	}

	return failedIds, nil
}
