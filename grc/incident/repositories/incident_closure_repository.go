package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	incidentDTOs "jethur/grc/incident/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type IncidentClosureRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(incidentID uint, params *incidentDTOs.IncidentClosureListQueryParamsDTO) (*incidentDTOs.IncidentClosureListResponseDTO, error)
	GetByID(id uint) (*models.IncidentClosure, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.IncidentClosure, error)
	Create(incidentClosure *models.IncidentClosure, tx ...utility.DatabaseTransaction) (uint, error)
	Update(incidentClosure *models.IncidentClosure, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	CreateLessons(lessons []*models.IncidentClosureLesson, tx ...utility.DatabaseTransaction) error
	DeleteLessons(lessonsIDs []*uint, tx ...utility.DatabaseTransaction) error
}

type IncidentClosureRepository struct {
	db *gorm.DB
}

func NewIncidentClosureRepository(db *gorm.DB) *IncidentClosureRepository {
	return &IncidentClosureRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *IncidentClosureRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *IncidentClosureRepository) GetAll(incidentID uint, params *incidentDTOs.IncidentClosureListQueryParamsDTO) (*incidentDTOs.IncidentClosureListResponseDTO, error) {

	var incidentClosures []*incidentDTOs.IncidentClosureListDTO
	var totalCount int64

	query := repo.db.Table("incident_closures").
		Select(`
        incident_closures.*,
        created_user.first_name_en AS created_user_first_name_en,
        created_user.first_name_ar AS created_user_first_name_ar,
        created_user.last_name_en AS created_user_last_name_en,
        created_user.last_name_ar AS created_user_last_name_ar,   
        updated_user.first_name_en AS updated_user_first_name_en,
        updated_user.first_name_ar AS updated_user_first_name_ar,
        updated_user.last_name_en AS updated_user_last_name_en,
        updated_user.last_name_ar AS updated_user_last_name_ar
    `).
		Joins(`
		LEFT JOIN users AS created_user ON created_user.id = incident_closures.created_by_id
        LEFT JOIN users AS updated_user ON updated_user.id = incident_closures.updated_by_id
    `)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
            (incident_closures.notes LIKE ?)`,
			keyword)
	}
	// Conditionally apply incident ID filter
	if incidentID != 0 {
		query = query.Where("incident_closures.incident_register_id = ?", incidentID)
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "notes" // Default sort by title if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&incidentClosures).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &incidentDTOs.IncidentClosureListResponseDTO{
		Items: incidentClosures,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *IncidentClosureRepository) GetByID(id uint) (*models.IncidentClosure, error) {
	var incidentClosure *models.IncidentClosure
	result := repo.db.First(&incidentClosure, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return incidentClosure, nil
}

func (repo *IncidentClosureRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.IncidentClosure, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var incidentClosure *models.IncidentClosure

	result := db.
		Preload("IncidentRegister").
		Preload("Lessons").
		First(&incidentClosure, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return incidentClosure, nil
}

func (repo *IncidentClosureRepository) Create(incidentClosure *models.IncidentClosure, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(incidentClosure).Error; err != nil {
		return 0, err
	}
	return incidentClosure.ID, nil
}

func (repo *IncidentClosureRepository) Update(incidentClosure *models.IncidentClosure, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(incidentClosure).Error; err != nil {
		return err
	}
	return nil
}

func (repo *IncidentClosureRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var incidentClosure *models.IncidentClosure
	if err := db.First(&incidentClosure, id).Error; err != nil {
		return err
	}

	result := db.Delete(&models.IncidentClosure{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *IncidentClosureRepository) CreateLessons(lessons []*models.IncidentClosureLesson, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&lessons).Error; err != nil {
		return err
	}
	return nil
}

func (repo *IncidentClosureRepository) DeleteLessons(lessonsIDs []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete compliance requirements by IDs
	if err := db.Where("id IN (?)", lessonsIDs).Delete(&models.IncidentClosureLesson{}).Error; err != nil {
		return err
	}

	return nil
}
