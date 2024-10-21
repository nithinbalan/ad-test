package services

import (
	"errors"
	"fmt"
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditFinalReportServiceInterface interface {
	Create(auditID uint, currentUserID uint) (uint, error)
	GetByAuditID(auditID uint) (*models.AuditFinalReport, error)
}

type AuditFinalReportService struct {
	repo                      repositories.AuditFinalReportRepositoryInterface
	activityLogSvc            logServices.ActivityLogServiceInterface
	auditSvc                  AuditServiceInterface
	auditDraftReportStatusSvc AuditDraftReportStatusServiceInterface
	auditDraftReportSvc       AuditDraftReportServiceInterface
	auditDraftFindingSvc      AuditDraftFindingServiceInterface
	userSvc                   organizationServices.UserServiceInterface
	auditFinalReportStatusSvc AuditFinalReportStatusServiceInterface
	auditFindingSvc           AuditFindingServiceInterface
}

func NewAuditFinalReportService(
	repo repositories.AuditFinalReportRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditSvc AuditServiceInterface,
	auditDraftReportStatusSvc AuditDraftReportStatusServiceInterface,
	auditDraftReportSvc AuditDraftReportServiceInterface,
	auditDraftFindingSvc AuditDraftFindingServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	auditFinalReportStatusSvc AuditFinalReportStatusServiceInterface,
	auditFindingSvc AuditFindingServiceInterface,
) AuditFinalReportServiceInterface {
	return &AuditFinalReportService{
		repo:                      repo,
		activityLogSvc:            activityLogSvc,
		auditSvc:                  auditSvc,
		auditDraftReportStatusSvc: auditDraftReportStatusSvc,
		auditDraftReportSvc:       auditDraftReportSvc,
		auditDraftFindingSvc:      auditDraftFindingSvc,
		userSvc:                   userSvc,
		auditFinalReportStatusSvc: auditFinalReportStatusSvc,
		auditFindingSvc:           auditFindingSvc,
	}
}

func (afrs *AuditFinalReportService) GetByAuditID(auditID uint) (*models.AuditFinalReport, error) {
	auditFinalReport, err := afrs.repo.GetByAuditID(auditID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, auditID)
		}
		return nil, err
	}
	return auditFinalReport, nil
}

