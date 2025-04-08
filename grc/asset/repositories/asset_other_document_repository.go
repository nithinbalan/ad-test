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
	GetAll(assetRegisterID uint, keyword string) ([]*dtos.AssetOtherDocumentListDTO, error)
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

func (repo *AssetOtherDocumentRepository) GetAll(assetRegisterID uint, keyword string) ([]*dtos.AssetOtherDocumentListDTO, error) {
	documents := []*dtos.AssetOtherDocumentListDTO{}

	query := repo.db.Table("asset_other_documents").
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
			created_by.id AS created_by_id,
			created_by.first_name_en AS created_by_first_name_en,
			created_by.first_name_ar AS created_by_first_name_ar,
			created_by.last_name_en AS created_by_last_name_en,
			created_by.last_name_ar AS created_by_last_name_ar,
			created_by.initial AS created_by_initial,
			created_by.color AS created_by_color,
			created_by.email AS created_by_email,
			created_by.profile_image_path AS created_by_profile_image_path,
			asset_other_documents.created_at,
			updated_by.id AS updated_by_id,
			updated_by.first_name_en AS updated_by_first_name_en,
			updated_by.first_name_ar AS updated_by_first_name_ar,
			updated_by.last_name_en AS updated_by_last_name_en,
			updated_by.last_name_ar AS updated_by_last_name_ar,
			updated_by.initial AS updated_by_initial,
			updated_by.color AS updated_by_color,
			updated_by.email AS updated_by_email,
			updated_by.profile_image_path AS updated_by_profile_image_path,
			asset_other_documents.updated_at
		`).
		Joins("LEFT JOIN asset_document_categories ON asset_document_categories.id = asset_other_documents.asset_document_category_id").
		Joins("LEFT JOIN users AS created_by ON created_by.id = asset_other_documents.created_by_id").
		Joins("LEFT JOIN users AS updated_by ON updated_by.id = asset_other_documents.updated_by_id").
		Order("asset_document_categories.title_en ASC")

	if keyword != "" {
		query = query.Where(`
			asset_other_documents.asset_register_id = ? AND 
			(asset_document_categories.title_en LIKE ? OR asset_document_categories.title_ar LIKE ?)`,
			assetRegisterID, "%"+keyword+"%", "%"+keyword+"%")
	} else {
		query = query.Where("asset_other_documents.asset_register_id = ?", assetRegisterID)
	}

	err := query.Find(&documents).Error
	return documents, err
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
