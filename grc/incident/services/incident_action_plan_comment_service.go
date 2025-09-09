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
	incidentDTOs "jethur/grc/incident/dtos"
	"jethur/grc/incident/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type IncidentActionPlanCommentServiceInterface interface {
	GetAll(incidentID uint, actionPlanID uint) ([]*incidentDTOs.IncidentActionPlanCommentListDTO, error)
	GetByID(incidentID uint, actionPlanID uint, id uint) (*models.IncidentActionPlanComment, error)
	GetFullDetailsByID(incidentID uint, actionPlanID uint, id uint, tx ...utility.DatabaseTransaction) (*incidentDTOs.IncidentActionPlanCommentDetailDTO, error)
	Create(incidentID uint, actionPlanID uint, currentUserID uint, dto *incidentDTOs.IncidentActionPlanCommentFormDTO) (uint, error)
	Update(incidentID uint, actionPlanID uint, id uint, currentUserID uint, dto *incidentDTOs.IncidentActionPlanCommentFormDTO) error
	Delete(incidentID uint, actionPlanID uint, ids []uint, currentUserID uint) ([]uint, error)

	GetDocumentPath(incidentID, actionPlanID, id uint, token string) (string, error)
}

type IncidentActionPlanCommentService struct {
	repo                  repositories.IncidentActionPlanCommentRepositoryInterface
	userSvc               organizationServices.UserServiceInterface
	activityLogSvc        logServices.ActivityLogServiceInterface
	incidentActionPlanSvc IncidentActionPlanServiceInterface
	notificationSvc       commonServices.NotificationServiceInterface
}

func NewIncidentActionPlanCommentService(repo repositories.IncidentActionPlanCommentRepositoryInterface, userSvc organizationServices.UserServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	incidentActionPlanSvc IncidentActionPlanServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) IncidentActionPlanCommentServiceInterface {
	return &IncidentActionPlanCommentService{
		repo:                  repo,
		userSvc:               userSvc,
		activityLogSvc:        activityLogSvc,
		incidentActionPlanSvc: incidentActionPlanSvc,
		notificationSvc:       notificationSvc,
	}
}

func (iapc *IncidentActionPlanCommentService) GetAll(incidentID uint, actionPlanID uint) ([]*incidentDTOs.IncidentActionPlanCommentListDTO, error) {
	if _, err := iapc.incidentActionPlanSvc.GetByID(incidentID, actionPlanID); err != nil {
		return nil, err
	}

	return iapc.repo.GetAll(actionPlanID)
}

func (iapc *IncidentActionPlanCommentService) GetByID(incidentID uint, actionPlanID uint, id uint) (*models.IncidentActionPlanComment, error) {
	if _, err := iapc.incidentActionPlanSvc.GetByID(incidentID, actionPlanID); err != nil {
		return nil, err
	}

	comment, err := iapc.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return comment, nil
}

