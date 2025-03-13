package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AssetPolicyChecklistRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetByID(id uint) (*models.AssetPolicyChecklist, error)
	Create(assetPolicyChecklist *models.AssetPolicyChecklist, tx ...utility.DatabaseTransaction) (uint, error)
	Update(checklist *models.AssetPolicyChecklist, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error

	DeleteAssetPolicyChecklists(ids []uint, tx ...utility.DatabaseTransaction) error
	BulkUpdateAssetPolicyChecklists(checklists []*models.AssetPolicyChecklist, tx ...utility.DatabaseTransaction) error
}

type AssetPolicyChecklistRepository struct {
	db *gorm.DB
}

func NewAssetPolicyChecklistRepository(db *gorm.DB) *AssetPolicyChecklistRepository {
	return &AssetPolicyChecklistRepository{db: db}
}

func (repo *AssetPolicyChecklistRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AssetPolicyChecklistRepository) GetByID(id uint) (*models.AssetPolicyChecklist, error) {
	var assetPolicyChecklist *models.AssetPolicyChecklist
	result := repo.db.First(&assetPolicyChecklist, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return assetPolicyChecklist, nil
}

func (repo *AssetPolicyChecklistRepository) Create(assetPolicyChecklist *models.AssetPolicyChecklist, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(assetPolicyChecklist).Error; err != nil {
		return 0, err
	}
	return assetPolicyChecklist.ID, nil
}

func (repo *AssetPolicyChecklistRepository) Update(checklist *models.AssetPolicyChecklist, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(checklist).Error
}

func (repo *AssetPolicyChecklistRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("id = ?", id).Delete(&models.AssetPolicyChecklist{}).Error
}

func (repo *AssetPolicyChecklistRepository) DeleteAssetPolicyChecklists(ids []uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete checklists in bulk
	return db.Where("id IN ?", ids).Delete(&models.AssetPolicyChecklist{}).Error
}

func (repo *AssetPolicyChecklistRepository) BulkUpdateAssetPolicyChecklists(checklists []*models.AssetPolicyChecklist, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Save each checklist individually in a loop
	for _, checklist := range checklists {
		if err := db.Save(checklist).Error; err != nil {
			return err
		}
	}

	return nil
}
