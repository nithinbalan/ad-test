package services

import (
	"errors"
	"fmt"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type SLACheckAnswerServiceInterface interface {
	GetByID(slaID uint, questionID uint, answerID uint) (*models.SLACheckAnswer, error)
	CreateAnswer(slaID uint, questionID uint, currentUserID uint, slaCheckAnswerFormDTO *complianceDTOs.SLACheckAnswerFormDTO) (uint, error)
	UpdateAnswer(slaID uint, questionID uint, answerID uint, currentUserID uint, slaCheckAnswerFormDTO *complianceDTOs.SLACheckAnswerFormDTO) error
	DeleteAnswer(slaID uint, questionID uint, answerID uint, currentUserID uint) (uint, error)
	AssessmentSummary(slaID uint) (*complianceDTOs.AssessmentSummaryDTO, error)
}

// Service and implementations
type SLACheckAnswerService struct {
	repo                repositories.SLACheckAnswerRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	slaRegisterSvc      SLARegisterServiceInterface
	complianceStatusSvc ComplianceStatusServiceInterface
	slaCheckQuestionSvc SLACheckQuestionServiceInterface
}

func NewSLACheckAnswerService(
	repo repositories.SLACheckAnswerRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	slaRegisterSvc SLARegisterServiceInterface,
	complianceStatusSvc ComplianceStatusServiceInterface,
	slaCheckQuestionSvc SLACheckQuestionServiceInterface,
) SLACheckAnswerServiceInterface {
	return &SLACheckAnswerService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		slaRegisterSvc:      slaRegisterSvc,
		complianceStatusSvc: complianceStatusSvc,
		slaCheckQuestionSvc: slaCheckQuestionSvc,
	}
}

func (scas SLACheckAnswerService) GetByID(slaID uint, questionID uint, answerID uint) (*models.SLACheckAnswer, error) {
	if _, err := scas.slaRegisterSvc.GetByID(slaID); err != nil {
		return nil, err
	}

	slaCheckAnswer, err := scas.repo.GetByID(answerID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, answerID)
		}
		return nil, err
	}
	return slaCheckAnswer, nil
}

func (scas *SLACheckAnswerService) CreateAnswer(slaID uint, questionID uint, currentUserID uint, slaCheckAnswerFormDTO *complianceDTOs.SLACheckAnswerFormDTO) (uint, error) {
	tx := scas.repo.BeginTransaction()

	if _, err := scas.slaRegisterSvc.GetByID(slaID); err != nil {
		return 0, err
	}

	// Fetch SLACheckQuestion from the database
	slaCheckQuestion, err := scas.slaCheckQuestionSvc.GetQuestionByID(questionID)
	if err != nil {
		return 0, err
	}

	if err := scas.validateSlaCheckAnswerFormDTO(slaCheckAnswerFormDTO); err != nil {
		return 0, err
	}

	slaCheckAnswer := &models.SLACheckAnswer{
		SLACheckQuestionID: questionID,
		ComplianceStatusID: slaCheckAnswerFormDTO.ComplianceStatusID,
		Comment:            slaCheckAnswerFormDTO.Comment,
		CreatedBy:          currentUserID,
	}

	answer, err := scas.repo.Create(slaCheckAnswer)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Update completed count of the SLACheckQuestionGroup
	err = scas.slaCheckQuestionSvc.IncrementCompletedCount(slaCheckQuestion.SLACheckQuestionGroupID, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	go func() {
		err := scas.activityLogSvc.LogActivity(
			models.CreateAction,
			models.SLACheckAnswer{}.TableName(),
			answer.ID,
			answer,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.SLACheckAnswer{}.TableName(),
				"recordID":  answer.ID,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return slaCheckAnswer.ID, nil
}

func (scas *SLACheckAnswerService) UpdateAnswer(slaID uint, questionID uint, answerID uint, currentUserID uint, slaCheckAnswerFormDTO *complianceDTOs.SLACheckAnswerFormDTO) error {
	if _, err := scas.slaRegisterSvc.GetByID(slaID); err != nil {
		return err
	}

	if _, err := scas.slaCheckQuestionSvc.GetQuestionByID(questionID); err != nil {
		return err
	}

	existingAnswer, err := scas.GetByID(slaID, questionID, answerID)
	if err != nil {
		return err
	}

	existingAnswer.ComplianceStatusID = slaCheckAnswerFormDTO.ComplianceStatusID
	existingAnswer.Comment = slaCheckAnswerFormDTO.Comment
	existingAnswer.UpdatedBy = &currentUserID

	// Update the SLACheckQuestionGroup
	err = scas.repo.Update(existingAnswer)
	if err != nil {
		return err
	}

	// Log the activity
	go func() {
		err := scas.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.SLACheckAnswer{}.TableName(),
			existingAnswer.ID,
			existingAnswer,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.SLACheckAnswer{}.TableName(),
				"recordID":  existingAnswer.ID,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (scqs *SLACheckAnswerService) DeleteAnswer(slaID uint, questionID uint, answerID uint, currentUserID uint) (uint, error) {
	if _, err := scqs.slaRegisterSvc.GetByID(slaID); err != nil {
		return 0, err
	}

	if _, err := scqs.slaCheckQuestionSvc.GetQuestionByID(questionID); err != nil {
		return 0, err
	}
	var failedId uint
	var successfulId uint

	answer, err := scqs.GetByID(slaID, questionID, answerID)
	if err != nil {
		return 0, err
	}

	err = scqs.repo.Delete(answerID)
	if err != nil {
		failedId = answerID
	} else {
		successfulId = answerID
	}

	if successfulId != 0 {
		go func() {
			err := scqs.activityLogSvc.LogActivity(
				models.DeleteAction,
				models.SLACheckAnswer{}.TableName(),
				successfulId,
				answer,
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

func (scas *SLACheckAnswerService) validateSlaCheckAnswerFormDTO(slaCheckAnswerFormDTO *complianceDTOs.SLACheckAnswerFormDTO) error {

	if _, err := scas.complianceStatusSvc.GetByID(slaCheckAnswerFormDTO.ComplianceStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("compliance status does not exist"), "not_exist", "compliance status")
	}

	return nil
}

func (scas SLACheckAnswerService) AssessmentSummary(slaID uint) (*complianceDTOs.AssessmentSummaryDTO, error) {
	if _, err := scas.slaRegisterSvc.GetByID(slaID); err != nil {
		return nil, err
	}

	assessmentSummary, err := scas.repo.AssessmentSummary(slaID)

	if err != nil {
		return nil, err
	}
	return assessmentSummary, nil
}
