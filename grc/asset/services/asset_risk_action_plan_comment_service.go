package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"jethur/config"
	assetDTOs "jethur/grc/asset/dtos"
	"jethur/grc/asset/repositories"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type AssetRiskActionPlanCommentServiceInterface interface {
	GetAll(riskID uint, actionPlanID uint) ([]*assetDTOs.AssetRiskActionPlanCommentListDTO, error)
	GetByID(riskID uint, actionPlanID uint, id uint) (*models.RiskActionPlanComment, error)
	GetFullDetailsByID(riskID uint, actionPlanID uint, id uint, tx ...utility.DatabaseTransaction) (*assetDTOs.AssetRiskActionPlanCommentDetailDTO, error)
	Create(riskID uint, actionPlanID uint, currentUserID uint, dto *assetDTOs.AssetRiskActionPlanCommentFormDTO) (uint, error)
	Update(riskID uint, actionPlanID uint, id uint, currentUserID uint, dto *assetDTOs.AssetRiskActionPlanCommentFormDTO) error
	Delete(riskID uint, actionPlanID uint, ids []uint, currentUserID uint) ([]uint, error)

	GetDocumentPath(riskID, actionPlanID, id uint, token string) (string, error)
}

type AssetRiskActionPlanCommentService struct {
	repo                   repositories.AssetRiskActionPlanCommentRepositoryInterface
	userSvc                organizationServices.UserServiceInterface
	activityLogSvc         logServices.ActivityLogServiceInterface
	assetRiskActionPlanSvc AssetRiskActionPlanServiceInterface
	notificationSvc        commonServices.NotificationServiceInterface
}

func NewAssetRiskActionPlanCommentService(
	repo repositories.AssetRiskActionPlanCommentRepositoryInterface,
	userSvc organizationServices.UserServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	assetRiskActionPlanSvc AssetRiskActionPlanServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) AssetRiskActionPlanCommentServiceInterface {
	return &AssetRiskActionPlanCommentService{
		repo:                   repo,
		userSvc:                userSvc,
		activityLogSvc:         activityLogSvc,
		assetRiskActionPlanSvc: assetRiskActionPlanSvc,
		notificationSvc:        notificationSvc,
	}
}

func (s *AssetRiskActionPlanCommentService) GetAll(riskID uint, actionPlanID uint) ([]*assetDTOs.AssetRiskActionPlanCommentListDTO, error) {
	if _, err := s.assetRiskActionPlanSvc.GetByID(riskID, actionPlanID); err != nil {
		return nil, err
	}
	return s.repo.GetAll(actionPlanID)
}

func (s *AssetRiskActionPlanCommentService) GetByID(riskID uint, actionPlanID uint, id uint) (*models.RiskActionPlanComment, error) {
	if _, err := s.assetRiskActionPlanSvc.GetByID(riskID, actionPlanID); err != nil {
		return nil, err
	}

	comment, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return comment, nil
}

func (s *AssetRiskActionPlanCommentService) GetFullDetailsByID(riskID uint, actionPlanID uint, id uint, tx ...utility.DatabaseTransaction) (*assetDTOs.AssetRiskActionPlanCommentDetailDTO, error) {
	if _, err := s.assetRiskActionPlanSvc.GetByID(riskID, actionPlanID); err != nil {
		return nil, err
	}

	var commentDetail assetDTOs.AssetRiskActionPlanCommentDetailDTO

	comment, err := s.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(comment, &commentDetail)
	if err != nil {
		return nil, err
	}

	commentDetail.CreatedAt = comment.CreatedAt
	return &commentDetail, nil
}

