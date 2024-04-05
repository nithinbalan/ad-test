package repositories

import (
	"errors"
	"jethur/grc/organization/models"
	organizationModels "jethur/grc/organization/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type StrategicProfileRepositoryInterface interface {
	GetByID(id uint) (*organizationModels.StrategicProfile, error)
	GetStrategicProfile() (*organizationModels.StrategicProfile, error)
	Update(strategicProfile *organizationModels.StrategicProfile) error
	UpdateValue(strategicProfileID uint, values []*organizationModels.StrategicProfileValue) error
}

type StrategicProfileRepository struct {
	db *gorm.DB
}

func NewStrategicProfileRepository(db *gorm.DB) *StrategicProfileRepository {
	return &StrategicProfileRepository{
		db: db,
	}
}

func (repo *StrategicProfileRepository) GetByID(id uint) (*organizationModels.StrategicProfile, error) {
	var strategicProfile *organizationModels.StrategicProfile
	result := repo.db.First(&strategicProfile, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return strategicProfile, result.Error
}

func (repo *StrategicProfileRepository) GetStrategicProfile() (*organizationModels.StrategicProfile, error) {
	var strategicProfile *organizationModels.StrategicProfile

	result := repo.db.Preload("StrategicProfileValues").First(&strategicProfile)

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

func (repo *StrategicProfileRepository) UpdateValue(strategicProfileID uint, values []*organizationModels.StrategicProfileValue) error {
	tx := repo.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Fetch existing values for the strategic profile
	var existingValues []organizationModels.StrategicProfileValue
	tx.Where("strategic_profile_id = ?", strategicProfileID).Find(&existingValues)

	// Map existing values by ID for easy lookup
	existingByID := make(map[uint]organizationModels.StrategicProfileValue)
	for _, val := range existingValues {
		existingByID[val.ID] = val
	}

	// Track IDs of values to keep to identify which ones to delete later
	toKeepIDs := make([]uint, 0)

	// Create or update values
	for _, value := range values {
		if value.ID > 0 && existingByID[value.ID].ID != 0 {
			// Update existing value
			if err := tx.Model(&organizationModels.StrategicProfileValue{}).Where("id = ?", value.ID).Updates(value).Error; err != nil {
				tx.Rollback()
				return err
			}
			toKeepIDs = append(toKeepIDs, value.ID)
		} else {
			// Create new value
			if err := tx.Create(&value).Error; err != nil {
				tx.Rollback()
				return err
			}
			toKeepIDs = append(toKeepIDs, value.ID)
		}
	}

	// Delete values not present in the incoming data
	if err := tx.Where("strategic_profile_id = ? AND id NOT IN ?", strategicProfileID, toKeepIDs).Delete(&models.StrategicProfileValue{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
