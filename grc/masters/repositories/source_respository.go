package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type SourceRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.Source, error)
	GetByID(id uint) (*models.Source, error)
	GetByType(typeStr string) (*models.Source, error)
}

// Implementation
type SourceRepository struct {
	db *gorm.DB
}

func NewSourceRepository(db *gorm.DB) *SourceRepository {
	return &SourceRepository{
		db: db,
	}
}

func (repo *SourceRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.Source, error) {
	var sources []*models.Source

	query := repo.db

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

	query = query.Order(sortBy + " " + sort)
	result := query.Find(&sources)

	return sources, result.Error
}

func (repo *SourceRepository) GetByID(id uint) (*models.Source, error) {
	var source *models.Source
	result := repo.db.First(&source, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return source, result.Error
}

func (repo *SourceRepository) GetByType(typeStr string) (*models.Source, error) {
	var source *models.Source

	result := repo.db.Where("type = ?", typeStr).First(&source)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Source type with type "+typeStr)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return source, nil
}
