package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type RoleRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.Role, error)
	GetByID(id uint) (*models.Role, error)
	Create(role *models.Role) (*models.Role, error)
	Update(id uint, role *models.Role) (*models.Role, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)

	CreateRolePermission(rolePermission *models.RolePermission, tx ...utility.DatabaseTransaction) (uint, error)
	DeleteRolePermissionsByRoleID(roleID uint, tx ...utility.DatabaseTransaction) error

	GetPermissionsForModule(roleID, moduleID uint) (map[string]bool, error)
	GetAllPermissionTypes() ([]string, error)

	GetSuperAdminRole() (*models.Role, error)
}

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{
		db: db,
	}
}

func (repo *RoleRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *RoleRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.Role, error) {
	var roles []*models.Role

	query := repo.db

	if keyword != "" {
		query = query.Where("(title_en LIKE ? OR title_ar LIKE ?)", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "title_en"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	if status == "active" || status == "inactive" {
		query = query.Where("status = ?", status)
	}

	query = query.Order(sortBy + " " + sort)
	result := query.Find(&roles)

	return roles, result.Error
}

func (repo *RoleRepository) GetByID(id uint) (*models.Role, error) {
	var role *models.Role

	result := repo.db.First(&role, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return role, result.Error
}

func (repo *RoleRepository) Create(role *models.Role) (*models.Role, error) {
	result := repo.db.Create(&role)

	return role, result.Error
}

func (repo *RoleRepository) Update(id uint, updatedRole *models.Role) (*models.Role, error) {
	var role *models.Role

	if err := repo.db.First(&role, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&role).Updates(updatedRole).Error; err != nil {
		return nil, err
	}

	return role, nil
}

func (repo *RoleRepository) Delete(id uint) error {

	var role *models.Role
	if err := repo.db.First(&role, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.Role{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *RoleRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.Role{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.Role{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil
}

func (repo *RoleRepository) Deactivate(ids []uint) ([]uint, error) {

	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.Role{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.Role{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil

}

func (repo *RoleRepository) ExistsByTitle(titleEn string, titleAr string) (bool, error) {
	var count int64

	// Check for the existence of titleEn or titleAr in a single query
	err := repo.db.Model(&models.Role{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	// Return true if either title exists, otherwise false
	return count > 0, nil
}

func (repo *RoleRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.Role{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}

func (repo *RoleRepository) DeleteRolePermissionsByRoleID(roleID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	result := db.Where("role_id = ?", roleID).Delete(&models.RolePermission{})
	return result.Error
}

func (repo *RoleRepository) CreateRolePermission(rolePermission *models.RolePermission, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	result := db.Create(rolePermission)
	if result.Error != nil {
		return 0, result.Error
	}
	return rolePermission.ID, nil
}

func (repo *RoleRepository) GetAllPermissionTypes() ([]string, error) {
	var permissionTypes []string
	err := repo.db.Table("role_permission_types").Select("type").Find(&permissionTypes).Error
	if err != nil {
		return nil, err
	}
	return permissionTypes, nil
}

func (repo *RoleRepository) GetPermissionsForModule(roleID, moduleID uint) (map[string]bool, error) {
	var permissionTypes []string
	permissions := map[string]bool{}

	err := repo.db.Table("role_permissions").
		Select("role_permission_types.type").
		Joins("JOIN role_permission_types ON role_permissions.role_permission_type_id = role_permission_types.id").
		Where("role_permissions.role_id = ? AND role_permissions.module_id = ? AND role_permissions.is_granted = ?", roleID, moduleID, 1).
		Pluck("role_permission_types.type", &permissionTypes).Error

	if err != nil {
		return nil, err
	}

	for _, permissionType := range permissionTypes {
		permissions[permissionType] = true
	}

	return permissions, nil
}

func (repo *RoleRepository) GetSuperAdminRole() (*models.Role, error) {
	var role models.Role
	err := repo.db.Where("is_super_admin = ?", true).Where("status = ?", "active").First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}
