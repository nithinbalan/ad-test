package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"jethur/config"
	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	"jethur/grc/business-continuity/repositories"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type TaeRecoveryActionPlanCommentServiceInterface interface {
	GetAll(testAndExerciseID uint, actionPlanID uint) ([]*businessContinuityDTOs.TaeRecoveryActionPlanCommentListDTO, error)
	GetByID(testAndExerciseID uint, actionPlanID uint, id uint) (*models.TaeRecoveryActionPlanComment, error)
	GetFullDetailsByID(testAndExerciseID uint, actionPlanID uint, id uint, tx ...utility.DatabaseTransaction) (*businessContinuityDTOs.TaeRecoveryActionPlanCommentDetailDTO, error)
	Create(testAndExerciseID uint, actionPlanID uint, currentUserID uint, dto *businessContinuityDTOs.TaeRecoveryActionPlanCommentFormDTO) (uint, error)
	Update(testAndExerciseID uint, actionPlanID uint, id uint, currentUserID uint, dto *businessContinuityDTOs.TaeRecoveryActionPlanCommentFormDTO) error
	Delete(testAndExerciseID uint, actionPlanID uint, ids []uint, currentUserID uint) ([]uint, error)

	GetDocumentPath(testAndExerciseID, actionPlanID, id uint, token string) (string, error)
}

type TaeRecoveryActionPlanCommentService struct {
	repo                                 repositories.TaeRecoveryActionPlanCommentRepositoryInterface
	userSvc                              organizationServices.UserServiceInterface
	activityLogSvc                       logServices.ActivityLogServiceInterface
	testAndExerciseRecoveryActionPlanSvc TestAndExerciseRecoveryActionPlanServiceInterface
	notificationSvc                      commonServices.NotificationServiceInterface
}

func NewTaeRecoveryActionPlanCommentService(repo repositories.TaeRecoveryActionPlanCommentRepositoryInterface, userSvc organizationServices.UserServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	testAndExerciseRecoveryActionPlanSvc TestAndExerciseRecoveryActionPlanServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) TaeRecoveryActionPlanCommentServiceInterface {
	return &TaeRecoveryActionPlanCommentService{
		repo:                                 repo,
		userSvc:                              userSvc,
		activityLogSvc:                       activityLogSvc,
		testAndExerciseRecoveryActionPlanSvc: testAndExerciseRecoveryActionPlanSvc,
		notificationSvc:                      notificationSvc,
	}
}

func (caps *TaeRecoveryActionPlanCommentService) GetAll(testAndExerciseID uint, actionPlanID uint) ([]*businessContinuityDTOs.TaeRecoveryActionPlanCommentListDTO, error) {
	if _, err := caps.testAndExerciseRecoveryActionPlanSvc.GetByID(testAndExerciseID, actionPlanID); err != nil {
		return nil, err
	}

	return caps.repo.GetAll(actionPlanID)
}

