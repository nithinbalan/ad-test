package repositories

import (
	"gorm.io/gorm"

	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
)

type OrganizationChartRepositoryInterface interface {
	GetAll() ([]*organizationDTOs.OrganizationChartUserListDTO, error)
	GetUserDetailsByID(id uint) (*models.User, error)
}

type OrganizationChartRepository struct {
	db *gorm.DB
}

func NewOrganizationChartRepository(db *gorm.DB) OrganizationChartRepositoryInterface {
	return &OrganizationChartRepository{db: db}
}

// repositories/user_repository.go

func (repo *OrganizationChartRepository) GetAll() ([]*organizationDTOs.OrganizationChartUserListDTO, error) {
	var users []*organizationDTOs.OrganizationChartUserListDTO

	err := repo.db.
		Table("users").
		Select(`
		users.id,
		users.first_name_en,
		users.last_name_en,
		users.first_name_ar,
		users.last_name_ar,
		users.email,
		users.mobile_number,
		users.profile_image_path,
		users.initial,
		users.color,

		org1.title_en AS organization_unit_1_title_en,
		org1.title_ar AS organization_unit_1_title_ar,

		org2.title_en AS organization_unit_2_title_en,
		org2.title_ar AS organization_unit_2_title_ar,

		designations.title_en AS designation_title_en,
		designations.title_ar AS designation_title_ar,

		roles_agg.roles_en,
		roles_agg.roles_ar
	`).
		Joins("LEFT JOIN organization_units org1 ON org1.id = users.organization_unit_1_id").
		Joins("LEFT JOIN organization_units org2 ON org2.id = users.organization_unit_2_id").
		Joins("LEFT JOIN designations ON designations.id = users.designation_id").
		Joins(`
		LEFT JOIN (
			SELECT 
				user_roles.user_id,
				STRING_AGG(roles.title_en, ', ') AS roles_en,
				STRING_AGG(roles.title_ar, ', ') AS roles_ar
			FROM user_roles
			JOIN roles ON roles.id = user_roles.role_id
			GROUP BY user_roles.user_id
		) roles_agg ON roles_agg.user_id = users.id
	`).
		Order("users.first_name_en ASC").
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (repo *OrganizationChartRepository) GetUserDetailsByID(id uint) (*models.User, error) {

	var user *models.User

	result := repo.db.
		Preload("OrganizationUnit1").
		Preload("OrganizationUnit2").
		Preload("OrganizationUnit3").
		Preload("OrganizationUnit4").
		Preload("OrganizationUnit5").
		Preload("UserStatus").
		Preload("Designation").
		Preload("UserRoles").
		Preload("UserSetting").
		Preload("UserSetting.Timezone").
		Preload("UserRoles.Role").
		Preload("UserRoles.Role.LicenseType").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("ReportingUser").
		Preload("LicenseType").
		Preload("DelegationUser").
		First(&user, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}
