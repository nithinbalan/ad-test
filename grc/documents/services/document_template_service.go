package services

import (
	"errors"
	"fmt"
	documentDTOs "jethur/grc/documents/dtos"
	"jethur/grc/documents/repositories"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"
	"log"
	"path/filepath"
	"strconv"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type DocumentTemplateServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*documentDTOs.DocumentTemplateListDTO, error)
	GetByID(id uint) (*models.DocumentTemplate, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*documentDTOs.DocumentTemplateDetailsDTO, error)
	Create(currentUserID uint, documentTemplateFormDTO *documentDTOs.DocumentTemplateFormDTO) (uint, error)
	Update(id uint, currentUserID uint, documentTemplateFormDTO *documentDTOs.DocumentTemplateFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

type DocumentTemplateService struct {
	repo            repositories.DocumentTemplateRepositoryInterface
	activityLogSvc  services.ActivityLogServiceInterface
	documentTypeSvc DocumentTypeServiceInterface
}

func NewDocumentTemplateService(
	repo repositories.DocumentTemplateRepositoryInterface,
	activityLogSvc services.ActivityLogServiceInterface,
	documentTypeSvc DocumentTypeServiceInterface) DocumentTemplateServiceInterface {
	return &DocumentTemplateService{
		repo:            repo,
		activityLogSvc:  activityLogSvc,
		documentTypeSvc: documentTypeSvc,
	}
}

func (dts *DocumentTemplateService) GetAll(sort string, sortBy string, keyword string, status string) ([]*documentDTOs.DocumentTemplateListDTO, error) {
	return dts.repo.GetAll(sort, sortBy, keyword, status)
}

func (dts *DocumentTemplateService) GetByID(id uint) (*models.DocumentTemplate, error) {
	documentTemplate, err := dts.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return documentTemplate, nil
}

func (dts *DocumentTemplateService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*documentDTOs.DocumentTemplateDetailsDTO, error) {
	var documentTempllateDTO documentDTOs.DocumentTemplateDetailsDTO
	documentTemplate, err := dts.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	err = mapstructure.Decode(documentTemplate, &documentTempllateDTO)
	if err != nil {
		return nil, err
	}

	documentTempllateDTO.CreatedAt = documentTemplate.CreatedAt
	documentTempllateDTO.UpdatedAt = documentTemplate.UpdatedAt

	return &documentTempllateDTO, nil
}

func (dts *DocumentTemplateService) Create(currentUserID uint, documentTemplateFormDTO *documentDTOs.DocumentTemplateFormDTO) (uint, error) {
	fmt.Println("123")
	tx := dts.repo.BeginTransaction()

	if err := dts.validateDocumentTemplateFormDTO(0, documentTemplateFormDTO); err != nil {
		return 0, err
	}

	documentTemplate := &models.DocumentTemplate{
		TitleEn:        documentTemplateFormDTO.TitleEn,
		TitleAr:        documentTemplateFormDTO.TitleAr,
		DocumentTypeID: documentTemplateFormDTO.DocumentTypeID,
		Status:         documentTemplateFormDTO.Status,
		CreatedByID:    currentUserID,
	}

	id, err := dts.repo.Create(documentTemplate, tx)
	if err != nil {
		return 0, err
	}

	//generate token
	dataToEncrypt := []byte(documentTemplateFormDTO.Document.FileName)
	encryptedFilename, err := utility.EncryptData(dataToEncrypt)
	if err != nil {
		return 0, err
	}

	//genereate URL
	url := utility.GenerateDocumentURL("document", "document-template", strconv.FormatUint(uint64(id), 10), "documents")

	documentTemplateDocument := &models.DocumentTemplateDocument{
		DocumentTemplateID: id,
		URL:                url,
		Name:               documentTemplateFormDTO.Document.FileName,
		Size:               documentTemplateFormDTO.Document.Size,
		Token:              encryptedFilename,
		OriginalName:       documentTemplateFormDTO.Document.OriginalName,
		Type:               documentTemplateFormDTO.Document.Type,
		CreatedByID:        currentUserID,
	}

	// Save document template document
	if err := dts.repo.CreateDocument(documentTemplateDocument, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	// Store document template document File using utility.StoreFile
	storagePath := filepath.Join("document", "document-template", strconv.FormatUint(uint64(id), 10), "documents")

	if err := utility.StoreFile(documentTemplateFormDTO.Document.FileName, storagePath); err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	go func() {
		err := dts.activityLogSvc.LogActivity(
			models.CreateAction,
			models.DocumentTemplate{}.TableName(),
			id,
			documentTemplateFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return id, nil
}

func (dts *DocumentTemplateService) Update(id uint, currentUserID uint, documentTemplateFormDTO *documentDTOs.DocumentTemplateFormDTO) error {
	tx := dts.repo.BeginTransaction()

	existingTemplate, err := dts.GetByID(id)
	if err != nil {
		return err
	}

	if err := dts.validateDocumentTemplateFormDTO(id, documentTemplateFormDTO); err != nil {
		return err
	}

	existingTemplate.TitleEn = documentTemplateFormDTO.TitleEn
	existingTemplate.TitleAr = documentTemplateFormDTO.TitleAr
	existingTemplate.DocumentTypeID = documentTemplateFormDTO.DocumentTypeID
	existingTemplate.Status = documentTemplateFormDTO.Status
	existingTemplate.UpdatedByID = &currentUserID

	err = dts.repo.Update(existingTemplate, tx)

	if err != nil {
		tx.Rollback()
		return err
	}

	existingTemplateDocument, err := dts.repo.GetDocumentByID(id)
	if err != nil {
		return err
	}

	if documentTemplateFormDTO.Document.IsNewFile != nil && *documentTemplateFormDTO.Document.IsNewFile {
		existingTemplateDocumentFileName := existingTemplateDocument.Name

		//generate token
		dataToEncrypt := []byte(documentTemplateFormDTO.Document.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return err
		}

		//genereate URL
		url := utility.GenerateDocumentURL("document", "document-template", strconv.FormatUint(uint64(id), 10), "documents")

		existingTemplateDocument.URL = url
		existingTemplateDocument.Name = documentTemplateFormDTO.Document.FileName
		existingTemplateDocument.Size = documentTemplateFormDTO.Document.Size
		existingTemplateDocument.Token = encryptedFilename
		existingTemplateDocument.OriginalName = documentTemplateFormDTO.Document.OriginalName
		existingTemplateDocument.Type = documentTemplateFormDTO.Document.Type
		existingTemplateDocument.UpdatedByID = &currentUserID

		if err := dts.repo.CreateDocument(existingTemplateDocument, tx); err != nil {
			tx.Rollback()
			return err
		}

		storagePath := filepath.Join("document", "document-template", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.RemoveFile(existingTemplateDocumentFileName, storagePath); err != nil {
			tx.Rollback()
			return err
		}

		if err := utility.StoreFile(documentTemplateFormDTO.Document.FileName, storagePath); err != nil {
			tx.Rollback()
			return err
		}

	}

	tx.Commit()
	return nil
}

func (dts *DocumentTemplateService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := dts.repo.BeginTransaction()
	var failedIds []uint
	var successfulIds []uint
	documentsToDelete := make(map[uint][]*documentDTOs.DocumentTemplateDocumentDTO)
	for _, id := range ids {
		documentTemplate, err := dts.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			continue
		}

		if documentTemplate.Document != nil {
			documentsToDelete[id] = append(documentsToDelete[id], documentTemplate.Document)
			if err := dts.repo.DeleteDocument(documentTemplate.Document.ID, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		err = dts.repo.Delete(id, tx)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}

	if len(successfulIds) > 0 {
		for documentTemplateID, documents := range documentsToDelete {
			for _, doc := range documents {
				storagePath := filepath.Join("document", "document-template", strconv.FormatUint(uint64(documentTemplateID), 10), "documents")
				if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
					//tx.Rollback()
					fmt.Printf("Failed to delete file at URL %s: %v\n", doc.URL, err)
				}
			}
			destinationPath := filepath.Join("document", "document-template", strconv.FormatUint(uint64(documentTemplateID), 10))
			if err := utility.RemoveFolder(destinationPath); err != nil {
				fmt.Printf("Failed to delete file at URL %s: %v\n", destinationPath, err)
			}
		}
	}

	tx.Commit()

	if len(successfulIds) > 0 {
		go func() {
			err := dts.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.DocumentTemplate{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (dts *DocumentTemplateService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := dts.repo.Activate(ids)
	if len(successfulIds) > 0 {
		go func() {
			err := dts.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.DocumentTemplate{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk activate activity: %v\n", err)
			}
		}()
	}
	return err
}

func (dts *DocumentTemplateService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := dts.repo.Deactivate(ids)
	if len(successfulIds) > 0 {
		go func() {
			err := dts.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.DocumentTemplate{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk deactivate activity: %v", err)
			}
		}()
	}
	return err
}

func (dts *DocumentTemplateService) validateDocumentTemplateFormDTO(id uint, documentTemplateFormDTO *documentDTOs.DocumentTemplateFormDTO) error {
	// Check for title uniqueness
	if id == 0 {
		// Create scenario
		if exists, err := dts.repo.ExistsByTitle(documentTemplateFormDTO.TitleEn, documentTemplateFormDTO.TitleAr); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	} else {
		// Update scenario
		if exists, err := dts.repo.ExistsByTitleExcludingID(id, documentTemplateFormDTO.TitleEn, documentTemplateFormDTO.TitleAr); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	if _, err := dts.documentTypeSvc.GetByID(documentTemplateFormDTO.DocumentTypeID); err != nil {
		return utility.NewValidationError(fmt.Errorf("document type does not exist"), "not_exist", "document type")
	}

	// if documentTemplateFormDTO.Document == nil {
	// 	return utility.NewValidationError(fmt.Errorf("document is required"), "required", "document")
	// }

	if documentTemplateFormDTO.Document.Type != ".docx" {
		return utility.NewValidationError(fmt.Errorf("only word files are allowed"), "invalid_word_file_type", "type")
	}

	return nil
}
