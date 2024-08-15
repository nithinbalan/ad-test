package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	incidentDTOs "jethur/grc/incident/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type IncidentInvestigationRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(incidentID uint, params *incidentDTOs.IncidentInvestigationListQueryParamsDTO) (*incidentDTOs.IncidentInvestigationListResponseDTO, error)
	GetByID(incidentID uint) (*models.IncidentInvestigation, error)
	GetFullDetailsByID(incidentID uint, tx ...utility.DatabaseTransaction) (*models.IncidentInvestigation, error)
	Create(incidentInvestigation *models.IncidentInvestigation, tx ...utility.DatabaseTransaction) (uint, error)
	Update(incidentInvestigation *models.IncidentInvestigation, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	CreateFindings(findings []*models.IncidentInvestigationFinding, tx ...utility.DatabaseTransaction) error
	CreateMembers(users []*models.IncidentInvestigationUser, tx ...utility.DatabaseTransaction) error
	DeleteMembers(investigationUserIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteFindings(investigationFindingIds []*uint, tx ...utility.DatabaseTransaction) error

	//Document
	CreateDocuments(docs []*models.IncidentInvestigationDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.IncidentInvestigationDocument, error)
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error
}

type IncidentInvestigationRepository struct {
	db *gorm.DB
}

func NewIncidentInvestigationRepository(db *gorm.DB) *IncidentInvestigationRepository {
	return &IncidentInvestigationRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *IncidentInvestigationRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *IncidentInvestigationRepository) GetAll(incidentID uint, params *incidentDTOs.IncidentInvestigationListQueryParamsDTO) (*incidentDTOs.IncidentInvestigationListResponseDTO, error) {
	var incidentInvestigations []*incidentDTOs.IncidentInvestigationListDTO
	var totalCount int64

	query := repo.db.Table("incident_investigations").
		Select(`
        incident_investigations.*,
		CONVERT(VARCHAR, incident_investigations.initiated_at, 23) AS initiated_at,
        incident_register.title as incident_title,
		team_lead.first_name_en AS team_lead_first_name_en,
		team_lead.first_name_ar AS team_lead_first_name_ar,
		team_lead.last_name_en AS team_lead_last_name_en,
		team_lead.last_name_ar AS team_lead_last_name_ar,
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
        LEFT JOIN incident_register ON incident_register.id = incident_investigations.incident_register_id
		LEFT JOIN users AS team_lead ON team_lead.id = incident_investigations.team_lead_id
		LEFT JOIN users AS created_user ON created_user.id = incident_investigations.created_by_id
		LEFT JOIN users AS updated_user ON updated_user.id = incident_investigations.updated_by_id
    `).Where(`incident_investigations.incident_register_id`, incidentID)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply other filters as needed (e.g., date filters)
	if params.InitiatedAt != "" {
		query = query.Where("initiated_at >= ?", params.InitiatedAt)
	}
	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(incident_investigations.description LIKE ? OR
			incident_register.title LIKE ? OR
			team_lead.first_name_en LIKE ? OR
			team_lead.first_name_ar LIKE ? OR
			team_lead.last_name_en LIKE ? OR
			team_lead.last_name_ar LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.last_name_ar LIKE ? OR
			updated_user.first_name_en LIKE ? OR
			updated_user.first_name_ar LIKE ? OR
			updated_user.last_name_en LIKE ? OR
			updated_user.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword)
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "title" // Default sort by title if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&incidentInvestigations).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &incidentDTOs.IncidentInvestigationListResponseDTO{
		Items: incidentInvestigations,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *IncidentInvestigationRepository) GetByID(incidentID uint) (*models.IncidentInvestigation, error) {
	var incidentInvestigation *models.IncidentInvestigation

	result := repo.db.Where("incident_register_id = ?", incidentID).First(&incidentInvestigation)
	if result.Error != nil {
		return nil, result.Error
	}

	return incidentInvestigation, nil

}

func (repo *IncidentInvestigationRepository) GetFullDetailsByID(incidentID uint, tx ...utility.DatabaseTransaction) (*models.IncidentInvestigation, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var incidentInvestigation *models.IncidentInvestigation

	result := db.Where("incident_register_id = ?", incidentID).
		Preload("IncidentRegister").
		Preload("TeamLead").
		Preload("Documents").
		Preload("Users.User").
		Preload("Findings").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&incidentInvestigation)

	if result.Error != nil {
		return nil, result.Error
	}

	return incidentInvestigation, nil
}

func (repo *IncidentInvestigationRepository) Create(incidentInvestigation *models.IncidentInvestigation, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(incidentInvestigation).Error; err != nil {
		return 0, err
	}

	return incidentInvestigation.ID, nil
}

func (repo *IncidentInvestigationRepository) Update(incidentInvestigation *models.IncidentInvestigation, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(incidentInvestigation).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IncidentInvestigationRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the Incident register
	if err := db.Delete(&models.IncidentInvestigation{}, id).Error; err != nil {
		return err
	}

	return nil
}

// Document func
func (repo *IncidentInvestigationRepository) GetDocumentByID(id uint) (*models.IncidentInvestigationDocument, error) {
	var document models.IncidentInvestigationDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *IncidentInvestigationRepository) CreateDocuments(docs []*models.IncidentInvestigationDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Save all documents in a single transaction
	if err := db.Create(&docs).Error; err != nil {
		return err
	}
	return nil
}

func (repo *IncidentInvestigationRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete documents by IDs
	if err := db.Where("id IN (?)", docIds).Delete(&models.IncidentInvestigationDocument{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IncidentInvestigationRepository) CreateFindings(findings []*models.IncidentInvestigationFinding, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&findings).Error; err != nil {
		return err
	}
	return nil
}

func (repo *IncidentInvestigationRepository) CreateMembers(users []*models.IncidentInvestigationUser, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&users).Error; err != nil {
		return err
	}
	return nil
}

func (repo *IncidentInvestigationRepository) DeleteMembers(investigationUserIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete users by IDs
	if err := db.Where("id IN (?)", investigationUserIds).Delete(&models.IncidentInvestigationUser{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IncidentInvestigationRepository) DeleteFindings(investigationFindingIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete compliance requirements by IDs
	if err := db.Where("id IN (?)", investigationFindingIds).Delete(&models.IncidentInvestigationFinding{}).Error; err != nil {
		return err
	}

	return nil
}
