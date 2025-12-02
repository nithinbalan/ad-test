package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"jethur/config"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/grc/risk/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type RiskActionPlanCommentServiceInterface interface {
	GetAll(riskID uint, actionPlanID uint) ([]*riskDTOs.RiskActionPlanCommentListDTO, error)
	GetByID(riskID uint, actionPlanID uint, id uint) (*models.RiskActionPlanComment, error)
	GetFullDetailsByID(riskID uint, actionPlanID uint, id uint, tx ...utility.DatabaseTransaction) (*riskDTOs.RiskActionPlanCommentDetailDTO, error)
	Create(riskID uint, actionPlanID uint, currentUserID uint, dto *riskDTOs.RiskActionPlanCommentFormDTO) (uint, error)
	Update(riskID uint, actionPlanID uint, id uint, currentUserID uint, dto *riskDTOs.RiskActionPlanCommentFormDTO) error
	Delete(riskID uint, actionPlanID uint, ids []uint, currentUserID uint) ([]uint, error)

	GetDocumentPath(riskID, actionPlanID, id uint, token string) (string, error)
}

type RiskActionPlanCommentService struct {
	repo              repositories.RiskActionPlanCommentRepositoryInterface
	userSvc           organizationServices.UserServiceInterface
	activityLogSvc    logServices.ActivityLogServiceInterface
	riskActionPlanSvc RiskActionPlanServiceInterface
	notificationSvc   commonServices.NotificationServiceInterface
}

func NewRiskActionPlanCommentService(
	repo repositories.RiskActionPlanCommentRepositoryInterface,
	userSvc organizationServices.UserServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	riskActionPlanSvc RiskActionPlanServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) RiskActionPlanCommentServiceInterface {
	return &RiskActionPlanCommentService{
		repo:              repo,
		userSvc:           userSvc,
		activityLogSvc:    activityLogSvc,
		riskActionPlanSvc: riskActionPlanSvc,
		notificationSvc:   notificationSvc,
	}
}

func (s *RiskActionPlanCommentService) GetAll(riskID uint, actionPlanID uint) ([]*riskDTOs.RiskActionPlanCommentListDTO, error) {
	if _, err := s.riskActionPlanSvc.GetByID(riskID, actionPlanID); err != nil {
		return nil, err
	}
	return s.repo.GetAll(actionPlanID)
}

