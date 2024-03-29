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
	GetAll(params *userDTOs.UserListQueryParamsDTO) (*userDTOs.UserListResponseDTO, error)
	GetByID(id uint) (*organizationModels.User, error)
	Create(user *organizationModels.User) (uint, error)
	Update(user *organizationModels.User) error
	UpdatePartial(id uint, updates map[string]interface{}) error
	GetByEmail(email string) (*organizationModels.User, error)
	GetUserSecurityByUserID(userId uint) (*organizationModels.UserSecurity, error)
	// Delete(id uint) error
	CreateUserSecurity(userSecurity *organizationModels.UserSecurity) error
	CreateUserSetting(userSetting *organizationModels.UserSetting) error
	UpdateUserSecurity(userSecurity *organizationModels.UserSecurity) error
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
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

func (repo *UserRepository) Create(user *organizationModels.User) (uint, error) {
	if err := repo.db.Create(user).Error; err != nil {
		return 0, err
	}

	return user.ID, nil
}

func (repo *UserRepository) Update(user *organizationModels.User) error {

	if err := repo.db.Save(user).Error; err != nil {
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

// CreateUserSecurity creates a new UserSecurity record associated with a user.
func (repo *UserRepository) CreateUserSecurity(userSecurity *organizationModels.UserSecurity) error {
	if err := repo.db.Create(userSecurity).Error; err != nil {
		return err
	}
	return nil
}

// CreateUserSetting creates a new UserSetting record.
func (repo *UserRepository) CreateUserSetting(userSetting *organizationModels.UserSetting) error {
	if err := repo.db.Create(userSetting).Error; err != nil {
		return err
	}
	return nil
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
