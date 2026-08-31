package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"jethur/config"
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type AuditFindingActionPlanCommentServiceInterface interface {
	GetAll(findingID uint, actionPlanID uint) ([]*auditDTOs.AuditFindingActionPlanCommentListDTO, error)
	GetByID(findingID uint, actionPlanID uint, id uint) (*models.AuditFindingActionPlanComment, error)
	GetFullDetailsByID(findingID uint, actionPlanID uint, id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditFindingActionPlanCommentDetailDTO, error)
	Create(findingID uint, actionPlanID uint, currentUserID uint, dto *auditDTOs.AuditFindingActionPlanCommentFormDTO) (uint, error)
	Update(findingID uint, actionPlanID uint, id uint, currentUserID uint, dto *auditDTOs.AuditFindingActionPlanCommentFormDTO) error
	Delete(findingID uint, actionPlanID uint, ids []uint, currentUserID uint) ([]uint, error)
	GetDocumentPath(findingID, actionPlanID, id uint, token string) (string, error)
}

type AuditFindingActionPlanCommentService struct {
	repo                      repositories.AuditFindingActionPlanCommentRepositoryInterface
	userSvc                   organizationServices.UserServiceInterface
	activityLogSvc            logServices.ActivityLogServiceInterface
	auditFindingActionPlanSvc AuditFindingActionPlanServiceInterface
	notificationSvc           commonServices.NotificationServiceInterface
}

func NewAuditFindingActionPlanCommentService(repo repositories.AuditFindingActionPlanCommentRepositoryInterface, userSvc organizationServices.UserServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditFindingActionPlanSvc AuditFindingActionPlanServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) AuditFindingActionPlanCommentServiceInterface {
	return &AuditFindingActionPlanCommentService{
		repo:                      repo,
		userSvc:                   userSvc,
		activityLogSvc:            activityLogSvc,
		auditFindingActionPlanSvc: auditFindingActionPlanSvc,
		notificationSvc:           notificationSvc,
	}
}

func (afaps *AuditFindingActionPlanCommentService) resolveFindingID(findingID uint, actionPlanID uint) (uint, error) {
	if findingID != 0 {
		return findingID, nil
	}
	return afaps.auditFindingActionPlanSvc.GetFindingIDByActionPlan(actionPlanID)
}

func (afaps *AuditFindingActionPlanCommentService) GetAll(findingID uint, actionPlanID uint) ([]*auditDTOs.AuditFindingActionPlanCommentListDTO, error) {
	findingID, err := afaps.resolveFindingID(findingID, actionPlanID)
	if err != nil {
		return nil, err
	}

	if _, err := afaps.auditFindingActionPlanSvc.GetByID(findingID, actionPlanID); err != nil {
		return nil, err
	}

	return afaps.repo.GetAll(actionPlanID)
}

func (afaps *AuditFindingActionPlanCommentService) GetByID(findingID uint, actionPlanID uint, id uint) (*models.AuditFindingActionPlanComment, error) {
	findingID, err := afaps.resolveFindingID(findingID, actionPlanID)
	if err != nil {
		return nil, err
	}

	if _, err := afaps.auditFindingActionPlanSvc.GetByID(findingID, actionPlanID); err != nil {
		return nil, err
	}

	comment, err := afaps.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return comment, nil
}

