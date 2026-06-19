package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

type AwarenessPriorityRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AwarenessPriority, error)
	GetByID(id uint) (*models.AwarenessPriority, error)
	GetByType(priority string) (*models.AwarenessPriority, error)
}

type AwarenessPriorityRepository struct {
	db *gorm.DB
}

func NewAwarenessPriorityRepository(db *gorm.DB) *AwarenessPriorityRepository {
	return &AwarenessPriorityRepository{
		db: db,
	}
}

func (repo *AwarenessPriorityRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AwarenessPriority, error) {
	var awarenessPriorities []*models.AwarenessPriority

	query := repo.db

	if keyword != "" {
		query = query.Where(
			"(title_en LIKE ? OR title_ar LIKE ?)",
			"%"+keyword+"%",
			"%"+keyword+"%",
		)
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

	query = query.Order(sortBy + " " + sort)

	result := query.Find(&awarenessPriorities)

	return awarenessPriorities, result.Error
}

func (repo *AwarenessPriorityRepository) GetByID(id uint) (*models.AwarenessPriority, error) {
	var awarenessPriority *models.AwarenessPriority

	result := repo.db.First(&awarenessPriority, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return awarenessPriority, nil
}

func (repo *AwarenessPriorityRepository) GetByType(priority string) (*models.AwarenessPriority, error) {
	var awarenessPriority *models.AwarenessPriority

	result := repo.db.Where("type = ?", priority).First(&awarenessPriority)

	if result.Error != nil {
		return nil, result.Error
	}

	return awarenessPriority, nil
}
