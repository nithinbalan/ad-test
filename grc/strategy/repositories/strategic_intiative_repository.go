package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/grc/strategy/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type StrategicIntiativeRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction

	GetAll(params *dtos.StrategicIntiativeListQueryParamsDTO) (*dtos.StrategicIntiativeListResponseDTO, error)
	GetByID(id uint) (*models.StrategicIntiative, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.StrategicIntiative, error)
	Create(strategicIntiative *models.StrategicIntiative, tx ...utility.DatabaseTransaction) (uint, error)
	Update(strategicIntiative *models.StrategicIntiative, tx ...utility.DatabaseTransaction) error
	Delete(id uint) error
}

// implementation
type StrategicIntiativeRepository struct {
	db *gorm.DB
}

func NewStrategicIntiativeRepository(db *gorm.DB) *StrategicIntiativeRepository {
	return &StrategicIntiativeRepository{
		db: db,
	}
}
func (repo *StrategicIntiativeRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *StrategicIntiativeRepository) GetAll(params *dtos.StrategicIntiativeListQueryParamsDTO) (*dtos.StrategicIntiativeListResponseDTO, error) {
	var intiatives []*dtos.StrategicIntiativeListDTO
	var totalCount int64

	query := repo.db.Table("strategic_intiatives").
		Select(`
			strategic_intiatives.*,

			responsible_user.first_name_en AS responsible_user_first_name_en,
			responsible_user.last_name_en AS responsible_user_last_name_en,
			responsible_user.first_name_ar AS responsible_user_first_name_ar,
			responsible_user.last_name_ar AS responsible_user_last_name_ar,
			responsible_user.initial AS responsible_user_initial,
			responsible_user.email AS responsible_user_email,
			responsible_user.color AS responsible_user_color,
			responsible_user.profile_image_path AS responsible_user_profile_image_path,

			responsible_unit.title_en AS responsible_unit_title_en,
			responsible_unit.title_ar AS responsible_unit_title_ar,

			owner_department.title_en AS owner_department_title_en,
			owner_department.title_ar AS owner_department_title_ar,

			strategic_initiative_statuses.title_en AS strategic_initiative_status_title_en,
			strategic_initiative_statuses.title_ar AS strategic_initiative_status_title_ar,
			strategic_initiative_statuses.color_code AS strategic_initiative_status_color_code,
			strategic_initiative_statuses.type AS strategic_initiative_status_type,

			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.email AS created_user_email,
			created_user.profile_image_path AS created_user_profile_image_path,

			updated_user.first_name_en AS updated_user_first_name_en,
			updated_user.first_name_ar AS updated_user_first_name_ar,
			updated_user.last_name_en AS updated_user_last_name_en,
			updated_user.last_name_ar AS updated_user_last_name_ar,
			updated_user.initial AS updated_user_initial,
			updated_user.color AS updated_user_color,
			updated_user.email AS updated_user_email,
			updated_user.profile_image_path AS updated_user_profile_image_path
		`).
		Joins("LEFT JOIN users AS responsible_user ON responsible_user.id = strategic_intiatives.responsible_user_id").
		Joins("LEFT JOIN organization_units AS responsible_unit ON responsible_unit.id = strategic_intiatives.responsible_unit_id").
		Joins("LEFT JOIN organization_units AS owner_department ON owner_department.id = strategic_intiatives.owner_department_id").
		Joins("LEFT JOIN strategic_initiative_statuses ON strategic_initiative_statuses.id = strategic_intiatives.strategic_initiative_status_id").
		Joins(`LEFT JOIN users AS created_user ON created_user.id = strategic_intiatives.created_by_id`).
		Joins(`LEFT JOIN users AS updated_user ON updated_user.id = strategic_intiatives.updated_by_id`)

	if params.Keyword != "" {
		kw := "%" + params.Keyword + "%"
		query = query.Where(`
			strategic_intiatives.title_en LIKE ? OR
			strategic_intiatives.title_ar LIKE ?
		`, kw, kw)
	}

	for col, val := range params.Filters {
		if len(val) > 0 {
			query = query.Where(col+" IN (?)", val)
		}
	}

	if params.StartDate != "" && params.EndDate != "" {
		query = query.Where("strategic_intiatives.created_at BETWEEN ? AND ?", params.StartDate+" 00:00:00", params.EndDate+" 23:59:59")
	}

	// Sorting
	query = query.Order(params.SortBy + " " + params.Sort)

	// Pagination
	query.Count(&totalCount)
	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize)

	if err := query.Find(&intiatives).Error; err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	response := &dtos.StrategicIntiativeListResponseDTO{
		Items: intiatives,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *StrategicIntiativeRepository) GetByID(id uint) (*models.StrategicIntiative, error) {
	var strategicIntiative *models.StrategicIntiative
	result := repo.db.First(&strategicIntiative, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return strategicIntiative, nil
}

func (repo *StrategicIntiativeRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.StrategicIntiative, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var strategicIntiative *models.StrategicIntiative

	result := db.
		Preload("ResponsibleUser").
		Preload("ResponsibleUnit").
		Preload("OwnerDepartment").
		Preload("StrategicInitiativeStatus").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&strategicIntiative, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return strategicIntiative, nil
}

func (repo *StrategicIntiativeRepository) Create(strategicIntiative *models.StrategicIntiative, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(strategicIntiative).Error; err != nil {
		return 0, err
	}

	return strategicIntiative.ID, nil
}

func (repo *StrategicIntiativeRepository) Update(strategicIntiative *models.StrategicIntiative, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Save(strategicIntiative).Error; err != nil {
		return err
	}

	return nil
}

func (repo *StrategicIntiativeRepository) Delete(id uint) error {
	var strategicIntiative *models.StrategicIntiative
	if err := repo.db.First(&strategicIntiative, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.StrategicIntiative{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
