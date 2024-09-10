package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type KRIValueUpdateRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(kriRegisterID uint, params *riskDTOs.KRIValueUpdateListQueryParamsDTO) (*riskDTOs.KRIValueUpdateListResponseDTO, error)
	Create(kriValueUpdate *models.KRIValueUpdate, tx ...utility.DatabaseTransaction) (uint, error)
	GetByID(id uint) (*models.KRIValueUpdate, error)
	Update(kriValueUpdate *models.KRIValueUpdate, tx ...utility.DatabaseTransaction) error
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.KRIValueUpdate, error)
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	//Document
	CreateDocuments(docs []*models.KRIValueUpdateDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.KRIValueUpdateDocument, error)
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error

	GetNextKRIValueUpdate(kriRegisterID uint, deletedIDs []uint, tx ...utility.DatabaseTransaction) (*models.KRIValueUpdate, error)
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
			risk_levels.title_en AS risk_level_title_en,
			risk_levels.title_ar AS risk_level_title_ar,
			risk_levels.color_code AS risk_level_color_code,
			created_user.id AS created_user_id,
        	created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.email AS created_user_email,
			created_user.profile_image_path AS created_user_profile_image_path,
			updated_user.id AS updated_user_id,
    		updated_user.first_name_en AS updated_user_first_name_en,
			updated_user.first_name_ar AS updated_user_first_name_ar,
			updated_user.last_name_en AS updated_user_last_name_en,
			updated_user.last_name_ar AS updated_user_last_name_ar,
			updated_user.initial AS updated_user_initial,
			updated_user.color AS updated_user_color,
			updated_user.email AS updated_user_email,
			updated_user.profile_image_path AS updated_user_profile_image_path       
        `).
		Joins(`
		    LEFT JOIN kri_register ON kri_register.id = kri_value_updates.kri_register_id
			LEFT JOIN risk_levels ON risk_levels.id = kri_register.risk_level_id
			LEFT JOIN users AS created_user ON created_user.id = kri_register.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = kri_register.updated_by_id
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
			kri_value_updates.high_threshold LIKE ? OR
			risk_levels.title_en LIKE ? OR
			risk_levels.title_ar LIKE ? OR
			kri_value_updates.remarks LIKE ? )`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "kri_value_updates.id" // Default sort by id if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&kriValueUpdates).Error
	if err != nil {
		return nil, err
	}

	// Fetch documents for each KRIValueUpdate
	for _, valueUpdate := range kriValueUpdates {
		var documents []*riskDTOs.KRIValueUpdateDocumentDTO
		err := repo.db.Model(&models.KRIValueUpdateDocument{}).
			Where("kri_value_update_id = ?", valueUpdate.ID).
			Find(&documents).Error
		if err != nil {
			return nil, err
		}
		valueUpdate.Documents = documents
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

	result := repo.db.First(&kriValueUpdate, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return kriValueUpdate, nil

}
func (repo *KRIValueUpdateRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.KRIValueUpdate, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var kriValueUpdate *models.KRIValueUpdate

	result := db.Preload("Documents").First(&kriValueUpdate, id)
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

func (repo *KRIValueUpdateRepository) CreateDocuments(docs []*models.KRIValueUpdateDocument, tx ...utility.DatabaseTransaction) error {
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

func (repo *KRIValueUpdateRepository) GetDocumentByID(id uint) (*models.KRIValueUpdateDocument, error) {
	var document models.KRIValueUpdateDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *KRIValueUpdateRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
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

func (repo *KRIValueUpdateRepository) GetNextKRIValueUpdate(kriRegisterID uint, deletedIDs []uint, tx ...utility.DatabaseTransaction) (*models.KRIValueUpdate, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var nextValueUpdate models.KRIValueUpdate
	query := `
		SELECT TOP 1 *
		FROM kri_value_updates
		WHERE kri_register_id = ? AND id NOT IN (?)
		ORDER BY id DESC
	`

	err := db.Raw(query, kriRegisterID, deletedIDs).Scan(&nextValueUpdate).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return &nextValueUpdate, nil
}
