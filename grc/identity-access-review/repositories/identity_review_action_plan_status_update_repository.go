package repositories

import (
	identityAccessReviewDTOs "jethur/grc/identity-access-review/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type IdentityReviewActionPlanStatusUpdateRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(identityReviewActionPlanID uint, sort string, sortBy string, keyword string) ([]*identityAccessReviewDTOs.IdentityReviewActionPlanStatusUpdateListDTO, error)
	GetByID(identityReviewActionPlanID uint, id uint) (*models.IdentityReviewActionPlanStatusUpdate, error)
	Create(identityReviewActionPlanStatusUpdate *models.IdentityReviewActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error)
	Update(identityReviewActionPlanStatusUpdate *models.IdentityReviewActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) error

	// Document
	CreateDocuments(docs []*models.IdentityReviewActionPlanStatusUpdateDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint, tx ...utility.DatabaseTransaction) (*models.IdentityReviewActionPlanStatusUpdateDocument, error)
	DeleteDocuments(docIds []uint, tx ...utility.DatabaseTransaction) error

	GetLatestIdentityReviewActionPlanStatusUpdate(identityReviewActionPlanID uint) (*models.IdentityReviewActionPlanStatusUpdate, error)
}

type IdentityReviewActionPlanStatusUpdateRepository struct {
	db *gorm.DB
}

func NewIdentityReviewActionPlanStatusUpdateRepository(db *gorm.DB) IdentityReviewActionPlanStatusUpdateRepositoryInterface {
	return &IdentityReviewActionPlanStatusUpdateRepository{
		db: db,
	}
}

func (repo *IdentityReviewActionPlanStatusUpdateRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *IdentityReviewActionPlanStatusUpdateRepository) GetAll(identityReviewActionPlanID uint, sort string, sortBy string, keyword string) ([]*identityAccessReviewDTOs.IdentityReviewActionPlanStatusUpdateListDTO, error) {
	var identityReviewActionPlanUpdates []*identityAccessReviewDTOs.IdentityReviewActionPlanStatusUpdateListDTO

	query := repo.db.Table("identity_review_action_plan_status_updates").
		Select(`
			identity_review_action_plan_status_updates.id,
			identity_review_action_plan_status_updates.identity_review_action_plan_id,
			identity_review_action_plans.title AS action_plan_title,
			identity_review_action_plan_status_updates.action_plan_status_id,
			action_plan_statuses.title_en AS action_plan_status_title_en,
			action_plan_statuses.title_ar AS action_plan_status_title_ar,
			action_plan_statuses.color_code AS action_plan_status_color_code,
			action_plan_statuses.type AS action_plan_status_type,
			identity_review_action_plan_status_updates.comment,
			identity_review_action_plan_status_updates.percentage,
			CONVERT(VARCHAR, identity_review_action_plan_status_updates.created_at, 120) AS created_at,
			CONVERT(VARCHAR, identity_review_action_plan_status_updates.updated_at, 120) AS updated_at,
			created_user.id AS created_user_id,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar
		`).
		Joins("JOIN identity_review_action_plans ON identity_review_action_plans.id = identity_review_action_plan_status_updates.identity_review_action_plan_id").
		Joins("JOIN action_plan_statuses ON action_plan_statuses.id = identity_review_action_plan_status_updates.action_plan_status_id").
		Joins("JOIN users AS created_user ON created_user.id = identity_review_action_plan_status_updates.created_by_id").
		Where("identity_review_action_plan_status_updates.identity_review_action_plan_id = ?", identityReviewActionPlanID)

	if keyword != "" {
		value := "%" + keyword + "%"
		query = query.Where("identity_review_action_plan_status_updates.comment LIKE ? OR action_plan_statuses.title_en LIKE ? OR action_plan_statuses.title_ar LIKE ?", value, value, value)
	}

	if err := query.Order(sortBy + " " + sort).Find(&identityReviewActionPlanUpdates).Error; err != nil {
		return nil, err
	}

	for _, update := range identityReviewActionPlanUpdates {
		if err := repo.db.Where("identity_review_action_plan_status_update_id = ?", update.ID).Find(&update.Documents).Error; err != nil {
			return nil, err
		}
	}

	return identityReviewActionPlanUpdates, nil
}

func (repo *IdentityReviewActionPlanStatusUpdateRepository) GetByID(identityReviewActionPlanID uint, id uint) (*models.IdentityReviewActionPlanStatusUpdate, error) {
	var identityReviewActionPlanStatusUpdate models.IdentityReviewActionPlanStatusUpdate
	if err := repo.db.
		Joins("CreatedByUser").
		Joins("UpdatedByUser").
		Preload("ActionPlanStatus").
		Preload("Documents").
		Where("identity_review_action_plan_status_updates.identity_review_action_plan_id = ? AND identity_review_action_plan_status_updates.id = ?", identityReviewActionPlanID, id).
		First(&identityReviewActionPlanStatusUpdate).Error; err != nil {
		return nil, err
	}

	return &identityReviewActionPlanStatusUpdate, nil
}

func (repo *IdentityReviewActionPlanStatusUpdateRepository) GetLatestIdentityReviewActionPlanStatusUpdate(identityReviewActionPlanID uint) (*models.IdentityReviewActionPlanStatusUpdate, error) {
	var identityReviewActionPlanStatusUpdate models.IdentityReviewActionPlanStatusUpdate
	if err := repo.db.
		Joins("CreatedByUser").
		Joins("UpdatedByUser").
		Preload("ActionPlanStatus").
		Preload("Documents").
		Where("identity_review_action_plan_status_updates.identity_review_action_plan_id = ?", identityReviewActionPlanID).
		Order("identity_review_action_plan_status_updates.created_at DESC").
		First(&identityReviewActionPlanStatusUpdate).Error; err != nil {
		return nil, err
	}

	return &identityReviewActionPlanStatusUpdate, nil
}

func (repo *IdentityReviewActionPlanStatusUpdateRepository) Create(identityReviewActionPlanStatusUpdate *models.IdentityReviewActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(identityReviewActionPlanStatusUpdate).Error; err != nil {
		return 0, err
	}

	return identityReviewActionPlanStatusUpdate.ID, nil
}

func (repo *IdentityReviewActionPlanStatusUpdateRepository) Update(identityReviewActionPlanStatusUpdate *models.IdentityReviewActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(identityReviewActionPlanStatusUpdate).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IdentityReviewActionPlanStatusUpdateRepository) CreateDocuments(docs []*models.IdentityReviewActionPlanStatusUpdateDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&docs).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IdentityReviewActionPlanStatusUpdateRepository) GetDocumentByID(id uint, tx ...utility.DatabaseTransaction) (*models.IdentityReviewActionPlanStatusUpdateDocument, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var document models.IdentityReviewActionPlanStatusUpdateDocument
	if err := db.First(&document, id).Error; err != nil {
		return nil, err
	}

	return &document, nil
}

func (repo *IdentityReviewActionPlanStatusUpdateRepository) DeleteDocuments(docIds []uint, tx ...utility.DatabaseTransaction) error {
	if len(docIds) == 0 {
		return nil
	}

	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Where("id IN (?)", docIds).Delete(&models.IdentityReviewActionPlanStatusUpdateDocument{}).Error; err != nil {
		return err
	}

	return nil
}
