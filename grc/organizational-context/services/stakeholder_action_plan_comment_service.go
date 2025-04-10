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
	organizationalContextDTOs "jethur/grc/organizational-context/dtos"
	"jethur/grc/organizational-context/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type StakeholderActionPlanCommentServiceInterface interface {
	GetAll(stakeholderID uint, actionPlanID uint) ([]*organizationalContextDTOs.StakeholderActionPlanCommentListDTO, error)
	GetByID(stakeholderID uint, actionPlanID uint, id uint) (*models.StakeholderActionPlanComment, error)
	GetFullDetailsByID(stakeholderID uint, actionPlanID uint, id uint, tx ...utility.DatabaseTransaction) (*organizationalContextDTOs.StakeholderActionPlanCommentDetailDTO, error)
	Create(stakeholderID uint, actionPlanID uint, currentUserID uint, dto *organizationalContextDTOs.StakeholderActionPlanCommentFormDTO) (uint, error)
	Update(stakeholderID uint, actionPlanID uint, id uint, currentUserID uint, dto *organizationalContextDTOs.StakeholderActionPlanCommentFormDTO) error
	Delete(stakeholderID uint, actionPlanID uint, ids []uint, currentUserID uint) ([]uint, error)

	GetDocumentPath(stakeholderID, actionPlanID, id uint, token string) (string, error)
}

type StakeholderActionPlanCommentService struct {
	repo                     repositories.StakeholderActionPlanCommentRepositoryInterface
	userSvc                  organizationServices.UserServiceInterface
	activityLogSvc           logServices.ActivityLogServiceInterface
	stakeholderActionPlanSvc StakeholderActionPlanServiceInterface
	notificationSvc          commonServices.NotificationServiceInterface
}

func NewStakeholderActionPlanCommentService(repo repositories.StakeholderActionPlanCommentRepositoryInterface, userSvc organizationServices.UserServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	stakeholderActionPlanSvc StakeholderActionPlanServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) StakeholderActionPlanCommentServiceInterface {
	return &StakeholderActionPlanCommentService{
		repo:                     repo,
		userSvc:                  userSvc,
		activityLogSvc:           activityLogSvc,
		stakeholderActionPlanSvc: stakeholderActionPlanSvc,
		notificationSvc:          notificationSvc,
	}
}

func (caps *StakeholderActionPlanCommentService) GetAll(stakeholderID uint, actionPlanID uint) ([]*organizationalContextDTOs.StakeholderActionPlanCommentListDTO, error) {
	if _, err := caps.stakeholderActionPlanSvc.GetByID(stakeholderID, actionPlanID); err != nil {
		return nil, err
	}

	return caps.repo.GetAll(actionPlanID)
}

