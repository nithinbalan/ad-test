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
	projectDTOs "jethur/grc/project/dtos"
	"jethur/grc/project/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ProjectChangeRequestReviewCommentServiceInterface interface {
	GetAll(projectID uint, changeRequestID uint, currentUserID uint, params *projectDTOs.ProjectChangeRequestReviewCommentListQueryParamsDTO) (*projectDTOs.ProjectChangeRequestReviewCommentListResponseDTO, error)
	GetByID(projectID uint, changeRequestID uint, commentID uint) (*models.ProjectChangeRequestReviewComment, error)
	GetFullDetailsByID(projectID uint, changeRequestID uint, commentID uint) (*projectDTOs.ProjectChangeRequestReviewCommentDetailDTO, error)
	Create(projectID uint, changeRequestID uint, currentUserID uint, dto *projectDTOs.ProjectChangeRequestReviewCommentFormDTO) (uint, error)
	Update(projectID uint, changeRequestID uint, id uint, currentUserID uint, dto *projectDTOs.ProjectChangeRequestReviewCommentFormDTO) error
	Delete(projectID uint, changeRequestID uint, ids []uint, currentUserID uint) ([]uint, error)
	Resolve(projectID uint, changeRequestID uint, commentID uint, currentUserID uint) error
	GetDocumentPath(projectID uint, changeRequestID uint, commentID uint, token string) (string, error)
}

type ProjectChangeRequestReviewCommentService struct {
	repo          repositories.ProjectChangeRequestReviewCommentRepositoryInterface
	userSvc       organizationServices.UserServiceInterface
	activityLog   logServices.ActivityLogServiceInterface
	projectCRSvc  ProjectChangeRequestServiceInterface
	commentStatus masterService.CommentStatusServiceInterface
}

func NewProjectChangeRequestReviewCommentService(
	repo repositories.ProjectChangeRequestReviewCommentRepositoryInterface,
	userSvc organizationServices.UserServiceInterface,
	activityLog logServices.ActivityLogServiceInterface,
	projectCRSvc ProjectChangeRequestServiceInterface,
	commentStatus masterService.CommentStatusServiceInterface,
) ProjectChangeRequestReviewCommentServiceInterface {
	return &ProjectChangeRequestReviewCommentService{
		repo:          repo,
		userSvc:       userSvc,
		activityLog:   activityLog,
		projectCRSvc:  projectCRSvc,
		commentStatus: commentStatus,
	}
}

func (pcrrcs *ProjectChangeRequestReviewCommentService) GetAll(projectID uint, changeRequestID uint, currentUserID uint, params *projectDTOs.ProjectChangeRequestReviewCommentListQueryParamsDTO) (*projectDTOs.ProjectChangeRequestReviewCommentListResponseDTO, error) {
	_, err := pcrrcs.projectCRSvc.GetByID(projectID, changeRequestID)
	if err != nil {
		return nil, err
	}

	comments, err := pcrrcs.repo.GetAll(changeRequestID, currentUserID, params)
	if err != nil {
		return nil, err
	}

	dtoMap := make(map[uint]*projectDTOs.ProjectChangeRequestReviewCommentListDTO)
	var roots []*projectDTOs.ProjectChangeRequestReviewCommentListDTO

	for _, comment := range comments.Comments {
		dto := &projectDTOs.ProjectChangeRequestReviewCommentListDTO{
			ID:                          comment.ID,
			Comment:                     comment.Comment,
			IsEdited:                    comment.IsEdited,
			CommentStatusID:             comment.CommentStatusID,
			CommentStatusTitleEn:        comment.CommentStatusTitleEn,
			CommentStatusTitleAr:        comment.CommentStatusTitleAr,
			CommentStatusType:           comment.CommentStatusType,
			CommentStatusColorCode:      comment.CommentStatusColorCode,
			ProjectChangeRequestID:      comment.ProjectChangeRequestID,
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
			Children:                    []*projectDTOs.ProjectChangeRequestReviewCommentListDTO{},
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

	response := &projectDTOs.ProjectChangeRequestReviewCommentListResponseDTO{
		Comments:      roots,
		OpenCount:     comments.OpenCount,
		ResolvedCount: comments.ResolvedCount,
	}

	return response, nil
}

func (pcrrcs *ProjectChangeRequestReviewCommentService) GetByID(projectID uint, changeRequestID uint, commentID uint) (*models.ProjectChangeRequestReviewComment, error) {
	if _, err := pcrrcs.projectCRSvc.GetByID(projectID, changeRequestID); err != nil {
		return nil, err
	}

	comment, err := pcrrcs.repo.GetByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, changeRequestID)
		}
		return nil, err
	}
	return comment, nil
}

