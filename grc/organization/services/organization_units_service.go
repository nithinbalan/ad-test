package services

import (
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
	Create(currentUserID uint, organizationUnitReportFormDTO *organizationDTOs.OrganizationUnitReportFormDTO) (uint, error)
	Delete(organizationUnitReportID uint, currentUserID uint) error
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

func (s *OrganizationUnitsService) Create(currentUserID uint, organizationUnitReportFormDTO *organizationDTOs.OrganizationUnitReportFormDTO) (uint, error) {
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

	organizationUnitReportID, err := s.repo.Create(organizationUnitReport, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
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

func (s *OrganizationUnitsService) Delete(organizationUnitReportID uint, currentUserID uint) error {

	existingReport, err := s.repo.GetByID(organizationUnitReportID)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(organizationUnitReportID); err != nil {
		return err
	}

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
	var report organizationDTOs.OrganizationUnitSummaryReportDTO

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

	report.Title = reportDetails.Title
	report.OrganizationUnitID = reportDetails.OrganizationUnitID
	report.OrganizationUnitTitleEn = reportDetails.OrganizationUnitTitleEn
	report.OrganizationUnitTitleAr = reportDetails.OrganizationUnitTitleAr
	report.OrganizationUnitColorCode = reportDetails.OrganizationUnitColorCode
	report.GeneratedOn = reportDetails.CreatedAt
	report.GeneratedByUserFirstNameEn = reportDetails.CreatedUserFirstNameAr
	report.GeneratedByUserFirstNameAr = reportDetails.CreatedUserFirstNameAr
	report.GeneratedByUserLastNameEn = reportDetails.CreatedUserLastNameEn
	report.GeneratedByUserLastNameAr = reportDetails.CreatedUserLastNameAr
	report.GeneratedByUserEmail = reportDetails.CreatedUserEmail
	report.GeneratedByUserColor = reportDetails.CreatedUserColor
	report.GeneratedByUserInitial = reportDetails.CreatedUserInitial
	report.GeneratedByUserProfileImagePath = reportDetails.CreatedUserProfileImagePath

	organizationUnitID := reportDetails.OrganizationUnitID

	orgUnit, err := s.organizationUnitSvc.GetByID(reportDetails.OrganizationUnitID)
	if err != nil {
		return nil, err
	}
	level := uint(orgUnit.OrganizationLevel.Level)

	if report.OrganizationUnitOverview, err = s.repo.GetOrganizationUnitOverview(organizationUnitID, level); err != nil {
		return nil, err
	}

	if report.KeyMatrix, err = s.repo.GetOrganizationUnitKeyMetrics(organizationUnitID); err != nil {
		return nil, err
	}

	if report.UserRoleDistribution, err = s.repo.GetUserRoleDistribution(organizationUnitID, level); err != nil {
		return nil, err
	}

	if report.OrganizationUnitUserDetails, err = s.repo.GetUsersByOrganizationUnitReport(organizationUnitID, level); err != nil {
		return nil, err
	}

	if report.RiskSummary, err = s.repo.GetRiskSummaryByOrganizationUnit(organizationUnitID); err != nil {
		return nil, err
	}

	if report.RiskRatingDistribution, err = s.repo.GetRiskRatingDistributionByOrganizationUnit(organizationUnitID); err != nil {
		return nil, err
	}

	if report.RiskMitigationStatusDistribution, err = s.repo.GetRiskMitigationStatusDistributionByOrganizationUnit(organizationUnitID); err != nil {
		return nil, err
	}

	if report.RiskDetails, err = s.repo.GetRiskByOrganizationUnit(organizationUnitID); err != nil {
		return nil, err
	}

	if report.PolicySummary, err = s.repo.GetPolicySummaryByOrganizationUnit(organizationUnitID, level); err != nil {
		return nil, err
	}

	if report.PolicyDetails, err = s.repo.GetPolicyByOrganizationUnit(organizationUnitID); err != nil {
		return nil, err
	}

	if report.ProcessHealthSummary, err = s.repo.GetProcessHealthSummaryByOrganizationUnit(organizationUnitID); err != nil {
		return nil, err
	}

	if report.ProcessDetails, err = s.repo.GetProcessByOrganizationUnit(organizationUnitID); err != nil {
		return nil, err
	}

	if report.ComplianceSummary, err = s.repo.GetComplianceSummaryByOrganizationUnit(organizationUnitID); err != nil {
		return nil, err
	}

	if report.ComplianceRequirementDetails, err = s.repo.GetComplianceRequirementsByOrganizationUnit(organizationUnitID); err != nil {
		return nil, err
	}

	report.AppUrl = os.Getenv("APP_URL")

	return &report, nil
}

func (s *OrganizationUnitsService) GetOrganizationUnitReportExport(organizationReportID uint, currentUserID uint) (string, error) {
	organizationReport, err := s.GetOrganizationUnitReport(organizationReportID)
	if err != nil {
		return "", err
	}

	organizationReport.GeneratedOn = time.Now().Format("2006-01-02 03:04 PM")

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
