package services

import (
	"encoding/json"
	"errors"
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/grc/organization/repositories"
	settingServices "jethur/grc/settings/services"
	"jethur/utility"
	"math"
	"os"
	"time"

	"gorm.io/gorm"
)

// Interface
type OrganizationUnitsServiceInterface interface {
	GetOrganizationUnitsHierarchies(sort string, sortBy string) ([]organizationDTOs.OrganizationUnitsDTO, error)
	GetUsersByOrganizationUnit(params *organizationDTOs.OrganizationUnitUserPrimaryDetailsQueryParamsDTO, orgUnitID uint) (*organizationDTOs.OrganizationUnitUserPrimaryDetailsResponseDTO, error)
	UpdatePrimaryOrganizationUnitHead(unitID uint, userID uint, currentUserID uint, departmentHeadUpdateDTO *organizationDTOs.DepartmentHeadUpdateDTO) error
	GetOrganizationUnitSummary(unitID uint) (*organizationDTOs.OrganizationUnitSummaryDTO, error)
	GetAllPoliciesByOrganizationUnit(unitID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error)
	GetAllRiskByOrganizationUnit(unitID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error)
	GetParentOrganizationUnitWithHead(unitID uint) (*organizationDTOs.ParentOrganizationUnitWithHeadDTO, error)

	GetAll(params *organizationDTOs.OrganizationUnitReportListQueryParamsDTO) (*organizationDTOs.OrganizationUnitReportListReponseDTO, error)
	CreateOrganizationUnitReport(currentUserID uint, organizationUnitReportFormDTO *organizationDTOs.OrganizationUnitReportFormDTO) (uint, error)
	DeleteOrganizationUnitReport(organizationUnitReportID uint, currentUserID uint) error
	GetOrganizationUnitReport(organizationUnitReportID uint) (*organizationDTOs.OrganizationUnitSummaryReportDTO, error)
	GetOrganizationUnitReportExport(organizationReportID uint, currentUserID uint) (string, error)
}

type OrganizationUnitsService struct {
	repo                 repositories.OrganizationUnitsRepositoryInterface
	activityLog          logServices.ActivityLogServiceInterface
	organizationLevelSvc settingServices.OrganizationLevelServiceInterface
	userRepo             repositories.UserRepositoryInterface
	userSvc              UserServiceInterface
	organizationUnitSvc  settingServices.OrganizationUnitServiceInterface
	reportTemplateSvc    settingServices.ReportTemplateSettingServiceInterface
}

func NewOrganizationUnitsService(
	repo repositories.OrganizationUnitsRepositoryInterface,
	activityLog logServices.ActivityLogServiceInterface,
	organizationLevelSvc settingServices.OrganizationLevelServiceInterface,
	userRepo repositories.UserRepositoryInterface,
	userSvc UserServiceInterface,
	organizationUnitSvc settingServices.OrganizationUnitServiceInterface,
	reportTemplateSvc settingServices.ReportTemplateSettingServiceInterface,
) *OrganizationUnitsService {
	return &OrganizationUnitsService{
		repo:                 repo,
		activityLog:          activityLog,
		organizationLevelSvc: organizationLevelSvc,
		userRepo:             userRepo,
		userSvc:              userSvc,
		organizationUnitSvc:  organizationUnitSvc,
		reportTemplateSvc:    reportTemplateSvc,
	}
}

func (s *OrganizationUnitsService) GetOrganizationUnitsHierarchies(sort string, sortBy string) ([]organizationDTOs.OrganizationUnitsDTO, error) {
	// Get all organization units
	units, err := s.repo.GetOrganizationUnitsHierarchies(sort, sortBy)
	if err != nil {
		return nil, err
	}

	// Create a map for quick lookup of units by ID
	unitMap := make(map[uint]*organizationDTOs.OrganizationUnitsDTO)
	var rootUnits []*organizationDTOs.OrganizationUnitsDTO

	// First pass: create all nodes
	for _, unit := range units {
		structure := &organizationDTOs.OrganizationUnitsDTO{
			ID:       unit.ID,
			TitleEn:  unit.TitleEn,
			TitleAr:  unit.TitleAr,
			Children: make([]organizationDTOs.OrganizationUnitsDTO, 0), // Initialize empty slice
		}

		// Safely handle nil OrganizationLevel
		if unit.OrganizationLevel != nil {
			structure.LevelTitleEn = unit.OrganizationLevel.TitleEn
			structure.LevelTitleAr = unit.OrganizationLevel.TitleAr
			structure.IsPrimary = unit.OrganizationLevel.IsPrimary
			structure.Level = uint8(unit.OrganizationLevel.Level)
		}
		unitMap[unit.ID] = structure

		// If it's a root unit (no parent), add to root units
		if unit.OrganizationUnitID == nil {
			rootUnits = append(rootUnits, structure)
		}
	}

	// Second pass: build the tree by adding children to their parents
	for _, unit := range units {
		if unit.OrganizationUnitID != nil {
			if parent, exists := unitMap[*unit.OrganizationUnitID]; exists {
				if child, childExists := unitMap[unit.ID]; childExists {
					parent.Children = append(parent.Children, *child)
				}
			}
		}
	}

	// Third pass: recursively copy the complete tree structure
	var copyTree func(*organizationDTOs.OrganizationUnitsDTO) organizationDTOs.OrganizationUnitsDTO
	copyTree = func(node *organizationDTOs.OrganizationUnitsDTO) organizationDTOs.OrganizationUnitsDTO {
		result := *node
		result.Children = make([]organizationDTOs.OrganizationUnitsDTO, 0)

		// Find all direct children from unitMap
		for _, unit := range units {
			if unit.OrganizationUnitID != nil && *unit.OrganizationUnitID == node.ID {
				if child, exists := unitMap[unit.ID]; exists {
					result.Children = append(result.Children, copyTree(child))
				}
			}
		}

		return result
	}

	// Convert to slice of DTOs with complete hierarchy
	result := make([]organizationDTOs.OrganizationUnitsDTO, len(rootUnits))
	for i, structure := range rootUnits {
		result[i] = copyTree(structure)
	}

	return result, nil
}

