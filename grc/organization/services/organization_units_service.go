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
type OrganizationUnitsServiceInterface interface {
	GetOrganizationUnitsHierarchies(sort string, sortBy string) ([]organizationDTOs.OrganizationUnitsDTO, error)
	GetUsersByOrganizationUnit(params *organizationDTOs.OrganizationUnitUserPrimaryDetailsQueryParamsDTO, orgUnitID uint) (*organizationDTOs.OrganizationUnitUserPrimaryDetailsResponseDTO, error)
	UpdatePrimaryOrganizationUnitHead(unitID uint, userID uint, currentUserID uint, departmentHeadUpdateDTO *organizationDTOs.DepartmentHeadUpdateDTO) error
	GetOrganizationUnitSummary(unitID uint) (*organizationDTOs.OrganizationUnitSummaryDTO, error)
	GetAllPoliciesByOrganizationUnit(unitID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error)
	GetAllRiskByOrganizationUnit(unitID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error)
	GetParentOrganizationUnitWithHead(unitID uint) (*organizationDTOs.ParentOrganizationUnitWithHeadDTO, error)
}

type OrganizationUnitsService struct {
	repo                 repositories.OrganizationUnitsRepositoryInterface
	activityLog          logServices.ActivityLogServiceInterface
	organizationLevelSvc settingServices.OrganizationLevelServiceInterface
	userRepo             repositories.UserRepositoryInterface
	userSvc              UserServiceInterface
	organizationUnitSvc  settingServices.OrganizationUnitServiceInterface
}

func NewOrganizationUnitsService(
	repo repositories.OrganizationUnitsRepositoryInterface,
	activityLog logServices.ActivityLogServiceInterface,
	organizationLevelSvc settingServices.OrganizationLevelServiceInterface,
	userRepo repositories.UserRepositoryInterface,
	userSvc UserServiceInterface,
	organizationUnitSvc settingServices.OrganizationUnitServiceInterface,
) *OrganizationUnitsService {
	return &OrganizationUnitsService{
		repo:                 repo,
		activityLog:          activityLog,
		organizationLevelSvc: organizationLevelSvc,
		userRepo:             userRepo,
		userSvc:              userSvc,
		organizationUnitSvc:  organizationUnitSvc,
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
func (s *OrganizationUnitsService) UpdatePrimaryOrganizationUnitHead(unitID uint, userID uint, currentUserID uint, departmentHeadUpdateDTO *organizationDTOs.DepartmentHeadUpdateDTO) error {

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
