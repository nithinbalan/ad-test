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

type ComplianceActionPlanCommentServiceInterface interface {
	GetAll(complianceID uint, actionPlanID uint) ([]*complianceDTOs.ComplianceActionPlanCommentListDTO, error)
	GetByID(complianceID uint, actionPlanID uint, id uint) (*models.ComplianceActionPlanComment, error)
	GetFullDetailsByID(complianceID uint, actionPlanID uint, id uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceActionPlanCommentDetailDTO, error)
	Create(complianceID uint, actionPlanID uint, currentUserID uint, dto *complianceDTOs.ComplianceActionPlanCommentFormDTO) (uint, error)
	Update(complianceID uint, actionPlanID uint, id uint, currentUserID uint, dto *complianceDTOs.ComplianceActionPlanCommentFormDTO) error
	Delete(complianceID uint, actionPlanID uint, ids []uint, currentUserID uint) ([]uint, error)

	GetDocumentPath(complianceID, actionPlanID, id uint, token string) (string, error)
}

type ComplianceActionPlanCommentService struct {
	repo                    repositories.ComplianceActionPlanCommentRepositoryInterface
	userSvc                 organizationServices.UserServiceInterface
	activityLogSvc          logServices.ActivityLogServiceInterface
	complianceActionPlanSvc ComplianceActionPlanServiceInterface
	notificationSvc         commonServices.NotificationServiceInterface
}

func NewComplianceActionPlanCommentService(repo repositories.ComplianceActionPlanCommentRepositoryInterface, userSvc organizationServices.UserServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	complianceActionPlanSvc ComplianceActionPlanServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) ComplianceActionPlanCommentServiceInterface {
	return &ComplianceActionPlanCommentService{
		repo:                    repo,
		userSvc:                 userSvc,
		activityLogSvc:          activityLogSvc,
		complianceActionPlanSvc: complianceActionPlanSvc,
		notificationSvc:         notificationSvc,
	}
}

func (caps *ComplianceActionPlanCommentService) GetAll(complianceID uint, actionPlanID uint) ([]*complianceDTOs.ComplianceActionPlanCommentListDTO, error) {
	if _, err := caps.complianceActionPlanSvc.GetByID(complianceID, actionPlanID); err != nil {
		return nil, err
	}

	return caps.repo.GetAll(actionPlanID)
}

