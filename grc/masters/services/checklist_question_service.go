package services

import (
	"fmt"
	"jethur/grc/logs/services"

	masterDTOs "jethur/grc/masters/dtos"
	"jethur/grc/masters/repositories"
	"jethur/grc/models"
	"jethur/utility"
)

type ChecklistQuestionServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, page int, pageSize int, filters map[string][]string) (*masterDTOs.ChecklistQuestionListResponseDTO, error)
	GetByID(id uint) (*models.ChecklistQuestion, error)
	Create(currentUserID uint, checklistQuestion *masterDTOs.ChecklistQuestionFormDTO) (uint, error)
	Update(id uint, currentUserID uint, checklistQuestion *masterDTOs.ChecklistQuestionFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

type ChecklistQuestionService struct {
	repo           repositories.ChecklistQuestionRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
	masterSvc      QuestionGroupServiceInterface
	checklistSvc   ChecklistServiceInterface
}

func NewChecklistQuestionService(
	repo repositories.ChecklistQuestionRepositoryInterface,
	activityLogSvc services.ActivityLogServiceInterface,
	masterSvc QuestionGroupServiceInterface,
	checklistSvc ChecklistServiceInterface,
) ChecklistQuestionServiceInterface {
	return &ChecklistQuestionService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
		masterSvc:      masterSvc,
		checklistSvc:   checklistSvc,
	}
}

func (cqs *ChecklistQuestionService) GetAll(sort string, sortBy string, keyword string, page int, pageSize int, filters map[string][]string) (*masterDTOs.ChecklistQuestionListResponseDTO, error) {
	return cqs.repo.GetAll(sort, sortBy, keyword, page, pageSize, filters)
}

func (cqs ChecklistQuestionService) GetByID(id uint) (*models.ChecklistQuestion, error) {
	return cqs.repo.GetByID(id)
}

func (cqs *ChecklistQuestionService) Create(currentUserID uint, checklistQuestionFormDTO *masterDTOs.ChecklistQuestionFormDTO) (uint, error) {

	if err := cqs.validateChecklistQuestionFormDTO(checklistQuestionFormDTO, nil); err != nil {
		return 0, err
	}

	checklistQuestion := &models.ChecklistQuestion{
		ChecklistID:     checklistQuestionFormDTO.ChecklistID,
		QuestionGroupID: checklistQuestionFormDTO.QuestionGroupID,
		QuestionEn:      checklistQuestionFormDTO.QuestionEn,
		QuestionAr:      checklistQuestionFormDTO.QuestionAr,
	}

	id, err := cqs.repo.Create(checklistQuestion)
	if err != nil {
		return 0, err
	}

	go func() {
		err := cqs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ChecklistQuestion{}.TableName(),
			id,
			checklistQuestionFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ChecklistQuestion{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (cqs *ChecklistQuestionService) Update(id uint, currentUserID uint, checklistQuestionFormDTO *masterDTOs.ChecklistQuestionFormDTO) error {
	existingChecklistQuestion, err := cqs.repo.GetByID(id)

	if err != nil {
		return err
	}

	if err := cqs.validateChecklistQuestionFormDTO(checklistQuestionFormDTO, &id); err != nil {
		return err
	}

	existingChecklistQuestion.ChecklistID = checklistQuestionFormDTO.ChecklistID
	existingChecklistQuestion.QuestionGroupID = checklistQuestionFormDTO.QuestionGroupID
	existingChecklistQuestion.QuestionEn = checklistQuestionFormDTO.QuestionEn
	existingChecklistQuestion.QuestionAr = checklistQuestionFormDTO.QuestionAr

	err = cqs.repo.Update(id, existingChecklistQuestion)

	if err != nil {
		return err
	}

	return nil

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
				models.DeleteAction, // Make sure this action constant is defined appropriately
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

func (cqs *ChecklistQuestionService) validateChecklistQuestionFormDTO(checklistQuestionFormDTO *masterDTOs.ChecklistQuestionFormDTO, id *uint) error {
	if _, err := cqs.checklistSvc.GetByID(checklistQuestionFormDTO.ChecklistID); err != nil {
		return utility.NewValidationError(fmt.Errorf("review checklist does not exist"), "not_exist", "checklist")
	}

	if _, err := cqs.masterSvc.GetByID(checklistQuestionFormDTO.QuestionGroupID); err != nil {
		return utility.NewValidationError(fmt.Errorf("question group does not exist"), "not_exist", "question group")
	}

	if id == nil {
		if exists, err := cqs.repo.ExistsByTitle(checklistQuestionFormDTO.QuestionEn, checklistQuestionFormDTO.QuestionAr, checklistQuestionFormDTO.ChecklistID, checklistQuestionFormDTO.QuestionGroupID); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the question must be unique"), "unique", "question")
		}

	} else {
		if exists, err := cqs.repo.ExistsByTitleExcludingID(*id, checklistQuestionFormDTO.QuestionEn, checklistQuestionFormDTO.QuestionAr, checklistQuestionFormDTO.ChecklistID, checklistQuestionFormDTO.QuestionGroupID); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the question must be unique"), "unique", "question")
		}
	}

	return nil
}
