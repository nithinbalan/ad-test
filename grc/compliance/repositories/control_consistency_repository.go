package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ControlConsistencyRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ControlConsistency, error)
	GetByID(id uint) (*models.ControlConsistency, error)
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
}

type ControlConsistencyRepository struct {
	db *gorm.DB
}

func NewControlConsistencyRepository(db *gorm.DB) *ControlConsistencyRepository {
	return &ControlConsistencyRepository{
		db: db,
	}
}

func (repo *ControlConsistencyRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ControlConsistency, error) {

	var controlConsistencies []*models.ControlConsistency

	query := repo.db.Model(&models.ControlConsistency{})

	if keyword != "" {
		query = query.Where("(title_en LIKE ? OR title_ar LIKE ?)", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "title_en"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	if status == "active" || status == "inactive" {
		query = query.Where("status = ?", status)
	}

	query = query.Preload("CreatedByUser").Order(sortBy + " " + sort)
	result := query.Find(&controlConsistencies)

	return controlConsistencies, result.Error
}

func (repo *ControlConsistencyRepository) GetByID(id uint) (*models.ControlConsistency, error) {
	var controlConsistency *models.ControlConsistency
	result := repo.db.Preload("CreatedByUser").First(&controlConsistency, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return controlConsistency, result.Error
}

func (repo *ControlConsistencyRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.ControlConsistency{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ControlConsistency{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *ControlConsistencyRepository) Deactivate(ids []uint) ([]uint, error) {

	var idsToUpdate []uint

	if err := repo.db.Model(&models.ControlConsistency{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ControlConsistency{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}


