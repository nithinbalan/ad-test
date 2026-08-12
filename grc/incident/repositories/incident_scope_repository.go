package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type IncidentScopeRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.IncidentScope, error)
	GetByID(id uint) (*models.IncidentScope, error)
	Create(scope *models.IncidentScope) (*models.IncidentScope, error)
	Update(id uint, scope *models.IncidentScope) (*models.IncidentScope, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// implementation
type IncidentScopeRepository struct {
	db *gorm.DB
}

func NewIncidentScopeRepository(db *gorm.DB) *IncidentScopeRepository {
	return &IncidentScopeRepository{
		db: db,
	}
}

func (repo *IncidentScopeRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.IncidentScope, error) {
	var incidentScopes []*models.IncidentScope
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
	result := query.Find(&incidentScopes)

	return incidentScopes, result.Error
}

func (repo *IncidentScopeRepository) GetByID(id uint) (*models.IncidentScope, error) {
	var scope *models.IncidentScope
	result := repo.db.First(&scope, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return scope, nil
}

func (repo *IncidentScopeRepository) Create(scope *models.IncidentScope) (*models.IncidentScope, error) {
	result := repo.db.Create(&scope)

	return scope, result.Error
}

func (repo *IncidentScopeRepository) Update(id uint, updatedScope *models.IncidentScope) (*models.IncidentScope, error) {
	var scope *models.IncidentScope
	if err := repo.db.First(&scope, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&scope).Updates(updatedScope).Error; err != nil {
		return nil, err
	}

	return scope, nil
}

func (repo *IncidentScopeRepository) Delete(id uint) error {
	var scope *models.IncidentScope
	if err := repo.db.First(&scope, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.IncidentScope{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *IncidentScopeRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	if err := repo.db.Model(&models.IncidentScope{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.IncidentScope{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}
	return idsToUpdate, nil
}

func (repo *IncidentScopeRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	if err := repo.db.Model(&models.IncidentScope{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.IncidentScope{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}
	return idsToUpdate, nil
}

func (repo *IncidentScopeRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.IncidentScope{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}

func (repo *IncidentScopeRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.IncidentScope{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}
