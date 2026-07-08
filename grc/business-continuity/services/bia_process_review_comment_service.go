package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"jethur/config"
	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	businessContinuityRepositories "jethur/grc/business-continuity/repositories"
	logServices "jethur/grc/logs/services"
	masterService "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type BiaProcessReviewCommentServiceInterface interface {
	GetAll(biaProcessID uint, currentUserID uint, params *businessContinuityDTOs.BiaProcessReviewCommentListQueryParamsDTO) (*businessContinuityDTOs.BiaProcessReviewCommentListResponseDTO, error)
	GetByID(biaProcessID uint, commentID uint) (*models.BiaProcessReviewComment, error)
	GetFullDetailsByID(biaProcessID uint, commentID uint) (*businessContinuityDTOs.BiaProcessReviewCommentDetailDTO, error)
	Create(biaProcessID uint, currentUserID uint, dto *businessContinuityDTOs.BiaProcessReviewCommentFormDTO) (uint, error)
	Update(biaProcessID uint, id uint, currentUserID uint, dto *businessContinuityDTOs.BiaProcessReviewCommentFormDTO) error
	Delete(biaProcessID uint, ids []uint, currentUserID uint) ([]uint, error)
	Resolve(biaProcessID uint, commentID uint, currentUserID uint) error
	GetDocumentPath(biaProcessID uint, commentID uint, token string) (string, error)
}

type BiaProcessReviewCommentService struct {
	repo             businessContinuityRepositories.BiaProcessReviewCommentRepositoryInterface
	userSvc          organizationServices.UserServiceInterface
	activityLogSvc   logServices.ActivityLogServiceInterface
	biaProcessSvc    BiaProcessServiceInterface
	commentStatusSvc masterService.CommentStatusServiceInterface
}

func NewBiaProcessReviewCommentService(
	repo businessContinuityRepositories.BiaProcessReviewCommentRepositoryInterface,
	userSvc organizationServices.UserServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	biaProcessSvc BiaProcessServiceInterface,
	commentStatusSvc masterService.CommentStatusServiceInterface,
) BiaProcessReviewCommentServiceInterface {
	return &BiaProcessReviewCommentService{
		repo:             repo,
		userSvc:          userSvc,
		activityLogSvc:   activityLogSvc,
		biaProcessSvc:    biaProcessSvc,
		commentStatusSvc: commentStatusSvc,
	}
}

func (svc *BiaProcessReviewCommentService) GetAll(biaProcessID uint, currentUserID uint, params *businessContinuityDTOs.BiaProcessReviewCommentListQueryParamsDTO) (*businessContinuityDTOs.BiaProcessReviewCommentListResponseDTO, error) {
	_, err := svc.biaProcessSvc.GetByID(biaProcessID)
	if err != nil {
		return nil, err
	}

	comments, err := svc.repo.GetAll(biaProcessID, currentUserID, params)
	if err != nil {
		return nil, err
	}

	dtoMap := make(map[uint]*businessContinuityDTOs.BiaProcessReviewCommentListDTO)
	var roots []*businessContinuityDTOs.BiaProcessReviewCommentListDTO

	for _, comment := range comments.Comments {
		dtoMap[comment.ID] = comment
	}

	for _, comment := range comments.Comments {
		dto := dtoMap[comment.ID]
		if comment.ParentCommentID != nil {
			if parent, ok := dtoMap[*comment.ParentCommentID]; ok {
				parent.Children = append(parent.Children, dto)
			}
		} else {
			roots = append(roots, dto)
		}
	}

	return &businessContinuityDTOs.BiaProcessReviewCommentListResponseDTO{
		Comments:      roots,
		OpenCount:     comments.OpenCount,
		ResolvedCount: comments.ResolvedCount,
	}, nil
}

func (svc *BiaProcessReviewCommentService) GetByID(biaProcessID uint, commentID uint) (*models.BiaProcessReviewComment, error) {
	if _, err := svc.biaProcessSvc.GetByID(biaProcessID); err != nil {
		return nil, err
	}

	comment, err := svc.repo.GetByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, biaProcessID)
		}
		return nil, err
	}
	return comment, nil
}

