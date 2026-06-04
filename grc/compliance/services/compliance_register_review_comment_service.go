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
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ComplianceRegisterReviewCommentServiceInterface interface {
	GetAll(complianceRegisterID uint, currentUserID uint, params *complianceDTOs.ComplianceRegisterReviewCommentListQueryParamsDTO) (*complianceDTOs.ComplianceRegisterReviewCommentListResponseDTO, error)
	GetByID(complianceRegisterID uint, commentID uint) (*models.ComplianceRegisterReviewComment, error)
	GetFullDetailsByID(complianceRegisterID uint, commentID uint) (*complianceDTOs.ComplianceRegisterReviewCommentDetailDTO, error)
	Create(complianceRegisterID uint, currentUserID uint, dto *complianceDTOs.ComplianceRegisterReviewCommentFormDTO) (uint, error)
	Update(complianceRegisterID uint, id uint, currentUserID uint, dto *complianceDTOs.ComplianceRegisterReviewCommentFormDTO) error
	Delete(complianceRegisterID uint, ids []uint, currentUserID uint) ([]uint, error)
	Resolve(complianceRegisterID uint, commentID uint, currentUserID uint) error
	GetDocumentPath(complianceRegisterID uint, commentID uint, token string) (string, error)
}

type ComplianceRegisterReviewCommentService struct {
	repo                  repositories.ComplianceRegisterReviewCommentRepositoryInterface
	userSvc               organizationServices.UserServiceInterface
	activityLogSvc        logServices.ActivityLogServiceInterface
	complianceRegisterSvc ComplianceRegisterServiceInterface
	commentStatusSvc      masterService.CommentStatusServiceInterface
}

func NewComplianceRegisterReviewCommentService(
	repo repositories.ComplianceRegisterReviewCommentRepositoryInterface,
	userSvc organizationServices.UserServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	complianceRegisterSvc ComplianceRegisterServiceInterface,
	commentStatusSvc masterService.CommentStatusServiceInterface,
) ComplianceRegisterReviewCommentServiceInterface {
	return &ComplianceRegisterReviewCommentService{
		repo:                  repo,
		userSvc:               userSvc,
		activityLogSvc:        activityLogSvc,
		complianceRegisterSvc: complianceRegisterSvc,
		commentStatusSvc:      commentStatusSvc,
	}
}

func (crrcs *ComplianceRegisterReviewCommentService) GetAll(complianceRegisterID uint, currentUserID uint, params *complianceDTOs.ComplianceRegisterReviewCommentListQueryParamsDTO) (*complianceDTOs.ComplianceRegisterReviewCommentListResponseDTO, error) {
	_, err := crrcs.complianceRegisterSvc.GetByID(complianceRegisterID)
	if err != nil {
		return nil, err
	}

	comments, err := crrcs.repo.GetAll(complianceRegisterID, currentUserID, params)
	if err != nil {
		return nil, err
	}
	dtoMap := make(map[uint]*complianceDTOs.ComplianceRegisterReviewCommentListDTO)
	var roots []*complianceDTOs.ComplianceRegisterReviewCommentListDTO

	for _, comment := range comments.Comments {
		dto := &complianceDTOs.ComplianceRegisterReviewCommentListDTO{
			ID:                          comment.ID,
			Comment:                     comment.Comment,
			IsEdited:                    comment.IsEdited,
			CommentStatusTitleEn:        comment.CommentStatusTitleEn,
			CommentStatusTitleAr:        comment.CommentStatusTitleAr,
			CommentStatusType:           comment.CommentStatusType,
			CommentStatusColorCode:      comment.CommentStatusColorCode,
			ComplianceRegisterID:        comment.ComplianceRegisterID,
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
			Documents:                   comment.Documents,
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

	response := &complianceDTOs.ComplianceRegisterReviewCommentListResponseDTO{
		Comments:      roots,
		OpenCount:     comments.OpenCount,
		ResolvedCount: comments.ResolvedCount,
	}

	return response, nil
}

func (crrcs *ComplianceRegisterReviewCommentService) GetByID(complianceRegisterID uint, commentID uint) (*models.ComplianceRegisterReviewComment, error) {
	if _, err := crrcs.complianceRegisterSvc.GetByID(complianceRegisterID); err != nil {
		return nil, err
	}

	comment, err := crrcs.repo.GetByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, complianceRegisterID)
		}
		return nil, err
	}
	return comment, nil
}

