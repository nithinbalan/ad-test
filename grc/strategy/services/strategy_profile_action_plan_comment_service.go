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
	strategyDTOs "jethur/grc/strategy/dtos"
	"jethur/grc/strategy/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type StrategyProfileActionPlanCommentServiceInterface interface {
	GetAll(strategyProfileID uint, actionPlanID uint) ([]*strategyDTOs.StrategyProfileActionPlanCommentListDTO, error)
	GetByID(strategyProfileID uint, actionPlanID uint, id uint) (*models.StrategyProfileActionPlanComment, error)
	GetFullDetailsByID(strategyProfileID uint, actionPlanID uint, id uint, tx ...utility.DatabaseTransaction) (*strategyDTOs.StrategyProfileActionPlanCommentDetailDTO, error)
	Create(strategyProfileID uint, actionPlanID uint, currentUserID uint, dto *strategyDTOs.StrategyProfileActionPlanCommentFormDTO) (uint, error)
	Update(strategyProfileID uint, actionPlanID uint, id uint, currentUserID uint, dto *strategyDTOs.StrategyProfileActionPlanCommentFormDTO) error
	Delete(strategyProfileID uint, actionPlanID uint, ids []uint, currentUserID uint) ([]uint, error)

	GetDocumentPath(strategyProfileID, actionPlanID, id uint, token string) (string, error)
}

type StrategyProfileActionPlanCommentService struct {
	repo                         repositories.StrategyProfileActionPlanCommentRepositoryInterface
	userSvc                      organizationServices.UserServiceInterface
	activityLogSvc               logServices.ActivityLogServiceInterface
	strategyProfileActionPlanSvc StrategyProfileActionPlanServiceInterface
	notificationSvc              commonServices.NotificationServiceInterface
}

func NewStrategyProfileActionPlanCommentService(repo repositories.StrategyProfileActionPlanCommentRepositoryInterface, userSvc organizationServices.UserServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	strategyProfileActionPlanSvc StrategyProfileActionPlanServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) StrategyProfileActionPlanCommentServiceInterface {
	return &StrategyProfileActionPlanCommentService{
		repo:                         repo,
		userSvc:                      userSvc,
		activityLogSvc:               activityLogSvc,
		strategyProfileActionPlanSvc: strategyProfileActionPlanSvc,
		notificationSvc:              notificationSvc,
	}
}

func (spaps *StrategyProfileActionPlanCommentService) GetAll(strategyProfileID uint, actionPlanID uint) ([]*strategyDTOs.StrategyProfileActionPlanCommentListDTO, error) {
	if _, err := spaps.strategyProfileActionPlanSvc.GetByID(strategyProfileID, actionPlanID); err != nil {
		return nil, err
	}

	return spaps.repo.GetAll(actionPlanID)
}

func (spaps *StrategyProfileActionPlanCommentService) GetByID(strategyProfileID uint, actionPlanID uint, id uint) (*models.StrategyProfileActionPlanComment, error) {
	if _, err := spaps.strategyProfileActionPlanSvc.GetByID(strategyProfileID, actionPlanID); err != nil {
		return nil, err
	}

	comment, err := spaps.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return comment, nil
}