func (svc *BiaProcessReviewCommentService) GetFullDetailsByID(biaProcessID uint, commentID uint) (*businessContinuityDTOs.BiaProcessReviewCommentDetailDTO, error) {
	if _, err := svc.biaProcessSvc.GetByID(biaProcessID); err != nil {
		fmt.Printf("BiaProcess not found for biaProcessID %d", biaProcessID)
		return nil, err
	}

	var commentDetail businessContinuityDTOs.BiaProcessReviewCommentDetailDTO

	comment, err := svc.repo.GetFullDetailsByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, biaProcessID)
		}
		return nil, err
	}

	err = mapstructure.Decode(comment, &commentDetail)
	if err != nil {
		fmt.Printf("Error decoding comment: %v\n", err)
		return nil, err
	}

	commentDetail.CreatedAt = comment.CreatedAt

	if len(comment.Documents) > 0 {
		var documentsDTO []*businessContinuityDTOs.BiaProcessReviewCommentDocumentDTO
		for _, doc := range comment.Documents {
			documentDTO := &businessContinuityDTOs.BiaProcessReviewCommentDocumentDTO{
				ID:           doc.ID,
				URL:          *doc.URL,
				Name:         doc.Name,
				Size:         *doc.Size,
				Token:        *doc.Token,
				OriginalName: *doc.OriginalName,
				Type:         *doc.Type,
			}
			documentsDTO = append(documentsDTO, documentDTO)
		}
		commentDetail.Documents = documentsDTO
	}

	return &commentDetail, nil
}