func (caps *StakeholderActionPlanCommentService) GetByID(stakeholderID uint, actionPlanID uint, id uint) (*models.StakeholderActionPlanComment, error) {
	if _, err := caps.stakeholderActionPlanSvc.GetByID(stakeholderID, actionPlanID); err != nil {
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

func (caps *StakeholderActionPlanCommentService) GetFullDetailsByID(stakeholderID uint, actionPlanID uint, id uint, tx ...utility.DatabaseTransaction) (*organizationalContextDTOs.StakeholderActionPlanCommentDetailDTO, error) {
	if _, err := caps.stakeholderActionPlanSvc.GetByID(stakeholderID, actionPlanID); err != nil {
		fmt.Printf("StakeholderActionPlan not found for stakeholderID %d, actionPlanID %d\n", stakeholderID, actionPlanID)
		return nil, err
	}

	var commentDetail organizationalContextDTOs.StakeholderActionPlanCommentDetailDTO

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

func (caps *StakeholderActionPlanCommentService) Create(stakeholderID uint, actionPlanID uint, currentUserID uint, stakeholderActionPlanCommentFormDTO *organizationalContextDTOs.StakeholderActionPlanCommentFormDTO) (uint, error) {
	tx := caps.repo.BeginTransaction()

	if _, err := caps.stakeholderActionPlanSvc.GetByID(stakeholderID, actionPlanID); err != nil {
		return 0, err
	}

	stakeholderActionPlan, err := caps.stakeholderActionPlanSvc.GetFullDetailsByID(stakeholderID, actionPlanID)
	if err != nil {
		return 0, err
	}

	comment := &models.StakeholderActionPlanComment{
		StakeholderActionPlanID: actionPlanID,
		Comment:                 stakeholderActionPlanCommentFormDTO.Comment,
		CreatedBYID:             currentUserID,
	}

	id, err := caps.repo.Create(comment)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.StakeholderActionPlanCommentDocument

	for _, doc := range stakeholderActionPlanCommentFormDTO.Documents {
		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL(
			"stakeholder", "stakeholder-register", strconv.FormatUint(uint64(stakeholderID), 10),
			"action-plans", strconv.FormatUint(uint64(actionPlanID), 10),
			"comments", strconv.FormatUint(uint64(id), 10),
			"documents",
		)

		document := &models.StakeholderActionPlanCommentDocument{
			StakeholderActionPlanCommentID: id,
			URL:                            url,
			Name:                           doc.FileName,
			Size:                           doc.Size,
			Token:                          encryptedFilename,
			OriginalName:                   doc.OriginalName,
			Type:                           doc.Type,
			CreatedBYID:                    currentUserID,
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
				"stakeholder", "stakeholder-register", strconv.FormatUint(uint64(stakeholderID), 10),
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
		existingComment, err := caps.GetFullDetailsByID(stakeholderID, actionPlanID, id)
		if err != nil {
			utility.LogError("Failed to fetch stakeholder action plan comment", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = caps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.StakeholderActionPlanComment{}.TableName(),
			id,
			existingComment,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.StakeholderActionPlanComment{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}

		if stakeholderActionPlan.ResponsibleUserID == currentUserID {
			// if responsible user added the comment, notification for review user
			if stakeholderActionPlan.ReviewUserID != nil {
				// prepare notification
				urlParam := fmt.Sprintf(`{"stakeholder_register_id":%d, "action_plan_id":%d}`, stakeholderID, id)
				queryParam := ""
				entity := "stakeholder-action-plan"
				action := "comment-added"

				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserIDs:              []uint{*stakeholderActionPlan.ReviewUserID},
					ModuleIdentifier:     "stakeholder",
					NotificationType:     string(commonDTOs.AssignmentType),
					HeadingEn:            "Comment Added For Stakeholder Action Plan",
					HeadingAr:            "تمت إضافة التعليق لخطة عمل الامتثال",
					TitleEn:              stakeholderActionPlan.Title,
					TitleAr:              stakeholderActionPlan.Title,
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					QueryParam:           queryParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				reviewUser, err := caps.userSvc.GetByID(*stakeholderActionPlan.ReviewUserID)
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
				subject := fmt.Sprintf("Comment Added for Stakeholder Action Plan:%s", stakeholderActionPlan.Title)
				appURL := utility.GetAppURL()
				emailData := map[string]interface{}{
					"FirstNameEn":            reviewUser.FirstNameEn,
					"LastNameEn":             reviewUser.LastNameEn,
					"StakeholderTitle":       stakeholderActionPlan.Stakeholder.TitleEn,
					"ActionPlanTitle":        stakeholderActionPlan.Title,
					"ActionPlanDescription":  stakeholderActionPlan.Description,
					"CommentedByFirstNameEn": assignedBy.FirstNameEn,
					"CommentedByLastNameEn":  assignedBy.LastNameEn,
					"CommentedOn":            time.Now().Format("2006-01-02"),
					"RedirectURL":            fmt.Sprintf("%sstakeholder/register/%d/action-plans?action_plan_id=%d", appURL, stakeholderID, id),
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              emailTo,
					Subject:              subject,
					TemplateRelativePath: "stakeholder/templates/email/stakeholder_action_plan_comment.html",
					EmailData:            emailData,
				}

				// send notification and mail
				err = caps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "stakeholder_action_plan")
				if err != nil {
					utility.LogError("Failed to send notifications", map[string]interface{}{
						"error": err.Error(),
					})
				}
			}
		} else if stakeholderActionPlan.ReviewUserID == &currentUserID {
			// if review user added the comment, notification for responsible user
			// prepare notification
			urlParam := fmt.Sprintf(`{"stakeholder_register_id":%d, "action_plan_id":%d}`, stakeholderID, id)
			queryParam := ""
			entity := "stakeholder-action-plan"
			action := "comment-added"

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{stakeholderActionPlan.ResponsibleUserID},
				ModuleIdentifier:     "stakeholder",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Comment Added For Stakeholder Action Plan",
				HeadingAr:            "تمت إضافة التعليق لخطة عمل الامتثال",
				TitleEn:              stakeholderActionPlan.Title,
				TitleAr:              stakeholderActionPlan.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			responsibleUser, err := caps.userSvc.GetByID(stakeholderActionPlan.ResponsibleUserID)
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
			subject := fmt.Sprintf("Comment Added for Stakeholder Action Plan:%s", stakeholderActionPlan.Title)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":            responsibleUser.FirstNameEn,
				"LastNameEn":             responsibleUser.LastNameEn,
				"StakeholderTitle":       stakeholderActionPlan.Stakeholder.TitleEn,
				"ActionPlanTitle":        stakeholderActionPlan.Title,
				"ActionPlanDescription":  stakeholderActionPlan.Description,
				"CommentedByFirstNameEn": assignedBy.FirstNameEn,
				"CommentedByLastNameEn":  assignedBy.LastNameEn,
				"AssignedOn":             time.Now().Format("2006-01-02"),
				"RedirectURL":            fmt.Sprintf("%sstakeholder/register/%d/action-plans?action_plan_id=%d", appURL, stakeholderID, id),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "stakeholder/templates/email/stakeholder_action_plan_comment.html",
				EmailData:            emailData,
			}

			// send notification and mail
			err = caps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "stakeholder_action_plan")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		} else {
			// if other user added the comment, notification for review user and responsible user
			// prepare notification
			urlParam := fmt.Sprintf(`{"stakeholder_register_id":%d, "action_plan_id":%d}`, stakeholderID, id)
			queryParam := ""
			entity := "stakeholder-action-plan"
			action := "comment-added"

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{stakeholderActionPlan.ResponsibleUserID},
				ModuleIdentifier:     "stakeholder",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Comment Added For Stakeholder Action Plan",
				HeadingAr:            "تمت إضافة التعليق لخطة عمل الامتثال",
				TitleEn:              stakeholderActionPlan.Title,
				TitleAr:              stakeholderActionPlan.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			responsibleUser, err := caps.userSvc.GetByID(stakeholderActionPlan.ResponsibleUserID)
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
			subject := fmt.Sprintf("Comment Added for Stakeholder Action Plan:%s", stakeholderActionPlan.Title)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":            responsibleUser.FirstNameEn,
				"LastNameEn":             responsibleUser.LastNameEn,
				"StakeholderTitle":       stakeholderActionPlan.Stakeholder.TitleEn,
				"ActionPlanTitle":        stakeholderActionPlan.Title,
				"ActionPlanDescription":  stakeholderActionPlan.Description,
				"CommentedByFirstNameEn": assignedBy.FirstNameEn,
				"CommentedByLastNameEn":  assignedBy.LastNameEn,
				"AssignedOn":             time.Now().Format("2006-01-02"),
				"RedirectURL":            fmt.Sprintf("%sstakeholder/register/%d/action-plans?action_plan_id=%d", appURL, stakeholderID, id),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "stakeholder/templates/email/stakeholder_action_plan_comment.html",
				EmailData:            emailData,
			}

			// send notification and mail
			err = caps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "stakeholder_action_plan")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
			if stakeholderActionPlan.ReviewUserID != nil {
				// prepare notification
				urlParam := fmt.Sprintf(`{"stakeholder_register_id":%d, "action_plan_id":%d}`, stakeholderID, id)
				queryParam := ""
				entity := "stakeholder-action-plan"
				action := "comment-added"

				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserIDs:              []uint{*stakeholderActionPlan.ReviewUserID},
					ModuleIdentifier:     "stakeholder",
					NotificationType:     string(commonDTOs.AssignmentType),
					HeadingEn:            "Comment Added For Stakeholder Action Plan",
					HeadingAr:            "تمت إضافة التعليق لخطة عمل الامتثال",
					TitleEn:              stakeholderActionPlan.Title,
					TitleAr:              stakeholderActionPlan.Title,
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					QueryParam:           queryParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				responsibleUser, err := caps.userSvc.GetByID(*stakeholderActionPlan.ReviewUserID)
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
				subject := fmt.Sprintf("Comment Added for Stakeholder Action Plan:%s", stakeholderActionPlan.Title)
				appURL := utility.GetAppURL()
				emailData := map[string]interface{}{
					"FirstNameEn":            responsibleUser.FirstNameEn,
					"LastNameEn":             responsibleUser.LastNameEn,
					"StakeholderTitle":       stakeholderActionPlan.Stakeholder.TitleEn,
					"ActionPlanTitle":        stakeholderActionPlan.Title,
					"ActionPlanDescription":  stakeholderActionPlan.Description,
					"CommentedByFirstNameEn": assignedBy.FirstNameEn,
					"CommentedByLastNameEn":  assignedBy.LastNameEn,
					"AssignedOn":             time.Now().Format("2006-01-02"),
					"RedirectURL":            fmt.Sprintf("%sstakeholder/register/%d/action-plans?action_plan_id=%d", appURL, stakeholderID, id),
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              emailTo,
					Subject:              subject,
					TemplateRelativePath: "stakeholder/templates/email/stakeholder_action_plan_comment.html",
					EmailData:            emailData,
				}

				// send notification and mail
				err = caps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "stakeholder_action_plan")
				if err != nil {
					utility.LogError("Failed to send notifications", map[string]interface{}{
						"error": err.Error(),
					})
				}
			}
		}
		// // prepare notification
		// urlParam := fmt.Sprintf(`{"stakeholder_register_id":%d, "action_plan_id":%d}`, stakeholderID, id)
		// queryParam := ""
		// entity := "stakeholder-action-plan"
		// action := "responsible-user-assign"

		// dbNotificationDTO := &commonDTOs.DBNotificationDTO{
		// 	UserID:               StakeholderActionPlanFormDTO.ResponsibleUserID,
		// 	ModuleIdentifier:     "stakeholder",
		// 	NotificationType:     string(commonDTOs.AssignmentType),
		// 	HeadingEn:            "Stakeholder Action Plan Responsibility Assigned",
		// 	HeadingAr:            "تم تعيين مسؤولية خطة عمل الامتثال",
		// 	TitleEn:              StakeholderActionPlanFormDTO.Title,
		// 	TitleAr:              StakeholderActionPlanFormDTO.Title,
		// 	Entity:               &entity,
		// 	Action:               &action,
		// 	URLParam:             urlParam,
		// 	QueryParam:           queryParam,
		// 	IsSystemNotification: false,
		// 	CreatedByID:          &currentUserID,
		// }

		// responsibleUser, err := caps.userSvc.GetByID(StakeholderActionPlanFormDTO.ResponsibleUserID)
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
		// subject := fmt.Sprintf("stakeholder Action Plan/%s: Responsibility Assigned", StakeholderActionPlanFormDTO.Title)
		// appURL := utility.GetAppURL()
		// emailData := map[string]interface{}{
		// 	"FirstNameEn":               responsibleUser.FirstNameEn,
		// 	"LastNameEn":                responsibleUser.LastNameEn,
		// 	"StakeholderReferenceNumber": stakeholder.ReferenceNumber,
		// 	"StakeholderTitle":           stakeholder.TitleEn,
		// 	"ActionPlanTitle":           StakeholderActionPlanFormDTO.Title,
		// 	"ActionPlanDescription":     StakeholderActionPlanFormDTO.Description,
		// 	"TargetDate":                dates.TargetDate.Format("2006-01-02"),
		// 	"AssignedByFirstNameEn":     assignedBy.FirstNameEn,
		// 	"AssignedByLastNameEn":      assignedBy.LastNameEn,
		// 	"AssignedOn":                time.Now().Format("2006-01-02"),
		// 	"RedirectURL":               fmt.Sprintf("%sstakeholder/register/%d/action-plans?action_plan_id=%d", appURL, stakeholderID, id),
		// }

		// emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
		// 	EmailTo:              emailTo,
		// 	Subject:              subject,
		// 	TemplateRelativePath: "stakeholder/templates/email/stakeholder_action_plan_create_responsible_user_mail.html",
		// 	EmailData:            emailData,
		// }

		// // send notification and mail
		// err = caps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "stakeholder_action_plan")
		// if err != nil {
		// 	utility.LogError("Failed to send notifications", map[string]interface{}{
		// 		"error": err.Error(),
		// 	})
		// }
	}()
	return id, nil
}

