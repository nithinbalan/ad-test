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
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type SLAActionPlanCommentServiceInterface interface {
	GetAll(slaID uint, actionPlanID uint) ([]*complianceDTOs.SLAActionPlanCommentListDTO, error)
	GetByID(slaID uint, actionPlanID uint, id uint) (*models.SLAActionPlanComment, error)
	GetFullDetailsByID(slaID uint, actionPlanID uint, id uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.SLAActionPlanCommentDetailDTO, error)
	Create(slaID uint, actionPlanID uint, currentUserID uint, dto *complianceDTOs.SLAActionPlanCommentFormDTO) (uint, error)
	Update(slaID uint, actionPlanID uint, id uint, currentUserID uint, dto *complianceDTOs.SLAActionPlanCommentFormDTO) error
	Delete(slaID uint, actionPlanID uint, ids []uint, currentUserID uint) ([]uint, error)

	GetDocumentPath(slaID, actionPlanID, id uint, token string) (string, error)
}

type SLAActionPlanCommentService struct {
	repo             repositories.SLAActionPlanCommentRepositoryInterface
	userSvc          organizationServices.UserServiceInterface
	activityLogSvc   logServices.ActivityLogServiceInterface
	slaActionPlanSvc SLAActionPlanServiceInterface
	notificationSvc  commonServices.NotificationServiceInterface
}

func NewSLAActionPlanCommentService(repo repositories.SLAActionPlanCommentRepositoryInterface,
	userSvc organizationServices.UserServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	slaActionPlanSvc SLAActionPlanServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) SLAActionPlanCommentServiceInterface {
	return &SLAActionPlanCommentService{
		repo:             repo,
		userSvc:          userSvc,
		activityLogSvc:   activityLogSvc,
		slaActionPlanSvc: slaActionPlanSvc,
		notificationSvc:  notificationSvc,
	}
}

func (svc *SLAActionPlanCommentService) GetAll(slaID uint, actionPlanID uint) ([]*complianceDTOs.SLAActionPlanCommentListDTO, error) {
	if _, err := svc.slaActionPlanSvc.GetByID(slaID, actionPlanID); err != nil {
		return nil, err
	}
	return svc.repo.GetAll(actionPlanID)
}

func (sapcs *SLAActionPlanCommentService) GetByID(slaID uint, actionPlanID uint, id uint) (*models.SLAActionPlanComment, error) {
	if _, err := sapcs.slaActionPlanSvc.GetByID(slaID, actionPlanID); err != nil {
		fmt.Printf("SLAActionPlan not found for slaID %d, actionPlanID %d\n", slaID, actionPlanID)
		return nil, err
	}

	comment, err := sapcs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return comment, nil
}

