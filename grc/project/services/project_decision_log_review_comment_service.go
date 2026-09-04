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

type ProjectDecisionLogReviewCommentServiceInterface interface {
	GetAll(projectID uint, decisionLogID uint, currentUserID uint, params *projectDTOs.ProjectDecisionLogReviewCommentListQueryParamsDTO) (*projectDTOs.ProjectDecisionLogReviewCommentListResponseDTO, error)
	GetByID(projectID uint, decisionLogID uint, commentID uint) (*models.ProjectDecisionLogReviewComment, error)
	GetFullDetailsByID(projectID uint, decisionLogID uint, commentID uint) (*projectDTOs.ProjectDecisionLogReviewCommentDetailDTO, error)
	Create(projectID uint, decisionLogID uint, currentUserID uint, dto *projectDTOs.ProjectDecisionLogReviewCommentFormDTO) (uint, error)
	Update(projectID uint, decisionLogID uint, id uint, currentUserID uint, dto *projectDTOs.ProjectDecisionLogReviewCommentFormDTO) error
	Delete(projectID uint, decisionLogID uint, ids []uint, currentUserID uint) ([]uint, error)
	Resolve(projectID uint, decisionLogID uint, commentID uint, currentUserID uint) error
	GetDocumentPath(projectID uint, decisionLogID uint, commentID uint, token string) (string, error)
}

type ProjectDecisionLogReviewCommentService struct {
	repo          repositories.ProjectDecisionLogReviewCommentRepositoryInterface
	userSvc       organizationServices.UserServiceInterface
	activityLog   logServices.ActivityLogServiceInterface
	projectDLSvc  ProjectDecisionLogServiceInterface
	commentStatus masterService.CommentStatusServiceInterface
}

func NewProjectDecisionLogReviewCommentService(
	repo repositories.ProjectDecisionLogReviewCommentRepositoryInterface,
	userSvc organizationServices.UserServiceInterface,
	activityLog logServices.ActivityLogServiceInterface,
	projectDLSvc ProjectDecisionLogServiceInterface,
	commentStatus masterService.CommentStatusServiceInterface,
) ProjectDecisionLogReviewCommentServiceInterface {
	return &ProjectDecisionLogReviewCommentService{
		repo:          repo,
		userSvc:       userSvc,
		activityLog:   activityLog,
		projectDLSvc:  projectDLSvc,
		commentStatus: commentStatus,
	}
}

func (pdlrcs *ProjectDecisionLogReviewCommentService) GetAll(projectID uint, decisionLogID uint, currentUserID uint, params *projectDTOs.ProjectDecisionLogReviewCommentListQueryParamsDTO) (*projectDTOs.ProjectDecisionLogReviewCommentListResponseDTO, error) {
	_, err := pdlrcs.projectDLSvc.GetByID(projectID, decisionLogID)
	if err != nil {
		return nil, err
	}

	comments, err := pdlrcs.repo.GetAll(decisionLogID, currentUserID, params)
	if err != nil {
		return nil, err
	}

	dtoMap := make(map[uint]*projectDTOs.ProjectDecisionLogReviewCommentListDTO)
	var roots []*projectDTOs.ProjectDecisionLogReviewCommentListDTO

	for _, comment := range comments.Comments {
		dto := &projectDTOs.ProjectDecisionLogReviewCommentListDTO{
			ID:                          comment.ID,
			Comment:                     comment.Comment,
			IsEdited:                    comment.IsEdited,
			CommentStatusID:             comment.CommentStatusID,
			CommentStatusTitleEn:        comment.CommentStatusTitleEn,
			CommentStatusTitleAr:        comment.CommentStatusTitleAr,
			CommentStatusType:           comment.CommentStatusType,
			CommentStatusColorCode:      comment.CommentStatusColorCode,
			ProjectDecisionLogID:        comment.ProjectDecisionLogID,
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
			Children:                    []*projectDTOs.ProjectDecisionLogReviewCommentListDTO{},
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

	response := &projectDTOs.ProjectDecisionLogReviewCommentListResponseDTO{
		Comments:      roots,
		OpenCount:     comments.OpenCount,
		ResolvedCount: comments.ResolvedCount,
	}

	return response, nil
}

func (pdlrcs *ProjectDecisionLogReviewCommentService) GetByID(projectID uint, decisionLogID uint, commentID uint) (*models.ProjectDecisionLogReviewComment, error) {
	if _, err := pdlrcs.projectDLSvc.GetByID(projectID, decisionLogID); err != nil {
		return nil, err
	}

	comment, err := pdlrcs.repo.GetByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, decisionLogID)
		}
		return nil, err
	}
	return comment, nil
}

