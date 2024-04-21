package services

import (
	"errors"
	"fmt"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/models"
	"jethur/grc/compliance/repositories"
	logModels "jethur/grc/logs/models"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/utility"

	"gorm.io/gorm"
)

type SLACheckQuestionServiceInterface interface {
	GetAllQuestionGroups(slaID uint, checkID uint) ([]*complianceDTOs.SLACheckQuestionGroupListDTO, error)
	UpdateQuestionGroup(slaID uint, checkID uint, groupID uint, currentUserID uint, slaCheckUpdateQuestionGroupDTO *complianceDTOs.SLACheckUpdateQuestionGroupDTO) error
	DeleteQuestionGroup(slaID uint, checkID uint, questionGroupID uint, currentUserID uint) (uint, error)
	GetQuestionsByGroupID(slaID uint, checkID uint, groupID uint) ([]models.SLACheckQuestion, error)
	GetQuestionByID(id uint) (*models.SLACheckQuestion, error)
	ChooseQuestions(slaID uint, checkID uint, currentUserID uint, checklistQuestionIds []uint) ([]uint, error)
	CreateQuestion(slaID uint, checkID uint, currentUserID uint, slaCheckQuestionFormDTO *complianceDTOs.SLACheckQuestionFormDTO) (uint, error)
	UpdateQuestion(slaID uint, checkID uint, groupID uint, currentUserID uint, slaCheckUpdateQuestionDTO *complianceDTOs.SLACheckUpdateQuestionDTO) error
	DeleteQuestion(slaID uint, checkId uint, questionId uint, currentUserID uint) (uint, error)
	IncrementCompletedCount(groupID uint) error
}

// Service and implementations
type SLACheckQuestionService struct {
	repo                 repositories.SLACheckQuestionRepositoryInterface
	activityLogSvc       logServices.ActivityLogServiceInterface
	slaRegisterSvc       SLARegisterServiceInterface
	checklistQuestionSvc masterServices.ChecklistQuestionServiceInterface
	questionGroupSvc     masterServices.QuestionGroupServiceInterface
}

func NewSLACheckQuestionService(
	repo repositories.SLACheckQuestionRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	slaRegisterSvc SLARegisterServiceInterface,
	checklistQuestionSvc masterServices.ChecklistQuestionServiceInterface,
	questionGroupSvc masterServices.QuestionGroupServiceInterface,
) SLACheckQuestionServiceInterface {
	return &SLACheckQuestionService{
		repo:                 repo,
		activityLogSvc:       activityLogSvc,
		slaRegisterSvc:       slaRegisterSvc,
		checklistQuestionSvc: checklistQuestionSvc,
		questionGroupSvc:     questionGroupSvc,
	}
}

func (scqs *SLACheckQuestionService) GetAllQuestionGroups(slaID uint, checkID uint) ([]*complianceDTOs.SLACheckQuestionGroupListDTO, error) {
	if _, err := scqs.slaRegisterSvc.GetByID(slaID); err != nil {
		return nil, err
	}
	return scqs.repo.GetAllQuestionGroups(checkID)
}