func (sapcs *SLAActionPlanCommentService) GetFullDetailsByID(slaID uint, actionPlanID uint, id uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.SLAActionPlanCommentDetailDTO, error) {
	if _, err := sapcs.slaActionPlanSvc.GetByID(slaID, actionPlanID); err != nil {
		return nil, err
	}

	var commentDetail complianceDTOs.SLAActionPlanCommentDetailDTO

	comment, err := sapcs.repo.GetFullDetailsByID(id, tx...)
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

func (sapcs *SLAActionPlanCommentService) Create(slaID uint, actionPlanID uint, currentUserID uint, slaActionPlanCommentFormDTO *complianceDTOs.SLAActionPlanCommentFormDTO) (uint, error) {
	tx := sapcs.repo.BeginTransaction()

	if _, err := sapcs.slaActionPlanSvc.GetByID(slaID, actionPlanID); err != nil {
		return 0, err
	}

	slaActionPlan, err := sapcs.slaActionPlanSvc.GetFullDetailsByID(slaID, actionPlanID)
	if err != nil {
		return 0, err
	}
	comment := &models.SLAActionPlanComment{
		SLAActionPlanID: actionPlanID,
		Comment:         slaActionPlanCommentFormDTO.Comment,
		CreatedByID:     currentUserID,
	}

	id, err := sapcs.repo.Create(comment)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.SLAActionPlanCommentDocument
	for _, doc := range slaActionPlanCommentFormDTO.Documents {
		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL(
			"compliance", "sla-register", strconv.FormatUint(uint64(slaID), 10),
			"action-plans", strconv.FormatUint(uint64(actionPlanID), 10),
			"comments", strconv.FormatUint(uint64(id), 10),
			"documents",
		)

		document := &models.SLAActionPlanCommentDocument{
			SLAActionPlanCommentID: id,
			URL:                    url,
			Name:                   doc.FileName,
			Size:                   doc.Size,
			Token:                  encryptedFilename,
			OriginalName:           doc.OriginalName,
			Type:                   doc.Type,
			CreatedByID:            currentUserID,
		}
		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := sapcs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}

		for _, doc := range documents {
			storagePath := filepath.Join("compliance", "sla-register", strconv.FormatUint(uint64(slaID), 10),
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
		existingComment, err := sapcs.GetFullDetailsByID(slaID, actionPlanID, id)
		if err != nil {
			utility.LogError("Failed to fetch SLA action plan comment", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		err = sapcs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.SLAActionPlanComment{}.TableName(),
			id,
			existingComment,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.SLAActionPlanComment{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}
		if slaActionPlan.ResponsibleUserID == currentUserID {
			if slaActionPlan.ReviewUserID != nil {
				urlParam := fmt.Sprintf(`{"sla_contract_id":%d, "action_plan_id":%d}`, slaID, actionPlanID)
				entity := "sla-register-action-plan"
				action := "comment-added"

				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserIDs:              []uint{*slaActionPlan.ReviewUserID},
					ModuleIdentifier:     "compliance",
					NotificationType:     string(commonDTOs.AssignmentType),
					HeadingEn:            "Comment Added For SLA Action Plan",
					HeadingAr:            "تمت إضافة التعليق لخطة عمل SLA",
					TitleEn:              slaActionPlan.Title,
					TitleAr:              slaActionPlan.Title,
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					QueryParam:           "",
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				reviewUser, err := sapcs.userSvc.GetByID(*slaActionPlan.ReviewUserID)
				if err != nil {
					utility.LogError("Failed to fetch review user", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}
				assignedBy, err := sapcs.userSvc.GetByID(currentUserID)
				if err != nil {
					utility.LogError("Failed to fetch commented user", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}

				emailTo := []string{reviewUser.Email}
				subject := fmt.Sprintf("Comment Added for SLA Action Plan: %s", slaActionPlan.Title)
				appURL := utility.GetAppURL()
				emailData := map[string]interface{}{
					"FirstNameEn":            reviewUser.FirstNameEn,
					"LastNameEn":             reviewUser.LastNameEn,
					"SLAReferenceNumber":     slaActionPlan.SlaRegister.ReferenceNumber,
					"SLATitle":               slaActionPlan.SlaRegister.TitleEn,
					"ActionPlanTitle":        slaActionPlan.Title,
					"ActionPlanDescription":  slaActionPlan.Description,
					"CommentedByFirstNameEn": assignedBy.FirstNameEn,
					"CommentedByLastNameEn":  assignedBy.LastNameEn,
					"CommentedOn":            time.Now().Format("2006-01-02"),
					"RedirectURL":            fmt.Sprintf("%scompliance/sla-contracts/%d/action-plans?action_plan_id=%d", appURL, slaID, actionPlanID),
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              emailTo,
					Subject:              subject,
					TemplateRelativePath: "compliance/templates/email/sla_action_plan_comment.html",
					EmailData:            emailData,
				}

				err = sapcs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "sla-register-action-plan")
				if err != nil {
					utility.LogError("Failed to send notifications", map[string]interface{}{
						"error": err.Error(),
					})
				}
			}
		} else if slaActionPlan.ReviewUserID != nil && *slaActionPlan.ReviewUserID == currentUserID {

			urlParam := fmt.Sprintf(`{"sla_contract_id":%d, "action_plan_id":%d}`, slaID, actionPlanID)
			entity := "sla-register-action-plan"
			action := "comment-added"
			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{slaActionPlan.ResponsibleUserID},
				ModuleIdentifier:     "compliance",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Comment Added For SLA Action Plan",
				HeadingAr:            "تمت إضافة التعليق لخطة عمل SLA",
				TitleEn:              slaActionPlan.Title,
				TitleAr:              slaActionPlan.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				QueryParam:           "",
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}
			responsibleUser, err := sapcs.userSvc.GetByID(slaActionPlan.ResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}
			assignedBy, err := sapcs.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch commented user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}
			emailTo := []string{responsibleUser.Email}
			subject := fmt.Sprintf("Comment Added for SLA Action Plan: %s", slaActionPlan.Title)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":            responsibleUser.FirstNameEn,
				"LastNameEn":             responsibleUser.LastNameEn,
				"SLAReferenceNumber":     slaActionPlan.SlaRegister.ReferenceNumber,
				"SLATitle":               slaActionPlan.SlaRegister.TitleEn,
				"ActionPlanTitle":        slaActionPlan.Title,
				"ActionPlanDescription":  slaActionPlan.Description,
				"CommentedByFirstNameEn": assignedBy.FirstNameEn,
				"CommentedByLastNameEn":  assignedBy.LastNameEn,
				"CommentedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":            fmt.Sprintf("%scompliance/sla-contracts/%d/action-plans?action_plan_id=%d", appURL, slaID, actionPlanID),
			}
			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "compliance/templates/email/sla_action_plan_comment.html",
				EmailData:            emailData,
			}

			err = sapcs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "sla_register_action_plan")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		} else {

			urlParam := fmt.Sprintf(`{"sla_contract_id":%d, "action_plan_id":%d}`, slaID, actionPlanID)
			entity := "sla-register-action-plan"
			action := "comment-added"

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{slaActionPlan.ResponsibleUserID},
				ModuleIdentifier:     "compliance",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Comment Added For SLA Action Plan",
				HeadingAr:            "تمت إضافة التعليق لخطة عمل SLA",
				TitleEn:              slaActionPlan.Title,
				TitleAr:              slaActionPlan.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				QueryParam:           "",
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			responsibleUser, err := sapcs.userSvc.GetByID(slaActionPlan.ResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}
			assignedBy, err := sapcs.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch commented user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}
			fmt.Println(&slaActionPlan.SlaRegister.ReferenceNumber, "ref")

			emailTo := []string{responsibleUser.Email}
			subject := fmt.Sprintf("Comment Added for SLA Action Plan: %s", slaActionPlan.Title)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":            responsibleUser.FirstNameEn,
				"LastNameEn":             responsibleUser.LastNameEn,
				"SLAReferenceNumber":     slaActionPlan.SlaRegister.ReferenceNumber,
				"SLATitle":               slaActionPlan.SlaRegister.TitleEn,
				"ActionPlanTitle":        slaActionPlan.Title,
				"ActionPlanDescription":  slaActionPlan.Description,
				"CommentedByFirstNameEn": assignedBy.FirstNameEn,
				"CommentedByLastNameEn":  assignedBy.LastNameEn,
				"AssignedOn":             time.Now().Format("2006-01-02"),
				"RedirectURL":            fmt.Sprintf("%scompliance/sla-contracts/%d/action-plans?action_plan_id=%d", appURL, slaID, actionPlanID),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "compliance/templates/email/sla_action_plan_comment.html",
				EmailData:            emailData,
			}

			err = sapcs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "sla-register-action-plan")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}

			if slaActionPlan.ReviewUserID != nil {

				urlParam := fmt.Sprintf(`{"sla_contract_id":%d, "action_plan_id":%d}`, slaID, actionPlanID)
				queryParam := ""
				entity := "sla-register-action-plan"
				action := "comment-added"

				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserIDs:              []uint{*slaActionPlan.ReviewUserID},
					ModuleIdentifier:     "compliance",
					NotificationType:     string(commonDTOs.AssignmentType),
					HeadingEn:            "Comment Added For SLA Action Plan",
					HeadingAr:            "تمت إضافة التعليق لخطة عمل SLA",
					TitleEn:              slaActionPlan.Title,
					TitleAr:              slaActionPlan.Title,
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					QueryParam:           queryParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				reviewUser, err := sapcs.userSvc.GetByID(*slaActionPlan.ReviewUserID)
				if err != nil {
					utility.LogError("Failed to fetch review user", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}
				assignedBy, err := sapcs.userSvc.GetByID(currentUserID)
				if err != nil {
					utility.LogError("Failed to fetch commented user", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}

				emailTo := []string{reviewUser.Email}
				subject := fmt.Sprintf("Comment Added for SLA Action Plan: %s", slaActionPlan.Title)
				appURL := utility.GetAppURL()
				emailData := map[string]interface{}{
					"FirstNameEn":            reviewUser.FirstNameEn,
					"LastNameEn":             reviewUser.LastNameEn,
					"SLAReferenceNumber":     slaActionPlan.SlaRegister.ReferenceNumber,
					"SLATitle":               slaActionPlan.SlaRegister.TitleEn,
					"ActionPlanTitle":        slaActionPlan.Title,
					"ActionPlanDescription":  slaActionPlan.Description,
					"CommentedByFirstNameEn": assignedBy.FirstNameEn,
					"CommentedByLastNameEn":  assignedBy.LastNameEn,
					"AssignedOn":             time.Now().Format("2006-01-02"),
					"RedirectURL":            fmt.Sprintf("%scompliance/sla-contracts/%d/action-plans?action_plan_id=%d", appURL, slaID, actionPlanID),
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              emailTo,
					Subject:              subject,
					TemplateRelativePath: "compliance/templates/email/sla_action_plan_comment.html",
					EmailData:            emailData,
				}

				err = sapcs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "sla_register_action_plan")
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

func (sapcs *SLAActionPlanCommentService) Update(slaID uint, actionPlanID uint, id uint, currentUserID uint, dto *complianceDTOs.SLAActionPlanCommentFormDTO) error {
	tx := sapcs.repo.BeginTransaction()

	if _, err := sapcs.slaActionPlanSvc.GetByID(slaID, actionPlanID); err != nil {
		return err
	}

	existingComment, err := sapcs.GetByID(slaID, actionPlanID, id)
	if err != nil {
		return err
	}

	oldComment, err := sapcs.GetFullDetailsByID(slaID, actionPlanID, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingComment.Comment = dto.Comment
	err = sapcs.repo.Update(existingComment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	var newDocuments []*models.SLAActionPlanCommentDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.SLAActionPlanCommentDocument

	for _, doc := range dto.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {
			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}

			url := utility.GenerateDocumentURL("compliance", "sla-register", strconv.FormatUint(uint64(slaID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

			document := &models.SLAActionPlanCommentDocument{
				SLAActionPlanCommentID: id,
				URL:                    url,
				Name:                   doc.FileName,
				Size:                   doc.Size,
				Token:                  encryptedFileName,
				OriginalName:           doc.OriginalName,
				Type:                   doc.Type,
				CreatedByID:            currentUserID,
			}
			newDocuments = append(newDocuments, document)
		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			if doc.ID != nil {
				existingDocument, err := sapcs.repo.GetDocumentByID(*doc.ID)
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
		if err := sapcs.repo.CreateDocuments(newDocuments, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIds) > 0 {
		if err := sapcs.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range newDocuments {
		storagePath := filepath.Join("compliance", "sla-register", strconv.FormatUint(uint64(slaID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("compliance", "sla-register", strconv.FormatUint(uint64(slaID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	go func() {
		existingComment, err := sapcs.GetFullDetailsByID(slaID, actionPlanID, id)
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
			err := sapcs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.SLAActionPlanComment{}.TableName(),
				id,
				changedData,
				currentUserID,
				"", "", "", nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.SLAActionPlanComment{}.TableName(),
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

func (sapcs *SLAActionPlanCommentService) Delete(slaID uint, actionPlanID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := sapcs.repo.BeginTransaction()

	if _, err := sapcs.slaActionPlanSvc.GetByID(slaID, actionPlanID); err != nil {
		return nil, err
	}

	var failedIDs []uint
	var successfulIDs []uint
	var deletedComments []map[string]interface{}
	for _, id := range ids {
		if _, err := sapcs.repo.GetFullDetailsByID(id); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}
		deletedComments = append(deletedComments, map[string]interface{}{
			"id":         id,
			"sla_action": actionPlanID,
		})
		if err := sapcs.repo.Delete(id, tx); err != nil {
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
			err := sapcs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.SLAActionPlanComment{}.TableName(),
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

func (sapcs *SLAActionPlanCommentService) GetDocumentPath(slaID, actionPlanID, id uint, token string) (string, error) {
	existingComment, err := sapcs.GetByID(slaID, actionPlanID, id)
	if err != nil {
		return "", err
	}
	if len(existingComment.Documents) == 0 {
		return "", fmt.Errorf("no documents found for the SLA action plan comment")
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	storagePath := filepath.Join("compliance", "sla-register", strconv.FormatUint(uint64(slaID), 10),
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
