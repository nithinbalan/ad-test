package services

import (
	"fmt"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/models"
	"jethur/grc/compliance/repositories"
	logModels "jethur/grc/logs/models"
	logServices "jethur/grc/logs/services"
	"jethur/utility"
)

type SLACheckAnswerServiceInterface interface {
	CreateAnswer(slaID uint, questionID uint, currentUserID uint, slaCheckAnswerFormDTO *complianceDTOs.SLACheckAnswerFormDTO) (uint, error)
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

func (scas *SLACheckAnswerService) CreateAnswer(slaID uint, questionID uint, currentUserID uint, slaCheckAnswerFormDTO *complianceDTOs.SLACheckAnswerFormDTO) (uint, error) {
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
		return 0, err
	}

	fmt.Println(answer.SLACheckQuestion)
	// Update completed count of the SLACheckQuestionGroup
	err = scas.slaCheckQuestionSvc.IncrementCompletedCount(slaCheckQuestion.SLACheckQuestionGroupID)
	if err != nil {
		return 0, err
	}

	go func() {
		err := scas.activityLogSvc.LogActivity(
			logModels.CreateAction,
			models.SLACheckAnswer{}.TableName(),
			answer.ID,
			answer,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    logModels.CreateAction,
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

func (crsus *SLACheckAnswerService) validateSlaCheckAnswerFormDTO(slaCheckAnswerFormDTO *complianceDTOs.SLACheckAnswerFormDTO) error {

	if _, err := crsus.complianceStatusSvc.GetByID(slaCheckAnswerFormDTO.ComplianceStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("compliance status does not exist"), "not_exist", "compliance status")
	}

	return nil
}
