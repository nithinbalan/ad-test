package services

import (
	"errors"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/repositories"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	"jethur/utility"
	"path/filepath"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type ComplianceReviewServiceInterface interface {
	GetAll(complianceID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.ComplianceReviewListDTO, error)
	Create(complianceID uint, currentUserID uint, complianceReviewFormDTO *complianceDTOs.ComplianceReviewFormDTO) (uint, error)
	Update(complianceID uint, id uint, currentUserID uint, ComplianceReviewFormDTO *complianceDTOs.ComplianceReviewFormDTO) error
}

// Service and implementations
type ComplianceReviewService struct {
	repo                      repositories.ComplianceReviewRepositoryInterface
	activityLogSvc            logServices.ActivityLogServiceInterface
	complianceRegisterSvc     ComplianceRegisterServiceInterface
	reviewFrequencySvc        masterServices.ReviewFrequencyServiceInterface
	complianceReviewStatusSvc ComplianceReviewStatusServiceInterface
}

func NewComplianceReviewService(
	repo repositories.ComplianceReviewRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	complianceRegisterSvc ComplianceRegisterServiceInterface,
	reviewFrequencySvc masterServices.ReviewFrequencyServiceInterface,
	complianceReviewStatusSvc ComplianceReviewStatusServiceInterface,
) ComplianceReviewServiceInterface {
	return &ComplianceReviewService{
		repo:                      repo,
		activityLogSvc:            activityLogSvc,
		complianceRegisterSvc:     complianceRegisterSvc,
		reviewFrequencySvc:        reviewFrequencySvc,
		complianceReviewStatusSvc: complianceReviewStatusSvc,
	}
}

func (crs *ComplianceReviewService) GetAll(complianceID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.ComplianceReviewListDTO, error) {
	if _, err := crs.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}

	return crs.repo.GetAll(complianceID, sort, sortBy, keyword)
}

func (crs ComplianceReviewService) GetByID(complianceID uint, id uint) (*models.ComplianceReview, error) {
	if _, err := crs.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}

	complianceReview, err := crs.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return complianceReview, nil
}

func (crs *ComplianceReviewService) Create(complianceID uint, currentUserID uint, complianceReviewFormDTO *complianceDTOs.ComplianceReviewFormDTO) (uint, error) {
	tx := crs.repo.BeginTransaction()

	complianceRegister, err := crs.complianceRegisterSvc.GetByID(complianceID)
	if err != nil {
		return 0, err
	}

	complianceReview := &models.ComplianceReview{
		ComplianceRegisterID: complianceID,
		Comment:              complianceReviewFormDTO.Comment,
		ReviewedBy:           currentUserID,
		ReviewedAt:           time.Now(),
	}

	id, err := crs.repo.Create(complianceReview,tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if complianceRegister.ReviewFrequencyID != nil {
		reviewFrequency, err := crs.reviewFrequencySvc.GetByID(*complianceRegister.ReviewFrequencyID)
		if err != nil {
			return 0, err
		}
		nextReviewDate := utility.CalculateNextReviewDate(reviewFrequency.Type)
		complianceRegister.NextReviewDate = &nextReviewDate
	}

	complianceReviewStatus, err := crs.complianceReviewStatusSvc.GetByType("reviewed")
	if err != nil {
		return 0, err
	}

	currentTime := time.Now()

	complianceRegister.LastReviewedDate = &currentTime
	complianceRegister.ComplianceReviewStatusID = &complianceReviewStatus.ID

	if err := crs.complianceRegisterSvc.ReviewStatusUpdate(complianceRegister, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.ComplianceReviewDocument

	for _, doc := range complianceReviewFormDTO.Documents {

		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL("compliances", strconv.FormatUint(uint64(complianceID), 10), "reviews", encryptedFilename)

		document := &models.ComplianceReviewDocument{
			ComplianceReviewID: id,
			URL:                url,
			Name:               doc.FileName,
			Size:               doc.Size,
			Token:              encryptedFilename,
			OriginalName:       doc.OriginalName,
			Type:               doc.Type,
			CreatedBy:          currentUserID,
		}

		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := crs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		//move file to storage
		for _, doc := range documents {
			storagePath := filepath.Join("compliance", "compliance-register", strconv.FormatUint(uint64(complianceID), 10), "review-documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}

		}
	}
	tx.Commit()

	go func() {
		err := crs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ComplianceReview{}.TableName(),
			id,
			complianceReview,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ComplianceReview{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (crs *ComplianceReviewService) Update(complianceID uint, id uint, currentUserID uint, complianceReviewFormDTO *complianceDTOs.ComplianceReviewFormDTO) error {
	tx := crs.repo.BeginTransaction()

	if _, err := crs.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return err
	}

	existingComplianceReview, err := crs.GetByID(complianceID, id)

	if err != nil {
		return err
	}

	existingComplianceReview.Comment = complianceReviewFormDTO.Comment
	existingComplianceReview.UpdatedBy = &currentUserID

	err = crs.repo.Update(existingComplianceReview,tx)

	if err != nil {
		tx.Rollback()
		return err
	}

	var documents []*models.ComplianceReviewDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.ComplianceReviewDocument

	for _, doc := range complianceReviewFormDTO.Documents {

		if doc.IsNewFile != nil && *doc.IsNewFile {
			// Encrypt the filename
			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)

			if err != nil {
				tx.Rollback()
				return err
			}

			//genereate URL
			url := utility.GenerateDocumentURL("compliances", strconv.FormatUint(uint64(complianceID), 10), "reviews", encryptedFileName)
			// Create a new Control test Document instance
			document := &models.ComplianceReviewDocument{
				ComplianceReviewID: id,
				URL:                url,
				Name:               doc.FileName,
				Size:               doc.Size,
				Token:              encryptedFileName,
				OriginalName:       doc.OriginalName,
				Type:               doc.Type,
				CreatedBy:          currentUserID,
			}

			// Append the document to the slice
			documents = append(documents, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			if doc.ID != nil {
				existingDocument, err := crs.repo.GetDocumentByID(*doc.ID)
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
		if err := crs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIds) > 0 {
		if err := crs.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range documents {
		storagePath := filepath.Join("compliance", "compliance-register", strconv.FormatUint(uint64(complianceID), 10), "review-documents")
		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("compliance", "compliance-register", strconv.FormatUint(uint64(complianceID), 10), "review-documents")
		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	return nil
}
