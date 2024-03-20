package services

import (
	"fmt"
	logModels "jethur/grc/logs/models"
	"jethur/grc/logs/services"
	checklistQuestionDTOs "jethur/grc/masters/dtos"
	"jethur/grc/masters/models"
	"jethur/grc/masters/repositories"
	"jethur/utility"
	"log"
)

type ChecklistQuestionServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, page int, pageSize int) (*checklistQuestionDTOs.ChecklistQuestionListResponseDTO, error)
	GetByID(id uint) (*models.ChecklistQuestion, error)
	Create(currentUserID uint, checklistQuestion *checklistQuestionDTOs.ChecklistQuestionFormDTO) (uint, error)
	Update(id uint, currentUserID uint, checklistQuestion *checklistQuestionDTOs.ChecklistQuestionFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type ChecklistQuestionService struct {
	repo           repositories.ChecklistQuestionRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
	checklistRepo  repositories.ChecklistRepositoryInterface
}

func NewChecklistQuestionService(
	repo repositories.ChecklistQuestionRepositoryInterface,
	activityLogSvc services.ActivityLogServiceInterface,
	checklistRepo repositories.ChecklistRepositoryInterface,
) ChecklistQuestionServiceInterface {
	return &ChecklistQuestionService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
		checklistRepo:  checklistRepo,
	}
}

func (cqs *ChecklistQuestionService) GetAll(sort string, sortBy string, keyword string, page int, pageSize int) (*checklistQuestionDTOs.ChecklistQuestionListResponseDTO, error) {
	return cqs.repo.GetAll(sort, sortBy, keyword, page, pageSize)
}

func (cqs ChecklistQuestionService) GetByID(id uint) (*models.ChecklistQuestion, error) {
	return cqs.repo.GetByID(id)
}

func (cqs *ChecklistQuestionService) Create(currentUserID uint, checklistQuestionFormDTO *checklistQuestionDTOs.ChecklistQuestionFormDTO) (uint, error) {

	// Validation for checklist
	checklist, err := cqs.checklistRepo.GetByID(checklistQuestionFormDTO.ChecklistID)
	if err != nil {
		return 0, utility.NewValidationError(fmt.Errorf("review checklist does not exist"), "not_exist", "checklist")
	}

	if exists, err := cqs.repo.ExistsByTitle(checklistQuestionFormDTO.QuestionEn, checklistQuestionFormDTO.QuestionAr, checklistQuestionFormDTO.ChecklistID); err != nil {
		return 0, err
	} else if exists {
		return 0, utility.NewValidationError(fmt.Errorf("the question must be unique"), "unique", "question")
	}

	if existsInLibrary, err := cqs.repo.ExistsQuestionLibrary(checklistQuestionFormDTO.QuestionLibraryID); err != nil {
		return 0, err
	} else if existsInLibrary {
		return 0, utility.NewValidationError(fmt.Errorf("the question already exists in this checklist"), "unique", "question")
	}
	_ = checklist

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

func (cqs *ChecklistQuestionService) Update(id uint, currentUserID uint, checklistQuestionFormDTO *checklistQuestionDTOs.ChecklistQuestionFormDTO) error {
	existingChecklistQuestion, err := cqs.repo.GetByID(id)

	if err != nil {
		return err
	}

	// Validation for checklist
	checklist, err := cqs.checklistRepo.GetByID(checklistQuestionFormDTO.ChecklistID)
	if err != nil {
		return utility.NewValidationError(fmt.Errorf("review checklist does not exist"), "not_exist", "checklist")
	}

	if checklistQuestionFormDTO.QuestionEn != existingChecklistQuestion.QuestionEn || checklistQuestionFormDTO.QuestionAr != existingChecklistQuestion.QuestionAr || checklistQuestionFormDTO.ChecklistID != existingChecklistQuestion.ChecklistID {
		if exists, err := cqs.repo.ExistsByTitleExcludingID(id, checklistQuestionFormDTO.QuestionEn, checklistQuestionFormDTO.QuestionAr, checklistQuestionFormDTO.ChecklistID); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the question must be unique"), "unique", "question")
		}
	}

	if existsInLibrary, err := cqs.repo.ExistsQuestionLibraryExcludingID(id, checklistQuestionFormDTO.QuestionLibraryID); err != nil {
		return err
	} else if existsInLibrary {
		return utility.NewValidationError(fmt.Errorf("the question already exists in this checklist"), "unique", "question")
	}

	_ = checklist

	err = cqs.repo.Update(id, checklistQuestionFormDTO)

	if err != nil {
		return err
	}

	fieldsToCompare := []string{
		"ChecklistID",
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
