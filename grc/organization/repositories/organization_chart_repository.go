package repositories

import (
	"gorm.io/gorm"

	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
)

type OrganizationChartRepositoryInterface interface {
	GetUsersForOrgChart() ([]*organizationDTOs.OrganizationChartUserDTO, error)
	GetUserDetailsByID(id uint) (*models.User, error)
}

type OrganizationChartRepository struct {
	db *gorm.DB
}

func NewOrganizationChartRepository(db *gorm.DB) OrganizationChartRepositoryInterface {
	return &OrganizationChartRepository{db: db}
}

func (repo *OrganizationChartRepository) GetUsersForOrgChart() ([]*organizationDTOs.OrganizationChartUserDTO, error) {
	var users []*organizationDTOs.OrganizationChartUserDTO

	err := repo.db.
		Table("users").
		Select(`
			users.id,
			users.first_name_en,
			users.last_name_en,
			users.first_name_ar,
			users.last_name_ar,
			users.email,
			designations.title_en AS designation,
			users.reporting_user_id
		`).
		Joins("LEFT JOIN designations ON designations.id = users.designation_id").
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (repo *OrganizationChartRepository) GetUserDetailsByID(id uint) (*models.User, error) {

	var user models.User

	err := repo.db.
		Preload("Designation").
		Preload("UserRoles").
		Preload("UserRoles.Role").         
		Preload("ReportingUser").   
		Preload("ReportingUser.Designation").
		Preload("OrganizationUnit1.OrganizationLevel").
		Preload("OrganizationUnit2.OrganizationLevel").
		Preload("OrganizationUnit3.OrganizationLevel").
		Preload("OrganizationUnit4.OrganizationLevel").
		Preload("OrganizationUnit5.OrganizationLevel").
		First(&user, id).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

