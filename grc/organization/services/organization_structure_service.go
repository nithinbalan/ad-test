package services

import (
	"errors"
	"fmt"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/grc/organization/repositories"
	settingServices "jethur/grc/settings/services"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interface
type OrganizationStructureServiceInterface interface {
	GetOrganizationStructure(sort string, sortBy string) ([]organizationDTOs.OrganizationStructureDTO, error)
	GetUsersByOrganizationUnit(params *organizationDTOs.OrganizationUnitUserPrimaryDetailsQueryParamsDTO, orgUnitID uint) (*organizationDTOs.OrganizationUnitUserPrimaryDetailsResponseDTO, error)
	GetOrganizationUnitDashboard(unitID uint) (*organizationDTOs.OrganizationUnitDashboardDTO, error)
	UpdateDepartmentHead(unitID uint, userID uint, currentUserID uint, departmentHeadUpdateDTO *organizationDTOs.DepartmentHeadUpdateDTO) error
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

func (s *OrganizationStructureService) GetOrganizationUnitDashboard(unitID uint) (*organizationDTOs.OrganizationUnitDashboardDTO, error) {
	dto := &organizationDTOs.OrganizationUnitDashboardDTO{
		ResponsibleUnitID: unitID,
	}

	// Compliance %
	cs, err := s.repo.GetComplianceStatusByUnit(unitID)
	if err != nil {
		return nil, err
	}
	dto.ComplianceStatus = cs

	// Risks
	riskCount, err := s.repo.GetOpenRiskCountByUnit(unitID)
	if err != nil {
		return nil, err
	}
	dto.OpenRiskCount = riskCount

	// Control Efficiency (DTO)
	ce, err := s.repo.GetControlEfficiencyByUnit(unitID)
	if err != nil {
		return nil, err
	}
	dto.OverallControlEfficiency = ce

	return dto, nil
}

func (s *OrganizationStructureService) UpdateDepartmentHead(unitID uint, userID uint, currentUserID uint, departmentHeadUpdateDTO *organizationDTOs.DepartmentHeadUpdateDTO) error {
	// Start transaction
	tx := s.repo.BeginTransaction()

	// Get primary organization level field
	primaryLevel, err := s.organizationLevelSvc.GetPrimaryLevel()
	if err != nil {
		return err
	}
	primaryField := fmt.Sprintf("organization_unit_%d_id", primaryLevel.Level)

	// If setting as head, first unset any existing head in this unit
	if err := s.repo.SetAllIsLatestFalse(unitID, primaryField, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Get the user
	user, err := s.userSvc.GetByID(userID)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &utility.ResourceNotFoundError{}
		}
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
