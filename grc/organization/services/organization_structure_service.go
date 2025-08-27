package services

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/grc/organization/repositories"
	settingServices "jethur/grc/settings/services"
	"jethur/utility"
)

// Interface
type OrganizationStructureServiceInterface interface {
	GetOrganizationStructure(sort string, sortBy string) ([]organizationDTOs.OrganizationStructureDTO, error)
	GetUsersByOrganizationUnit(params *organizationDTOs.OrganizationUnitUserPrimaryDetailsQueryParamsDTO, orgUnitID uint) (*organizationDTOs.OrganizationUnitUserPrimaryDetailsResponseDTO, error)
	UpdatePrimaryOrganizationUnitHead(unitID uint, userID uint, currentUserID uint, departmentHeadUpdateDTO *organizationDTOs.DepartmentHeadUpdateDTO) error
	GetOrganizationUnitSummary(unitID uint) (*organizationDTOs.OrganizationUnitSummaryDTO, error)
	GetAllPoliciesByOrganizationUnit(unitID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error)
	GetAllRiskByOrganizationUnit(unitID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error)
	GetParentOrganizationUnitWithHead(unitID uint) (*organizationDTOs.ParentOrganizationUnitWithHeadDTO, error)
}

type OrganizationStructureService struct {
	repo                 repositories.OrganizationStructureRepositoryInterface
	activityLog          logServices.ActivityLogServiceInterface
	organizationLevelSvc settingServices.OrganizationLevelServiceInterface
	userRepo             repositories.UserRepositoryInterface
	userSvc              UserServiceInterface
}

func NewOrganizationStructureService(
	repo repositories.OrganizationStructureRepositoryInterface,
	activityLog logServices.ActivityLogServiceInterface,
	organizationLevelSvc settingServices.OrganizationLevelServiceInterface,
	userRepo repositories.UserRepositoryInterface,
	userSvc UserServiceInterface,
) *OrganizationStructureService {
	return &OrganizationStructureService{
		repo:                 repo,
		activityLog:          activityLog,
		organizationLevelSvc: organizationLevelSvc,
		userRepo:             userRepo,
		userSvc:              userSvc,
	}
}

func (s *OrganizationStructureService) GetOrganizationStructure(sort string, sortBy string) ([]organizationDTOs.OrganizationStructureDTO, error) {
	// Get all organization units
	units, err := s.repo.GetOrganizationStructure(sort, sortBy)
	if err != nil {
		return nil, err
	}

	// Create a map for quick lookup of units by ID
	unitMap := make(map[uint]*organizationDTOs.OrganizationStructureDTO)
	var rootUnits []*organizationDTOs.OrganizationStructureDTO

	// First pass: create all nodes
	for _, unit := range units {
		structure := &organizationDTOs.OrganizationStructureDTO{
			ID:      unit.ID,
			TitleEn: unit.TitleEn,
			TitleAr: unit.TitleAr,
		}
		unitMap[unit.ID] = structure

		// If it's a root unit, add to root units
		if unit.OrganizationUnitID == nil {
			rootUnits = append(rootUnits, structure)
		}
	}

	// Second pass: build the tree
	for _, unit := range units {
		if unit.OrganizationUnitID != nil {
			if parent, exists := unitMap[*unit.OrganizationUnitID]; exists {
				parent.Children = append(parent.Children, *unitMap[unit.ID])
			}
		}
	}

	// Convert to slice of DTOs
	result := make([]organizationDTOs.OrganizationStructureDTO, len(rootUnits))
	for i, structure := range rootUnits {
		result[i] = *structure
	}

	return result, nil
}

func (s *OrganizationStructureService) GetUsersByOrganizationUnit(params *organizationDTOs.OrganizationUnitUserPrimaryDetailsQueryParamsDTO, orgUnitID uint) (*organizationDTOs.OrganizationUnitUserPrimaryDetailsResponseDTO, error) {
	var primaryOrgUnitField string

	// Get the primary organization level
	primaryLevel, err := s.organizationLevelSvc.GetPrimaryLevel()
	if err != nil {
		return nil, err
	}

	// Determine the field name based on the primary level
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

	return s.repo.GetUsersByOrganizationUnit(params, orgUnitID, primaryOrgUnitField)
}

func (s *OrganizationStructureService) GetOrganizationUnitSummary(unitID uint) (*organizationDTOs.OrganizationUnitSummaryDTO, error) {
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
				CompliancePct: pct,
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
func (s *OrganizationStructureService) UpdatePrimaryOrganizationUnitHead(unitID uint, userID uint, currentUserID uint, departmentHeadUpdateDTO *organizationDTOs.DepartmentHeadUpdateDTO) error {

	// Get primary organization level field
	primaryLevel, err := s.organizationLevelSvc.GetPrimaryLevel()
	if err != nil {
		return err
	}
	primaryField := fmt.Sprintf("organization_unit_%d_id", primaryLevel.Level)

	// Start transaction
	tx := s.repo.BeginTransaction()

	// If setting as head, first unset any existing head in this unit
	if err := s.repo.SetAllIsLatestFalse(unitID, primaryField, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Get the user
	user, err := s.userSvc.GetByID(userID)
	if err != nil {
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
			nil,
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

func (s *OrganizationStructureService) GetAllPoliciesByOrganizationUnit(unitID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error) {
	return s.repo.GetAllPoliciesByOrganizationUnit(unitID, params)
}

func (s *OrganizationStructureService) GetAllRiskByOrganizationUnit(unitID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error) {
	return s.repo.GetAllRiskByOrganizationUnit(unitID, params)
}

func (s *OrganizationStructureService) GetParentOrganizationUnitWithHead(unitID uint) (*organizationDTOs.ParentOrganizationUnitWithHeadDTO, error) {

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