func (iapc *IncidentActionPlanCommentService) GetFullDetailsByID(incidentID uint, actionPlanID uint, id uint, tx ...utility.DatabaseTransaction) (*incidentDTOs.IncidentActionPlanCommentDetailDTO, error) {
	if _, err := iapc.incidentActionPlanSvc.GetByID(incidentID, actionPlanID); err != nil {
		fmt.Printf("IncidentActionPlan not found for incidentID %d, actionPlanID %d\n", incidentID, actionPlanID)
		return nil, err
	}

	var commentDetail incidentDTOs.IncidentActionPlanCommentDetailDTO

	comment, err := iapc.repo.GetFullDetailsByID(id, tx...)
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

func (iapc *IncidentActionPlanCommentService) Create(incidentID uint, actionPlanID uint, currentUserID uint, incidentActionPlanCommentFormDTO *incidentDTOs.IncidentActionPlanCommentFormDTO) (uint, error) {
	tx := iapc.repo.BeginTransaction()

	if _, err := iapc.incidentActionPlanSvc.GetByID(incidentID, actionPlanID); err != nil {
		return 0, err
	}

	incidentActionPlan, err := iapc.incidentActionPlanSvc.GetFullDetailsByID(incidentID, actionPlanID)

	if err != nil {
		return 0, err
	}

	comment := &models.IncidentActionPlanComment{
		IncidentActionPlanID: actionPlanID,
		Comment:              incidentActionPlanCommentFormDTO.Comment,
		CreatedBYID:          currentUserID,
	}

	id, err := iapc.repo.Create(comment)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.IncidentActionPlanCommentDocument

	for _, doc := range incidentActionPlanCommentFormDTO.Documents {
		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL(
			"incident", "incident-register", strconv.FormatUint(uint64(incidentID), 10),
			"action-plans", strconv.FormatUint(uint64(actionPlanID), 10),
			"comments", strconv.FormatUint(uint64(id), 10),
			"documents",
		)

		document := &models.IncidentActionPlanCommentDocument{
			IncidentActionPlanCommentID: id,
			URL:                         url,
			Name:                        doc.FileName,
			Size:                        doc.Size,
			Token:                       encryptedFilename,
			OriginalName:                doc.OriginalName,
			Type:                        doc.Type,
			CreatedBYID:                 currentUserID,
		}

		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := iapc.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}

		for _, doc := range documents {
			storagePath := filepath.Join(
				"incident", "incident-register", strconv.FormatUint(uint64(incidentID), 10),
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
		existingComment, err := iapc.GetFullDetailsByID(incidentID, actionPlanID, id)
		if err != nil {
			utility.LogError("Failed to fetch incident action plan comment", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = iapc.activityLogSvc.LogActivity(
			models.CreateAction,
			models.IncidentActionPlanComment{}.TableName(),
			id,
			existingComment,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.IncidentActionPlanComment{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
		if uint(incidentActionPlan.ResponsibleUser.ID) == currentUserID {
			// if responsible user added the comment, notification for review user
			if incidentActionPlan.ReviewUser.ID != 0 {

				// prepare notification
				urlParam := fmt.Sprintf(`{"incident_id":%d, "action_plan_id":%d}`, incidentID, id)
				queryParam := ""
				entity := "incident-action-plan"
				action := "comment-added"

				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserIDs:              []uint{uint(incidentActionPlan.ReviewUser.ID)},
					ModuleIdentifier:     "incident",
					NotificationType:     string(commonDTOs.AssignmentType),
					HeadingEn:            "Comment Added For Incident Action Plan",
					HeadingAr:            "تمت إضافة التعليق لخطة عمل الامتثال",
					TitleEn:              incidentActionPlan.Title,
					TitleAr:              incidentActionPlan.Title,
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					QueryParam:           queryParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				reviewUser, err := iapc.userSvc.GetByID(uint(incidentActionPlan.ReviewUser.ID))
				if err != nil {
					utility.LogError("Failed to fetch review user", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}
				assignedBy, err := iapc.userSvc.GetByID(currentUserID)
				if err != nil {
					utility.LogError("Failed to fetch commented user", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}

				emailTo := []string{reviewUser.Email}
				subject := fmt.Sprintf("Comment Added for Incident Action Plan:%s", incidentActionPlan.Title)
				appURL := utility.GetAppURL()
				emailData := map[string]interface{}{
					"FirstNameEn":             reviewUser.FirstNameEn,
					"LastNameEn":              reviewUser.LastNameEn,
					"IncidentTitle":           incidentActionPlan.IncidentRegister.Title,
					"IncidentReferenceNumber": incidentActionPlan.IncidentRegister.ReferenceNumber,
					"ActionPlanTitle":         incidentActionPlan.Title,
					"ActionPlanDescription":   incidentActionPlan.Description,
					"CommentedByFirstNameEn":  assignedBy.FirstNameEn,
					"CommentedByLastNameEn":   assignedBy.LastNameEn,
					"CommentedOn":             time.Now().Format("2006-01-02"),
					"RedirectURL":             fmt.Sprintf("%s/incident/register/%d/action-plans?action_plan_id=%d", appURL, incidentID, id),
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              emailTo,
					Subject:              subject,
					TemplateRelativePath: "incident/templates/email/incident_action_plan_comment.html",
					EmailData:            emailData,
				}

				// send notification and mail
				err = iapc.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "incident_action_plan")
				if err != nil {
					utility.LogError("Failed to send notifications", map[string]interface{}{
						"error": err.Error(),
					})
				}
			}
		} else if uint(incidentActionPlan.ReviewUser.ID) == currentUserID {
			// if review user added the comment, notification for responsible user
			// prepare notification
			urlParam := fmt.Sprintf(`{"incident_id":%d, "action_plan_id":%d}`, incidentID, id)
			queryParam := ""
			entity := "incident-action-plan"
			action := "comment-added"

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{uint(incidentActionPlan.ResponsibleUser.ID)},
				ModuleIdentifier:     "incident",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Comment Added For Incident Action Plan",
				HeadingAr:            "تمت إضافة التعليق لخطة عمل الامتثال",
				TitleEn:              incidentActionPlan.Title,
				TitleAr:              incidentActionPlan.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			responsibleUser, err := iapc.userSvc.GetByID(uint(incidentActionPlan.ResponsibleUser.ID))
			if err != nil {
				utility.LogError("Failed to fetch responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}
			assignedBy, err := iapc.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch commented user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			emailTo := []string{responsibleUser.Email}
			subject := fmt.Sprintf("Comment Added for Incident Action Plan:%s", incidentActionPlan.Title)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":             responsibleUser.FirstNameEn,
				"LastNameEn":              responsibleUser.LastNameEn,
				"IncidentTitle":           incidentActionPlan.IncidentRegister.Title,
				"IncidentReferenceNumber": incidentActionPlan.IncidentRegister.ReferenceNumber,
				"ActionPlanTitle":         incidentActionPlan.Title,
				"ActionPlanDescription":   incidentActionPlan.Description,
				"CommentedByFirstNameEn":  assignedBy.FirstNameEn,
				"CommentedByLastNameEn":   assignedBy.LastNameEn,
				"AssignedOn":              time.Now().Format("2006-01-02"),
				"RedirectURL":             fmt.Sprintf("%s/incident/register/%d/action-plans?action_plan_id=%d", appURL, incidentID, id),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "incident/templates/email/incident_action_plan_comment.html",
				EmailData:            emailData,
			}

			// send notification and mail
			err = iapc.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "incident_action_plan")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		} else {
			// if other user added the comment, notification for review user and responsible user
			// prepare notification
			urlParam := fmt.Sprintf(`{"incident_id":%d, "action_plan_id":%d}`, incidentID, id)
			queryParam := ""
			entity := "incident-action-plan"
			action := "comment-added"

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{uint(incidentActionPlan.ResponsibleUser.ID)},
				ModuleIdentifier:     "incident",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Comment Added For Incident Action Plan",
				HeadingAr:            "تمت إضافة التعليق لخطة عمل الامتثال",
				TitleEn:              incidentActionPlan.Title,
				TitleAr:              incidentActionPlan.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			responsibleUser, err := iapc.userSvc.GetByID(uint(incidentActionPlan.ResponsibleUser.ID))
			if err != nil {
				utility.LogError("Failed to fetch responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}
			assignedBy, err := iapc.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch commented user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			emailTo := []string{responsibleUser.Email}
			subject := fmt.Sprintf("Comment Added for Incident Action Plan:%s", incidentActionPlan.Title)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":             responsibleUser.FirstNameEn,
				"LastNameEn":              responsibleUser.LastNameEn,
				"IncidentTitle":           incidentActionPlan.IncidentRegister.Title,
				"IncidentReferenceNumber": incidentActionPlan.IncidentRegister.ReferenceNumber,

				"ActionPlanTitle":        incidentActionPlan.Title,
				"ActionPlanDescription":  incidentActionPlan.Description,
				"CommentedByFirstNameEn": assignedBy.FirstNameEn,
				"CommentedByLastNameEn":  assignedBy.LastNameEn,
				"AssignedOn":             time.Now().Format("2006-01-02"),
				"RedirectURL":            fmt.Sprintf("%s/incident/register/%d/action-plans?action_plan_id=%d", appURL, incidentID, id),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "incident/templates/email/incident_action_plan_comment.html",
				EmailData:            emailData,
			}

			// send notification and mail
			err = iapc.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "incident_action_plan")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
			if incidentActionPlan.ReviewUser.ID != 0 {
				// prepare notification
				urlParam := fmt.Sprintf(`{"incident_id":%d, "action_plan_id":%d}`, incidentID, id)
				queryParam := ""
				entity := "incident-action-plan"
				action := "comment-added"

				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserIDs:              []uint{uint(incidentActionPlan.ReviewUser.ID)},
					ModuleIdentifier:     "incident",
					NotificationType:     string(commonDTOs.AssignmentType),
					HeadingEn:            "Comment Added For Incident Action Plan",
					HeadingAr:            "تمت إضافة التعليق لخطة عمل الامتثال",
					TitleEn:              incidentActionPlan.Title,
					TitleAr:              incidentActionPlan.Title,
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					QueryParam:           queryParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				responsibleUser, err := iapc.userSvc.GetByID(uint(incidentActionPlan.ReviewUser.ID))
				if err != nil {
					utility.LogError("Failed to fetch review user", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}
				assignedBy, err := iapc.userSvc.GetByID(currentUserID)
				if err != nil {
					utility.LogError("Failed to fetch commented user", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}

				emailTo := []string{responsibleUser.Email}
				subject := fmt.Sprintf("Comment Added for Incident Action Plan:%s", incidentActionPlan.Title)
				appURL := utility.GetAppURL()
				emailData := map[string]interface{}{
					"FirstNameEn":             responsibleUser.FirstNameEn,
					"LastNameEn":              responsibleUser.LastNameEn,
					"IncidentTitle":           incidentActionPlan.IncidentRegister.Title,
					"IncidentReferenceNumber": incidentActionPlan.IncidentRegister.ReferenceNumber,

					"ActionPlanTitle":        incidentActionPlan.Title,
					"ActionPlanDescription":  incidentActionPlan.Description,
					"CommentedByFirstNameEn": assignedBy.FirstNameEn,
					"CommentedByLastNameEn":  assignedBy.LastNameEn,
					"AssignedOn":             time.Now().Format("2006-01-02"),
					"RedirectURL":            fmt.Sprintf("%s/incident/register/%d/action-plans?action_plan_id=%d", appURL, incidentID, id),
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              emailTo,
					Subject:              subject,
					TemplateRelativePath: "incident/templates/email/incident_action_plan_comment.html",
					EmailData:            emailData,
				}

				// send notification and mail
				err = iapc.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "incident_action_plan")
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
func (iapc *IncidentActionPlanCommentService) Update(incidentID uint, actionPlanID uint, id uint, currentUserID uint, incidentActionPlanCommentFormDTO *incidentDTOs.IncidentActionPlanCommentFormDTO) error {
	tx := iapc.repo.BeginTransaction()

	if _, err := iapc.incidentActionPlanSvc.GetByID(incidentID, actionPlanID); err != nil {
		return err
	}

	existingComment, err := iapc.GetByID(incidentID, actionPlanID, id)
	if err != nil {
		return err
	}

	oldComment, err := iapc.GetFullDetailsByID(incidentID, actionPlanID, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingComment.Comment = incidentActionPlanCommentFormDTO.Comment
	err = iapc.repo.Update(existingComment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	var newDocuments []*models.IncidentActionPlanCommentDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.IncidentActionPlanCommentDocument

	for _, doc := range incidentActionPlanCommentFormDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {

			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}

			url := utility.GenerateDocumentURL("incident", "incident-register", strconv.FormatUint(uint64(incidentID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

			document := &models.IncidentActionPlanCommentDocument{
				IncidentActionPlanCommentID: id,
				URL:                         url,
				Name:                        doc.FileName,
				Size:                        doc.Size,
				Token:                       encryptedFileName,
				OriginalName:                doc.OriginalName,
				Type:                        doc.Type,
				CreatedBYID:                 currentUserID,
			}

			newDocuments = append(newDocuments, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			if doc.ID != nil {
				existingDocument, err := iapc.repo.GetDocumentByID(*doc.ID)
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
		if err := iapc.repo.CreateDocuments(newDocuments, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIds) > 0 {
		if err := iapc.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range newDocuments {
		storagePath := filepath.Join("incident", "incident-register", strconv.FormatUint(uint64(incidentID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("incident", "incident-register", strconv.FormatUint(uint64(incidentID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	go func() {
		existingComment, err := iapc.GetFullDetailsByID(incidentID, actionPlanID, id)
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
			err := iapc.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.IncidentActionPlanComment{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.IncidentActionPlanComment{}.TableName(),
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

func (iapc *IncidentActionPlanCommentService) Delete(incidentID uint, actionPlanID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := iapc.repo.BeginTransaction()
	if _, err := iapc.incidentActionPlanSvc.GetByID(incidentID, actionPlanID); err != nil {
		return nil, err
	}

	var failedIDs []uint
	var successfulIDs []uint
	var deletedComments []map[string]interface{}
	for _, id := range ids {
		if _, err := iapc.repo.GetFullDetailsByID(id); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		deletedComments = append(deletedComments, map[string]interface{}{
			"id":              id,
			"incident_action": actionPlanID,
		})

		err := iapc.repo.Delete(id, tx)
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
			err := iapc.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.IncidentActionPlanComment{}.TableName(),
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

func (iapc *IncidentActionPlanCommentService) GetDocumentPath(incidentID, actionPlanID, id uint, token string) (string, error) {
	existingComment, err := iapc.GetByID(incidentID, actionPlanID, id)
	if err != nil {
		return "", err
	}
	if len(existingComment.Documents) == 0 {
		return "", fmt.Errorf("no documents found for the incident action plan comment")
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	storagePath := filepath.Join("incident", "incident-register", strconv.FormatUint(uint64(incidentID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}
