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
	controlDTOs "jethur/grc/organization/dtos"
	"jethur/grc/organization/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ControlActionPlanCommentServiceInterface interface {
	GetAll(controlID uint, actionPlanID uint) ([]*controlDTOs.ControlActionPlanCommentListDTO, error)
	GetByID(controlID uint, actionPlanID uint, id uint) (*models.ControlActionPlanComment, error)
	GetFullDetailsByID(controlID uint, actionPlanID uint, id uint, tx ...utility.DatabaseTransaction) (*controlDTOs.ControlActionPlanCommentDetailDTO, error)
	Create(controlID uint, actionPlanID uint, currentUserID uint, dto *controlDTOs.ControlActionPlanCommentFormDTO) (uint, error)
	Update(controlID uint, actionPlanID uint, id uint, currentUserID uint, dto *controlDTOs.ControlActionPlanCommentFormDTO) error
	Delete(controlID uint, actionPlanID uint, ids []uint, currentUserID uint) ([]uint, error)

	GetDocumentPath(controlID, actionPlanID, id uint, token string) (string, error)
}

type ControlActionPlanCommentService struct {
	repo                 repositories.ControlActionPlanCommentRepositoryInterface
	userSvc              UserServiceInterface
	activityLogSvc       logServices.ActivityLogServiceInterface
	controlActionPlanSvc ControlActionPlanServiceInterface
	notificationSvc      commonServices.NotificationServiceInterface
}

func NewControlActionPlanCommentService(repo repositories.ControlActionPlanCommentRepositoryInterface, userSvc UserServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	controlActionPlanSvc ControlActionPlanServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) ControlActionPlanCommentServiceInterface {
	return &ControlActionPlanCommentService{
		repo:                 repo,
		userSvc:              userSvc,
		activityLogSvc:       activityLogSvc,
		controlActionPlanSvc: controlActionPlanSvc,
		notificationSvc:      notificationSvc,
	}
}

func (caps *ControlActionPlanCommentService) GetAll(controlID uint, actionPlanID uint) ([]*controlDTOs.ControlActionPlanCommentListDTO, error) {
	if _, err := caps.controlActionPlanSvc.GetByID(controlID, actionPlanID); err != nil {
		return nil, err
	}

	return caps.repo.GetAll(actionPlanID)
}

