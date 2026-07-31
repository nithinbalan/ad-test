package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"jethur/config"
	assetRiskDTOs "jethur/grc/asset/dtos"
	"jethur/grc/asset/repositories"
	logServices "jethur/grc/logs/services"
	masterService "jethur/grc/masters/services"
	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type AssetRiskResponseReviewCommentServiceInterface interface {
	GetAll(riskID uint, currentUserID uint, params *assetRiskDTOs.AssetRiskResponseReviewCommentListQueryParamsDTO) (*assetRiskDTOs.AssetRiskResponseReviewCommentListResponseDTO, error)
	GetByID(riskID uint, commentID uint) (*models.RiskResponseReviewComment, error)
	GetFullDetailsByID(riskID uint, commentID uint) (*assetRiskDTOs.AssetRiskResponseReviewCommentDetailDTO, error)
	Create(riskID uint, currentUserID uint, dto *assetRiskDTOs.AssetRiskResponseReviewCommentFormDTO) (uint, error)
	Update(riskID uint, id uint, currentUserID uint, dto *assetRiskDTOs.AssetRiskResponseReviewCommentFormDTO) error
	Delete(riskID uint, ids []uint, currentUserID uint) ([]uint, error)
	Resolve(riskID uint, commentID uint, currentUserID uint) error
	GetDocumentPath(riskID uint, commentID uint, token string) (string, error)
}

type AssetRiskResponseReviewCommentService struct {
	repo             repositories.AssetRiskResponseReviewCommentRepositoryInterface
	activityLogSvc   logServices.ActivityLogServiceInterface
	riskResponseSvc  AssetRiskResponseServiceInterface
	commentStatusSvc masterService.CommentStatusServiceInterface
}

func NewAssetRiskResponseReviewCommentService(
	repo repositories.AssetRiskResponseReviewCommentRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	riskResponseSvc AssetRiskResponseServiceInterface,
	commentStatusSvc masterService.CommentStatusServiceInterface,
) AssetRiskResponseReviewCommentServiceInterface {
	return &AssetRiskResponseReviewCommentService{
		repo:             repo,
		activityLogSvc:   activityLogSvc,
		riskResponseSvc:  riskResponseSvc,
		commentStatusSvc: commentStatusSvc,
	}
}

func (arrcs *AssetRiskResponseReviewCommentService) GetAll(riskID uint, currentUserID uint, params *assetRiskDTOs.AssetRiskResponseReviewCommentListQueryParamsDTO) (*assetRiskDTOs.AssetRiskResponseReviewCommentListResponseDTO, error) {
	riskResponse, err := arrcs.riskResponseSvc.GetByRiskRegisterID(riskID)
	if err != nil {
		return nil, err
	}

	comments, err := arrcs.repo.GetAll(riskResponse.ID, currentUserID, params)
	if err != nil {
		return nil, err
	}

	dtoMap := make(map[uint]*assetRiskDTOs.AssetRiskResponseReviewCommentListDTO)
	var roots []*assetRiskDTOs.AssetRiskResponseReviewCommentListDTO

	for _, comment := range comments.Comments {
		dtoMap[comment.ID] = comment
	}

	for _, comment := range comments.Comments {
		if comment.ParentCommentID != nil {
			if parent, ok := dtoMap[*comment.ParentCommentID]; ok {
				parent.Children = append(parent.Children, comment)
				continue
			}
		}
		roots = append(roots, comment)
	}

	return &assetRiskDTOs.AssetRiskResponseReviewCommentListResponseDTO{
		Comments:      roots,
		OpenCount:     comments.OpenCount,
		ResolvedCount: comments.ResolvedCount,
	}, nil
}

func (arrcs *AssetRiskResponseReviewCommentService) GetByID(riskID uint, commentID uint) (*models.RiskResponseReviewComment, error) {
	if _, err := arrcs.riskResponseSvc.GetByRiskRegisterID(riskID); err != nil {
		return nil, err
	}

	comment, err := arrcs.repo.GetByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, commentID)
		}
		return nil, err
	}
	return comment, nil
}

