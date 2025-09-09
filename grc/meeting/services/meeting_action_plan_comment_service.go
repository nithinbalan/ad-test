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
	meetingDTOs "jethur/grc/meeting/dtos"
	"jethur/grc/meeting/repositories"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type MeetingActionPlanCommentServiceInterface interface {
	GetAll(meetingID uint, actionPlanID uint) ([]*meetingDTOs.MeetingActionPlanCommentListDTO, error)
	GetByID(meetingID uint, actionPlanID uint, id uint) (*models.MeetingActionPlanComment, error)
	GetFullDetailsByID(meetingID uint, actionPlanID uint, id uint, tx ...utility.DatabaseTransaction) (*meetingDTOs.MeetingActionPlanCommentDetailDTO, error)
	Create(meetingID uint, actionPlanID uint, currentUserID uint, dto *meetingDTOs.MeetingActionPlanCommentFormDTO) (uint, error)
	Update(meetingID uint, actionPlanID uint, id uint, currentUserID uint, dto *meetingDTOs.MeetingActionPlanCommentFormDTO) error
	Delete(meetingID uint, actionPlanID uint, ids []uint, currentUserID uint) ([]uint, error)

	GetDocumentPath(meetingID, actionPlanID, id uint, token string) (string, error)
}

type MeetingActionPlanCommentService struct {
	repo                 repositories.MeetingActionPlanCommentRepositoryInterface
	userSvc              organizationServices.UserServiceInterface
	activityLogSvc       logServices.ActivityLogServiceInterface
	meetingActionPlanSvc MeetingActionPlanServiceInterface
	notificationSvc      commonServices.NotificationServiceInterface
}

func NewMeetingActionPlanCommentService(repo repositories.MeetingActionPlanCommentRepositoryInterface, userSvc organizationServices.UserServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	meetingActionPlanSvc MeetingActionPlanServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) MeetingActionPlanCommentServiceInterface {
	return &MeetingActionPlanCommentService{
		repo:                 repo,
		userSvc:              userSvc,
		activityLogSvc:       activityLogSvc,
		meetingActionPlanSvc: meetingActionPlanSvc,
		notificationSvc:      notificationSvc,
	}
}

func (mapc *MeetingActionPlanCommentService) GetAll(meetingID uint, actionPlanID uint) ([]*meetingDTOs.MeetingActionPlanCommentListDTO, error) {
	if _, err := mapc.meetingActionPlanSvc.GetByID(meetingID, actionPlanID); err != nil {
		return nil, err
	}

	return mapc.repo.GetAll(actionPlanID)
}

func (mapc *MeetingActionPlanCommentService) GetByID(meetingID uint, actionPlanID uint, id uint) (*models.MeetingActionPlanComment, error) {
	if _, err := mapc.meetingActionPlanSvc.GetByID(meetingID, actionPlanID); err != nil {
		return nil, err
	}

	comment, err := mapc.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return comment, nil
}

