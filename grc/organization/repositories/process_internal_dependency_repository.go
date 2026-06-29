package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	orgDTOs "jethur/grc/organization/dtos"

	"gorm.io/gorm"
)

type ProcessInternalDependencyRepositoryInterface interface {
	GetAll(processID uint, params *orgDTOs.ProcessInternalDependencyListQueryParamsDTO) (*orgDTOs.ProcessInternalDependencyListResponseDTO, error)
	GetByID(id uint) (*models.ProcessInternalDependency, error)
	GetFullDetailsByID(id uint) (*models.ProcessInternalDependency, error)
	Create(item *models.ProcessInternalDependency) (uint, error)
	Update(item *models.ProcessInternalDependency) error
	Delete(id uint) error
}

type ProcessInternalDependencyRepository struct {
	db *gorm.DB
}

func NewProcessInternalDependencyRepository(db *gorm.DB) *ProcessInternalDependencyRepository {
	return &ProcessInternalDependencyRepository{
		db: db,
	}
}

func (repo *ProcessInternalDependencyRepository) GetAll(processID uint, params *orgDTOs.ProcessInternalDependencyListQueryParamsDTO) (*orgDTOs.ProcessInternalDependencyListResponseDTO, error) {
	var result []*orgDTOs.ProcessInternalDependencyListDTO
	var totalCount int64

	query := repo.db.Table("process_internal_dependencies").
		Select(`
			process_internal_dependencies.id,
			process_internal_dependencies.process_id,
			process_internal_dependencies.supporting_unit_id,
			supporting_unit.title_en AS supporting_unit_title_en,
			supporting_unit.title_ar AS supporting_unit_title_ar,
			process_internal_dependencies.dependency_tier_id,
			process_dependency_tiers.title_en AS dependency_tier_title_en,
			process_dependency_tiers.title_ar AS dependency_tier_title_ar,
			process_internal_dependencies.dependency_percentage_from,
			process_internal_dependencies.dependency_percentage_to,
			process_internal_dependencies.required_support,
			process_internal_dependencies.notes,
			process_internal_dependencies.created_by_id,
			created_user.id AS created_user_id,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.email AS created_user_email,
			created_user.profile_image_path AS created_user_profile_image_path,
			process_internal_dependencies.updated_by_id,
			updated_user.id AS updated_user_id,
			updated_user.first_name_en AS updated_user_first_name_en,
			updated_user.first_name_ar AS updated_user_first_name_ar,
			updated_user.last_name_en AS updated_user_last_name_en,
			updated_user.last_name_ar AS updated_user_last_name_ar,
			updated_user.initial AS updated_user_initial,
			updated_user.color AS updated_user_color,
			updated_user.email AS updated_user_email,
			updated_user.profile_image_path AS updated_user_profile_image_path,
			CONVERT(VARCHAR, process_internal_dependencies.created_at, 23) AS created_at,
			CONVERT(VARCHAR, process_internal_dependencies.updated_at, 23) AS updated_at
		`).
		Joins(`LEFT JOIN organization_units AS supporting_unit ON supporting_unit.id = process_internal_dependencies.supporting_unit_id`).
		Joins(`LEFT JOIN process_dependency_tiers ON process_dependency_tiers.id = process_internal_dependencies.dependency_tier_id`).
		Joins(`LEFT JOIN users AS created_user ON created_user.id = process_internal_dependencies.created_by_id`).
		Joins(`LEFT JOIN users AS updated_user ON updated_user.id = process_internal_dependencies.updated_by_id`).
		Where("process_internal_dependencies.process_id = ?", processID)

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			supporting_unit.title_en LIKE ? OR
			supporting_unit.title_ar LIKE ? OR
			process_dependency_tiers.title_en LIKE ? OR
			process_internal_dependencies.required_support LIKE ? OR
			process_internal_dependencies.notes LIKE ?
		`, keyword, keyword, keyword, keyword, keyword)
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

	return &orgDTOs.ProcessInternalDependencyListResponseDTO{
		Items: result,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}, nil
}

func (repo *ProcessInternalDependencyRepository) GetByID(id uint) (*models.ProcessInternalDependency, error) {
	var item *models.ProcessInternalDependency
	err := repo.db.First(&item, id).Error
	return item, err
}

func (repo *ProcessInternalDependencyRepository) GetFullDetailsByID(id uint) (*models.ProcessInternalDependency, error) {
	var item *models.ProcessInternalDependency

	result := repo.db.
		Preload("SupportingUnit").
		Preload("DependencyTier").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&item, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return item, nil
}

func (repo *ProcessInternalDependencyRepository) Create(item *models.ProcessInternalDependency) (uint, error) {
	err := repo.db.Create(item).Error
	return item.ID, err
}

func (repo *ProcessInternalDependencyRepository) Update(item *models.ProcessInternalDependency) error {
	return repo.db.Save(item).Error
}

func (repo *ProcessInternalDependencyRepository) Delete(id uint) error {
	return repo.db.Delete(&models.ProcessInternalDependency{}, id).Error
}
