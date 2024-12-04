package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type OrganizationLevelRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll() ([]*models.OrganizationLevel, error)
	GetByID(id uint) (*models.OrganizationLevel, error)
	CreateOrganizationLevels(organizationLevels []*models.OrganizationLevel) ([]*models.OrganizationLevel, error)
	Create(orgLevel *models.OrganizationLevel) (*models.OrganizationLevel, error)
	Update(id uint, updatedOrgLevel *models.OrganizationLevel) (*models.OrganizationLevel, error)
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	GetPrimaryLevel() (*models.OrganizationLevel, error)
	DeleteOrganizationUnits(ids []uint, tx ...utility.DatabaseTransaction) error

	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)

	ExistsByLevel(level uint8) (bool, error)
	ExistsByLevelExcludingID(id uint, level uint8) (bool, error)
}

type OrganizationLevelRepository struct {
	db *gorm.DB
}

func NewOrganizationLevelRepository(db *gorm.DB) *OrganizationLevelRepository {
	return &OrganizationLevelRepository{
		db: db,
	}
}

func (repo *OrganizationLevelRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *OrganizationLevelRepository) GetAll() ([]*models.OrganizationLevel, error) {
	var organizationLevels []*models.OrganizationLevel

	query := repo.db

	result := query.Find(&organizationLevels)

	return organizationLevels, result.Error
}

func (repo *OrganizationLevelRepository) GetByID(id uint) (*models.OrganizationLevel, error) {
	var organizationLevel *models.OrganizationLevel

	result := repo.db.
		Preload("OrganizationUnits").
		First(&organizationLevel, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return organizationLevel, nil

}

func (repo *OrganizationLevelRepository) CreateOrganizationLevels(organizationLevels []*models.OrganizationLevel) ([]*models.OrganizationLevel, error) {
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

func (repo *OrganizationLevelRepository) Create(orgLevel *models.OrganizationLevel) (*models.OrganizationLevel, error) {

	if err := repo.db.Create(orgLevel).Error; err != nil {
		return nil, err
	}
	return orgLevel, nil
}

func (repo *OrganizationLevelRepository) Update(id uint, updatedOrgLevel *models.OrganizationLevel) (*models.OrganizationLevel, error) {
	var orgLevel *models.OrganizationLevel

	if err := repo.db.First(&orgLevel, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&orgLevel).Updates(updatedOrgLevel).Error; err != nil {
		return nil, err
	}

	return orgLevel, nil
}

func (repo *OrganizationLevelRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var orgLevel *models.OrganizationLevel
	if err := db.First(&orgLevel, id).Error; err != nil {
		return err
	}

	result := db.Delete(&models.OrganizationLevel{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *OrganizationLevelRepository) DeleteOrganizationUnits(ids []uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("id IN ?", ids).Delete(&models.OrganizationUnit{}).Error
}

func (repo *OrganizationLevelRepository) ExistsByTitle(titleEn string, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.OrganizationLevel{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}

func (repo *OrganizationLevelRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.OrganizationLevel{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}

func (repo *OrganizationLevelRepository) ExistsByLevel(level uint8) (bool, error) {
	var count int64

	err := repo.db.Model(&models.OrganizationLevel{}).
		Where("level = ?", level).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}

func (repo *OrganizationLevelRepository) ExistsByLevelExcludingID(id uint, level uint8) (bool, error) {
	var count int64

	err := repo.db.Model(&models.OrganizationLevel{}).
		Where("id <> ? AND level = ?", id, level).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