func (caps *ComplianceActionPlanCommentService) GetByID(complianceID uint, actionPlanID uint, id uint) (*models.ComplianceActionPlanComment, error) {
	if _, err := caps.complianceActionPlanSvc.GetByID(complianceID, actionPlanID); err != nil {
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

func (caps *ComplianceActionPlanCommentService) GetFullDetailsByID(complianceID uint, actionPlanID uint, id uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceActionPlanCommentDetailDTO, error) {
	if _, err := caps.complianceActionPlanSvc.GetByID(complianceID, actionPlanID); err != nil {
		fmt.Printf("ComplianceActionPlan not found for complianceID %d, actionPlanID %d\n", complianceID, actionPlanID)
		return nil, err
	}

	var commentDetail complianceDTOs.ComplianceActionPlanCommentDetailDTO

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

func (caps *ComplianceActionPlanCommentService) Create(complianceID uint, actionPlanID uint, currentUserID uint, complianceActionPlanCommentFormDTO *complianceDTOs.ComplianceActionPlanCommentFormDTO) (uint, error) {
	tx := caps.repo.BeginTransaction()

	if _, err := caps.complianceActionPlanSvc.GetByID(complianceID, actionPlanID); err != nil {
		return 0, err
	}

	complianceActionPlan, err := caps.complianceActionPlanSvc.GetFullDetailsByID(complianceID, actionPlanID)
	if err != nil {
		return 0, err
	}

	comment := &models.ComplianceActionPlanComment{
		ComplianceActionPlanID: actionPlanID,
		Comment:                complianceActionPlanCommentFormDTO.Comment,
		CreatedBYID:            currentUserID,
	}

	id, err := caps.repo.Create(comment)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.ComplianceActionPlanCommentDocument

	for _, doc := range complianceActionPlanCommentFormDTO.Documents {
		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL(
			"compliance", "compliance-register", strconv.FormatUint(uint64(complianceID), 10),
			"action-plans", strconv.FormatUint(uint64(actionPlanID), 10),
			"comments", strconv.FormatUint(uint64(id), 10),
			"documents",
		)

		document := &models.ComplianceActionPlanCommentDocument{
			ComplianceActionPlanCommentID: id,
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
				"compliance", "compliance-register", strconv.FormatUint(uint64(complianceID), 10),
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
		existingComment, err := caps.GetFullDetailsByID(complianceID, actionPlanID, id)
		if err != nil {
			utility.LogError("Failed to fetch compliance action plan comment", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = caps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ComplianceActionPlanComment{}.TableName(),
			id,
			existingComment,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ComplianceActionPlanComment{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}

		if complianceActionPlan.ResponsibleUserID == currentUserID {
			// if responsible user added the comment, notification for review user
			if complianceActionPlan.ReviewUserID != nil {
				// prepare notification
				urlParam := fmt.Sprintf(`{"compliance_register_id":%d, "action_plan_id":%d}`, complianceID, id)
				queryParam := ""
				entity := "compliance-action-plan"
				action := "comment-added"

				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserID:               *complianceActionPlan.ReviewUserID,
					ModuleIdentifier:     "compliance",
					NotificationType:     string(commonDTOs.AssignmentType),
					HeadingEn:            "Comment Added For Compliance Action Plan",
					HeadingAr:            "تمت إضافة التعليق لخطة عمل الامتثال",
					TitleEn:              complianceActionPlan.Title,
					TitleAr:              complianceActionPlan.Title,
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					QueryParam:           queryParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				reviewUser, err := caps.userSvc.GetByID(*complianceActionPlan.ReviewUserID)
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
				subject := fmt.Sprintf("Comment Added for Compliance Action Plan:%s", complianceActionPlan.Title)
				appURL := utility.GetAppURL()
				emailData := map[string]interface{}{
					"FirstNameEn":               reviewUser.FirstNameEn,
					"LastNameEn":                reviewUser.LastNameEn,
					"ComplianceReferenceNumber": complianceActionPlan.ComplianceRegister.ReferenceNumber,
					"ComplianceTitle":           complianceActionPlan.ComplianceRegister.TitleEn,
					"ActionPlanTitle":           complianceActionPlan.Title,
					"ActionPlanDescription":     complianceActionPlan.Description,
					"CommentedByFirstNameEn":    assignedBy.FirstNameEn,
					"CommentedByLastNameEn":     assignedBy.LastNameEn,
					"CommentedOn":               time.Now().Format("2006-01-02"),
					"RedirectURL":               fmt.Sprintf("%scompliance/register/%d/action-plans?action_plan_id=%d", appURL, complianceID, id),
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              emailTo,
					Subject:              subject,
					TemplateRelativePath: "compliance/templates/email/compliance_action_plan_comment.html",
					EmailData:            emailData,
				}

				// send notification and mail
				err = caps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "compliance_action_plan")
				if err != nil {
					utility.LogError("Failed to send notifications", map[string]interface{}{
						"error": err.Error(),
					})
				}
			}
		} else if complianceActionPlan.ReviewUserID == &currentUserID {
			// if review user added the comment, notification for responsible user
			// prepare notification
			urlParam := fmt.Sprintf(`{"compliance_register_id":%d, "action_plan_id":%d}`, complianceID, id)
			queryParam := ""
			entity := "compliance-action-plan"
			action := "comment-added"

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserID:               complianceActionPlan.ResponsibleUserID,
				ModuleIdentifier:     "compliance",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Comment Added For Compliance Action Plan",
				HeadingAr:            "تمت إضافة التعليق لخطة عمل الامتثال",
				TitleEn:              complianceActionPlan.Title,
				TitleAr:              complianceActionPlan.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			responsibleUser, err := caps.userSvc.GetByID(complianceActionPlan.ResponsibleUserID)
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
			subject := fmt.Sprintf("Comment Added for Compliance Action Plan:%s", complianceActionPlan.Title)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":               responsibleUser.FirstNameEn,
				"LastNameEn":                responsibleUser.LastNameEn,
				"ComplianceReferenceNumber": complianceActionPlan.ComplianceRegister.ReferenceNumber,
				"ComplianceTitle":           complianceActionPlan.ComplianceRegister.TitleEn,
				"ActionPlanTitle":           complianceActionPlan.Title,
				"ActionPlanDescription":     complianceActionPlan.Description,
				"CommentedByFirstNameEn":    assignedBy.FirstNameEn,
				"CommentedByLastNameEn":     assignedBy.LastNameEn,
				"AssignedOn":                time.Now().Format("2006-01-02"),
				"RedirectURL":               fmt.Sprintf("%scompliance/register/%d/action-plans?action_plan_id=%d", appURL, complianceID, id),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "compliance/templates/email/compliance_action_plan_comment.html",
				EmailData:            emailData,
			}

			// send notification and mail
			err = caps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "compliance_action_plan")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		} else {
			// if other user added the comment, notification for review user and responsible user
			// prepare notification
			urlParam := fmt.Sprintf(`{"compliance_register_id":%d, "action_plan_id":%d}`, complianceID, id)
			queryParam := ""
			entity := "compliance-action-plan"
			action := "comment-added"

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserID:               complianceActionPlan.ResponsibleUserID,
				ModuleIdentifier:     "compliance",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Comment Added For Compliance Action Plan",
				HeadingAr:            "تمت إضافة التعليق لخطة عمل الامتثال",
				TitleEn:              complianceActionPlan.Title,
				TitleAr:              complianceActionPlan.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			responsibleUser, err := caps.userSvc.GetByID(complianceActionPlan.ResponsibleUserID)
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
			subject := fmt.Sprintf("Comment Added for Compliance Action Plan:%s", complianceActionPlan.Title)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":               responsibleUser.FirstNameEn,
				"LastNameEn":                responsibleUser.LastNameEn,
				"ComplianceReferenceNumber": complianceActionPlan.ComplianceRegister.ReferenceNumber,
				"ComplianceTitle":           complianceActionPlan.ComplianceRegister.TitleEn,
				"ActionPlanTitle":           complianceActionPlan.Title,
				"ActionPlanDescription":     complianceActionPlan.Description,
				"CommentedByFirstNameEn":    assignedBy.FirstNameEn,
				"CommentedByLastNameEn":     assignedBy.LastNameEn,
				"AssignedOn":                time.Now().Format("2006-01-02"),
				"RedirectURL":               fmt.Sprintf("%scompliance/register/%d/action-plans?action_plan_id=%d", appURL, complianceID, id),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "compliance/templates/email/compliance_action_plan_comment.html",
				EmailData:            emailData,
			}

			// send notification and mail
			err = caps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "compliance_action_plan")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
			if complianceActionPlan.ReviewUserID != nil {
				// prepare notification
				urlParam := fmt.Sprintf(`{"compliance_register_id":%d, "action_plan_id":%d}`, complianceID, id)
				queryParam := ""
				entity := "compliance-action-plan"
				action := "comment-added"

				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserID:               *complianceActionPlan.ReviewUserID,
					ModuleIdentifier:     "compliance",
					NotificationType:     string(commonDTOs.AssignmentType),
					HeadingEn:            "Comment Added For Compliance Action Plan",
					HeadingAr:            "تمت إضافة التعليق لخطة عمل الامتثال",
					TitleEn:              complianceActionPlan.Title,
					TitleAr:              complianceActionPlan.Title,
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					QueryParam:           queryParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				responsibleUser, err := caps.userSvc.GetByID(*complianceActionPlan.ReviewUserID)
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
				subject := fmt.Sprintf("Comment Added for Compliance Action Plan:%s", complianceActionPlan.Title)
				appURL := utility.GetAppURL()
				emailData := map[string]interface{}{
					"FirstNameEn":               responsibleUser.FirstNameEn,
					"LastNameEn":                responsibleUser.LastNameEn,
					"ComplianceReferenceNumber": complianceActionPlan.ComplianceRegister.ReferenceNumber,
					"ComplianceTitle":           complianceActionPlan.ComplianceRegister.TitleEn,
					"ActionPlanTitle":           complianceActionPlan.Title,
					"ActionPlanDescription":     complianceActionPlan.Description,
					"CommentedByFirstNameEn":    assignedBy.FirstNameEn,
					"CommentedByLastNameEn":     assignedBy.LastNameEn,
					"AssignedOn":                time.Now().Format("2006-01-02"),
					"RedirectURL":               fmt.Sprintf("%scompliance/register/%d/action-plans?action_plan_id=%d", appURL, complianceID, id),
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              emailTo,
					Subject:              subject,
					TemplateRelativePath: "compliance/templates/email/compliance_action_plan_comment.html",
					EmailData:            emailData,
				}

				// send notification and mail
				err = caps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "compliance_action_plan")
				if err != nil {
					utility.LogError("Failed to send notifications", map[string]interface{}{
						"error": err.Error(),
					})
				}
			}
		}
		// // prepare notification
		// urlParam := fmt.Sprintf(`{"compliance_register_id":%d, "action_plan_id":%d}`, complianceID, id)
		// queryParam := ""
		// entity := "compliance-action-plan"
		// action := "responsible-user-assign"

		// dbNotificationDTO := &commonDTOs.DBNotificationDTO{
		// 	UserID:               ComplianceActionPlanFormDTO.ResponsibleUserID,
		// 	ModuleIdentifier:     "compliance",
		// 	NotificationType:     string(commonDTOs.AssignmentType),
		// 	HeadingEn:            "Compliance Action Plan Responsibility Assigned",
		// 	HeadingAr:            "تم تعيين مسؤولية خطة عمل الامتثال",
		// 	TitleEn:              ComplianceActionPlanFormDTO.Title,
		// 	TitleAr:              ComplianceActionPlanFormDTO.Title,
		// 	Entity:               &entity,
		// 	Action:               &action,
		// 	URLParam:             urlParam,
		// 	QueryParam:           queryParam,
		// 	IsSystemNotification: false,
		// 	CreatedByID:          &currentUserID,
		// }

		// responsibleUser, err := caps.userSvc.GetByID(ComplianceActionPlanFormDTO.ResponsibleUserID)
		// if err != nil {
		// 	utility.LogError("Failed to fetch responsible user", map[string]interface{}{
		// 		"error": err.Error(),
		// 	})
		// 	return
		// }
		// assignedBy, err := caps.userSvc.GetByID(currentUserID)
		// if err != nil {
		// 	utility.LogError("Failed to fetch assigning user", map[string]interface{}{
		// 		"error": err.Error(),
		// 	})
		// 	return
		// }

		// emailTo := []string{responsibleUser.Email}
		// subject := fmt.Sprintf("compliance Action Plan/%s: Responsibility Assigned", ComplianceActionPlanFormDTO.Title)
		// appURL := utility.GetAppURL()
		// emailData := map[string]interface{}{
		// 	"FirstNameEn":               responsibleUser.FirstNameEn,
		// 	"LastNameEn":                responsibleUser.LastNameEn,
		// 	"ComplianceReferenceNumber": compliance.ReferenceNumber,
		// 	"ComplianceTitle":           compliance.TitleEn,
		// 	"ActionPlanTitle":           ComplianceActionPlanFormDTO.Title,
		// 	"ActionPlanDescription":     ComplianceActionPlanFormDTO.Description,
		// 	"TargetDate":                dates.TargetDate.Format("2006-01-02"),
		// 	"AssignedByFirstNameEn":     assignedBy.FirstNameEn,
		// 	"AssignedByLastNameEn":      assignedBy.LastNameEn,
		// 	"AssignedOn":                time.Now().Format("2006-01-02"),
		// 	"RedirectURL":               fmt.Sprintf("%scompliance/register/%d/action-plans?action_plan_id=%d", appURL, complianceID, id),
		// }

		// emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
		// 	EmailTo:              emailTo,
		// 	Subject:              subject,
		// 	TemplateRelativePath: "compliance/templates/email/compliance_action_plan_create_responsible_user_mail.html",
		// 	EmailData:            emailData,
		// }

		// // send notification and mail
		// err = caps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "compliance_action_plan")
		// if err != nil {
		// 	utility.LogError("Failed to send notifications", map[string]interface{}{
		// 		"error": err.Error(),
		// 	})
		// }
	}()
	return id, nil
}