func (rapcs *AssetRiskActionPlanCommentService) Create(riskID uint, actionPlanID uint, currentUserID uint, dto *assetDTOs.AssetRiskActionPlanCommentFormDTO) (id uint, err error) {
	existingActionPlan, err := rapcs.assetRiskActionPlanSvc.GetByID(riskID, actionPlanID)
	if err != nil {
		return 0, err
	}
	if riskID == 0 {
		riskID = existingActionPlan.RiskRegisterID
	}

	riskActionPlan, err := rapcs.assetRiskActionPlanSvc.GetFullDetailsByID(riskID, actionPlanID, 0, false)
	if err != nil {
		return 0, err
	}

	tx := rapcs.repo.BeginTransaction()

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

	comment := &models.RiskActionPlanComment{
		RiskActionPlanID: actionPlanID,
		Comment:          dto.Comment,
		CreatedByID:      currentUserID,
	}

	id, err = rapcs.repo.Create(comment, tx)
	if err != nil {
		return 0, err
	}

	var documents []*models.RiskActionPlanCommentDocument

	for _, doc := range dto.Documents {
		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			return 0, err
		}

		url := utility.GenerateDocumentURL(
			"asset", "asset-risk-register", strconv.FormatUint(uint64(riskID), 10),
			"action-plans", strconv.FormatUint(uint64(actionPlanID), 10),
			"comments", strconv.FormatUint(uint64(id), 10),
			"documents",
		)

		document := &models.RiskActionPlanCommentDocument{
			RiskActionPlanCommentID: id,
			URL:                     url,
			Name:                    doc.FileName,
			Size:                    doc.Size,
			Token:                   encryptedFilename,
			OriginalName:            doc.OriginalName,
			Type:                    doc.Type,
			CreatedByID:             currentUserID,
		}

		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := rapcs.repo.CreateDocuments(documents, tx); err != nil {
			return 0, err
		}

		for _, doc := range documents {
			storagePath := filepath.Join(
				"asset", "asset-risk-register", strconv.FormatUint(uint64(riskID), 10),
				"action-plans", strconv.FormatUint(uint64(actionPlanID), 10),
				"comments", strconv.FormatUint(uint64(id), 10),
				"documents",
			)

			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				return 0, err
			}
		}
	}

	go func() {
		existingComment, logErr := rapcs.GetFullDetailsByID(riskID, actionPlanID, id)
		if logErr != nil {
			utility.LogError("Failed to fetch asset risk action plan comment", map[string]interface{}{
				"error": logErr.Error(),
			})
			return
		}

		logErr = rapcs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.RiskActionPlanComment{}.TableName(),
			id,
			existingComment,
			currentUserID,
			"asset", "asset_risk", "action_plan_comment", nil,
		)
		if logErr != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.RiskActionPlanComment{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     logErr.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}

		notifiedByResponsible := riskActionPlan.ResponsibleUserID == currentUserID
		notifiedByReviewer := riskActionPlan.ReviewUserID != nil && *riskActionPlan.ReviewUserID == currentUserID

		if notifiedByResponsible {
			// if responsible user added the comment, notification for review user
			if riskActionPlan.ReviewUserID != nil {
				rapcs.notifyActionPlanComment(riskID, actionPlanID, currentUserID, *riskActionPlan.ReviewUserID, riskActionPlan)
			}
		} else if notifiedByReviewer {
			// if review user added the comment, notification for responsible user
			rapcs.notifyActionPlanComment(riskID, actionPlanID, currentUserID, riskActionPlan.ResponsibleUserID, riskActionPlan)
		} else {
			// if other user added the comment, notification for review user and responsible user
			rapcs.notifyActionPlanComment(riskID, actionPlanID, currentUserID, riskActionPlan.ResponsibleUserID, riskActionPlan)
			if riskActionPlan.ReviewUserID != nil {
				rapcs.notifyActionPlanComment(riskID, actionPlanID, currentUserID, *riskActionPlan.ReviewUserID, riskActionPlan)
			}
		}
	}()

	return id, nil
}

