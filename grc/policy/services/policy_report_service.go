package services

import (
	"encoding/json"
	"errors"
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	documentDTOs "jethur/grc/documents/dtos"
	"math"

	documentServices "jethur/grc/documents/services"
	organizationDTOs "jethur/grc/organization/dtos"
	settinsSvc "jethur/grc/settings/services"
	"os"
	"time"

	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	policyDTOs "jethur/grc/policy/dtos"
	"jethur/grc/policy/repositories"

	"jethur/utility"

	"gorm.io/gorm"
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

	documentID := policyReportFormDTO.DocumentID

	document, err := prs.documentSvc.GetFullDetailsByID(documentID, nil)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	primaryOrgUnitField, err := prs.GetPrimaryOrgField()
	if err != nil {
		return 0, err

	}

	clauseParams := &documentDTOs.DocumentClauseListQueryParamsDTO{
		SortBy: "created_at",
		Sort:   "asc",
	}

	isDocumentPublic := document.DocumentAccessType.Type == "public"

	overview, err := prs.repo.GetPolicyOverview(documentID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	overviewDTO := &models.PolicyReportOverview{
		PolicyReportID:              policyReportID,
		ReferenceNumber:             overview.ReferenceNumber,
		PolicyTitleEn:               overview.TitleEn,
		PolicyTitleAr:               overview.TitleAr,
		DescriptionEn:               overview.DescriptionEn,
		DescriptionAr:               overview.DescriptionAr,
		ResponsibleUnitTitleEn:      overview.ResponsibleUnitTitleEn,
		ResponsibleUnitTitleAr:      overview.ResponsibleUnitTitleAr,
		ResponsibleUserNameEn:       fmt.Sprintf("%s %s", overview.ResponsibleUserFirstNameEn, overview.ResponsibleUserLastNameEn),
		ResponsibleUserNameAr:       fmt.Sprintf("%s %s", overview.ResponsibleUserFirstNameAr, overview.ResponsibleUserLastNameAr),
		PublishedDate:               overview.PublishedAt,
		NextReviewDate:              overview.NextReviewDate,
		DocumentVersion:             overview.DocumentVersion,
		DocumentStatusTitleEn:       overview.DocumentStatusTitleEn,
		DocumentStatusTitleAr:       overview.DocumentStatusTitleAr,
		DocumentStatusColorCode:     overview.DocumentStatusColorCode,
		DocumentAccessTypeTitleEn:   overview.DocumentAccessTypeTitleEn,
		DocumentAccessTypeTitleAr:   overview.DocumentAccessTypeTitleAr,
		DocumentAccessTypeColorCode: overview.DocumentAccessTypeColorCode,
		CreatedByID:                 currentUserID,
	}

	if err := prs.repo.CreatePolicyReportOverview(overviewDTO, tx); err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to insert policy_report_overivews: %w", err)
	}

	keyMetrics, err := prs.repo.GetPolicyKeyMetrics(documentID, primaryOrgUnitField, isDocumentPublic)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	keyMetricsDTO := &models.PolicyReportKeyMetrics{
		PolicyReportID:                 policyReportID,
		ClausesDefinedCount:            int(keyMetrics.ClausesDefined),
		TotalAcknowledgementCount:      int(keyMetrics.TotalAcknowledgments),
		AcknowledgementsCompletedCount: int(keyMetrics.AcknowledgmentsCompleted),
		AcknowledgementsPendingCount:   int(keyMetrics.AcknowledgmentsPending),
		AcknowledgementPercentage:      float32(keyMetrics.AcknowledgmentRate),
		RisksReportedCount:             int(keyMetrics.RiskReported),
		AssociatedStandardsCount:       int(keyMetrics.AssociatedStandardsCount),
		AssociatedClausesCount:         int(keyMetrics.AssociatedClausesCount),
		StrategicObjectivesLinkedCount: int(keyMetrics.StrategicObjectivesLinked),
		ApprovalCycleCount:             int(keyMetrics.ApprovalCycle),
		CreatedByID:                    currentUserID,
	}

	if err := prs.repo.CreatePolicyReportKeyMetrics(keyMetricsDTO, tx); err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to insert policy_report_key_metrics: %w", err)
	}

	clauses, err := prs.documentClauseSvc.GetAll(clauseParams, documentID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	for _, clause := range clauses {

		clauseControlsJSON, err := json.Marshal(clause.Controls)
		if err != nil {
			return 0, fmt.Errorf("failed to marshal clause controls: %w", err)
		}

		clauseControlsStr := string(clauseControlsJSON)

		clauseDTO := &models.PolicyReportClause{
			PolicyReportID: policyReportID,
			ClauseNumber:   clause.ClauseNumber,
			ClauseTitleEn:  clause.ClauseTitleEn,
			ClauseTitleAr:  clause.ClauseTitleAr,
			Description:    clause.Description,
			Controls:       &clauseControlsStr,
			CreatedByID:    currentUserID,
		}

		if err := prs.repo.CreatePolicyReportClause(clauseDTO, tx); err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to insert policy_report_clauses: %w", err)
		}
	}

	acknowledgedUsers, err := prs.repo.GetPolicyAcknowledgementList(documentID, document.LatestDocumentVersion.ID, primaryOrgUnitField, isDocumentPublic)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	for _, acknowledgedUser := range acknowledgedUsers {
		acknowledgedUserDTO := &models.PolicyReportAcknowledgementUser{
			PolicyReportID:          policyReportID,
			UserNameEn:              fmt.Sprintf("%s %s", acknowledgedUser.FirstNameEn, acknowledgedUser.LastNameEn),
			UserNameAr:              fmt.Sprintf("%s %s", acknowledgedUser.FirstNameAr, acknowledgedUser.LastNameAr),
			OrganizationUnitTitleEn: acknowledgedUser.OrganizationUnitTitleEn,
			OrganizationUnitTitleAr: acknowledgedUser.OrganizationUnitTitleAr,
			Acknowledged:            acknowledgedUser.Acknowledged,
			AcknowledgedOn:          acknowledgedUser.AcknowledgedOn,
			Comment:                 acknowledgedUser.Comment,
			CreatedByID:             currentUserID,
		}

		if err := prs.repo.CreatePolicyReportAcknowledgmentUser(acknowledgedUserDTO, tx); err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to insert policy_report_clauses: %w", err)
		}
	}

	acknowledgedOrganizationUnits, err := prs.repo.GetPolicyAcknowldegementOrganizationUnit(documentID, document.LatestDocumentVersion.ID, primaryOrgUnitField, isDocumentPublic)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	for _, acknowledgedOrganizationUnit := range acknowledgedOrganizationUnits {
		acknowledgedOrganizationUnitDTO := &models.PolicyReportAcknowledgementOrganizationUnit{
			PolicyReportID:          policyReportID,
			OrganizationUnitTitleEn: acknowledgedOrganizationUnit.OrganizationUnitTitleEn,
			OrganizationUnitTitleAr: acknowledgedOrganizationUnit.OrganizationUnitTitleAr,
			TotalUsers:              int(acknowledgedOrganizationUnit.TotalUsers),
			AcknowledgedUsers:       int(acknowledgedOrganizationUnit.AcknowledgedUsers),
			NotAcknowledgedUsers:    int(acknowledgedOrganizationUnit.NotAcknowledgedUsers),
			AcknowledgedPercent:     float32(acknowledgedOrganizationUnit.AcknowledgedPercent),
			CreatedByID:             currentUserID,
		}

		if err := prs.repo.CreatePolicyReportAcknowledgmentOrganizationUnit(acknowledgedOrganizationUnitDTO, tx); err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to insert policy_report_acknowledgement_organization_units: %w", err)
		}
	}

	associatedRisks, err := prs.repo.GetPolicyAssociatedRisks(documentID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	for _, risk := range associatedRisks {
		riskDTO := &models.PolicyReportRisk{
			PolicyReportID:        policyReportID,
			RiskTitle:             risk.Title,
			Description:           risk.Description,
			RiskCategoryTitleEn:   risk.RiskCategoryTitleEn,
			RiskCategoryTitleAr:   risk.RiskCategoryTitleAr,
			RiskStatusTitleEn:     risk.RiskStatusTitleEn,
			RiskStatusTitleAr:     risk.RiskStatusTitleAr,
			RiskStatusColorCode:   risk.RiskStatusColorCode,
			RiskLikelihoodTitleEn: risk.CurrentLikelihoodTitleEn,
			RiskLikelihoodTitleAr: risk.CurrentLikelihoodTitleAr,
			RiskImpactTitleEn:     risk.CurrentImpactTitleEn,
			RiskImpactTitleAr:     risk.CurrentImpactTitleAr,
			RiskRatingTitleEn:     risk.CurrentRiskRatingTitleEn,
			RiskRatingTitleAr:     risk.CurrentRiskRatingTitleAr,
			RiskRatingColorCode:   risk.CurrentRiskRatingColorCode,
			CreatedByID:           currentUserID,
		}

		if err := prs.repo.CreatePolicyReportRisk(riskDTO, tx); err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to insert policy_report_risks: %w", err)
		}
	}

	associatedStrategicObjectives, err := prs.repo.GetPolicyAssociatedtrategicObjectives(documentID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	for _, strategicObjective := range associatedStrategicObjectives {

		strategicObjectiveDTO := &models.PolicyReportStrategicObjective{
			PolicyReportID:                    policyReportID,
			StrategicObjectiveTitleEn:         strategicObjective.TitleEn,
			StrategicObjectiveTitleAr:         strategicObjective.TitleAr,
			DescriptionEn:                     strategicObjective.DescriptionEn,
			DescriptionAr:                     strategicObjective.DescriptionAr,
			StartDate:                         strategicObjective.StartDate,
			EndDate:                           strategicObjective.EndDate,
			StrategicObjectiveCategoryTitleEn: strategicObjective.StrategicObjectiveCategoryTitleEn,
			StrategicObjectiveCategoryTitleAr: strategicObjective.StrategicObjectiveCategoryTitleAr,
			ResponsibleUserNameEn:             fmt.Sprintf("%s %s", *strategicObjective.ResponsibleUserFirstNameEn, *strategicObjective.ResponsibleUserLastNameEn),
			ResponsibleUserNameAr:             fmt.Sprintf("%s %s", *strategicObjective.ResponsibleUserFirstNameAr, *strategicObjective.ResponsibleUserLastNameAr),
			ResponsibleUnitTitleEn:            *strategicObjective.ResponsibleUnitTitleEn,
			ResponsibleUnitTitleAr:            *strategicObjective.ResponsibleUnitTitleAr,
			CreatedByID:                       currentUserID,
		}

		if err := prs.repo.CreatePolicyReportStrategicObjective(strategicObjectiveDTO, tx); err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to insert policy_report_strategic_objectives: %w", err)
		}
	}

	associatedCompliances, err := prs.repo.GetPolicyAssocatedCompliance(documentID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	for _, compliance := range associatedCompliances {
		complianceClausesJSON, err := json.Marshal(compliance.Clauses)
		if err != nil {
			return 0, fmt.Errorf("failed to marshal compliance clauses: %w", err)
		}

		complianceClausesStr := string(complianceClausesJSON)

		complianceDTO := &models.PolicyReportCompliance{
			PolicyReportID:    policyReportID,
			ComplianceTitleEn: compliance.ComplianceRegisterTitleEn,
			ComplianceTitleAr: compliance.ComplianceRegisterTitleAr,
			Clauses:           &complianceClausesStr,
			CreatedByID:       currentUserID,
		}

		if err := prs.repo.CreatePolicyReportCompliance(complianceDTO, tx); err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to insert policy_report_compliances: %w", err)
		}
	}

	documentVersions, err := prs.repo.GetPolicyVersions(documentID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	for _, version := range documentVersions {
		var firstNameEn, lastNameEn, firstNameAr, lastNameAr string

		if version.PublishedByFirstNameEn != nil {
			firstNameEn = *version.PublishedByFirstNameEn
		}
		if version.PublishedByLastNameEn != nil {
			lastNameEn = *version.PublishedByLastNameEn
		}
		if version.PublishedByFirstNameAr != nil {
			firstNameAr = *version.PublishedByFirstNameAr
		}
		if version.PublishedByLastNameAr != nil {
			lastNameAr = *version.PublishedByLastNameAr
		}

		userNameEn := fmt.Sprintf("%s %s", firstNameEn, lastNameEn)
		userNameAr := fmt.Sprintf("%s %s", firstNameAr, lastNameAr)

		versionDTO := &models.PolicyReportDocumentVersion{
			PolicyReportID:      policyReportID,
			DocumentVersion:     version.Version,
			PublishedAt:         version.PublishedAt,
			PublishedUserNameEn: &userNameEn,
			PublishedUserNameAr: &userNameAr,
			CreatedByID:         currentUserID,
		}

		if err := prs.repo.CreatePolicyReportDocumentVersion(versionDTO, tx); err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to insert policy_report_document_versions: %w", err)
		}
	}

	approvalHistories, err := prs.repo.GetPolicyApprovalHistory(document.LatestDocumentVersion.ID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	for _, history := range approvalHistories {
		historyDTO := &models.PolicyReportDocumentApprovalHistory{
			PolicyReportID:     policyReportID,
			ApprovedUserNameEn: fmt.Sprintf("%s %s", history.ApprovedByFirstNameEn, history.ApprovedByLastNameEn),
			ApprovedUserNameAr: fmt.Sprintf("%s %s", history.ApprovedByFirstNameAr, history.ApprovedByLastNameAr),
			DesignationTitleEn: history.DesignationTitleEn,
			DesignationTitleAr: history.DesignationTitleAr,
			ActionTitleEn:      history.ActionTitleEn,
			ActionTitleAr:      history.ActionTitleAr,
			ActionColorCode:    history.ActionColorCode,
			ActionDate:         history.ActionDate,
			Comment:            &history.Comment,
			CreatedByID:        currentUserID,
		}

		if err := prs.repo.CreatePolicyReportDocumentApprovalHistory(historyDTO, tx); err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to insert policy_report_document_approval_histories: %w", err)
		}
	}

	documentAccesses, err := prs.repo.GePolicyAccessUsers(documentID, primaryOrgUnitField)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	for _, access := range documentAccesses {
		accessDTO := &models.PolicyReportDocumentAccess{
			PolicyReportID:          policyReportID,
			UserNameEn:              fmt.Sprintf("%s %s", access.FirstNameEn, access.LastNameEn),
			UserNameAr:              fmt.Sprintf("%s %s", access.FirstNameAr, access.LastNameAr),
			DesignationTitleEn:      access.DesignationTitleEn,
			DesignationTitleAr:      access.DesignationTitleAr,
			OrganizationUnitTitleEn: access.OrganizationUnitTitleEn,
			OrganizationUnitTitleAr: access.OrganizationUnitTitleAr,
			CreatedByID:             currentUserID,
		}

		if err := prs.repo.CreatePolicyReportDocumentAccess(accessDTO, tx); err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to insert policy_report_document_accesses: %w", err)

		}
	}

	documentAttachments, err := prs.repo.GetPolicyAttachments(documentID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	for _, attchment := range documentAttachments {
		attchmentDTO := &models.PolicyReportDocumentAttachment{
			PolicyReportID:      policyReportID,
			DocumentTitleEn:     attchment.DocumentTitleEn,
			DocumentTitleAr:     attchment.DocumentTitleAr,
			DocumentTypeTitleEn: attchment.DocumentTypeTitleEn,
			DocumentTypeTitleAr: attchment.DocumentTypeTitleAr,
			LinkedDate:          attchment.CreatedAt,
			CreatedByID:         currentUserID,
		}

		if err := prs.repo.CreatePolicyReportDocumentAttachment(attchmentDTO, tx); err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to insert policy_report_document_attachments: %w", err)
		}
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
	tx := prs.repo.BeginTransaction()

	existingReport, err := prs.repo.GetByID(policyReportID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, policyReportID)
		}
		return err
	}

	if err := prs.repo.DeletePolicyReportOverview(policyReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := prs.repo.DeletePolicyReportKeyMetrics(policyReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := prs.repo.DeletePolicyReportClause(policyReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := prs.repo.DeletePolicyReportAcknowledgmentUser(policyReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := prs.repo.DeletePolicyReportAcknowledgmentOrganizationUnit(policyReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := prs.repo.DeletePolicyReportOverview(policyReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := prs.repo.DeletePolicyReportRisk(policyReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := prs.repo.DeletePolicyReportStrategicObjective(policyReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := prs.repo.DeletePolicyReportCompliance(policyReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := prs.repo.DeletePolicyReportDocumentVersion(policyReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := prs.repo.DeletePolicyReportDocumentAccess(policyReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := prs.repo.DeletePolicyReportDocumentApprovalHistory(policyReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := prs.repo.DeletePolicyReportDocumentAttachment(policyReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := prs.repo.Delete(policyReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

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

	reportDetails, err := prs.repo.GetFullDetailsByID(policyReportID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, policyReportID)
		}
		return nil, err
	}

	if _, err := prs.reportTemplateSvc.GetImagePath("logo"); err != nil {
		return nil, utility.NewForbiddenError("logo_image_missing")
	}
	if _, err := prs.reportTemplateSvc.GetImagePath("cover-page-bg"); err != nil {
		return nil, utility.NewForbiddenError("cover_image_missing")
	}

	overview, err := prs.repo.GetPolicyReportOverview(policyReportID)
	if err != nil {
		return nil, err
	}

	keyMetrics, err := prs.repo.GetPolicyReportKeyMetrics(policyReportID)
	if err != nil {
		return nil, err
	}

	clauses, err := prs.repo.GetPolicyReportClause(policyReportID)
	if err != nil {
		return nil, err
	}

	for i := range clauses {
		if clauses[i].Controls != nil {
			var controls []*organizationDTOs.ControlMinimalListDTO
			if err := json.Unmarshal([]byte(*clauses[i].Controls), &controls); err == nil {
				clauses[i].ControlsParsed = controls
			} else {
				clauses[i].ControlsParsed = []*organizationDTOs.ControlMinimalListDTO{}
			}
		} else {
			clauses[i].ControlsParsed = []*organizationDTOs.ControlMinimalListDTO{}
		}
	}

	acknowledgedUsers, err := prs.repo.GetPolicyReportAcknowledgmentUser(policyReportID)
	if err != nil {
		return nil, err
	}

	acknowledgedUnits, err := prs.repo.GetPolicyReportAcknowledgmentOrganizationUnit(policyReportID)
	if err != nil {
		return nil, err
	}

	var totalUsers, totalAcknowledged, totalNotAcknowledged int

	for _, unit := range acknowledgedUnits {
		totalUsers += unit.TotalUsers
		totalAcknowledged += unit.AcknowledgedUsers
		totalNotAcknowledged += unit.NotAcknowledgedUsers
	}

	var totalPercent float32
	if totalUsers > 0 {
		percent := float64(totalAcknowledged) * 100.0 / float64(totalUsers)
		totalPercent = float32(math.Round(percent*100) / 100)
	}

	acknowledgmentUnitSummary := &policyDTOs.PolicyReportAcknowledgementOrganizationUnitSummaryResponseDTO{
		Items:                     acknowledgedUnits,
		TotalUsers:                totalUsers,
		TotalAcknowledgedUsers:    totalAcknowledged,
		TotalNotAcknowledgedUsers: totalNotAcknowledged,
		TotalAcknowledgedPercent:  totalPercent,
	}

	risks, err := prs.repo.GetPolicyReportRisk(policyReportID)
	if err != nil {
		return nil, err
	}

	strategicObjectives, err := prs.repo.GetPolicyReportStrategicObjective(policyReportID)
	if err != nil {
		return nil, err
	}

	compliances, err := prs.repo.GetPolicyReportCompliance(policyReportID)
	if err != nil {
		return nil, err
	}

	for i := range compliances {
		if compliances[i].Clauses != nil {
			var clauses []*commonDTOs.AssociatedClauseDTO
			if err := json.Unmarshal([]byte(*compliances[i].Clauses), &clauses); err == nil {
				compliances[i].ClausesParsed = clauses
			} else {
				compliances[i].ClausesParsed = []*commonDTOs.AssociatedClauseDTO{}
			}
		} else {
			compliances[i].ClausesParsed = []*commonDTOs.AssociatedClauseDTO{}
		}
	}

	documentVersions, err := prs.repo.GetPolicyReportDocumentVersion(policyReportID)
	if err != nil {
		return nil, err
	}

	documentAccesses, err := prs.repo.GetPolicyReportDocumentAccess(policyReportID)
	if err != nil {
		return nil, err
	}

	approvalHistories, err := prs.repo.GetPolicyReportDocumentApprovalHistory(policyReportID)
	if err != nil {
		return nil, err
	}

	documentAttachments, err := prs.repo.GetPolicyReportDocumentAttachment(policyReportID)
	if err != nil {
		return nil, err
	}

	response := &policyDTOs.PolicySummaryReportDTO{
		ID:                                    reportDetails.ID,
		Title:                                 reportDetails.Title,
		DocumentTitleEn:                       reportDetails.DocumentTitleEn,
		DocumentTitleAr:                       reportDetails.DocumentTitleAr,
		CreatedAt:                             reportDetails.CreatedAt,
		CreatedUserFirstNameEn:                reportDetails.CreatedUserFirstNameEn,
		CreatedUserLastNameEn:                 reportDetails.CreatedUserLastNameEn,
		CreatedUserFirstNameAr:                reportDetails.CreatedUserFirstNameAr,
		CreatedUserLastNameAr:                 reportDetails.CreatedUserLastNameAr,
		PolicyOverview:                        overview,
		PolicyKeyMetrics:                      keyMetrics,
		ClauseList:                            clauses,
		AcknowledgementUserList:               acknowledgedUsers,
		OrganizatationUnitAcknowledgementList: acknowledgmentUnitSummary,
		AssociatedRiskList:                    risks,
		AssociatedStrategicObjectiveLIst:      strategicObjectives,
		AssociatedComplianceList:              compliances,
		PolicyVersionList:                     documentVersions,
		PolicyAccessUserList:                  documentAccesses,
		PolicyApprovalHistroyList:             approvalHistories,
		PolicyAttachmentList:                  documentAttachments,
		AppUrl:                                os.Getenv("APP_URL"),
	}

	return response, nil
}

func (prs *PolicyReportService) GetPolicySummaryReportExport(policyReportID uint, currentUserID uint) (string, error) {
	policyReport, err := prs.GetPolicySummaryReport(policyReportID)
	if err != nil {
		return "", err
	}

	policyReport.CreatedAt = time.Now().Format("2006-01-02 03:04 PM")

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
