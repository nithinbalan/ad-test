package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"jethur/config"
	logServices "jethur/grc/logs/services"
	masterService "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/grc/risk/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type RiskResponseReviewCommentServiceInterface interface {
	GetAll(riskResponseID uint, currentUserID uint, params *riskDTOs.RiskResponseReviewCommentListQueryParamsDTO) (*riskDTOs.RiskResponseReviewCommentListResponseDTO, error)
	GetByID(riskResponseID uint, commentID uint) (*models.RiskResponseReviewComment, error)
	GetFullDetailsByID(riskResponseID uint, commentID uint) (*riskDTOs.RiskResponseReviewCommentDetailDTO, error)
	Create(riskResponseID uint, currentUserID uint, dto *riskDTOs.RiskResponseReviewCommentFormDTO) (uint, error)
	Update(riskResponseID uint, id uint, currentUserID uint, dto *riskDTOs.RiskResponseReviewCommentFormDTO) error
	Delete(riskResponseID uint, ids []uint, currentUserID uint) ([]uint, error)
	Resolve(riskResponseID uint, commentID uint, currentUserID uint) error
	GetDocumentPath(riskResponseID uint, commentID uint, token string) (string, error)
}

type RiskResponseReviewCommentService struct {
	repo             repositories.RiskResponseReviewCommentRepositoryInterface
	userSvc          organizationServices.UserServiceInterface
	activityLogSvc   logServices.ActivityLogServiceInterface
	riskResponseSvc  RiskResponseServiceInterface
	commentStatusSvc masterService.CommentStatusServiceInterface
}

func NewRiskResponseReviewCommentService(
	repo repositories.RiskResponseReviewCommentRepositoryInterface,
	userSvc organizationServices.UserServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	riskResponseSvc RiskResponseServiceInterface,
	commentStatusSvc masterService.CommentStatusServiceInterface,
) RiskResponseReviewCommentServiceInterface {
	return &RiskResponseReviewCommentService{
		repo:             repo,
		userSvc:          userSvc,
		activityLogSvc:   activityLogSvc,
		riskResponseSvc:  riskResponseSvc,
		commentStatusSvc: commentStatusSvc,
	}
}

func (rrcs *RiskResponseReviewCommentService) GetAll(riskResponseID uint, currentUserID uint, params *riskDTOs.RiskResponseReviewCommentListQueryParamsDTO) (*riskDTOs.RiskResponseReviewCommentListResponseDTO, error) {
	_, err := rrcs.riskResponseSvc.GetByID(riskResponseID)
	if err != nil {
		return nil, err
	}

	comments, err := rrcs.repo.GetAll(riskResponseID, currentUserID, params)
	if err != nil {
		return nil, err
	}
	dtoMap := make(map[uint]*riskDTOs.RiskResponseReviewCommentListDTO)
	var roots []*riskDTOs.RiskResponseReviewCommentListDTO

	for _, comment := range comments.Comments {
		dto := &riskDTOs.RiskResponseReviewCommentListDTO{
			ID:                          comment.ID,
			Comment:                     comment.Comment,
			IsEdited:                    comment.IsEdited,
			CommentStatusTitleEn:        comment.CommentStatusTitleEn,
			CommentStatusTitleAr:        comment.CommentStatusTitleAr,
			CommentStatusType:           comment.CommentStatusType,
			CommentStatusColorCode:      comment.CommentStatusColorCode,
			RiskResponseID:              comment.RiskResponseID,
			Level:                       comment.Level,
			CreatedUserID:               comment.CreatedUserID,
			CreatedUserFirstNameEn:      comment.CreatedUserFirstNameEn,
			CreatedUserFirstNameAr:      comment.CreatedUserFirstNameAr,
			CreatedUserLastNameEn:       comment.CreatedUserLastNameEn,
			CreatedUserLastNameAr:       comment.CreatedUserLastNameAr,
			CreatedUserInitial:          comment.CreatedUserInitial,
			CreatedUserColor:            comment.CreatedUserColor,
			CreatedUserEmail:            comment.CreatedUserEmail,
			CreatedUserProfileImagePath: comment.CreatedUserProfileImagePath,
			CreatedAt:                   comment.CreatedAt,
			ParentCommentID:             comment.ParentCommentID,
		}
		dtoMap[comment.ID] = dto
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

	response := &riskDTOs.RiskResponseReviewCommentListResponseDTO{
		Comments:      roots,
		OpenCount:     comments.OpenCount,
		ResolvedCount: comments.ResolvedCount,
	}

	return response, nil
}

func (rrcs *RiskResponseReviewCommentService) GetByID(riskResponseID uint, commentID uint) (*models.RiskResponseReviewComment, error) {
	if _, err := rrcs.riskResponseSvc.GetByID(riskResponseID); err != nil {
		return nil, err
	}

	comment, err := rrcs.repo.GetByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, riskResponseID)
		}
		return nil, err
	}
	return comment, nil
}

