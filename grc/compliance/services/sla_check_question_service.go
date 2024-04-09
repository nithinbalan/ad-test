package services

import (
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/models"
	"jethur/grc/compliance/repositories"
	logModels "jethur/grc/logs/models"
	logServices "jethur/grc/logs/services"
	"jethur/utility"
)

type SLACheckQuestionServiceInterface interface {
	// GetAll(slaID uint, checkID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.SLACheckQuestionListDTO, error)
	// ChooseQuestion(ids []uint, currentUserID uint) error
	Create(slaID uint, checkID uint, currentUserID uint, slaCheckQuestionFormDTO *complianceDTOs.SLACheckQuestionFormDTO) (uint, error)
	Update(slaID uint, checkID uint, id uint, currentUserID uint, slaCheckQuestionFormDTO *complianceDTOs.SLACheckQuestionFormDTO) error
	// Delete(slaID uint, checkID uint, id uint, currentUserID uint) error
}

// Service and implementations
type SLACheckQuestionService struct {
	repo           repositories.SLACheckQuestionRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
	slaRegisterSvc SLARegisterServiceInterface
}

func NewSLACheckQuestionService(
	repo repositories.SLACheckQuestionRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	slaRegisterSvc SLARegisterServiceInterface,
) SLACheckQuestionServiceInterface {
	return &SLACheckQuestionService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
		slaRegisterSvc: slaRegisterSvc,
	}
}

// func (sqs *SLACheckQuestionService) GetAll(slaID uint, checkID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.SLACheckQuestionListDTO, error) {
// 	// if _, err := sqs.complianceRequirementSvc.GetByID(slaID, checkID); err != nil {
// 	// 	return nil, err
// 	// }

// 	return sqs.repo.GetAll(checkID, sort, sortBy, keyword)
// }

// func (sqs SLACheckQuestionService) ChooseQuestion(slaID uint, checkID uint, id uint) (*models.SLACheckQuestion, error) {
// 	if _, err := sqs.complianceRequirementSvc.GetByID(slaID, checkID); err != nil {
// 		return nil, err
// 	}

// 	slaCheckQuestion, err := sqs.repo.GetByID(id)

// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, utility.NewResourceNotFoundError(err, id)
// 		}
// 		return nil, err
// 	}
// 	return slaCheckQuestion, nil
// }

func (sqs *SLACheckQuestionService) Create(slaID uint, checkID uint, currentUserID uint, slaCheckQuestionFormDTO *complianceDTOs.SLACheckQuestionFormDTO) (uint, error) {
	if _, err := sqs.slaRegisterSvc.GetByID(slaID); err != nil {
		return 0, err
	}

	slaCheckQuestionGroup := &models.SLACheckQuestionGroup{
		SLACheckID: checkID,
		TitleEn:    slaCheckQuestionFormDTO.QuestionGroupTitleEn,
		TitleAr:    slaCheckQuestionFormDTO.QuestionGroupTitleAr,
		CreatedBy:  currentUserID,
	}

	slaCheckQuestionGroupId, err := sqs.repo.CreateQuestionGroup(slaCheckQuestionGroup)
	if err != nil {
		return 0, err
	}

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
