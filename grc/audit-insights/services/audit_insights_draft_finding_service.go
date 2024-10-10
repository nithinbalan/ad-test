package services

import (
	"jethur/grc/audit-insights/dtos"
	"jethur/grc/audit-insights/repositories"
	auditDTOs "jethur/grc/audit/dtos"
	auditServices "jethur/grc/audit/services"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"
)

type AuditInsightsDraftFindingServiceInterface interface {
	Accept(auditDraftFindingID uint, currentUserID uint, auditInsightsDraftFindingApproveFormDTO *dtos.AuditInsightsDraftFindingApproveFormDTO) error
	Challenge(auditDraftFindingID uint, currentUserID uint, auditInsightsDraftFindingChallengeFormDTO *dtos.AuditInsightsDraftFindingChallengeFormDTO) error
}

type AuditInsightsDraftFindingService struct {
	repo                             repositories.AuditInsightsDraftFindingRepositoryInterface
	activityLogSvc                   logServices.ActivityLogServiceInterface
	userSvc                          organizationServices.UserServiceInterface
	auditDraftReportSvc              auditServices.AuditDraftReportServiceInterface
	auditDraftReportStatusSvc        auditServices.AuditDraftReportStatusServiceInterface
	auditDraftFindingSvc             auditServices.AuditDraftFindingServiceInterface
	auditDraftFindingStatusSvc       auditServices.AuditDraftFindingStatusServiceInterface
	auditDraftFindingConversationSvc auditServices.AuditDraftFindingConversationServiceInterface
	auditDraftFindingActionPlanSvc   auditServices.AuditDraftFindingActionPlanServiceInterface
}

func NewAuditInsightsDraftFindingService(
	repo repositories.AuditInsightsDraftFindingRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	auditDraftReportSvc auditServices.AuditDraftReportServiceInterface,
	auditDraftReportStatusSvc auditServices.AuditDraftReportStatusServiceInterface,
	auditDraftFindingSvc auditServices.AuditDraftFindingServiceInterface,
	auditDraftFindingStatusSvc auditServices.AuditDraftFindingStatusServiceInterface,
	auditDraftFindingConversationSvc auditServices.AuditDraftFindingConversationServiceInterface,
	auditDraftFindingActionPlanSvc auditServices.AuditDraftFindingActionPlanServiceInterface,
) AuditInsightsDraftFindingServiceInterface {
	return &AuditInsightsDraftFindingService{
		repo:                             repo,
		activityLogSvc:                   activityLogSvc,
		userSvc:                          userSvc,
		auditDraftReportSvc:              auditDraftReportSvc,
		auditDraftReportStatusSvc:        auditDraftReportStatusSvc,
		auditDraftFindingSvc:             auditDraftFindingSvc,
		auditDraftFindingStatusSvc:       auditDraftFindingStatusSvc,
		auditDraftFindingConversationSvc: auditDraftFindingConversationSvc,
		auditDraftFindingActionPlanSvc:   auditDraftFindingActionPlanSvc,
	}
}

