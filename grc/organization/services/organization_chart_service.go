package services

import (
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationRepositories "jethur/grc/organization/repositories"
	settingDTOs "jethur/grc/settings/dtos"
	settingService "jethur/grc/settings/services"
	"jethur/utility"
)

type OrganizationChartServiceInterface interface {
	GetAll(params *organizationDTOs.OrganizationChartUserListQueryParamsDTO) ([]*organizationDTOs.OrganizationChartUserListDTO, error)
	GetUserFullDetailsByID(id uint) (*organizationDTOs.OrganizationChartUserDetailsDTO, error)
	ChangeReportingUser(userID uint, changeReportingUserFormDTO *organizationDTOs.ChangeReportingUserFormDTO) error
}

type OrganizationChartService struct {
	repo                 organizationRepositories.OrganizationChartRepositoryInterface
	organizationLevelSvc settingService.OrganizationLevelServiceInterface
}

func NewOrganizationChartService(
	repo organizationRepositories.OrganizationChartRepositoryInterface,
	organizationLevelSvc settingService.OrganizationLevelServiceInterface,
) OrganizationChartServiceInterface {
	return &OrganizationChartService{
		repo:                 repo,
		organizationLevelSvc: organizationLevelSvc,
	}
}

func (ocs *OrganizationChartService) GetAll(params *organizationDTOs.OrganizationChartUserListQueryParamsDTO) ([]*organizationDTOs.OrganizationChartUserListDTO, error) {
	primaryLevel, err := ocs.organizationLevelSvc.GetPrimaryLevel()
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
		return nil, fmt.Errorf("unsupported primary organization level")
	}

	return ocs.repo.GetAll(params, primaryOrgUnitField)
}

func (ocs *OrganizationChartService) GetUserFullDetailsByID(id uint) (*organizationDTOs.OrganizationChartUserDetailsDTO, error) {

	user, err := ocs.repo.GetUserDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	var organizationChartUserDetailsDTO organizationDTOs.OrganizationChartUserDetailsDTO

	if err := mapstructure.Decode(user, &organizationChartUserDetailsDTO); err != nil {
		return nil, err
	}
	primaryLevel, err := ocs.organizationLevelSvc.GetPrimaryLevel()
	if err != nil {
		return nil, err
	}

	var primaryOrgUnit *models.OrganizationUnit

	switch primaryLevel.Level {
	case 1:
		primaryOrgUnit = user.OrganizationUnit1

	case 2:
		primaryOrgUnit = user.OrganizationUnit2

	case 3:
		primaryOrgUnit = user.OrganizationUnit3

	case 4:
		primaryOrgUnit = user.OrganizationUnit4

	case 5:
		primaryOrgUnit = user.OrganizationUnit5
	default:
		return nil, fmt.Errorf("invalid primary level: %d", primaryLevel.Level)
	}

	if primaryOrgUnit != nil {
		organizationChartUserDetailsDTO.PrimaryOrganizationUnitID = &primaryOrgUnit.ID
		organizationChartUserDetailsDTO.PrimaryOrganizationUnitEn = primaryOrgUnit.TitleEn
		organizationChartUserDetailsDTO.PrimaryOrganizationUnitAr = primaryOrgUnit.TitleAr
	}

	roles := make([]settingDTOs.RoleDTO, 0, len(*user.UserRoles))
	for _, userRole := range *user.UserRoles {
		roles = append(roles, settingDTOs.RoleDTO{
			ID:           userRole.Role.ID,
			TitleEn:      userRole.Role.TitleEn,
			TitleAr:      userRole.Role.TitleAr,
			IsSuperAdmin: userRole.Role.IsSuperAdmin,
			Status:       userRole.Role.Status,
			Type:         userRole.Role.Type,
		})
	}
	organizationChartUserDetailsDTO.Roles = &roles

	return &organizationChartUserDetailsDTO, nil
}

func (ocs *OrganizationChartService) ChangeReportingUser(userID uint, changeReportingUserFormDTO *organizationDTOs.ChangeReportingUserFormDTO) error {
	if changeReportingUserFormDTO.ReportingUserID != nil && *changeReportingUserFormDTO.ReportingUserID == userID {
		return utility.NewForbiddenError("user cannot report to himself")
	}

	err := ocs.repo.UpdateReportingUser(userID, changeReportingUserFormDTO.ReportingUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, userID)
		}
		return err
	}

	return nil
}
