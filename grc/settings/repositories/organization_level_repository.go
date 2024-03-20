package repositories

import (
	"errors"
	"jethur/grc/settings/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type OrganizationLevelRepositoryInterface interface {
	GetAll() ([]*models.OrganizationLevel, error)
	GetByID(id uint) (*models.OrganizationLevel, error)
	Create(organizationLevels []*models.OrganizationLevel) ([]*models.OrganizationLevel, error)
	Update(organizationLevels []*models.OrganizationLevel) ([]*models.OrganizationLevel, error)
	Delete(id uint) error
}

// implementation
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
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return organizationLevel, result.Error
}

func (repo *OrganizationLevelRepository) Create(organizationLevels []*models.OrganizationLevel) ([]*models.OrganizationLevel, error) {
	result := repo.db.Create(&organizationLevels)

	return organizationLevels, result.Error
}

func (repo *OrganizationLevelRepository) Update(organizationLevels []*models.OrganizationLevel) ([]*models.OrganizationLevel, error) {
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
			err := tx.Model(&models.OrganizationLevel{}).Where("id = ?", orgLevel.ID).Updates(orgLevel).Error
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

func (repo *OrganizationLevelRepository) Delete(id uint) error {
	var organizationLevel *models.OrganizationLevel

	if err := repo.db.First(&organizationLevel, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.OrganizationLevel{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
