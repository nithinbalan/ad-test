package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type EntryTypeRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.EntryType, error)
	GetByID(id uint) (*models.EntryType, error)
	GetByType(status string) (*models.EntryType, error)
}

// implementation
type EntryTypeRepository struct {
	db *gorm.DB
}

func NewEntryTypeRepository(db *gorm.DB) *EntryTypeRepository {
	return &EntryTypeRepository{
		db: db,
	}
}

func (repo *EntryTypeRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.EntryType, error) {

	var entryTypees []*models.EntryType

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
	result := query.Find(&entryTypees)

	return entryTypees, result.Error
}

func (repo *EntryTypeRepository) GetByID(id uint) (*models.EntryType, error) {
	var entryType *models.EntryType
	result := repo.db.First(&entryType, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return entryType, result.Error
}

func (repo *EntryTypeRepository) GetByType(status string) (*models.EntryType, error) {
	var entryType *models.EntryType

	result := repo.db.Where("type = ?", status).First(&entryType)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Entry Type with type "+status)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return entryType, nil
}
