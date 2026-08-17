package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"jethur/config"
	incidentDTOs "jethur/grc/incident/dtos"
	incidentRepositories "jethur/grc/incident/repositories"
	logServices "jethur/grc/logs/services"
	masterService "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type IncidentClosureReviewCommentServiceInterface interface {
	GetAll(incidentClosureID uint, currentUserID uint, params *incidentDTOs.IncidentClosureReviewCommentListQueryParamsDTO) (*incidentDTOs.IncidentClosureReviewCommentListResponseDTO, error)
	GetByID(incidentClosureID uint, commentID uint) (*models.IncidentClosureReviewComment, error)
	GetFullDetailsByID(incidentClosureID uint, commentID uint) (*incidentDTOs.IncidentClosureReviewCommentDetailDTO, error)
	Create(incidentClosureID uint, currentUserID uint, dto *incidentDTOs.IncidentClosureReviewCommentFormDTO) (uint, error)
	Update(incidentClosureID uint, id uint, currentUserID uint, dto *incidentDTOs.IncidentClosureReviewCommentFormDTO) error
	Delete(incidentClosureID uint, ids []uint, currentUserID uint) ([]uint, error)
	Resolve(incidentClosureID uint, commentID uint, currentUserID uint) error
	GetDocumentPath(incidentClosureID uint, commentID uint, token string) (string, error)
}

type IncidentClosureReviewCommentService struct {
	repo               incidentRepositories.IncidentClosureReviewCommentRepositoryInterface
	userSvc            organizationServices.UserServiceInterface
	activityLogSvc     logServices.ActivityLogServiceInterface
	incidentClosureSvc IncidentClosureServiceInterface
	commentStatusSvc   masterService.CommentStatusServiceInterface
}

func NewIncidentClosureReviewCommentService(
	repo incidentRepositories.IncidentClosureReviewCommentRepositoryInterface,
	userSvc organizationServices.UserServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	incidentClosureSvc IncidentClosureServiceInterface,
	commentStatusSvc masterService.CommentStatusServiceInterface,
) IncidentClosureReviewCommentServiceInterface {
	return &IncidentClosureReviewCommentService{
		repo:               repo,
		userSvc:            userSvc,
		activityLogSvc:     activityLogSvc,
		incidentClosureSvc: incidentClosureSvc,
		commentStatusSvc:   commentStatusSvc,
	}
}

func (icrcs *IncidentClosureReviewCommentService) GetAll(incidentClosureID uint, currentUserID uint, params *incidentDTOs.IncidentClosureReviewCommentListQueryParamsDTO) (*incidentDTOs.IncidentClosureReviewCommentListResponseDTO, error) {
	incidentClosure, err := icrcs.incidentClosureSvc.GetByClosureID(incidentClosureID)
	if err != nil {
		return nil, err
	}

	comments, err := icrcs.repo.GetAll(incidentClosure.ID, currentUserID, params)
	if err != nil {
		return nil, err
	}
	dtoMap := make(map[uint]*incidentDTOs.IncidentClosureReviewCommentListDTO)
	var roots []*incidentDTOs.IncidentClosureReviewCommentListDTO

	for _, comment := range comments.Comments {
		dto := &incidentDTOs.IncidentClosureReviewCommentListDTO{
			ID:                          comment.ID,
			Comment:                     comment.Comment,
			IsEdited:                    comment.IsEdited,
			CommentStatusID:             comment.CommentStatusID,
			CommentStatusTitleEn:        comment.CommentStatusTitleEn,
			CommentStatusTitleAr:        comment.CommentStatusTitleAr,
			CommentStatusType:           comment.CommentStatusType,
			CommentStatusColorCode:      comment.CommentStatusColorCode,
			IncidentClosureID:           comment.IncidentClosureID,
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

	response := &incidentDTOs.IncidentClosureReviewCommentListResponseDTO{
		Comments:      roots,
		OpenCount:     comments.OpenCount,
		ResolvedCount: comments.ResolvedCount,
	}

	return response, nil
}

func (icrcs *IncidentClosureReviewCommentService) GetByID(incidentClosureID uint, commentID uint) (*models.IncidentClosureReviewComment, error) {
	incidentClosure, err := icrcs.incidentClosureSvc.GetByClosureID(incidentClosureID)
	if err != nil {
		return nil, err
	}

	comment, err := icrcs.repo.GetByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, incidentClosureID)
		}
		return nil, err
	}

	if comment.IncidentClosureID != incidentClosure.ID {
		return nil, utility.NewResourceNotFoundError(errors.New("resource not found"), incidentClosureID)
	}

	return comment, nil
}

