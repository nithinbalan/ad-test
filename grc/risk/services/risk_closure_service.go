package services

import (
	"errors"
	"path/filepath"
	"strconv"

	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"
	riskRepositories "jethur/grc/risk/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type RiskClosureServiceInterface interface {
	GetAll(riskID uint) ([]*riskDTOs.RiskClosureListDTO, error)
	GetFullDetailsByID(riskID uint, id uint, tx ...utility.DatabaseTransaction) (*riskDTOs.RiskClosureDetailsDTO, error)
	Create(riskID uint, currentUserID uint, closureFormDTO *riskDTOs.RiskClosureFormDTO) (uint, error)
	Update(riskID uint, id uint, currentUserID uint, closureFormDTO *riskDTOs.RiskClosureFormDTO) error
	Delete(riskID uint, id uint, currentUserID uint) (uint, error)
}

type RiskClosureService struct {
	repo            riskRepositories.RiskClosureRepositoryInterface
	activityLogSvc  logServices.ActivityLogServiceInterface
	riskRegisterSvc RiskRegisterServiceInterface
}

func NewRiskClosureService(
	repo riskRepositories.RiskClosureRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	riskRegisterSvc RiskRegisterServiceInterface,
) RiskClosureServiceInterface {
	return &RiskClosureService{
		repo:            repo,
		activityLogSvc:  activityLogSvc,
		riskRegisterSvc: riskRegisterSvc,
	}
}

func (rcs *RiskClosureService) GetAll(riskID uint) ([]*riskDTOs.RiskClosureListDTO, error) {
	return rcs.repo.GetAll(riskID)
}

func (rcs *RiskClosureService) GetFullDetailsByID(riskID uint, id uint, tx ...utility.DatabaseTransaction) (*riskDTOs.RiskClosureDetailsDTO, error) {
	if _, err := rcs.riskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}

	closureDetails, err := rcs.repo.GetFullDetailsByID(riskID, id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, riskID)
		}
		return nil, err
	}

	var closureDetailsDTO riskDTOs.RiskClosureDetailsDTO
	err = mapstructure.Decode(closureDetails, &closureDetailsDTO)
	if err != nil {
		return nil, err
	}

	closureDetailsDTO.CreatedAt = closureDetails.CreatedAt
	closureDetailsDTO.UpdatedAt = closureDetails.UpdatedAt

	return &closureDetailsDTO, nil
}

func (rcs RiskClosureService) GetByID(riskID uint, id uint) (*models.RiskClosure, error) {
	if _, err := rcs.riskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}

	riskClosure, err := rcs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return riskClosure, nil
}