func (crrcs *ComplianceRegisterReviewCommentService) GetFullDetailsByID(complianceRegisterID uint, commentID uint) (*complianceDTOs.ComplianceRegisterReviewCommentDetailDTO, error) {
	if _, err := crrcs.complianceRegisterSvc.GetByID(complianceRegisterID); err != nil {
		fmt.Printf("Compliance Register not found for complianceRegisterID %d", complianceRegisterID)
		return nil, err
	}

	var commentDetail complianceDTOs.ComplianceRegisterReviewCommentDetailDTO

	comment, err := crrcs.repo.GetFullDetailsByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, complianceRegisterID)
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
		var documentsDTO []*complianceDTOs.ComplianceRegisterReviewCommentDocumentDTO
		for _, doc := range comment.Documents {
			documentDTO := &complianceDTOs.ComplianceRegisterReviewCommentDocumentDTO{
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

func (crrcs *ComplianceRegisterReviewCommentService) Create(complianceRegisterID uint, currentUserID uint, commentFormDTO *complianceDTOs.ComplianceRegisterReviewCommentFormDTO) (uint, error) {
	tx := crrcs.repo.BeginTransaction()

	complianceRegister, err := crrcs.complianceRegisterSvc.GetByID(complianceRegisterID)
	if err != nil {
		return 0, err
	}

	commentStatus, err := crrcs.commentStatusSvc.GetByType("open")
	if err != nil {
		return 0, err
	}

	comment := &models.ComplianceRegisterReviewComment{
		ComplianceRegisterID: complianceRegisterID,
		Comment:               commentFormDTO.Comment,
		CommentStatusID:       commentStatus.ID,
		Level:                 complianceRegister.NextReviewLevel,
		ParentCommentID:       commentFormDTO.ParentCommentID,
		CreatedBYID:           currentUserID,
	}

	id, err := crrcs.repo.Create(comment, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.ComplianceRegisterReviewCommentDocument
	for _, doc := range commentFormDTO.Documents {
		// Encrypt the filename
		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		// Construct the URL for the document
		url := utility.GenerateDocumentURL("compliance", "compliance-register", strconv.FormatUint(uint64(complianceRegisterID), 10), "review-comments", strconv.FormatUint(uint64(id), 10), "documents")

		document := &models.ComplianceRegisterReviewCommentDocument{
			ComplianceRegisterReviewCommentID: id,
			URL:                                &url,
			Name:                               doc.FileName,
			Size:                               &doc.Size,
			Token:                              &encryptedFilename,
			OriginalName:                       &doc.OriginalName,
			Type:                               &doc.Type,
			CreatedByID:                        currentUserID,
			CreatedAt:                          time.Now(),
		}

		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := crrcs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		for _, doc := range documents {
			storagePath := filepath.Join("compliance", "compliance-register", strconv.FormatUint(uint64(complianceRegisterID), 10), "review-comments", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}

	tx.Commit()
	go func() {
		existingComment, err := crrcs.GetFullDetailsByID(complianceRegisterID, id)
		if err != nil {
			utility.LogError("Failed to fetch comment", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = crrcs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ComplianceRegisterReviewComment{}.TableName(),
			id,
			existingComment,
			currentUserID,
			"compliance", "compliance_register", "compliance_register_review_comment", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ComplianceRegisterReviewComment{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()
	return id, nil
}

func (crrcs *ComplianceRegisterReviewCommentService) Update(complianceRegisterID uint, id uint, currentUserID uint, commentFormDTO *complianceDTOs.ComplianceRegisterReviewCommentFormDTO) error {
	tx := crrcs.repo.BeginTransaction()

	if _, err := crrcs.complianceRegisterSvc.GetByID(complianceRegisterID); err != nil {
		return err
	}

	existingComment, err := crrcs.GetByID(complianceRegisterID, id)
	if err != nil {
		return err
	}

	oldComment, err := crrcs.GetFullDetailsByID(complianceRegisterID, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingComment.Comment = commentFormDTO.Comment
	existingComment.IsEdited = true
	now := time.Now()
	existingComment.UpdatedBYID = &currentUserID
	existingComment.UpdatedAt = now
	err = crrcs.repo.Update(existingComment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Handle documents
	var newDocuments []*models.ComplianceRegisterReviewCommentDocument
	var deleteDocumentIDs []*uint
	var documentsToDelete []*models.ComplianceRegisterReviewCommentDocument

	for _, doc := range commentFormDTO.Documents {
		// NEW FILE
		if doc.IsNewFile != nil && *doc.IsNewFile {
			encryptedFilename, err := utility.EncryptData([]byte(doc.FileName))
			if err != nil {
				tx.Rollback()
				return err
			}

			url := utility.GenerateDocumentURL(
				"compliance",
				"compliance-register",
				strconv.FormatUint(uint64(complianceRegisterID), 10),
				"review-comments",
				strconv.FormatUint(uint64(id), 10),
				"documents",
			)

			document := &models.ComplianceRegisterReviewCommentDocument{
				ComplianceRegisterReviewCommentID: id,
				URL:                                &url,
				Name:                               doc.FileName,
				Size:                               &doc.Size,
				Token:                              &encryptedFilename,
				OriginalName:                       &doc.OriginalName,
				Type:                               &doc.Type,
				CreatedByID:                        currentUserID,
				CreatedAt:                          now,
			}

			newDocuments = append(newDocuments, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile && doc.ID != nil {
			// DELETE FILE
			existingDocument, err := crrcs.repo.GetDocumentByID(*doc.ID)
			if err != nil {
				tx.Rollback()
				return err
			}

			deleteDocumentIDs = append(deleteDocumentIDs, &existingDocument.ID)
			documentsToDelete = append(documentsToDelete, existingDocument)
		}
	}

	if len(newDocuments) > 0 {
		if err := crrcs.repo.CreateDocuments(newDocuments, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIDs) > 0 {
		if err := crrcs.repo.DeleteDocuments(deleteDocumentIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// File system operations
	for _, doc := range newDocuments {
		storagePath := filepath.Join("compliance", "compliance-register", strconv.FormatUint(uint64(complianceRegisterID), 10), "review-comments", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range documentsToDelete {
		storagePath := filepath.Join("compliance", "compliance-register", strconv.FormatUint(uint64(complianceRegisterID), 10), "review-comments", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	go func() {
		existingComment, err := crrcs.GetFullDetailsByID(complianceRegisterID, id)
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
			err := crrcs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.ComplianceRegisterReviewComment{}.TableName(),
				id,
				changedData,
				currentUserID,
				"compliance", "compliance_register", "compliance_register_review_comment", nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.ComplianceRegisterReviewComment{}.TableName(),
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

func (crrcs *ComplianceRegisterReviewCommentService) Delete(complianceRegisterID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := crrcs.repo.BeginTransaction()
	if _, err := crrcs.complianceRegisterSvc.GetByID(complianceRegisterID); err != nil {
		return nil, err
	}

	var failedIDs []uint
	var successfulIDs []uint
	var deletedComments []map[string]interface{}
	var documentsToDelete []*models.ComplianceRegisterReviewCommentDocument

	for _, id := range ids {
		comment, err := crrcs.repo.GetFullDetailsByID(id)
		if err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		deletedComments = append(deletedComments, map[string]interface{}{
			"id":                    id,
			"compliance_register_id": complianceRegisterID,
		})

		// Collect documents for deletion
		if len(comment.Documents) > 0 {
			documentsToDelete = append(documentsToDelete, comment.Documents...)
		}

		err = crrcs.repo.Delete(id, tx)
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
		storagePath := filepath.Join("compliance", "compliance-register", strconv.FormatUint(uint64(complianceRegisterID), 10), "review-comments", strconv.FormatUint(uint64(doc.ComplianceRegisterReviewCommentID), 10), "documents")
		if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return failedIDs, err
		}
	}

	tx.Commit()
	if len(successfulIDs) > 0 {
		go func() {
			err := crrcs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.ComplianceRegisterReviewComment{}.TableName(),
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

func (crrcs *ComplianceRegisterReviewCommentService) Resolve(complianceRegisterID uint, commentID uint, currentUserID uint) error {
	tx := crrcs.repo.BeginTransaction()

	if _, err := crrcs.complianceRegisterSvc.GetByID(complianceRegisterID); err != nil {
		return err
	}

	existingComment, err := crrcs.GetByID(complianceRegisterID, commentID)
	if err != nil {
		return err
	}

	commentStatus, err := crrcs.commentStatusSvc.GetByType("resolved")
	if err != nil {
		return err
	}

	now := time.Now()

	existingComment.CommentStatusID = commentStatus.ID
	existingComment.UpdatedBYID = &currentUserID
	existingComment.UpdatedAt = now
	err = crrcs.repo.Update(existingComment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (crrcs *ComplianceRegisterReviewCommentService) GetDocumentPath(complianceRegisterID uint, commentID uint, token string) (string, error) {
	comment, err := crrcs.repo.GetFullDetailsByID(commentID)
	if err != nil {
		return "", err
	}
	if len(comment.Documents) == 0 {
		return "", fmt.Errorf("no documents found for the compliance register review comment")
	}
	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}
	storagePath := filepath.Join("compliance", "compliance-register", strconv.FormatUint(uint64(complianceRegisterID), 10), "review-comments", strconv.FormatUint(uint64(commentID), 10), "documents")
	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}
