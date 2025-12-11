package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

type HealthCheckRepository interface {
	GetSystemPreference() (*models.SystemPreference, error)
}

type healthCheckRepo struct {
	db *gorm.DB
}

func NewHealthCheckRepository(db *gorm.DB) HealthCheckRepository {
	return &healthCheckRepo{db: db}
}

func (r *healthCheckRepo) GetSystemPreference() (*models.SystemPreference, error) {
	var pref models.SystemPreference
	err := r.db.First(&pref).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &pref, err
}
