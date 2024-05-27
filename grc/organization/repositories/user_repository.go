package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	userDTOs "jethur/grc/organization/dtos"
	"jethur/utility"
	"time"

	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *userDTOs.UserListQueryParamsDTO) (*userDTOs.UserListResponseDTO, error)
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
	GetUserLogins(page int, pageSize int) (*userDTOs.UserLoginListResponseDTO, error)
	CreateUserLogin(userLogin *models.UserLogin, tx ...utility.DatabaseTransaction) error
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
			user_statuses.title_en AS user_status_title_en,
			user_statuses.title_ar AS user_status_title_ar,
			user_statuses.color_code AS user_status_color_code
		`).
		Joins("INNER JOIN organization_units org1 ON org1.id = users.organization_unit_1_id").
		Joins("LEFT JOIN organization_units org2 ON org2.id = users.organization_unit_2_id").
		Joins("LEFT JOIN organization_units org3 ON org3.id = users.organization_unit_3_id").
		Joins("LEFT JOIN organization_units org4 ON org4.id = users.organization_unit_4_id").
		Joins("LEFT JOIN organization_units org5 ON org5.id = users.organization_unit_5_id").
		Joins("LEFT JOIN designations ON designations.id = users.designation_id").
		Joins("INNER JOIN user_statuses ON user_statuses.id = users.user_status_id").
		Joins("INNER JOIN (SELECT user_roles.user_id, STRING_AGG(roles.title_en, ', ') AS roles_en, STRING_AGG(roles.title_ar, ', ') AS roles_ar FROM user_roles JOIN roles ON roles.id = user_roles.role_id GROUP BY user_roles.user_id) roles_agg ON roles_agg.user_id = users.id")

		// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply quick filters
	for _, filter := range params.QuickFilter {
		switch filter {
		case "recently_added":
			thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
			query = query.Where("users.created_at >= ?", thirtyDaysAgo)
		case "admin_users":
			query = query.Where("roles_agg.roles_en = ? OR roles_agg.roles_ar = ?", "Admin", "مسؤول")
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
			 user_statuses.title_ar LIKE ?)`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the queryc
	err := query.Offset(offset).Limit(params.PageSize).Find(&users).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

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
		Preload("UserRoles.Role").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&user, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (repo *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user *models.User

	result := repo.db.Where("email = ?", email).First(&user)

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

func (repo *UserRepository) GetUserLogins(page int, pageSize int) (*userDTOs.UserLoginListResponseDTO, error) {
	var userLogins []*userDTOs.UserLoginListDTO
	var totalCount int64

	// Initialize the query
	query := repo.db.Table("user_logins").
		Select("id, login_at, login_status, ip_address, location, device_info, browser_info")

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	offset := (page - 1) * pageSize

	// Execute the query
	err := query.Offset(offset).Limit(pageSize).Find(&userLogins).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(pageSize) - 1) / int64(pageSize))

	response := &userDTOs.UserLoginListResponseDTO{
		Items: userLogins,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: page,
			PageSize:    pageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *UserRepository) CreateUserLogin(userLogin *models.UserLogin, tx ...utility.DatabaseTransaction) error {
	db := repo.db

	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(userLogin).Error; err != nil {
		return err
	}
	return nil
}
