package repositories

import (
	"errors"
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/grc/strategy/dtos"
	"jethur/utility"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// Interfaces
type StrategicObjectiveRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction

	GetAll(params *dtos.StrategicObjectiveListQueryParamsDTO) (*dtos.StrategicObjectiveListResponseDTO, error)
	GetByID(id uint) (*models.StrategicObjective, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.StrategicObjective, error)
	Create(strategicObjective *models.StrategicObjective, tx ...utility.DatabaseTransaction) (uint, error)
	Update(strategicObjective *models.StrategicObjective, tx ...utility.DatabaseTransaction) error
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
	GenerateReferenceNumber() (string, error)
	GetFocusAreasByStrategicObjectiveID(strategicObjectiveID uint) ([]*dtos.StrategicObjectiveFocusAreaDTO, error)
	GetRiskSummary(strategicObjectiveID uint) (*dtos.StrategicObjectiveRiskSummaryDTO, error)
}

// implementation
type StrategicObjectiveRepository struct {
	db *gorm.DB
}

func NewStrategicObjectiveRepository(db *gorm.DB) *StrategicObjectiveRepository {
	return &StrategicObjectiveRepository{
		db: db,
	}
}
func (repo *StrategicObjectiveRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *StrategicObjectiveRepository) GetAll(params *dtos.StrategicObjectiveListQueryParamsDTO) (*dtos.StrategicObjectiveListResponseDTO, error) {
	var objectives []*dtos.StrategicObjectiveListDTO
	var totalCount int64

	query := repo.db.Table("strategic_objectives").
		Select(`
			strategic_objectives.*,
			strategic_objective_categories.title_en AS strategic_objective_category_title_en,
			strategic_objective_categories.title_ar AS strategic_objective_category_title_ar,

			organization_units.title_en AS responsible_unit_title_en,
			organization_units.title_ar AS responsible_unit_title_ar,
			responsible_user.first_name_en AS responsible_user_first_name_en,
			responsible_user.last_name_en AS responsible_user_last_name_en,
			responsible_user.first_name_ar AS responsible_user_first_name_ar,
			responsible_user.last_name_ar AS responsible_user_last_name_ar,
			responsible_user.initial AS responsible_user_initial,
			responsible_user.email AS responsible_user_email,
			responsible_user.color AS responsible_user_color,
			responsible_user.profile_image_path AS responsible_user_profile_image_path,
			review_user.first_name_en AS review_user_first_name_en,
			review_user.first_name_ar AS review_user_first_name_ar,
			review_user.last_name_en AS review_user_last_name_en,
			review_user.last_name_ar AS review_user_last_name_ar,
			review_user.initial AS review_user_initial,
			review_user.color AS review_user_color,
			review_user.email AS review_user_email,
			review_user.profile_image_path AS review_user_profile_image_path,

			review_frequencies.title_en AS review_frequency_title_en,
			review_frequencies.title_ar AS review_frequency_title_ar,

			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.email AS created_user_email,
			created_user.profile_image_path AS created_user_profile_image_path,

			updated_user.first_name_en AS updated_user_first_name_en,
			updated_user.first_name_ar AS updated_user_first_name_ar,
			updated_user.last_name_en AS updated_user_last_name_en,
			updated_user.last_name_ar AS updated_user_last_name_ar,
			updated_user.initial AS updated_user_initial,
			updated_user.color AS updated_user_color,
			updated_user.email AS updated_user_email,
			updated_user.profile_image_path AS updated_user_profile_image_path
		`).
		Joins(`LEFT JOIN strategic_objective_categories ON strategic_objective_categories.id = strategic_objectives.strategic_objective_category_id`).
		Joins(`LEFT JOIN organization_units ON organization_units.id = strategic_objectives.responsible_unit_id`).
		Joins("LEFT JOIN users AS responsible_user ON responsible_user.id = strategic_objectives.responsible_user_id").
		Joins(`LEFT JOIN users AS review_user ON review_user.id = strategic_objectives.review_user_id`).
		Joins(`LEFT JOIN review_frequencies ON review_frequencies.id = strategic_objectives.review_frequency_id`).
		Joins(`LEFT JOIN users AS created_user ON created_user.id = strategic_objectives.created_by_id`).
		Joins(`LEFT JOIN users AS updated_user ON updated_user.id = strategic_objectives.updated_by_id`)

	if params.Keyword != "" {
		kw := "%" + params.Keyword + "%"
		query = query.Where(`
			strategic_objectives.reference_number LIKE ? OR
			strategic_objectives.title_en LIKE ? OR
			strategic_objectives.title_ar LIKE ? OR
			strategic_objective_categories.title_en LIKE ? OR
			strategic_objective_categories.title_ar LIKE ?
		`, kw, kw, kw, kw, kw)
	}

	for col, val := range params.Filters {
		if len(val) > 0 {
			query = query.Where(col+" IN (?)", val)
		}
	}

	if params.StartDate != "" && params.EndDate != "" {
		query = query.Where("strategic_objectives.created_at BETWEEN ? AND ?", params.StartDate+" 00:00:00", params.EndDate+" 23:59:59")
	}

	// Sorting
	query = query.Order(params.SortBy + " " + params.Sort)

	// Pagination
	query.Count(&totalCount)
	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize)

	if err := query.Find(&objectives).Error; err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	response := &dtos.StrategicObjectiveListResponseDTO{
		Items: objectives,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *StrategicObjectiveRepository) GetByID(id uint) (*models.StrategicObjective, error) {
	var strategicObjective *models.StrategicObjective
	result := repo.db.First(&strategicObjective, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return strategicObjective, nil
}

func (repo *StrategicObjectiveRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.StrategicObjective, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var strategicObjective *models.StrategicObjective

	result := db.
		Preload("StrategicObjectiveCategory").
		Preload("ResponsibleUser").
		Preload("ResponsibleUnit").
		Preload("ReviewUser").
		Preload("ReviewFrequency").
		Preload("StrategicObjectiveKPIs").
		Preload("StrategicObjectiveKPIs.KPIRegister").
		Preload("StrategicObjectiveKPIs.KPIRegister.KPICategory").
		Preload("StrategicObjectiveKPIs.KPIRegister.KPIStatus").
		Preload("StrategicObjectiveKPIs.KPIRegister.ResponsibleUser").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&strategicObjective, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return strategicObjective, nil
}

func (repo *StrategicObjectiveRepository) Create(strategicObjective *models.StrategicObjective, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(strategicObjective).Error; err != nil {
		return 0, err
	}

	return strategicObjective.ID, nil
}

func (repo *StrategicObjectiveRepository) Update(strategicObjective *models.StrategicObjective, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Save(strategicObjective).Error; err != nil {
		return err
	}

	return nil
}

func (repo *StrategicObjectiveRepository) Delete(id uint) error {
	var strategicObjective *models.StrategicObjective
	if err := repo.db.First(&strategicObjective, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.StrategicObjective{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *StrategicObjectiveRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.StrategicObjective{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.StrategicObjective{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *StrategicObjectiveRepository) Deactivate(ids []uint) ([]uint, error) {

	var idsToUpdate []uint

	if err := repo.db.Model(&models.StrategicObjective{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.StrategicObjective{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil

}

func (repo *StrategicObjectiveRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.StrategicObjective{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}

	return count > 0, nil
}

func (repo *StrategicObjectiveRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.StrategicObjective{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}

func (repo *StrategicObjectiveRepository) GenerateReferenceNumber() (string, error) {
	var lastReferenceNumber string

	err := repo.db.Model(&models.StrategicObjective{}).
		Where("reference_number IS NOT NULL").
		Select("reference_number").
		Order("reference_number DESC").
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

	newReferenceNumber := fmt.Sprintf("%s-%d-%04d", "SO", currentYear, lastReferenceID+1)
	return newReferenceNumber, nil
}

func (repo *StrategicObjectiveRepository) GetFocusAreasByStrategicObjectiveID(strategicObjectiveID uint) ([]*dtos.StrategicObjectiveFocusAreaDTO, error) {
	var focusAreas []*dtos.StrategicObjectiveFocusAreaDTO
	err := repo.db.Table("strategy_profile_strategic_objectives as spso").
		Select(`DISTINCT fa.id, fa.title_en, fa.title_ar, fa.status`).
		Joins(`LEFT JOIN strategy_profile_focus_areas spfa ON spfa.id = spso.strategy_profile_focus_area_id`).
		Joins(`LEFT JOIN focus_areas fa ON fa.id = spfa.focus_area_id`).
		Where("spso.strategic_objective_id = ? AND fa.id IS NOT NULL", strategicObjectiveID).
		Scan(&focusAreas).Error
	if err != nil {
		return nil, err
	}
	return focusAreas, nil
}

func (repo *StrategicObjectiveRepository) GetRiskSummary(strategicObjectiveID uint) (*dtos.StrategicObjectiveRiskSummaryDTO, error) {
	result := &dtos.StrategicObjectiveRiskSummaryDTO{}

	// Get total risks count
	var totalRisks int64
	err := repo.db.Model(&models.RiskStrategicObjective{}).
		Where("strategic_objective_id = ?", strategicObjectiveID).
		Count(&totalRisks).Error
	if err != nil {
		return nil, err
	}

	totalRisksInt := int(totalRisks)
	result.TotalRisks = &totalRisksInt

	// Count high risks
	var highRisks int64
	err = repo.db.Model(&models.RiskStrategicObjective{}).
		Joins("JOIN risk_register rr ON risk_strategic_objectives.risk_register_id = rr.id").
		Joins("JOIN risk_matrix_risk_rating_values rmrrv ON rr.current_risk_rating_value_id = rmrrv.id").
		Joins("JOIN risk_matrix_risk_ratings rmrr ON rmrrv.risk_matrix_risk_rating_id = rmrr.id").
		Where("risk_strategic_objectives.strategic_objective_id = ?", strategicObjectiveID).
		Where(`
        rmrr.id IN (
            SELECT TOP 1 rmrr2.id 
            FROM risk_matrix_risk_ratings rmrr2
            ORDER BY 
                CASE LOWER(rmrr2.type)
                    WHEN 'very_high' THEN 1
                    WHEN 'high' THEN 2
                    WHEN 'medium' THEN 3
                    WHEN 'low' THEN 4
                    WHEN 'very_low' THEN 5
                    ELSE 6
                END
        )
    `).
		Count(&highRisks).Error

	if err != nil {
		return nil, fmt.Errorf("failed to count high risks: %w", err)
	}

	highRisksInt := int(highRisks)
	result.HighRisks = &highRisksInt

	// Count risks that have controls mapped to them
	var controlledRisks int64
	err = repo.db.Model(&models.RiskStrategicObjective{}).
		Joins("JOIN risk_register rr ON risk_strategic_objectives.risk_register_id = rr.id").
		Joins("JOIN risk_controls rc ON rr.id = rc.risk_register_id").
		Where("risk_strategic_objectives.strategic_objective_id = ?", strategicObjectiveID).
		Count(&controlledRisks).Error

	if err != nil {
		return nil, fmt.Errorf("failed to count controlled risks: %w", err)
	}

	// Calculate coverage (if applicable)
	var coverage float64
	if totalRisks > 0 {
		coverage = float64(controlledRisks) / float64(totalRisks) * 100
	}
	result.Coverage = &coverage

	// Also store the count of controlled risks
	controlledRisksInt := int(controlledRisks)
	result.ControllRiskCount = &controlledRisksInt

	// Get last reviewed date from kpi_register's LastReviewedDate for KPIs mapped to this strategic objective
	var lastReviewed *time.Time
	err = repo.db.Model(&models.KPIRegister{}).
		Select("kpi_register.last_reviewed_date").
		Joins("JOIN strategy_profile_kpis spk ON spk.kpi_id = kpi_register.id").
		Joins("JOIN strategy_profile_strategic_objectives spso ON spk.strategy_profile_strategic_objective_id = spso.id").
		Where("spso.strategic_objective_id = ? AND kpi_register.last_reviewed_date IS NOT NULL", strategicObjectiveID).
		Order("kpi_register.last_reviewed_date DESC").
		Limit(1).
		Scan(&lastReviewed).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to get last reviewed date: %w", err)
	}

	// Only format and set LastReviewed if we have a valid date
	if lastReviewed != nil && !lastReviewed.IsZero() {
		lastReviewedStr := lastReviewed.Format("02 Jan 2006")
		result.LastReviewed = &lastReviewedStr
	} else {
		result.LastReviewed = nil
	}

	// Get top risk (highest severity)
	err = repo.db.Raw(`
		SELECT 
			rr.id,
			rr.title,
			rmrr.type AS rating,
			rr.description,
			rmi.title_en AS impact_title_en,
			rmi.title_ar AS impact_title_ar,
			rml.title_en AS likelihood_title_en,
			rml.title_ar AS likelihood_title_ar
			FROM risk_strategic_objectives rso
			JOIN risk_register rr ON rso.risk_register_id = rr.id
			JOIN risk_matrix_risk_rating_values rmrrv ON rr.current_risk_rating_value_id = rmrrv.id
			JOIN risk_matrix_risk_ratings rmrr ON rmrrv.risk_matrix_risk_rating_id = rmrr.id
			LEFT JOIN risk_matrix_impacts rmi ON rr.current_impact_id = rmi.id
			LEFT JOIN risk_matrix_likelihoods rml ON rr.current_likelihood_id = rml.id
			WHERE rso.strategic_objective_id = ?
			ORDER BY 
				CASE 
				WHEN LOWER(rmrr.type) = 'very_high' THEN 1
				WHEN LOWER(rmrr.type) = 'high' THEN 2
				WHEN LOWER(rmrr.type) = 'medium' THEN 3
				WHEN LOWER(rmrr.type) = 'low' THEN 4
				WHEN LOWER(rmrr.type) = 'very_low' THEN 5
				ELSE 6
			END,
			rr.updated_at DESC
		OFFSET 0 ROWS FETCH NEXT 1 ROWS ONLY
		`, strategicObjectiveID).Scan(&result.TopRisk).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return result, nil
}
