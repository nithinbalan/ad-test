package services

import (
	"errors"
	"fmt"
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	complianceDTOs "jethur/grc/compliance/dtos"
	documentDTOs "jethur/grc/documents/dtos"
	incidentDTOs "jethur/grc/incident/dtos"
	logServices "jethur/grc/logs/services"
	masterDTOs "jethur/grc/masters/dtos"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationServices "jethur/grc/organization/services"
	riskDTOs "jethur/grc/risk/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
	"jethur/utility"
	"time"

	"gorm.io/gorm"
)

type AuditDraftReportServiceInterface interface {
	GetFullDetailsByID(auditID uint, currentUserID uint, source string, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditDraftReportDetailDTO, error)
	GetByID(id uint) (*models.AuditDraftReport, error)
	GetByAuditID(auditID uint) (*models.AuditDraftReport, error) //used in auditInsights
	Create(auditID uint, currentUserID uint, auditDraftReportFormDTO *auditDTOs.AuditDraftReportFormDTO) (uint, error)
	Update(auditID uint, id uint, currentUserID uint, auditDraftReportFormDTO *auditDTOs.AuditDraftReportFormDTO) error
	Submit(auditID uint, id uint, currentUserID uint) error
	Complete(auditID uint, id uint, currentUserID uint) error
	Delete(auditID uint, currentUserID uint) error
	StatusUpdate(id uint, currentUserID uint, auditDraftReportStatusID uint, nextReviewLevel int, tx utility.DatabaseTransaction) error
}

type AuditDraftReportService struct {
	repo                      repositories.AuditDraftReportRepositoryInterface
	activityLogSvc            logServices.ActivityLogServiceInterface
	auditSvc                  AuditServiceInterface
	auditDraftReportStatusSvc AuditDraftReportStatusServiceInterface
	auditDraftFindingSvc      AuditDraftFindingServiceInterface
	userSvc                   organizationServices.UserServiceInterface
	notificationSvc           commonServices.NotificationServiceInterface
}

func NewAuditDraftReportService(
	repo repositories.AuditDraftReportRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditSvc AuditServiceInterface,
	auditDraftReportStatusSvc AuditDraftReportStatusServiceInterface,
	auditDraftFindingSvc AuditDraftFindingServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) AuditDraftReportServiceInterface {
	return &AuditDraftReportService{
		repo:                      repo,
		activityLogSvc:            activityLogSvc,
		auditSvc:                  auditSvc,
		auditDraftReportStatusSvc: auditDraftReportStatusSvc,
		auditDraftFindingSvc:      auditDraftFindingSvc,
		userSvc:                   userSvc,
		notificationSvc:           notificationSvc,
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

func (adrs *AuditDraftReportService) GetByAuditID(auditID uint) (*models.AuditDraftReport, error) {
	auditDraftReport, err := adrs.repo.GetByAuditID(auditID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, auditID)
		}
		return nil, err
	}
	return auditDraftReport, nil
}