func (arrcs *AssetRiskResponseReviewCommentService) GetFullDetailsByID(riskID uint, commentID uint) (*assetRiskDTOs.AssetRiskResponseReviewCommentDetailDTO, error) {
	if _, err := arrcs.riskResponseSvc.GetByRiskRegisterID(riskID); err != nil {
		return nil, err
	}

	var commentDetail assetRiskDTOs.AssetRiskResponseReviewCommentDetailDTO

	comment, err := arrcs.repo.GetFullDetailsByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, commentID)
		}
		return nil, err
	}

	if err := mapstructure.Decode(comment, &commentDetail); err != nil {
		return nil, err
	}

	commentDetail.CreatedAt = comment.CreatedAt

	if len(comment.Documents) > 0 {
		var documentsDTO []*assetRiskDTOs.AssetRiskResponseReviewCommentDocumentDTO
		for _, doc := range comment.Documents {
			documentsDTO = append(documentsDTO, &assetRiskDTOs.AssetRiskResponseReviewCommentDocumentDTO{
				ID:           doc.ID,
				URL:          *doc.URL,
				Name:         doc.Name,
				Size:         *doc.Size,
				Token:        *doc.Token,
				OriginalName: *doc.OriginalName,
				Type:         *doc.Type,
			})
		}
		commentDetail.Documents = documentsDTO
	}

	return &commentDetail, nil
}

func (arrcs *AssetRiskResponseReviewCommentService) Create(riskID uint, currentUserID uint, commentFormDTO *assetRiskDTOs.AssetRiskResponseReviewCommentFormDTO) (id uint, err error) {
	tx := arrcs.repo.BeginTransaction()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil {
			tx.Rollback()
			return
		}

		err = tx.Commit()
	}()

	riskResponse, err := arrcs.riskResponseSvc.GetByRiskRegisterID(riskID)
	if err != nil {
		return 0, err
	}

	commentStatus, err := arrcs.commentStatusSvc.GetByType("open")
	if err != nil {
		return 0, err
	}

	comment := &models.RiskResponseReviewComment{
		RiskResponseID:  riskResponse.ID,
		Comment:         commentFormDTO.Comment,
		CommentStatusID: commentStatus.ID,
		Level:           riskResponse.NextReviewLevel,
		ParentCommentID: commentFormDTO.ParentCommentID,
		CreatedBYID:     currentUserID,
	}

	id, err = arrcs.repo.Create(comment, tx)
	if err != nil {
		return 0, err
	}

	var documents []*models.RiskResponseReviewCommentDocument
	for _, doc := range commentFormDTO.Documents {
		encryptedFilename, err := utility.EncryptData([]byte(doc.FileName))
		if err != nil {
			return 0, err
		}

		url := utility.GenerateDocumentURL("asset", "risk-response", strconv.FormatUint(uint64(riskID), 10), "review-comments", strconv.FormatUint(uint64(id), 10), "documents")

		documents = append(documents, &models.RiskResponseReviewCommentDocument{
			RiskResponseReviewCommentID: id,
			URL:                         &url,
			Name:                        doc.FileName,
			Size:                        &doc.Size,
			Token:                       &encryptedFilename,
			OriginalName:                &doc.OriginalName,
			Type:                        &doc.Type,
			CreatedByID:                 currentUserID,
			CreatedAt:                   time.Now(),
		})
	}

	if len(documents) > 0 {
		if err := arrcs.repo.CreateDocuments(documents, tx); err != nil {
			return 0, err
		}
		for _, doc := range documents {
			storagePath := filepath.Join("asset", "risk-response", strconv.FormatUint(uint64(riskID), 10), "review-comments", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				return 0, err
			}
		}
	}

	go func() {
		existingComment, logErr := arrcs.GetFullDetailsByID(riskID, id)
		if logErr != nil {
			utility.LogError("Failed to fetch comment", map[string]interface{}{
				"error": logErr.Error(),
			})
			return
		}

		logErr = arrcs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.RiskResponseReviewComment{}.TableName(),
			id,
			existingComment,
			currentUserID,
			"asset", "asset_risk", "risk_response_review_comment", nil,
		)
		if logErr != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":   models.CreateAction,
				"recordID": id,
				"error":    logErr.Error(),
			})
		}
	}()

	return id, nil
}

