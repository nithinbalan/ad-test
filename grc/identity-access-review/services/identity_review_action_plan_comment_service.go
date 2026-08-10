package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"jethur/config"
	"jethur/grc/identity-access-review/dtos"
	"jethur/grc/identity-access-review/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type IdentityReviewActionPlanCommentServiceInterface interface {
	GetAll(identityReviewUserID uint, actionPlanID uint) ([]*dtos.IdentityReviewActionPlanCommentListDTO, error)
	GetByID(identityReviewUserID uint, actionPlanID uint, id uint) (*models.IdentityReviewActionPlanComment, error)
	GetFullDetailsByID(identityReviewUserID uint, actionPlanID uint, id uint, tx ...utility.DatabaseTransaction) (*dtos.IdentityReviewActionPlanCommentDetailDTO, error)
	Create(identityReviewUserID uint, actionPlanID uint, currentUserID uint, dto *dtos.IdentityReviewActionPlanCommentFormDTO) (uint, error)
	Update(identityReviewUserID uint, actionPlanID uint, id uint, currentUserID uint, dto *dtos.IdentityReviewActionPlanCommentFormDTO) error
	Delete(identityReviewUserID uint, actionPlanID uint, ids []uint, currentUserID uint) ([]uint, error)

	GetDocumentPath(identityReviewUserID, actionPlanID, id uint, token string) (string, error)
}

type IdentityReviewActionPlanCommentService struct {
	repo                        repositories.IdentityReviewActionPlanCommentRepositoryInterface
	activityLogSvc              logServices.ActivityLogServiceInterface
	identityReviewActionPlanSvc IdentityReviewActionPlanServiceInterface
}

func NewIdentityReviewActionPlanCommentService(repo repositories.IdentityReviewActionPlanCommentRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	identityReviewActionPlanSvc IdentityReviewActionPlanServiceInterface,
) IdentityReviewActionPlanCommentServiceInterface {
	return &IdentityReviewActionPlanCommentService{
		repo:                        repo,
		activityLogSvc:              activityLogSvc,
		identityReviewActionPlanSvc: identityReviewActionPlanSvc,
	}
}

func (irapcs *IdentityReviewActionPlanCommentService) GetAll(identityReviewUserID uint, actionPlanID uint) ([]*dtos.IdentityReviewActionPlanCommentListDTO, error) {
	if _, err := irapcs.identityReviewActionPlanSvc.GetByID(identityReviewUserID, actionPlanID); err != nil {
		return nil, err
	}

	return irapcs.repo.GetAll(actionPlanID)
}

func (irapcs *IdentityReviewActionPlanCommentService) GetByID(identityReviewUserID uint, actionPlanID uint, id uint) (*models.IdentityReviewActionPlanComment, error) {
	if _, err := irapcs.identityReviewActionPlanSvc.GetByID(identityReviewUserID, actionPlanID); err != nil {
		return nil, err
	}

	comment, err := irapcs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	if comment.IdentityReviewActionPlanID != actionPlanID {
		return nil, utility.NewResourceNotFoundError(gorm.ErrRecordNotFound, id)
	}
	return comment, nil
}

func (irapcs *IdentityReviewActionPlanCommentService) GetFullDetailsByID(identityReviewUserID uint, actionPlanID uint, id uint, tx ...utility.DatabaseTransaction) (*dtos.IdentityReviewActionPlanCommentDetailDTO, error) {
	if _, err := irapcs.identityReviewActionPlanSvc.GetByID(identityReviewUserID, actionPlanID); err != nil {
		return nil, err
	}

	var commentDetail dtos.IdentityReviewActionPlanCommentDetailDTO

	comment, err := irapcs.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	if comment.IdentityReviewActionPlanID != actionPlanID {
		return nil, utility.NewResourceNotFoundError(gorm.ErrRecordNotFound, id)
	}

	if err := mapstructure.Decode(comment, &commentDetail); err != nil {
		return nil, err
	}

	commentDetail.CreatedAt = comment.CreatedAt
	return &commentDetail, nil
}

