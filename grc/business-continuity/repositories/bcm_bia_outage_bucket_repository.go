package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type BCMBIAOutageBucketRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.BCMBIAOutageBucket, error)
	GetByID(id uint) (*models.BCMBIAOutageBucket, error)
}

// implementation
type BCMBIAOutageBucketRepository struct {
	db *gorm.DB
}

func NewBCMBIAOutageBucketRepository(db *gorm.DB) *BCMBIAOutageBucketRepository {
	return &BCMBIAOutageBucketRepository{
		db: db,
	}
}

func (repo *BCMBIAOutageBucketRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.BCMBIAOutageBucket, error) {
	var bcmBiaOutageBuckets []*models.BCMBIAOutageBucket

	query := repo.db

	if keyword != "" {
		query = query.Where("(title_en LIKE ? OR title_ar LIKE ?)", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "title_en"
	}
	// Validate sortBy against allowlist to prevent SQL injection
	allowedSortBy := map[string]bool{"title_en": true, "title_ar": true, "status": true, "created_at": true, "updated_at": true}
	if !allowedSortBy[sortBy] {
		sortBy = "title_en"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	if status == "active" || status == "inactive" {
		query = query.Where("status = ?", status)
	}

	query = query.Order(sortBy + " " + sort)
	result := query.Find(&bcmBiaOutageBuckets)

	return bcmBiaOutageBuckets, result.Error
}

func (repo *BCMBIAOutageBucketRepository) GetByID(id uint) (*models.BCMBIAOutageBucket, error) {
	var bcmBiaOutageBucket models.BCMBIAOutageBucket
	result := repo.db.First(&bcmBiaOutageBucket, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &bcmBiaOutageBucket, nil
}
