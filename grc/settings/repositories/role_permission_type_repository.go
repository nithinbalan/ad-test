package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type RolePermissionTypeRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string) ([]*models.RolePermissionType, error)
	GetByID(id uint) (*models.RolePermissionType, error)
	GetByType(permissionType string) (*models.RolePermissionType, error)
}

type RolePermissionTypeRepository struct {
	db *gorm.DB
}

func NewRolePermissionTypeRepository(db *gorm.DB) *RolePermissionTypeRepository {
	return &RolePermissionTypeRepository{
		db: db,
	}
}

func (repo *RolePermissionTypeRepository) GetAll(sort string, sortBy string, keyword string) ([]*models.RolePermissionType, error) {
	var rolePermissionTypes []*models.RolePermissionType

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

	query = query.Order(sortBy + " " + sort)
	result := query.Find(&rolePermissionTypes)

	return rolePermissionTypes, result.Error
}

func (repo *RolePermissionTypeRepository) GetByID(id uint) (*models.RolePermissionType, error) {
	var rolePermissionType *models.RolePermissionType
	result := repo.db.First(&rolePermissionType, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return rolePermissionType, nil
}

func (repo *RolePermissionTypeRepository) GetByType(permissionType string) (*models.RolePermissionType, error) {
	var rolePermissionType *models.RolePermissionType

	result := repo.db.Where("type = ?", permissionType).First(&rolePermissionType)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Role Permission Type with type "+permissionType)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return rolePermissionType, nil
}
