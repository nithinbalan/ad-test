package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type StrategicInitiativeStatusRepositoryInterface interface {
	List(sort string, sortBy string, keyword string, status string) ([]*models.StrategicInitiativeStatus, error)
	GetByType(statusType string) (*models.StrategicInitiativeStatus, error)
	GetByID(id uint) (*models.StrategicInitiativeStatus, error)
}

type StrategicInitiativeStatusRepository struct {
	db *gorm.DB
}

func NewStrategicInitiativeStatusRepository(db *gorm.DB) *StrategicInitiativeStatusRepository {
	return &StrategicInitiativeStatusRepository{db: db}
}

func (repo *StrategicInitiativeStatusRepository) List(sort string, sortBy string, keyword string, status string) ([]*models.StrategicInitiativeStatus, error) {
	var rows []*models.StrategicInitiativeStatus

	query := repo.db.Model(&models.StrategicInitiativeStatus{})

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
	if err := query.Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (repo *StrategicInitiativeStatusRepository) GetByID(id uint) (*models.StrategicInitiativeStatus, error) {
	var s models.StrategicInitiativeStatus
	err := repo.db.First(&s, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(err, id)
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (repo *StrategicInitiativeStatusRepository) GetByType(statusType string) (*models.StrategicInitiativeStatus, error) {
	var s models.StrategicInitiativeStatus
	err := repo.db.Where("type = ? AND status = ?", statusType, "active").First(&s).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(err, "strategic initiative status type "+statusType)
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}
