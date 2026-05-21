package repositories

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/utility"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type KRIRegisterRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *riskDTOs.KRIRegisterListQueryParamsDTO) (*riskDTOs.KRIRegisterListResponseDTO, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.KRIRegister, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.KRIRegister, error)
	Create(kriRegister *models.KRIRegister, tx ...utility.DatabaseTransaction) (uint, error)
	Update(kriRegister *models.KRIRegister, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	ForceDelete(id uint, tx ...utility.DatabaseTransaction) error
	ExistsByTitle(Title string, KRICategoryID uint) (bool, error)
	ExistsByTitleExcludingID(id uint, Title string, KRICategoryID uint) (bool, error)
	GenerateReferenceNumber() (string, error)
	CreateInitialKriValueUpdate(kriValueUpdate *models.KRIValueUpdate, tx ...utility.DatabaseTransaction) (uint, error)
	DeleteKriValueUpdates(ids []*uint, tx ...utility.DatabaseTransaction) error
	GetTrendAnalysis(kriId uint, limit int) ([]*models.KRIValueUpdate, error)
	GetFrequencyTrendAnalysis(kriId uint, limit int) ([]*models.KRIFrequencyValue, error)
	CreateFrequencyValues(frequencyValues []*models.KRIFrequencyValue, tx ...utility.DatabaseTransaction) error
	DeleteFrequencyValuesByKRIID(kriRegisterID uint, tx ...utility.DatabaseTransaction) error
	GetFrequencyValuesByKRIID(kriRegisterID uint, tx ...utility.DatabaseTransaction) ([]*models.KRIFrequencyValue, error)
	GetAllFrequencyValues(kriRegisterID uint) ([]*riskDTOs.KRIFrequencyValueListDTO, error)
}

type KRIRegisterRepository struct {
	db *gorm.DB
}

func NewKRIRegisterRepository(db *gorm.DB) *KRIRegisterRepository {
	return &KRIRegisterRepository{
		db: db,
	}
}

