package repositories

import (
	"errors"
	"jethur/grc/organization/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type UserStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.UserStatus, error)
	GetByID(id uint) (*models.UserStatus, error)
}

// implementation
type UserStatusRepository struct {
	db *gorm.DB
}

func NewUserStatusRepository(db *gorm.DB) *UserStatusRepository {
	return &UserStatusRepository{
		db: db,
	}
}

func (repo *UserStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.UserStatus, error) {

	var userStatuses []*models.UserStatus

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
	result := query.Find(&userStatuses)

	return userStatuses, result.Error
}

func (repo *UserStatusRepository) GetByID(id uint) (*models.UserStatus, error) {
	var status *models.UserStatus
	result := repo.db.First(&status, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return status, result.Error
}
