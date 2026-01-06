package services

import (
	"errors"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationRepositories "jethur/grc/organization/repositories"
	"jethur/utility"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type OrganizationChartServiceInterface interface {
	GetUsers() ([]*organizationDTOs.OrganizationChartUserListDTO, error)
	GetUserFullDetailsByID(id uint) (*organizationDTOs.OrganizationChartUserDetailsDTO, error)
}

type OrganizationChartService struct {
	repo organizationRepositories.OrganizationChartRepositoryInterface
}

func NewOrganizationChartService(
	repo organizationRepositories.OrganizationChartRepositoryInterface,
) OrganizationChartServiceInterface {
	return &OrganizationChartService{repo: repo}
}

func (ocs *OrganizationChartService) GetUsers() ([]*organizationDTOs.OrganizationChartUserListDTO, error) {
	return ocs.repo.GetAll()
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
	if user.Designation != nil {
		organizationChartUserDetailsDTO.Designation = user.Designation.TitleEn
	}

	if user.PrimaryOrganizationUnit != nil {
		organizationChartUserDetailsDTO.PrimaryOrganizationID = &user.PrimaryOrganizationUnit.ID
		organizationChartUserDetailsDTO.PrimaryOrganization = user.PrimaryOrganizationUnit.TitleEn
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


