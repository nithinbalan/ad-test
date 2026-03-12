package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type CommitteeTermRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string, typeFilter string) ([]*models.CommitteeTerm, error)
	GetByID(id uint) (*models.CommitteeTerm, error)
	GetByType(typeFilter string) (*models.CommitteeTerm, error)
}

// implementation
type CommitteeTermRepository struct {
	db *gorm.DB
}

func NewCommitteeTermRepository(db *gorm.DB) *CommitteeTermRepository {
	return &CommitteeTermRepository{
		db: db,
	}
}

func (repo *CommitteeTermRepository) GetAll(sort string, sortBy string, keyword string, status string, typeFilter string) ([]*models.CommitteeTerm, error) {
	var committeeTerms []*models.CommitteeTerm

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

	if typeFilter == "temporary" || typeFilter == "permanent" {
		query = query.Where("type = ?", typeFilter)
	}

	query = query.Order(sortBy + " " + sort)
	result := query.Find(&committeeTerms)

	return committeeTerms, result.Error
}

func (repo *CommitteeTermRepository) GetByID(id uint) (*models.CommitteeTerm, error) {
	var committeeTerm *models.CommitteeTerm

	result := repo.db.First(&committeeTerm, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return committeeTerm, result.Error
}

func (repo *CommitteeTermRepository) GetByType(typeFilter string) (*models.CommitteeTerm, error) {
	var committeeTerm *models.CommitteeTerm

	result := repo.db.Where("type = ?", typeFilter).First(&committeeTerm)

	if result.Error != nil {
		return nil, result.Error
	}

	return committeeTerm, nil
}
