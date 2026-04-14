package services

import (
	"errors"
	"fmt"
	"jethur/config"
	logServices "jethur/grc/logs/services"
	masterService "jethur/grc/masters/services"
	"jethur/grc/models"
	"jethur/grc/organization/dtos"
	"jethur/grc/organization/repositories"
	"jethur/utility"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type InternalCommitteeReviewCommentServiceInterface interface {
	GetAll(committeeID uint, currentUserID uint, params *dtos.InternalCommitteeReviewCommentListQueryParamsDTO) (*dtos.InternalCommitteeReviewCommentListResponseDTO, error)
	GetByID(committeeID uint, commentID uint) (*models.InternalCommitteeReviewComment, error)
	GetFullDetailsByID(committeeID uint, commentID uint) (*dtos.InternalCommitteeReviewCommentDetailsDTO, error)
	Create(committeeID uint, currentUserID uint, dto *dtos.InternalCommitteeReviewCommentFormDTO) (uint, error)
	Update(committeeID uint, id uint, currentUserID uint, dto *dtos.InternalCommitteeReviewCommentFormDTO) error
	Delete(committeeID uint, ids []uint, currentUserID uint) ([]uint, error)
	Resolve(committeeID uint, commentID uint, currentUserID uint) error
	GetDocumentPath(committeeID uint, commentID uint, token string) (string, error)
}

type InternalCommitteeReviewCommentService struct {
	repo             repositories.InternalCommitteeReviewCommentRepositoryInterface
	userSvc          UserServiceInterface
	activityLogSvc   logServices.ActivityLogServiceInterface
	committeeSvc     InternalCommitteeServiceInterface
	commentStatusSvc masterService.CommentStatusServiceInterface
}

func NewInternalCommitteeReviewCommentService(
	repo repositories.InternalCommitteeReviewCommentRepositoryInterface,
	userSvc UserServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	committeeSvc InternalCommitteeServiceInterface,
	commentStatusSvc masterService.CommentStatusServiceInterface,
) InternalCommitteeReviewCommentServiceInterface {
	return &InternalCommitteeReviewCommentService{
		repo:             repo,
		userSvc:          userSvc,
		activityLogSvc:   activityLogSvc,
		committeeSvc:     committeeSvc,
		commentStatusSvc: commentStatusSvc,
	}
}

func (service *InternalCommitteeReviewCommentService) GetAll(committeeID uint, currentUserID uint, params *dtos.InternalCommitteeReviewCommentListQueryParamsDTO) (*dtos.InternalCommitteeReviewCommentListResponseDTO, error) {
	// Get latest version ID for the committee
	committee, err := service.committeeSvc.GetLatestVersionByID(committeeID)
	if err != nil {
		return nil, err
	}

	var latestVersionID uint
	for _, version := range committee.InternalCommitteeVersions {
		if version.IsLatestVersion != nil && *version.IsLatestVersion {
			latestVersionID = version.ID
			break
		}
	}

	if latestVersionID == 0 {
		return nil, errors.New("no latest version found for committee")
	}

	comments, err := service.repo.GetAll(latestVersionID, currentUserID, params)
	if err != nil {
		return nil, err
	}

	// Build hierarchy
	dtoMap := make(map[uint]*dtos.InternalCommitteeReviewCommentListDTO)
	var roots []*dtos.InternalCommitteeReviewCommentListDTO

	for _, comment := range comments.Comments {
		dto := &dtos.InternalCommitteeReviewCommentListDTO{
			ID:                          comment.ID,
			Comment:                     comment.Comment,
			Level:                       comment.Level,
			IsEdited:                    comment.IsEdited,
			CommentStatusID:             comment.CommentStatusID,
			CommentStatusTitleEn:        comment.CommentStatusTitleEn,
			CommentStatusTitleAr:        comment.CommentStatusTitleAr,
			CommentStatusColorCode:      comment.CommentStatusColorCode,
			InternalCommitteeVersionID:  comment.InternalCommitteeVersionID,
			CreatedAt:                   comment.CreatedAt,
			CreatedUserID:               comment.CreatedUserID,
			CreatedUserFirstNameEn:      comment.CreatedUserFirstNameEn,
			CreatedUserFirstNameAr:      comment.CreatedUserFirstNameAr,
			CreatedUserLastNameEn:       comment.CreatedUserLastNameEn,
			CreatedUserLastNameAr:       comment.CreatedUserLastNameAr,
			CreatedUserInitial:          comment.CreatedUserInitial,
			CreatedUserColor:            comment.CreatedUserColor,
			CreatedUserEmail:            comment.CreatedUserEmail,
			CreatedUserProfileImagePath: comment.CreatedUserProfileImagePath,
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

	return &dtos.InternalCommitteeReviewCommentListResponseDTO{
		Comments:      roots,
		OpenCount:     comments.OpenCount,
		ResolvedCount: comments.ResolvedCount,
	}, nil
}

func (service *InternalCommitteeReviewCommentService) GetByID(committeeID uint, commentID uint) (*models.InternalCommitteeReviewComment, error) {
	// Verify committee exists
	_, err := service.committeeSvc.GetByID(committeeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, committeeID)
		}
		return nil, err
	}

	comment, err := service.repo.GetByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, commentID)
		}
		return nil, err
	}

	return comment, nil
}