func (spaps *StrategyProfileActionPlanCommentService) GetFullDetailsByID(strategyProfileID uint, actionPlanID uint, id uint, tx ...utility.DatabaseTransaction) (*strategyDTOs.StrategyProfileActionPlanCommentDetailDTO, error) {
	if _, err := spaps.strategyProfileActionPlanSvc.GetByID(strategyProfileID, actionPlanID); err != nil {
		fmt.Printf("Strategy Profile Action Plan not found for strategyProfileID %d, actionPlanID %d\n", strategyProfileID, actionPlanID)
		return nil, err
	}

	var commentDetail strategyDTOs.StrategyProfileActionPlanCommentDetailDTO

	comment, err := spaps.repo.GetFullDetailsByID(id, tx...)
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

func (spaps *StrategyProfileActionPlanCommentService) Create(strategyProfileID uint, actionPlanID uint, currentUserID uint, strategyProfileActionPlanFormDTO *strategyDTOs.StrategyProfileActionPlanCommentFormDTO) (uint, error) {
	tx := spaps.repo.BeginTransaction()

	if _, err := spaps.strategyProfileActionPlanSvc.GetByID(strategyProfileID, actionPlanID); err != nil {
		return 0, err
	}

	strategyProfileActionPlan, err := spaps.strategyProfileActionPlanSvc.GetFullDetailsByID(strategyProfileID, actionPlanID)
	if err != nil {
		return 0, err
	}

	comment := &models.StrategyProfileActionPlanComment{
		StrategyProfileActionPlanID: actionPlanID,
		Comment:                     strategyProfileActionPlanFormDTO.Comment,
		CreatedBYID:                 currentUserID,
	}

	id, err := spaps.repo.Create(comment)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.StrategyProfileActionPlanCommentDocument

	for _, doc := range strategyProfileActionPlanFormDTO.Documents {
		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL(
			"strategy", "strategy-profile", strconv.FormatUint(uint64(strategyProfileID), 10),
			"action-plans", strconv.FormatUint(uint64(actionPlanID), 10),
			"comments", strconv.FormatUint(uint64(id), 10),
			"documents",
		)

		document := &models.StrategyProfileActionPlanCommentDocument{
			StrategyProfileActionPlanCommentID: id,
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
		if err := spaps.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}

		for _, doc := range documents {
			storagePath := filepath.Join(
				"strategy", "strategy-profile", strconv.FormatUint(uint64(strategyProfileID), 10),
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
		existingComment, err := spaps.GetFullDetailsByID(strategyProfileID, actionPlanID, id)
		if err != nil {
			utility.LogError("Failed to fetch strategy action plan comment", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = spaps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.StrategyProfileActionPlanComment{}.TableName(),
			id,
			existingComment,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.StrategyProfileActionPlanComment{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}

		if strategyProfileActionPlan.ResponsibleUserID == &currentUserID {
			// if responsible user added the comment, notification for review user
			if strategyProfileActionPlan.ReviewUserID != nil {
				// prepare notification
				urlParam := fmt.Sprintf(`{"strategy_profile_id":%d, "action_plan_id":%d}`, strategyProfileID, id)
				queryParam := ""
				entity := "strategy-profile-action-plan"
				action := "comment-added"

				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserIDs:              []uint{*strategyProfileActionPlan.ReviewUserID},
					ModuleIdentifier:     "strategy",
					NotificationType:     string(commonDTOs.AssignmentType),
					HeadingEn:            "Comment Added For StrategyProfile Action Plan",
					HeadingAr:            "تمت إضافة التعليق لخطة عمل الامتثال",
					TitleEn:              strategyProfileActionPlan.Title,
					TitleAr:              strategyProfileActionPlan.Title,
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					QueryParam:           queryParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				reviewUser, err := spaps.userSvc.GetByID(*strategyProfileActionPlan.ReviewUserID)
				if err != nil {
					utility.LogError("Failed to fetch review user", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}
				assignedBy, err := spaps.userSvc.GetByID(currentUserID)
				if err != nil {
					utility.LogError("Failed to fetch commented user", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}

				emailTo := []string{reviewUser.Email}
				subject := fmt.Sprintf("Comment Added for Strategy Profile Action Plan:%s", strategyProfileActionPlan.Title)
				appURL := utility.GetAppURL()
				emailData := map[string]interface{}{
					"FirstNameEn":            reviewUser.FirstNameEn,
					"LastNameEn":             reviewUser.LastNameEn,
					"StrategyProfileTitle":   strategyProfileActionPlan.StrategyProfile.Title,
					"ActionPlanTitle":        strategyProfileActionPlan.Title,
					"CommentedByFirstNameEn": assignedBy.FirstNameEn,
					"CommentedByLastNameEn":  assignedBy.LastNameEn,
					"CommentedOn":            time.Now().Format("2006-01-02"),
					"RedirectURL":            fmt.Sprintf("%sstrategy/strategy-profile/%d/action-plans?action_plan_id=%d", appURL, strategyProfileID, id),
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              emailTo,
					Subject:              subject,
					TemplateRelativePath: "strategy/templates/email/action_plan_comment.html",
					EmailData:            emailData,
				}

				// send notification and mail
				err = spaps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "strategy_profile_action_plan")
				if err != nil {
					utility.LogError("Failed to send notifications", map[string]interface{}{
						"error": err.Error(),
					})
				}
			}
		} else if strategyProfileActionPlan.ReviewUserID == &currentUserID {
			// if review user added the comment, notification for responsible user
			// prepare notification
			urlParam := fmt.Sprintf(`{"strategy_profile_id":%d, "action_plan_id":%d}`, strategyProfileID, id)
			queryParam := ""
			entity := "strategy-profile-action-plan"
			action := "comment-added"

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{*strategyProfileActionPlan.ResponsibleUserID},
				ModuleIdentifier:     "strategy",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Comment Added For Strategy Profile Action Plan",
				HeadingAr:            "تمت إضافة التعليق لخطة عمل الامتثال",
				TitleEn:              strategyProfileActionPlan.Title,
				TitleAr:              strategyProfileActionPlan.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			responsibleUser, err := spaps.userSvc.GetByID(*strategyProfileActionPlan.ResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}
			assignedBy, err := spaps.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch commented user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			emailTo := []string{responsibleUser.Email}
			subject := fmt.Sprintf("Comment Added for StrategyProfile Action Plan:%s", strategyProfileActionPlan.Title)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":            responsibleUser.FirstNameEn,
				"LastNameEn":             responsibleUser.LastNameEn,
				"StrategyProfileTitle":   strategyProfileActionPlan.StrategyProfile.Title,
				"ActionPlanTitle":        strategyProfileActionPlan.Title,
				"CommentedByFirstNameEn": assignedBy.FirstNameEn,
				"CommentedByLastNameEn":  assignedBy.LastNameEn,
				"CommentedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":            fmt.Sprintf("%sstrategy/strategy-profile/%d/action-plans?action_plan_id=%d", appURL, strategyProfileID, id),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "strategy/templates/email/action_plan_comment.html",
				EmailData:            emailData,
			}

			// send notification and mail
			err = spaps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "strategy_profile_action_plan")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		} else {
			// if other user added the comment, notification for review user and responsible user
			// prepare notification
			urlParam := fmt.Sprintf(`{"strategy_profile_id":%d, "action_plan_id":%d}`, strategyProfileID, id)
			queryParam := ""
			entity := "strategy-profile-action-plan"
			action := "comment-added"

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{*strategyProfileActionPlan.ResponsibleUserID},
				ModuleIdentifier:     "strategy",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Comment Added For Strategy Profile Action Plan",
				HeadingAr:            "تمت إضافة التعليق لخطة عمل الامتثال",
				TitleEn:              strategyProfileActionPlan.Title,
				TitleAr:              strategyProfileActionPlan.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			responsibleUser, err := spaps.userSvc.GetByID(*strategyProfileActionPlan.ResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}
			assignedBy, err := spaps.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch commented user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			emailTo := []string{responsibleUser.Email}
			subject := fmt.Sprintf("Comment Added for Strategy Profile Action Plan:%s", strategyProfileActionPlan.Title)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":            responsibleUser.FirstNameEn,
				"LastNameEn":             responsibleUser.LastNameEn,
				"StrategyProfileTitle":   strategyProfileActionPlan.StrategyProfile.Title,
				"ActionPlanTitle":        strategyProfileActionPlan.Title,
				"CommentedByFirstNameEn": assignedBy.FirstNameEn,
				"CommentedByLastNameEn":  assignedBy.LastNameEn,
				"CommentedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":            fmt.Sprintf("%sstrategy/strategy-profile/%d/action-plans?action_plan_id=%d", appURL, strategyProfileID, id),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "strategy/templates/email/action_plan_comment.html",
				EmailData:            emailData,
			}

			// send notification and mail
			err = spaps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "strategy_profile_action_plan")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
			if strategyProfileActionPlan.ReviewUserID != nil {
				// prepare notification
				urlParam := fmt.Sprintf(`{"strategy_profile_id":%d, "action_plan_id":%d}`, strategyProfileID, id)
				queryParam := ""
				entity := "strategy-profile-action-plan"
				action := "comment-added"

				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserIDs:              []uint{*strategyProfileActionPlan.ReviewUserID},
					ModuleIdentifier:     "strategy",
					NotificationType:     string(commonDTOs.AssignmentType),
					HeadingEn:            "Comment Added For Strategy Profile Action Plan",
					HeadingAr:            "تمت إضافة التعليق لخطة عمل الامتثال",
					TitleEn:              strategyProfileActionPlan.Title,
					TitleAr:              strategyProfileActionPlan.Title,
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					QueryParam:           queryParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				responsibleUser, err := spaps.userSvc.GetByID(*strategyProfileActionPlan.ReviewUserID)
				if err != nil {
					utility.LogError("Failed to fetch review user", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}
				assignedBy, err := spaps.userSvc.GetByID(currentUserID)
				if err != nil {
					utility.LogError("Failed to fetch commented user", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}

				emailTo := []string{responsibleUser.Email}
				subject := fmt.Sprintf("Comment Added for StrategyProfile Action Plan:%s", strategyProfileActionPlan.Title)
				appURL := utility.GetAppURL()
				emailData := map[string]interface{}{
					"FirstNameEn":            responsibleUser.FirstNameEn,
					"LastNameEn":             responsibleUser.LastNameEn,
					"StrategyProfileTitle":   strategyProfileActionPlan.StrategyProfile.Title,
					"ActionPlanTitle":        strategyProfileActionPlan.Title,
					"CommentedByFirstNameEn": assignedBy.FirstNameEn,
					"CommentedByLastNameEn":  assignedBy.LastNameEn,
					"CommentedOn":            time.Now().Format("2006-01-02"),
					"RedirectURL":            fmt.Sprintf("%sstrategy/strategy-profile/%d/action-plans?action_plan_id=%d", appURL, strategyProfileID, id),
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              emailTo,
					Subject:              subject,
					TemplateRelativePath: "strategy/templates/email/action_plan_comment.html",
					EmailData:            emailData,
				}

				// send notification and mail
				err = spaps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "strategy_profile_action_plan")
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

func (spaps *StrategyProfileActionPlanCommentService) Update(strategyProfileID uint, actionPlanID uint, id uint, currentUserID uint, strategyProfileActionPlanFormDTO *strategyDTOs.StrategyProfileActionPlanCommentFormDTO) error {
	tx := spaps.repo.BeginTransaction()

	if _, err := spaps.strategyProfileActionPlanSvc.GetByID(strategyProfileID, actionPlanID); err != nil {
		return err
	}

	existingComment, err := spaps.GetByID(strategyProfileID, actionPlanID, id)
	if err != nil {
		return err
	}

	oldComment, err := spaps.GetFullDetailsByID(strategyProfileID, actionPlanID, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingComment.Comment = strategyProfileActionPlanFormDTO.Comment
	err = spaps.repo.Update(existingComment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	var newDocuments []*models.StrategyProfileActionPlanCommentDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.StrategyProfileActionPlanCommentDocument

	for _, doc := range strategyProfileActionPlanFormDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {

			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}

			url := utility.GenerateDocumentURL("strategy", "strategy-profile", strconv.FormatUint(uint64(strategyProfileID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

			document := &models.StrategyProfileActionPlanCommentDocument{
				StrategyProfileActionPlanCommentID: id,
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
				existingDocument, err := spaps.repo.GetDocumentByID(*doc.ID)
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
		if err := spaps.repo.CreateDocuments(newDocuments, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIds) > 0 {
		if err := spaps.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range newDocuments {
		storagePath := filepath.Join("strategy", "strategy-profile", strconv.FormatUint(uint64(strategyProfileID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("strategy", "strategy-profile", strconv.FormatUint(uint64(strategyProfileID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	go func() {
		existingComment, err := spaps.GetFullDetailsByID(strategyProfileID, actionPlanID, id)
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
			err := spaps.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.StrategyProfileActionPlanComment{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.StrategyProfileActionPlanComment{}.TableName(),
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

func (spaps *StrategyProfileActionPlanCommentService) Delete(strategyProfileID uint, actionPlanID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := spaps.repo.BeginTransaction()
	if _, err := spaps.strategyProfileActionPlanSvc.GetByID(strategyProfileID, actionPlanID); err != nil {
		return nil, err
	}

	var failedIDs []uint
	var successfulIDs []uint
	var deletedComments []map[string]interface{}
	for _, id := range ids {
		if _, err := spaps.repo.GetFullDetailsByID(id); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		deletedComments = append(deletedComments, map[string]interface{}{
			"id":                      id,
			"strategy_profile_action": actionPlanID,
		})

		err := spaps.repo.Delete(id, tx)
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
			err := spaps.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.StrategyProfileActionPlanComment{}.TableName(),
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

func (spaps *StrategyProfileActionPlanCommentService) GetDocumentPath(strategyProfileID, actionPlanID, id uint, token string) (string, error) {
	existingComment, err := spaps.GetByID(strategyProfileID, actionPlanID, id)
	if err != nil {
		return "", err
	}
	if len(existingComment.Documents) == 0 {
		return "", fmt.Errorf("no documents found for the strategy profile action plan comment")
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	storagePath := filepath.Join("strategy", "strategy-profile", strconv.FormatUint(uint64(strategyProfileID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}
