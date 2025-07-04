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

type RiskAppetiteStatementRepositoryInterface interface {
	GetAll(params *riskDTOs.RiskAppetiteStatementListQueryParamsDTO) (*riskDTOs.RiskAppetiteStatementListResponseDTO, error)
	GetAllByCategories(params *riskDTOs.RiskAppetiteStatementListQueryParamsDTO) (*riskDTOs.RiskAppetiteStatementListGroupedResponseDTO, error)
	GetByID(id uint) (*models.RiskAppetiteStatement, error)
	GetFullDetailsByID(id uint) (*models.RiskAppetiteStatement, error)
	Create(RiskAppetiteStatement *models.RiskAppetiteStatement) (uint, error)
	Update(RiskAppetiteStatement *models.RiskAppetiteStatement) error
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByStatement(StatementEn, StatementAr string, RiskCategoryID uint) (bool, error)
	ExistsByStatementExcludingID(id uint, StatementEn, StatementAr string, RiskCategoryID uint) (bool, error)
	GenerateReferenceNumber() (string, error)

	GetRiskAppetiteSummary(params *riskDTOs.RiskAppetiteStatementSummaryQueryParamsDTO) ([]riskDTOs.RiskAppetiteStatementSummaryDTO, error)
	GetDashboardSummary() (*riskDTOs.RiskAppetiteDashboardSummaryDTO, error)
	GetMappedRisks(statementID uint, params *riskDTOs.MappedRiskListQueryParamsDTO) (*riskDTOs.MappedRiskListResponseDTO, error)
}

type RiskAppetiteStatementRepository struct {
	db *gorm.DB
}

func NewRiskAppetiteStatementRepository(db *gorm.DB) *RiskAppetiteStatementRepository {
	return &RiskAppetiteStatementRepository{
		db: db,
	}
}

