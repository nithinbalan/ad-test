package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"jethur/config"

	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	projectDTOs "jethur/grc/project/dtos"
	"jethur/grc/project/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ProjectActionPlanCommentServiceInterface interface {
	GetAll(projectID uint, actionPlanID uint) ([]*projectDTOs.ProjectActionPlanCommentListDTO, error)
	GetByID(projectID uint, actionPlanID uint, id uint) (*models.ProjectActionPlanComment, error)
	GetFullDetailsByID(projectID uint, actionPlanID uint, id uint, tx ...utility.DatabaseTransaction) (*projectDTOs.ProjectActionPlanCommentDetailDTO, error)
	Create(projectID uint, actionPlanID uint, currentUserID uint, dto *projectDTOs.ProjectActionPlanCommentFormDTO) (uint, error)
	Update(projectID uint, actionPlanID uint, id uint, currentUserID uint, dto *projectDTOs.ProjectActionPlanCommentFormDTO) error
	Delete(projectID uint, actionPlanID uint, ids []uint, currentUserID uint) ([]uint, error)

	GetDocumentPath(projectID, actionPlanID, id uint, token string) (string, error)
}

type ProjectActionPlanCommentService struct {
	repo                 repositories.ProjectActionPlanCommentRepositoryInterface
	userSvc              organizationServices.UserServiceInterface
	activityLogSvc       logServices.ActivityLogServiceInterface
	projectActionPlanSvc ProjectActionPlanServiceInterface
}

func NewProjectActionPlanCommentService(repo repositories.ProjectActionPlanCommentRepositoryInterface, userSvc organizationServices.UserServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	projectActionPlanSvc ProjectActionPlanServiceInterface,
) ProjectActionPlanCommentServiceInterface {
	return &ProjectActionPlanCommentService{
		repo:                 repo,
		userSvc:              userSvc,
		activityLogSvc:       activityLogSvc,
		projectActionPlanSvc: projectActionPlanSvc,
	}
}

func (caps *ProjectActionPlanCommentService) GetAll(projectID uint, actionPlanID uint) ([]*projectDTOs.ProjectActionPlanCommentListDTO, error) {
	if _, err := caps.projectActionPlanSvc.GetByID(projectID, actionPlanID); err != nil {
		return nil, err
	}

	return caps.repo.GetAll(actionPlanID)
}

func (caps *ProjectActionPlanCommentService) GetByID(projectID uint, actionPlanID uint, id uint) (*models.ProjectActionPlanComment, error) {
	if _, err := caps.projectActionPlanSvc.GetByID(projectID, actionPlanID); err != nil {
		return nil, err
	}

	comment, err := caps.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return comment, nil
}

func (caps *ProjectActionPlanCommentService) GetFullDetailsByID(projectID uint, actionPlanID uint, id uint, tx ...utility.DatabaseTransaction) (*projectDTOs.ProjectActionPlanCommentDetailDTO, error) {
	if _, err := caps.projectActionPlanSvc.GetByID(projectID, actionPlanID); err != nil {
		fmt.Printf("ProjectActionPlan not found for projectID %d, actionPlanID %d\n", projectID, actionPlanID)
		return nil, err
	}

	var commentDetail projectDTOs.ProjectActionPlanCommentDetailDTO

	comment, err := caps.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(comment, &commentDetail)
	if err != nil {
		fmt.Printf("Error decoding comment: %v\n", err)
		return nil, err
	}

	commentDetail.CreatedAt = comment.CreatedAt
	return &commentDetail, nil
}

