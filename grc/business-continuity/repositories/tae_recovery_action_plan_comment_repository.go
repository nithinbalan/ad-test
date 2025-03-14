package repositories

import (
	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type TaeRecoveryActionPlanCommentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(actionPlanID uint) ([]*businessContinuityDTOs.TaeRecoveryActionPlanCommentListDTO, error)
	GetByID(id uint) (*models.TaeRecoveryActionPlanComment, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.TaeRecoveryActionPlanComment, error)
	Create(comment *models.TaeRecoveryActionPlanComment, tx ...utility.DatabaseTransaction) (uint, error)
	Update(comment *models.TaeRecoveryActionPlanComment, tx ...utility.DatabaseTransaction) error
	Delete(actionPlanID uint, tx ...utility.DatabaseTransaction) error

	CreateDocuments(docs []*models.TaeRecoveryActionPlanCommentDocument, tx ...utility.DatabaseTransaction) error
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.TaeRecoveryActionPlanCommentDocument, error)
}

type TaeRecoveryActionPlanCommentRepository struct {
	db *gorm.DB
}

func NewTaeRecoveryActionPlanCommentRepository(db *gorm.DB) TaeRecoveryActionPlanCommentRepositoryInterface {
	return &TaeRecoveryActionPlanCommentRepository{
		db: db,
	}
}

func (repo *TaeRecoveryActionPlanCommentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *TaeRecoveryActionPlanCommentRepository) GetAll(actionPlanID uint) ([]*businessContinuityDTOs.TaeRecoveryActionPlanCommentListDTO, error) {
	var comments []*businessContinuityDTOs.TaeRecoveryActionPlanCommentListDTO
	query := repo.db.Table("tae_recovery_action_plan_comments").
		Select(`
	tae_recovery_action_plan_comments.id,
	tae_recovery_action_plan_comments.comment,
	tae_recovery_action_plan_comments.test_and_exercise_recovery_action_plan_id,
	tae_recovery_action_plan_comments.created_at,
	tae_recovery_action_plan_comments.created_by_id,
	users.first_name_en AS created_user_first_name_en,
	users.first_name_ar AS created_user_first_name_ar,
	users.last_name_en AS created_user_last_name_en,
	users.last_name_ar AS created_user_last_name_ar,
	users.initial AS created_user_initial,
	users.color AS created_user_color,
	users.email AS created_user_email,
	users.profile_image_path AS created_user_profile_image_path
`).
		Joins("LEFT JOIN users ON users.id = tae_recovery_action_plan_comments.created_by_id").
		Where("tae_recovery_action_plan_comments.test_and_exercise_recovery_action_plan_id = ?", actionPlanID)

	if err := query.Find(&comments).Error; err != nil {
		return nil, err
	}

	for _, comment := range comments {
		var documents []*businessContinuityDTOs.TaeRecoveryActionPlanCommentDocumentDTO
		repo.db.Model(&models.TaeRecoveryActionPlanCommentDocument{}).
			Where("tae_recovery_action_plan_comment_id = ?", comment.ID).
			Find(&documents)
		comment.Documents = documents
	}

	return comments, nil
}

func (repo *TaeRecoveryActionPlanCommentRepository) GetByID(id uint) (*models.TaeRecoveryActionPlanComment, error) {
	var taeRecoveryActionPlanComment *models.TaeRecoveryActionPlanComment
	result := repo.db.Preload("Documents").First(&taeRecoveryActionPlanComment, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return taeRecoveryActionPlanComment, nil
}

func (repo *TaeRecoveryActionPlanCommentRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.TaeRecoveryActionPlanComment, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var comment models.TaeRecoveryActionPlanComment
	result := db.
		Preload("TaeRecoveryActionPlan").
		Preload("CreatedByUser").
		Preload("Documents").
		Where("id = ?", id).
		First(&comment)

	if result.Error != nil {
		return nil, result.Error
	}

	return &comment, nil
}

func (repo *TaeRecoveryActionPlanCommentRepository) Create(comment *models.TaeRecoveryActionPlanComment, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(comment).Error; err != nil {
		return 0, err
	}

	return comment.ID, nil
}

func (repo *TaeRecoveryActionPlanCommentRepository) Update(comment *models.TaeRecoveryActionPlanComment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Model(&models.TaeRecoveryActionPlanComment{}).Where("id = ?", comment.ID).Updates(comment).Error; err != nil {
		return err
	}

	return nil
}

func (repo *TaeRecoveryActionPlanCommentRepository) Delete(actionPlanID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id = ?", actionPlanID).Delete(&models.TaeRecoveryActionPlanComment{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *TaeRecoveryActionPlanCommentRepository) CreateDocuments(docs []*models.TaeRecoveryActionPlanCommentDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(&docs).Error; err != nil {
		return err
	}
	return nil
}

func (repo *TaeRecoveryActionPlanCommentRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", docIds).Delete(&models.TaeRecoveryActionPlanCommentDocument{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *TaeRecoveryActionPlanCommentRepository) GetDocumentByID(id uint) (*models.TaeRecoveryActionPlanCommentDocument, error) {
	var document models.TaeRecoveryActionPlanCommentDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}