func (repo *RiskAppetiteStatementRepository) GetAll(params *riskDTOs.RiskAppetiteStatementListQueryParamsDTO) (*riskDTOs.RiskAppetiteStatementListResponseDTO, error) {

	var riskAppetiteStatements []*riskDTOs.RiskAppetiteStatementListDTO
	var totalCount int64

	query := repo.db.Table("risk_appetite_statements").
		Select(`
            risk_appetite_statements.*,
			risk_categories.title_en AS risk_category_title_en,
			risk_categories.title_ar AS risk_category_title_ar,
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
		    LEFT JOIN risk_categories ON risk_categories.id = risk_appetite_statements.risk_category_id
			LEFT JOIN users AS created_user ON created_user.id = risk_appetite_statements.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = risk_appetite_statements.updated_by_id
		`)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply quick filters
	for _, filter := range params.QuickFilter {
		query = query.Where("risk_categories.title_en = ? OR risk_categories.title_ar = ?", filter, filter)
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		query = query.Where(`
			(risk_appetite_statements.statement_en LIKE ? OR
			risk_appetite_statements.statement_ar LIKE ? OR
			risk_categories.title_en LIKE ? OR
			risk_categories.title_ar LIKE ?)`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "risk_appetite_statements.statement_en" // Default sort by risk_appetite_statements.statement_en if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the queryc
	err := query.Offset(offset).Limit(params.PageSize).Find(&riskAppetiteStatements).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &riskDTOs.RiskAppetiteStatementListResponseDTO{
		Items: riskAppetiteStatements,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}

	//@todo The response will be changed based on risk implementation

	return response, nil
}

func (repo *RiskAppetiteStatementRepository) GetAllByCategories(params *riskDTOs.RiskAppetiteStatementListQueryParamsDTO) (*riskDTOs.RiskAppetiteStatementListGroupedResponseDTO, error) {
	var totalCount int64
	var response []*riskDTOs.RiskAppetiteStatementGroupedDTO

	query := repo.db.Table("risk_appetite_statements").
		Select(`
            risk_appetite_statements.*,
			risk_categories.title_en AS risk_category_title_en,
			risk_categories.title_ar AS risk_category_title_ar,
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
		    LEFT JOIN risk_categories ON risk_categories.id = risk_appetite_statements.risk_category_id
			LEFT JOIN users AS created_user ON created_user.id = risk_appetite_statements.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = risk_appetite_statements.updated_by_id
		`)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply quick filters
	for _, filter := range params.QuickFilter {
		query = query.Where("risk_categories.title_en = ? OR risk_categories.title_ar = ?", filter, filter)
	}
	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		query = query.Where(`
			(risk_appetite_statements.statement_en LIKE ? OR
			risk_appetite_statements.statement_ar LIKE ? OR
			risk_categories.title_en LIKE ? OR
			risk_categories.title_ar LIKE ?)`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "risk_appetite_statements.statement_en"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Add OFFSET and FETCH NEXT for pagination
	query = query.Offset(offset).Limit(params.PageSize)

	rows, err := query.Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	riskStatements := make(map[uint]*riskDTOs.RiskAppetiteStatementGroupedDTO)

	for rows.Next() {
		var riskStatement riskDTOs.RiskAppetiteStatementListDTO
		err := repo.db.ScanRows(rows, &riskStatement)
		if err != nil {
			return nil, err
		}

		if _, ok := riskStatements[riskStatement.RiskCategoryID]; !ok {
			riskStatements[riskStatement.RiskCategoryID] = &riskDTOs.RiskAppetiteStatementGroupedDTO{
				RiskCategoryID:      riskStatement.RiskCategoryID,
				RiskCategoryTitleEn: riskStatement.RiskCategoryTitleEn,
				RiskCategoryTitleAr: riskStatement.RiskCategoryTitleAr,
			}
		}
		riskStatements[riskStatement.RiskCategoryID].Statements = append(riskStatements[riskStatement.RiskCategoryID].Statements, &riskStatement)
	}

	//@to be completed after risk evaluation.
	for _, grouped := range riskStatements {
		for _, statement := range grouped.Statements {
			// Initialize the Risks field to an empty slice

			var risks []*riskDTOs.RiskRegisterListDTO
			err := repo.db.Table("risk_register").
				Limit(2).Find(&risks).Error
			if err != nil {
				return nil, err
			}
			if len(risks) == 0 {
				statement.Risks = []*riskDTOs.RiskRegisterListDTO{}
			} else {
				statement.Risks = risks
			}
		}
		response = append(response, grouped)
	}

	if response == nil {
		response = []*riskDTOs.RiskAppetiteStatementGroupedDTO{}
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
	return &riskDTOs.RiskAppetiteStatementListGroupedResponseDTO{
		Items:      response,
		Pagination: paginationInfo,
	}, nil
}

func (repo *RiskAppetiteStatementRepository) GetByID(id uint) (*models.RiskAppetiteStatement, error) {
	var riskAppetiteStatement *models.RiskAppetiteStatement

	result := repo.db.First(&riskAppetiteStatement, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return riskAppetiteStatement, nil

}

func (repo *RiskAppetiteStatementRepository) GetFullDetailsByID(id uint) (*models.RiskAppetiteStatement, error) {
	var riskAppetiteStatement *models.RiskAppetiteStatement // Declare the variable inside the function

	result := repo.db.
		Preload("RiskCategory").
		First(&riskAppetiteStatement, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return riskAppetiteStatement, nil
}

func (repo *RiskAppetiteStatementRepository) Create(riskAppetiteStatement *models.RiskAppetiteStatement) (uint, error) {

	if err := repo.db.Create(riskAppetiteStatement).Error; err != nil {
		return 0, err
	}

	return riskAppetiteStatement.ID, nil
}

func (repo *RiskAppetiteStatementRepository) Update(riskAppetiteStatement *models.RiskAppetiteStatement) error {

	if err := repo.db.Save(riskAppetiteStatement).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskAppetiteStatementRepository) Delete(id uint) error {

	var riskAppetiteStatement *models.RiskAppetiteStatement
	if err := repo.db.First(&riskAppetiteStatement, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.RiskAppetiteStatement{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *RiskAppetiteStatementRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.RiskAppetiteStatement{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.RiskAppetiteStatement{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil
}

func (repo *RiskAppetiteStatementRepository) Deactivate(ids []uint) ([]uint, error) {

	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.RiskAppetiteStatement{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.RiskAppetiteStatement{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil

}

func (repo *RiskAppetiteStatementRepository) ExistsByStatement(StatementEn, StatementAr string, RiskCategoryID uint) (bool, error) {
	var count int64

	// Check for the existence of Title or QuestionAr in a single query
	err := repo.db.Model(&models.RiskAppetiteStatement{}).
		Where("(statement_en = ? OR statement_ar = ?) AND risk_category_id = ?", StatementEn, StatementAr, RiskCategoryID).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	// Return true if either title exists, otherwise false
	return count > 0, nil
}

func (repo *RiskAppetiteStatementRepository) ExistsByStatementExcludingID(id uint, StatementEn, StatementAr string, RiskCategoryID uint) (bool, error) {
	var count int64
	err := repo.db.Model(&models.RiskAppetiteStatement{}).
		Where("id <> ? AND (statement_en = ? OR statement_ar = ?) AND risk_category_id = ?", id, StatementEn, StatementAr, RiskCategoryID).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}

func (repo *RiskAppetiteStatementRepository) GenerateReferenceNumber() (string, error) {
	var lastReferenceNumber string
	err := repo.db.Model(&models.RiskAppetiteStatement{}).
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
	newReferenceNumber := fmt.Sprintf("%s-%d-%04d", "RAS", currentYear, lastReferenceID+1)

	return newReferenceNumber, nil
}

func (repo *RiskAppetiteStatementRepository) GetRiskAppetiteSummary(
	params *riskDTOs.RiskAppetiteStatementSummaryQueryParamsDTO,
) ([]riskDTOs.RiskAppetiteStatementSummaryDTO, error) {

	var summaries []riskDTOs.RiskAppetiteStatementSummaryDTO

	query := repo.db.
		Table("risk_appetite_statements AS ras").
		Select(`
			ras.id AS id,
			rc.title_en AS risk_category_title_en,
			rc.title_ar AS risk_category_title_ar,
			ras.statement_en AS risk_appetite_statement_title_en,
			ras.statement_ar AS risk_appetite_statement_title_ar,
			COUNT(DISTINCT re.risk_register_id) AS risks_mapped,
			SUM(CASE WHEN ratl.type = 'above' THEN 1 ELSE 0 END) AS above_count,
			SUM(CASE WHEN ratl.type = 'within' THEN 1 ELSE 0 END) AS within_count,
			SUM(CASE WHEN ratl.type = 'below' THEN 1 ELSE 0 END) AS below_count
		`).
		Joins(`LEFT JOIN risk_categories rc ON rc.id = ras.risk_category_id`).
		Joins(`LEFT JOIN risk_appetite_alignments raa ON raa.risk_appetite_statement_id = ras.id`).
		Joins(`LEFT JOIN risk_evaluations re ON re.id = raa.risk_evaluation_id AND re.type = 'inherent'`).
		Joins(`LEFT JOIN risk_appetite_tolerance_levels ratl ON ratl.id = raa.risk_appetite_tolerance_level_id`).
		Group(`ras.id, rc.title_en, rc.title_ar, ras.statement_en, ras.statement_ar`)

	// Apply dynamic filters
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply quick filters
	for _, filter := range params.QuickFilter {
		query = query.Where(
			"rc.title_en = ? OR rc.title_ar = ?",
			filter, filter,
		)
	}

	// Apply search keyword
	if params.Keyword != "" {
		like := "%" + params.Keyword + "%"
		query = query.Where(`
			ras.statement_en LIKE ? OR
			ras.statement_ar LIKE ? OR
			rc.title_en LIKE ? OR
			rc.title_ar LIKE ?`,
			like, like, like, like,
		)
	}

	// Sorting
	sortBy := params.SortBy
	if sortBy == "" {
		sortBy = "ras.statement_en"
	}
	sortOrder := params.Sort
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}
	query = query.Order(sortBy + " " + sortOrder)

	// Fetch all results
	if err := query.Scan(&summaries).Error; err != nil {
		return nil, err
	}

	return summaries, nil
}

func (repo *RiskAppetiteStatementRepository) GetDashboardSummary() (*riskDTOs.RiskAppetiteDashboardSummaryDTO, error) {
	var summary riskDTOs.RiskAppetiteDashboardSummaryDTO

	var activeCount int64
	var totalRisks int64
	var alignedRisks int64

	// 1. Active statements
	if err := repo.db.
		Table("risk_appetite_statements").
		Where("status = ?", "active").
		Count(&activeCount).Error; err != nil {
		return nil, err
	}
	summary.ActiveStatementsCount = int(activeCount)

	// 2. Total risks
	if err := repo.db.
		Table("risk_register").
		Count(&totalRisks).Error; err != nil {
		return nil, err
	}
	summary.TotalRisksCount = int(totalRisks)

	// 3. Aligned risks (only inherent)
	if err := repo.db.
		Table("risk_appetite_alignments AS raa").
		Joins("INNER JOIN risk_evaluations re ON re.id = raa.risk_evaluation_id").
		Where("re.type = ?", "inherent").
		Select("COUNT(DISTINCT re.risk_register_id)").
		Count(&alignedRisks).Error; err != nil {
		return nil, err
	}
	summary.AlignedRisksCount = int(alignedRisks)

	// 4–6. Tolerance counts (only inherent)
	type CountResult struct {
		Type  string
		Count int
	}
	var counts []CountResult
	if err := repo.db.
		Table("risk_appetite_alignments AS raa").
		Joins("INNER JOIN risk_evaluations re ON re.id = raa.risk_evaluation_id").
		Joins("LEFT JOIN risk_appetite_tolerance_levels ratl ON ratl.id = raa.risk_appetite_tolerance_level_id").
		Where("re.type = ?", "inherent").
		Select("ratl.type, COUNT(DISTINCT re.risk_register_id) AS count").
		Group("ratl.type").
		Scan(&counts).Error; err != nil {
		return nil, err
	}
	for _, c := range counts {
		switch c.Type {
		case "above":
			summary.AboveToleranceRisksCount = c.Count
		case "within":
			summary.WithinToleranceRisksCount = c.Count
		case "below":
			summary.BelowToleranceRisksCount = c.Count
		}
	}

	return &summary, nil
}

func (repo *RiskAppetiteStatementRepository) GetMappedRisks(statementID uint, params *riskDTOs.MappedRiskListQueryParamsDTO) (*riskDTOs.MappedRiskListResponseDTO, error) {

	var risks []*riskDTOs.MappedRiskListDTO
	var totalCount int64

	query := repo.db.Table("risk_appetite_alignments AS raa").
		Select(`
			rr.id AS risk_id,
			rr.reference_number,
			rr.title,
			rr.description,
			rc.title_en AS risk_category_title_en,
			rc.title_ar AS risk_category_title_ar,
			rsrc.title_en AS risk_source_title_en,
			rsrc.title_ar AS risk_source_title_ar,
			ou.title_en AS owner_unit_title_en,
			ou.title_ar AS owner_unit_title_ar,
			u.id AS owner_id,
			u.first_name_en AS owner_first_name_en,
			u.first_name_ar AS owner_first_name_ar,
			u.last_name_en AS owner_last_name_en,
			u.last_name_ar AS owner_last_name_ar,
			u.initial AS owner_initial,
			u.color AS owner_color,
			u.email AS owner_email,
			u.profile_image_path AS owner_profile_image_path,
			rml.title_en AS risk_matrix_likelihood_title_en,
			rml.title_ar AS risk_matrix_likelihood_title_ar,
			rml.score AS risk_matrix_likelihood_score,
			rmi.title_en AS risk_matrix_impact_title_en,
			rmi.title_ar AS risk_matrix_impact_title_ar,
			rmi.score AS risk_matrix_impact_score,
			rmrr.title_en AS risk_matrix_risk_rating_title_en,
			rmrr.title_ar AS risk_matrix_risk_rating_title_ar,
			rmrr.type AS risk_matrix_risk_rating_type,
			rmrr.color_code AS risk_matrix_risk_rating_color_code,
			rstat.title_en AS risk_status_title_en,
			rstat.title_ar AS risk_status_title_ar,
			rstat.color_code AS risk_status_color_code,
			ratl.type AS appetite_alignment
		`).
		Joins(`INNER JOIN risk_evaluations re ON re.id = raa.risk_evaluation_id AND re.type = 'inherent'`).
		Joins(`INNER JOIN risk_register rr ON rr.id = re.risk_register_id`).
		Joins(`LEFT JOIN risk_categories rc ON rc.id = rr.risk_category_id`).
		Joins(`LEFT JOIN risk_sources rsrc ON rsrc.id = rr.risk_source_id`).
		Joins(`LEFT JOIN organization_units ou ON ou.id = rr.responsible_unit_id`).
		Joins(`LEFT JOIN users u ON u.id = rr.responsible_user_id`).
		Joins(`LEFT JOIN risk_matrix_likelihoods rml ON rml.id = re.risk_matrix_likelihood_id`).
		Joins(`LEFT JOIN risk_matrix_impacts rmi ON rmi.id = re.risk_matrix_impact_id`).
		Joins(`LEFT JOIN risk_matrix_risk_ratings rmrr ON rmrr.id = re.risk_matrix_risk_rating_id`).
		Joins(`LEFT JOIN risk_statuses rstat ON rstat.id = rr.risk_status_id`).
		Joins(`LEFT JOIN risk_appetite_tolerance_levels ratl ON ratl.id = raa.risk_appetite_tolerance_level_id`).
		Where("raa.risk_appetite_statement_id = ?", statementID)

	// Count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Keyword search
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(rr.reference_number LIKE ? OR
			rr.title LIKE ? OR
			rr.description LIKE ? OR
			rc.title_en LIKE ? OR
			rc.title_ar LIKE ? OR
			rsrc.title_en LIKE ? OR
			rsrc.title_ar LIKE ? OR
			ou.title_en LIKE ? OR
			ou.title_ar LIKE ? OR
			u.first_name_en LIKE ? OR
			u.first_name_ar LIKE ? OR
			u.last_name_en LIKE ? OR
			u.last_name_ar LIKE ? OR
			rstat.title_en LIKE ? OR
			rstat.title_ar LIKE ? OR
			rstat.color_code LIKE ? OR
			ratl.type LIKE ?)
		`,
			keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword)
	}

	// Sorting
	if params.SortBy == "" {
		params.SortBy = "rr.title"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}
	query = query.Order(params.SortBy + " " + params.Sort)

	// Pagination
	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	if err := query.Offset(offset).Limit(params.PageSize).Scan(&risks).Error; err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &riskDTOs.MappedRiskListResponseDTO{
		Items: risks,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}

	return response, nil
}
