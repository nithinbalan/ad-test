package services

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"jethur/config"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/project/dtos"
	"jethur/grc/project/repositories"
	"jethur/utility"
)

type ProjectSecurityGateChecklistCommentServiceInterface interface {
	GetByID(securityGateID uint, checklistID uint, id uint, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateChecklistComment, error)
	Create(securityGateID uint, checklistID uint, currentUserID uint, form *dtos.ProjectSecurityGateChecklistCommentFormDTO, historyID *uint, tx ...utility.DatabaseTransaction) (uint, error)
	Update(securityGateID uint, checklistID uint, id uint, currentUserID uint, form *dtos.ProjectSecurityGateChecklistCommentFormDTO) error
	Delete(securityGateID uint, checklistID uint, ids []uint, currentUserID uint) ([]uint, error)
	GetDocumentPath(securityGateID uint, checklistID uint, id uint, token string) (string, error)
}

type ProjectSecurityGateChecklistCommentService struct {
	repo         repositories.ProjectSecurityGateChecklistCommentRepositoryInterface
	checklistSvc ProjectSecurityGateChecklistServiceInterface
	activityLog  logServices.ActivityLogServiceInterface
}

func NewProjectSecurityGateChecklistCommentService(
	repo repositories.ProjectSecurityGateChecklistCommentRepositoryInterface,
	checklistSvc ProjectSecurityGateChecklistServiceInterface,
	activityLog logServices.ActivityLogServiceInterface,
) ProjectSecurityGateChecklistCommentServiceInterface {
	return &ProjectSecurityGateChecklistCommentService{
		repo:         repo,
		checklistSvc: checklistSvc,
		activityLog:  activityLog,
	}
}

func (psgccs *ProjectSecurityGateChecklistCommentService) GetByID(securityGateID uint, checklistID uint, id uint, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateChecklistComment, error) {
	if _, err := psgccs.checklistSvc.GetByID(securityGateID, checklistID, tx...); err != nil {
		return nil, err
	}

	comment, err := psgccs.repo.GetByID(id, tx...)
	if err != nil {
		return nil, err
	}

	if comment.ChecklistID != checklistID {
		return nil, &utility.ResourceNotFoundError{}
	}

	return comment, nil
}

func (psgccs *ProjectSecurityGateChecklistCommentService) GetDocumentPath(securityGateID uint, checklistID uint, id uint, token string) (string, error) {
	comment, err := psgccs.GetByID(securityGateID, checklistID, id)
	if err != nil {
		return "", err
	}

	if len(comment.Documents) == 0 {
		return "", &utility.ResourceNotFoundError{}
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", &utility.ResourceNotFoundError{}
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", &utility.ResourceNotFoundError{}
	}

	storagePath := filepath.Join(
		"project",
		"security-gates", strconv.FormatUint(uint64(securityGateID), 10),
		"checklists", strconv.FormatUint(uint64(checklistID), 10),
		"comments", strconv.FormatUint(uint64(id), 10),
		"documents",
	)

	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", &utility.ResourceNotFoundError{}
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}

