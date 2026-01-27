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
	violationReportDTOs "jethur/grc/violation-report/dtos"
	"jethur/grc/violation-report/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ViolationReportActionPlanCommentServiceInterface interface {
	GetAll(violationReportID uint, actionPlanID uint) ([]*violationReportDTOs.ViolationReportActionPlanCommentListDTO, error)
	GetByID(violationReportID uint, actionPlanID uint, id uint) (*models.ViolationReportActionPlanComment, error)
	GetFullDetailsByID(violationReportID uint, actionPlanID uint, id uint, tx ...utility.DatabaseTransaction) (*violationReportDTOs.ViolationReportActionPlanCommentDetailDTO, error)
	Create(violationReportID uint, actionPlanID uint, currentUserID uint, dto *violationReportDTOs.ViolationReportActionPlanCommentFormDTO) (uint, error)
	Update(violationReportID uint, actionPlanID uint, id uint, currentUserID uint, dto *violationReportDTOs.ViolationReportActionPlanCommentFormDTO) error
	Delete(violationReportID uint, actionPlanID uint, ids []uint, currentUserID uint) ([]uint, error)

	GetDocumentPath(violationReportID, actionPlanID, id uint, token string) (string, error)
}

type ViolationReportActionPlanCommentService struct {
	repo                         repositories.ViolationReportActionPlanCommentRepositoryInterface
	userSvc                      organizationServices.UserServiceInterface
	activityLogSvc               logServices.ActivityLogServiceInterface
	violationReportActionPlanSvc ViolationReportActionPlanServiceInterface
	notificationSvc              commonServices.NotificationServiceInterface
}

func NewViolationReportActionPlanCommentService(repo repositories.ViolationReportActionPlanCommentRepositoryInterface, userSvc organizationServices.UserServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	violationReportActionPlanSvc ViolationReportActionPlanServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) ViolationReportActionPlanCommentServiceInterface {
	return &ViolationReportActionPlanCommentService{
		repo:                         repo,
		userSvc:                      userSvc,
		activityLogSvc:               activityLogSvc,
		violationReportActionPlanSvc: violationReportActionPlanSvc,
		notificationSvc:              notificationSvc,
	}
}

func (vrapc *ViolationReportActionPlanCommentService) GetAll(violationReportID uint, actionPlanID uint) ([]*violationReportDTOs.ViolationReportActionPlanCommentListDTO, error) {
	if _, err := vrapc.violationReportActionPlanSvc.GetByID(violationReportID, actionPlanID); err != nil {
		return nil, err
	}

	return vrapc.repo.GetAll(actionPlanID)
}

func (vrapc *ViolationReportActionPlanCommentService) GetByID(violationReportID uint, actionPlanID uint, id uint) (*models.ViolationReportActionPlanComment, error) {
	if _, err := vrapc.violationReportActionPlanSvc.GetByID(violationReportID, actionPlanID); err != nil {
		return nil, err
	}

	comment, err := vrapc.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return comment, nil
}