func (caps *ProjectActionPlanCommentService) Create(projectID uint, actionPlanID uint, currentUserID uint, projectActionPlanCommentFormDTO *projectDTOs.ProjectActionPlanCommentFormDTO) (uint, error) {
	tx := caps.repo.BeginTransaction()

	if _, err := caps.projectActionPlanSvc.GetByID(projectID, actionPlanID); err != nil {
		return 0, err
	}

	comment := &models.ProjectActionPlanComment{
		ProjectActionPlanID: actionPlanID,
		Comment:             projectActionPlanCommentFormDTO.Comment,
		CreatedBYID:         currentUserID,
	}

	id, err := caps.repo.Create(comment)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.ProjectActionPlanCommentDocument

	for _, doc := range projectActionPlanCommentFormDTO.Documents {
		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL(
			"project", "project", strconv.FormatUint(uint64(projectID), 10),
			"action-plans", strconv.FormatUint(uint64(actionPlanID), 10),
			"comments", strconv.FormatUint(uint64(id), 10),
			"documents",
		)

		document := &models.ProjectActionPlanCommentDocument{
			ProjectActionPlanCommentID: id,
			URL:                        url,
			Name:                       doc.FileName,
			Size:                       doc.Size,
			Token:                      encryptedFilename,
			OriginalName:               doc.OriginalName,
			Type:                       doc.Type,
			CreatedBYID:                currentUserID,
		}

		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := caps.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}

		for _, doc := range documents {
			storagePath := filepath.Join(
				"project", "project", strconv.FormatUint(uint64(projectID), 10),
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
		existingComment, err := caps.GetFullDetailsByID(projectID, actionPlanID, id)
		if err != nil {
			utility.LogError("Failed to fetch project action plan comment", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = caps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ProjectActionPlanComment{}.TableName(),
			id,
			existingComment,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ProjectActionPlanComment{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()
	return id, nil
}

func (caps *ProjectActionPlanCommentService) Update(projectID uint, actionPlanID uint, id uint, currentUserID uint, projectActionPlanCommentFormDTO *projectDTOs.ProjectActionPlanCommentFormDTO) error {
	tx := caps.repo.BeginTransaction()

	if _, err := caps.projectActionPlanSvc.GetByID(projectID, actionPlanID); err != nil {
		return err
	}

	existingComment, err := caps.GetByID(projectID, actionPlanID, id)
	if err != nil {
		return err
	}

	oldComment, err := caps.GetFullDetailsByID(projectID, actionPlanID, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingComment.Comment = projectActionPlanCommentFormDTO.Comment
	err = caps.repo.Update(existingComment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	var newDocuments []*models.ProjectActionPlanCommentDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.ProjectActionPlanCommentDocument

	for _, doc := range projectActionPlanCommentFormDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {

			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}

			url := utility.GenerateDocumentURL("project", "project", strconv.FormatUint(uint64(projectID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

			document := &models.ProjectActionPlanCommentDocument{
				ProjectActionPlanCommentID: id,
				URL:                        url,
				Name:                       doc.FileName,
				Size:                       doc.Size,
				Token:                      encryptedFileName,
				OriginalName:               doc.OriginalName,
				Type:                       doc.Type,
				CreatedBYID:                currentUserID,
			}

			newDocuments = append(newDocuments, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			if doc.ID != nil {
				existingDocument, err := caps.repo.GetDocumentByID(*doc.ID)
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
		if err := caps.repo.CreateDocuments(newDocuments, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIds) > 0 {
		if err := caps.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range newDocuments {
		storagePath := filepath.Join("project", "project", strconv.FormatUint(uint64(projectID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("project", "project", strconv.FormatUint(uint64(projectID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	go func() {
		existingComment, err := caps.GetFullDetailsByID(projectID, actionPlanID, id)
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
			err := caps.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.ProjectActionPlanComment{}.TableName(),
				id,
				changedData,
				currentUserID,
				"", "", "", nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.ProjectActionPlanComment{}.TableName(),
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

func (caps *ProjectActionPlanCommentService) Delete(projectID uint, actionPlanID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := caps.repo.BeginTransaction()
	if _, err := caps.projectActionPlanSvc.GetByID(projectID, actionPlanID); err != nil {
		return nil, err
	}

	var failedIDs []uint
	var successfulIDs []uint
	var deletedComments []map[string]interface{}
	for _, id := range ids {
		if _, err := caps.repo.GetFullDetailsByID(id); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		deletedComments = append(deletedComments, map[string]interface{}{
			"id":                  id,
			"project_action_plan": actionPlanID,
		})

		err := caps.repo.Delete(id, tx)
		if err != nil {
			failedIDs = append(failedIDs, id)
		}
	}

	if len(failedIDs) > 0 {
		tx.Rollback()
		return failedIDs, nil
	}

	tx.Commit()
	if len(successfulIDs) > 0 {
		go func() {
			err := caps.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.ProjectActionPlanComment{}.TableName(),
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

func (caps *ProjectActionPlanCommentService) GetDocumentPath(projectID, actionPlanID, id uint, token string) (string, error) {
	existingComment, err := caps.GetByID(projectID, actionPlanID, id)
	if err != nil {
		return "", err
	}
	if len(existingComment.Documents) == 0 {
		return "", fmt.Errorf("no documents found for the project action plan comment")
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	storagePath := filepath.Join("project", "project", strconv.FormatUint(uint64(projectID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}
