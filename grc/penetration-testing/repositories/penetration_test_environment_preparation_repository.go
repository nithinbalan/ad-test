package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type PenetrationTestEnvironmentPreparationRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	ExistsByPenetrationTestID(penetrationTestID uint, tx ...utility.DatabaseTransaction) (bool, error)
	GetByID(penetrationTestID uint) (*models.PenetrationTestEnvironmentPreparation, error)
	GetFullDetailsByID(penetrationTestID uint, tx ...utility.DatabaseTransaction) (*models.PenetrationTestEnvironmentPreparation, error)
	Create(environmentPreparation *models.PenetrationTestEnvironmentPreparation, tx ...utility.DatabaseTransaction) (uint, error)
	Update(environmentPreparation *models.PenetrationTestEnvironmentPreparation, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
}

type PenetrationTestEnvironmentPreparationRepository struct {
	db *gorm.DB
}

func NewPenetrationTestEnvironmentPreparationRepository(db *gorm.DB) *PenetrationTestEnvironmentPreparationRepository {
	return &PenetrationTestEnvironmentPreparationRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *PenetrationTestEnvironmentPreparationRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *PenetrationTestEnvironmentPreparationRepository) GetByID(penetrationTestID uint) (*models.PenetrationTestEnvironmentPreparation, error) {
	var environmentPreparation *models.PenetrationTestEnvironmentPreparation
	result := repo.db.Where("penetration_test_id = ?", penetrationTestID).First(&environmentPreparation)

	if result.Error != nil {
		return nil, result.Error
	}

	return environmentPreparation, nil
}

func (repo *PenetrationTestEnvironmentPreparationRepository) GetFullDetailsByID(penetrationTestID uint, tx ...utility.DatabaseTransaction) (*models.PenetrationTestEnvironmentPreparation, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var environmentPreparation *models.PenetrationTestEnvironmentPreparation

	result := db.Where("penetration_test_id = ?", penetrationTestID).
		Preload("EnvironmentReadiness").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&environmentPreparation)

	if result.Error != nil {
		return nil, result.Error
	}

	return environmentPreparation, nil
}

func (repo *PenetrationTestEnvironmentPreparationRepository) Create(environmentPreparation *models.PenetrationTestEnvironmentPreparation, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(environmentPreparation).Error; err != nil {
		return 0, err
	}
	return environmentPreparation.ID, nil
}

func (repo *PenetrationTestEnvironmentPreparationRepository) Update(environmentPreparation *models.PenetrationTestEnvironmentPreparation, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(environmentPreparation).Error; err != nil {
		return err
	}
	return nil
}

func (repo *PenetrationTestEnvironmentPreparationRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var environmentPreparation *models.PenetrationTestEnvironmentPreparation
	if err := db.First(&environmentPreparation, id).Error; err != nil {
		return err
	}

	result := db.Delete(&models.PenetrationTestEnvironmentPreparation{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *PenetrationTestEnvironmentPreparationRepository) ExistsByPenetrationTestID(penetrationTestID uint, tx ...utility.DatabaseTransaction) (bool, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var count int64
	err := db.Model(&models.PenetrationTestEnvironmentPreparation{}).
		Where("penetration_test_id = ?", penetrationTestID).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
