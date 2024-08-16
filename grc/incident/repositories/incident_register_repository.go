package repositories

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	incidentDTOs "jethur/grc/incident/dtos"
	"jethur/grc/models"
	"jethur/utility"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type IncidentRegisterRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *incidentDTOs.IncidentRegisterListQueryParamsDTO) (*incidentDTOs.IncidentRegisterListResponseDTO, error)
	GetByID(id uint) (*models.IncidentRegister, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.IncidentRegister, error)
	Create(incidentRegister *models.IncidentRegister, tx ...utility.DatabaseTransaction) (uint, error)
	Update(incidentRegister *models.IncidentRegister, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	GenerateReferenceNumber() (string, error)
	CreateIncidentComplianceRequirements(complianceRequirements []*models.IncidentComplianceRequirement, tx ...utility.DatabaseTransaction) error
	CreateIncidentControls(controls []*models.IncidentControl, tx ...utility.DatabaseTransaction) error
	DeleteIncidentControls(controlIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteIncidentComplianceRequirements(requirementIds []*uint, tx ...utility.DatabaseTransaction) error
	CreateIncidentRisks(risks []*models.IncidentRisk, tx ...utility.DatabaseTransaction) error
	DeleteIncidentRisks(riskIds []*uint, tx ...utility.DatabaseTransaction) error
	CreateIncidentProcesses(processes []*models.IncidentProcess, tx ...utility.DatabaseTransaction) error
	DeleteIncidentProcesses(processIds []*uint, tx ...utility.DatabaseTransaction) error

	//Document
	CreateDocuments(docs []*models.IncidentDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.IncidentDocument, error)
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error

	CreateIncidentSLAs(slaRegisters []*models.IncidentSLA, tx ...utility.DatabaseTransaction) error
	DeleteIncidentSLAs(slaRegisterIds []*uint, tx ...utility.DatabaseTransaction) error
}

type IncidentRegisterRepository struct {
	db *gorm.DB
}

func NewIncidentRegisterRepository(db *gorm.DB) *IncidentRegisterRepository {
	return &IncidentRegisterRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *IncidentRegisterRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *IncidentRegisterRepository) GetAll(params *incidentDTOs.IncidentRegisterListQueryParamsDTO) (*incidentDTOs.IncidentRegisterListResponseDTO, error) {
	var incidentRegisters []*incidentDTOs.IncidentRegisterListDTO
	var totalCount int64

	query := repo.db.Table("incident_register").
		Select(`
            incident_register.*,
			CONVERT(VARCHAR, incident_register.incident_at, 23) AS incident_at,
			CONVERT(VARCHAR, incident_register.reported_at, 23) AS reported_at,
            incident_categories.title_en AS incident_category_title_en,
            incident_categories.title_ar AS incident_category_title_ar,
            incident_types.title_en AS incident_type_title_en,
            incident_types.title_ar AS incident_type_title_ar,
            incident_statuses.title_en AS incident_status_title_en,
            incident_statuses.title_ar AS incident_status_title_ar,
			severity_levels.title_en AS severity_level_title_en,
			severity_levels.title_ar AS severity_level_title_ar,
			severity_levels.color_code AS severity_level_color_code,
			organization_units.title_en AS responsible_unit_title_en,
			organization_units.title_ar AS responsible_unit_title_ar,
			responsible_user.first_name_en AS responsible_user_first_name_en,
            responsible_user.first_name_ar AS responsible_user_first_name_ar,
            responsible_user.last_name_en AS responsible_user_last_name_en,
            responsible_user.last_name_ar AS responsible_user_last_name_ar,
			responsible_user.initial AS responsible_user_initial,
            responsible_user.color AS responsible_user_color,
            responsible_user.email AS responsible_user_email,
            responsible_user.profile_image_path AS responsible_user_profile_image_path,
            reported_user.first_name_en AS reported_by_first_name_en,
            reported_user.first_name_ar AS reported_by_first_name_ar,
            reported_user.last_name_en AS reported_by_last_name_en,
            reported_user.last_name_ar AS reported_by_last_name_ar,
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
            LEFT JOIN incident_categories ON incident_categories.id = incident_register.incident_category_id
            LEFT JOIN incident_sub_categories ON incident_sub_categories.id = incident_register.incident_sub_category_id
            LEFT JOIN incident_types ON incident_types.id = incident_register.incident_type_id
            LEFT JOIN incident_statuses ON incident_statuses.id = incident_register.incident_status_id
		    LEFT JOIN severity_levels ON severity_levels.id = incident_register.severity_level_id
            LEFT JOIN organization_units ON organization_units.id = incident_register.responsible_unit_id
		    LEFT JOIN users AS responsible_user ON responsible_user.id = incident_register.responsible_user_id
			LEFT JOIN users AS reported_user ON reported_user.id = incident_register.reported_by
			LEFT JOIN users AS created_user ON created_user.id = incident_register.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = incident_register.updated_by_id
        `)
	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply other filters as needed (e.g., date filters)
	if params.ReportedAt != "" {
		query = query.Where("reported_at >= ?", params.ReportedAt)
	}
	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(incident_register.title LIKE ? OR
			incident_register.description LIKE ? OR
			incident_register.reference_number LIKE ? OR
			incident_categories.title_en LIKE ? OR
			incident_categories.title_ar LIKE ? OR
			incident_types.title_en LIKE ? OR
			incident_types.title_ar LIKE ? OR
			incident_statuses.title_en LIKE ? OR
			incident_statuses.title_ar LIKE ? OR
			severity_levels.title_en LIKE ? OR
			severity_levels.title_ar LIKE ? OR
			organization_units.title_en LIKE ? OR
			organization_units.title_ar LIKE ? OR
			responsible_user.first_name_en LIKE ? OR
			responsible_user.first_name_ar LIKE ? OR
			responsible_user.last_name_en LIKE ? OR
			responsible_user.last_name_ar LIKE ? OR
			responsible_user.email LIKE ? OR
            responsible_user.initial LIKE ? OR
            responsible_user.color LIKE ? OR
            responsible_user.profile_image_path LIKE ? OR
			reported_user.first_name_en LIKE ? OR
			reported_user.first_name_ar LIKE ? OR
			reported_user.last_name_en LIKE ? OR
			reported_user.last_name_ar LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.last_name_ar LIKE ? OR
			updated_user.id AS updated_user_id,
    		updated_user.first_name_en LIKE ? OR
			updated_user.first_name_ar LIKE ? OR
			updated_user.last_name_en LIKE ? OR
			updated_user.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
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
	err := query.Offset(offset).Limit(params.PageSize).Find(&incidentRegisters).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &incidentDTOs.IncidentRegisterListResponseDTO{
		Items: incidentRegisters,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *IncidentRegisterRepository) GetByID(id uint) (*models.IncidentRegister, error) {
	var incident *models.IncidentRegister

	result := repo.db.First(&incident, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return incident, nil

}

func (repo *IncidentRegisterRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.IncidentRegister, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var incident *models.IncidentRegister

	result := db.
		Preload("IncidentCategory").
		Preload("IncidentType").
		Preload("IncidentStatus").
		Preload("ReportedByUser").
		Preload("SeverityLevel").
		Preload("Documents").
		Preload("Controls").
		Preload("ComplianceRequirements").
		Preload("Risks").
		Preload("Processes").
		Preload("SLARegisters").
		Preload("ResponsibleUnit.ParentUnit.ParentUnit.ParentUnit.ParentUnit.ParentUnit").
		Preload("ResponsibleUser").
		Preload("ResponsibleUser.OrganizationUnit1").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&incident, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return incident, nil
}

func (repo *IncidentRegisterRepository) Create(incidentRegister *models.IncidentRegister, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(incidentRegister).Error; err != nil {
		return 0, err
	}

	return incidentRegister.ID, nil
}

func (repo *IncidentRegisterRepository) Update(incidentRegister *models.IncidentRegister, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(incidentRegister).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IncidentRegisterRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the Incident register
	if err := db.Delete(&models.IncidentRegister{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IncidentRegisterRepository) GenerateReferenceNumber() (string, error) {
	var lastReferenceNumber string
	err := repo.db.Model(&models.IncidentRegister{}).
		Select("reference_number").
		Order("reference_number desc").
		Limit(1).
		Pluck("reference_number", &lastReferenceNumber).
		Error

	if err != nil {
		return "", err
	}

	lastDigits := "0000"
	if lastReferenceNumber != "" && len(lastReferenceNumber) >= 4 {
		lastDigits = lastReferenceNumber[len(lastReferenceNumber)-4:]
	}

	currentYear := time.Now().Year()
	lastReferenceID, err := strconv.ParseUint(lastDigits, 10, 64)
	if err != nil {
		return "", err
	}
	newReferenceNumber := fmt.Sprintf("%s-%d-%04d", "IR", currentYear, lastReferenceID+1)

	return newReferenceNumber, nil

}

// Document func
func (repo *IncidentRegisterRepository) GetDocumentByID(id uint) (*models.IncidentDocument, error) {
	var document models.IncidentDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *IncidentRegisterRepository) CreateDocuments(docs []*models.IncidentDocument, tx ...utility.DatabaseTransaction) error {
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

func (repo *IncidentRegisterRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete documents by IDs
	if err := db.Where("id IN (?)", docIds).Delete(&models.IncidentDocument{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IncidentRegisterRepository) CreateIncidentComplianceRequirements(complianceRequirements []*models.IncidentComplianceRequirement, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&complianceRequirements).Error; err != nil {
		return err
	}
	return nil
}

func (repo *IncidentRegisterRepository) CreateIncidentControls(controls []*models.IncidentControl, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&controls).Error; err != nil {
		return err
	}
	return nil
}

func (repo *IncidentRegisterRepository) DeleteIncidentControls(controlIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete controls by IDs
	if err := db.Where("id IN (?)", controlIds).Delete(&models.IncidentControl{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IncidentRegisterRepository) DeleteIncidentComplianceRequirements(requirementIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete compliance requirements by IDs
	if err := db.Where("id IN (?)", requirementIds).Delete(&models.IncidentComplianceRequirement{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IncidentRegisterRepository) CreateIncidentRisks(risks []*models.IncidentRisk, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&risks).Error; err != nil {
		return err
	}
	return nil
}

func (repo *IncidentRegisterRepository) DeleteIncidentRisks(riskIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete risks by IDs
	if err := db.Where("id IN (?)", riskIds).Delete(&models.IncidentRisk{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IncidentRegisterRepository) CreateIncidentProcesses(processes []*models.IncidentProcess, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&processes).Error; err != nil {
		return err
	}
	return nil
}

func (repo *IncidentRegisterRepository) DeleteIncidentProcesses(processIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete processes by IDs
	if err := db.Where("id IN (?)", processIds).Delete(&models.IncidentProcess{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IncidentRegisterRepository) CreateIncidentSLAs(slaRegisters []*models.IncidentSLA, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&slaRegisters).Error; err != nil {
		return err
	}
	return nil
}

func (repo *IncidentRegisterRepository) DeleteIncidentSLAs(slaRegisterIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete SLARegister by IDs
	if err := db.Where("id IN (?)", slaRegisterIds).Delete(&models.IncidentSLA{}).Error; err != nil {
		return err
	}

	return nil
}
