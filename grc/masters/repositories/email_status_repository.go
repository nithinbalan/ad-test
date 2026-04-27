package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

type EmailStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.EmailStatus, error)
	GetByID(id uint) (*models.EmailStatus, error)
	GetByType(status string) (*models.EmailStatus, error)
}

type EmailStatusRepository struct {
	db *gorm.DB
}

func NewEmailStatusRepository(db *gorm.DB) *EmailStatusRepository {
	return &EmailStatusRepository{
		db: db,
	}
}

func (repo *EmailStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.EmailStatus, error) {
	var statuses []*models.EmailStatus
	query := repo.db.Model(&models.EmailStatus{})

	if keyword != "" {
		keyword = "%" + keyword + "%"
		query = query.Where("title_en LIKE ? OR title_ar LIKE ?", keyword, keyword)
	}

	if status != "" {
		query = query.Where("type = ?", status)
	}

	if sortBy != "" {
		if sort == "" {
			sort = "ASC"
		}
		query = query.Order(sortBy + " " + sort)
	}

	err := query.Find(&statuses).Error
	if err != nil {
		return nil, err
	}
	return statuses, nil
}

func (repo *EmailStatusRepository) GetByID(id uint) (*models.EmailStatus, error) {
	var status *models.EmailStatus
	err := repo.db.First(&status, id).Error
	if err != nil {
		return nil, err
	}
	return status, nil
}

func (repo *EmailStatusRepository) GetByType(statusType string) (*models.EmailStatus, error) {
	var status *models.EmailStatus
	err := repo.db.Where("type = ?", statusType).First(&status).Error
	if err != nil {
		return nil, err
	}
	return status, nil
}
