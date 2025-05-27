package services

import (
	"errors"
	"fmt"
	"time"

	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/grc/reports/dtos"
	"jethur/grc/reports/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type GRCReportWorkflowReviewCommentServiceInterface interface {
	GetAll(grcReportID uint) ([]*dtos.GRCReportWorkflowReviewCommentListDTO, error)
	GetByID(grcReportID uint, id uint) (*models.GRCReportWorkflowReviewComment, error)
	GetFullDetailsByID(grcReportID uint, id uint, tx ...utility.DatabaseTransaction) (*dtos.GRCReportWorkflowReviewCommentDetailDTO, error)
	Create(grcReportID uint, currentUserID uint, dto *dtos.GRCReportWorkflowReviewCommentFormDTO) (uint, error)
	Update(grcReportID uint, id uint, currentUserID uint, dto *dtos.GRCReportWorkflowReviewCommentFormDTO) error
	Delete(grcReportID uint, ids []uint, currentUserID uint) ([]uint, error)

	GRCReportSectionCommentResolve(grcReportID uint, id uint, currentUserID uint) error
	GRCReportSectionCommentClose(grcReportID uint, id uint, currentUserID uint) error

	GRCReportSectionReviewed(grcReportID uint, currentUserID uint, commentActionForm *dtos.GRCReportWorkflowReviewCommentFormDTO) error
	GRCReportSectionApproved(grcReportID uint, currentUserID uint, commentActionForm *dtos.GRCReportWorkflowReviewCommentFormDTO) error
}

type GRCReportWorkflowReviewCommentService struct {
	repo                       repositories.GRCReportWorkflowReviewCommentRepositoryInterface
	grcReportSvc               GRCReportServiceInterface
	grcReportWorkflowSvc       GRCReportWorkflowServiceInterface
	grcReportWorkflowStatusSvc GRCReportWorkflowStatusServiceInterface
	grcReportCommentStatusSvc  GRCReportCommentStatusServiceInterface
	activityLogSvc             logServices.ActivityLogServiceInterface
	userSvc                    organizationServices.UserServiceInterface
	userSettingSvc             organizationServices.UserSettingServiceInterface
	notificationSvc            commonServices.NotificationServiceInterface
}

func NewGRCReportWorkflowReviewCommentService(repo repositories.GRCReportWorkflowReviewCommentRepositoryInterface,
	grcReportSvc GRCReportServiceInterface,
	grcReportWorkflowSvc GRCReportWorkflowServiceInterface,
	grcReportWorkflowStatusSvc GRCReportWorkflowStatusServiceInterface,
	grcReportCommentStatusSvc GRCReportCommentStatusServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) *GRCReportWorkflowReviewCommentService {
	return &GRCReportWorkflowReviewCommentService{
		repo:                       repo,
		grcReportSvc:               grcReportSvc,
		grcReportWorkflowSvc:       grcReportWorkflowSvc,
		grcReportWorkflowStatusSvc: grcReportWorkflowStatusSvc,
		grcReportCommentStatusSvc:  grcReportCommentStatusSvc,
		activityLogSvc:             activityLogSvc,
		userSvc:                    userSvc,
		userSettingSvc:             userSettingSvc,
		notificationSvc:            notificationSvc,
	}
}

func (grwrcs *GRCReportWorkflowReviewCommentService) GetAll(grcReportID uint) ([]*dtos.GRCReportWorkflowReviewCommentListDTO, error) {
	return grwrcs.repo.GetAll(grcReportID)
}

func (grwrcs *GRCReportWorkflowReviewCommentService) GetByID(grcReportID uint, id uint) (*models.GRCReportWorkflowReviewComment, error) {
	comment, err := grwrcs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return comment, nil
}

func (grwrcs *GRCReportWorkflowReviewCommentService) GetFullDetailsByID(grcReportID uint, id uint, tx ...utility.DatabaseTransaction) (*dtos.GRCReportWorkflowReviewCommentDetailDTO, error) {
	var commentDetail dtos.GRCReportWorkflowReviewCommentDetailDTO

	comment, err := grwrcs.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(comment, &commentDetail)
	if err != nil {
		return nil, err
	}

	commentDetail.CreatedAt = comment.CreatedAt
	return &commentDetail, nil
}

