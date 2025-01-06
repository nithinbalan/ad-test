package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	incidentDTOs "jethur/grc/incident/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type IncidentCommunicationLogRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(incidentID uint, params *incidentDTOs.IncidentCommunicationLogListQueryParamsDTO) (*incidentDTOs.IncidentCommunicationLogListResponseDTO, error)
	GetByID(id uint) (*models.IncidentCommunicationLog, error)
	GetFullDetailsByID(id uint) (*models.IncidentCommunicationLog, error)
	Create(incidentCommunicationLog *models.IncidentCommunicationLog, tx ...utility.DatabaseTransaction) (uint, error)
	Update(incidentCommunicationLog *models.IncidentCommunicationLog, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	GetCountByRiskID(incidentID uint, tx ...utility.DatabaseTransaction) (int64, error)
}

type IncidentCommunicationLogRepository struct {
	db *gorm.DB
}

func NewIncidentCommunicationLogRepository(db *gorm.DB) *IncidentCommunicationLogRepository {
	return &IncidentCommunicationLogRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *IncidentCommunicationLogRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *IncidentCommunicationLogRepository) GetAll(incidentID uint, params *incidentDTOs.IncidentCommunicationLogListQueryParamsDTO) (*incidentDTOs.IncidentCommunicationLogListResponseDTO, error) {

	var incidentCommunicationLogs []*incidentDTOs.IncidentCommunicationLogListDTO
	var totalCount int64

	query := repo.db.Table("incident_communication_logs").
		Select(`
        incident_communication_logs.*,
        created_user.id AS created_user_id,
        created_user.first_name_en AS created_user_first_name_en,
        created_user.first_name_ar AS created_user_first_name_ar,
        created_user.last_name_en AS created_user_last_name_en,
        created_user.last_name_ar AS created_user_last_name_ar,   
        updated_user.id AS updated_user_id,
    	updated_user.first_name_en AS updated_user_first_name_en,
        updated_user.first_name_ar AS updated_user_first_name_ar,
        updated_user.last_name_en AS updated_user_last_name_en,
        updated_user.last_name_ar AS updated_user_last_name_ar 
    `).
		Joins(`
        LEFT JOIN users AS created_user ON created_user.id = incident_communication_logs.created_by_id
        LEFT JOIN users AS updated_user ON updated_user.id = incident_communication_logs.updated_by_id
    `).Where(`incident_communication_logs.incident_register_id`, incidentID)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply target_date filters
	if params.CommunicatedAt != "" {
		query = query.Where("communicated_at = ?", params.CommunicatedAt)
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
            (incident_communication_logs.communicated_by LIKE ? OR
            incident_communication_logs.communicated_to LIKE ? OR
            incident_communication_logs.communication_method LIKE ?)`,
			keyword, keyword, keyword)
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize)

	err := query.Find(&incidentCommunicationLogs).Error
	if err != nil {
		return nil, err
	}

	totalPages := 1
	totalPages = int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &incidentDTOs.IncidentCommunicationLogListResponseDTO{
		Items: incidentCommunicationLogs,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *IncidentCommunicationLogRepository) GetByID(id uint) (*models.IncidentCommunicationLog, error) {
	var incidentCommunicationLog *models.IncidentCommunicationLog
	result := repo.db.First(&incidentCommunicationLog, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return incidentCommunicationLog, nil
}

func (repo *IncidentCommunicationLogRepository) GetFullDetailsByID(id uint) (*models.IncidentCommunicationLog, error) {
	var incidentCommunicationLog *models.IncidentCommunicationLog

	result := repo.db.
		Preload("IncidentRegister").
		First(&incidentCommunicationLog, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return incidentCommunicationLog, nil
}

func (repo *IncidentCommunicationLogRepository) Create(incidentCommunicationLog *models.IncidentCommunicationLog, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(incidentCommunicationLog).Error; err != nil {
		return 0, err
	}
	return incidentCommunicationLog.ID, nil
}

func (repo *IncidentCommunicationLogRepository) Update(incidentCommunicationLog *models.IncidentCommunicationLog, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(incidentCommunicationLog).Error; err != nil {
		return err
	}
	return nil
}

func (repo *IncidentCommunicationLogRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var incidentCommunicationLog *models.IncidentCommunicationLog
	if err := db.First(&incidentCommunicationLog, id).Error; err != nil {
		return err
	}

	result := db.Delete(&models.IncidentCommunicationLog{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *IncidentCommunicationLogRepository) GetCountByRiskID(incidentID uint, tx ...utility.DatabaseTransaction) (int64, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var count int64
	if err := db.Model(&models.IncidentCommunicationLog{}).
		Where("incident_register_id = ?", incidentID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
