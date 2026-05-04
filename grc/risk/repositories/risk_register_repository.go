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
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskRegister, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskRegister, error)
	Create(riskRegister *models.RiskRegister, tx ...utility.DatabaseTransaction) (uint, error)
	Update(riskRegister *models.RiskRegister, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	ForceDelete(id uint, tx ...utility.DatabaseTransaction) error
	GenerateReferenceNumber() (string, error)

	//Document
	CreateDocuments(docs []*models.RiskDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.RiskDocument, error)
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error

	//Associated Items
	CreateRiskSLAs(slaRegisters []*models.RiskSLA, tx ...utility.DatabaseTransaction) error
	DeleteRiskSLAs(slaRegisterIds []*uint, tx ...utility.DatabaseTransaction) error
	GetControlSummary(riskID uint, controlType string) ([]models.RiskControlAnalysis, error)

	GetAllRiskTitles() ([]string, error)
	GetAllIdentifiedRisks(params *riskDTOs.RiskRegisterListQueryParamsDTO) (*riskDTOs.RiskRegisterListResponseDTO, error)

	//
	DeleteByRiskResponseID(riskResponseID uint, tx ...utility.DatabaseTransaction) error
	DeleteByRiskResponseReviewCommentID(riskResponseReviewCommentID uint, tx ...utility.DatabaseTransaction) error
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
			CONVERT(VARCHAR, risk_register.reported_at, 23) AS reported_at,
            risk_register.created_at,
            risk_register.updated_at,
			risk_register.risk_code,
			CONVERT(VARCHAR, risk_register.registration_date, 23) AS registration_date,
			risk_register.is_strategic,
			risk_register.is_operational,
			risk_register.expected_cause,
			risk_register.consequence,
			risk_register.is_internal,
			risk_register.is_external,
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
			risk_statuses.type AS risk_status_type,
            review_frequencies.title_en AS review_frequency_title_en,
            review_frequencies.title_ar AS review_frequency_title_ar,
			severity_levels.title_en AS severity_level_title_en,
            severity_levels.title_ar AS severity_level_title_ar,
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
			divisions.id AS division_id,
            divisions.title_en AS division_title_en,
            divisions.title_ar AS division_title_ar,
			sections.id AS section_id,
            sections.title_en AS section_title_en,
            sections.title_ar AS section_title_ar

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
			LEFT JOIN severity_levels ON severity_levels.id = risk_register.severity_level_id
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
			LEFT JOIN organization_units AS divisions ON divisions.id = risk_register.division_id
			LEFT JOIN organization_units AS sections ON sections.id = risk_register.section_id
        `).Where("risk_statuses.type NOT IN (?)", []string{"identified", "submitted", "rejected"})

	if params.RiskAppetiteToleranceLevelID != 0 {
		query = query.Joins(`
				LEFT JOIN risk_evaluations ON risk_evaluations.risk_register_id = risk_register.id
				LEFT JOIN risk_appetite_alignments ON risk_appetite_alignments.risk_evaluation_id = risk_evaluations.id
			`).Where("risk_appetite_alignments.risk_appetite_tolerance_level_id = ?", params.RiskAppetiteToleranceLevelID)
	}
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
		query = query.Where("reported_at <= ?", params.ReportedAtTo)
	}

	if params.NextReviewDateFrom != "" {
		query = query.Where("next_review_date >= ?", params.NextReviewDateFrom)
	}

	if params.NextReviewDateTo != "" {
		query = query.Where("next_review_date <= ?", params.NextReviewDateTo)
	}

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("risk_register.reported_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	// Apply quick filters
	for _, filter := range params.QuickFilter {
		switch filter {
		case "open":
			query = query.Where("risk_statuses.type = ? ", "open")
		case "identified":
			query = query.Where("risk_statuses.type = ? ", "identified")
		case "critical_risk":
			// Fetch the highest available risk rating dynamically
			var highestRating string
			err := repo.db.Table("risk_matrix_risk_ratings").
				Select("type").
				Where("status = ?", "active").
				Order("CASE type WHEN 'very_high' THEN 1 WHEN 'high' THEN 2 WHEN 'medium' THEN 3 WHEN 'low' THEN 4 WHEN 'very_low' THEN 5 END").
				Limit(1).Pluck("type", &highestRating).Error
			if err != nil || highestRating == "" {
				return nil, fmt.Errorf("failed to fetch highest risk rating: %v", err)
			}
			query = query.Where("current_risk_rating.type = ?", highestRating)

		case "open_actions":
			query = query.Joins(`
				LEFT JOIN (
					SELECT DISTINCT risk_register_id
					FROM risk_action_plans
					JOIN action_plan_statuses ON action_plan_statuses.id = risk_action_plans.action_plan_status_id
					WHERE action_plan_statuses.is_closed = 0
				) AS open_actions ON open_actions.risk_register_id = risk_register.id
			`).
				Where("open_actions.risk_register_id IS NOT NULL")

		case "newly_identified":
			twoWeeksAgo := time.Now().AddDate(0, 0, -14).Format("2006-01-02 00:00:00")
			query = query.Where("risk_register.created_at >= ?", twoWeeksAgo)

		case "strategic_risks":
			query = query.Joins(`
				LEFT JOIN (
					SELECT DISTINCT risk_register_id
					FROM risk_strategic_objectives
				) AS strategic_objectives ON strategic_objectives.risk_register_id = risk_register.id
			`).Where("strategic_objectives.risk_register_id IS NOT NULL")

		case "high_financial_impact":
			query = query.Where("risk_categories.title_en = ? ", "Financial Risk")

		case "closed_resolved":
			query = query.Where("risk_statuses.type IN (?)", []string{"closed", "resolved"})
		case "open_close":
			query = query.Where("risk_statuses.type IN (?)", []string{"closed", "open"})
		case "kri_enabled":
			query = query.Joins(`
				LEFT JOIN (
					SELECT DISTINCT risk_register_id
					FROM risk_kris
				) AS kri_enabled ON kri_enabled.risk_register_id = risk_register.id
			`).Where("kri_enabled.risk_register_id IS NOT NULL")

		case "kri_not_enabled":
			query = query.Joins(`
				LEFT JOIN (
					SELECT DISTINCT risk_register_id
					FROM risk_kris
				) AS kri_risks ON kri_risks.risk_register_id = risk_register.id
			`).Where("kri_risks.risk_register_id IS NULL")

		}

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
			organization_units.title_en LIKE ? OR
			organization_units.title_ar LIKE ? OR
			current_risk_rating.title_en LIKE ? OR
			current_risk_rating.title_ar LIKE ? OR
			current_impact.title_en LIKE ? OR
			current_impact.title_ar LIKE ? OR
			current_likelihood.title_en LIKE ? OR
			current_likelihood.title_ar LIKE ? OR
			reported_at LIKE ? OR
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
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	if params.Value != 0 {
		query = query.Where("current_risk_rating_value.value = ?", params.Value)
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	// Count total results
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply pagination only if not exporting
	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	if err := query.Find(&riskRegisters).Error; err != nil {
		return nil, err
	}

	if params.Export {
		return &riskDTOs.RiskRegisterListResponseDTO{
			Items: riskRegisters,
		}, nil
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

func (repo *RiskRegisterRepository) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskRegister, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var risk *models.RiskRegister

	result := db.
		First(&risk, id)
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
		Preload("ResponsibleUnit").
		Preload("ResponsibleUser").
		Preload("ReviewFrequency").
		Preload("RiskStatus").
		Preload("Documents").
		Preload("SLARegisters").
		Preload("KRIs").
		Preload("Assets").
		Preload("Controls").
		Preload("Controls.Control").
		Preload("Controls.Control.ControlEffectivenessRating").
		Preload("CurrentLikelihood").
		Preload("CurrentImpact").
		Preload("CurrentRiskRating").
		Preload("CurrentRiskRatingValue").
		Preload("InherentLikelihood").
		Preload("InherentImpact").
		Preload("InherentRiskRating").
		Preload("InherentRiskRatingValue").
		Preload("ResidualLikelihood").
		Preload("ResidualImpact").
		Preload("ResidualRiskRating").
		Preload("ResidualRiskRatingValue").
		Preload("RiskResponseStrategy").
		Preload("SeverityLevel").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("SubmittedByUser").
		Preload("Division").
		Preload("Section").
		Preload("ReviewedByUser").
		Preload("RiskWorkflows").
		Preload("RiskWorkflows.RiskWorkflowParticipants").
		Preload("RiskWorkflows.RiskWorkflowReviews").
		Preload("AuditableItems", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC").Limit(5) // Fetch the latest 5 auditable items
		}).
		Preload("AuditableItems.ResponsibleUser").
		Preload("AuditableItems.ResponsibleUnit").
		Preload("AuditableItems.AuditorRiskRating").
		Preload("AuditableItems.CreatedByUser").
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
		Where("reference_number LIKE ?", "IR-%").
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

func (repo *RiskRegisterRepository) CreateRiskSLAs(slaRegisters []*models.RiskSLA, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&slaRegisters).Error; err != nil {
		return err
	}
	return nil
}

func (repo *RiskRegisterRepository) DeleteRiskSLAs(slaRegisterIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete SLARegister by IDs
	if err := db.Where("id IN (?)", slaRegisterIds).Delete(&models.RiskSLA{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskRegisterRepository) GetControlSummary(riskID uint, controlType string) ([]models.RiskControlAnalysis, error) {
	var riskControlAnalyses []models.RiskControlAnalysis

	result := repo.db.
		Preload("ControlEfficiencyStatus").
		Preload("ControlImplementationStatus").
		Joins("JOIN risk_controls ON risk_controls.id = risk_control_analyses.risk_control_id").
		Where("risk_controls.risk_register_id = ?", riskID).
		Where("risk_control_analyses.type = ?", controlType).
		Find(&riskControlAnalyses)
	if result.Error != nil {
		return nil, result.Error
	}

	return riskControlAnalyses, nil
}

func (repo *RiskRegisterRepository) GetAllRiskTitles() ([]string, error) {
	var titles []string
	err := repo.db.
		Table("risk_register").
		Pluck("title", &titles).Error
	return titles, err
}

func (repo *RiskRegisterRepository) GetAllIdentifiedRisks(params *riskDTOs.RiskRegisterListQueryParamsDTO) (*riskDTOs.RiskRegisterListResponseDTO, error) {
	var riskRegisters []*riskDTOs.RiskRegisterListDTO
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
			risk_register.risk_code,
			CONVERT(VARCHAR, risk_register.registration_date, 23) AS registration_date,
			risk_register.is_strategic,
			risk_register.is_operational,
			risk_register.expected_cause,
			risk_register.consequence,
			risk_register.is_internal,
			risk_register.is_external,
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
			risk_statuses.type AS risk_status_type,
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
			divisions.id AS division_id,
            divisions.title_en AS division_title_en,
            divisions.title_ar AS division_title_ar,
						sections.id AS section_id,
            sections.title_en AS section_title_en,
            sections.title_ar AS section_title_ar
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
			LEFT JOIN organization_units AS divisions ON divisions.id = risk_register.division_id
						 LEFT JOIN organization_units AS sections ON sections.id = risk_register.section_id
        `).Where("risk_statuses.type IN (?)", []string{"identified", "submitted", "rejected"})

	if params.RiskAppetiteToleranceLevelID != 0 {
		query = query.Joins(`
				LEFT JOIN risk_evaluations ON risk_evaluations.risk_register_id = risk_register.id
				LEFT JOIN risk_appetite_alignments ON risk_appetite_alignments.risk_evaluation_id = risk_evaluations.id
			`).Where("risk_appetite_alignments.risk_appetite_tolerance_level_id = ?", params.RiskAppetiteToleranceLevelID)
	}
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
		query = query.Where("reported_at <= ?", params.ReportedAtTo)
	}

	if params.NextReviewDateFrom != "" {
		query = query.Where("next_review_date >= ?", params.NextReviewDateFrom)
	}

	if params.NextReviewDateTo != "" {
		query = query.Where("next_review_date <= ?", params.NextReviewDateTo)
	}

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("risk_register.reported_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	// Apply quick filters
	for _, filter := range params.QuickFilter {
		switch filter {
		case "open":
			query = query.Where("risk_statuses.type = ? ", "open")
		case "identified":
			query = query.Where("risk_statuses.type = ? ", "identified")
		case "critical_risk":
			// Fetch the highest available risk rating dynamically
			var highestRating string
			err := repo.db.Table("risk_matrix_risk_ratings").
				Select("type").
				Where("status = ?", "active").
				Order("CASE type WHEN 'very_high' THEN 1 WHEN 'high' THEN 2 WHEN 'medium' THEN 3 WHEN 'low' THEN 4 WHEN 'very_low' THEN 5 END").
				Limit(1).Pluck("type", &highestRating).Error
			if err != nil || highestRating == "" {
				return nil, fmt.Errorf("failed to fetch highest risk rating: %v", err)
			}
			query = query.Where("current_risk_rating.type = ?", highestRating)

		case "open_actions":
			query = query.Joins(`
				LEFT JOIN (
					SELECT DISTINCT risk_register_id
					FROM risk_action_plans
					JOIN action_plan_statuses ON action_plan_statuses.id = risk_action_plans.action_plan_status_id
					WHERE action_plan_statuses.is_closed = 0
				) AS open_actions ON open_actions.risk_register_id = risk_register.id
			`).
				Where("open_actions.risk_register_id IS NOT NULL")

		case "newly_identified":
			twoWeeksAgo := time.Now().AddDate(0, 0, -14).Format("2006-01-02 00:00:00")
			query = query.Where("risk_register.created_at >= ?", twoWeeksAgo)

		case "strategic_risks":
			query = query.Joins(`
				LEFT JOIN (
					SELECT DISTINCT risk_register_id
					FROM risk_strategic_objectives
				) AS strategic_objectives ON strategic_objectives.risk_register_id = risk_register.id
			`).Where("strategic_objectives.risk_register_id IS NOT NULL")

		case "high_financial_impact":
			query = query.Where("risk_categories.title_en = ? ", "Financial Risk")

		case "closed_resolved":
			query = query.Where("risk_statuses.type IN (?)", []string{"closed", "resolved"})
		case "kri_enabled":
			query = query.Joins(`
				LEFT JOIN (
					SELECT DISTINCT risk_register_id
					FROM risk_kris
				) AS kri_enabled ON kri_enabled.risk_register_id = risk_register.id
			`).Where("kri_enabled.risk_register_id IS NOT NULL")

		case "kri_not_enabled":
			query = query.Joins(`
				LEFT JOIN (
					SELECT DISTINCT risk_register_id
					FROM risk_kris
				) AS kri_risks ON kri_risks.risk_register_id = risk_register.id
			`).Where("kri_risks.risk_register_id IS NULL")

		}

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
			organization_units.title_en LIKE ? OR
			organization_units.title_ar LIKE ? OR
			current_risk_rating.title_en LIKE ? OR
			current_risk_rating.title_ar LIKE ? OR
			current_impact.title_en LIKE ? OR
			current_impact.title_ar LIKE ? OR
			current_likelihood.title_en LIKE ? OR
			current_likelihood.title_ar LIKE ? OR
			reported_at LIKE ? OR
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
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	if params.Value != 0 {
		query = query.Where("current_risk_rating_value.value = ?", params.Value)
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	// Count total results
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply pagination only if not exporting
	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	if err := query.Debug().Find(&riskRegisters).Error; err != nil {
		return nil, err
	}

	if params.Export {
		return &riskDTOs.RiskRegisterListResponseDTO{
			Items: riskRegisters,
		}, nil
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

func (repo *RiskRegisterRepository) ForceDelete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	queries := []string{

		// --- Deepest children first ---
		`DELETE audit_plan_auditable_item_risk_controls
		FROM audit_plan_auditable_item_risk_controls
		JOIN audit_plan_auditable_item_risks
		ON audit_plan_auditable_item_risks.id = audit_plan_auditable_item_risk_controls.audit_plan_auditable_item_risk_id
		WHERE audit_plan_auditable_item_risks.risk_register_id = ?`,

		`DELETE audit_plan_auditable_item_risks
		FROM audit_plan_auditable_item_risks
		WHERE audit_plan_auditable_item_risks.risk_register_id = ?`,

		`DELETE audit_finding_risks
		FROM audit_finding_risks
		WHERE audit_finding_risks.risk_register_id = ?`,

		`DELETE audit_finding_action_plan_status_update_documents
		FROM audit_finding_action_plan_status_update_documents
		JOIN audit_finding_action_plan_status_updates ON audit_finding_action_plan_status_updates.id = audit_finding_action_plan_status_update_documents.audit_finding_action_plan_status_update_id
		JOIN audit_finding_action_plans ON audit_finding_action_plans.id = audit_finding_action_plan_status_updates.audit_finding_action_plan_id
		JOIN audit_findings ON audit_findings.id = audit_finding_action_plans.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.risk_register_id = ?`,

		`DELETE audit_finding_action_plan_status_updates
		FROM audit_finding_action_plan_status_updates
		JOIN audit_finding_action_plans ON audit_finding_action_plans.id = audit_finding_action_plan_status_updates.audit_finding_action_plan_id
		JOIN audit_findings ON audit_findings.id = audit_finding_action_plans.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.risk_register_id = ?`,

		`DELETE audit_finding_action_plan_comment_documents
		FROM audit_finding_action_plan_comment_documents
		JOIN audit_finding_action_plan_comments ON audit_finding_action_plan_comments.id = audit_finding_action_plan_comment_documents.audit_finding_action_plan_comment_id
		JOIN audit_finding_action_plans ON audit_finding_action_plans.id = audit_finding_action_plan_comments.audit_finding_action_plan_id
		JOIN audit_findings ON audit_findings.id = audit_finding_action_plans.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.risk_register_id = ?`,

		`DELETE audit_finding_action_plan_comments
		FROM audit_finding_action_plan_comments
		JOIN audit_finding_action_plans ON audit_finding_action_plans.id = audit_finding_action_plan_comments.audit_finding_action_plan_id
		JOIN audit_findings ON audit_findings.id = audit_finding_action_plans.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.risk_register_id = ?`,

		`DELETE tasks
		FROM tasks
		JOIN audit_finding_action_plans ON audit_finding_action_plans.id = tasks.audit_finding_action_plan_id
		JOIN audit_findings ON audit_findings.id = audit_finding_action_plans.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.risk_register_id = ?`,

		`DELETE audit_finding_action_plans
		FROM audit_finding_action_plans
		JOIN audit_findings ON audit_findings.id = audit_finding_action_plans.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.risk_register_id = ?`,

		`DELETE audit_finding_assets
		FROM audit_finding_assets
		JOIN audit_findings ON audit_findings.id = audit_finding_assets.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.risk_register_id = ?`,

		`DELETE audit_finding_associated_documents
		FROM audit_finding_associated_documents
		JOIN audit_findings ON audit_findings.id = audit_finding_associated_documents.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.risk_register_id = ?`,

		`DELETE audit_finding_compliance_requirements
		FROM audit_finding_compliance_requirements
		JOIN audit_findings ON audit_findings.id = audit_finding_compliance_requirements.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.risk_register_id = ?`,

		`DELETE audit_finding_controls
		FROM audit_finding_controls
		JOIN audit_findings ON audit_findings.id = audit_finding_controls.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.risk_register_id = ?`,

		`DELETE audit_finding_documents
		FROM audit_finding_documents
		JOIN audit_findings ON audit_findings.id = audit_finding_documents.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.risk_register_id = ?`,

		`DELETE audit_finding_incidents
		FROM audit_finding_incidents
		JOIN audit_findings ON audit_findings.id = audit_finding_incidents.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.risk_register_id = ?`,

		`DELETE audit_finding_kpis
		FROM audit_finding_kpis
		JOIN audit_findings ON audit_findings.id = audit_finding_kpis.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.risk_register_id = ?`,

		`DELETE audit_finding_processes
		FROM audit_finding_processes
		JOIN audit_findings ON audit_findings.id = audit_finding_processes.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.risk_register_id = ?`,

		`DELETE audit_finding_risks
		FROM audit_finding_risks
		JOIN audit_findings ON audit_findings.id = audit_finding_risks.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.risk_register_id = ?`,

		`DELETE audit_finding_root_cause_analyses
		FROM audit_finding_root_cause_analyses
		JOIN audit_findings ON audit_findings.id = audit_finding_root_cause_analyses.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.risk_register_id = ?`,

		`DELETE audit_finding_strategic_objectives
		FROM audit_finding_strategic_objectives
		JOIN audit_findings ON audit_findings.id = audit_finding_strategic_objectives.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.risk_register_id = ?`,

		`DELETE audit_findings
		FROM audit_findings
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.risk_register_id = ?`,

		`DELETE audit_draft_finding_conversation_documents
		FROM audit_draft_finding_conversation_documents
		JOIN audit_draft_finding_conversations ON audit_draft_finding_conversations.id = audit_draft_finding_conversation_documents.audit_draft_finding_conversation_id
		JOIN audit_draft_findings ON audit_draft_findings.id = audit_draft_finding_conversations.audit_draft_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_draft_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.risk_register_id = ?`,

		`DELETE audit_draft_finding_conversations
		FROM audit_draft_finding_conversations
		JOIN audit_draft_findings ON audit_draft_findings.id = audit_draft_finding_conversations.audit_draft_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_draft_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.risk_register_id = ?`,

		`DELETE audit_draft_finding_documents
		FROM audit_draft_finding_documents
		JOIN audit_draft_findings ON audit_draft_findings.id = audit_draft_finding_documents.audit_draft_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_draft_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.risk_register_id = ?`,

		`DELETE audit_draft_finding_action_plans
		FROM audit_draft_finding_action_plans
		JOIN audit_draft_findings ON audit_draft_findings.id = audit_draft_finding_action_plans.audit_draft_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_draft_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.risk_register_id = ?`,

		`DELETE audit_draft_findings
		FROM audit_draft_findings
		JOIN audit_checklists ON audit_checklists.id = audit_draft_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.risk_register_id = ?`,

		`DELETE audit_checklist_documents
		FROM audit_checklist_documents
		JOIN audit_checklists ON audit_checklists.id = audit_checklist_documents.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.risk_register_id = ?`,

		`DELETE audit_checklists
		FROM audit_checklists
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.risk_register_id = ?`,

		`DELETE audit_plan_checklists
		FROM audit_plan_checklists
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.risk_register_id = ?`,

		`DELETE audit_plan_auditable_item_risk_controls
		FROM audit_plan_auditable_item_risk_controls
		JOIN audit_plan_auditable_item_risks ON audit_plan_auditable_item_risks.id = audit_plan_auditable_item_risk_controls.audit_plan_auditable_item_risk_id
		JOIN audit_plan_auditable_items ON audit_plan_auditable_items.id = audit_plan_auditable_item_risks.audit_plan_auditable_item_id
		JOIN auditable_items ON auditable_items.id = audit_plan_auditable_items.auditable_item_id
		WHERE auditable_items.risk_register_id = ?`,

		`DELETE audit_plan_auditable_item_risks
		FROM audit_plan_auditable_item_risks
		JOIN audit_plan_auditable_items ON audit_plan_auditable_items.id = audit_plan_auditable_item_risks.audit_plan_auditable_item_id
		JOIN auditable_items ON auditable_items.id = audit_plan_auditable_items.auditable_item_id
		WHERE auditable_items.risk_register_id = ?`,

		`DELETE audit_plan_auditable_items
		FROM audit_plan_auditable_items
		JOIN auditable_items ON auditable_items.id = audit_plan_auditable_items.auditable_item_id
		WHERE auditable_items.risk_register_id = ?`,

		`DELETE audit_schedule_auditable_items
		FROM audit_schedule_auditable_items
		JOIN auditable_items ON auditable_items.id = audit_schedule_auditable_items.auditable_item_id
		WHERE auditable_items.risk_register_id = ?`,

		`DELETE audit_plan_schedule_auditable_items
		FROM audit_plan_schedule_auditable_items
		JOIN auditable_items ON auditable_items.id = audit_plan_schedule_auditable_items.auditable_item_id
		WHERE auditable_items.risk_register_id = ?`,

		`DELETE audit_document_request_auditable_items
		FROM audit_document_request_auditable_items
		JOIN auditable_items ON auditable_items.id = audit_document_request_auditable_items.auditable_item_id
		WHERE auditable_items.risk_register_id = ?`,

		`DELETE risk_workflow_reviews
		FROM risk_workflow_reviews
		JOIN risk_workflows ON risk_workflows.id = risk_workflow_reviews.risk_workflow_id
		WHERE risk_workflows.risk_register_id = ?`,

		`DELETE risk_workflow_participants
		FROM risk_workflow_participants
		JOIN risk_workflows ON risk_workflows.id = risk_workflow_participants.risk_workflow_id
		WHERE risk_workflows.risk_register_id = ?`,

		`DELETE risk_impact_analysis_guidelines
		FROM risk_impact_analysis_guidelines
		JOIN risk_impact_analyses ON risk_impact_analyses.id = risk_impact_analysis_guidelines.risk_impact_analysis_id
		JOIN risk_evaluations ON risk_evaluations.id = risk_impact_analyses.risk_evaluation_id
		WHERE risk_evaluations.risk_register_id = ?`,

		`DELETE risk_impact_analyses
		FROM risk_impact_analyses
		JOIN risk_evaluations ON risk_evaluations.id = risk_impact_analyses.risk_evaluation_id
		WHERE risk_evaluations.risk_register_id = ?`,

		`DELETE risk_appetite_alignments
		FROM risk_appetite_alignments
		JOIN risk_evaluations ON risk_evaluations.id = risk_appetite_alignments.risk_evaluation_id
		WHERE risk_evaluations.risk_register_id = ?`,

		`DELETE risk_action_plan_status_update_documents
		FROM risk_action_plan_status_update_documents
		JOIN risk_action_plan_status_updates ON risk_action_plan_status_updates.id = risk_action_plan_status_update_documents.risk_action_plan_status_update_id
		JOIN risk_action_plans ON risk_action_plans.id = risk_action_plan_status_updates.risk_action_plan_id
		WHERE risk_action_plans.risk_register_id = ?`,

		`DELETE risk_action_plan_status_updates
		FROM risk_action_plan_status_updates
		JOIN risk_action_plans ON risk_action_plans.id = risk_action_plan_status_updates.risk_action_plan_id
		WHERE risk_action_plans.risk_register_id = ?`,

		`DELETE risk_action_plan_controls
		FROM risk_action_plan_controls
		JOIN risk_action_plans ON risk_action_plans.id = risk_action_plan_controls.risk_action_plan_id
		WHERE risk_action_plans.risk_register_id = ?`,

		`DELETE risk_action_plan_comment_documents
		FROM risk_action_plan_comment_documents
		JOIN risk_action_plan_comments ON risk_action_plan_comments.id = risk_action_plan_comment_documents.risk_action_plan_comment_id
		JOIN risk_action_plans ON risk_action_plans.id = risk_action_plan_comments.risk_action_plan_id
		WHERE risk_action_plans.risk_register_id = ?`,

		`DELETE risk_action_plan_comments
		FROM risk_action_plan_comments
		JOIN risk_action_plans ON risk_action_plans.id = risk_action_plan_comments.risk_action_plan_id
		WHERE risk_action_plans.risk_register_id = ?`,

		`DELETE risk_control_analyses
		FROM risk_control_analyses
		JOIN risk_controls ON risk_controls.id = risk_control_analyses.risk_control_id
		WHERE risk_controls.risk_register_id = ?`,

		`DELETE control_test_documents
		FROM control_test_documents
		JOIN control_tests ON control_tests.id = control_test_documents.control_test_id
		WHERE control_tests.risk_register_id = ?`,

		`DELETE risk_closure_documents
		FROM risk_closure_documents
		JOIN risk_closures ON risk_closures.id = risk_closure_documents.risk_closure_id
		WHERE risk_closures.risk_register_id = ?`,

		`DELETE FROM third_party_risks
		WHERE third_party_risks.risk_register_id = ?`,

		`DELETE FROM risk_workflows
		WHERE risk_workflows.risk_register_id = ?`,

		`DELETE FROM risk_timelines
		WHERE risk_timelines.risk_register_id = ?`,

		`DELETE FROM risk_threats
		WHERE risk_threats.risk_register_id = ?`,

		`DELETE FROM risk_strategic_objectives
		WHERE risk_strategic_objectives.risk_register_id = ?`,

		`DELETE FROM risk_slas
		WHERE risk_slas.risk_register_id = ?`,

		`DELETE FROM risk_root_cause_analysis
		WHERE risk_root_cause_analysis.risk_register_id = ?`,

		`DELETE FROM risk_review_comments
		WHERE risk_review_comments.risk_register_id = ?`,

		`DELETE FROM risk_processes
		WHERE risk_processes.risk_register_id = ?`,

		`DELETE FROM risk_kris
		WHERE risk_kris.risk_register_id = ?`,

		`DELETE FROM risk_evaluations
		WHERE risk_evaluations.risk_register_id = ?`,

		`DELETE FROM risk_documents
		WHERE risk_documents.risk_register_id = ?`,

		`DELETE FROM risk_controls
		WHERE risk_controls.risk_register_id = ?`,

		`DELETE FROM risk_compliance_requirements
		WHERE risk_compliance_requirements.risk_register_id = ?`,

		`DELETE FROM risk_comments
		WHERE risk_comments.risk_register_id = ?`,

		`DELETE FROM risk_closures
		WHERE risk_closures.risk_register_id = ?`,

		`DELETE FROM risk_associated_documents
		WHERE risk_associated_documents.risk_register_id = ?`,

		`DELETE FROM risk_assets
		WHERE risk_assets.risk_register_id = ?`,

		`DELETE FROM risk_action_plans
		WHERE risk_action_plans.risk_register_id = ?`,

		`DELETE FROM project_risks
		WHERE project_risks.risk_register_id = ?`,

		`DELETE FROM policy_exception_risks
		WHERE policy_exception_risks.risk_register_id = ?`,

		`DELETE FROM meeting_risks
		WHERE meeting_risks.risk_register_id = ?`,

		`DELETE FROM kpi_risks
		WHERE kpi_risks.risk_register_id = ?`,

		`DELETE FROM incident_risks
		WHERE incident_risks.risk_register_id = ?`,

		`DELETE FROM control_tests
		WHERE control_tests.risk_register_id = ?`,

		`DELETE FROM bcp_risks
		WHERE bcp_risks.risk_register_id = ?`,

		`DELETE FROM auditable_items
		WHERE auditable_items.risk_register_id = ?`,
	}
	// Delete the Risk register

	return db.Transaction(func(tx *gorm.DB) error {
		for _, q := range queries {
			if err := tx.Exec(q, id).Error; err != nil {
				return err
			}
		}
		return tx.Delete(&models.RiskRegister{}, id).Error
	})

}

func (repo *RiskRegisterRepository) DeleteByRiskResponseID(riskResponseID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Get workflow IDs for this risk response
	var workflowIDs []uint
	if err := db.Model(&models.RiskResponseWorkflow{}).Where("risk_response_id = ?", riskResponseID).Pluck("id", &workflowIDs).Error; err != nil {
		return err
	}

	if len(workflowIDs) == 0 {
		return nil
	}

	// Delete workflow participants
	if err := db.Where("risk_response_workflow_id IN (?)", workflowIDs).Delete(&models.RiskResponseWorkflowParticipant{}).Error; err != nil {
		return err
	}

	// Delete workflow reviews
	if err := db.Where("risk_response_workflow_id IN (?)", workflowIDs).Delete(&models.RiskResponseWorkflowReview{}).Error; err != nil {
		return err
	}

	// Delete workflows
	if err := db.Where("id IN (?)", workflowIDs).Delete(&models.RiskResponseWorkflow{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskRegisterRepository) DeleteByRiskResponseReviewCommentID(riskResponseID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Get comment IDs for this risk response
	var commentIDs []uint
	if err := db.Model(&models.RiskResponseReviewComment{}).Where("risk_response_id = ?", riskResponseID).Pluck("id", &commentIDs).Error; err != nil {
		return err
	}

	if len(commentIDs) == 0 {
		return nil
	}

	// Delete comment documents
	if err := db.Where("risk_response_review_comment_id IN (?)", commentIDs).Delete(&models.RiskResponseReviewCommentDocument{}).Error; err != nil {
		return err
	}

	// Delete comments
	if err := db.Where("id IN (?)", commentIDs).Delete(&models.RiskResponseReviewComment{}).Error; err != nil {
		return err
	}

	return nil
}
