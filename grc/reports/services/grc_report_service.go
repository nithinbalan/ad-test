package services

import (
	"encoding/json"
	"errors"
	"fmt"
	kpiDTOs "jethur/grc/kpi/dtos"
	kpiServices "jethur/grc/kpi/services"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/reports/dtos"
	"jethur/grc/reports/repositories"
	riskServices "jethur/grc/risk/services"
	settingServices "jethur/grc/settings/services"
	"jethur/utility"
	"math"
	"strings"

	"gorm.io/gorm"
)

type GRCReportServiceInterface interface {
	GetAll(params *dtos.GRCReportListQueryParamsDTO) (*dtos.GRCReportListResponseDTO, error)
	GetGRCReportDetails(grcReportID uint) (*dtos.GRCReportDetailDTO, error)
	Create(currentUserID uint, grcReportFormDTO *dtos.GRCReportFormDTO) (uint, error)
	GetSectionByID(grcReportID uint, sectionID uint, tx ...utility.DatabaseTransaction) (*models.GRCReportSection, error)
	UpdateSection(grcID uint, id uint, currentUserID uint, grcReporSectionFormDTO *dtos.GRCReportSectionFormDTO) error
	DeleteGRCReportSection(grcReportID uint, sectionID uint, currentUserID uint) error
	DeleteGRCReport(grcReportID uint, currentUserID uint) error
}

type GRCReportService struct {
	repo                    repositories.GRCReportRepositoryInterface
	activityLogSvc          logServices.ActivityLogServiceInterface
	organizationLevelSvc    settingServices.OrganizationLevelServiceInterface
	riskResponseStrategySvc riskServices.RiskResponseStrategyServiceInterface
	grcDashboardSvc         kpiServices.GRCDashboardServiceInterface
}

func NewGRCReportService(
	repo repositories.GRCReportRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	organizationLevelSvc settingServices.OrganizationLevelServiceInterface,
	riskResponseStrategySvc riskServices.RiskResponseStrategyServiceInterface,
	grcDashboardSvc kpiServices.GRCDashboardServiceInterface,
) GRCReportServiceInterface {
	return &GRCReportService{
		repo:                    repo,
		activityLogSvc:          activityLogSvc,
		organizationLevelSvc:    organizationLevelSvc,
		riskResponseStrategySvc: riskResponseStrategySvc,
		grcDashboardSvc:         grcDashboardSvc,
	}
}

func (grs *GRCReportService) GetAll(params *dtos.GRCReportListQueryParamsDTO) (*dtos.GRCReportListResponseDTO, error) {
	return grs.repo.GetAll(params)
}

func (grs *GRCReportService) GetGRCReportDetails(grcReportID uint) (*dtos.GRCReportDetailDTO, error) {
	grcReport, err := grs.repo.GetGRCReportByID(grcReportID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, grcReportID)
		}
		return nil, err
	}

	sections, err := grs.repo.GetGRCReportSections(grcReportID)
	if err != nil {
		return nil, err
	}

	governanceSummary, err := grs.repo.GetGRCReportGovernanceSummary(grcReportID)
	if err != nil {
		return nil, err
	}

	documents, err := grs.repo.GetGRCReportGovernanceDocuments(grcReportID)
	if err != nil {
		return nil, err
	}

	tasks, err := grs.repo.GetGRCReportGovernanceTasks(grcReportID)
	if err != nil {
		return nil, err
	}

	// 🔽 New Risk-Related Fetches
	riskSummary, err := grs.repo.GetGRCReportRiskSummary(grcReportID)
	if err != nil {
		return nil, err
	}

	topFiveRisks, err := grs.repo.GetGRCReportRiskTopFive(grcReportID)
	if err != nil {
		return nil, err
	}

	kris, err := grs.repo.GetGRCReportRiskKRIs(grcReportID)
	if err != nil {
		return nil, err
	}

	riskControls, err := grs.repo.GetGRCReportRiskControls(grcReportID)
	if err != nil {
		return nil, err
	}

	// Construct Response DTO
	response := &dtos.GRCReportDetailDTO{
		ID:                     grcReport.ID,
		Title:                  grcReport.Title,
		StartDate:              grcReport.StartDate,
		EndDate:                grcReport.EndDate,
		CreatedAt:              grcReport.CreatedAt,
		CreatedUserFirstNameEn: grcReport.CreatedUserFirstNameEn,
		CreatedUserFirstNameAr: grcReport.CreatedUserFirstNameAr,
		CreatedUserLastNameEn:  grcReport.CreatedUserLastNameEn,
		CreatedUserLastNameAr:  grcReport.CreatedUserLastNameAr,
		Sections:               sections,
		GovernanceSummary:      *governanceSummary,
		Documents:              documents,
		Tasks:                  tasks,
		RiskSummary:            *riskSummary,
		TopFiveRisks:           topFiveRisks,
		RiskKRIs:               kris,
		RiskControls:           riskControls,
	}

	return response, nil
}

