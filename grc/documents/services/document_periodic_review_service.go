package services

import (
	"errors"
	"fmt"
	"jethur/config"
	documentDTOs "jethur/grc/documents/dtos"
	documentRepositories "jethur/grc/documents/repositories"
	"jethur/grc/logs/services"
	masterService "jethur/grc/masters/services"
	"jethur/grc/models"
	"jethur/utility"
	"os"
	"path/filepath"
	"strconv"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type DocumentPeriodicReviewServiceInterface interface {
	GetAll(params *documentDTOs.DocumentPeriodicReviewQueryParamsDTO, documentId uint) (*documentDTOs.DocumentPeriodicReviewListResponseDTO, error)
	GetByID(id uint) (*models.DocumentReview, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*documentDTOs.DocumentPeriodicReviewDetailsDTO, error)
	Create(currentUserID uint, documentID uint, formDTO *documentDTOs.DocumentPeriodicReviewFormDTO) (uint, error)
	Update(id uint, currentUserID uint, formDTO *documentDTOs.DocumentPeriodicReviewFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)

	GetDocumentPath(id uint, token string, tx ...utility.DatabaseTransaction) (string, error)
}

type DocumentPeriodicReviewService struct {
	repo                    documentRepositories.DocumentPeriodicReviewRepositoryInterface
	activityLogSvc          services.ActivityLogServiceInterface
	documentSvc             DocumentServiceInterface
	documentReviewStatusSvc masterService.ReviewStatusServiceInterface
}

func NewDocumentPeriodicReviewService(
	repo documentRepositories.DocumentPeriodicReviewRepositoryInterface,
	activityLogSvc services.ActivityLogServiceInterface,
	documentSvc DocumentServiceInterface,
	documentReviewStatusSvc masterService.ReviewStatusServiceInterface,
) DocumentPeriodicReviewServiceInterface {
	return &DocumentPeriodicReviewService{
		repo:                    repo,
		activityLogSvc:          activityLogSvc,
		documentSvc:             documentSvc,
		documentReviewStatusSvc: documentReviewStatusSvc,
	}
}

func (dprs *DocumentPeriodicReviewService) GetAll(params *documentDTOs.DocumentPeriodicReviewQueryParamsDTO, documentId uint) (*documentDTOs.DocumentPeriodicReviewListResponseDTO, error) {
	return dprs.repo.GetAll(params, documentId)
}

func (dprs *DocumentPeriodicReviewService) GetByID(id uint) (*models.DocumentReview, error) {
	review, err := dprs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return review, nil
}

func (dprs *DocumentPeriodicReviewService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*documentDTOs.DocumentPeriodicReviewDetailsDTO, error) {
	var detailsDTO documentDTOs.DocumentPeriodicReviewDetailsDTO
	review, err := dprs.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	err = mapstructure.Decode(review, &detailsDTO)
	if err != nil {
		return nil, err
	}
	detailsDTO.CreatedAt = review.CreatedAt
	detailsDTO.UpdatedAt = review.UpdatedAt
	return &detailsDTO, nil
}

func (dprs *DocumentPeriodicReviewService) Create(currentUserID uint, documentID uint, formDTO *documentDTOs.DocumentPeriodicReviewFormDTO) (uint, error) {
	tx := dprs.repo.BeginTransaction()

	dates, err := dprs.parseControlDates(formDTO)
	if err != nil {
		return 0, err
	}

	// Update the document's review status to "reviewed"
	reviewStatus, err := dprs.documentReviewStatusSvc.GetByType("reviewed")
	if err != nil {
		return 0, err
	}

	document, err := dprs.documentSvc.GetByID(documentID)
	if err != nil {
		return 0, err
	}

	review := &models.DocumentReview{
		DocumentID:       documentID,
		ReviewDate:       *dates.ReviewDate,
		ISChangeRequired: formDTO.ISChangeRequired,
		Description:      formDTO.Description,
		CreatedBYID:      currentUserID,
	}

	id, err := dprs.repo.Create(review, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	document.ReviewStatusID = &reviewStatus.ID

	if err := dprs.documentSvc.UpdateDocument(document, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.DocumentReviewSupportDocument
	for _, doc := range formDTO.SupportDocuments {
		dataToEncrypt := []byte(doc.FileName)
		encryptedFileName, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		url := utility.GenerateDocumentURL("document", "periodic-reviews", strconv.FormatUint(uint64(id), 10), "documents", encryptedFileName)
		document := &models.DocumentReviewSupportDocument{
			DocumentReviewID: id,
			URL:              url,
			Name:             doc.FileName,
			Size:             doc.Size,
			Token:            encryptedFileName,
			OriginalName:     &doc.OriginalName,
			Type:             doc.Type,
		}
		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := dprs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()

	for _, doc := range documents {
		storagePath := filepath.Join("document", "periodic-reviews", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			return 0, err
		}
	}

	return id, nil
}

func (dprs *DocumentPeriodicReviewService) Update(id uint, currentUserID uint, formDTO *documentDTOs.DocumentPeriodicReviewFormDTO) error {
	tx := dprs.repo.BeginTransaction()

	existing, err := dprs.GetByID(id)
	if err != nil {
		return err
	}

	dates, err := dprs.parseControlDates(formDTO)
	if err != nil {
		return err
	}

	existing.ReviewDate = *dates.ReviewDate
	existing.ISChangeRequired = formDTO.ISChangeRequired
	existing.Description = formDTO.Description
	existing.UpdatedBYID = &currentUserID

	if err := dprs.repo.Update(existing, tx); err != nil {
		tx.Rollback()
		return err
	}

	var documents []*models.DocumentReviewSupportDocument
	for _, doc := range formDTO.SupportDocuments {
		dataToEncrypt := []byte(doc.FileName)
		encryptedFileName, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return err
		}
		url := utility.GenerateDocumentURL("document", "periodic-reviews", strconv.FormatUint(uint64(id), 10), "documents", encryptedFileName)
		document := &models.DocumentReviewSupportDocument{
			DocumentReviewID: id,
			URL:              url,
			Name:             doc.FileName,
			Size:             doc.Size,
			Token:            encryptedFileName,
			OriginalName:     &doc.OriginalName,
			Type:             doc.Type,
		}
		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := dprs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	for _, doc := range documents {
		storagePath := filepath.Join("document", "periodic-reviews", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			fmt.Printf("Failed to store file %s: %v\n", doc.Name, err)
		}
	}

	return nil
}

func (dprs *DocumentPeriodicReviewService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := dprs.repo.BeginTransaction()

	var failedIds []uint
	var successfulIds []uint

	for _, id := range ids {
		if err := dprs.repo.Delete(id, tx); err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}

	tx.Commit()

	return failedIds, nil
}

func (dprs *DocumentPeriodicReviewService) GetDocumentPath(id uint, token string, tx ...utility.DatabaseTransaction) (string, error) {
	// Retrieve the full details of the document review
	existingDocumentReview, err := dprs.GetFullDetailsByID(id, tx...)
	if err != nil {
		return "", err
	}

	// Decrypt the provided token to get the file name
	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	// Construct the storage path
	storagePath := filepath.Join("document", "periodic-reviews", strconv.FormatUint(uint64(existingDocumentReview.ID), 10), "documents")

	// Retrieve the base storage path from environment variables
	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)

	// Check if the file exists at the constructed path
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	// Return the full path to the file
	return fullPath, nil
}

func (dprs *DocumentPeriodicReviewService) parseControlDates(dto *documentDTOs.DocumentPeriodicReviewFormDTO) (*documentDTOs.DocumentPeriodicReviewParsedDates, error) {
	parsed := &documentDTOs.DocumentPeriodicReviewParsedDates{}

	var err error

	parsed.ReviewDate, err = utility.ParseDate(dto.ReviewDate)
	if err != nil {
		return nil, utility.NewValidationError(fmt.Errorf("invalid last review date"), "invalid_date", "last review date")
	}

	return parsed, nil
}
