package repositories

import (
	"errors"
	commonDTOs "jethur/grc/common/dtos"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/models"
	"jethur/utility"
	"time"

	"gorm.io/gorm"
)

type ComplianceRequirementRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *complianceDTOs.ComplianceRequirementListQueryParamsDTO) (*complianceDTOs.ComplianceRequirementListResponseDTO, error)
	GetRequirementsByComplianceID(params *complianceDTOs.ComplianceRequirementListQueryParamsDTO) (*complianceDTOs.ComplianceRequirementByComplianceIDListResponseDTO, error)
	GetByID(id uint) (*models.ComplianceRequirement, error)
	GetFullDetailsByID(id uint) (*models.ComplianceRequirement, error)
	Create(complianceRequirement *models.ComplianceRequirement, tx ...utility.DatabaseTransaction) (uint, error)
	Update(complianceRequirement *models.ComplianceRequirement, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	GetOrCreateComplianceClause(complianceClause *models.ComplianceClause, tx ...utility.DatabaseTransaction) (*models.ComplianceClause, error)
	CountRequirementsByClauseID(clauseID uint) (int64, error)
	DeleteClause(clauseID uint, tx ...utility.DatabaseTransaction) error
}

type ComplianceRequirementRepository struct {
	db *gorm.DB
}

func NewComplianceRequirementRepository(db *gorm.DB) *ComplianceRequirementRepository {
	return &ComplianceRequirementRepository{
		db: db,
	}
}

