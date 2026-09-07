package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type PenetrationTestClosureRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetByID(penetrationTestID uint) (*models.PenetrationTestClosure, error)
	GetFullDetailsByID(penetrationTestID uint, tx ...utility.DatabaseTransaction) (*models.PenetrationTestClosure, error)
	Create(closure *models.PenetrationTestClosure, tx ...utility.DatabaseTransaction) (uint, error)
	Update(closure *models.PenetrationTestClosure, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
}

type PenetrationTestClosureRepository struct {
	db *gorm.DB
}

func NewPenetrationTestClosureRepository(db *gorm.DB) *PenetrationTestClosureRepository {
	return &PenetrationTestClosureRepository{
		db: db,
	}
}

func (repo *PenetrationTestClosureRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *PenetrationTestClosureRepository) GetByID(penetrationTestID uint) (*models.PenetrationTestClosure, error) {
	var closure *models.PenetrationTestClosure
	result := repo.db.Where("penetration_test_id = ?", penetrationTestID).First(&closure)

	if result.Error != nil {
		return nil, result.Error
	}

	return closure, nil
}

func (repo *PenetrationTestClosureRepository) GetFullDetailsByID(penetrationTestID uint, tx ...utility.DatabaseTransaction) (*models.PenetrationTestClosure, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var closure *models.PenetrationTestClosure

	result := db.Where("penetration_test_id = ?", penetrationTestID).
		Preload("ClosureOwner").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&closure)

	if result.Error != nil {
		return nil, result.Error
	}

	return closure, nil
}

func (repo *PenetrationTestClosureRepository) Create(closure *models.PenetrationTestClosure, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(closure).Error; err != nil {
		return 0, err
	}
	return closure.ID, nil
}

func (repo *PenetrationTestClosureRepository) Update(closure *models.PenetrationTestClosure, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(closure).Error; err != nil {
		return err
	}
	return nil
}

func (repo *PenetrationTestClosureRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var closure *models.PenetrationTestClosure
	if err := db.First(&closure, id).Error; err != nil {
		return err
	}

	result := db.Delete(&models.PenetrationTestClosure{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