func (s *RiskActionPlanCommentService) GetByID(riskID uint, actionPlanID uint, id uint) (*models.RiskActionPlanComment, error) {
	if _, err := s.riskActionPlanSvc.GetByID(riskID, actionPlanID); err != nil {
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

func (s *RiskActionPlanCommentService) GetFullDetailsByID(riskID uint, actionPlanID uint, id uint, tx ...utility.DatabaseTransaction) (*riskDTOs.RiskActionPlanCommentDetailDTO, error) {
	if _, err := s.riskActionPlanSvc.GetByID(riskID, actionPlanID); err != nil {
		return nil, err
	}

	var commentDetail riskDTOs.RiskActionPlanCommentDetailDTO

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

func (rapcs *RiskActionPlanCommentService) Create(riskID uint, actionPlanID uint, currentUserID uint, dto *riskDTOs.RiskActionPlanCommentFormDTO) (uint, error) {
	tx := rapcs.repo.BeginTransaction()

	if _, err := rapcs.riskActionPlanSvc.GetByID(riskID, actionPlanID); err != nil {
		return 0, err
	}

	riskActionPlan, err := rapcs.riskActionPlanSvc.GetFullDetailsByID(riskID, actionPlanID)
	if err != nil {
		return 0, err
	}

	comment := &models.RiskActionPlanComment{
		RiskActionPlanID: actionPlanID,
		Comment:          dto.Comment,
		CreatedByID:      currentUserID,
	}

	id, err := rapcs.repo.Create(comment)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.RiskActionPlanCommentDocument

	for _, doc := range dto.Documents {
		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL(
			"risk", "risk-register", strconv.FormatUint(uint64(riskID), 10),
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
			tx.Rollback()
			return 0, err
		}

		for _, doc := range documents {
			storagePath := filepath.Join(
				"risk", "risk-register", strconv.FormatUint(uint64(riskID), 10),
				"action-plans", strconv.FormatUint(uint64(actionPlanID), 10),
				"comments", strconv.FormatUint(uint64(id), 10),
				"documents",
			)

			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}

	tx.Commit()
	go func() {
		existingComment, err := rapcs.GetFullDetailsByID(riskID, actionPlanID, id)
		if err != nil {
			utility.LogError("Failed to fetch risk action plan comment", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = rapcs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.RiskActionPlanComment{}.TableName(),
			id,
			existingComment,
			currentUserID,
			"risk", "risk_register", "action_plan_comment", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.RiskActionPlanComment{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}

		if riskActionPlan.ResponsibleUserID == currentUserID {
			// if responsible user added the comment, notification for review user
			if riskActionPlan.ReviewUserID != nil {
				// prepare notification
				urlParam := fmt.Sprintf(`{"risk_register_id":%d, "action_plan_id":%d}`, riskID, actionPlanID)
				queryParam := ""
				entity := "risk-action-plan"
				action := "comment-added"

				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserIDs:              []uint{*riskActionPlan.ReviewUserID},
					ModuleIdentifier:     "risk",
					NotificationType:     string(commonDTOs.AssignmentType),
					HeadingEn:            "Comment Added For Risk Action Plan",
					HeadingAr:            "تمت إضافة تعليق لخطة العمل الخاصة بالمخاطر",
					TitleEn:              riskActionPlan.Title,
					TitleAr:              riskActionPlan.Title,
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					QueryParam:           queryParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				reviewUser, err := rapcs.userSvc.GetByID(*riskActionPlan.ReviewUserID)
				if err != nil {
					utility.LogError("Failed to fetch review user", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}
				assignedBy, err := rapcs.userSvc.GetByID(currentUserID)
				if err != nil {
					utility.LogError("Failed to fetch commented user", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}

				emailTo := []string{reviewUser.Email}
				subject := fmt.Sprintf("Comment Added for Risk Action Plan:%s", riskActionPlan.Title)
				appURL := utility.GetAppURL()
				emailData := map[string]interface{}{
					"FirstNameEn":            reviewUser.FirstNameEn,
					"LastNameEn":             reviewUser.LastNameEn,
					"RiskReferenceNumber":    riskActionPlan.RiskRegister.ReferenceNumber,
					"RiskTitle":              riskActionPlan.RiskRegister.Title,
					"ActionPlanTitle":        riskActionPlan.Title,
					"ActionPlanDescription":  riskActionPlan.Description,
					"CommentedByFirstNameEn": assignedBy.FirstNameEn,
					"CommentedByLastNameEn":  assignedBy.LastNameEn,
					"CommentedOn":            time.Now().Format("2006-01-02"),
					"RedirectURL":            fmt.Sprintf("%srisk/risk-register/%d/risk-response", appURL, riskID),
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              emailTo,
					Subject:              subject,
					TemplateRelativePath: "risk/templates/email/risk_action_plan_comment_mail.html",
					EmailData:            emailData,
				}

				// send notification and mail
				err = rapcs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "risk_action_plan")
				if err != nil {
					utility.LogError("Failed to send notifications", map[string]interface{}{
						"error": err.Error(),
					})
				}

			}
		} else if riskActionPlan.ReviewUserID == &currentUserID {
			// if review user added the comment, notification for responsible user
			// prepare notification
			urlParam := fmt.Sprintf(`{"risk_register_id":%d, "action_plan_id":%d}`, riskID, actionPlanID)
			queryParam := ""
			entity := "risk-action-plan"
			action := "comment-added"

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{riskActionPlan.ResponsibleUserID},
				ModuleIdentifier:     "risk",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Comment Added For Risk Action Plan",
				HeadingAr:            "تمت إضافة تعليق لخطة العمل الخاصة بالمخاطر",
				TitleEn:              riskActionPlan.Title,
				TitleAr:              riskActionPlan.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			responsibleUser, err := rapcs.userSvc.GetByID(riskActionPlan.ResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}
			assignedBy, err := rapcs.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch commented user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			emailTo := []string{responsibleUser.Email}
			subject := fmt.Sprintf("Comment Added for Risk Action Plan:%s", riskActionPlan.Title)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":            responsibleUser.FirstNameEn,
				"LastNameEn":             responsibleUser.LastNameEn,
				"RiskReferenceNumber":    riskActionPlan.RiskRegister.ReferenceNumber,
				"RiskTitle":              riskActionPlan.RiskRegister.Title,
				"ActionPlanTitle":        riskActionPlan.Title,
				"ActionPlanDescription":  riskActionPlan.Description,
				"CommentedByFirstNameEn": assignedBy.FirstNameEn,
				"CommentedByLastNameEn":  assignedBy.LastNameEn,
				"CommentedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":            fmt.Sprintf("%srisk/risk-register/%d/risk-response", appURL, riskID),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "risk/templates/email/risk_action_plan_comment_mail.html",
				EmailData:            emailData,
			}

			// send notification and mail
			err = rapcs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "risk_action_plan")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		} else {
			// if other user added the comment, notification for review user and responsible user
			// prepare notification
			urlParam := fmt.Sprintf(`{"risk_register_id":%d, "action_plan_id":%d}`, riskID, actionPlanID)
			queryParam := ""
			entity := "risk-action-plan"
			action := "comment-added"

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{riskActionPlan.ResponsibleUserID},
				ModuleIdentifier:     "risk",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Comment Added For Risk Action Plan",
				HeadingAr:            "تمت إضافة تعليق لخطة العمل الخاصة بالمخاطر",
				TitleEn:              riskActionPlan.Title,
				TitleAr:              riskActionPlan.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			responsibleUser, err := rapcs.userSvc.GetByID(riskActionPlan.ResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}
			assignedBy, err := rapcs.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch commented user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			emailTo := []string{responsibleUser.Email}
			subject := fmt.Sprintf("Comment Added for Risk Action Plan:%s", riskActionPlan.Title)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":            responsibleUser.FirstNameEn,
				"LastNameEn":             responsibleUser.LastNameEn,
				"RiskReferenceNumber":    riskActionPlan.RiskRegister.ReferenceNumber,
				"RiskTitle":              riskActionPlan.RiskRegister.Title,
				"ActionPlanTitle":        riskActionPlan.Title,
				"ActionPlanDescription":  riskActionPlan.Description,
				"CommentedByFirstNameEn": assignedBy.FirstNameEn,
				"CommentedByLastNameEn":  assignedBy.LastNameEn,
				"CommentedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":            fmt.Sprintf("%srisk/risk-register/%d/risk-response", appURL, riskID),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "risk/templates/email/risk_action_plan_comment_mail.html",
				EmailData:            emailData,
			}

			// send notification and mail
			err = rapcs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "risk_action_plan")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
			if riskActionPlan.ReviewUserID != nil {
				// prepare notification
				urlParam := fmt.Sprintf(`{"risk_register_id":%d, "action_plan_id":%d}`, riskID, actionPlanID)
				queryParam := ""
				entity := "risk-action-plan"
				action := "comment-added"

				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserIDs:              []uint{*riskActionPlan.ReviewUserID},
					ModuleIdentifier:     "risk",
					NotificationType:     string(commonDTOs.AssignmentType),
					HeadingEn:            "Comment Added For Risk Action Plan",
					HeadingAr:            "تمت إضافة تعليق لخطة العمل الخاصة بالمخاطر",
					TitleEn:              riskActionPlan.Title,
					TitleAr:              riskActionPlan.Title,
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					QueryParam:           queryParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				responsibleUser, err := rapcs.userSvc.GetByID(*riskActionPlan.ReviewUserID)
				if err != nil {
					utility.LogError("Failed to fetch review user", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}
				assignedBy, err := rapcs.userSvc.GetByID(currentUserID)
				if err != nil {
					utility.LogError("Failed to fetch commented user", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}

				emailTo := []string{responsibleUser.Email}
				subject := fmt.Sprintf("Comment Added for Risk Action Plan:%s", riskActionPlan.Title)
				appURL := utility.GetAppURL()
				emailData := map[string]interface{}{
					"FirstNameEn":            responsibleUser.FirstNameEn,
					"LastNameEn":             responsibleUser.LastNameEn,
					"RiskReferenceNumber":    riskActionPlan.RiskRegister.ReferenceNumber,
					"RiskTitle":              riskActionPlan.RiskRegister.Title,
					"ActionPlanTitle":        riskActionPlan.Title,
					"ActionPlanDescription":  riskActionPlan.Description,
					"CommentedByFirstNameEn": assignedBy.FirstNameEn,
					"CommentedByLastNameEn":  assignedBy.LastNameEn,
					"CommentedOn":            time.Now().Format("2006-01-02"),
					"RedirectURL":            fmt.Sprintf("%srisk/risk-register/%d/risk-response", appURL, riskID),
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              emailTo,
					Subject:              subject,
					TemplateRelativePath: "risk/templates/email/risk_action_plan_comment_mail.html",
					EmailData:            emailData,
				}

				// send notification and mail
				err = rapcs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "risk_action_plan")
				if err != nil {
					utility.LogError("Failed to send notifications", map[string]interface{}{
						"error": err.Error(),
					})
				}
			}
		}
	}()

	return id, nil
}

func (rapcs *RiskActionPlanCommentService) Update(riskID uint, actionPlanID uint, id uint, currentUserID uint, riskActionPlanCommentFormDTO *riskDTOs.RiskActionPlanCommentFormDTO) error {
	tx := rapcs.repo.BeginTransaction()

	if _, err := rapcs.riskActionPlanSvc.GetByID(riskID, actionPlanID); err != nil {
		return err
	}

	existingComment, err := rapcs.GetByID(riskID, actionPlanID, id)
	if err != nil {
		return err
	}

	oldComment, err := rapcs.GetFullDetailsByID(riskID, actionPlanID, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingComment.Comment = riskActionPlanCommentFormDTO.Comment
	err = rapcs.repo.Update(existingComment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	var newDocuments []*models.RiskActionPlanCommentDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.RiskActionPlanCommentDocument

	for _, doc := range riskActionPlanCommentFormDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {
			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}

			url := utility.GenerateDocumentURL("risk", "risk-register", strconv.FormatUint(uint64(riskID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

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
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIds) > 0 {
		if err := rapcs.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range newDocuments {
		storagePath := filepath.Join("risk", "risk-register", strconv.FormatUint(uint64(riskID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("risk", "risk-register", strconv.FormatUint(uint64(riskID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	go func() {
		existingComment, err := rapcs.GetFullDetailsByID(riskID, actionPlanID, id)
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
			err := rapcs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.RiskActionPlanComment{}.TableName(),
				id,
				changedData,
				currentUserID,
				"risk", "risk_register", "action_plan_comment", nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.RiskActionPlanComment{}.TableName(),
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

func (s *RiskActionPlanCommentService) Delete(riskID uint, actionPlanID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := s.repo.BeginTransaction()
	if _, err := s.riskActionPlanSvc.GetByID(riskID, actionPlanID); err != nil {
		return nil, err
	}

	var failedIDs []uint
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

		err := s.repo.Delete(id, tx)
		if err != nil {
			failedIDs = append(failedIDs, id)
		}
	}

	if len(failedIDs) > 0 {
		tx.Rollback()
		return failedIDs, nil
	}

	tx.Commit()
	if len(successfulIDs) > 0 {
		go func() {
			err := s.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.RiskActionPlanComment{}.TableName(),
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

func (s *RiskActionPlanCommentService) GetDocumentPath(riskID, actionPlanID, id uint, token string) (string, error) {
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
		"risk", "risk-register", strconv.FormatUint(uint64(riskID), 10),
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