func (caps *ComplianceActionPlanCommentService) Update(complianceID uint, actionPlanID uint, id uint, currentUserID uint, complianceActionPlanCommentFormDTO *complianceDTOs.ComplianceActionPlanCommentFormDTO) error {
	tx := caps.repo.BeginTransaction()

	if _, err := caps.complianceActionPlanSvc.GetByID(complianceID, actionPlanID); err != nil {
		return err
	}

	existingComment, err := caps.GetByID(complianceID, actionPlanID, id)
	if err != nil {
		return err
	}

	oldComment, err := caps.GetFullDetailsByID(complianceID, actionPlanID, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingComment.Comment = complianceActionPlanCommentFormDTO.Comment
	err = caps.repo.Update(existingComment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	var newDocuments []*models.ComplianceActionPlanCommentDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.ComplianceActionPlanCommentDocument

	for _, doc := range complianceActionPlanCommentFormDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {

			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}

			url := utility.GenerateDocumentURL("compliance", "compliance-register", strconv.FormatUint(uint64(complianceID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

			document := &models.ComplianceActionPlanCommentDocument{
				ComplianceActionPlanCommentID: id,
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
		storagePath := filepath.Join("compliance", "compliance-register", strconv.FormatUint(uint64(complianceID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("compliance", "compliance-register", strconv.FormatUint(uint64(complianceID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	go func() {
		existingComment, err := caps.GetFullDetailsByID(complianceID, actionPlanID, id)
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
				models.ComplianceActionPlanComment{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.ComplianceActionPlanComment{}.TableName(),
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

func (caps *ComplianceActionPlanCommentService) Delete(complianceID uint, actionPlanID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := caps.repo.BeginTransaction()
	if _, err := caps.complianceActionPlanSvc.GetByID(complianceID, actionPlanID); err != nil {
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
			"compliance_action": actionPlanID,
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

func (caps *ComplianceActionPlanCommentService) GetDocumentPath(complianceID, actionPlanID, id uint, token string) (string, error) {
	existingComment, err := caps.GetByID(complianceID, actionPlanID, id)
	if err != nil {
		return "", err
	}
	if len(existingComment.Documents) == 0 {
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

	storagePath := filepath.Join("compliance", "compliance-register", strconv.FormatUint(uint64(complianceID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}
