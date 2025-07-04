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

	GetAllRiskAppetiteStatments(params *riskDTOs.RiskAppetiteStatementSummaryQueryParamsDTO) ([]riskDTOs.RiskAppetiteStatementSummaryDTO, error)
	GetRiskAppetiteStatementSummary() (*riskDTOs.RiskAppetiteDashboardSummaryDTO, error)
	GetAllRisks(statementID uint, params *riskDTOs.MappedRiskListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error)
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

func (repo *RiskAppetiteStatementRepository) GetAllRiskAppetiteStatments(
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

	// Apply search keyword
	if params.Keyword != "" {
		like := "%" + params.Keyword + "%"
		query = query.Where(`
		rc.title_en LIKE ? OR
			rc.titlet_ar LIKE ? OR
			ras.statement_en LIKE ? OR
			ras.statement_ar LIKE ? OR
			rc.title_en LIKE ? OR
			rc.title_ar LIKE ?`,
			like, like, like, like, like, like,
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

func (repo *RiskAppetiteStatementRepository) GetRiskAppetiteStatementSummary() (*riskDTOs.RiskAppetiteDashboardSummaryDTO, error) {
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
func (repo *RiskAppetiteStatementRepository) GetAllRisks(statementID uint, params *riskDTOs.MappedRiskListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error) {

	var risks []*commonDTOs.AssociatedRiskListDTO
	var totalCount int64

	query := repo.db.Table("risk_register").
		Select(`
			risk_register.id,
			risk_register.reference_number,
			risk_register.title,
			risk_register.description,
			CONVERT(VARCHAR, risk_register.reported_at, 23) AS reported_at,
			risk_register.created_at,
			risk_register.updated_at,
			risk_register.current_risk_rating_value_id,
			risk_categories.title_en AS risk_category_title_en,
			risk_categories.title_ar AS risk_category_title_ar,
			risk_sources.title_en AS risk_source_title_en,
			risk_sources.title_ar AS risk_source_title_ar,
			organization_units.id AS responsible_unit_id,
			organization_units.title_en AS responsible_unit_title_en,
			organization_units.title_ar AS responsible_unit_title_ar,
			users.id AS responsible_user_id,
			users.first_name_en AS responsible_user_first_name_en,
			users.first_name_ar AS responsible_user_first_name_ar,
			users.last_name_en AS responsible_user_last_name_en,
			users.last_name_ar AS responsible_user_last_name_ar,
			users.initial AS responsible_user_initial,
			users.color AS responsible_user_color,
			users.email AS responsible_user_email,
			users.profile_image_path AS responsible_user_profile_image_path,
			risk_statuses.title_en AS risk_status_title_en,
			risk_statuses.title_ar AS risk_status_title_ar,
			risk_response_strategies.title_en AS risk_response_strategy_title_en,
			risk_response_strategies.title_ar AS risk_response_strategy_title_ar,
			risk_statuses.color_code AS risk_status_color_code,
			review_frequencies.title_en AS review_frequency_title_en,
			review_frequencies.title_ar AS review_frequency_title_ar,
			CONVERT(VARCHAR, risk_register.next_review_date, 23) AS next_review_date,
			risk_register.last_review_note,
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
			updated_user.profile_image_path AS updated_user_profile_image_path,
			current_likelihood.title_en AS current_likelihood_title_en,
			current_likelihood.title_ar AS current_likelihood_title_ar,
			current_impact.title_en AS current_impact_title_en,
			current_impact.title_ar AS current_impact_title_ar,
			current_risk_rating.title_en AS current_risk_rating_title_en,
			current_risk_rating.title_ar AS current_risk_rating_title_ar,
			current_risk_rating_value.value AS current_risk_rating_value,
			current_risk_rating.color_code AS current_risk_rating_color_code,
			inherent_likelihood.title_en AS inherent_likelihood_title_en,
			inherent_likelihood.title_ar AS inherent_likelihood_title_ar,
			inherent_impact.title_en AS inherent_impact_title_en,
			inherent_impact.title_ar AS inherent_impact_title_ar,
			inherent_risk_rating.title_en AS inherent_risk_rating_title_en,
			inherent_risk_rating.title_ar AS inherent_risk_rating_title_ar,
			inherent_risk_rating_value.value AS inherent_risk_rating_value,
			inherent_risk_rating.color_code AS inherent_risk_rating_color_code,
			residual_likelihood.title_en AS residual_likelihood_title_en,
			residual_likelihood.title_ar AS residual_likelihood_title_ar,
			residual_impact.title_en AS residual_impact_title_en,
			residual_impact.title_ar AS residual_impact_title_ar,
			residual_risk_rating.title_en AS residual_risk_rating_title_en,
			residual_risk_rating.title_ar AS residual_risk_rating_title_ar,
			residual_risk_rating_value.value AS residual_risk_rating_value,
			residual_risk_rating.color_code AS residual_risk_rating_color_code,
			ratl.type AS appetite_alignment
		`).
		Joins(`
			LEFT JOIN risk_categories ON risk_categories.id = risk_register.risk_category_id
			LEFT JOIN risk_sources ON risk_sources.id = risk_register.risk_source_id
			LEFT JOIN organization_units ON organization_units.id = risk_register.responsible_unit_id
			LEFT JOIN users ON users.id = risk_register.responsible_user_id
			LEFT JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id
			LEFT JOIN risk_response_strategies ON risk_response_strategies.id = risk_register.risk_response_strategy_id
			LEFT JOIN users AS created_user ON created_user.id = risk_register.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = risk_register.updated_by_id
			LEFT JOIN review_frequencies ON review_frequencies.id = risk_register.review_frequency_id
			LEFT JOIN risk_matrix_likelihoods AS current_likelihood ON current_likelihood.id = risk_register.current_likelihood_id
			LEFT JOIN risk_matrix_impacts AS current_impact ON current_impact.id = risk_register.current_impact_id
			LEFT JOIN risk_matrix_risk_ratings AS current_risk_rating ON current_risk_rating.id = risk_register.current_risk_rating_id
			LEFT JOIN risk_matrix_risk_rating_values AS current_risk_rating_value ON current_risk_rating_value.id = risk_register.current_risk_rating_value_id
			LEFT JOIN risk_matrix_likelihoods AS inherent_likelihood ON inherent_likelihood.id = risk_register.inherent_likelihood_id
			LEFT JOIN risk_matrix_impacts AS inherent_impact ON inherent_impact.id = risk_register.inherent_impact_id
			LEFT JOIN risk_matrix_risk_ratings AS inherent_risk_rating ON inherent_risk_rating.id = risk_register.inherent_risk_rating_id
			LEFT JOIN risk_matrix_risk_rating_values AS inherent_risk_rating_value ON inherent_risk_rating_value.id = risk_register.inherent_risk_rating_value_id
			LEFT JOIN risk_matrix_likelihoods AS residual_likelihood ON residual_likelihood.id = risk_register.residual_likelihood_id
			LEFT JOIN risk_matrix_impacts AS residual_impact ON residual_impact.id = risk_register.residual_impact_id
			LEFT JOIN risk_matrix_risk_ratings AS residual_risk_rating ON residual_risk_rating.id = risk_register.residual_risk_rating_id
			LEFT JOIN risk_matrix_risk_rating_values AS residual_risk_rating_value ON residual_risk_rating_value.id = risk_register.residual_risk_rating_value_id
			INNER JOIN risk_evaluations re ON re.risk_register_id = risk_register.id AND re.type = 'inherent'
			INNER JOIN risk_appetite_alignments raa ON raa.risk_evaluation_id = re.id
			LEFT JOIN risk_appetite_tolerance_levels ratl ON ratl.id = raa.risk_appetite_tolerance_level_id
		`).
		Where("raa.risk_appetite_statement_id = ?", statementID)

	// Apply filters, keyword, etc (copy same logic as in your GetAll)
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(risk_register.title LIKE ? OR
			risk_register.description LIKE ? OR
			risk_register.reference_number LIKE ?)
		`, keyword, keyword, keyword)
	}

	if params.SortBy == "" {
		params.SortBy = "risk_register.title"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}
	query = query.Order(params.SortBy + " " + params.Sort)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	if err := query.Find(&risks).Error; err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	return &commonDTOs.AssociatedRiskListResponseDTO{
		Items: risks,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}, nil
}
