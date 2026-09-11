package repositories

import (
	"jethur/grc/models"
	tiDTOs "jethur/grc/threat-intelligence/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type ThreatIntelligenceRiskLevelRepositoryInterface interface {
	GetAll(params *tiDTOs.ThreatIntelligenceRiskLevelListQueryParamsDTO) ([]*tiDTOs.ThreatIntelligenceRiskLevelDTO, error)
	GetByID(id uint) (*models.ThreatIntelligenceRiskLevel, error)
	Create(riskLevel *models.ThreatIntelligenceRiskLevel) (uint, error)
	Update(riskLevel *models.ThreatIntelligenceRiskLevel) error
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
	ExistsByType(riskLevelType string) (bool, error)
	ExistsByTypeExcludingID(id uint, riskLevelType string) (bool, error)
}

// implementation
type ThreatIntelligenceRiskLevelRepository struct {
	db *gorm.DB
}

func NewThreatIntelligenceRiskLevelRepository(db *gorm.DB) *ThreatIntelligenceRiskLevelRepository {
	return &ThreatIntelligenceRiskLevelRepository{
		db: db,
	}
}

func (repo *ThreatIntelligenceRiskLevelRepository) GetAll(params *tiDTOs.ThreatIntelligenceRiskLevelListQueryParamsDTO) ([]*tiDTOs.ThreatIntelligenceRiskLevelDTO, error) {
	var riskLevels []*tiDTOs.ThreatIntelligenceRiskLevelDTO

	query := repo.db.Table("threat_intelligence_risk_levels").
		Select("threat_intelligence_risk_levels.*")

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	if params.Keyword != "" {
		query = query.Where(`
			(threat_intelligence_risk_levels.title_en LIKE ? OR
				threat_intelligence_risk_levels.title_ar LIKE ? OR
				threat_intelligence_risk_levels.type LIKE ?)`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	if params.Status == "active" || params.Status == "inactive" {
		query = query.Where("threat_intelligence_risk_levels.status = ?", params.Status)
	}

	if params.SortBy == "" {
		params.SortBy = "threat_intelligence_risk_levels.title_en"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	result := query.Find(&riskLevels)

	return riskLevels, result.Error
}

func (repo *ThreatIntelligenceRiskLevelRepository) GetByID(id uint) (*models.ThreatIntelligenceRiskLevel, error) {
	var riskLevel *models.ThreatIntelligenceRiskLevel
	result := repo.db.Preload("ThreatIntelligenceRiskLevelValue").First(&riskLevel, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return riskLevel, result.Error
}

func (repo *ThreatIntelligenceRiskLevelRepository) Create(riskLevel *models.ThreatIntelligenceRiskLevel) (uint, error) {
	if err := repo.db.Create(riskLevel).Error; err != nil {
		return 0, err
	}

	return riskLevel.ID, nil
}

func (repo *ThreatIntelligenceRiskLevelRepository) Update(updatedRiskLevel *models.ThreatIntelligenceRiskLevel) error {
	if err := repo.db.Save(updatedRiskLevel).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ThreatIntelligenceRiskLevelRepository) Delete(id uint) error {
	var riskLevel *models.ThreatIntelligenceRiskLevel
	if err := repo.db.First(&riskLevel, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.ThreatIntelligenceRiskLevel{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *ThreatIntelligenceRiskLevelRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.ThreatIntelligenceRiskLevel{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ThreatIntelligenceRiskLevel{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *ThreatIntelligenceRiskLevelRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.ThreatIntelligenceRiskLevel{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ThreatIntelligenceRiskLevel{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *ThreatIntelligenceRiskLevelRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.ThreatIntelligenceRiskLevel{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}

func (repo *ThreatIntelligenceRiskLevelRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.ThreatIntelligenceRiskLevel{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}

func (repo *ThreatIntelligenceRiskLevelRepository) ExistsByType(riskLevelType string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.ThreatIntelligenceRiskLevel{}).
		Where("type = ?", riskLevelType).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}

func (repo *ThreatIntelligenceRiskLevelRepository) ExistsByTypeExcludingID(id uint, riskLevelType string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.ThreatIntelligenceRiskLevel{}).
		Where("id <> ? AND type = ?", id, riskLevelType).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