func (s *OrganizationUnitsService) GetUsersByOrganizationUnit(params *organizationDTOs.OrganizationUnitUserPrimaryDetailsQueryParamsDTO, orgUnitID uint) (*organizationDTOs.OrganizationUnitUserPrimaryDetailsResponseDTO, error) {
	orgUnit, err := s.organizationUnitSvc.GetByID(orgUnitID)
	if err != nil {
		return nil, err
	}
	level := uint(orgUnit.OrganizationLevel.Level)
	// Pass the level value directly to the repository
	return s.repo.GetUsersByOrganizationUnit(params, orgUnitID, level)
}

func (s *OrganizationUnitsService) GetOrganizationUnitSummary(unitID uint) (*organizationDTOs.OrganizationUnitSummaryDTO, error) {
	dto := &organizationDTOs.OrganizationUnitSummaryDTO{
		ResponsibleUnitID: unitID,
	}

	// Step 1+2: Get compliance summary
	compliant, total, statuses, err := s.repo.GetComplianceSummaryByUnit(unitID)
	if err != nil {
		return nil, err
	}

	// Step 3: If no counts → no compliance_status
	if total == 0 {
		dto.ComplianceStatus = nil
	} else {
		// Step 4: Compute percentage
		pct := float64(compliant) * 100.0 / float64(total)

		// Step 5: Match against thresholds
		var matched *organizationDTOs.ComplianceStatusCountDTO
		for _, st := range statuses {
			if st.Type == "exempt" || st.Type == "not_applicable" {
				continue
			}
			if st.ScoreFrom != nil && st.ScoreTo != nil {
				if int(pct) >= *st.ScoreFrom && int(pct) <= *st.ScoreTo {
					matched = &st
					break
				}
			}
		}

		if matched != nil {
			dto.ComplianceStatus = &organizationDTOs.OverallComplianceStatusDTO{
				TitleEn:       matched.TitleEn,
				TitleAr:       matched.TitleAr,
				ColorCode:     matched.ColorCode,
				CompliancePct: math.Round(pct*100) / 100,
			}
		}
	}

	// Risks
	if riskCount, err := s.repo.GetOpenRiskCountByUnit(unitID); err == nil {
		dto.OpenRiskCount = riskCount
	}

	// Control Efficiency
	if ce, err := s.repo.GetControlEfficiencyByUnit(unitID); err == nil {
		dto.OverallControlEfficiency = ce
	}

	return dto, nil
}

