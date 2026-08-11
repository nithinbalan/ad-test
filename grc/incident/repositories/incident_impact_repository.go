package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type IncidentImpactRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.IncidentImpact, error)
	GetByID(id uint) (*models.IncidentImpact, error)
	Create(impact *models.IncidentImpact) (*models.IncidentImpact, error)
	Update(id uint, impact *models.IncidentImpact) (*models.IncidentImpact, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// implementation
type IncidentImpactRepository struct {
	db *gorm.DB
}

func NewIncidentImpactRepository(db *gorm.DB) *IncidentImpactRepository {
	return &IncidentImpactRepository{
		db: db,
	}
}

func (repo *IncidentImpactRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.IncidentImpact, error) {
	var incidentImpacts []*models.IncidentImpact
	query := repo.db

	if keyword != "" {
		query = query.Where("(title_en LIKE ? OR title_ar LIKE ?)", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "created_at"
	}
	if sort != "asc" && sort != "desc" {
		sort = "desc"
	}

	if status == "active" || status == "inactive" {
		query = query.Where("status = ?", status)
	}

	query = query.Order(sortBy + " " + sort)
	result := query.Find(&incidentImpacts)

	return incidentImpacts, result.Error
}

func (repo *IncidentImpactRepository) GetByID(id uint) (*models.IncidentImpact, error) {
	var impact *models.IncidentImpact
	result := repo.db.First(&impact, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return impact, nil
}

func (repo *IncidentImpactRepository) Create(impact *models.IncidentImpact) (*models.IncidentImpact, error) {
	result := repo.db.Create(&impact)

	return impact, result.Error
}

func (repo *IncidentImpactRepository) Update(id uint, updatedImpact *models.IncidentImpact) (*models.IncidentImpact, error) {
	var impact *models.IncidentImpact
	if err := repo.db.First(&impact, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&impact).Updates(updatedImpact).Error; err != nil {
		return nil, err
	}

	return impact, nil
}

func (repo *IncidentImpactRepository) Delete(id uint) error {
	var impact *models.IncidentImpact
	if err := repo.db.First(&impact, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.IncidentImpact{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *IncidentImpactRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	if err := repo.db.Model(&models.IncidentImpact{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.IncidentImpact{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}
	return idsToUpdate, nil
}

func (repo *IncidentImpactRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	if err := repo.db.Model(&models.IncidentImpact{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.IncidentImpact{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}
	return idsToUpdate, nil
}

func (repo *IncidentImpactRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.IncidentImpact{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}

func (repo *IncidentImpactRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.IncidentImpact{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}
