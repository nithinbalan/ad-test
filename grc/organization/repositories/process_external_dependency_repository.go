package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	orgDTOs "jethur/grc/organization/dtos"

	"gorm.io/gorm"
)

type ProcessExternalDependencyRepositoryInterface interface {
	GetAll(processID uint, params *orgDTOs.ProcessExternalDependencyListQueryParamsDTO) (*orgDTOs.ProcessExternalDependencyListResponseDTO, error)
	GetByID(id uint) (*models.ProcessExternalDependency, error)
	GetFullDetailsByID(id uint) (*models.ProcessExternalDependency, error)
	GetVendorRTOBySLARegisterID(slaRegisterID uint) (*int, *uint, error)
	Create(item *models.ProcessExternalDependency) (uint, error)
	Update(item *models.ProcessExternalDependency) error
	Delete(id uint) error
}

type ProcessExternalDependencyRepository struct {
	db *gorm.DB
}

func NewProcessExternalDependencyRepository(db *gorm.DB) *ProcessExternalDependencyRepository {
	return &ProcessExternalDependencyRepository{db: db}
}

func (repo *ProcessExternalDependencyRepository) GetAll(processID uint, params *orgDTOs.ProcessExternalDependencyListQueryParamsDTO) (*orgDTOs.ProcessExternalDependencyListResponseDTO, error) {
	var result []*orgDTOs.ProcessExternalDependencyListDTO
	var totalCount int64

	query := repo.db.Table("process_external_dependencies").
		Select(`
			process_external_dependencies.id,
			process_external_dependencies.process_id,
			process_external_dependencies.entity_name,
			process_external_dependencies.entity_type_id,
			entity_type.title_en AS entity_type_title_en,
			entity_type.title_ar AS entity_type_title_ar,
			process_external_dependencies.dependency_tier_id,
			process_dependency_tiers.title_en AS dependency_tier_title_en,
			process_dependency_tiers.title_ar AS dependency_tier_title_ar,
			process_external_dependencies.needed_support_description,
			process_external_dependencies.vendor_rto,
			process_external_dependencies.vendor_rto_unit_id,
			vendor_rto_unit.title_en AS vendor_rto_unit_title_en,
			vendor_rto_unit.title_ar AS vendor_rto_unit_title_ar,
			process_external_dependencies.contact_number,
			process_external_dependencies.contact_email,
			process_external_dependencies.created_by_id,
			created_user.id AS created_user_id,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.email AS created_user_email,
			created_user.profile_image_path AS created_user_profile_image_path,
			process_external_dependencies.updated_by_id,
			updated_user.id AS updated_user_id,
			updated_user.first_name_en AS updated_user_first_name_en,
			updated_user.first_name_ar AS updated_user_first_name_ar,
			updated_user.last_name_en AS updated_user_last_name_en,
			updated_user.last_name_ar AS updated_user_last_name_ar,
			updated_user.initial AS updated_user_initial,
			updated_user.color AS updated_user_color,
			updated_user.email AS updated_user_email,
			updated_user.profile_image_path AS updated_user_profile_image_path,
			CONVERT(VARCHAR, process_external_dependencies.created_at, 23) AS created_at,
			CONVERT(VARCHAR, process_external_dependencies.updated_at, 23) AS updated_at
		`).
		Joins(`LEFT JOIN process_external_dependency_entity_types AS entity_type ON entity_type.id = process_external_dependencies.entity_type_id`).
		Joins(`LEFT JOIN process_dependency_tiers ON process_dependency_tiers.id = process_external_dependencies.dependency_tier_id`).
		Joins(`LEFT JOIN business_continuity_duration_units AS vendor_rto_unit ON vendor_rto_unit.id = process_external_dependencies.vendor_rto_unit_id`).
		Joins(`LEFT JOIN users AS created_user ON created_user.id = process_external_dependencies.created_by_id`).
		Joins(`LEFT JOIN users AS updated_user ON updated_user.id = process_external_dependencies.updated_by_id`).
		Where("process_external_dependencies.process_id = ?", processID)

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			process_external_dependencies.entity_name LIKE ? OR
			entity_type.title_en LIKE ? OR
			process_dependency_tiers.title_en LIKE ? OR
			process_external_dependencies.needed_support_description LIKE ?
		`, keyword, keyword, keyword, keyword)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	if err := query.Find(&result).Error; err != nil {
		return nil, err
	}

	totalPages := 1
	if !params.Export {
		totalPages = int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	}

	return &orgDTOs.ProcessExternalDependencyListResponseDTO{
		Items: result,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}, nil
}

func (repo *ProcessExternalDependencyRepository) GetByID(id uint) (*models.ProcessExternalDependency, error) {
	var item *models.ProcessExternalDependency
	err := repo.db.First(&item, id).Error
	return item, err
}

func (repo *ProcessExternalDependencyRepository) GetFullDetailsByID(id uint) (*models.ProcessExternalDependency, error) {
	var item *models.ProcessExternalDependency

	result := repo.db.
		Preload("Process").
		Preload("EntityType").
		Preload("DependencyTier").
		Preload("VendorRTOUnit").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&item, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return item, nil
}

func (repo *ProcessExternalDependencyRepository) GetVendorRTOBySLARegisterID(slaRegisterID uint) (*int, *uint, error) {
	var slaRegister models.SLARegister

	result := repo.db.
		Table("sla_register").
		Select("rto_value", "rto_unit_id").
		Where("id = ?", slaRegisterID).
		First(&slaRegister)

	if result.Error != nil {
		return nil, nil, result.Error
	}

	return slaRegister.RTOValue, slaRegister.RTOUnitID, nil
}

func (repo *ProcessExternalDependencyRepository) Create(item *models.ProcessExternalDependency) (uint, error) {
	err := repo.db.Create(item).Error
	return item.ID, err
}

func (repo *ProcessExternalDependencyRepository) Update(item *models.ProcessExternalDependency) error {
	return repo.db.Save(item).Error
}

func (repo *ProcessExternalDependencyRepository) Delete(id uint) error {
	return repo.db.Delete(&models.ProcessExternalDependency{}, id).Error
}