func (mapc *MeetingActionPlanCommentService) GetFullDetailsByID(meetingID uint, actionPlanID uint, id uint, tx ...utility.DatabaseTransaction) (*meetingDTOs.MeetingActionPlanCommentDetailDTO, error) {
	if _, err := mapc.meetingActionPlanSvc.GetByID(meetingID, actionPlanID); err != nil {
		fmt.Printf("MeetingActionPlan not found for meetingID %d, actionPlanID %d\n", meetingID, actionPlanID)
		return nil, err
	}

	var commentDetail meetingDTOs.MeetingActionPlanCommentDetailDTO

	comment, err := mapc.repo.GetFullDetailsByID(id, tx...)
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

func (mapc *MeetingActionPlanCommentService) Create(meetingID uint, actionPlanID uint, currentUserID uint, meetingActionPlanCommentFormDTO *meetingDTOs.MeetingActionPlanCommentFormDTO) (uint, error) {
	tx := mapc.repo.BeginTransaction()

	if _, err := mapc.meetingActionPlanSvc.GetByID(meetingID, actionPlanID); err != nil {
		return 0, err
	}

	meetingActionPlan, err := mapc.meetingActionPlanSvc.GetFullDetailsByID(meetingID, actionPlanID)
	if err != nil {
		return 0, err
	}

	comment := &models.MeetingActionPlanComment{
		MeetingActionPlanID: actionPlanID,
		Comment:             meetingActionPlanCommentFormDTO.Comment,
		CreatedBYID:         currentUserID,
	}

	id, err := mapc.repo.Create(comment)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.MeetingActionPlanCommentDocument

	for _, doc := range meetingActionPlanCommentFormDTO.Documents {
		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL(
			"meeting", "meetings", strconv.FormatUint(uint64(meetingID), 10),
			"action-plans", strconv.FormatUint(uint64(actionPlanID), 10),
			"comments", strconv.FormatUint(uint64(id), 10),
			"documents",
		)

		document := &models.MeetingActionPlanCommentDocument{
			MeetingActionPlanCommentID: id,
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
		if err := mapc.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}

		for _, doc := range documents {
			storagePath := filepath.Join(
				"meeting", "meetings", strconv.FormatUint(uint64(meetingID), 10),
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
		existingComment, err := mapc.GetFullDetailsByID(meetingID, actionPlanID, id)
		if err != nil {
			utility.LogError("Failed to fetch meeting action plan comment", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = mapc.activityLogSvc.LogActivity(
			models.CreateAction,
			models.MeetingActionPlanComment{}.TableName(),
			id,
			existingComment,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.MeetingActionPlanComment{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}

		if meetingActionPlan.ResponsibleUserID == currentUserID {
			// if responsible user added the comment, notification for review user
			if meetingActionPlan.ReviewUser.ID != 0 {
				// prepare notification
				urlParam := fmt.Sprintf(`{"meeting_id":%d, "action_plan_id":%d}`, meetingID, id)
				queryParam := ""
				entity := "meeting-action-plan"
				action := "comment-added"

				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserIDs:              []uint{uint(meetingActionPlan.ReviewUser.ID)},
					ModuleIdentifier:     "meeting",
					NotificationType:     string(commonDTOs.AssignmentType),
					HeadingEn:            "Comment Added For Meeting Action Plan",
					HeadingAr:            "تمت إضافة التعليق لخطة عمل الامتثال",
					TitleEn:              meetingActionPlan.Title,
					TitleAr:              meetingActionPlan.Title,
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					QueryParam:           queryParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				reviewUser, err := mapc.userSvc.GetByID(uint(meetingActionPlan.ReviewUser.ID))
				if err != nil {
					utility.LogError("Failed to fetch review user", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}
				assignedBy, err := mapc.userSvc.GetByID(currentUserID)
				if err != nil {
					utility.LogError("Failed to fetch commented user", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}

				emailTo := []string{reviewUser.Email}
				subject := fmt.Sprintf("Comment Added for Meeting Action Plan:%s", meetingActionPlan.Title)
				appURL := utility.GetAppURL()
				emailData := map[string]interface{}{
					"FirstNameEn":            reviewUser.FirstNameEn,
					"LastNameEn":             reviewUser.LastNameEn,
					"MeetingReferenceNumber": meetingActionPlan.Meeting.ReferenceNumber,
					"MeetingTitle":           meetingActionPlan.Meeting.Title,
					"ActionPlanTitle":        meetingActionPlan.Title,
					"ActionPlanDescription":  meetingActionPlan.Description,
					"CommentedByFirstNameEn": assignedBy.FirstNameEn,
					"CommentedByLastNameEn":  assignedBy.LastNameEn,
					"CommentedOn":            time.Now().Format("2006-01-02"),
					"RedirectURL":            fmt.Sprintf("%s/management-review/meetings/%d/action-plans?action_plan_id=%d", appURL, meetingID, id),
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              emailTo,
					Subject:              subject,
					TemplateRelativePath: "meeting/templates/email/meeting_action_plan_comment.html",
					EmailData:            emailData,
				}

				// send notification and mail
				err = mapc.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "meeting_action_plan")
				if err != nil {
					utility.LogError("Failed to send notifications", map[string]interface{}{
						"error": err.Error(),
					})
				}
			}
		} else if uint(meetingActionPlan.ReviewUser.ID) == currentUserID {
			// if review user added the comment, notification for responsible user
			// prepare notification
			urlParam := fmt.Sprintf(`{"meeting_id":%d, "action_plan_id":%d}`, meetingID, id)
			queryParam := ""
			entity := "meeting-action-plan"
			action := "comment-added"

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{meetingActionPlan.ResponsibleUserID},
				ModuleIdentifier:     "meeting",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Comment Added For Meeting Action Plan",
				HeadingAr:            "تمت إضافة التعليق لخطة عمل الامتثال",
				TitleEn:              meetingActionPlan.Title,
				TitleAr:              meetingActionPlan.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			responsibleUser, err := mapc.userSvc.GetByID(meetingActionPlan.ResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}
			assignedBy, err := mapc.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch commented user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			emailTo := []string{responsibleUser.Email}
			subject := fmt.Sprintf("Comment Added for Meeting Action Plan:%s", meetingActionPlan.Title)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":            responsibleUser.FirstNameEn,
				"LastNameEn":             responsibleUser.LastNameEn,
				"MeetingReferenceNumber": meetingActionPlan.Meeting.ReferenceNumber,
				"MeetingTitle":           meetingActionPlan.Meeting.Title,
				"ActionPlanTitle":        meetingActionPlan.Title,
				"ActionPlanDescription":  meetingActionPlan.Description,
				"CommentedByFirstNameEn": assignedBy.FirstNameEn,
				"CommentedByLastNameEn":  assignedBy.LastNameEn,
				"AssignedOn":             time.Now().Format("2006-01-02"),
				"RedirectURL":            fmt.Sprintf("%s/management-review/meetings/%d/action-plans?action_plan_id=%d", appURL, meetingID, id),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "meeting/templates/email/meeting_action_plan_comment.html",
				EmailData:            emailData,
			}

			// send notification and mail
			err = mapc.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "meeting_action_plan")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		} else {
			// if other user added the comment, notification for review user and responsible user
			// prepare notification
			urlParam := fmt.Sprintf(`{"meeting_id":%d, "action_plan_id":%d}`, meetingID, id)
			queryParam := ""
			entity := "meeting-action-plan"
			action := "comment-added"

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{meetingActionPlan.ResponsibleUserID},
				ModuleIdentifier:     "meeting",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Comment Added For Meeting Action Plan",
				HeadingAr:            "تمت إضافة التعليق لخطة عمل الامتثال",
				TitleEn:              meetingActionPlan.Title,
				TitleAr:              meetingActionPlan.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			responsibleUser, err := mapc.userSvc.GetByID(meetingActionPlan.ResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}
			assignedBy, err := mapc.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch commented user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			emailTo := []string{responsibleUser.Email}
			subject := fmt.Sprintf("Comment Added for Meeting Action Plan:%s", meetingActionPlan.Title)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":            responsibleUser.FirstNameEn,
				"LastNameEn":             responsibleUser.LastNameEn,
				"MeetingReferenceNumber": meetingActionPlan.Meeting.ReferenceNumber,
				"MeetingTitle":           meetingActionPlan.Meeting.Title,
				"ActionPlanTitle":        meetingActionPlan.Title,
				"ActionPlanDescription":  meetingActionPlan.Description,
				"CommentedByFirstNameEn": assignedBy.FirstNameEn,
				"CommentedByLastNameEn":  assignedBy.LastNameEn,
				"AssignedOn":             time.Now().Format("2006-01-02"),
				"RedirectURL":            fmt.Sprintf("%s/management-review/meetings/%d/action-plans?action_plan_id=%d", appURL, meetingID, id),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "meeting/templates/email/meeting_action_plan_comment.html",
				EmailData:            emailData,
			}

			// send notification and mail
			err = mapc.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "meeting_action_plan")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
			if meetingActionPlan.ReviewUser.ID != 0 {
				// prepare notification
				urlParam := fmt.Sprintf(`{"meeting_id":%d, "action_plan_id":%d}`, meetingID, id)
				queryParam := ""
				entity := "meeting-action-plan"
				action := "comment-added"

				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserIDs:              []uint{uint(meetingActionPlan.ReviewUser.ID)},
					ModuleIdentifier:     "meeting",
					NotificationType:     string(commonDTOs.AssignmentType),
					HeadingEn:            "Comment Added For Meeting Action Plan",
					HeadingAr:            "تمت إضافة التعليق لخطة عمل الامتثال",
					TitleEn:              meetingActionPlan.Title,
					TitleAr:              meetingActionPlan.Title,
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					QueryParam:           queryParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				responsibleUser, err := mapc.userSvc.GetByID(uint(meetingActionPlan.ReviewUser.ID))
				if err != nil {
					utility.LogError("Failed to fetch review user", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}
				assignedBy, err := mapc.userSvc.GetByID(currentUserID)
				if err != nil {
					utility.LogError("Failed to fetch commented user", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}

				emailTo := []string{responsibleUser.Email}
				subject := fmt.Sprintf("Comment Added for Meeting Action Plan:%s", meetingActionPlan.Title)
				appURL := utility.GetAppURL()
				emailData := map[string]interface{}{
					"FirstNameEn":            responsibleUser.FirstNameEn,
					"LastNameEn":             responsibleUser.LastNameEn,
					"MeetingReferenceNumber": meetingActionPlan.Meeting.ReferenceNumber,
					"MeetingTitle":           meetingActionPlan.Meeting.Title,
					"ActionPlanTitle":        meetingActionPlan.Title,
					"ActionPlanDescription":  meetingActionPlan.Description,
					"CommentedByFirstNameEn": assignedBy.FirstNameEn,
					"CommentedByLastNameEn":  assignedBy.LastNameEn,
					"AssignedOn":             time.Now().Format("2006-01-02"),
					"RedirectURL":            fmt.Sprintf("%s/management-review/meetings/%d/action-plans?action_plan_id=%d", appURL, meetingID, id),
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              emailTo,
					Subject:              subject,
					TemplateRelativePath: "meeting/templates/email/meeting_action_plan_comment.html",
					EmailData:            emailData,
				}

				// send notification and mail
				err = mapc.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "meeting_action_plan")
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

func (mapc *MeetingActionPlanCommentService) Update(meetingID uint, actionPlanID uint, id uint, currentUserID uint, meetingActionPlanCommentFormDTO *meetingDTOs.MeetingActionPlanCommentFormDTO) error {
	tx := mapc.repo.BeginTransaction()

	if _, err := mapc.meetingActionPlanSvc.GetByID(meetingID, actionPlanID); err != nil {
		return err
	}

	existingComment, err := mapc.GetByID(meetingID, actionPlanID, id)
	if err != nil {
		return err
	}

	oldComment, err := mapc.GetFullDetailsByID(meetingID, actionPlanID, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingComment.Comment = meetingActionPlanCommentFormDTO.Comment
	err = mapc.repo.Update(existingComment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	var newDocuments []*models.MeetingActionPlanCommentDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.MeetingActionPlanCommentDocument

	for _, doc := range meetingActionPlanCommentFormDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {

			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}

			url := utility.GenerateDocumentURL("meeting", "meetings", strconv.FormatUint(uint64(meetingID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

			document := &models.MeetingActionPlanCommentDocument{
				MeetingActionPlanCommentID: id,
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
				existingDocument, err := mapc.repo.GetDocumentByID(*doc.ID)
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
		if err := mapc.repo.CreateDocuments(newDocuments, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIds) > 0 {
		if err := mapc.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range newDocuments {
		storagePath := filepath.Join("meeting", "meetings", strconv.FormatUint(uint64(meetingID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("meeting", "meetings", strconv.FormatUint(uint64(meetingID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	go func() {
		existingComment, err := mapc.GetFullDetailsByID(meetingID, actionPlanID, id)
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
			err := mapc.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.MeetingActionPlanComment{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.MeetingActionPlanComment{}.TableName(),
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

func (mapc *MeetingActionPlanCommentService) Delete(meetingID uint, actionPlanID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := mapc.repo.BeginTransaction()
	if _, err := mapc.meetingActionPlanSvc.GetByID(meetingID, actionPlanID); err != nil {
		return nil, err
	}

	var failedIDs []uint
	var successfulIDs []uint
	var deletedComments []map[string]interface{}
	for _, id := range ids {
		if _, err := mapc.repo.GetFullDetailsByID(id); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		deletedComments = append(deletedComments, map[string]interface{}{
			"id":             id,
			"meeting_action": actionPlanID,
		})

		err := mapc.repo.Delete(id, tx)
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
			err := mapc.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.MeetingActionPlanComment{}.TableName(),
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

func (mapc *MeetingActionPlanCommentService) GetDocumentPath(meetingID, actionPlanID, id uint, token string) (string, error) {
	existingComment, err := mapc.GetByID(meetingID, actionPlanID, id)
	if err != nil {
		return "", err
	}
	if len(existingComment.Documents) == 0 {
		return "", fmt.Errorf("no documents found for the meeting action plan comment")
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	storagePath := filepath.Join("meeting", "meetings", strconv.FormatUint(uint64(meetingID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}
