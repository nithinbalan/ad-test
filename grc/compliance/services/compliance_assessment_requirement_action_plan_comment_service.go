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

type ComplianceAssessmentRequirementActionPlanCommentServiceInterface interface {
	GetAll(assessmentID uint, requirementID uint, actionPlanID uint) ([]*complianceDTOs.ComplianceActionPlanCommentListDTO, error)
	GetByID(assessmentID uint, requirementID uint, actionPlanID uint, id uint) (*models.ComplianceActionPlanComment, error)
	GetFullDetailsByID(assessmentID, requirementID uint, actionPlanID uint, id uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceActionPlanCommentDetailDTO, error)
	Create(assessmentID uint, requirementID uint, actionPlanID uint, currentUserID uint, dto *complianceDTOs.ComplianceActionPlanCommentFormDTO) (uint, error)
	Update(assessmentID uint, requirementID uint, actionPlanID uint, id uint, currentUserID uint, dto *complianceDTOs.ComplianceActionPlanCommentFormDTO) error
	Delete(assessmentID uint, requirementID uint, actionPlanID uint, ids []uint, currentUserID uint) ([]uint, error)

	GetDocumentPath(assessmentID, requirementID, actionPlanID, id uint, token string) (string, error)
}

type ComplianceAssessmentRequirementActionPlanCommentService struct {
	repo                    repositories.ComplianceActionPlanCommentRepositoryInterface
	userSvc                 organizationServices.UserServiceInterface
	activityLogSvc          logServices.ActivityLogServiceInterface
	complianceActionPlanSvc ComplianceAssessmentRequirementActionPlanServiceInterface
	notificationSvc         commonServices.NotificationServiceInterface
}

func NewComplianceAssessmentRequirementActionPlanCommentService(repo repositories.ComplianceActionPlanCommentRepositoryInterface, userSvc organizationServices.UserServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	complianceActionPlanSvc ComplianceAssessmentRequirementActionPlanServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) ComplianceAssessmentRequirementActionPlanCommentServiceInterface {
	return &ComplianceAssessmentRequirementActionPlanCommentService{
		repo:                    repo,
		userSvc:                 userSvc,
		activityLogSvc:          activityLogSvc,
		complianceActionPlanSvc: complianceActionPlanSvc,
		notificationSvc:         notificationSvc,
	}
}

func (carapcs *ComplianceAssessmentRequirementActionPlanCommentService) GetAll(assessmentID uint, requirementID uint, actionPlanID uint) ([]*complianceDTOs.ComplianceActionPlanCommentListDTO, error) {
	if _, err := carapcs.complianceActionPlanSvc.GetByID(assessmentID, requirementID, actionPlanID); err != nil {
		return nil, err
	}

	return carapcs.repo.GetAll(actionPlanID)
}

func (carapcs *ComplianceAssessmentRequirementActionPlanCommentService) GetByID(assessmentID uint, requirementID uint, actionPlanID uint, id uint) (*models.ComplianceActionPlanComment, error) {
	if _, err := carapcs.complianceActionPlanSvc.GetByID(assessmentID, requirementID, actionPlanID); err != nil {
		return nil, err
	}

	comment, err := carapcs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return comment, nil
}

