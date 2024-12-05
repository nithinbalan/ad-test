package services

import (
	"errors"
	"fmt"
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"path/filepath"
	"strconv"
	"time"

	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type AuditChecklistServiceInterface interface {
	GetAll(auditId uint, params *auditDTOs.AuditChecklistListQueryParamsDTO) ([]*auditDTOs.AuditChecklistListDTO, error)
	GetByID(auditId uint, id uint) (*models.AuditChecklist, error)
	GetFullDetailsByID(auditId uint, id uint) (*auditDTOs.AuditChecklistDetailDTO, error)
	GetDraftFindingsDetailsByChecklistID(auditID uint, checklistID uint) ([]*auditDTOs.AuditDraftFindingDetailDTO, error)
	Create(auditId uint, currentUserID uint, auditChecklistFormDTO *auditDTOs.AuditChecklistFormDTO) (uint, error)
	Update(auditId uint, id uint, currentUserID uint, auditChecklistFormDTO *auditDTOs.AuditChecklistFormDTO) error
	Delete(auditId uint, ids []uint, currentUserID uint) ([]uint, error)
	Answer(auditId uint, id uint, currentUserID uint, auditChecklistAnswerFormDTO *auditDTOs.AuditChecklistAnswerFormDTO) error
	DeleteAnswer(auditId uint, ids []uint, currentUserID uint) ([]uint, error)
}

type AuditChecklistService struct {
	repo                      repositories.AuditChecklistRepositoryInterface
	activityLogSvc            logServices.ActivityLogServiceInterface
	auditSvc                  AuditServiceInterface
	auditableItemSvc          AuditableItemServiceInterface
	conformanceSvc            masterServices.ConformanceServiceInterface
	statusSvc                 masterServices.StatusServiceInterface
	auditPlanAuditableItemSvc AuditPlanAuditableItemServiceInterface
}

func NewAuditChecklistService(
	repo repositories.AuditChecklistRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditSvc AuditServiceInterface,
	auditableItemSvc AuditableItemServiceInterface,
	conformanceSvc masterServices.ConformanceServiceInterface,
	statusSvc masterServices.StatusServiceInterface,
	auditPlanAuditableItemSvc AuditPlanAuditableItemServiceInterface,
) AuditChecklistServiceInterface {
	return &AuditChecklistService{
		repo:                      repo,
		activityLogSvc:            activityLogSvc,
		auditSvc:                  auditSvc,
		auditableItemSvc:          auditableItemSvc,
		conformanceSvc:            conformanceSvc,
		statusSvc:                 statusSvc,
		auditPlanAuditableItemSvc: auditPlanAuditableItemSvc,
	}
}

func (acs *AuditChecklistService) GetAll(auditId uint, params *auditDTOs.AuditChecklistListQueryParamsDTO) ([]*auditDTOs.AuditChecklistListDTO, error) {
	if _, err := acs.auditSvc.GetByID(auditId); err != nil {
		return nil, err
	}
	return acs.repo.GetAll(auditId, params)
}

func (acs *AuditChecklistService) GetByID(auditId uint, id uint) (*models.AuditChecklist, error) {
	if _, err := acs.auditSvc.GetByID(auditId); err != nil {
		return nil, err
	}

	auditChecklist, err := acs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return auditChecklist, nil
}

func (acs *AuditChecklistService) GetFullDetailsByID(auditId uint, id uint) (*auditDTOs.AuditChecklistDetailDTO, error) {
	if _, err := acs.auditSvc.GetByID(auditId); err != nil {
		return nil, err
	}
	var auditChecklitDetailDTO auditDTOs.AuditChecklistDetailDTO

	checklist, err := acs.repo.GetFullDetailsByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(checklist, &auditChecklitDetailDTO)
	if err != nil {
		return nil, err
	}

	//todo :- also add extra fields which required in details screen

	return &auditChecklitDetailDTO, nil
}

func (acs *AuditChecklistService) GetDraftFindingsDetailsByChecklistID(auditId uint, checklistID uint) ([]*auditDTOs.AuditDraftFindingDetailDTO, error) {
	if _, err := acs.GetByID(auditId, checklistID); err != nil {
		return nil, err
	}
	auditDraftFindings, err := acs.repo.GetDraftFindingsDetailsByChecklistID(checklistID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Return an empty array and no error if no records are found
			return []*auditDTOs.AuditDraftFindingDetailDTO{}, nil
		}
		return nil, err
	}

	// Convert the slice of AuditDraftFinding to a slice of AuditDraftFindingDetailDTO
	var auditDraftFindingDetailDTOs []*auditDTOs.AuditDraftFindingDetailDTO
	for _, finding := range auditDraftFindings {
		var dto auditDTOs.AuditDraftFindingDetailDTO
		err = mapstructure.Decode(finding, &dto)
		if err != nil {
			return nil, err
		}
		auditDraftFindingDetailDTOs = append(auditDraftFindingDetailDTOs, &dto)
	}

	// If there are no findings, return an empty slice
	if len(auditDraftFindingDetailDTOs) == 0 {
		return []*auditDTOs.AuditDraftFindingDetailDTO{}, nil
	}

	return auditDraftFindingDetailDTOs, nil
}

