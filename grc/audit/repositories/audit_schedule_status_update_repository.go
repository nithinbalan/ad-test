package repositories

import (
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditScheduleStatusUpdateRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(auditScheduleID uint, sort string, sortBy string, keyword string) ([]*auditDTOs.AuditScheduleStatusUpdateListDTO, error)
	GetByID(id uint) (*models.AuditScheduleStatusUpdate, error)
	Create(auditScheduleStatusUpdate *models.AuditScheduleStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error)
	Update(auditScheduleStatusUpdate *models.AuditScheduleStatusUpdate, tx ...utility.DatabaseTransaction) error
	GetLatestAuditScheduleStatusUpdate(auditScheduleID uint) (*auditDTOs.AuditScheduleStatusUpdateListDTO, error)

	// Document methods
	CreateDocuments(docs []*models.AuditScheduleStatusUpdateDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.AuditScheduleStatusUpdateDocument, error)
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error
}

type AuditScheduleStatusUpdateRepository struct {
	db *gorm.DB
}

func NewAuditScheduleStatusUpdateRepository(db *gorm.DB) *AuditScheduleStatusUpdateRepository {
	return &AuditScheduleStatusUpdateRepository{db: db}
}

// BeginTransaction starts a new transaction.
func (repo *AuditScheduleStatusUpdateRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AuditScheduleStatusUpdateRepository) GetAll(auditScheduleID uint, sort string, sortBy string, keyword string) ([]*auditDTOs.AuditScheduleStatusUpdateListDTO, error) {
	var updates []*auditDTOs.AuditScheduleStatusUpdateListDTO

	query := repo.db.Table("audit_schedule_status_updates").
		Select(`
            audit_schedule_status_updates.*,
            audit_schedules.title AS audit_schedule_title,
            audit_schedule_statuses.title_en AS audit_schedule_status_title_en,
            audit_schedule_statuses.title_ar AS audit_schedule_status_title_ar,
            audit_schedule_statuses.color_code AS audit_schedule_status_color_code,
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
            LEFT JOIN audit_schedules ON audit_schedules.id = audit_schedule_status_updates.audit_schedule_id
            LEFT JOIN audit_schedule_statuses ON audit_schedule_statuses.id = audit_schedule_status_updates.audit_schedule_status_id
            LEFT JOIN users AS created_user ON created_user.id = audit_schedule_status_updates.created_by_id
            LEFT JOIN users AS updated_user ON updated_user.id = audit_schedule_status_updates.updated_by_id
        `).
		Where("audit_schedule_status_updates.audit_schedule_id = ?", auditScheduleID)

	if keyword != "" {
		query = query.Where(`audit_schedule_status_updates.comment LIKE ?`, "%"+keyword+"%")
	}

	if sortBy == "" {
		sortBy = "audit_schedule_status_updates.id"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	query = query.Order(sortBy + " " + sort)

	err := query.Find(&updates).Error

	// Fetching documents for each audit schedule status updates
	for _, update := range updates {
		var documents []*models.AuditScheduleStatusUpdateDocument
		repo.db.Model(&models.AuditScheduleStatusUpdateDocument{}).Where("audit_schedule_status_update_id = ?", update.ID).Find(&documents)
		update.Documents = documents
	}
	return updates, err
}

func (repo *AuditScheduleStatusUpdateRepository) GetByID(id uint) (*models.AuditScheduleStatusUpdate, error) {
	var update models.AuditScheduleStatusUpdate
	result := repo.db.Preload("Documents").First(&update, id)
	return &update, result.Error
}

func (repo *AuditScheduleStatusUpdateRepository) Create(auditScheduleStatusUpdate *models.AuditScheduleStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(auditScheduleStatusUpdate).Error; err != nil {
		return 0, err
	}
	return auditScheduleStatusUpdate.ID, nil
}

func (repo *AuditScheduleStatusUpdateRepository) Update(auditScheduleStatusUpdate *models.AuditScheduleStatusUpdate, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(auditScheduleStatusUpdate).Error
}

func (repo *AuditScheduleStatusUpdateRepository) GetLatestAuditScheduleStatusUpdate(auditScheduleID uint) (*auditDTOs.AuditScheduleStatusUpdateListDTO, error) {
	var update *auditDTOs.AuditScheduleStatusUpdateListDTO
	query := repo.db.Table("audit_schedule_status_updates").
		Select(`
            audit_schedule_status_updates.*,
            audit_schedules.title AS audit_schedule_title,
            audit_schedule_statuses.title_en AS audit_schedule_status_title_en,
            audit_schedule_statuses.title_ar AS audit_schedule_status_title_ar,
            audit_schedule_statuses.color_code AS audit_schedule_status_color_code,
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
            LEFT JOIN audit_schedules ON audit_schedules.id = audit_schedule_status_updates.audit_schedule_id
            LEFT JOIN audit_schedule_statuses ON audit_schedule_statuses.id = audit_schedule_status_updates.audit_schedule_status_id
            LEFT JOIN users AS created_user ON created_user.id = audit_schedule_status_updates.created_by_id
            LEFT JOIN users AS updated_user ON updated_user.id = audit_schedule_status_updates.updated_by_id
        `).
		Where("audit_schedule_status_updates.audit_schedule_id = ?", auditScheduleID).
		Order("audit_schedule_status_updates.created_at DESC").
		First(&update)

	if query.Error != nil {
		return nil, query.Error
	}

	return update, nil
}

func (repo *AuditScheduleStatusUpdateRepository) CreateDocuments(docs []*models.AuditScheduleStatusUpdateDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&docs).Error
}

func (repo *AuditScheduleStatusUpdateRepository) GetDocumentByID(id uint) (*models.AuditScheduleStatusUpdateDocument, error) {
	var document models.AuditScheduleStatusUpdateDocument
	err := repo.db.Where("id = ?", id).First(&document).Error
	return &document, err
}

func (repo *AuditScheduleStatusUpdateRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("id IN (?)", docIds).Delete(&models.AuditScheduleStatusUpdateDocument{}).Error
}
