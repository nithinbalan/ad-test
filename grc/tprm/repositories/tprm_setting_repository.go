package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

type TPRMSettingRepositoryInterface interface {
	GetTPRMSetting() (*models.TPRMSetting, error)
	Update(setting models.TPRMSetting) error
}

type TPRMSettingRepository struct {
	db *gorm.DB
}

func NewTPRMSettingRepository(db *gorm.DB) *TPRMSettingRepository {
	return &TPRMSettingRepository{
		db: db,
	}
}

func (repo *TPRMSettingRepository) GetTPRMSetting() (*models.TPRMSetting, error) {
	var setting *models.TPRMSetting
	result := repo.db.
		First(&setting)

	if result.Error != nil {
		return nil, result.Error
	}

	return setting, nil
}

func (repo *TPRMSettingRepository) Update(setting models.TPRMSetting) error {
	db := repo.db
	if err := db.Save(setting).Error; err != nil {
		return err
	}

	return nil
}
