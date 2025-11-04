package services

import (
	"fmt"
	documentDTOs "jethur/grc/documents/dtos"
	documentServices "jethur/grc/documents/services"
	settinsSvc "jethur/grc/settings/services"
	"os"
	"time"

	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	policyDTOs "jethur/grc/policy/dtos"
	"jethur/grc/policy/repositories"

	"jethur/utility"
)

type PolicyReportServiceInterface interface {
	GetAll(params *policyDTOs.PolicyReportListQueryParamsDTO) (*policyDTOs.PolicyReportLisResponseDTO, error)
	Create(currentUserID uint, PolicyReportFormDTO *policyDTOs.PolicyReportFormDTO) (uint, error)
	Delete(policyReportID uint, currentUserID uint) error

	GetPolicySummaryReport(policyReportID uint) (*policyDTOs.PolicySummaryReportDTO, error)
	GetPolicySummaryReportExport(policyReportID uint, currentUserID uint) (string, error)
}

type PolicyReportService struct {
	repo              repositories.PolicyReportRepositoryInterface
	activityLogSvc    logServices.ActivityLogServiceInterface
	documentSvc       documentServices.DocumentServiceInterface
	documentClauseSvc documentServices.DocumentClauseServiceInterface
	orgLevelSvc       settinsSvc.OrganizationLevelServiceInterface
	reportTemplateSvc settinsSvc.ReportTemplateSettingServiceInterface
}

func NewPolicyReportService(
	repo repositories.PolicyReportRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	documentSvc documentServices.DocumentServiceInterface,
	documentClauseSvc documentServices.DocumentClauseServiceInterface,
	orgLevelSvc settinsSvc.OrganizationLevelServiceInterface,
	reportTemplateSvc settinsSvc.ReportTemplateSettingServiceInterface,

) PolicyReportServiceInterface {
	return &PolicyReportService{
		repo:              repo,
		activityLogSvc:    activityLogSvc,
		documentSvc:       documentSvc,
		documentClauseSvc: documentClauseSvc,
		orgLevelSvc:       orgLevelSvc,
		reportTemplateSvc: reportTemplateSvc,
	}
}

func (prs *PolicyReportService) GetAll(params *policyDTOs.PolicyReportListQueryParamsDTO) (*policyDTOs.PolicyReportLisResponseDTO, error) {
	return prs.repo.GetAll(params)
}

func (prs *PolicyReportService) Create(currentUserID uint, policyReportFormDTO *policyDTOs.PolicyReportFormDTO) (uint, error) {
	tx := prs.repo.BeginTransaction()

	if err := prs.validatePolicyReportFormDTO(policyReportFormDTO); err != nil {
		return 0, err
	}

	dates, err := prs.parsePolicyReportDates(policyReportFormDTO)
	if err != nil {
		return 0, err
	}

	policyReport := &models.PolicyReport{
		Title:       policyReportFormDTO.Title,
		StartDate:   *dates.StartDate,
		EndDate:     *dates.EndDate,
		DocumentID:  policyReportFormDTO.DocumentID,
		CreatedByID: currentUserID,
	}

	policyReportID, err := prs.repo.Create(policyReport, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	go func() {
		err := prs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.PolicyReport{}.TableName(),
			policyReportID,
			policyReport,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{"error": err.Error()})
		}
	}()

	return policyReportID, nil
}

func (prs *PolicyReportService) Delete(policyReportID uint, currentUserID uint) error {

	existingReport, err := prs.repo.GetByID(policyReportID)
	if err != nil {
		return err
	}

	if err := prs.repo.Delete(policyReportID); err != nil {
		return err
	}

	go func() {
		err := prs.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.OrganizationUnitReport{}.TableName(),
			policyReportID,
			existingReport,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{"error": err.Error()})
		}
	}()

	return nil
}

