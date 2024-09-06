package services

import (
	"errors"
	"fmt"
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organiztionDTOs "jethur/grc/organization/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditDraftReportServiceInterface interface {
	GetFullDetailsByID(auditID uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditDraftReportDetailDTO, error)
	GetByID(id uint) (*models.AuditDraftReport, error)
	Create(auditID uint, currentUserID uint, auditDraftReportFormDTO *auditDTOs.AuditDraftReportFormDTO) (uint, error)
	Update(auditID uint, id uint, currentUserID uint, auditDraftReportFormDTO *auditDTOs.AuditDraftReportFormDTO) error
	Delete(auditID uint, ids []uint, currentUserID uint) ([]uint, error)
}

type AuditDraftReportService struct {
	repo           repositories.AuditDraftReportRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
	auditSvc       AuditServiceInterface
}

func NewAuditDraftReportService(
	repo repositories.AuditDraftReportRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditSvc AuditServiceInterface,
) AuditDraftReportServiceInterface {
	return &AuditDraftReportService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
		auditSvc:       auditSvc,
	}
}

func (adrs *AuditDraftReportService) GetByID(id uint) (*models.AuditDraftReport, error) {
	auditDraftReport, err := adrs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return auditDraftReport, nil
}

func (adrs *AuditDraftReportService) GetFullDetailsByID(auditID uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditDraftReportDetailDTO, error) {
	if _, err := adrs.auditSvc.GetByID(auditID); err != nil {
		return nil, err
	}

	auditDraftReport, err := adrs.repo.GetFullDetailsByID(auditID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, auditID)
		}
		return nil, err
	}

	// Map the data to the DTO
	auditReportDTO := &auditDTOs.AuditDraftReportDetailDTO{
		ID:                 auditDraftReport.ID,
		AuditID:            auditDraftReport.AuditID,
		Title:              auditDraftReport.Audit.AuditPlan.Title,
		ReferenceNumber:    auditDraftReport.Audit.AuditPlan.ReferenceNumber,
		ReportDateString:   auditDraftReport.ReportDate.Format("2006-01-02"),
		PreparedBy:         mapUserToUserMinimalDTO(auditDraftReport.Audit.CreatedByUser),
		PreparedAt:         auditDraftReport.CreatedAt.String(),
		ExecutiveSummary:   *auditDraftReport.ExecutiveSummary,
		TeamLead:           mapUserToUserMinimalDTO(auditDraftReport.Audit.AuditPlan.TeamLead),
		Conclusion:         *auditDraftReport.Conclusion,
		Auditors:           []*auditDTOs.AuditPlanAuditMemberDTO{}, // Initialize as empty slice
		AuditDraftFindings: []*auditDTOs.AuditDraftFindingDetailDTO{},
	}

	// Map auditors
	for _, auditor := range auditDraftReport.Audit.AuditPlan.Members {
		auditReportDTO.Auditors = append(auditReportDTO.Auditors, &auditDTOs.AuditPlanAuditMemberDTO{
			ID:          &auditor.ID,
			AuditPlanID: auditor.AuditPlanID,
			UserID:      auditor.UserID,
			User:        mapUserToUserMinimalDTO(auditor.User),
		})
	}
	// Map findings
	for _, finding := range auditDraftReport.Audit.AuditDraftFindings {
		findingDTO := auditDTOs.AuditDraftFindingDetailDTO{
			ID:                        finding.ID,
			Title:                     finding.Title,
			Description:               finding.Description,
			AuditID:                   finding.AuditID,
			AuditableItemID:           finding.AuditableItemID,
			AuditChecklistID:          finding.AuditChecklistID,
			AuditFindingCategoryID:    finding.AuditFindingCategoryID,
			Recommendations:           finding.Recommendations,
			AuditDraftFindingStatusID: finding.AuditDraftFindingStatusID,
			CreatedAt:                 finding.CreatedAt,
			UpdatedAt:                 finding.UpdatedAt,
			Documents:                 []*auditDTOs.AuditDraftFindingDocumentDTO{},
			// CreatedByUser:             mapUserToUserMinimalDTO(finding.CreatedByUser),
			// UpdatedByUser:             mapUserToUserMinimalDTO(finding.UpdatedByUser),
		}
		if finding.AuditDraftFindingStatus != nil {
			// Map AuditDraftFindingStatus to DTO
			findingDTO.AuditDraftFindingStatus = auditDTOs.AuditDraftFindingStatusDTO{
				ID:          finding.AuditDraftFindingStatus.ID,
				TitleEn:     finding.AuditDraftFindingStatus.TitleEn,
				TitleAr:     finding.AuditDraftFindingStatus.TitleAr,
				ColorCode:   finding.AuditDraftFindingStatus.ColorCode,
				Type:        finding.AuditDraftFindingStatus.Type,
				IsFinalized: finding.AuditDraftFindingStatus.IsFinalized,
			}
		}
		// Map documents
		for _, document := range finding.Documents {
			// Map the document model to the DTO
			documentDTO := &auditDTOs.AuditDraftFindingDocumentDTO{
				ID:           &document.ID,
				URL:          &document.URL,
				Name:         document.Name,
				Token:        document.Token,
				OriginalName: document.OriginalName,
				Type:         document.Type,
				Size:         document.Size,
			}
			findingDTO.Documents = append(findingDTO.Documents, documentDTO)
		}
		// for _, conversation := range finding.Conversations {
		// 	findingDTO.Conversations = append(findingDTO.Conversations, ConversationDTO{
		// 		Message: conversation.Message,
		// 	})
		// }

		// // Map action plans
		// for _, actionPlan := range finding.ActionPlans {
		// 	findingDTO.ActionPlans = append(findingDTO.ActionPlans, ActionPlanDTO{
		// 		Plan:           actionPlan.Plan,
		// 		ResponsibleUser: actionPlan.ResponsibleUser,
		// 		TargetDate:     actionPlan.TargetDate.Format("2006-01-02"),
		// 	})
		// }

		auditReportDTO.AuditDraftFindings = append(auditReportDTO.AuditDraftFindings, &findingDTO)
	}

	return auditReportDTO, nil
}