func (grs *GRCReportService) GetSectionByID(grcReportID uint, sectionID uint, tx ...utility.DatabaseTransaction) (*models.GRCReportSection, error) {
	section, err := grs.repo.GetSectionByID(sectionID, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, sectionID)
		}
		return nil, err
	}
	return section, nil
}

func (grs *GRCReportService) Create(currentUserID uint, grcReportFormDTO *dtos.GRCReportFormDTO) (uint, error) {
	tx := grs.repo.BeginTransaction()

	dates, err := grs.parseGRCReportDates(grcReportFormDTO)
	if err != nil {
		return 0, err
	}

	grcReport := &models.GRCReport{
		Title:         grcReportFormDTO.Title,
		StartDate:     *dates.StartDate,
		EndDate:       *dates.EndDate,
		ReviewUserID:  &grcReportFormDTO.ReviewUserID,
		ApproveUserID: &grcReportFormDTO.ApproveUserID,
		CreatedBYID:   currentUserID,
	}

	grcReportID, err := grs.repo.Create(grcReport, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	sectionTypes := []string{
		"main-executive-summary", "governance-summary", "executive-summary-governance", "policy-approved-overview", "governance-tasks", "conclusion-governance", "main-conclusion",
	}

	for i, sectionType := range sectionTypes {
		titleEn, titleAr := formatGRCSectionTitle(sectionType)
		section := models.GRCReportSection{
			GRCReportID:   grcReportID,
			Type:          sectionType,
			TitleEn:       titleEn,
			TitleAr:       titleAr,
			DescriptionEn: "",
			DescriptionAr: "",
			OrderNum:      i + 1,
			CreatedBYID:   currentUserID,
		}
		if err := grs.repo.CreateSection(&section, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	//governance report section
	documents, err := grs.repo.GetPolicyDocumentByDateRange(*dates.StartDate, *dates.EndDate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	for _, document := range documents {
		summary, err := grs.GetDocumentAcknowledgmentSummary(document)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to get acknowledgment summary for document ID %d: %w", document.ID, err)
		}

		grcDocument := &models.GRCReportGovernanceDocument{
			GRCReportID:             grcReportID,
			ReferenceNumber:         document.ReferenceNumber,
			DocumentTitleEn:         document.TitleEn,
			DocumentTitleAr:         document.TitleAr,
			ResponsibleUserNameEn:   fmt.Sprintf("%s %s", document.ResponsibleUser.FirstNameEn, document.ResponsibleUser.LastNameEn),
			ResponsibleUserNameAr:   fmt.Sprintf("%s %s", document.ResponsibleUser.FirstNameAr, document.ResponsibleUser.LastNameAr),
			ResponsibleUnitTitleEn:  document.ResponsibleUnit.TitleEn,
			ResponsibleUnitTitleAr:  document.ResponsibleUnit.TitleAr,
			TotalEmployees:          int(summary.TotalUsers),
			AcknowledgedCount:       int(summary.AcknowledgedUsers),
			PendingCount:            int(summary.NotAcknowledgedUsers),
			DocumentStatusTitleEn:   document.DocumentStatus.TitleEn,
			DocumentStatusTitleAr:   document.DocumentStatus.TitleAr,
			DocumentStatusColorCode: document.DocumentStatus.ColorCode,
			CreatedBYID:             currentUserID,
		}

		if err := grs.repo.CreateGRCReportGovernanceDocument(grcDocument, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tasks, err := grs.repo.GetTasksByDateRange(*dates.StartDate, *dates.EndDate, tx)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to fetch grc tasks: %w", err)
	}

	for _, task := range tasks {
		var severityLevelTitleEn, severityLevelTitleAr, severityLevelColorCode *string
		if task.SeverityLevel != nil {
			severityLevelTitleEn = &task.SeverityLevel.TitleEn
			severityLevelTitleAr = &task.SeverityLevel.TitleAr
			severityLevelColorCode = &task.SeverityLevel.ColorCode
		}
		taskRecord := &models.GRCReportGovernanceTask{
			GRCReportID:            grcReportID,
			TaskTitle:              task.Title,
			TaskDescription:        task.Description,
			ResponsibleUserNameEn:  fmt.Sprintf("%s %s", task.ResponsibleUser.FirstNameEn, task.ResponsibleUser.LastNameEn),
			ResponsibleUserNameAr:  fmt.Sprintf("%s %s", task.ResponsibleUser.FirstNameAr, task.ResponsibleUser.LastNameAr),
			TaskStatusTitleEn:      task.TaskStatus.TitleEn,
			TaskStatusTitleAr:      task.TaskStatus.TitleAr,
			TaskStatusColorCode:    task.TaskStatus.ColorCode,
			TargetDate:             task.TargetDate,
			SeverityLevelTitleEn:   severityLevelTitleEn,
			SeverityLevelTitleAr:   severityLevelTitleAr,
			SeverityLevelColorCode: severityLevelColorCode,
			CreatedBYID:            currentUserID,
		}

		if err := grs.repo.CreateGRCReportGovernanceTask(taskRecord, tx); err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to insert grc_report_task: %w", err)
		}
	}

	// Build Governance Summary - using category titles
	var totalPolicies int
	categoryTitleSet := make(map[string]struct{})

	for _, doc := range documents {
		if doc.DocumentType != nil && doc.DocumentType.Type == "policy" {
			if doc.DocumentCategory != nil {
				title := doc.DocumentCategory.TitleEn
				categoryTitleSet[title] = struct{}{}
			}
			totalPolicies++
		}

	}

	var categoryTitles []string
	for title := range categoryTitleSet {
		categoryTitles = append(categoryTitles, title)
	}
	policiesInSystem := strings.Join(categoryTitles, ", ")

	// Count total users in system for roles defined
	userCount, err := grs.repo.AllUserCount(tx)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Policy Approved Rate
	var publishedPolicies int
	for _, doc := range documents {
		if doc.DocumentType != nil && doc.DocumentType.Type == "policy" {
			if doc.DocumentStatus != nil && doc.DocumentStatus.Type == "published" {
				publishedPolicies++
			}
		}
	}

	var policyApprovedRate *float64
	if totalPolicies > 0 {
		rate := (float64(publishedPolicies) / float64(totalPolicies)) * 100
		rounded := math.Round(rate*100) / 100 // round to 2 decimal places
		policyApprovedRate = &rounded
	}

	// Fetch data
	params := kpiDTOs.GRCDashboardListQueryParamsDTO{
		DateFrom: dates.StartDate.Format("2006-01-02"),
		DateTo:   dates.EndDate.Format("2006-01-02"),
	}

	documentOverview, err := grs.grcDashboardSvc.GetGRCDocumentOverview(&params)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch document overview: %w", err)
	}

	documentTypeChart, err := grs.grcDashboardSvc.GetGRCDocumentTypeChart(&params)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch document type chart: %w", err)
	}

	topDocuments, err := grs.grcDashboardSvc.GetGRCTopPublishedDocuments(&params)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch top published documents: %w", err)
	}

	overviewJSON, err := json.Marshal(documentOverview)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal document overview: %w", err)
	}

	typeChartJSON, err := json.Marshal(documentTypeChart)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal document type chart: %w", err)
	}

	topDocsJSON, err := json.Marshal(topDocuments)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal top published documents: %w", err)
	}

	overviewStr := string(overviewJSON)
	typeChartStr := string(typeChartJSON)
	topDocsStr := string(topDocsJSON)

	grcSummary := &models.GRCReportGovernanceSummary{
		GRCReportID:           grcReportID,
		ReviewUserID:          &grcReportFormDTO.ReviewUserID,
		ApproveUserID:         &grcReportFormDTO.ApproveUserID,
		PoliciesInSystem:      policiesInSystem,
		TotalPolicies:         totalPolicies,
		RolesDefined:          int(userCount),
		PolicyApprovedRate:    policyApprovedRate,
		Highlights:            nil,
		Challenges:            nil,
		DocumentOverviewChart: &overviewStr,
		DocumentTypeChart:     &typeChartStr,
		TopPublishedDocuments: &topDocsStr,
		CreatedBYID:           currentUserID,
	}

	if err := grs.repo.CreateGRCReportGovernanceSummary(grcSummary, tx); err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to insert grc_report_governance_summary: %w", err)
	}

	// risk report section
	risks, err := grs.repo.GetRisksByDateRange(*dates.StartDate, *dates.EndDate, tx)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch risks: %w", err)
	}

	var totalRisks int
	criticalCategorySet := make(map[string]struct{})
	var riskScoreSummaryParts []string
	var riskKRIs []*models.RiskKri
	totalRisks = len(risks)

	highestRatingType, err := grs.repo.GetHighestActiveRiskRatingType(tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	for _, risk := range risks {
		if risk.CurrentRiskRating != nil && risk.CurrentRiskRating.Type == highestRatingType {
			if risk.RiskCategory != nil {
				criticalCategorySet[risk.RiskCategory.TitleEn] = struct{}{}
			}
		}
		if risk.CurrentRiskRating != nil {
			part := fmt.Sprintf("- %s risk with strategy: %s", risk.CurrentRiskRating.TitleEn, risk.RiskResponseStrategy.TitleEn)
			riskScoreSummaryParts = append(riskScoreSummaryParts, part)
		}
		if risk.KRIs != nil {
			riskKRIs = append(riskKRIs, risk.KRIs...)
		}
	}

	strategies, err := grs.riskResponseStrategySvc.GetAll("", "", "", "active")
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to fetch active risk response strategies: %w", err)
	}
	var riskStrategies []string
	for _, s := range strategies {
		riskStrategies = append(riskStrategies, s.TitleEn)
	}

	var criticalCategories []string
	for title := range criticalCategorySet {
		criticalCategories = append(criticalCategories, title)
	}

	var kriSummaries []string
	for _, rk := range riskKRIs {
		if rk.KriRegister != nil && rk.KriRegister.Title != "" {
			value := "- " + rk.KriRegister.Title
			if rk.KriRegister.CurrentValue != nil {
				value += fmt.Sprintf(" %d%%", *rk.KriRegister.CurrentValue)
			}
			kriSummaries = append(kriSummaries, value)
		}
	}

	strPtr := func(s string) *string { return &s }

	summary := &models.GRCReportRiskSummary{
		GRCReportID:         grcReportID,
		RiskIdentifiedCount: totalRisks,
		CriticalRisks:       strPtr(strings.Join(criticalCategories, ", ")),
		TreatmentStrategies: strPtr(strings.Join(riskStrategies, ", ")),
		RiskScoreSummary:    strPtr(strings.Join(riskScoreSummaryParts, "\n")),
		KRIs:                strPtr(strings.Join(kriSummaries, "\n")),
		CreatedBYID:         currentUserID,
	}

	if err := grs.repo.CreateGRCReportRiskSummary(summary, tx); err != nil {
		return 0, fmt.Errorf("failed to insert grc_report_risk_summary: %w", err)
	}

	// Top 5 Risks based on Current Risk Rating Value (descending order)
	topFiveRisks, err := grs.repo.GetTopFiveRisksByRatingValue(*dates.StartDate, *dates.EndDate, tx)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to fetch top five risks: %w", err)
	}

	for _, risk := range topFiveRisks {
		controlTitles, err := grs.repo.GetControlTitlesByRiskID(risk.ID, tx)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to fetch control titles for risk %s: %w", risk.ReferenceNumber, err)
		}
		existingControls := strings.Join(controlTitles, ", ")

		topFive := &models.GRCReportRiskTopFive{
			GRCReportID:                 grcReportID,
			RiskReferenceNumber:         risk.ReferenceNumber,
			RiskCategoryTitleEn:         *risk.RiskCategoryTitleEn,
			RiskCategoryTitleAr:         *risk.RiskCategoryTitleAr,
			Title:                       risk.Title,
			Description:                 risk.Description,
			ResponsibleUserNameEn:       *risk.ResponsibleUserNameEn,
			ResponsibleUserNameAr:       *risk.ResponsibleUserNameAr,
			ResponsibleUnitTitleEn:      *risk.ResponsibleUnitTitleEn,
			ResponsibleUnitTitleAr:      *risk.ResponsibleUnitTitleAr,
			IsOpen:                      strings.ToLower(*risk.RiskStatusType) == "open",
			InherentRiskRatingTitleEn:   risk.InherentRiskRatingTitleEn,
			InherentRiskRatingTitleAr:   risk.InherentRiskRatingTitleAr,
			InherentRiskRatingColorCode: risk.InherentRiskRatingColorCode,
			InherentRiskRatingScore:     risk.InherentRiskRatingScore,
			InherentLikelihoodTitleEn:   risk.InherentLikelihoodTitleEn,
			InherentLikelihoodTitleAr:   risk.InherentLikelihoodTitleAr,
			InherentImpactTitleEn:       risk.InherentImpactTitleEn,
			InherentImpactTitleAr:       risk.InherentImpactTitleAr,
			ResidualRiskRatingTitleEn:   risk.ResidualRiskRatingTitleEn,
			ResidualRiskRatingTitleAr:   risk.ResidualRiskRatingTitleAr,
			ResidualRiskRatingColorCode: risk.ResidualRiskRatingColorCode,
			ResidualRiskRatingScore:     risk.ResidualRiskRatingScore,
			ResidualLikelihoodTitleEn:   risk.ResidualLikelihoodTitleEn,
			ResidualLikelihoodTitleAr:   risk.ResidualLikelihoodTitleAr,
			ResidualImpactTitleEn:       risk.ResidualImpactTitleEn,
			ResidualImpactTitleAr:       risk.ResidualImpactTitleAr,
			ExistingControls:            &existingControls,
			RiskResponseStrategyTitleEn: risk.RiskResponseStrategyTitleEn,
			RiskResponseStrategyTitleAr: risk.RiskResponseStrategyTitleAr,
			ApprovalNote:                nil,
			CreatedBYID:                 currentUserID,
		}
		if err := grs.repo.CreateGRCReportRiskTopFive(topFive, tx); err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to insert top five risk: %w", err)
		}
	}

	kriRecords, err := grs.repo.GetKRIsByDateRange(*dates.StartDate, *dates.EndDate, tx)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to fetch KRIs: %w", err)
	}

	for _, kri := range kriRecords {
		reportKRI := &models.GRCReportRiskKRI{
			GRCReportID:              grcReportID,
			RiskRegisterTitle:        kri.RiskRegisterTitle,
			KRIRegisterTitle:         kri.KRIRegisterTitle,
			KRIRegisterDescription:   kri.KRIRegisterDescription,
			LowThreshold:             kri.LowThreshold,
			MediumThreshold:          kri.MediumThreshold,
			CurrentValue:             kri.CurrentValue,
			RiskStatusTitleEn:        kri.RiskStatusTitleEn,
			RiskStatusTitleAr:        kri.RiskStatusTitleAr,
			RiskStatusTitleColorCode: kri.RiskStatusTitleColorCode,
			CreatedBYID:              currentUserID,
		}

		if err := grs.repo.CreateGRCReportRiskKRI(reportKRI, tx); err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to insert KRI into report: %w", err)
		}
	}

	riskControls, err := grs.repo.GetRiskControlsByDateRange(*dates.StartDate, *dates.EndDate, tx)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to fetch risk controls: %w", err)
	}

	for _, control := range riskControls {
		controlRecord := &models.GRCReportRiskControl{
			GRCReportID:                          grcReportID,
			ControlCode:                          control.ControlCode,
			ControlCategoryTitleEn:               control.ControlCategoryTitleEn,
			ControlCategoryTitleAr:               control.ControlCategoryTitleAr,
			Title:                                control.Title,
			Description:                          control.Description,
			ResponsibleUserNameEn:                control.ResponsibleUserNameEn,
			ResponsibleUserNameAr:                control.ResponsibleUserNameAr,
			ResponsibleUnitTitleEn:               control.ResponsibleUnitTitleEn,
			ResponsibleUnitTitleAr:               control.ResponsibleUnitTitleAr,
			ControlImplementationStatusTitleEn:   control.ControlImplementationStatusTitleEn,
			ControlImplementationStatusTitleAr:   control.ControlImplementationStatusTitleAr,
			ControlImplementationStatusColorCode: control.ControlImplementationStatusColorCode,
			ControlEfficiencyStatusTitleEn:       control.ControlEfficiencyStatusTitleEn,
			ControlEfficiencyStatusTitleAr:       control.ControlEfficiencyStatusTitleAr,
			ControlEfficiencyStatusColorCode:     control.ControlEfficiencyStatusColorCode,
			TestFrequencyTitleEn:                 control.TestFrequencyTitleEn,
			TestFrequencyTitleAr:                 control.TestFrequencyTitleAr,
			NextTestDate:                         control.NextTestDate,
			TestResponsibleUserNameEn:            control.TestResponsibleUserNameEn,
			TestResponsibleUserNameAr:            control.TestResponsibleUserNameAr,
			CreatedBYID:                          currentUserID,
		}

		if err := grs.repo.CreateGRCReportRiskControl(controlRecord, tx); err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to insert risk control record: %w", err)
		}
	}

	tx.Commit()

	go func() {
		err := grs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.GRCReport{}.TableName(),
			grcReportID,
			grcReport,
			currentUserID,
			nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{"error": err.Error()})
		}
	}()

	return grcReportID, nil
}

