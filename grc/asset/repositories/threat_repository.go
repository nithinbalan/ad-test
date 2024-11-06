package repositories

import (
	"jethur/grc/asset/dtos"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type ThreatRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *dtos.ThreatListQueryParamsDTO) (*dtos.ThreatListResponseDTO, error)
	GetByID(id uint) (*models.Threat, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.Threat, error)
	Create(threat *models.Threat, tx ...utility.DatabaseTransaction) (uint, error)
	Update(threat *models.Threat, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error

	GetDocumentByID(id uint) (*models.ThreatDocument, error)
	CreateDocuments(docs []*models.ThreatDocument, tx ...utility.DatabaseTransaction) error
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error
}

// Implementation
type ThreatRepository struct {
	db *gorm.DB
}

func NewThreatRepository(db *gorm.DB) *ThreatRepository {
	return &ThreatRepository{
		db: db,
	}
}

func (repo *ThreatRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ThreatRepository) GetAll(params *dtos.ThreatListQueryParamsDTO) (*dtos.ThreatListResponseDTO, error) {
	var threats []*dtos.ThreatListDTO
	var totalCount int64

	query := repo.db.Table("threats").
		Select(`
			threats.id,
			threats.title_en,
			threats.title_ar,
			threat_types.title_en AS threat_type_title_en,
			threat_types.title_ar AS threat_type_title_ar,
			severity_levels.title_en AS severity_level_title_en,
			severity_levels.title_ar AS severity_level_title_ar,
			threat_statuses.title_en AS threat_status_title_en,
			threat_statuses.title_ar AS threat_status_title_ar,
			organization_units.title_en AS responsible_unit_title_en,
			organization_units.title_ar AS responsible_unit_title_ar,
			responsible_user.id AS responsible_user_id,
			responsible_user.first_name_en AS responsible_user_first_name_en,
			responsible_user.first_name_ar AS responsible_user_first_name_ar,
			responsible_user.last_name_en AS responsible_user_last_name_en,
			responsible_user.last_name_ar AS responsible_user_last_name_ar,
			responsible_user.initial AS responsible_user_initial,
			responsible_user.color AS responsible_user_color,
			responsible_user.email AS responsible_user_email,
			responsible_user.profile_image_path AS responsible_user_profile_image_path,
			threats.description,
			threats.source,
			CONVERT(VARCHAR, threats.detection_date, 23) AS detection_date,
			CONVERT(VARCHAR, threats.resolution_date, 23) AS resolution_date,
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
		Joins("LEFT JOIN threat_types ON threat_types.id = threats.threat_type_id").
		Joins("LEFT JOIN severity_levels ON severity_levels.id = threats.severity_level_id").
		Joins("LEFT JOIN threat_statuses ON threat_statuses.id = threats.threat_status_id").
		Joins("LEFT JOIN organization_units ON organization_units.id = threats.responsible_unit_id").
		Joins("LEFT JOIN users AS responsible_user ON responsible_user.id = threats.responsible_user_id").
		Joins("LEFT JOIN users AS created_user ON created_user.id = threats.created_by_id").
		Joins("LEFT JOIN users AS updated_user ON updated_user.id = threats.updated_by_id")

	// Apply filters
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where("threats."+columnName+" IN (?)", values)
		}
	}

	// Apply date filters
	if params.DetectionDateFrom != "" {
		query = query.Where("threats.detection_date >= ?", params.DetectionDateFrom)
	}
	if params.DetectionDateTo != "" {
		query = query.Where("threats.detection_date <= ?", params.DetectionDateTo)
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		query = query.Where(`
			(threats.title_en LIKE ? OR
			threats.title_ar LIKE ? OR
			threat_types.title_en LIKE ? OR
			threat_types.title_ar LIKE ? OR
			severity_levels.title_en LIKE ? OR
			severity_levels.title_ar LIKE ? OR
			threat_statuses.title_en LIKE ? OR
			threat_statuses.title_ar LIKE ? OR
			organization_units.title_en LIKE ? OR
			organization_units.title_ar LIKE ? OR
			responsible_user.first_name LIKE ? OR
			responsible_user.last_name LIKE ?)`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply sorting and pagination
	query = query.Order(params.SortBy + " " + params.Sort)
	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	// Execute query
	err := query.Find(&threats).Error
	if err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := 1
	if !params.Export {
		totalPages = int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	}

	// Prepare response
	response := &dtos.ThreatListResponseDTO{
		Items: threats,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}

	return response, nil
}

func (repo *ThreatRepository) GetByID(id uint) (*models.Threat, error) {
	var threat *models.Threat
	result := repo.db.First(&threat, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return threat, nil
}

func (repo *ThreatRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.Threat, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var threat models.Threat

	result := db.
		Preload("ThreatType").
		Preload("SeverityLevel").
		Preload("ThreatStatus").
		Preload("ResponsibleUser").
		Preload("ResponsibleUnit").
		Preload("Documents").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&threat, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &threat, nil
}

func (repo *ThreatRepository) Create(threat *models.Threat, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(threat).Error; err != nil {
		return 0, err
	}
	return threat.ID, nil
}

func (repo *ThreatRepository) Update(threat *models.Threat, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(threat).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ThreatRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var threat *models.Threat
	if err := db.First(&threat, id).Error; err != nil {
		return err
	}
	return repo.db.Delete(&models.Threat{}, id).Error
}

func (repo *ThreatRepository) GetDocumentByID(id uint) (*models.ThreatDocument, error) {
	var document models.ThreatDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *ThreatRepository) CreateDocuments(docs []*models.ThreatDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&docs).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ThreatRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Where("id IN (?)", docIds).Delete(&models.ThreatDocument{}).Error; err != nil {
		return err
	}
	return nil
}