func (s *OrganizationUnitsService) UpdatePrimaryOrganizationUnitHead(unitID uint, userID uint, currentUserID uint, departmentHeadUpdateDTO *organizationDTOs.DepartmentHeadUpdateDTO) error {

	// Get primary organization level field
	primaryLevel, err := s.organizationLevelSvc.GetPrimaryLevel()
	if err != nil {
		return err
	}
	primaryField := fmt.Sprintf("organization_unit_%d_id", primaryLevel.Level)

	// Get the user first to check if they're already the head
	user, err := s.userSvc.GetByID(userID)
	if err != nil {
		return err
	}

	// Check if user is already head and belongs to the same unit
	userPrimaryUnitID, err := s.userRepo.GetUserPrimaryUnitID(primaryField, userID)
	if err != nil {
		return err
	}

	// If user is already head of this unit, no need to update
	if user.IsHead && userPrimaryUnitID == unitID {
		return nil
	}

	// Start transaction
	tx := s.repo.BeginTransaction()

	// If setting as head, first unset any existing head in this unit
	if err := s.repo.SetAllIsLatestFalse(unitID, primaryField, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Old data for logging
	oldData := map[string]interface{}{
		"is_head":    user.IsHead,
		"updated_by": user.UpdatedByID,
	}

	// Update user head status
	user.IsHead = true
	user.UpdatedByID = &currentUserID

	// Save the changes
	if err := s.userRepo.Update(user, tx); err != nil {
		tx.Rollback()
		return err
	}

	// New data for logging
	newData := map[string]interface{}{
		"is_head":    user.IsHead,
		"updated_by": user.UpdatedByID,
	}

	// Commit first before logging
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	// Log the activity asynchronously
	go func() {
		err := s.activityLog.LogActivity(
			models.UpdateAction,
			models.User{}.TableName(),
			user.ID,
			map[string]interface{}{
				"old": oldData,
				"new": newData,
			},
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.User{}.TableName(),
				"recordID":  user.ID,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (s *OrganizationUnitsService) GetAllPoliciesByOrganizationUnit(unitID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error) {
	return s.repo.GetAllPoliciesByOrganizationUnit(unitID, params)
}

func (s *OrganizationUnitsService) GetAllRiskByOrganizationUnit(unitID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error) {
	return s.repo.GetAllRiskByOrganizationUnit(unitID, params)
}

func (s *OrganizationUnitsService) GetParentOrganizationUnitWithHead(unitID uint) (*organizationDTOs.ParentOrganizationUnitWithHeadDTO, error) {

	// Get primary organization level field
	primaryLevel, err := s.organizationLevelSvc.GetPrimaryLevel()
	if err != nil {
		return nil, err
	}
	primaryField := fmt.Sprintf("organization_unit_%d_id", primaryLevel.Level)

	// Get parent organization unit with head
	result, err := s.repo.GetParentOrganizationUnitWithHead(unitID, primaryField)
	if err != nil {
		return nil, fmt.Errorf("failed to get parent organization unit: %w", err)
	}

	return result, nil
}

func (s *OrganizationUnitsService) GetAll(params *organizationDTOs.OrganizationUnitReportListQueryParamsDTO) (*organizationDTOs.OrganizationUnitReportListReponseDTO, error) {
	return s.repo.GetAll(params)
}

func (s *OrganizationUnitsService) CreateOrganizationUnitReport(currentUserID uint, organizationUnitReportFormDTO *organizationDTOs.OrganizationUnitReportFormDTO) (uint, error) {
	tx := s.repo.BeginTransaction()

	if err := s.validateOrganizatoinUnitReportFormDTO(organizationUnitReportFormDTO); err != nil {
		return 0, err
	}

	dates, err := s.parseOrganizationUnitReportDates(organizationUnitReportFormDTO)
	if err != nil {
		return 0, err
	}

	organizationUnitReport := &models.OrganizationUnitReport{
		Title:              organizationUnitReportFormDTO.Title,
		StartDate:          *dates.StartDate,
		EndDate:            *dates.EndDate,
		OrganizationUnitID: organizationUnitReportFormDTO.OrganizationUnitID,
		CreatedByID:        currentUserID,
	}

	organizationUnitReportID, err := s.repo.CreateOrganizationUnitReport(organizationUnitReport, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	organizationUnitID := organizationUnitReportFormDTO.OrganizationUnitID

	orgUnit, err := s.organizationUnitSvc.GetByID(organizationUnitID)
	if err != nil {
		return 0, err
	}
	orgLevel := uint(orgUnit.OrganizationLevel.Level)

	overview, err := s.repo.GetOrganizationUnitOverview(organizationUnitID, orgLevel)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if overview.OrganizationUnitHeadNameEn != nil && *overview.OrganizationUnitHeadNameEn == "" {
		overview.OrganizationUnitHeadNameEn = nil
	}
	if overview.OrganizationUnitHeadNameAr != nil && *overview.OrganizationUnitHeadNameAr == "" {
		overview.OrganizationUnitHeadNameAr = nil
	}

	overviewDTO := &models.OrganizationUnitReportOverview{
		OrganizationUnitReportID:   organizationUnitReportID,
		OrganizationUnitTitleEn:    overview.OrganizationUnitTitleEn,
		OrganizationUnitTitleAr:    overview.OrganizationUnitTitleAr,
		OrganizationUnitHeadNameEn: overview.OrganizationUnitHeadNameEn,
		OrganizationUnitHeadNameAr: overview.OrganizationUnitHeadNameAr,
		OrganizationUnitHeadEmail:  overview.OrganizationUnitHeadEmail,
		TotalEmployees:             int(overview.TotalEmployees),
		CreatedByID:                currentUserID,
	}

	if err := s.repo.CreateOrganizationUnitReportOverview(overviewDTO, tx); err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to insert organization_unit_report_overview: %w", err)
	}

	keyMetrics, err := s.repo.GetOrganizationUnitKeyMetrics(organizationUnitID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	keyMetricsDTO := &models.OrganizationUnitReportKeyMetrics{
		OrganizationUnitReportID: organizationUnitReportID,
		TotalUsers:               int(overview.TotalEmployees),
		TotalRisks:               int(keyMetrics.TotalRisks),
		TotalPolicies:            int(keyMetrics.TotalPolicies),
		TotalProcesses:           int(keyMetrics.TotalProcesses),
		TotalCompliances:         int(keyMetrics.TotalCompliances),
		CreatedByID:              currentUserID,
	}

	if err := s.repo.CreateOrganizationUnitReportKeyMetrics(keyMetricsDTO, tx); err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to insert organization_unit_report_key_metrics: %w", err)
	}

	userRoles, err := s.repo.GetUserRoleDistribution(organizationUnitID, orgLevel)
	if err != nil {
		tx.Rollback()
		return 0, nil
	}

	for _, role := range userRoles {
		userRoleDTO := &models.OrganizationUnitReportUserRole{
			OrganizationUnitReportID: organizationUnitReportID,
			RoleTitleEn:              role.RoleTitleEn,
			RoleTitleAr:              role.RoleTitleAr,
			Count:                    int(role.Count),
			CreatedByID:              currentUserID,
		}

		if err := s.repo.CreateOrganizationUnitReportUserRole(userRoleDTO, tx); err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to insert organization_unit_report_user_roles: %w", err)

		}
	}

	users, err := s.repo.GetUsersByOrganizationUnitReport(organizationUnitID, orgLevel)
	if err != nil {
		tx.Rollback()
		return 0, nil
	}

	for _, user := range users {
		userDTO := &models.OrganizationUnitReportUser{
			OrganizationUnitReportID: organizationUnitReportID,
			UserNameEn:               fmt.Sprintf("%s %s", user.FirstNameEn, user.LastNameEn),
			UserNameAr:               fmt.Sprintf("%s %s", user.FirstNameAr, user.LastNameAr),
			DesignationTitleEn:       user.DesignationTitleEn,
			DesignationTitleAr:       user.DesignationTitleAr,
			RoleTitleEn:              user.RoleTitleEn,
			RoleTitleAr:              user.RoleTitleAr,
			CreatedByID:              currentUserID,
		}

		if err := s.repo.CreateOrganizationUnitReportUser(userDTO, tx); err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to insert organization_unit_report_users: %w", err)
		}
	}

	riskSummary, err := s.repo.GetRiskSummaryByOrganizationUnit(organizationUnitID)
	if err != nil {
		tx.Rollback()
		return 0, nil
	}

	riskSummaryDTO := &models.OrganizationUnitReportRiskSummary{
		OrganizationUnitReportID: organizationUnitReportID,
		TotalRiskCount:           int(riskSummary.TotalRisks),
		OpenRiskCount:            int(riskSummary.OpenRisks),
		ClosedRiskCount:          int(riskSummary.ClosedRisks),
		ClosurePercentage:        float32(riskSummary.ClosurePercentage),
		CreatedByID:              currentUserID,
	}

	if err := s.repo.CreateOrganizationUnitReportRiskSummary(riskSummaryDTO, tx); err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to insert organization_unit_report_risk_summary: %w", err)
	}

	riskRatings, err := s.repo.GetRiskRatingDistributionByOrganizationUnit(organizationUnitID)
	if err != nil {
		tx.Rollback()
		return 0, nil
	}

	for _, riskRating := range riskRatings {
		riskRatingDTO := &models.OrganizationUnitReportRiskRating{
			OrganizationUnitReportID: organizationUnitReportID,
			RiskRatingTitleEn:        riskRating.TitleEn,
			RiskRatingTitleAr:        riskRating.TitleAr,
			RiskRatingColorCode:      riskRating.ColorCode,
			Count:                    int(riskRating.Count),
			CreatedByID:              currentUserID,
		}

		if err := s.repo.CreateOrganizationUnitReportRiskRating(riskRatingDTO, tx); err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to insert organization_unit_report_risk_ratings: %w", err)

		}
	}

	riskMitigations, err := s.repo.GetRiskMitigationStatusDistributionByOrganizationUnit(organizationUnitID)
	if err != nil {
		tx.Rollback()
		return 0, nil
	}

	for _, riskMitigation := range riskMitigations {
		riskMitigationDTO := &models.OrganizationUnitReportRiskMitigation{
			OrganizationUnitReportID:      organizationUnitReportID,
			RiskActionPlanStatusTitleEn:   riskMitigation.TitleEn,
			RiskActionPlanStatusTitleAr:   riskMitigation.TitleAr,
			RiskActionPlanStatusColorCode: riskMitigation.ColorCode,
			Count:                         int(riskMitigation.Count),
			CreatedByID:                   currentUserID,
		}

		if err := s.repo.CreateOrganizationUnitReportRiskMitigation(riskMitigationDTO, tx); err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to insert organization_unit_report_risk_mitigations: %w", err)

		}
	}

	risks, err := s.repo.GetRiskByOrganizationUnit(organizationUnitID)
	if err != nil {
		tx.Rollback()
		return 0, nil
	}

	for _, risk := range risks {
		riskDTO := &models.OrganizationUnitReportRisk{
			OrganizationUnitReportID: organizationUnitReportID,
			RiskTitle:                risk.Title,
			Description:              risk.Description,
			RiskCategoryTitleEn:      risk.RiskCategoryTitleEn,
			RiskCategoryTitleAr:      risk.RiskCategoryTitleAr,
			RiskStatusTitleEn:        risk.RiskStatusTitleEn,
			RiskStatusTitleAr:        risk.RiskStatusTitleAr,
			RiskStatusColorCode:      risk.RiskStatusColorCode,
			RiskLikelihoodTitleEn:    risk.CurrentLikelihoodTitleEn,
			RiskLikelihoodTitleAr:    risk.CurrentLikelihoodTitleAr,
			RiskImpactTitleEn:        risk.CurrentImpactTitleEn,
			RiskImpactTitleAr:        risk.CurrentImpactTitleAr,
			RiskRatingTitleEn:        risk.CurrentRiskRatingTitleEn,
			RiskRatingTitleAr:        risk.CurrentRiskRatingTitleAr,
			RiskRatingColorCode:      risk.CurrentRiskRatingColorCode,
			CreatedByID:              currentUserID,
		}

		if err := s.repo.CreateOrganizationUnitReportRisk(riskDTO, tx); err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to insert organization_unit_report_risks: %w", err)

		}
	}

	policySummary, err := s.repo.GetPolicySummaryByOrganizationUnit(organizationUnitID, orgLevel)
	if err != nil {
		tx.Rollback()
		return 0, nil
	}

	policySummaryDTO := &models.OrganizationUnitReportPolicySummary{
		OrganizationUnitReportID: organizationUnitReportID,
		TotalPolicyCount:         int(policySummary.TotalPolicies),
		ActivePolicyCount:        int(policySummary.ActivePolicies),
		UnderReviewPolicyCount:   int(policySummary.PoliciesUnderReview),
		PolicyHealthPercent:      float32(policySummary.PolicyHealthPercentage),
		PolicyAwarenessPercent:   float32(policySummary.PolicyAwarenessPercentage),
		CreatedByID:              currentUserID,
	}

	if err := s.repo.CreateOrganizationUnitReportPolicySummary(policySummaryDTO, tx); err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to insert organization_unit_report_policy_summary: %w", err)
	}

	policies, err := s.repo.GetPolicyByOrganizationUnit(organizationUnitID)
	if err != nil {
		tx.Rollback()
		return 0, nil
	}

	for _, policy := range policies {
		policyDTO := &models.OrganizationUnitReportPolicy{
			OrganizationUnitReportID: organizationUnitReportID,
			DocumentTitleEn:          policy.TitleEn,
			DocumentTitleAr:          policy.TitleAr,
			DocumentCategoryTitleEn:  policy.DocumentCategoryTitleEn,
			DocumentCategoryTitleAr:  policy.DocumentCategoryTitleAr,
			DocumentTypeTitleEn:      policy.DocumentTypeTitleEn,
			DocumentTypeTitleAr:      policy.DocumentTypeTitleAr,
			ResponsibleUserNameEn:    fmt.Sprintf("%s %s", policy.ResponsibleUserFirstNameEn, policy.ResponsibleUserLastNameEn),
			ResponsibleUserNameAr:    fmt.Sprintf("%s %s", policy.ResponsibleUserFirstNameAr, policy.ResponsibleUserLastNameAr),
			Acknowledged:             policy.Acknowledged,
			CreatedByID:              currentUserID,
		}

		if err := s.repo.CreateOrganizationUnitReportPolicy(policyDTO, tx); err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to insert organization_unit_report_policies: %w", err)

		}
	}

	processSummary, err := s.repo.GetProcessHealthSummaryByOrganizationUnit(organizationUnitID)
	if err != nil {
		tx.Rollback()
		return 0, nil
	}

	processSummaryDTO := &models.OrganizationUnitReportProcessSummary{
		OrganizationUnitReportID: organizationUnitReportID,
		ActiveProcessCount:       int(processSummary.ActiveProcesses),
		TotalProcessCount:        int(processSummary.TotalProcesses),
		ActivePercentage:         float32(processSummary.ActivePercentage),
		CreatedByID:              currentUserID,
	}

	if err := s.repo.CreateOrganizationUnitReportProcessSummary(processSummaryDTO, tx); err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to insert organization_unit_report_process_summary: %w", err)
	}

	processes, err := s.repo.GetProcessByOrganizationUnit(organizationUnitID)
	if err != nil {
		tx.Rollback()
		return 0, nil
	}

	for _, process := range processes {
		processDTO := &models.OrganizationUnitReportProcess{
			OrganizationUnitReportID: organizationUnitReportID,
			ProcessTitleEn:           process.TitleEn,
			ProcessTitleAr:           process.TitleAr,
			ProcessGroupTitleEn:      process.ProcessGroupTitleEn,
			ProcessGroupTitleAr:      process.ProcessGroupTitleAr,
			ProcessOwnerNameEn:       fmt.Sprintf("%s %s", process.OwnerFirstNameEn, process.OwnerLastNameEn),
			ProcessOwnerNameAR:       fmt.Sprintf("%s %s", process.OwnerFirstNameAr, process.OwnerLastNameAr),
			ResponsibleUnitTitleEn:   process.ResponsibleUnitTitleEn,
			ResponsibleUnitTitleAr:   process.ResponsibleUnitTitleAr,
			RiskRatingTitleEn:        process.RiskRatingTitleEn,
			RiskRatingTitleAr:        process.RiskRatingTitleAr,
			RiskRatingColorCode:      process.RiskRatingColorCode,
			LifecycleStatusTitleEn:   process.LifecycleStatusTitleEn,
			LifecycleStatusTitleAr:   process.LifecycleStatusTitleAr,
			LifecycleStatusColorCode: process.LifecycleStatusColorCode,
			CreatedByID:              currentUserID,
		}

		if err := s.repo.CreateOrganizationUnitReportProcess(processDTO, tx); err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to insert organization_unit_report_risks: %w", err)
		}
	}

	complianceSummary, err := s.repo.GetComplianceSummaryByOrganizationUnit(organizationUnitID)
	if err != nil {
		tx.Rollback()
		return 0, nil
	}

	complianceStatusesJSON, err := json.Marshal(complianceSummary.ComplianceStatuses)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal compliance statuses: %w", err)
	}

	complianceStatusesStr := string(complianceStatusesJSON)

	complianceSummaryDTO := &models.OrganizationUnitReportComplianceSummary{
		OrganizationUnitReportID: organizationUnitReportID,
		TotalRequirementCount:    int(complianceSummary.TotalRequirements),
		ComplianceStatuses:       &complianceStatusesStr,
		CreatedByID:              currentUserID,
	}

	if err := s.repo.CreateOrganizationUnitReportComplianceSummary(complianceSummaryDTO, tx); err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to insert organizatoin_unit_report_compliance_summary: %w", err)
	}

	complianceRequirements, err := s.repo.GetComplianceRequirementsByOrganizationUnit(organizationUnitID)
	if err != nil {
		tx.Rollback()
		return 0, nil
	}

	for _, complianceRequirement := range complianceRequirements {
		complianceRequirementDTO := &models.OrganizationUnitReportComplianceRequirement{
			OrganizationUnitReportID:             organizationUnitReportID,
			ComplianceRequirementTitleEn:         complianceRequirement.TitleEn,
			ComplianceRequirementTitleAr:         complianceRequirement.TitleAr,
			ComplianceRequirementCategoryTitleEn: complianceRequirement.ComplianceRequirementCategoryTitleEn,
			ComplianceRequirementCategoryTitleAr: complianceRequirement.ComplianceRequirementCategoryTitleAr,
			ComplianceStatusTitleEn:              complianceRequirement.ComplianceStatusTitleEn,
			ComplianceStatusTitleAr:              complianceRequirement.ComplianceStatusTitleAr,
			ComplianceStatusColorCode:            complianceRequirement.ComplianceStatusColorCode,
			SeverityLevelTitleEn:                 complianceRequirement.SeverityLevelTitleEn,
			SeverityLevelTitleAr:                 complianceRequirement.SeverityLevelTitleAr,
			SeverityLevelColorCode:               complianceRequirement.SeverityLevelColorCode,
			CreatedByID:                          currentUserID,
		}

		if err := s.repo.CreateOrganizationUnitReportComplianceRequirement(complianceRequirementDTO, tx); err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to insert organization_unit_report_compliances: %w", err)
		}
	}

	tx.Commit()

	go func() {
		err := s.activityLog.LogActivity(
			models.CreateAction,
			models.OrganizationUnitReport{}.TableName(),
			organizationUnitReportID,
			organizationUnitReport,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{"error": err.Error()})
		}
	}()

	return organizationUnitReportID, nil
}