func (afaps *AuditFindingActionPlanCommentService) GetFullDetailsByID(findingID uint, actionplanID uint, id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditFindingActionPlanCommentDetailDTO, error) {
	findingID, err := afaps.resolveFindingID(findingID, actionplanID)
	if err != nil {
		return nil, err
	}

	if _, err := afaps.auditFindingActionPlanSvc.GetByID(findingID, actionplanID); err != nil {
		fmt.Printf("auditFindingActionPlan for findingID %d, actionPlanID %d\n", findingID, actionplanID)
		return nil, err
	}

	var commentDetail auditDTOs.AuditFindingActionPlanCommentDetailDTO
	comment, err := afaps.repo.GetFullDetailsByID(id, tx...)
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

func (afaps *AuditFindingActionPlanCommentService) Create(findingID uint, actionPlanID uint, currentUserID uint, auditFindingActionPlanCommentFromDTO *auditDTOs.AuditFindingActionPlanCommentFormDTO) (uint, error) {
	tx := afaps.repo.BeginTransaction()

	findingID, err := afaps.resolveFindingID(findingID, actionPlanID)
	if err != nil {
		return 0, err
	}

	if _, err := afaps.auditFindingActionPlanSvc.GetByID(findingID, actionPlanID); err != nil {
		return 0, err
	}

	auditFindingActionPlan, err := afaps.auditFindingActionPlanSvc.GetFullDetailsByID(findingID, actionPlanID)
	if err != nil {
		return 0, err
	}

	comment := &models.AuditFindingActionPlanComment{
		AuditFindingActionPlanID: actionPlanID,
		Comment:                  auditFindingActionPlanCommentFromDTO.Comment,
		CreatedBYID:              currentUserID,
	}
	id, err := afaps.repo.Create(comment)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	fmt.Println(id)
	var documents []*models.AuditFindingActionPlanCommentDocument

	for _, doc := range auditFindingActionPlanCommentFromDTO.Documents {
		dataToEncrypt := []byte(doc.FileName)
		encryptedFileName, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL(
			"audit", "audit-findings", strconv.FormatUint(uint64(findingID), 10),
			"action-plans", strconv.FormatUint(uint64(actionPlanID), 10),
			"comments", strconv.FormatUint(uint64(id), 10),
			"documents",
		)

		document := &models.AuditFindingActionPlanCommentDocument{
			AuditFindingActionPlanCommentID: id,
			URL:                             url,
			Name:                            doc.FileName,
			Size:                            doc.Size,
			Token:                           encryptedFileName,
			OriginalName:                    doc.OriginalName,
			Type:                            doc.Type,
			CreatedBYID:                     currentUserID,
		}

		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := afaps.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}

		for _, doc := range documents {
			storagePath := filepath.Join(
				"audit", "audit-findings", strconv.FormatUint(uint64(findingID), 10),
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
		existingComment, err := afaps.GetFullDetailsByID(findingID, actionPlanID, id)
		if err != nil {
			utility.LogError("Failed to fetch audit finding action plan comment", map[string]interface{}{"error": err.Error()})
			return
		}

		err = afaps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditFindingActionPlanComment{}.TableName(),
			id,
			existingComment,
			currentUserID,
			"",
			"",
			"",
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AuditFindingActionPlanComment{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}

		if *auditFindingActionPlan.ResponsibleUserID == currentUserID {
			if auditFindingActionPlan.ReviewUserID != nil {
				// prepare notification
				urlParam := fmt.Sprintf(`{"audit_finding_id":%d, "action_plan_id":%d}`, findingID, actionPlanID)
				queryParam := ""
				entity := "audit-finding-action-plan"
				action := "comment-added"

				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserIDs:              []uint{*auditFindingActionPlan.ReviewUserID},
					ModuleIdentifier:     "audit",
					NotificationType:     string(commonDTOs.AssignmentType),
					HeadingEn:            "Comment Added For Audit Findings Action Plan",
					HeadingAr:            "تمت إضافة تعليق إلى خطة عمل نتائج التدقيق",
					TitleEn:              auditFindingActionPlan.Title,
					TitleAr:              auditFindingActionPlan.Title,
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					QueryParam:           queryParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				reviewUser, err := afaps.userSvc.GetByID(*auditFindingActionPlan.ReviewUserID)
				if err != nil {
					utility.LogError("Failed to fetch review user", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}
				assignedBy, err := afaps.userSvc.GetByID(currentUserID)
				if err != nil {
					utility.LogError("Failed to fetch commented user", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}

				emailTo := []string{reviewUser.Email}
				subject := fmt.Sprintf("Comment Added for Audit Finding Action Plan:%s", auditFindingActionPlan.Title)
				appURL := utility.GetAppURL()
				emailData := map[string]interface{}{
					"FirstNameEn":                 reviewUser.FirstNameEn,
					"LastNameEn":                  reviewUser.LastNameEn,
					"AuditFindingReferenceNumber": auditFindingActionPlan.AuditFinding.ReferenceNumber,
					"AuditFindingTitle":           auditFindingActionPlan.AuditFinding.Title,
					"ActionPlanTitle":             auditFindingActionPlan.Title,
					"ActionPlanDescription":       auditFindingActionPlan.Description,
					"CommentedByFirstNameEn":      assignedBy.FirstNameEn,
					"CommentedByLastNameEn":       assignedBy.LastNameEn,
					"CommentedOn":                 time.Now().Format("2006-01-02"),
					"RedirectURL":                 fmt.Sprintf("%s/audit-insight/action-plans?action_plan_id=%d", appURL, actionPlanID),
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              emailTo,
					Subject:              subject,
					TemplateRelativePath: "audit/templates/email/audit_finding_action_plan_comment.html",
					EmailData:            emailData,
				}

				// send notification and mail
				err = afaps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "audit_finding_action_plan")
				if err != nil {
					utility.LogError("Failed to send notifications", map[string]interface{}{
						"error": err.Error(),
					})
				}
			}
		} else if auditFindingActionPlan.ReviewUserID != nil && *auditFindingActionPlan.ReviewUserID == currentUserID {
			// prepare notification
			urlParam := fmt.Sprintf(`{"audit_finding_id":%d, "action_plan_id":%d}`, findingID, actionPlanID)
			queryParam := ""
			entity := "audit-finding-action-plan"
			action := "comment-added"

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{*auditFindingActionPlan.ResponsibleUserID},
				ModuleIdentifier:     "audit",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Comment Added For Audit Findings Action Plan",
				HeadingAr:            "تمت إضافة تعليق إلى خطة عمل نتائج التدقيق",
				TitleEn:              auditFindingActionPlan.Title,
				TitleAr:              auditFindingActionPlan.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			responsibleUser, err := afaps.userSvc.GetByID(*auditFindingActionPlan.ResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}
			assignedBy, err := afaps.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch commented user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			emailTo := []string{responsibleUser.Email}
			subject := fmt.Sprintf("Comment Added for Audit Finding Action Plan:%s", auditFindingActionPlan.Title)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":                 responsibleUser.FirstNameEn,
				"LastNameEn":                  responsibleUser.LastNameEn,
				"AuditFindingReferenceNumber": auditFindingActionPlan.AuditFinding.ReferenceNumber,
				"AuditFindingTitle":           auditFindingActionPlan.AuditFinding.Title,
				"ActionPlanTitle":             auditFindingActionPlan.Title,
				"ActionPlanDescription":       auditFindingActionPlan.Description,
				"CommentedByFirstNameEn":      assignedBy.FirstNameEn,
				"CommentedByLastNameEn":       assignedBy.LastNameEn,
				"CommentedOn":                 time.Now().Format("2006-01-02"),
				"RedirectURL":                 fmt.Sprintf("%s/audit-insight/action-plans?action_plan_id=%d", appURL, actionPlanID),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "audit/templates/email/audit_finding_action_plan_comment.html",
				EmailData:            emailData,
			}

			err = afaps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "audit_finding_action_plan")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		} else {
			// if other user added the comment, notification for review user and responsible user
			// prepare notification
			urlParam := fmt.Sprintf(`{"audit_finding_id":%d, "action_plan_id":%d}`, findingID, actionPlanID)
			queryParam := ""
			entity := "audit-finding-action-plan"
			action := "comment-added"

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{*auditFindingActionPlan.ResponsibleUserID},
				ModuleIdentifier:     "audit",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Comment Added For Audit Findings Action Plan",
				HeadingAr:            "تمت إضافة تعليق إلى خطة عمل نتائج التدقيق",
				TitleEn:              auditFindingActionPlan.Title,
				TitleAr:              auditFindingActionPlan.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			responsibleUser, err := afaps.userSvc.GetByID(*auditFindingActionPlan.ResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}
			assignedBy, err := afaps.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch commented user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			emailTo := []string{responsibleUser.Email}
			subject := fmt.Sprintf("Comment Added for Audit Finding Action Plan:%s", auditFindingActionPlan.Title)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":                 responsibleUser.FirstNameEn,
				"LastNameEn":                  responsibleUser.LastNameEn,
				"AuditFindingReferenceNumber": auditFindingActionPlan.AuditFinding.ReferenceNumber,
				"AuditFindingTitle":           auditFindingActionPlan.AuditFinding.Title,
				"ActionPlanTitle":             auditFindingActionPlan.Title,
				"ActionPlanDescription":       auditFindingActionPlan.Description,
				"CommentedByFirstNameEn":      assignedBy.FirstNameEn,
				"CommentedByLastNameEn":       assignedBy.LastNameEn,
				"CommentedOn":                 time.Now().Format("2006-01-02"),
				"RedirectURL":                 fmt.Sprintf("%s/audit-insight/action-plans?action_plan_id=%d", appURL, actionPlanID),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "audit/templates/email/audit_finding_action_plan_comment.html",
				EmailData:            emailData,
			}

			// send notification and mail
			err = afaps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "audit_finding_action_plan")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}

			if auditFindingActionPlan.ReviewUserID != nil {
				// prepare notification
				urlParam := fmt.Sprintf(`{"audit_fiding_id":%d, "action_plan_id":%d}`, findingID, actionPlanID)
				queryParam := ""
				entity := "audit-finding-action-plan"
				action := "comment-added"

				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserIDs:              []uint{*auditFindingActionPlan.ReviewUserID},
					ModuleIdentifier:     "audit",
					NotificationType:     string(commonDTOs.AssignmentType),
					HeadingEn:            "Comment Added For Audit Findings Action Plan",
					HeadingAr:            "تمت إضافة تعليق إلى خطة عمل نتائج التدقيق",
					TitleEn:              auditFindingActionPlan.Title,
					TitleAr:              auditFindingActionPlan.Title,
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					QueryParam:           queryParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				responsibleUser, err := afaps.userSvc.GetByID(*auditFindingActionPlan.ReviewUserID)
				if err != nil {
					utility.LogError("Failed to fetch review user", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}
				assignedBy, err := afaps.userSvc.GetByID(currentUserID)
				if err != nil {
					utility.LogError("Failed to fetch commented user", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}

				emailTo := []string{responsibleUser.Email}
				subject := fmt.Sprintf("Comment Added for Audit Finding Action Plan:%s", auditFindingActionPlan.Title)
				appURL := utility.GetAppURL()
				emailData := map[string]interface{}{
					"FirstNameEn":                 responsibleUser.FirstNameEn,
					"LastNameEn":                  responsibleUser.LastNameEn,
					"AuditFindingReferenceNumber": auditFindingActionPlan.AuditFinding.ReferenceNumber,
					"AuditFindingTitle":           auditFindingActionPlan.AuditFinding.Title,
					"ActionPlanTitle":             auditFindingActionPlan.Title,
					"ActionPlanDescription":       auditFindingActionPlan.Description,
					"CommentedByFirstNameEn":      assignedBy.FirstNameEn,
					"CommentedByLastNameEn":       assignedBy.LastNameEn,
					"CommentedOn":                 time.Now().Format("2006-01-02"),
					"RedirectURL":                 fmt.Sprintf("%s/audit-insight/action-plans?action_plan_id=%d", appURL, actionPlanID),
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              emailTo,
					Subject:              subject,
					TemplateRelativePath: "audit/templates/email/audit_finding_action_plan_comment.html",
					EmailData:            emailData,
				}

				// send notification and mail
				err = afaps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "audit_finding_action_plan")
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

func (afaps *AuditFindingActionPlanCommentService) Update(findingID uint, actionPlanID uint, id uint, currentUserID uint, AuditFindingActionPlanCommentFormDTO *auditDTOs.AuditFindingActionPlanCommentFormDTO) error {
	tx := afaps.repo.BeginTransaction()

	findingID, err := afaps.resolveFindingID(findingID, actionPlanID)
	if err != nil {
		return err
	}

	if _, err := afaps.auditFindingActionPlanSvc.GetByID(findingID, actionPlanID); err != nil {
		return err
	}

	existingComment, err := afaps.GetByID(findingID, actionPlanID, id)
	if err != nil {
		return err
	}

	oldComment, err := afaps.GetFullDetailsByID(findingID, actionPlanID, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingComment.Comment = AuditFindingActionPlanCommentFormDTO.Comment
	err = afaps.repo.Update(existingComment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	var newDocuments []*models.AuditFindingActionPlanCommentDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.AuditFindingActionPlanCommentDocument

	for _, doc := range AuditFindingActionPlanCommentFormDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {

			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}

			url := utility.GenerateDocumentURL("audit", "audit-findings", strconv.FormatUint(uint64(findingID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

			document := &models.AuditFindingActionPlanCommentDocument{
				AuditFindingActionPlanCommentID: id,
				URL:                             url,
				Name:                            doc.FileName,
				Size:                            doc.Size,
				Token:                           encryptedFileName,
				OriginalName:                    doc.OriginalName,
				Type:                            doc.Type,
				CreatedBYID:                     currentUserID,
			}

			newDocuments = append(newDocuments, document)
		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			if doc.ID != nil {
				existingDocument, err := afaps.repo.GetDocumentByID(*doc.ID)
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
		if err := afaps.repo.CreateDocuments(newDocuments, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIds) > 0 {
		if err := afaps.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range newDocuments {
		storagePath := filepath.Join("audit", "audit-findings", strconv.FormatUint(uint64(findingID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("audit", "audit-findings", strconv.FormatUint(uint64(findingID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()

	go func() {
		existingComment, err := afaps.GetFullDetailsByID(findingID, actionPlanID, id)
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
			err := afaps.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.AuditFindingActionPlanComment{}.TableName(),
				id,
				changedData,
				currentUserID,
				"",
				"",
				"",
				nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.AuditFindingActionPlanComment{}.TableName(),
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

func (afaps *AuditFindingActionPlanCommentService) Delete(findingID uint, actionPlanID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := afaps.repo.BeginTransaction()

	findingID, err := afaps.resolveFindingID(findingID, actionPlanID)
	if err != nil {
		return nil, err
	}

	if _, err := afaps.auditFindingActionPlanSvc.GetByID(findingID, actionPlanID); err != nil {
		return nil, err
	}

	var failedIDs []uint
	var successfulIDs []uint
	var deletedComments []map[string]interface{}
	for _, id := range ids {
		if _, err := afaps.repo.GetFullDetailsByID(id); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		deletedComments = append(deletedComments, map[string]interface{}{
			"id":                id,
			"compliance_action": actionPlanID,
		})

		err := afaps.repo.Delete(id, tx)
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
			err := afaps.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.ComplianceActionPlanComment{}.TableName(),
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

func (afaps *AuditFindingActionPlanCommentService) GetDocumentPath(findingID, actionPlanID, id uint, token string) (string, error) {
	findingID, err := afaps.resolveFindingID(findingID, actionPlanID)
	if err != nil {
		return "", err
	}

	existingComment, err := afaps.GetByID(findingID, actionPlanID, id)
	if err != nil {
		return "", err
	}
	if len(*existingComment.Documents) == 0 {
		return "", fmt.Errorf("no documents found for the compliance action plan comment")
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	storagePath := filepath.Join("audit", "audit-findings", strconv.FormatUint(uint64(findingID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}