func (aidfs *AuditInsightsDraftFindingService) Accept(auditDraftFindingID uint, currentUserID uint, auditInsightsDraftFindingApproveFormDTO *dtos.AuditInsightsDraftFindingApproveFormDTO) error {
	existingAuditDraftFinding, err := aidfs.auditDraftFindingSvc.GetFullDetailsByID(auditDraftFindingID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, auditDraftFindingID)
	}
	orgUnitID, err := aidfs.userSvc.GetUserPrimaryUnitID(currentUserID)
	if err != nil {
		return err
	}
	if *existingAuditDraftFinding.AuditableItem.ResponsibleUnitID != orgUnitID {
		return utility.NewForbiddenError("You are not authorized to accept this draft finding. Only users from the responsible unit can perform this action.")
	}

	// Retrieve the AuditDraftReport using AuditID from AuditDraftFinding.Audit
	auditID := existingAuditDraftFinding.AuditID
	auditDraftReport, err := aidfs.auditDraftReportSvc.GetByAuditID(auditID)
	if err != nil {
		return err
	}

	// Check if AuditDraftReportStatus is 'InDepartmentReview'
	inReviewStatus, err := aidfs.auditDraftReportStatusSvc.GetByType("in-department-review")
	if err != nil {
		return err
	}
	if auditDraftReport.AuditDraftReportStatusID != inReviewStatus.ID {
		return utility.NewForbiddenError("Action not allowed: Only audit draft reports with an 'In Department Review' status can be accepted.")
	}

	auditDraftFindingStatus, err := aidfs.auditDraftFindingStatusSvc.GetByType("open")
	if err != nil {
		return err
	}
	if existingAuditDraftFinding.AuditDraftFindingStatusID != auditDraftFindingStatus.ID {
		return utility.NewForbiddenError("Action not allowed: Only audit draft findings with an 'open' status can be processed. Please ensure the audit draft finding is open before proceeding.")
	}
	auditDraftFindingStatusAccepted, err := aidfs.auditDraftFindingStatusSvc.GetByType("accepted")
	if err != nil {
		return err
	}

	// Update Audit Draft Finding Status
	if err := aidfs.auditDraftFindingSvc.StatusUpdate(auditDraftFindingID, currentUserID, auditDraftFindingStatusAccepted.ID); err != nil {
		return err
	}
	isAccepted := false
	isDismissed := false
	auditDraftFindingConversationFormDTO := auditDTOs.AuditDraftFindingConversationFormDTO{
		Comment:     auditInsightsDraftFindingApproveFormDTO.Comment,
		Type:        auditInsightsDraftFindingApproveFormDTO.Type,
		IsAccepted:  &isAccepted,
		IsDismissed: &isDismissed,
	}
	if _, err := aidfs.auditDraftFindingConversationSvc.Create(auditDraftFindingID, currentUserID, &auditDraftFindingConversationFormDTO); err != nil {
		return err
	}
	auditDraftFindingActionPlanFormDTO := auditDTOs.AuditDraftFindingActionPlanFormDTO{
		Title:             auditInsightsDraftFindingApproveFormDTO.Title,
		Description:       auditInsightsDraftFindingApproveFormDTO.Description,
		ResponsibleUserID: auditInsightsDraftFindingApproveFormDTO.ResponsibleUserID,
		TargetDate:        auditInsightsDraftFindingApproveFormDTO.TargetDate,
	}

	if _, err := aidfs.auditDraftFindingActionPlanSvc.Create(auditDraftFindingID, currentUserID, &auditDraftFindingActionPlanFormDTO); err != nil {
		return err
	}
	//todo - notification to user

	go func() {
		err := aidfs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditDraftFindingConversation{}.TableName(),
			auditDraftFindingID,
			auditInsightsDraftFindingApproveFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AuditDraftFindingConversation{}.TableName(),
				"recordID":  auditDraftFindingID,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (aidfs *AuditInsightsDraftFindingService) Challenge(auditDraftFindingID uint, currentUserID uint, auditInsightsDraftFindingChallengeFormDTO *dtos.AuditInsightsDraftFindingChallengeFormDTO) error {
	existingAuditDraftFinding, err := aidfs.auditDraftFindingSvc.GetFullDetailsByID(auditDraftFindingID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, auditDraftFindingID)
	}
	orgUnitID, err := aidfs.userSvc.GetUserPrimaryUnitID(currentUserID)
	if err != nil {
		return err
	}

	if *existingAuditDraftFinding.AuditableItem.ResponsibleUnitID != orgUnitID {
		return utility.NewForbiddenError("You are not authorized to accept this draft finding. Only users from the responsible unit can perform this action.")
	}

	// Retrieve the AuditDraftReport using AuditID from AuditDraftFinding.Audit
	auditID := existingAuditDraftFinding.AuditID
	auditDraftReport, err := aidfs.auditDraftReportSvc.GetByAuditID(auditID)
	if err != nil {
		return err
	}

	// Check if AuditDraftReportStatus is 'InDepartmentReview'
	inReviewStatus, err := aidfs.auditDraftReportStatusSvc.GetByType("in-department-review")
	if err != nil {
		return err
	}
	if auditDraftReport.AuditDraftReportStatusID != inReviewStatus.ID {
		return utility.NewForbiddenError("Action not allowed: Only audit draft reports with an 'In Department Review' status can be challenged.")
	}

	auditDraftFindingStatus, err := aidfs.auditDraftFindingStatusSvc.GetByType("open")
	if err != nil {
		return err
	}
	if existingAuditDraftFinding.AuditDraftFindingStatusID != auditDraftFindingStatus.ID {
		return utility.NewForbiddenError("Action not allowed: Only audit draft findings with an 'open' status can be processed. Please ensure the audit draft finding is open before proceeding.")
	}
	auditDraftFindingStatusChallenged, err := aidfs.auditDraftFindingStatusSvc.GetByType("challenged")
	if err != nil {
		return err
	}

	// Update Audit Draft Finding Status
	if err := aidfs.auditDraftFindingSvc.StatusUpdate(auditDraftFindingID, currentUserID, auditDraftFindingStatusChallenged.ID); err != nil {
		return err
	}
	isAccepted := false
	isDismissed := false
	auditDraftFindingConversationFormDTO := auditDTOs.AuditDraftFindingConversationFormDTO{
		Comment:     auditInsightsDraftFindingChallengeFormDTO.Comment,
		Type:        auditInsightsDraftFindingChallengeFormDTO.Type,
		IsAccepted:  &isAccepted,
		IsDismissed: &isDismissed,
	}
	if _, err := aidfs.auditDraftFindingConversationSvc.Create(auditDraftFindingID, currentUserID, &auditDraftFindingConversationFormDTO); err != nil {
		return err
	}
	//todo - notification to user

	go func() {
		err := aidfs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditDraftFindingConversation{}.TableName(),
			auditDraftFindingID,
			auditInsightsDraftFindingChallengeFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AuditDraftFindingConversation{}.TableName(),
				"recordID":  auditDraftFindingID,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}
