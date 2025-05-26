package services

import (
	"encoding/json"
	"errors"
	"fmt"
	complianceDTOs "jethur/grc/compliance/dtos"
	complianceServices "jethur/grc/compliance/services"
	kpiDTOs "jethur/grc/kpi/dtos"
	kpiServices "jethur/grc/kpi/services"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	userDTOs "jethur/grc/organization/dtos"
	"jethur/grc/reports/dtos"
	"jethur/grc/reports/repositories"
	riskDTOs "jethur/grc/risk/dtos"
	riskServices "jethur/grc/risk/services"
	settingDTOs "jethur/grc/settings/dtos"
	settingServices "jethur/grc/settings/services"
	"jethur/utility"
	"math"
	"strings"
	"time"

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

	UpdateHighlights(grcID uint, currentUserID uint, grcReportHighlightChallengeFormDTO *dtos.GRCReportHighlightChallengeFormDTO) error
	UpdateChallenges(grcID uint, currentUserID uint, grcReportHighlightChallengeFormDTO *dtos.GRCReportHighlightChallengeFormDTO) error
	UpdateRiskApprovalNote(grcID uint, riskID uint, currentUserID uint, grcReportHighlightChallengeFormDTO *dtos.GRCReportHighlightChallengeFormDTO) error
}

type GRCReportService struct {
	repo                    repositories.GRCReportRepositoryInterface
	activityLogSvc          logServices.ActivityLogServiceInterface
	organizationLevelSvc    settingServices.OrganizationLevelServiceInterface
	riskResponseStrategySvc riskServices.RiskResponseStrategyServiceInterface
	grcDashboardSvc         kpiServices.GRCDashboardServiceInterface
	riskHeatMapSvc          riskServices.RiskHeatmapServiceInterface
	complianceStatusSvc     complianceServices.ComplianceStatusServiceInterface
	organizationUnitSvc     settingServices.OrganizationUnitServiceInterface
}

