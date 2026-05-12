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

type RiskUpdatePlanReviewCommentServiceInterface interface {
	GetAll(riskRegisterUpdatePlanID uint, currentUserID uint, params *riskDTOs.RiskUpdatePlanReviewCommentListQueryParamsDTO) (*riskDTOs.RiskUpdatePlanReviewCommentListResponseDTO, error)
	GetByID(riskRegisterUpdatePlanID uint, commentID uint) (*models.RiskUpdatePlanReviewComment, error)
	GetFullDetailsByID(riskRegisterUpdatePlanID uint, commentID uint) (*riskDTOs.RiskUpdatePlanReviewCommentDetailDTO, error)
	Create(riskRegisterUpdatePlanID uint, currentUserID uint, dto *riskDTOs.RiskUpdatePlanReviewCommentFormDTO) (uint, error)
	Update(riskRegisterUpdatePlanID uint, id uint, currentUserID uint, dto *riskDTOs.RiskUpdatePlanReviewCommentFormDTO) error
	Delete(riskRegisterUpdatePlanID uint, ids []uint, currentUserID uint) ([]uint, error)
	Resolve(riskRegisterUpdatePlanID uint, commentID uint, currentUserID uint) error
	GetDocumentPath(riskRegisterUpdatePlanID uint, commentID uint, token string) (string, error)
}

type RiskUpdatePlanReviewCommentService struct {
	repo                      repositories.RiskUpdatePlanReviewCommentRepositoryInterface
	userSvc                   organizationServices.UserServiceInterface
	activityLogSvc            logServices.ActivityLogServiceInterface
	riskRegisterUpdatePlanSvc RiskRegisterUpdatePlanServiceInterface
	commentStatusSvc          masterService.CommentStatusServiceInterface
}

func NewRiskUpdatePlanReviewCommentService(
	repo repositories.RiskUpdatePlanReviewCommentRepositoryInterface,
	userSvc organizationServices.UserServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	riskRegisterUpdatePlanSvc RiskRegisterUpdatePlanServiceInterface,
	commentStatusSvc masterService.CommentStatusServiceInterface,
) RiskUpdatePlanReviewCommentServiceInterface {
	return &RiskUpdatePlanReviewCommentService{
		repo:                      repo,
		userSvc:                   userSvc,
		activityLogSvc:            activityLogSvc,
		riskRegisterUpdatePlanSvc: riskRegisterUpdatePlanSvc,
		commentStatusSvc:          commentStatusSvc,
	}
}

func (ruprcs *RiskUpdatePlanReviewCommentService) GetAll(riskRegisterUpdatePlanID uint, currentUserID uint, params *riskDTOs.RiskUpdatePlanReviewCommentListQueryParamsDTO) (*riskDTOs.RiskUpdatePlanReviewCommentListResponseDTO, error) {
	_, err := ruprcs.riskRegisterUpdatePlanSvc.GetByID(riskRegisterUpdatePlanID)
	if err != nil {
		return nil, err
	}

	comments, err := ruprcs.repo.GetAll(riskRegisterUpdatePlanID, currentUserID, params)
	if err != nil {
		return nil, err
	}
	dtoMap := make(map[uint]*riskDTOs.RiskUpdatePlanReviewCommentListDTO)
	var roots []*riskDTOs.RiskUpdatePlanReviewCommentListDTO

	for _, comment := range comments.Comments {
		dto := &riskDTOs.RiskUpdatePlanReviewCommentListDTO{
			ID:                          comment.ID,
			Comment:                     comment.Comment,
			IsEdited:                    comment.IsEdited,
			CommentStatusTitleEn:        comment.CommentStatusTitleEn,
			CommentStatusTitleAr:        comment.CommentStatusTitleAr,
			CommentStatusType:           comment.CommentStatusType,
			CommentStatusColorCode:      comment.CommentStatusColorCode,
			RiskRegisterUpdatePlanID:    comment.RiskRegisterUpdatePlanID,
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

	response := &riskDTOs.RiskUpdatePlanReviewCommentListResponseDTO{
		Comments:      roots,
		OpenCount:     comments.OpenCount,
		ResolvedCount: comments.ResolvedCount,
	}

	return response, nil
}

func (ruprcs *RiskUpdatePlanReviewCommentService) GetByID(riskRegisterUpdatePlanID uint, commentID uint) (*models.RiskUpdatePlanReviewComment, error) {
	if _, err := ruprcs.riskRegisterUpdatePlanSvc.GetByID(riskRegisterUpdatePlanID); err != nil {
		return nil, err
	}

	comment, err := ruprcs.repo.GetByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, riskRegisterUpdatePlanID)
		}
		return nil, err
	}
	return comment, nil
}