func (adrs *AuditDraftReportService) Create(auditID uint, currentUserID uint, auditDraftReportFormDTO *auditDTOs.AuditDraftReportFormDTO) (uint, error) {
	tx := adrs.repo.BeginTransaction()
	if _, err := adrs.auditSvc.GetByID(auditID); err != nil {
		return 0, err
	}
	dates, err := adrs.parseAuditDraftReportDates(auditDraftReportFormDTO)
	if err != nil {
		return 0, err
	}

	auditDraftReport := &models.AuditDraftReport{
		AuditID:          auditID,
		ExecutiveSummary: auditDraftReportFormDTO.ExecutiveSummary,
		Conclusion:       auditDraftReportFormDTO.Conclusion,
		ReportDate:       *dates.ReportDate,
		CreatedBYID:      currentUserID,
	}

	id, err := adrs.repo.Create(auditDraftReport, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	go func() {
		err := adrs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditDraftReport{}.TableName(),
			id,
			auditDraftReportFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AuditDraftReport{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (adrs *AuditDraftReportService) Update(auditID uint, id uint, currentUserID uint, auditDraftReportFormDTO *auditDTOs.AuditDraftReportFormDTO) error {
	tx := adrs.repo.BeginTransaction()
	if _, err := adrs.auditSvc.GetByID(auditID); err != nil {
		return err
	}
	existingAuditDraftReport, err := adrs.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}

	existingAuditDraftReport.ExecutiveSummary = auditDraftReportFormDTO.ExecutiveSummary
	existingAuditDraftReport.Conclusion = auditDraftReportFormDTO.Conclusion
	existingAuditDraftReport.UpdatedBYID = &currentUserID

	err = adrs.repo.Update(existingAuditDraftReport, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		err := adrs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.AuditDraftReport{}.TableName(),
			id,
			auditDraftReportFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.AuditDraftReport{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (adrs *AuditDraftReportService) Delete(auditID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := adrs.repo.BeginTransaction()
	var failedIDs []uint
	var successfulIDs []uint
	if _, err := adrs.auditSvc.GetByID(auditID); err != nil {
		return nil, err
	}
	for _, id := range ids {
		_, err := adrs.repo.GetFullDetailsByID(auditID, tx)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				failedIDs = append(failedIDs, id)
				continue
			}
			tx.Rollback()
			return nil, err
		}

		if err := adrs.repo.Delete(id, tx); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		successfulIDs = append(successfulIDs, id)
	}

	if len(successfulIDs) > 0 {
		tx.Commit()

		go func() {
			for _, id := range successfulIDs {
				err := adrs.activityLogSvc.LogActivity(
					models.DeleteAction,
					models.AuditDraftReport{}.TableName(),
					id,
					nil,
					currentUserID,
					nil,
				)

				if err != nil {
					fields := map[string]interface{}{
						"action":    models.DeleteAction,
						"tableName": models.AuditDraftReport{}.TableName(),
						"recordID":  id,
						"userID":    currentUserID,
						"error":     err.Error(),
					}

					utility.LogError("Failed to log activity", fields)
				}
			}
		}()
	} else {
		tx.Rollback()
	}

	return failedIDs, nil
}

func (adrs *AuditDraftReportService) parseAuditDraftReportDates(dto *auditDTOs.AuditDraftReportFormDTO) (*auditDTOs.AuditDraftReportParsedDates, error) {
	parsed := &auditDTOs.AuditDraftReportParsedDates{}

	var err error
	parsed.ReportDate, err = utility.ParseDate(dto.ReportDate)
	if err != nil {
		return nil, utility.NewValidationError(fmt.Errorf("invalid report date"), "invalid_date", "report date")
	}

	return parsed, nil
}

// Helper function to map User to UserMinimalDTO
func mapUserToUserMinimalDTO(user *models.User) organiztionDTOs.UserMinimalDTO {
	if user == nil {
		return organiztionDTOs.UserMinimalDTO{}
	}
	return organiztionDTOs.UserMinimalDTO{
		ID:               int(user.ID),
		FirstNameEN:      user.FirstNameEn,
		LastNameEN:       user.LastNameEn,
		FirstNameAR:      user.FirstNameAr,
		LastNameAR:       user.LastNameAr,
		ProfileImagePath: user.ProfileImagePath,
		Email:            user.Email,
		Initial:          user.Initial,
		Color:            user.Color,
	}
}