func (irapcs *IdentityReviewActionPlanCommentService) Create(identityReviewUserID uint, actionPlanID uint, currentUserID uint, identityReviewActionPlanCommentFormDTO *dtos.IdentityReviewActionPlanCommentFormDTO) (uint, error) {
	tx := irapcs.repo.BeginTransaction()

	if _, err := irapcs.identityReviewActionPlanSvc.GetByID(identityReviewUserID, actionPlanID); err != nil {
		return 0, err
	}

	comment := &models.IdentityReviewActionPlanComment{
		IdentityReviewActionPlanID: actionPlanID,
		Comment:                    identityReviewActionPlanCommentFormDTO.Comment,
		CreatedByID:                currentUserID,
	}

	id, err := irapcs.repo.Create(comment, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.IdentityReviewActionPlanCommentDocument

	for _, doc := range identityReviewActionPlanCommentFormDTO.Documents {
		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL(
			"identity-access-review", "identity-review-users", strconv.FormatUint(uint64(identityReviewUserID), 10),
			"action-plans", strconv.FormatUint(uint64(actionPlanID), 10),
			"comments", strconv.FormatUint(uint64(id), 10),
			"documents",
		)

		document := &models.IdentityReviewActionPlanCommentDocument{
			IdentityReviewActionPlanCommentID: id,
			URL:                               url,
			Name:                              doc.FileName,
			Size:                              doc.Size,
			Token:                             encryptedFilename,
			OriginalName:                      doc.OriginalName,
			Type:                              doc.Type,
			CreatedByID:                       currentUserID,
		}

		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := irapcs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}

		for _, doc := range documents {
			storagePath := filepath.Join(
				"identity-access-review", "identity-review-users", strconv.FormatUint(uint64(identityReviewUserID), 10),
				"action-plans", strconv.FormatUint(uint64(actionPlanID), 10),
				"comments", strconv.FormatUint(uint64(id), 10),
				"documents",
			)

			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}

	tx.Commit()
	go func() {
		existingComment, err := irapcs.GetFullDetailsByID(identityReviewUserID, actionPlanID, id)
		if err != nil {
			utility.LogError("Failed to fetch identity review action plan comment", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = irapcs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.IdentityReviewActionPlanComment{}.TableName(),
			id,
			existingComment,
			currentUserID,
			"identity_access_review", "identity_review_user", "action_plan_comment", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.IdentityReviewActionPlanComment{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()
	return id, nil
}

func (irapcs *IdentityReviewActionPlanCommentService) Update(identityReviewUserID uint, actionPlanID uint, id uint, currentUserID uint, identityReviewActionPlanCommentFormDTO *dtos.IdentityReviewActionPlanCommentFormDTO) error {
	tx := irapcs.repo.BeginTransaction()

	if _, err := irapcs.identityReviewActionPlanSvc.GetByID(identityReviewUserID, actionPlanID); err != nil {
		return err
	}

	existingComment, err := irapcs.GetByID(identityReviewUserID, actionPlanID, id)
	if err != nil {
		return err
	}

	oldComment, err := irapcs.GetFullDetailsByID(identityReviewUserID, actionPlanID, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingComment.Comment = identityReviewActionPlanCommentFormDTO.Comment
	err = irapcs.repo.Update(existingComment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	var newDocuments []*models.IdentityReviewActionPlanCommentDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.IdentityReviewActionPlanCommentDocument

	for _, doc := range identityReviewActionPlanCommentFormDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {

			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}

			url := utility.GenerateDocumentURL("identity-access-review", "identity-review-users", strconv.FormatUint(uint64(identityReviewUserID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

			document := &models.IdentityReviewActionPlanCommentDocument{
				IdentityReviewActionPlanCommentID: id,
				URL:                               url,
				Name:                              doc.FileName,
				Size:                              doc.Size,
				Token:                             encryptedFileName,
				OriginalName:                      doc.OriginalName,
				Type:                              doc.Type,
				CreatedByID:                       currentUserID,
			}

			newDocuments = append(newDocuments, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			if doc.ID != nil {
				existingDocument, err := irapcs.repo.GetDocumentByID(*doc.ID)
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

	if len(newDocuments) > 0 {
		if err := irapcs.repo.CreateDocuments(newDocuments, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIds) > 0 {
		if err := irapcs.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range newDocuments {
		storagePath := filepath.Join("identity-access-review", "identity-review-users", strconv.FormatUint(uint64(identityReviewUserID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("identity-access-review", "identity-review-users", strconv.FormatUint(uint64(identityReviewUserID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	go func() {
		existingComment, err := irapcs.GetFullDetailsByID(identityReviewUserID, actionPlanID, id)
		if err != nil {
			return
		}

		changedData := make(map[string]interface{})

		oldData := map[string]interface{}{
			"comment": oldComment.Comment,
		}

		newData := map[string]interface{}{
			"comment": existingComment.Comment,
		}

		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}

		if len(changedData) > 0 {
			err := irapcs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.IdentityReviewActionPlanComment{}.TableName(),
				id,
				changedData,
				currentUserID,
				"identity_access_review", "identity_review_user", "action_plan_comment", nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.IdentityReviewActionPlanComment{}.TableName(),
					"recordID":  id,
					"userID":    currentUserID,
					"error":     err.Error(),
				}

				utility.LogError("Failed to log activity", fields)
			}
		}
	}()
	return nil
}

func (irapcs *IdentityReviewActionPlanCommentService) Delete(identityReviewUserID uint, actionPlanID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := irapcs.repo.BeginTransaction()
	if _, err := irapcs.identityReviewActionPlanSvc.GetByID(identityReviewUserID, actionPlanID); err != nil {
		return nil, err
	}

	var failedIDs []uint
	var successfulIDs []uint
	var deletedComments []map[string]interface{}
	for _, id := range ids {
		existingComment, err := irapcs.repo.GetFullDetailsByID(id)
		if err != nil || existingComment.IdentityReviewActionPlanID != actionPlanID {
			failedIDs = append(failedIDs, id)
			continue
		}

		deletedComments = append(deletedComments, map[string]interface{}{
			"id":                     id,
			"identity_review_action": actionPlanID,
		})

		var documentIDs []*uint
		for _, document := range existingComment.Documents {
			documentIDs = append(documentIDs, &document.ID)
		}
		if len(documentIDs) > 0 {
			if err := irapcs.repo.DeleteDocuments(documentIDs, tx); err != nil {
				failedIDs = append(failedIDs, id)
				continue
			}
		}

		if err := irapcs.repo.Delete(id, tx); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		successfulIDs = append(successfulIDs, id)
	}

	if len(failedIDs) > 0 {
		tx.Rollback()
		return failedIDs, nil
	}

	tx.Commit()
	if len(successfulIDs) > 0 {
		go func() {
			err := irapcs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.IdentityReviewActionPlanComment{}.TableName(),
				successfulIDs,
				deletedComments,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}
	return failedIDs, nil
}

func (irapcs *IdentityReviewActionPlanCommentService) GetDocumentPath(identityReviewUserID, actionPlanID, id uint, token string) (string, error) {
	existingComment, err := irapcs.GetByID(identityReviewUserID, actionPlanID, id)
	if err != nil {
		return "", err
	}
	if len(existingComment.Documents) == 0 {
		return "", fmt.Errorf("no documents found for the identity review action plan comment")
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	storagePath := filepath.Join("identity-access-review", "identity-review-users", strconv.FormatUint(uint64(identityReviewUserID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}
