package repositories

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	riskDTOs "jethur/grc/strategy/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type KPIValueUpdateRepositoryInterface interface {
	SetAllIsLatestFalse(kpiRegisterID uint, tx ...utility.DatabaseTransaction) error
	BeginTransaction() utility.DatabaseTransaction
	GetAll(kpiRegisterID uint, params *riskDTOs.KPIValueUpdateListQueryParamsDTO) (*riskDTOs.KPIValueUpdateListResponseDTO, error)
	Create(kpiValueUpdate *models.KPIValueUpdate, tx ...utility.DatabaseTransaction) (uint, error)
	GetByID(id uint) (*models.KPIValueUpdate, error)
	Update(kpiValueUpdate *models.KPIValueUpdate, tx ...utility.DatabaseTransaction) error
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.KPIValueUpdate, error)
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	//Document
	CreateDocuments(docs []*models.KPIValueUpdateDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.KPIValueUpdateDocument, error)
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error

	CreateKPIValueUpdateInputData(inputs []*models.KPIValueUpdateDataInput, tx ...utility.DatabaseTransaction) error
	DeleteKPIValueUpdateInputData(kpiRegisterValueUpdateId uint, tx ...utility.DatabaseTransaction) error

	GetNextKPIValueUpdate(kpiRegisterID uint, deletedIDs []uint, tx ...utility.DatabaseTransaction) (*models.KPIValueUpdate, error)
}

type KPIValueUpdateRepository struct {
	db *gorm.DB
}

