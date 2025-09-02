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
	threatDTOs "jethur/grc/asset/dtos"
	"jethur/grc/asset/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ThreatActionPlanCommentServiceInterface interface {
	GetAll(threatID uint, actionPlanID uint) ([]*threatDTOs.ThreatActionPlanCommentListDTO, error)
	GetByID(threatID uint, actionPlanID uint, id uint) (*models.ThreatActionPlanComment, error)
	GetFullDetailsByID(threatID uint, actionPlanID uint, id uint, tx ...utility.DatabaseTransaction) (*threatDTOs.ThreatActionPlanCommentDetailDTO, error)
	Create(threatID uint, actionPlanID uint, currentUserID uint, dto *threatDTOs.ThreatActionPlanCommentFormDTO) (uint, error)
	Update(threatID uint, actionPlanID uint, id uint, currentUserID uint, dto *threatDTOs.ThreatActionPlanCommentFormDTO) error
	Delete(threatID uint, actionPlanID uint, ids []uint, currentUserID uint) ([]uint, error)

	GetDocumentPath(threatID, actionPlanID, id uint, token string) (string, error)
}

type ThreatActionPlanCommentService struct {
	repo                    repositories.ThreatActionPlanCommentRepositoryInterface
	userSvc                 organizationServices.UserServiceInterface
	activityLogSvc          logServices.ActivityLogServiceInterface
	threatActionPlanSvc ThreatActionPlanServiceInterface
	notificationSvc         commonServices.NotificationServiceInterface
}

func NewThreatActionPlanCommentService(repo repositories.ThreatActionPlanCommentRepositoryInterface, userSvc organizationServices.UserServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	threatActionPlanSvc ThreatActionPlanServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) ThreatActionPlanCommentServiceInterface {
	return &ThreatActionPlanCommentService{
		repo:                    repo,
		userSvc:                 userSvc,
		activityLogSvc:          activityLogSvc,
	    threatActionPlanSvc: threatActionPlanSvc,
		notificationSvc:         notificationSvc,
	}
}

func (caps *ThreatActionPlanCommentService) GetAll(threatID uint, actionPlanID uint) ([]*threatDTOs.ThreatActionPlanCommentListDTO, error) {
	if _, err := caps.threatActionPlanSvc.GetByID(threatID, actionPlanID); err != nil {
		return nil, err
	}

	return caps.repo.GetAll(actionPlanID)
}

func (caps *ThreatActionPlanCommentService) GetByID(threatID uint, actionPlanID uint, id uint) (*models.ThreatActionPlanComment, error) {
	if _, err := caps.threatActionPlanSvc.GetByID(threatID, actionPlanID); err != nil {
		return nil, err
	}

	comment, err := caps.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return comment, nil
}

func (caps *ThreatActionPlanCommentService) GetFullDetailsByID(threatID uint, actionPlanID uint, id uint, tx ...utility.DatabaseTransaction) (*threatDTOs.ThreatActionPlanCommentDetailDTO, error) {
	if _, err := caps.threatActionPlanSvc.GetByID(threatID, actionPlanID); err != nil {
		fmt.Printf("ThreatActionPlan not found for threatID %d, actionPlanID %d\n", threatID, actionPlanID)
		return nil, err
	}

	var commentDetail threatDTOs.ThreatActionPlanCommentDetailDTO

	comment, err := caps.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(comment, &commentDetail)
	if err != nil {
		fmt.Printf("Error decoding comment: %v\n", err)
		return nil, err
	}

	commentDetail.CreatedAt = comment.CreatedAt
	return &commentDetail, nil
}

