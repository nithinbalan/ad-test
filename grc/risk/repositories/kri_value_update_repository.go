package repositories

import (
	"jethur/config"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/utility"
	"os"
	"path/filepath"
	"strconv"

	"gorm.io/gorm"
)

type KRIValueUpdateRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(kriRegisterID uint, params *riskDTOs.KRIValueUpdateListQueryParamsDTO) (*riskDTOs.KRIValueUpdateListResponseDTO, error)
	Create(kriValueUpdate *models.KRIValueUpdate, tx ...utility.DatabaseTransaction) (uint, error)
	Update(kriValueUpdate *models.KRIValueUpdate, tx ...utility.DatabaseTransaction) error
	GetByID(id uint) (*models.KRIValueUpdate, error)
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	//Document
	SaveKRIValueUpdateDocuments(docs []*models.KRIValueUpdateDocument, tx ...utility.DatabaseTransaction) error
	GetKRIValueUpdateDocumentByID(id uint) (*models.KRIValueUpdateDocument, error)
	StoreKRIValueUpdateDocumentFile(kriRegisterID uint, updatesId uint, fileName string) error
	RemoveKRIValueUpdateDocumentFile(kriRegisterID uint, updatesId uint, fileNames string) error
	DeleteKRIValueUpdateDocumentID(docId uint, tx ...utility.DatabaseTransaction) error
	DeleteKRIValueUpdateDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error
}

type KRIValueUpdateRepository struct {
	db *gorm.DB
}

func NewKRIValueUpdateRepository(db *gorm.DB) *KRIValueUpdateRepository {
	return &KRIValueUpdateRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *KRIValueUpdateRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *KRIValueUpdateRepository) GetAll(kriRegisterID uint, params *riskDTOs.KRIValueUpdateListQueryParamsDTO) (*riskDTOs.KRIValueUpdateListResponseDTO, error) {

	var kriValueUpdates []*riskDTOs.KRIValueUpdateListDTO
	var totalCount int64

	query := repo.db.Table("kri_value_updates").
		Select(`
            kri_value_updates.*,
			kri_register.title AS kri_register,
			created_user.first_name_en AS created_user_first_name_en,
            created_user.first_name_ar AS created_user_first_name_ar,
            created_user.last_name_en AS created_user_last_name_en,
            created_user.last_name_ar AS created_user_last_name_ar,   
			updated_user.first_name_en AS updated_user_first_name_en,
            updated_user.first_name_ar AS updated_user_first_name_ar,
            updated_user.last_name_en AS updated_user_last_name_en,
            updated_user.last_name_ar AS updated_user_last_name_ar          
        `).
		Joins(`
		    LEFT JOIN kri_register ON kri_register.id = kri_value_updates.kri_register_id
			LEFT JOIN users AS created_user ON created_user.id = kri_register.created_by
			LEFT JOIN users AS updated_user ON updated_user.id = kri_register.updated_by
		`).
		Where("kri_value_updates.kri_register_id = ?", kriRegisterID)

	// Apply LastReviewedDate filters
	if params.ReviewedDateFrom != "" {
		query = query.Where("kri_value_updates.reviewed_date >= ?", params.ReviewedDateFrom)
	}
	if params.ReviewedDateTo != "" {
		query = query.Where("kri_value_updates.reviewed_date <= ?", params.ReviewedDateTo)
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		query = query.Where(`
			(kri_register.title LIKE ? OR
			kri_value_updates.value LIKE ? OR
			kri_value_updates.low_threshold LIKE ? OR
			kri_value_updates.medium_threshold LIKE ? OR
			organization_units.high_threshold LIKE ? OR
			organization_units.remarks LIKE ? )`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "kri_value_updates.id" // Default sort by title if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the queryc
	err := query.Offset(offset).Limit(params.PageSize).Find(&kriValueUpdates).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &riskDTOs.KRIValueUpdateListResponseDTO{
		Items: kriValueUpdates,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *KRIValueUpdateRepository) Create(kriValueUpdate *models.KRIValueUpdate, tx ...utility.DatabaseTransaction) (uint, error) {

	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(kriValueUpdate).Error; err != nil {
		return 0, err
	}

	return kriValueUpdate.ID, nil
}

func (repo *KRIValueUpdateRepository) GetByID(id uint) (*models.KRIValueUpdate, error) {
	var kriValueUpdate *models.KRIValueUpdate

	result := repo.db.Preload("KRIValueUpdateDocuments").First(&kriValueUpdate, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return kriValueUpdate, nil

}

func (repo *KRIValueUpdateRepository) Update(kriValueUpdate *models.KRIValueUpdate, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(kriValueUpdate).Error; err != nil {
		return err
	}
	return nil
}

func (repo *KRIValueUpdateRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var kriValueUpdate *models.KRIValueUpdate
	if err := db.First(&kriValueUpdate, id).Error; err != nil {
		return err
	}

	result := db.Delete(&models.KRIValueUpdate{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *KRIValueUpdateRepository) SaveKRIValueUpdateDocuments(docs []*models.KRIValueUpdateDocument, tx ...utility.DatabaseTransaction) error {
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

func (repo *KRIValueUpdateRepository) GetKRIValueUpdateDocumentByID(id uint) (*models.KRIValueUpdateDocument, error) {
	var document models.KRIValueUpdateDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *KRIValueUpdateRepository) StoreKRIValueUpdateDocumentFile(kriRegisterID uint, updatesId uint, fileName string) error {
	storagePath := config.GetEnv("STORAGE_PATH", "")
	temporaryStoragePath := filepath.Join(storagePath, "tmp", fileName)
	directoryPath := filepath.Join(storagePath, "risk", "kri-register", strconv.FormatUint(uint64(kriRegisterID), 10), "updates", strconv.FormatUint(uint64(updatesId), 10), "documents")
	destinationPath := filepath.Join(directoryPath, fileName)
	if err := os.MkdirAll(directoryPath, os.ModePerm); err != nil {
		return err
	}
	if err := os.Rename(temporaryStoragePath, destinationPath); err != nil {
		return err
	}
	return nil
}

func (repo *KRIValueUpdateRepository) RemoveKRIValueUpdateDocumentFile(kriRegisterID uint, updatesId uint, fileName string) error {
	storagePath := config.GetEnv("STORAGE_PATH", "")
	filePath := filepath.Join(storagePath, "risk", "kri-register", strconv.FormatUint(uint64(kriRegisterID), 10), "updates", strconv.FormatUint(uint64(updatesId), 10), "documents", fileName)
	if err := os.Remove(filePath); err != nil {
		return err
	}
	return nil
}

func (repo *KRIValueUpdateRepository) DeleteKRIValueUpdateDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete documents by IDs
	if err := db.Where("id IN (?)", docIds).Delete(&models.KRIValueUpdateDocument{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *KRIValueUpdateRepository) DeleteKRIValueUpdateDocumentID(docID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Delete(&models.KRIValueUpdateDocument{}, docID).Error; err != nil {
		return err
	}

	return nil
}
