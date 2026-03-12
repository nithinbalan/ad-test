package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type CommitteeTypeRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string, typeFilter string) ([]*models.CommitteeType, error)
	GetByID(id uint) (*models.CommitteeType, error)
	GetByType(typeFilter string) (*models.CommitteeType, error)
}

// implementation
type CommitteeTypeRepository struct {
	db *gorm.DB
}

func NewCommitteeTypeRepository(db *gorm.DB) *CommitteeTypeRepository {
	return &CommitteeTypeRepository{
		db: db,
	}
}

func (repo *CommitteeTypeRepository) GetAll(sort string, sortBy string, keyword string, status string, typeFilter string) ([]*models.CommitteeType, error) {
	var committeeTypes []*models.CommitteeType

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

	if typeFilter == "statutory_reference" || typeFilter == "no_statutory_reference" {
		query = query.Where("type = ?", typeFilter)
	}

	query = query.Order(sortBy + " " + sort)
	result := query.Find(&committeeTypes)

	return committeeTypes, result.Error
}

func (repo *CommitteeTypeRepository) GetByID(id uint) (*models.CommitteeType, error) {
	var committeeType *models.CommitteeType

	result := repo.db.First(&committeeType, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return committeeType, result.Error
}

func (repo *CommitteeTypeRepository) GetByType(typeFilter string) (*models.CommitteeType, error) {
	var committeeType *models.CommitteeType

	result := repo.db.Where("type = ?", typeFilter).First(&committeeType)

	if result.Error != nil {
		return nil, result.Error
	}

	return committeeType, nil
}