func (caps *ThreatActionPlanCommentService) Create(threatID uint, actionPlanID uint, currentUserID uint, threatActionPlanCommentFormDTO *threatDTOs.ThreatActionPlanCommentFormDTO) (uint, error) {
	tx := caps.repo.BeginTransaction()

	if _, err := caps.threatActionPlanSvc.GetByID(threatID, actionPlanID); err != nil {
		return 0, err
	}

	threatActionPlan, err := caps.threatActionPlanSvc.GetFullDetailsByID(threatID, actionPlanID)


	if err != nil {
		return 0, err
	}

	comment := &models.ThreatActionPlanComment{
		ThreatActionPlanID: actionPlanID,
		Comment:                threatActionPlanCommentFormDTO.Comment,
		CreatedBYID:            currentUserID,
	}

	id, err := caps.repo.Create(comment)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.ThreatActionPlanCommentDocument

	for _, doc := range threatActionPlanCommentFormDTO.Documents {
		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL(
			"threat", "threat", strconv.FormatUint(uint64(threatID), 10),
			"action-plans", strconv.FormatUint(uint64(actionPlanID), 10),
			"comments", strconv.FormatUint(uint64(id), 10),
			"documents",
		)

		document := &models.ThreatActionPlanCommentDocument{
			ThreatActionPlanCommentID: id,
			URL:                           url,
			Name:                          doc.FileName,
			Size:                          doc.Size,
			Token:                         encryptedFilename,
			OriginalName:                  doc.OriginalName,
			Type:                          doc.Type,
			CreatedBYID:                   currentUserID,
		}

		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := caps.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}

		for _, doc := range documents {
			storagePath := filepath.Join(
				"threat", "threat", strconv.FormatUint(uint64(threatID), 10),
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
		existingComment, err := caps.GetFullDetailsByID(threatID, actionPlanID, id)
		if err != nil {
			utility.LogError("Failed to fetch threat action plan comment", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = caps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ThreatActionPlanComment{}.TableName(),
			id,
			existingComment,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ThreatActionPlanComment{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
		if uint(threatActionPlan.ResponsibleUser.ID) == currentUserID {
			// if responsible user added the comment, notification for review user
			if threatActionPlan.ReviewUser != nil {
				   

				// prepare notification
				urlParam := fmt.Sprintf(`{"_threat_id":%d, "action_plan_id":%d}`, threatID, id)
				queryParam := ""
				entity := "threat-action-plan"
				action := "comment-added"

				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserIDs:              []uint{uint(threatActionPlan.ReviewUser.ID)},
					ModuleIdentifier:     "threat",
					NotificationType:     string(commonDTOs.AssignmentType),
					HeadingEn:            "Comment Added For Threat Action Plan",
					HeadingAr:            "تمت إضافة التعليق لخطة عمل الامتثال",
					TitleEn:              threatActionPlan.Title,
					TitleAr:              threatActionPlan.Title,
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					QueryParam:           queryParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				reviewUser, err := caps.userSvc.GetByID(uint(threatActionPlan.ReviewUser.ID))
				if err != nil {
					utility.LogError("Failed to fetch review user", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}
				assignedBy, err := caps.userSvc.GetByID(currentUserID)
				if err != nil {
					utility.LogError("Failed to fetch commented user", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}

				emailTo := []string{reviewUser.Email}
				subject := fmt.Sprintf("Comment Added for Threat Action Plan:%s", threatActionPlan.Title)
				appURL := utility.GetAppURL()
				emailData := map[string]interface{}{
					"FirstNameEn":               reviewUser.FirstNameEn,
					"LastNameEn":                reviewUser.LastNameEn,
					"ThreatReferenceNumber": threatActionPlan.Threat.ID,
					"ThreatTitle":           threatActionPlan.Threat.TitleEn,
					"ActionPlanTitle":           threatActionPlan.Title,
					"ActionPlanDescription":     threatActionPlan.Description,
					"CommentedByFirstNameEn":    assignedBy.FirstNameEn,
					"CommentedByLastNameEn":     assignedBy.LastNameEn,
					"CommentedOn":               time.Now().Format("2006-01-02"),
					"RedirectURL":               fmt.Sprintf("%s/threat/register/%d/action-plans?action_plan_id=%d", appURL, threatID, id),
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              emailTo,
					Subject:              subject,
					TemplateRelativePath: "threat/templates/email/threat_action_plan_comment.html",
					EmailData:            emailData,
				}

				// send notification and mail
				err = caps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "threat_action_plan")
				if err != nil {
					utility.LogError("Failed to send notifications", map[string]interface{}{
						"error": err.Error(),
					})
				}
			}
		} else if  threatActionPlan.ReviewUser != nil && uint(threatActionPlan.ReviewUser.ID) == currentUserID {
			// if review user added the comment, notification for responsible user
			// prepare notification
			urlParam := fmt.Sprintf(`{"threat_id":%d, "action_plan_id":%d}`, threatID, id)
			queryParam := ""
			entity := "threat-action-plan"
			action := "comment-added"

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{uint(threatActionPlan.ResponsibleUser.ID)},
				ModuleIdentifier:     "threat",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Comment Added For Threat Action Plan",
				HeadingAr:            "تمت إضافة التعليق لخطة عمل الامتثال",
				TitleEn:              threatActionPlan.Title,
				TitleAr:              threatActionPlan.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			responsibleUser, err := caps.userSvc.GetByID(uint(threatActionPlan.ResponsibleUser.ID))
			if err != nil {
				utility.LogError("Failed to fetch responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}
			assignedBy, err := caps.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch commented user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			emailTo := []string{responsibleUser.Email}
			subject := fmt.Sprintf("Comment Added for Threat Action Plan:%s", threatActionPlan.Title)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":               responsibleUser.FirstNameEn,
				"LastNameEn":                responsibleUser.LastNameEn,
				"ThreatReferenceNumber": threatActionPlan.Threat.ID,
				"ThreatTitle":           threatActionPlan.Threat.TitleEn,
				"ActionPlanTitle":           threatActionPlan.Title,
				"ActionPlanDescription":     threatActionPlan.Description,
				"CommentedByFirstNameEn":    assignedBy.FirstNameEn,
				"CommentedByLastNameEn":     assignedBy.LastNameEn,
				"AssignedOn":                time.Now().Format("2006-01-02"),
				"RedirectURL":               fmt.Sprintf("%s/threat/register/%d/action-plans?action_plan_id=%d", appURL, threatID, id),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "threat/templates/email/threat_action_plan_comment.html",
				EmailData:            emailData,
			}

			// send notification and mail
			err = caps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "threat_action_plan")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		} else {
			// if other user added the comment, notification for review user and responsible user
			// prepare notification
			urlParam := fmt.Sprintf(`{"threat_id":%d, "action_plan_id":%d}`, threatID, id)
			queryParam := ""
			entity := "threat-action-plan"
			action := "comment-added"

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{uint(threatActionPlan.ResponsibleUser.ID)},
				ModuleIdentifier:     "threat",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Comment Added For Threat Action Plan",
				HeadingAr:            "تمت إضافة التعليق لخطة عمل الامتثال",
				TitleEn:              threatActionPlan.Title,
				TitleAr:              threatActionPlan.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			responsibleUser, err := caps.userSvc.GetByID(uint(threatActionPlan.ResponsibleUser.ID))
			if err != nil {
				utility.LogError("Failed to fetch responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}
			assignedBy, err := caps.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch commented user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			emailTo := []string{responsibleUser.Email}
			subject := fmt.Sprintf("Comment Added for Threat Action Plan:%s", threatActionPlan.Title)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":               responsibleUser.FirstNameEn,
				"LastNameEn":                responsibleUser.LastNameEn,
				"ThreatReferenceNumber": threatActionPlan.Threat.ID,
				"ThreatTitle":           threatActionPlan.Threat.TitleEn,
				"ActionPlanTitle":           threatActionPlan.Title,
				"ActionPlanDescription":     threatActionPlan.Description,
				"CommentedByFirstNameEn":    assignedBy.FirstNameEn,
				"CommentedByLastNameEn":     assignedBy.LastNameEn,
				"AssignedOn":                time.Now().Format("2006-01-02"),
				"RedirectURL":               fmt.Sprintf("%s/threat/register/%d/action-plans?action_plan_id=%d", appURL, threatID, id),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "threat/templates/email/threat_action_plan_comment.html",
				EmailData:            emailData,
			}

			// send notification and mail
			err = caps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "threat_action_plan")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
			if threatActionPlan.ReviewUser != nil {
				// prepare notification
				urlParam := fmt.Sprintf(`{"threat_id":%d, "action_plan_id":%d}`, threatID, id)
				queryParam := ""
				entity := "threat-action-plan"
				action := "comment-added"

				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserIDs:              []uint{uint(threatActionPlan.ReviewUser.ID)},
					ModuleIdentifier:     "threat",
					NotificationType:     string(commonDTOs.AssignmentType),
					HeadingEn:            "Comment Added For Threat Action Plan",
					HeadingAr:            "تمت إضافة التعليق لخطة عمل الامتثال",
					TitleEn:              threatActionPlan.Title,
					TitleAr:              threatActionPlan.Title,
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					QueryParam:           queryParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				responsibleUser, err := caps.userSvc.GetByID(uint(threatActionPlan.ReviewUser.ID))
				if err != nil {
					utility.LogError("Failed to fetch review user", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}
				assignedBy, err := caps.userSvc.GetByID(currentUserID)
				if err != nil {
					utility.LogError("Failed to fetch commented user", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}

				emailTo := []string{responsibleUser.Email}
				subject := fmt.Sprintf("Comment Added for Threat Action Plan:%s", threatActionPlan.Title)
				appURL := utility.GetAppURL()
				emailData := map[string]interface{}{
					"FirstNameEn":               responsibleUser.FirstNameEn,
					"LastNameEn":                responsibleUser.LastNameEn,
					"ThreatReferenceNumber": threatActionPlan.Threat.ID,
					"ThreatTitle":           threatActionPlan.Threat.TitleEn,
					"ActionPlanTitle":           threatActionPlan.Title,
					"ActionPlanDescription":     threatActionPlan.Description,
					"CommentedByFirstNameEn":    assignedBy.FirstNameEn,
					"CommentedByLastNameEn":     assignedBy.LastNameEn,
					"AssignedOn":                time.Now().Format("2006-01-02"),
					"RedirectURL":               fmt.Sprintf("%s/threat/register/%d/action-plans?action_plan_id=%d", appURL, threatID, id),
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              emailTo,
					Subject:              subject,
					TemplateRelativePath: "threat/templates/email/threat_action_plan_comment.html",
					EmailData:            emailData,
				}

				// send notification and mail
				err = caps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "threat_action_plan")
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
func (caps *ThreatActionPlanCommentService) Update(threatID uint, actionPlanID uint, id uint, currentUserID uint, threatActionPlanCommentFormDTO *threatDTOs.ThreatActionPlanCommentFormDTO) error {
	tx := caps.repo.BeginTransaction()

	if _, err := caps.threatActionPlanSvc.GetByID(threatID, actionPlanID); err != nil {
		return err
	}

	existingComment, err := caps.GetByID(threatID, actionPlanID, id)
	if err != nil {
		return err
	}

	oldComment, err := caps.GetFullDetailsByID(threatID, actionPlanID, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingComment.Comment = threatActionPlanCommentFormDTO.Comment
	err = caps.repo.Update(existingComment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	var newDocuments []*models.ThreatActionPlanCommentDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.ThreatActionPlanCommentDocument

	for _, doc := range threatActionPlanCommentFormDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {

			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}

			url := utility.GenerateDocumentURL("threat", "threat", strconv.FormatUint(uint64(threatID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

			document := &models.ThreatActionPlanCommentDocument{
				ThreatActionPlanCommentID: id,
				URL:                           url,
				Name:                          doc.FileName,
				Size:                          doc.Size,
				Token:                         encryptedFileName,
				OriginalName:                  doc.OriginalName,
				Type:                          doc.Type,
				CreatedBYID:                   currentUserID,
			}

			newDocuments = append(newDocuments, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			if doc.ID != nil {
				existingDocument, err := caps.repo.GetDocumentByID(*doc.ID)
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
		if err := caps.repo.CreateDocuments(newDocuments, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIds) > 0 {
		if err := caps.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range newDocuments {
		storagePath := filepath.Join("threat", "threat", strconv.FormatUint(uint64(threatID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("threat", "threat", strconv.FormatUint(uint64(threatID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	go func() {
		existingComment, err := caps.GetFullDetailsByID(threatID, actionPlanID, id)
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
			err := caps.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.ThreatActionPlanComment{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.ThreatActionPlanComment{}.TableName(),
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

func (caps *ThreatActionPlanCommentService) Delete(threatID uint, actionPlanID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := caps.repo.BeginTransaction()
	if _, err := caps.threatActionPlanSvc.GetByID(threatID, actionPlanID); err != nil {
		return nil, err
	}

	var failedIDs []uint
	var successfulIDs []uint
	var deletedComments []map[string]interface{}
	for _, id := range ids {
		if _, err := caps.repo.GetFullDetailsByID(id); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		deletedComments = append(deletedComments, map[string]interface{}{
			"id":                id,
			"threat_action": actionPlanID,
		})

		err := caps.repo.Delete(id, tx)
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
			err := caps.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.ThreatActionPlanComment{}.TableName(),
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

func (caps *ThreatActionPlanCommentService) GetDocumentPath(threatID, actionPlanID, id uint, token string) (string, error) {
	existingComment, err := caps.GetByID(threatID, actionPlanID, id)
	if err != nil {
		return "", err
	}
	if len(existingComment.Documents) == 0 {
		return "", fmt.Errorf("no documents found for the threat action plan comment")
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	storagePath := filepath.Join("threat", "threat", strconv.FormatUint(uint64(threatID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}