func (ruprcs *RiskUpdatePlanReviewCommentService) GetFullDetailsByID(riskRegisterUpdatePlanID uint, commentID uint) (*riskDTOs.RiskUpdatePlanReviewCommentDetailDTO, error) {
	if _, err := ruprcs.riskRegisterUpdatePlanSvc.GetByID(riskRegisterUpdatePlanID); err != nil {
		fmt.Printf("Risk Register Update Plan not found for riskRegisterUpdatePlanID %d", riskRegisterUpdatePlanID)
		return nil, err
	}

	var commentDetail riskDTOs.RiskUpdatePlanReviewCommentDetailDTO

	comment, err := ruprcs.repo.GetFullDetailsByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, riskRegisterUpdatePlanID)
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
		var documentsDTO []*riskDTOs.RiskUpdatePlanReviewCommentDocumentDTO
		for _, doc := range comment.Documents {
			documentDTO := &riskDTOs.RiskUpdatePlanReviewCommentDocumentDTO{
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

func (ruprcs *RiskUpdatePlanReviewCommentService) Create(riskRegisterUpdatePlanID uint, currentUserID uint, commentFormDTO *riskDTOs.RiskUpdatePlanReviewCommentFormDTO) (uint, error) {
	tx := ruprcs.repo.BeginTransaction()

	riskRegisterUpdatePlan, err := ruprcs.riskRegisterUpdatePlanSvc.GetByID(riskRegisterUpdatePlanID)
	if err != nil {
		return 0, err
	}

	commentStatus, err := ruprcs.commentStatusSvc.GetByType("open")
	if err != nil {
		return 0, err
	}

	comment := &models.RiskUpdatePlanReviewComment{
		RiskRegisterUpdatePlanID: riskRegisterUpdatePlanID,
		Comment:                  commentFormDTO.Comment,
		CommentStatusID:          commentStatus.ID,
		Level:                    riskRegisterUpdatePlan.NextReviewLevel,
		ParentCommentID:          commentFormDTO.ParentCommentID,
		CreatedBYID:              currentUserID,
	}

	id, err := ruprcs.repo.Create(comment, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.RiskUpdatePlanReviewCommentDocument
	for _, doc := range commentFormDTO.Documents {
		// Encrypt the filename
		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		// Construct the URL for the document
		url := utility.GenerateDocumentURL("risk", "risk-update-plan", strconv.FormatUint(uint64(riskRegisterUpdatePlanID), 10), "review-comments", strconv.FormatUint(uint64(id), 10), "documents")

		document := &models.RiskUpdatePlanReviewCommentDocument{
			RiskUpdatePlanReviewCommentID: id,
			URL:                           &url,
			Name:                          doc.FileName,
			Size:                          &doc.Size,
			Token:                         &encryptedFilename,
			OriginalName:                  &doc.OriginalName,
			Type:                          &doc.Type,
			CreatedByID:                   currentUserID,
			CreatedAt:                     time.Now(),
		}

		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := ruprcs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		for _, doc := range documents {
			storagePath := filepath.Join("risk", "risk-update-plan", strconv.FormatUint(uint64(riskRegisterUpdatePlanID), 10), "review-comments", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}

	tx.Commit()
	go func() {
		existingComment, err := ruprcs.GetFullDetailsByID(riskRegisterUpdatePlanID, id)
		if err != nil {
			utility.LogError("Failed to fetch comment", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = ruprcs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.RiskUpdatePlanReviewComment{}.TableName(),
			id,
			existingComment,
			currentUserID,
			"risk", "risk_register", "risk_update_plan_review_comment", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.RiskUpdatePlanReviewComment{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()
	return id, nil
}

func (ruprcs *RiskUpdatePlanReviewCommentService) Update(riskRegisterUpdatePlanID uint, id uint, currentUserID uint, commentFormDTO *riskDTOs.RiskUpdatePlanReviewCommentFormDTO) error {
	tx := ruprcs.repo.BeginTransaction()

	if _, err := ruprcs.riskRegisterUpdatePlanSvc.GetByID(riskRegisterUpdatePlanID); err != nil {
		return err
	}

	existingComment, err := ruprcs.GetByID(riskRegisterUpdatePlanID, id)
	if err != nil {
		return err
	}

	oldComment, err := ruprcs.GetFullDetailsByID(riskRegisterUpdatePlanID, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingComment.Comment = commentFormDTO.Comment
	existingComment.IsEdited = true
	now := time.Now()
	existingComment.UpdatedBYID = &currentUserID
	existingComment.UpdatedAt = &now
	err = ruprcs.repo.Update(existingComment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Handle documents
	var newDocuments []*models.RiskUpdatePlanReviewCommentDocument
	var deleteDocumentIDs []*uint
	var documentsToDelete []*models.RiskUpdatePlanReviewCommentDocument

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
				"risk-update-plan",
				strconv.FormatUint(uint64(riskRegisterUpdatePlanID), 10),
				"review-comments",
				strconv.FormatUint(uint64(id), 10),
				"documents",
			)

			document := &models.RiskUpdatePlanReviewCommentDocument{
				RiskUpdatePlanReviewCommentID: id,
				URL:                           &url,
				Name:                          doc.FileName,
				Size:                          &doc.Size,
				Token:                         &encryptedFilename,
				OriginalName:                  &doc.OriginalName,
				Type:                          &doc.Type,
				CreatedByID:                   currentUserID,
				CreatedAt:                     now,
			}

			newDocuments = append(newDocuments, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile && doc.ID != nil {
			// DELETE FILE
			existingDocument, err := ruprcs.repo.GetDocumentByID(*doc.ID)
			if err != nil {
				tx.Rollback()
				return err
			}

			deleteDocumentIDs = append(deleteDocumentIDs, &existingDocument.ID)
			documentsToDelete = append(documentsToDelete, existingDocument)
		}
	}

	if len(newDocuments) > 0 {
		if err := ruprcs.repo.CreateDocuments(newDocuments, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIDs) > 0 {
		if err := ruprcs.repo.DeleteDocuments(deleteDocumentIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// File system operations
	for _, doc := range newDocuments {
		storagePath := filepath.Join("risk", "risk-update-plan", strconv.FormatUint(uint64(riskRegisterUpdatePlanID), 10), "review-comments", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range documentsToDelete {
		storagePath := filepath.Join("risk", "risk-update-plan", strconv.FormatUint(uint64(riskRegisterUpdatePlanID), 10), "review-comments", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	go func() {
		existingComment, err := ruprcs.GetFullDetailsByID(riskRegisterUpdatePlanID, id)
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
			err := ruprcs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.RiskUpdatePlanReviewComment{}.TableName(),
				id,
				changedData,
				currentUserID,
				"risk", "risk_register", "risk_update_plan_review_comment", nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.RiskUpdatePlanReviewComment{}.TableName(),
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

func (ruprcs *RiskUpdatePlanReviewCommentService) Delete(riskRegisterUpdatePlanID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := ruprcs.repo.BeginTransaction()
	if _, err := ruprcs.riskRegisterUpdatePlanSvc.GetByID(riskRegisterUpdatePlanID); err != nil {
		return nil, err
	}

	var failedIDs []uint
	var successfulIDs []uint
	var deletedComments []map[string]interface{}
	var documentsToDelete []*models.RiskUpdatePlanReviewCommentDocument

	for _, id := range ids {
		comment, err := ruprcs.repo.GetFullDetailsByID(id)
		if err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		deletedComments = append(deletedComments, map[string]interface{}{
			"id":                        id,
			"risk_register_update_plan": riskRegisterUpdatePlanID,
		})

		// Collect documents for deletion
		if len(comment.Documents) > 0 {
			documentsToDelete = append(documentsToDelete, comment.Documents...)
		}

		err = ruprcs.repo.Delete(id, tx)
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
		storagePath := filepath.Join("risk", "risk-update-plan", strconv.FormatUint(uint64(riskRegisterUpdatePlanID), 10), "review-comments", strconv.FormatUint(uint64(doc.RiskUpdatePlanReviewCommentID), 10), "documents")
		if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return failedIDs, err
		}
	}

	tx.Commit()
	if len(successfulIDs) > 0 {
		go func() {
			err := ruprcs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.RiskUpdatePlanReviewComment{}.TableName(),
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

func (ruprcs *RiskUpdatePlanReviewCommentService) Resolve(riskRegisterUpdatePlanID uint, commentID uint, currentUserID uint) error {
	tx := ruprcs.repo.BeginTransaction()

	if _, err := ruprcs.riskRegisterUpdatePlanSvc.GetByID(riskRegisterUpdatePlanID); err != nil {
		return err
	}

	existingComment, err := ruprcs.GetByID(riskRegisterUpdatePlanID, commentID)
	if err != nil {
		return err
	}

	commentStatus, err := ruprcs.commentStatusSvc.GetByType("resolved")
	if err != nil {
		return err
	}

	now := time.Now()

	existingComment.CommentStatusID = commentStatus.ID
	existingComment.UpdatedBYID = &currentUserID
	existingComment.UpdatedAt = &now
	err = ruprcs.repo.Update(existingComment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (ruprcs *RiskUpdatePlanReviewCommentService) GetDocumentPath(riskRegisterUpdatePlanID uint, commentID uint, token string) (string, error) {
	comment, err := ruprcs.repo.GetFullDetailsByID(commentID)
	if err != nil {
		return "", err
	}
	if len(comment.Documents) == 0 {
		return "", fmt.Errorf("no documents found for the risk update plan review comment")
	}
	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}
	storagePath := filepath.Join("risk", "risk-update-plan", strconv.FormatUint(uint64(riskRegisterUpdatePlanID), 10), "review-comments", strconv.FormatUint(uint64(commentID), 10), "documents")
	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}
