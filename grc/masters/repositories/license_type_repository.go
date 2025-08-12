package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type LicenseTypeRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.LicenseType, error)
	GetByID(id uint) (*models.LicenseType, error)
	GetByType(status string) (*models.LicenseType, error)
}

// implementation
type LicenseTypeRepository struct {
	db *gorm.DB
}

func NewLicenseTypeRepository(db *gorm.DB) *LicenseTypeRepository {
	return &LicenseTypeRepository{
		db: db,
	}
}

func (repo *LicenseTypeRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.LicenseType, error) {
	var licenseTypes []*models.LicenseType
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
	result := query.Find(&licenseTypes)

	return licenseTypes, result.Error
}

func (repo *LicenseTypeRepository) GetByID(id uint) (*models.LicenseType, error) {
	var category *models.LicenseType
	result := repo.db.First(&category, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return category, nil
}

func (repo *LicenseTypeRepository) GetByType(status string) (*models.LicenseType, error) {
	var licenseType *models.LicenseType

	result := repo.db.Where("type = ?", status).First(&licenseType)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "License Type with type "+status)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return licenseType, nil
}
