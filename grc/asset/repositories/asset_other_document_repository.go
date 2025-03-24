package repositories

import (
	"errors"
	"fmt"
	"jethur/grc/asset/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AssetOtherDocumentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(assetRegisterID uint) ([]*dtos.AssetOtherDocumentListDTO, error)
	GetByID(id uint) (*models.AssetOtherDocument, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AssetOtherDocument, error)
	Create(doc *models.AssetOtherDocument, tx ...utility.DatabaseTransaction) (uint, error)
	Update(doc *models.AssetOtherDocument) error
	CreateAssetOtherDocuments(documents []*models.AssetOtherDocument, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	DeleteAssetOtherDocuments(documentIds []*uint, tx ...utility.DatabaseTransaction) error
}

type AssetOtherDocumentRepository struct {
	db *gorm.DB
}

func NewAssetOtherDocumentRepository(db *gorm.DB) *AssetOtherDocumentRepository {
	return &AssetOtherDocumentRepository{db: db}
}

func (repo *AssetOtherDocumentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AssetOtherDocumentRepository) GetAllGroupedByCategory(assetRegisterID uint) (map[uint][]*dtos.AssetOtherDocumentListDTO, error) {
	var documents []*dtos.AssetOtherDocumentListDTO

	err := repo.db.Table("asset_other_documents").
		Select(`
			asset_other_documents.id,
			asset_other_documents.asset_register_id,
			asset_other_documents.asset_document_category_id,
			asset_document_categories.title_en AS asset_document_category_title_en,
			asset_document_categories.title_ar AS asset_document_category_title_ar,
			asset_other_documents.url,
			asset_other_documents.name,
			asset_other_documents.size,
			asset_other_documents.token,
			asset_other_documents.original_name,
			asset_other_documents.type,
			asset_other_documents.created_at,
			asset_other_documents.updated_at
		`).
		Joins("LEFT JOIN asset_document_categories ON asset_document_categories.id = asset_other_documents.asset_document_category_id").
		Where("asset_other_documents.asset_register_id = ?", assetRegisterID).
		Order("asset_other_documents.created_at DESC").
		Find(&documents).Error

	if err != nil {
		return nil, err
	}

	// Group by asset_document_category_id
	grouped := make(map[uint][]*dtos.AssetOtherDocumentListDTO)
	for _, doc := range documents {
		grouped[doc.AssetDocumentCategoryID] = append(grouped[doc.AssetDocumentCategoryID], doc)
	}

	return grouped, nil
}

func (repo *AssetOtherDocumentRepository) GetByID(id uint) (*models.AssetOtherDocument, error) {
	var document *models.AssetOtherDocument
	result := repo.db.First(&document, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return document, nil
}

func (repo *AssetOtherDocumentRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AssetOtherDocument, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var document *models.AssetOtherDocument
	err := db.Preload("AssetRegister").
		Preload("AssetDocumentCategory").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Where("id = ?", id).
		First(&document).Error

	if err != nil {
		fmt.Printf("Error fetching full details for AssetOtherDocument id %d: %v\n", id, err)
		return nil, err
	}
	return document, nil
}

func (repo *AssetOtherDocumentRepository) Create(doc *models.AssetOtherDocument, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(doc).Error; err != nil {
		return 0, err
	}
	return doc.ID, nil
}
func (repo *AssetOtherDocumentRepository) CreateAssetOtherDocuments(documents []*models.AssetOtherDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&documents).Error; err != nil {
		return err
	}
	return nil
}

func (repo *AssetOtherDocumentRepository) DeleteAssetOtherDocuments(documentIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", documentIds).Delete(&models.AssetOtherDocument{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AssetOtherDocumentRepository) Update(doc *models.AssetOtherDocument) error {
	if err := repo.db.Save(doc).Error; err != nil {
		return err
	}
	return nil
}

func (repo *AssetOtherDocumentRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var doc *models.AssetOtherDocument
	if err := db.First(&doc, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, id)
		}
		return err
	}
	return repo.db.Delete(&models.AssetOtherDocument{}, id).Error
}
