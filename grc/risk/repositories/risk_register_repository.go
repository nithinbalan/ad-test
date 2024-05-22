package repositories

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/utility"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type RiskRegisterRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *riskDTOs.RiskRegisterListQueryParamsDTO) (*riskDTOs.RiskRegisterListResponseDTO, error)
	GetByID(id uint) (*models.RiskRegister, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskRegister, error)
	Create(riskRegister *models.RiskRegister, tx ...utility.DatabaseTransaction) (uint, error)
	Update(riskRegister *models.RiskRegister, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	// DeleteRiskControls(controlIds []*uint, tx ...utility.DatabaseTransaction) error
	// DeleteRiskComplianceRequirements(requirementIds []*uint, tx ...utility.DatabaseTransaction) error
	GenerateReferenceNumber() (string, error)

	//Document
	CreateDocuments(docs []*models.RiskDocument, tx ...utility.DatabaseTransaction) error
	// SaveRiskComplianceRequirements(complianceRequirements []*models.RiskComplianceRequirement, tx ...utility.DatabaseTransaction) error
	// SaveRiskControls(controls []*models.RiskControl, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.RiskDocument, error)
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error
}

type RiskRegisterRepository struct {
	db *gorm.DB
}

func NewRiskRegisterRepository(db *gorm.DB) *RiskRegisterRepository {
	return &RiskRegisterRepository{
		db: db,
	}
}

