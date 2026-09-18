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

type ProjectReviewCommentServiceInterface interface {
	GetAll(projectID uint, currentUserID uint, params *projectDTOs.ProjectReviewCommentListQueryParamsDTO) (*projectDTOs.ProjectReviewCommentListResponseDTO, error)
	GetByID(projectID uint, commentID uint) (*models.ProjectReviewComment, error)
	GetFullDetailsByID(projectID uint, commentID uint) (*projectDTOs.ProjectReviewCommentDetailDTO, error)
	Create(projectID uint, currentUserID uint, dto *projectDTOs.ProjectReviewCommentFormDTO) (uint, error)
	Update(projectID uint, id uint, currentUserID uint, dto *projectDTOs.ProjectReviewCommentFormDTO) error
	Delete(projectID uint, ids []uint, currentUserID uint) ([]uint, error)
	Resolve(projectID uint, commentID uint, currentUserID uint) error
	GetDocumentPath(projectID uint, commentID uint, token string) (string, error)
}

type ProjectReviewCommentService struct {
	repo          repositories.ProjectReviewCommentRepositoryInterface
	userSvc       organizationServices.UserServiceInterface
	activityLog   logServices.ActivityLogServiceInterface
	projectSvc    ProjectServiceInterface
	commentStatus masterService.CommentStatusServiceInterface
}

func NewProjectReviewCommentService(
	repo repositories.ProjectReviewCommentRepositoryInterface,
	userSvc organizationServices.UserServiceInterface,
	activityLog logServices.ActivityLogServiceInterface,
	projectSvc ProjectServiceInterface,
	commentStatus masterService.CommentStatusServiceInterface,
) ProjectReviewCommentServiceInterface {
	return &ProjectReviewCommentService{
		repo:          repo,
		userSvc:       userSvc,
		activityLog:   activityLog,
		projectSvc:    projectSvc,
		commentStatus: commentStatus,
	}
}

func (prcs *ProjectReviewCommentService) GetAll(projectID uint, currentUserID uint, params *projectDTOs.ProjectReviewCommentListQueryParamsDTO) (*projectDTOs.ProjectReviewCommentListResponseDTO, error) {
	_, err := prcs.projectSvc.GetByID(projectID)
	if err != nil {
		return nil, err
	}

	comments, err := prcs.repo.GetAll(projectID, currentUserID, params)
	if err != nil {
		return nil, err
	}

	dtoMap := make(map[uint]*projectDTOs.ProjectReviewCommentListDTO)
	var roots []*projectDTOs.ProjectReviewCommentListDTO

	for _, comment := range comments.Comments {
		dto := &projectDTOs.ProjectReviewCommentListDTO{
			ID:                          comment.ID,
			Comment:                     comment.Comment,
			IsEdited:                    comment.IsEdited,
			CommentStatusID:             comment.CommentStatusID,
			CommentStatusTitleEn:        comment.CommentStatusTitleEn,
			CommentStatusTitleAr:        comment.CommentStatusTitleAr,
			CommentStatusType:           comment.CommentStatusType,
			CommentStatusColorCode:      comment.CommentStatusColorCode,
			ProjectID:                   comment.ProjectID,
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
			Children:                    []*projectDTOs.ProjectReviewCommentListDTO{},
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

	response := &projectDTOs.ProjectReviewCommentListResponseDTO{
		Comments:      roots,
		OpenCount:     comments.OpenCount,
		ResolvedCount: comments.ResolvedCount,
	}

	return response, nil
}

func (prcs *ProjectReviewCommentService) GetByID(projectID uint, commentID uint) (*models.ProjectReviewComment, error) {
	if _, err := prcs.projectSvc.GetByID(projectID); err != nil {
		return nil, err
	}

	comment, err := prcs.repo.GetByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, projectID)
		}
		return nil, err
	}
	return comment, nil
}