func (prs *PolicyReportService) GetPolicySummaryReport(policyReportID uint) (*policyDTOs.PolicySummaryReportDTO, error) {

	tx := prs.repo.BeginTransaction()

	var policySummaryReport policyDTOs.PolicySummaryReportDTO

	reportDetails, err := prs.repo.GetFullDetailsByID(policyReportID, tx)
	if err != nil {
		return nil, err
	}

	if _, err := prs.reportTemplateSvc.GetImagePath("logo"); err != nil {
		return nil, utility.NewForbiddenError("logo_image_missing")
	}
	if _, err := prs.reportTemplateSvc.GetImagePath("cover-page-bg"); err != nil {
		return nil, utility.NewForbiddenError("cover_image_missing")
	}

	policySummaryReport.Title = reportDetails.Title
	policySummaryReport.DocumentTitleEn = reportDetails.DocumentTitleEn
	policySummaryReport.DocumentTitleAr = reportDetails.DocumentTitleAr
	policySummaryReport.GeneratedOn = reportDetails.CreatedAt
	policySummaryReport.GeneratedByUserFirstNameEn = reportDetails.CreatedUserFirstNameEn
	policySummaryReport.GeneratedByUserFirstNameAr = reportDetails.CreatedUserFirstNameAr
	policySummaryReport.GeneratedByUserLastNameEn = reportDetails.CreatedUserLastNameAr
	policySummaryReport.GeneratedByUserLastNameAr = reportDetails.CreatedUserLastNameAr
	policySummaryReport.GeneratedByUserEmail = reportDetails.CreatedUserEmail
	policySummaryReport.GeneratedByUserColor = reportDetails.CreatedUserColor
	policySummaryReport.GeneratedByUserInitial = reportDetails.CreatedUserInitial
	policySummaryReport.GeneratedByUserProfileImagePath = reportDetails.CreatedUserProfileImagePath

	documentID := reportDetails.DocumentID

	fmt.Print(documentID)

	// sample params for clauselist
	clauseParams := &documentDTOs.DocumentClauseListQueryParamsDTO{
		SortBy: "created_at",
		Sort:   "asc",
	}

	if policySummaryReport.PolicyOverview, err = prs.repo.GetPolicyOverview(documentID); err != nil {
		return nil, err
	}

	if policySummaryReport.PolicyKeySummary, err = prs.GetPolicyKeySummary(documentID); err != nil {
		return nil, err
	}

	if policySummaryReport.ClauseList, err = prs.documentClauseSvc.GetAll(clauseParams, documentID); err != nil {
		return nil, err
	}

	if policySummaryReport.AcknowledgementUserList, err = prs.GetPolicyAcknowldegementList(documentID); err != nil {
		return nil, err
	}

	if policySummaryReport.OrganizatationUnitAcknowledgementList, err = prs.GetPolicyAcknowldegementOrganizationUnitSummary(documentID); err != nil {
		return nil, err
	}

	if policySummaryReport.AssociatedRiskList, err = prs.repo.GetPolicyAssociatedRisks(documentID); err != nil {
		return nil, err
	}

	if policySummaryReport.AssociatedStrategicObjectiveLIst, err = prs.repo.GetPolicyAssociatedtrategicObjectives(documentID); err != nil {
		return nil, err
	}

	if policySummaryReport.AssociatedComplianceList, err = prs.repo.GetPolicyAssocatedCompliance(documentID); err != nil {
		return nil, err
	}

	if policySummaryReport.PolicyVersionList, err = prs.repo.GetPolicyVersions(documentID); err != nil {
		return nil, err
	}

	document, err := prs.documentSvc.GetFullDetailsByID(documentID, nil)
	if err != nil {
		return nil, err
	}

	if policySummaryReport.PolicyApprovalHistroyList, err = prs.repo.GetPolicyApprovalHistory(document.LatestDocumentVersion.ID); err != nil {
		return nil, err
	}

	primaryOrgUnitField, err := prs.GetPrimaryOrgField()
	if err != nil {
		return nil, err

	}

	if policySummaryReport.PolicyAccessUserList, err = prs.repo.GePolicyAccessUsers(documentID, primaryOrgUnitField); err != nil {
		return nil, err
	}

	if policySummaryReport.PolicyAttachmentList, err = prs.repo.GetPolicyAttachments(documentID); err != nil {
		return nil, err
	}

	policySummaryReport.AppUrl = os.Getenv("APP_URL")

	return &policySummaryReport, nil
}