func (acs *AuditChecklistService) Create(auditId uint, currentUserID uint, auditChecklistFormDTO *auditDTOs.AuditChecklistFormDTO) (uint, error) {
	tx := acs.repo.BeginTransaction()

	if _, err := acs.auditSvc.GetByID(auditId); err != nil {
		return 0, err
	}

	if err := acs.validateAuditChecklistFormDTO(auditChecklistFormDTO); err != nil {
		tx.Rollback()
		return 0, err
	}
	status, err := acs.statusSvc.GetByType("pending")
	if err != nil {
		return 0, err
	}
	auditChecklist := &models.AuditChecklist{
		AuditID:         auditId,
		AuditableItemID: auditChecklistFormDTO.AuditableItemID,
		Question:        auditChecklistFormDTO.Question,
		StatusID:        status.ID,
		CreatedBYID:     currentUserID,
	}

	id, err := acs.repo.Create(auditChecklist, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	go func() {
		err := acs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditChecklist{}.TableName(),
			id,
			auditChecklistFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AuditChecklist{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (acs *AuditChecklistService) Update(auditId uint, id uint, currentUserID uint, auditChecklistFormDTO *auditDTOs.AuditChecklistFormDTO) error {
	tx := acs.repo.BeginTransaction()

	if _, err := acs.auditSvc.GetByID(auditId); err != nil {
		return err
	}

	existingAuditChecklist, err := acs.GetByID(auditId, id)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}

	if err := acs.validateAuditChecklistFormDTO(auditChecklistFormDTO); err != nil {
		tx.Rollback()
		return err
	}
	existingAuditChecklist.AuditableItemID = auditChecklistFormDTO.AuditableItemID
	existingAuditChecklist.Question = auditChecklistFormDTO.Question
	existingAuditChecklist.UpdatedBYID = &currentUserID

	err = acs.repo.Update(existingAuditChecklist, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		err := acs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.AuditChecklist{}.TableName(),
			id,
			auditChecklistFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.AuditChecklist{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (acs *AuditChecklistService) Delete(auditId uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := acs.repo.BeginTransaction()

	if _, err := acs.auditSvc.GetByID(auditId); err != nil {
		return nil, err
	}

	var failedIDs []uint
	var successfulIDs []uint
	documentsToDelete := make(map[uint][]*auditDTOs.AuditChecklistDocumentDTO)

	for _, id := range ids {

		if _, err := acs.GetByID(auditId, id); err != nil {
			failedIDs = append(failedIDs, id)
		}
		existingAuditChecklistFullDetails, err := acs.GetFullDetailsByID(auditId, id)
		if err != nil {
			failedIDs = append(failedIDs, id)
			return nil, utility.NewResourceNotFoundError(err, id)
		}

		if len(existingAuditChecklistFullDetails.Documents) > 0 {
			var deleteDocumentIds []*uint
			for _, doc := range existingAuditChecklistFullDetails.Documents {
				if doc != nil {
					documentsToDelete[id] = append(documentsToDelete[id], doc)
					deleteDocumentIds = append(deleteDocumentIds, doc.ID)
				}
			}
			if err := acs.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		if err := acs.repo.Delete(id, tx); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		successfulIDs = append(successfulIDs, id)
	}
	if len(successfulIDs) > 0 {
		for checkID, documents := range documentsToDelete {
			for _, doc := range documents {
				fmt.Println("filename=", doc.Name)
				storagePath := filepath.Join("audit", "audit-checklist", strconv.FormatUint(uint64(checkID), 10), "documents")
				if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
					tx.Rollback()
					fmt.Printf("Failed to delete file at URL %s: %v\n", *doc.URL, err)
				}
			}
		}
	}
	tx.Commit()
	if len(successfulIDs) > 0 {

		go func() {
			for _, id := range successfulIDs {
				err := acs.activityLogSvc.LogActivity(
					models.DeleteAction,
					models.AuditChecklist{}.TableName(),
					id,
					nil,
					currentUserID,
					nil,
				)

				if err != nil {
					fields := map[string]interface{}{
						"action":    models.DeleteAction,
						"tableName": models.AuditChecklist{}.TableName(),
						"recordID":  id,
						"userID":    currentUserID,
						"error":     err.Error(),
					}

					utility.LogError("Failed to log activity", fields)
				}
			}
		}()
	} else {
		tx.Rollback()
	}

	return failedIDs, nil
}

func (acs *AuditChecklistService) Answer(auditId uint, id uint, currentUserID uint, auditChecklistAnswerFormDTO *auditDTOs.AuditChecklistAnswerFormDTO) error {
	tx := acs.repo.BeginTransaction()

	if _, err := acs.auditSvc.GetByID(auditId); err != nil {
		return err
	}

	existingAuditChecklist, err := acs.GetByID(auditId, id)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}

	if err := acs.validateAuditChecklistAnswerFormDTO(auditChecklistAnswerFormDTO); err != nil {
		tx.Rollback()
		return err
	}
	now := time.Now()
	status, err := acs.statusSvc.GetByType("completed")
	if err != nil {
		return err
	}
	existingAuditChecklist.ConformanceID = &auditChecklistAnswerFormDTO.ConformanceID
	existingAuditChecklist.Comment = auditChecklistAnswerFormDTO.Comment
	existingAuditChecklist.AnsweredBYID = &currentUserID
	existingAuditChecklist.StatusID = status.ID
	existingAuditChecklist.AnsweredAt = &now

	err = acs.repo.Update(existingAuditChecklist, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	
	var documents []*models.AuditChecklistDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.AuditChecklistDocument

	for _, doc := range auditChecklistAnswerFormDTO.Documents {

		if doc.IsNewFile != nil && *doc.IsNewFile {
			// Encrypt the filename
			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)

			if err != nil {
				tx.Rollback()
				return err
			}

			//genereate URL
			url := utility.GenerateDocumentURL("audit", "audit-checklist", strconv.FormatUint(uint64(id), 10), "documents")
			// Create a new Control test Document instance
			document := &models.AuditChecklistDocument{
				AuditChecklistID: id,
				URL:                                  url,
				Name:                                 doc.FileName,
				Size:                                 doc.Size,
				Token:                                encryptedFileName,
				OriginalName:                         doc.OriginalName,
				Type:                                 doc.Type,
				CreatedBYID:                          currentUserID,
			}

			// Append the document to the slice
			documents = append(documents, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			if doc.ID != nil {
				existingDocument, err := acs.repo.GetDocumentByID(*doc.ID)
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return utility.NewResourceNotFoundError(err, id)
					}
					return err
				}

				documentsToDelete = append(documentsToDelete, existingDocument)
				deleteDocumentIds = append(deleteDocumentIds, &existingDocument.ID)
			}
		}
	}

	if len(documents) > 0 {
		if err := acs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIds) > 0 {
		if err := acs.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range documents {
		storagePath := filepath.Join("audit", "audit-checklist", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("audit", "audit-checklist", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	go func() {
		err := acs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.AuditChecklist{}.TableName(),
			id,
			auditChecklistAnswerFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.AuditChecklist{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (acs *AuditChecklistService) DeleteAnswer(auditId uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := acs.repo.BeginTransaction()

	if _, err := acs.auditSvc.GetByID(auditId); err != nil {
		return nil, err
	}

	var failedIDs []uint
	var successfulIDs []uint
	status, err := acs.statusSvc.GetByType("pending")
	if err != nil {
		return nil, err
	}
	documentsToDelete := make(map[uint][]*auditDTOs.AuditChecklistDocumentDTO)
	for _, id := range ids {
		existingAuditChecklist, err := acs.GetByID(auditId, id)
		if err != nil {
			failedIDs = append(failedIDs, id)
			return nil, utility.NewResourceNotFoundError(err, id)
		}

		existingAuditChecklist.ConformanceID = nil
		existingAuditChecklist.Comment = nil
		existingAuditChecklist.AnsweredBYID  = nil
		existingAuditChecklist.StatusID = status.ID
		existingAuditChecklist.AnsweredAt = nil
		existingAuditChecklistFullDetails, err := acs.GetFullDetailsByID(auditId, id)
		if err != nil {
			failedIDs = append(failedIDs, id)
			return nil, utility.NewResourceNotFoundError(err, id)
		}

		if len(existingAuditChecklistFullDetails.Documents) > 0 {
			var deleteDocumentIds []*uint
			for _, doc := range existingAuditChecklistFullDetails.Documents {
				if doc != nil {
					documentsToDelete[id] = append(documentsToDelete[id], doc)
					deleteDocumentIds = append(deleteDocumentIds, doc.ID)
				}
			}
			if err := acs.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		err = acs.repo.Update(existingAuditChecklist, tx)
		if err != nil {
			return nil, err
		}

		successfulIDs = append(successfulIDs, id)
	}
	if len(successfulIDs) > 0 {
		for checkID, documents := range documentsToDelete {
			for _, doc := range documents {
				fmt.Println("filename=", doc.Name)
				storagePath := filepath.Join("audit", "audit-checklist", strconv.FormatUint(uint64(checkID), 10), "documents")
				if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
					tx.Rollback()
					fmt.Printf("Failed to delete file at URL %s: %v\n", *doc.URL, err)
				}
			}
		}
	}
	tx.Commit()

	if len(successfulIDs) > 0 {

		go func() {
			for _, id := range successfulIDs {
				err := acs.activityLogSvc.LogActivity(
					models.DeleteAction,
					models.AuditChecklist{}.TableName(),
					id,
					nil,
					currentUserID,
					nil,
				)

				if err != nil {
					fields := map[string]interface{}{
						"action":    models.DeleteAction,
						"tableName": models.AuditChecklist{}.TableName(),
						"recordID":  id,
						"userID":    currentUserID,
						"error":     err.Error(),
					}

					utility.LogError("Failed to log activity", fields)
				}
			}
		}()
	} else {
		tx.Rollback()
	}

	return failedIDs, nil
}

func (acs *AuditChecklistService) validateAuditChecklistFormDTO(dto *auditDTOs.AuditChecklistFormDTO) error {
	if _, err := acs.auditPlanAuditableItemSvc.GetByAuditableItemID(*dto.AuditableItemID); err != nil {
		return utility.NewValidationError(fmt.Errorf("auditable item does not exist in auidt plan"), "not_exist", "auditable item")
	}

	return nil
}

func (acs *AuditChecklistService) validateAuditChecklistAnswerFormDTO(dto *auditDTOs.AuditChecklistAnswerFormDTO) error {
	if _, err := acs.conformanceSvc.GetByID(dto.ConformanceID); err != nil {
		return utility.NewValidationError(fmt.Errorf("conformance does not exist"), "not_exist", "conformance")
	}

	return nil
}