func (scqs *SLACheckQuestionService) UpdateQuestionGroup(slaID uint, checkID uint, groupID uint, currentUserID uint, slaCheckUpdateQuestionGroupDTO *complianceDTOs.SLACheckUpdateQuestionGroupDTO) error {
	if _, err := scqs.slaRegisterSvc.GetByID(slaID); err != nil {
		return err
	}
	// Fetch the SLACheckQuestionGroup by ID
	group, err := scqs.repo.GetQuestionGroupByID(groupID)
	if err != nil {
		return err
	}

	// Update the titles if provided
	if slaCheckUpdateQuestionGroupDTO.QuestionGroupTitleEn != "" {
		group.TitleEn = slaCheckUpdateQuestionGroupDTO.QuestionGroupTitleEn
	}
	if slaCheckUpdateQuestionGroupDTO.QuestionGroupTitleAr != "" {
		group.TitleAr = slaCheckUpdateQuestionGroupDTO.QuestionGroupTitleAr
	}

	// Update the SLACheckQuestionGroup
	err = scqs.repo.UpdateQuestionGroup(group)
	if err != nil {
		return err
	}

	// Log the activity
	go func() {
		err := scqs.activityLogSvc.LogActivity(
			logModels.UpdateAction,
			models.SLACheckQuestionGroup{}.TableName(),
			groupID,
			group,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    logModels.UpdateAction,
				"tableName": models.SLACheckQuestionGroup{}.TableName(),
				"recordID":  groupID,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (scqs *SLACheckQuestionService) DeleteQuestionGroup(slaID uint, checkID uint, questionGroupID uint, currentUserID uint) (uint, error) {
	if _, err := scqs.slaRegisterSvc.GetByID(slaID); err != nil {
		return 0, err
	}
	var failedQuestionIDs []uint
	var successfulQuestionIDs []uint
	var failedGroupID uint
	var successfulGroupID uint

	// Fetch the SLACheckQuestionGroup by ID
	questionGroup, err := scqs.repo.GetQuestionGroupByID(questionGroupID)
	if err != nil {
		return 0, err
	}

	// Fetch all questions under the question group
	questions, err := scqs.repo.GetQuestionsByGroupID(questionGroupID)
	if err != nil {
		return 0, err
	}

	// Delete each question
	for _, question := range questions {
		err = scqs.repo.DeleteQuestion(question.ID)
		if err != nil {
			failedQuestionIDs = append(failedQuestionIDs, question.ID)
		} else {
			successfulQuestionIDs = append(successfulQuestionIDs, question.ID)
		}
	}

	// If all questions were successfully deleted, delete the question group
	if len(failedQuestionIDs) == 0 {
		err = scqs.repo.DeleteQuestionGroup(questionGroupID)
		if err != nil {
			failedGroupID = questionGroupID
		} else {
			successfulGroupID = questionGroupID
		}
	}

	// Log the activity for successful deletion
	if len(successfulQuestionIDs) > 0 {
		go func() {
			err := scqs.activityLogSvc.LogBulkActivities(
				logModels.DeleteAction,
				models.SLACheckQuestion{}.TableName(),
				successfulQuestionIDs,
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

	if successfulGroupID != 0 {
		go func() {
			err := scqs.activityLogSvc.LogActivity(
				logModels.DeleteAction,
				models.SLACheckQuestionGroup{}.TableName(),
				successfulGroupID,
				questionGroup,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log activity: %v\n", err)
			}
		}()
	}

	return failedGroupID, nil
}

func (scqs *SLACheckQuestionService) GetQuestionsByGroupID(slaID uint, checkID uint, groupID uint) ([]models.SLACheckQuestion, error) {
	if _, err := scqs.slaRegisterSvc.GetByID(slaID); err != nil {
		return nil, err
	}
	return scqs.repo.GetQuestionsByGroupID(groupID)
}

func (scqs SLACheckQuestionService) GetQuestionByID(id uint) (*models.SLACheckQuestion, error) {
	slaCheckQuestion, err := scqs.repo.GetQuestionByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return slaCheckQuestion, nil
}

func (scqs SLACheckQuestionService) ChooseQuestions(slaID uint, checkID uint, currentUserID uint, checklistQuestionIds []uint) ([]uint, error) {
	if _, err := scqs.slaRegisterSvc.GetByID(slaID); err != nil {
		return nil, err
	}

	var failedIds []uint
	var successfulIds []uint
	for _, checklistQuestionID := range checklistQuestionIds {
		checklistQuestion, err := scqs.checklistQuestionSvc.GetByID(checklistQuestionID)
		if err != nil {
			return nil, err
		}

		// Check if a similar question already exists under a similar group
		existingQuestion, err := scqs.repo.GetSimilarQuestion(checkID, checklistQuestion.QuestionEn, checklistQuestion.QuestionAr)
		if err != nil {
			return nil, err
		}

		if existingQuestion != nil {
			continue
		}

		// Check if a group with the same title already exists for the given check ID
		existingGroup, err := scqs.repo.GetSimilarQuestionGroup(checklistQuestion.QuestionGroup.TitleEn, checklistQuestion.QuestionGroup.TitleAr, checkID)
		if err != nil {
			return nil, err
		}

		var slaCheckQuestionGroupId uint
		if existingGroup != nil {
			// Use the existing group ID
			slaCheckQuestionGroupId = existingGroup.ID
		} else {
			// Create a new SLACheckQuestionGroup
			slaCheckQuestionGroup := &models.SLACheckQuestionGroup{
				SLACheckID: checkID,
				TitleEn:    checklistQuestion.QuestionGroup.TitleEn,
				TitleAr:    checklistQuestion.QuestionGroup.TitleAr,
				CreatedBy:  currentUserID,
			}

			// Create the new group
			groupID, err := scqs.repo.CreateQuestionGroup(slaCheckQuestionGroup)
			if err != nil {
				return nil, err
			}
			slaCheckQuestionGroupId = groupID
		}

		// Create SLACheckQuestion
		slaCheckQuestion := &models.SLACheckQuestion{
			SLACheckQuestionGroupID: slaCheckQuestionGroupId,
			TitleEn:                 checklistQuestion.QuestionEn,
			TitleAr:                 checklistQuestion.QuestionAr,
			CreatedBy:               currentUserID,
		}

		// Update total count of the SLACheckQuestionGroup
		err = scqs.repo.IncrementTotalCount(slaCheckQuestionGroupId)
		if err != nil {
			return nil, err
		}

		_, err = scqs.repo.CreateQuestion(slaCheckQuestion)

		if err != nil {
			failedIds = append(failedIds, checklistQuestionID)
		} else {
			successfulIds = append(successfulIds, checklistQuestionID)
		}
	}

	if len(successfulIds) > 0 {
		go func() {
			err := scqs.activityLogSvc.LogBulkActivities(
				logModels.CreateAction,
				models.SLACheckQuestion{}.TableName(),
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

	return failedIds, nil
}

func (scqs *SLACheckQuestionService) CreateQuestion(slaID uint, checkID uint, currentUserID uint, slaCheckQuestionFormDTO *complianceDTOs.SLACheckQuestionFormDTO) (uint, error) {
	if _, err := scqs.slaRegisterSvc.GetByID(slaID); err != nil {
		return 0, err
	}

	var slaCheckQuestionGroupId uint
	if slaCheckQuestionFormDTO.QuestionGroupID != 0 {
		// Use the provided question group ID to fetch the title
		group, err := scqs.questionGroupSvc.GetByID(slaCheckQuestionFormDTO.QuestionGroupID)
		if err != nil {
			return 0, err
		}
		// Use the title from the existing group
		slaCheckQuestionFormDTO.QuestionGroupTitleEn = group.TitleEn
		slaCheckQuestionFormDTO.QuestionGroupTitleAr = group.TitleAr
	}
	// Create a new SLACheckQuestionGroup
	slaCheckQuestionGroup := &models.SLACheckQuestionGroup{
		SLACheckID: checkID,
		TitleEn:    slaCheckQuestionFormDTO.QuestionGroupTitleEn,
		TitleAr:    slaCheckQuestionFormDTO.QuestionGroupTitleAr,
		CreatedBy:  currentUserID,
	}

	groupID, err := scqs.repo.CreateQuestionGroup(slaCheckQuestionGroup)
	if err != nil {
		return 0, err
	}

	slaCheckQuestionGroupId = groupID

	slaCheckQuestion := &models.SLACheckQuestion{
		SLACheckQuestionGroupID: slaCheckQuestionGroupId,
		TitleEn:                 slaCheckQuestionFormDTO.QuestionTitleEn,
		TitleAr:                 slaCheckQuestionFormDTO.QuestionTitleAr,
		CreatedBy:               currentUserID,
	}

	slaCheckQuestionId, err := scqs.repo.CreateQuestion(slaCheckQuestion)
	if err != nil {
		return 0, err
	}

	// Update total count of the SLACheckQuestionGroup
	err = scqs.repo.IncrementTotalCount(slaCheckQuestionGroupId)
	if err != nil {
		return 0, err
	}

	go func() {
		err := scqs.activityLogSvc.LogActivity(
			logModels.CreateAction,
			models.SLACheckQuestion{}.TableName(),
			slaCheckQuestionId,
			slaCheckQuestion,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    logModels.CreateAction,
				"tableName": models.SLACheckQuestion{}.TableName(),
				"recordID":  slaCheckQuestionId,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return slaCheckQuestionGroupId, nil
}

func (scqs *SLACheckQuestionService) UpdateQuestion(slaID uint, checkID uint, questionID uint, currentUserID uint, slaCheckUpdateQuestionDTO *complianceDTOs.SLACheckUpdateQuestionDTO) error {
	if _, err := scqs.slaRegisterSvc.GetByID(slaID); err != nil {
		return err
	}
	// Fetch the SLACheckQuestion by ID
	question, err := scqs.GetQuestionByID(questionID)

	if err != nil {
		return err
	}

	// Update the titles if provided
	if slaCheckUpdateQuestionDTO.QuestionTitleEn != "" {
		question.TitleEn = slaCheckUpdateQuestionDTO.QuestionTitleEn
	}
	if slaCheckUpdateQuestionDTO.QuestionTitleAr != "" {
		question.TitleAr = slaCheckUpdateQuestionDTO.QuestionTitleAr
	}

	// Update the SLACheckQuestionGroup
	err = scqs.repo.UpdateQuestion(question)
	if err != nil {
		return err
	}

	// Log the activity
	go func() {
		err := scqs.activityLogSvc.LogActivity(
			logModels.UpdateAction,
			models.SLACheckQuestion{}.TableName(),
			questionID,
			question,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    logModels.UpdateAction,
				"tableName": models.SLACheckQuestion{}.TableName(),
				"recordID":  questionID,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (scqs *SLACheckQuestionService) DeleteQuestion(slaID uint, checkId uint, questionId uint, currentUserID uint) (uint, error) {
	if _, err := scqs.slaRegisterSvc.GetByID(slaID); err != nil {
		return 0, err
	}
	var failedId uint
	var successfulId uint

	question, err := scqs.GetQuestionByID(questionId)
	if err != nil {
		return 0, err
	}

	err = scqs.repo.DeleteQuestion(questionId)
	if err != nil {
		failedId = questionId
	} else {
		successfulId = questionId
	}

	if successfulId != 0 {
		go func() {
			err := scqs.activityLogSvc.LogActivity(
				logModels.DeleteAction,
				models.SLACheckQuestion{}.TableName(),
				successfulId,
				question,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log activity: %v\n", err)
			}
		}()
	}

	return failedId, nil
}

func (scqs *SLACheckQuestionService) IncrementCompletedCount(groupID uint) error {
	return scqs.repo.IncrementCompletedCount(groupID)
}
