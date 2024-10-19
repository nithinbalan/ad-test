package repositories

import (
	"errors"
	commonDTOs "jethur/grc/common/dtos"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/models"
	"jethur/utility"
	"math"
	"time"

	"gorm.io/gorm"
)

type ComplianceRequirementRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *complianceDTOs.ComplianceRequirementListQueryParamsDTO) (*complianceDTOs.ComplianceRequirementListResponseDTO, error)
	GetRequirementsByComplianceID(complianceID uint, params *complianceDTOs.ComplianceRequirementListQueryParamsDTO) (*complianceDTOs.ComplianceRequirementByComplianceIDListResponseDTO, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.ComplianceRequirement, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ComplianceRequirement, error)
	Create(complianceRequirement *models.ComplianceRequirement, tx ...utility.DatabaseTransaction) (uint, error)
	Update(complianceRequirement *models.ComplianceRequirement, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	GetOrCreateComplianceClause(complianceClause *models.ComplianceClause, tx ...utility.DatabaseTransaction) (*models.ComplianceClause, error)
	CountRequirementsByClauseID(clauseID uint, tx ...utility.DatabaseTransaction) (int64, error)
	DeleteClause(clauseID uint, tx ...utility.DatabaseTransaction) error

	CreateComplianceRequirementControls(controls []*models.ComplianceControl, tx ...utility.DatabaseTransaction) error
	DeleteComplianceRequirementControls(controlIds []*uint, tx ...utility.DatabaseTransaction) error
	CreateComplianceRequirementStrategicObjectives(objectives []*models.ComplianceRequirementStrategicObjective, tx ...utility.DatabaseTransaction) error
	DeleteComplianceRequirementStrategicObjectives(objectiveIds []*uint, tx ...utility.DatabaseTransaction) error

	CreateInitialComplianceRequirementStatusUpdate(complianceRequirementStatusUpdate *models.ComplianceRequirementStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error)
	GetStatusUpdateByRequirementID(requirementID uint, tx ...utility.DatabaseTransaction) ([]*models.ComplianceRequirementStatusUpdate, error)
	DeleteComplianceRequirementStatusUpdates(ids []*uint, tx ...utility.DatabaseTransaction) error

	CreateComplianceRequirementDocuments(documents []*models.ComplianceRequirementAssociatedDocument, tx ...utility.DatabaseTransaction) error
	DeleteComplianceRequirementDocuments(documentIds []*uint, tx ...utility.DatabaseTransaction) error

	CalculateAverageScore(complianceRegisterID uint, tx ...utility.DatabaseTransaction) (float64, error)
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
			CONVERT(VARCHAR, compliance_requirements.deadline, 23) AS deadline,
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
			responsible_user.id AS responsible_user_id,
			responsible_user.first_name_en AS responsible_user_first_name_en,
            responsible_user.first_name_ar AS responsible_user_first_name_ar,
            responsible_user.last_name_en AS responsible_user_last_name_en,
            responsible_user.last_name_ar AS responsible_user_last_name_ar,
			responsible_user.initial AS responsible_user_initial,
            responsible_user.color AS responsible_user_color,
            responsible_user.email AS responsible_user_email,
            responsible_user.profile_image_path AS responsible_user_profile_image_path,
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
			LEFT JOIN compliance_requirement_categories ON compliance_requirement_categories.id = compliance_requirements.compliance_requirement_category_id
			LEFT JOIN severity_levels ON severity_levels.id = compliance_requirements.severity_level_id
			LEFT JOIN compliance_statuses ON compliance_statuses.id = compliance_requirements.compliance_status_id
			LEFT JOIN users AS responsible_user ON responsible_user.id = compliance_requirements.responsible_user_id
			LEFT JOIN users AS created_user ON created_user.id = compliance_requirements.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = compliance_requirements.updated_by_id
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
			responsible_user.last_name_ar LIKE ? OR
            responsible_user.initial LIKE ? OR
            responsible_user.color LIKE ? OR
            responsible_user.profile_image_path LIKE ? OR
			responsible_user.email LIKE ?)`,
			keyword, keyword, keyword, keyword,
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

func (repo *ComplianceRequirementRepository) GetRequirementsByComplianceID(complianceID uint, params *complianceDTOs.ComplianceRequirementListQueryParamsDTO) (*complianceDTOs.ComplianceRequirementByComplianceIDListResponseDTO, error) {
	var response []*complianceDTOs.ComplianceRequirementByComplianceIDListDTO
	var totalCount int64

	query := repo.db.Table("compliance_requirements").
		Select(`
			compliance_clauses.id AS clause_id,
			compliance_clauses.clause_no AS compliance_clause_no,
			compliance_clauses.title AS compliance_clause_title,
			compliance_requirements.*,
			CONVERT(VARCHAR, compliance_requirements.deadline, 23) AS deadline,
			compliance_requirement_categories.title_en AS compliance_requirement_category_title_en,
			compliance_requirement_categories.title_ar AS compliance_requirement_category_title_ar,
			severity_levels.title_en AS severity_level_title_en,
			severity_levels.title_ar AS severity_level_title_ar,
			severity_levels.color_code AS severity_level_color_code,
			compliance_statuses.title_en AS compliance_status_title_en,
			compliance_statuses.title_ar AS compliance_status_title_ar,
			compliance_statuses.color_code AS compliance_status_color_code,
			responsible_user.id AS responsible_user_id,
			responsible_user.first_name_en AS responsible_user_first_name_en,
            responsible_user.first_name_ar AS responsible_user_first_name_ar,
            responsible_user.last_name_en AS responsible_user_last_name_en,
            responsible_user.last_name_ar AS responsible_user_last_name_ar,
			responsible_user.initial AS responsible_user_initial,
            responsible_user.color AS responsible_user_color,
            responsible_user.email AS responsible_user_email,
            responsible_user.profile_image_path AS responsible_user_profile_image_path,
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
			LEFT JOIN compliance_requirement_categories ON compliance_requirement_categories.id = compliance_requirements.compliance_requirement_category_id
			LEFT JOIN severity_levels ON severity_levels.id = compliance_requirements.severity_level_id
			LEFT JOIN compliance_statuses ON compliance_statuses.id = compliance_requirements.compliance_status_id
			LEFT JOIN users AS responsible_user ON responsible_user.id = compliance_requirements.responsible_user_id
			LEFT JOIN users AS created_user ON created_user.id = compliance_requirements.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = compliance_requirements.updated_by_id
			LEFT JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		`).Where(`compliance_clauses.compliance_register_id = ?`, complianceID)

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	if params.DeadlineFrom != "" {
		query = query.Where("deadline >= ?", params.DeadlineFrom)
	}
	if params.DeadlineTo != "" {
		query = query.Where("deadline <= ?", params.DeadlineTo)
	}

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

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

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
			responsible_user.first_name_ar LIKE ? OR
			responsible_user.last_name_en LIKE ? OR
			responsible_user.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword,
			keyword, keyword)
	}

	if params.SortBy == "" {
		params.SortBy = "compliance_clauses.clause_no"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}
	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize)

	rows, err := query.Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	clauseRequirements := make(map[uint]*complianceDTOs.ComplianceRequirementByComplianceIDListDTO)

	for rows.Next() {
		var complianceReq complianceDTOs.ComplianceRequirementListDTO
		err := repo.db.ScanRows(rows, &complianceReq)
		if err != nil {
			return nil, err
		}

		if _, ok := clauseRequirements[complianceReq.ClauseID]; !ok {
			clauseRequirements[complianceReq.ClauseID] = &complianceDTOs.ComplianceRequirementByComplianceIDListDTO{
				ClauseID:    complianceReq.ClauseID,
				ClauseNo:    complianceReq.ComplianceClauseNo,
				ClauseTitle: complianceReq.ComplianceClauseTitle,
			}
		}

		clauseRequirements[complianceReq.ClauseID].Requirements = append(clauseRequirements[complianceReq.ClauseID].Requirements, &complianceReq)
	}

	for _, v := range clauseRequirements {
		response = append(response, v)
	}

	// Ensure response is not null
	if response == nil {
		response = []*complianceDTOs.ComplianceRequirementByComplianceIDListDTO{}
	}

	// Calculate total pages for pagination
	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	// Create pagination info
	paginationInfo := commonDTOs.PaginationInfoDTO{
		CurrentPage: params.Page,
		PageSize:    params.PageSize,
		TotalPages:  totalPages,
		TotalCount:  totalCount,
	}

	// Wrap response and pagination info in a single DTO
	return &complianceDTOs.ComplianceRequirementByComplianceIDListResponseDTO{
		Items:      response,
		Pagination: paginationInfo,
	}, nil
}

func (repo *ComplianceRequirementRepository) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.ComplianceRequirement, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var complianceRequirement *models.ComplianceRequirement
	result := db.
		Preload("ComplianceClause").
		First(&complianceRequirement, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return complianceRequirement, nil
}

func (repo *ComplianceRequirementRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ComplianceRequirement, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var complianceRequirement *models.ComplianceRequirement

	result := db.
		Preload("ComplianceClause").
		Preload("ComplianceStatus").
		Preload("ComplianceRequirementCategory").
		Preload("SeverityLevel").
		Preload("ResponsibleUser").
		Preload("Controls").
		Preload("Controls.Control.ControlCategory").
		Preload("Controls.Control.ControlType").
		Preload("StrategicObjectives.StrategicObjective").
		Preload("Documents").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
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

func (repo *ComplianceRequirementRepository) CreateInitialComplianceRequirementStatusUpdate(complianceRequirementStatusUpdate *models.ComplianceRequirementStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(complianceRequirementStatusUpdate).Error; err != nil {
		return 0, err
	}
	return complianceRequirementStatusUpdate.ID, nil
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

func (repo *ComplianceRequirementRepository) CountRequirementsByClauseID(clauseID uint, tx ...utility.DatabaseTransaction) (int64, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var count int64

	if err := db.Model(&models.ComplianceRequirement{}).Where("compliance_clause_id = ?", clauseID).Count(&count).Error; err != nil {
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

func (repo *ComplianceRequirementRepository) CreateComplianceRequirementControls(controls []*models.ComplianceControl, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&controls).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ComplianceRequirementRepository) DeleteComplianceRequirementControls(controlIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete controls by IDs
	if err := db.Where("id IN (?)", controlIds).Delete(&models.ComplianceControl{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ComplianceRequirementRepository) CreateComplianceRequirementStrategicObjectives(objectives []*models.ComplianceRequirementStrategicObjective, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&objectives).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ComplianceRequirementRepository) DeleteComplianceRequirementStrategicObjectives(objectiveIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete strategic objectives by IDs
	if err := db.Where("id IN (?)", objectiveIds).Delete(&models.ComplianceRequirementStrategicObjective{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ComplianceRequirementRepository) CreateComplianceRequirementDocuments(documents []*models.ComplianceRequirementAssociatedDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&documents).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ComplianceRequirementRepository) DeleteComplianceRequirementDocuments(documentIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", documentIds).Delete(&models.ComplianceRequirementAssociatedDocument{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ComplianceRequirementRepository) DeleteComplianceRequirementStatusUpdates(ids []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Where("compliance_requirement_status_update_id IN (?)", ids).Delete(&models.ComplianceRequirementStatusUpdateDocument{}).Error; err != nil {
		return err
	}

	if err := db.Where("id IN (?)", ids).Delete(&models.ComplianceRequirementStatusUpdate{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ComplianceRequirementRepository) GetStatusUpdateByRequirementID(requirementID uint, tx ...utility.DatabaseTransaction) ([]*models.ComplianceRequirementStatusUpdate, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var requirementStatusUpdates []*models.ComplianceRequirementStatusUpdate
	result := db.Where("compliance_requirement_id = ?", requirementID).Find(&requirementStatusUpdates)

	if result.Error != nil {
		return nil, result.Error
	}

	return requirementStatusUpdates, nil
}

func (repo *ComplianceRequirementRepository) CalculateAverageScore(complianceRegisterID uint, tx ...utility.DatabaseTransaction) (float64, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var totalScore int
	var requirementCount int

	// Join ComplianceRequirement with ComplianceClause and ComplianceStatus to filter by ComplianceRegister ID
	// and exclude requirements with status "exempt" or "not_applicable"
	err := db.Table("compliance_requirements").
		Joins("JOIN compliance_clauses ON compliance_requirements.compliance_clause_id = compliance_clauses.id").
		Joins("JOIN compliance_statuses ON compliance_requirements.compliance_status_id = compliance_statuses.id").
		Where("compliance_clauses.compliance_register_id = ?", complianceRegisterID).
		Where("compliance_statuses.type NOT IN ('exempt', 'not_applicable')").
		Select("SUM(compliance_requirements.score) AS total_score, COUNT(compliance_requirements.id) AS requirement_count").
		Row().
		Scan(&totalScore, &requirementCount)
	if err != nil {
		return 0, err
	}

	// Calculate the average score
	var averageScore float64
	if requirementCount > 0 {
		averageScore = float64(totalScore) / float64(requirementCount)
	} else {
		averageScore = 0
	}

	// Round the average score to 2 decimal places
	roundedScore := math.Round(averageScore*100) / 100

	return roundedScore, nil
}
