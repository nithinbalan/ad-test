package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ControlAutomationRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ControlAutomation, error)
	GetByID(id uint) (*models.ControlAutomation, error)
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
}

type ControlAutomationRepository struct {
	db *gorm.DB
}

func NewControlAutomationRepository(db *gorm.DB) *ControlAutomationRepository {
	return &ControlAutomationRepository{
		db: db,
	}
}

func (repo *ControlAutomationRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ControlAutomation, error) {

	var controlAutomations []*models.ControlAutomation

	query := repo.db.Model(&models.ControlAutomation{})

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
	result := query.Find(&controlAutomations)

	return controlAutomations, result.Error
}

func (repo *ControlAutomationRepository) GetByID(id uint) (*models.ControlAutomation, error) {
	var controlAutomation *models.ControlAutomation
	result := repo.db.Preload("CreatedByUser").First(&controlAutomation, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return controlAutomation, result.Error
}

func (repo *ControlAutomationRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.ControlAutomation{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ControlAutomation{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *ControlAutomationRepository) Deactivate(ids []uint) ([]uint, error) {

	var idsToUpdate []uint

	if err := repo.db.Model(&models.ControlAutomation{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ControlAutomation{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}