func (repo *ComplianceRequirementRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ComplianceRequirementRepository) GetAll(params *complianceDTOs.ComplianceRequirementListQueryParamsDTO) (*complianceDTOs.ComplianceRequirementListResponseDTO, error) {
	var complianceRequirements []*complianceDTOs.ComplianceRequirementListDTO
	var totalCount int64

	query := repo.db.Table("compliance_requirements").
		Select(`
			compliance_requirements.*,
			compliance_clauses.clause_no AS compliance_clause_no,
			compliance_clauses.title AS compliance_clause_title,
			compliance_requirement_categories.title_en AS compliance_requirement_category_title_en,
			compliance_requirement_categories.title_ar AS compliance_requirement_category_title_ar,
			severity_levels.title_en AS severity_level_title_en,
			severity_levels.title_ar AS severity_level_title_ar,
			severity_levels.color_code AS severity_level_color_code,
			compliance_statuses.title_en AS compliance_status_title_en,
			compliance_statuses.title_ar AS compliance_status_title_ar,
			compliance_statuses.color_code AS compliance_status_color_code,
			responsible_user.first_name_en AS responsible_user_first_name_en,
            responsible_user.first_name_ar AS responsible_user_first_name_ar,
            responsible_user.last_name_en AS responsible_user_last_name_en,
            responsible_user.last_name_ar AS responsible_user_last_name_ar,
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
			LEFT JOIN compliance_requirement_categories ON compliance_requirement_categories.id = compliance_requirements.compliance_requirement_category_id
			LEFT JOIN severity_levels ON severity_levels.id = compliance_requirements.severity_level_id
			LEFT JOIN compliance_statuses ON compliance_statuses.id = compliance_requirements.compliance_status_id
			LEFT JOIN users AS responsible_user ON responsible_user.id = compliance_requirements.responsible_user_id
			LEFT JOIN users AS created_user ON created_user.id = compliance_requirements.created_by
			LEFT JOIN users AS updated_user ON updated_user.id = compliance_requirements.updated_by
			LEFT JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		`)
	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply deadline filters
	if params.DeadlineFrom != "" {
		query = query.Where("deadline >= ?", params.DeadlineFrom)
	}
	if params.DeadlineTo != "" {
		query = query.Where("deadline <= ?", params.DeadlineTo)
	}

	// Apply quick filters
	for _, filter := range params.QuickFilter {
		switch filter {
		case "severe":
			query = query.Where("severity_levels.type = ?", "high")
		case "non_compliant":
			query = query.Where("compliance_statuses.type = ?", "not_implemented")
		case "non_compliant_due":
			query = query.Where("compliance_statuses.type = ? AND compliance_requirements.deadline < ?", "not_implemented", time.Now())
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
			(compliance_clauses.clause_no LIKE ? OR
			compliance_clauses.title LIKE ? OR
			compliance_requirements.title LIKE ? OR
			compliance_requirement_categories.title_en LIKE ? OR
			compliance_requirement_categories.title_ar LIKE ? OR
			severity_levels.title_en LIKE ? OR
			severity_levels.title_ar LIKE ? OR
			compliance_requirements.deadline LIKE ? OR
			compliance_statuses.title_en LIKE ? OR
			compliance_statuses.title_ar LIKE ? OR
			responsible_user.first_name_en LIKE ? OR
			responsible_user.first_name_en LIKE ? OR
			responsible_user.first_name_ar LIKE ? OR
			responsible_user.last_name_en LIKE ? OR
			responsible_user.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword,
			keyword, keyword, keyword)
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
	err := query.Offset(offset).Limit(params.PageSize).Find(&complianceRequirements).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &complianceDTOs.ComplianceRequirementListResponseDTO{
		Items: complianceRequirements,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *ComplianceRequirementRepository) GetRequirementsByComplianceID(params *complianceDTOs.ComplianceRequirementListQueryParamsDTO) (*complianceDTOs.ComplianceRequirementByComplianceIDListResponseDTO, error) {
	// Initialize the response DTO with an empty map
	response := &complianceDTOs.ComplianceRequirementByComplianceIDListResponseDTO{
		Items: make(map[string][]*complianceDTOs.ComplianceRequirementListDTO),
	}

	query := repo.db.Table("compliance_requirements").
		Select(`
			compliance_requirements.*,
			compliance_clauses.id AS clause_id,
			compliance_clauses.clause_no AS compliance_clause_no,
			compliance_clauses.title AS compliance_clause_title,
			compliance_requirement_categories.title_en AS compliance_requirement_category_title_en,
			compliance_requirement_categories.title_ar AS compliance_requirement_category_title_ar,
			severity_levels.title_en AS severity_level_title_en,
			severity_levels.title_ar AS severity_level_title_ar,
			severity_levels.color_code AS severity_level_color_code,
			compliance_statuses.title_en AS compliance_status_title_en,
			compliance_statuses.title_ar AS compliance_status_title_ar,
			compliance_statuses.color_code AS compliance_status_color_code,
			responsible_user.first_name_en AS responsible_user_first_name_en,
			responsible_user.first_name_ar AS responsible_user_first_name_ar,
			responsible_user.last_name_en AS responsible_user_last_name_en,
			responsible_user.last_name_ar AS responsible_user_last_name_ar,
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
			LEFT JOIN compliance_requirement_categories ON compliance_requirement_categories.id = compliance_requirements.compliance_requirement_category_id
			LEFT JOIN severity_levels ON severity_levels.id = compliance_requirements.severity_level_id
			LEFT JOIN compliance_statuses ON compliance_statuses.id = compliance_requirements.compliance_status_id
			LEFT JOIN users AS responsible_user ON responsible_user.id = compliance_requirements.responsible_user_id
			LEFT JOIN users AS created_user ON created_user.id = compliance_requirements.created_by
			LEFT JOIN users AS updated_user ON updated_user.id = compliance_requirements.updated_by
			LEFT JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		`)
	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply deadline filters
	if params.DeadlineFrom != "" {
		query = query.Where("deadline >= ?", params.DeadlineFrom)
	}
	if params.DeadlineTo != "" {
		query = query.Where("deadline <= ?", params.DeadlineTo)
	}

	// Apply quick filters
	for _, filter := range params.QuickFilter {
		switch filter {
		case "severe":
			query = query.Where("severity_levels.type = ?", "high")
		case "non_compliant":
			query = query.Where("compliance_statuses.type = ?", "not_implemented")
		case "non_compliant_due":
			query = query.Where("compliance_statuses.type = ? AND compliance_requirements.deadline < ?", "not_implemented", time.Now())
		}
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(compliance_clauses.clause_no LIKE ? OR
			compliance_clauses.title LIKE ? OR
			compliance_requirements.title LIKE ? OR
			compliance_requirement_categories.title_en LIKE ? OR
			compliance_requirement_categories.title_ar LIKE ? OR
			severity_levels.title_en LIKE ? OR
			severity_levels.title_ar LIKE ? OR
			compliance_requirements.deadline LIKE ? OR
			compliance_statuses.title_en LIKE ? OR
			compliance_statuses.title_ar LIKE ? OR
			responsible_user.first_name_en LIKE ? OR
			responsible_user.first_name_en LIKE ? OR
			responsible_user.first_name_ar LIKE ? OR
			responsible_user.last_name_en LIKE ? OR
			responsible_user.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword,
			keyword, keyword, keyword)
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "title" // Default sort by title if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	// Execute the query
	rows, err := query.Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate through the rows and populate the response DTO
	for rows.Next() {
		var complianceReq complianceDTOs.ComplianceRequirementListDTO
		err := repo.db.ScanRows(rows, &complianceReq)
		if err != nil {
			return nil, err
		}
		// response.Items[complianceReq.ClauseID] = append(response.Items[complianceReq.ClauseID], &complianceReq)
		response.Items[complianceReq.ComplianceClauseTitle] = append(response.Items[complianceReq.ComplianceClauseTitle], &complianceReq)

	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return response, nil
}

func (repo *ComplianceRequirementRepository) GetByID(id uint) (*models.ComplianceRequirement, error) {
	var complianceRequirement *models.ComplianceRequirement
	result := repo.db.First(&complianceRequirement, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return complianceRequirement, nil
}

func (repo *ComplianceRequirementRepository) GetFullDetailsByID(id uint) (*models.ComplianceRequirement, error) {
	var complianceRequirement *models.ComplianceRequirement

	result := repo.db.
		Preload("ComplianceClause").
		Preload("ComplianceStatus").
		Preload("ComplianceRequirementCategory").
		Preload("SeverityLevel").
		Preload("ResponsibleUser").
		First(&complianceRequirement, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return complianceRequirement, nil
}

func (repo *ComplianceRequirementRepository) Create(complianceRequirement *models.ComplianceRequirement, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(complianceRequirement).Error; err != nil {
		return 0, err
	}
	return complianceRequirement.ID, nil
}

func (repo *ComplianceRequirementRepository) Update(complianceRequirement *models.ComplianceRequirement, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(complianceRequirement).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ComplianceRequirementRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var complianceRequirement *models.ComplianceRequirement

	if err := db.First(&complianceRequirement, id).Error; err != nil {
		return err
	}

	result := db.Delete(&models.ComplianceRequirement{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *ComplianceRequirementRepository) GetOrCreateComplianceClause(complianceClause *models.ComplianceClause, tx ...utility.DatabaseTransaction) (*models.ComplianceClause, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var existingClause models.ComplianceClause

	err := db.Where("clause_no = ? AND title = ? AND compliance_register_id = ?", complianceClause.ClauseNo, complianceClause.Title, complianceClause.ComplianceRegisterID).First(&existingClause).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := db.Create(complianceClause).Error; err != nil {
				return nil, err
			}
			return complianceClause, nil
		}
		return nil, err
	}

	// If the clause exists, return it
	return &existingClause, nil
}

func (repo *ComplianceRequirementRepository) CountRequirementsByClauseID(clauseID uint) (int64, error) {
	var count int64

	if err := repo.db.Model(&models.ComplianceRequirement{}).Where("compliance_clause_id = ?", clauseID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *ComplianceRequirementRepository) DeleteClause(clauseID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Delete(&models.ComplianceClause{}, clauseID).Error; err != nil {
		return err
	}

	return nil
}
