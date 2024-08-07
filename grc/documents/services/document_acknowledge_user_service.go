package services

import (
	"fmt"
	"jethur/grc/documents/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"
)

type DocumentAcknowledgeUserServiceInterface interface {
	Create(docuemntID uint, currentUserID uint) (uint, error)
	Delete(docuemntID uint, currentUserID uint) (uint, error)
}

type DocumentAcknowledgeUserService struct {
	repo           repositories.DocumentAcknowledgeUserRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
	documentSvc    DocumentServiceInterface
}

func NewDocumentAcknowledgeUserService(repo repositories.DocumentAcknowledgeUserRepositoryInterface, activityLogSvc logServices.ActivityLogServiceInterface, documentSvc DocumentServiceInterface) DocumentAcknowledgeUserServiceInterface {
	return &DocumentAcknowledgeUserService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
		documentSvc:    documentSvc,
	}
}

func (daus *DocumentAcknowledgeUserService) Create(docuemntID uint, currentUserID uint) (uint, error) {
	tx := daus.repo.BeginTransaction()

	_, err := daus.documentSvc.GetByID(docuemntID)
	if err != nil {
		return 0, err
	}
	// Check if the entry already exists
	existing, err := daus.repo.Exists(docuemntID, currentUserID, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if existing {
		tx.Rollback()
		return 0, utility.NewForbiddenError("document acknowledgment by user already exists")
	}

	documentAcknowledgeUser := &models.DocumentAcknowledgeUser{
		DocumentID:  docuemntID,
		UserID:      currentUserID,
		CreatedBYID: currentUserID,
	}

	id, err := daus.repo.Create(documentAcknowledgeUser, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	go func() {
		err := daus.activityLogSvc.LogActivity(
			models.CreateAction,
			models.DocumentAcknowledgeUser{}.TableName(),
			documentAcknowledgeUser.ID,
			nil,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.DocumentAcknowledgeUser{}.TableName(),
				"recordID":  id.ID,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id.ID, nil
}

func (daus *DocumentAcknowledgeUserService) Delete(docuemntID uint, currentUserID uint) (uint, error) {
	tx := daus.repo.BeginTransaction()
	var failedId uint

	_, err := daus.documentSvc.GetByID(docuemntID)
	if err != nil {
		return 0, err
	}
	// Check if the entry already exists
	exists, err := daus.repo.Exists(docuemntID, currentUserID, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if !exists {
		tx.Rollback()
		return 0, utility.NewForbiddenError("document acknowledgment by user does not exist")
	}

	if err = daus.repo.Delete(uint(docuemntID), currentUserID, tx); err != nil {
		failedId = uint(docuemntID)
		return failedId, err
	}
	tx.Commit()
	// Log successful delete activity
	go func() {
		err := daus.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.DocumentAcknowledgeUser{}.TableName(),
			uint(docuemntID),
			nil,
			currentUserID,
			nil,
		)
		if err != nil {
			// Handle error if logging fails
			fmt.Printf("Failed to log delete activity: %v\n", err)
		}
	}()

	return failedId, nil
}