func (caps *ControlActionPlanCommentService) GetByID(controlID uint, actionPlanID uint, id uint) (*models.ControlActionPlanComment, error) {
	if _, err := caps.controlActionPlanSvc.GetByID(controlID, actionPlanID); err != nil {
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

func (caps *ControlActionPlanCommentService) GetFullDetailsByID(controlID uint, actionPlanID uint, id uint, tx ...utility.DatabaseTransaction) (*controlDTOs.ControlActionPlanCommentDetailDTO, error) {
	if _, err := caps.controlActionPlanSvc.GetByID(controlID, actionPlanID); err != nil {
		fmt.Printf("ControlActionPlan not found for controlID %d, actionPlanID %d\n", controlID, actionPlanID)
		return nil, err
	}

	var commentDetail controlDTOs.ControlActionPlanCommentDetailDTO

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

func (caps *ControlActionPlanCommentService) Create(controlID uint, actionPlanID uint, currentUserID uint, controlActionPlanCommentFormDTO *controlDTOs.ControlActionPlanCommentFormDTO) (uint, error) {
	tx := caps.repo.BeginTransaction()

	if _, err := caps.controlActionPlanSvc.GetByID(controlID, actionPlanID); err != nil {
		return 0, err
	}

	controlActionPlan, err := caps.controlActionPlanSvc.GetFullDetailsByID(controlID, actionPlanID)
	if err != nil {
		return 0, err
	}

	comment := &models.ControlActionPlanComment{
		ControlActionPlanID: actionPlanID,
		Comment:             controlActionPlanCommentFormDTO.Comment,
		CreatedBYID:         currentUserID,
	}

	id, err := caps.repo.Create(comment)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.ControlActionPlanCommentDocument

	for _, doc := range controlActionPlanCommentFormDTO.Documents {
		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL(
			"control", "control-register", strconv.FormatUint(uint64(controlID), 10),
			"action-plans", strconv.FormatUint(uint64(actionPlanID), 10),
			"comments", strconv.FormatUint(uint64(id), 10),
			"documents",
		)

		document := &models.ControlActionPlanCommentDocument{
			ControlActionPlanCommentID: id,
			URL:                        url,
			Name:                       doc.FileName,
			Size:                       doc.Size,
			Token:                      encryptedFilename,
			OriginalName:               doc.OriginalName,
			Type:                       doc.Type,
			CreatedBYID:                currentUserID,
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
				"control", "control-register", strconv.FormatUint(uint64(controlID), 10),
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
		existingComment, err := caps.GetFullDetailsByID(controlID, actionPlanID, id)
		if err != nil {
			utility.LogError("Failed to fetch control action plan comment", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = caps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ControlActionPlanComment{}.TableName(),
			id,
			existingComment,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ControlActionPlanComment{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}

		if controlActionPlan.ResponsibleUserID == currentUserID {
			// if responsible user added the comment, notification for review user
			if controlActionPlan.ReviewUser != nil && controlActionPlan.ReviewUser.ID != 0{
				// prepare notification
				urlParam := fmt.Sprintf(`{"control_register_id":%d, "action_plan_id":%d}`, controlID, id)
				queryParam := ""
				entity := "control-action-plan"
				action := "comment-added"

				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserIDs:              []uint{uint(controlActionPlan.ReviewUser.ID)},
					ModuleIdentifier:     "controls",
					NotificationType:     string(commonDTOs.AssignmentType),
					HeadingEn:            "Comment Added For Control Action Plan",
					HeadingAr:            "تمت إضافة التعليق لخطة عمل الامتثال",
					TitleEn:              controlActionPlan.Title,
					TitleAr:              controlActionPlan.Title,
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					QueryParam:           queryParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				reviewUser, err := caps.userSvc.GetByID(uint(controlActionPlan.ReviewUser.ID))
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
				subject := fmt.Sprintf("Comment Added for Control Action Plan:%s", controlActionPlan.Title)
				appURL := utility.GetAppURL()
				emailData := map[string]interface{}{
					"FirstNameEn":            reviewUser.FirstNameEn,
					"LastNameEn":             reviewUser.LastNameEn,
					"ControlReferenceNumber": controlActionPlan.Control.ReferenceNumber,
					"ControlTitle":           controlActionPlan.Control.TitleEn,
					"ActionPlanTitle":        controlActionPlan.Title,
					"ActionPlanDescription":  controlActionPlan.Description,
					"CommentedByFirstNameEn": assignedBy.FirstNameEn,
					"CommentedByLastNameEn":  assignedBy.LastNameEn,
					"CommentedOn":            time.Now().Format("2006-01-02"),
					"RedirectURL":            fmt.Sprintf("%s/compliance/controls/%d/action-plans?action_plan_id=%d", appURL, controlID, id),
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              emailTo,
					Subject:              subject,
					TemplateRelativePath: "organization/templates/email/control_action_plan_comment.html",
					EmailData:            emailData,
				}

				// send notification and mail
				err = caps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "control_action_plan")
				if err != nil {
					utility.LogError("Failed to send notifications", map[string]interface{}{
						"error": err.Error(),
					})
				}
			}
		} else if uint(controlActionPlan.ReviewUser.ID) == currentUserID {
			// if review user added the comment, notification for responsible user
			// prepare notification
			urlParam := fmt.Sprintf(`{"control_register_id":%d, "action_plan_id":%d}`, controlID, id)
			queryParam := ""
			entity := "control-action-plan"
			action := "comment-added"

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{controlActionPlan.ResponsibleUserID},
				ModuleIdentifier:     "controls",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Comment Added For Control Action Plan",
				HeadingAr:            "تمت إضافة التعليق لخطة عمل الامتثال",
				TitleEn:              controlActionPlan.Title,
				TitleAr:              controlActionPlan.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			responsibleUser, err := caps.userSvc.GetByID(controlActionPlan.ResponsibleUserID)
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
			subject := fmt.Sprintf("Comment Added for Control Action Plan:%s", controlActionPlan.Title)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":            responsibleUser.FirstNameEn,
				"LastNameEn":             responsibleUser.LastNameEn,
				"ControlReferenceNumber": controlActionPlan.Control.ReferenceNumber,
				"ControlTitle":           controlActionPlan.Control.TitleEn,
				"ActionPlanTitle":        controlActionPlan.Title,
				"ActionPlanDescription":  controlActionPlan.Description,
				"CommentedByFirstNameEn": assignedBy.FirstNameEn,
				"CommentedByLastNameEn":  assignedBy.LastNameEn,
				"AssignedOn":             time.Now().Format("2006-01-02"),
				"RedirectURL":            fmt.Sprintf("%s/compliance/controls/%d/action-plans?action_plan_id=%d", appURL, controlID, id),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "organization/templates/email/control_action_plan_comment.html",
				EmailData:            emailData,
			}

			// send notification and mail
			err = caps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "control_action_plan")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		} else {
			// if other user added the comment, notification for review user and responsible user
			// prepare notification
			urlParam := fmt.Sprintf(`{"control_register_id":%d, "action_plan_id":%d}`, controlID, id)
			queryParam := ""
			entity := "control-action-plan"
			action := "comment-added"

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{controlActionPlan.ResponsibleUserID},
				ModuleIdentifier:     "controls",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Comment Added For Control Action Plan",
				HeadingAr:            "تمت إضافة التعليق لخطة عمل الامتثال",
				TitleEn:              controlActionPlan.Title,
				TitleAr:              controlActionPlan.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			responsibleUser, err := caps.userSvc.GetByID(controlActionPlan.ResponsibleUserID)
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
			subject := fmt.Sprintf("Comment Added for Control Action Plan:%s", controlActionPlan.Title)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":            responsibleUser.FirstNameEn,
				"LastNameEn":             responsibleUser.LastNameEn,
				"ControlReferenceNumber": controlActionPlan.Control.ReferenceNumber,
				"ControlTitle":           controlActionPlan.Control.TitleEn,
				"ActionPlanTitle":        controlActionPlan.Title,
				"ActionPlanDescription":  controlActionPlan.Description,
				"CommentedByFirstNameEn": assignedBy.FirstNameEn,
				"CommentedByLastNameEn":  assignedBy.LastNameEn,
				"AssignedOn":             time.Now().Format("2006-01-02"),
				"RedirectURL":            fmt.Sprintf("%s/compliance/controls/%d/action-plans?action_plan_id=%d", appURL, controlID, id),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "organization/templates/email/control_action_plan_comment.html",
				EmailData:            emailData,
			}

			// send notification and mail
			err = caps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "control_action_plan")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
			if controlActionPlan.ReviewUser.ID != 0 {
				// prepare notification
				urlParam := fmt.Sprintf(`{"control_register_id":%d, "action_plan_id":%d}`, controlID, id)
				queryParam := ""
				entity := "control-action-plan"
				action := "comment-added"

				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserIDs:              []uint{uint(controlActionPlan.ReviewUser.ID)},
					ModuleIdentifier:     "controls",
					NotificationType:     string(commonDTOs.AssignmentType),
					HeadingEn:            "Comment Added For Control Action Plan",
					HeadingAr:            "تمت إضافة التعليق لخطة عمل الامتثال",
					TitleEn:              controlActionPlan.Title,
					TitleAr:              controlActionPlan.Title,
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					QueryParam:           queryParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				responsibleUser, err := caps.userSvc.GetByID(uint(controlActionPlan.ReviewUser.ID))
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
				subject := fmt.Sprintf("Comment Added for Control Action Plan:%s", controlActionPlan.Title)
				appURL := utility.GetAppURL()
				emailData := map[string]interface{}{
					"FirstNameEn":            responsibleUser.FirstNameEn,
					"LastNameEn":             responsibleUser.LastNameEn,
					"ControlReferenceNumber": controlActionPlan.Control.ReferenceNumber,
					"ControlTitle":           controlActionPlan.Control.TitleEn,
					"ActionPlanTitle":        controlActionPlan.Title,
					"ActionPlanDescription":  controlActionPlan.Description,
					"CommentedByFirstNameEn": assignedBy.FirstNameEn,
					"CommentedByLastNameEn":  assignedBy.LastNameEn,
					"AssignedOn":             time.Now().Format("2006-01-02"),
					"RedirectURL":            fmt.Sprintf("%s/compliance/controls/%d/action-plans?action_plan_id=%d", appURL, controlID, id),
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              emailTo,
					Subject:              subject,
					TemplateRelativePath: "organization/templates/email/control_action_plan_comment.html",
					EmailData:            emailData,
				}

				// send notification and mail
				err = caps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "control_action_plan")
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

func (caps *ControlActionPlanCommentService) Update(controlID uint, actionPlanID uint, id uint, currentUserID uint, controlActionPlanCommentFormDTO *controlDTOs.ControlActionPlanCommentFormDTO) error {
	tx := caps.repo.BeginTransaction()

	if _, err := caps.controlActionPlanSvc.GetByID(controlID, actionPlanID); err != nil {
		return err
	}

	existingComment, err := caps.GetByID(controlID, actionPlanID, id)
	if err != nil {
		return err
	}

	oldComment, err := caps.GetFullDetailsByID(controlID, actionPlanID, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingComment.Comment = controlActionPlanCommentFormDTO.Comment
	err = caps.repo.Update(existingComment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	var newDocuments []*models.ControlActionPlanCommentDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.ControlActionPlanCommentDocument

	for _, doc := range controlActionPlanCommentFormDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {

			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}

			url := utility.GenerateDocumentURL("control", "control-register", strconv.FormatUint(uint64(controlID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

			document := &models.ControlActionPlanCommentDocument{
				ControlActionPlanCommentID: id,
				URL:                        url,
				Name:                       doc.FileName,
				Size:                       doc.Size,
				Token:                      encryptedFileName,
				OriginalName:               doc.OriginalName,
				Type:                       doc.Type,
				CreatedBYID:                currentUserID,
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
		storagePath := filepath.Join("control", "control-register", strconv.FormatUint(uint64(controlID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("control", "control-register", strconv.FormatUint(uint64(controlID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	go func() {
		existingComment, err := caps.GetFullDetailsByID(controlID, actionPlanID, id)
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
				models.ControlActionPlanComment{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.ControlActionPlanComment{}.TableName(),
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

func (caps *ControlActionPlanCommentService) Delete(controlID uint, actionPlanID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := caps.repo.BeginTransaction()
	if _, err := caps.controlActionPlanSvc.GetByID(controlID, actionPlanID); err != nil {
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
			"id":             id,
			"control_action": actionPlanID,
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
				models.ControlActionPlanComment{}.TableName(),
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

func (caps *ControlActionPlanCommentService) GetDocumentPath(controlID, actionPlanID, id uint, token string) (string, error) {
	existingComment, err := caps.GetByID(controlID, actionPlanID, id)
	if err != nil {
		return "", err
	}
	if len(existingComment.Documents) == 0 {
		return "", fmt.Errorf("no documents found for the control action plan comment")
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	storagePath := filepath.Join("control", "control-register", strconv.FormatUint(uint64(controlID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}