func (arrcs *AssetRiskResponseReviewCommentService) Update(riskID uint, id uint, currentUserID uint, commentFormDTO *assetRiskDTOs.AssetRiskResponseReviewCommentFormDTO) (err error) {
	tx := arrcs.repo.BeginTransaction()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil {
			tx.Rollback()
			return
		}

		err = tx.Commit()
	}()

	if _, err = arrcs.riskResponseSvc.GetByRiskRegisterID(riskID, tx); err != nil {
		return err
	}

	existingComment, err := arrcs.repo.GetByID(id)
	if err != nil {
		return err
	}

	oldComment, err := arrcs.repo.GetFullDetailsByID(id)
	if err != nil {
		return err
	}

	now := time.Now()
	existingComment.Comment = commentFormDTO.Comment
	existingComment.IsEdited = true
	existingComment.UpdatedBYID = &currentUserID
	existingComment.UpdatedAt = now

	if err := arrcs.repo.Update(existingComment, tx); err != nil {
		return err
	}

	var newDocuments []*models.RiskResponseReviewCommentDocument
	var deleteDocumentIDs []*uint
	var documentsToDelete []*models.RiskResponseReviewCommentDocument

	for _, doc := range commentFormDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {
			encryptedFilename, err := utility.EncryptData([]byte(doc.FileName))
			if err != nil {
				return err
			}

			url := utility.GenerateDocumentURL(
				"asset",
				"risk-response",
				strconv.FormatUint(uint64(riskID), 10),
				"review-comments",
				strconv.FormatUint(uint64(id), 10),
				"documents",
			)

			newDocuments = append(newDocuments, &models.RiskResponseReviewCommentDocument{
				RiskResponseReviewCommentID: id,
				URL:                         &url,
				Name:                        doc.FileName,
				Size:                        &doc.Size,
				Token:                       &encryptedFilename,
				OriginalName:                &doc.OriginalName,
				Type:                        &doc.Type,
				CreatedByID:                 currentUserID,
				CreatedAt:                   now,
			})
		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile && doc.ID != nil {
			existingDocument, err := arrcs.repo.GetDocumentByID(*doc.ID)
			if err != nil {
				return err
			}

			deleteDocumentIDs = append(deleteDocumentIDs, &existingDocument.ID)
			documentsToDelete = append(documentsToDelete, existingDocument)
		}
	}

	if len(newDocuments) > 0 {
		if err := arrcs.repo.CreateDocuments(newDocuments, tx); err != nil {
			return err
		}
	}

	if len(deleteDocumentIDs) > 0 {
		if err := arrcs.repo.DeleteDocuments(deleteDocumentIDs, tx); err != nil {
			return err
		}
	}

	for _, doc := range newDocuments {
		storagePath := filepath.Join("asset", "risk-response", strconv.FormatUint(uint64(riskID), 10), "review-comments", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			return err
		}
	}

	for _, doc := range documentsToDelete {
		storagePath := filepath.Join("asset", "risk-response", strconv.FormatUint(uint64(riskID), 10), "review-comments", strconv.FormatUint(uint64(doc.RiskResponseReviewCommentID), 10), "documents")
		if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
			return err
		}
	}

	go func() {
		newComment, logErr := arrcs.repo.GetFullDetailsByID(id)
		if logErr != nil {
			return
		}

		changes := utility.CompareFieldChanges(
			map[string]interface{}{"comment": oldComment.Comment},
			map[string]interface{}{"comment": newComment.Comment},
		)

		if len(changes) > 0 {
			if logErr := arrcs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.RiskResponseReviewComment{}.TableName(),
				id,
				changes,
				currentUserID,
				"asset", "asset_risk", "risk_response_review_comment", nil,
			); logErr != nil {
				utility.LogError("Failed to log activity", map[string]interface{}{
					"action":   models.UpdateAction,
					"recordID": id,
					"error":    logErr.Error(),
				})
			}
		}
	}()

	return nil
}

