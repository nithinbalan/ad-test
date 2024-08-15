package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

type OrganizationLevelRepositoryInterface interface {
	GetAll() ([]*models.OrganizationLevel, error)
	GetByID(id uint) (*models.OrganizationLevel, error)
	Create(organizationLevels []*models.OrganizationLevel) ([]*models.OrganizationLevel, error)
	GetPrimaryLevel() (*models.OrganizationLevel, error)
}

type OrganizationLevelRepository struct {
	db *gorm.DB
}

func NewOrganizationLevelRepository(db *gorm.DB) *OrganizationLevelRepository {
	return &OrganizationLevelRepository{
		db: db,
	}
}

func (repo *OrganizationLevelRepository) GetAll() ([]*models.OrganizationLevel, error) {
	var organizationLevels []*models.OrganizationLevel

	query := repo.db

	result := query.Find(&organizationLevels)

	return organizationLevels, result.Error
}

func (repo *OrganizationLevelRepository) GetByID(id uint) (*models.OrganizationLevel, error) {
	var organizationLevel *models.OrganizationLevel

	result := repo.db.First(&organizationLevel, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return organizationLevel, nil

}

func (repo *OrganizationLevelRepository) Create(organizationLevels []*models.OrganizationLevel) ([]*models.OrganizationLevel, error) {
	var updatedOrganizationLevels []*models.OrganizationLevel
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
			err := tx.Save(orgLevel).Error
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
		updatedOrganizationLevels = append(updatedOrganizationLevels, orgLevel)
		idsForRetention = append(idsForRetention, orgLevel.ID)
	}

	// Delete any records not in the list of idsForRetention
	if err := tx.Not("id", idsForRetention).Delete(&models.OrganizationLevel{}).Error; err != nil {
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
func (repo *OrganizationLevelRepository) GetPrimaryLevel() (*models.OrganizationLevel, error) {
	var primaryLevel models.OrganizationLevel

	// Fetch the primary organization level
	err := repo.db.
		Where("is_primary = ?", true).
		First(&primaryLevel).
		Error
	if err != nil {
		return nil, err
	}
	return &primaryLevel, nil
}
