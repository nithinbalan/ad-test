package repositories

import (
	"jethur/grc/models"
	"jethur/grc/settings/dtos"

	"gorm.io/gorm"
)

type WorkflowCustomRoleRepositoryInterface interface {
	GetAll(params *dtos.WorkflowCustomRoleListQueryParamsDTO) ([]*models.WorkflowCustomRole, error)
	GetByType(roleType string) (*models.WorkflowCustomRole, error)
}

type WorkflowCustomRoleRepository struct {
	db *gorm.DB
}

func NewWorkflowCustomRoleRepository(db *gorm.DB) WorkflowCustomRoleRepositoryInterface {
	return &WorkflowCustomRoleRepository{
		db: db,
	}
}

func (repo *WorkflowCustomRoleRepository) GetAll(params *dtos.WorkflowCustomRoleListQueryParamsDTO) ([]*models.WorkflowCustomRole, error) {
	var roles []*models.WorkflowCustomRole

	query := repo.db.Table("workflow_custom_roles").
		Select(`
			workflow_custom_roles.id,
			workflow_custom_roles.title_en,
			workflow_custom_roles.title_ar,
			workflow_custom_roles.type,
			workflow_custom_roles.module_identifier
		`)

	if len(params.QuickFilter) > 0 {
		query = query.Where("workflow_custom_roles.type IN ?", params.QuickFilter)
	}

	result := query.Find(&roles)
	return roles, result.Error
}

func (repo *WorkflowCustomRoleRepository) GetByType(roleType string) (*models.WorkflowCustomRole, error) {
	var role models.WorkflowCustomRole

	result := repo.db.
		Table("workflow_custom_roles").
		Where("type = ?", roleType).
		First(&role)

	if result.Error != nil {
		return nil, result.Error
	}

	return &role, nil
}
