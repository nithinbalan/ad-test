package repositories

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	kpiDTOs "jethur/grc/kpi-register/dtos"
	"jethur/grc/models"

	"jethur/utility"

	"gorm.io/gorm"
)

type KPIValueUpdateRepositoryInterface interface {
	SetAllIsLatestFalse(kpiRegisterID uint, tx ...utility.DatabaseTransaction) error
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *kpiDTOs.KPIValueUpdateListQueryParamsDTO, currentUserID uint) (*kpiDTOs.KPIValueUpdateListResponseDTO, error)
	GetByID(id uint) (*models.KPIValueUpdate, error)
	Update(kpiValueUpdate *models.KPIValueUpdate, tx ...utility.DatabaseTransaction) error
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.KPIValueUpdate, error)
	//Document
	CreateDocuments(docs []*models.KPIValueUpdateDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.KPIValueUpdateDocument, error)
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error

	CreateKPIValueUpdateInputData(inputs []*models.KPIValueUpdateDataInput, tx ...utility.DatabaseTransaction) error
	DeleteKPIValueUpdateInputData(kpiRegisterValueUpdateId uint, tx ...utility.DatabaseTransaction) error

	GetNextKPIValueUpdate(kpiRegisterID uint, deletedIDs []uint, tx ...utility.DatabaseTransaction) (*models.KPIValueUpdate, error)
	GetStrategicInitiativeIDsByKPIRegisterID(kpiRegisterID uint, tx ...utility.DatabaseTransaction) ([]uint, error)
	GetLatestValueByKPIRegisterID(kpiRegisterID uint, tx ...utility.DatabaseTransaction) (*float64, error)
	CreateKPIValueUpdateHistory(history *models.KPIValueUpdateHistory, tx ...utility.DatabaseTransaction) error
	DeleteKPIValueUpdateHistories(kpiValueUpdateID uint, tx ...utility.DatabaseTransaction) error
	CreateKPIValueUpdateHistoryInputs(inputs []*models.KPIValueUpdateHistoryInput, tx ...utility.DatabaseTransaction) error
	DeleteKPIValueUpdateHistoryInputsByKPIValueUpdateID(kpiValueUpdateID uint, tx ...utility.DatabaseTransaction) error

	CreateKPIValueUpdateReview(review *models.KPIValueUpdateReview, tx ...utility.DatabaseTransaction) error
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
func (repo *KPIValueUpdateRepository) GetAll(params *kpiDTOs.KPIValueUpdateListQueryParamsDTO, currentUserID uint) (*kpiDTOs.KPIValueUpdateListResponseDTO, error) {
	var kpiValueUpdates []*kpiDTOs.KPIValueUpdateListDTO
	var totalCount int64

	query := repo.db.Table("kpi_value_updates").
		Select(`
			kpi_value_updates.*,
			CONVERT(VARCHAR, kpi_value_updates.reviewed_date, 23) AS reviewed_date,
			CONVERT(VARCHAR, kpi_value_updates.start_date, 23) AS start_date,
			CONVERT(VARCHAR, kpi_value_updates.end_date, 23) AS end_date,
			kpi_register.title AS kpi_register_title,
			kpi_categories.title_en AS kpi_category_title_en,
			kpi_categories.title_ar AS kpi_category_title_ar,
			review_user.id AS review_user_id,
			review_user.first_name_en AS review_user_first_name_en,
			review_user.first_name_ar AS review_user_first_name_ar,
			review_user.last_name_en AS review_user_last_name_en,
			review_user.last_name_ar AS review_user_last_name_ar,
			review_user.initial AS review_user_initial,
			review_user.color AS review_user_color,
			review_user.email AS review_user_email,
			review_user.profile_image_path AS review_user_profile_image_path,
			kpi_review_statuses.id AS kpi_review_status_id,
			kpi_review_statuses.title_en AS kpi_review_status_title_en,
			kpi_review_statuses.title_ar AS kpi_review_status_title_ar,
			kpi_review_statuses.type AS kpi_review_status_type,
			kpi_review_statuses.color_code AS kpi_review_status_color_code,
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
			LEFT JOIN kpi_categories ON kpi_categories.id = kpi_register.kpi_category_id
			LEFT JOIN kpi_review_statuses ON kpi_review_statuses.id = kpi_value_updates.kpi_review_status_id
			LEFT JOIN users AS review_user ON review_user.id = kpi_value_updates.review_user_id
			LEFT JOIN users AS created_user ON created_user.id = kpi_value_updates.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = kpi_value_updates.updated_by_id
		`)

	if params.IsPassiveUser {
		query = query.Where("(kpi_value_updates.review_user_id = ? OR kpi_register.review_user_id = ?)", currentUserID, currentUserID)
	}

	// APPLY DYNAMIC FILTERS
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			switch columnName {
			case "kpi_register_id":
				query = query.Where("kpi_value_updates.kpi_register_id IN (?)", values)
			case "kpi_category_id":
				query = query.Where("kpi_register.kpi_category_id IN (?)", values)
			case "responsible_user_id":
				query = query.Where("kpi_register.responsible_user_id IN (?)", values)
			case "review_user_id":
				query = query.Where("kpi_value_updates.review_user_id IN (?)", values)
			case "kpi_review_status_id":
				query = query.Where("kpi_value_updates.kpi_review_status_id IN (?)", values)
			default:
				query = query.Where(columnName+" IN (?)", values)
			}
		}
	}

	// APPLY KEYWORD SEARCH
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
		(kpi_register.title LIKE ? OR
		kpi_categories.title_en LIKE ? OR
		kpi_categories.title_ar LIKE ? OR
		CAST(kpi_value_updates.value AS NVARCHAR(50)) LIKE ? OR
		kpi_value_updates.remarks LIKE ? OR
		CAST(kpi_value_updates.reviewed_date AS NVARCHAR(30)) LIKE ? OR
		kpi_review_statuses.title_en LIKE ? OR
		kpi_review_statuses.title_ar LIKE ? OR
		kpi_review_statuses.type LIKE ? OR
		review_user.first_name_en LIKE ? OR
		review_user.first_name_ar LIKE ? OR
		review_user.last_name_en LIKE ? OR
		review_user.last_name_ar LIKE ? OR
		review_user.email LIKE ? OR
		created_user.first_name_en LIKE ? OR
		created_user.first_name_ar LIKE ? OR
		created_user.last_name_en LIKE ? OR
		created_user.last_name_ar LIKE ? OR
		created_user.email LIKE ? )`,
			keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword,
		)
	}

	// QUICK FILTERS USING REVIEW STATUS TYPE
	for _, filter := range params.QuickFilter {
		switch filter {
		case "overdue":
			query = query.Where("kpi_review_statuses.type = ?", "overdue")
		case "awaiting-update":
			query = query.Where("kpi_review_statuses.type = ?", "awaiting-update")
		case "updated":
			query = query.Where("kpi_review_statuses.type = ?", "updated")
		}
	}

	// COUNT TOTAL RECORDS
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// APPLY SORTING
	if params.SortBy == "" {
		params.SortBy = "kpi_value_updates.id"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "desc"
	}
	query = query.Order(params.SortBy + " " + params.Sort)

	// APPLY PAGINATION
	if !params.Export && params.Page > 0 && params.PageSize > 0 {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	// EXECUTE QUERY
	if err := query.Find(&kpiValueUpdates).Error; err != nil {
		return nil, err
	}

	totalPages := 1
	if !params.Export && params.PageSize > 0 {
		totalPages = int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	}

	response := &kpiDTOs.KPIValueUpdateListResponseDTO{
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

	result := db.Preload("Documents").
		Preload("KPIRegister").
		Preload("KPIRegister.KPICategory").
		Preload("KPIRegister.KPIStatus").
		Preload("KPIRegister.KPIReviewStatus").
		Preload("KPIRegister.ResponsibleUser").
		Preload("KPIRegister.ReviewUser").
		Preload("KPIRegister.CalculationType").
		Preload("KPIRegister.ReviewFrequency").
		Preload("KPIRegister.DataInputs").
		Preload("KPIRegister.KPIAchievementLogic").
		Preload("Unit").
		Preload("DataInputs").
		Preload("KPIReviewStatus").
		Preload("ReviewUser").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("ValueUpdateHistories.CreatedByUser").
		Preload("ValueUpdateHistories.Unit").
		Preload("ValueUpdateHistories.KPIValueUpdateHistoryInputs.CreatedByUser").
		Preload("ValueUpdateReviews").
		Preload("ValueUpdateReviews.KPIReviewStatus").
		Preload("ValueUpdateReviews.CreatedByUser").
		First(&kpiValueUpdate, id)
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

func (repo *KPIValueUpdateRepository) GetStrategicInitiativeIDsByKPIRegisterID(kpiRegisterID uint, tx ...utility.DatabaseTransaction) ([]uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var initiativeIDs []uint
	err := db.
		Table("strategy_profile_strategic_initiative_kpis").
		Where("kpi_register_id = ?", kpiRegisterID).
		Pluck("strategy_profile_strategic_initiative_id", &initiativeIDs).Error

	return initiativeIDs, err
}

func (repo *KPIValueUpdateRepository) GetLatestValueByKPIRegisterID(kpiRegisterID uint, tx ...utility.DatabaseTransaction) (*float64, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var kpiValueUpdate models.KPIValueUpdate
	err := db.Where("kpi_register_id = ? AND is_latest = ?", kpiRegisterID, true).First(&kpiValueUpdate).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return kpiValueUpdate.Value, nil
}

func (repo *KPIValueUpdateRepository) CreateKPIValueUpdateHistory(history *models.KPIValueUpdateHistory, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(&history).Error; err != nil {
		return err
	}

	return nil
}

func (repo *KPIValueUpdateRepository) DeleteKPIValueUpdateHistories(kpiValueUpdateID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Delete(&models.KPIValueUpdateHistory{}, "kpi_value_update_id = ?", kpiValueUpdateID).Error; err != nil {
		return err
	}

	return nil
}

func (repo *KPIValueUpdateRepository) CreateKPIValueUpdateHistoryInputs(inputs []*models.KPIValueUpdateHistoryInput, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(&inputs).Error; err != nil {
		return err
	}

	return nil
}

func (repo *KPIValueUpdateRepository) DeleteKPIValueUpdateHistoryInputsByKPIValueUpdateID(kpiValueUpdateID uint, tx ...utility.DatabaseTransaction) error {

	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Exec(`
		DELETE hInputs
		FROM kpi_value_update_history_inputs hInputs
		INNER JOIN kpi_value_update_histories h
			ON h.id = hInputs.kpi_value_update_history_id
		WHERE h.kpi_value_update_id = ?
	`, kpiValueUpdateID).Error
}

func (repo *KPIValueUpdateRepository) CreateKPIValueUpdateReview(review *models.KPIValueUpdateReview, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(review).Error; err != nil {
		return err
	}

	return nil
}