func (grwrcs *GRCReportWorkflowReviewCommentService) Create(grcReportID uint, currentUserID uint, formDTO *dtos.GRCReportWorkflowReviewCommentFormDTO) (uint, error) {
	tx := grwrcs.repo.BeginTransaction()

	grcReportWorkflow, err := grwrcs.grcReportWorkflowSvc.GetByType(grcReportID, formDTO.Type)
	if err != nil {
		return 0, err
	}

	commentStatus, err := grwrcs.grcReportCommentStatusSvc.GetByType("new")
	if err != nil {
		return 0, err
	}

	comment := &models.GRCReportWorkflowReviewComment{
		GRCReportWorkflowID:      grcReportWorkflow.ID,
		GRCReportCommentStatusID: commentStatus.ID,
		Comment:                  formDTO.Comment,
		CreatedBYID:              currentUserID,
	}

	id, err := grwrcs.repo.Create(comment)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	go func() {
		existingComment, err := grwrcs.GetFullDetailsByID(grcReportID, id)
		if err != nil {
			utility.LogError("Failed to fetch risk comment", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = grwrcs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.GRCReportWorkflowReviewComment{}.TableName(),
			id,
			existingComment,
			currentUserID,
			nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	return id, nil
}

func (grwrcs *GRCReportWorkflowReviewCommentService) Update(grcReportID uint, id uint, currentUserID uint, formDTO *dtos.GRCReportWorkflowReviewCommentFormDTO) error {
	tx := grwrcs.repo.BeginTransaction()

	existingComment, err := grwrcs.GetByID(grcReportID, id)
	if err != nil {
		return err
	}

	oldComment, err := grwrcs.GetFullDetailsByID(grcReportID, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingComment.Comment = formDTO.Comment
	err = grwrcs.repo.Update(existingComment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	go func() {
		existingComment, err := grwrcs.GetFullDetailsByID(grcReportID, id)
		if err != nil {
			return
		}

		changedData := utility.CompareFieldChanges(
			map[string]interface{}{"comment": oldComment.Comment},
			map[string]interface{}{"comment": existingComment.Comment},
		)

		if len(changedData) > 0 {
			err := grwrcs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.GRCReportWorkflowReviewComment{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to log activity", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()
	return nil
}

func (grwrcs *GRCReportWorkflowReviewCommentService) Delete(grcReportID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := grwrcs.repo.BeginTransaction()

	var failedIDs []uint
	var successfulIDs []uint
	var deletedComments []map[string]interface{}

	for _, id := range ids {
		if _, err := grwrcs.repo.GetFullDetailsByID(id); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		deletedComments = append(deletedComments, map[string]interface{}{
			"id":            id,
			"grc_report_id": grcReportID,
		})

		err := grwrcs.repo.Delete(id, tx)
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
			err := grwrcs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.GRCReportWorkflowReviewComment{}.TableName(),
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

func (grwrcs *GRCReportWorkflowReviewCommentService) GRCReportSectionCommentResolve(grcReportID uint, id uint, currentUserID uint) error {

	tx := grwrcs.repo.BeginTransaction()

	commentStatus, err := grwrcs.grcReportCommentStatusSvc.GetByType("resolved")
	if err != nil {
		tx.Rollback()
		return err
	}

	existingComment, err := grwrcs.GetByID(grcReportID, id)
	if err != nil {
		return err
	}

	oldComment, err := grwrcs.GetFullDetailsByID(grcReportID, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingComment.GRCReportCommentStatusID = commentStatus.ID
	existingComment.ResolvedBYID = &currentUserID
	err = grwrcs.repo.Update(existingComment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	go func() {
		existingComment, err := grwrcs.GetFullDetailsByID(grcReportID, id)
		if err != nil {
			return
		}

		changedData := utility.CompareFieldChanges(
			map[string]interface{}{"grc_workflow_comment_status": oldComment.GRCReportCommentStatus.TitleEn},
			map[string]interface{}{"grc_workflow_comment_status": existingComment.GRCReportCommentStatus.TitleEn},
		)

		if len(changedData) > 0 {
			err := grwrcs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.GRCReportWorkflowReviewComment{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to log activity", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()
	return nil

}

func (grwrcs *GRCReportWorkflowReviewCommentService) GRCReportSectionCommentClose(grcReportID uint, id uint, currentUserID uint) error {

	tx := grwrcs.repo.BeginTransaction()

	commentStatus, err := grwrcs.grcReportCommentStatusSvc.GetByType("closed")
	if err != nil {
		tx.Rollback()
		return err
	}

	existingComment, err := grwrcs.GetByID(grcReportID, id)
	if err != nil {
		return err
	}

	oldComment, err := grwrcs.GetFullDetailsByID(grcReportID, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingComment.GRCReportCommentStatusID = commentStatus.ID
	err = grwrcs.repo.Update(existingComment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	go func() {
		existingComment, err := grwrcs.GetFullDetailsByID(grcReportID, id)
		if err != nil {
			return
		}

		changedData := utility.CompareFieldChanges(
			map[string]interface{}{"grc_workflow_comment_status": oldComment.GRCReportCommentStatus.TitleEn},
			map[string]interface{}{"grc_workflow_comment_status": existingComment.GRCReportCommentStatus.TitleEn},
		)

		if len(changedData) > 0 {
			err := grwrcs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.GRCReportWorkflowReviewComment{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to log activity", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()
	return nil

}

func (grwrcs *GRCReportWorkflowReviewCommentService) GRCReportSectionReviewed(grcReportID uint, currentUserID uint, commentActionForm *dtos.GRCReportWorkflowReviewCommentFormDTO) error {
	tx := grwrcs.repo.BeginTransaction()

	reviewedStatus, err := grwrcs.grcReportWorkflowStatusSvc.GetByType("pending-approval")
	if err != nil {
		tx.Rollback()
		return err
	}

	commentStatus, err := grwrcs.grcReportCommentStatusSvc.GetByType("new")
	if err != nil {
		tx.Rollback()
		return err
	}

	grcReportWorkflow, err := grwrcs.grcReportWorkflowSvc.GetByType(grcReportID, commentActionForm.Type)
	if err != nil {
		tx.Rollback()
		return err
	}

	grcReportWorkflowDetails, err := grwrcs.grcReportWorkflowSvc.GetFullDetailsByID(grcReportWorkflow.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	comment := &models.GRCReportWorkflowReviewComment{
		GRCReportWorkflowID:       grcReportWorkflow.ID,
		Comment:                   commentActionForm.Comment,
		GRCReportCommentStatusID:  commentStatus.ID,
		GRCReportWorkflowStatusID: &reviewedStatus.ID,
		CreatedBYID:               currentUserID,
	}

	_, err = grwrcs.repo.Create(comment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	resolvedAt := time.Now()
	// Update workflow status to reviewed
	grcReportWorkflow.GRCReportWorkflowStatusID = reviewedStatus.ID
	grcReportWorkflow.ReviewedAt = &resolvedAt
	err = grwrcs.repo.UpdateSectionWorkflowStatus(grcReportWorkflow, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Log the action
	go func() {
		err := grwrcs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.GRCReport{}.TableName(),
			grcReportID,
			nil,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.GRCReport{}.TableName(),
				"recordID":  grcReportID,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}

		responsibleUserID := grcReportWorkflow.ApproveUserID
		responsibleUser, err := grwrcs.userSvc.GetByID(uint(*responsibleUserID))
		if err != nil {
			utility.LogError("Failed to fetch responsible user", map[string]interface{}{
				"userID": responsibleUserID,
				"error":  err.Error(),
			})
			return
		}

		urlParam := fmt.Sprintf(`{"grc_report_id":%d}`, grcReportID)
		entity := "grc-report-section-submitted"
		action := "approve-grc-report-section"
		appURL := utility.GetAppURL()

		dbNotification := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{*responsibleUserID},
			ModuleIdentifier:     "grc_report",
			NotificationType:     string(commonDTOs.InformationType),
			HeadingEn:            "GRC Report Section Submitted",
			HeadingAr:            "تم تقديم قسم تقرير GRC",
			TitleEn:              grcReportWorkflowDetails.GRCReport.Title,
			TitleAr:              grcReportWorkflowDetails.GRCReport.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		emailNotification := &commonDTOs.EmailNotificationDTO{
			EmailTo:              []string{responsibleUser.Email},
			Subject:              fmt.Sprintf("GRC Report Section Submitted: %s", grcReportWorkflowDetails.GRCReport.Title),
			TemplateRelativePath: "report/templates/email/grc_report_section_submitted.html",
			EmailData: map[string]interface{}{
				"FirstNameEn":    responsibleUser.FirstNameEn,
				"LastNameEn":     responsibleUser.LastNameEn,
				"GRCReportTitle": grcReportWorkflowDetails.GRCReport.Title,
				"Type":           grcReportWorkflow.Type,
				"RedirectURL":    fmt.Sprintf("%sreport/grc-reports/%d/overview", appURL, grcReportID),
			},
		}

		// Send Notifications
		if err := grwrcs.notificationSvc.SendNotification(dbNotification, emailNotification, "grc_report"); err != nil {
			utility.LogError("Failed to send document request notifications", map[string]interface{}{
				"grcReportID": grcReportID,
				"userID":      responsibleUserID,
				"error":       err.Error(),
			})
		}
	}()

	return nil
}

func (grwrcs *GRCReportWorkflowReviewCommentService) GRCReportSectionApproved(grcReportID uint, currentUserID uint, commentActionForm *dtos.GRCReportWorkflowReviewCommentFormDTO) error {
	tx := grwrcs.repo.BeginTransaction()

	approveStatus, err := grwrcs.grcReportWorkflowStatusSvc.GetByType("approved")
	if err != nil {
		tx.Rollback()
		return err
	}

	commentStatus, err := grwrcs.grcReportCommentStatusSvc.GetByType("new")
	if err != nil {
		tx.Rollback()
		return err
	}

	grcReportWorkflow, err := grwrcs.grcReportWorkflowSvc.GetByType(grcReportID, commentActionForm.Type)
	if err != nil {
		tx.Rollback()
		return err
	}

	comment := &models.GRCReportWorkflowReviewComment{
		GRCReportWorkflowID:       grcReportWorkflow.ID,
		Comment:                   commentActionForm.Comment,
		GRCReportCommentStatusID:  commentStatus.ID,
		GRCReportWorkflowStatusID: &approveStatus.ID,
		CreatedBYID:               currentUserID,
	}

	_, err = grwrcs.repo.Create(comment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	approvedAt := time.Now()
	// Update workflow status to approve
	grcReportWorkflow.GRCReportWorkflowStatusID = approveStatus.ID
	grcReportWorkflow.ApprovedAt = &approvedAt
	err = grwrcs.repo.UpdateSectionWorkflowStatus(grcReportWorkflow, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}
