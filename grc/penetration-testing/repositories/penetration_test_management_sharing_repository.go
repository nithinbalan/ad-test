package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type PenetrationTestManagementSharingRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetByID(penetrationTestID uint) (*models.PenetrationTestManagementSharing, error)
	GetFullDetailsByID(penetrationTestID uint, tx ...utility.DatabaseTransaction) (*models.PenetrationTestManagementSharing, error)
	Create(managementSharing *models.PenetrationTestManagementSharing, tx ...utility.DatabaseTransaction) (uint, error)
	Update(managementSharing *models.PenetrationTestManagementSharing, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
}

type PenetrationTestManagementSharingRepository struct {
	db *gorm.DB
}

func NewPenetrationTestManagementSharingRepository(db *gorm.DB) *PenetrationTestManagementSharingRepository {
	return &PenetrationTestManagementSharingRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *PenetrationTestManagementSharingRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *PenetrationTestManagementSharingRepository) GetByID(penetrationTestID uint) (*models.PenetrationTestManagementSharing, error) {
	var managementSharing *models.PenetrationTestManagementSharing
	result := repo.db.Where("penetration_test_id = ?", penetrationTestID).First(&managementSharing)

	if result.Error != nil {
		return nil, result.Error
	}

	return managementSharing, nil
}

func (repo *PenetrationTestManagementSharingRepository) GetFullDetailsByID(penetrationTestID uint, tx ...utility.DatabaseTransaction) (*models.PenetrationTestManagementSharing, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var managementSharing *models.PenetrationTestManagementSharing

	result := db.Where("penetration_test_id = ?", penetrationTestID).
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&managementSharing)

	if result.Error != nil {
		return nil, result.Error
	}

	return managementSharing, nil
}

func (repo *PenetrationTestManagementSharingRepository) Create(managementSharing *models.PenetrationTestManagementSharing, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(managementSharing).Error; err != nil {
		return 0, err
	}
	return managementSharing.ID, nil
}

func (repo *PenetrationTestManagementSharingRepository) Update(managementSharing *models.PenetrationTestManagementSharing, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(managementSharing).Error; err != nil {
		return err
	}
	return nil
}

func (repo *PenetrationTestManagementSharingRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var managementSharing *models.PenetrationTestManagementSharing
	if err := db.First(&managementSharing, id).Error; err != nil {
		return err
	}

	result := db.Delete(&models.PenetrationTestManagementSharing{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