func (psgccs *ProjectSecurityGateChecklistCommentService) Create(securityGateID uint, checklistID uint, currentUserID uint, form *dtos.ProjectSecurityGateChecklistCommentFormDTO, historyID *uint, tx ...utility.DatabaseTransaction) (id uint, err error) {
	ownsTx := len(tx) == 0 || tx[0] == nil

	var dbTx utility.DatabaseTransaction
	if ownsTx {
		dbTx = psgccs.repo.BeginTransaction()
	} else {
		dbTx = tx[0]
	}

	defer func() {
		if p := recover(); p != nil {
			if ownsTx {
				dbTx.Rollback()
			}
			panic(p)
		}

		if err != nil && ownsTx {
			dbTx.Rollback()
		}
	}()

	if _, err = psgccs.checklistSvc.GetByID(securityGateID, checklistID, dbTx); err != nil {
		return 0, err
	}

	comment := &models.ProjectSecurityGateChecklistComment{
		ChecklistID: checklistID,
		HistoryID:   historyID,
		Comment:     form.Comment,
		CreatedBYID: currentUserID,
		CreatedAt:   time.Now(),
	}

	if err = psgccs.repo.Create(comment, dbTx); err != nil {
		return 0, err
	}

	id = comment.ID

	var documents []*models.ProjectSecurityGateChecklistCommentDocument

	for _, doc := range form.Documents {
		encryptedFileName, encErr := utility.EncryptData([]byte(doc.FileName))
		if encErr != nil {
			err = encErr
			return 0, err
		}

		url := utility.GenerateDocumentURL(
			"project",
			"security-gates", strconv.FormatUint(uint64(securityGateID), 10),
			"checklists", strconv.FormatUint(uint64(checklistID), 10),
			"comments", strconv.FormatUint(uint64(id), 10),
			"documents",
		)

		documents = append(documents, &models.ProjectSecurityGateChecklistCommentDocument{
			ProjectSecurityGateChecklistCommentID: id,
			URL:                                   url,
			Name:                                  doc.FileName,
			Size:                                  doc.Size,
			Token:                                 encryptedFileName,
			OriginalName:                          doc.OriginalName,
			Type:                                  doc.Type,
			CreatedBYID:                           currentUserID,
			CreatedAt:                             time.Now(),
		})
	}

	if len(documents) > 0 {
		if err = psgccs.repo.CreateDocuments(documents, dbTx); err != nil {
			return 0, err
		}

		storagePath := filepath.Join(
			"project",
			"security-gates", strconv.FormatUint(uint64(securityGateID), 10),
			"checklists", strconv.FormatUint(uint64(checklistID), 10),
			"comments", strconv.FormatUint(uint64(id), 10),
			"documents",
		)

		for _, document := range documents {
			if err = utility.StoreFile(document.Name, storagePath); err != nil {
				return 0, err
			}
		}
	}

	if ownsTx {
		if err = dbTx.Commit(); err != nil {
			return 0, err
		}
	}

	go func() {
		activityErr := psgccs.activityLog.LogActivity(
			models.CreateAction,
			models.ProjectSecurityGateChecklistComment{}.TableName(),
			id,
			comment,
			currentUserID,
			"project", "project_security_gate", "project_security_gate_checklist_comment_create", nil,
		)
		if activityErr != nil {
			log.Printf("Failed to log activity: %v", activityErr)
		}
	}()

	return id, err
}

func (psgccs *ProjectSecurityGateChecklistCommentService) Update(securityGateID uint, checklistID uint, id uint, currentUserID uint, form *dtos.ProjectSecurityGateChecklistCommentFormDTO) (err error) {
	tx := psgccs.repo.BeginTransaction()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil {
			tx.Rollback()
			return
		}
	}()

	comment, err := psgccs.GetByID(securityGateID, checklistID, id, tx)
	if err != nil {
		return err
	}

	oldData := map[string]interface{}{
		"comment": comment.Comment,
	}

	comment.Comment = form.Comment

	newData := map[string]interface{}{
		"comment": comment.Comment,
	}
	changes := utility.CompareFieldChanges(oldData, newData)

	if err = psgccs.repo.Update(comment, tx); err != nil {
		return err
	}

	var newDocuments []*models.ProjectSecurityGateChecklistCommentDocument
	var deleteDocumentIDs []*uint
	var documentsToDelete []*models.ProjectSecurityGateChecklistCommentDocument

	for _, doc := range form.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {
			encryptedFileName, encErr := utility.EncryptData([]byte(doc.FileName))
			if encErr != nil {
				err = encErr
				return err
			}

			url := utility.GenerateDocumentURL(
				"project",
				"security-gates", strconv.FormatUint(uint64(securityGateID), 10),
				"checklists", strconv.FormatUint(uint64(checklistID), 10),
				"comments", strconv.FormatUint(uint64(comment.ID), 10),
				"documents",
			)

			newDocuments = append(newDocuments, &models.ProjectSecurityGateChecklistCommentDocument{
				ProjectSecurityGateChecklistCommentID: comment.ID,
				URL:                                   url,
				Name:                                  doc.FileName,
				Size:                                  doc.Size,
				Token:                                 encryptedFileName,
				OriginalName:                          doc.OriginalName,
				Type:                                  doc.Type,
				CreatedBYID:                           currentUserID,
				CreatedAt:                             time.Now(),
			})
		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			if doc.ID != nil {
				existingDocument, docErr := psgccs.repo.GetDocumentByID(*doc.ID, tx)
				if docErr != nil {
					err = docErr
					return err
				}
				if existingDocument.ProjectSecurityGateChecklistCommentID != comment.ID {
					err = &utility.ResourceNotFoundError{}
					return err
				}

				documentsToDelete = append(documentsToDelete, existingDocument)
				deleteDocumentIDs = append(deleteDocumentIDs, &existingDocument.ID)
			}
		}
	}

	if len(newDocuments) > 0 {
		if err = psgccs.repo.CreateDocuments(newDocuments, tx); err != nil {
			return err
		}
	}

	if len(deleteDocumentIDs) > 0 {
		if err = psgccs.repo.DeleteDocuments(deleteDocumentIDs, tx); err != nil {
			return err
		}
	}

	storagePath := filepath.Join(
		"project",
		"security-gates", strconv.FormatUint(uint64(securityGateID), 10),
		"checklists", strconv.FormatUint(uint64(checklistID), 10),
		"comments", strconv.FormatUint(uint64(comment.ID), 10),
		"documents",
	)

	for _, document := range newDocuments {
		if err = utility.StoreFile(document.Name, storagePath); err != nil {
			return err
		}
	}

	for _, document := range documentsToDelete {
		if err = utility.RemoveFile(document.Name, storagePath); err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	go func() {
		if len(changes) == 0 {
			return
		}

		activityErr := psgccs.activityLog.LogActivity(
			models.UpdateAction,
			models.ProjectSecurityGateChecklistComment{}.TableName(),
			comment.ID,
			changes,
			currentUserID,
			"project", "project_security_gate", "project_security_gate_checklist_comment_update", nil,
		)
		if activityErr != nil {
			log.Printf("Failed to log update activity: %v", activityErr)
		}
	}()

	return err
}

