package repositories

import (
	bcmDTOs "jethur/grc/business-continuity/dtos"
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type BCMBIAImpactLevelRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*bcmDTOs.BCMBIAImpactLevelDTO, error)
	GetByID(id uint) (*models.BCMBIAImpactLevel, error)
	GetFullDetailsByID(id uint) (*models.BCMBIAImpactLevel, error)
	GetByImpactAnalysisLevelID(impactAnalysisLevelID uint) (*models.BCMBIAImpactLevel, error)
	GetByOutageBucketID(outageBucketID uint) (*models.BCMBIAImpactLevel, error)
	Update(bcmBiaImpactLevel *models.BCMBIAImpactLevel) error
}

// implementation
type BCMBIAImpactLevelRepository struct {
	db *gorm.DB
}

func NewBCMBIAImpactLevelRepository(db *gorm.DB) *BCMBIAImpactLevelRepository {
	return &BCMBIAImpactLevelRepository{
		db: db,
	}
}

func (repo *BCMBIAImpactLevelRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*bcmDTOs.BCMBIAImpactLevelDTO, error) {
	var bcmBiaImpactLevels []*bcmDTOs.BCMBIAImpactLevelDTO

	query := repo.db.Table("bcm_bia_impact_levels").
		Select(`
			bcm_bia_impact_levels.id,
			bcm_bia_impact_levels.impact_analysis_level_id,
			impact_level.title_en AS impact_analysis_level_title_en,
			impact_level.title_ar AS impact_analysis_level_title_ar,
			impact_level.type AS impact_analysis_level_type,
			impact_level.color_code AS impact_analysis_level_color_code,
			bcm_bia_impact_levels.outage_bucket_id,
			outage_bucket.title_en AS outage_bucket_title_en,
			outage_bucket.title_ar AS outage_bucket_title_ar,
			outage_bucket.from_value AS outage_bucket_from_value,
			outage_bucket.to_value AS outage_bucket_to_value,
			outage_bucket.duration_unit_id AS outage_bucket_duration_unit_id,
			outage_duration_unit.title_en AS outage_bucket_duration_unit_title_en,
			outage_duration_unit.title_ar AS outage_bucket_duration_unit_title_ar,
			outage_duration_unit.minutes AS outage_bucket_duration_unit_minutes,
			outage_duration_unit.type AS outage_bucket_duration_unit_type,
			bcm_bia_impact_levels.classification_en,
			bcm_bia_impact_levels.classification_ar,
			bcm_bia_impact_levels.is_critical,
			bcm_bia_impact_levels.rto_from_value,
			bcm_bia_impact_levels.rto_to_value,
			bcm_bia_impact_levels.rto_condition,
			bcm_bia_impact_levels.rto_unit_id,
			rto_unit.title_en AS rto_unit_title_en,
			rto_unit.title_ar AS rto_unit_title_ar,
			rto_unit.minutes AS rto_unit_minutes,
			rto_unit.type AS rto_unit_type,
			bcm_bia_impact_levels.mtpd_from_value,
			bcm_bia_impact_levels.mtpd_to_value,
			bcm_bia_impact_levels.mtpd_condition,
			bcm_bia_impact_levels.mtpd_unit_id,
			mtpd_unit.title_en AS mtpd_unit_title_en,
			mtpd_unit.title_ar AS mtpd_unit_title_ar,
			mtpd_unit.minutes AS mtpd_unit_minutes,
			mtpd_unit.type AS mtpd_unit_type,
			bcm_bia_impact_levels.rpo_from_value,
			bcm_bia_impact_levels.rpo_to_value,
			bcm_bia_impact_levels.rpo_condition,
			bcm_bia_impact_levels.rpo_unit_id,
			rpo_unit.title_en AS rpo_unit_title_en,
			rpo_unit.title_ar AS rpo_unit_title_ar,
			rpo_unit.minutes AS rpo_unit_minutes,
			rpo_unit.type AS rpo_unit_type,
			bcm_bia_impact_levels.rsp_from_value,
			bcm_bia_impact_levels.rsp_to_value,
			bcm_bia_impact_levels.rsp_condition,
			bcm_bia_impact_levels.rsp_unit_id,
			rsp_unit.title_en AS rsp_unit_title_en,
			rsp_unit.title_ar AS rsp_unit_title_ar,
			rsp_unit.minutes AS rsp_unit_minutes,
			rsp_unit.type AS rsp_unit_type,
			bcm_bia_impact_levels.[order],
			bcm_bia_impact_levels.status
		`).
		Joins("LEFT JOIN bcm_impact_analysis_levels AS impact_level ON impact_level.id = bcm_bia_impact_levels.impact_analysis_level_id").
		Joins("LEFT JOIN bcm_bia_outage_buckets AS outage_bucket ON outage_bucket.id = bcm_bia_impact_levels.outage_bucket_id").
		Joins("LEFT JOIN business_continuity_duration_units AS outage_duration_unit ON outage_duration_unit.id = outage_bucket.duration_unit_id").
		Joins("LEFT JOIN business_continuity_duration_units AS rto_unit ON rto_unit.id = bcm_bia_impact_levels.rto_unit_id").
		Joins("LEFT JOIN business_continuity_duration_units AS mtpd_unit ON mtpd_unit.id = bcm_bia_impact_levels.mtpd_unit_id").
		Joins("LEFT JOIN business_continuity_duration_units AS rpo_unit ON rpo_unit.id = bcm_bia_impact_levels.rpo_unit_id").
		Joins("LEFT JOIN business_continuity_duration_units AS rsp_unit ON rsp_unit.id = bcm_bia_impact_levels.rsp_unit_id")

	if keyword != "" {
		query = query.Where("(bcm_bia_impact_levels.classification_en LIKE ? OR bcm_bia_impact_levels.classification_ar LIKE ?)", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "order"
	}
	// Validate sortBy against allowlist to prevent SQL injection
	allowedSortBy := map[string]bool{"order": true, "classification_en": true, "classification_ar": true, "status": true, "created_at": true, "updated_at": true}
	if !allowedSortBy[sortBy] {
		sortBy = "order"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	if status == "active" || status == "inactive" {
		query = query.Where("bcm_bia_impact_levels.status = ?", status)
	}

	orderColumn := sortBy
	if orderColumn == "order" {
		orderColumn = "[order]"
	}
	query = query.Order("bcm_bia_impact_levels." + orderColumn + " " + sort)
	result := query.Find(&bcmBiaImpactLevels)

	return bcmBiaImpactLevels, result.Error
}

func (repo *BCMBIAImpactLevelRepository) GetByID(id uint) (*models.BCMBIAImpactLevel, error) {
	var bcmBiaImpactLevel models.BCMBIAImpactLevel
	result := repo.db.First(&bcmBiaImpactLevel, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &bcmBiaImpactLevel, nil
}

func (repo *BCMBIAImpactLevelRepository) GetFullDetailsByID(id uint) (*models.BCMBIAImpactLevel, error) {
	var bcmBiaImpactLevel *models.BCMBIAImpactLevel
	result := repo.db.
		Preload("ImpactAnalysisLevel").
		Preload("OutageBucket").
		Preload("OutageBucket.DurationUnit").
		Preload("RtoUnit").
		Preload("MtpdUnit").
		Preload("RpoUnit").
		Preload("RspUnit").
		First(&bcmBiaImpactLevel, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return bcmBiaImpactLevel, nil
}

func (repo *BCMBIAImpactLevelRepository) Update(bcmBiaImpactLevel *models.BCMBIAImpactLevel) error {
	if err := repo.db.Save(bcmBiaImpactLevel).Error; err != nil {
		return err
	}
	return nil
}

func (repo *BCMBIAImpactLevelRepository) GetByImpactAnalysisLevelID(impactAnalysisLevelID uint) (*models.BCMBIAImpactLevel, error) {
	var bcmBiaImpactLevel *models.BCMBIAImpactLevel

	result := repo.db.Where("impact_analysis_level_id = ?", impactAnalysisLevelID).First(&bcmBiaImpactLevel)

	if result.Error != nil {
		return nil, result.Error
	}

	return bcmBiaImpactLevel, nil
}

func (repo *BCMBIAImpactLevelRepository) GetByOutageBucketID(outageBucketID uint) (*models.BCMBIAImpactLevel, error) {
	var bcmBiaImpactLevel *models.BCMBIAImpactLevel

	result := repo.db.Where("outage_bucket_id = ?", outageBucketID).First(&bcmBiaImpactLevel)
	if result.Error != nil {
		return nil, result.Error
	}

	return bcmBiaImpactLevel, nil
}
