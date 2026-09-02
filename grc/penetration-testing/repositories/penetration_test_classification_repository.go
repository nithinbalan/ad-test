package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type PenetrationTestClassificationRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PenetrationTestClassification, error)
	GetByID(id uint) (*models.PenetrationTestClassification, error)
}

type PenetrationTestClassificationRepository struct {
	db *gorm.DB
}

func NewPenetrationTestClassificationRepository(db *gorm.DB) *PenetrationTestClassificationRepository {
	return &PenetrationTestClassificationRepository{db: db}
}

func (repo *PenetrationTestClassificationRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PenetrationTestClassification, error) {
	var classifications []*models.PenetrationTestClassification
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
	return classifications, query.Order(sortBy + " " + sort).Find(&classifications).Error
}

func (repo *PenetrationTestClassificationRepository) GetByID(id uint) (*models.PenetrationTestClassification, error) {
	var value *models.PenetrationTestClassification
	result := repo.db.First(&value, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return value, nil
}
