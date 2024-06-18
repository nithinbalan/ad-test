package services

import (
	"errors"
	"fmt"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"
	"path/filepath"
	"strconv"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type SLACheckAnswerServiceInterface interface {
	GetByID(slaID uint, questionID uint, answerID uint) (*models.SLACheckAnswer, error)
	GetFullDetailsByID(slaID uint, questionID uint, answerID uint) (*complianceDTOs.SLACheckAnswerDetailDTO, error)
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

func (scas SLACheckAnswerService) GetFullDetailsByID(slaID uint, questionID uint, answerID uint) (*complianceDTOs.SLACheckAnswerDetailDTO, error) {
	if _, err := scas.slaRegisterSvc.GetByID(slaID); err != nil {
		return nil, err
	}

	var slaCheckAnswerDTO complianceDTOs.SLACheckAnswerDetailDTO

	answer, err := scas.repo.GetFullDetailsByID(answerID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, answerID)
		}
		return nil, err
	}

	err = mapstructure.Decode(answer, &slaCheckAnswerDTO)
	if err != nil {
		return nil, err
	}

	return &slaCheckAnswerDTO, nil
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
		CreatedByID:        currentUserID,
	}

	answer, err := scas.repo.Create(slaCheckAnswer)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.SLACheckAnswerDocument

	for _, doc := range slaCheckAnswerFormDTO.Documents {

		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL("sla-register", strconv.FormatUint(uint64(slaID), 10), "questions", strconv.FormatUint(uint64(questionID), 10), "answers", encryptedFilename)

		document := &models.SLACheckAnswerDocument{
			SLACheckAnswerID: answer.ID,
			URL:              url,
			Name:             doc.FileName,
			Size:             doc.Size,
			Token:            encryptedFilename,
			OriginalName:     doc.OriginalName,
			Type:             doc.Type,
			CreatedByID:      currentUserID,
		}

		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := scas.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		//move file to storage
		for _, doc := range documents {
			storagePath := filepath.Join("compliance", "sla-register", strconv.FormatUint(uint64(slaID), 10), "questions", strconv.FormatUint(uint64(questionID), 10), "answers", strconv.FormatUint(uint64(answer.ID), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}

		}
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
	tx := scas.repo.BeginTransaction()

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
	existingAnswer.UpdatedByID = &currentUserID

	// Update the answer
	err = scas.repo.Update(existingAnswer)
	if err != nil {
		return err
	}

	var documents []*models.SLACheckAnswerDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.SLACheckAnswerDocument

	for _, doc := range slaCheckAnswerFormDTO.Documents {

		if doc.IsNewFile != nil && *doc.IsNewFile {
			// Encrypt the filename
			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)

			if err != nil {
				tx.Rollback()
				return err
			}

			//genereate URL
			url := utility.GenerateDocumentURL("sla-register", strconv.FormatUint(uint64(slaID), 10), "questions", strconv.FormatUint(uint64(questionID), 10), "answers", encryptedFileName)
			// Create a new Control test Document instance
			document := &models.SLACheckAnswerDocument{
				SLACheckAnswerID: answerID,
				URL:              url,
				Name:             doc.FileName,
				Size:             doc.Size,
				Token:            encryptedFileName,
				OriginalName:     doc.OriginalName,
				Type:             doc.Type,
				CreatedByID:      currentUserID,
			}

			// Append the document to the slice
			documents = append(documents, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			if doc.ID != nil {
				existingDocument, err := scas.repo.GetDocumentByID(*doc.ID)
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return utility.NewResourceNotFoundError(err, doc.ID)
					}
					return err
				}

				documentsToDelete = append(documentsToDelete, existingDocument)
				deleteDocumentIds = append(deleteDocumentIds, &existingDocument.ID)
			}
		}
	}

	if len(documents) > 0 {
		if err := scas.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIds) > 0 {
		if err := scas.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range documents {
		storagePath := filepath.Join("compliance", "sla-register", strconv.FormatUint(uint64(slaID), 10), "questions", strconv.FormatUint(uint64(questionID), 10), "answers", strconv.FormatUint(uint64(answerID), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("compliance", "sla-register", strconv.FormatUint(uint64(slaID), 10), "questions", strconv.FormatUint(uint64(questionID), 10), "answers", strconv.FormatUint(uint64(answerID), 10), "documents")
		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
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

// func (scas *SLACheckAnswerService) DeleteAnswer(slaID uint, questionID uint, answerID uint, currentUserID uint) (uint, error) {
// 	if _, err := scas.slaRegisterSvc.GetByID(slaID); err != nil {
// 		return 0, err
// 	}

// 	if _, err := scas.slaCheckQuestionSvc.GetQuestionByID(questionID); err != nil {
// 		return 0, err
// 	}
// 	var failedId uint
// 	var successfulId uint

// 	answer, err := scas.GetByID(slaID, questionID, answerID)
// 	if err != nil {
// 		return 0, err
// 	}

// 	err = scas.repo.Delete(answerID)
// 	if err != nil {
// 		failedId = answerID
// 	} else {
// 		successfulId = answerID
// 	}

// 	if successfulId != 0 {
// 		go func() {
// 			err := scas.activityLogSvc.LogActivity(
// 				models.DeleteAction,
// 				models.SLACheckAnswer{}.TableName(),
// 				successfulId,
// 				answer,
// 				currentUserID,
// 				nil,
// 			)
// 			if err != nil {
// 				fmt.Printf("Failed to log activity: %v\n", err)
// 			}
// 		}()
// 	}

// 	return failedId, nil
// }

func (scas *SLACheckAnswerService) DeleteAnswer(slaID uint, questionID uint, answerID uint, currentUserID uint) (uint, error) {
	// tx := scas.repo.BeginTransaction()

	_, err := scas.slaRegisterSvc.GetByID(slaID)
	if err != nil {
		return 0, err
	}

	var failedId uint
	var successfulId uint
	var deleteDocumentIds []*uint

	answer, err := scas.GetFullDetailsByID(slaID, questionID, answerID)
	if err != nil {
		return 0, err
	}

	if answer.Documents != nil && len(answer.Documents) > 0 {
		for _, document := range answer.Documents {
			deleteDocumentIds = append(deleteDocumentIds, document.ID)
		}

		if len(deleteDocumentIds) > 0 {
			if err := scas.repo.DeleteDocuments(deleteDocumentIds); err != nil {
				// tx.Rollback()
				return 0, err
			}
		}
	}

	if err := scas.repo.Delete(answerID); err != nil {
		failedId = answerID
		// tx.Rollback()
	} else {
		successfulId = answerID
		fmt.Println(answer.Documents)
		for _, doc := range answer.Documents {
			storagePath := filepath.Join("compliance", "sla-register", strconv.FormatUint(uint64(slaID), 10), "questions", strconv.FormatUint(uint64(questionID), 10), "answers", strconv.FormatUint(uint64(answerID), 10), "documents")
			if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
				return 0, err
				// tx.Rollback()
			}
		}

	}

	// if successfulId > 0 {
	// 	for _, doc := range answer.Documents {
	// 		storagePath := filepath.Join("compliance", "sla-register", strconv.FormatUint(uint64(slaID), 10), "questions", strconv.FormatUint(uint64(questionID), 10), "answers", strconv.FormatUint(uint64(answerID), 10), "documents")
	// 		if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
	// 			tx.Rollback()
	// 		}
	// 	}

	// }

	// tx.Commit()

	go func() {
		err := scas.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.SLACheckAnswer{}.TableName(),
			successfulId,
			nil,
			currentUserID,
			nil,
		)
		if err != nil {
			fmt.Printf("Failed to log bulk delete activity: %v\n", err)
		}
	}()

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