func (adrs *AuditDraftReportService) GetFullDetailsByID(auditID uint, currentUserID uint, source string, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditDraftReportDetailDTO, error) {
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
	createdByUser := mapUserToUserMinimalDTO(auditDraftReport.CreatedByUser)

	var updatedByUser *organizationDTOs.UserMinimalDTO
	var submittedByUser *organizationDTOs.UserMinimalDTO
	var submittedAt *string
	var auditeeReviewAubmittedByUser *organizationDTOs.UserMinimalDTO
	var auditeeReviewAubmittedAt *string

	// Check if SubmittedByUser and SubmittedAt are present before accessing them
	if auditDraftReport.InternalReviewSubmittedByUser != nil {
		submittedByUserDTO := mapUserToUserMinimalDTO(auditDraftReport.InternalReviewSubmittedByUser)
		submittedByUser = &submittedByUserDTO
	}
	if auditDraftReport.UpdatedByUser != nil {
		updatedByUserDTO := mapUserToUserMinimalDTO(auditDraftReport.UpdatedByUser)
		updatedByUser = &updatedByUserDTO
	}
	if auditDraftReport.InternalReviewSubmittedAt != nil {
		formattedSubmittedAt := auditDraftReport.InternalReviewSubmittedAt.Format("2006-01-02")
		submittedAt = &formattedSubmittedAt
	}
	if auditDraftReport.AuditeeReviewSubmittedByUser != nil {
		auditeeSubmittedByUserDTO := mapUserToUserMinimalDTO(auditDraftReport.AuditeeReviewSubmittedByUser)
		auditeeReviewAubmittedByUser = &auditeeSubmittedByUserDTO
	}
	if auditDraftReport.AuditeeReviewSubmittedAt != nil {
		auditeeFormattedSubmittedAt := auditDraftReport.AuditeeReviewSubmittedAt.Format("2006-01-02")
		auditeeReviewAubmittedAt = &auditeeFormattedSubmittedAt
	}
	// Map the AuditDraftReportStatus
	auditDraftReportStatusDTO := auditDTOs.AuditDraftReportStatusDTO{
		ID:        auditDraftReport.AuditDraftReportStatus.ID,
		TitleEn:   auditDraftReport.AuditDraftReportStatus.TitleEn,
		TitleAr:   auditDraftReport.AuditDraftReportStatus.TitleAr,
		Type:      auditDraftReport.AuditDraftReportStatus.Type,
		ColorCode: auditDraftReport.AuditDraftReportStatus.ColorCode,
	}
	// Declare variables to hold the counts for resolved/unresolved comments
	var executiveSummaryResolvedCount, executiveSummaryOpenCount int
	var conclusionResolvedCount, conclusionOpenCount int
	// Map to store the count of resolved and unresolved comments by type
	commentCountByType := make(map[string]map[string]int)

	for _, comment := range auditDraftReport.AuditDraftReportComments {
		// Initialize the inner map for a new type
		if _, exists := commentCountByType[comment.Type]; !exists {
			commentCountByType[comment.Type] = map[string]int{
				"resolved":   0,
				"unresolved": 0,
			}
		}

		// Increment the count based on whether the comment is resolved or not
		if comment.IsResolved {
			commentCountByType[comment.Type]["resolved"]++
		} else {
			commentCountByType[comment.Type]["unresolved"]++
		}
	}

	// Assigning the counts for "executive-summary"
	if count, exists := commentCountByType["executive-summary"]; exists {
		executiveSummaryResolvedCount = count["resolved"]
		executiveSummaryOpenCount = count["unresolved"]
	}

	// Assigning the counts for "conclusion"
	if count, exists := commentCountByType["conclusion"]; exists {
		conclusionResolvedCount = count["resolved"]
		conclusionOpenCount = count["unresolved"]
	}
	// Map the data to the DTO
	auditReportDTO := &auditDTOs.AuditDraftReportDetailDTO{
		ID:                            auditDraftReport.ID,
		AuditID:                       auditDraftReport.AuditID,
		Title:                         auditDraftReport.Audit.AuditPlan.Title,
		ReferenceNumber:               auditDraftReport.Audit.AuditPlan.ReferenceNumber,
		ReportDateString:              auditDraftReport.ReportDate.Format("2006-01-02"),
		PreparedBy:                    mapUserToUserMinimalDTO(auditDraftReport.Audit.CreatedByUser),
		PreparedAt:                    auditDraftReport.CreatedAt.Format("2006-01-02"),
		InternalReviewSubmittedByUser: submittedByUser,
		InternalReviewSubmittedAt:     submittedAt,
		AuditeeReviewSubmittedByUser:  auditeeReviewAubmittedByUser,
		AuditeeReviewSubmittedAt:      auditeeReviewAubmittedAt,
		ExecutiveSummary:              *auditDraftReport.ExecutiveSummary,
		TeamLead:                      mapUserToUserMinimalDTO(auditDraftReport.Audit.AuditPlan.TeamLead),
		Conclusion:                    *auditDraftReport.Conclusion,
		AuditDraftReportStatus:        &auditDraftReportStatusDTO,             // Initialize as empty slice
		Auditors:                      []*auditDTOs.AuditPlanAuditMemberDTO{}, // Initialize as empty slice
		ReviewUsers:                   []*auditDTOs.AuditReportWorkflowParticipantDTO{},
		AuditDraftFindings:            []*auditDTOs.AuditDraftFindingReportDetailDTO{},
		CreatedByUser:                 createdByUser,
		UpdatedByUser:                 updatedByUser,
		ExecutiveSummaryResolvedCount: executiveSummaryResolvedCount,
		ExecutiveSummaryOpenCount:     executiveSummaryOpenCount,
		ConclusionResolvedCount:       conclusionResolvedCount,
		ConclusionOpenCount:           conclusionOpenCount,
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
	// Map review users (AuditReportWorkflowParticipants)
	for _, workflow := range auditDraftReport.AuditReportWorkflows {
		for _, participant := range workflow.AuditReportWorkflowParticipants {
			auditReportDTO.ReviewUsers = append(auditReportDTO.ReviewUsers, &auditDTOs.AuditReportWorkflowParticipantDTO{
				ID:                    participant.ID,
				AuditReportWorkflowID: participant.AuditReportWorkflowID,
				UserID:                participant.UserID,
				User:                  mapUserToUserMinimalDTO(participant.User),
			})
		}
	}
	if source == "audit-insight" {
		orgUnitID, err := adrs.userSvc.GetUserPrimaryUnitID(currentUserID)
		if err != nil {
			return nil, err
		}
		for _, finding := range auditDraftReport.Audit.AuditDraftFindings {
			if finding.AuditableItem != nil && finding.AuditableItem.ResponsibleUnit != nil && finding.AuditableItem.ResponsibleUnit.ID == orgUnitID {
				var acceptedByUser *organizationDTOs.UserMinimalDTO = nil
				var challengedByUser *organizationDTOs.UserMinimalDTO = nil
				if finding.AcceptedByUser != nil {
					acceptedByUserDTO := mapUserToUserMinimalDTO(finding.AcceptedByUser)
					acceptedByUser = &acceptedByUserDTO
				}
				if finding.ChallengedByUser != nil {
					challengedByUserDTO := mapUserToUserMinimalDTO(finding.ChallengedByUser)
					challengedByUser = &challengedByUserDTO
				}
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
					AcceptedByUser:                 acceptedByUser,
					AcceptedAt:                     finding.AcceptedAt,
					ChallengedByUser:               challengedByUser,
					ChallengedAt:                   finding.ChallengedAt,
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
					ResponsibleUnitID:       finding.AuditableItem.ResponsibleUnitID,
					ResponsibleUserID:       finding.AuditableItem.ResponsibleUserID,
				}
				if finding.AuditableItem.Process != nil {
					findingDTO.AuditableItem.Process = &organizationDTOs.ProcessDetailDTO{
						TitleEn: safeStringPointer(&finding.AuditableItem.Process.TitleEn),
						TitleAr: safeStringPointer(&finding.AuditableItem.Process.TitleAr),
					}
				}
				if finding.AuditableItem.Control != nil {
					findingDTO.AuditableItem.Control = &organizationDTOs.ControlDetailDTO{
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
				if finding.AuditableItem.ResponsibleUser != nil {
					userDTO := mapUserToUserMinimalDTO(finding.AuditableItem.ResponsibleUser)
					findingDTO.AuditableItem.ResponsibleUser = &userDTO
				}
				if finding.AuditableItem.ResponsibleUnit != nil {
					findingDTO.AuditableItem.ResponsibleUnit = settingsDTOs.OrganizationUnitDTO{
						ID:        int(finding.AuditableItem.ResponsibleUnit.ID),
						TitleEn:   finding.AuditableItem.ResponsibleUnit.TitleEn,
						TitleAr:   finding.AuditableItem.ResponsibleUnit.TitleAr,
						ColorCode: finding.AuditableItem.ResponsibleUnit.ColorCode,
						Status:    finding.AuditableItem.ResponsibleUnit.Status,
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
						Comment:       conversation.Comment,
						Type:          conversation.Type,
						CreatedAt:     conversation.CreatedAt,
						CreatedByUser: mapUserToUserMinimalDTO(&conversation.CreatedByUser),
					})
				}

				// Map action plans
				for _, actionPlan := range finding.AuditDraftFindingActionPlans {
					targetDate := actionPlan.TargetDate.Format("2006-01-02")
					findingDTO.AuditDraftFindingActionPlans = append(findingDTO.AuditDraftFindingActionPlans, &auditDTOs.AuditDraftFindingActionPlanDetailDTO{
						ID:                  int(actionPlan.ID),
						AuditDraftFindingID: actionPlan.AuditDraftFindingID,
						Title:               actionPlan.Title,
						Description:         actionPlan.Description,
						ResponsibleUser:     mapUserToUserMinimalDTO(actionPlan.ResponsibleUser),
						TargetDateString:    &targetDate,
					})
				}
				auditReportDTO.AuditDraftFindings = append(auditReportDTO.AuditDraftFindings, &findingDTO)
			}
		}
	} else {
		// Map findings
		for _, finding := range auditDraftReport.Audit.AuditDraftFindings {
			var acceptedByUser *organizationDTOs.UserMinimalDTO = nil
			var challengedByUser *organizationDTOs.UserMinimalDTO = nil
			if finding.AcceptedByUser != nil {
				acceptedByUserDTO := mapUserToUserMinimalDTO(finding.AcceptedByUser)
				acceptedByUser = &acceptedByUserDTO
			}
			if finding.ChallengedByUser != nil {
				challengedByUserDTO := mapUserToUserMinimalDTO(finding.ChallengedByUser)
				challengedByUser = &challengedByUserDTO
			}
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
				AcceptedByUser:                 acceptedByUser,
				AcceptedAt:                     finding.AcceptedAt,
				ChallengedByUser:               challengedByUser,
				ChallengedAt:                   finding.ChallengedAt,
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
					ResponsibleUnitID:       finding.AuditableItem.ResponsibleUnitID,
					ResponsibleUserID:       finding.AuditableItem.ResponsibleUserID,
				}
				if finding.AuditableItem.Process != nil {
					findingDTO.AuditableItem.Process = &organizationDTOs.ProcessDetailDTO{
						TitleEn: safeStringPointer(&finding.AuditableItem.Process.TitleEn),
						TitleAr: safeStringPointer(&finding.AuditableItem.Process.TitleAr),
					}
				}
				if finding.AuditableItem.Control != nil {
					findingDTO.AuditableItem.Control = &organizationDTOs.ControlDetailDTO{
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
				if finding.AuditableItem.ResponsibleUser != nil {
					userDTO := mapUserToUserMinimalDTO(finding.AuditableItem.ResponsibleUser)
					findingDTO.AuditableItem.ResponsibleUser = &userDTO
				}
				if finding.AuditableItem.ResponsibleUnit != nil {
					findingDTO.AuditableItem.ResponsibleUnit = settingsDTOs.OrganizationUnitDTO{
						ID:        int(finding.AuditableItem.ResponsibleUnit.ID),
						TitleEn:   finding.AuditableItem.ResponsibleUnit.TitleEn,
						TitleAr:   finding.AuditableItem.ResponsibleUnit.TitleAr,
						ColorCode: finding.AuditableItem.ResponsibleUnit.ColorCode,
						Status:    finding.AuditableItem.ResponsibleUnit.Status,
					}
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
					Comment:       conversation.Comment,
					Type:          conversation.Type,
					CreatedAt:     conversation.CreatedAt,
					CreatedByUser: mapUserToUserMinimalDTO(&conversation.CreatedByUser),
				})
			}

			// Map action plans
			for _, actionPlan := range finding.AuditDraftFindingActionPlans {
				targetDate := actionPlan.TargetDate.Format("2006-01-02")
				findingDTO.AuditDraftFindingActionPlans = append(findingDTO.AuditDraftFindingActionPlans, &auditDTOs.AuditDraftFindingActionPlanDetailDTO{
					ID:                  int(actionPlan.ID),
					AuditDraftFindingID: actionPlan.AuditDraftFindingID,
					Title:               actionPlan.Title,
					Description:         actionPlan.Description,
					ResponsibleUser:     mapUserToUserMinimalDTO(actionPlan.ResponsibleUser),
					TargetDateString:    &targetDate,
				})
			}

			auditReportDTO.AuditDraftFindings = append(auditReportDTO.AuditDraftFindings, &findingDTO)
		}
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
	auditDraftReportStatus, err := adrs.auditDraftReportStatusSvc.GetByType("draft")
	if err != nil {
		return 0, err
	}

	auditDraftReport := &models.AuditDraftReport{
		AuditID:                  auditID,
		ExecutiveSummary:         auditDraftReportFormDTO.ExecutiveSummary,
		Conclusion:               auditDraftReportFormDTO.Conclusion,
		ReportDate:               *dates.ReportDate,
		AuditDraftReportStatusID: auditDraftReportStatus.ID,
		CreatedBYID:              currentUserID,
	}

	id, err := adrs.repo.Create(auditDraftReport, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	go func() {
		existingAuditDraftReport, err := adrs.repo.GetFullDetailsByID(auditID)
		if err != nil {
			return
		}

		err = adrs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditDraftReport{}.TableName(),
			id,
			existingAuditDraftReport,
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
	oldAuditDraftReport, err := adrs.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
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
		existingAuditDraftReport, err := adrs.repo.GetByID(id)
		if err != nil {
			return
		}
		// Prepare log data
		changedData := make(map[string]interface{})
		//log changes
		oldData := map[string]interface{}{
			"executive_summary": oldAuditDraftReport.ExecutiveSummary,
			"conclusion":        oldAuditDraftReport.Conclusion,
		}

		newData := map[string]interface{}{
			"executive_summary": existingAuditDraftReport.ExecutiveSummary,
			"conclusion":        existingAuditDraftReport.Conclusion,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := adrs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.AuditDraftReport{}.TableName(),
				id,
				changedData,
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
		}
	}()

	return nil
}

func (adrs *AuditDraftReportService) Delete(auditID uint, currentUserID uint) error {
	tx := adrs.repo.BeginTransaction()

	var deletedAuditDraftReport []map[string]interface{}

	auditDraftReport, err := adrs.repo.GetFullDetailsByID(auditID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := adrs.repo.Delete(auditDraftReport.ID, tx); err != nil {
		tx.Rollback()
		return err
	}

	deletedAuditDraftReport = append(deletedAuditDraftReport, map[string]interface{}{
		"id":         auditDraftReport.ID,
		"audit_plan": auditDraftReport.Audit.AuditPlan.Title,
		"audit_id":   auditDraftReport.AuditID,
	})

	tx.Commit()

	go func() {
		err := adrs.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.AuditDraftReport{}.TableName(),
			auditDraftReport.ID,
			deletedAuditDraftReport,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.DeleteAction,
				"tableName": models.AuditDraftReport{}.TableName(),
				"recordID":  auditDraftReport.ID,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (adrs *AuditDraftReportService) Submit(auditID uint, id uint, currentUserID uint) error {
	tx := adrs.repo.BeginTransaction()
	if _, err := adrs.auditSvc.GetByID(auditID); err != nil {
		return err
	}
	existingAuditDraftReport, err := adrs.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}

	auditDraftReportStatusApproved, err := adrs.auditDraftReportStatusSvc.GetByType("approved")
	if err != nil {
		return err
	}
	if existingAuditDraftReport.AuditDraftReportStatusID != auditDraftReportStatusApproved.ID {
		tx.Rollback()
		return utility.NewForbiddenError("The audit draft report has not been approved.")
	}

	auditDraftReportStatus, err := adrs.auditDraftReportStatusSvc.GetByType("in-department-review")
	if err != nil {
		return err
	}
	now := time.Now()
	existingAuditDraftReport.AuditeeReviewSubmittedAt = &now
	existingAuditDraftReport.AuditDraftReportStatusID = auditDraftReportStatus.ID
	existingAuditDraftReport.AuditeeReviewSubmittedBYID = &currentUserID
	existingAuditDraftReport.UpdatedBYID = &currentUserID
	existingAuditDraftReport.UpdatedAt = &now

	err = adrs.repo.Update(existingAuditDraftReport, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	//todo - notification to auditable item responsible user where findings captured.

	tx.Commit()

	go func() {
		urlParam := fmt.Sprintf(`{"id":%d}`, id)
		entity := "audit-draft-report"
		action := "review-user-assign"

		existingAuditDraftReportFullDetails, err := adrs.repo.GetByAuditID(auditID)
		if err != nil {
			utility.LogError("Failed to fetch review user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		if existingAuditDraftReportFullDetails.Audit.AuditPlan.AuditableItems != nil {
			for _, item := range existingAuditDraftReportFullDetails.Audit.AuditPlan.AuditableItems {
				if item.AuditableItem.ResponsibleUser == nil {
					utility.LogError("Failed to fetch assigning user", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}
				responsibleUser := item.AuditableItem.ResponsibleUser
				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserID:               responsibleUser.ID,
					ModuleIdentifier:     "audit",
					NotificationType:     string(commonDTOs.ReviewType),
					HeadingEn:            "Audit Draft Report Department Review",
					HeadingAr:            "مراجعة قسم تقرير مسودة التدقيق",
					TitleEn:              existingAuditDraftReport.Audit.AuditPlan.Title,
					TitleAr:              existingAuditDraftReport.Audit.AuditPlan.Title,
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				assignedBy, err := adrs.userSvc.GetByID(currentUserID)
				if err != nil {
					utility.LogError("Failed to fetch assigning user", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}

				emailTo := []string{responsibleUser.Email}
				subject := fmt.Sprintf("Audit Draft Report /%s: Department Review Assigned", existingAuditDraftReport.Audit.AuditPlan.Title)
				appURL := utility.GetAppURL()
				emailData := map[string]interface{}{
					"FirstNameEn":           responsibleUser.FirstNameEn,
					"LastNameEn":            responsibleUser.LastNameEn,
					"AuditTitle":            existingAuditDraftReport.Audit.AuditPlan.Title,
					"Reference Number":      existingAuditDraftReport.Audit.AuditPlan.ReferenceNumber,
					"AssignedByFirstNameEn": assignedBy.FirstNameEn,
					"AssignedByLastNameEn":  assignedBy.LastNameEn,
					"AssignedOn":            time.Now().Format("2006-01-02"),
					"RedirectURL":           fmt.Sprintf("%s/audit-insight/draft-report/%d/details", appURL, id),
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              emailTo,
					Subject:              subject,
					TemplateRelativePath: "audit/templates/email/audit_draft_report_department_review_user_mail.html",
					EmailData:            emailData,
				}

				// Send Notification and Email
				err = adrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "audit")
				if err != nil {
					utility.LogError("Failed to send notifications", map[string]interface{}{
						"error": err.Error(),
					})
				}
			}
		}

		err = adrs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.AuditDraftReport{}.TableName(),
			id,
			nil,
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

// Status update mainly using for review - AuditDraftReportReviewService
func (adrs *AuditDraftReportService) StatusUpdate(id uint, currentUserID uint, auditDraftReportStatusID uint, nextReviewLevel int, tx utility.DatabaseTransaction) error {
	existingAuditDraftReport, err := adrs.GetByID(id)
	if err != nil {
		return err
	}

	now := time.Now()

	existingAuditDraftReport.InternalReviewSubmittedAt = &now
	existingAuditDraftReport.AuditDraftReportStatusID = auditDraftReportStatusID
	existingAuditDraftReport.InternalReviewSubmittedBYID = &currentUserID
	existingAuditDraftReport.NextReviewLevel = nextReviewLevel
	existingAuditDraftReport.UpdatedBYID = &currentUserID

	err = adrs.repo.Update(existingAuditDraftReport, tx)

	if err != nil {
		return err
	}

	return nil
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
func mapUserToUserMinimalDTO(user *models.User) organizationDTOs.UserMinimalDTO {
	if user == nil {
		return organizationDTOs.UserMinimalDTO{}
	}
	return organizationDTOs.UserMinimalDTO{
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

func (adrs *AuditDraftReportService) Complete(auditID uint, id uint, currentUserID uint) error {
	tx := adrs.repo.BeginTransaction()
	if _, err := adrs.auditSvc.GetByID(auditID); err != nil {
		return err
	}
	existingAuditDraftReport, err := adrs.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}

	auditDraftReportStatus, err := adrs.auditDraftReportStatusSvc.GetByType("in-department-review")
	if err != nil {
		return err
	}

	if existingAuditDraftReport.AuditDraftReportStatusID != auditDraftReportStatus.ID {
		tx.Rollback()
		return utility.NewForbiddenError("The audit draft report has not been in department review.")
	}
	auditDraftFindings, err := adrs.auditDraftFindingSvc.GetAllFindingsByAuditID(existingAuditDraftReport.AuditID)
	if err != nil {
		return err
	}
	// Check if all findings are finalized
	for _, finding := range auditDraftFindings {
		// Ensure the status of each finding is finalized
		if finding.AuditDraftFindingStatus != nil && !finding.AuditDraftFindingStatus.IsFinalized {
			tx.Rollback()
			return utility.NewForbiddenError(fmt.Sprintf("Finding '%s' is not finalized.", finding.Title))
		}
	}
	auditDraftReportStatusCompleted, err := adrs.auditDraftReportStatusSvc.GetByType("completed")
	if err != nil {
		return err
	}
	now := time.Now()
	existingAuditDraftReport.AuditDraftReportStatusID = auditDraftReportStatusCompleted.ID
	existingAuditDraftReport.UpdatedBYID = &currentUserID
	existingAuditDraftReport.UpdatedAt = &now

	err = adrs.repo.Update(existingAuditDraftReport, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	//todo - notification to auditable item responsible user where findings captured.

	tx.Commit()

	go func() {
		err := adrs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.AuditDraftReport{}.TableName(),
			id,
			nil,
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