func (pcrrcs *ProjectChangeRequestReviewCommentService) GetFullDetailsByID(projectID uint, changeRequestID uint, commentID uint) (*projectDTOs.ProjectChangeRequestReviewCommentDetailDTO, error) {
	if _, err := pcrrcs.projectCRSvc.GetByID(projectID, changeRequestID); err != nil {
		return nil, err
	}

	var commentDetail projectDTOs.ProjectChangeRequestReviewCommentDetailDTO

	comment, err := pcrrcs.repo.GetFullDetailsByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, changeRequestID)
		}
		return nil, err
	}

	err = mapstructure.Decode(comment, &commentDetail)
	if err != nil {
		return nil, err
	}

	commentDetail.CreatedAt = comment.CreatedAt
	commentDetail.UpdatedAt = comment.UpdatedAt

	if len(comment.Documents) > 0 {
		var documentsDTO []*projectDTOs.ProjectChangeRequestReviewCommentDocumentDTO
		for _, doc := range comment.Documents {
			url := utility.GenerateDocumentURL("projects", strconv.FormatUint(uint64(projectID), 10), "change-requests", strconv.FormatUint(uint64(changeRequestID), 10), "review-comments", strconv.FormatUint(uint64(commentID), 10), "documents")
			documentDTO := &projectDTOs.ProjectChangeRequestReviewCommentDocumentDTO{
				ID:           doc.ID,
				URL:          url,
				Name:         doc.Name,
				Size:         doc.Size,
				Token:        doc.Token,
				OriginalName: doc.OriginalName,
				Type:         doc.Type,
			}
			documentsDTO = append(documentsDTO, documentDTO)
		}
		commentDetail.Documents = documentsDTO
	}

	return &commentDetail, nil
}