func (arrcs *AssetRiskResponseReviewCommentService) Delete(riskID uint, ids []uint, currentUserID uint) (failedIDs []uint, err error) {
	tx := arrcs.repo.BeginTransaction()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil || len(failedIDs) > 0 {
			tx.Rollback()
			return
		}

		err = tx.Commit()
	}()

	if _, err = arrcs.riskResponseSvc.GetByRiskRegisterID(riskID, tx); err != nil {
		return nil, err
	}

	var successfulIDs []uint
	var deletedComments []map[string]interface{}
	var documentsToDelete []*models.RiskResponseReviewCommentDocument

	for _, id := range ids {
		comment, getErr := arrcs.repo.GetFullDetailsByID(id)
		if getErr != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		deletedComments = append(deletedComments, map[string]interface{}{
			"id":               id,
			"risk_register_id": riskID,
		})

		if len(comment.Documents) > 0 {
			documentsToDelete = append(documentsToDelete, comment.Documents...)
		}

		if deleteErr := arrcs.repo.Delete(id, tx); deleteErr != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		successfulIDs = append(successfulIDs, id)
	}

	if len(failedIDs) > 0 {
		return failedIDs, nil
	}

	for _, doc := range documentsToDelete {
		storagePath := filepath.Join("asset", "risk-response", strconv.FormatUint(uint64(riskID), 10), "review-comments", strconv.FormatUint(uint64(doc.RiskResponseReviewCommentID), 10), "documents")
		if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
			return failedIDs, err
		}
	}

	if len(successfulIDs) > 0 {
		go func() {
			if logErr := arrcs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.RiskResponseReviewComment{}.TableName(),
				successfulIDs,
				deletedComments,
				currentUserID,
				nil,
			); logErr != nil {
				utility.LogError("Failed to log bulk delete activity", map[string]interface{}{
					"error": logErr.Error(),
				})
			}
		}()
	}

	return failedIDs, nil
}

func (arrcs *AssetRiskResponseReviewCommentService) Resolve(riskID uint, commentID uint, currentUserID uint) (err error) {
	tx := arrcs.repo.BeginTransaction()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil {
			tx.Rollback()
			return
		}

		err = tx.Commit()
	}()

	if _, err = arrcs.riskResponseSvc.GetByRiskRegisterID(riskID, tx); err != nil {
		return err
	}

	existingComment, err := arrcs.repo.GetByID(commentID)
	if err != nil {
		return err
	}

	commentStatus, err := arrcs.commentStatusSvc.GetByType("resolved")
	if err != nil {
		return err
	}

	now := time.Now()
	existingComment.CommentStatusID = commentStatus.ID
	existingComment.UpdatedBYID = &currentUserID
	existingComment.UpdatedAt = now

	return arrcs.repo.Update(existingComment, tx)
}

func (arrcs *AssetRiskResponseReviewCommentService) GetDocumentPath(riskID uint, commentID uint, token string) (string, error) {
	comment, err := arrcs.repo.GetFullDetailsByID(commentID)
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

	storagePath := filepath.Join("asset", "risk-response", strconv.FormatUint(uint64(riskID), 10), "review-comments", strconv.FormatUint(uint64(commentID), 10), "documents")
	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}