func (psgccs *ProjectSecurityGateChecklistCommentService) Delete(securityGateID uint, checklistID uint, ids []uint, currentUserID uint) (failedIDs []uint, err error) {
	tx := psgccs.repo.BeginTransaction()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil || len(failedIDs) > 0 {
			tx.Rollback()
			return
		}
	}()

	if _, err = psgccs.checklistSvc.GetByID(securityGateID, checklistID, tx); err != nil {
		return failedIDs, err
	}

	var successfulIDs []uint
	var deletedComments []map[string]interface{}

	for _, id := range ids {
		comment, fetchErr := psgccs.repo.GetByID(id, tx)
		if fetchErr != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		if comment.ChecklistID != checklistID {
			failedIDs = append(failedIDs, id)
			continue
		}

		if deleteErr := psgccs.repo.DeleteDocumentsByCommentID(id, tx); deleteErr != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		if deleteErr := psgccs.repo.Delete(id, tx); deleteErr != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		successfulIDs = append(successfulIDs, id)
		deletedComments = append(deletedComments, map[string]interface{}{
			"id":           id,
			"checklist_id": comment.ChecklistID,
			"comment":      comment.Comment,
		})
	}

	if len(failedIDs) > 0 {
		return failedIDs, err
	}

	if err = tx.Commit(); err != nil {
		return failedIDs, err
	}

	for _, id := range successfulIDs {
		folderPath := filepath.Join(
			"project",
			"security-gates", strconv.FormatUint(uint64(securityGateID), 10),
			"checklists", strconv.FormatUint(uint64(checklistID), 10),
			"comments", strconv.FormatUint(uint64(id), 10),
		)

		if removeErr := utility.RemoveFolder(folderPath); removeErr != nil {
			utility.LogError("Failed to remove checklist comment storage folder", map[string]interface{}{
				"commentID": id,
				"path":      folderPath,
				"error":     removeErr.Error(),
			})
		}
	}

	if len(successfulIDs) > 0 {
		go func() {
			activityErr := psgccs.activityLog.LogBulkActivities(
				models.DeleteAction,
				models.ProjectSecurityGateChecklistComment{}.TableName(),
				successfulIDs,
				deletedComments,
				currentUserID,
				nil,
			)
			if activityErr != nil {
				log.Printf("Failed to log delete activity: %v", activityErr)
			}
		}()
	}

	return failedIDs, err
}
