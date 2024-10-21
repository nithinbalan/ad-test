package services

import (
	"errors"
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
		return 0, err
	}
	// Check if a final report already exists for this audit
	existingReport, err := afrs.repo.GetByAuditID(auditID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// If an existing report is found, do not allow creation
	if existingReport != nil {
		tx.Rollback()
		return 0, utility.NewForbiddenError("An audit final report already exists for this audit.")
	}
	auditDraftReport, err := afrs.auditDraftReportSvc.GetByAuditID(auditID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}


	// Check if AuditDraftReportStatus is 'completed'
	completedStatus, err := afrs.auditDraftReportStatusSvc.GetByType("completed")
	if err != nil {
		return 0, err
	}
	if auditDraftReport.AuditDraftReportStatusID != completedStatus.ID {
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
		return 0, err
	}
	referenceNumber, err := afrs.repo.GenerateReferenceNumber()
	if err != nil {
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
	var findings []*models.AuditFinding

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
			// Add a general "Audit Procedure" section first
			reportSubSections = append(reportSubSections, &models.AuditFinalReportSection{
				AuditFinalReportID:        id,
				AuditFinalReportSectionID: &mainSection.ID,
				TitleEn:                   "Scope of the Audit",
				TitleAr:                   "ملخص تنفيذي",
				// DescriptionEn:             auditDraftReport.Audit.AuditPlan.Purpose,
				// DescriptionAr:             auditDraftReport.Audit.AuditPlan.Purpose,
				Type:                      section.Type,
				Order:                     float64Ptr(subOrder + 0.1),
				CreatedBYID:               currentUserID,
			})
		
			// Add a section for each auditable item
			subOrder += 0.1 
			for _, auditableItem := range auditDraftReport.Audit.AuditPlan.AuditableItems {
				// Create a section for each auditable item
				reportSubSections = append(reportSubSections, &models.AuditFinalReportSection{
					AuditFinalReportID:        id,
					AuditFinalReportSectionID: &mainSection.ID, 
					TitleEn:                   "Scope of the Audit",
					TitleAr:                   "البند القابل للتدقيق ",
					DescriptionEn:             auditableItem.AuditableItem.TitleEn, 
					DescriptionAr:             auditableItem.AuditableItem.TitleAr, 
					Type:                      section.Type,
					Order:                     float64Ptr(subOrder + 0.1),
					CreatedBYID:               currentUserID,
				})
				subOrder += 0.1
			}
		}
		

		if section.Type == "findings" {
			auditDraftFindings, err := afrs.auditDraftFindingSvc.GetAllFindingsByAuditID(auditID)
			if err != nil {
				tx.Rollback()
				return 0, err
			}
			for _, draftFinding := range auditDraftFindings {
				if draftFinding.AuditDraftFindingStatus != nil && draftFinding.AuditDraftFindingStatus.IsFinalized {
					finding := &models.AuditFinding{
						AuditID:                draftFinding.AuditID,
						Title:                  draftFinding.Title,
						Description:            draftFinding.Description,
						AuditableItemID:        *draftFinding.AuditableItemID,
						AuditFindingCategoryID: draftFinding.AuditFindingCategoryID,
						SeverityLevelID:        *draftFinding.SeverityLevelID,
						Recommendations:        draftFinding.Recommendations,
						AuditFindingStatusID:   draftFinding.AuditDraftFindingStatusID,
						CreatedBYID:            currentUserID,
					}
					findings = append(findings, finding)

					reportSubSections = append(reportSubSections, &models.AuditFinalReportSection{
						AuditFinalReportID:        id,
						AuditFinalReportSectionID: &mainSection.ID,
						TitleEn:                   draftFinding.Title,
						TitleAr:                   draftFinding.Title,
						DescriptionEn:             draftFinding.Description,
						DescriptionAr:             draftFinding.Description,
						Type:                      section.Type,
						Order:                     float64Ptr(subOrder + 0.1),
						CreatedBYID:               currentUserID,
					})
				}
			}
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

	//Bulk insert findings
	if len(findings) > 0 {
		if err := afrs.repo.CreateAuditFindings(findings, tx); err != nil {
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
