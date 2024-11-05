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
	GetAll(params *dtos.ThreatListQueryParamsDTO) (*dtos.ThreatListResponseDTO, error)
	GetByID(id uint) (*models.Threat, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.Threat, error)
	Create(threat *models.Threat) (uint, error)
	Update(threat *models.Threat) error
	Delete(id uint) error
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
			responsible_user.first_name AS responsible_user_first_name,
			responsible_user.last_name AS responsible_user_last_name,
			threats.description_en,
			threats.description_ar,
			CONVERT(VARCHAR, threats.detection_date, 23) AS detection_date,
			CONVERT(VARCHAR, threats.resolution_date, 23) AS resolution_date,
			threats.created_by_id,
			threats.created_at,
			threats.updated_by_id,
			threats.updated_at
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

	var threat *models.Threat

	result := db.
		Preload("ThreatType").
		Preload("SeverityLevel").
		Preload("ThreatStatus").
		Preload("ResponsibleUser").
		Preload("ResponsibleUser.OrganizationUnit").
		Preload("ResponsibleUnit").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&threat, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return threat, nil
}

func (repo *ThreatRepository) Create(threat *models.Threat) (uint, error) {
	if err := repo.db.Create(threat).Error; err != nil {
		return 0, err
	}
	return threat.ID, nil
}

func (repo *ThreatRepository) Update(threat *models.Threat) error {
	if err := repo.db.Save(threat).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ThreatRepository) Delete(id uint) error {
	var threat *models.Threat
	if err := repo.db.First(&threat, id).Error; err != nil {
		return err
	}
	return repo.db.Delete(&models.Threat{}, id).Error
}