func (icrcs *IncidentClosureReviewCommentService) GetFullDetailsByID(incidentClosureID uint, commentID uint) (*incidentDTOs.IncidentClosureReviewCommentDetailDTO, error) {
	incidentClosure, err := icrcs.incidentClosureSvc.GetByClosureID(incidentClosureID)
	if err != nil {
		return nil, err
	}

	var commentDetail incidentDTOs.IncidentClosureReviewCommentDetailDTO

	comment, err := icrcs.repo.GetFullDetailsByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, incidentClosureID)
		}
		return nil, err
	}

	if comment.IncidentClosureID != incidentClosure.ID {
		return nil, utility.NewResourceNotFoundError(errors.New("resource not found"), incidentClosureID)
	}

	err = mapstructure.Decode(comment, &commentDetail)
	if err != nil {
		return nil, err
	}

	commentDetail.CreatedAt = comment.CreatedAt

	if len(comment.Documents) > 0 {
		var documentsDTO []*incidentDTOs.IncidentClosureReviewCommentDocumentDTO
		for _, doc := range comment.Documents {
			documentDTO := &incidentDTOs.IncidentClosureReviewCommentDocumentDTO{
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

func (icrcs *IncidentClosureReviewCommentService) Create(incidentClosureID uint, currentUserID uint, commentFormDTO *incidentDTOs.IncidentClosureReviewCommentFormDTO) (uint, error) {
	tx := icrcs.repo.BeginTransaction()

	incidentClosure, err := icrcs.incidentClosureSvc.GetByClosureID(incidentClosureID)
	if err != nil {
		return 0, err
	}

	commentStatus, err := icrcs.commentStatusSvc.GetByType("open")
	if err != nil {
		return 0, err
	}

	comment := &models.IncidentClosureReviewComment{
		IncidentClosureID: incidentClosure.ID,
		Comment:           commentFormDTO.Comment,
		CommentStatusID:   commentStatus.ID,
		Level:             incidentClosure.NextReviewLevel,
		ParentCommentID:   commentFormDTO.ParentCommentID,
		CreatedBYID:       currentUserID,
	}

	id, err := icrcs.repo.Create(comment, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.IncidentClosureReviewCommentDocument
	for _, doc := range commentFormDTO.Documents {
		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL("incident", "incident-closure", strconv.FormatUint(uint64(incidentClosure.ID), 10), "review-comments", strconv.FormatUint(uint64(id), 10), "documents")

		document := &models.IncidentClosureReviewCommentDocument{
			IncidentClosureReviewCommentID: id,
			URL:                            &url,
			Name:                           doc.FileName,
			Size:                           &doc.Size,
			Token:                          &encryptedFilename,
			OriginalName:                   &doc.OriginalName,
			Type:                           &doc.Type,
			CreatedByID:                    currentUserID,
			CreatedAt:                      time.Now(),
		}

		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := icrcs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		for _, doc := range documents {
			storagePath := filepath.Join("incident", "incident-closure", strconv.FormatUint(uint64(incidentClosure.ID), 10), "review-comments", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}

	tx.Commit()
	go func() {
		existingComment, err := icrcs.GetFullDetailsByID(incidentClosureID, id)
		if err != nil {
			utility.LogError("Failed to fetch comment", map[string]interface{}{"error": err.Error()})
			return
		}

		err = icrcs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.IncidentClosureReviewComment{}.TableName(),
			id,
			existingComment,
			currentUserID,
			"incident", "incident_register", "incident_closure_review_comment", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.IncidentClosureReviewComment{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()
	return id, nil
}

func (icrcs *IncidentClosureReviewCommentService) Update(incidentClosureID uint, id uint, currentUserID uint, commentFormDTO *incidentDTOs.IncidentClosureReviewCommentFormDTO) error {
	tx := icrcs.repo.BeginTransaction()

	incidentClosure, err := icrcs.incidentClosureSvc.GetByClosureID(incidentClosureID)
	if err != nil {
		return err
	}

	existingComment, err := icrcs.GetByID(incidentClosureID, id)
	if err != nil {
		return err
	}

	oldComment, err := icrcs.GetFullDetailsByID(incidentClosureID, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingComment.Comment = commentFormDTO.Comment
	existingComment.IsEdited = true
	now := time.Now()
	existingComment.UpdatedBYID = &currentUserID
	existingComment.UpdatedAt = &now
	err = icrcs.repo.Update(existingComment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	var newDocuments []*models.IncidentClosureReviewCommentDocument
	var deleteDocumentIDs []*uint
	var documentsToDelete []*models.IncidentClosureReviewCommentDocument

	for _, doc := range commentFormDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {
			encryptedFilename, err := utility.EncryptData([]byte(doc.FileName))
			if err != nil {
				tx.Rollback()
				return err
			}

			url := utility.GenerateDocumentURL(
				"incident",
				"incident-closure",
				strconv.FormatUint(uint64(incidentClosure.ID), 10),
				"review-comments",
				strconv.FormatUint(uint64(id), 10),
				"documents",
			)

			document := &models.IncidentClosureReviewCommentDocument{
				IncidentClosureReviewCommentID: id,
				URL:                            &url,
				Name:                           doc.FileName,
				Size:                           &doc.Size,
				Token:                          &encryptedFilename,
				OriginalName:                   &doc.OriginalName,
				Type:                           &doc.Type,
				CreatedByID:                    currentUserID,
				CreatedAt:                      now,
			}

			newDocuments = append(newDocuments, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile && doc.ID != nil {
			existingDocument, err := icrcs.repo.GetDocumentByID(*doc.ID)
			if err != nil {
				tx.Rollback()
				return err
			}

			deleteDocumentIDs = append(deleteDocumentIDs, &existingDocument.ID)
			documentsToDelete = append(documentsToDelete, existingDocument)
		}
	}

	if len(newDocuments) > 0 {
		if err := icrcs.repo.CreateDocuments(newDocuments, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIDs) > 0 {
		if err := icrcs.repo.DeleteDocuments(deleteDocumentIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range newDocuments {
		storagePath := filepath.Join("incident", "incident-closure", strconv.FormatUint(uint64(incidentClosure.ID), 10), "review-comments", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range documentsToDelete {
		storagePath := filepath.Join("incident", "incident-closure", strconv.FormatUint(uint64(incidentClosure.ID), 10), "review-comments", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	go func() {
		existingComment, err := icrcs.GetFullDetailsByID(incidentClosureID, id)
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
			err := icrcs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.IncidentClosureReviewComment{}.TableName(),
				id,
				changedData,
				currentUserID,
				"incident", "incident_register", "incident_closure_review_comment", nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.IncidentClosureReviewComment{}.TableName(),
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

func (icrcs *IncidentClosureReviewCommentService) Delete(incidentClosureID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := icrcs.repo.BeginTransaction()
	incidentClosure, err := icrcs.incidentClosureSvc.GetByClosureID(incidentClosureID)
	if err != nil {
		return nil, err
	}

	var failedIDs []uint
	var successfulIDs []uint
	var deletedComments []map[string]interface{}
	var documentsToDelete []*models.IncidentClosureReviewCommentDocument

	for _, id := range ids {
		comment, err := icrcs.repo.GetFullDetailsByID(id)
		if err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		deletedComments = append(deletedComments, map[string]interface{}{
			"id":                  id,
			"incident_closure_id": incidentClosureID,
		})

		if len(comment.Documents) > 0 {
			documentsToDelete = append(documentsToDelete, comment.Documents...)
			var docIDs []*uint
			for _, doc := range comment.Documents {
				docID := doc.ID
				docIDs = append(docIDs, &docID)
			}
			if err := icrcs.repo.DeleteDocuments(docIDs, tx); err != nil {
				failedIDs = append(failedIDs, id)
				continue
			}
		}

		err = icrcs.repo.Delete(id, tx)
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
		storagePath := filepath.Join("incident", "incident-closure", strconv.FormatUint(uint64(incidentClosure.ID), 10), "review-comments", strconv.FormatUint(uint64(doc.IncidentClosureReviewCommentID), 10), "documents")
		if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return failedIDs, err
		}
	}

	tx.Commit()
	if len(successfulIDs) > 0 {
		go func() {
			err := icrcs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.IncidentClosureReviewComment{}.TableName(),
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

func (icrcs *IncidentClosureReviewCommentService) Resolve(incidentClosureID uint, commentID uint, currentUserID uint) error {
	tx := icrcs.repo.BeginTransaction()

	if _, err := icrcs.incidentClosureSvc.GetByClosureID(incidentClosureID); err != nil {
		return err
	}

	existingComment, err := icrcs.GetByID(incidentClosureID, commentID)
	if err != nil {
		return err
	}

	commentStatus, err := icrcs.commentStatusSvc.GetByType("resolved")
	if err != nil {
		return err
	}

	now := time.Now()

	existingComment.CommentStatusID = commentStatus.ID
	existingComment.UpdatedBYID = &currentUserID
	existingComment.UpdatedAt = &now
	err = icrcs.repo.Update(existingComment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (icrcs *IncidentClosureReviewCommentService) GetDocumentPath(incidentClosureID uint, commentID uint, token string) (string, error) {
	comment, err := icrcs.repo.GetFullDetailsByID(commentID)
	if err != nil {
		return "", err
	}
	if len(comment.Documents) == 0 {
		return "", fmt.Errorf("no documents found for the incident closure review comment")
	}
	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}
	storagePath := filepath.Join("incident", "incident-closure", strconv.FormatUint(uint64(comment.IncidentClosureID), 10), "review-comments", strconv.FormatUint(uint64(commentID), 10), "documents")
	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}