func (pcrrcs *ProjectChangeRequestReviewCommentService) Create(projectID uint, changeRequestID uint, currentUserID uint, commentFormDTO *projectDTOs.ProjectChangeRequestReviewCommentFormDTO) (uint, error) {
	tx := pcrrcs.repo.BeginTransaction()

	pcr, err := pcrrcs.projectCRSvc.GetByID(projectID, changeRequestID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	commentStatus, err := pcrrcs.commentStatus.GetByType("open")
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	comment := &models.ProjectChangeRequestReviewComment{
		ProjectChangeRequestID: changeRequestID,
		Comment:                commentFormDTO.Comment,
		CommentStatusID:        commentStatus.ID,
		Level:                  pcr.NextReviewLevel,
		ParentCommentID:        commentFormDTO.ParentCommentID,
		CreatedBYID:            currentUserID,
	}

	id, err := pcrrcs.repo.Create(comment, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.ProjectChangeRequestReviewCommentDocument
	for _, doc := range commentFormDTO.Documents {
		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		document := &models.ProjectChangeRequestReviewCommentDocument{
			ProjectChangeRequestReviewCommentID: id,
			Name:                                doc.FileName,
			Size:                                doc.Size,
			Token:                               encryptedFilename,
			OriginalName:                        doc.OriginalName,
			Type:                                doc.Type,
			CreatedBYID:                         currentUserID,
			CreatedAt:                           time.Now(),
		}

		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := pcrrcs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		for _, doc := range documents {
			storagePath := filepath.Join("project", "change-requests", strconv.FormatUint(uint64(changeRequestID), 10), "review-comments", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	go func() {
		_ = pcrrcs.activityLog.LogActivity(
			models.CreateAction,
			models.ProjectChangeRequestReviewComment{}.TableName(),
			id,
			comment,
			currentUserID,
			"project", "project_register", "project_change_request_review_comment", nil,
		)
	}()

	return id, nil
}

func (pcrrcs *ProjectChangeRequestReviewCommentService) Update(projectID uint, changeRequestID uint, id uint, currentUserID uint, commentFormDTO *projectDTOs.ProjectChangeRequestReviewCommentFormDTO) error {
	tx := pcrrcs.repo.BeginTransaction()

	if _, err := pcrrcs.projectCRSvc.GetByID(projectID, changeRequestID); err != nil {
		tx.Rollback()
		return err
	}

	existingComment, err := pcrrcs.GetByID(projectID, changeRequestID, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingComment.Comment = commentFormDTO.Comment
	existingComment.IsEdited = true
	now := time.Now()
	existingComment.UpdatedBYID = &currentUserID
	existingComment.UpdatedAt = &now

	if err := pcrrcs.repo.Update(existingComment, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	go func() {
		_ = pcrrcs.activityLog.LogActivity(
			models.UpdateAction,
			models.ProjectChangeRequestReviewComment{}.TableName(),
			id,
			existingComment,
			currentUserID,
			"project", "project_register", "project_change_request_review_comment", nil,
		)
	}()

	return nil
}

func (pcrrcs *ProjectChangeRequestReviewCommentService) Delete(projectID uint, changeRequestID uint, ids []uint, currentUserID uint) ([]uint, error) {
	if _, err := pcrrcs.projectCRSvc.GetByID(projectID, changeRequestID); err != nil {
		return nil, err
	}

	tx := pcrrcs.repo.BeginTransaction()

	var failedIDs []uint
	var successfulIDs []uint
	documentsToDelete := make(map[uint][]*projectDTOs.ProjectChangeRequestReviewCommentDocumentDTO)
	deletedComments := make(map[uint]*projectDTOs.ProjectChangeRequestReviewCommentDetailDTO)

	for _, id := range ids {
		comment, err := pcrrcs.GetFullDetailsByID(projectID, changeRequestID, id)
		if err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		if len(comment.Documents) > 0 {
			var deleteDocIds []*uint
			for _, doc := range comment.Documents {
				documentsToDelete[id] = append(documentsToDelete[id], doc)
				docID := doc.ID
				deleteDocIds = append(deleteDocIds, &docID)
			}
			if len(deleteDocIds) > 0 {
				if err := pcrrcs.repo.DeleteDocuments(deleteDocIds, tx); err != nil {
					failedIDs = append(failedIDs, id)
					continue
				}
			}
		}

		if err := pcrrcs.repo.Delete(id, tx); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		deletedComments[id] = comment
		successfulIDs = append(successfulIDs, id)
	}

	if len(failedIDs) > 0 {
		tx.Rollback()
		return failedIDs, nil
	}

	if len(successfulIDs) > 0 {
		for commentID, documents := range documentsToDelete {
			for _, doc := range documents {
				storagePath := filepath.Join("project", "change-requests", strconv.FormatUint(uint64(changeRequestID), 10), "review-comments", strconv.FormatUint(uint64(commentID), 10), "documents")
				if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
					tx.Rollback()
					return nil, err
				}
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	if len(successfulIDs) > 0 {
		go func() {
			_ = pcrrcs.activityLog.LogBulkActivities(
				models.DeleteAction,
				models.ProjectChangeRequestReviewComment{}.TableName(),
				successfulIDs,
				deletedComments,
				currentUserID,
				nil,
			)
		}()
	}

	return failedIDs, nil
}

func (pcrrcs *ProjectChangeRequestReviewCommentService) Resolve(projectID uint, changeRequestID uint, commentID uint, currentUserID uint) error {
	tx := pcrrcs.repo.BeginTransaction()

	if _, err := pcrrcs.projectCRSvc.GetByID(projectID, changeRequestID); err != nil {
		tx.Rollback()
		return err
	}

	existingComment, err := pcrrcs.GetByID(projectID, changeRequestID, commentID)
	if err != nil {
		tx.Rollback()
		return err
	}

	resolvedStatus, err := pcrrcs.commentStatus.GetByType("resolved")
	if err != nil {
		tx.Rollback()
		return err
	}

	now := time.Now()
	existingComment.CommentStatusID = resolvedStatus.ID
	existingComment.UpdatedBYID = &currentUserID
	existingComment.UpdatedAt = &now

	if err := pcrrcs.repo.Update(existingComment, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (pcrrcs *ProjectChangeRequestReviewCommentService) GetDocumentPath(projectID uint, changeRequestID uint, commentID uint, token string) (string, error) {
	comment, err := pcrrcs.repo.GetFullDetailsByID(commentID)
	if err != nil {
		return "", err
	}
	if len(comment.Documents) == 0 {
		return "", fmt.Errorf("no documents found for this review comment")
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	storagePath := filepath.Join("project", "change-requests", strconv.FormatUint(uint64(changeRequestID), 10), "review-comments", strconv.FormatUint(uint64(commentID), 10), "documents")
	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}
