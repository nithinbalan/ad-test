package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type BrandingRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.Branding, error)
	GetFullDetails() (*models.Branding, error)
	Update(branding *models.Branding, tx ...utility.DatabaseTransaction) error
}

type BrandingRepository struct {
	db *gorm.DB
}

func NewBrandingRepository(db *gorm.DB) BrandingRepositoryInterface {
	return &BrandingRepository{
		db: db,
	}
}

func (repo *BrandingRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *BrandingRepository) GetFullDetails() (*models.Branding, error) {
	var branding models.Branding

	err := repo.db.
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&branding).Error

	if err != nil {
		return nil, err
	}

	return &branding, nil
}

func (repo *BrandingRepository) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.Branding, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var branding models.Branding

	if err := db.First(&branding, id).Error; err != nil {
		return nil, err
	}

	return &branding, nil
}

func (repo *BrandingRepository) Update(branding *models.Branding, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(branding).Error; err != nil {
		return err
	}

	return nil
}