func (svc *BiaProcessReviewCommentService) Create(biaProcessID uint, currentUserID uint, commentFormDTO *businessContinuityDTOs.BiaProcessReviewCommentFormDTO) (uint, error) {
	tx := svc.repo.BeginTransaction()

	biaProcess, err := svc.biaProcessSvc.GetByID(biaProcessID)
	if err != nil {
		return 0, err
	}

	commentStatus, err := svc.commentStatusSvc.GetByType("open")
	if err != nil {
		return 0, err
	}

	comment := &models.BiaProcessReviewComment{
		BiaProcessID:    biaProcessID,
		Comment:         commentFormDTO.Comment,
		CommentStatusID: commentStatus.ID,
		Level:           biaProcess.NextReviewLevel,
		ParentCommentID: commentFormDTO.ParentCommentID,
		CreatedBYID:     currentUserID,
	}

	id, err := svc.repo.Create(comment, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.BiaProcessReviewCommentDocument
	for _, doc := range commentFormDTO.Documents {
		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL("business-continuity", "bia-process", strconv.FormatUint(uint64(biaProcessID), 10), "review-comments", strconv.FormatUint(uint64(id), 10), "documents")

		document := &models.BiaProcessReviewCommentDocument{
			BiaProcessReviewCommentID: id,
			URL:                       &url,
			Name:                      doc.FileName,
			Size:                      &doc.Size,
			Token:                     &encryptedFilename,
			OriginalName:              &doc.OriginalName,
			Type:                      &doc.Type,
			CreatedByID:               currentUserID,
			CreatedAt:                 time.Now(),
		}
		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := svc.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		for _, doc := range documents {
			storagePath := filepath.Join("business-continuity", "bia-process", strconv.FormatUint(uint64(biaProcessID), 10), "review-comments", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}

	tx.Commit()
	go func() {
		existingComment, err := svc.GetFullDetailsByID(biaProcessID, id)
		if err != nil {
			utility.LogError("Failed to fetch comment", map[string]interface{}{"error": err.Error()})
			return
		}

		err = svc.activityLogSvc.LogActivity(
			models.CreateAction,
			models.BiaProcessReviewComment{}.TableName(),
			id,
			existingComment,
			currentUserID,
			"business-continuity", "bia_process", "bia_process_review_comment", nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.BiaProcessReviewComment{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}
	}()
	return id, nil
}

func (svc *BiaProcessReviewCommentService) Update(biaProcessID uint, id uint, currentUserID uint, commentFormDTO *businessContinuityDTOs.BiaProcessReviewCommentFormDTO) error {
	tx := svc.repo.BeginTransaction()

	if _, err := svc.biaProcessSvc.GetByID(biaProcessID); err != nil {
		return err
	}

	existingComment, err := svc.GetByID(biaProcessID, id)
	if err != nil {
		return err
	}

	oldComment, err := svc.GetFullDetailsByID(biaProcessID, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingComment.Comment = commentFormDTO.Comment
	existingComment.IsEdited = true
	now := time.Now()
	existingComment.UpdatedBYID = &currentUserID
	existingComment.UpdatedAt = now
	if err := svc.repo.Update(existingComment, tx); err != nil {
		tx.Rollback()
		return err
	}

	var newDocuments []*models.BiaProcessReviewCommentDocument
	var deleteDocumentIDs []*uint
	var documentsToDelete []*models.BiaProcessReviewCommentDocument

	for _, doc := range commentFormDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {
			encryptedFilename, err := utility.EncryptData([]byte(doc.FileName))
			if err != nil {
				tx.Rollback()
				return err
			}

			url := utility.GenerateDocumentURL("business-continuity", "bia-process", strconv.FormatUint(uint64(biaProcessID), 10), "review-comments", strconv.FormatUint(uint64(id), 10), "documents")

			document := &models.BiaProcessReviewCommentDocument{
				BiaProcessReviewCommentID: id,
				URL:                       &url,
				Name:                      doc.FileName,
				Size:                      &doc.Size,
				Token:                     &encryptedFilename,
				OriginalName:              &doc.OriginalName,
				Type:                      &doc.Type,
				CreatedByID:               currentUserID,
				CreatedAt:                 now,
			}
			newDocuments = append(newDocuments, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile && doc.ID != nil {
			existingDocument, err := svc.repo.GetDocumentByID(*doc.ID)
			if err != nil {
				tx.Rollback()
				return err
			}
			deleteDocumentIDs = append(deleteDocumentIDs, &existingDocument.ID)
			documentsToDelete = append(documentsToDelete, existingDocument)
		}
	}

	if len(newDocuments) > 0 {
		if err := svc.repo.CreateDocuments(newDocuments, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	if len(deleteDocumentIDs) > 0 {
		if err := svc.repo.DeleteDocuments(deleteDocumentIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range newDocuments {
		storagePath := filepath.Join("business-continuity", "bia-process", strconv.FormatUint(uint64(biaProcessID), 10), "review-comments", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}
	for _, doc := range documentsToDelete {
		storagePath := filepath.Join("business-continuity", "bia-process", strconv.FormatUint(uint64(biaProcessID), 10), "review-comments", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	go func() {
		newComment, err := svc.GetFullDetailsByID(biaProcessID, id)
		if err != nil {
			return
		}

		changedData := make(map[string]interface{})
		oldData := map[string]interface{}{"comment": oldComment.Comment}
		newData := map[string]interface{}{"comment": newComment.Comment}

		for k, v := range utility.CompareFieldChanges(oldData, newData) {
			changedData[k] = v
		}

		if len(changedData) > 0 {
			err := svc.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.BiaProcessReviewComment{}.TableName(),
				id,
				changedData,
				currentUserID,
				"business-continuity", "bia_process", "bia_process_review_comment", nil,
			)
			if err != nil {
				utility.LogError("Failed to log activity", map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.BiaProcessReviewComment{}.TableName(),
					"recordID":  id,
					"userID":    currentUserID,
					"error":     err.Error(),
				})
			}
		}
	}()
	return nil
}

func (svc *BiaProcessReviewCommentService) Delete(biaProcessID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := svc.repo.BeginTransaction()
	if _, err := svc.biaProcessSvc.GetByID(biaProcessID); err != nil {
		return nil, err
	}

	var failedIDs []uint
	var successfulIDs []uint
	var deletedComments []map[string]interface{}
	var documentsToDelete []*models.BiaProcessReviewCommentDocument

	for _, id := range ids {
		comment, err := svc.repo.GetFullDetailsByID(id)
		if err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		deletedComments = append(deletedComments, map[string]interface{}{
			"id":          id,
			"bia_process": biaProcessID,
		})

		if len(comment.Documents) > 0 {
			documentsToDelete = append(documentsToDelete, comment.Documents...)
		}

		err = svc.repo.Delete(id, tx)
		if err != nil {
			failedIDs = append(failedIDs, id)
		} else {
			successfulIDs = append(successfulIDs, id)
		}
	}

	if len(failedIDs) > 0 {
		tx.Rollback()
		return failedIDs, nil
	}

	for _, doc := range documentsToDelete {
		storagePath := filepath.Join("business-continuity", "bia-process", strconv.FormatUint(uint64(biaProcessID), 10), "review-comments", strconv.FormatUint(uint64(doc.BiaProcessReviewCommentID), 10), "documents")
		if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return failedIDs, err
		}
	}

	tx.Commit()
	if len(successfulIDs) > 0 {
		go func() {
			err := svc.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.BiaProcessReviewComment{}.TableName(),
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

func (svc *BiaProcessReviewCommentService) Resolve(biaProcessID uint, commentID uint, currentUserID uint) error {
	tx := svc.repo.BeginTransaction()

	if _, err := svc.biaProcessSvc.GetByID(biaProcessID); err != nil {
		return err
	}

	existingComment, err := svc.GetByID(biaProcessID, commentID)
	if err != nil {
		return err
	}

	commentStatus, err := svc.commentStatusSvc.GetByType("resolved")
	if err != nil {
		return err
	}

	now := time.Now()
	existingComment.CommentStatusID = commentStatus.ID
	existingComment.UpdatedBYID = &currentUserID
	existingComment.UpdatedAt = now
	if err := svc.repo.Update(existingComment, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (svc *BiaProcessReviewCommentService) GetDocumentPath(biaProcessID uint, commentID uint, token string) (string, error) {
	comment, err := svc.repo.GetFullDetailsByID(commentID)
	if err != nil {
		return "", err
	}
	if len(comment.Documents) == 0 {
		return "", fmt.Errorf("no documents found for the bia process review comment")
	}
	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}
	storagePath := filepath.Join("business-continuity", "bia-process", strconv.FormatUint(uint64(biaProcessID), 10), "review-comments", strconv.FormatUint(uint64(commentID), 10), "documents")
	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}