func (rrcs *RiskResponseReviewCommentService) GetFullDetailsByID(riskResponseID uint, commentID uint) (*riskDTOs.RiskResponseReviewCommentDetailDTO, error) {
	if _, err := rrcs.riskResponseSvc.GetByID(riskResponseID); err != nil {
		fmt.Printf("Risk Response not found for riskResponseID %d", riskResponseID)
		return nil, err
	}

	var commentDetail riskDTOs.RiskResponseReviewCommentDetailDTO

	comment, err := rrcs.repo.GetFullDetailsByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, riskResponseID)
		}
		return nil, err
	}

	err = mapstructure.Decode(comment, &commentDetail)
	if err != nil {
		fmt.Printf("Error decoding comment: %v\n", err)
		return nil, err
	}

	commentDetail.CreatedAt = comment.CreatedAt

	// Map documents
	if len(comment.Documents) > 0 {
		var documentsDTO []*riskDTOs.RiskResponseReviewCommentDocumentDTO
		for _, doc := range comment.Documents {
			documentDTO := &riskDTOs.RiskResponseReviewCommentDocumentDTO{
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

func (rrcs *RiskResponseReviewCommentService) Create(riskResponseID uint, currentUserID uint, commentFormDTO *riskDTOs.RiskResponseReviewCommentFormDTO) (uint, error) {
	tx := rrcs.repo.BeginTransaction()

	riskResponse, err := rrcs.riskResponseSvc.GetByID(riskResponseID)
	if err != nil {
		return 0, err
	}

	commentStatus, err := rrcs.commentStatusSvc.GetByType("open")
	if err != nil {
		return 0, err
	}

	comment := &models.RiskResponseReviewComment{
		RiskResponseID:  riskResponseID,
		Comment:         commentFormDTO.Comment,
		CommentStatusID: commentStatus.ID,
		Level:           riskResponse.NextReviewLevel,
		ParentCommentID: commentFormDTO.ParentCommentID,
		CreatedBYID:     currentUserID,
	}

	id, err := rrcs.repo.Create(comment, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.RiskResponseReviewCommentDocument
	for _, doc := range commentFormDTO.Documents {
		// Encrypt the filename
		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		// Construct the URL for the document
		url := utility.GenerateDocumentURL("risk", "risk-response", strconv.FormatUint(uint64(riskResponseID), 10), "review-comments", strconv.FormatUint(uint64(id), 10), "documents")

		document := &models.RiskResponseReviewCommentDocument{
			RiskResponseReviewCommentID: id,
			URL:                         &url,
			Name:                        doc.FileName,
			Size:                        &doc.Size,
			Token:                       &encryptedFilename,
			OriginalName:                &doc.OriginalName,
			Type:                        &doc.Type,
			CreatedByID:                 currentUserID,
			CreatedAt:                   time.Now(),
		}

		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := rrcs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		for _, doc := range documents {
			storagePath := filepath.Join("risk", "risk-response", strconv.FormatUint(uint64(riskResponseID), 10), "review-comments", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}

	tx.Commit()
	go func() {
		existingComment, err := rrcs.GetFullDetailsByID(riskResponseID, id)
		if err != nil {
			utility.LogError("Failed to fetch comment", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = rrcs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.RiskResponseReviewComment{}.TableName(),
			id,
			existingComment,
			currentUserID,
			"risk", "risk_register", "risk_response_review_comment", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.RiskResponseReviewComment{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()
	return id, nil
}

func (rrcs *RiskResponseReviewCommentService) Update(riskResponseID uint, id uint, currentUserID uint, commentFormDTO *riskDTOs.RiskResponseReviewCommentFormDTO) error {
	tx := rrcs.repo.BeginTransaction()

	if _, err := rrcs.riskResponseSvc.GetByID(riskResponseID); err != nil {
		return err
	}

	existingComment, err := rrcs.GetByID(riskResponseID, id)
	if err != nil {
		return err
	}

	oldComment, err := rrcs.GetFullDetailsByID(riskResponseID, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingComment.Comment = commentFormDTO.Comment
	existingComment.IsEdited = true
	now := time.Now()
	existingComment.UpdatedBYID = &currentUserID
	existingComment.UpdatedAt = now
	err = rrcs.repo.Update(existingComment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Handle documents
	var newDocuments []*models.RiskResponseReviewCommentDocument
	var deleteDocumentIDs []*uint
	var documentsToDelete []*models.RiskResponseReviewCommentDocument

	for _, doc := range commentFormDTO.Documents {
		// NEW FILE
		if doc.IsNewFile != nil && *doc.IsNewFile {
			encryptedFilename, err := utility.EncryptData([]byte(doc.FileName))
			if err != nil {
				tx.Rollback()
				return err
			}

			url := utility.GenerateDocumentURL(
				"risk",
				"risk-response",
				strconv.FormatUint(uint64(riskResponseID), 10),
				"review-comments",
				strconv.FormatUint(uint64(id), 10),
				"documents",
			)

			document := &models.RiskResponseReviewCommentDocument{
				RiskResponseReviewCommentID: id,
				URL:                         &url,
				Name:                        doc.FileName,
				Size:                        &doc.Size,
				Token:                       &encryptedFilename,
				OriginalName:                &doc.OriginalName,
				Type:                        &doc.Type,
				CreatedByID:                 currentUserID,
				CreatedAt:                   now,
			}

			newDocuments = append(newDocuments, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile && doc.ID != nil {
			// DELETE FILE
			existingDocument, err := rrcs.repo.GetDocumentByID(*doc.ID)
			if err != nil {
				tx.Rollback()
				return err
			}

			deleteDocumentIDs = append(deleteDocumentIDs, &existingDocument.ID)
			documentsToDelete = append(documentsToDelete, existingDocument)
		}
	}

	if len(newDocuments) > 0 {
		if err := rrcs.repo.CreateDocuments(newDocuments, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIDs) > 0 {
		if err := rrcs.repo.DeleteDocuments(deleteDocumentIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// File system operations
	for _, doc := range newDocuments {
		storagePath := filepath.Join("risk", "risk-response", strconv.FormatUint(uint64(riskResponseID), 10), "review-comments", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range documentsToDelete {
		storagePath := filepath.Join("risk", "risk-response", strconv.FormatUint(uint64(riskResponseID), 10), "review-comments", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	go func() {
		existingComment, err := rrcs.GetFullDetailsByID(riskResponseID, id)
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
			err := rrcs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.RiskResponseReviewComment{}.TableName(),
				id,
				changedData,
				currentUserID,
				"risk", "risk_register", "risk_response_review_comment", nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.RiskResponseReviewComment{}.TableName(),
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

func (rrcs *RiskResponseReviewCommentService) Delete(riskResponseID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := rrcs.repo.BeginTransaction()
	if _, err := rrcs.riskResponseSvc.GetByID(riskResponseID); err != nil {
		return nil, err
	}

	var failedIDs []uint
	var successfulIDs []uint
	var deletedComments []map[string]interface{}
	var documentsToDelete []*models.RiskResponseReviewCommentDocument

	for _, id := range ids {
		comment, err := rrcs.repo.GetFullDetailsByID(id)
		if err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		deletedComments = append(deletedComments, map[string]interface{}{
			"id":            id,
			"risk_response": riskResponseID,
		})

		// Collect documents for deletion
		if len(comment.Documents) > 0 {
			documentsToDelete = append(documentsToDelete, comment.Documents...)
		}

		err = rrcs.repo.Delete(id, tx)
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

	// Delete files from filesystem after successful DB deletion
	for _, doc := range documentsToDelete {
		storagePath := filepath.Join("risk", "risk-response", strconv.FormatUint(uint64(riskResponseID), 10), "review-comments", strconv.FormatUint(uint64(doc.RiskResponseReviewCommentID), 10), "documents")
		if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return failedIDs, err
		}
	}

	tx.Commit()
	if len(successfulIDs) > 0 {
		go func() {
			err := rrcs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.RiskResponseReviewComment{}.TableName(),
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

func (rrcs *RiskResponseReviewCommentService) Resolve(riskResponseID uint, commentID uint, currentUserID uint) error {
	tx := rrcs.repo.BeginTransaction()

	if _, err := rrcs.riskResponseSvc.GetByID(riskResponseID); err != nil {
		return err
	}

	existingComment, err := rrcs.GetByID(riskResponseID, commentID)
	if err != nil {
		return err
	}

	commentStatus, err := rrcs.commentStatusSvc.GetByType("resolved")
	if err != nil {
		return err
	}

	now := time.Now()

	existingComment.CommentStatusID = commentStatus.ID
	existingComment.UpdatedBYID = &currentUserID
	existingComment.UpdatedAt = now
	err = rrcs.repo.Update(existingComment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (rrcs *RiskResponseReviewCommentService) GetDocumentPath(riskResponseID uint, commentID uint, token string) (string, error) {
	comment, err := rrcs.repo.GetFullDetailsByID(commentID)
	if err != nil {
		return "", err
	}
	if len(comment.Documents) == 0 {
		return "", fmt.Errorf("no documents found for the risk response review comment")
	}
	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}
	storagePath := filepath.Join("risk", "risk-response", strconv.FormatUint(uint64(riskResponseID), 10), "review-comments", strconv.FormatUint(uint64(commentID), 10), "documents")
	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}
