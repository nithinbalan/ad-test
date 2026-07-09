package repositories

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	userDTOs "jethur/grc/organization/dtos"
	settingDTOs "jethur/grc/settings/dtos"
	"jethur/utility"
	"strings"
	"time"

	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *userDTOs.UserListQueryParamsDTO) (*userDTOs.UserListResponseDTO, error)
	GetUsersWithPrimaryUnit(primaryLevelID uint, params *userDTOs.PrimaryUnitUserListQueryParamsDTO) (*userDTOs.PrimaryUnitUserListResponseDTO, error)
	GetAllUsers(primaryOrgUnitField string) ([]*userDTOs.UserPrimaryDetailsDTO, error)
	GetByID(id uint) (*models.User, error)
	GetFullDetailsByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Create(user *models.User, tx ...utility.DatabaseTransaction) (uint, error)
	Update(user *models.User, tx ...utility.DatabaseTransaction) error
	UpdatePartial(id uint, updates map[string]interface{}, tx ...utility.DatabaseTransaction) error
	Delete(id uint) error
	ExistsByEmail(email string) (bool, error)
	ExistsByEmailExcludingID(id uint, email string) (bool, error)

	//User security
	CreateUserSecurity(userSecurity models.UserSecurity, tx ...utility.DatabaseTransaction) error

	//User settings
	CreateUserSetting(userSetting models.UserSetting, tx ...utility.DatabaseTransaction) error

	//User roles
	CreateUserRole(userRoles []models.UserRole, tx ...utility.DatabaseTransaction) error
	UpdateUserRoles(userID uint, desiredRoleIDs []uint, tx ...utility.DatabaseTransaction) error

	GetUserLoginHistory(userID uint, params *userDTOs.UserLoginHistoryListQueryParamsDTO) (*userDTOs.UserLoginHistoryListResponseDTO, error)
	CreateUserLoginHistory(userLogin *models.UserLoginHistory, tx ...utility.DatabaseTransaction) error
	GetUserPrimaryUnitID(primaryOrgUnitField string, userID uint) (uint, error)
	GetUserIdsByRole(roleID uint, tx ...utility.DatabaseTransaction) ([]uint, error)

	GetUserNotifications(userID uint, params *userDTOs.UserNotificationListQueryParamsDTO) (*userDTOs.UserNotificationListResponseDTO, error)
	GetUserNewNotificationCount(userID uint) (int64, error)
	MarkAllNotificationsAsRead(userID uint) (int64, error)
	MarkSpecificNotificationsAsRead(userID uint, notificationIDs []uint) (int64, error)
	MarkAllNotificationsAsUnread(userID uint) (int64, error)
	MarkSpecificNotificationsAsUnread(userID uint, notificationIDs []uint) (int64, error)

	GetPermissionsForModule(roleID, moduleID uint) (map[string]bool, error)
	GetUserRoles(userID uint) ([]models.UserRole, error)
	GetUserCounts() (int64, int64, int64, error)

	UpdateUserPrimaryDepartment(userID, orgUnitID uint, fields map[string]*uint, tx ...utility.DatabaseTransaction) error
	UpdateResponsibleUnits(userID, newDepartmentID uint, tx ...utility.DatabaseTransaction) error
	UpdateResponsibleUsers(oldUserID, newUserID uint, tx ...utility.DatabaseTransaction) error
	UpdateResponsibleUserTasks(oldUserID, newUserID uint, tx ...utility.DatabaseTransaction) error

	CountTasksByUser(userID uint, moduleID uint) (int, error)
	GetUserModulesWithOpenTasks(userID uint) ([]settingDTOs.ModulePermissionListDTO, error)

	GetUsersByOrganizationUnit(primaryOrgUnitField string, departmentUnitID uint) ([]*models.User, error)
	IsSuperAdmin(userID uint) (bool, error)
	GetSuperAdmin(tx ...utility.DatabaseTransaction) (*models.User, error)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *UserRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *UserRepository) GetAll(params *userDTOs.UserListQueryParamsDTO) (*userDTOs.UserListResponseDTO, error) {
	var users []*userDTOs.UserListDTO
	var totalCount int64

	query := repo.db.Table("users").
		Select(`
			users.*,
			org1.title_en AS organization_unit_1_title_en,
			org2.title_en AS organization_unit_2_title_en,
			org3.title_en AS organization_unit_3_title_en,
			org4.title_en AS organization_unit_4_title_en,
			org5.title_en AS organization_unit_5_title_en,
			org1.title_ar AS organization_unit_1_title_ar,
			org2.title_ar AS organization_unit_2_title_ar,
			org3.title_ar AS organization_unit_3_title_ar,
			org4.title_ar AS organization_unit_4_title_ar,
			org5.title_ar AS organization_unit_5_title_ar,
			designations.title_en AS designation_title_en,
			designations.title_ar AS designation_title_ar,
			roles_agg.roles_en,
			roles_agg.roles_ar,
			license_types.title_en as license_type_title_en,
			license_types.title_ar as license_type_title_ar,
			license_types.type as license_type_type,
			license_types.color_code as license_type_color_code,
			user_statuses.title_en AS user_status_title_en,
			user_statuses.title_ar AS user_status_title_ar,
			user_statuses.color_code AS user_status_color_code,
			user_statuses.is_active AS user_status_is_active,
			(SELECT COUNT(*) FROM users u2 WHERE u2.reporting_user_id = users.id) AS reporting_users_count

		`).
		Joins("INNER JOIN organization_units org1 ON org1.id = users.organization_unit_1_id").
		Joins("LEFT JOIN organization_units org2 ON org2.id = users.organization_unit_2_id").
		Joins("LEFT JOIN organization_units org3 ON org3.id = users.organization_unit_3_id").
		Joins("LEFT JOIN organization_units org4 ON org4.id = users.organization_unit_4_id").
		Joins("LEFT JOIN organization_units org5 ON org5.id = users.organization_unit_5_id").
		Joins("LEFT JOIN designations ON designations.id = users.designation_id").
		Joins("LEFT JOIN license_types ON license_types.id = users.license_type_id").
		Joins("LEFT JOIN user_statuses ON user_statuses.id = users.user_status_id").
		Joins("LEFT JOIN (SELECT user_roles.user_id, STRING_AGG(roles.title_en, ', ') AS roles_en, STRING_AGG(roles.title_ar, ', ') AS roles_ar FROM user_roles JOIN roles ON roles.id = user_roles.role_id GROUP BY user_roles.user_id) roles_agg ON roles_agg.user_id = users.id").
		Joins("LEFT JOIN users AS reported_to_user ON reported_to_user.id = users.reporting_user_id")

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			if columnName == "reporting_user_id" {
				if len(values) == 1 && strings.ToLower(values[0]) == "null" {
					query = query.Where("users.reporting_user_id IS NULL")
				} else {
					query = query.Where("users.reporting_user_id IN (?)", values)
				}
				continue
			}
			query = query.Where("users."+columnName+" IN (?)", values)
		}
	}

	if len(params.UserIDs) > 0 {
		query = query.Where("users.id IN (?)", params.UserIDs)
	}

	if len(params.RoleIDs) > 0 {
		query = query.Where("EXISTS (SELECT 1 FROM user_roles WHERE user_roles.user_id = users.id AND user_roles.role_id IN (?))", params.RoleIDs)
	}

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("users.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	for _, filter := range params.QuickFilter {
		switch filter {
		case "active":
			query = query.Where("user_statuses.is_active = ?", true)
		case "inactive":
			query = query.Where("user_statuses.is_active = ?", false)
		case "locked":
			query = query.Where("user_statuses.type = ?", "locked")
		case "super_admin":
			query = query.Joins("JOIN user_roles ur ON ur.user_id = users.id").
				Joins("JOIN roles r ON r.id = ur.role_id").
				Where("r.is_super_admin = ?", true)
		case "recent":
			tenDaysAgo := time.Now().AddDate(0, 0, -10)
			query = query.Where("users.created_at >= ?", tenDaysAgo)
		case "logged_in_today":
			query = query.Where("CONVERT(date, users.last_login_at) = CONVERT(date, GETDATE())")
		case "logged_in_this_week":
			query = query.Where("DATEPART(week, users.last_login_at) = DATEPART(week, GETDATE())")
		case "not_logged_in":
			query = query.Where("users.last_login_at IS NULL")
		}
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		query = query.Where(`
        (users.first_name_en LIKE ? OR
         users.last_name_en LIKE ? OR
         users.first_name_ar LIKE ? OR
         users.last_name_ar LIKE ? OR
         users.email LIKE ? OR
         users.mobile_number LIKE ? OR
         designations.title_en LIKE ? OR
         designations.title_ar LIKE ? OR
         user_statuses.title_en LIKE ? OR
         user_statuses.title_ar LIKE ? OR
         roles_agg.roles_en LIKE ? OR
         roles_agg.roles_ar LIKE ?)
        `,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	err := query.Find(&users).Error
	if err != nil {
		return nil, err
	}

	totalPages := 1
	if !params.Export {
		totalPages = int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	}

	response := &userDTOs.UserListResponseDTO{
		Items: users,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}
func (repo *UserRepository) GetUsersWithPrimaryUnit(primaryLevelID uint, params *userDTOs.PrimaryUnitUserListQueryParamsDTO) (*userDTOs.PrimaryUnitUserListResponseDTO, error) {
	var users []*userDTOs.PrimaryUnitUserListDTO
	var totalCount int64

	query := repo.db.Table("users").
		Select(`
			users.id,
			users.first_name_en,
			users.last_name_en,
			users.first_name_ar,
			users.last_name_ar,
			users.email,
			users.mobile_number,
			users.is_ad_user,

			designations.title_en AS designation_title_en,
			license_types.type AS license_type_type,
			user_statuses.is_active AS user_status_is_active
		`).
		Joins("INNER JOIN organization_units org1 ON org1.id = users.organization_unit_1_id").
		Joins("LEFT JOIN organization_units org2 ON org2.id = users.organization_unit_2_id").
		Joins("LEFT JOIN organization_units org3 ON org3.id = users.organization_unit_3_id").
		Joins("LEFT JOIN organization_units org4 ON org4.id = users.organization_unit_4_id").
		Joins("LEFT JOIN organization_units org5 ON org5.id = users.organization_unit_5_id").
		Joins("LEFT JOIN designations ON designations.id = users.designation_id").
		Joins("LEFT JOIN license_types ON license_types.id = users.license_type_id").
		Joins("LEFT JOIN user_statuses ON user_statuses.id = users.user_status_id").
		Where(`
			users.user_status_id = ? AND (
				org1.organization_level_id = ? OR
				org2.organization_level_id = ? OR
				org3.organization_level_id = ? OR
				org4.organization_level_id = ? OR
				org5.organization_level_id = ?
			)
		`,
			1,
			primaryLevelID,
			primaryLevelID,
			primaryLevelID,
			primaryLevelID,
			primaryLevelID,
		)

	// Keyword search
	if params.Keyword != "" {
		like := "%" + params.Keyword + "%"
		query = query.Where(`
		(
			users.first_name_en LIKE ? OR
			users.last_name_en LIKE ? OR
			users.first_name_ar LIKE ? OR
			users.last_name_ar LIKE ? OR
			users.email LIKE ? OR
			users.mobile_number LIKE ? OR
			designations.title_en LIKE ?
		)
		`,
			like, like, like, like,
			like, like,
			like,
		)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Sort
	query = query.Order(params.SortBy + " " + params.Sort)

	// Pagination
	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize)

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &userDTOs.PrimaryUnitUserListResponseDTO{
		Items: users,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}

	return response, nil
}

func (repo *UserRepository) GetAllUsers(primaryOrgUnitField string) ([]*userDTOs.UserPrimaryDetailsDTO, error) {
	var users []*userDTOs.UserPrimaryDetailsDTO

	// Perform the query
	err := repo.db.
		Table("users").
		Select(`
			users.id,
			users.first_name_en,
			users.last_name_en,
			users.email,
			designations.title_en AS designation,
			primary_unit.title_en AS primary_unit,
			users.mobile_number,
			users.profile_image_path`).
		Joins("LEFT JOIN designations ON designations.id = users.designation_id").
		Joins(fmt.Sprintf("LEFT JOIN organization_units AS primary_unit ON primary_unit.id = users.%s", primaryOrgUnitField)).
		Where("users.user_status_id = ?", 1).
		Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (repo *UserRepository) GetByID(id uint) (*models.User, error) {
	var user *models.User

	result := repo.db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (repo *UserRepository) GetFullDetailsByID(id uint) (*models.User, error) {
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

func (repo *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user *models.User

	result := repo.db.
		Where("LOWER(email) = LOWER(?)", email).
		First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil

}

func (repo *UserRepository) Create(user *models.User, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(user).Error; err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (repo *UserRepository) Update(user *models.User, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(user).Error; err != nil {
		return err
	}

	return nil
}

func (repo *UserRepository) UpdatePartial(id uint, updates map[string]interface{}, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Model(&models.User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return err
	}

	return nil
}

func (repo *UserRepository) Delete(id uint) error {
	tx := repo.db.Begin()

	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Where("user_id = ?", id).Delete(&models.UserSecurity{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("user_id = ?", id).Delete(&models.UserSetting{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("user_id = ?", id).Delete(&models.UserRole{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(&models.User{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (repo *UserRepository) ExistsByEmail(email string) (bool, error) {
	var count int64

	// Check for the existence of titleEn or titleAr in a single query
	err := repo.db.Model(&models.User{}).
		Where("email = ?", email).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	// Return true if either email exists, otherwise false
	return count > 0, nil
}

func (repo *UserRepository) ExistsByEmailExcludingID(id uint, email string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.User{}).
		Where("id <> ? AND email = ?", id, email).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}

func (repo *UserRepository) CreateUserSecurity(userSecurity models.UserSecurity, tx ...utility.DatabaseTransaction) error {
	db := repo.db

	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(&userSecurity).Error; err != nil {
		return err
	}

	return nil
}

func (repo *UserRepository) CreateUserSetting(userSetting models.UserSetting, tx ...utility.DatabaseTransaction) error {
	db := repo.db

	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(&userSetting).Error; err != nil {
		return err
	}

	return nil
}

func (repo *UserRepository) CreateUserRole(userRoles []models.UserRole, tx ...utility.DatabaseTransaction) error {
	db := repo.db

	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.CreateInBatches(userRoles, 100).Error; err != nil {
		return err
	}

	return nil
}

func (repo *UserRepository) UpdateUserRoles(userID uint, desiredRoleIDs []uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db

	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// @review
	// Role id comparison can be moved to service layer
	// So the id to be updated (desiredIds) can be used for logs purpsose also
	// Here just update role with ids.

	// Fetch current roles of the user
	var currentRoles []models.UserRole
	if err := db.Where("user_id = ?", userID).Find(&currentRoles).Error; err != nil {
		return err
	}

	currentRoleIDs := make(map[uint]bool)
	for _, role := range currentRoles {
		currentRoleIDs[role.RoleID] = true
	}

	// Determine roles to add
	for _, roleID := range desiredRoleIDs {
		if _, exists := currentRoleIDs[roleID]; !exists {
			if err := db.Create(&models.UserRole{UserID: userID, RoleID: roleID}).Error; err != nil {
				return err
			}
		}
	}

	// Determine roles to remove
	for _, role := range currentRoles {
		if !contains(desiredRoleIDs, role.RoleID) {
			if err := db.Where("user_id = ? AND role_id = ?", userID, role.RoleID).Delete(&models.UserRole{}).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// @Review
// Move this funciton in ot utilty helper, so we cna use everywhere, (if it is common otherwise is fine)
// Choose name which is easy to understand as a common purpsose.
func contains(slice []uint, val uint) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func (repo *UserRepository) GetUserLoginHistory(userID uint, params *userDTOs.UserLoginHistoryListQueryParamsDTO) (*userDTOs.UserLoginHistoryListResponseDTO, error) {
	var userLogins []*userDTOs.UserLoginHistoryListDTO
	var totalCount int64

	query := repo.db.Table("user_login_history").
		Select("id, login_at, login_status, ip_address, location, device_info, browser_info").
		Where("user_login_history.user_id = ?", userID)

	if params.Keyword != "" {
		query = query.Where("(login_status LIKE ? OR location LIKE ? OR ip_address LIKE ? OR device_info LIKE ? OR browser_info LIKE ?)", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	// Apply implemented_at filters
	if params.LoginAtFrom != "" {
		query = query.Where("implemented_at >= ?", params.LoginAtFrom)
	}
	if params.LoginAtTo != "" {
		query = query.Where("implemented_at <= ?", params.LoginAtTo)
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	offset := (params.Page - 1) * params.PageSize

	err := query.Offset(offset).Limit(params.PageSize).Find(&userLogins).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &userDTOs.UserLoginHistoryListResponseDTO{
		Items: userLogins,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *UserRepository) CreateUserLoginHistory(userLogin *models.UserLoginHistory, tx ...utility.DatabaseTransaction) error {
	db := repo.db

	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(userLogin).Error; err != nil {
		return err
	}
	return nil
}

func (repo *UserRepository) GetUserPrimaryUnitID(primaryOrgUnitField string, userID uint) (uint, error) {
	var organizationUnitID uint
	err := repo.db.
		Table("users").
		Select(primaryOrgUnitField).
		Where("id = ?", userID).
		Scan(&organizationUnitID).
		Error
	if err != nil {
		return 0, err
	}

	return organizationUnitID, nil
}

func (repo *UserRepository) GetUserIdsByRole(roleID uint, tx ...utility.DatabaseTransaction) ([]uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var userIDs []uint
	err := db.
		Model(&models.User{}).
		Joins("JOIN user_roles on user_roles.user_id = users.id").
		Where("user_roles.role_id = ?", roleID).
		Pluck("users.id", &userIDs).Error
	if err != nil {
		return nil, err
	}
	return userIDs, nil
}

func (repo *UserRepository) GetUserNotifications(userID uint, params *userDTOs.UserNotificationListQueryParamsDTO) (*userDTOs.UserNotificationListResponseDTO, error) {
	var notifications []*userDTOs.UserNotificationListDTO
	var totalCount int64

	// Update all notifications to set is_new = false for the given user
	if err := repo.db.Table("user_notifications").
		Where("user_notifications.user_id = ? AND user_notifications.is_new = ?", userID, true).
		Update("is_new", false).Error; err != nil {
		return nil, err
	}

	query := repo.db.Table("user_notifications").
		Select(`
			user_notifications.id,
			user_notifications.user_id,
        	users.first_name_en AS user_first_name_en,
			users.first_name_ar AS user_first_name_ar,
			users.last_name_en AS user_last_name_en,
			users.last_name_ar AS user_last_name_ar,
			users.initial AS user_initial,
       		users.color AS user_color,
        	users.email AS user_email,
        	users.profile_image_path AS user_profile_image_path,
			user_notifications.module_id,
			user_notifications.notification_type,
			user_notifications.reference_number,
			user_notifications.heading_en,
			user_notifications.heading_ar,
			user_notifications.title_en,
			user_notifications.title_ar,
			user_notifications.entity,
			user_notifications.action,
			user_notifications.url_param,
			user_notifications.query_param,
			user_notifications.is_new,
			user_notifications.is_read,
			user_notifications.is_system_notification,
			user_notifications.created_at,
			created_user.id AS created_user_id,
        	created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
       		created_user.color AS created_user_color,
        	created_user.email AS created_user_email,
        	created_user.profile_image_path AS created_user_profile_image_path,
			modules.identifier AS module_identifier,
			modules.title_en AS module_title_en,
			modules.title_ar AS module_title_ar
		`).
		Joins("LEFT JOIN users ON users.id = user_notifications.user_id").
		Joins("LEFT JOIN users AS created_user ON created_user.id = user_notifications.created_by_id").
		Joins("LEFT JOIN modules ON modules.id = user_notifications.module_id").
		Where("user_notifications.user_id = ?", userID)

	// Apply dynamic filters to the query
	if len(params.ModuleIDs) > 0 {
		query = query.Where("user_notifications.module_id IN (?)", params.ModuleIDs)
	}

	// Apply quick filters
	for _, filter := range params.QuickFilter {
		switch filter {
		case "new":
			query = query.Where("user_notifications.is_new = ?", true)
		case "read":
			query = query.Where("user_notifications.is_read = ?", true)
		case "unread":
			query = query.Where("user_notifications.is_read = ?", false)

		}
	}

	if params.Keyword != "" {
		query = query.Where(`
			(user_notifications.notification_type LIKE ? OR
			 user_notifications.reference_number LIKE ? OR
			 user_notifications.heading_en LIKE ? OR
			 user_notifications.heading_ar LIKE ? OR
			 user_notifications.title_en LIKE ? OR
			 user_notifications.title_ar LIKE ? OR
			 modules.identifier LIKE ? OR
			 modules.title_en LIKE ? OR
			 modules.title_ar LIKE ?)`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply sorting
	if params.SortBy != "" && params.Sort != "" {
		query = query.Order(params.SortBy + " " + params.Sort)
	} else {
		query = query.Order("user_notifications.created_at DESC") // Default sorting
	}

	// Apply pagination
	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize)

	// Execute the query
	err := query.Find(&notifications).Error
	if err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	// Construct the response
	response := &userDTOs.UserNotificationListResponseDTO{
		Items: notifications,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}

	return response, nil
}

func (repo *UserRepository) GetUserNewNotificationCount(userID uint) (int64, error) {
	var count int64

	err := repo.db.Table("user_notifications").
		Where("user_notifications.user_id = ? AND user_notifications.is_new = ?", userID, true).
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (repo *UserRepository) MarkAllNotificationsAsRead(userID uint) (int64, error) {
	result := repo.db.Model(&models.UserNotification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true)

	return result.RowsAffected, result.Error
}

func (repo *UserRepository) MarkSpecificNotificationsAsRead(userID uint, notificationIDs []uint) (int64, error) {
	result := repo.db.Model(&models.UserNotification{}).
		Where("user_id = ? AND id IN (?) AND is_read = ?", userID, notificationIDs, false).
		Update("is_read", true)

	return result.RowsAffected, result.Error
}

func (repo *UserRepository) MarkAllNotificationsAsUnread(userID uint) (int64, error) {
	result := repo.db.Model(&models.UserNotification{}).
		Where("user_id = ? AND is_read = ?", userID, true).
		Update("is_read", false)

	return result.RowsAffected, result.Error
}

func (repo *UserRepository) MarkSpecificNotificationsAsUnread(userID uint, notificationIDs []uint) (int64, error) {
	result := repo.db.Model(&models.UserNotification{}).
		Where("user_id = ? AND id IN (?) AND is_read = ?", userID, notificationIDs, true).
		Update("is_read", false)

	return result.RowsAffected, result.Error
}

func (repo *UserRepository) GetPermissionsForModule(roleID, moduleID uint) (map[string]bool, error) {
	var permissionTypes []string
	permissions := map[string]bool{}

	err := repo.db.Table("role_permissions").
		Select("role_permission_types.type").
		Joins("JOIN role_permission_types ON role_permissions.role_permission_type_id = role_permission_types.id").
		Where("role_permissions.role_id = ? AND role_permissions.module_id = ?", roleID, moduleID).
		Pluck("role_permission_types.type", &permissionTypes).Error

	if err != nil {
		return nil, err
	}

	for _, permissionType := range permissionTypes {
		permissions[permissionType] = true
	}

	return permissions, nil
}

func (repo *UserRepository) GetUserRoles(userID uint) ([]models.UserRole, error) {
	var userRoles []models.UserRole
	err := repo.db.Where("user_id = ?", userID).Find(&userRoles).Error
	if err != nil {
		return nil, err
	}
	return userRoles, nil
}

func (repo *UserRepository) GetUserCounts() (activeCount, passiveCount, totalActiveCount int64, err error) {
	type Counts struct {
		ActiveCount      int64
		PassiveCount     int64
		TotalActiveUsers int64
	}

	var counts Counts
	err = repo.db.Table("users").
		Select(`
			SUM(CASE WHEN lt.type = 'active' THEN 1 ELSE 0 END) AS active_count,
			SUM(CASE WHEN lt.type = 'passive' THEN 1 ELSE 0 END) AS passive_count,
			SUM(CASE WHEN users.user_status_id IS NOT NULL AND us.is_active = ? THEN 1 ELSE 0 END) AS total_active_users
		`, true).
		Joins("JOIN license_types lt ON lt.id = users.license_type_id").
		Joins("JOIN user_statuses us ON us.id = users.user_status_id").
		Scan(&counts).Error
	if err != nil {
		return 0, 0, 0, err
	}

	return counts.ActiveCount, counts.PassiveCount, counts.TotalActiveUsers, nil
}

func (repo *UserRepository) UpdateUserPrimaryDepartment(userID, orgUnitID uint, fields map[string]*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	for columnName, value := range fields {
		updateQuery := fmt.Sprintf("UPDATE users SET %s = ? WHERE id = ?", columnName)
		if err := db.Exec(updateQuery, value, userID).Error; err != nil {
			return fmt.Errorf("failed to update table users: %w", err)
		}
	}
	return nil
}

func (repo *UserRepository) UpdateResponsibleUnits(userID, newDepartmentID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var tables []string
	query := `
		SELECT TABLE_NAME
    FROM INFORMATION_SCHEMA.COLUMNS
    WHERE COLUMN_NAME IN ('responsible_user_id', 'responsible_unit_id')
    GROUP BY TABLE_NAME
    HAVING COUNT(DISTINCT COLUMN_NAME) = 2;
	`

	err := db.Raw(query).Scan(&tables).Error
	if err != nil {
		return fmt.Errorf("failed to fetch table list: %w", err)
	}

	for _, table := range tables {
		updateQuery := fmt.Sprintf("UPDATE %s SET responsible_unit_id = ? WHERE responsible_user_id = ?", table)
		if err := db.Exec(updateQuery, newDepartmentID, userID).Error; err != nil {
			return fmt.Errorf("failed to update table %s: %w", table, err)
		}
	}
	return nil
}

func (repo *UserRepository) UpdateResponsibleUsers(oldUserID, newUserID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var tables []string
	query := `
		SELECT table_name
		FROM information_schema.columns
		WHERE column_name IN ('responsible_user_id', 'responsible_unit_id')
		GROUP BY table_name
		HAVING COUNT(DISTINCT column_name) = 2;
	`

	err := db.Raw(query).Scan(&tables).Error
	if err != nil {
		return fmt.Errorf("failed to fetch table list: %w", err)
	}

	for _, table := range tables {
		updateQuery := fmt.Sprintf("UPDATE %s SET responsible_user_id = ? WHERE responsible_user_id = ?", table)
		if err := db.Exec(updateQuery, newUserID, oldUserID).Error; err != nil {
			return fmt.Errorf("failed to update table %s: %w", table, err)
		}
	}

	return nil
}

func (repo *UserRepository) UpdateResponsibleUserTasks(oldUserID, newUserID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var statusIDs []int
	err := db.Table("task_statuses").
		Where("is_open = ?", true).
		Pluck("id", &statusIDs).Error
	if err != nil {
		return fmt.Errorf("failed to fetch task status IDs: %w", err)
	}

	err = db.Exec(`
		UPDATE tasks
		SET responsible_user_id = ?
		WHERE responsible_user_id = ? AND task_status_id IN (?)`,
		newUserID, oldUserID, statusIDs).Error

	if err != nil {
		return fmt.Errorf("failed to update tasks: %w", err)
	}
	return nil
}

func (r *UserRepository) CountTasksByUser(userID uint, moduleID uint) (int, error) {
	var count int64
	err := r.db.Table("tasks t").
		Joins("JOIN task_statuses ts ON ts.id = t.task_status_id").
		Where("t.responsible_user_id = ?", userID).
		Where("t.module_id = ?", moduleID).
		Where("ts.type NOT IN ?", []string{"closed"}).
		Count(&count).Error

	return int(count), err
}

func (r *UserRepository) GetUserModulesWithOpenTasks(userID uint) ([]settingDTOs.ModulePermissionListDTO, error) {
	var submodules []settingDTOs.ModulePermissionListDTO

	query := `
	SELECT 
		m.id, 
		m.identifier, 
		m.title_en, 
		m.title_ar, 
		m.[order], 
		m.icon, 
		m.route, 
		m.is_enable,
		COUNT(t.id) AS open_task_count
	FROM modules m
	INNER JOIN tasks t 
		ON t.module_id = m.id
		AND t.responsible_user_id = ?
	INNER JOIN task_statuses ts 
		ON ts.id = t.task_status_id
		AND ts.type NOT IN ('closed')
	INNER JOIN role_permissions rp 
		ON rp.module_id = m.id 
		AND rp.role_permission_type_id = (SELECT id FROM role_permission_types WHERE type = 'view')
		AND rp.is_granted = 1
	INNER JOIN user_roles ur 
		ON ur.role_id = rp.role_id 
		AND ur.user_id = ?
	WHERE m.module_id IS NOT NULL
	GROUP BY 
		m.id, m.identifier, m.title_en, m.title_ar, m.[order], 
		m.icon, m.route, m.is_enable
	ORDER BY m.[order]
	`

	err := r.db.Raw(query, userID, userID).Scan(&submodules).Error
	if err != nil {
		return nil, err
	}

	return submodules, nil
}
func (repo *UserRepository) GetUsersByOrganizationUnit(primaryOrgUnitField string, departmentUnitID uint) ([]*models.User, error) {
	var users []*models.User

	// Build the query dynamically using the column name
	err := repo.db.
		Where(primaryOrgUnitField+" = ?", departmentUnitID).
		Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (repo *UserRepository) IsSuperAdmin(userID uint) (bool, error) {
	var count int64
	err := repo.db.Table("user_roles").
		Joins("JOIN roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ? AND roles.is_super_admin = ?", userID, true).
		Count(&count).Error
	return count > 0, err
}

func (repo *UserRepository) GetSuperAdmin(tx ...utility.DatabaseTransaction) (*models.User, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var user models.User

	err := db.
		Model(&models.User{}).
		Joins("INNER JOIN user_roles ur ON ur.user_id = users.id").
		Joins("INNER JOIN roles r ON r.id = ur.role_id").
		Where("r.type = ?", "super-admin").
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}