func (s *OrganizationUnitsService) DeleteOrganizationUnitReport(organizationUnitReportID uint, currentUserID uint) error {
	tx := s.repo.BeginTransaction()

	existingReport, err := s.repo.GetByID(organizationUnitReportID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, organizationUnitReportID)
		}
		return err
	}

	if err := s.repo.DeleteOrganizationUnitReportOverview(organizationUnitReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.repo.DeleteOrganizationUnitReportKeyMetrics(organizationUnitReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.repo.DeleteOrganizationUnitReportUserRole(organizationUnitReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.repo.DeleteOrganizationUnitReportUser(organizationUnitReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.repo.DeleteOrganizationUnitReportRisk(organizationUnitReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.repo.DeleteOrganizationUnitReportPolicy(organizationUnitReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.repo.DeleteOrganizationUnitReportProcess(organizationUnitReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.repo.DeleteOrganizationUnitReportComplianceRequirement(organizationUnitReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.repo.DeleteOrganizationUnitReportRiskSummary(organizationUnitReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.repo.DeleteOrganizationUnitReportRiskRating(organizationUnitReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.repo.DeleteOrganizationUnitReportRiskMitigation(organizationUnitReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.repo.DeleteOrganizationUnitReportPolicySummary(organizationUnitReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.repo.DeleteOrganizationUnitReportProcessSummary(organizationUnitReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.repo.DeleteOrganizationUnitReportComplianceSummary(organizationUnitReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.repo.DeleteOrganizationUnitReport(organizationUnitReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		err := s.activityLog.LogActivity(
			models.DeleteAction,
			models.OrganizationUnitReport{}.TableName(),
			organizationUnitReportID,
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

func (s *OrganizationUnitsService) GetOrganizationUnitReport(organizationUnitReportID uint) (*organizationDTOs.OrganizationUnitSummaryReportDTO, error) {

	reportDetails, err := s.repo.GetByID(organizationUnitReportID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, organizationUnitReportID)
		}
		return nil, err
	}

	if _, err := s.reportTemplateSvc.GetImagePath("logo"); err != nil {
		return nil, utility.NewForbiddenError("logo_image_missing")
	}
	if _, err := s.reportTemplateSvc.GetImagePath("cover-page-bg"); err != nil {
		return nil, utility.NewForbiddenError("cover_image_missing")
	}

	organizationUnitOverview, err := s.repo.GetOrganizationUnitReportOverview(organizationUnitReportID)
	if err != nil {
		return nil, err
	}

	keyMetrics, err := s.repo.GetOrganizationUnitReportKeyMetrics(organizationUnitReportID)
	if err != nil {
		return nil, err
	}

	userRoles, err := s.repo.GetOrganizationUnitReportUserRole(organizationUnitReportID)
	if err != nil {
		return nil, err
	}

	users, err := s.repo.GetOrganizationUnitReportUser(organizationUnitReportID)
	if err != nil {
		return nil, err
	}

	riskSummary, err := s.repo.GetOrganizationUnitReportRiskSummary(organizationUnitReportID)
	if err != nil {
		return nil, err
	}

	riskRatings, err := s.repo.GetOrganizationUnitReportRiskRating(organizationUnitReportID)
	if err != nil {
		return nil, err
	}

	riskMitigations, err := s.repo.GetOrganizationUnitReportRiskMitigation(organizationUnitReportID)
	if err != nil {
		return nil, err
	}

	risks, err := s.repo.GetOrganizationUnitReportRisk(organizationUnitReportID)
	if err != nil {
		return nil, err
	}

	policySummary, err := s.repo.GetOrganizationUnitReportPolicySummary(organizationUnitReportID)
	if err != nil {
		return nil, err
	}

	policies, err := s.repo.GetOrganizationUnitReportPolicy(organizationUnitReportID)
	if err != nil {
		return nil, err
	}

	processSummary, err := s.repo.GetOrganizationUnitReportProcessSummary(organizationUnitReportID)
	if err != nil {
		return nil, err
	}

	processes, err := s.repo.GetOrganizationUnitReportProcess(organizationUnitReportID)
	if err != nil {
		return nil, err
	}

	complianceSummary, err := s.repo.GetOrganizationUnitReportComplianceSummary(organizationUnitReportID)
	if err != nil {
		return nil, err
	}

	if complianceSummary != nil {
		if complianceSummary.ComplianceStatuses != nil {
			var statuses []*organizationDTOs.OrganizationUnitDistributionItemDTO

			if err := json.Unmarshal([]byte(*complianceSummary.ComplianceStatuses), &statuses); err == nil {
				complianceSummary.ComplianceStatusesParsed = statuses
			} else {
				complianceSummary.ComplianceStatusesParsed = []*organizationDTOs.OrganizationUnitDistributionItemDTO{}
			}
		} else {
			complianceSummary.ComplianceStatusesParsed = []*organizationDTOs.OrganizationUnitDistributionItemDTO{}
		}
	}

	complianceRequirments, err := s.repo.GetOrganizationUnitReportComplianceRequirement(organizationUnitReportID)
	if err != nil {
		return nil, err
	}

	response := &organizationDTOs.OrganizationUnitSummaryReportDTO{
		ID:                               reportDetails.ID,
		Title:                            reportDetails.Title,
		OrganizationUnitID:               reportDetails.OrganizationUnitID,
		OrganizationUnitTitleEn:          reportDetails.OrganizationUnitTitleEn,
		OrganizationUnitTitleAr:          reportDetails.OrganizationUnitTitleAr,
		OrganizationUnitColorCode:        reportDetails.OrganizationUnitColorCode,
		CreatedAt:                        reportDetails.CreatedAt,
		CreatedUserFirstNameEn:           reportDetails.CreatedUserFirstNameEn,
		CreatedUserFirstNameAr:           reportDetails.CreatedUserFirstNameAr,
		CreatedUserLastNameEn:            reportDetails.CreatedUserLastNameEn,
		CreatedUserLastNameAr:            reportDetails.CreatedUserLastNameAr,
		OrganizationUnitOverview:         organizationUnitOverview,
		KeyMetrics:                       keyMetrics,
		UserRoleDistribution:             userRoles,
		OrganizationUnitUserDetails:      users,
		RiskSummary:                      riskSummary,
		RiskRatingDistribution:           riskRatings,
		RiskMitigationStatusDistribution: riskMitigations,
		RiskDetails:                      risks,
		PolicySummary:                    policySummary,
		PolicyDetails:                    policies,
		ProcessHealthSummary:             processSummary,
		ProcessDetails:                   processes,
		ComplianceSummary:                complianceSummary,
		ComplianceRequirementDetails:     complianceRequirments,
		AppUrl:                           os.Getenv("APP_URL"),
	}

	return response, nil
}

func (s *OrganizationUnitsService) GetOrganizationUnitReportExport(organizationReportID uint, currentUserID uint) (string, error) {
	organizationReport, err := s.GetOrganizationUnitReport(organizationReportID)
	if err != nil {
		return "", err
	}

	organizationReport.CreatedAt = time.Now().Format("2006-01-02 03:04 PM")

	// Main content template
	templatePaths := []string{"organization/templates/pdf/organization_unit_summary_report.html"}
	htmlContent, err := utility.GeneratePDFContent(organizationReport, templatePaths, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate HTML content: %w", err)
	}

	// Header template - process with data to resolve template variables
	headerTemplatePaths := []string{"organization/templates/pdf/organization_unit_summary_report_header.html"}
	headerContent, err := utility.GeneratePDFContent(organizationReport, headerTemplatePaths, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate header content: %w", err)
	}

	// Footer template - process with data to resolve template variables
	footerTemplatePaths := []string{"organization/templates/pdf/organization_unit_summary_report_footer.html"}
	footerContent, err := utility.GeneratePDFContent(organizationReport, footerTemplatePaths, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate footer content: %w", err)
	}
	coverTemplatePaths := []string{"organization/templates/pdf/organization_unit_summary_report_cover.html"}
	coverContent, err := utility.GeneratePDFContent(organizationReport, coverTemplatePaths, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate footer content: %w", err)
	}

	// Convert HTML → PDF with header and footer
	pdfFileName := fmt.Sprintf("DepartmentSummaryReport_%d.pdf", organizationReportID)
	if err := utility.GeneratePDFWithCoverHeaderFooter(coverContent, htmlContent, pdfFileName, headerContent, footerContent); err != nil {
		return "", fmt.Errorf("failed to generate PDF: %w", err)
	}

	return pdfFileName, nil
}

func (s *OrganizationUnitsService) parseOrganizationUnitReportDates(dto *organizationDTOs.OrganizationUnitReportFormDTO) (*organizationDTOs.OrganizationUnitReportParsedDates, error) {
	parsed := &organizationDTOs.OrganizationUnitReportParsedDates{}

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

func (s *OrganizationUnitsService) validateOrganizatoinUnitReportFormDTO(dto *organizationDTOs.OrganizationUnitReportFormDTO) error {

	if _, err := s.organizationUnitSvc.GetByID(dto.OrganizationUnitID); err != nil {
		return utility.NewValidationError(fmt.Errorf("organization unit does not exist for ID: %d", dto.OrganizationUnitID), "not_exist", "organization unit")
	}

	return nil
}
