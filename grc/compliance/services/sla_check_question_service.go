package services

import (
	"fmt"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/models"
	"jethur/grc/compliance/repositories"
	logModels "jethur/grc/logs/models"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/utility"
)

type SLACheckQuestionServiceInterface interface {
	// GetAll(slaID uint, checkID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.SLACheckQuestionListDTO, error)
	ChooseQuestion(slaID uint, checkID uint, currentUserID uint, checklistQuestionIds []uint) ([]uint, error)
	Create(slaID uint, checkID uint, currentUserID uint, slaCheckQuestionFormDTO *complianceDTOs.SLACheckQuestionFormDTO) (uint, error)
	// Update(slaID uint, checkID uint, id uint, currentUserID uint, slaCheckQuestionFormDTO *complianceDTOs.SLACheckQuestionFormDTO) error
	// Delete(slaID uint, checkID uint, id uint, currentUserID uint) error
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

// func (sqs *SLACheckQuestionService) GetAll(slaID uint, checkID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.SLACheckQuestionListDTO, error) {
// 	// if _, err := sqs.complianceRequirementSvc.GetByID(slaID, checkID); err != nil {
// 	// 	return nil, err
// 	// }

// 	return sqs.repo.GetAll(checkID, sort, sortBy, keyword)
// }

func (sqs SLACheckQuestionService) ChooseQuestion(slaID uint, checkID uint, currentUserID uint, checklistQuestionIds []uint) ([]uint, error) {
	if _, err := sqs.slaRegisterSvc.GetByID(slaID); err != nil {
		return nil, err
	}

	var failedIds []uint
	var successfulIds []uint
	for _, checklistQuestionID := range checklistQuestionIds {
		checklistQuestion, err := sqs.checklistQuestionSvc.GetByID(checklistQuestionID)
		if err != nil {
			return nil, err
		}

		// Check if a similar question already exists under a similar group
		existingQuestion, err := sqs.repo.GetSimilarQuestion(checkID, checklistQuestion)
		if err != nil {
			return nil, err
		}

		if existingQuestion != nil {
			// Skip adding the question
			successfulIds = append(successfulIds, checklistQuestionID)
			continue
		}

		// Check if a group with the same title already exists for the given check ID
		existingGroup, err := sqs.repo.GetQuestionGroupByTitleAndCheckID(checklistQuestion.QuestionGroup.TitleEn, checkID)
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
			groupID, err := sqs.repo.CreateQuestionGroup(slaCheckQuestionGroup)
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

		_, err = sqs.repo.CreateQuestion(slaCheckQuestion)

		if err != nil {
			failedIds = append(failedIds, checklistQuestionID)
		} else {
			successfulIds = append(successfulIds, checklistQuestionID)
		}
	}

	if len(successfulIds) > 0 {
		go func() {
			err := sqs.activityLogSvc.LogBulkActivities(
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

func (sqs *SLACheckQuestionService) Create(slaID uint, checkID uint, currentUserID uint, slaCheckQuestionFormDTO *complianceDTOs.SLACheckQuestionFormDTO) (uint, error) {
	if _, err := sqs.slaRegisterSvc.GetByID(slaID); err != nil {
		return 0, err
	}

	var slaCheckQuestionGroupId uint
	if slaCheckQuestionFormDTO.QuestionGroupID != 0 {
		// Use the provided question group ID to fetch the title
		group, err := sqs.questionGroupSvc.GetByID(slaCheckQuestionFormDTO.QuestionGroupID)
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

	groupID, err := sqs.repo.CreateQuestionGroup(slaCheckQuestionGroup)
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

	slaCheckQuestionId, err := sqs.repo.CreateQuestion(slaCheckQuestion)
	if err != nil {
		return 0, err
	}

	go func() {
		err := sqs.activityLogSvc.LogActivity(
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

func (sqs *SLACheckQuestionService) Update(slaID uint, checkID uint, id uint, currentUserID uint, slaCheckQuestionFormDTO *complianceDTOs.SLACheckQuestionFormDTO) error {
	if _, err := sqs.slaRegisterSvc.GetByID(slaID); err != nil {
		return err
	}

	// existingSLACheckQuestion, err := sqs.GetByID(slaID, checkID, id)

	// if err != nil {
	// 	return err
	// }

	// if err := sqs.validateSLACheckQuestionFormDTO(slaCheckQuestionFormDTO); err != nil {
	// 	return err
	// }

	// // Update SLACheckQuestion
	// existingSLACheckQuestion.ComplianceStatusID = slaCheckQuestionFormDTO.ComplianceStatusID
	// existingSLACheckQuestion.Comment = slaCheckQuestionFormDTO.Comment
	// existingSLACheckQuestion.UpdatedBy = &currentUserID

	// err = sqs.repo.Update(existingSLACheckQuestion)

	// if err != nil {
	// 	return err
	// }

	return nil
}

// func (sqs *SLACheckQuestionService) validateSLACheckQuestionFormDTO(slaCheckQuestionFormDTO *complianceDTOs.SLACheckQuestionFormDTO) error {

// 	if _, err := sqs.complianceStatusSvc.GetByID(slaCheckQuestionFormDTO.ComplianceStatusID); err != nil {
// 		return utility.NewValidationError(fmt.Errorf("compliance status does not exist"), "not_exist", "compliance status")
// 	}

// 	return nil
// }