func (vrapc *ViolationReportActionPlanCommentService) GetFullDetailsByID(violationReportID uint, actionPlanID uint, id uint, tx ...utility.DatabaseTransaction) (*violationReportDTOs.ViolationReportActionPlanCommentDetailDTO, error) {
	if _, err := vrapc.violationReportActionPlanSvc.GetByID(violationReportID, actionPlanID); err != nil {
		fmt.Printf("ViolationReportActionPlan not found for violationReportID %d, actionPlanID %d\n", violationReportID, actionPlanID)
		return nil, err
	}

	var commentDetail violationReportDTOs.ViolationReportActionPlanCommentDetailDTO

	comment, err := vrapc.repo.GetFullDetailsByID(id, tx...)
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

func (vrapc *ViolationReportActionPlanCommentService) Create(violationReportID uint, actionPlanID uint, currentUserID uint, violationReportActionPlanCommentFormDTO *violationReportDTOs.ViolationReportActionPlanCommentFormDTO) (uint, error) {
	tx := vrapc.repo.BeginTransaction()

	if _, err := vrapc.violationReportActionPlanSvc.GetByID(violationReportID, actionPlanID); err != nil {
		return 0, err
	}

	violationReportActionPlan, err := vrapc.violationReportActionPlanSvc.GetFullDetailsByID(violationReportID, actionPlanID)
	if err != nil {
		return 0, err
	}

	comment := &models.ViolationReportActionPlanComment{
		ViolationReportActionPlanID: actionPlanID,
		Comment:                     violationReportActionPlanCommentFormDTO.Comment,
		CreatedBYID:                 currentUserID,
	}

	id, err := vrapc.repo.Create(comment)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.ViolationReportActionPlanCommentDocument

	for _, doc := range violationReportActionPlanCommentFormDTO.Documents {
		// Skip empty or invalid documents
		if doc.FileName == "" || doc.Size <= 0 {
			continue
		}

		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL(
			"violation-report", "violation-report-register", strconv.FormatUint(uint64(violationReportID), 10),
			"action-plans", strconv.FormatUint(uint64(actionPlanID), 10),
			"comments", strconv.FormatUint(uint64(id), 10),
			"documents",
		)

		document := &models.ViolationReportActionPlanCommentDocument{
			ViolationReportActionPlanCommentID: id,
			URL:                                url,
			Name:                               doc.FileName,
			Size:                               doc.Size,
			Token:                              encryptedFilename,
			OriginalName:                       doc.OriginalName,
			Type:                               doc.Type,
			CreatedBYID:                        currentUserID,
		}

		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := vrapc.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}

		for _, doc := range documents {
			storagePath := filepath.Join(
				"violation-report", "violation-report-register", strconv.FormatUint(uint64(violationReportID), 10),
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
		existingComment, err := vrapc.GetFullDetailsByID(violationReportID, actionPlanID, id)
		if err != nil {
			utility.LogError("Failed to fetch violation report action plan comment", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = vrapc.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ViolationReportActionPlanComment{}.TableName(),
			id,
			existingComment,
			currentUserID,
			"violation_report", "violation_report_register", "action_plan_comment", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ViolationReportActionPlanComment{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
		if uint(violationReportActionPlan.ResponsibleUser.ID) == currentUserID {
			// if responsible user added the comment, notification for review user
			if violationReportActionPlan.ReviewUser != nil && violationReportActionPlan.ReviewUser.ID != 0 {

				// prepare notification
				urlParam := fmt.Sprintf(`{"violation_report_id":%d, "action_plan_id":%d}`, violationReportID, id)
				queryParam := ""
				entity := "violation-report-action-plan"
				action := "comment-added"

				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserIDs:              []uint{uint(violationReportActionPlan.ReviewUser.ID)},
					ModuleIdentifier:     "violation_report",
					NotificationType:     string(commonDTOs.AssignmentType),
					HeadingEn:            "Comment Added For Violation Report Action Plan",
					HeadingAr:            "تمت إضافة التعليق لخطة عمل تقرير المخالفة",
					TitleEn:              violationReportActionPlan.Title,
					TitleAr:              violationReportActionPlan.Title,
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					QueryParam:           queryParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				reviewUser, err := vrapc.userSvc.GetByID(uint(violationReportActionPlan.ReviewUser.ID))
				if err != nil {
					utility.LogError("Failed to fetch review user", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}
				assignedBy, err := vrapc.userSvc.GetByID(currentUserID)
				if err != nil {
					utility.LogError("Failed to fetch commented user", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}

				emailTo := []string{reviewUser.Email}
				subject := fmt.Sprintf("Comment Added for Violation Report Action Plan:%s", violationReportActionPlan.Title)
				appURL := utility.GetAppURL()
				emailData := map[string]interface{}{
					"FirstNameEn":                    reviewUser.FirstNameEn,
					"LastNameEn":                     reviewUser.LastNameEn,
					"ViolationReportTitle":           violationReportActionPlan.ViolationReport.TitleEn,
					"ViolationReportReferenceNumber": violationReportActionPlan.ViolationReport.ReferenceNumber,
					"ActionPlanTitle":                violationReportActionPlan.Title,
					"ActionPlanDescription":          violationReportActionPlan.Description,
					"CommentedByFirstNameEn":         assignedBy.FirstNameEn,
					"CommentedByLastNameEn":          assignedBy.LastNameEn,
					"CommentedOn":                    time.Now().Format("2006-01-02"),
					"RedirectURL":                    fmt.Sprintf("%s/violation-report/register/%d/action-plans?action_plan_id=%d", appURL, violationReportID, id),
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              emailTo,
					Subject:              subject,
					TemplateRelativePath: "violation-report/templates/email/violation_report_action_plan_comment.html",
					EmailData:            emailData,
				}

				// send notification and mail
				err = vrapc.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "violation_report_action_plan")
				if err != nil {
					utility.LogError("Failed to send notifications", map[string]interface{}{
						"error": err.Error(),
					})
				}
			}
		} else if violationReportActionPlan.ReviewUser != nil && uint(violationReportActionPlan.ReviewUser.ID) == currentUserID {
			// if review user added the comment, notification for responsible user
			// prepare notification
			urlParam := fmt.Sprintf(`{"violation_report_id":%d, "action_plan_id":%d}`, violationReportID, id)
			queryParam := ""
			entity := "violation-report-action-plan"
			action := "comment-added"

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{uint(violationReportActionPlan.ResponsibleUser.ID)},
				ModuleIdentifier:     "violation_report",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Comment Added For Violation Report Action Plan",
				HeadingAr:            "تمت إضافة التعليق لخطة عمل تقرير المخالفة",
				TitleEn:              violationReportActionPlan.Title,
				TitleAr:              violationReportActionPlan.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			responsibleUser, err := vrapc.userSvc.GetByID(uint(violationReportActionPlan.ResponsibleUser.ID))
			if err != nil {
				utility.LogError("Failed to fetch responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}
			assignedBy, err := vrapc.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch commented user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			emailTo := []string{responsibleUser.Email}
			subject := fmt.Sprintf("Comment Added for Violation Report Action Plan:%s", violationReportActionPlan.Title)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":                    responsibleUser.FirstNameEn,
				"LastNameEn":                     responsibleUser.LastNameEn,
				"ViolationReportTitle":           violationReportActionPlan.ViolationReport.TitleEn,
				"ViolationReportReferenceNumber": violationReportActionPlan.ViolationReport.ReferenceNumber,
				"ActionPlanTitle":                violationReportActionPlan.Title,
				"ActionPlanDescription":          violationReportActionPlan.Description,
				"CommentedByFirstNameEn":         assignedBy.FirstNameEn,
				"CommentedByLastNameEn":          assignedBy.LastNameEn,
				"AssignedOn":                     time.Now().Format("2006-01-02"),
				"RedirectURL":                    fmt.Sprintf("%s/violation-report/register/%d/action-plans?action_plan_id=%d", appURL, violationReportID, id),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "violation-report/templates/email/violation_report_action_plan_comment.html",
				EmailData:            emailData,
			}

			// send notification and mail
			err = vrapc.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "violation_report_action_plan")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		} else {
			// if other user added the comment, notification for review user and responsible user
			// prepare notification
			urlParam := fmt.Sprintf(`{"violation_report_id":%d, "action_plan_id":%d}`, violationReportID, id)
			queryParam := ""
			entity := "violation-report-action-plan"
			action := "comment-added"

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{uint(violationReportActionPlan.ResponsibleUser.ID)},
				ModuleIdentifier:     "violation_report",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Comment Added For Violation Report Action Plan",
				HeadingAr:            "تمت إضافة التعليق لخطة عمل تقرير المخالفة",
				TitleEn:              violationReportActionPlan.Title,
				TitleAr:              violationReportActionPlan.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			responsibleUser, err := vrapc.userSvc.GetByID(uint(violationReportActionPlan.ResponsibleUser.ID))
			if err != nil {
				utility.LogError("Failed to fetch responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}
			assignedBy, err := vrapc.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch commented user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			emailTo := []string{responsibleUser.Email}
			subject := fmt.Sprintf("Comment Added for Violation Report Action Plan:%s", violationReportActionPlan.Title)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":                    responsibleUser.FirstNameEn,
				"LastNameEn":                     responsibleUser.LastNameEn,
				"ViolationReportTitle":           violationReportActionPlan.ViolationReport.TitleEn,
				"ViolationReportReferenceNumber": violationReportActionPlan.ViolationReport.ReferenceNumber,

				"ActionPlanTitle":        violationReportActionPlan.Title,
				"ActionPlanDescription":  violationReportActionPlan.Description,
				"CommentedByFirstNameEn": assignedBy.FirstNameEn,
				"CommentedByLastNameEn":  assignedBy.LastNameEn,
				"AssignedOn":             time.Now().Format("2006-01-02"),
				"RedirectURL":            fmt.Sprintf("%s/violation-report/register/%d/action-plans?action_plan_id=%d", appURL, violationReportID, id),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "violation-report/templates/email/violation_report_action_plan_comment.html",
				EmailData:            emailData,
			}

			// send notification and mail
			err = vrapc.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "violation_report_action_plan")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
			if violationReportActionPlan.ReviewUser != nil && violationReportActionPlan.ReviewUser.ID != 0 {
				// prepare notification
				urlParam := fmt.Sprintf(`{"violation_report_id":%d, "action_plan_id":%d}`, violationReportID, id)
				queryParam := ""
				entity := "violation-report-action-plan"
				action := "comment-added"

				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserIDs:              []uint{uint(violationReportActionPlan.ReviewUser.ID)},
					ModuleIdentifier:     "violation_report",
					NotificationType:     string(commonDTOs.AssignmentType),
					HeadingEn:            "Comment Added For Violation Report Action Plan",
					HeadingAr:            "تمت إضافة التعليق لخطة عمل تقرير المخالفة",
					TitleEn:              violationReportActionPlan.Title,
					TitleAr:              violationReportActionPlan.Title,
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					QueryParam:           queryParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				responsibleUser, err := vrapc.userSvc.GetByID(uint(violationReportActionPlan.ReviewUser.ID))
				if err != nil {
					utility.LogError("Failed to fetch review user", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}
				assignedBy, err := vrapc.userSvc.GetByID(currentUserID)
				if err != nil {
					utility.LogError("Failed to fetch commented user", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}

				emailTo := []string{responsibleUser.Email}
				subject := fmt.Sprintf("Comment Added for Violation Report Action Plan:%s", violationReportActionPlan.Title)
				appURL := utility.GetAppURL()
				emailData := map[string]interface{}{
					"FirstNameEn":                    responsibleUser.FirstNameEn,
					"LastNameEn":                     responsibleUser.LastNameEn,
					"ViolationReportTitle":           violationReportActionPlan.ViolationReport.TitleEn,
					"ViolationReportReferenceNumber": violationReportActionPlan.ViolationReport.ReferenceNumber,

					"ActionPlanTitle":        violationReportActionPlan.Title,
					"ActionPlanDescription":  violationReportActionPlan.Description,
					"CommentedByFirstNameEn": assignedBy.FirstNameEn,
					"CommentedByLastNameEn":  assignedBy.LastNameEn,
					"AssignedOn":             time.Now().Format("2006-01-02"),
					"RedirectURL":            fmt.Sprintf("%s/violation-report/register/%d/action-plans?action_plan_id=%d", appURL, violationReportID, id),
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              emailTo,
					Subject:              subject,
					TemplateRelativePath: "violation-report/templates/email/violation_report_action_plan_comment.html",
					EmailData:            emailData,
				}

				// send notification and mail
				err = vrapc.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "violation_report_action_plan")
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

func (vrapc *ViolationReportActionPlanCommentService) Update(violationReportID uint, actionPlanID uint, id uint, currentUserID uint, violationReportActionPlanCommentFormDTO *violationReportDTOs.ViolationReportActionPlanCommentFormDTO) error {
	tx := vrapc.repo.BeginTransaction()

	if _, err := vrapc.violationReportActionPlanSvc.GetByID(violationReportID, actionPlanID); err != nil {
		return err
	}

	existingComment, err := vrapc.GetByID(violationReportID, actionPlanID, id)
	if err != nil {
		return err
	}

	oldComment, err := vrapc.GetFullDetailsByID(violationReportID, actionPlanID, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingComment.Comment = violationReportActionPlanCommentFormDTO.Comment
	err = vrapc.repo.Update(existingComment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	var newDocuments []*models.ViolationReportActionPlanCommentDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.ViolationReportActionPlanCommentDocument

	for _, doc := range violationReportActionPlanCommentFormDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {
			// Skip empty or invalid documents
			if doc.FileName == "" || doc.Size <= 0 {
				continue
			}

			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}

			url := utility.GenerateDocumentURL("violation-report", "violation-report-register", strconv.FormatUint(uint64(violationReportID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

			document := &models.ViolationReportActionPlanCommentDocument{
				ViolationReportActionPlanCommentID: id,
				URL:                                url,
				Name:                               doc.FileName,
				Size:                               doc.Size,
				Token:                              encryptedFileName,
				OriginalName:                       doc.OriginalName,
				Type:                               doc.Type,
				CreatedBYID:                        currentUserID,
			}

			newDocuments = append(newDocuments, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			if doc.ID != nil {
				existingDocument, err := vrapc.repo.GetDocumentByID(*doc.ID)
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
		if err := vrapc.repo.CreateDocuments(newDocuments, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIds) > 0 {
		if err := vrapc.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range newDocuments {
		storagePath := filepath.Join("violation-report", "violation-report-register", strconv.FormatUint(uint64(violationReportID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("violation-report", "violation-report-register", strconv.FormatUint(uint64(violationReportID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	go func() {
		existingComment, err := vrapc.GetFullDetailsByID(violationReportID, actionPlanID, id)
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
			err := vrapc.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.ViolationReportActionPlanComment{}.TableName(),
				id,
				changedData,
				currentUserID,
				"violation_report", "violation_report_register", "action_plan_comment", nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.ViolationReportActionPlanComment{}.TableName(),
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

func (vrapc *ViolationReportActionPlanCommentService) Delete(violationReportID uint, actionPlanID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := vrapc.repo.BeginTransaction()
	if _, err := vrapc.violationReportActionPlanSvc.GetByID(violationReportID, actionPlanID); err != nil {
		return nil, err
	}

	var failedIDs []uint
	var successfulIDs []uint
	var deletedComments []map[string]interface{}
	for _, id := range ids {
		if _, err := vrapc.repo.GetFullDetailsByID(id); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		deletedComments = append(deletedComments, map[string]interface{}{
			"id":                              id,
			"violation_report_action_plan_id": actionPlanID,
		})

		err := vrapc.repo.Delete(id, tx)
		if err != nil {
			failedIDs = append(failedIDs, id)
		} else {
			successfulIDs = append(successfulIDs, id)
		}
	}

	if len(failedIDs) > 0 {
		tx.Rollback()
		return failedIDs, nil
	}

	tx.Commit()
	if len(successfulIDs) > 0 {
		go func() {
			err := vrapc.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.ViolationReportActionPlanComment{}.TableName(),
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

func (vrapc *ViolationReportActionPlanCommentService) GetDocumentPath(violationReportID, actionPlanID, id uint, token string) (string, error) {
	existingComment, err := vrapc.GetByID(violationReportID, actionPlanID, id)
	if err != nil {
		return "", err
	}
	if len(existingComment.Documents) == 0 {
		return "", fmt.Errorf("no documents found for the violation report action plan comment")
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	storagePath := filepath.Join("violation-report", "violation-report-register", strconv.FormatUint(uint64(violationReportID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}
