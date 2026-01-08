package services

import (
	"errors"
	"fmt"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationRepositories "jethur/grc/organization/repositories"
	settingService "jethur/grc/settings/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type OrganizationChartServiceInterface interface {
	GetAll(params *organizationDTOs.OrganizationChartUserListQueryParamsDTO) ([]*organizationDTOs.OrganizationChartUserListDTO, error)
	GetUserFullDetailsByID(id uint) (*organizationDTOs.OrganizationChartUserDetailsDTO, error)
	ChangeReportingUser(userID uint,changeReportingUserFormDTO *organizationDTOs.ChangeReportingUserFormDTO) error
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
		organizationChartUserDetailsDTO.PrimaryOrganizationID = &primaryOrgUnit.ID
		organizationChartUserDetailsDTO.PrimaryOrganizationEn = primaryOrgUnit.TitleEn
		organizationChartUserDetailsDTO.PrimaryOrganizationAr = primaryOrgUnit.TitleAr
	}

	if user.UserRoles != nil {
		organizationChartUserDetailsDTO.SystemRoles = make([]organizationDTOs.OrganizationChartRoleDTO, 0)

		for _, ur := range *user.UserRoles {
			if ur.Role.ID != 0 {
				organizationChartUserDetailsDTO.SystemRoles = append(organizationChartUserDetailsDTO.SystemRoles,
					organizationDTOs.OrganizationChartRoleDTO{
						ID:      ur.Role.ID,
						TitleEn: ur.Role.TitleEn,
						TitleAr: ur.Role.TitleAr,
					},
				)
			}
		}
	}

	return &organizationChartUserDetailsDTO, nil
}

func (ocs *OrganizationChartService) ChangeReportingUser(userID uint,changeReportingUserFormDTO *organizationDTOs.ChangeReportingUserFormDTO) error {
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
