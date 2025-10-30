package services

import (
	"fmt"
	documentDTOs "jethur/grc/documents/dtos"
	documentServices "jethur/grc/documents/services"
	settinsSvc "jethur/grc/settings/services"

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
}

type PolicyReportService struct {
	repo              repositories.PolicyReportRepositoryInterface
	activityLogSvc    logServices.ActivityLogServiceInterface
	documentSvc       documentServices.DocumentServiceInterface
	documentClauseSvc documentServices.DocumentClauseServiceInterface
	orgLevelSvc       settinsSvc.OrganizationLevelServiceInterface
}

func NewPolicyReportService(
	repo repositories.PolicyReportRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	documentSvc documentServices.DocumentServiceInterface,
	documentClauseSvc documentServices.DocumentClauseServiceInterface,
	orgLevelSvc settinsSvc.OrganizationLevelServiceInterface,
) PolicyReportServiceInterface {
	return &PolicyReportService{
		repo:              repo,
		activityLogSvc:    activityLogSvc,
		documentSvc:       documentSvc,
		documentClauseSvc: documentClauseSvc,
		orgLevelSvc:       orgLevelSvc,
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

	policySummaryReport.Title = reportDetails.Title
	policySummaryReport.DocumentTitleEn = reportDetails.Document.TitleEn
	policySummaryReport.DocumentTitleAr = reportDetails.Document.TitleAr
	policySummaryReport.GeneratedOn = reportDetails.CreatedAt

	if reportDetails.CreatedByUser != nil {
		policySummaryReport.GeneratedByUserFirstNameEn = reportDetails.CreatedByUser.FirstNameEn
		policySummaryReport.GeneratedByUserFirstNameAr = reportDetails.CreatedByUser.FirstNameAr
		policySummaryReport.GeneratedByUserLastNameEn = reportDetails.CreatedByUser.LastNameEn
		policySummaryReport.GeneratedByUserLastNameAr = reportDetails.CreatedByUser.LastNameAr
		policySummaryReport.GeneratedByUserEmail = reportDetails.CreatedByUser.Email
		policySummaryReport.GeneratedByUserColor = reportDetails.CreatedByUser.Color
		policySummaryReport.GeneratedByUserInitial = reportDetails.CreatedByUser.Initial
		if reportDetails.CreatedByUser.ProfileImagePath != nil {
			policySummaryReport.GeneratedByUserProfileImagePath = *reportDetails.CreatedByUser.ProfileImagePath
		}
	}

	if policySummaryReport.PolicyOverview, err = prs.repo.GetPolicyOverview(reportDetails.DocumentID); err != nil {
		return nil, err
	}

	if policySummaryReport.PolicyKeySummary, err = prs.repo.GetPolicyKeySummary(reportDetails.DocumentID); err != nil {
		return nil, err
	}

	// sample params for clauselist
	clauseParams := &documentDTOs.DocumentClauseListQueryParamsDTO{
		SortBy: "created_at",
		Sort:   "asc",
	}

	if policySummaryReport.ClauseList, err = prs.documentClauseSvc.GetAll(clauseParams, reportDetails.DocumentID); err != nil {
		return nil, err
	}

	if policySummaryReport.AcknowledgementUserList, err = prs.GetPolicyAcknowldegementList(reportDetails.DocumentID); err != nil {
		return nil, err
	}

	if policySummaryReport.OrganizatationUnitAcknowledgementList, err = prs.GetPolicyAcknowldegementOrganizationUnitSummary(reportDetails.DocumentID); err != nil {
		return nil, err
	}

	if policySummaryReport.AssociatedRiskList, err = prs.repo.GetPolicyAssociatedRisks(reportDetails.DocumentID); err != nil {
		return nil, err
	}

	if policySummaryReport.AssociatedStrategicObjectiveLIst, err = prs.repo.GetPolicyAssociatedtrategicObjectives(reportDetails.DocumentID); err != nil {
		return nil, err
	}

	if policySummaryReport.AssociatedComplianceRequirmentLIst, err = prs.repo.GetPolicyAssociatedComplianceRequirements(reportDetails.DocumentID); err != nil {
		return nil, err
	}

	if policySummaryReport.PolicyVersionList, err = prs.repo.GetPolicyVersions(reportDetails.DocumentID); err != nil {
		return nil, err
	}

	document, err := prs.documentSvc.GetFullDetailsByID(reportDetails.DocumentID, nil)
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

	if policySummaryReport.PolicyAccessUserList, err = prs.repo.GePolicyAccessUsers(reportDetails.DocumentID, primaryOrgUnitField); err != nil {
		return nil, err
	}

	if policySummaryReport.PolicyAttachmentList, err = prs.repo.GetPolicyAttachments(reportDetails.DocumentID); err != nil {
		return nil, err
	}

	return &policySummaryReport, nil
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