func (carapcs *ComplianceAssessmentRequirementActionPlanCommentService) GetFullDetailsByID(assessmentID uint, requirementID uint, actionPlanID uint, id uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceActionPlanCommentDetailDTO, error) {
	if _, err := carapcs.complianceActionPlanSvc.GetByID(assessmentID, requirementID, actionPlanID); err != nil {
		fmt.Printf("ComplianceActionPlan not found for assessmentID %d, actionPlanID %d\n", assessmentID, actionPlanID)
		return nil, err
	}

	var commentDetail complianceDTOs.ComplianceActionPlanCommentDetailDTO

	comment, err := carapcs.repo.GetFullDetailsByID(id, tx...)
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

func (carapcs *ComplianceAssessmentRequirementActionPlanCommentService) Create(assessmentID uint, requirementID uint, actionPlanID uint, currentUserID uint, complianceActionPlanCommentFormDTO *complianceDTOs.ComplianceActionPlanCommentFormDTO) (uint, error) {
	tx := carapcs.repo.BeginTransaction()

	if _, err := carapcs.complianceActionPlanSvc.GetByID(assessmentID, requirementID, actionPlanID); err != nil {
		return 0, err
	}

	complianceActionPlan, err := carapcs.complianceActionPlanSvc.GetFullDetailsByID(assessmentID, requirementID, actionPlanID)
	if err != nil {
		return 0, err
	}

	comment := &models.ComplianceActionPlanComment{
		ComplianceActionPlanID: actionPlanID,
		Comment:                complianceActionPlanCommentFormDTO.Comment,
		CreatedBYID:            currentUserID,
	}

	id, err := carapcs.repo.Create(comment)
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
			"compliance", "compliance-register", strconv.FormatUint(uint64(assessmentID), 10),
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
		if err := carapcs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}

		for _, doc := range documents {
			storagePath := filepath.Join(
				"compliance", "compliance-register", strconv.FormatUint(uint64(assessmentID), 10),
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
		existingComment, err := carapcs.GetFullDetailsByID(assessmentID, requirementID, actionPlanID, id)
		if err != nil {
			utility.LogError("Failed to fetch compliance action plan comment", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = carapcs.activityLogSvc.LogActivity(
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
				urlParam := fmt.Sprintf(`{"compliance_assessment_id":%d,"requirement_id":%d, "action_plan_id":%d}`, assessmentID, requirementID, id)
				queryParam := ""
				entity := "compliance-assessment-requirement-action-plan"
				action := "comment-added"

				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserID:               *complianceActionPlan.ReviewUserID,
					ModuleIdentifier:     "compliance",
					NotificationType:     string(commonDTOs.AssignmentType),
					HeadingEn:            "Comment Added For Compliance Assessment Requirement Action Plan",
					HeadingAr:            "تمت إضافة التعليق لخطة عمل متطلبات تقييم الامتثال",
					TitleEn:              complianceActionPlan.Title,
					TitleAr:              complianceActionPlan.Title,
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					QueryParam:           queryParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				reviewUser, err := carapcs.userSvc.GetByID(*complianceActionPlan.ReviewUserID)
				if err != nil {
					utility.LogError("Failed to fetch review user", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}
				assignedBy, err := carapcs.userSvc.GetByID(currentUserID)
				if err != nil {
					utility.LogError("Failed to fetch commented user", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}

				emailTo := []string{reviewUser.Email}
				subject := fmt.Sprintf("Comment Added for Compliance Assessment Requirement Action Plan:%s", complianceActionPlan.Title)
				appURL := utility.GetAppURL()
				emailData := map[string]interface{}{
					"FirstNameEn":            reviewUser.FirstNameEn,
					"LastNameEn":             reviewUser.LastNameEn,
					"Requirement":            complianceActionPlan.ComplianceAssessmentRequirement.ComplianceRequirement.TitleEn,
					"ActionPlanTitle":        complianceActionPlan.Title,
					"ActionPlanDescription":  complianceActionPlan.Description,
					"CommentedByFirstNameEn": assignedBy.FirstNameEn,
					"CommentedByLastNameEn":  assignedBy.LastNameEn,
					"CommentedOn":            time.Now().Format("2006-01-02"),
					"RedirectURL":            fmt.Sprintf("%scompliance/assessments/%d/overview?action_plan_id=%d", appURL, assessmentID, actionPlanID),
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              emailTo,
					Subject:              subject,
					TemplateRelativePath: "compliance/templates/email/compliance_assessment_requirement_action_plan_comment.html",
					EmailData:            emailData,
				}

				// send notification and mail
				err = carapcs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "compliance_assessment_requirement")
				if err != nil {
					utility.LogError("Failed to send notifications", map[string]interface{}{
						"error": err.Error(),
					})
				}
			}
		} else if complianceActionPlan.ReviewUserID == &currentUserID {
			// if review user added the comment, notification for responsible user
			// prepare notification
			urlParam := fmt.Sprintf(`{"compliance_assessment_id":%d,"requirement_id":%d, "action_plan_id":%d}`, assessmentID, requirementID, id)
			queryParam := ""
			entity := "compliance-assessment-requirement-action-plan"
			action := "comment-added"

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserID:               complianceActionPlan.ResponsibleUserID,
				ModuleIdentifier:     "compliance",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Comment Added For Compliance Assessment Requirement Action Plan",
				HeadingAr:            "تمت إضافة التعليق لخطة عمل متطلبات تقييم الامتثال",
				TitleEn:              complianceActionPlan.Title,
				TitleAr:              complianceActionPlan.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			responsibleUser, err := carapcs.userSvc.GetByID(complianceActionPlan.ResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}
			assignedBy, err := carapcs.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch commented user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			emailTo := []string{responsibleUser.Email}
			subject := fmt.Sprintf("Comment Added for Compliance Assessment Requirement Action Plan:%s", complianceActionPlan.Title)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":            responsibleUser.FirstNameEn,
				"LastNameEn":             responsibleUser.LastNameEn,
				"Requirement":            complianceActionPlan.ComplianceAssessmentRequirement.ComplianceRequirement.TitleEn,
				"ActionPlanTitle":        complianceActionPlan.Title,
				"ActionPlanDescription":  complianceActionPlan.Description,
				"CommentedByFirstNameEn": assignedBy.FirstNameEn,
				"CommentedByLastNameEn":  assignedBy.LastNameEn,
				"AssignedOn":             time.Now().Format("2006-01-02"),
				"RedirectURL":            fmt.Sprintf("%scompliance/assessments/%d/overview?action_plan_id=%d", appURL, assessmentID, actionPlanID),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "compliance/templates/email/compliance_assessment_requirement_action_plan_comment.html",
				EmailData:            emailData,
			}

			// send notification and mail
			err = carapcs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "compliance_assessment_requirement")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		} else {
			// if other user added the comment, notification for review user and responsible user
			// prepare notification
			urlParam := fmt.Sprintf(`{"compliance_assessment_id":%d,"requirement_id":%d, "action_plan_id":%d}`, assessmentID, requirementID, id)
			queryParam := ""
			entity := "compliance-assessment-requirement-action-plan"
			action := "comment-added"

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserID:               complianceActionPlan.ResponsibleUserID,
				ModuleIdentifier:     "compliance",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Comment Added For Compliance Assessment Requirement Action Plan",
				HeadingAr:            "تمت إضافة التعليق لخطة عمل متطلبات تقييم الامتثال",
				TitleEn:              complianceActionPlan.Title,
				TitleAr:              complianceActionPlan.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			responsibleUser, err := carapcs.userSvc.GetByID(complianceActionPlan.ResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}
			assignedBy, err := carapcs.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch commented user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			emailTo := []string{responsibleUser.Email}
			subject := fmt.Sprintf("Comment Added for Compliance Assessment Requirement Action Plan:%s", complianceActionPlan.Title)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":            responsibleUser.FirstNameEn,
				"LastNameEn":             responsibleUser.LastNameEn,
				"Requirement":            complianceActionPlan.ComplianceAssessmentRequirement.ComplianceRequirement.TitleEn,
				"ActionPlanTitle":        complianceActionPlan.Title,
				"ActionPlanDescription":  complianceActionPlan.Description,
				"CommentedByFirstNameEn": assignedBy.FirstNameEn,
				"CommentedByLastNameEn":  assignedBy.LastNameEn,
				"AssignedOn":             time.Now().Format("2006-01-02"),
				"RedirectURL":            fmt.Sprintf("%scompliance/assessments/%d/overview?action_plan_id=%d", appURL, assessmentID, actionPlanID),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "compliance/templates/email/compliance_assessment_requirement_action_plan_comment.html",
				EmailData:            emailData,
			}

			// send notification and mail
			err = carapcs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "compliance_assessment_requirement")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
			if complianceActionPlan.ReviewUserID != nil {
				// prepare notification
				urlParam := fmt.Sprintf(`{"compliance_assessment_id":%d,"requirement_id":%d, "action_plan_id":%d}`, assessmentID, requirementID, id)
				queryParam := ""
				entity := "compliance-assessment-requirement-action-plan"
				action := "comment-added"

				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserID:               *complianceActionPlan.ReviewUserID,
					ModuleIdentifier:     "compliance",
					NotificationType:     string(commonDTOs.AssignmentType),
					HeadingEn:            "Comment Added For Compliance Assessment Requirement Action Plan",
					HeadingAr:            "تمت إضافة التعليق لخطة عمل متطلبات تقييم الامتثال",
					TitleEn:              complianceActionPlan.Title,
					TitleAr:              complianceActionPlan.Title,
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					QueryParam:           queryParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				responsibleUser, err := carapcs.userSvc.GetByID(*complianceActionPlan.ReviewUserID)
				if err != nil {
					utility.LogError("Failed to fetch review user", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}
				assignedBy, err := carapcs.userSvc.GetByID(currentUserID)
				if err != nil {
					utility.LogError("Failed to fetch commented user", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}

				emailTo := []string{responsibleUser.Email}
				subject := fmt.Sprintf("Comment Added for Compliance Assessment Requirement Action Plan:%s", complianceActionPlan.Title)
				appURL := utility.GetAppURL()
				emailData := map[string]interface{}{
					"FirstNameEn":            responsibleUser.FirstNameEn,
					"LastNameEn":             responsibleUser.LastNameEn,
					"Requirement":            complianceActionPlan.ComplianceAssessmentRequirement.ComplianceRequirement.TitleEn,
					"ActionPlanTitle":        complianceActionPlan.Title,
					"ActionPlanDescription":  complianceActionPlan.Description,
					"CommentedByFirstNameEn": assignedBy.FirstNameEn,
					"CommentedByLastNameEn":  assignedBy.LastNameEn,
					"AssignedOn":             time.Now().Format("2006-01-02"),
					"RedirectURL":            fmt.Sprintf("%scompliance/assessments/%d/overview?action_plan_id=%d", appURL, assessmentID, actionPlanID),
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              emailTo,
					Subject:              subject,
					TemplateRelativePath: "compliance/templates/email/compliance_assessment_requirement_action_plan_comment.html",
					EmailData:            emailData,
				}

				// send notification and mail
				err = carapcs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "compliance_assessment_requirement")
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

func (carapcs *ComplianceAssessmentRequirementActionPlanCommentService) Update(assessmentID uint, requirementID uint, actionPlanID uint, id uint, currentUserID uint, complianceActionPlanCommentFormDTO *complianceDTOs.ComplianceActionPlanCommentFormDTO) error {
	tx := carapcs.repo.BeginTransaction()

	if _, err := carapcs.complianceActionPlanSvc.GetByID(assessmentID, requirementID, actionPlanID); err != nil {
		return err
	}

	existingComment, err := carapcs.GetByID(assessmentID, requirementID, actionPlanID, id)
	if err != nil {
		return err
	}

	oldComment, err := carapcs.GetFullDetailsByID(assessmentID, requirementID, actionPlanID, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingComment.Comment = complianceActionPlanCommentFormDTO.Comment
	err = carapcs.repo.Update(existingComment, tx)
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

			url := utility.GenerateDocumentURL("compliance", "compliance-register", strconv.FormatUint(uint64(assessmentID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

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
				existingDocument, err := carapcs.repo.GetDocumentByID(*doc.ID)
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
		if err := carapcs.repo.CreateDocuments(newDocuments, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIds) > 0 {
		if err := carapcs.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range newDocuments {
		storagePath := filepath.Join("compliance", "compliance-register", strconv.FormatUint(uint64(assessmentID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("compliance", "compliance-register", strconv.FormatUint(uint64(assessmentID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	go func() {
		existingComment, err := carapcs.GetFullDetailsByID(assessmentID, requirementID, actionPlanID, id)
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
			err := carapcs.activityLogSvc.LogActivity(
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

func (carapcs *ComplianceAssessmentRequirementActionPlanCommentService) Delete(assessmentID uint, requirementID uint, actionPlanID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := carapcs.repo.BeginTransaction()
	if _, err := carapcs.complianceActionPlanSvc.GetByID(assessmentID, requirementID, actionPlanID); err != nil {
		return nil, err
	}

	var failedIDs []uint
	var successfulIDs []uint
	var deletedComments []map[string]interface{}
	for _, id := range ids {
		if _, err := carapcs.repo.GetFullDetailsByID(id); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		deletedComments = append(deletedComments, map[string]interface{}{
			"id":                id,
			"compliance_action": actionPlanID,
		})

		err := carapcs.repo.Delete(id, tx)
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
			err := carapcs.activityLogSvc.LogBulkActivities(
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

func (carapcs *ComplianceAssessmentRequirementActionPlanCommentService) GetDocumentPath(assessmentID, requirementID, actionPlanID, id uint, token string) (string, error) {
	existingComment, err := carapcs.GetByID(assessmentID, requirementID, actionPlanID, id)
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

	storagePath := filepath.Join("compliance", "compliance-register", strconv.FormatUint(uint64(assessmentID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "comments", strconv.FormatUint(uint64(id), 10), "documents")

	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}
