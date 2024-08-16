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
	GenerateReferenceNumber() (string, error)

	//Document
	CreateDocuments(docs []*models.RiskDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.RiskDocument, error)
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error

	//Associated Items
	CreateRiskComplianceRequirements(complianceRequirements []*models.RiskComplianceRequirement, tx ...utility.DatabaseTransaction) error
	CreateRiskControls(controls []*models.RiskControl, tx ...utility.DatabaseTransaction) error
	DeleteRiskControls(controlIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteRiskComplianceRequirements(requirementIds []*uint, tx ...utility.DatabaseTransaction) error
	CreateRiskProcesses(processes []*models.RiskProcess, tx ...utility.DatabaseTransaction) error
	DeleteRiskProcesses(processIds []*uint, tx ...utility.DatabaseTransaction) error
	CreateRiskStrategicObjectives(strategicObjectives []*models.RiskStrategicObjective, tx ...utility.DatabaseTransaction) error
	DeleteRiskStrategicObjectives(strategicObjectiveIds []*uint, tx ...utility.DatabaseTransaction) error
	CreateRiskSLAs(slaRegisters []*models.RiskSLA, tx ...utility.DatabaseTransaction) error
	DeleteRiskSLAs(slaRegisterIds []*uint, tx ...utility.DatabaseTransaction) error

	GetControlSummary(riskID uint, controlType string) ([]models.RiskControlAnalysis, error)
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
            updated_user.id AS updated_user_id,
    		updated_user.first_name_en AS updated_user_first_name_en,
            updated_user.first_name_ar AS updated_user_first_name_ar,
            updated_user.last_name_en AS updated_user_last_name_en,
            updated_user.last_name_ar AS updated_user_last_name_ar,
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
            residual_risk_rating.color_code AS residual_risk_rating_color_code
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
		query = query.Where("reported_at <= ?", params.ReportedAtTo)
	}

	if params.NextReviewDateFrom != "" {
		query = query.Where("next_review_date >= ?", params.NextReviewDateFrom)
	}

	if params.NextReviewDateTo != "" {
		query = query.Where("next_review_date <= ?", params.NextReviewDateTo)
	}

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
			updated_user.id AS updated_user_id,
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

	if params.Value != 0 {
		query = query.Where("current_risk_rating_value.value = ?", params.Value)
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

	result := repo.db.
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
		Preload("Controls").
		Preload("Controls.Control.ControlCategory").
		Preload("Controls.Control.ControlType").
		Preload("ComplianceRequirements").
		Preload("ComplianceRequirements.ComplianceRequirement").
		Preload("Processes").
		Preload("StrategicObjectives").
		Preload("SLARegisters").
		Preload("KRIs").
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

func (repo *RiskRegisterRepository) CreateRiskComplianceRequirements(complianceRequirements []*models.RiskComplianceRequirement, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&complianceRequirements).Error; err != nil {
		return err
	}
	return nil
}

func (repo *RiskRegisterRepository) CreateRiskControls(controls []*models.RiskControl, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&controls).Error; err != nil {
		return err
	}
	return nil
}

func (repo *RiskRegisterRepository) DeleteRiskComplianceRequirements(requirementIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", requirementIds).Delete(&models.RiskComplianceRequirement{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskRegisterRepository) DeleteRiskControls(controlIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", controlIds).Delete(&models.RiskControl{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskRegisterRepository) CreateRiskProcesses(processes []*models.RiskProcess, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&processes).Error; err != nil {
		return err
	}
	return nil
}

func (repo *RiskRegisterRepository) DeleteRiskProcesses(processIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete processes by IDs
	if err := db.Where("id IN (?)", processIds).Delete(&models.RiskProcess{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskRegisterRepository) CreateRiskStrategicObjectives(strategicObjectives []*models.RiskStrategicObjective, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&strategicObjectives).Error; err != nil {
		return err
	}
	return nil
}

func (repo *RiskRegisterRepository) DeleteRiskStrategicObjectives(strategicObjectiveIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete strategicObjectives by IDs
	if err := db.Where("id IN (?)", strategicObjectiveIds).Delete(&models.RiskStrategicObjective{}).Error; err != nil {
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