func (rapcs *AssetRiskActionPlanCommentService) notifyActionPlanComment(riskID uint, actionPlanID uint, currentUserID uint, recipientUserID uint, riskActionPlan *assetDTOs.AssetRiskActionPlanDetailDTO) {
	urlParam := fmt.Sprintf(`{"risk_register_id":%d, "action_plan_id":%d}`, riskID, actionPlanID)
	queryParam := ""
	entity := "asset-risk-action-plan"
	action := "comment-added"

	dbNotificationDTO := &commonDTOs.DBNotificationDTO{
		UserIDs:              []uint{recipientUserID},
		ModuleIdentifier:     "asset_risk",
		NotificationType:     string(commonDTOs.AssignmentType),
		HeadingEn:            "Comment Added For ISMS Risk Action Plan",
		HeadingAr:            "تمت إضافة تعليق لخطة عمل مخاطر ISMS",
		TitleEn:              riskActionPlan.Title,
		TitleAr:              riskActionPlan.Title,
		Entity:               &entity,
		Action:               &action,
		URLParam:             urlParam,
		QueryParam:           queryParam,
		IsSystemNotification: false,
		CreatedByID:          &currentUserID,
	}

	recipientUser, err := rapcs.userSvc.GetByID(recipientUserID)
	if err != nil {
		utility.LogError("Failed to fetch recipient user", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	commentedBy, err := rapcs.userSvc.GetByID(currentUserID)
	if err != nil {
		utility.LogError("Failed to fetch commented user", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	actionPlanDescription := "-"
	if riskActionPlan.Description != "" {
		actionPlanDescription = riskActionPlan.Description
	}

	emailTo := []string{recipientUser.Email}
	subject := fmt.Sprintf("ISMS Risk Action Plan: %s - Comment Added", riskActionPlan.Title)
	appURL := utility.GetAppURL()
	emailData := map[string]interface{}{
		"FirstNameEn":            recipientUser.FirstNameEn,
		"LastNameEn":             recipientUser.LastNameEn,
		"RiskReferenceNumber":    riskActionPlan.RiskRegister.ReferenceNumber,
		"RiskTitle":              riskActionPlan.RiskRegister.Title,
		"ActionPlanTitle":        riskActionPlan.Title,
		"ActionPlanDescription":  actionPlanDescription,
		"CommentedByFirstNameEn": commentedBy.FirstNameEn,
		"CommentedByLastNameEn":  commentedBy.LastNameEn,
		"CommentedOn":            time.Now().Format("2006-01-02"),
		"RedirectURL":            fmt.Sprintf("%s/asset/risk-mitigations?id=%d", appURL, actionPlanID),
	}

	emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
		EmailTo:              emailTo,
		Subject:              subject,
		TemplateRelativePath: "risk/templates/email/risk_action_plan_comment_mail.html",
		EmailData:            emailData,
	}

	if err := rapcs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "asset_risk_action_plan"); err != nil {
		utility.LogError("Failed to send notifications", map[string]interface{}{
			"error": err.Error(),
		})
	}
}

func (rapcs *AssetRiskActionPlanCommentService) Update(riskID uint, actionPlanID uint, id uint, currentUserID uint, formDTO *assetDTOs.AssetRiskActionPlanCommentFormDTO) (err error) {
	existingActionPlan, err := rapcs.assetRiskActionPlanSvc.GetByID(riskID, actionPlanID)
	if err != nil {
		return err
	}
	if riskID == 0 {
		riskID = existingActionPlan.RiskRegisterID
	}

	existingComment, err := rapcs.GetByID(riskID, actionPlanID, id)
	if err != nil {
		return err
	}

	oldComment, err := rapcs.GetFullDetailsByID(riskID, actionPlanID, id)
	if err != nil {
		return err
	}

	tx := rapcs.repo.BeginTransaction()

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

	existingComment.Comment = formDTO.Comment
	err = rapcs.repo.Update(existingComment, tx)
	if err != nil {
		return err
	}

	var newDocuments []*models.RiskActionPlanCommentDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.RiskActionPlanCommentDocument

	for _, doc := range formDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {
			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				return err
			}

			url := utility.GenerateDocumentURL("asset", "asset-risk-register", strconv.FormatUint(uint64(riskID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

			document := &models.RiskActionPlanCommentDocument{
				RiskActionPlanCommentID: id,
				URL:                     url,
				Name:                    doc.FileName,
				Size:                    doc.Size,
				Token:                   encryptedFileName,
				OriginalName:            doc.OriginalName,
				Type:                    doc.Type,
				CreatedByID:             currentUserID,
			}

			newDocuments = append(newDocuments, document)
		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			if doc.ID != nil {
				existingDocument, err := rapcs.repo.GetDocumentByID(*doc.ID)
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return utility.NewResourceNotFoundError(err, id)
					}
					return err
				}

				documentsToDelete = append(documentsToDelete, existingDocument)
				deleteDocumentIds = append(deleteDocumentIds, &existingDocument.ID)
			}
		}
	}

	if len(newDocuments) > 0 {
		if err := rapcs.repo.CreateDocuments(newDocuments, tx); err != nil {
			return err
		}
	}

	if len(deleteDocumentIds) > 0 {
		if err := rapcs.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			return err
		}
	}

	for _, doc := range newDocuments {
		storagePath := filepath.Join("asset", "asset-risk-register", strconv.FormatUint(uint64(riskID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("asset", "asset-risk-register", strconv.FormatUint(uint64(riskID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			return err
		}
	}

	go func() {
		existingComment, logErr := rapcs.GetFullDetailsByID(riskID, actionPlanID, id)
		if logErr != nil {
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
			logErr := rapcs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.RiskActionPlanComment{}.TableName(),
				id,
				changedData,
				currentUserID,
				"asset", "asset_risk", "action_plan_comment", nil,
			)
			if logErr != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.RiskActionPlanComment{}.TableName(),
					"recordID":  id,
					"userID":    currentUserID,
					"error":     logErr.Error(),
				}

				utility.LogError("Failed to log activity", fields)
			}
		}
	}()
	return nil
}

func (s *AssetRiskActionPlanCommentService) Delete(riskID uint, actionPlanID uint, ids []uint, currentUserID uint) (failedIDs []uint, err error) {
	if _, err := s.assetRiskActionPlanSvc.GetByID(riskID, actionPlanID); err != nil {
		return nil, err
	}

	tx := s.repo.BeginTransaction()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if len(failedIDs) > 0 {
			tx.Rollback()
			return
		}

		tx.Commit()
	}()

	var successfulIDs []uint
	var deletedComments []map[string]interface{}
	for _, id := range ids {
		if _, err := s.repo.GetFullDetailsByID(id); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		deletedComments = append(deletedComments, map[string]interface{}{
			"id":          id,
			"risk_action": actionPlanID,
		})

		if err := s.repo.Delete(id, tx); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		successfulIDs = append(successfulIDs, id)
	}

	if len(successfulIDs) > 0 {
		go func() {
			logErr := s.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.RiskActionPlanComment{}.TableName(),
				successfulIDs,
				deletedComments,
				currentUserID,
				nil,
			)
			if logErr != nil {
				utility.LogError("Failed to log bulk delete activity", map[string]interface{}{
					"error": logErr.Error(),
				})
			}
		}()
	}
	return failedIDs, nil
}

func (s *AssetRiskActionPlanCommentService) GetDocumentPath(riskID, actionPlanID, id uint, token string) (string, error) {
	existingComment, err := s.GetByID(riskID, actionPlanID, id)
	if err != nil {
		return "", err
	}
	if len(existingComment.Documents) == 0 {
		return "", fmt.Errorf("no documents found for the risk action plan comment")
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	storagePath := filepath.Join(
		"asset", "asset-risk-register", strconv.FormatUint(uint64(riskID), 10),
		"action-plans", strconv.FormatUint(uint64(actionPlanID), 10),
		"comments", strconv.FormatUint(uint64(id), 10),
		"documents",
	)

	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}
