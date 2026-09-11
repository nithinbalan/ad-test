package repositories

import (
	"jethur/grc/models"
	tiDTOs "jethur/grc/threat-intelligence/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type ThreatIntelligenceLikelihoodRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*tiDTOs.ThreatIntelligenceLikelihoodDTO, error)
	GetByID(id uint) (*models.ThreatIntelligenceLikelihood, error)
	Create(likelihood *models.ThreatIntelligenceLikelihood) (uint, error)
	Update(likelihood *models.ThreatIntelligenceLikelihood) error
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// implementation
type ThreatIntelligenceLikelihoodRepository struct {
	db *gorm.DB
}

func NewThreatIntelligenceLikelihoodRepository(db *gorm.DB) *ThreatIntelligenceLikelihoodRepository {
	return &ThreatIntelligenceLikelihoodRepository{
		db: db,
	}
}

func (repo *ThreatIntelligenceLikelihoodRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*tiDTOs.ThreatIntelligenceLikelihoodDTO, error) {
	var likelihoods []*tiDTOs.ThreatIntelligenceLikelihoodDTO

	query := repo.db.Table("threat_intelligence_likelihoods").
		Select("threat_intelligence_likelihoods.*")

	if keyword != "" {
		query = query.Where(`
				(threat_intelligence_likelihoods.title_en LIKE ? OR
				threat_intelligence_likelihoods.title_ar LIKE ?)`,
			"%"+keyword+"%", "%"+keyword+"%")
	}

	if status == "active" || status == "inactive" {
		query = query.Where("threat_intelligence_likelihoods.status = ?", status)
	}

	if sortBy == "" {
		sortBy = "threat_intelligence_likelihoods.score"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	query = query.Order(sortBy + " " + sort)
	result := query.Find(&likelihoods)

	return likelihoods, result.Error
}

func (repo *ThreatIntelligenceLikelihoodRepository) GetByID(id uint) (*models.ThreatIntelligenceLikelihood, error) {
	var likelihood *models.ThreatIntelligenceLikelihood
	result := repo.db.First(&likelihood, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return likelihood, nil
}

func (repo *ThreatIntelligenceLikelihoodRepository) Create(likelihood *models.ThreatIntelligenceLikelihood) (uint, error) {
	if err := repo.db.Create(likelihood).Error; err != nil {
		return 0, err
	}

	return likelihood.ID, nil
}

func (repo *ThreatIntelligenceLikelihoodRepository) Update(updatedLikelihood *models.ThreatIntelligenceLikelihood) error {
	if err := repo.db.Save(updatedLikelihood).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ThreatIntelligenceLikelihoodRepository) Delete(id uint) error {
	var likelihood *models.ThreatIntelligenceLikelihood
	if err := repo.db.First(&likelihood, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.ThreatIntelligenceLikelihood{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *ThreatIntelligenceLikelihoodRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.ThreatIntelligenceLikelihood{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ThreatIntelligenceLikelihood{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *ThreatIntelligenceLikelihoodRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.ThreatIntelligenceLikelihood{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ThreatIntelligenceLikelihood{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *ThreatIntelligenceLikelihoodRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.ThreatIntelligenceLikelihood{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}

func (repo *ThreatIntelligenceLikelihoodRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.ThreatIntelligenceLikelihood{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
