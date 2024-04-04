package repositories

import (
	"errors"
	commonDTOs "jethur/grc/common/dtos"
	userDTOs "jethur/grc/organization/dtos"
	organizationModels "jethur/grc/organization/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *userDTOs.UserListQueryParamsDTO) (*userDTOs.UserListResponseDTO, error)
	GetByID(id uint) (*organizationModels.User, error)
	GetByEmail(email string) (*organizationModels.User, error)
	Create(user *organizationModels.User, tx ...utility.DatabaseTransaction) (uint, error)
	Update(user *organizationModels.User, tx ...utility.DatabaseTransaction) error
	UpdatePartial(id uint, updates map[string]interface{}) error
	Delete(id uint) error
	ExistsByEmail(email string) (bool, error)
	ExistsByEmailExcludingID(id uint, email string) (bool, error)

	//User security
	GetUserSecurityByUserID(userId uint) (*organizationModels.UserSecurity, error)
	CreateUserSecurity(userSecurity organizationModels.UserSecurity, tx ...utility.DatabaseTransaction) error
	UpdateUserSecurity(userSecurity *organizationModels.UserSecurity) error

	//User settings
	CreateUserSetting(userSetting organizationModels.UserSetting, tx ...utility.DatabaseTransaction) error

	//User roles
	CreateUserRole(userRoles []organizationModels.UserRole, tx ...utility.DatabaseTransaction) error
	UpdateUserRoles(userID uint, desiredRoleIDs []uint, tx ...utility.DatabaseTransaction) error
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
		Joins("INNER JOIN designations ON designations.id = users.designation_id").
		Joins("INNER JOIN user_statuses ON user_statuses.id = users.user_status_id").
		Joins("INNER JOIN (SELECT user_roles.user_id, STRING_AGG(roles.title_en, ', ') AS roles_en, STRING_AGG(roles.title_ar, ', ') AS roles_ar FROM user_roles JOIN roles ON roles.id = user_roles.role_id GROUP BY user_roles.user_id) roles_agg ON roles_agg.user_id = users.id")

		// Apply search keyword if exists
	if params.Keyword != "" {
		query = query.Where(`
			(users.first_name_en LIKE ? OR
			 users.last_name_en LIKE ? OR
			 users.first_name_ar LIKE ? OR
			 users.last_name_ar LIKE ? OR
			 users.email LIKE ? OR
			 users.mobile_number LIKE ? OR
			 org1.title_en LIKE ? OR
			 org1.title_ar LIKE ? OR
			 org2.title_en LIKE ? OR
			 org2.title_ar LIKE ? OR
			 org3.title_en LIKE ? OR
			 org3.title_ar LIKE ? OR
			 org4.title_en LIKE ? OR
			 org4.title_ar LIKE ? OR
			 org5.title_en LIKE ? OR
			 org5.title_ar LIKE ? OR
			 designations.title_en LIKE ? OR
			 designations.title_ar LIKE ? OR
			 user_statuses.title_en LIKE ? OR
			 user_statuses.title_ar LIKE ?)`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}
	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "first_name_en" // Default sort by title_en if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	//query = query.Order(params.SortBy + " " + params.Sort)

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

func (repo *UserRepository) GetByID(id uint) (*organizationModels.User, error) {
	var user *organizationModels.User
	result := repo.db.First(&user, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return user, result.Error
}

func (repo *UserRepository) Create(user *organizationModels.User, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(user).Error; err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (repo *UserRepository) Update(user *organizationModels.User, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(user).Error; err != nil {
		return err
	}

	return nil
}

func (repo *UserRepository) Delete(id uint) error {

	// Start a transaction
	tx := repo.db.Begin()

	// Check for errors starting the transaction
	if tx.Error != nil {
		return tx.Error
	}

	// Attempt to delete related records first to maintain referential integrity
	if err := tx.Where("user_id = ?", id).Delete(&organizationModels.UserSecurity{}).Error; err != nil {
		tx.Rollback() // Rollback if there's an error
		return err
	}

	if err := tx.Where("user_id = ?", id).Delete(&organizationModels.UserSetting{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("user_id = ?", id).Delete(&organizationModels.UserRole{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Finally, attempt to delete the user
	if err := tx.Delete(&organizationModels.User{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction if all deletions were successful
	return tx.Commit().Error
}

func (repo *UserRepository) CreateUserSecurity(userSecurity organizationModels.UserSecurity, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&userSecurity).Error; err != nil {
		return err
	}
	return nil
}
func (repo *UserRepository) UpdatePartial(id uint, updates map[string]interface{}) error {
	if err := repo.db.Model(&organizationModels.User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return err
	}
	return nil
}

func (repo *UserRepository) GetByEmail(email string) (*organizationModels.User, error) {
	var user *organizationModels.User
	result := repo.db.Where("email = ?", email).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, email)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return user, nil

}

func (repo *UserRepository) CreateUserSetting(userSetting organizationModels.UserSetting, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&userSetting).Error; err != nil {
		return err
	}
	return nil
}

func (repo *UserRepository) CreateUserRole(userRoles []organizationModels.UserRole, tx ...utility.DatabaseTransaction) error {
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

	// Fetch current roles of the user
	var currentRoles []organizationModels.UserRole
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
			if err := db.Create(&organizationModels.UserRole{UserID: userID, RoleID: roleID}).Error; err != nil {
				return err
			}
		}
	}

	// Determine roles to remove
	for _, role := range currentRoles {
		if !contains(desiredRoleIDs, role.RoleID) {
			if err := db.Where("user_id = ? AND role_id = ?", userID, role.RoleID).Delete(&organizationModels.UserRole{}).Error; err != nil {
				return err
			}
		}
	}

	return nil

}

func (repo *UserRepository) ExistsByEmail(email string) (bool, error) {
	var count int64

	// Check for the existence of titleEn or titleAr in a single query
	err := repo.db.Model(&organizationModels.User{}).
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

	err := repo.db.Model(&organizationModels.User{}).
		Where("id <> ? AND email = ?", id, email).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}

// Helper function to check if a slice contains a given uint value
func contains(slice []uint, val uint) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
func (repo *UserRepository) GetUserSecurityByUserID(userId uint) (*organizationModels.UserSecurity, error) {
	var userSecurity *organizationModels.UserSecurity
	result := repo.db.Where("user_id = ?", userId).First(&userSecurity)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, userId)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return userSecurity, nil

}

func (repo *UserRepository) UpdateUserSecurity(userSecurity *organizationModels.UserSecurity) error {

	if err := repo.db.Save(userSecurity).Error; err != nil {
		return err
	}

	return nil
}
