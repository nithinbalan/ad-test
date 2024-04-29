package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

type StrategicProfileRepositoryInterface interface {
	Get() (*models.StrategicProfile, error)
	GetStrategicProfileDetails() (*models.StrategicProfile, error)
	Update(strategicProfile *models.StrategicProfile) error
	UpdateValues(values []*models.StrategicProfileValue) error
}

type StrategicProfileRepository struct {
	db *gorm.DB
}

func NewStrategicProfileRepository(db *gorm.DB) *StrategicProfileRepository {
	return &StrategicProfileRepository{
		db: db,
	}
}

func (repo *StrategicProfileRepository) Get() (*models.StrategicProfile, error) {
	var strategicProfile *models.StrategicProfile

	result := repo.db.First(&strategicProfile)

	if result.Error != nil {
		return nil, result.Error
	}

	return strategicProfile, nil
}

func (repo *StrategicProfileRepository) GetStrategicProfileDetails() (*models.StrategicProfile, error) {
	var strategicProfile *models.StrategicProfile

	result := repo.db.Preload("StrategicProfileValues").First(&strategicProfile)

	if result.Error != nil {
		return nil, result.Error
	}

	return strategicProfile, nil
}

func (repo *StrategicProfileRepository) Update(strategicProfile *models.StrategicProfile) error {
	if err := repo.db.Save(strategicProfile).Error; err != nil {
		return err
	}

	return nil
}

func (repo *StrategicProfileRepository) UpdateValues(values []*models.StrategicProfileValue) error {
	tx := repo.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Fetch existing values for the strategic profile
	var existingValues []models.StrategicProfileValue
	tx.Where("strategic_profile_id = ?", values[0].StrategicProfileID).Find(&existingValues)

	// Map existing values by ID for easy lookup
	existingByID := make(map[uint]models.StrategicProfileValue)
	for _, val := range existingValues {
		existingByID[val.ID] = val
	}

	// Track IDs of values to keep to identify which ones to delete later
	toKeepIDs := make([]uint, 0)

	// Create or update values
	for _, value := range values {
		if value.ID > 0 && existingByID[value.ID].ID != 0 {
			// Update existing value
			if err := tx.Model(&models.StrategicProfileValue{}).Where("id = ?", value.ID).Updates(value).Error; err != nil {
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
	if err := tx.Where("strategic_profile_id = ? AND id NOT IN ?", values[0].StrategicProfileID, toKeepIDs).Delete(&models.StrategicProfileValue{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
