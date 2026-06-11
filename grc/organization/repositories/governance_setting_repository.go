package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

type GovernanceSettingsRepositoryInterface interface {
	Get() (*models.GovernanceSetting, error)
	Update(setting *models.GovernanceSetting) (*models.GovernanceSetting, error)
}

type GovernanceSettingsRepository struct {
	db *gorm.DB
}

func NewGovernanceSettingsRepository(db *gorm.DB) *GovernanceSettingsRepository {
	return &GovernanceSettingsRepository{
		db: db,
	}
}

func (repo *GovernanceSettingsRepository) Get() (*models.GovernanceSetting, error) {
	var governanceSetting *models.GovernanceSetting

	result := repo.db.
		Preload("ApproveDepartment").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&governanceSetting)

	if result.Error != nil {
		return nil, result.Error
	}

	return governanceSetting, nil
}

func (repo *GovernanceSettingsRepository) Update(updatedSetting *models.GovernanceSetting) (*models.GovernanceSetting, error) {
	var setting models.GovernanceSetting

	if err := repo.db.First(&setting).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&setting).Updates(updatedSetting).Error; err != nil {
		return nil, err
	}

	return &setting, nil
}
