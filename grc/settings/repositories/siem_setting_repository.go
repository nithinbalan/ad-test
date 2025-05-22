package repositories

import (
	"errors"
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interface
type SIEMSettingRepositoryInterface interface {
	GetFirst() (*models.SIEMSetting, error)
	Create(siem *models.SIEMSetting) error
	Update(siem *models.SIEMSetting) error
}

// Implementation
type SIEMSettingRepository struct {
	db *gorm.DB
}

func NewSIEMSettingRepository(db *gorm.DB) *SIEMSettingRepository {
	return &SIEMSettingRepository{
		db: db,
	}
}

func (repo *SIEMSettingRepository) GetFirst() (*models.SIEMSetting, error) {
	var setting models.SIEMSetting

	result := repo.db.First(&setting)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil // <- clean way to indicate "no existing record"
	}
	return &setting, result.Error
}

func (repo *SIEMSettingRepository) Create(siem *models.SIEMSetting) error {
	return repo.db.Create(siem).Error
}

func (repo *SIEMSettingRepository) Update(siem *models.SIEMSetting) error {
	return repo.db.Save(siem).Error
}
