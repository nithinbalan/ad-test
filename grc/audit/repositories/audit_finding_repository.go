package repositories

import (
	"fmt"
	auditDTOs "jethur/grc/audit/dtos"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/utility"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type AuditFindingRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *auditDTOs.AuditFindingListQueryParamsDTO) (*auditDTOs.AuditFindingListResponseDTO, error)
	GetByID(id uint) (*models.AuditFinding, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditFinding, error)
	Create(auditFinding *models.AuditFinding, tx ...utility.DatabaseTransaction) (uint, error)
	Update(auditFinding *models.AuditFinding, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	GenerateReferenceNumber(tx ...utility.DatabaseTransaction) (string, error)

	CreateDocuments(docs []*models.AuditFindingDocument, tx ...utility.DatabaseTransaction) error
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.AuditFindingDocument, error)

	GetAuditFindingCountByCategory(auditID uint) ([]auditDTOs.AuditFindingCategoryCountResponseDTO, error)
}

type AuditFindingRepository struct {
	db *gorm.DB
}

func NewAuditFindingRepository(db *gorm.DB) *AuditFindingRepository {
	return &AuditFindingRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *AuditFindingRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AuditFindingRepository) GetAll(params *auditDTOs.AuditFindingListQueryParamsDTO) (*auditDTOs.AuditFindingListResponseDTO, error) {
	var auditFindings []*auditDTOs.AuditFindingListDTO
	var totalCount int64

	query := repo.db.Table("audit_findings").
		Select(`
			audit_findings.*,
			audits.id AS audit_id,
			audit_plans.title AS audit_title,
			auditable_items.title_en AS auditable_item_title_en,
			auditable_items.title_ar AS auditable_item_title_ar,
			audit_finding_categories.title_en AS audit_finding_category_title_en,
			audit_finding_categories.title_ar AS audit_finding_category_title_ar,
			severity_levels.title_en AS severity_level_title_en,
			severity_levels.title_ar AS severity_level_title_ar,
			severity_levels.color_code AS severity_level_color_code,
			audit_finding_statuses.title_en AS audit_finding_status_title_en,
			audit_finding_statuses.title_ar AS audit_finding_status_title_ar,
			audit_finding_statuses.color_code AS audit_finding_status_color_code,
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
		Joins(`
			LEFT JOIN audits ON audits.id = audit_findings.audit_id
			LEFT JOIN audit_plans ON audit_plans.id = audits.audit_plan_id
			LEFT JOIN auditable_items ON auditable_items.id = audit_findings.auditable_item_id
			LEFT JOIN audit_finding_categories ON audit_finding_categories.id = audit_findings.audit_finding_category_id
			LEFT JOIN severity_levels ON severity_levels.id = audit_findings.severity_level_id
			LEFT JOIN audit_finding_statuses ON audit_finding_statuses.id = audit_findings.audit_finding_status_id
			LEFT JOIN users AS created_user ON created_user.id = audit_findings.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = audit_findings.updated_by_id
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
			(auditable_items.title_en LIKE ? OR
			audit_finding_categories.title_en LIKE ? OR
			severity_levels.title_en LIKE ? OR
			audit_finding_statuses.title_en LIKE ? OR
			created_user.first_name_en LIKE ? OR
			updated_user.first_name_en LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword)
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "audit_findings.created_at" // Default sort by created_at if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&auditFindings).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &auditDTOs.AuditFindingListResponseDTO{
		Items: auditFindings,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *AuditFindingRepository) GetByID(id uint) (*models.AuditFinding, error) {
	var auditFinding *models.AuditFinding

	result := repo.db.Preload("Documents").First(&auditFinding, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return auditFinding, nil
}

func (repo *AuditFindingRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditFinding, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var auditFinding *models.AuditFinding

	result := db.
		Preload("AuditableItem").
		Preload("AuditableItem.Process").
		Preload("AuditableItem.Control").
		Preload("AuditableItem.Document").
		Preload("AuditableItem.IncidentRegister").
		Preload("AuditableItem.ComplianceRequirement").
		Preload("AuditableItem.RiskRegister").
		Preload("AuditableItem.User").
		Preload("AuditableItem.ResponsibleUnit").
		Preload("AuditableItem.ResponsibleUser").
		Preload("AuditFindingCategory").
		Preload("SeverityLevel").
		Preload("AuditFindingStatus").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("Documents").
		First(&auditFinding, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return auditFinding, nil
}

func (repo *AuditFindingRepository) Create(auditFinding *models.AuditFinding, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(auditFinding).Error; err != nil {
		return 0, err
	}

	return auditFinding.ID, nil
}

func (repo *AuditFindingRepository) Update(auditFinding *models.AuditFinding, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(auditFinding).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditFindingRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the AuditFinding
	if err := db.Delete(&models.AuditFinding{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditFindingRepository) GenerateReferenceNumber(tx ...utility.DatabaseTransaction) (string, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var lastReferenceNumber string
	err := db.Model(&models.AuditFinding{}).
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

func (repo *AuditFindingRepository) CreateDocuments(docs []*models.AuditFindingDocument, tx ...utility.DatabaseTransaction) error {
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

func (repo *AuditFindingRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete documents by IDs
	if err := db.Where("id IN (?)", docIds).Delete(&models.AuditFindingDocument{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditFindingRepository) GetDocumentByID(id uint) (*models.AuditFindingDocument, error) {
	var document models.AuditFindingDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *AuditFindingRepository) GetAuditFindingCountByCategory(auditID uint) ([]auditDTOs.AuditFindingCategoryCountResponseDTO, error) {
	categoryCounts := []auditDTOs.AuditFindingCategoryCountResponseDTO{}

	query := `
        SELECT 
            audit_finding_categories.id, 
            audit_finding_categories.title_en, 
            audit_finding_categories.title_ar, 
            audit_finding_categories.color_code, 
            COUNT(audit_findings.id) AS count
        FROM 
            audit_findings
        JOIN 
            audit_finding_categories ON audit_finding_categories.id = audit_findings.audit_finding_category_id
        WHERE 
            audit_findings.audit_id = ?
        GROUP BY 
            audit_finding_categories.id, 
            audit_finding_categories.title_en, 
            audit_finding_categories.title_ar,
			audit_finding_categories.color_code
    `

	err := repo.db.Raw(query, auditID).Scan(&categoryCounts).Error
	if err != nil {
		return nil, err
	}

	return categoryCounts, nil
}