func (caps *StakeholderActionPlanCommentService) Update(stakeholderID uint, actionPlanID uint, id uint, currentUserID uint, stakeholderActionPlanCommentFormDTO *organizationalContextDTOs.StakeholderActionPlanCommentFormDTO) error {
	tx := caps.repo.BeginTransaction()

	if _, err := caps.stakeholderActionPlanSvc.GetByID(stakeholderID, actionPlanID); err != nil {
		return err
	}

	existingComment, err := caps.GetByID(stakeholderID, actionPlanID, id)
	if err != nil {
		return err
	}

	oldComment, err := caps.GetFullDetailsByID(stakeholderID, actionPlanID, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingComment.Comment = stakeholderActionPlanCommentFormDTO.Comment
	err = caps.repo.Update(existingComment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	var newDocuments []*models.StakeholderActionPlanCommentDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.StakeholderActionPlanCommentDocument

	for _, doc := range stakeholderActionPlanCommentFormDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {

			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}

			url := utility.GenerateDocumentURL("stakeholder", "stakeholder-register", strconv.FormatUint(uint64(stakeholderID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

			document := &models.StakeholderActionPlanCommentDocument{
				StakeholderActionPlanCommentID: id,
				URL:                            url,
				Name:                           doc.FileName,
				Size:                           doc.Size,
				Token:                          encryptedFileName,
				OriginalName:                   doc.OriginalName,
				Type:                           doc.Type,
				CreatedBYID:                    currentUserID,
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
		storagePath := filepath.Join("stakeholder", "stakeholder-register", strconv.FormatUint(uint64(stakeholderID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("stakeholder", "stakeholder-register", strconv.FormatUint(uint64(stakeholderID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	go func() {
		existingComment, err := caps.GetFullDetailsByID(stakeholderID, actionPlanID, id)
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
				models.StakeholderActionPlanComment{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.StakeholderActionPlanComment{}.TableName(),
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

func (caps *StakeholderActionPlanCommentService) Delete(stakeholderID uint, actionPlanID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := caps.repo.BeginTransaction()
	if _, err := caps.stakeholderActionPlanSvc.GetByID(stakeholderID, actionPlanID); err != nil {
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
			"id":                 id,
			"stakeholder_action": actionPlanID,
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
				models.StakeholderActionPlanComment{}.TableName(),
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

func (caps *StakeholderActionPlanCommentService) GetDocumentPath(stakeholderID, actionPlanID, id uint, token string) (string, error) {
	existingComment, err := caps.GetByID(stakeholderID, actionPlanID, id)
	if err != nil {
		return "", err
	}
	if len(existingComment.Documents) == 0 {
		return "", fmt.Errorf("no documents found for the stakeholder action plan comment")
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	storagePath := filepath.Join("stakeholder", "stakeholder-register", strconv.FormatUint(uint64(stakeholderID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}