func (repo *RiskRegisterRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *RiskRegisterRepository) GetAll(params *riskDTOs.RiskRegisterListQueryParamsDTO) (*riskDTOs.RiskRegisterListResponseDTO, error) {
	var riskRegisters []*riskDTOs.RiskRegisterListDTO
	var totalCount int64

	query := repo.db.Table("risk_register").
		Select(`
            risk_register.id,
            risk_register.reference_number,
            risk_register.title,
            risk_register.description,
            risk_register.reported_at,
            risk_register.created_at,
            risk_register.updated_at,
            risk_categories.title_en AS risk_category_title_en,
            risk_categories.title_ar AS risk_category_title_ar,
            risk_sources.title_en AS risk_source_title_en,
            risk_sources.title_ar AS risk_source_title_ar,
            organization_units.title_en AS owner_unit_title_en,
            organization_units.title_ar AS owner_unit_title_ar,
            users.first_name_en AS owner_first_name_en,
            users.first_name_ar AS owner_first_name_ar,
            users.last_name_en AS owner_last_name_en,
            users.last_name_ar AS owner_last_name_ar,
            risk_statuses.title_en AS risk_status_title_en,
            risk_statuses.title_ar AS risk_status_title_ar,
            risk_statuses.color_code AS risk_status_color_code,
            review_frequencies.title_en AS review_frequency_title_en,
            review_frequencies.title_ar AS review_frequency_title_ar,
            risk_register.next_review_date,
            risk_register.comments,
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
            LEFT JOIN risk_categories ON risk_categories.id = risk_register.risk_category_id
            LEFT JOIN risk_sources ON risk_sources.id = risk_register.risk_source_id
            LEFT JOIN organization_units ON organization_units.id = risk_register.owner_unit_id
            LEFT JOIN users ON users.id = risk_register.owner_id
            LEFT JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id
            LEFT JOIN users AS created_user ON created_user.id = risk_register.created_by
            LEFT JOIN users AS updated_user ON updated_user.id = risk_register.updated_by
            LEFT JOIN review_frequencies ON review_frequencies.id = risk_register.review_frequency_id
        `)
	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	if params.ReportedAtFrom != "" {
		query = query.Where("reported_at >= ?", params.ReportedAtFrom)
	}

	if params.ReportedAtTo != "" {
		query = query.Where("reported_at >= ?", params.ReportedAtTo)
	}

	if params.NextReviewDateFrom != "" {
		query = query.Where("reported_at >= ?", params.NextReviewDateFrom)
	}

	if params.NextReviewDateTo != "" {
		query = query.Where("reported_at >= ?", params.NextReviewDateTo)
	}

	//@todo dynamic filters
	// Apply quick filters
	for _, filter := range params.QuickFilter {
		switch filter {
		case "open_high_risk":
			query = query.Where("risk_statuses.type = ?", "open")

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
			(risk_register.title LIKE ? OR
			risk_register.description LIKE ? OR
			risk_register.reference_number LIKE ? OR
			risk_categories.title_en LIKE ? OR
			risk_categories.title_ar LIKE ? OR
			risk_sources.title_en LIKE ? OR
			risk_sources.title_ar LIKE ? OR
			users.first_name_en LIKE ? OR
			users.first_name_ar LIKE ? OR
			users.last_name_en LIKE ? OR
			users.last_name_ar LIKE ? OR
			risk_statuses.title_en LIKE ? OR
			risk_statuses.title_ar LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.last_name_ar LIKE ? OR
			updated_user.first_name_en LIKE ? OR
			updated_user.first_name_ar LIKE ? OR
			updated_user.last_name_en LIKE ? OR
			updated_user.last_name_ar LIKE ? OR
			review_frequencies.title_en LIKE ? OR
			review_frequencies.title_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
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
	err := query.Offset(offset).Limit(params.PageSize).Find(&riskRegisters).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &riskDTOs.RiskRegisterListResponseDTO{
		Items: riskRegisters,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *RiskRegisterRepository) GetByID(id uint) (*models.RiskRegister, error) {
	var risk *models.RiskRegister

	result := repo.db.First(&risk, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return risk, nil

}

func (repo *RiskRegisterRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskRegister, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var risk *models.RiskRegister

	result := db.
		Preload("RiskCategory").
		Preload("RiskSource").
		Preload("OwnerUnit").
		Preload("Owner").
		Preload("ReviewFrequency").
		Preload("RiskStatus").
		Preload("Documents").
		// Preload("Controls").
		// Preload("ComplianceRequirements").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&risk, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return risk, nil
}

func (repo *RiskRegisterRepository) Create(riskRegister *models.RiskRegister, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(riskRegister).Error; err != nil {
		return 0, err
	}

	return riskRegister.ID, nil
}

func (repo *RiskRegisterRepository) Update(riskRegister *models.RiskRegister, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(riskRegister).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskRegisterRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the Risk register
	if err := db.Delete(&models.RiskRegister{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskRegisterRepository) GenerateReferenceNumber() (string, error) {
	var lastReferenceNumber string
	err := repo.db.Model(&models.RiskRegister{}).
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
	newReferenceNumber := fmt.Sprintf("%s-%d-%04d", "RR", currentYear, lastReferenceID+1)

	return newReferenceNumber, nil
}

func (repo *RiskRegisterRepository) GetDocumentByID(id uint) (*models.RiskDocument, error) {
	var document models.RiskDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *RiskRegisterRepository) CreateDocuments(docs []*models.RiskDocument, tx ...utility.DatabaseTransaction) error {
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

func (repo *RiskRegisterRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete documents by IDs
	if err := db.Where("id IN (?)", docIds).Delete(&models.RiskDocument{}).Error; err != nil {
		return err
	}

	return nil
}

// func (repo *RiskRegisterRepository) SaveRiskComplianceRequirements(complianceRequirements []*models.RiskComplianceRequirement, tx ...utility.DatabaseTransaction) error {
// 	db := repo.db
// 	if len(tx) > 0 && tx[0] != nil {
// 		db = tx[0].DB()
// 	}
// 	if err := db.Create(&complianceRequirements).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (repo *RiskRegisterRepository) SaveRiskControls(controls []*models.RiskControl, tx ...utility.DatabaseTransaction) error {
// 	db := repo.db
// 	if len(tx) > 0 && tx[0] != nil {
// 		db = tx[0].DB()
// 	}
// 	if err := db.Create(&controls).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (repo *RiskRegisterRepository) DeleteRiskControls(controlIds []*uint, tx ...utility.DatabaseTransaction) error {
// 	db := repo.db
// 	if len(tx) > 0 && tx[0] != nil {
// 		db = tx[0].DB()
// 	}

// 	// Delete controls by IDs
// 	if err := db.Where("id IN (?)", controlIds).Delete(&models.RiskControl{}).Error; err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (repo *RiskRegisterRepository) DeleteRiskComplianceRequirements(requirementIds []*uint, tx ...utility.DatabaseTransaction) error {
// 	db := repo.db
// 	if len(tx) > 0 && tx[0] != nil {
// 		db = tx[0].DB()
// 	}

// 	// Delete compliance requirements by IDs
// 	if err := db.Where("id IN (?)", requirementIds).Delete(&models.RiskComplianceRequirement{}).Error; err != nil {
// 		return err
// 	}

// 	return nil
// }
