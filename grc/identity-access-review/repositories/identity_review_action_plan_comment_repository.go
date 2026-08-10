package repositories

import (
	identityAccessReviewDTOs "jethur/grc/identity-access-review/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type IdentityReviewActionPlanCommentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(actionPlanID uint) ([]*identityAccessReviewDTOs.IdentityReviewActionPlanCommentListDTO, error)
	GetByID(id uint) (*models.IdentityReviewActionPlanComment, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.IdentityReviewActionPlanComment, error)
	Create(comment *models.IdentityReviewActionPlanComment, tx ...utility.DatabaseTransaction) (uint, error)
	Update(comment *models.IdentityReviewActionPlanComment, tx ...utility.DatabaseTransaction) error
	Delete(actionPlanID uint, tx ...utility.DatabaseTransaction) error

	CreateDocuments(docs []*models.IdentityReviewActionPlanCommentDocument, tx ...utility.DatabaseTransaction) error
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.IdentityReviewActionPlanCommentDocument, error)
}

type IdentityReviewActionPlanCommentRepository struct {
	db *gorm.DB
}

func NewIdentityReviewActionPlanCommentRepository(db *gorm.DB) IdentityReviewActionPlanCommentRepositoryInterface {
	return &IdentityReviewActionPlanCommentRepository{
		db: db,
	}
}

func (repo *IdentityReviewActionPlanCommentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *IdentityReviewActionPlanCommentRepository) GetAll(actionPlanID uint) ([]*identityAccessReviewDTOs.IdentityReviewActionPlanCommentListDTO, error) {
	var comments []*identityAccessReviewDTOs.IdentityReviewActionPlanCommentListDTO
	query := repo.db.Table("identity_review_action_plan_comments").
		Select(`
	identity_review_action_plan_comments.id,
	identity_review_action_plan_comments.comment,
	identity_review_action_plan_comments.identity_review_action_plan_id,
	identity_review_action_plan_comments.created_at,
	identity_review_action_plan_comments.created_by_id,
	users.first_name_en AS created_user_first_name_en,
	users.first_name_ar AS created_user_first_name_ar,
	users.last_name_en AS created_user_last_name_en,
	users.last_name_ar AS created_user_last_name_ar,
	users.initial AS created_user_initial,
	users.color AS created_user_color,
	users.email AS created_user_email,
	users.profile_image_path AS created_user_profile_image_path
`).
		Joins("LEFT JOIN users ON users.id = identity_review_action_plan_comments.created_by_id").
		Where("identity_review_action_plan_comments.identity_review_action_plan_id = ?", actionPlanID)

	if err := query.Find(&comments).Error; err != nil {
		return nil, err
	}

	for _, comment := range comments {
		var documents []*identityAccessReviewDTOs.IdentityReviewActionPlanCommentDocumentDTO
		repo.db.Model(&models.IdentityReviewActionPlanCommentDocument{}).
			Where("identity_review_action_plan_comment_id = ?", comment.ID).
			Find(&documents)
		comment.Documents = documents
	}

	return comments, nil
}

func (repo *IdentityReviewActionPlanCommentRepository) GetByID(id uint) (*models.IdentityReviewActionPlanComment, error) {
	var identityReviewActionPlanComment *models.IdentityReviewActionPlanComment
	result := repo.db.Preload("Documents").First(&identityReviewActionPlanComment, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return identityReviewActionPlanComment, nil
}

func (repo *IdentityReviewActionPlanCommentRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.IdentityReviewActionPlanComment, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var comment models.IdentityReviewActionPlanComment
	result := db.
		Preload("IdentityReviewActionPlan").
		Preload("CreatedByUser").
		Preload("Documents").
		Where("id = ?", id).
		First(&comment)

	if result.Error != nil {
		return nil, result.Error
	}

	return &comment, nil
}

func (repo *IdentityReviewActionPlanCommentRepository) Create(comment *models.IdentityReviewActionPlanComment, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(comment).Error; err != nil {
		return 0, err
	}

	return comment.ID, nil
}

func (repo *IdentityReviewActionPlanCommentRepository) Update(comment *models.IdentityReviewActionPlanComment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Model(&models.IdentityReviewActionPlanComment{}).Where("id = ?", comment.ID).Updates(comment).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IdentityReviewActionPlanCommentRepository) Delete(actionPlanID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id = ?", actionPlanID).Delete(&models.IdentityReviewActionPlanComment{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IdentityReviewActionPlanCommentRepository) CreateDocuments(docs []*models.IdentityReviewActionPlanCommentDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(&docs).Error; err != nil {
		return err
	}
	return nil
}

func (repo *IdentityReviewActionPlanCommentRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", docIds).Delete(&models.IdentityReviewActionPlanCommentDocument{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IdentityReviewActionPlanCommentRepository) GetDocumentByID(id uint) (*models.IdentityReviewActionPlanCommentDocument, error) {
	var document models.IdentityReviewActionPlanCommentDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}
