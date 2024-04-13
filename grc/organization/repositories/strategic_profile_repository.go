package repositories

import (
	organizationModels "jethur/grc/organization/models"

	"gorm.io/gorm"
)

type StrategicProfileRepositoryInterface interface {
	Get() (*organizationModels.StrategicProfile, error)
	GetStrategicProfileDetails() (*organizationModels.StrategicProfile, error)
	Update(strategicProfile *organizationModels.StrategicProfile) error
}

type StrategicProfileRepository struct {
	db *gorm.DB
}

func NewStrategicProfileRepository(db *gorm.DB) *StrategicProfileRepository {
	return &StrategicProfileRepository{
		db: db,
	}
}

func (repo *StrategicProfileRepository) Get() (*organizationModels.StrategicProfile, error) {
	var strategicProfile *organizationModels.StrategicProfile
	result := repo.db.First(&strategicProfile)
	if result.Error != nil {
		return nil, result.Error
	}

	// @review
	// please move this error return to service, here just return error as it is
	// everywhere we used to follow method you used, but we have to move to service layer from new apis
	return strategicProfile, nil
}

func (repo *StrategicProfileRepository) GetStrategicProfileDetails() (*organizationModels.StrategicProfile, error) {
	var strategicProfile *organizationModels.StrategicProfile

	result := repo.db.First(&strategicProfile)

	if result.Error != nil {
		return nil, result.Error
	}

	return strategicProfile, nil
}

func (repo *StrategicProfileRepository) Update(strategicProfile *organizationModels.StrategicProfile) error {
	if err := repo.db.Save(strategicProfile).Error; err != nil {
		return err
	}

	return nil
}