func (afrs *AuditFinalReportService) Create(auditID uint, currentUserID uint) (uint, error) {
	tx := afrs.repo.BeginTransaction()
	if _, err := afrs.auditSvc.GetByID(auditID); err != nil {
		tx.Rollback()
		return 0, err
	}
	// Check if a final report already exists for this audit
	// existingReport, err := afrs.repo.GetByAuditID(auditID)
	// if err != nil {
	//     tx.Rollback()
	//     return 0, err
	// }

	// // If an existing report is found, rollback and return an error
	// if existingReport != nil {
	//     tx.Rollback()
	//     return 0, utility.NewForbiddenError("An audit final report already exists for this audit.")
	// }
	auditDraftReport, err := afrs.auditDraftReportSvc.GetByAuditID(auditID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Check if AuditDraftReportStatus is 'completed'
	completedStatus, err := afrs.auditDraftReportStatusSvc.GetByType("completed")
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if auditDraftReport.AuditDraftReportStatusID != completedStatus.ID {
		tx.Rollback()
		return 0, utility.NewForbiddenError("Action not allowed: Only audit draft reports with an 'Completed' status can be accepted.")
	}
	audit, err := afrs.auditSvc.GetFullDetailsByID(auditID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Check if the current user is the TeamLead
	isAuthorizedUser := false
	if uint(audit.AuditPlan.TeamLead.ID) == currentUserID {
		isAuthorizedUser = true
	}

	// Check if the current user is one of the AuditPlan Members
	if !isAuthorizedUser {
		for _, member := range audit.AuditPlan.Members {
			if member.UserID == currentUserID {
				isAuthorizedUser = true
				break
			}
		}
	}

	// If the current user is not authorized, rollback and return forbidden error
	if !isAuthorizedUser {
		tx.Rollback()
		return 0, utility.NewForbiddenError("The current user is not authorized to perform this action.")
	}
	auditFinalReportStatus, err := afrs.auditFinalReportStatusSvc.GetByType("draft")
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	referenceNumber, err := afrs.repo.GenerateReferenceNumber()
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	auditFinalReport := &models.AuditFinalReport{
		AuditID:                  auditID,
		ReferenceNumber:          referenceNumber,
		AuditFinalReportStatusID: auditFinalReportStatus.ID,
		CreatedBYID:              currentUserID,
	}

	id, err := afrs.repo.Create(auditFinalReport, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	// Mapping main section titles to types
	mainSections := []struct {
		Title string
		Type  string
	}{
		{"Audit Report Overview", "overview"},
		{"Audit Findings", "finding-summary"},
		{"Findings", "findings"},
		{"Conclusions", "conclusion"},
		{"Appendix", "custom"},
	}

	var reportSections []*models.AuditFinalReportSection
	order := 1
	// Step 1: Insert main sections first
	for _, sectionInfo := range mainSections {
		// Create a new section for each main section
		section := &models.AuditFinalReportSection{
			AuditFinalReportID: id,
			TitleEn:            sectionInfo.Title,
			TitleAr:            sectionInfo.Title,
			Type:               sectionInfo.Type,
			Order:              float64Ptr(float64(order)),
			CreatedBYID:        currentUserID,
		}

		// Add the section to the slice
		reportSections = append(reportSections, section)
		order++
	}
	// Insert the main sections into the database using your repository method
	if err := afrs.repo.CreateAuditFinalReportSections(reportSections, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	// Step 2: For each main section, insert sub-sections based on their type
	var reportSubSections []*models.AuditFinalReportSection

	for _, section := range reportSections {
		mainSection, err := afrs.repo.GetSectionByType(section.Type, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		subOrder := *section.Order // Use the main section order as the base
		if section.Type == "overview" && auditDraftReport.ExecutiveSummary != nil {
			reportSubSections = append(reportSubSections, &models.AuditFinalReportSection{
				AuditFinalReportID:        id,
				AuditFinalReportSectionID: &mainSection.ID,
				TitleEn:                   "Executive Summary",
				TitleAr:                   "ملخص تنفيذي",
				DescriptionEn:             auditDraftReport.ExecutiveSummary,
				DescriptionAr:             auditDraftReport.ExecutiveSummary,
				Type:                      section.Type,
				Order:                     float64Ptr(subOrder + 0.1),
				CreatedBYID:               currentUserID,
			})
		}
		if section.Type == "overview" && auditDraftReport.Audit.AuditPlan != nil {
			reportSubSections = append(reportSubSections, &models.AuditFinalReportSection{
				AuditFinalReportID:        id,
				AuditFinalReportSectionID: &mainSection.ID,
				TitleEn:                   "Audit Objective",
				TitleAr:                   "ملخص تنفيذي",
				DescriptionEn:             auditDraftReport.Audit.AuditPlan.Objective,
				DescriptionAr:             auditDraftReport.Audit.AuditPlan.Objective,
				Type:                      section.Type,
				Order:                     float64Ptr(subOrder + 0.1),
				CreatedBYID:               currentUserID,
			})
		}
		if section.Type == "overview" && auditDraftReport.Audit.AuditPlan != nil {
			reportSubSections = append(reportSubSections, &models.AuditFinalReportSection{
				AuditFinalReportID:        id,
				AuditFinalReportSectionID: &mainSection.ID,
				TitleEn:                   "Audit Procedure",
				TitleAr:                   "ملخص تنفيذي",
				DescriptionEn:             auditDraftReport.Audit.AuditPlan.Purpose,
				DescriptionAr:             auditDraftReport.Audit.AuditPlan.Purpose,
				Type:                      section.Type,
				Order:                     float64Ptr(subOrder + 0.1),
				CreatedBYID:               currentUserID,
			})
		}
		if section.Type == "overview" && auditDraftReport.Audit.AuditPlan.AuditableItems != nil {
			// Prepare an HTML list for the auditable items
			auditableItemsListEn := "<ul>"
			auditableItemsListAr := "<ul>"

			// Helper function to add items to the list
			addToList := func(titleEn, titleAr string) {
				auditableItemsListEn += fmt.Sprintf("<li>%s</li>", titleEn)
				auditableItemsListAr += fmt.Sprintf("<li>%s</li>", titleAr)
			}

			for _, auditableItem := range auditDraftReport.Audit.AuditPlan.AuditableItems {
				titleEn, titleAr := "", ""

				switch auditableItem.AuditableItem.Type {
				case "process":
					titleEn = auditableItem.AuditableItem.Process.TitleEn
					titleAr = auditableItem.AuditableItem.Process.TitleAr

				case "control":
					titleEn = auditableItem.AuditableItem.Control.TitleEn
					titleAr = auditableItem.AuditableItem.Control.TitleAr

				case "document":
					titleEn = auditableItem.AuditableItem.Document.TitleEn
					titleAr = auditableItem.AuditableItem.Document.TitleAr

				case "incident":
					titleEn = auditableItem.AuditableItem.IncidentRegister.Title
					titleAr = auditableItem.AuditableItem.IncidentRegister.Title

				case "compliance_requirement":
					titleEn = auditableItem.AuditableItem.ComplianceRequirement.Title
					titleAr = auditableItem.AuditableItem.ComplianceRequirement.Title

				case "risk":
					titleEn = auditableItem.AuditableItem.RiskRegister.Title
					titleAr = auditableItem.AuditableItem.RiskRegister.Title

				case "user":
					titleEn = auditableItem.AuditableItem.User.FirstNameEn
					titleAr = auditableItem.AuditableItem.User.FirstNameAr

				case "custom":
					if auditableItem.AuditableItem.TitleEn != nil && auditableItem.AuditableItem.TitleAr != nil {
						titleEn = *auditableItem.AuditableItem.TitleEn
						titleAr = *auditableItem.AuditableItem.TitleAr
					}
				}

				// Only add to the list if titles are not empty
				if titleEn != "" || titleAr != "" {
					addToList(titleEn, titleAr)
				}
			}

			// Close the HTML lists
			auditableItemsListEn += "</ul>"
			auditableItemsListAr += "</ul>"

			// Add a single section with the full list of auditable items in both languages
			reportSubSections = append(reportSubSections, &models.AuditFinalReportSection{
				AuditFinalReportID:        id,
				AuditFinalReportSectionID: &mainSection.ID,
				TitleEn:                   "Scope of the Audit",
				TitleAr:                   "نطاق التدقيق",
				DescriptionEn:             &auditableItemsListEn,
				DescriptionAr:             &auditableItemsListAr,
				Type:                      section.Type,
				Order:                     float64Ptr(subOrder + 0.1),
				CreatedBYID:               currentUserID,
			})

		}

		if section.Type == "conclusion" && auditDraftReport.Conclusion != nil {
			reportSubSections = append(reportSubSections, &models.AuditFinalReportSection{
				AuditFinalReportID:        id,
				AuditFinalReportSectionID: &mainSection.ID,
				TitleEn:                   "Concluding Statement",
				TitleAr:                   "خاتمة",
				DescriptionEn:             auditDraftReport.Conclusion,
				DescriptionAr:             auditDraftReport.Conclusion,
				Type:                      section.Type,
				Order:                     float64Ptr(subOrder + 0.1),
				CreatedBYID:               currentUserID,
			})
		}
	}
	// findings only copying from draft findings to AuditFinding, not saving in to AuditFinalReportSection
	auditDraftFindings, err := afrs.auditDraftFindingSvc.GetAllFindingsByAuditID(auditID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	// Initialize slices for findings and action plans
	var actionPlans []*models.AuditFindingActionPlan

	for _, draftFinding := range auditDraftFindings {
		if draftFinding.AuditDraftFindingStatus != nil && draftFinding.AuditDraftFindingStatus.IsFinalized {
			// Create AuditFinding
			findingDTO := &auditDTOs.AuditFindingFormDTO{
				AuditID:                draftFinding.AuditID,
				Title:                  draftFinding.Title,
				Description:            draftFinding.Description,
				AuditableItemID:        *draftFinding.AuditableItemID,
				AuditFindingCategoryID: draftFinding.AuditFindingCategoryID,
				SeverityLevelID:        *draftFinding.SeverityLevelID,
				Recommendations:        draftFinding.Recommendations,
			}
			findingID, err := afrs.auditFindingSvc.Create(currentUserID, findingDTO)
			if err != nil {
				tx.Rollback()
				return 0, err
			}
			// Copy action plans from draft findings to final findings
			for _, draftActionPlan := range draftFinding.AuditDraftFindingActionPlans {
				actionPlan := &models.AuditFindingActionPlan{
					Title:             draftActionPlan.Title,
					AuditFindingID:    findingID,
					Description:       &draftActionPlan.Description,
					ResponsibleUserID: *draftActionPlan.ResponsibleUserID,
					TargetDate:        draftActionPlan.TargetDate,
					CreatedBYID:       currentUserID,
					// SeverityLevelID:    draftActionPlan.SeverityLevelID,
					ActionPlanStatusID: 1, //todo -- add ActionPlanStatusID
				}
				// Attach action plan to its respective AuditFinding
				actionPlans = append(actionPlans, actionPlan)
			}
		}
	}

	// Bulk insert findings
	if len(actionPlans) > 0 {
		// Bulk insert action plans
		if err := afrs.repo.CreateAuditFindingActionPlans(actionPlans, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	//Bulk insert AuditFinalReportSection
	if len(reportSubSections) > 0 {
		if err := afrs.repo.CreateAuditFinalReportSections(reportSubSections, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	tx.Commit()

	go func() {
		err := afrs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditFinalReport{}.TableName(),
			id,
			nil,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AuditFinalReport{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

// Helper function to return a pointer to a float64 value
func float64Ptr(val float64) *float64 {
	return &val
}
