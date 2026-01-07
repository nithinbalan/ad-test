package repositories

import (
	"fmt"
	"gorm.io/gorm"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	"strings"
)

type OrganizationChartRepositoryInterface interface {
	GetAll(params *organizationDTOs.OrganizationChartUserListQueryParamsDTO, primaryOrgUnitField string) ([]*organizationDTOs.OrganizationChartUserListDTO, error)
	GetUserDetailsByID(id uint) (*models.User, error)
	UpdateReportingUser(userID uint, reportingUserID *uint) error
}

type OrganizationChartRepository struct {
	db *gorm.DB
}

func NewOrganizationChartRepository(db *gorm.DB) OrganizationChartRepositoryInterface {
	return &OrganizationChartRepository{db: db}
}

func (repo *OrganizationChartRepository) GetAll(params *organizationDTOs.OrganizationChartUserListQueryParamsDTO, primaryOrgUnitField string) ([]*organizationDTOs.OrganizationChartUserListDTO, error) {

	var users []*organizationDTOs.OrganizationChartUserListDTO
	query := repo.db.
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
			designations.title_en AS designation,
			primary_unit.id       AS primary_organization_unit_id,
			reporting_user.id AS reporting_user_id,
		    reporting_user.first_name_en AS reporting_user_first_name_en,
		    reporting_user.last_name_en AS reporting_user_last_name_en,
			reporting_user.first_name_ar AS reporting_user_first_name_ar,
		    reporting_user.last_name_ar AS reporting_user_last_name_ar,
			roles_agg.roles_en,
			roles_agg.roles_ar
		`).
		Joins(fmt.Sprintf("LEFT JOIN organization_units AS primary_unit ON primary_unit.id = users.%s", primaryOrgUnitField)).
		Joins("LEFT JOIN designations ON designations.id = users.designation_id").
		Joins("LEFT JOIN user_statuses ON user_statuses.id = users.user_status_id").
		Joins("LEFT JOIN users AS reporting_user ON reporting_user.id = users.reporting_user_id").
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
		`).Where("user_statuses.is_active = ?", true)

	for column, values := range params.Filters {
		if len(values) == 0 {
			continue
		}

		if column == "reporting_user_id" {
			if len(values) == 1 && strings.ToLower(values[0]) == "null" {
				query = query.Where("users.reporting_user_id IS NULL")
			} else {
				query = query.Where("users.reporting_user_id IN (?)", values)
			}
			continue
		}
       
		query = query.Where("users."+column+" IN (?)", values)
	}

	err := query.
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
		Preload("UserRoles.Role").
		Preload("ReportingUser").
		First(&user, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (repo *OrganizationChartRepository) UpdateReportingUser(userID uint, reportingUserID *uint) error {

	result := repo.db.
		Model(&models.User{}).
		Where("id = ?", userID).
		Update("reporting_user_id", reportingUserID)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
