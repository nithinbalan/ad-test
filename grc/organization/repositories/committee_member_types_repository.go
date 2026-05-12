package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

type CommitteeMemberTypeRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.CommitteeMemberType, error)
	GetByID(id uint) (*models.CommitteeMemberType, error)
}

type CommitteeMemberTypeRepository struct {
	db *gorm.DB
}

func NewCommitteeMemberTypeRepository(db *gorm.DB) *CommitteeMemberTypeRepository {
	return &CommitteeMemberTypeRepository{
		db: db,
	}
}

func (repo *CommitteeMemberTypeRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.CommitteeMemberType, error) {
	var committeeMemberTypes []*models.CommitteeMemberType

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
	result := query.Find(&committeeMemberTypes)

	return committeeMemberTypes, result.Error
}

func (repo *CommitteeMemberTypeRepository) GetByID(id uint) (*models.CommitteeMemberType, error) {
	var committeeMemberType *models.CommitteeMemberType

	result := repo.db.First(&committeeMemberType, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return committeeMemberType, nil
}
