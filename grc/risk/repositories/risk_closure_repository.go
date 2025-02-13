package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskClosureRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetFullDetailsByID(riskID uint, tx ...utility.DatabaseTransaction) (*models.RiskClosure, error)
	GetByID(id uint) (*models.RiskClosure, error)
	Create(riskClosure *models.RiskClosure, tx ...utility.DatabaseTransaction) (uint, error)
	Update(riskClosure *models.RiskClosure, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error

	CreateDocuments(docs []*models.RiskClosureDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.RiskClosureDocument, error)
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error
}

type RiskClosureRepository struct {
	db *gorm.DB
}

func NewRiskClosureRepository(db *gorm.DB) RiskClosureRepositoryInterface {
	return &RiskClosureRepository{
		db: db,
	}
}

func (repo *RiskClosureRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *RiskClosureRepository) GetFullDetailsByID(riskID uint, tx ...utility.DatabaseTransaction) (*models.RiskClosure, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var riskClosure models.RiskClosure
	result := db.
		Preload("RiskRegister").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("Documents").
		Where("risk_register_id = ?", riskID).
		First(&riskClosure)

	if result.Error != nil {
		return nil, result.Error
	}

	return &riskClosure, nil
}

func (repo *RiskClosureRepository) GetByID(id uint) (*models.RiskClosure, error) {
	var riskClosure models.RiskClosure
	result := repo.db.First(&riskClosure, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &riskClosure, nil
}

func (repo *RiskClosureRepository) Create(riskClosure *models.RiskClosure, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(riskClosure).Error; err != nil {
		return 0, err
	}
	return riskClosure.ID, nil
}

func (repo *RiskClosureRepository) Update(riskClosure *models.RiskClosure, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Save(riskClosure).Error; err != nil {
		return err
	}
	return nil
}

func (repo *RiskClosureRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Delete(&models.RiskClosure{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (repo *RiskClosureRepository) CreateDocuments(docs []*models.RiskClosureDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Save all documents in a single transaction
	if err := db.Create(&docs).Error; err != nil {
		return err
	}
	return nil
}

func (repo *RiskClosureRepository) GetDocumentByID(id uint) (*models.RiskClosureDocument, error) {
	var document models.RiskClosureDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *RiskClosureRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete documents by IDs
	if err := db.Where("id IN (?)", docIds).Delete(&models.RiskClosureDocument{}).Error; err != nil {
		return err
	}

	return nil
}
