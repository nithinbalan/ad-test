package services

import (
	"fmt"
	logModels "jethur/grc/logs/models"
	"jethur/grc/logs/services"
	masterDTOs "jethur/grc/masters/dtos"
	"jethur/grc/masters/models"
	"jethur/grc/masters/repositories"
	"jethur/utility"
	"log"
)

type ChecklistQuestionServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, page int, pageSize int) (*masterDTOs.ChecklistQuestionListResponseDTO, error)
	GetByID(id uint) (*models.ChecklistQuestion, error)
	Create(currentUserID uint, checklistQuestion *masterDTOs.ChecklistQuestionFormDTO) (uint, error)
	Update(id uint, currentUserID uint, checklistQuestion *masterDTOs.ChecklistQuestionFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type ChecklistQuestionService struct {
	repo               repositories.ChecklistQuestionRepositoryInterface
	activityLogSvc     services.ActivityLogServiceInterface
	masterSvc          QuestionGroupServiceInterface
	checklistSvc       ChecklistServiceInterface
	questionLibrarySvc QuestionLibraryServiceInterface
}

func NewChecklistQuestionService(
	repo repositories.ChecklistQuestionRepositoryInterface,
	activityLogSvc services.ActivityLogServiceInterface,
	masterSvc QuestionGroupServiceInterface,
	checklistSvc ChecklistServiceInterface,
	questionLibrarySvc QuestionLibraryServiceInterface,
) ChecklistQuestionServiceInterface {
	return &ChecklistQuestionService{
		repo:               repo,
		activityLogSvc:     activityLogSvc,
		masterSvc:          masterSvc,
		checklistSvc:       checklistSvc,
		questionLibrarySvc: questionLibrarySvc,
	}
}

func (cqs *ChecklistQuestionService) GetAll(sort string, sortBy string, keyword string, page int, pageSize int) (*masterDTOs.ChecklistQuestionListResponseDTO, error) {
	return cqs.repo.GetAll(sort, sortBy, keyword, page, pageSize)
}

func (cqs ChecklistQuestionService) GetByID(id uint) (*models.ChecklistQuestion, error) {
	return cqs.repo.GetByID(id)
}

func (cqs *ChecklistQuestionService) Create(currentUserID uint, checklistQuestionFormDTO *masterDTOs.ChecklistQuestionFormDTO) (uint, error) {

	if err := cqs.validateChecklistQuestionFormDTO(checklistQuestionFormDTO, nil); err != nil {
		return 0, err
	}

	id, err := cqs.repo.Create(checklistQuestionFormDTO)
	if err != nil {
		return 0, err
	}

	go func() {
		err := cqs.activityLogSvc.LogActivity(
			logModels.CreateAction,
			models.ChecklistQuestion{}.TableName(),
			id,
			checklistQuestionFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			// @todo - save to log file
			log.Printf("Failed to log activity: %v", err)
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

	err = cqs.repo.Update(id, checklistQuestionFormDTO)

	if err != nil {
		return err
	}

	fieldsToCompare := []string{
		"ChecklistID",
		"QuestionGroupID",
		"QuestionLibraryID",
		"QuestionEn",
		"QuestionAr"}
	changes, err := utility.CompareFields(existingChecklistQuestion, checklistQuestionFormDTO, fieldsToCompare)

	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := cqs.activityLogSvc.LogActivity(
			logModels.UpdateAction,
			models.ChecklistQuestion{}.TableName(),
			id,
			changedDataJSON,
			currentUserID,
			nil,
		)
		if err != nil {
			// @todo - save to log file
			fmt.Printf("Failed to log bulk delete activity: %v\n", err)
		}
	}()
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
				logModels.DeleteAction, // Make sure this action constant is defined appropriately
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

	if checklistQuestionFormDTO.QuestionLibraryID != 0 {
		if _, err := cqs.questionLibrarySvc.GetByID(checklistQuestionFormDTO.QuestionLibraryID); err != nil {
			return utility.NewValidationError(fmt.Errorf("question library does not exist"), "not_exist", "question library")
		}
	}

	if id == nil {
		if exists, err := cqs.repo.ExistsByTitle(checklistQuestionFormDTO.QuestionEn, checklistQuestionFormDTO.QuestionAr, checklistQuestionFormDTO.ChecklistID, checklistQuestionFormDTO.QuestionGroupID); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the question must be unique"), "unique", "question")
		}

		if existsInLibrary, err := cqs.repo.ExistsQuestionLibrary(checklistQuestionFormDTO.QuestionLibraryID, checklistQuestionFormDTO.QuestionEn, checklistQuestionFormDTO.QuestionAr); err != nil {
			return err
		} else if existsInLibrary {
			return utility.NewValidationError(fmt.Errorf("the question already exists in this checklist"), "unique", "question")
		}

	} else {
		if exists, err := cqs.repo.ExistsByTitleExcludingID(*id, checklistQuestionFormDTO.QuestionEn, checklistQuestionFormDTO.QuestionAr, checklistQuestionFormDTO.ChecklistID, checklistQuestionFormDTO.QuestionGroupID); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the question must be unique"), "unique", "question")
		}

		if existsInLibrary, err := cqs.repo.ExistsQuestionLibraryExcludingID(*id, checklistQuestionFormDTO.QuestionEn, checklistQuestionFormDTO.QuestionAr, checklistQuestionFormDTO.QuestionLibraryID); err != nil {
			return err
		} else if existsInLibrary {
			return utility.NewValidationError(fmt.Errorf("the question already exists in this checklist"), "unique", "question")
		}

	}

	return nil
}