func (repo *KRIRegisterRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *KRIRegisterRepository) GetAll(params *riskDTOs.KRIRegisterListQueryParamsDTO) (*riskDTOs.KRIRegisterListResponseDTO, error) {
	var kriRegisters []*riskDTOs.KRIRegisterListDTO
	var totalCount int64

	query := repo.db.Table("kri_register").
		Select(`
			kri_register.*,
			CONVERT(VARCHAR, kri_register.last_reviewed_date, 23) AS last_reviewed_date,
			CONVERT(VARCHAR, kri_register.next_review_date, 23) AS next_review_date,
			kri_categories.title_en AS kri_category_title_en,
			kri_categories.title_ar AS kri_category_title_ar,
			organization_units.title_en AS responsible_unit_title_en,
			organization_units.title_ar AS responsible_unit_title_ar,
			responsible_user.id AS responsible_user_id,
			responsible_user.first_name_en AS responsible_user_first_name_en,
			responsible_user.first_name_ar AS responsible_user_first_name_ar,
			responsible_user.last_name_en AS responsible_user_last_name_en,
			responsible_user.last_name_ar AS responsible_user_last_name_ar,
			responsible_user.initial AS responsible_user_initial,
			responsible_user.color AS responsible_user_color,
			responsible_user.email AS responsible_user_email,
			responsible_user.profile_image_path AS responsible_user_profile_image_path,
			review_frequencies.title_en AS review_frequency_title_en,
			review_frequencies.title_ar AS review_frequency_title_ar,
			risk_levels.title_en AS risk_level_title_en,
			risk_levels.title_ar AS risk_level_title_ar,
			risk_levels.color_code AS risk_level_color_code,
			units.title_en AS unit_title_en,
			units.title_ar AS unit_title_ar,
			units.code_en AS unit_code_en,
			units.code_ar AS unit_code_ar,
			units.color_code AS unit_color_code,
			level_of_measurements.title_en AS level_of_measurement_title_en,
			level_of_measurements.title_ar AS level_of_measurement_title_ar,
			risk_responsible_units.title_en AS risk_responsible_unit_title_en,
			risk_responsible_units.title_ar AS risk_responsible_unit_title_ar,
			kri_sub_categories.title_en AS kri_sub_category_title_en,
			kri_sub_categories.title_ar AS kri_sub_category_title_ar,
			distribution_statuses.title_en AS distribution_status_title_en,
			distribution_statuses.title_ar AS distribution_status_title_ar,
			kri_register.is_lower_better,
			kri_register.is_higher_better,
			workflow_statuses.id AS workflow_status_id,
			workflow_statuses.title_en AS workflow_status_title_en,
			workflow_statuses.title_ar AS workflow_status_title_ar,
			workflow_statuses.type AS workflow_status_type,
			workflow_statuses.color_code AS workflow_status_color_code,
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
			LEFT JOIN kri_categories ON kri_categories.id = kri_register.kri_category_id
			LEFT JOIN organization_units ON organization_units.id = kri_register.responsible_unit_id
			LEFT JOIN users AS responsible_user ON responsible_user.id = kri_register.responsible_user_id
			LEFT JOIN review_frequencies ON review_frequencies.id = kri_register.review_frequency_id
			LEFT JOIN risk_levels ON risk_levels.id = kri_register.risk_level_id
			LEFT JOIN units ON units.id = kri_register.unit_id
			LEFT JOIN level_of_measurements ON level_of_measurements.id = kri_register.level_of_measurement_id
			LEFT JOIN organization_units AS risk_responsible_units ON risk_responsible_units.id = kri_register.risk_responsible_unit_id
			LEFT JOIN kri_sub_categories ON kri_sub_categories.id = kri_register.kri_sub_category_id
			LEFT JOIN distribution_statuses ON distribution_statuses.id = kri_register.distribution_status_id
			LEFT JOIN workflow_statuses ON workflow_statuses.id = kri_register.workflow_status_id
			LEFT JOIN users AS created_user ON created_user.id = kri_register.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = kri_register.updated_by_id
		`)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply LastReviewDate filters
	if params.LastReviewedDateFrom != "" {
		query = query.Where("kri_register.last_reviewed_date >= ?", params.LastReviewedDateFrom)
	}
	if params.LastReviewedDateTo != "" {
		query = query.Where("kri_register.last_reviewed_date <= ?", params.LastReviewedDateTo)
	}

	// Apply next_review_date filters
	if params.NextReviewDateFrom != "" {
		query = query.Where("kri_register.next_review_date >= ?", params.NextReviewDateFrom)
	}
	if params.NextReviewDateTo != "" {
		query = query.Where("kri_register.next_review_date <= ?", params.NextReviewDateTo)
	}

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("kri_register.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}
	// Apply quick filters
	for _, filter := range params.QuickFilter {
		switch filter {
		case "high_risk_level":
			query = query.Where("risk_levels.type = ?", "high")
		case "medium_risk_level":
			query = query.Where("risk_levels.type = ?", "medium")
		case "low_risk_level":
			query = query.Where("risk_levels.type = ?", "low")
		case "overdue_for_review":
			query = query.Where("kri_register.next_review_date < ?", time.Now())
		}
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		searchColumns := []string{
			"kri_register.title",
			"kri_register.calculation_type",
			"kri_register.description",
			"kri_register.reference_number",
			"kri_register.low_threshold",
			"kri_register.medium_threshold",
			"kri_register.high_threshold",
			"kri_register.current_value",
			"kri_register.annual_target",
			"kri_register.review_target",
			"kri_categories.title_en",
			"kri_categories.title_ar",
			"units.title_en",
			"units.title_ar",
			"units.code_en",
			"units.code_ar",
			"organization_units.title_en",
			"organization_units.title_ar",
			"review_frequencies.title_en",
			"review_frequencies.title_ar",
			"risk_levels.title_en",
			"risk_levels.title_ar",
			"level_of_measurements.title_en",
			"level_of_measurements.title_ar",
			"risk_responsible_units.title_en",
			"risk_responsible_units.title_ar",
			"kri_sub_categories.title_en",
			"kri_sub_categories.title_ar",
			"distribution_statuses.title_en",
			"distribution_statuses.title_ar",
			"responsible_user.first_name_en",
			"responsible_user.first_name_ar",
			"responsible_user.last_name_en",
			"responsible_user.last_name_ar",
			"created_user.first_name_en",
			"created_user.first_name_ar",
			"created_user.last_name_en",
			"created_user.last_name_ar",
			"updated_user.first_name_en",
			"updated_user.first_name_ar",
			"updated_user.last_name_en",
			"updated_user.last_name_ar",
		}

		conditions := make([]string, 0, len(searchColumns))
		args := make([]interface{}, 0, len(searchColumns))
		for _, col := range searchColumns {
			conditions = append(conditions, col+" LIKE ?")
			args = append(args, keyword)
		}

		query = query.Where(strings.Join(conditions, " OR "), args...)
	}

	// Sort by the specified column and order
	query = query.Order(params.SortBy + " " + params.Sort)

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply pagination only if not exporting
	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	if err := query.Find(&kriRegisters).Error; err != nil {
		return nil, err
	}

	if params.Export {
		return &riskDTOs.KRIRegisterListResponseDTO{
			Items: kriRegisters,
		}, nil
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &riskDTOs.KRIRegisterListResponseDTO{
		Items: kriRegisters,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *KRIRegisterRepository) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.KRIRegister, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var kriRegister *models.KRIRegister

	result := db.First(&kriRegister, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return kriRegister, nil

}

func (repo *KRIRegisterRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.KRIRegister, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var kriRegister *models.KRIRegister // Declare the variable inside the function

	result := db.
		Preload("KRICategory").
		Preload("ResponsibleUnit").
		Preload("ResponsibleUser").
		Preload("ReviewFrequency").
		Preload("RiskLevel").
		Preload("Unit").
		Preload("LevelOfMeasurement").
		Preload("RiskResponsibleUnit").
		Preload("KRISubCategory").
		Preload("DistributionStatus").
		Preload("ValueUpdates").
		Preload("ValueUpdates.KRIFrequencyValue").
		Preload("ValueUpdates.KRIFrequencyValue.RiskToleranceStatus").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("SubmittedBy").
		Preload("ApprovedBy").
		Preload("WorkflowStatus").
		Preload("KRIRegisterWorkflows").
		Preload("KRIRegisterWorkflows.KRIRegisterWorkflowParticipants").
		Preload("KRIRegisterWorkflows.KRIRegisterWorkflowReviews").
		Preload("KRIFrequencyValues").
		Preload("KRIFrequencyValues.RiskToleranceStatus").
		First(&kriRegister, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return kriRegister, nil
}

func (repo *KRIRegisterRepository) Create(kriRegister *models.KRIRegister, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(kriRegister).Error; err != nil {
		return 0, err
	}
	return kriRegister.ID, nil
}

func (repo *KRIRegisterRepository) Update(kriRegister *models.KRIRegister, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(kriRegister).Error; err != nil {
		return err
	}
	return nil
}

func (repo *KRIRegisterRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	result := db.Delete(&models.KRIRegister{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *KRIRegisterRepository) ExistsByTitle(Title string, KRICategoryID uint) (bool, error) {
	var count int64

	// Check for the existence of Title or QuestionAr in a single query
	err := repo.db.Model(&models.KRIRegister{}).
		Where("title = ? AND kri_category_id = ? ", Title, KRICategoryID).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	// Return true if either title exists, otherwise false
	return count > 0, nil

}

func (repo *KRIRegisterRepository) ExistsByTitleExcludingID(id uint, Title string, KRICategoryID uint) (bool, error) {
	var count int64
	err := repo.db.Model(&models.KRIRegister{}).
		Where("id <> ? AND title = ? AND kri_category_id = ?", id, Title, KRICategoryID).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}

func (repo *KRIRegisterRepository) GenerateReferenceNumber() (string, error) {
	var lastReferenceNumber string
	err := repo.db.Model(&models.KRIRegister{}).
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
	newReferenceNumber := fmt.Sprintf("%s-%d-%04d", "KRI", currentYear, lastReferenceID+1)

	return newReferenceNumber, nil
}

func (repo *KRIRegisterRepository) CreateInitialKriValueUpdate(kriValueUpdate *models.KRIValueUpdate, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(kriValueUpdate).Error; err != nil {
		return 0, err
	}
	return kriValueUpdate.ID, nil
}

func (repo *KRIRegisterRepository) DeleteKriValueUpdates(ids []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", ids).Delete(&models.KRIValueUpdate{}).Error; err != nil {
		return err
	}

	return nil
}

// Old implementation - kept for reference
func (repo *KRIRegisterRepository) GetTrendAnalysis(kriID uint, limit int) ([]*models.KRIValueUpdate, error) {
	var kriValueUpdate []*models.KRIValueUpdate

	result := repo.db.Where("kri_register_id = ?", kriID).
		Order("reviewed_date DESC").
		Limit(limit).
		Find(&kriValueUpdate)

	if result.Error != nil {
		return nil, result.Error
	}

	return kriValueUpdate, nil
}

// New implementation - fetches from kri_frequency_value with month-wise data and risk tolerance status
func (repo *KRIRegisterRepository) GetFrequencyTrendAnalysis(kriID uint, limit int) ([]*models.KRIFrequencyValue, error) {
	var frequencyValues []*models.KRIFrequencyValue

	result := repo.db.Preload("RiskToleranceStatus").
		Where("kri_register_id = ? AND achievement_value IS NOT NULL", kriID).
		Order("year ASC, month_en DESC").
		Limit(limit).
		Find(&frequencyValues)

	if result.Error != nil {
		return nil, result.Error
	}

	return frequencyValues, nil
}

func (repo *KRIRegisterRepository) ForceDelete(id uint, tx ...utility.DatabaseTransaction) error {

	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	queries := []string{
		`DELETE kud
		 FROM kri_value_update_documents AS kud
		 JOIN kri_value_updates AS ku ON ku.id = kud.kri_value_update_id
		 WHERE ku.kri_register_id = ?`,

		`DELETE FROM kri_value_updates WHERE kri_register_id = ?`,
		`DELETE FROM kri_frequency_values WHERE kri_register_id = ?`,
		`DELETE krwr
		 FROM kri_register_workflow_reviews AS krwr
		 JOIN kri_register_workflows AS krw ON krw.id = krwr.kri_register_workflow_id
		 WHERE krw.kri_register_id = ?`,
		`DELETE krwp
		 FROM kri_register_workflow_participants AS krwp
		 JOIN kri_register_workflows AS krw ON krw.id = krwp.kri_register_workflow_id
		 WHERE krw.kri_register_id = ?`,
		`DELETE FROM kri_register_workflows WHERE kri_register_id = ?`,
		`DELETE FROM kri_register_review_comments WHERE kri_register_id = ?`,
		`DELETE FROM risk_kris WHERE kri_register_id = ?`,
		`DELETE FROM kri_processes WHERE kri_register_id = ?`,
		`DELETE FROM kri_kpis WHERE kri_register_id = ?`,
		`DELETE FROM kri_incidents WHERE kri_register_id = ?`,
		`DELETE FROM kri_controls WHERE kri_register_id = ?`,
		`DELETE FROM kri_compliance_requirements WHERE kri_register_id = ?`,
		`DELETE FROM kri_associated_documents WHERE kri_register_id = ?`,
		`DELETE FROM asset_kris WHERE kri_register_id = ?`,
	}

	for _, q := range queries {
		res := db.Exec(q, id)
		if res.Error != nil {
			return res.Error
		}
	}

	// Force hard delete (ignore soft delete)
	return db.Delete(&models.KRIRegister{}, id).Error
}

func (repo *KRIRegisterRepository) CreateFrequencyValues(frequencyValues []*models.KRIFrequencyValue, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if len(frequencyValues) == 0 {
		return nil
	}

	if err := db.Create(&frequencyValues).Error; err != nil {
		return err
	}

	return nil
}

func (repo *KRIRegisterRepository) DeleteFrequencyValuesByKRIID(kriRegisterID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("kri_register_id = ?", kriRegisterID).Delete(&models.KRIFrequencyValue{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *KRIRegisterRepository) GetFrequencyValuesByKRIID(kriRegisterID uint, tx ...utility.DatabaseTransaction) ([]*models.KRIFrequencyValue, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var frequencyValues []*models.KRIFrequencyValue
	if err := db.Where("kri_register_id = ?", kriRegisterID).Find(&frequencyValues).Error; err != nil {
		return nil, err
	}

	return frequencyValues, nil
}

func (repo *KRIRegisterRepository) GetAllFrequencyValues(kriRegisterID uint) ([]*riskDTOs.KRIFrequencyValueListDTO, error) {
	var frequencyValues []*riskDTOs.KRIFrequencyValueListDTO

	query := repo.db.Table("kri_frequency_values").
		Select(`
			id,
			kri_register_id,
			month_en,
			month_ar,
			type,
			target_value,
			actual_value,
			achievement_value,
			CONVERT(VARCHAR, created_at, 120) AS created_at
		`).
		Where("kri_register_id = ?", kriRegisterID).
		Order("id ASC")

	if err := query.Find(&frequencyValues).Error; err != nil {
		return nil, err
	}

	return frequencyValues, nil
}