func (pdlrcs *ProjectDecisionLogReviewCommentService) GetFullDetailsByID(projectID uint, decisionLogID uint, commentID uint) (*projectDTOs.ProjectDecisionLogReviewCommentDetailDTO, error) {
	if _, err := pdlrcs.projectDLSvc.GetByID(projectID, decisionLogID); err != nil {
		return nil, err
	}

	var commentDetail projectDTOs.ProjectDecisionLogReviewCommentDetailDTO

	comment, err := pdlrcs.repo.GetFullDetailsByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, decisionLogID)
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
		var documentsDTO []*projectDTOs.ProjectDecisionLogReviewCommentDocumentDTO
		for _, doc := range comment.Documents {
			url := utility.GenerateDocumentURL("projects", strconv.FormatUint(uint64(projectID), 10), "decision-logs", strconv.FormatUint(uint64(decisionLogID), 10), "review-comments", strconv.FormatUint(uint64(commentID), 10), "documents")
			documentDTO := &projectDTOs.ProjectDecisionLogReviewCommentDocumentDTO{
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

func (pdlrcs *ProjectDecisionLogReviewCommentService) Create(projectID uint, decisionLogID uint, currentUserID uint, commentFormDTO *projectDTOs.ProjectDecisionLogReviewCommentFormDTO) (uint, error) {
	tx := pdlrcs.repo.BeginTransaction()

	pdl, err := pdlrcs.projectDLSvc.GetByID(projectID, decisionLogID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	commentStatus, err := pdlrcs.commentStatus.GetByType("open")
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	comment := &models.ProjectDecisionLogReviewComment{
		ProjectDecisionLogID: decisionLogID,
		Comment:              commentFormDTO.Comment,
		CommentStatusID:      commentStatus.ID,
		Level:                pdl.NextReviewLevel,
		ParentCommentID:      commentFormDTO.ParentCommentID,
		CreatedBYID:          currentUserID,
	}

	id, err := pdlrcs.repo.Create(comment, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.ProjectDecisionLogReviewCommentDocument
	for _, doc := range commentFormDTO.Documents {
		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		document := &models.ProjectDecisionLogReviewCommentDocument{
			ProjectDecisionLogReviewCommentID: id,
			Name:                              doc.FileName,
			Size:                              doc.Size,
			Token:                             encryptedFilename,
			OriginalName:                      doc.OriginalName,
			Type:                              doc.Type,
			CreatedBYID:                       currentUserID,
			CreatedAt:                         time.Now(),
		}

		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := pdlrcs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		for _, doc := range documents {
			storagePath := filepath.Join("project", strconv.FormatUint(uint64(projectID), 10), "decision-logs", strconv.FormatUint(uint64(decisionLogID), 10), "review-comments", strconv.FormatUint(uint64(id), 10), "documents")
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
		_ = pdlrcs.activityLog.LogActivity(
			models.CreateAction,
			models.ProjectDecisionLogReviewComment{}.TableName(),
			id,
			comment,
			currentUserID,
			"project", "project_register", "project_decision_log_review_comment", nil,
		)
	}()

	return id, nil
}

func (pdlrcs *ProjectDecisionLogReviewCommentService) Update(projectID uint, decisionLogID uint, id uint, currentUserID uint, commentFormDTO *projectDTOs.ProjectDecisionLogReviewCommentFormDTO) error {
	tx := pdlrcs.repo.BeginTransaction()

	if _, err := pdlrcs.projectDLSvc.GetByID(projectID, decisionLogID); err != nil {
		tx.Rollback()
		return err
	}

	existingComment, err := pdlrcs.GetByID(projectID, decisionLogID, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingComment.Comment = commentFormDTO.Comment
	existingComment.IsEdited = true
	now := time.Now()
	existingComment.UpdatedBYID = &currentUserID
	existingComment.UpdatedAt = &now

	if err := pdlrcs.repo.Update(existingComment, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	go func() {
		_ = pdlrcs.activityLog.LogActivity(
			models.UpdateAction,
			models.ProjectDecisionLogReviewComment{}.TableName(),
			id,
			existingComment,
			currentUserID,
			"project", "project_register", "project_decision_log_review_comment", nil,
		)
	}()

	return nil
}

func (pdlrcs *ProjectDecisionLogReviewCommentService) Delete(projectID uint, decisionLogID uint, ids []uint, currentUserID uint) ([]uint, error) {
	if _, err := pdlrcs.projectDLSvc.GetByID(projectID, decisionLogID); err != nil {
		return nil, err
	}

	tx := pdlrcs.repo.BeginTransaction()

	var failedIDs []uint
	var successfulIDs []uint
	documentsToDelete := make(map[uint][]*projectDTOs.ProjectDecisionLogReviewCommentDocumentDTO)
	deletedComments := make(map[uint]*projectDTOs.ProjectDecisionLogReviewCommentDetailDTO)

	for _, id := range ids {
		comment, err := pdlrcs.GetFullDetailsByID(projectID, decisionLogID, id)
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
				if err := pdlrcs.repo.DeleteDocuments(deleteDocIds, tx); err != nil {
					failedIDs = append(failedIDs, id)
					continue
				}
			}
		}

		if err := pdlrcs.repo.Delete(id, tx); err != nil {
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
				storagePath := filepath.Join("project", strconv.FormatUint(uint64(projectID), 10), "decision-logs", strconv.FormatUint(uint64(decisionLogID), 10), "review-comments", strconv.FormatUint(uint64(commentID), 10), "documents")
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
			_ = pdlrcs.activityLog.LogBulkActivities(
				models.DeleteAction,
				models.ProjectDecisionLogReviewComment{}.TableName(),
				successfulIDs,
				deletedComments,
				currentUserID,
				nil,
			)
		}()
	}

	return failedIDs, nil
}

func (pdlrcs *ProjectDecisionLogReviewCommentService) Resolve(projectID uint, decisionLogID uint, commentID uint, currentUserID uint) error {
	tx := pdlrcs.repo.BeginTransaction()

	if _, err := pdlrcs.projectDLSvc.GetByID(projectID, decisionLogID); err != nil {
		tx.Rollback()
		return err
	}

	existingComment, err := pdlrcs.GetByID(projectID, decisionLogID, commentID)
	if err != nil {
		tx.Rollback()
		return err
	}

	resolvedStatus, err := pdlrcs.commentStatus.GetByType("resolved")
	if err != nil {
		tx.Rollback()
		return err
	}

	now := time.Now()
	existingComment.CommentStatusID = resolvedStatus.ID
	existingComment.UpdatedBYID = &currentUserID
	existingComment.UpdatedAt = &now

	if err := pdlrcs.repo.Update(existingComment, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (pdlrcs *ProjectDecisionLogReviewCommentService) GetDocumentPath(projectID uint, decisionLogID uint, commentID uint, token string) (string, error) {
	comment, err := pdlrcs.repo.GetFullDetailsByID(commentID)
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

	storagePath := filepath.Join("project", strconv.FormatUint(uint64(projectID), 10), "decision-logs", strconv.FormatUint(uint64(decisionLogID), 10), "review-comments", strconv.FormatUint(uint64(commentID), 10), "documents")
	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}
