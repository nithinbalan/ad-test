package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type DataSourceRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.DataSource, error)
	GetByID(id uint) (*models.DataSource, error)
	GetByType(status string) (*models.DataSource, error)
}

// implementation
type DataSourceRepository struct {
	db *gorm.DB
}

func NewDataSourceRepository(db *gorm.DB) *DataSourceRepository {
	return &DataSourceRepository{
		db: db,
	}
}

func (repo *DataSourceRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.DataSource, error) {

	var dataSourcees []*models.DataSource

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
	result := query.Find(&dataSourcees)

	return dataSourcees, result.Error
}

func (repo *DataSourceRepository) GetByID(id uint) (*models.DataSource, error) {
	var dataSource *models.DataSource
	result := repo.db.First(&dataSource, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return dataSource, result.Error
}

func (repo *DataSourceRepository) GetByType(status string) (*models.DataSource, error) {
	var dataSource *models.DataSource

	result := repo.db.Where("type = ?", status).First(&dataSource)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Data Source with type "+status)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return dataSource, nil
}