func (rcs *RiskClosureService) Create(riskID uint, currentUserID uint, closureFormDTO *riskDTOs.RiskClosureFormDTO) (uint, error) {
	tx := rcs.repo.BeginTransaction()

	if _, err := rcs.riskRegisterSvc.GetByID(riskID); err != nil {
		return 0, err
	}

	closure := &models.RiskClosure{
		RiskRegisterID: riskID,
		Justification:  closureFormDTO.Justification,
		LessonLearned:  closureFormDTO.LessonLearned,
		WhatWorkedWell: closureFormDTO.WhatWorkedWell,
		Improvements:   closureFormDTO.Improvements,
		Takeaways:      closureFormDTO.Takeaways,
		CreatedBYID:    currentUserID,
	}

	id, err := rcs.repo.Create(closure, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.RiskClosureDocument
	for _, doc := range closureFormDTO.Documents {

		// Encrypt the filename
		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		// Construct the URL for the document
		url := utility.GenerateDocumentURL("risk", "risk-register", strconv.FormatUint(uint64(riskID), 10), "closure", strconv.FormatUint(uint64(id), 10), "documents")

		// Create a new RiskDocument instance
		document := &models.RiskClosureDocument{
			RiskClosureID: id,
			URL:           url,
			Name:          doc.FileName,
			Size:          doc.Size,
			Token:         encryptedFilename,
			OriginalName:  doc.OriginalName,
			Type:          doc.Type,
			CreatedBYID:   currentUserID,
		}

		// Append the document to the slice
		documents = append(documents, document)
	}
	if len(documents) > 0 {
		// Save the RiskDocument
		if err := rcs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		for _, doc := range documents {
			storagePath := filepath.Join("risk", "risk-register", strconv.FormatUint(uint64(riskID), 10), "closure", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}

	tx.Commit()

	// Log activity
	go func() {
		rcs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.RiskClosure{}.TableName(),
			id,
			closure,
			currentUserID,
			nil,
		)
	}()

	return id, nil
}

// Update an existing risk closure using DTO.
func (rcs *RiskClosureService) Update(riskID uint, id uint, currentUserID uint, closureFormDTO *riskDTOs.RiskClosureFormDTO) error {
	tx := rcs.repo.BeginTransaction()

	if _, err := rcs.riskRegisterSvc.GetByID(riskID); err != nil {
		return err
	}

	existingClosure, err := rcs.GetByID(riskID, id)
	if err != nil {
		return err
	}

	existingClosure.Justification = closureFormDTO.Justification
	existingClosure.LessonLearned = closureFormDTO.LessonLearned
	existingClosure.WhatWorkedWell = closureFormDTO.WhatWorkedWell
	existingClosure.Improvements = closureFormDTO.Improvements
	existingClosure.Takeaways = closureFormDTO.Takeaways
	existingClosure.UpdatedBYID = &currentUserID

	if err := rcs.repo.Update(existingClosure, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Save Documents
	// documents update
	var documents []*models.RiskClosureDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.RiskClosureDocument

	// Iterate over each document in the array
	for _, doc := range closureFormDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {

			dataToEncrypt := []byte(doc.FileName)
			encryptedFilename, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}

			url := utility.GenerateDocumentURL("risk", "risk-register", strconv.FormatUint(uint64(riskID), 10), "closure", strconv.FormatUint(uint64(id), 10), "documents")
			// Create a new RiskDocument instance
			document := &models.RiskClosureDocument{
				RiskClosureID: id,
				URL:           url,
				Name:          doc.FileName,
				Size:          doc.Size,
				Token:         encryptedFilename,
				OriginalName:  doc.OriginalName,
				Type:          doc.Type,
				CreatedBYID:   currentUserID,
			}
			documents = append(documents, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			// Handle deleted file case
			if doc.ID != nil {
				existingDocument, err := rcs.repo.GetDocumentByID(*doc.ID)
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
		if err := rcs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIds) > 0 {
		if err := rcs.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(documents) > 0 {
		for _, doc := range documents {
			storagePath := filepath.Join("risk", "risk-register", strconv.FormatUint(uint64(riskID), 10), "closure", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	if len(documentsToDelete) > 0 {
		for _, document := range documentsToDelete {
			storagePath := filepath.Join("risk", "risk-register", strconv.FormatUint(uint64(riskID), 10), "closure", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.RemoveFile(document.Name, storagePath); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	tx.Commit()

	// Log activity
	go func() {
		rcs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.RiskClosure{}.TableName(),
			id,
			existingClosure,
			currentUserID,
			nil,
		)
	}()

	return nil
}

func (rcs *RiskClosureService) Delete(riskID uint, id uint, currentUserID uint) (uint, error) {
	tx := rcs.repo.BeginTransaction()

	if _, err := rcs.riskRegisterSvc.GetByID(riskID); err != nil {
		return 0, err
	}

	riskClosure, err := rcs.GetFullDetailsByID(riskID, id, tx)
	if err != nil {
		return 0, err
	}

	var deleteDocumentIds []*uint
	documentsToDelete := make([]*riskDTOs.RiskClosureDocumentDTO, 0)
	if len(riskClosure.Documents) > 0 {
		for _, doc := range riskClosure.Documents {
			if doc != nil {
				documentsToDelete = append(documentsToDelete, doc)
				deleteDocumentIds = append(deleteDocumentIds, &doc.ID)
			}
		}

		if len(deleteDocumentIds) > 0 {
			if err := rcs.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}

	if err := rcs.repo.Delete(id, tx); err != nil {
		tx.Rollback()
		return id, err
	}

	// Remove files associated with documents
	for _, doc := range documentsToDelete {
		storagePath := filepath.Join("risk", "risk-register", strconv.FormatUint(uint64(riskID), 10), "closure", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
		}

		destinationPath := filepath.Join("risk", "risk-register", strconv.FormatUint(uint64(riskID), 10), "closure", strconv.FormatUint(uint64(id), 10))
		if err := utility.RemoveFolder(destinationPath); err != nil {
			tx.Rollback()
		}
	}

	tx.Commit()

	// Log activity
	go func() {
		rcs.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.RiskClosure{}.TableName(),
			id,
			nil,
			currentUserID,
			nil,
		)
	}()

	return 0, nil
}
