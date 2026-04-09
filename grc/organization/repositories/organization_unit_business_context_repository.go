package repositories

import (
	"errors"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"

	"gorm.io/gorm"
)

// Interface
type OrganizationUnitBusinessContextRepositoryInterface interface {
	GetAll(departmentID uint, params *organizationDTOs.OrganizationUnitBusinessContextListQueryParamsDTO) (*organizationDTOs.OrganizationUnitBusinessContextListResponseDTO, error)
	GetByID(departmentID uint) (*models.OrganizationUnitBusinessContext, error)
	GetFullDetailsByID(departmentID uint) (*models.OrganizationUnitBusinessContext, error)
	Create(context *models.OrganizationUnitBusinessContext) (uint, error)
	Update(departmentID uint, updatedContext *models.OrganizationUnitBusinessContext) error
	Delete(departmentID uint) error
	ExistsByDepartmentID(departmentID uint) (bool, error)
	ExistsByDepartmentIDExcludingID(id uint, departmentID uint) (bool, error)
}

// Implementation
type OrganizationUnitBusinessContextRepository struct {
	db *gorm.DB
}

func NewOrganizationUnitBusinessContextRepository(db *gorm.DB) *OrganizationUnitBusinessContextRepository {
	return &OrganizationUnitBusinessContextRepository{
		db: db,
	}
}

func (repo *OrganizationUnitBusinessContextRepository) GetAll(departmentID uint, params *organizationDTOs.OrganizationUnitBusinessContextListQueryParamsDTO) (*organizationDTOs.OrganizationUnitBusinessContextListResponseDTO, error) {
	var context *organizationDTOs.OrganizationUnitBusinessContextListDTO

	query := repo.db.Table("organization_unit_business_contexts AS oubc").
		Select(`
			oubc.id,
			oubc.department_id,
			ou.title_en AS department_title_en,
			ou.title_ar AS department_title_ar,
			oubc.internal_audit_note,
			oubc.external_context,
			oubc.annual_objective_description,
			oubc.created_at,
			oubc.updated_at,
			cu.id AS created_user_id,
			cu.first_name_en AS created_user_first_name_en,
			cu.first_name_ar AS created_user_first_name_ar,
			cu.last_name_en AS created_user_last_name_en,
			cu.last_name_ar AS created_user_last_name_ar,
			cu.initial AS created_user_initial,
			cu.color AS created_user_color,
			cu.email AS created_user_email,
			cu.profile_image_path AS created_user_profile_image_path,
			uu.id AS updated_user_id,
			uu.first_name_en AS updated_user_first_name_en,
			uu.first_name_ar AS updated_user_first_name_ar,
			uu.last_name_en AS updated_user_last_name_en,
			uu.last_name_ar AS updated_user_last_name_ar,
			uu.initial AS updated_user_initial,
			uu.color AS updated_user_color,
			uu.email AS updated_user_email,
			uu.profile_image_path AS updated_user_profile_image_path
		`).
		Joins("LEFT JOIN organization_units AS ou ON oubc.department_id = ou.id").
		Joins("LEFT JOIN users AS cu ON oubc.created_by_id = cu.id").
		Joins("LEFT JOIN users AS uu ON oubc.updated_by_id = uu.id").
		Where("oubc.department_id = ?", departmentID)

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("(oubc.internal_audit_note LIKE ? OR oubc.external_context LIKE ? OR oubc.annual_objective_description LIKE ?)",
			keyword, keyword, keyword)
	}

	sortBy := params.SortBy
	if sortBy == "" {
		sortBy = "oubc.created_at"
	}

	sort := params.Sort
	if sort != "asc" && sort != "desc" {
		sort = "desc"
	}

	query = query.Order(sortBy + " " + sort)

	if err := query.First(&context).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Return empty response when no record found
			return &organizationDTOs.OrganizationUnitBusinessContextListResponseDTO{}, nil
		}
		return nil, err
	}

	response := &organizationDTOs.OrganizationUnitBusinessContextListResponseDTO{
		OrganizationUnitBusinessContextListDTO: context,
	}

	return response, nil
}

func (repo *OrganizationUnitBusinessContextRepository) GetByID(departmentID uint) (*models.OrganizationUnitBusinessContext, error) {
	var context models.OrganizationUnitBusinessContext
	err := repo.db.
		Where("department_id = ?", departmentID).
		First(&context).Error

	if err != nil {
		return nil, err
	}

	return &context, nil
}

func (repo *OrganizationUnitBusinessContextRepository) GetFullDetailsByID(departmentID uint) (*models.OrganizationUnitBusinessContext, error) {
	var context models.OrganizationUnitBusinessContext
	err := repo.db.
		Preload("Department").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Where("department_id = ?", departmentID).
		First(&context).Error

	if err != nil {
		return nil, err
	}

	return &context, nil
}

func (repo *OrganizationUnitBusinessContextRepository) Create(context *models.OrganizationUnitBusinessContext) (uint, error) {
	result := repo.db.Create(&context)
	return context.ID, result.Error
}

func (repo *OrganizationUnitBusinessContextRepository) Update(departmentID uint, updatedContext *models.OrganizationUnitBusinessContext) error {
	var context *models.OrganizationUnitBusinessContext

	if err := repo.db.Where("department_id = ?", departmentID).First(&context).Error; err != nil {
		return err
	}

	if err := repo.db.Model(&context).Updates(updatedContext).Error; err != nil {
		return err
	}

	return nil
}

func (repo *OrganizationUnitBusinessContextRepository) Delete(departmentID uint) error {
	var context *models.OrganizationUnitBusinessContext
	if err := repo.db.Where("department_id = ?", departmentID).First(&context).Error; err != nil {
		return err
	}

	return repo.db.Delete(&context).Error
}

func (repo *OrganizationUnitBusinessContextRepository) ExistsByDepartmentID(departmentID uint) (bool, error) {
	var count int64
	err := repo.db.Model(&models.OrganizationUnitBusinessContext{}).
		Where("department_id = ?", departmentID).
		Count(&count).Error

	return count > 0, err
}

func (repo *OrganizationUnitBusinessContextRepository) ExistsByDepartmentIDExcludingID(id uint, departmentID uint) (bool, error) {
	var count int64
	err := repo.db.Model(&models.OrganizationUnitBusinessContext{}).
		Where("department_id = ? AND id != ?", departmentID, id).
		Count(&count).Error

	return count > 0, err
}