func (caps *TaeRecoveryActionPlanCommentService) GetByID(testAndExerciseID uint, actionPlanID uint, id uint) (*models.TaeRecoveryActionPlanComment, error) {
	if _, err := caps.testAndExerciseRecoveryActionPlanSvc.GetByID(testAndExerciseID, actionPlanID); err != nil {
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

func (caps *TaeRecoveryActionPlanCommentService) GetFullDetailsByID(testAndExerciseID uint, actionPlanID uint, id uint, tx ...utility.DatabaseTransaction) (*businessContinuityDTOs.TaeRecoveryActionPlanCommentDetailDTO, error) {
	if _, err := caps.testAndExerciseRecoveryActionPlanSvc.GetByID(testAndExerciseID, actionPlanID); err != nil {
		fmt.Printf("TaeRecoveryActionPlan not found for testAndExerciseID %d, actionPlanID %d\n", testAndExerciseID, actionPlanID)
		return nil, err
	}

	var commentDetail businessContinuityDTOs.TaeRecoveryActionPlanCommentDetailDTO

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

func (caps *TaeRecoveryActionPlanCommentService) Create(testAndExerciseID uint, actionPlanID uint, currentUserID uint, taeRecoveryActionPlanFormDTO *businessContinuityDTOs.TaeRecoveryActionPlanCommentFormDTO) (uint, error) {
	tx := caps.repo.BeginTransaction()

	if _, err := caps.testAndExerciseRecoveryActionPlanSvc.GetByID(testAndExerciseID, actionPlanID); err != nil {
		return 0, err
	}

	taeRecoveryActionPlan, err := caps.testAndExerciseRecoveryActionPlanSvc.GetFullDetailsByID(testAndExerciseID, actionPlanID)
	if err != nil {
		return 0, err
	}

	comment := &models.TaeRecoveryActionPlanComment{
		TestAndExerciseRecoveryActionPlanID: actionPlanID,
		Comment:                             taeRecoveryActionPlanFormDTO.Comment,
		CreatedBYID:                         currentUserID,
	}

	id, err := caps.repo.Create(comment)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.TaeRecoveryActionPlanCommentDocument

	for _, doc := range taeRecoveryActionPlanFormDTO.Documents {
		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL(
			"business-continuity", "test-and-exercise", strconv.FormatUint(uint64(testAndExerciseID), 10),
			"action-plans", strconv.FormatUint(uint64(actionPlanID), 10),
			"comments", strconv.FormatUint(uint64(id), 10),
			"documents",
		)

		document := &models.TaeRecoveryActionPlanCommentDocument{
			TaeRecoveryActionPlanCommentID: id,
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
				"business-continuity", "test-and-exercise", strconv.FormatUint(uint64(testAndExerciseID), 10),
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
		existingComment, err := caps.GetFullDetailsByID(testAndExerciseID, actionPlanID, id)
		if err != nil {
			utility.LogError("Failed to fetch test and exercise recovery action plan comment", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = caps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.TaeRecoveryActionPlanComment{}.TableName(),
			id,
			existingComment,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.TaeRecoveryActionPlanComment{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}

		if taeRecoveryActionPlan.ResponsibleUserID == currentUserID {
			// if responsible user added the comment, notification for review user
			if taeRecoveryActionPlan.ReviewUserID != nil {
				// prepare notification
				urlParam := fmt.Sprintf(`{"test_and_exercise_id":%d, "action_plan_id":%d}`, testAndExerciseID, id)
				queryParam := ""
				entity := "recovery-action-plan"
				action := "comment-added"

				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserIDs:              []uint{*taeRecoveryActionPlan.ReviewUserID},
					ModuleIdentifier:     "business_continuity",
					NotificationType:     string(commonDTOs.AssignmentType),
					HeadingEn:            "Comment Added For TaeRecovery Action Plan",
					HeadingAr:            "تمت إضافة التعليق لخطة عمل الامتثال",
					TitleEn:              taeRecoveryActionPlan.Title,
					TitleAr:              taeRecoveryActionPlan.Title,
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					QueryParam:           queryParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				reviewUser, err := caps.userSvc.GetByID(*taeRecoveryActionPlan.ReviewUserID)
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
				subject := fmt.Sprintf("Comment Added for TaeRecovery Action Plan:%s", taeRecoveryActionPlan.Title)
				appURL := utility.GetAppURL()
				emailData := map[string]interface{}{
					"FirstNameEn":            reviewUser.FirstNameEn,
					"LastNameEn":             reviewUser.LastNameEn,
					"ReferenceNumber":        taeRecoveryActionPlan.TestAndExercise.ReferenceNumber,
					"TestAndExerciseTitle":   taeRecoveryActionPlan.TestAndExercise.BcpTestAndExercise.Title,
					"ActionPlanTitle":        taeRecoveryActionPlan.Title,
					"CommentedByFirstNameEn": assignedBy.FirstNameEn,
					"CommentedByLastNameEn":  assignedBy.LastNameEn,
					"CommentedOn":            time.Now().Format("2006-01-02"),
					"RedirectURL":            fmt.Sprintf("%sbusiness-continuity/test-and-exercise/%d/action-plans?action_plan_id=%d", appURL, testAndExerciseID, id),
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              emailTo,
					Subject:              subject,
					TemplateRelativePath: "business-continuity/templates/email/action_plan_comment.html",
					EmailData:            emailData,
				}

				// send notification and mail
				err = caps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "test_and_exercise")
				if err != nil {
					utility.LogError("Failed to send notifications", map[string]interface{}{
						"error": err.Error(),
					})
				}
			}
		} else if taeRecoveryActionPlan.ReviewUserID == &currentUserID {
			// if review user added the comment, notification for responsible user
			// prepare notification
			urlParam := fmt.Sprintf(`{"test_and_exercise_id":%d, "action_plan_id":%d}`, testAndExerciseID, id)
			queryParam := ""
			entity := "recovery-action-plan"
			action := "comment-added"

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{taeRecoveryActionPlan.ResponsibleUserID},
				ModuleIdentifier:     "business_continuity",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Comment Added For TaeRecovery Action Plan",
				HeadingAr:            "تمت إضافة التعليق لخطة عمل الامتثال",
				TitleEn:              taeRecoveryActionPlan.Title,
				TitleAr:              taeRecoveryActionPlan.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			responsibleUser, err := caps.userSvc.GetByID(taeRecoveryActionPlan.ResponsibleUserID)
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
			subject := fmt.Sprintf("Comment Added for TaeRecovery Action Plan:%s", taeRecoveryActionPlan.Title)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":            responsibleUser.FirstNameEn,
				"LastNameEn":             responsibleUser.LastNameEn,
				"ReferenceNumber":        taeRecoveryActionPlan.TestAndExercise.ReferenceNumber,
				"TestAndExerciseTitle":   taeRecoveryActionPlan.TestAndExercise.BcpTestAndExercise.Title,
				"ActionPlanTitle":        taeRecoveryActionPlan.Title,
				"CommentedByFirstNameEn": assignedBy.FirstNameEn,
				"CommentedByLastNameEn":  assignedBy.LastNameEn,
				"AssignedOn":             time.Now().Format("2006-01-02"),
				"RedirectURL":            fmt.Sprintf("%sbusiness-continuity/test-and-exercise/%d/action-plans?action_plan_id=%d", appURL, testAndExerciseID, id),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "business-continuity/templates/email/action_plan_comment.html",
				EmailData:            emailData,
			}

			// send notification and mail
			err = caps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "test_and_exercise")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		} else {
			// if other user added the comment, notification for review user and responsible user
			// prepare notification
			urlParam := fmt.Sprintf(`{"test_and_exercise_id":%d, "action_plan_id":%d}`, testAndExerciseID, id)
			queryParam := ""
			entity := "recovery-action-plan"
			action := "comment-added"

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{taeRecoveryActionPlan.ResponsibleUserID},
				ModuleIdentifier:     "business_continuity",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Comment Added For TaeRecovery Action Plan",
				HeadingAr:            "تمت إضافة التعليق لخطة عمل الامتثال",
				TitleEn:              taeRecoveryActionPlan.Title,
				TitleAr:              taeRecoveryActionPlan.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			responsibleUser, err := caps.userSvc.GetByID(taeRecoveryActionPlan.ResponsibleUserID)
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
			subject := fmt.Sprintf("Comment Added for TaeRecovery Action Plan:%s", taeRecoveryActionPlan.Title)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":            responsibleUser.FirstNameEn,
				"LastNameEn":             responsibleUser.LastNameEn,
				"ReferenceNumber":        taeRecoveryActionPlan.TestAndExercise.ReferenceNumber,
				"TestAndExerciseTitle":   taeRecoveryActionPlan.TestAndExercise.BcpTestAndExercise.Title,
				"ActionPlanTitle":        taeRecoveryActionPlan.Title,
				"CommentedByFirstNameEn": assignedBy.FirstNameEn,
				"CommentedByLastNameEn":  assignedBy.LastNameEn,
				"AssignedOn":             time.Now().Format("2006-01-02"),
				"RedirectURL":            fmt.Sprintf("%sbusiness-continuity/test-and-exercise/%d/action-plans?action_plan_id=%d", appURL, testAndExerciseID, id),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "business-continuity/templates/email/action_plan_comment.html",
				EmailData:            emailData,
			}

			// send notification and mail
			err = caps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "test_and_exercise")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
			if taeRecoveryActionPlan.ReviewUserID != nil {
				// prepare notification
				urlParam := fmt.Sprintf(`{"test_and_exercise_id":%d, "action_plan_id":%d}`, testAndExerciseID, id)
				queryParam := ""
				entity := "recovery-action-plan"
				action := "comment-added"

				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserIDs:              []uint{*taeRecoveryActionPlan.ReviewUserID},
					ModuleIdentifier:     "business_continuity",
					NotificationType:     string(commonDTOs.AssignmentType),
					HeadingEn:            "Comment Added For TaeRecovery Action Plan",
					HeadingAr:            "تمت إضافة التعليق لخطة عمل الامتثال",
					TitleEn:              taeRecoveryActionPlan.Title,
					TitleAr:              taeRecoveryActionPlan.Title,
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					QueryParam:           queryParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				responsibleUser, err := caps.userSvc.GetByID(*taeRecoveryActionPlan.ReviewUserID)
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
				subject := fmt.Sprintf("Comment Added for TaeRecovery Action Plan:%s", taeRecoveryActionPlan.Title)
				appURL := utility.GetAppURL()
				emailData := map[string]interface{}{
					"FirstNameEn":            responsibleUser.FirstNameEn,
					"LastNameEn":             responsibleUser.LastNameEn,
					"ReferenceNumber":        taeRecoveryActionPlan.TestAndExercise.ReferenceNumber,
					"TestAndExerciseTitle":   taeRecoveryActionPlan.TestAndExercise.BcpTestAndExercise.Title,
					"ActionPlanTitle":        taeRecoveryActionPlan.Title,
					"CommentedByFirstNameEn": assignedBy.FirstNameEn,
					"CommentedByLastNameEn":  assignedBy.LastNameEn,
					"AssignedOn":             time.Now().Format("2006-01-02"),
					"RedirectURL":            fmt.Sprintf("%sbusiness-continuity/test-and-exercise/%d/action-plans?action_plan_id=%d", appURL, testAndExerciseID, id),
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              emailTo,
					Subject:              subject,
					TemplateRelativePath: "business-continuity/templates/email/action_plan_comment.html",
					EmailData:            emailData,
				}

				// send notification and mail
				err = caps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "test_and_exercise")
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

func (caps *TaeRecoveryActionPlanCommentService) Update(testAndExerciseID uint, actionPlanID uint, id uint, currentUserID uint, taeRecoveryActionPlanFormDTO *businessContinuityDTOs.TaeRecoveryActionPlanCommentFormDTO) error {
	tx := caps.repo.BeginTransaction()

	if _, err := caps.testAndExerciseRecoveryActionPlanSvc.GetByID(testAndExerciseID, actionPlanID); err != nil {
		return err
	}

	existingComment, err := caps.GetByID(testAndExerciseID, actionPlanID, id)
	if err != nil {
		return err
	}

	oldComment, err := caps.GetFullDetailsByID(testAndExerciseID, actionPlanID, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingComment.Comment = taeRecoveryActionPlanFormDTO.Comment
	err = caps.repo.Update(existingComment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	var newDocuments []*models.TaeRecoveryActionPlanCommentDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.TaeRecoveryActionPlanCommentDocument

	for _, doc := range taeRecoveryActionPlanFormDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {

			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}

			url := utility.GenerateDocumentURL("business-continuity", "test-and-exercise", strconv.FormatUint(uint64(testAndExerciseID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

			document := &models.TaeRecoveryActionPlanCommentDocument{
				TaeRecoveryActionPlanCommentID: id,
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
		storagePath := filepath.Join("business-continuity", "test-and-exercise", strconv.FormatUint(uint64(testAndExerciseID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("business-continuity", "test-and-exercise", strconv.FormatUint(uint64(testAndExerciseID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	go func() {
		existingComment, err := caps.GetFullDetailsByID(testAndExerciseID, actionPlanID, id)
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
				models.TaeRecoveryActionPlanComment{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.TaeRecoveryActionPlanComment{}.TableName(),
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

func (caps *TaeRecoveryActionPlanCommentService) Delete(testAndExerciseID uint, actionPlanID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := caps.repo.BeginTransaction()
	if _, err := caps.testAndExerciseRecoveryActionPlanSvc.GetByID(testAndExerciseID, actionPlanID); err != nil {
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
			"id":                  id,
			"tae_recovery_action": actionPlanID,
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
				models.TaeRecoveryActionPlanComment{}.TableName(),
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

func (caps *TaeRecoveryActionPlanCommentService) GetDocumentPath(testAndExerciseID, actionPlanID, id uint, token string) (string, error) {
	existingComment, err := caps.GetByID(testAndExerciseID, actionPlanID, id)
	if err != nil {
		return "", err
	}
	if len(existingComment.Documents) == 0 {
		return "", fmt.Errorf("no documents found for the test and exercise recovery action plan comment")
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	storagePath := filepath.Join("business-continuity", "test-and-exercise", strconv.FormatUint(uint64(testAndExerciseID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}
