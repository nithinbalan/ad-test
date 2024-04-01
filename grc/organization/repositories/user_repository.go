package repositories

import (
	"errors"
	"jethur/grc/common/db"
	commonDTOs "jethur/grc/common/dtos"
	userDTOs "jethur/grc/organization/dtos"
	organizationModels "jethur/grc/organization/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	BeginTransaction() db.DatabaseTransaction
	GetAll(params *userDTOs.UserListQueryParamsDTO) (*userDTOs.UserListResponseDTO, error)
	GetByID(id uint) (*organizationModels.User, error)
	Create(tx db.DatabaseTransaction, user *organizationModels.User) (uint, error)
	Update(tx db.DatabaseTransaction, user *organizationModels.User) error
	UpdatePartial(id uint, updates map[string]interface{}) error
	GetByEmail(email string) (*organizationModels.User, error)
	GetUserSecurityByUserID(userId uint) (*organizationModels.UserSecurity, error)
	// Delete(id uint) error
	CreateUserSecurity(tx db.DatabaseTransaction, userSecurity organizationModels.UserSecurity) error
	UpdateUserSecurity(userSecurity *organizationModels.UserSecurity) error
	CreateUserSetting(tx db.DatabaseTransaction, userSetting organizationModels.UserSetting) error
	CreateUserRole(tx db.DatabaseTransaction, userRoles []organizationModels.UserRole) error
	UpdateUserRoles(tx db.DatabaseTransaction, userID uint, desiredRoleIDs []uint) error
	ExistsByEmail(email string) (bool, error)
	ExistsByEmailExcludingID(id uint, email string) (bool, error)
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
func (repo *UserRepository) BeginTransaction() db.DatabaseTransaction {
	return db.NewGormDatabaseTransaction(repo.db)
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
			designations.title_ar AS designation_title_ar
		`).
		Joins("LEFT JOIN organization_units org1 ON org1.id = users.organization_unit_1_id").
		Joins("LEFT JOIN organization_units org2 ON org2.id = users.organization_unit_2_id").
		Joins("LEFT JOIN organization_units org3 ON org3.id = users.organization_unit_3_id").
		Joins("LEFT JOIN organization_units org4 ON org4.id = users.organization_unit_4_id").
		Joins("LEFT JOIN organization_units org5 ON org5.id = users.organization_unit_5_id").
		Joins("LEFT JOIN designations ON designations.id = users.designation_id")

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

func (repo *UserRepository) Create(tx db.DatabaseTransaction, user *organizationModels.User) (uint, error) {
	var db *gorm.DB
	if tx != nil {
		db = tx.DB() // Use the transaction's DB instance if available
	} else {
		db = repo.db // Fall back to the repository's DB instance
	}
	if err := db.Create(user).Error; err != nil {
		return 0, err
	}

	return user.ID, nil
}

func (repo *UserRepository) Update(tx db.DatabaseTransaction, user *organizationModels.User) error {
	var db *gorm.DB
	if tx != nil {
		db = tx.DB() // Use the transaction's DB instance if available
	} else {
		db = repo.db // Fall back to the repository's DB instance
	}
	if err := db.Save(user).Error; err != nil {
		return err
	}

	return nil
}

func (repo *UserRepository) CreateUserSecurity(tx db.DatabaseTransaction, userSecurity organizationModels.UserSecurity) error {
	var db *gorm.DB
	if tx != nil {
		db = tx.DB() // Use the transaction's DB instance if available
	} else {
		db = repo.db // Fall back to the repository's DB instance
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

func (repo *UserRepository) CreateUserSetting(tx db.DatabaseTransaction, userSetting organizationModels.UserSetting) error {
	var db *gorm.DB
	if tx != nil {
		db = tx.DB() // Use the transaction's DB instance if available
	} else {
		db = repo.db // Fall back to the repository's DB instance
	}
	if err := db.Create(&userSetting).Error; err != nil {
		return err
	}
	return nil
}

func (repo *UserRepository) CreateUserRole(tx db.DatabaseTransaction, userRoles []organizationModels.UserRole) error {
	var db *gorm.DB
	if tx != nil {
		db = tx.DB() // Use the transaction's DB instance if available
	} else {
		db = repo.db // Fall back to the repository's DB instance
	}
	// Perform the bulk insert in a transaction

	if err := db.CreateInBatches(userRoles, 100).Error; err != nil {
		// If there's an error, roll back the transaction
		return err
	}
	// If the insert is successful, commit the transaction
	return nil

}

func (repo *UserRepository) UpdateUserRoles(tx db.DatabaseTransaction, userID uint, desiredRoleIDs []uint) error {
	var db *gorm.DB
	if tx != nil {
		db = tx.DB() // Use the transaction's DB instance if available
	} else {
		db = repo.db // Fall back to the repository's DB instance
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
			// This role is not assigned to the user, so add it
			if err := db.Create(&organizationModels.UserRole{UserID: userID, RoleID: roleID}).Error; err != nil {
				return err // Rollback on error
			}
		}
		// If the role exists, do nothing (it's already assigned)
	}

	// Determine roles to remove
	for _, role := range currentRoles {
		if !contains(desiredRoleIDs, role.RoleID) {
			// This current role is not in the desired set, so remove it
			if err := db.Where("user_id = ? AND role_id = ?", userID, role.RoleID).Delete(&organizationModels.UserRole{}).Error; err != nil {
				return err // Rollback on error
			}
		}
	}

	return nil // Commit the transaction if all operations were successful

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