func (grs *GRCReportService) UpdateSection(grcReportID uint, sectionID uint, currentUserID uint, grcReporSectionFormDTO *dtos.GRCReportSectionFormDTO) error {
	tx := grs.repo.BeginTransaction()

	existingReportSection, err := grs.GetSectionByID(grcReportID, sectionID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Capture old values before update
	oldValues := map[string]interface{}{
		"title_en":       existingReportSection.TitleEn,
		"title_ar":       existingReportSection.TitleAr,
		"description_en": existingReportSection.DescriptionEn,
		"description_ar": existingReportSection.DescriptionAr,
	}

	// Update section details
	existingReportSection.TitleEn = grcReporSectionFormDTO.TitleEn
	existingReportSection.TitleAr = grcReporSectionFormDTO.TitleAr
	existingReportSection.DescriptionEn = *grcReporSectionFormDTO.DescriptionEn
	existingReportSection.DescriptionAr = *grcReporSectionFormDTO.DescriptionAr
	existingReportSection.UpdatedBYID = &currentUserID

	err = grs.repo.UpdateSection(existingReportSection, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Capture new values after update
	newValues := map[string]interface{}{
		"title_en":       existingReportSection.TitleEn,
		"title_ar":       existingReportSection.TitleAr,
		"description_en": existingReportSection.DescriptionEn,
		"description_ar": existingReportSection.DescriptionAr,
	}

	// Identify changes
	changedData := utility.CompareFieldChanges(oldValues, newValues)

	// Log activity only if there are changes
	if len(changedData) > 0 {
		go func() {
			err := grs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.GRCReportSection{}.TableName(),
				sectionID,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to log activity", map[string]interface{}{"error": err.Error()})
			}
		}()
	}

	return nil
}

func (grs *GRCReportService) DeleteGRCReportSection(grcReportID uint, sectionID uint, currentUserID uint) error {
	tx := grs.repo.BeginTransaction()

	section, err := grs.GetSectionByID(grcReportID, sectionID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete the section
	if err := grs.repo.DeleteGRCReportSection(sectionID, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Reorder remaining sections
	if err := grs.repo.ReorderGRCReportSections(grcReportID, section.OrderNum, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	go func() {
		err := grs.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.GRCReportSection{}.TableName(),
			sectionID,
			section,
			currentUserID,
			nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{"error": err.Error()})
		}
	}()
	return nil
}

func (grs *GRCReportService) DeleteGRCReport(grcReportID uint, currentUserID uint) error {
	tx := grs.repo.BeginTransaction()

	existingReport, err := grs.repo.GetGRCReportByID(grcReportID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := grs.repo.DeleteGRCReportGovernanceTasks(grcReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := grs.repo.DeleteGRCReportGovernanceDocuments(grcReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := grs.repo.DeleteGRCReportSections(grcReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := grs.repo.DeleteGRCReport(grcReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	go func() {
		err := grs.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.GRCReport{}.TableName(),
			grcReportID,
			existingReport,
			currentUserID,
			nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{"error": err.Error()})
		}
	}()

	tx.Commit()
	return nil
}

func formatGRCSectionTitle(sectionType string) (string, string) {
	switch sectionType {
	case "main-executive-summary":
		return "Main Executive Summary", "الملخص التنفيذي الرئيسي"
	case "governance-summary":
		return "Governance Summary", "ملخص الحوكمة"
	case "executive-summary-governance":
		return "Executive Summary - Governance", "الملخص التنفيذي - الحوكمة"
	case "policy-approved-overview":
		return "Policy Approved Overview", "نظرة عامة على السياسات المعتمدة"
	case "governance-tasks":
		return "Governance Tasks", "مهام الحوكمة"
	case "conclusion-governance":
		return "Conclusion - Governance", "الخاتمة - الحوكمة"
	case "main-conclusion":
		return "Main Conclusion", "الخاتمة الرئيسية"
	default:
		return sectionType, sectionType
	}
}

func (grs *GRCReportService) parseGRCReportDates(dto *dtos.GRCReportFormDTO) (*dtos.GRCReportParsedDates, error) {
	parsed := &dtos.GRCReportParsedDates{}

	var err error

	parsed.StartDate, err = utility.ParseDate(dto.StartDate)
	if err != nil {
		return nil, utility.NewValidationError(fmt.Errorf("invalid start date"), "invalid_date", "start date")
	}

	parsed.EndDate, err = utility.ParseDate(dto.EndDate)
	if err != nil {
		return nil, utility.NewValidationError(fmt.Errorf("invalid end date"), "invalid_date", "end date")
	}

	return parsed, nil
}

func (grs *GRCReportService) GetDocumentAcknowledgmentSummary(document *models.Document) (*dtos.DocumentAcknowledgmentSummaryDTO, error) {

	if document.DocumentAccessType.Type == "public" {
		return grs.repo.GetPublicDocumentAcknowledgmentSummary(document.ID)
	}

	primaryLevel, err := grs.organizationLevelSvc.GetPrimaryLevel()
	if err != nil {
		return nil, err
	}

	var primaryOrgUnitField string
	switch primaryLevel.Level {
	case 1:
		primaryOrgUnitField = "organization_unit_1_id"
	case 2:
		primaryOrgUnitField = "organization_unit_2_id"
	case 3:
		primaryOrgUnitField = "organization_unit_3_id"
	case 4:
		primaryOrgUnitField = "organization_unit_4_id"
	case 5:
		primaryOrgUnitField = "organization_unit_5_id"
	default:
		return nil, fmt.Errorf("unsupported primary level")
	}

	return grs.repo.GetDocumentAcknowledgmentSummary(document.ID, primaryOrgUnitField)
}
