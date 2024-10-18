package repositories

import (
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditPlanCommencementLetterRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(auditPlanID uint, keyword string) ([]*auditDTOs.AuditPlanCommencementLetterListDTO, error)
	GetByID(id uint) (*models.AuditPlanCommencementLetter, error)
	GetFullDetailsByID(auditPlanID uint, letterID uint, tx ...utility.DatabaseTransaction) (*models.AuditPlanCommencementLetter, error)
	Create(auditPlanCommencementLetter *models.AuditPlanCommencementLetter, tx ...utility.DatabaseTransaction) (uint, error)
	Update(commencementLetter *models.AuditPlanCommencementLetter, tx ...utility.DatabaseTransaction) error
	Delete(id uint) error

	CreateToUser(toUser *models.AuditPlanCommencementLetterToUser, tx ...utility.DatabaseTransaction) error
	CreateCCUser(ccUser *models.AuditPlanCommencementLetterCCUser, tx ...utility.DatabaseTransaction) error
}

type AuditPlanCommencementLetterRepository struct {
	db *gorm.DB
}

func NewAuditPlanCommencementLetterRepository(db *gorm.DB) *AuditPlanCommencementLetterRepository {
	return &AuditPlanCommencementLetterRepository{
		db: db,
	}
}

func (repo *AuditPlanCommencementLetterRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AuditPlanCommencementLetterRepository) GetAll(auditPlanID uint, keyword string) ([]*auditDTOs.AuditPlanCommencementLetterListDTO, error) {
	var commencementLetters []*auditDTOs.AuditPlanCommencementLetterListDTO

	query := repo.db.Table("audit_plan_commencement_letters").
		Select(`
		audit_plan_commencement_letters.id,
		audit_plan_commencement_letters.audit_plan_id,
		organization_units.title_en AS responsible_unit_title_en,
		organization_units.title_ar AS responsible_unit_title_ar,
		STRING_AGG(to_user.email, ', ') AS to_user_emails,
		from_user.id AS from_user_id,
		from_user.first_name_en AS from_user_first_name_en,
		from_user.first_name_ar AS from_user_first_name_ar,
		from_user.last_name_en AS from_user_last_name_en,
		from_user.last_name_ar AS from_user_last_name_ar,
		from_user.initial AS from_user_initial,
		from_user.color AS from_user_color,
		from_user.email AS from_user_email,
		from_user.profile_image_path AS from_user_profile_image_path,
		audit_plan_commencement_letters.created_at,
		audit_plan_commencement_letters.updated_at,
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
		Joins(`LEFT JOIN organization_units ON organization_units.id = audit_plan_commencement_letters.responsible_unit_id`).
		Joins(`LEFT JOIN users AS from_user ON from_user.id = audit_plan_commencement_letters.from_user_id`).
		Joins(`LEFT JOIN users AS created_user ON created_user.id = audit_plan_commencement_letters.created_by_id`).
		Joins(`LEFT JOIN users AS updated_user ON updated_user.id = audit_plan_commencement_letters.updated_by_id`).
		Joins(`LEFT JOIN audit_plan_commencement_letter_to_user_ids ON audit_plan_commencement_letters.id = audit_plan_commencement_letter_to_user_ids.audit_plan_commencement_letter_id`).
		Joins(`LEFT JOIN users AS to_user ON to_user.id = audit_plan_commencement_letter_to_user_ids.user_id`).
		Where("audit_plan_commencement_letters.audit_plan_id = ?", auditPlanID).
		Group(`
		audit_plan_commencement_letters.id,
		audit_plan_commencement_letters.audit_plan_id,
		organization_units.title_en,
		organization_units.title_ar,
		from_user.id,
		from_user.first_name_en,
		from_user.first_name_ar,
		from_user.last_name_en,
		from_user.last_name_ar,
		from_user.initial,
		from_user.color,
		from_user.email,
		from_user.profile_image_path,
		audit_plan_commencement_letters.created_at,
		audit_plan_commencement_letters.updated_at,
		created_user.id,
		created_user.first_name_en,
		created_user.first_name_ar,
		created_user.last_name_en,
		created_user.last_name_ar,
		created_user.initial,
		created_user.color,
		created_user.email,
		created_user.profile_image_path,
		updated_user.id,
		updated_user.first_name_en,
		updated_user.first_name_ar,
		updated_user.last_name_en,
		updated_user.last_name_ar,
		updated_user.initial,
		updated_user.color,
		updated_user.email,
		updated_user.profile_image_path
	`)

	// Apply optional keyword filter
	if keyword != "" {
		search := "%" + keyword + "%"
		query = query.Where("organization_units.title_en LIKE ? OR organization_units.title_ar LIKE ?", search, search)
	}

	// Execute the query
	if err := query.Find(&commencementLetters).Error; err != nil {
		return nil, err
	}

	return commencementLetters, nil
}

func (repo *AuditPlanCommencementLetterRepository) GetByID(id uint) (*models.AuditPlanCommencementLetter, error) {
	var commencementLetter *models.AuditPlanCommencementLetter

	result := repo.db.First(&commencementLetter, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return commencementLetter, nil
}

func (repo *AuditPlanCommencementLetterRepository) GetFullDetailsByID(auditPlanID uint, letterID uint, tx ...utility.DatabaseTransaction) (*models.AuditPlanCommencementLetter, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var commencementLetter *models.AuditPlanCommencementLetter

	result := db.
		Preload("ToUser.User").
		Preload("CCUsers.User").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("ResponsibleUnit").
		Preload("FromUser.Designation").
		Where("audit_plan_id = ? AND id = ?", auditPlanID, letterID).
		First(&commencementLetter)

	if result.Error != nil {
		return nil, result.Error
	}

	return commencementLetter, nil
}

func (repo *AuditPlanCommencementLetterRepository) Create(auditPlanCommencementLetter *models.AuditPlanCommencementLetter, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(auditPlanCommencementLetter).Error; err != nil {
		return 0, err
	}

	return auditPlanCommencementLetter.ID, nil
}

func (repo *AuditPlanCommencementLetterRepository) Update(commencementLetter *models.AuditPlanCommencementLetter, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(commencementLetter).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditPlanCommencementLetterRepository) Delete(id uint) error {

	if err := repo.db.Delete(&models.AuditPlanCommencementLetter{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditPlanCommencementLetterRepository) CreateToUser(toUser *models.AuditPlanCommencementLetterToUser, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(toUser).Error; err != nil {
		return err
	}
	return nil
}

func (repo *AuditPlanCommencementLetterRepository) CreateCCUser(ccUser *models.AuditPlanCommencementLetterCCUser, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(ccUser).Error; err != nil {
		return err
	}
	return nil
}