func (prcs *ProjectReviewCommentService) GetFullDetailsByID(projectID uint, commentID uint) (*projectDTOs.ProjectReviewCommentDetailDTO, error) {
	if _, err := prcs.projectSvc.GetByID(projectID); err != nil {
		return nil, err
	}

	var commentDetail projectDTOs.ProjectReviewCommentDetailDTO

	comment, err := prcs.repo.GetFullDetailsByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, projectID)
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
		var documentsDTO []*projectDTOs.ProjectReviewCommentDocumentDTO
		for _, doc := range comment.Documents {
			url := utility.GenerateDocumentURL("projects", strconv.FormatUint(uint64(projectID), 10), "review-comments", strconv.FormatUint(uint64(commentID), 10), "documents")
			documentDTO := &projectDTOs.ProjectReviewCommentDocumentDTO{
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

func (prcs *ProjectReviewCommentService) Create(projectID uint, currentUserID uint, commentFormDTO *projectDTOs.ProjectReviewCommentFormDTO) (uint, error) {
	tx := prcs.repo.BeginTransaction()

	project, err := prcs.projectSvc.GetByID(projectID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	commentStatus, err := prcs.commentStatus.GetByType("open")
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	comment := &models.ProjectReviewComment{
		ProjectID:       projectID,
		Comment:         commentFormDTO.Comment,
		CommentStatusID: commentStatus.ID,
		Level:           project.NextReviewLevel,
		ParentCommentID: commentFormDTO.ParentCommentID,
		CreatedBYID:     currentUserID,
	}

	id, err := prcs.repo.Create(comment, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.ProjectReviewCommentDocument
	for _, doc := range commentFormDTO.Documents {
		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		document := &models.ProjectReviewCommentDocument{
			ProjectReviewCommentID: id,
			Name:                   doc.FileName,
			Size:                   doc.Size,
			Token:                  encryptedFilename,
			OriginalName:           doc.OriginalName,
			Type:                   doc.Type,
			CreatedBYID:            currentUserID,
			CreatedAt:              time.Now(),
		}

		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := prcs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		for _, doc := range documents {
			storagePath := filepath.Join("project", strconv.FormatUint(uint64(projectID), 10), "review-comments", strconv.FormatUint(uint64(id), 10), "documents")
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
		_ = prcs.activityLog.LogActivity(
			models.CreateAction,
			models.ProjectReviewComment{}.TableName(),
			id,
			comment,
			currentUserID,
			"project", "project_register", "project_review_comment", nil,
		)
	}()

	return id, nil
}

func (prcs *ProjectReviewCommentService) Update(projectID uint, id uint, currentUserID uint, commentFormDTO *projectDTOs.ProjectReviewCommentFormDTO) error {
	tx := prcs.repo.BeginTransaction()

	if _, err := prcs.projectSvc.GetByID(projectID); err != nil {
		tx.Rollback()
		return err
	}

	existingComment, err := prcs.GetByID(projectID, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingComment.Comment = commentFormDTO.Comment
	existingComment.IsEdited = true
	now := time.Now()
	existingComment.UpdatedBYID = &currentUserID
	existingComment.UpdatedAt = &now

	if err := prcs.repo.Update(existingComment, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	go func() {
		_ = prcs.activityLog.LogActivity(
			models.UpdateAction,
			models.ProjectReviewComment{}.TableName(),
			id,
			existingComment,
			currentUserID,
			"project", "project_register", "project_review_comment", nil,
		)
	}()

	return nil
}

func (prcs *ProjectReviewCommentService) Delete(projectID uint, ids []uint, currentUserID uint) ([]uint, error) {
	if _, err := prcs.projectSvc.GetByID(projectID); err != nil {
		return nil, err
	}

	tx := prcs.repo.BeginTransaction()

	var failedIDs []uint
	var successfulIDs []uint
	documentsToDelete := make(map[uint][]*projectDTOs.ProjectReviewCommentDocumentDTO)
	deletedComments := make(map[uint]*projectDTOs.ProjectReviewCommentDetailDTO)

	for _, id := range ids {
		comment, err := prcs.GetFullDetailsByID(projectID, id)
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
				if err := prcs.repo.DeleteDocuments(deleteDocIds, tx); err != nil {
					failedIDs = append(failedIDs, id)
					continue
				}
			}
		}

		if err := prcs.repo.Delete(id, tx); err != nil {
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
				storagePath := filepath.Join("project", strconv.FormatUint(uint64(projectID), 10), "review-comments", strconv.FormatUint(uint64(commentID), 10), "documents")
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
			_ = prcs.activityLog.LogBulkActivities(
				models.DeleteAction,
				models.ProjectReviewComment{}.TableName(),
				successfulIDs,
				deletedComments,
				currentUserID,
				nil,
			)
		}()
	}

	return failedIDs, nil
}

func (prcs *ProjectReviewCommentService) Resolve(projectID uint, commentID uint, currentUserID uint) error {
	tx := prcs.repo.BeginTransaction()

	if _, err := prcs.projectSvc.GetByID(projectID); err != nil {
		tx.Rollback()
		return err
	}

	existingComment, err := prcs.GetByID(projectID, commentID)
	if err != nil {
		tx.Rollback()
		return err
	}

	resolvedStatus, err := prcs.commentStatus.GetByType("resolved")
	if err != nil {
		tx.Rollback()
		return err
	}

	now := time.Now()
	existingComment.CommentStatusID = resolvedStatus.ID
	existingComment.UpdatedBYID = &currentUserID
	existingComment.UpdatedAt = &now

	if err := prcs.repo.Update(existingComment, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (prcs *ProjectReviewCommentService) GetDocumentPath(projectID uint, commentID uint, token string) (string, error) {
	comment, err := prcs.repo.GetFullDetailsByID(commentID)
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

	storagePath := filepath.Join("project", strconv.FormatUint(uint64(projectID), 10), "review-comments", strconv.FormatUint(uint64(commentID), 10), "documents")
	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}
