package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"

	"gorm.io/gorm"
)

type OrganizationUnitAttributeRepositoryInterface interface {
	GetAll(departmentID uint, params *organizationDTOs.OrganizationUnitAttributeListQueryParamsDTO) (*organizationDTOs.OrganizationUnitAttributeListResponseDTO, error)
	GetByID(id uint) (*models.OrganizationUnitAttribute, error)
	GetFullDetailsByID(id uint) (*models.OrganizationUnitAttribute, error)
	Create(attribute *models.OrganizationUnitAttribute) (uint, error)
	Update(id uint, attribute *models.OrganizationUnitAttribute) error
	Delete(ids []uint) error
}

type OrganizationUnitAttributeRepository struct {
	db *gorm.DB
}

func NewOrganizationUnitAttributeRepository(db *gorm.DB) *OrganizationUnitAttributeRepository {
	return &OrganizationUnitAttributeRepository{
		db: db,
	}
}

func (repo *OrganizationUnitAttributeRepository) GetAll(departmentID uint, params *organizationDTOs.OrganizationUnitAttributeListQueryParamsDTO) (*organizationDTOs.OrganizationUnitAttributeListResponseDTO, error) {
	var attributes []*organizationDTOs.OrganizationUnitAttributeListDTO
	var totalCount int64

	query := repo.db.Table("organization_unit_attributes AS oua").
		Select(`
			oua.id,
			oua.department_id,
			ou.title_en AS department_title_en,
			ou.title_ar AS department_title_ar,
			oua.title,
			oua.description,
			oua.created_at,
			oua.updated_at,
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
		Joins("LEFT JOIN organization_units AS ou ON oua.department_id = ou.id").
		Joins("LEFT JOIN users AS cu ON oua.created_by_id = cu.id").
		Joins("LEFT JOIN users AS uu ON oua.updated_by_id = uu.id").
		Where("oua.department_id = ?", departmentID)

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("(oua.title LIKE ? OR oua.description LIKE ?)", keyword, keyword)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	sortBy := params.SortBy
	if sortBy == "" {
		sortBy = "oua.created_at"
	}

	sort := params.Sort
	if sort != "asc" && sort != "desc" {
		sort = "desc"
	}

	query = query.Order(sortBy + " " + sort)

	if params.Page > 0 && params.PageSize > 0 {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	if err := query.Scan(&attributes).Error; err != nil {
		return nil, err
	}

	totalPages := 1
	if params.Page > 0 && params.PageSize > 0 {
		totalPages = int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	}

	response := &organizationDTOs.OrganizationUnitAttributeListResponseDTO{
		Items: attributes,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}

	return response, nil
}

func (repo *OrganizationUnitAttributeRepository) GetByID(id uint) (*models.OrganizationUnitAttribute, error) {
	var attribute models.OrganizationUnitAttribute
	err := repo.db.
		Where("id = ?", id).
		First(&attribute).Error

	if err != nil {
		return nil, err
	}

	return &attribute, nil
}

func (repo *OrganizationUnitAttributeRepository) GetFullDetailsByID(id uint) (*models.OrganizationUnitAttribute, error) {
	var attribute models.OrganizationUnitAttribute
	err := repo.db.
		Preload("Department").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Where("id = ?", id).
		First(&attribute).Error

	if err != nil {
		return nil, err
	}

	return &attribute, nil
}

func (repo *OrganizationUnitAttributeRepository) Create(attribute *models.OrganizationUnitAttribute) (uint, error) {
	result := repo.db.Create(&attribute)
	return attribute.ID, result.Error
}

func (repo *OrganizationUnitAttributeRepository) Update(id uint, attribute *models.OrganizationUnitAttribute) error {
	return repo.db.Model(&models.OrganizationUnitAttribute{}).
		Where("id = ?", id).
		Updates(attribute).Error
}

func (repo *OrganizationUnitAttributeRepository) Delete(ids []uint) error {
	return repo.db.Where("id IN ?", ids).Delete(&models.OrganizationUnitAttribute{}).Error
}