func NewGRCReportService(
	repo repositories.GRCReportRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	organizationLevelSvc settingServices.OrganizationLevelServiceInterface,
	riskResponseStrategySvc riskServices.RiskResponseStrategyServiceInterface,
	grcDashboardSvc kpiServices.GRCDashboardServiceInterface,
	riskHeatMapSvc riskServices.RiskHeatmapServiceInterface,
	complianceStatusSvc complianceServices.ComplianceStatusServiceInterface,
	organizationUnitSvc settingServices.OrganizationUnitServiceInterface,
) GRCReportServiceInterface {
	return &GRCReportService{
		repo:                    repo,
		activityLogSvc:          activityLogSvc,
		organizationLevelSvc:    organizationLevelSvc,
		riskResponseStrategySvc: riskResponseStrategySvc,
		grcDashboardSvc:         grcDashboardSvc,
		riskHeatMapSvc:          riskHeatMapSvc,
		complianceStatusSvc:     complianceStatusSvc,
		organizationUnitSvc:     organizationUnitSvc,
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

	complianceSummary, err := grs.repo.GetGRCReportComplianceSummary(grcReportID)
	if err != nil {
		return nil, err
	}

	complianceSummaryDTO := mapComplianceSummaryToDTO(complianceSummary)

	slas, err := grs.repo.GetGRCReportSLA(grcReportID)
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
		ComplianceSummary:      complianceSummaryDTO,
		SLA:                    slas,
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
		"main-executive-summary",
		"governance-summary",
		"executive-summary-governance",
		"policy-approved-overview",
		"governance-tasks",
		"conclusion-governance",
		"executive-summary-compliance",
		"compliance-summary-self-assessment",
		"overall-compliance-status-by-department-self-assessment",
		"sla-summary",
		"conclusion-recommendations-compliance",
		"integrated-grc-performance",
		"risk-management-summary",
		"executive-summary-risk",
		"risk-governance-context",
		"top-5-risks",
		"key-risk-indicators-kris-performance",
		"risk-control-matrix",
		"risk-monitoring-reporting-communication",
		"recommendations-action-plans",
		"conclusion-recommendations-risk",
		"main-conclusion",
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
			AcknowledgeDueDate:      document.LatestDocumentVersion.AcknowledgeDueDate,
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

	heatmapParams := riskDTOs.RiskHeatmapQueryParamsDTO{
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
	var riskKRIs []*models.RiskKri
	var htmlRiskSummary strings.Builder
	var htmlRiskKriSummary strings.Builder
	totalRisks = len(risks)

	highestRatingType, err := grs.repo.GetHighestActiveRiskRatingType(tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	ratingCount := make(map[string]int)

	for _, risk := range risks {
		if risk.CurrentRiskRating != nil && risk.CurrentRiskRating.Type == highestRatingType {
			if risk.RiskCategory != nil {
				criticalCategorySet[risk.RiskCategory.TitleEn] = struct{}{}
			}
		}

		if risk.CurrentRiskRating != nil {
			rating := risk.CurrentRiskRating.TitleEn
			ratingCount[rating]++
		}
		if risk.KRIs != nil {
			riskKRIs = append(riskKRIs, risk.KRIs...)
		}
	}
	htmlRiskSummary.WriteString("<ul>\n")

	for rating, count := range ratingCount {
		line := fmt.Sprintf("<li>%d %s-rated risks</li>\n", count, rating)
		htmlRiskSummary.WriteString(line)
	}

	htmlRiskSummary.WriteString("</ul>")

	riskSummary := htmlRiskSummary.String()

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

	htmlRiskKriSummary.WriteString("<ul>\n")
	for _, rk := range riskKRIs {
		if rk.KriRegister != nil && rk.KriRegister.Title != "" {
			if rk.KriRegister.CurrentValue != nil {
				kriLine := fmt.Sprintf("<li>%s abc %d%%</li>\n", rk.KriRegister.Title, *rk.KriRegister.CurrentValue)

				htmlRiskKriSummary.WriteString(kriLine)
			}
		}
	}

	htmlRiskKriSummary.WriteString("</ul>")
	riskKriSummary := htmlRiskKriSummary.String()

	riskOverviewChart, err := grs.grcDashboardSvc.GetGRCRiskOverview(&params)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch document overview: %w", err)
	}

	riskHeatmapChart, err := grs.riskHeatMapSvc.GetRiskHeatmapByRiskMatrix(&heatmapParams)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch document type chart: %w", err)
	}

	topRisks, err := grs.grcDashboardSvc.GetGRCTopRisks(&params)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch top published documents: %w", err)
	}

	risktrend, err := grs.repo.GetRiskTrend(*dates.StartDate, *dates.EndDate, tx)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch risk trend: %w", err)
	}

	riskOverviewJSON, err := json.Marshal(riskOverviewChart)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal risk overview: %w", err)
	}

	riskHeatmapJSON, err := json.Marshal(riskHeatmapChart)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal risk overview: %w", err)
	}

	topRiskJSON, err := json.Marshal(topRisks)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal top risk: %w", err)
	}

	riskTrendJSON, err := json.Marshal(risktrend)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal top risk: %w", err)
	}

	riskOverviewStr := string(riskOverviewJSON)
	riskHeatmapStr := string(riskHeatmapJSON)
	topRiskStr := string(topRiskJSON)
	riskTrendStr := string(riskTrendJSON)

	strPtr := func(s string) *string { return &s }

	summary := &models.GRCReportRiskSummary{
		GRCReportID:         grcReportID,
		RiskIdentifiedCount: totalRisks,
		CriticalRisks:       strPtr(strings.Join(criticalCategories, ", ")),
		TreatmentStrategies: strPtr(strings.Join(riskStrategies, ", ")),
		RiskScoreSummary:    &riskSummary,
		KRIs:                &riskKriSummary,
		RiskOverviewChart:   &riskOverviewStr,
		TopRisks:            &topRiskStr,
		RiskHeatmapChart:    &riskHeatmapStr,
		RiskTrend:           &riskTrendStr,
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

	//compliance and sla section
	assessments, err := grs.repo.GetAllComplianceAssessments(*dates.StartDate, *dates.EndDate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Group assessments by organization unit
	orgUnitAssessments := make(map[uint][]models.ComplianceAssessment)
	for _, assessment := range assessments {
		orgUnitAssessments[assessment.PrimaryUnitID] = append(orgUnitAssessments[assessment.PrimaryUnitID], *assessment)
	}

	// Calculate scores for each department
	departmentScores := make(map[uint]float64)
	for orgUnitID, assessments := range orgUnitAssessments {
		var totalScore float64
		var assessmentCount int

		// Calculate average score for each assessment
		for _, assessment := range assessments {
			var requirementScore float64
			var requirementCount int

			// Calculate average score of requirements for this assessment
			for _, requirement := range assessment.AssessmentRequirements {
				if requirement.Score != nil && requirement.ComplianceStatus != nil &&
					requirement.ComplianceStatus.Type != "not_applicable" {
					requirementScore += float64(*requirement.Score)
					requirementCount++
				}
			}

			// Add assessment average to total if it has valid requirements
			if requirementCount > 0 {
				totalScore += requirementScore / float64(requirementCount)
				assessmentCount++
			}
		}

		// Calculate department's average score
		if assessmentCount > 0 {
			departmentScores[orgUnitID] = totalScore / float64(assessmentCount)
		}
	}

	// Group departments by parent unit and calculate parent unit scores
	parentUnitScores := make(map[uint]float64)
	parentUnitDepartments := make(map[uint][]uint)

	for orgUnitID := range departmentScores {
		orgUnit, err := grs.organizationUnitSvc.GetByID(orgUnitID)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		// Add department to its parent's list
		if orgUnit.OrganizationUnitID != nil {
			parentUnitDepartments[*orgUnit.OrganizationUnitID] = append(parentUnitDepartments[*orgUnit.OrganizationUnitID], orgUnitID)
		}
	}

	// Calculate parent unit scores
	for parentID, departmentIDs := range parentUnitDepartments {
		var totalScore float64
		var deptCount int

		for _, deptID := range departmentIDs {
			if score, exists := departmentScores[deptID]; exists {
				totalScore += score
				deptCount++
			}
		}

		if deptCount > 0 {
			parentUnitScores[parentID] = totalScore / float64(deptCount)
		}
	}

	// Create summary entries for parent units first
	parentUnitSummaries := make(map[uint]uint) // parentUnitID -> summaryID
	for parentID, score := range parentUnitScores {
		summary := &models.GRCReportComplianceSummary{
			GRCReportID:            grcReportID,
			OrgUnitID:              parentID,
			OverallComplianceScore: &score,
			ReviewUserID:           &grcReportFormDTO.ReviewUserID,
			ApproveUserID:          &grcReportFormDTO.ApproveUserID,
			CreatedBYID:            currentUserID,
		}

		summaryID, err := grs.repo.CreateGRCReportComplianceSummary(summary, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		parentUnitSummaries[parentID] = summaryID
	}

	// Create department entries linked to parent summaries
	for orgUnitID, score := range departmentScores {
		orgUnit, err := grs.repo.GetOrganizationUnitByID(orgUnitID)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		if orgUnit.OrganizationUnitID != nil {
			parentSummaryID := parentUnitSummaries[*orgUnit.OrganizationUnitID]

			complianceStatusID, err := grs.getComplianceStatusByScore(score)
			if err != nil {
				tx.Rollback()
				return 0, err
			}

			department := &models.GRCReportComplianceDepartment{
				GRCReportID:                  grcReportID,
				GRCReportComplianceSummaryID: parentSummaryID,
				OrgUnitID:                    orgUnitID,
				ComplianceScore:              score,
				ComplianceStatusID:           *complianceStatusID,
				CreatedBYID:                  currentUserID,
			}

			_, err = grs.repo.CreateGRCReportComplianceDepartment(department, tx)
			if err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}

	// Create SLA entries
	slas, err := grs.repo.GetAllSLAs(tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	for _, sla := range slas {
		var issueDate, expiryDate *time.Time
		if len(sla.Documents) > 0 {
			issueDate = sla.Documents[0].IssueDate
			expiryDate = sla.Documents[0].ExpiryDate
		}

		var parentOrgUnitID *uint
		if sla.ResponsibleUnit != nil && sla.ResponsibleUnit.OrganizationUnitID != nil {
			parentOrgUnitID = sla.ResponsibleUnit.OrganizationUnitID
		} else {
			parentOrgUnitID = &sla.ResponsibleUnitID
		}

		grcReportSLA := &models.GRCReportSLA{
			GRCReportID:              grcReportID,
			ParentOrganizationUnitID: *parentOrgUnitID,
			OrganizationUnitID:       sla.ResponsibleUnitID,
			ReferenceNumber:          sla.ReferenceNumber,
			TitleEn:                  sla.TitleEn,
			TitleAr:                  &sla.TitleAr,
			ComplianceStatusID:       sla.ComplianceStatusID,
			LifecycleStatusID:        sla.LifecycleStatusID,
			IssueDate:                issueDate,
			ExpiryDate:               expiryDate,
			CreatedBYID:              currentUserID,
		}

		_, err = grs.repo.CreateGRCReportSLA(grcReportSLA, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
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
		return "Policy Approved And Notified Overview", "نظرة عامة على السياسات المعتمدة والمعلنة"
	case "governance-tasks":
		return "Governance Tasks & Action Plans", "مهام الحوكمة وخطط العمل"
	case "conclusion-governance":
		return "Conclusion - Governance", "الخاتمة - الحوكمة"
	case "executive-summary-compliance":
		return "Executive Summary Compliance", "الملخص التنفيذي - الامتثال"
	case "compliance-summary-self-assessment":
		return "Compliance Summary – Self Assessment", "ملخص الامتثال - التقييم الذاتي"
	case "overall-compliance-status-by-department-self-assessment":
		return "Overall Compliance Status by Department - Self Assessment", "حالة الامتثال العامة حسب الإدارات - التقييم الذاتي"
	case "conclusion-recommendations-compliance":
		return "Conclusion & Recommendations", "الخاتمة والتوصيات - الامتثال"
	case "integrated-grc-performance":
		return "Integrated GRC Performance", "الأداء المتكامل للحوكمة والامتثال"
	case "risk-management-summary":
		return "Risk Management Summary", "ملخص إدارة المخاطر"
	case "executive-summary-risk":
		return "Executive Summary", "الملخص التنفيذي - المخاطر"
	case "risk-governance-context":
		return "Risk Governance & Context", "حوكمة المخاطر والسياق"
	case "top-5-risks":
		return "Top 5 Risks", "أهم خمسة مخاطر"
	case "key-risk-indicators-kris-performance":
		return "Key Risk Indicators (KRIs) & Performance", "مؤشرات المخاطر الرئيسية والأداء"
	case "risk-control-matrix":
		return "Risk Control", "مصفوفة تحكم المخاطر"
	case "risk-monitoring-reporting-communication":
		return "Risk Monitoring, Reporting & Communication", "مراقبة المخاطر والتقارير والاتصالات"
	case "recommendations-action-plans":
		return "Recommendations & Action Plans", "التوصيات وخطط العمل"
	case "conclusion-recommendations-risk":
		return "Conclusion & Recommendations", "الخاتمة والتوصيات - المخاطر"
	case "main-conclusion":
		return "Main Conclusion", "الخاتمة الرئيسية"
	case "sla-summary":
		return "SLA Summary", "ملخص الخدمات المتاحة"
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

func (grs *GRCReportService) getComplianceStatusByScore(score float64) (*uint, error) {
	// Get all compliance statuses
	sort := "asc"
	sortBy := "id"
	keyword := ""
	status := ""
	statuses, err := grs.complianceStatusSvc.GetAll(sort, sortBy, keyword, status)
	if err != nil {
		return nil, err
	}

	// Find the status that matches the score range
	for _, status := range statuses {
		if status.ScoreFrom != nil && status.ScoreTo != nil {
			if score >= float64(*status.ScoreFrom) && score <= float64(*status.ScoreTo) {
				return &status.ID, nil
			}
		}
	}

	return nil, nil
}

func mapComplianceSummaryToDTO(complianceSummary []*models.GRCReportComplianceSummary) []*dtos.GRCReportComplianceSummaryDTO {
	complianceSummaryDTOs := []*dtos.GRCReportComplianceSummaryDTO{}
	for _, summary := range complianceSummary {
		// Map the summary's organization unit
		summaryOrgUnit := &settingDTOs.OrganizationUnitDTO{
			ID:      int(summary.OrganizationUnit.ID),
			TitleEn: summary.OrganizationUnit.TitleEn,
			TitleAr: summary.OrganizationUnit.TitleAr,
		}

		complianceSummaryDTO := &dtos.GRCReportComplianceSummaryDTO{
			ID:                     summary.ID,
			GRCReportID:            summary.GRCReportID,
			OrganizationUnit:       summaryOrgUnit,
			OverallComplianceScore: *summary.OverallComplianceScore,
			ReviewUser: &userDTOs.UserMinimalDTO{
				ID:          int(summary.ReviewUser.ID),
				FirstNameEN: summary.ReviewUser.FirstNameEn,
				FirstNameAR: summary.ReviewUser.FirstNameAr,
				LastNameEN:  summary.ReviewUser.LastNameEn,
				LastNameAR:  summary.ReviewUser.LastNameAr,
			},
			ApproveUser: &userDTOs.UserMinimalDTO{
				ID:          int(summary.ApproveUser.ID),
				FirstNameEN: summary.ApproveUser.FirstNameEn,
				FirstNameAR: summary.ApproveUser.FirstNameAr,
				LastNameEN:  summary.ApproveUser.LastNameEn,
				LastNameAR:  summary.ApproveUser.LastNameAr,
			},
			Departments: []*dtos.GRCReportComplianceDepartmentDTO{},
		}

		// Map each department with its own organization unit
		for _, department := range summary.GRCReportComplianceDepartments {
			// Create a new organization unit DTO for each department
			departmentOrgUnit := &settingDTOs.OrganizationUnitDTO{
				ID:      int(department.OrganizationUnit.ID),
				TitleEn: department.OrganizationUnit.TitleEn,
				TitleAr: department.OrganizationUnit.TitleAr,
			}

			complianceSummaryDTO.Departments = append(complianceSummaryDTO.Departments, &dtos.GRCReportComplianceDepartmentDTO{
				ID:               department.ID,
				OrganizationUnit: departmentOrgUnit,
				ComplianceScore:  department.ComplianceScore,
				ComplianceStatus: &complianceDTOs.ComplianceStatusDTO{
					ID:        int(department.ComplianceStatus.ID),
					TitleEn:   department.ComplianceStatus.TitleEn,
					TitleAr:   department.ComplianceStatus.TitleAr,
					ColorCode: department.ComplianceStatus.ColorCode,
				},
			})
		}
		complianceSummaryDTOs = append(complianceSummaryDTOs, complianceSummaryDTO)
	}
	return complianceSummaryDTOs
}

func (grs *GRCReportService) UpdateHighlights(grcReportID uint, currentUserID uint, grcReportHighlightChallengeFormDTO *dtos.GRCReportHighlightChallengeFormDTO) error {
	tx := grs.repo.BeginTransaction()

	existingReportGovernanceSummary, err := grs.repo.GetGRCReportGovernanceSummaryByID(grcReportID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update highlight details
	existingReportGovernanceSummary.Highlights = grcReportHighlightChallengeFormDTO.Description
	existingReportGovernanceSummary.UpdatedBYID = &currentUserID

	err = grs.repo.UpdateGRCReportGovernanceSummary(existingReportGovernanceSummary, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (grs *GRCReportService) UpdateChallenges(grcReportID uint, currentUserID uint, grcReportHighlightChallengeFormDTO *dtos.GRCReportHighlightChallengeFormDTO) error {
	tx := grs.repo.BeginTransaction()

	existingReportGovernanceSummary, err := grs.repo.GetGRCReportGovernanceSummaryByID(grcReportID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update challenge details
	existingReportGovernanceSummary.Challenges = grcReportHighlightChallengeFormDTO.Description
	existingReportGovernanceSummary.UpdatedBYID = &currentUserID

	err = grs.repo.UpdateGRCReportGovernanceSummary(existingReportGovernanceSummary, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (grs *GRCReportService) UpdateRiskApprovalNote(grcReportID uint, reportID uint, currentUserID uint, grcReportHighlightChallengeFormDTO *dtos.GRCReportHighlightChallengeFormDTO) error {
	tx := grs.repo.BeginTransaction()

	existingReportRisk, err := grs.repo.GetGRCReportRiskByID(grcReportID, reportID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update challenge details
	existingReportRisk.ApprovalNote = grcReportHighlightChallengeFormDTO.Description

	err = grs.repo.UpdateGRCReportRisk(existingReportRisk, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}
