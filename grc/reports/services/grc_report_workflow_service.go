package services

import (
	"errors"
	"fmt"
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

type GRCReportWorkflowServiceInterface interface {
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*dtos.GRCReportWorkflowDetailDTO, error)
	GetByType(id uint, status string) (*models.GRCReportWorkflow, error)
	UpdateGRCReportSectionWorkflowReviewUser(reportID uint, currentUserID uint, grcReportWorkflowFormDTO *dtos.GRCReportWorkflowFormDTO) (uint, error)
	UpdateGRCReportSectionWorkflowApproveUser(reportID uint, currentUserID uint, grcReportWorkflowFormDTO *dtos.GRCReportWorkflowFormDTO) (uint, error)
	GRCReportSectionSubmit(reportID uint, currentUserID uint, formDTO *dtos.GRCReportSectionReviewApproveFormFormDTO) error
}

type GRCReportWorkflowService struct {
	repo                       repositories.GRCReportWorkflowRepositoryInterface
	activityLogSvc             logServices.ActivityLogServiceInterface
	grcReportWorkflowStatusSvc GRCReportWorkflowStatusServiceInterface
	grcReportSvc               GRCReportServiceInterface
	userSvc                    organizationServices.UserServiceInterface
	userSettingSvc             organizationServices.UserSettingServiceInterface
	notificationSvc            commonServices.NotificationServiceInterface
}

func NewGRCReportWorkflowService(
	repo repositories.GRCReportWorkflowRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	grcReportWorkflowStatusSvc GRCReportWorkflowStatusServiceInterface,
	grcReportSvc GRCReportServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) GRCReportWorkflowServiceInterface {
	return &GRCReportWorkflowService{
		repo:                       repo,
		activityLogSvc:             activityLogSvc,
		grcReportWorkflowStatusSvc: grcReportWorkflowStatusSvc,
		grcReportSvc:               grcReportSvc,
		userSvc:                    userSvc,
		userSettingSvc:             userSettingSvc,
		notificationSvc:            notificationSvc,
	}
}

func (grws GRCReportWorkflowService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*dtos.GRCReportWorkflowDetailDTO, error) {
	// Initialize your DTO
	var stakeholderDTO dtos.GRCReportWorkflowDetailDTO

	grcReportWorkflow, err := grws.repo.GetFullDetailsByID(id, tx...)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Map the result to your DTO
	err = mapstructure.Decode(grcReportWorkflow, &stakeholderDTO)
	if err != nil {
		return nil, err
	}

	stakeholderDTO.CreatedAt = grcReportWorkflow.CreatedAt
	stakeholderDTO.UpdatedAt = grcReportWorkflow.UpdatedAt

	return &stakeholderDTO, nil
}

func (grws *GRCReportWorkflowService) GetByType(id uint, status string) (*models.GRCReportWorkflow, error) {
	workflow, err := grws.repo.GetByType(id, status)

	if err != nil {
		return nil, err
	}

	return workflow, nil
}

func (grws *GRCReportWorkflowService) UpdateGRCReportSectionWorkflowReviewUser(reportID uint, currentUserID uint, grcReportWorkflowFormDTO *dtos.GRCReportWorkflowFormDTO) (uint, error) {

	tx := grws.repo.BeginTransaction()

	grcReportWorkflow, err := grws.GetByType(reportID, grcReportWorkflowFormDTO.Type)
	if err != nil {
		return 0, err
	}

	grcReportWorkflow.ReviewUserID = &grcReportWorkflowFormDTO.UserID
	grcReportWorkflow.UpdatedBYID = &currentUserID

	if err := grws.repo.Update(grcReportWorkflow); err != nil {
		tx.Rollback()
		return grcReportWorkflow.ID, err
	}

	tx.Commit()
	return grcReportWorkflow.ID, nil
}

func (grws *GRCReportWorkflowService) UpdateGRCReportSectionWorkflowApproveUser(reportID uint, currentUserID uint, grcReportWorkflowFormDTO *dtos.GRCReportWorkflowFormDTO) (uint, error) {

	tx := grws.repo.BeginTransaction()

	grcReportWorkflow, err := grws.GetByType(reportID, grcReportWorkflowFormDTO.Type)
	if err != nil {
		return 0, err
	}

	grcReportWorkflow.ApproveUserID = &grcReportWorkflowFormDTO.UserID
	grcReportWorkflow.UpdatedBYID = &currentUserID

	if err := grws.repo.Update(grcReportWorkflow); err != nil {
		tx.Rollback()
		return grcReportWorkflow.ID, err
	}

	tx.Commit()
	return grcReportWorkflow.ID, nil
}

func (grws *GRCReportWorkflowService) GRCReportSectionSubmit(grcReportID uint, currentUserID uint, formDTO *dtos.GRCReportSectionReviewApproveFormFormDTO) error {
	_, err := grws.grcReportSvc.GetGRCReportDetails(grcReportID)
	if err != nil {
		return err
	}
	existingReportWorkflow, err := grws.repo.GetByType(grcReportID, formDTO.Type)
	if err != nil {
		return utility.NewResourceNotFoundError(err, grcReportID)
	}

	grcReportWorkflowDetails, err := grws.repo.GetFullDetailsByID(existingReportWorkflow.ID)
	if err != nil {
		return err
	}

	tx := grws.repo.BeginTransaction()

	workflowStatus, err := grws.grcReportWorkflowStatusSvc.GetByType("in-review")
	if err != nil {
		return err
	}

	existingReportWorkflow.GRCReportWorkflowStatusID = workflowStatus.ID
	existingReportWorkflow.UpdatedBYID = &currentUserID

	if err := grws.repo.Update(existingReportWorkflow); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Log the action
	go func() {
		err := grws.activityLogSvc.LogActivity(
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

		responsibleUserID := existingReportWorkflow.ReviewUserID
		responsibleUser, err := grws.userSvc.GetByID(uint(*responsibleUserID))
		if err != nil {
			utility.LogError("Failed to fetch responsible user", map[string]interface{}{
				"userID": responsibleUserID,
				"error":  err.Error(),
			})
			return
		}

		urlParam := fmt.Sprintf(`{"grc_report_id":%d}`, grcReportID)
		entity := "grc-report-section-submitted"
		action := "review-grc-report-section"
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
				"Type":           existingReportWorkflow.Type,
				"RedirectURL":    fmt.Sprintf("%sreport/grc-reports/%d/overview", appURL, grcReportID),
			},
		}

		// Send Notifications
		if err := grws.notificationSvc.SendNotification(dbNotification, emailNotification, "grc_report"); err != nil {
			utility.LogError("Failed to send document request notifications", map[string]interface{}{
				"grcReportID": grcReportID,
				"userID":      responsibleUserID,
				"error":       err.Error(),
			})
		}
	}()

	return nil
}
