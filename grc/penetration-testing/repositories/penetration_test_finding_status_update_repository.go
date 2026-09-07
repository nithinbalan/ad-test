package repositories

import (
	"jethur/grc/models"
	penetrationTestDTOs "jethur/grc/penetration-testing/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type PenetrationTestFindingStatusUpdateRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(findingID uint, sort string, sortBy string, keyword string) ([]*penetrationTestDTOs.PenetrationTestFindingStatusUpdateListDTO, error)
	GetByID(findingID uint, id uint, tx ...utility.DatabaseTransaction) (*models.PenetrationTestFindingStatusUpdate, error)
	GetLatest(findingID uint) (*penetrationTestDTOs.PenetrationTestFindingStatusUpdateListDTO, error)
	Create(value *models.PenetrationTestFindingStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error)
	Update(value *models.PenetrationTestFindingStatusUpdate, tx ...utility.DatabaseTransaction) error
	CreateDocuments(docs []*models.PenetrationTestFindingStatusUpdateDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.PenetrationTestFindingStatusUpdateDocument, error)
	DeleteDocuments(docIDs []*uint, tx ...utility.DatabaseTransaction) error
}

type PenetrationTestFindingStatusUpdateRepository struct {
	db *gorm.DB
}

func NewPenetrationTestFindingStatusUpdateRepository(db *gorm.DB) *PenetrationTestFindingStatusUpdateRepository {
	return &PenetrationTestFindingStatusUpdateRepository{db: db}
}

func (repo *PenetrationTestFindingStatusUpdateRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *PenetrationTestFindingStatusUpdateRepository) GetAll(findingID uint, sort string, sortBy string, keyword string) ([]*penetrationTestDTOs.PenetrationTestFindingStatusUpdateListDTO, error) {
	var values []*penetrationTestDTOs.PenetrationTestFindingStatusUpdateListDTO
	query := repo.db.Table("penetration_test_finding_status_updates").Select(`
		penetration_test_finding_status_updates.*,
		action_plan_statuses.title_en AS action_plan_status_title_en,
		action_plan_statuses.title_ar AS action_plan_status_title_ar,
		action_plan_statuses.color_code AS action_plan_status_color_code,
		action_plan_statuses.type AS action_plan_status_type,
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
	`).Joins("LEFT JOIN action_plan_statuses ON action_plan_statuses.id = penetration_test_finding_status_updates.action_plan_status_id").Joins("LEFT JOIN users AS created_user ON created_user.id = penetration_test_finding_status_updates.created_by_id").Joins("LEFT JOIN users AS updated_user ON updated_user.id = penetration_test_finding_status_updates.updated_by_id").Where("penetration_test_finding_id = ?", findingID)
	if keyword != "" {
		query = query.Where("comment LIKE ? OR action_plan_statuses.title_en LIKE ? OR action_plan_statuses.title_ar LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "created_at"
	}
	allowedSortColumns := map[string]bool{"created_at": true, "percentage": true, "action_plan_status_id": true}
	if !allowedSortColumns[sortBy] {
		sortBy = "created_at"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	result := query.Order("penetration_test_finding_status_updates." + sortBy + " " + sort).Find(&values)
	for _, value := range values {
		var documents []*models.PenetrationTestFindingStatusUpdateDocument
		repo.db.Where("penetration_test_finding_status_update_id = ?", value.ID).Find(&documents)
		value.Documents = documents
	}
	return values, result.Error
}

func (repo *PenetrationTestFindingStatusUpdateRepository) GetByID(findingID uint, id uint, tx ...utility.DatabaseTransaction) (*models.PenetrationTestFindingStatusUpdate, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var value models.PenetrationTestFindingStatusUpdate
	if err := db.Preload("ActionPlanStatus").Preload("CreatedByUser").Preload("UpdatedByUser").Preload("Documents").Where("penetration_test_finding_id = ? AND id = ?", findingID, id).First(&value).Error; err != nil {
		return nil, err
	}
	return &value, nil
}

func (repo *PenetrationTestFindingStatusUpdateRepository) GetLatest(findingID uint) (*penetrationTestDTOs.PenetrationTestFindingStatusUpdateListDTO, error) {
	var value penetrationTestDTOs.PenetrationTestFindingStatusUpdateListDTO
	result := repo.db.Table("penetration_test_finding_status_updates").Select(`
		penetration_test_finding_status_updates.*,
		action_plan_statuses.title_en AS action_plan_status_title_en,
		action_plan_statuses.title_ar AS action_plan_status_title_ar,
		action_plan_statuses.color_code AS action_plan_status_color_code,
		action_plan_statuses.type AS action_plan_status_type,
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
	`).Joins("LEFT JOIN action_plan_statuses ON action_plan_statuses.id = penetration_test_finding_status_updates.action_plan_status_id").Joins("LEFT JOIN users AS created_user ON created_user.id = penetration_test_finding_status_updates.created_by_id").Joins("LEFT JOIN users AS updated_user ON updated_user.id = penetration_test_finding_status_updates.updated_by_id").Where("penetration_test_finding_id = ?", findingID).Order("penetration_test_finding_status_updates.created_at DESC").First(&value)
	return &value, result.Error
}

func (repo *PenetrationTestFindingStatusUpdateRepository) Create(value *models.PenetrationTestFindingStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(value).Error; err != nil {
		return 0, err
	}
	return value.ID, nil
}

func (repo *PenetrationTestFindingStatusUpdateRepository) Update(value *models.PenetrationTestFindingStatusUpdate, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(value).Error
}

func (repo *PenetrationTestFindingStatusUpdateRepository) CreateDocuments(docs []*models.PenetrationTestFindingStatusUpdateDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&docs).Error
}

func (repo *PenetrationTestFindingStatusUpdateRepository) GetDocumentByID(id uint) (*models.PenetrationTestFindingStatusUpdateDocument, error) {
	var document models.PenetrationTestFindingStatusUpdateDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *PenetrationTestFindingStatusUpdateRepository) DeleteDocuments(docIDs []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("id IN (?)", docIDs).Delete(&models.PenetrationTestFindingStatusUpdateDocument{}).Error
}
