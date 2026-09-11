package repositories

import (
	"jethur/grc/models"
	tiDTOs "jethur/grc/threat-intelligence/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type ThreatIntelligenceSeverityRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*tiDTOs.ThreatIntelligenceSeverityDTO, error)
	GetByID(id uint) (*models.ThreatIntelligenceSeverity, error)
	Create(severity *models.ThreatIntelligenceSeverity) (uint, error)
	Update(severity *models.ThreatIntelligenceSeverity) error
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// implementation
type ThreatIntelligenceSeverityRepository struct {
	db *gorm.DB
}

func NewThreatIntelligenceSeverityRepository(db *gorm.DB) *ThreatIntelligenceSeverityRepository {
	return &ThreatIntelligenceSeverityRepository{
		db: db,
	}
}

func (repo *ThreatIntelligenceSeverityRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*tiDTOs.ThreatIntelligenceSeverityDTO, error) {
	var severities []*tiDTOs.ThreatIntelligenceSeverityDTO

	query := repo.db.Table("threat_intelligence_severities").
		Select("threat_intelligence_severities.*")

	if keyword != "" {
		query = query.Where(`
				(threat_intelligence_severities.title_en LIKE ? OR
				threat_intelligence_severities.title_ar LIKE ?)`,
			"%"+keyword+"%", "%"+keyword+"%")
	}

	if status == "active" || status == "inactive" {
		query = query.Where("threat_intelligence_severities.status = ?", status)
	}

	if sortBy == "" {
		sortBy = "threat_intelligence_severities.score"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	query = query.Order(sortBy + " " + sort)
	result := query.Find(&severities)

	return severities, result.Error
}

func (repo *ThreatIntelligenceSeverityRepository) GetByID(id uint) (*models.ThreatIntelligenceSeverity, error) {
	var severity *models.ThreatIntelligenceSeverity
	result := repo.db.First(&severity, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return severity, nil
}

func (repo *ThreatIntelligenceSeverityRepository) Create(severity *models.ThreatIntelligenceSeverity) (uint, error) {
	if err := repo.db.Create(severity).Error; err != nil {
		return 0, err
	}

	return severity.ID, nil
}

func (repo *ThreatIntelligenceSeverityRepository) Update(updatedSeverity *models.ThreatIntelligenceSeverity) error {
	if err := repo.db.Save(updatedSeverity).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ThreatIntelligenceSeverityRepository) Delete(id uint) error {
	var severity *models.ThreatIntelligenceSeverity
	if err := repo.db.First(&severity, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.ThreatIntelligenceSeverity{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *ThreatIntelligenceSeverityRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.ThreatIntelligenceSeverity{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ThreatIntelligenceSeverity{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *ThreatIntelligenceSeverityRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.ThreatIntelligenceSeverity{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ThreatIntelligenceSeverity{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *ThreatIntelligenceSeverityRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.ThreatIntelligenceSeverity{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}

func (repo *ThreatIntelligenceSeverityRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.ThreatIntelligenceSeverity{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
