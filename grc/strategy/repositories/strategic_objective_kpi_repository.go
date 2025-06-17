package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/grc/strategy/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type StrategicObjectiveKPIRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAllStrategicObjectiveKPI(strategicObjectiveID uint, params *dtos.StrategicObjectiveKPIListQueryParamsDTO) (*dtos.StrategicObjectiveKPIListResponseDTO, error)
	CreateStrategicObjectiveKPIs(mappings []models.StrategicObjectiveKPI, tx utility.DatabaseTransaction) error
	DeleteSingleStrategicObjectiveKPI(strategicObjectiveID uint, kpiRegisterID uint, tx utility.DatabaseTransaction) error
}

// implementation
type StrategicObjectiveKPIRepository struct {
	db *gorm.DB
}

func NewStrategicObjectiveKPIRepository(db *gorm.DB) *StrategicObjectiveKPIRepository {
	return &StrategicObjectiveKPIRepository{
		db: db,
	}
}
func (repo *StrategicObjectiveKPIRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *StrategicObjectiveKPIRepository) GetAllStrategicObjectiveKPI(
	strategicObjectiveID uint,
	params *dtos.StrategicObjectiveKPIListQueryParamsDTO,
) (*dtos.StrategicObjectiveKPIListResponseDTO, error) {

	var kpis []*dtos.StrategicObjectiveKPIListDTO
	var totalCount int64

	// ✅ Use valid base table & join chain
	query := repo.db.Table("strategic_objective_kpis").
		Select(`
			kpi_register.*,
			kpi_categories.title_en AS kpi_category_title_en,
			kpi_categories.title_ar AS kpi_category_title_ar,
			responsible_user.id AS responsible_user_id,
			responsible_user.first_name_en AS responsible_user_first_name_en,
			responsible_user.first_name_ar AS responsible_user_first_name_ar,
			responsible_user.last_name_en AS responsible_user_last_name_en,
			responsible_user.last_name_ar AS responsible_user_last_name_ar,
			responsible_user.initial AS responsible_user_initial,
			responsible_user.color AS responsible_user_color,
			responsible_user.email AS responsible_user_email,
			responsible_user.profile_image_path AS responsible_user_profile_image_path,
			kpi_statuses.title_en AS kpi_status_title_en,
			kpi_statuses.title_ar AS kpi_status_title_ar,
			kpi_statuses.type AS kpi_status_type,
			kpi_statuses.color_code AS kpi_status_color_code,
			created_user.id AS created_user_id,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.email AS created_user_email,
			created_user.profile_image_path AS created_user_profile_image_path,
			updated_user.id AS updated_user_id,
			updated_user.first_name_en AS updated_user_first_name_en,
			updated_user.first_name_ar AS updated_user_first_name_ar,
			updated_user.last_name_en AS updated_user_last_name_en,
			updated_user.last_name_ar AS updated_user_last_name_ar,
			updated_user.initial AS updated_user_initial,
			updated_user.color AS updated_user_color,
			updated_user.email AS updated_user_email,
			updated_user.profile_image_path AS updated_user_profile_image_path
		`).
		Joins(`INNER JOIN kpi_register ON kpi_register.id = strategic_objective_kpis.kpi_register_id`).
		Joins(`LEFT JOIN kpi_categories ON kpi_categories.id = kpi_register.kpi_category_id`).
		Joins(`LEFT JOIN kpi_statuses ON kpi_statuses.id = kpi_register.kpi_status_id`).
		Joins(`LEFT JOIN users AS responsible_user ON responsible_user.id = kpi_register.responsible_user_id`).
		Joins(`LEFT JOIN users AS created_user ON created_user.id = kpi_register.created_by_id`).
		Joins(`LEFT JOIN users AS updated_user ON updated_user.id = kpi_register.updated_by_id`).
		Where("strategic_objective_kpis.strategic_objective_id = ?", strategicObjectiveID)

	// ✅ Other filters (excluding soID)
	for columnName, values := range params.Filters {
		if columnName == "strategic_objective_id" {
			continue
		}
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// ✅ Keyword search
	if params.Keyword != "" {
		kw := "%" + params.Keyword + "%"
		query = query.Where(`
			kpi_register.title LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.last_name_en LIKE ? OR
			updated_user.first_name_en LIKE ? OR
			updated_user.last_name_en LIKE ?`,
			kw, kw, kw, kw, kw)
	}

	// ✅ Ensure safe base-column ordering!
	sortBy := "kpi_register.title" // fallback
	if params.SortBy != "" {
		// Only allow safe columns from kpi_register or real joined columns.
		sortBy = params.SortBy
	}
	query = query.Order(sortBy + " " + params.Sort)

	// ✅ Total count
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// ✅ Pagination
	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	if err := query.Find(&kpis).Error; err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	return &dtos.StrategicObjectiveKPIListResponseDTO{
		Items: kpis,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}, nil
}

func (repo *StrategicObjectiveKPIRepository) CreateStrategicObjectiveKPIs(mappings []models.StrategicObjectiveKPI, tx utility.DatabaseTransaction) error {
	db := tx.DB()
	return db.Create(&mappings).Error
}

func (repo *StrategicObjectiveKPIRepository) DeleteSingleStrategicObjectiveKPI(strategicObjectiveID uint, kpiRegisterID uint, tx utility.DatabaseTransaction) error {
	db := tx.DB()
	return db.
		Where("strategic_objective_id = ? AND kpi_register_id = ?", strategicObjectiveID, kpiRegisterID).
		Delete(&models.StrategicObjectiveKPI{}).Error
}