func NewKPIValueUpdateRepository(db *gorm.DB) *KPIValueUpdateRepository {
	return &KPIValueUpdateRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *KPIValueUpdateRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *KPIValueUpdateRepository) GetAll(kpiRegisterID uint, params *riskDTOs.KPIValueUpdateListQueryParamsDTO) (*riskDTOs.KPIValueUpdateListResponseDTO, error) {
	var kpiValueUpdates []*riskDTOs.KPIValueUpdateListDTO
	var totalCount int64

	query := repo.db.Table("kpi_value_updates").
		Select(`
            kpi_value_updates.*,
			CONVERT(VARCHAR, kpi_value_updates.reviewed_date, 23) AS reviewed_date,
			kpi_register.title AS kpi_register,
			review_user.id AS review_user_id,
			review_user.first_name_en AS review_user_first_name_en,
			review_user.first_name_ar AS review_user_first_name_ar,
			review_user.last_name_en AS review_user_last_name_en,
			review_user.last_name_ar AS review_user_last_name_ar,
			review_user.initial AS review_user_initial,
			review_user.color AS review_user_color,
			review_user.email AS review_user_email,
			review_user.profile_image_path AS review_user_profile_image_path,
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
		    LEFT JOIN kpi_register ON kpi_register.id = kpi_value_updates.kpi_register_id
			LEFT JOIN users AS review_user ON review_user.id = kpi_value_updates.review_user_id
			LEFT JOIN users AS created_user ON created_user.id = kpi_value_updates.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = kpi_value_updates.updated_by_id
		`).
		Where("kpi_value_updates.kpi_register_id = ?", kpiRegisterID)

	// Apply LastReviewedDate filters
	if params.ReviewedDateFrom != "" {
		query = query.Where("kpi_value_updates.reviewed_date >= ?", params.ReviewedDateFrom)
	}
	if params.ReviewedDateTo != "" {
		query = query.Where("kpi_value_updates.reviewed_date <= ?", params.ReviewedDateTo)
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		query = query.Where(`
			(kpi_register.title LIKE ? OR
			kpi_value_updates.value LIKE ? OR
			kpi_value_updates.remarks LIKE ? )`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "kpi_value_updates.id" // Default sort by id if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&kpiValueUpdates).Error
	if err != nil {
		return nil, err
	}

	// Fetch documents for each KPIValueUpdate
	for _, valueUpdate := range kpiValueUpdates {
		var documents []*riskDTOs.KPIValueUpdateDocumentDTO
		err := repo.db.Model(&models.KPIValueUpdateDocument{}).
			Where("kpi_value_update_id = ?", valueUpdate.ID).
			Find(&documents).Error
		if err != nil {
			return nil, err
		}
		valueUpdate.Documents = documents
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &riskDTOs.KPIValueUpdateListResponseDTO{
		Items: kpiValueUpdates,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *KPIValueUpdateRepository) Create(kpiValueUpdate *models.KPIValueUpdate, tx ...utility.DatabaseTransaction) (uint, error) {

	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(kpiValueUpdate).Error; err != nil {
		return 0, err
	}

	return kpiValueUpdate.ID, nil
}

func (repo *KPIValueUpdateRepository) GetByID(id uint) (*models.KPIValueUpdate, error) {
	var kpiValueUpdate *models.KPIValueUpdate

	result := repo.db.First(&kpiValueUpdate, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return kpiValueUpdate, nil

}
func (repo *KPIValueUpdateRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.KPIValueUpdate, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var kpiValueUpdate *models.KPIValueUpdate

	result := db.Preload("Documents").Preload("Unit").Preload("DataInputs").Preload("ReviewUser").First(&kpiValueUpdate, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return kpiValueUpdate, nil

}

func (repo *KPIValueUpdateRepository) Update(kpiValueUpdate *models.KPIValueUpdate, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(kpiValueUpdate).Error; err != nil {
		return err
	}
	return nil
}

func (repo *KPIValueUpdateRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Get the record being deleted
	var kpiValueUpdate models.KPIValueUpdate
	if err := db.First(&kpiValueUpdate, id).Error; err != nil {
		return err
	}

	// Get the KPI register ID before deletion
	kpiRegisterID := kpiValueUpdate.KPIRegisterID

	// Delete the record
	result := db.Delete(&models.KPIValueUpdate{}, id)
	if result.Error != nil {
		return result.Error
	}

	// If the deleted record was marked as latest, update the most recent remaining record
	if kpiValueUpdate.IsLatest {
		var latestUpdate models.KPIValueUpdate
		err := db.Where("kpi_register_id = ? AND id != ?", kpiRegisterID, id).
			Order("created_at DESC").
			First(&latestUpdate).Error

		if err == nil {
			// Update the most recent record to be marked as latest
			err = db.Model(&models.KPIValueUpdate{}).
				Where("id = ?", latestUpdate.ID).
				Update("is_latest", true).Error
			if err != nil {
				return fmt.Errorf("failed to update latest record: %w", err)
			}
		}
	}

	return nil
}

func (repo *KPIValueUpdateRepository) CreateDocuments(docs []*models.KPIValueUpdateDocument, tx ...utility.DatabaseTransaction) error {
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

func (repo *KPIValueUpdateRepository) GetDocumentByID(id uint) (*models.KPIValueUpdateDocument, error) {
	var document models.KPIValueUpdateDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *KPIValueUpdateRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete documents by IDs
	if err := db.Where("id IN (?)", docIds).Delete(&models.KPIValueUpdateDocument{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *KPIValueUpdateRepository) GetNextKPIValueUpdate(kpiRegisterID uint, deletedIDs []uint, tx ...utility.DatabaseTransaction) (*models.KPIValueUpdate, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var nextValueUpdate models.KPIValueUpdate
	query := `
		SELECT TOP 1 *
		FROM kpi_value_updates
		WHERE kpi_register_id = ? AND id NOT IN (?)
		ORDER BY id DESC
	`

	err := db.Raw(query, kpiRegisterID, deletedIDs).Scan(&nextValueUpdate).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return &nextValueUpdate, nil
}

func (repo *KPIValueUpdateRepository) SetAllIsLatestFalse(kpiRegisterID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Model(&models.KPIValueUpdate{}).
		Where("kpi_register_id = ? AND is_latest = 1", kpiRegisterID).
		Update("is_latest", false).Error
}

func (repo *KPIValueUpdateRepository) CreateKPIValueUpdateInputData(inputDatas []*models.KPIValueUpdateDataInput, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(&inputDatas).Error; err != nil {
		return err
	}

	return nil
}

func (repo *KPIValueUpdateRepository) DeleteKPIValueUpdateInputData(valueUpdateId uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("kpi_value_update_id = ?", valueUpdateId).
		Delete(&models.KPIValueUpdateDataInput{}).Error; err != nil {
		return err
	}

	return nil
}
