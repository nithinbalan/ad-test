package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type BcpPhaseRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.BcpPhase, error)
	GetByID(id uint) (*models.BcpPhase, error)
}

// implementation
type BcpPhaseRepository struct {
	db *gorm.DB
}

func NewBcpPhaseRepository(db *gorm.DB) *BcpPhaseRepository {
	return &BcpPhaseRepository{
		db: db,
	}
}

func (repo *BcpPhaseRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.BcpPhase, error) {
	var bcpPhases []*models.BcpPhase
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
	result := query.Find(&bcpPhases)

	return bcpPhases, result.Error
}

func (repo *BcpPhaseRepository) GetByID(id uint) (*models.BcpPhase, error) {
	var category *models.BcpPhase
	result := repo.db.First(&category, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return category, nil
}
