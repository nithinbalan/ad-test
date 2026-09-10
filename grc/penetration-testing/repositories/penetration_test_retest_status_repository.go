package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type PenetrationTestRetestStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PenetrationTestRetestStatus, error)
	GetByID(id uint) (*models.PenetrationTestRetestStatus, error)
	GetByType(statusType string) (*models.PenetrationTestRetestStatus, error)
}

type PenetrationTestRetestStatusRepository struct{ db *gorm.DB }

func NewPenetrationTestRetestStatusRepository(db *gorm.DB) *PenetrationTestRetestStatusRepository {
	return &PenetrationTestRetestStatusRepository{db: db}
}

func (repo *PenetrationTestRetestStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PenetrationTestRetestStatus, error) {
	var values []*models.PenetrationTestRetestStatus
	query := repo.db
	if keyword != "" {
		query = query.Where("(title_en LIKE ? OR title_ar LIKE ?)", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if sortBy == "" {
		sortBy = "title_en"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	return values, query.Order(sortBy + " " + sort).Find(&values).Error
}

func (repo *PenetrationTestRetestStatusRepository) GetByID(id uint) (*models.PenetrationTestRetestStatus, error) {
	var value *models.PenetrationTestRetestStatus
	result := repo.db.First(&value, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return value, nil
}

func (repo *PenetrationTestRetestStatusRepository) GetByType(statusType string) (*models.PenetrationTestRetestStatus, error) {
	var value *models.PenetrationTestRetestStatus
	result := repo.db.Where("type = ?", statusType).First(&value)
	if result.Error != nil {
		return nil, result.Error
	}
	return value, nil
}
