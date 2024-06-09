package repositories

import (
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskTreatmentGuidelineRepositoryInterface interface {
	GetAll(params *riskDTOs.RiskTreatmentGuidelineListQueryParamsDTO) ([]*riskDTOs.RiskTreatmentGuidelineListDTO, error)
	GetByID(id uint) (*models.RiskTreatmentGuideline, error)
	Create(guideline *models.RiskTreatmentGuideline) (uint, error)
	Update(guideline *models.RiskTreatmentGuideline, tx ...utility.DatabaseTransaction) error
	Delete(id uint) error
	ExistsByTimeframe(timeframeEn, timeframeAr string) (bool, error)
	ExistsByTimeframeExcludingID(id uint, timeframeEn, timeframeAr string) (bool, error)
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
}

type RiskTreatmentGuidelineRepository struct {
	db *gorm.DB
}

func NewRiskTreatmentGuidelineRepository(db *gorm.DB) *RiskTreatmentGuidelineRepository {
	return &RiskTreatmentGuidelineRepository{
		db: db,
	}
}

func (repo *RiskTreatmentGuidelineRepository) GetAll(params *riskDTOs.RiskTreatmentGuidelineListQueryParamsDTO) ([]*riskDTOs.RiskTreatmentGuidelineListDTO, error) {
	var guidelines []*riskDTOs.RiskTreatmentGuidelineListDTO

	query := repo.db.Table("risk_treatment_guidelines").
		Select(`
			risk_treatment_guidelines.id,
			risk_treatment_guidelines.risk_level_id,
			risk_treatment_guidelines.timeframe_en,
			risk_treatment_guidelines.timeframe_ar,
			risk_treatment_guidelines.status,
			risk_levels.title_en AS risk_level_title_en,
			risk_levels.title_ar AS risk_level_title_ar,
			risk_levels.color_code AS risk_level_color_code
		`).
		Joins(`
			LEFT JOIN risk_levels ON risk_levels.id = risk_treatment_guidelines.risk_level_id
		`)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		query = query.Where(`
			(risk_treatment_guidelines.timeframe_en LIKE ? OR
			risk_treatment_guidelines.timeframe_ar LIKE ? OR
			risk_levels.title_en LIKE ? OR
			risk_levels.title_ar LIKE ?)`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "risk_treatment_guidelines.timeframe_en" // Default sort by timeframe_en if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	// Execute the query
	err := query.Find(&guidelines).Error
	if err != nil {
		return nil, err
	}

	return guidelines, nil
}

func (repo *RiskTreatmentGuidelineRepository) GetByID(id uint) (*models.RiskTreatmentGuideline, error) {
	var guideline *models.RiskTreatmentGuideline

	result := repo.db.
		First(&guideline, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return guideline, nil
}

func (repo *RiskTreatmentGuidelineRepository) Create(guideline *models.RiskTreatmentGuideline) (uint, error) {
	if err := repo.db.Create(guideline).Error; err != nil {
		return 0, err
	}
	return guideline.ID, nil
}

func (repo *RiskTreatmentGuidelineRepository) Update(guideline *models.RiskTreatmentGuideline, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(guideline).Error; err != nil {
		return err
	}
	return nil
}

func (repo *RiskTreatmentGuidelineRepository) Delete(id uint) error {
	var guideline *models.RiskTreatmentGuideline
	if err := repo.db.First(&guideline, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.RiskTreatmentGuideline{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *RiskTreatmentGuidelineRepository) ExistsByTimeframe(timeframeEn, timeframeAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.RiskTreatmentGuideline{}).
		Where("timeframe_en = ? OR timeframe_ar = ?", timeframeEn, timeframeAr).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}

func (repo *RiskTreatmentGuidelineRepository) ExistsByTimeframeExcludingID(id uint, timeframeEn, timeframeAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.RiskTreatmentGuideline{}).
		Where("id <> ? AND (timeframe_en = ? OR timeframe_ar = ?)", id, timeframeEn, timeframeAr).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}

func (repo *RiskTreatmentGuidelineRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.RiskTreatmentGuideline{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.RiskTreatmentGuideline{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *RiskTreatmentGuidelineRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.RiskTreatmentGuideline{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.RiskTreatmentGuideline{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}
