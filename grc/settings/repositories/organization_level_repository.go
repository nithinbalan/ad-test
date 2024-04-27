package repositories

import (
	settingsModel "jethur/models/settings"

	"gorm.io/gorm"
)

type OrganizationLevelRepositoryInterface interface {
	GetAll() ([]*settingsModel.OrganizationLevel, error)
	Create(organizationLevels []*settingsModel.OrganizationLevel) ([]*settingsModel.OrganizationLevel, error)
}

type OrganizationLevelRepository struct {
	db *gorm.DB
}

func NewOrganizationLevelRepository(db *gorm.DB) *OrganizationLevelRepository {
	return &OrganizationLevelRepository{
		db: db,
	}
}

func (repo *OrganizationLevelRepository) GetAll() ([]*settingsModel.OrganizationLevel, error) {
	var organizationLevels []*settingsModel.OrganizationLevel

	query := repo.db

	result := query.Find(&organizationLevels)

	return organizationLevels, result.Error
}

func (repo *OrganizationLevelRepository) Create(organizationLevels []*settingsModel.OrganizationLevel) ([]*settingsModel.OrganizationLevel, error) {
	var updatedOrganizationLevels []*settingsModel.OrganizationLevel
	var idsForRetention []uint

	// Start a new transaction
	tx := repo.db.Begin()

	// Error handling for rollback in case of failure
	defer func() {
		if repo := recover(); repo != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, err
	}

	for _, orgLevel := range organizationLevels {
		if orgLevel.ID == 0 {
			// ID is not present, create a new record
			err := tx.Create(&orgLevel).Error
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		} else {
			// ID is present, update the existing record
			err := tx.Model(&settingsModel.OrganizationLevel{}).Where("id = ?", orgLevel.ID).Updates(orgLevel).Error
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
		updatedOrganizationLevels = append(updatedOrganizationLevels, orgLevel)
		idsForRetention = append(idsForRetention, orgLevel.ID)
	}

	// Delete any records not in the list of idsForRetention
	if err := tx.Not("id", idsForRetention).Delete(&settingsModel.OrganizationLevel{}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return updatedOrganizationLevels, nil
}