func (service *InternalCommitteeReviewCommentService) GetFullDetailsByID(committeeID uint, commentID uint) (*dtos.InternalCommitteeReviewCommentDetailsDTO, error) {
	if _, err := service.committeeSvc.GetByID(committeeID); err != nil {
		fmt.Printf("Committee not found for committeeID %d", committeeID)
		return nil, err
	}

	var commentDetail dtos.InternalCommitteeReviewCommentDetailsDTO

	comment, err := service.repo.GetFullDetailsByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, commentID)
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
		var documentsDTO []*dtos.InternalCommitteeReviewCommentDocumentDTO
		for _, doc := range comment.Documents {
			documentDTO := &dtos.InternalCommitteeReviewCommentDocumentDTO{
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

func (service *InternalCommitteeReviewCommentService) Create(committeeID uint, currentUserID uint, formDTO *dtos.InternalCommitteeReviewCommentFormDTO) (uint, error) {
	// Get latest version ID for the committee
	committee, err := service.committeeSvc.GetLatestVersionByID(committeeID)
	if err != nil {
		return 0, err
	}

	var latestVersionID uint
	var nextReviewLevel *int
	for _, version := range committee.InternalCommitteeVersions {
		if version.IsLatestVersion != nil && *version.IsLatestVersion {
			latestVersionID = version.ID
			nextReviewLevel = version.NextReviewLevel
			break
		}
	}

	if latestVersionID == 0 {
		return 0, errors.New("no latest version found for committee")
	}

	// Get open status
	openStatus, err := service.commentStatusSvc.GetByType("open")
	if err != nil {
		return 0, err
	}

	// Validate parent comment if provided
	if formDTO.ParentCommentID != nil {
		parentComment, err := service.repo.GetByID(*formDTO.ParentCommentID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return 0, utility.NewResourceNotFoundError(err, *formDTO.ParentCommentID)
			}
			return 0, err
		}
		if parentComment.InternalCommitteeVersionID != latestVersionID {
			return 0, errors.New("parent comment belongs to different committee version")
		}
	}

	tx := service.repo.BeginTransaction()

	comment := &models.InternalCommitteeReviewComment{
		InternalCommitteeID:        committeeID,
		Comment:                    formDTO.Comment,
		Level:                      nextReviewLevel,
		CommentStatusID:            openStatus.ID,
		InternalCommitteeVersionID: latestVersionID,
		ParentCommentID:            formDTO.ParentCommentID,
		CreatedBYID:                currentUserID,
		CreatedAt:                  time.Now(),
	}

	id, err := service.repo.Create(comment, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.InternalCommitteeReviewCommentDocument
	for _, doc := range formDTO.Documents {
		// Encrypt the filename
		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		// Construct the URL for the document
		url := utility.GenerateDocumentURL("organization", "internal-committee", strconv.FormatUint(uint64(committeeID), 10), "review-comments", strconv.FormatUint(uint64(id), 10), "documents")

		document := &models.InternalCommitteeReviewCommentDocument{
			InternalCommitteeReviewCommentID: id,
			URL:                              &url,
			Name:                             doc.FileName,
			Size:                             &doc.Size,
			Token:                            &encryptedFilename,
			OriginalName:                     &doc.OriginalName,
			Type:                             &doc.Type,
			CreatedByID:                      currentUserID,
			CreatedAt:                        time.Now(),
		}

		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := service.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		for _, doc := range documents {
			storagePath := filepath.Join("organization", "internal-committee", strconv.FormatUint(uint64(committeeID), 10), "review-comments", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}

	tx.Commit()

	go func() {
		existingComment, err := service.GetFullDetailsByID(committeeID, id)
		if err != nil {
			utility.LogError("Failed to fetch comment", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = service.activityLogSvc.LogActivity(
			models.CreateAction,
			models.InternalCommitteeReviewComment{}.TableName(),
			id,
			existingComment,
			currentUserID,
			"organization", "committee_management", "internal_committee_review_comment", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.InternalCommitteeReviewComment{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()
	return id, nil
}

func (service *InternalCommitteeReviewCommentService) Update(committeeID uint, commentID uint, currentUserID uint, formDTO *dtos.InternalCommitteeReviewCommentFormDTO) error {
	if _, err := service.committeeSvc.GetByID(committeeID); err != nil {
		return err
	}

	existingComment, err := service.GetByID(committeeID, commentID)
	if err != nil {
		return err
	}

	oldComment, err := service.GetFullDetailsByID(committeeID, commentID)
	if err != nil {
		return err
	}

	// Check if user can update this comment
	if existingComment.CreatedBYID != currentUserID {
		return utility.NewForbiddenError("you can only update your own comments")
	}

	// Check if comment is resolved (resolved comments cannot be edited)
	commentStatus, err := service.commentStatusSvc.GetByID(existingComment.CommentStatusID)
	if err != nil {
		return err
	}
	if commentStatus.Type == "resolved" {
		return utility.NewForbiddenError("resolved comments cannot be edited")
	}

	tx := service.repo.BeginTransaction()
	now := time.Now()

	existingComment.Comment = formDTO.Comment
	existingComment.UpdatedBYID = &currentUserID
	existingComment.UpdatedAt = &now
	existingComment.IsEdited = true

	err = service.repo.Update(existingComment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Handle documents
	var newDocuments []*models.InternalCommitteeReviewCommentDocument
	var deleteDocumentIDs []*uint
	var documentsToDelete []*models.InternalCommitteeReviewCommentDocument

	for _, doc := range formDTO.Documents {
		// NEW FILE
		if doc.IsNewFile != nil && *doc.IsNewFile {
			encryptedFilename, err := utility.EncryptData([]byte(doc.FileName))
			if err != nil {
				tx.Rollback()
				return err
			}

			url := utility.GenerateDocumentURL(
				"organization",
				"internal-committee",
				strconv.FormatUint(uint64(committeeID), 10),
				"review-comments",
				strconv.FormatUint(uint64(commentID), 10),
				"documents",
			)

			document := &models.InternalCommitteeReviewCommentDocument{
				InternalCommitteeReviewCommentID: commentID,
				URL:                              &url,
				Name:                             doc.FileName,
				Size:                             &doc.Size,
				Token:                            &encryptedFilename,
				OriginalName:                     &doc.OriginalName,
				Type:                             &doc.Type,
				CreatedByID:                      currentUserID,
				CreatedAt:                        now,
			}

			newDocuments = append(newDocuments, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile && doc.ID != nil {
			// DELETE FILE
			existingDocument, err := service.repo.GetDocumentByID(*doc.ID)
			if err != nil {
				tx.Rollback()
				return err
			}

			deleteDocumentIDs = append(deleteDocumentIDs, &existingDocument.ID)
			documentsToDelete = append(documentsToDelete, existingDocument)
		}
	}

	if len(newDocuments) > 0 {
		if err := service.repo.CreateDocuments(newDocuments, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIDs) > 0 {
		if err := service.repo.DeleteDocuments(deleteDocumentIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// File system operations
	for _, doc := range newDocuments {
		storagePath := filepath.Join("organization", "internal-committee", strconv.FormatUint(uint64(committeeID), 10), "review-comments", strconv.FormatUint(uint64(commentID), 10), "documents")
		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range documentsToDelete {
		storagePath := filepath.Join("organization", "internal-committee", strconv.FormatUint(uint64(committeeID), 10), "review-comments", strconv.FormatUint(uint64(commentID), 10), "documents")
		if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	go func() {
		existingComment, err := service.GetFullDetailsByID(committeeID, commentID)
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
			err := service.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.InternalCommitteeReviewComment{}.TableName(),
				commentID,
				changedData,
				currentUserID,
				"organization", "committee_management", "internal_committee_review_comment", nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.InternalCommitteeReviewComment{}.TableName(),
					"recordID":  commentID,
					"userID":    currentUserID,
					"error":     err.Error(),
				}

				utility.LogError("Failed to log activity", fields)
			}
		}
	}()
	return nil
}

func (service *InternalCommitteeReviewCommentService) Delete(committeeID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := service.repo.BeginTransaction()
	if _, err := service.committeeSvc.GetByID(committeeID); err != nil {
		return nil, err
	}

	var failedIDs []uint
	var successfulIDs []uint
	var deletedComments []map[string]interface{}
	var documentsToDelete []*models.InternalCommitteeReviewCommentDocument

	for _, id := range ids {
		comment, err := service.repo.GetFullDetailsByID(id)
		if err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		deletedComments = append(deletedComments, map[string]interface{}{
			"id":        id,
			"committee": committeeID,
		})

		// Collect documents for deletion
		if len(comment.Documents) > 0 {
			documentsToDelete = append(documentsToDelete, comment.Documents...)
		}

		err = service.repo.Delete(id, tx)
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
		storagePath := filepath.Join("organization", "internal-committee", strconv.FormatUint(uint64(committeeID), 10), "review-comments", strconv.FormatUint(uint64(doc.InternalCommitteeReviewCommentID), 10), "documents")
		if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return failedIDs, err
		}
	}

	tx.Commit()
	if len(successfulIDs) > 0 {
		go func() {
			err := service.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.InternalCommitteeReviewComment{}.TableName(),
				successfulIDs,
				deletedComments,
				currentUserID,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to log bulk delete activity", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}()
	}
	return failedIDs, nil
}

func (service *InternalCommitteeReviewCommentService) Resolve(committeeID uint, commentID uint, currentUserID uint) error {
	tx := service.repo.BeginTransaction()

	if _, err := service.committeeSvc.GetByID(committeeID); err != nil {
		return err
	}

	existingComment, err := service.GetByID(committeeID, commentID)
	if err != nil {
		return err
	}

	commentStatus, err := service.commentStatusSvc.GetByType("resolved")
	if err != nil {
		return err
	}

	existingComment.CommentStatusID = commentStatus.ID
	existingComment.UpdatedBYID = &currentUserID
	err = service.repo.Update(existingComment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// Helper functions

func (service *InternalCommitteeReviewCommentService) buildCommentHierarchy(comments []*dtos.InternalCommitteeReviewCommentListDTO) []*dtos.InternalCommitteeReviewCommentListDTO {
	var rootComments []*dtos.InternalCommitteeReviewCommentListDTO
	commentMap := make(map[uint]*dtos.InternalCommitteeReviewCommentListDTO)

	// Create map of all comments
	for _, comment := range comments {
		commentMap[comment.ID] = comment
	}

	// Build hierarchy
	for _, comment := range comments {
		if comment.ParentCommentID == nil {
			rootComments = append(rootComments, comment)
		} else {
			if parent, exists := commentMap[*comment.ParentCommentID]; exists {
				parent.Children = append(parent.Children, comment)
			}
		}
	}

	return rootComments
}

func (service *InternalCommitteeReviewCommentService) GetDocumentPath(committeeID uint, commentID uint, token string) (string, error) {

	comment, err := service.repo.GetFullDetailsByID(commentID)
	if err != nil {
		return "", err
	}
	if len(comment.Documents) == 0 {
		return "", fmt.Errorf("no documents found for the audit final report review comment")
	}
	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}
	storagePath := filepath.Join("organization", "internal-committee", strconv.FormatUint(uint64(committeeID), 10), "review-comments", strconv.FormatUint(uint64(commentID), 10), "documents")
	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}