func (prs *PolicyReportService) GetPolicySummaryReportExport(policyReportID uint, currentUserID uint) (string, error) {
	policyReport, err := prs.GetPolicySummaryReport(policyReportID)
	if err != nil {
		return "", err
	}

	policyReport.GeneratedOn = time.Now().Format("2006-01-02 03:04 PM")

	// Main content template
	templatePaths := []string{"policy/templates/pdf/policy_summary_report.html"}
	htmlContent, err := utility.GeneratePDFContent(policyReport, templatePaths, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate HTML content: %w", err)
	}

	// Header template - process with data to resolve template variables
	headerTemplatePaths := []string{"policy/templates/pdf/policy_summary_report_header.html"}
	headerContent, err := utility.GeneratePDFContent(policyReport, headerTemplatePaths, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate header content: %w", err)
	}

	// Footer template - process with data to resolve template variables
	footerTemplatePaths := []string{"policy/templates/pdf/policy_summary_report_footer.html"}
	footerContent, err := utility.GeneratePDFContent(policyReport, footerTemplatePaths, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate footer content: %w", err)
	}
	coverTemplatePaths := []string{"policy/templates/pdf/policy_summary_report_cover.html"}
	coverContent, err := utility.GeneratePDFContent(policyReport, coverTemplatePaths, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate footer content: %w", err)
	}

	// Convert HTML → PDF with header and footer
	pdfFileName := fmt.Sprintf("PolicySummaryReport_%d.pdf", policyReportID)
	if err := utility.GeneratePDFWithCoverHeaderFooter(coverContent, htmlContent, pdfFileName, headerContent, footerContent); err != nil {
		return "", fmt.Errorf("failed to generate PDF: %w", err)
	}

	return pdfFileName, nil
}

func (prs *PolicyReportService) GetPolicyKeySummary(documentID uint) (*policyDTOs.PolicyReportKeySummaryDTO, error) {
	document, err := prs.documentSvc.GetFullDetailsByID(documentID, nil)
	if err != nil {
		return nil, err
	}

	primaryOrgUnitField, err := prs.GetPrimaryOrgField()
	if err != nil {
		return nil, err

	}
	if document.DocumentAccessType.Type == "public" {
		return prs.repo.GetPolicyKeySummary(documentID, primaryOrgUnitField, true)

	}
	return prs.repo.GetPolicyKeySummary(documentID, primaryOrgUnitField, false)
}

func (prs *PolicyReportService) GetPolicyAcknowldegementList(documentID uint) ([]*policyDTOs.PolicyAcknowledgementListDTO, error) {
	document, err := prs.documentSvc.GetFullDetailsByID(documentID, nil)
	if err != nil {
		return nil, err
	}

	primaryOrgUnitField, err := prs.GetPrimaryOrgField()
	if err != nil {
		return nil, err

	}
	if document.DocumentAccessType.Type == "public" {
		return prs.repo.GetPolicyAcknowledgementList(documentID, document.LatestDocumentVersion.ID, primaryOrgUnitField, true)

	}
	return prs.repo.GetPolicyAcknowledgementList(documentID, document.LatestDocumentVersion.ID, primaryOrgUnitField, false)

}

func (prs *PolicyReportService) GetPolicyAcknowldegementOrganizationUnitSummary(documentID uint) (*policyDTOs.PolicyAcknowledgementOrganizatationUnitSummaryDTO, error) {
	document, err := prs.documentSvc.GetFullDetailsByID(documentID, nil)
	if err != nil {
		return nil, err
	}

	primaryOrgUnitField, err := prs.GetPrimaryOrgField()
	if err != nil {
		return nil, err

	}

	if document.DocumentAccessType.Type == "public" {
		return prs.repo.GetPolicyAcknowldegementOrganizationUnitSummary(documentID, document.LatestDocumentVersion.ID, primaryOrgUnitField, true)
	}
	return prs.repo.GetPolicyAcknowldegementOrganizationUnitSummary(documentID, document.LatestDocumentVersion.ID, primaryOrgUnitField, false)

}

func (prs *PolicyReportService) GetPrimaryOrgField() (string, error) {
	primaryLevel, err := prs.orgLevelSvc.GetPrimaryLevel()
	if err != nil {
		return "", err
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
		return "", fmt.Errorf("unsupported primary level")
	}

	return primaryOrgUnitField, nil
}

func (prs *PolicyReportService) parsePolicyReportDates(dto *policyDTOs.PolicyReportFormDTO) (*policyDTOs.PolicyReportParsedDates, error) {
	parsed := &policyDTOs.PolicyReportParsedDates{}

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

func (prs *PolicyReportService) validatePolicyReportFormDTO(dto *policyDTOs.PolicyReportFormDTO) error {

	if _, err := prs.documentSvc.GetByID(dto.DocumentID); err != nil {
		return utility.NewValidationError(fmt.Errorf("policy does not exist for ID: %d", dto.DocumentID), "not_exist", "document")
	}

	return nil
}
