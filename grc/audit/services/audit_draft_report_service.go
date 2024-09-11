package services

import (
	"errors"
	"fmt"
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	complianceDTOs "jethur/grc/compliance/dtos"
	documentDTOs "jethur/grc/documents/dtos"
	incidentDTOs "jethur/grc/incident/dtos"
	logServices "jethur/grc/logs/services"
	masterDTOs "jethur/grc/masters/dtos"
	"jethur/grc/models"
	organiztionDTOs "jethur/grc/organization/dtos"
	riskDTOs "jethur/grc/risk/dtos"
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
			// Return nil instead of an error when no data is found
			return nil, nil
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
		PreparedAt:         auditDraftReport.CreatedAt.Format("2006-01-02"),
		ExecutiveSummary:   *auditDraftReport.ExecutiveSummary,
		TeamLead:           mapUserToUserMinimalDTO(auditDraftReport.Audit.AuditPlan.TeamLead),
		Conclusion:         *auditDraftReport.Conclusion,
		Auditors:           []*auditDTOs.AuditPlanAuditMemberDTO{}, // Initialize as empty slice
		AuditDraftFindings: []*auditDTOs.AuditDraftFindingReportDetailDTO{},
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
		findingDTO := auditDTOs.AuditDraftFindingReportDetailDTO{
			ID:                             finding.ID,
			Title:                          finding.Title,
			Description:                    finding.Description,
			AuditID:                        finding.AuditID,
			AuditableItemID:                finding.AuditableItemID,
			AuditChecklistID:               finding.AuditChecklistID,
			AuditFindingCategoryID:         finding.AuditFindingCategoryID,
			Recommendations:                finding.Recommendations,
			AuditDraftFindingStatusID:      finding.AuditDraftFindingStatusID,
			CreatedAt:                      finding.CreatedAt,
			UpdatedAt:                      finding.UpdatedAt,
			Documents:                      []*auditDTOs.AuditDraftFindingDocumentDTO{},
			AuditDraftFindingConversations: []*auditDTOs.AuditDraftFindingConversationDetailDTO{},
			AuditDraftFindingActionPlans:   []*auditDTOs.AuditDraftFindingActionPlanDetailDTO{},
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
		// Map AuditableItem to DTO
		if finding.AuditableItem != nil {
			findingDTO.AuditableItem = auditDTOs.AuditableItemDetailDTO{
				ID:                      finding.AuditableItem.ID,
				Type:                    finding.AuditableItem.Type,
				TitleEn:                 finding.AuditableItem.TitleEn,
				TitleAr:                 finding.AuditableItem.TitleAr,
				Comment:                 &finding.AuditableItem.Comment,
				ProcessID:               finding.AuditableItem.ProcessID,
				ControlID:               finding.AuditableItem.ControlID,
				DocumentID:              finding.AuditableItem.DocumentID,
				IncidentRegisterID:      finding.AuditableItem.IncidentRegisterID,
				ComplianceRequirementID: finding.AuditableItem.ComplianceRequirementID,
				RiskRegisterID:          finding.AuditableItem.RiskRegisterID,
			}
			if finding.AuditableItem.Process != nil {
				findingDTO.AuditableItem.Process = &organiztionDTOs.ProcessDetailDTO{
					TitleEn: safeStringPointer(&finding.AuditableItem.Process.TitleEn),
					TitleAr: safeStringPointer(&finding.AuditableItem.Process.TitleAr),
				}
			}
			if finding.AuditableItem.Control != nil {
				findingDTO.AuditableItem.Control = &organiztionDTOs.ControlDetailDTO{
					TitleEn: safeStringPointer(&finding.AuditableItem.Control.TitleEn),
					TitleAr: safeStringPointer(&finding.AuditableItem.Control.TitleAr),
				}
			}
			if finding.AuditableItem.Document != nil {
				findingDTO.AuditableItem.Document = &documentDTOs.DocumentDetailsDTO{
					TitleEn: safeStringPointer(&finding.AuditableItem.Document.TitleEn),
					TitleAr: safeStringPointer(&finding.AuditableItem.Document.TitleAr),
				}
			}
			if finding.AuditableItem.ComplianceRequirement != nil {
				findingDTO.AuditableItem.ComplianceRequirement = &complianceDTOs.ComplianceRequirementDetailDTO{
					Title: safeStringPointer(&finding.AuditableItem.ComplianceRequirement.Title),
				}
			}
			if finding.AuditableItem.IncidentRegister != nil {
				findingDTO.AuditableItem.IncidentRegister = &incidentDTOs.IncidentRegisterDetailDTO{
					Title: safeStringPointer(&finding.AuditableItem.IncidentRegister.Title),
				}
			}
			if finding.AuditableItem.RiskRegister != nil {
				findingDTO.AuditableItem.RiskRegister = &riskDTOs.RiskRegisterDetailDTO{
					Title: safeStringPointer(&finding.AuditableItem.RiskRegister.Title),
				}
			}
			if finding.AuditableItem.User != nil {
				userDTO := mapUserToUserMinimalDTO(finding.AuditableItem.User)
				findingDTO.AuditableItem.User = &userDTO
			}

		}

		// Map AuditFindingCategory to DTO
		if finding.AuditFindingCategory != nil {
			findingDTO.AuditFindingCategory = auditDTOs.AuditFindingCategoryDTO{
				ID:            int(finding.AuditFindingCategory.ID),
				TitleEn:       finding.AuditFindingCategory.TitleEn,
				TitleAr:       finding.AuditFindingCategory.TitleAr,
				DescriptionEn: *finding.AuditFindingCategory.DescriptionEn,
				DescriptionAr: *finding.AuditFindingCategory.DescriptionAr,
			}
		}

		// Map SeverityLevel to DTO
		if finding.SeverityLevel != nil {
			findingDTO.SeverityLevel = masterDTOs.SeverityLevelDTO{
				ID:        finding.SeverityLevel.ID,
				TitleEn:   finding.SeverityLevel.TitleEn,
				TitleAr:   finding.SeverityLevel.TitleAr,
				ColorCode: finding.SeverityLevel.ColorCode,
				Status:    finding.SeverityLevel.Status,
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
		for _, conversation := range finding.AuditDraftFindingConversations {
			findingDTO.AuditDraftFindingConversations = append(findingDTO.AuditDraftFindingConversations, &auditDTOs.AuditDraftFindingConversationDetailDTO{
				ID:            int(conversation.ID),
				Description:   conversation.Description,
				Type:          conversation.Type,
				CreatedAt:     conversation.CreatedAt,
				CreatedByUser: mapUserToUserMinimalDTO(&conversation.CreatedByUser),
			})
		}

		// Map action plans
		for _, actionPlan := range finding.AuditDraftFindingActionPlans {
			targetDate := actionPlan.TargetDate.Format("2006-01-02")
			findingDTO.AuditDraftFindingActionPlans = append(findingDTO.AuditDraftFindingActionPlans, &auditDTOs.AuditDraftFindingActionPlanDetailDTO{
				ID:               int(actionPlan.ID),
				Title:            actionPlan.Title,
				Description:      actionPlan.Description,
				ResponsibleUser:  mapUserToUserMinimalDTO(actionPlan.ResponsibleUser),
				TargetDateString: &targetDate,
			})
		}

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

// Helper function to safely dereference string pointers
func safeStringPointer(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}
